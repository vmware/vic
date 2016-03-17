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

// NewVirtualSerialPort returns VirtualSerialPort spec.
func NewVirtualSerialPort() *types.VirtualSerialPort {
	defer trace.End(trace.Begin(""))

	return &types.VirtualSerialPort{
		VirtualDevice: types.VirtualDevice{},
		YieldOnPoll:   true,
	}
}

func (s *VirtualMachineConfigSpec) addVirtualSerialPort(device *types.VirtualSerialPort, debug bool, connected bool) *VirtualMachineConfigSpec {
	device.GetVirtualDevice().Key = s.generateNextKey()

	// Set serial device's backing to a datastore file when debug is true
	// We then instruct Linux kernel to use that as a serial console
	if !debug {
		device.GetVirtualDevice().Backing = &types.VirtualSerialPortURIBackingInfo{
			VirtualDeviceURIBackingInfo: types.VirtualDeviceURIBackingInfo{
				Direction:  string(types.VirtualDeviceURIBackingOptionDirectionClient),
				ServiceURI: s.ConnectorURI(),
			},
		}

		device.GetVirtualDevice().Connectable = &types.VirtualDeviceConnectInfo{
			Connected:         connected,
			StartConnected:    connected,
			AllowGuestControl: false,
		}
	} else {
		device.GetVirtualDevice().Backing = &types.VirtualSerialPortFileBackingInfo{
			VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
				FileName: fmt.Sprintf("%s/%s/%[2]s.debug", s.VMPathName(), s.ID()),
			},
		}
	}

	s.DeviceChange = append(s.DeviceChange,
		&types.VirtualDeviceConfigSpec{
			Operation: types.VirtualDeviceConfigSpecOperationAdd,
			Device:    device,
		},
	)
	return s
}

func (s *VirtualMachineConfigSpec) removeVirtualSerialPort(device types.BaseVirtualDevice) *VirtualMachineConfigSpec {
	s.DeviceChange = append(s.DeviceChange,
		&types.VirtualDeviceConfigSpec{
			Operation: types.VirtualDeviceConfigSpecOperationRemove,
			Device:    device,
		},
	)
	return s
}

// AddVirtualSerialPort adds a virtual serial port.
func (s *VirtualMachineConfigSpec) AddVirtualSerialPort(device *types.VirtualSerialPort) *VirtualMachineConfigSpec {
	defer trace.End(trace.Begin(s.ID()))

	return s.addVirtualSerialPort(device, false, false)
}

// AddVirtualConnectedSerialPort adds a connected virtual serial port.
func (s *VirtualMachineConfigSpec) AddVirtualConnectedSerialPort(device *types.VirtualSerialPort) *VirtualMachineConfigSpec {
	defer trace.End(trace.Begin(s.ID()))

	return s.addVirtualSerialPort(device, false, true)
}

// AddVirtualDebugSerialPort adds a file backed virtual serial port.
func (s *VirtualMachineConfigSpec) AddVirtualDebugSerialPort(device *types.VirtualSerialPort) *VirtualMachineConfigSpec {
	defer trace.End(trace.Begin(s.ID()))

	return s.addVirtualSerialPort(device, true, false)
}

// RemoveVirtualSerialPort removes a virtual serial port.
func (s *VirtualMachineConfigSpec) RemoveVirtualSerialPort(device *types.VirtualSerialPort) *VirtualMachineConfigSpec {
	defer trace.End(trace.Begin(s.ID()))

	return s.removeVirtualSerialPort(device)
}
