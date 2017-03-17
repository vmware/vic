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
	"net/url"
	"os"
	"path"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
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

	dm, err := disk.NewDiskManager(op, s)
	if err != nil {
		return nil, err
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
func (v *VolumeStore) volDiskDsURL(ID string) (string, error) {
	// XXX this could be hidden in a helper.  We shouldn't use rooturl outside the datastore struct
	return path.Join(v.ds.RootURL.String(), v.volDirPath(ID), ID+".vmdk"), nil
}

func (v *VolumeStore) VolumeCreate(op trace.Operation, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*storage.Volume, error) {

	// Create the volume directory in the store.
	if _, err := v.ds.Mkdir(op, false, v.volDirPath(ID)); err != nil {
		return nil, err
	}

	// Get the path to the disk in datastore uri format
	volDiskDsURL, err := v.volDiskDsURL(ID)
	if err != nil {
		return nil, err
	}

	// Create the disk
	vmdisk, err := v.dm.CreateAndAttach(op, volDiskDsURL, "", int64(capacityKB), os.O_RDWR)
	if err != nil {
		return nil, err
	}
	defer v.dm.Detach(op, vmdisk)

	vol, err := storage.NewVolume(store, ID, info, vmdisk)
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
		volDiskDsURL, err := v.volDiskDsURL(ID)
		if err != nil {
			return nil, err
		}

		dev, err := disk.NewVirtualDisk(volDiskDsURL)
		if err != nil {
			return nil, err
		}

		metaDataDir := v.volMetadataDirPath(ID)
		meta, err := getMetadata(op, v.ds, metaDataDir)
		if err != nil {
			return nil, err
		}

		vol, err := storage.NewVolume(v.SelfLink, ID, meta, dev)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, vol)
	}

	return volumes, nil
}
