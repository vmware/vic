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
	"strconv"
	"testing"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/event"
	"github.com/vmware/vic/pkg/vsphere/session"

	"github.com/stretchr/testify/assert"
)

// used to test callbacks
var callcount int

func newVMMO() *types.ManagedObjectReference {
	return &types.ManagedObjectReference{Value: "101", Type: "vm"}
}

func TestNewEventManager(t *testing.T) {
	// TODO: use simulator
	// for now create uninitilized session
	s := &session.Session{}
	mgr := NewEventManager(s)
	assert.NotNil(t, mgr)
}

func TestMonitoredObject(t *testing.T) {

	s := &session.Session{}
	mgr := NewEventManager(s)
	mo := newVMMO()

	mgr.AddMonitoredObject(mo.String())
	mos := mgr.monitoredObjects()
	assert.Equal(t, 1, len(mos))
	mgr.RemoveMonitoredObject(mo.String())
	mos = mgr.monitoredObjects()
	assert.Equal(t, 0, len(mos))
}

func TestBlacklist(t *testing.T) {

	s := &session.Session{}
	mgr := NewEventManager(s)
	mo := newVMMO()
	mos := mo.String()

	mgr.Blacklist(mos)
	assert.Equal(t, 1, mgr.BlacklistCount())
	assert.True(t, mgr.blacklisted(*mo))
	mgr.Unblacklist(mos)
	assert.Equal(t, 0, mgr.BlacklistCount())
	assert.False(t, mgr.blacklisted(*mo))
}

func TestRegistration(t *testing.T) {
	s := &session.Session{}
	mgr := NewEventManager(s)

	mgr.Register("test", callMe)
	assert.Equal(t, 1, mgr.RegistryCount())
	assert.NotNil(t, mgr.Registry())
	assert.Equal(t, 1, len(mgr.Registry()))
	mgr.Unregister("FooBar")
	assert.Equal(t, 1, mgr.RegistryCount())
	mgr.Unregister("test")
	assert.Equal(t, 0, mgr.RegistryCount())

}

func TestEvented(t *testing.T) {
	s := &session.Session{}
	mgr := NewEventManager(s)
	callcount = 0

	// register local callback
	mgr.Register("test", callMe)

	// test lifecycle events
	page := eventPage(3, true)
	evented(mgr, page)
	assert.Equal(t, 3, callcount)

	// non-lifecycle events
	page = eventPage(2, false)
	evented(mgr, page)
	assert.Equal(t, 3, callcount)

}

// will test basic failure of no callbacks configured
// additional testing requires additions to the simulator
func TestStartFailure(t *testing.T) {
	s := &session.Session{}
	mgr := NewEventManager(s)
	assert.Error(t, mgr.Start())

}

// simple callback counter
func callMe(ie event.Event, s *session.Session) {
	callcount++
}

// utility function to mock a vsphere event
//
// size is the number of events to create
// lifeCycle is true when we want to generate state events
// lifeCycle events == poweredOn, poweredOff, etc..
func eventPage(size int, lifeCycle bool) []types.BaseEvent {
	page := make([]types.BaseEvent, 0, size)
	moid := 100
	for i := 0; i < size; i++ {
		var eve types.BaseEvent
		moid++
		vm := types.ManagedObjectReference{Value: strconv.Itoa(moid), Type: "vm"}
		if lifeCycle {
			eve = types.BaseEvent(&types.VmPoweredOnEvent{VmEvent: types.VmEvent{Event: types.Event{Vm: &types.VmEventArgument{Vm: vm}}}})
		} else {
			eve = types.BaseEvent(&types.VmReconfiguredEvent{VmEvent: types.VmEvent{Event: types.Event{Vm: &types.VmEventArgument{Vm: vm}}}})
		}

		page = append(page, eve)
	}

	return page
}
