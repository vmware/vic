// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/portlayer/event"
	"github.com/vmware/vic/lib/portlayer/event/events"
)

var containerEvents []events.Event

func TestEventedState(t *testing.T) {
	// poweredOn event
	event := events.ContainerPoweredOn
	assert.EqualValues(t, StateStarting, eventedState(event, StateStarting))
	assert.EqualValues(t, StateRunning, eventedState(event, StateRunning))
	assert.EqualValues(t, StateRunning, eventedState(event, StateStopped))
	assert.EqualValues(t, StateRunning, eventedState(event, StateSuspended))

	// powerOff event
	event = events.ContainerPoweredOff
	assert.EqualValues(t, StateStopping, eventedState(event, StateStopping))
	assert.EqualValues(t, StateStopped, eventedState(event, StateStopped))
	assert.EqualValues(t, StateStopped, eventedState(event, StateRunning))

	// suspended event
	event = events.ContainerSuspended
	assert.EqualValues(t, StateSuspending, eventedState(event, StateSuspending))
	assert.EqualValues(t, StateSuspended, eventedState(event, StateSuspended))
	assert.EqualValues(t, StateSuspended, eventedState(event, StateRunning))

	// removed event
	event = events.ContainerRemoved
	assert.EqualValues(t, StateRemoved, eventedState(event, StateRunning))
	assert.EqualValues(t, StateRemoved, eventedState(event, StateStopped))
	assert.EqualValues(t, StateRemoving, eventedState(event, StateRemoving))
}

func TestPublishContainerEvent(t *testing.T) {

	NewContainerCache()
	containerEvents = make([]events.Event, 0)
	Config = Configuration{}

	mgr := event.NewEventManager()
	Config.EventManager = mgr
	mgr.Subscribe(events.NewEventType(events.ContainerEvent{}).Topic(), "testing", containerCallback)

	// create new running container and place in cache
	id := "123439"
	container := newTestContainer(id)
	addTestVM(container)
	container.SetState(StateRunning)
	Containers.Put(container)

	publishContainerEvent(id, time.Now().UTC(), events.ContainerPoweredOff)
	time.Sleep(time.Millisecond * 30)

	assert.Equal(t, 1, len(containerEvents))
	assert.Equal(t, id, containerEvents[0].Reference())
	assert.Equal(t, events.ContainerPoweredOff, containerEvents[0].String())
}

func containerCallback(ee events.Event) {
	containerEvents = append(containerEvents, ee)
}

func TestVMRemovedEventCallback(t *testing.T) {

	NewContainerCache()
	containerEvents = make([]events.Event, 0)
	Config = Configuration{}

	mgr := event.NewEventManager()
	Config.EventManager = mgr
	// subscribe the exec layer to the event stream for Vm events
	mgr.Subscribe(events.NewEventType(events.ContainerEvent{}).Topic(), "testing", func(e events.Event) {
		if c := Containers.Container(e.Reference()); c != nil {
			c.OnEvent(e)
		}
	})

	// create new running container and place in cache
	id := "123439"
	container := newTestContainer(id)
	addTestVM(container)
	container.SetState(StateRunning)
	Containers.Put(container)

	publishContainerEvent(id, time.Now().UTC(), events.ContainerRemoved)
	time.Sleep(time.Millisecond * 30)
	assert.True(t, Containers.Container(id) == nil, "Container should be removed")

	addTestVM(container)
	Containers.put(container)
	container.vm.EnterFixingState()
	publishContainerEvent(id, time.Now().UTC(), events.ContainerRemoved)
	time.Sleep(time.Millisecond * 30)
	assert.True(t, Containers.Container(id) != nil, "Container should not be removed in fixing status")

	container.vm.LeaveFixingState()
	publishContainerEvent(id, time.Now().UTC(), events.ContainerRemoved)
	time.Sleep(time.Millisecond * 30)
	assert.True(t, Containers.Container(id) == nil, "Container should be removed if not in fixing status")

}
