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

package tether

import (
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/ip"
)

type ExecutorConfig struct {
	// The name of the system
	Name string `vic:"0.1" scope:"read-only" key:"common/name"`

	// ID corresponds to that of the primary session
	ID string `vic:"0.1" scope:"read-only" key:"common/id"`

	// Debug is a numeric level that controls extent of debugging
	DebugLevel int `vic:"0.1" scope:"read-only" key:"common/diagnostics/debug"`

	// Exclusive access to childPidTable
	pidMutex sync.Mutex

	// Set of child PIDs created by us.
	pids map[int]*SessionConfig

	// Sessions is the set of sessions currently hosted by this executor
	// These are keyed by session ID
	Sessions map[string]*SessionConfig `vic:"0.1" scope:"read-only" key:"sessions"`

	// Maps the mount name to the detail mount specification
	Mounts map[string]executor.MountSpec `vic:"0.1" scope:"read-only" key:"mounts"`

	// This describes an executors presence on a network, and contains sufficient
	// information to configure the interface in the guest.
	Networks map[string]*NetworkEndpoint `vic:"0.1" scope:"read-only" key:"networks"`

	// Key is the host key used during communicate back with the Interaction endpoint if any
	// Used if the in-guest tether is responsible for authenticating the connection
	Key []byte `vic:"0.1" scope:"read-only" key:"key"`
}

// SessionConfig defines the content of a session - this maps to the root of a process tree
// inside an executor
// This is close to but not perfectly aligned with the new docker/docker/daemon/execdriver/driver:CommonProcessConfig
type SessionConfig struct {
	// Protects the structure
	m sync.Mutex

	// The primary session may have the same ID as the executor owning it
	executor.Common `vic:"0.1" scope:"read-only" key:"common"`
	executor.Detail `vic:"0.1" scope:"read-write" key:"detail"`

	// Diagnostics holds basic diagnostics data
	Diagnostics executor.Diagnostics `vic:"0.1" scope:"read-write" key:"diagnostics"`

	// The primary process for the session
	Cmd exec.Cmd `vic:"0.1" scope:"read-only" key:"cmd" recurse:"depth=2,nofollow"`

	// The exit status of the process, if any
	ExitStatus int `vic:"0.1" scope:"read-write" key:"status"`

	Started string `vic:"0.1" scope:"read-write" key:"started"`

	// Allow attach
	Attach bool `vic:"0.1" scope:"read-only" key:"attach"`

	// Allocate a tty or not
	Tty bool `vic:"0.1" scope:"read-only" key:"tty"`

	// Restart controls whether a process gets relaunched if it exists
	Restart bool `vic:"0.1" scope:"read-only" key:"restart"`

	// StopSignal is the signal name or number used to stop a container
	StopSignal string `vic:"0.1" scope:"read-only" key:"stopSignal"`

	// User and group for setuid programs
	User  string `vic:"0.1" scope:"read-only" key:"user"`
	Group string `vic:"0.1" scope:"read-only" key:"group"`

	// if there's a pty then we need additional management data
	Pty       *os.File
	Outwriter dio.DynamicMultiWriter
	Errwriter dio.DynamicMultiWriter
	Reader    dio.DynamicMultiReader
}

type NetworkEndpoint struct {
	// Common.Name - the nic alias requested (only one name and one alias possible in linux)
	// Common.ID - pci slot of the vnic allowing for interface identifcation in-guest
	executor.Common

	// Whether this endpoint's IP was specified by the client (true if it was)
	Static bool `vic:"0.1" scope:"read-only" key:"static"`

	// IP address to assign
	IP *net.IPNet `vic:"0.1" scope:"read-only" key:"ip"`

	// Actual IP address assigned
	Assigned net.IPNet `vic:"0.1" scope:"read-write" key:"assigned"`

	// The network in which this information should be interpreted. This is embedded directly rather than
	// as a pointer so that we can ensure the data is consistent
	Network executor.ContainerNetwork `vic:"0.1" scope:"read-only" key:"network"`

	// DHCP runtime info
	DHCP *DHCPInfo `vic:"0.1" scope:"read-only" recurse:"depth=0"`
}

func (e *NetworkEndpoint) IsDynamic() bool {
	return !e.Static && (e.IP == nil || ip.IsUnspecifiedIP(e.IP.IP))
}

type DHCPInfo struct {
	Assigned    net.IPNet
	Nameservers []net.IP
	Gateway     net.IPNet
}
