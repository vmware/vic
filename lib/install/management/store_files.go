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
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/lib/portlayer/storage/vsphere"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	"golang.org/x/net/context"
)

const (
	volumeRoot = "volumes"
)

func (d *Dispatcher) IsVSAN(ds *object.Datastore) bool {
	var isVSAN bool
	dsType, err := ds.Type(d.ctx)
	if err != nil {
		log.Debugf("Failed to get datastore type of %s, treat it as normal datastore", ds.Name())
	}
	if dsType == types.HostFileSystemVolumeFileSystemTypeVsan {
		log.Debugf("Datastore %s is VSAN", ds.Name())
		isVSAN = true
	}
	return isVSAN
}

func (d *Dispatcher) DeleteDataStores(vchVM *vm.VirtualMachine, conf *metadata.VirtualContainerHostConfigSpec) error {
	ds, err := d.session.Finder.Datastore(d.ctx, conf.ImageStores[0].Host)
	if err != nil {
		err = errors.Errorf("Failed to find image datastore %s", conf.ImageStores[0].Host)
		return err
	}
	d.session.Datastore = ds
	isVSAN := d.IsVSAN(ds)

	path, err := d.getVCHRootDir(vchVM)
	if err != nil {
		return err
	}

	var errs []string
	var emptyImages bool
	var emptyVolumes bool
	log.Infof("Removing images")
	if emptyImages, err = d.deleteImages(ds, path, isVSAN); err != nil {
		errs = append(errs, err.Error())
	}
	log.Infof("Removing volumes")
	if emptyVolumes, err = d.deleteVolumes(ds, path, isVSAN); err != nil {
		errs = append(errs, err.Error())
	} else if !emptyVolumes {
		log.Infof("Volumes directory %s is not empty, to delete with --force specified", path)
	}

	if emptyImages && emptyVolumes {
		// if not empty, don't try to delete parent directory here
		log.Infof("Removing root directory")
		if err = d.deleteParent(ds, path, isVSAN); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func (d *Dispatcher) deleteParent(ds *object.Datastore, root string, isVSAN bool) error {
	_, err := d.deleteDatastoreFiles(ds, root, isVSAN, true)
	return err
}

func (d *Dispatcher) deleteImages(ds *object.Datastore, root string, isVSAN bool) (bool, error) {
	p := path.Join(root, vsphere.StorageImageDir)
	// alway forcing delete images
	return d.deleteDatastoreFiles(ds, p, isVSAN, true)
}

func (d *Dispatcher) deleteVolumes(ds *object.Datastore, root string, isVSAN bool) (bool, error) {
	p := path.Join(root, volumeRoot)
	// if not forced delete, leave volumes there. Cause user data can be persisted in volumes
	return d.deleteDatastoreFiles(ds, p, isVSAN, d.force)
}

func (d *Dispatcher) deleteDatastoreFiles(ds *object.Datastore, path string, isVSAN bool, force bool) (bool, error) {
	var empty bool
	dsPath := ds.Path(path)

	res, err := d.lsFolder(ds, dsPath)
	if err != nil {
		if !types.IsFileNotFound(err) {
			err = errors.Errorf("Failed to browse folder %s, %s", dsPath, err)
			return empty, err
		}
		log.Debugf("Folder %s is not found")
		empty = true
		return empty, nil
	}
	if len(res.File) > 0 && !force {
		log.Debugf("Folder %s is not empty, leave it there", dsPath)
		return empty, nil
	}
	if isVSAN {
		err = d.deleteVSANFiles(ds, path, d.force)
	} else {
		err = d.deleteVMFSFiles(ds, dsPath, d.force)
	}
	if err == nil {
		empty = true
	}
	return empty, err
}

func (d *Dispatcher) deleteVMFSFiles(ds *object.Datastore, dsPath string, force bool) error {
	m := object.NewFileManager(ds.Client())

	if _, err := tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return m.DeleteDatastoreFile(ctx, dsPath, d.session.Datacenter)
	}); err != nil {
		log.Debugf("Failed to delete %s, %s", dsPath, err)
	}
	return nil
}

func (d *Dispatcher) deleteVSANFiles(ds *object.Datastore, path string, force bool) error {
	nm := object.NewDatastoreNamespaceManager(ds.Client())
	return nm.DeleteDirectory(d.ctx, d.session.Datacenter, path)
}

func (d *Dispatcher) lsFolder(ds *object.Datastore, dsPath string) (*types.HostDatastoreBrowserSearchResults, error) {
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
	parent := vsphere.StorageParentDir
	uuid, err := vchVM.UUID(d.ctx)
	if err != nil {
		err = errors.Errorf("Failed to get VCH UUID, %s", err)
		return "", err
	}
	return path.Join(parent, uuid), nil
}
