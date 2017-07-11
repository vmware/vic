// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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

package vsphere

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/lib/portlayer/storage/compute"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/govmomi/vim25/mo"
)

const VolumesDir = "volumes"

// VolumeStore caches Volume references to volumes in the system.
type VolumeStore struct {
	// helper to the backend
	ds *datastore.Helper

	// wraps our vmdks and filesystem primitives.
	dm *disk.Manager

	// Service url to this VolumeStore
	SelfLink *url.URL
}

func NewVolumeStore(op trace.Operation, storeName string, s *session.Session, ds *datastore.Helper) (*VolumeStore, error) {
	// Create the volume dir if it doesn't already exist
	if _, err := ds.Mkdir(op, true, VolumesDir); err != nil && !os.IsExist(err) {
		return nil, err
	}

	dm, err := disk.NewDiskManager(op, s, storage.Config.ContainerView)
	if err != nil {
		return nil, err
	}

	if DetachAll {
		if err = dm.DetachAll(op); err != nil {
			return nil, err
		}
	}

	u, err := util.VolumeStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	v := &VolumeStore{
		dm:       dm,
		ds:       ds,
		SelfLink: u,
	}

	return v, nil
}

// Returns the path to the vol relative to the given store.  The dir structure
// for a vol in the datastore is `<configured datastore path>/volumes/<vol ID>/<vol ID>.vmkd`.
// Everything up to "volumes" is taken care of by the datastore wrapper.
func (v *VolumeStore) volDirPath(ID string) string {
	return path.Join(VolumesDir, ID)
}

// Returns the path to the metadata directory for a volume
func (v *VolumeStore) volMetadataDirPath(ID string) string {
	return path.Join(v.volDirPath(ID), metaDataDir)
}

// Returns the path to the vmdk itself (in datastore URL format)
func (v *VolumeStore) volDiskDSPath(ID string) *object.DatastorePath {
	return &object.DatastorePath{
		Datastore: v.ds.RootURL.Datastore,
		Path:      path.Join(v.ds.RootURL.Path, v.volDirPath(ID), ID+".vmdk"),
	}
}

func (v *VolumeStore) VolumeCreate(op trace.Operation, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*storage.Volume, error) {

	// Create the volume directory in the store.
	if _, err := v.ds.Mkdir(op, false, v.volDirPath(ID)); err != nil {
		return nil, err
	}

	// Get the path to the disk in datastore uri format
	volDiskDSPath := v.volDiskDSPath(ID)

	config := disk.NewPersistentDisk(volDiskDSPath).WithCapacity(int64(capacityKB))
	// Create the disk
	vmdisk, err := v.dm.CreateAndAttach(op, config)
	if err != nil {
		return nil, err
	}
	defer v.dm.Detach(op, vmdisk.VirtualDiskConfig)
	vol, err := storage.NewVolume(store, ID, info, vmdisk, executor.CopyNew)
	if err != nil {
		return nil, err
	}

	// Make the filesystem and set its label
	if err = vmdisk.Mkfs(vol.Label); err != nil {
		return nil, err
	}

	// Persist the metadata
	metaDataDir := v.volMetadataDirPath(ID)
	if err = writeMetadata(op, v.ds, metaDataDir, info); err != nil {
		return nil, err
	}

	op.Infof("volumestore: %s (%s)", ID, vol.SelfLink)
	return vol, nil
}

func (v *VolumeStore) VolumeDestroy(op trace.Operation, vol *storage.Volume) error {
	volDir := v.volDirPath(vol.ID)

	op.Infof("VolumeStore: Deleting %s", volDir)
	if err := v.ds.Rm(op, volDir); err != nil {
		op.Errorf("VolumeStore: delete error: %s", err.Error())
		return err
	}
	return nil
}

func (v *VolumeStore) VolumeGet(op trace.Operation, ID string) (*storage.Volume, error) {
	// We can't get the volume directly without looking up what datastore it's on.
	return nil, fmt.Errorf("not supported: use VolumesList")
}

