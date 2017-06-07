// Copyright 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package vsan is a workaround for vsan DOM object leaking issue caused by FileManager.DeleteDatastoreFile, see github issue #3787 and bugzilla issue #1808703
// This file used draft vSphere API, which is subject to change in the future, so this workaround should be removed as soon as the DOM leaking issue is fixed by vSAN.
package vsan

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

var (
	SyncedDomCache  = newDomCache()
	refreshInterval = 5 * time.Minute
)

func init() {
	go SyncedDomCache.run()
}

type DomDeleteError struct {
	Err         error
	FailedUuids []string
	// if delete dom objects returns no error, but with some of uuid delete failure, this result will return detail error information for each uuid
	Result []types.HostVsanInternalSystemDeleteVsanObjectsResult
}

func (e DomDeleteError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("Failed to delete uuids: %s", e.FailedUuids)
}

type DSDomCache interface {
	Refresh(ctx context.Context) error
	DeleteVMDKDoms(ctx context.Context, paths []string) ([]string, error)
	CleanOrphanDoms(ctx context.Context) ([]string, error)
}

type dsWaitGroup struct {
	dsRef string
	wg    *sync.WaitGroup
}

type syncedDomCache struct {
	// datastore map - this cache can support multiple vsan datastores
	dsMap map[string]DSDomCache
	m     sync.RWMutex
	// refresh request channel, dsWaitGroup.wg.Done() will be called after this datastore's refresh is done
	refresh chan *dsWaitGroup
	// channel for stop request
	stop chan bool
	// finished refresh channel, after this request is accepted, all waitgroups for this datastore refresh will be invoked
	doneOnce chan string
	// wait refresh request channel, which does not trigger new refresh, just wait for current one if there is
	waitCurrentRefresh chan *dsWaitGroup
	// map to hold all waitgroups for one datastore refresh
	refreshings map[string][]*sync.WaitGroup
}

func newDomCache() *syncedDomCache {
	return &syncedDomCache{
		dsMap:              make(map[string]DSDomCache),
		refresh:            make(chan *dsWaitGroup, 10),
		stop:               make(chan bool, 1),
		doneOnce:           make(chan string, 10),
		waitCurrentRefresh: make(chan *dsWaitGroup, 10),
		refreshings:        make(map[string][]*sync.WaitGroup),
	}
}

func (c *syncedDomCache) closeChannels() {
	defer trace.End(trace.Begin(""))
	close(c.refresh)
	close(c.stop)
	close(c.doneOnce)
	close(c.waitCurrentRefresh)
	c.m.Lock()
	defer c.m.Unlock()
	log.Debugf("Finish all waitgroups")
	for i := range c.refreshings {
		if len(c.refreshings[i]) == 0 {
			continue
		}
		c.doneWaitGroups(c.refreshings[i])
	}
}

func (c *syncedDomCache) doneWaitGroups(wgs []*sync.WaitGroup) {
	defer trace.End(trace.Begin(fmt.Sprintf("%#v waitgroups is done", wgs)))
	for i := range wgs {
		wgs[i].Done()
	}
}

func (c *syncedDomCache) datastoreDom(ds string) DSDomCache {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.dsMap[ds]
}

// AddDomCache create dom cache for Datastore
func (c *syncedDomCache) AddDomCache(ctx context.Context, ds *object.Datastore) error {
	if ds == nil {
		defer trace.End(trace.Begin(fmt.Sprintf("datastore is empty")))
		return nil
	}
	defer trace.End(trace.Begin(ds.Reference().String()))

	ref := ds.Reference().String()
	if dsc := c.datastoreDom(ref); dsc != nil {
		log.Debugf("Dom cache for datastore %s already exists", ds.InventoryPath)
		return nil
	}

	if dsType, _ := ds.Type(ctx); dsType != types.HostFileSystemVolumeFileSystemTypeVsan {
		log.Debugf("datastore %s is not vsan, no need to build vsan dom cache", ds.InventoryPath)
		return nil
	}

	hosts, err := ds.AttachedHosts(ctx)
	if err != nil {
		err = errors.Errorf("failed to get attached hosts for datastore %s: %s", ds.InventoryPath, err)
		log.Error(err)
		return err
	}
	if len(hosts) == 0 {
		err = errors.Errorf("no host attached to datastore %s", ds.InventoryPath)
		log.Error(err)
		return err
	}

	// vsan object can be query and removed from any host whatever the owner is, so here just get the first host
	hc := hosts[0].ConfigManager()
	var h mo.HostSystem

	err = hc.Properties(ctx, hc.Reference(), []string{"configManager.vsanInternalSystem"}, &h)
	if err != nil {
		err = errors.Errorf("failed to get vsanInternalSystem for host %s: %s", hosts[0].Reference(), err)
		log.Error(err)
		return err
	}

	if h.ConfigManager.VsanInternalSystem == nil {
		err = errors.Errorf("vsanInternalSystem of host %s is empty", hosts[0].Reference())
		log.Error(err)
		return err
	}

	hvis := NewHostVsanInternalSystem(ds.Client(), *h.ConfigManager.VsanInternalSystem)
	dsc := &vsanDSDomCache{
		ds:    ds,
		hvis:  hvis,
		uuids: make(map[string]string),
		paths: make(map[string]string),
	}
	log.Debugf("Dom cache for datastore %s is created, start to refresh", ds.InventoryPath)
	c.m.Lock()
	c.dsMap[ds.Reference().String()] = dsc
	c.m.Unlock()

	c.Refresh(ds.Reference().String(), nil)
	return nil
}

