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

package vsphere

import (
	"fmt"
	"net/url"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"

	"golang.org/x/net/context"
)

const volumesDir = "volumes"

// VolumeStore caches Volume references to volumes in the system.
type VolumeStore struct {

	// maps datastore uri (volume store) to datastore
	ds map[url.URL]*datastore

	// wraps our vmdks and filesystem primitives.
	dm *disk.Manager

	sess *session.Session
}

func NewVolumeStore(ctx context.Context, s *session.Session) (*VolumeStore, error) {
	dm, err := disk.NewDiskManager(ctx, s)
	if err != nil {
		return nil, err
	}

	v := &VolumeStore{
		dm:   dm,
		sess: s,
		ds:   make(map[url.URL]*datastore),
	}

	return v, nil
}

// AddStore adds a volumestore by uri.
//
// ds is the Datastore volumes
// parentDir is the path to parent directory on the datastore.  The path must
//       exist. The resulting path will be parentDir/VIC/<vch uuid>/volumes.
// storeName is the name used to refer to the ds + parentDir tupple.
//
// returns the URL used to refer to the volume store
func (v *VolumeStore) AddStore(ctx context.Context, ds *object.Datastore, parentDir, storeName string) (*url.URL, error) {
	u, err := util.VolumeStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	if _, ok := v.ds[*u]; ok {
		return nil, fmt.Errorf("volumestore (%s) already added", u.String())
	}

	// Root our datastore by the directory structured above.
	p := path.Join(parentDir, storageParentDir, volumesDir)
	d, err := newDatastore(ctx, v.sess, ds, p)
	if err != nil {
		return nil, fmt.Errorf("volumestore (%s:%s) error: %s", storeName, p, err)
	}

	v.ds[*u] = d
	return u, nil
}

func (v *VolumeStore) getDatastore(store *url.URL) (*datastore, error) {
	// find the datastore
	dstore, ok := v.ds[*store]
	if !ok {
		return nil, fmt.Errorf("volumestore (%s) not found", store.String())
	}

	return dstore, nil
}

// Returns the path to the vol relative to the given store.  The dir structure
// for a vol in the datastore is `<root>/VIC/<vch uuid>/volumes/<vol ID>/<vol ID>.vmkd`.
// Everything up to "volumes" is taken care of by the datastore wrapper.
func (v *VolumeStore) volDirPath(ID string) string {
	return ID
}

// Returns the path to the vmdk itself (in datastore URL format)
func (v *VolumeStore) volDiskDsURL(store *url.URL, ID string) (string, error) {
	// find the datastore
	dstore, err := v.getDatastore(store)
	if err != nil {
		return "", err
	}

	// XXX this could be hidden in a helper.  We shouldn't use rooturl outside the datastore struct
	return path.Join(dstore.rooturl, v.volDirPath(ID), ID+".vmdk"), nil
}

func (v *VolumeStore) VolumeCreate(ctx context.Context, ID string, store *url.URL, capacityMB uint64, info map[string][]byte) (*storage.Volume, error) {

	// find the datastore
	dstore, err := v.getDatastore(store)
	if err != nil {
		return nil, err
	}

	// Create the volume directory in the store.
	_, err = dstore.Mkdir(ctx, false, ID)
	if err != nil {
		return nil, err
	}

	// Get the path to the disk in datastore uri format
	var volDiskDsURL string
	volDiskDsURL, err = v.volDiskDsURL(store, ID)
	if err != nil {
		return nil, err
	}

	// Create the disk
	vmdisk, err := v.dm.CreateAndAttach(ctx, volDiskDsURL, "", int64(capacityMB), os.O_RDWR)
	if err != nil {
		return nil, err
	}
	defer v.dm.Detach(ctx, vmdisk)

	// Make the filesystem and set its label to the volume ID
	if err = vmdisk.Mkfs(ID); err != nil {
		return nil, err
	}

	// XXX persist the metadata

	vol, err := storage.NewVolume(store, ID, vmdisk)
	if err != nil {
		return nil, err
	}

	log.Infof("volumestore: %s (%s)", ID, vol.SelfLink)
	return vol, nil
}

func (v *VolumeStore) VolumeDestroy(ctx context.Context, ID string) error {
	return fmt.Errorf("TBD.  Not supported yet")
}

func (v *VolumeStore) VolumeGet(ctx context.Context, ID string) (*storage.Volume, error) {
	// We can't get the volume directly without looking up what datastore it's on.
	return nil, fmt.Errorf("not supported: use VolumesList")
}

func (v *VolumeStore) VolumesList(ctx context.Context) ([]*storage.Volume, error) {

	volumes := []*storage.Volume{}

	for volStore, vols := range v.ds {

		store := volStore

		res, err := vols.Ls(ctx, "")
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
			volDiskDsURL, err := v.volDiskDsURL(&store, ID)
			if err != nil {
				return nil, err
			}

			dev, err := disk.NewVirtualDisk(volDiskDsURL)
			if err != nil {
				return nil, err
			}

			vol, err := storage.NewVolume(&store, ID, dev)
			if err != nil {
				return nil, err
			}

			volumes = append(volumes, vol)
		}

	}

	return volumes, nil
}
