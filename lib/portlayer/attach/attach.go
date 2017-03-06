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

package attach

import (
	"fmt"
	"net"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"

	log "github.com/Sirupsen/logrus"
)

func lookupVCHIP() (net.IP, error) {
	// FIXME: THERE MUST BE ANOTHER WAY
	// following is from Create@exec.go
	ips, err := net.LookupIP(constants.ManagementHostName)
	if err != nil {
		return nil, err
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("No IP found on %s", constants.ManagementHostName)
	}

	if len(ips) > 1 {
		return nil, fmt.Errorf("Multiple IPs found on %s: %#v", constants.ManagementHostName, ips)
	}
	return ips[0], nil
}

func toggle(handle *exec.Handle, connected bool) (*exec.Handle, error) {
	// get the virtual device list
	devices := object.VirtualDeviceList(handle.Config.Hardware.Device)

	// select the virtual serial ports
	serials := devices.SelectByBackingInfo((*types.VirtualSerialPortURIBackingInfo)(nil))
	if len(serials) == 0 {
		return nil, fmt.Errorf("Unable to find a device with desired backing")
	}
	if len(serials) > 1 {
		return nil, fmt.Errorf("Multiple matches found with desired backing")
	}
	serial := serials[0]

	ip, err := lookupVCHIP()
	if err != nil {
		return nil, err
	}

	log.Debugf("Found a device with desired backing: %#v", serial)

	c := serial.GetVirtualDevice().Connectable
	b := serial.GetVirtualDevice().Backing.(*types.VirtualSerialPortURIBackingInfo)

	serviceURI := fmt.Sprintf("tcp://%s:%d", ip, constants.SerialOverLANPort)
	proxyURI := fmt.Sprintf("telnet://%s:%d", ip, constants.SerialOverLANPort)

	if b.ProxyURI == proxyURI && c.Connected == connected {
		log.Debugf("Already in the desired state (connected: %t, proxyURI: %s)", connected, proxyURI)
		return handle, nil
	}

	if b.ServiceURI == serviceURI && c.Connected == connected {
		log.Debugf("Already in the desired state (connected: %t, serviceURI: %s)", connected, serviceURI)
		return handle, nil
	}

	// set the values
	log.Debugf("Setting Connected to %t", connected)
	c.Connected = connected
	if connected && handle.ExecConfig.Sessions[handle.ExecConfig.ID].Attach {
		log.Debugf("Setting the start connected state to %t", connected)
		c.StartConnected = handle.ExecConfig.Sessions[handle.ExecConfig.ID].Attach
	}

	log.Debugf("Setting ServiceURI to %s", serviceURI)
	b.ServiceURI = serviceURI

	log.Debugf("Setting the ProxyURI to %s", proxyURI)
	b.ProxyURI = proxyURI

	config := &types.VirtualDeviceConfigSpec{
		Device:    serial,
		Operation: types.VirtualDeviceConfigSpecOperationEdit,
	}
	handle.Spec.DeviceChange = append(handle.Spec.DeviceChange, config)

	// iterate over Sessions and set their RunBlock property to connected
	// if attach happens before start then this property will be persist in the vmx
	// if attash happens after start then this propery will be thrown away by commit (one simply cannot change ExtraConfig if the vm is powered on)
	for _, session := range handle.ExecConfig.Sessions {
		session.RunBlock = connected
	}

	return handle, nil
}

// Join adds network backed serial port to the caller and configures them
func Join(h interface{}) (interface{}, error) {
	defer trace.End(trace.Begin(""))

	handle, ok := h.(*exec.Handle)
	if !ok {
		return nil, fmt.Errorf("Type assertion failed for %#+v", handle)
	}

	// Tether serial port - backed by network
	serial := &types.VirtualSerialPort{
		VirtualDevice: types.VirtualDevice{
			Backing: &types.VirtualSerialPortURIBackingInfo{
				VirtualDeviceURIBackingInfo: types.VirtualDeviceURIBackingInfo{
					Direction: string(types.VirtualDeviceURIBackingOptionDirectionClient),
					ProxyURI:  fmt.Sprintf("telnet://0.0.0.0:%d", constants.SerialOverLANPort),
					// Set it to 0.0.0.0 during Join call, VCH IP will be set when we call Bind
					ServiceURI: fmt.Sprintf("tcp://0.0.0.0:%d", constants.SerialOverLANPort),
				},
			},
			Connectable: &types.VirtualDeviceConnectInfo{
				Connected:         false,
				StartConnected:    false,
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

	return handle, nil
}

// Bind sets the *Connected fields of the VirtualSerialPort
func Bind(h interface{}) (interface{}, error) {
	defer trace.End(trace.Begin(""))

	handle, ok := h.(*exec.Handle)
	if !ok {
		return nil, fmt.Errorf("Type assertion failed for %#+v", handle)
	}
	if handle.MigrationError != nil {
		return nil, fmt.Errorf("Migration failed %s", handle.MigrationError)
	}
	return toggle(handle, true)
}

// Unbind unsets the *Connected fields of the VirtualSerialPort
func Unbind(h interface{}) (interface{}, error) {
	defer trace.End(trace.Begin(""))

	handle, ok := h.(*exec.Handle)
	if !ok {
		return nil, fmt.Errorf("Type assertion failed for %#+v", handle)
	}
	if handle.MigrationError != nil {
		return nil, fmt.Errorf("Migration failed %s", handle.MigrationError)
	}
	return toggle(handle, false)
}
