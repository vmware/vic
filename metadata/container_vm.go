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

package metadata

import (
	"net/url"
	"os/exec"
)

// Common data between managed entities, across execution environments
type Common struct {
	// A reference to the components hosting execution environment, if any
	ExecutionEnvironment string
	// Unambiguous ID with meaning in the context of its hosting execution environment
	ID string
	// Convenience field to record a human readable name
	Name string
	// Freeform notes related to the entity
	Notes string
}

// MountSpec details a mount that must be executed within the executor
type MountSpec struct {
	// The label that identifies which attached storage is the target
	Label string
	// The path in the executor at which this should be mounted
	Path string
	// Freeform mode string, which could translate directly to mount options
	// We may want to turn this into a more structured form eventually
	Mode string
}

// SessionControl defines the set of operations that can have a custom signal mapping
// on a per session basis
type SessionControl int

// ContainerVM holds that data tightly associated with a containerVM, but that should not
// be visible to the guest. This is the external complement to ExecutorConfig.
type ContainerVM struct {
	Common

	// The version of the bootstrap image that this container was booted from.
	Version string

	// Name aliases for this specific container, Maps alias to unambiguous name
	// This uses unambiguous name rather than reified network endpoint to persist
	// the intent rather than a point-in-time manifesting of that intent.
	Aliases map[string]string
	// The location of the interaction service that the tether should connect to. Examples:
	// * tcp://x.x.x.x:2377
	// * vmci://moid - should this be an moid or a VMCI CID? Does one insulate us from reboots?
	Interaction url.URL
}

// ExecutorConfig holds the data tightly associated with an Executor. This is distinct from Sessions
// in that there is no process inherently associated - this is closer to a ThreadPool than a Thread and
// is the owner of the shared filesystem environment. This is the guest visible complement to ContainerVM.
type ExecutorConfig struct {
	Common

	// Maps the mount name to the detail mount specification
	Mounts map[string]MountSpec
	// This describes an executors presence on a network, and contains sufficient
	// information to configure the interface in the guest.
	Networks map[string]NetworkEndpoint
}

// SessionConfig defines the content of a session - this maps to the root of a process tree
// inside an executor
// This is close to but not perfectly aligned with the new docker/docker/daemon/execdriver/driver:CommonProcessConfig
type SessionConfig struct {
	// The primary session may have the same ID as the executor owning it
	Common

	// The primary process for the session - this is an exec.Cmd only for now as it contains everything
	// we need... along with some we don't
	Cmd exec.Cmd
	// Allocate a tty or not
	Tty bool
	// Maps the intent to the signal for this specific app
	Signals map[SessionControl]int
	// Use struct composition to add in the guest specific portions
	// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
	// ulimits
	// user
	// rootfs - within the container context
}