func (v *VolumeStore) VolumesList(op trace.Operation) ([]*storage.Volume, error) {
	volumes := []*storage.Volume{}

	res, err := v.ds.Ls(op, VolumesDir)
	if err != nil {
		return nil, fmt.Errorf("error listing vols: %s", err)
	}

	for _, f := range res.File {
		file, ok := f.(*types.FileInfo)
		if !ok {
			continue
		}

		ID := file.Path

		// Get the path to the disk in datastore uri format
		volDiskDSPath := v.volDiskDSPath(ID)

		config := disk.NewPersistentDisk(volDiskDSPath)
		dev, err := disk.NewVirtualDisk(config, v.dm.Disks)
		if err != nil {
			return nil, err
		}

		metaDataDir := v.volMetadataDirPath(ID)
		meta, err := getMetadata(op, v.ds, metaDataDir)
		if err != nil {
			return nil, err
		}

		vol, err := storage.NewVolume(v.SelfLink, ID, meta, dev, executor.CopyNew)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, vol)
	}

	return volumes, nil
}

// Export reads the delta between child and parent volume layers, returning
// the difference as a tar archive.
//
// store - the volume store containing the two layers
// id - must inherit from ancestor if ancestor is specified
// ancestor - the volume layer up the chain against which to diff
// spec - describes filters on paths found in the data (include, exclude, strip)
// data - set to true to include file data in the tar archive, false to include headers only
func (v *VolumeStore) Export(op trace.Operation, store *url.URL, id, ancestor string, spec *archive.FilterSpec, data bool) (io.ReadCloser, error) {
	_, err := util.VolumeStoreName(store)
	if err != nil {
		return nil, err
	}

	mounts := []*object.DatastorePath{}
	cleanFunc := func() {
		for _, mount := range mounts {
			if err := v.dm.UnmountAndDetach(op, mount, !persistent); err != nil {
				op.Infof("Error cleaning up disk: %s", err.Error())
			}
		}
	}

	c := v.volDiskDSPath(id)
	childFs, err := v.dm.AttachAndMount(op, c, !persistent)
	if err != nil {
		return nil, err
	}
	mounts = append(mounts, c)

	ancestorFs := ancestor
	if ancestor != "" {
		a := v.volDiskDSPath(ancestor)
		ancestorFs, err = v.dm.AttachAndMount(op, a, !persistent)
		if err != nil {
			cleanFunc()
			return nil, err
		}
		mounts = append(mounts, a)
	}

	tar, err := archive.Diff(op, childFs, ancestorFs, spec, data)
	if err != nil {
		cleanFunc()
		return nil, err
	}

	// wrap in a cleanReader so we can cleanup after the stream finishes
	return &cleanReader{
		ReadCloser: tar,
		clean:      cleanFunc,
	}, nil
}

func (v *VolumeStore) Import(op trace.Operation, store *url.URL, id string, spec *archive.FilterSpec, tarstream io.ReadCloser) error {
	_, err := util.VolumeStoreName(store)
	if err != nil {
		return err
	}

	diskRefPath := v.volDiskDSPath(id)

	mountPath, err := v.dm.AttachAndMount(op, diskRefPath, persistent)
	if err != nil {
		return err
	}
	defer func() {
		err := v.dm.UnmountAndDetach(op, diskRefPath, persistent)
		if err != nil {
			op.Infof("Error cleaning up disk: %s", err.Error())
		}
	}()

	return archive.Unpack(op, tarstream, spec, mountPath)
}

func (v *VolumeStore) StatPath(op trace.Operation, storeId, deviceId, target string) (*compute.FileStat, error) {
	diskDsURI := v.volDiskDSPath(deviceId)

	// check if the disk is in use first
	config := disk.NewPersistentDisk(diskDsURI)
	// filter powered off vms
	filter := func(vm *mo.VirtualMachine) bool {
		return vm.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOn
	}

	vms, err := v.dm.InUse(op, config, filter)
	if err != nil {
		return nil, err
	}
	if vms != nil {
		return nil, &storage.ErrDiskInUse{}
	}

	mountPath, err := v.dm.AttachAndMount(op, diskDsURI, persistent)
	if err != nil {
		return nil, err
	}

	defer func() {
		e1 := v.dm.UnmountAndDetach(op, diskDsURI, persistent)
		if e1 != nil {
			op.Errorf(e1.Error())
		}
	}()

	return compute.InspectFileStat(filepath.Join(mountPath, target))
}
