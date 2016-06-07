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

package simulator

import (
	"fmt"
	"testing"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
	"golang.org/x/net/context"
)

func TestCreateVmESX(t *testing.T) {
	content := esx.ServiceContent
	s := New(NewServiceInstance(content, esx.RootFolder))

	ts := s.NewServer()
	defer ts.Close()

	ctx := context.Background()

	c, err := govmomi.NewClient(ctx, ts.URL, true)
	if err != nil {
		t.Fatal(err)
	}

	dc := object.NewDatacenter(c.Client, esx.Datacenter.Reference())

	folders, err := dc.Folders(ctx)
	if err != nil {
		t.Fatal(err)
	}

	pool := object.NewResourcePool(c.Client, esx.ResourcePool.Reference())

	spec := types.VirtualMachineConfigSpec{
		// Note: real ESX allows the VM to be created without a GuestId,
		// but will power on will fail.
		GuestId: string(types.VirtualMachineGuestOsIdentifierOtherGuest),
	}

	steps := []func(){
		func() {
			spec.Name = "test"
		},
		func() {
			spec.Files = &types.VirtualMachineFileInfo{
				VmPathName: fmt.Sprintf("[] %s/%s.vmx", spec.Name, spec.Name),
			}
		},
	}

	// expecting CreateVM to fail until all steps are taken
	for _, step := range steps {
		task, cerr := folders.VmFolder.CreateVM(ctx, spec, pool, nil)
		if cerr != nil {
			t.Fatal(err)
		}

		_, cerr = task.WaitForResult(ctx, nil)
		if cerr == nil {
			t.Error("expected error")
		}

		step()
	}

	task, err := folders.VmFolder.CreateVM(ctx, spec, pool, nil)
	if err != nil {
		t.Fatal(err)
	}

	info, err := task.WaitForResult(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	vm := object.NewVirtualMachine(c.Client, info.Result.(types.ManagedObjectReference))

	name, err := vm.ObjectName(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if name != spec.Name {
		t.Errorf("name=%s", name)
	}

	state, err := vm.PowerState(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if state != types.VirtualMachinePowerStatePoweredOff {
		t.Errorf("state=%s", state)
	}

	ops := []struct {
		method func(context.Context) (*object.Task, error)
		state  types.VirtualMachinePowerState
		fail   bool
	}{
		// Off -> On  == ok
		{vm.PowerOn, types.VirtualMachinePowerStatePoweredOn, false},
		// On  -> On  == fail
		{vm.PowerOn, types.VirtualMachinePowerStatePoweredOn, true},
		// On  -> Off == ok
		{vm.PowerOff, types.VirtualMachinePowerStatePoweredOff, false},
		// Off -> Off == fail
		{vm.PowerOff, types.VirtualMachinePowerStatePoweredOff, true},
		// Off -> On  == ok
		{vm.PowerOn, types.VirtualMachinePowerStatePoweredOn, false},
		// Destroy == fail (power is On)
		{vm.Destroy, types.VirtualMachinePowerStatePoweredOn, true},
		// On  -> Off == ok
		{vm.PowerOff, types.VirtualMachinePowerStatePoweredOff, false},
		// Destroy == ok (power is Off)
		{vm.Destroy, "", false},
	}

	for _, op := range ops {
		task, err = op.method(ctx)
		if err != nil {
			t.Fatal(err)
		}

		err = task.Wait(ctx)
		if op.fail {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
		}

		if len(op.state) != 0 {
			state, err = vm.PowerState(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if state != op.state {
				t.Errorf("state=%s", state)
			}
		}
	}
}
