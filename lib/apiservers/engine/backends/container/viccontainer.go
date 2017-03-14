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

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
)

// VicContainer is VIC's abridged version of Docker's container object.
type VicContainer struct {
	Name        string
	ImageID     string
	ContainerID string
	Config      *containertypes.Config //Working copy of config (with overrides from container create)
	HostConfig  *containertypes.HostConfig

	m     sync.RWMutex
	execs map[string]*types.ExecConfig
}

// NewVicContainer returns a reference to a new VicContainer
func NewVicContainer() *VicContainer {
	return &VicContainer{
		Config: &containertypes.Config{},
		execs:  make(map[string]*types.ExecConfig),
	}
}

// Add adds a new exec configuration to the container.
func (v *VicContainer) Add(id string, Config *types.ExecConfig) {
	v.m.Lock()
	v.execs[id] = Config
	v.m.Unlock()
}

// Get returns an exec configuration by its id.
func (v *VicContainer) Get(id string) *types.ExecConfig {
	v.m.RLock()
	res := v.execs[id]
	v.m.RUnlock()
	return res
}

// Delete removes an exec configuration from the container.
func (v *VicContainer) Delete(id string) {
	v.m.Lock()
	delete(v.execs, id)
	v.m.Unlock()
}

// List returns the list of exec ids in the container.
func (v *VicContainer) List() []string {
	var IDs []string
	v.m.RLock()
	for id := range v.execs {
		IDs = append(IDs, id)
	}
	v.m.RUnlock()
	return IDs
}
