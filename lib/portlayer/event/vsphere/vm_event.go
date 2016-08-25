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
	"github.com/vmware/vic/lib/portlayer/event"

	"github.com/vmware/govmomi/vim25/types"
)

type VMEvent struct {
	eventType string
	eventID   int
	message   string
	ref       string
}

func NewVMEvent(be types.BaseEvent) *VMEvent {
	var eventType string
	// vm events that we care about
	switch be.(type) {
	case *types.VmPoweredOnEvent:
		eventType = event.ContainerPoweredOn
	case *types.VmPoweredOffEvent:
		eventType = event.ContainerPoweredOff
	case *types.VmSuspendedEvent:
		eventType = event.ContainerSuspended
	case *types.VmRemovedEvent:
		eventType = event.ContainerRemoved
	case *types.VmGuestShutdownEvent:
		eventType = event.ContainerShutdown
	}
	e := be.GetEvent()
	return &VMEvent{
		eventType: eventType,
		eventID:   int(e.Key),
		message:   e.FullFormattedMessage,
		ref:       e.Vm.Vm.String(),
	}
}

func (vme *VMEvent) EventID() int {
	return vme.eventID
}

// return event type / description
func (vme *VMEvent) String() string {
	return vme.eventType
}

func (vme *VMEvent) Message() string {
	return vme.message
}

func (vme *VMEvent) Reference() string {
	return vme.ref
}
