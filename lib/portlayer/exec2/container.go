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

import (
	"time"

	"github.com/vmware/vic/pkg/vsphere/vm"
)

// All data in this struct must be data that is either immutable
// or can be relied upon without having to query either the container guest
// or the underlying infrastructure. Some of this state will be updated by events
type container struct {
	ConstantConfig

	vm       vm.VirtualMachine
	runState RunState

	config          Config
	processCfg      ProcessConfig // container main process
	childProcessCfg []ProcessConfig

	filesToCopy []FileToCopy // cache if copy while stopped
}

// config that will be applied to a container on commit
// Needs to be public as it will be shared by net, storage etc
type PendingCommit struct {
	ConstantConfig

	runState       RunState
	config         Config
	processConfig  ProcessConfig
	childProcesses []ProcessConfig
	filesToCopy    []FileToCopy
}

// config state that cannot change for the lifetime of the container
type ConstantConfig struct {
	CID     ID
	Created time.Time
	Limits  ResourceLimits
}

// variable container configuration state
type Config struct {
	Name    string
	WorkDir string
}

// configuration state of a container process
type ProcessConfig struct {
	PID      ID
	ExecPath string
	ExecArgs string
}

func NewProcessConfig(execPath string, execArgs string) ProcessConfig {
	return ProcessConfig{PID: GenerateID(), ExecArgs: execArgs, ExecPath: execPath}
}

type ProcessStatus int

const (
	Started = iota
	Exited
)

// runtime status of a container process
type ProcessRunState struct {
	PID        ID
	Status     ProcessStatus
	GuestPid   int
	ExitCode   int
	ExitErr    string
	StartedAt  time.Time
	FinishedAt time.Time
}

type FileToCopy struct {
	targetName string
	targetDir  string
	data       []byte
}

type ResourceLimits struct {
	MemoryMb int
	CPUMhz   int
}
