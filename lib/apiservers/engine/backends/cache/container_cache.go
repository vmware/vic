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

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/pkg/truncindex"

	"github.com/vmware/vic/lib/apiservers/engine/backends/container"
)

// Tracks our container info from calls
type CCache struct {
	m sync.RWMutex

	idIndex          *truncindex.TruncIndex
	containersByID   map[string]*container.VicContainer
	containersByName map[string]*container.VicContainer
}

var containerCache *CCache

func init() {
	containerCache = &CCache{
		idIndex:          truncindex.NewTruncIndex([]string{}),
		containersByID:   make(map[string]*container.VicContainer),
		containersByName: make(map[string]*container.VicContainer),
	}
}

// ContainerCache returns a reference to the container cache
func ContainerCache() *CCache {
	return containerCache
}

func (cc *CCache) GetContainer(nameOrID string) *container.VicContainer {
	cc.m.RLock()
	defer cc.m.RUnlock()

	// get the full ID if we only have a prefix
	if cid, err := cc.idIndex.Get(nameOrID); err == nil {
		nameOrID = cid
	}

	if container := cc.getContainerByID(nameOrID); container != nil {
		return container
	}

	return cc.getContainerByName(nameOrID)
}

func (cc *CCache) getContainerByID(id string) *container.VicContainer {
	if container, exist := cc.containersByID[id]; exist {
		return container
	}
	return nil
}

func (cc *CCache) getContainerByName(name string) *container.VicContainer {
	if container, exist := cc.containersByName[name]; exist {
		return container
	}
	return nil
}

func (cc *CCache) SaveContainer(container *container.VicContainer) {
	cc.m.Lock()
	defer cc.m.Unlock()

	// TODO(jzt): this probably shouldn't assume a valid container ID
	if err := cc.idIndex.Add(container.ContainerID); err != nil {
		log.Warnf("Error inserting ID into index: %s", err)
	}
	cc.containersByID[container.ContainerID] = container
	cc.containersByName[container.Name] = container
}
