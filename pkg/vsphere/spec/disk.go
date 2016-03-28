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

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/trace"
)

const (
	//from portlayer/vsphere/storage/store.go
	defaultCapacityInKB = 8 * 1024 * 1024
)

// NewVirtualDisk returns a new disk attached to the controller
func NewVirtualDisk(controller types.BaseVirtualController) *types.VirtualDisk {

	defer trace.End(trace.Begin(""))

	return &types.VirtualDisk{
		VirtualDevice: types.VirtualDevice{
			ControllerKey: controller.GetVirtualController().Key,
			// Zero value is a valid UnitNumber. Set it to -1 so that ESXi/vCenter can assign one for us.
			UnitNumber: -1,
		},
	}
}

// NewVirtualSCSIDisk returns a new disk attached to the SCSI controller
func NewVirtualSCSIDisk(controller types.VirtualSCSIController) *types.VirtualDisk {
	defer trace.End(trace.Begin(""))

	return NewVirtualDisk(&controller)
}

// NewVirtualIDEDisk returns a new disk attached to the IDE controller
func NewVirtualIDEDisk(controller types.VirtualIDEController) *types.VirtualDisk {
	defer trace.End(trace.Begin(""))

	return NewVirtualDisk(&controller)
}

// AddVirtualDisk adds a virtual disk to a virtual machine.
func (s *VirtualMachineConfigSpec) AddVirtualDisk(device *types.VirtualDisk) *VirtualMachineConfigSpec {
	defer trace.End(trace.Begin(s.ID()))

	device.GetVirtualDevice().Key = s.generateNextKey()

	device.CapacityInKB = defaultCapacityInKB

	moref := s.Datastore.Reference()
	device.GetVirtualDevice().Backing = &types.VirtualDiskFlatVer2BackingInfo{
		DiskMode:        string(types.VirtualDiskModePersistent),
		ThinProvisioned: types.NewBool(true),

		VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
			FileName:  s.Datastore.Path(fmt.Sprintf("%s/%[1]s.vmdk", s.ID())),
			Datastore: &moref,
		},

		Parent: &types.VirtualDiskFlatVer2BackingInfo{
			VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
				FileName: s.Datastore.Path(fmt.Sprintf("VIC/%s/%s/%[2]s.vmdk", s.ImageStoreName(), s.ParentImageID())),
			},
		},
	}

	return s.AddAndCreateVirtualDevice(device)
}

// RemoveVirtualDisk remvoes the virtual disk from a virtual machine.
func (s *VirtualMachineConfigSpec) RemoveVirtualDisk(device *types.VirtualDisk) *VirtualMachineConfigSpec {
	defer trace.End(trace.Begin(s.ID()))

	return s.RemoveAndDestroyVirtualDevice(device)
}
