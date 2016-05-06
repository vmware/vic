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

package spec

import (
	"fmt"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
)

// VirtualMachineConfigSpecConfig holds the config values
type VirtualMachineConfigSpecConfig struct {
	// ID of the VM
	ID string

	// ParentImageID of the VM
	ParentImageID string

	// Name of the VM
	Name string

	// Number of CPUs
	NumCPUs int32
	// Memory - in MB
	MemoryMB int64

	// VMFork enabled
	VMForkEnabled bool

	// datastore path of the media file we boot from
	BootMediaPath string

	// datastore path of the VM
	VMPathName string

	// URI of the network serial port
	ConnectorURI string

	// Name of the network
	NetworkName string

	// Name of the image store
	ImageStoreName string

	// Temporary
	Metadata metadata.ExecutorConfig
}

// VirtualMachineConfigSpec type
type VirtualMachineConfigSpec struct {
	*session.Session

	*types.VirtualMachineConfigSpec

	config *VirtualMachineConfigSpecConfig

	// internal value to keep track of next ID
	key int32
}

// NewVirtualMachineConfigSpec returns a VirtualMachineConfigSpec
func NewVirtualMachineConfigSpec(ctx context.Context, session *session.Session, config *VirtualMachineConfigSpecConfig) (*VirtualMachineConfigSpec, error) {
	defer trace.End(trace.Begin(config.ID))

	VMPathName := config.VMPathName
	if !session.IsVSAN(ctx) {
		// VMFS requires the full path to vmx or everything but the datastore is ignored
		VMPathName = fmt.Sprintf("%s/%s/%[2]s.vmx", config.VMPathName, config.ID)
	}

	log.Debugf("Adding metadata to the configspec: %+v", config.Metadata)
	// TEMPORARY

	spec := &types.VirtualMachineConfigSpec{
		Name: config.ID,
		Files: &types.VirtualMachineFileInfo{
			VmPathName: VMPathName,
		},
		NumCPUs:             config.NumCPUs,
		CpuHotAddEnabled:    &config.VMForkEnabled, // this disables vNUMA when true
		MemoryMB:            config.MemoryMB,
		MemoryHotAddEnabled: &config.VMForkEnabled,

		// needed to cause the disk uuid to propogate into linux for presentation via /dev/disk/by-id/
		ExtraConfig: []types.BaseOptionValue{
			// lets us see the UUID for the containerfs disk (hidden from daemon)
			&types.OptionValue{Key: "disk.EnableUUID", Value: "true"},
			// needed to avoid the questions that occur when attaching multiple disks with the same uuid (bugzilla 1362918)
			&types.OptionValue{Key: "answer.msg.disk.duplicateUUID", Value: "Yes"},
			&types.OptionValue{Key: "answer.msg.serial.file.open", Value: "Replace"},

			&types.OptionValue{Key: "sched.mem.lpage.maxSharedPages", Value: "256"},
			// seems to be needed to avoid children hanging shortly after fork
			&types.OptionValue{Key: "vmotion.checkpointSVGAPrimarySize", Value: "4194304"},

			// trying this out - if it works then we need to determine if we can rely on serial0 being the correct index.
			&types.OptionValue{Key: "serial0.hardwareFlowControl", Value: "TRUE"},

			// https://enatai-jira.eng.vmware.com/browse/BON-257
			&types.OptionValue{Key: "memory.noHotAddOver4GB", Value: "FALSE"},
			&types.OptionValue{Key: "memory.maxGrow", Value: "512"},

			// http://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=2030189
			&types.OptionValue{Key: "tools.remindInstall", Value: "FALSE"},
			&types.OptionValue{Key: "tools.upgrade.policy", Value: "manual"},
		},
	}

	// encode the config as optionvalues
	cfg := map[string]string{}
	extraconfig.Encode(extraconfig.MapSink(cfg), config.Metadata)
	metaCfg := extraconfig.OptionValueFromMap(cfg)

	// merge it with the sec
	spec.ExtraConfig = append(spec.ExtraConfig, metaCfg...)

	return &VirtualMachineConfigSpec{
		Session:                  session,
		VirtualMachineConfigSpec: spec,
		config: config,
	}, nil
}

