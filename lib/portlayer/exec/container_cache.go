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
	"regexp"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/lib/portlayer/event/collector/vsphere"
	"github.com/vmware/vic/lib/portlayer/event/events"
)

/*
* ContainerCache will provide an in-memory cache of containerVMs.  It will
* be refreshed on portlayer start and updated via container lifecycle
* operations (start, stop, rm) and well as in response to infrastructure
* events
 */
type containerCache struct {
	m sync.RWMutex

	// cache by container id
	cache map[string]*Container
}

var containers *containerCache

func NewContainerCache() {
	// cache by the container ID and the vsphere
	// managed object reference
	containers = &containerCache{
		cache: make(map[string]*Container),
	}
	// register with the event manager
	if VCHConfig.EventManager != nil {
		// subscribe for VM Events
		VCHConfig.EventManager.Subscribe(events.NewEventType(vsphere.VmEvent{}).Topic(), "containerCache", callback)
	}
}

func (conCache *containerCache) Container(idOrRef string) *Container {
	conCache.m.RLock()
	defer conCache.m.RUnlock()
	// find by id or moref
	container := conCache.cache[idOrRef]
	return container
}

func (conCache *containerCache) Containers(allStates bool) []*Container {
	conCache.m.RLock()
	defer conCache.m.RUnlock()
	// cache contains 2 items for each container
	capacity := len(conCache.cache) / 2
	containers := make([]*Container, 0, capacity)

	for id := range conCache.cache {
		// is the key a proper ID?
		if isContainerID(id) {
			container := conCache.cache[id]
			if allStates {
				containers = append(containers, container)
			} else if container.State == StateRunning {
				// only include running...
				containers = append(containers, container)

			}
		}
	}
	return containers
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
	if container != nil {
		delete(conCache.cache, container.ExecConfig.ID)
		if container.vm != nil {
			delete(conCache.cache, container.vm.Reference().String())
		}
	}

}

func isContainerID(id string) bool {
	// ID should only be lower case chars & numbers
	match, _ := regexp.Match("^[a-z0-9]*$", []byte(id))
	return match
}

// callback for the event stream
func callback(ie events.Event) {
	// grab the container from the cache
	container := containers.Container(ie.Reference())
	if container != nil {

		newState := eventedState(ie.String(), container.State)
		// do we have a state change
		if newState != container.State {
			switch newState {
			case StateStopping,
				StateRunning,
				StateStopped,
				StateSuspended:
				log.Debugf("Container(%s) state set to %s via event activity", container.ExecConfig.ID, newState.String())
				container.State = newState
			case StateRemoved:
				log.Debugf("Container(%s) %s via event activity", container.ExecConfig.ID, newState.String())
				// remove from cache
				containers.Remove(container.ExecConfig.ID)

			}
		}
	}

	return
}

// eventedState will determine the target container
// state based on the current container state and the event
func eventedState(e string, current State) State {
	switch e {
	case events.ContainerShutdown:
		// no need to worry about existing state
		return StateStopping
	case events.ContainerPoweredOn:
		// are we in the process of starting
		if current != StateStarting {
			return StateRunning
		}
	case events.ContainerPoweredOff:
		// are we in the process of stopping
		if current != StateStopping {
			return StateStopped
		}
	case events.ContainerSuspended:
		// are we in the process of suspending
		if current != StateSuspending {
			return StateSuspended
		}
	case events.ContainerRemoved:
		if current != StateRemoving {
			return StateRemoved
		}
	}
	return current
}
