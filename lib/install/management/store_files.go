// Copyright 2016 VMware, Inc. All Rights Reserved.
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

package management

import (
	"bytes"
	"fmt"
	"path"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/lib/portlayer/storage/vsphere"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	"golang.org/x/net/context"
)

const (
	volumeRoot = "volumes"
)

func (d *Dispatcher) DeleteStores(vchVM *vm.VirtualMachine, conf *metadata.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin(""))

	ds := d.session.Datastore

	p, err := d.getVCHRootDir(vchVM) // p would be path but there's an imported package called path
	if err != nil {
		return err
	}

	var errs []string
	var emptyImages bool
	var emptyVolumes bool
	log.Infof("Removing images")
	if emptyImages, err = d.deleteImages(ds, p); err != nil {
		errs = append(errs, err.Error())
	}
	emptyVolumes, err = d.deleteDatastoreFiles(ds, path.Join(p, volumeRoot), d.force)

	if emptyImages && emptyVolumes {
		// if not empty, don't try to delete parent directory here
		log.Debugf("Removing stores directory")
		if err = d.deleteParent(ds, p); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func (d *Dispatcher) deleteParent(ds *object.Datastore, root string) error {
	defer trace.End(trace.Begin(""))

	_, err := d.deleteDatastoreFiles(ds, root, true)
	return err
}

func (d *Dispatcher) deleteImages(ds *object.Datastore, root string) (bool, error) {
	defer trace.End(trace.Begin(""))

	p := path.Join(root, vsphere.StorageImageDir)
	// alway forcing delete images
	return d.deleteDatastoreFiles(ds, p, true)
}

func (d *Dispatcher) deleteDatastoreFiles(ds *object.Datastore, path string, force bool) (bool, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("path %q, force %t", path, force)))

	var empty bool
	dsPath := ds.Path(path)

	res, err := d.lsFolder(ds, dsPath)
	if err != nil {
		if !types.IsFileNotFound(err) {
			err = errors.Errorf("Failed to browse folder %q: %s", dsPath, err)
			return empty, err
		}
		log.Debugf("Folder %q is not found", dsPath)
		empty = true
		return empty, nil
	}
	if len(res.File) > 0 && !force {
		log.Debugf("Folder %q is not empty, leave it there", dsPath)
		return empty, nil
	}

	m := object.NewFileManager(ds.Client())
	if d.isVSAN(ds) {
		if err = d.deleteFilesIteratively(m, ds, dsPath); err != nil {
			return empty, err
		}
		return true, nil
	}

	if err = d.deleteVMFSFiles(m, ds, dsPath); err != nil {
		return empty, err
	}
	return true, nil
}

func (d *Dispatcher) isVSAN(ds *object.Datastore) bool {
	dsType, _ := ds.Type(d.ctx)

	return dsType == types.HostFileSystemVolumeFileSystemTypeVsan
}

func (d *Dispatcher) deleteFilesIteratively(m *object.FileManager, ds *object.Datastore, dsPath string) error {
	defer trace.End(trace.Begin(dsPath))

	// Get sorted result to make sure children files listed ahead of folder. Then we can empty folder before delete it
	// This function specifically designed for vSan, as vSan sometimes will throw error to delete folder is the folder is not empty
	res, err := d.getSortedChildren(ds, dsPath)
	if err != nil {
		if !types.IsFileNotFound(err) {
			err = errors.Errorf("Failed to browse sub folders %q: %s", dsPath, err)
			return err
		}
		log.Debugf("Folder %q is not found", dsPath)
		return nil
	}

	for _, path := range res {
		if err = d.deleteVMFSFiles(m, ds, path); err != nil {
			return err
		}
	}
	return d.deleteVMFSFiles(m, ds, dsPath)
}

func (d *Dispatcher) deleteVMFSFiles(m *object.FileManager, ds *object.Datastore, dsPath string) error {
	defer trace.End(trace.Begin(dsPath))

	if _, err := tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return m.DeleteDatastoreFile(ctx, dsPath, d.session.Datacenter)
	}); err != nil {
		log.Debugf("Failed to delete %q: %s", dsPath, err)
	}
	return nil
}

// getSortedChildren returns all children under datastore path in reversed order.
func (d *Dispatcher) getSortedChildren(ds *object.Datastore, dsPath string) ([]string, error) {
	res, err := d.lsSubFolder(ds, dsPath)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, dir := range res.HostDatastoreBrowserSearchResults {
		for _, f := range dir.File {
			dsf, ok := f.(*types.FileInfo)
			if !ok {
				continue
			}
			result = append(result, path.Join(dir.FolderPath, dsf.Path))
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(result)))
	return result, nil
}

