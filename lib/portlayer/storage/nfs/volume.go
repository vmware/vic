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
	"bytes"
	"io"
	"net/url"
	"path"

	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"io/ioutil"
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

	//volume directory path
	volumeStoreDir string

	//nfs target for filesystem interaction
	Target NFSTarget

	// Service URL to this Volumestore
	SelfLink *url.URL
}

func NewVolumeStore(op trace.Operation, storeName string, nfsTargetURL *url.URL) (*NFSv3VolumeStore, error) {
	// XXX: Potential credential leak, maybe log less...
	op.Infof("Creating datastore (%s) at target (%q)", storeName, nfsTargetURL.String())

	// Create a target which can interact with the nfs server.
	target, err := NewNFSv3Target(op, nfsTargetURL, path.Join(nfsTargetURL.Path, VolumesDir))
	if err != nil {
		return nil, err
	}

	target, err = target.Open()
	if err != nil {
		return nil, err
	}

	defer target.Close()

	if _, err := target.MkDir(nfsTargetURL.Path, true); !os.IsExist(err) {
		return nil, err
	}

	v := &NFSv3VolumeStore{
		Name:           storeName,
		volumeStoreDir: path.Join(nfsTargetURL.Path),
		SelfLink:       nfsTargetURL,
		Target:         target,
	}

	return v, nil
}

// Returns the path to the vol relative to the given store.  The dir structure
// for a vol in a nfs store is `<configured nfs server path>/volumes/<vol ID>/<volume contents>`.
func (v *NFSv3VolumeStore) volDirPath(ID string) string {
	return path.Join(v.volumeStoreDir, VolumesDir, ID)
}

func (v *NFSv3VolumeStore) volumesDir() string {
	return path.Join(v.volumeStoreDir, VolumesDir)
}

// Returns the path to the metadata directory for a volume
func (v *NFSv3VolumeStore) volMetadataDirPath(ID string) string {
	return path.Join(v.volumeStoreDir, metadataDir, ID)
}

// Returns the path to the specified volumes directory on the nfs target
func (v *NFSv3VolumeStore) volumeURL(ID string) string {
	return path.Join(v.volumeStoreDir, v.volDirPath(ID))
}

// Creates a volume directory and volume object for NFS based volumes
func (v *NFSv3VolumeStore) VolumeCreate(op trace.Operation, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*storage.Volume, error) {
	target, err := v.Target.Open()
	if err != nil {
		return nil, err
	}

	defer target.Close()
	if _, err := target.MkDir(v.volumeURL(ID), false); err != nil {
		return nil, err
	}

	backing := NewNFSVolumeDevice(*v.Target.EndPoint(), v.volDirPath(ID))

	vol, err := storage.NewVolume(v.Target.EndPoint(), ID, info, backing)
	if err != nil {
		return nil, err
	}

	if err := writeMetadata(op, v.volMetadataDirPath(ID), info, target); err != nil {
		return nil, err
	}

	op.Infof("nfs volume (%s) successfully created on volume store (%s)", ID, v.Name)
	return vol, nil
}

// Removes a volume and all of it's contents from the nfs store. We already know via the cache if it is in use.
func (v *NFSv3VolumeStore) VolumeDestroy(op trace.Operation, vol *storage.Volume) error {
	op.Infof("Attempting to remove volume (%s) and it's metadata from volume store (%s)", vol.ID, v.Name)
	target, err := v.Target.Open()
	if err != nil {
		return err
	}
	defer target.Close()

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

	return nil
}

func (v *NFSv3VolumeStore) VolumesList(op trace.Operation) ([]*storage.Volume, error) {
	volumes := make([]*storage.Volume, 0)

	target, err := v.Target.Open()
	if err != nil {
		return nil, err
	}
	defer target.Close()

	volFileInfo, err := target.ReadDir(v.volumesDir())
	if err != nil {
		return nil, err
	}

	for _, volFileInfo := range volFileInfo {

		volMetadata, err := getMetadata(op, v.volMetadataDirPath(volFileInfo.Name()), target)
		if err != nil {
			return nil, err //perhaps wait until the end of the loop to report this?
		}

		volDeviceBacking := NewNFSVolumeDevice(target.EndPoint(), v.volDirPath(volFileInfo.Name()))

		vol, err := storage.NewVolume(v.Target.EndPoint(), volFileInfo.Name(), volMetadata, volDeviceBacking)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, vol)
	}

	return volumes, nil
}

func writeMetadata(op trace.Operation, metadataPath string, info map[string][]byte, target NFSTarget) error {
	op.Infof("Attempting to write metadata to (%s)", metadataPath)
	for fileName, data := range info {
		targetPath := path.Join(metadataPath, fileName)
		r := bytes.NewReader(data)
		//NOTE: If we begin to allow readonly behavior then we must make sure this path is writable...
		err := target.Write(op, targetPath, r)
		if err != nil {
			//XXX: Should we attempt to clean up what we wrote in this case?
			return err
		}
	}
	op.Infof("Successfully wrote metadata to (%s)", metadataPath)
	return nil
}

func getMetadata(op trace.Operation, metadataPath string, target NFSTarget) (map[string][]byte, error) {
	op.Infof("Attempting to volume metadata at (%s)", metadataPath)
	metadataInfo := make(map[string][]byte)
	dataKeys, err := target.ReadDir(metadataPath)
	if err != nil {
		return nil, err
	}

	for _, metadataFile := range dataKeys {
		targetPath := path.Join(metadataPath, metadataFile.Name())
		byteReader, err := target.Read(targetPath)
		if err != nil {
			return nil, err
		}

		dataBlob, err := ioutil.ReadAll(byteReader)
		if err != nil {
			return nil, err
		}

		metadataInfo[metadataFile.Name()] = dataBlob
	}

	op.Infof("Successfully read volume metadata at (%s)", metadataPath)
	return metadataInfo, nil
}
