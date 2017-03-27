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
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
)

type ResourcePool struct {
	mo.ResourcePool
}

func NewResourcePool() *ResourcePool {
	return &ResourcePool{
		ResourcePool: esx.ResourcePool,
	}
}

func NewResourceConfigSpec() types.ResourceConfigSpec {
	spec := types.ResourceConfigSpec{
		CpuAllocation:    new(types.ResourceAllocationInfo),
		MemoryAllocation: new(types.ResourceAllocationInfo),
	}

	return spec
}

func (p *ResourcePool) setDefaultConfig(c types.BaseResourceAllocationInfo) {
	info := c.GetResourceAllocationInfo()

	if info.Shares == nil {
		info.Shares = new(types.SharesInfo)
	}

	if info.Shares.Level == "" {
		info.Shares.Level = types.SharesLevelNormal
	}

	if info.ExpandableReservation == nil {
		info.ExpandableReservation = types.NewBool(false)
	}
}

func (p *ResourcePool) CreateResourcePool(c *types.CreateResourcePool) soap.HasFault {
	body := &methods.CreateResourcePoolBody{}

	if e := Map.FindByName(c.Name, p.ResourcePool.ResourcePool); e != nil {
		body.Fault_ = Fault("", &types.DuplicateName{
			Name:   e.Entity().Name,
			Object: e.Reference(),
		})

		return body
	}

	child := NewResourcePool()

	child.Name = c.Name
	child.Owner = p.Owner
	child.Summary.GetResourcePoolSummary().Name = c.Name
	child.Config.CpuAllocation = c.Spec.CpuAllocation
	child.Config.MemoryAllocation = c.Spec.MemoryAllocation
	child.Config.Entity = c.Spec.Entity

	p.setDefaultConfig(child.Config.CpuAllocation)
	p.setDefaultConfig(child.Config.MemoryAllocation)

	Map.PutEntity(p, Map.NewEntity(child))

	p.ResourcePool.ResourcePool = append(p.ResourcePool.ResourcePool, child.Reference())

	body.Res = &types.CreateResourcePoolResponse{
		Returnval: child.Reference(),
	}

	return body
}

type destroyPoolTask struct {
	*ResourcePool
}

func (c *destroyPoolTask) Run(task *Task) (types.AnyType, types.BaseMethodFault) {
	if c.Parent.Type != "ResourcePool" {
		// Can't destroy the root pool
		return nil, &types.InvalidArgument{}
	}

	p := Map.Get(*c.Parent).(*ResourcePool)

	rp := &p.ResourcePool
	// Remove child reference from rp
	rp.ResourcePool = RemoveReference(c.Reference(), rp.ResourcePool)

	// The grandchildren become children of the parent (rp)
	//..........................................hello........hello........hello..........
	rp.ResourcePool = append(rp.ResourcePool, c.ResourcePool.ResourcePool.ResourcePool...)

	// And VMs move to the parent
	vms := c.ResourcePool.ResourcePool.Vm
	for _, vm := range vms {
		Map.Get(vm).(*VirtualMachine).ResourcePool = &rp.Self
	}

	rp.Vm = append(rp.Vm, vms...)

	Map.Remove(c.Reference())

	return nil, nil
}

func (p *ResourcePool) DestroyTask(c *types.Destroy_Task) soap.HasFault {
	r := &methods.Destroy_TaskBody{}

	task := NewTask(&destroyPoolTask{p})

	r.Res = &types.Destroy_TaskResponse{
		Returnval: task.Self,
	}

	task.Run()

	return r
}
