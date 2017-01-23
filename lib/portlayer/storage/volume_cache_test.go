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
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"

	"context"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/trace"
)

type MockVolumeStore struct {
	// id -> volume
	db map[string]*Volume
}

func NewMockVolumeStore() *MockVolumeStore {
	m := &MockVolumeStore{
		db: make(map[string]*Volume),
	}

	return m
}

func (m *MockVolumeStore) VolumeStoresList(op trace.Operation) (map[string]url.URL, error) {
	return nil, nil
}

// Creates a volume on the given volume store, of the given size, with the given metadata.
func (m *MockVolumeStore) VolumeCreate(op trace.Operation, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*Volume, error) {
	storeName, err := util.VolumeStoreName(store)
	if err != nil {
		return nil, err
	}

	selfLink, err := util.VolumeURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	vol := &Volume{
		ID:       ID,
		Store:    store,
		SelfLink: selfLink,
	}

	m.db[ID] = vol

	return vol, nil
}

// Get an existing volume via it's ID and volume store.
func (m *MockVolumeStore) VolumeGet(op trace.Operation, ID string) (*Volume, error) {
	vol, ok := m.db[ID]
	if !ok {
		return nil, os.ErrNotExist
	}

	return vol, nil
}

// Destroys a volume
func (m *MockVolumeStore) VolumeDestroy(op trace.Operation, vol *Volume) error {
	if _, ok := m.db[vol.ID]; !ok {
		return os.ErrNotExist
	}

	delete(m.db, vol.ID)

	return nil
}

// Lists all volumes on the given volume store`
func (m *MockVolumeStore) VolumesList(op trace.Operation) ([]*Volume, error) {
	var i int
	list := make([]*Volume, len(m.db))
	for _, v := range m.db {
		t := *v
		list[i] = &t
		i++
	}

	return list, nil
}

func TestVolumeCreateGetListAndDelete(t *testing.T) {
	op := trace.NewOperation(context.Background(), "test")

	exec.NewContainerCache()

	mvs := NewMockVolumeStore()
	v := NewVolumeLookupCache(op)
	storeURL, err := v.AddStore(op, "testStore", mvs)
	if !assert.NoError(t, err) {
		return
	}

	inVols := make(map[string]*Volume)
	inVolsM := &sync.Mutex{}

	wg := &sync.WaitGroup{}
	createFn := func(i int) {
		defer wg.Done()

		id := fmt.Sprintf("ID-%d", i)

		// Write to the datastore
		vol, err := v.VolumeCreate(op, id, storeURL, 0, nil)
		if !assert.NoError(t, err) || !assert.NotNil(t, vol) {
			return
		}

		inVolsM.Lock()
		inVols[id] = vol
		inVolsM.Unlock()
	}

	// Create a set of volumes
	numVolumes := 5
	wg.Add(numVolumes)
	for i := 0; i < numVolumes; i++ {
		go createFn(i)
	}
	wg.Wait()

	getFn := func(inVol *Volume) {
		vol, err := v.VolumeGet(op, inVol.ID)
		if !assert.NoError(t, err) || !assert.NotNil(t, vol) {
			return
		}

		if !assert.Equal(t, inVol, vol) {
			return
		}
		wg.Done()
	}

	wg.Add(numVolumes)
	for _, inVol := range inVols {
		getFn(inVol)
	}
	wg.Wait()

	volumeList, err := v.VolumesList(op)
	if !assert.NoError(t, err) || !assert.Equal(t, numVolumes, len(volumeList)) {
		return
	}

	// Test that the list returned by VolumeList matches our inVols list.  Then
	// delete each vol via the cache, then check the datastore to ensure it's
	// empty
	for _, outVol := range volumeList {
		if !assert.Equal(t, inVols[outVol.ID], outVol) {
			return
		}

		if err = v.VolumeDestroy(op, outVol.ID); !assert.NoError(t, err) {
			return
		}
	}

	// check the datastore is empty.
	if !assert.Empty(t, mvs.db) {
		return
	}
}

// Create 2 store caches but use the same backing datastore.  Create images
// with the first cache, then get the image with the second.  This simulates
// restart since the second cache is empty and has to go to the backing store.
func TestVolumeCacheRestart(t *testing.T) {
	mvs := NewMockVolumeStore()
	op := trace.NewOperation(context.Background(), "test")

	firstCache := NewVolumeLookupCache(op)
	storeURL, err := firstCache.AddStore(op, "testStore", mvs)
	if !assert.NoError(t, err) || !assert.NotNil(t, storeURL) {
		return
	}

	// Create a set of volumes
	inVols := make(map[string]*Volume)
	for i := 1; i < 50; i++ {
		id := fmt.Sprintf("ID-%d", i)

		// Write to the datastore
		vol, err := firstCache.VolumeCreate(op, id, storeURL, 0, nil)
		if !assert.NoError(t, err) || !assert.NotNil(t, vol) {
			return
		}

		inVols[id] = vol
	}

	secondCache := NewVolumeLookupCache(op)
	if !assert.NotNil(t, secondCache) {
		return
	}

	_, err = secondCache.AddStore(op, "testStore", mvs)
	if !assert.NoError(t, err) || !assert.NotNil(t, storeURL) {
		return
	}

	// get the vols from the second cache to ensure it goes to the ds
	for _, expectedVol := range inVols {
		vol, err := secondCache.VolumeGet(op, expectedVol.ID)
		if !assert.NoError(t, err) || !assert.Equal(t, expectedVol, vol) {
			return
		}
	}
}
