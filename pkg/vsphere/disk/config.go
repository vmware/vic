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

package disk

import (
	"fmt"
	"hash/fnv"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/fs"
)

type VirtualDiskConfig struct {
	// The URI in the datastore this disk can be found with
	DatastoreURI *object.DatastorePath

	// The URI in the datastore to the parent of this disk
	ParentDatastoreURI *object.DatastorePath

	// The size of the disk
	CapacityInKB int64

	// Underlying filesystem
	Filesystem fs.Filesystem

	// Storage provisioning mode
	ProvisionType types.VirtualDiskType

	// Base disk UUID
	UUID string

	DiskMode types.VirtualDiskMode
}

func NewPersistentDisk(URI *object.DatastorePath) *VirtualDiskConfig {
	return &VirtualDiskConfig{
		DatastoreURI:  URI,
		DiskMode:      types.VirtualDiskModeIndependent_persistent,
		Filesystem:    fs.NewExt4(),
		ProvisionType: types.VirtualDiskTypeThin,
	}
}

func NewNonPersistentDisk(URI *object.DatastorePath) *VirtualDiskConfig {
	return &VirtualDiskConfig{
		DatastoreURI:  URI,
		DiskMode:      types.VirtualDiskModeIndependent_nonpersistent,
		Filesystem:    fs.NewExt4(),
		ProvisionType: types.VirtualDiskTypeThin,
	}
}

func (d *VirtualDiskConfig) WithParent(parent *object.DatastorePath) *VirtualDiskConfig {
	d.ParentDatastoreURI = parent

	return d
}

func (d *VirtualDiskConfig) WithFilesystem(ftype fs.FilesystemType) *VirtualDiskConfig {
	d.Filesystem = fs.GetFilesystem(ftype)
	return d
}

func (d *VirtualDiskConfig) WithCapacity(capacity int64) *VirtualDiskConfig {
	d.CapacityInKB = capacity

	return d
}

func (d *VirtualDiskConfig) WithProvisionType(ptype types.VirtualDiskType) *VirtualDiskConfig {
	d.ProvisionType = ptype

	return d
}

// WithUUID can only be set on the base disk layer due to disklib bug
// TODO: add an error mechanism for validating conditional settings like this
func (d *VirtualDiskConfig) WithUUID(uuid string) *VirtualDiskConfig {
	d.UUID = uuid

	return d
}

func (d *VirtualDiskConfig) Hash() uint64 {
	key := fmt.Sprintf("%s-%t", d.DatastoreURI, d.IsPersistent())

	hash := fnv.New64a()
	hash.Write([]byte(key))

	return hash.Sum64()
}

func (d *VirtualDiskConfig) IsPersistent() bool {
	return d.DiskMode == types.VirtualDiskModeIndependent_persistent || d.DiskMode == types.VirtualDiskModePersistent
}

func (d *VirtualDiskConfig) IsThinProvisioned() *bool {
	return types.NewBool(d.ProvisionType == types.VirtualDiskTypeThin)
}
