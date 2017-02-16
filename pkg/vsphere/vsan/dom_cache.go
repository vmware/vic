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
	"regexp"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

const (
	vsanNamePattern = "vsan:.*/"
	vmdkFileSuffix  = ".vmdk"
)

var (
	vsanNameRegex, _ = regexp.Compile(vsanNamePattern)
)

// vsanDatastoreCache is the map from vmdk file path to vsan dom object uuid, all vmdk files in the datastore will be cached
type vsanDSDomCache struct {
	ds   *object.Datastore
	hvis *HostVsanInternalSystem
	// vsan dom object uuid: vsan dom object path, this path does not have vsan namespace prefix in this format /vmfs/volumes/vsan:52932941b44e2147-f1490d38c9730c6d/.
	// uuids keeps all vsan dom objects in this datastore, to improve VMDK uuid query performance, because API does not provide any way to filter based on object path
	uuids map[string]string
	// vsan dom object path:object uuid. This path has same truncation with the above map
	paths map[string]string
	m     *sync.Mutex
}

// DeleteVMDKDoms deletes vmdk dom objects if the vmdk file exists in dom cache, if not, return undeleted files
func (v *vsanDSDomCache) DeleteVMDKDoms(ctx context.Context, paths []string) ([]string, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("paths: %s", paths)))
	var uuids []string
	var leftVMDKs []string

	v.m.Lock()
	for _, path := range paths {
		if !strings.HasSuffix(path, vmdkFileSuffix) {
			log.Debugf("non vmdk file %s, no need to delete here", path)
			continue
		}
		if uuid, ok := v.paths[path]; ok {
			uuids = append(uuids, uuid)
		} else {
			leftVMDKs = append(leftVMDKs, path)
		}
	}
	v.m.Unlock()

	force := true
	deleteResult, err := v.hvis.DeleteVsanObjects(ctx, uuids, &force)
	if err != nil {
		return leftVMDKs, DomDeleteError{
			Err:         err,
			FailedUuids: uuids,
		}
	}
	var deletedDoms []string
	var failedIds []string
	for _, r := range deleteResult {
		if r.Success {
			deletedDoms = append(deletedDoms, r.Uuid)
		} else {
			failedIds = append(failedIds, r.Uuid)
		}
	}
	v.m.Lock()
	defer v.m.Unlock()
	v.removeFromCache(deletedDoms)

	if len(failedIds) > 0 {
		return leftVMDKs, DomDeleteError{
			FailedUuids: failedIds,
			Result:      deleteResult,
		}
	}
	return leftVMDKs, nil
}

// removeFromCache deletes obsolete objects from cache.
// this method does not lock cache
func (v *vsanDSDomCache) removeFromCache(uuids []string) {
	for _, uuid := range uuids {
		vmdk := v.uuids[uuid]
		delete(v.uuids, uuid)
		delete(v.paths, vmdk)
	}
}

// Refresh searches dom objects from vsan datastore, and build reverse index for vmdk files
// The vmdk file format is removed vsan datastore header, e.g. /vmfs/volumes/vsan:52932941b44e2147-f1490d38c9730c6d/, to make it searchable through vmfs file path
func (v *vsanDSDomCache) Refresh(ctx context.Context) error {
	defer trace.End(trace.Begin(fmt.Sprintf("%s", v.ds.Reference().String())))
	v.m.Lock()
	defer v.m.Unlock()

	uuids, err := v.hvis.QueryVsanObjectUuidsByFilter(ctx, nil, 0, 0)
	if err != nil {
		return err
	}

	// do not query existing uuids, as we don't rename vmdk file
	setIds := make(map[string]struct{}, len(uuids))
	for _, uuid := range uuids {
		setIds[uuid] = struct{}{}
	}

	var obsoletes []string
	for exist := range v.uuids {
		if _, ok := setIds[exist]; ok {
			delete(setIds, exist)
		} else {
			obsoletes = append(obsoletes, exist)
		}
	}
	// remove deleted uuids
	v.removeFromCache(obsoletes)
	leftIds := make([]string, len(setIds))
	i := 0
	for left := range setIds {
		leftIds[i] = left
		i++
	}
	mAttrs, err := v.hvis.GetVsanObjExtAttrs(ctx, leftIds)
	if err != nil {
		return err
	}
	// fill cache
	for key, val := range mAttrs {
		p, err := v.truncateFilePath(val.Path)
		if err != nil {
			return err
		}
		if p == "" {
			continue
		}
		v.uuids[key] = p
	}

	// create reverse index
	for key, path := range v.uuids {
		v.paths[path] = key
	}
	return nil
}

func (v *vsanDSDomCache) CleanOrphanDoms(ctx context.Context) ([]string, error) {
	defer trace.End(trace.Begin(""))

	// query file manager to see if the vmdk file exists
	orphanVMDKs, err := v.queryOrphanVMDKs(ctx)
	if err != nil {
		err = errors.Errorf("failed to get vmdk file information: %s", err)
		return nil, err
	}
	log.Debugf("Found orphan vmdks: %s", orphanVMDKs)
	return v.DeleteVMDKDoms(ctx, orphanVMDKs)
}

func (v *vsanDSDomCache) queryOrphanVMDKs(ctx context.Context) ([]string, error) {
	defer trace.End(trace.Begin(v.ds.Reference().String()))
	var vmdks []string
	v.m.Lock()
	for k := range v.paths {
		vmdks = append(vmdks, k)
	}
	v.m.Unlock()

	var orphanVMDKs []string
	for _, vmdk := range vmdks {
		if _, err := v.ds.Stat(ctx, vmdk); err != nil {
			switch err.(type) {
			case object.DatastoreNoSuchDirectoryError,
				object.DatastoreNoSuchFileError:
				log.Debugf("vmdk %s is not found in datastore: %s", vmdk, err)
				orphanVMDKs = append(orphanVMDKs, vmdk)
			default:
				return orphanVMDKs, err
			}
		}
	}
	return orphanVMDKs, nil
}

// truncateFilePath removes vsan namespace prefix, return "" if not vmdk file
func (v *vsanDSDomCache) truncateFilePath(path string) (string, error) {
	defer trace.End(trace.Begin(path))

	if !strings.HasSuffix(path, vmdkFileSuffix) {
		log.Debugf("non vmdk file %s, no need to cache", path)
		return "", nil
	}

	loc := vsanNameRegex.FindStringIndex(path)
	if len(loc) == 0 {
		err := errors.Errorf("pattern %s is not found in path %s", vsanNamePattern, path)
		return "", err
	}
	return path[loc[1]:], nil
}
