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
	"path"
	"strings"
	"time"

	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type VirtualMachine struct {
	mo.VirtualMachine
}

func NewVirtualMachine(spec *types.VirtualMachineConfigSpec) (*VirtualMachine, types.BaseMethodFault) {
	vm := &VirtualMachine{}

	if spec.Name == "" {
		return nil, &types.InvalidVmConfig{Property: "configSpec.name"}
	}

	if spec.Files == nil || spec.Files.VmPathName == "" {
		return nil, &types.InvalidVmConfig{Property: "configSpec.files.vmPathName"}
	}

	dsPath := spec.Files.VmPathName
	if strings.HasSuffix(dsPath, ".vmx") {
		dsPath = path.Dir(dsPath)
	}

	vm.Config = &types.VirtualMachineConfigInfo{}

	defaults := types.VirtualMachineConfigSpec{
		NumCPUs:           1,
		NumCoresPerSocket: 1,
		MemoryMB:          32,
		Uuid:              "TODO",
		Version:           "vmx-11",
		Files: &types.VirtualMachineFileInfo{
			SnapshotDirectory: dsPath,
			SuspendDirectory:  dsPath,
			LogDirectory:      dsPath,
		},
	}

	err := vm.configure(&defaults)
	if err != nil {
		return nil, err
	}

	vm.Runtime.PowerState = types.VirtualMachinePowerStatePoweredOff
	vm.Summary.Runtime = vm.Runtime

	err = vm.configure(spec)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func (vm *VirtualMachine) configure(spec *types.VirtualMachineConfigSpec) types.BaseMethodFault {
	apply := []struct {
		src string
		dst *string
	}{
		{spec.Name, &vm.Name},
		{spec.Name, &vm.Config.Name},
		{spec.GuestId, &vm.Config.GuestId},
		{spec.GuestId, &vm.Config.GuestFullName},
		{spec.Uuid, &vm.Config.Uuid},
		{spec.Version, &vm.Config.Version},
		{spec.Files.VmPathName, &vm.Config.Files.VmPathName},
		{spec.Files.SnapshotDirectory, &vm.Config.Files.SnapshotDirectory},
		{spec.Files.LogDirectory, &vm.Config.Files.LogDirectory},
	}

	for _, f := range apply {
		if f.src != "" {
			*f.dst = f.src
		}
	}

	if spec.MemoryMB != 0 {
		vm.Config.Hardware.MemoryMB = int32(spec.MemoryMB)
		vm.Summary.Config.MemorySizeMB = vm.Config.Hardware.MemoryMB
	}

	if spec.NumCPUs != 0 {
		vm.Config.Hardware.NumCPU = spec.NumCPUs
		vm.Summary.Config.NumCpu = vm.Config.Hardware.NumCPU
	}

	vm.Config.Modified = time.Now()

	return nil
}

type powerVMTask struct {
	*VirtualMachine

	state types.VirtualMachinePowerState
}

func (c *powerVMTask) Run(task *Task) (types.AnyType, types.BaseMethodFault) {
	if c.VirtualMachine.Runtime.PowerState == c.state {
		return nil, &types.InvalidPowerState{
			RequestedState: c.state,
			ExistingState:  c.VirtualMachine.Runtime.PowerState,
		}
	}

	c.VirtualMachine.Runtime.PowerState = c.state
	c.VirtualMachine.Summary.Runtime.PowerState = c.state

	return nil, nil
}

func (vm *VirtualMachine) PowerOnVMTask(c *types.PowerOnVM_Task) soap.HasFault {
	r := &methods.PowerOnVM_TaskBody{}

	task := NewTask(&powerVMTask{vm, types.VirtualMachinePowerStatePoweredOn})

	r.Res = &types.PowerOnVM_TaskResponse{
		Returnval: task.Self,
	}

	task.Run()

	return r
}

func (vm *VirtualMachine) PowerOffVMTask(c *types.PowerOffVM_Task) soap.HasFault {
	r := &methods.PowerOffVM_TaskBody{}

	task := NewTask(&powerVMTask{vm, types.VirtualMachinePowerStatePoweredOff})

	r.Res = &types.PowerOffVM_TaskResponse{
		Returnval: task.Self,
	}

	task.Run()

	return r
}

type destroyVMTask struct {
	*VirtualMachine
}

func (c *destroyVMTask) Run(task *Task) (types.AnyType, types.BaseMethodFault) {
	if c.VirtualMachine.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
		return nil, &types.InvalidPowerState{
			RequestedState: types.VirtualMachinePowerStatePoweredOff,
			ExistingState:  c.VirtualMachine.Runtime.PowerState,
		}
	}

	// TODO: remove references from HostSystem and Datastore
	Map.Remove(c.Reference())

	return nil, nil
}

func (vm *VirtualMachine) DestroyTask(c *types.Destroy_Task) soap.HasFault {
	r := &methods.Destroy_TaskBody{}

	task := NewTask(&destroyVMTask{vm})

	r.Res = &types.Destroy_TaskResponse{
		Returnval: task.Self,
	}

	task.Run()

	return r
}
