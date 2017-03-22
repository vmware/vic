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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/truncindex"

	"github.com/vmware/vic/lib/apiservers/engine/backends/container"
)

// Tracks our container info from calls
type CCache struct {
	m sync.RWMutex

	idIndex            *truncindex.TruncIndex
	containersByID     map[string]*container.VicContainer
	containersByName   map[string]*container.VicContainer
	containersByExecID map[string]*container.VicContainer
}

var containerCache *CCache

func init() {
	containerCache = &CCache{
		idIndex:            truncindex.NewTruncIndex([]string{}),
		containersByID:     make(map[string]*container.VicContainer),
		containersByName:   make(map[string]*container.VicContainer),
		containersByExecID: make(map[string]*container.VicContainer),
	}
}

// ContainerCache returns a reference to the container cache
func ContainerCache() *CCache {
	return containerCache
}

func (cc *CCache) getContainer(nameOrID string) *container.VicContainer {
	// get the full ID if we only have a prefix
	if cid, err := cc.idIndex.Get(nameOrID); err == nil {
		nameOrID = cid
	}

	if container, exist := cc.containersByID[nameOrID]; exist {
		return container
	}

	if container, exist := cc.containersByName[nameOrID]; exist {
		return container
	}
	return nil
}

func (cc *CCache) GetContainer(nameOrID string) *container.VicContainer {
	cc.m.RLock()
	defer cc.m.RUnlock()

	return cc.getContainer(nameOrID)
}

func (cc *CCache) AddContainer(container *container.VicContainer) {
	cc.m.Lock()
	defer cc.m.Unlock()

	// TODO(jzt): this probably shouldn't assume a valid container ID
	if err := cc.idIndex.Add(container.ContainerID); err != nil {
		log.Warnf("Error adding ID into index: %s", err)
	}
	cc.containersByID[container.ContainerID] = container
	cc.containersByName[container.Name] = container
}

func (cc *CCache) DeleteContainer(nameOrID string) {
	cc.m.Lock()
	defer cc.m.Unlock()

	container := cc.getContainer(nameOrID)
	if container == nil {
		return
	}

	delete(cc.containersByID, container.ContainerID)
	delete(cc.containersByName, container.Name)

	if err := cc.idIndex.Delete(container.ContainerID); err != nil {
		log.Warnf("Error deleting ID from index: %s", err)
	}

	// remove exec references
	for _, id := range container.List() {
		container.Delete(id)
	}
}

func (cc *CCache) AddExecToContainer(container *container.VicContainer, eid string, config *types.ExecConfig) {
	cc.m.Lock()
	defer cc.m.Unlock()

	// ignore if we already have it
	if _, ok := cc.containersByExecID[eid]; ok {
		return
	}

	container.Add(eid, config)
	cc.containersByExecID[eid] = container
}

func (cc *CCache) GetContainerFromExec(eid string) *container.VicContainer {
	cc.m.RLock()
	defer cc.m.RUnlock()

	if container, exist := cc.containersByExecID[eid]; exist {
		return container
	}
	return nil
}
