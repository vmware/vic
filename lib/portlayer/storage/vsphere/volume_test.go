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
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	portlayer "github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"golang.org/x/net/context"
)

func TestVolumeCreateListAndRestart(t *testing.T) {
	client := datastore.Session(context.TODO(), t)
	if client == nil {
		return
	}

	ctx := context.TODO()

	// Create the backing store on vsphere
	vsVolumeStore, err := NewVolumeStore(ctx, client)
	if !assert.NoError(t, err) || !assert.NotNil(t, vsVolumeStore) {
		return
	}

	// Root our datastore
	testStorePath := datastore.TestName("voltest")
	ds, err := datastore.NewHelper(ctx, client, client.Datastore, testStorePath)
	if !assert.NoError(t, err) || !assert.NotNil(t, ds) {
		return
	}

	// Add a volume store and give it a name ("testStoreName")
	volumeStore, err := vsVolumeStore.AddStore(ctx, ds, "testStoreName")
	if !assert.NoError(t, err) || !assert.NotNil(t, volumeStore) {
		return
	}

	// test we can list it
	m, err := vsVolumeStore.VolumeStoresList(ctx)
	if !assert.NoError(t, err) || !assert.NotNil(t, m) {
		return
	}

	// test the returned url matches
	s, ok := m["testStoreName"]
	if !assert.True(t, ok) || !assert.Equal(t, testStorePath, filepath.Base(s.String())) {
		return
	}

	// Clean up the mess
	defer func() {
		fm := object.NewFileManager(client.Vim25())
		tasks.WaitForResult(context.TODO(), func(ctx context.Context) (tasks.Task, error) {
			return fm.DeleteDatastoreFile(ctx, client.Datastore.Path(testStorePath), client.Datacenter)
		})
	}()

	// Create the cache
	cache, err := portlayer.NewVolumeLookupCache(ctx, vsVolumeStore)
	if !assert.NoError(t, err) || !assert.NotNil(t, cache) {
		return
	}

	// Create the volumes (in parallel)
	numVols := 5
	wg := &sync.WaitGroup{}
	wg.Add(numVols)
	volumes := make(map[string]*portlayer.Volume)
	for i := 0; i < numVols; i++ {
		go func(idx int) {
			defer wg.Done()
			ID := fmt.Sprintf("testvolume-%d", idx)

			// add some metadata if i is even
			var info map[string][]byte

			if idx%2 == 0 {
				info = make(map[string][]byte)
				info[ID] = []byte(ID)
			}

			outVol, err := cache.VolumeCreate(ctx, ID, volumeStore, 10240, info)
			if !assert.NoError(t, err) || !assert.NotNil(t, outVol) {
				return
			}

			volumes[ID] = outVol
		}(i)
	}

	wg.Wait()

	// list using the datastore (skipping the cache)
	outVols, err := vsVolumeStore.VolumesList(ctx)
	if !assert.NoError(t, err) || !assert.NotNil(t, outVols) || !assert.Equal(t, numVols, len(outVols)) {
		return
	}

	for _, outVol := range outVols {
		if !assert.Equal(t, volumes[outVol.ID], outVol) {
			return
		}
	}

	// Test restart

	// Create a new vs and cache to the same datastore (simulating restart) and compare
	secondVStore, err := NewVolumeStore(ctx, client)
	if !assert.NoError(t, err) || !assert.NotNil(t, vsVolumeStore) {
		return
	}

	volumeStore, err = secondVStore.AddStore(ctx, ds, "testStoreName")
	if !assert.NoError(t, err) || !assert.NotNil(t, volumeStore) {
		return
	}
	secondCache, err := portlayer.NewVolumeLookupCache(ctx, secondVStore)
	if !assert.NoError(t, err) || !assert.NotNil(t, cache) {
		return
	}

	secondOutVols, err := secondCache.VolumesList(ctx)
	if !assert.NoError(t, err) || !assert.NotNil(t, secondOutVols) || !assert.Equal(t, numVols, len(secondOutVols)) {
		return
	}

	for _, outVol := range secondOutVols {
		// XXX we could compare the Volumes, but the paths are different the
		// second time around on vsan since the vsan UUID is not included.
		if !assert.NotEmpty(t, volumes[outVol.ID].Device.DiskPath()) {
			return
		}
	}
}