func (c *syncedDomCache) getNonRefreshingDS(requests []*dsWaitGroup) []*dsWaitGroup {
	defer trace.End(trace.Begin(fmt.Sprintf("%#v", requests)))

	c.m.Lock()
	defer c.m.Unlock()
	var ret []*dsWaitGroup
	for _, request := range requests {
		if _, ok := c.refreshings[request.dsRef]; ok {
			// this datastore is refreshing, add wait group into the list
			if request.wg != nil {
				log.Debugf("datastore %s is refreshing, append waiting group only", request.dsRef)
				c.refreshings[request.dsRef] = append(c.refreshings[request.dsRef], request.wg)
			}
		} else {
			ret = append(ret, request)
		}
	}
	return ret
}

// doRefresh should be called by run only to make sure concurrency is handled properly
func (c *syncedDomCache) doRefresh(ctx context.Context, requests []*dsWaitGroup) error {
	defer trace.End(trace.Begin(fmt.Sprintf("%#v", requests)))
	dsws := c.getNonRefreshingDS(requests)
	if len(dsws) == 0 {
		log.Debugf("Requested datastores %#v is refreshing", requests)
		return nil
	}

	c.m.Lock()
	defer c.m.Unlock()
	var errs []string
	for _, dsw := range dsws {
		dsc, ok := c.dsMap[dsw.dsRef]
		if !ok {
			err := errors.Errorf("requested datastore %s does not exist", c.dsMap[dsw.dsRef])
			log.Error(err)
			errs = append(errs, err.Error())
			continue
		}
		if dsw.wg != nil {
			c.refreshings[dsw.dsRef] = append(c.refreshings[dsw.dsRef], dsw.wg)
			log.Debugf("waiting groups for datastore %s: %#v", dsw.dsRef, c.refreshings[dsw.dsRef])
		} else {
			log.Debugf("wait group for datastore %s is empty", dsw.dsRef)
			if _, ok := c.refreshings[dsw.dsRef]; !ok {
				// no refreshing item yet
				c.refreshings[dsw.dsRef] = nil
			}
		}
		go func(dsRef string) {
			if err := dsc.Refresh(ctx); err != nil {
				log.Error(err)
			}
			c.doneOnce <- dsRef
		}(dsw.dsRef)
	}
	if len(errs) > 0 {
		return errors.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}

// Start cache refresh after first refresh request coming. After that, refresh once every 5 minutes.
// Allow to inject more refresh request between timer timeout. But if the request comes while refresh is running, append the refresh waitgroup only, without new refresh triggered
func (c *syncedDomCache) run() {
	var tick <-chan time.Time
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case r := <-c.refresh:
			log.Debugf("start refresh datastore %s", r.dsRef)
			c.doRefresh(ctx, []*dsWaitGroup{r})
		case <-c.stop:
			log.Debugf("Stop vsan dom cache refresh")
			c.closeChannels()
			return
		case ds := <-c.doneOnce:
			log.Debugf("datastore %s is refreshed", ds)
			func() {
				c.m.Lock()
				defer c.m.Unlock()

				if _, ok := c.refreshings[ds]; ok {
					c.doneWaitGroups(c.refreshings[ds])
				}
				delete(c.refreshings, ds)
				if len(c.refreshings) == 0 {
					log.Debugf("All refresh finished, restart timer")
					tick = time.After(refreshInterval)
				}
			}()
		case w := <-c.waitCurrentRefresh:
			log.Debug("waiting refresh for datastore %s", w.dsRef)
			func() {
				c.m.RLock()
				defer c.m.RUnlock()
				if _, ok := c.refreshings[w.dsRef]; ok {
					if w.wg != nil {
						c.refreshings[w.dsRef] = append(c.refreshings[w.dsRef], w.wg)
					} else {
						log.Errorf("waiting request for datastore %s should have non-nil waitgroup.", w.dsRef)
					}
				} else {
					log.Debugf("No refresh running for datastore %s, done waitgroup immediately", w.dsRef)
					w.wg.Done()
				}
			}()
		case <-tick:
			log.Debugf("Timer timout, restart refresh for all datastores")
			reqs := func() []*dsWaitGroup {
				c.m.Lock()
				defer c.m.Unlock()
				var requests []*dsWaitGroup
				for k := range c.dsMap {
					r := &dsWaitGroup{
						dsRef: k,
						wg:    nil,
					}
					requests = append(requests, r)
				}
				return requests
			}()
			c.doRefresh(ctx, reqs)
		}
	}
}

