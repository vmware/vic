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

package vicbackends

import (
	//	"fmt"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	strfmt "github.com/go-swagger/go-swagger/strfmt"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
)

const (
	containerID      = "7a4ae369f2862b5fd9fdf08e87dd582c41e3e5623e33d20c02accec25ab94b0b"
	imageID          = "6a07b6ab2384072474d0cd00eee9fdb936afe8cf43d012be42fd1f6f1f22b163"
	hostArch         = "x86"
	hostKernel       = "0.3"
	hostOs           = "Linux"
	resCPUCount      = 2
	resMemory        = 1024
	containerLogPath = "/vmfs/volumes/datastore1/7a4ae369f2862b5fd9fdf08e87dd582c41e3e5623e33d20c02accec25ab94b0b"
	containerState   = "RUNNING"
	procEnv          = "PATH=/bin"
	procErrMsg       = ""
	procExecArgs     = ""
	procExecPath     = "/bin/bash"
	procExitCode     = 0
	procStatus       = "RUNNING"
	procPid          = 1
	procWorkDir      = "/"
	scopeID          = "bridge"
	scopeGateway     = "127.0.0.1"
	scopeDefault     = "default"
	scopeName        = "default"
	scopeType        = "bridge"
	scopeSubnet      = "255.255.255.0"
	mountPoint       = "/test"
	mountRW          = true
)

var (
	scopeDNS    = []string{"192.168.218.1", "192.168.100.1"}
	mountLabels = []string{"test", "dummy"}
)

func init() {
	// Silence all non-fatal logs during unit testing
	log.SetLevel(log.FatalLevel)
}

//----------
// Mock data
//----------

func newStringPointer(str string) *string {
	ns := new(string)
	*ns = str
	return ns
}

func newInt64Pointer(num int64) *int64 {
	nn := new(int64)
	*nn = num
	return nn
}

func newInt32Pointer(num int32) *int32 {
	nn := new(int32)
	*nn = num
	return nn
}

func newBoolPointer(val bool) *bool {
	nb := new(bool)
	*nb = val
	return nb
}

func newTimePointer(val time.Time) *strfmt.DateTime {
	nt := new(strfmt.DateTime)
	*nt = strfmt.DateTime(val)
	return nt
}

func getMockFullHostConfig() *models.HostConfig {
	return &models.HostConfig{
		Architecture:  newStringPointer(hostArch),
		KernelVersion: newStringPointer(hostKernel),
		Ostype:        newStringPointer(hostOs),
		Reservation: &models.ReservationConfig{
			CPUCount:    newInt64Pointer(resCPUCount),
			MemoryLimit: newInt64Pointer(resMemory),
		},
	}
}

func getMockFullContainerConfig(timeVal time.Time) *models.ContainerConfig {
	return &models.ContainerConfig{
		AttachStderr: newBoolPointer(true),
		AttachStdin:  newBoolPointer(true),
		AttachStdout: newBoolPointer(true),
		ConsoleSize: &models.ContainerConfigConsoleSize{
			Height: newInt64Pointer(0),
			Width:  newInt64Pointer(0),
		},
		ContainerID: newStringPointer(containerID),
		Created:     newTimePointer(timeVal),
		HostName:    newStringPointer(containerID),
		ImageID:     newStringPointer(imageID),
		Labels: map[string]string{
			"label1": "dummy",
		},
		LogPath:   newStringPointer(containerLogPath),
		Name:      []string{"dummy-name"},
		OpenStdin: newBoolPointer(false),
		Reservation: &models.ReservationConfig{
			CPUCount:    newInt64Pointer(2),
			MemoryLimit: newInt64Pointer(1024),
		},
		RestartCount: newInt32Pointer(2),
		State:        newStringPointer(containerState),
		Tty:          newBoolPointer(true),
	}
}

func getMockFullProcessConfig(timeVal time.Time) *models.ProcessConfig {
	return &models.ProcessConfig{
		Env:        newStringPointer(procEnv),
		ErrorMsg:   newStringPointer(procErrMsg),
		ExecArgs:   newStringPointer(procExecArgs),
		ExecPath:   newStringPointer(procExecPath),
		ExitCode:   newInt32Pointer(procExitCode),
		Finished:   newTimePointer(timeVal),
		Pid:        newInt32Pointer(procPid),
		Started:    newTimePointer(timeVal),
		Status:     newStringPointer(procStatus),
		WorkingDir: newStringPointer(procWorkDir),
	}
}

func getMockFullScopeConfig() *models.ScopeConfig {
	newScope := &models.ScopeConfig{
		Gateway:   newStringPointer(scopeGateway),
		ID:        newStringPointer(scopeID),
		IPAM:      nil,
		Name:      scopeDefault,
		ScopeType: scopeType,
		Subnet:    newStringPointer(scopeSubnet),
	}

	newScope.DNS = make([]string, len(scopeDNS))
	copy(newScope.DNS, scopeDNS)

	return newScope
}

