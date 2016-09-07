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
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"

	"golang.org/x/net/context"
)

var volumesDir = StorageParentDir + "/volumes"

// VolumeStore caches Volume references to volumes in the system.
type VolumeStore struct {

	// maps datastore uri (volume store) to datastore
	ds map[url.URL]*datastore.Helper

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
		ds:   make(map[url.URL]*datastore.Helper),
	}

	return v, nil
}

// AddStore adds a volumestore by uri.
//
// ds is the Datastore (+ path) volumes will be created under.
//       The resulting path will be parentDir/VIC/volumes.
// storeName is the name used to refer to the datastore + path (ds above).
//
// returns the URL used to refer to the volume store
func (v *VolumeStore) AddStore(ctx context.Context, ds *datastore.Helper, storeName string) (*url.URL, error) {
	u, err := util.VolumeStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	if _, ok := v.ds[*u]; ok {
		return nil, fmt.Errorf("volumestore (%s) already added", u.String())
	}

	if _, err = ds.Mkdir(ctx, true, volumesDir); err != nil && !os.IsExist(err) {
		return nil, err
	}

	v.ds[*u] = ds
	return u, nil
}

func (v *VolumeStore) VolumeStoresList(ctx context.Context) (map[string]url.URL, error) {
	m := make(map[string]url.URL)
	for u, ds := range v.ds {

		// Get the ds:// URL given the datastore url ("[datastore] /path")
		dsurl, err := datastore.ToURL(ds.RootURL)
		if err != nil {
			return nil, err
		}

		// from the storage url, get the store name
		storeName, err := util.VolumeStoreName(&u)
		if err != nil {
			return nil, err
		}

		m[storeName] = *dsurl
	}

	return m, nil
}

func (v *VolumeStore) getDatastore(store *url.URL) (*datastore.Helper, error) {
	// find the datastore
	dstore, ok := v.ds[*store]
	if !ok {
		return nil, VolumeStoreNotFoundError{msg: fmt.Sprintf("volumestore (%s) not found", store.String())}
	}

	return dstore, nil
}

// Returns the path to the vol relative to the given store.  The dir structure
// for a vol in the datastore is `<configured datastore path>/volumes/<vol ID>/<vol ID>.vmkd`.
// Everything up to "volumes" is taken care of by the datastore wrapper.
func (v *VolumeStore) volDirPath(ID string) string {
	return path.Join(volumesDir, ID)
}

// Returns the path to the metadata directory for a volume
func (v *VolumeStore) volMetadataDirPath(ID string) string {
	return path.Join(v.volDirPath(ID), metaDataDir)
}

// Returns the path to the vmdk itself (in datastore URL format)
func (v *VolumeStore) volDiskDsURL(store *url.URL, ID string) (string, error) {
	// find the datastore
	dstore, err := v.getDatastore(store)
	if err != nil {
		return "", err
	}

	// XXX this could be hidden in a helper.  We shouldn't use rooturl outside the datastore struct
	return path.Join(dstore.RootURL, v.volDirPath(ID), ID+".vmdk"), nil
}

func (v *VolumeStore) VolumeCreate(ctx context.Context, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*storage.Volume, error) {

	// find the datastore
	dstore, err := v.getDatastore(store)
	if err != nil {
		return nil, err
	}

	// Create the volume directory in the store.
	_, err = dstore.Mkdir(ctx, false, v.volDirPath(ID))
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
	vmdisk, err := v.dm.CreateAndAttach(ctx, volDiskDsURL, "", int64(capacityKB), os.O_RDWR)
	if err != nil {
		return nil, err
	}
	defer v.dm.Detach(ctx, vmdisk)

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
	if err = writeMetadata(ctx, dstore, metaDataDir, info); err != nil {
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

		res, err := vols.Ls(ctx, volumesDir)
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

			metaDataDir := v.volMetadataDirPath(ID)
			meta, err := getMetadata(ctx, vols, metaDataDir)
			if err != nil {
				return nil, err
			}

			vol, err := storage.NewVolume(&store, ID, meta, dev)
			if err != nil {
				return nil, err
			}

			volumes = append(volumes, vol)
		}

	}

	return volumes, nil
}

//Custom Error Types

// VolumeStoreNotFoundError : custom error type for when we fail to find a target volume store
type VolumeStoreNotFoundError struct {
	msg string
}

func (e VolumeStoreNotFoundError) Error() string {
	return e.msg
}

// VolumeExistsError : custom error type for when a create operation targets and already occupied ID
type VolumeExistsError struct {
	msg string
}

func (e VolumeExistsError) Error() string {
	return e.msg
}
