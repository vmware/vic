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

package nfs

import (
	"io"
	"net/url"
	"path"

	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"os"
)

const (
	// This is the same as our vmdk approach namespacing the data further on the filesystem
	VolumesDir = "volumes"

	// path that namespaces the metadata for a specific volume
	metadataDir = "volumedata"
)

type NFSv3VolumeStore struct {
	//volume store name
	Name string

	not      // Service URL to this Volumestore
	SelfLink *url.URL
}

func NewVolumeStore(op trace.Operation, storeName string, nfsTargetUrl *url.URL) (*NFSv3VolumeStore, error) {
	// XXX: Potential credential leak, maybe log less...
	op.Infof("Creating datastore (%s) at target (%q)", storeName, nfsTargetUrl.String())

	// Create a target which can interact with the nfs server.
	target, err := NewNFSv3Target(op, nfsTargetUrl)
	if err != nil {
		return nil, err
	}

	// we need to check for os.IsExist here in case the volume store is already present.
	// XXX: we must make sure that os.IsExist is a valid check here. implementation for the target
	// should return this.
	if _, err := target.MkDir(nfsTargetUrl.Path, true); !os.IsExist(err) {
		return nil, err
	}

	v := &NFSv3VolumeStore{
		Name:     storeName,
		SelfLink: nfsTargetUrl,
	}

	return v, nil
}

// Returns the path to the vol relative to the given store.  The dir structure
// for a vol in a nfs store is `<configured nfs server path>/volumes/<vol ID>/<volume contents>`.
func (v *NFSv3VolumeStore) volDirPath(ID string) string {
	return path.Join(v.SelfLink.Path, VolumesDir, ID)
}

// Returns the path to the metadata directory for a volume
func (v *NFSv3VolumeStore) volMetadataDirPath(ID string) string {
	return path.Join(v.SelfLink.Path, metadataDir, ID)
}

// Returns the path to the specified volumes directory on the nfs target
func (v *NFSv3VolumeStore) volumeURL(ID string) string {
	return path.Join(v.SelfLink.String(), v.volDirPath(ID))
}

// Creates a volume directory and volume object for NFS based volumes
func (v *NFSv3VolumeStore) VolumeCreate(op trace.Operation, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*storage.Volume, error) {
	target, err := OpenNFSv3Target(v.SelfLink)
	if err != nil {
		return err
	}

	if _, err := v.nfsTarget.MkDir(v.volumeURL(ID), false); err != nil {
		return nil, err
	}

	volBacking := NewNFSVolumeBacking(v.volDirPath(ID))

	vol := storage.NewNFSVolume(v.SelfLink, ID, info, v.volDirPath(ID))

	if err := writeMetadata(op, v, ID, info); err != nil {
		return err
	}

	op.Infof("nfs volume (%s) successfully created on volume store (%s)", ID, v.Name)

	CloseNFSv3Target(&target)
	return vol, nil
}

// Removes a volume and all of it's contents from the nfs store. We already know via the cache if it is in use.
func (v *NFSv3VolumeStore) VolumeDestroy(op trace.Operation, vol *storage.Volume) error {
	target, err := OpenNFSv3Target(v.SelfLink)
	if err != nil {
		return err
	}

	volTargetPath := v.volDirPath(vol.ID)

	op.Infof("Attempting to remove volume (%s) and it's metadata from volume store (%s)", vol.ID, v.Name)

	//remove volume directory and children
	if err := target.RemoveAll(v.volDirPath(vol.ID)); err != nil {
		op.Errorf("Failed to remove volume (%s) on volume store (%s) due to error (%s)", vol.ID, v.Name, err)
		return err
	}

	//XXX: what should we do if we lose the nfs server between volume and metadata deletion?

	//remove volume metadata directory and children
	if err := target.RemoveAll(v.volMetadataDirPath(vol.ID)); err != nil {
		op.Errorf("Failed to remove metadata for volume (%s) at path (%q) on volume store (%s)", vol.ID, v.volDirPath(vol.ID), v.Name)
		//FIXME: Should we bail here? the volume is gone at this point...
	}

	if err := CloseNFSv3Target(target); err != nil {
		return err
	}

	return nil
}

func (v *NFSv3VolumeStore) VolumesList(op trace.Operation) ([]*Volume, error) {
	return nil, nil
}

func writeMetadata(op trace.Operation, store storage.VolumeStorer, ID string, info map[string][]byte) error {
	//implement this later
	//will use target.Write
	return nil
}

func getMetadata() (map[string][]byte, error) {
	//will implement this later
	//will use target.Write
	return nil, nil
}