func (d *Dispatcher) lsSubFolder(ds *object.Datastore, dsPath string) (*types.ArrayOfHostDatastoreBrowserSearchResults, error) {
	defer trace.End(trace.Begin(dsPath))

	spec := types.HostDatastoreBrowserSearchSpec{
		MatchPattern: []string{"*"},
	}

	b, err := ds.Browser(d.ctx)
	if err != nil {
		return nil, err
	}

	task, err := b.SearchDatastoreSubFolders(d.ctx, dsPath, &spec)
	if err != nil {
		return nil, err
	}

	info, err := task.WaitForResult(d.ctx, nil)
	if err != nil {
		return nil, err
	}

	res := info.Result.(types.ArrayOfHostDatastoreBrowserSearchResults)
	return &res, nil
}

func (d *Dispatcher) lsFolder(ds *object.Datastore, dsPath string) (*types.HostDatastoreBrowserSearchResults, error) {
	defer trace.End(trace.Begin(dsPath))

	spec := types.HostDatastoreBrowserSearchSpec{
		MatchPattern: []string{"*"},
	}

	b, err := ds.Browser(d.ctx)
	if err != nil {
		return nil, err
	}

	task, err := b.SearchDatastore(d.ctx, dsPath, &spec)
	if err != nil {
		return nil, err
	}

	info, err := task.WaitForResult(d.ctx, nil)
	if err != nil {
		return nil, err
	}

	res := info.Result.(types.HostDatastoreBrowserSearchResults)
	return &res, nil
}

func (d *Dispatcher) getVCHRootDir(vchVM *vm.VirtualMachine) (string, error) {
	defer trace.End(trace.Begin(""))

	parent := vsphere.StorageParentDir
	uuid, err := vchVM.UUID(d.ctx)
	if err != nil {
		err = errors.Errorf("Failed to get VCH UUID: %s", err)
		return "", err
	}
	return path.Join(parent, uuid), nil
}

func (d *Dispatcher) createVolumeStores(conf *metadata.VirtualContainerHostConfigSpec) error {
	for _, url := range conf.VolumeLocations {
		ds, err := d.session.Finder.Datastore(d.ctx, url.Host)
		if err != nil {
			return errors.Errorf("Could not retrieve datastore with host %q due to error %s", url.Host, err)
		}
		nds, err := datastore.NewHelper(d.ctx, d.session, ds, url.Path)
		if err != nil {
			return errors.Errorf("Could not create volume store due to error: %s", err)
		}
		// FIXME: (GitHub Issue #1301) this is not valid URL syntax and should be translated appropriately when time allows
		url.Path = nds.RootURL
	}
	return nil
}

// returns # of removed stores
func (d *Dispatcher) deleteVolumeStoreIfForced(conf *metadata.VirtualContainerHostConfigSpec) (removed int) {
	defer trace.End(trace.Begin(""))
	removed = 0

	if !d.force {
		if len(conf.VolumeLocations) == 0 {
			return 0
		}

		volumeStores := new(bytes.Buffer)
		for label, url := range conf.VolumeLocations {
			volumeStores.WriteString(fmt.Sprintf("\t%s: %s\n", label, url.Path))
		}
		log.Warnf("Since --force was not specified, the following volume stores will not be removed. Use the vSphere UI to delete content you do not wish to keep.\n%q", volumeStores.String())
		return 0
	}

	log.Infoln("Removing volume stores...")
	for label, url := range conf.VolumeLocations {
		// FIXME: url is being encoded by the portlayer incorrectly, so we have to convert url.Path to the right url.URL object
		dsURL, err := datastore.ToURL(url.Path)
		if err != nil {
			log.Warnf("Didn't receive an expected volume store path format: %q", url.Path)
			continue
		}

		log.Debugf("Provided datastore URL: %q\nParsed volume store path: %q", url.Path, dsURL.Path)

		log.Infof("Deleting volume store %q on Datastore %q at path %q", label, dsURL.Host, dsURL.Path)

		datastores, err := d.session.Finder.DatastoreList(d.ctx, dsURL.Host)

		if err != nil {
			log.Errorf("Error finding datastore %q: %s", dsURL.Host, err)
			continue
		}
		if len(datastores) > 1 {
			foundDatastores := new(bytes.Buffer)
			for _, d := range datastores {
				foundDatastores.WriteString(fmt.Sprintf("\n%s\n", d.InventoryPath))
			}
			log.Errorf("Ambiguous datastore name (%q) provided. Results were: %q", dsURL.Host, foundDatastores)
			continue
		}

		datastore := datastores[0]
		if _, err := d.deleteDatastoreFiles(datastore, dsURL.Path, d.force); err != nil {
			log.Errorf("Failed to delete volume store %q on Datastore %q at path %q", label, dsURL.Host, dsURL.Path)
		} else {
			removed++
		}
	}
	return removed

}