func getMockFullVolumeConfig() *models.VolumeConfig {
	newVolume := &models.VolumeConfig{
		MountPoint: newStringPointer(mountPoint),
		ReadWrite:  newBoolPointer(mountRW),
	}

	newVolume.Label = make([]string, len(mountLabels))
	copy(newVolume.Label, mountLabels)

	return newVolume
}

func getMockFullContainerInfo(timeVal time.Time) *models.ContainerInfo {
	info := &models.ContainerInfo{}
	info.HostConfig = getMockFullHostConfig()
	info.ContainerConfig = getMockFullContainerConfig(timeVal)
	info.ProcessConfig = getMockFullProcessConfig(timeVal)

	info.ScopeConfig = make([]*models.ScopeConfig, 1)
	info.ScopeConfig[0] = getMockFullScopeConfig()

	info.VolumeConfig = make([]*models.VolumeConfig, 1)
	info.VolumeConfig[0] = getMockFullVolumeConfig()

	return info
}

func getMockEmptyContainerInfo() *models.ContainerInfo {
	info := &models.ContainerInfo{}
	info.HostConfig = &models.HostConfig{}
	info.ContainerConfig = &models.ContainerConfig{}
	info.ProcessConfig = &models.ProcessConfig{}
	info.ScopeConfig = make([]*models.ScopeConfig, 0)
	info.VolumeConfig = make([]*models.VolumeConfig, 0)

	return info
}

//----------
// Unit tests
//----------

// Tests empty swagger container info to container inspect json
// conversion.
func TestEmptyContainerInfoToContainerInspect(t *testing.T) {
	info := getMockEmptyContainerInfo()

	inspectJSON, err := containerInfoToDockerContainerInspect(containerID, info)

	assert.NotEqual(t, inspectJSON, nil, "Failed to convert empty container info to inspect json")
	assert.Equal(t, err, nil, "Received error getting container inspect json with empty container info: %s", err)
}

// Tests full swagger container info to container inspect json
// conversion.
func TestFullContainerInfoToContainerInspect(t *testing.T) {
	now := time.Now()
	info := getMockFullContainerInfo(now)

	inspectJSON, err := containerInfoToDockerContainerInspect(containerID, info)

	assert.NotEqual(t, inspectJSON, nil, "Failed to convert empty container info to inspect json")
	assert.Equal(t, err, nil, "Received error getting container inspect json: %s", err)
}

// Test partial swagger container info to container inspect json
// conversion.
func TestPartialContainerInfoToContainerInspect(t *testing.T) {
	now := time.Now()

	//empty out host config
	{
		info := getMockFullContainerInfo(now)
		info.HostConfig = nil

		inspectJSON, err := containerInfoToDockerContainerInspect(containerID, info)

		assert.NotEqual(t, inspectJSON, nil, "Failed to get back a container inspect json")
		assert.Equal(t, err, nil, "Received error getting container inspect json when HostConfig is nil: %s", err)
	}

	//empty out process config
	{
		info := getMockFullContainerInfo(now)
		info.ProcessConfig = nil

		inspectJSON, err := containerInfoToDockerContainerInspect(containerID, info)

		assert.NotEqual(t, inspectJSON, nil, "Failed to get back a container inspect json")
		assert.Equal(t, err, nil, "Received error getting container inspect json when ProcessConfig is nil: %s", err)
	}

	//empty out container config
	{
		info := getMockFullContainerInfo(now)
		info.ContainerConfig = nil

		inspectJSON, err := containerInfoToDockerContainerInspect(containerID, info)

		assert.NotEqual(t, inspectJSON, nil, "Failed to get back a container inspect json")
		assert.NotEqual(t, err, nil, "Expected error for inspect json when ContainerConfig is nil: %s", err)
	}

	//empty out scope config
	{
		info := getMockFullContainerInfo(now)
		info.ScopeConfig = nil

		inspectJSON, err := containerInfoToDockerContainerInspect(containerID, info)

		assert.NotEqual(t, inspectJSON, nil, "Failed to get back a container inspect json")
		assert.Equal(t, err, nil, "Received error getting container inspect json when ScopeConfig is nil: %s", err)
	}

	//empty out volume config
	{
		info := getMockFullContainerInfo(now)
		info.VolumeConfig = nil

		inspectJSON, err := containerInfoToDockerContainerInspect(containerID, info)

		assert.NotEqual(t, inspectJSON, nil, "Failed to get back a container inspect json")
		assert.Equal(t, err, nil, "Received error getting container inspect json when VolumeConfig is nil: %s", err)
	}
}
