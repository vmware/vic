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

package event

import (
	"github.com/vmware/vic/pkg/vsphere/session"
)

type EventManager interface {

	// AddMonitoredObject will add the object for event listening
	AddMonitoredObject(ref string) error

	// RemoveMonitoredObject will remove the object from event listening
	RemoveMonitoredObject(ref string)

	// Register will for event callbacks
	Register(caller string, callback func(Event, *session.Session))

	// Unregsiter from event callbacks
	Unregister(caller string)

	// Registry will return the callback map
	Registry() map[string]func(Event, *session.Session)

	// RegistryCount returns the count of callbacks
	RegistryCount() int

	// Start listening for events
	Start() error

	// Stop listening for events
	Stop()

	// Blacklist will identify an object to be omitted from event callbacks
	Blacklist(ref string) error

	// Unblacklist will remove an object that was previously blacklisted
	Unblacklist(ref string)
}
