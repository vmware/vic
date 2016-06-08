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
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
	"github.com/vmware/vic/pkg/vsphere/simulator/vc"
	"golang.org/x/net/context"
)

func TestCreateVm(t *testing.T) {
	models := []func() (*vmCreateModel, error){
		vmCreateModelWithVC,
		vmCreateModelWithESX,
	}

	ctx := context.Background()

	for _, model := range models {
		m, err := model()
		defer m.Server.Close()

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
					VmPathName: fmt.Sprintf("[%s] %s/%s.vmx", m.Datastore.Name(), spec.Name, spec.Name),
				}
			},
		}

		vmFolder := m.Folders.VmFolder
		// expecting CreateVM to fail until all steps are taken
		for _, step := range steps {
			task, err := vmFolder.CreateVM(ctx, spec, m.Pool, m.Host)
			if err != nil {
				t.Fatal(err)
			}

			_, err = task.WaitForResult(ctx, nil)
			if err == nil {
				t.Error("expected error")
			}

			step()
		}

		task, err := vmFolder.CreateVM(ctx, spec, m.Pool, m.Host)
		if err != nil {
			t.Fatal(err)
		}

		info, err := task.WaitForResult(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}

		vm := object.NewVirtualMachine(m.Client.Client, info.Result.(types.ManagedObjectReference))

		name, err := vm.ObjectName(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if name != spec.Name {
			t.Errorf("name=%s", name)
		}

		_, err = vm.Device(ctx)
		if err != nil {
			t.Fatal(err)
		}

		recreate := func(context.Context) (*object.Task, error) {
			return vmFolder.CreateVM(ctx, spec, m.Pool, nil)
		}

		ops := []struct {
			method func(context.Context) (*object.Task, error)
			state  types.VirtualMachinePowerState
			fail   bool
		}{
			// Powered off by default
			{nil, types.VirtualMachinePowerStatePoweredOff, false},
			// Create with same .vmx path should fail
			{recreate, "", true},
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

		for i, op := range ops {
			if op.method != nil {
				task, err = op.method(ctx)
				if err != nil {
					t.Fatal(err)
				}

				err = task.Wait(ctx)
				if op.fail {
					if err == nil {
						t.Errorf("%d: expected error", i)
					}
				} else {
					if err != nil {
						t.Errorf("%d: %s", i, err)
					}
				}
			}

			if len(op.state) != 0 {
				state, err := vm.PowerState(ctx)
				if err != nil {
					t.Fatal(err)
				}

				if state != op.state {
					t.Errorf("state=%s", state)
				}
			}
		}
	}
}

type vmCreateModel struct {
	Service *Service
	Server  *Server

	Client     *govmomi.Client
	Finder     *find.Finder
	Datacenter *object.Datacenter
	Folders    *object.DatacenterFolders
	Cluster    *object.ClusterComputeResource
	Pool       *object.ResourcePool
	Host       *object.HostSystem
	Datastore  *object.Datastore
}

func vmCreateModelWithESX() (*vmCreateModel, error) {
	m := new(vmCreateModel)

	m.Service = New(NewServiceInstance(esx.ServiceContent, esx.RootFolder))
	m.Server = m.Service.NewServer()

	ctx := context.Background()

	c, err := govmomi.NewClient(ctx, m.Server.URL, true)
	if err != nil {
		return nil, err
	}

	m.Client = c

	m.Finder = find.NewFinder(c.Client, false)

	m.Datacenter, err = m.Finder.DefaultDatacenter(ctx)
	if err != nil {
		return nil, err
	}

	m.Finder.SetDatacenter(m.Datacenter)

	m.Folders, err = m.Datacenter.Folders(ctx)
	if err != nil {
		return nil, err
	}

	m.Host, err = m.Finder.DefaultHostSystem(ctx)
	if err != nil {
		return nil, err
	}

	m.Pool, err = m.Finder.DefaultResourcePool(ctx)
	if err != nil {
		return nil, err
	}

	_, err = m.Server.TempDatastore(ctx, m.Host)
	if err != nil {
		return nil, err
	}

	m.Datastore, err = m.Finder.DefaultDatastore(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func vmCreateModelWithVC() (*vmCreateModel, error) {
	m := new(vmCreateModel)

	m.Service = New(NewServiceInstance(vc.ServiceContent, vc.RootFolder))

	m.Server = m.Service.NewServer()

	ctx := context.Background()

	c, err := govmomi.NewClient(ctx, m.Server.URL, true)
	if err != nil {
		return nil, err
	}

	m.Client = c

	f := object.NewRootFolder(c.Client)

	m.Datacenter, err = f.CreateDatacenter(ctx, "dc1")
	if err != nil {
		return nil, err
	}

	m.Finder = find.NewFinder(c.Client, false)

	m.Datacenter, err = m.Finder.DefaultDatacenter(ctx)
	if err != nil {
		return nil, err
	}

	m.Finder.SetDatacenter(m.Datacenter)

	m.Folders, err = m.Datacenter.Folders(ctx)
	if err != nil {
		return nil, err
	}

	m.Cluster, err = m.Folders.HostFolder.CreateCluster(ctx, "cluster1", types.ClusterConfigSpecEx{})
	if err != nil {
		return nil, err
	}

	m.Pool, err = m.Finder.ResourcePool(ctx, "*/*")
	if err != nil {
		return nil, err
	}

	var hosts []*object.HostSystem

	for i := 0; i < 3; i++ {
		spec := types.HostConnectSpec{
			HostName: fmt.Sprintf("host-%d", i),
		}

		task, cerr := m.Cluster.AddHost(ctx, spec, true, nil, nil)
		if cerr != nil {
			return nil, cerr
		}

		info, cerr := task.WaitForResult(ctx, nil)
		if cerr != nil {
			return nil, cerr
		}

		hosts = append(hosts, object.NewHostSystem(c.Client, info.Result.(types.ManagedObjectReference)))
	}

	_, err = m.Server.TempDatastore(ctx, hosts...)
	if err != nil {
		return nil, err
	}

	m.Datastore, err = m.Finder.DefaultDatastore(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}