// Refresh trigger vsan dom cache refresh. This method will return immediately. wg.Done will be called after the refresh is finished, if wg is not nil
func (c *syncedDomCache) Refresh(dsRef string, wg *sync.WaitGroup) {
	c.refresh <- &dsWaitGroup{
		dsRef: dsRef,
		wg:    wg,
	}
}

// Stop DOM Cache refresh
func (c *syncedDomCache) Stop(dsRef string, wg *sync.WaitGroup) {
	c.stop <- true
}

func (c *syncedDomCache) waitRefresh(dsRef string) {
	defer trace.End(trace.Begin(dsRef))

	var w sync.WaitGroup
	w.Add(1)
	c.waitCurrentRefresh <- &dsWaitGroup{
		dsRef: dsRef,
		wg:    &w,
	}
	w.Wait()
	return
}

// SyncDeleteVMDKDom deletes vsan dom object through HostVsanInternalSystem. If vmdk file is not found in dom cache, wait current refresh and refresh once again if oneMoreRefresh is true, and retry delete. Returns undeleted files
// path is datastore path format, with namespace uuid as root dir
// This method should be called after vmdk file is deleted
func (c *syncedDomCache) SyncDeleteVMDKDoms(ctx context.Context, ds *object.Datastore, paths []string, oneMoreRefresh bool) ([]string, error) {
	if ds == nil {
		defer trace.End(trace.Begin(fmt.Sprintf("datastore is empty, paths: %s, one more refresh: %t", paths, oneMoreRefresh)))
		return nil, nil
	}
	defer trace.End(trace.Begin(fmt.Sprintf("datastore: %s, paths: %s, one more refresh: %t", ds.Reference(), paths, oneMoreRefresh)))

	leftPaths, err := c.deleteVMDKDoms(ctx, ds, paths)
	if err != nil {
		return leftPaths, err
	}

	if len(leftPaths) == 0 {
		return leftPaths, nil
	}

	// wait current refresh finish
	c.waitRefresh(ds.Reference().String())
	log.Debugf("Current refresh for datastore %s is done", ds.Reference())

	if oneMoreRefresh {
		// refresh and try again
		var wg sync.WaitGroup
		wg.Add(1)

		log.Debugf("Refresh datastore %s dom cache one more time", ds.Reference())
		c.Refresh(ds.Reference().String(), &wg)
		wg.Wait()
		log.Debugf("Refresh finished")
	}
	if leftPaths, err = c.deleteVMDKDoms(ctx, ds, leftPaths); err != nil {
		return leftPaths, err
	}
	if len(leftPaths) > 0 {
		log.Debugf("vmdk files dom objects %s are not found, ignore the request", leftPaths)
	}
	return nil, nil
}

// DeleteVMDKDoms deletes vmdk dom objects if the vmdk file exists in dom cache, if not, return undeleted files
// path is datastore path format, with namespace uuid as root dir
// This method should be called after vmdk file is deleted
func (c *syncedDomCache) deleteVMDKDoms(ctx context.Context, ds *object.Datastore, paths []string) ([]string, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("datastore %s, paths: %s", ds.Reference(), paths)))

	dsc := c.datastoreDom(ds.Reference().String())
	if dsc == nil {
		err := errors.Errorf("datastore %s is not cached", ds.Reference())
		log.Error(err)
		return paths, err
	}
	return dsc.DeleteVMDKDoms(ctx, paths)
}

// SyncCleanOrphanDoms deletes vsan dom objects without vmdk file backed in one datastore. This method will wait current refresh and refresh once again if oneMoreRefresh is true, and try cleanup.
func (c *syncedDomCache) SyncCleanOrphanDoms(ctx context.Context, ds *object.Datastore, oneMoreRefresh bool) ([]string, error) {
	if ds == nil {
		defer trace.End(trace.Begin(fmt.Sprintf("datastore is empty, one more refresh: %t", oneMoreRefresh)))
		return nil, nil
	}
	defer trace.End(trace.Begin(fmt.Sprintf("datastore: %s, one more refresh: %t", ds.Reference(), oneMoreRefresh)))

	// wait current refresh finish
	c.waitRefresh(ds.Reference().String())
	log.Debugf("Current refresh for datastore %s is done", ds.Reference())

	if oneMoreRefresh {
		// refresh and try again
		var wg sync.WaitGroup
		wg.Add(1)

		log.Debugf("Refresh datastore %s dom cache one more time", ds.Reference())
		c.Refresh(ds.Reference().String(), &wg)
		wg.Wait()
		log.Debugf("Refresh finished")
	}
	return c.cleanOrphanDoms(ctx, ds)
}

// cleanOrphanDoms deletes vsan dom objects without vmdk file backed in one datastore.
func (c *syncedDomCache) cleanOrphanDoms(ctx context.Context, ds *object.Datastore) ([]string, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("datastore %s", ds.Reference())))

	dsc := c.datastoreDom(ds.Reference().String())
	if dsc == nil {
		err := errors.Errorf("datastore %s is not cached", ds.Reference())
		log.Error(err)
		return nil, err
	}
	return dsc.CleanOrphanDoms(ctx)
}
