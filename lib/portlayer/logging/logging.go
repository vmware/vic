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

package logging

import (
	"fmt"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"
)

// Join adds two file backed serial port and configures them
func Join(h interface{}) (interface{}, error) {
	defer trace.End(trace.Begin(""))

	handle, ok := h.(*exec.Handle)
	if !ok {
		return nil, fmt.Errorf("Type assertion failed for %#+v", handle)
	}

	// make sure a spec exists
	handle.SetSpec(nil)

	VMPathName := handle.Spec.VMPathName()
	VMName := handle.Spec.Spec().Name

	for _, logFile := range []string{"tether.debug", "output.log"} {
		filename := fmt.Sprintf("%s/%s/%s", VMPathName, VMName, logFile)

		// Debug and log serial ports - backed by datastore file
		serial := &types.VirtualSerialPort{
			VirtualDevice: types.VirtualDevice{
				Backing: &types.VirtualSerialPortFileBackingInfo{
					VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
						FileName: filename,
					},
				},
				Connectable: &types.VirtualDeviceConnectInfo{
					Connected:         true,
					StartConnected:    true,
					AllowGuestControl: true,
				},
			},
			YieldOnPoll: true,
		}
		config := &types.VirtualDeviceConfigSpec{
			Device:    serial,
			Operation: types.VirtualDeviceConfigSpecOperationAdd,
		}
		handle.Spec.DeviceChange = append(handle.Spec.DeviceChange, config)
	}

	return handle, nil
}

// TODO: We can't really toggle the logging ports so bind/unbind are NOOP

// Bind sets the *Connected fields of the VirtualSerialPort
func Bind(h interface{}) (interface{}, error) {
	defer trace.End(trace.Begin(""))

	return h, nil
}

// Unbind unsets the *Connected fields of the VirtualSerialPort
func Unbind(h interface{}) (interface{}, error) {
	defer trace.End(trace.Begin(""))

	return h, nil
}
