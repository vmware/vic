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

package management

import (
	"context"
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

func (d *Dispatcher) createResourcePool(conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) (*object.ResourcePool, error) {
	defer trace.End(trace.Begin(""))

	d.vchPoolPath = path.Join(settings.ResourcePoolPath, conf.Name)

	rp, err := d.session.Finder.ResourcePool(d.ctx, d.vchPoolPath)
	if err != nil {
		_, ok := err.(*find.NotFoundError)
		if !ok {
			err = errors.Errorf("Failed to query compute resource (%q): %q", d.vchPoolPath, err)
			return nil, err
		}
	} else {
		conf.ComputeResources = append(conf.ComputeResources, rp.Reference())
		return rp, nil
	}

	log.Infof("Creating Resource Pool %q", conf.Name)
	// TODO: expose the limits and reservation here via options
	resSpec := types.ResourceConfigSpec{
		CpuAllocation:    &types.ResourceAllocationInfo{},
		MemoryAllocation: &types.ResourceAllocationInfo{},
	}
	setResources(resSpec.CpuAllocation.GetResourceAllocationInfo(), settings.VCHSize.CPU)
	setResources(resSpec.MemoryAllocation.GetResourceAllocationInfo(), settings.VCHSize.Memory)
	rp, err = d.session.Pool.Create(d.ctx, conf.Name, resSpec)
	if err != nil {
		log.Debugf("Failed to create resource pool %q: %s", d.vchPoolPath, err)
		return nil, err
	}

	conf.ComputeResources = append(conf.ComputeResources, rp.Reference())
	return rp, nil
}

func setResources(spec *types.ResourceAllocationInfo, resource types.ResourceAllocationInfo) {
	*spec = resource
	if spec.Limit == 0 {
		spec.Limit = -1
	}
	// FIXME: govmomi omitempty
	if spec.Reservation == 0 {
		spec.Reservation = 1
	}
	if spec.Shares == nil {
		spec.Shares = &types.SharesInfo{
			Level: types.SharesLevelNormal,
		}
	}
	if spec.ExpandableReservation == nil {
		spec.ExpandableReservation = types.NewBool(true)
	}
}

func (d *Dispatcher) destroyResourcePoolIfEmpty(conf *config.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin(""))

	log.Infof("Removing Resource Pool %q", conf.Name)

	if d.parentResourcepool == nil {
		log.Warnf("Did not find parent VCH resource pool")
		return nil
	}
	var vms []*vm.VirtualMachine
	var err error
	if vms, err = d.parentResourcepool.GetChildrenVMs(d.ctx, d.session); err != nil {
		err = errors.Errorf("Unable to get children vm of resource pool %q: %s", d.parentResourcepool.Name(), err)
		return err
	}
	if len(vms) != 0 {
		err = errors.Errorf("Resource pool is not empty: %q", d.parentResourcepool.Name())
		return err
	}
	if _, err := tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
		return d.parentResourcepool.Destroy(ctx)
	}); err != nil {
		return err
	}
	return nil
}

func (d *Dispatcher) findResourcePool(path string) (*object.ResourcePool, error) {
	defer trace.End(trace.Begin(path))
	rp, err := d.session.Finder.ResourcePool(d.ctx, path)
	if err != nil {
		_, ok := err.(*find.NotFoundError)
		if !ok {
			err = errors.Errorf("Failed to query resource pool %q: %s", path, err)
			return nil, err
		}
		return nil, nil
	}
	return rp, nil
}

func (d *Dispatcher) getPoolResourceSettings(pool *object.ResourcePool) (*config.Resources, error) {
	var p mo.ResourcePool
	ps := []string{"config.cpuAllocation", "config.memoryAllocation"}

	if err := pool.Properties(d.ctx, pool.Reference(), ps, &p); err != nil {
		return nil, err
	}
	res := &config.Resources{}
	cpu := p.Config.CpuAllocation.GetResourceAllocationInfo()
	if cpu != nil {
		res.CPU = *cpu
		// handle default value
		validate.HandleDefaultSettings(&res.CPU)
	}
	memory := p.Config.MemoryAllocation.GetResourceAllocationInfo()
	if memory != nil {
		res.Memory = *memory
		validate.HandleDefaultSettings(&res.Memory)
	}
	return res, nil
}

func updateResourcePoolConfig(ctx context.Context, pool *object.ResourcePool, name string, size *config.Resources) error {
	defer trace.End(trace.Begin(fmt.Sprintf("cpu %#v, memory: %#v", size.CPU, size.Memory)))
	resSpec := &types.ResourceConfigSpec{
		CpuAllocation:    &types.ResourceAllocationInfo{},
		MemoryAllocation: &types.ResourceAllocationInfo{},
	}
	setResources(resSpec.CpuAllocation.GetResourceAllocationInfo(), size.CPU)
	setResources(resSpec.MemoryAllocation.GetResourceAllocationInfo(), size.Memory)
	return pool.UpdateConfig(ctx, name, resSpec)
}
