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

package exec2

// ContainerLifecycle represents operations concerned with the creation, modification
// and deletion of containers.
type ContainerLifecycle interface {

	// CreateContainer creates a new container representation and returns a Handle
	// The Handle can be used to configure the container before its actually created
	// Calling Commit on the Handle will create the container
	CreateContainer(name string) (Handle, error)

	// GetHandle allows for an existing container to be modified
	// The Handle can be used to reconfigure the container
	// Calling Commit on the Handle will apply the reconfiguration
	// Commit will fail if another client committed a modification after GetHandle was called
	GetHandle(cid ID) (Handle, error)

	// CopyTo copies a file into the container represented by the handle
	// If the container is stopped, the file will be copied in when it next runs
	CopyTo(handle Handle, targetDir string, fname string, data []byte) (Handle, error)

	// SetEntryPoint sets the entry point for the container
	// This is the executable, the lifecycle of which is tied to the container lifecycle
	SetEntryPoint(handle Handle, workDir string, execPath string, args string) (Handle, error)

	// ExecProcess executes a child process in the container
	// A child process will exit when the main process exits
	ExecProcess(handle Handle, execPath string, execArgs string) (Handle, error)

	// SetLimits sets resource limits on the container
	// A value of -1 implies a default value, not unlimited
	// New limits will be ignored if committed to a running container
	SetLimits(handle Handle, memoryMb int, cpuMhz int) (Handle, error)

	// SetRunState allows for the running state of a container to be modified
	// Created is not a valid state and will return an error
	SetRunState(handle Handle, runState RunState) (Handle, error)

	// Commit applies changes made to the Handle to either a new or running container
	// Commit will fail if another client committed a modification after GetHandle was called
	// Commit blocks until all changes have been committed
	Commit(handle Handle) (ID, error)

	// DestroyContainer destroys an existing container
	DestroyContainer(cid ID) error
}

// ContainerQuery represents queries that can be made against a Container or its Processes
type ContainerQuery interface {
	// ListContainers lists all container IDs for a given state
	// If forState is nil, all containers are returned
	ListContainers(forState RunState) ([]ID, error)

	// GetConfig returns container and process config
	// This is an efficient call as GetConfig will not query the infrastructure
	GetConfig(cID ID) (ContainerConfig, error)

	// GetState returns the current state of the container and its main process
	// This call will query the infrastructure and container to return a snapshot
	GetState(cID ID) (ContainerState, error)

	// GetChildProcesses returns the current state of container child processes
	// This call will query the container to return a snapshot
	GetChildProcesses(cid ID) ([]ProcessState, error)

	// CopyFrom copies file data out of a running container
	// Returns an error if the container is not running
	CopyFrom(cid ID, sourceDir string, fname string) ([]byte, error)
}

// RunState represents the running state of a container
type RunState int

const (
	_ RunState = iota
	Created
	Running
	Stopped
)

// ContainerConfig is a type representing the configuration of a container and its processes
type ContainerConfig struct {
	ConstantConfig
	Config
	process    ProcessConfig
	childProcs []ProcessConfig
}

// ContainerState is a type representing the runtime state of a container and its processes
type ContainerState struct {
	status  RunState
	process ProcessState
}

// ProcessState is the runtime state of a process in a container
type ProcessState struct {
	ProcessRunState
}
