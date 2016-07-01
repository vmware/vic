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

package storage

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/vmware/vic/lib/portlayer/util"

	"golang.org/x/net/context"
)

type Disk interface {
	MountPath() (string, error)
	DiskPath() string
	//FIXME: Add a capacity and populate it.
}

// VolumeStorer is an interface to create, remove, enumerate, and get Volumes.
type VolumeStorer interface {
	// Creates a volume on the given volume store, of the given size, with the given metadata.
	VolumeCreate(ctx context.Context, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*Volume, error)

	// Get an existing volume via it's ID.
	VolumeGet(ctx context.Context, ID string) (*Volume, error)

	// Destroys a volume
	VolumeDestroy(ctx context.Context, ID string) error

	// Lists all volumes
	VolumesList(ctx context.Context) ([]*Volume, error)

	//Modifies the config spec of a container to attach a volume
}

// Volume is the handle to identify a volume on the backing store.  The URI
// namespace used to identify the Volume in the storage layer has the following
// path scheme:
//
// `/storage/volumes/<volume store identifier, usually the vch uuid>/<volume id>`
//
type Volume struct {
	// Identifies the volume
	ID string

	// Label is the computed label of the Volume.  This is set by the runtime.
	Label string

	// The volumestore the volume lives on. (e.g the datastore + vch + configured vol directory)
	Store *url.URL

	// Metadata the volume is included with.  Is persisted along side the volume vmdk.
	Info map[string][]byte

	// Namespace in the storage layer to look up this volume.
	SelfLink *url.URL

	// Backing device
	Device Disk
}

// NewVolume creates a Volume
func NewVolume(store *url.URL, ID string, device Disk) (*Volume, error) {
	storeName, err := util.VolumeStoreName(store)
	if err != nil {
		return nil, err
	}

	selflink, err := util.VolumeURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	// Set the label to the md5 of the ID

	vol := &Volume{
		ID:       ID,
		Label:    label(ID),
		Store:    store,
		SelfLink: selflink,
		Device:   device,
	}

	return vol, nil
}

// given an ID, compute the volume's label
func label(ID string) string {

	// e2label's manpage says the label size is 16 chars
	m := md5.Sum([]byte(ID))
	return fmt.Sprintf("%x", m)[:16]
}

func (v *Volume) Parse(u *url.URL) error {
	// Check the path isn't malformed.
	if !filepath.IsAbs(u.Path) {
		return errors.New("invalid uri path")
	}

	segments := strings.Split(filepath.Clean(u.Path), "/")[1:]

	if segments[0] != util.StorageURLPath {
		return errors.New("not a storage path")
	}

	if len(segments) < 3 {
		return errors.New("uri path mismatch")
	}

	store, err := util.VolumeStoreNameToURL(segments[2])
	if err != nil {
		return err
	}

	id := segments[3]

	var SelfLink url.URL
	SelfLink = *u

	v.ID = id
	v.SelfLink = &SelfLink
	v.Store = store

	return nil
}
