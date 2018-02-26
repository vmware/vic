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

package mock

import (
	"fmt"
	"net/url"
	"os"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/vic/lib/portlayer/storage/volume"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

type MockVolumeStore struct {
	// id -> volume
	db map[string]*volume.Volume
}

func NewMockVolumeStore() *MockVolumeStore {
	m := &MockVolumeStore{
		db: make(map[string]*volume.Volume),
	}

	return m
}

// Creates a volume on the given volume store, of the given size, with the given metadata.
func (m *MockVolumeStore) VolumeCreate(op trace.Operation, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*volume.Volume, error) {
	storeName, err := util.VolumeStoreName(store)
	if err != nil {
		return nil, err
	}

	selfLink, err := util.VolumeURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	vol := &volume.Volume{
		ID:       ID,
		Store:    store,
		SelfLink: selfLink,
	}

	m.db[ID] = vol

	return vol, nil
}

// Get an existing volume via it's ID and volume store.
func (m *MockVolumeStore) VolumeGet(op trace.Operation, ID string) (*volume.Volume, error) {
	vol, ok := m.db[ID]
	if !ok {
		return nil, os.ErrNotExist
	}

	return vol, nil
}

// Destroys a volume
func (m *MockVolumeStore) VolumeDestroy(op trace.Operation, vol *volume.Volume) error {
	if _, ok := m.db[vol.ID]; !ok {
		return os.ErrNotExist
	}

	delete(m.db, vol.ID)

	return nil
}

func (m *MockVolumeStore) VolumeStoresList(op trace.Operation) (map[string]url.URL, error) {
	return nil, fmt.Errorf("not implemented")
}

// Lists all volumes on the given volume store`
func (m *MockVolumeStore) VolumesList(op trace.Operation) ([]*volume.Volume, error) {
	var i int
	list := make([]*volume.Volume, len(m.db))
	for _, v := range m.db {
		t := *v
		list[i] = &t
		i++
	}

	return list, nil
}

func (m *MockVolumeStore) URL(op trace.Operation, id string) (*url.URL, error) {
	return nil, nil
}

func (m *MockVolumeStore) Owners(op trace.Operation, url *url.URL, filter func(vm *mo.VirtualMachine) bool) ([]*vm.VirtualMachine, error) {
	return nil, nil
}
