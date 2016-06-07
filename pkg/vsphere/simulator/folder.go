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
	"math/rand"
	"sync"

	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type Folder struct {
	mo.Folder

	m sync.Mutex
}

func (f *Folder) putChild(o mo.Entity) {
	Map.PutEntity(f, o)

	f.m.Lock()
	defer f.m.Unlock()

	f.ChildEntity = append(f.ChildEntity, o.Reference())
}

func (f *Folder) hasChildType(kind string) bool {
	for _, t := range f.ChildType {
		if t == kind {
			return true
		}
	}
	return false
}

func (f *Folder) typeNotSupported() *soap.Fault {
	return Fault(fmt.Sprintf("%s supports types: %#v", f.Self, f.ChildType), &types.NotSupported{})
}

type addStandaloneHostTask struct {
	*Folder

	req *types.AddStandaloneHost_Task
}

func (add *addStandaloneHostTask) Run(task *Task) (types.AnyType, types.BaseMethodFault) {
	host, err := CreateStandaloneHost(add.Folder, add.req.Spec)
	if err != nil {
		return nil, err
	}

	if add.req.AddConnected {
		host.Runtime.ConnectionState = types.HostSystemConnectionStateConnected
	}

	return host.Reference(), nil
}

func (f *Folder) AddStandaloneHostTask(a *types.AddStandaloneHost_Task) soap.HasFault {
	r := &methods.AddStandaloneHost_TaskBody{}

	if f.hasChildType("ComputeResource") && f.hasChildType("Folder") {
		task := NewTask(&addStandaloneHostTask{f, a})

		r.Res = &types.AddStandaloneHost_TaskResponse{
			Returnval: task.Self,
		}

		task.Run()
	} else {
		r.Fault_ = f.typeNotSupported()
	}

	return r
}

func (f *Folder) CreateFolder(c *types.CreateFolder) soap.HasFault {
	r := &methods.CreateFolderBody{}

	if f.hasChildType("Folder") {
		folder := &Folder{}

		folder.Name = c.Name
		folder.ChildType = f.ChildType
		folder.ChildEntity = f.ChildEntity

		f.putChild(folder)

		r.Res = &types.CreateFolderResponse{
			Returnval: folder.Self,
		}
	} else {
		r.Fault_ = f.typeNotSupported()
	}

	return r
}

func (f *Folder) CreateDatacenter(c *types.CreateDatacenter) soap.HasFault {
	r := &methods.CreateDatacenterBody{}

	if f.hasChildType("Datacenter") && f.hasChildType("Folder") {
		dc := &mo.Datacenter{}

		dc.Name = c.Name

		f.putChild(dc)

		createDatacenterFolders(dc, true)

		r.Res = &types.CreateDatacenterResponse{
			Returnval: dc.Self,
		}
	} else {
		r.Fault_ = f.typeNotSupported()
	}

	return r
}

func (f *Folder) CreateClusterEx(c *types.CreateClusterEx) soap.HasFault {
	r := &methods.CreateClusterExBody{}

	if f.hasChildType("ComputeResource") && f.hasChildType("Folder") {
		cluster, err := CreateClusterComputeResource(f, c.Name, c.Spec)
		if err != nil {
			r.Fault_ = Fault("", err)
			return r
		}

		r.Res = &types.CreateClusterExResponse{
			Returnval: cluster.Self,
		}
	} else {
		r.Fault_ = f.typeNotSupported()
	}

	return r
}

type createVMTask struct {
	*Folder

	req *types.CreateVM_Task
}

func (c *createVMTask) Run(task *Task) (types.AnyType, types.BaseMethodFault) {
	vm, err := NewVirtualMachine(&c.req.Config)
	if err != nil {
		return nil, err
	}

	vm.ResourcePool = &c.req.Pool
	if c.req.Host == nil {
		hostFolder := Map.getEntityDatacenter(c.Folder).HostFolder
		hosts := Map.Get(hostFolder).(*Folder).ChildEntity
		host := hosts[rand.Intn(len(hosts))]
		vm.Runtime.Host = &host
	} else {
		vm.Runtime.Host = c.req.Host
	}

	c.Folder.putChild(vm)

	return vm.Reference(), nil
}

func (f *Folder) CreateVMTask(c *types.CreateVM_Task) soap.HasFault {
	r := &methods.CreateVM_TaskBody{}

	task := NewTask(&createVMTask{f, c})

	r.Res = &types.CreateVM_TaskResponse{
		Returnval: task.Self,
	}

	task.Run()

	return r
}
