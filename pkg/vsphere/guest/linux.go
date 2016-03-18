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

package guest

import (
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/spec"
	"golang.org/x/net/context"
)

const (
	linuxGuestID = "other3xLinux64Guest"

	scsiBusNumber = 0
	scsiKey       = 100
	ideKey        = 200
)

// LinuxGuestType type
type LinuxGuestType struct {
	*spec.VirtualMachineConfigSpec
}

// NewLinuxGuest returns a new Linux guest spec with predefined values
func NewLinuxGuest(ctx context.Context, session *session.Session, config *spec.VirtualMachineConfigSpecConfig) *LinuxGuestType {
	s := spec.NewVirtualMachineConfigSpec(ctx, session, config)

	// SCSI controller
	scsi := spec.NewVirtualSCSIController(scsiBusNumber, scsiKey)
	// PV SCSI controller
	pv := spec.NewParaVirtualSCSIController(scsi)
	s.AddParaVirtualSCSIController(pv)

	// IDE controller
	ide := spec.NewVirtualIDEController(ideKey)
	s.AddVirtualIDEController(ide)

	// CDROM
	cdrom := spec.NewVirtualCdrom(ide)
	s.AddVirtualCdrom(cdrom)

	// NIC
	vmxnet3 := spec.NewVirtualVmxnet3()
	s.AddVirtualVmxnet3(vmxnet3)

	// Tether serial port - backed by network
	serial := spec.NewVirtualSerialPort()
	s.AddVirtualConnectedSerialPort(serial)

	// Debug serial port - backed by datastore file
	debugserial := spec.NewVirtualSerialPort()
	s.AddVirtualDebugSerialPort(debugserial)

	// Set the guest id
	s.GuestId = linuxGuestID

	return &LinuxGuestType{
		VirtualMachineConfigSpec: s,
	}
}

// GuestID returns the guest id of the linux guest
func (l *LinuxGuestType) GuestID() string {
	return l.VirtualMachineConfigSpec.GuestId
}

// Spec returns the underlying types.VirtualMachineConfigSpec to the caller
func (l *LinuxGuestType) Spec() *types.VirtualMachineConfigSpec {
	return l.VirtualMachineConfigSpec.VirtualMachineConfigSpec
}