// AddVirtualDevice appends an Add operation to the DeviceChange list
func (s *VirtualMachineConfigSpec) AddVirtualDevice(device types.BaseVirtualDevice) *VirtualMachineConfigSpec {
	s.DeviceChange = append(s.DeviceChange,
		&types.VirtualDeviceConfigSpec{
			Operation: types.VirtualDeviceConfigSpecOperationAdd,
			Device:    device,
		},
	)
	return s
}

// AddAndCreateVirtualDevice appends an Add operation to the DeviceChange list
func (s *VirtualMachineConfigSpec) AddAndCreateVirtualDevice(device types.BaseVirtualDevice) *VirtualMachineConfigSpec {
	s.DeviceChange = append(s.DeviceChange,
		&types.VirtualDeviceConfigSpec{
			Operation:     types.VirtualDeviceConfigSpecOperationAdd,
			FileOperation: types.VirtualDeviceConfigSpecFileOperationCreate,
			Device:        device,
		},
	)
	return s
}

// RemoveVirtualDevice appends a Remove operation to the DeviceChange list
func (s *VirtualMachineConfigSpec) RemoveVirtualDevice(device types.BaseVirtualDevice) *VirtualMachineConfigSpec {
	s.DeviceChange = append(s.DeviceChange,
		&types.VirtualDeviceConfigSpec{
			Operation: types.VirtualDeviceConfigSpecOperationRemove,
			Device:    device,
		},
	)
	return s
}

// RemoveAndDestroyVirtualDevice appends a Remove operation to the DeviceChange list
func (s *VirtualMachineConfigSpec) RemoveAndDestroyVirtualDevice(device types.BaseVirtualDevice) *VirtualMachineConfigSpec {
	s.DeviceChange = append(s.DeviceChange,
		&types.VirtualDeviceConfigSpec{
			Operation:     types.VirtualDeviceConfigSpecOperationRemove,
			FileOperation: types.VirtualDeviceConfigSpecFileOperationDestroy,

			Device: device,
		},
	)
	return s
}

// Name returns the name of the VM
func (s *VirtualMachineConfigSpec) Name() string {
	defer trace.End(trace.Begin(s.config.Name))

	return s.config.Name
}

// ID returns the ID of the VM
func (s *VirtualMachineConfigSpec) ID() string {
	defer trace.End(trace.Begin(s.config.ID))

	return s.config.ID
}

// ParentImageID returns the ID of the image that VM is based on
func (s *VirtualMachineConfigSpec) ParentImageID() string {
	defer trace.End(trace.Begin(s.config.ParentImageID))

	return s.config.ParentImageID
}

// BootMediaPath returns the image path
func (s *VirtualMachineConfigSpec) BootMediaPath() string {
	defer trace.End(trace.Begin(s.config.ID))

	return s.config.BootMediaPath
}

// VMPathName returns the VM folder path
func (s *VirtualMachineConfigSpec) VMPathName() string {
	defer trace.End(trace.Begin(s.config.ID))

	return s.config.VMPathName
}

// NetworkName returns the network name
func (s *VirtualMachineConfigSpec) NetworkName() string {
	defer trace.End(trace.Begin(s.config.ID))

	return s.config.NetworkName
}

// ConnectorURI returns the connector URI
func (s *VirtualMachineConfigSpec) ConnectorURI() string {
	defer trace.End(trace.Begin(s.config.ID))

	return s.config.ConnectorURI
}

// ImageStoreName returns the image store name
func (s *VirtualMachineConfigSpec) ImageStoreName() string {
	defer trace.End(trace.Begin(s.config.ID))

	return s.config.ImageStoreName
}

func (s *VirtualMachineConfigSpec) generateNextKey() int32 {

	s.key -= 10
	return s.key
}
