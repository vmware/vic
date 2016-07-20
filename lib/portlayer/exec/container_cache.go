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

package exec

import (
	"sync"
)

/*
* ContainerCache will provide an in-memory cache of containerVMs.  It will
* be refreshed on portlayer start and updated via container lifecycle
* operations (start, stop, rm) and well as in response to infrastructure
* events
 */
type containerCache struct {
	m sync.RWMutex

	// cache maps containerVM to container ID && infrastructure ID
	cache map[string]*Container
}

var containers *containerCache

func NewContainerCache() {
	containers = &containerCache{
		cache: make(map[string]*Container),
	}
}

// returns a reference to a specific container
//
// idOrRef will either be a container ID or an infra
// reference to the containerVM.  In the case of vSphere
// the infra reference is a MoRef
//
func (conCache *containerCache) Container(idOrRef string) *Container {
	conCache.m.RLock()
	defer conCache.m.RUnlock()
	// find by id
	container := conCache.cache[idOrRef]
	return container
}

// puts a container in the cache and will overwrite an existing container
func (conCache *containerCache) Put(container *Container) {
	conCache.m.Lock()
	defer conCache.m.Unlock()

	// add pointer to cache by container ID
	conCache.cache[container.ExecConfig.ID] = container
	//if we have a moref then add to infra ID
	if container.vm != nil {
		conCache.cache[container.vm.Reference().String()] = container
	}
}
func (conCache *containerCache) Remove(idOrRef string) {
	conCache.m.Lock()
	defer conCache.m.Unlock()
	// find by id
	container := conCache.cache[idOrRef]
	// container := cache.Container(idOrRef)
	if container != nil {
		delete(conCache.cache, container.ExecConfig.ID)
		if container.vm != nil {
			delete(conCache.cache, container.vm.Reference().String())
		}
	}

}
