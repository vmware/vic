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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/test/env"
	"golang.org/x/net/context"
)

func TestVirtualMachineConfigSpec(t *testing.T) {

	ctx := context.Background()

	sessionconfig := &session.Config{
		Service:        env.URL(t),
		Insecure:       true,
		Keepalive:      time.Duration(5) * time.Minute,
		DatacenterPath: "",
		DatastorePath:  "/ha-datacenter/datastore/*",
		HostPath:       "/ha-datacenter/host/*/*",
		NetworkPath:    "/ha-datacenter/network/VM Network",
		PoolPath:       "/ha-datacenter/host/*/Resources",
	}

	s, err := session.NewSession(sessionconfig).Create(ctx)
	if err != nil {
		t.Logf("%+v", err.Error())
		if _, ok := err.(*find.MultipleFoundError); !ok {
			t.Errorf(err.Error())
		} else {
			t.SkipNow()
		}
	}
	defer s.Logout(ctx)

	specconfig := &VirtualMachineConfigSpecConfig{
		NumCPUs:       2,
		MemoryMB:      2048,
		VMForkEnabled: true,

		ConnectorURI: "tcp://1.2.3.4:9876",

		ID: "zombie_attack",

		BootMediaPath: s.Datastore.Path("brainz.iso"),
		VMPathName:    fmt.Sprintf("[%s]", s.Datastore.Name()),
		NetworkName:   strings.Split(s.Network.Reference().Value, "-")[0],
	}
	// FIXME: find a better way to pass those
	var scsibus int32
	var scsikey int32 = 100
	var idekey int32 = 200

	root, _ := NewVirtualMachineConfigSpec(ctx, s, specconfig)
	scsi := NewVirtualSCSIController(scsibus, scsikey)

	pv := NewParaVirtualSCSIController(scsi)
	root.AddParaVirtualSCSIController(pv)

	bl := NewVirtualBusLogicController(scsi)
	root.AddVirtualBusLogicController(bl)

	ll := NewVirtualLsiLogicController(scsi)
	root.AddVirtualLsiLogicController(ll)

	ls := NewVirtualLsiLogicSASController(scsi)
	root.AddVirtualLsiLogicSASController(ls)
	///
	ide := NewVirtualIDEController(idekey)
	root.AddVirtualIDEController(ide)

	cdrom := NewVirtualCdrom(ide)
	root.AddVirtualCdrom(cdrom)

	floppy := NewVirtualFloppy(ide)
	root.AddVirtualFloppy(floppy)

	vmxnet3 := NewVirtualVmxnet3()
	root.AddVirtualVmxnet3(vmxnet3)

	pcnet32 := NewVirtualPCNet32()
	root.AddVirtualPCNet32(pcnet32)

	e1000 := NewVirtualE1000()
	root.AddVirtualE1000(e1000)

	serial := NewVirtualSerialPort()
	root.AddVirtualSerialPort(serial)

	debugserial := NewVirtualSerialPort()
	root.AddVirtualFileSerialPort(debugserial, "debug")

	for i := 0; i < len(root.DeviceChange); i++ {
		t.Logf("%+v", root.DeviceChange[i].GetVirtualDeviceConfigSpec().Device)
	}

}

func TestCollectSlotNumbers(t *testing.T) {
	s := &VirtualMachineConfigSpec{
		config: &VirtualMachineConfigSpecConfig{
			ID: "foo",
		},
		VirtualMachineConfigSpec: &types.VirtualMachineConfigSpec{},
	}

	slots := s.CollectSlotNumbers()
	assert.Empty(t, slots)

	s.AddVirtualVmxnet3(NewVirtualVmxnet3())
	s.DeviceChange[0].GetVirtualDeviceConfigSpec().Device.GetVirtualDevice().SlotInfo = &types.VirtualDevicePciBusSlotInfo{PciSlotNumber: 32}
	slots = s.CollectSlotNumbers()
	assert.EqualValues(t, []int32{32}, slots)

	// add a device without a slot number
	s.AddVirtualVmxnet3(NewVirtualVmxnet3())
	slots = s.CollectSlotNumbers()
	assert.EqualValues(t, []int32{32}, slots)

	// add another device with slot number
	s.AddVirtualVmxnet3(NewVirtualVmxnet3())
	s.DeviceChange[len(s.DeviceChange)-1].GetVirtualDeviceConfigSpec().Device.GetVirtualDevice().SlotInfo = &types.VirtualDevicePciBusSlotInfo{PciSlotNumber: 33}
	slots = s.CollectSlotNumbers()
	assert.EqualValues(t, []int32{32, 33}, slots)
}

func TestFindSlotNumber(t *testing.T) {
	allSlots := make(map[int32]bool)
	for s := pciSlotNumberBegin; s != pciSlotNumberEnd; s += pciSlotNumberInc {
		allSlots[s] = true
	}

	// missing first slot
	missingFirstSlot := make(map[int32]bool)
	for s := pciSlotNumberBegin + pciSlotNumberInc; s != pciSlotNumberEnd; s += pciSlotNumberInc {
		missingFirstSlot[s] = true
	}

	// missing last slot
	missingLastSlot := make(map[int32]bool)
	for s := pciSlotNumberBegin; s != pciSlotNumberEnd-pciSlotNumberInc; s += pciSlotNumberInc {
		missingLastSlot[s] = true
	}

	// missing a slot in the middle
	var missingSlot int32
	missingMiddleSlot := make(map[int32]bool)
	for s := pciSlotNumberBegin; s != pciSlotNumberEnd-pciSlotNumberInc; s += pciSlotNumberInc {
		if pciSlotNumberBegin+(2*pciSlotNumberInc) == s {
			missingSlot = s
			continue
		}
		missingMiddleSlot[s] = true
	}

	var tests = []struct {
		slots map[int32]bool
		out   int32
	}{
		{make(map[int32]bool), pciSlotNumberBegin},
		{allSlots, spec.NilSlot},
		{missingFirstSlot, pciSlotNumberBegin},
		{missingLastSlot, pciSlotNumberEnd - pciSlotNumberInc},
		{missingMiddleSlot, missingSlot},
	}

	for _, te := range tests {
		if s := findSlotNumber(te.slots); s != te.out {
			t.Fatalf("findSlotNumber(%v) => %d, want %d", te.slots, s, te.out)
		}
	}
}
