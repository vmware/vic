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
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
)

type VirtualMachine struct {
	mo.VirtualMachine

	log *log.Logger
}

func NewVirtualMachine(spec *types.VirtualMachineConfigSpec) (*VirtualMachine, types.BaseMethodFault) {
	vm := &VirtualMachine{}

	if spec.Name == "" {
		return nil, &types.InvalidVmConfig{Property: "configSpec.name"}
	}

	if spec.Files == nil || spec.Files.VmPathName == "" {
		return nil, &types.InvalidVmConfig{Property: "configSpec.files.vmPathName"}
	}

	vm.Config = &types.VirtualMachineConfigInfo{}
	vm.Summary.Guest = &types.VirtualMachineGuestSummary{}
	vm.Summary.Storage = &types.VirtualMachineStorageSummary{}

	// Add the default devices
	devices, _ := object.VirtualDeviceList(esx.VirtualDevice).ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)

	if !strings.HasSuffix(spec.Files.VmPathName, ".vmx") {
		spec.Files.VmPathName = path.Join(spec.Files.VmPathName, spec.Name+".vmx")
	}

	dsPath := path.Dir(spec.Files.VmPathName)

	defaults := types.VirtualMachineConfigSpec{
		NumCPUs:           1,
		NumCoresPerSocket: 1,
		MemoryMB:          32,
		Uuid:              uuid.New().String(),
		Version:           "vmx-11",
		Files: &types.VirtualMachineFileInfo{
			SnapshotDirectory: dsPath,
			SuspendDirectory:  dsPath,
			LogDirectory:      dsPath,
		},
		DeviceChange: devices,
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
	err := vm.configureDevices(spec)
	if err != nil {
		return err
	}

	apply := []struct {
		src string
		dst *string
	}{
		{spec.Name, &vm.Name},
		{spec.Name, &vm.Config.Name},
		{spec.Name, &vm.Summary.Config.Name},
		{spec.GuestId, &vm.Config.GuestId},
		{spec.GuestId, &vm.Config.GuestFullName},
		{spec.GuestId, &vm.Summary.Guest.GuestId},
		{spec.GuestId, &vm.Summary.Config.GuestId},
		{spec.GuestId, &vm.Summary.Config.GuestFullName},
		{spec.Uuid, &vm.Config.Uuid},
		{spec.Version, &vm.Config.Version},
		{spec.Files.VmPathName, &vm.Config.Files.VmPathName},
		{spec.Files.SnapshotDirectory, &vm.Config.Files.SnapshotDirectory},
		{spec.Files.LogDirectory, &vm.Config.Files.LogDirectory},
	}

	for _, f := range apply {
		if *f.dst == "" {
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

func (vm *VirtualMachine) useDatastore(name string) *Datastore {
	host := Map.Get(*vm.Runtime.Host).(*HostSystem)

	ds := Map.FindByName(name, host.Datastore).(*Datastore)

	if Map.FindByName(name, vm.Datastore) == nil {
		vm.Datastore = append(vm.Datastore, ds.Reference())
	}

	return ds
}

func (vm *VirtualMachine) setLog(w io.Writer) {
	vm.log = log.New(w, "vmx ", log.Flags())
}

func (vm *VirtualMachine) createFile(spec string, name string) (*os.File, types.BaseMethodFault) {
	p, fault := parseDatastorePath(spec)
	if fault != nil {
		return nil, fault
	}

	ds := vm.useDatastore(p.Datastore)

	file := path.Join(ds.Info.GetDatastoreInfo().Url, p.Path)

	if name != "" {
		if path.Ext(file) != "" {
			file = path.Dir(file)
		}

		file = path.Join(file, name)
	}

	dir := path.Dir(file)

	_ = os.MkdirAll(dir, 0700)

	_, err := os.Stat(file)
	if err == nil {
		return nil, &types.FileAlreadyExists{
			FileFault: types.FileFault{
				File: file,
			},
		}
	}

	f, err := os.Create(file)
	if err != nil {
		return nil, &types.FileFault{
			File: file,
		}
	}

	return f, nil
}

func (vm *VirtualMachine) create(spec *types.VirtualMachineConfigSpec) types.BaseMethodFault {
	files := []struct {
		spec string
		name string
		use  func(w io.Writer)
	}{
		{vm.Config.Files.VmPathName, "", nil},
		{vm.Config.Files.VmPathName, fmt.Sprintf("%s.nvram", vm.Name), nil},
		{vm.Config.Files.LogDirectory, "vmware.log", vm.setLog},
	}

	for _, file := range files {
		f, err := vm.createFile(file.spec, file.name)
		if err != nil {
			return err
		}

		if file.use != nil {
			file.use(f)
		} else {
			_ = f.Close()
		}
	}

	vm.log.Print("created")

	return nil
}

func (vm *VirtualMachine) configureDevices(spec *types.VirtualMachineConfigSpec) types.BaseMethodFault {
	devices := object.VirtualDeviceList(vm.Config.Hardware.Device)

	for i, change := range spec.DeviceChange {
		dspec := change.GetVirtualDeviceConfigSpec()
		device := dspec.Device.GetVirtualDevice()
		invalid := &types.InvalidDeviceSpec{DeviceIndex: int32(i)}

		switch dspec.Operation {
		case types.VirtualDeviceConfigSpecOperationAdd:
			if devices.FindByKey(device.Key) != nil {
				return invalid
			}
			devices = append(devices, device)
		}
	}

	vm.Config.Hardware.Device = []types.BaseVirtualDevice(devices)

	return nil
}

type powerVMTask struct {
	*VirtualMachine

	state types.VirtualMachinePowerState
}

func (c *powerVMTask) Run(task *Task) (types.AnyType, types.BaseMethodFault) {
	c.log.Printf("running power task: requesting %s, existing %s",
		c.state, c.VirtualMachine.Runtime.PowerState)

	if c.VirtualMachine.Runtime.PowerState == c.state {
		return nil, &types.InvalidPowerState{
			RequestedState: c.state,
			ExistingState:  c.VirtualMachine.Runtime.PowerState,
		}
	}

	c.VirtualMachine.Runtime.PowerState = c.state
	c.VirtualMachine.Summary.Runtime.PowerState = c.state

	bt := &c.VirtualMachine.Summary.Runtime.BootTime
	if c.state == types.VirtualMachinePowerStatePoweredOn {
		now := time.Now()
		*bt = &now
	} else {
		*bt = nil
	}

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
