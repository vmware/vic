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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/event"
	"github.com/vmware/vic/pkg/uid"
)

func TestContainerCache(t *testing.T) {
	NewContainerCache()
	containerID := "1234"
	id := uid.Parse(containerID)

	// create a new container
	NewContainer(id)
	// shouldn't have a container as it's not commited
	assert.Equal(t, len(containers.cache), 0)

	// create a new container
	container := &Container{ExecConfig: &executor.ExecutorConfig{}}
	container.ExecConfig.ID = containerID
	// put it in the cache
	containers.Put(container)
	// Get it
	cachedContainer := containers.Container(containerID)
	// did we find it?
	assert.NotNil(t, cachedContainer)
	// do we have this one in the cache?
	assert.Equal(t, cachedContainer.ExecConfig.ID, containerID)
	// remove the container
	containers.Remove(containerID)
	assert.Equal(t, len(containers.cache), 0)
	// remove non-existent container
	containers.Remove("blahblah")
}

func TestIsContainerID(t *testing.T) {
	validID := "aba122943"
	invalidID := "ABC-XZ_@"

	assert.True(t, isContainerID(validID))
	assert.False(t, isContainerID(invalidID))
}

func TestPoweredOnEvents(t *testing.T) {
	// if container is starting then viewed that poweredOn event is part of PL activity
	po := event.ContainerPoweredOn
	assert.EqualValues(t, StateStarting, eventedState(po, StateStarting))
	// if container is running and poweredOn event then it's either PL activity or it's been handled
	assert.EqualValues(t, StateRunning, eventedState(po, StateRunning))
	// if container stopped and then poweredOn event recieved then return state running
	assert.EqualValues(t, StateRunning, eventedState(po, StateStopped))
	// if container suspended and then poweredOn event recieved then return state running
	assert.EqualValues(t, StateRunning, eventedState(po, StateSuspended))
}

func TestPoweredOffEvents(t *testing.T) {
	// if container is stopping then viewed that poweredOff event is part of PL activity
	po := event.ContainerPoweredOff
	assert.EqualValues(t, StateStopping, eventedState(po, StateStopping))
	// if container is stopped and poweredOff event then it's either PL activity or it's been handled
	assert.EqualValues(t, StateStopped, eventedState(po, StateStopped))
	// if container running and then poweredOff event recieved then return state stopped
	assert.EqualValues(t, StateStopped, eventedState(po, StateRunning))
}

func TestSuspendedEvents(t *testing.T) {
	// if container is suspending (pause) then viewed that suspended event is part of PL activity
	se := event.ContainerSuspended
	assert.EqualValues(t, StateSuspending, eventedState(se, StateSuspending))
	// if container is suspeded (paused) and suspended event then it's either PL activity or it's been handled
	assert.EqualValues(t, StateSuspended, eventedState(se, StateSuspended))
	// if container running and then suspended event recieved then return state stopped
	assert.EqualValues(t, StateSuspended, eventedState(se, StateRunning))
}
