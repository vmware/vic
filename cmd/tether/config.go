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

package main

import (
	"os"
	"os/exec"
	"sync"

	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/dio"
)

type ExecutorConfig struct {
	// ID corresponds to that of the primary session
	ID string `vic:"0.1" scope:"read-only" key:"common/id"`

	// Exclusive access to childPidTable
	pidMutex sync.Mutex

	// Set of child PIDs created by us.
	pids map[int]*SessionConfig

	// Sessions is the set of sessions currently hosted by this executor
	// These are keyed by session ID
	Sessions map[string]*SessionConfig `vic:"0.1" scope:"read-only" key:"sessions"`

	// Maps the mount name to the detail mount specification
	Mounts map[string]metadata.MountSpec `vic:"0.1" scope:"read-only" key:"mounts"`

	// This describes an executors presence on a network, and contains sufficient
	// information to configure the interface in the guest.
	Networks map[string]*metadata.NetworkEndpoint `vic:"0.1" scope:"read-only" key:"networks"`

	// Key is the host key used during communicate back with the Interaction endpoint if any
	// Used if the in-guest tether is responsible for authenticating the connection
	Key []byte `vic:"0.1" scope:"read-only" key:"key"`
}

// SessionConfig defines the content of a session - this maps to the root of a process tree
// inside an executor
// This is close to but not perfectly aligned with the new docker/docker/daemon/execdriver/driver:CommonProcessConfig
type SessionConfig struct {
	// The primary session may have the same ID as the executor owning it
	metadata.Common `vic:"0.1" scope:"read-only" key:"common"`

	// The primary process for the session
	Cmd exec.Cmd `vic:"0.1" scope:"read-only" key:"cmd" recurse:"depth=2,nofollow"`

	// The exit status of the process, if any
	ExitStatus int `vic:"0.1" scope:"read-write" key:"status"`

	Started string `vic:"0.1" scope:"read-write" key:"started"`

	// Allow attach
	Attach bool `vic:"0.1" scope:"read-only" key:"attach"`

	// Allocate a tty or not
	Tty bool `vic:"0.1" scope:"read-only" key:"tty"`

	// if there's a pty then we need additional management data
	pty       *os.File
	outwriter dio.DynamicMultiWriter
	errwriter dio.DynamicMultiWriter
	reader    dio.DynamicMultiReader
}
