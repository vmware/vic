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

package cache

import (
	"sync"
)

// LCache is an in-memory cache to account for existing image layers
// It is used primarily by imagec when coordinating layer downloads
// The cache is initially hydrated by way of the image cache at startup
type LCache struct {
	m sync.RWMutex
	// layers maps from layer ID to boolean
	// true indicates that the layer is currently downloading
	// false indicates that the layer has been downloaded and written to the portlayer
	layers map[string]bool
}

// LayerNotFoundError is returned when a layer does not exist in the cache
type LayerNotFoundError struct{}

func (e LayerNotFoundError) Error() string {
	return "Layer does not exist"
}

var (
	layerCache = &LCache{layers: make(map[string]bool)}
)

// LayerCache returns a reference to the layer cache
func LayerCache() *LCache {
	return layerCache
}

// AddExisting adds an existing layer to the cache
func (lc *LCache) AddExisting(id string) {
	lc.m.Lock()
	defer lc.m.Unlock()

	lc.layers[id] = false
}

// AddNew adds a new layer to the cache
func (lc *LCache) AddNew(id string) {
	lc.m.Lock()
	defer lc.m.Unlock()

	lc.layers[id] = true // a newly cached layer begins in the downloading state
}

// Remove removes a layer from the cache
func (lc *LCache) Remove(id string) {
	lc.m.Lock()
	defer lc.m.Unlock()

	delete(lc.layers, id)
}

// Commit marks a layer as downloaded
func (lc *LCache) Commit(id string) {
	lc.m.Lock()
	defer lc.m.Unlock()

	lc.layers[id] = false
}

// IsDownloading returns true if the layer is currently downloading, false otherwise
func (lc *LCache) IsDownloading(id string) (bool, error) {
	lc.m.RLock()
	defer lc.m.RUnlock()

	downloading, ok := lc.layers[id]
	if !ok {
		return false, LayerNotFoundError{}
	}
	return downloading, nil
}
