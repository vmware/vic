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

package types

import (
	"os"
	"os/exec"
	"sync"

	"github.com/vmware/vic/pkg/dio"
)

type ExecutorConfig struct {
	// Sessions is the set of sessions currently hosted by this executor
	// These are keyed by session ID
	Sessions map[string]*SessionConfig `vic:"0.1" scope:"read-only" key:"sessions"`
}

type SessionConfig struct {
	// Exclusive access to the structure
	sync.Mutex `vic:"0.1" scope:"read-only" recurse:"depth=0"`

	// List of environment variable to set in the container
	Env []string
	// Command to run when starting the container
	Cmd []string
	// Current directory (PWD) in the command will be launched
	WorkingDir string

	// Allow attach or not
	Attach bool `vic:"0.1" scope:"read-only" key:"attach"`

	// Open stdin or not
	OpenStdin bool `vic:"0.1" scope:"read-only" key:"openstdin"`

	// Delay launching the Cmd until an attach request comes
	RunBlock bool `vic:"0.1" scope:"read-only" key:"runblock"`

	// Allocate a tty or not
	Tty bool `vic:"0.1" scope:"read-only" key:"tty"`

	// Restart controls whether a process gets relaunched if it exists
	Restart bool `vic:"0.1" scope:"read-only" key:"restart"`

	// StopSignal is the signal name or number used to stop a container
	StopSignal string `vic:"0.1" scope:"read-only" key:"stopSignal"`

	// User and group for setuid programs
	User  string `vic:"0.1" scope:"read-only" key:"user"`
	Group string `vic:"0.1" scope:"read-only" key:"group"`
}

type Session struct {
	// Exclusive access to the structure
	sync.Mutex `vic:"0.1" scope:"read-only" recurse:"depth=0"`

	*SessionConfig

	Interaction

	ID string

	// The primary process for the session
	Cmd exec.Cmd `vic:"0.1" scope:"read-only" key:"cmd" recurse:"depth=2,nofollow"`

	// The exit status of the process, if any
	ExitStatus int `vic:"0.1" scope:"read-write" key:"status"`

	// This indicates the launch status
	Started string `vic:"0.1" scope:"read-write" key:"started"`

	// RessurectionCount is a log of how many times the entity has been restarted due
	// to error exit
	ResurrectionCount int `vic:"0.1" scope:"read-write" key:"resurrections"`
}

type Interaction struct {
	// Exclusive access to the structure
	sync.Mutex `vic:"0.1" scope:"read-only" recurse:"depth=0"`

	Done <-chan struct{}

	Pty *os.File `vic:"0.1" scope:"read-only" recurse:"depth=0"`

	Outwriter dio.DynamicMultiWriter `vic:"0.1" scope:"read-only" recurse:"depth=0"`
	Errwriter dio.DynamicMultiWriter `vic:"0.1" scope:"read-only" recurse:"depth=0"`
	Reader    dio.DynamicMultiReader `vic:"0.1" scope:"read-only" recurse:"depth=0"`
}
