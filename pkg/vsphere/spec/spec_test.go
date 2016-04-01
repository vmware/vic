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
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

func URL(t *testing.T) string {
	s := os.Getenv("TEST_URL")
	if s == "" {
		t.SkipNow()
	}
	return s
}

func TestVirtualMachineConfigSpec(t *testing.T) {

	ctx := context.Background()

	sessionconfig := &session.Config{
		Service:        URL(t),
		Insecure:       true,
		Keepalive:      time.Duration(5) * time.Minute,
		DatacenterPath: "",
		DatastorePath:  "/ha-datacenter/datastore/*",
		HostPath:       "/ha-datacenter/host/*/*",
		NetworkPath:    "/ha-datacenter/network/*",
		PoolPath:       "/ha-datacenter/host/*/Resources",
	}

	session, err := session.NewSession(sessionconfig).Create(ctx)
	if err != nil {
		t.Logf("%+v", err.Error())
		if _, ok := err.(*find.MultipleFoundError); !ok {
			t.Errorf(err.Error())
		} else {
			t.SkipNow()
		}
	}
	defer session.Logout(ctx)

	specconfig := &VirtualMachineConfigSpecConfig{
		NumCPUs:       2,
		MemoryMB:      2048,
		VMForkEnabled: true,

		ConnectorURI: "tcp://1.2.3.4:9876",

		ID: "zombie_attack",

		BootMediaPath: session.Datastore.Path("brainz.iso"),
		VMPathName:    fmt.Sprintf("[%s]", session.Datastore.Name()),
		NetworkName:   strings.Split(session.Network.Reference().Value, "-")[0],
	}
	// FIXME: find a better way to pass those
	scsibus := 0
	scsikey := 100
	idekey := 200

	root, _ := NewVirtualMachineConfigSpec(ctx, session, specconfig)
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
