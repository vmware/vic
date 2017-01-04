// Copyright 2017 VMware, Inc. All Rights Reserved.
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
	"github.com/vmware/vic/lib/portlayer/event/events"

	"github.com/vmware/govmomi/vim25/types"
)

type HostEvent struct {
	*events.BaseEvent
}

func NewHostEvent(be types.BaseEvent) *HostEvent {
	var ee string
	// host events that we care about
	switch be.(type) {
	case *types.EnteringMaintenanceModeEvent:
		ee = events.HostEnteringMaintenanceMode
	case *types.EnteredMaintenanceModeEvent:
		ee = events.HostEnteredMaintenanceMode
	case *types.ExitMaintenanceModeEvent:
		ee = events.HostExitMaintenanceMode
	default:
		panic("Unknown event")
	}
	e := be.GetEvent()
	return &HostEvent{
		&events.BaseEvent{
			Event:       ee,
			ID:          int(e.Key),
			Detail:      e.FullFormattedMessage,
			Ref:         e.Host.Host.String(),
			CreatedTime: e.CreatedTime,
		},
	}

}

func (he *HostEvent) Topic() string {
	if he.Type == "" {
		he.Type = events.NewEventType(he)
	}
	return he.Type.Topic()
}
