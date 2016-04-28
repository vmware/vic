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

package container

import (
	"sync"

	"github.com/docker/docker/runconfig"
	containertypes "github.com/docker/engine-api/types/container"
)

// This is VIC's abridged version of Docker's container object.
type VicContainer struct {
	*runconfig.StreamConfig

	ID          string
	ContainerID string
	Config      *containertypes.Config
}

// Tracks our container info from calls
type Cache struct {
	m sync.RWMutex

	containerStore map[string]*VicContainer
}

var cache *Cache

func init() {
	cache = &Cache{containerStore: make(map[string]*VicContainer)}
}

func NewVicContainer() *VicContainer {
	return &VicContainer{
		StreamConfig: runconfig.NewStreamConfig(),
		Config:       &containertypes.Config{},
	}
}

func GetCache() *Cache {
	return cache
}

func (cc *Cache) GetContainerByName(name string) *VicContainer {
	cc.m.RLock()
	defer cc.m.RUnlock()

	if container, exist := cc.containerStore[name]; exist {
		return container
	}

	return nil
}

func (cc *Cache) SaveContainer(name string, container *VicContainer) {
	cc.m.Lock()
	defer cc.m.Unlock()

	cc.containerStore[name] = container
}
