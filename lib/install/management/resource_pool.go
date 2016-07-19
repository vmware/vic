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
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	"golang.org/x/net/context"
)

func (d *Dispatcher) createResourcePool(conf *metadata.VirtualContainerHostConfigSpec, settings *data.InstallerData) (*object.ResourcePool, error) {
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
		CpuAllocation: &types.ResourceAllocationInfo{
			Shares: &types.SharesInfo{
				Level: types.SharesLevelNormal,
			},
			ExpandableReservation: types.NewBool(true),
			Limit: -1,
			// FIXME: govmomi omitempty
			Reservation: 1,
		},
		MemoryAllocation: &types.ResourceAllocationInfo{
			Shares: &types.SharesInfo{
				Level: types.SharesLevelNormal,
			},
			ExpandableReservation: types.NewBool(true),
			Limit: -1,
			// FIXME: govmomi omitempty
			Reservation: 1,
		},
	}

	rp, err = d.session.Pool.Create(d.ctx, conf.Name, resSpec)
	if err != nil {
		log.Debugf("Failed to create resource pool %q: %s", d.vchPoolPath, err)
		return nil, err
	}

	conf.ComputeResources = append(conf.ComputeResources, rp.Reference())
	return rp, nil
}

func (d *Dispatcher) destroyResourcePoolIfEmpty(conf *metadata.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin(""))

	log.Infof("Removing Resource Pool %q", conf.Name)

	rpRef := conf.ComputeResources[len(conf.ComputeResources)-1]
	rp := compute.NewResourcePool(d.ctx, d.session, rpRef)

	var vms []*vm.VirtualMachine
	var err error
	if vms, err = rp.GetChildrenVMs(d.ctx, d.session); err != nil {
		err = errors.Errorf("Unable to get children vm of resource pool %q: %s", rp.Name(), err)
		return err
	}
	if len(vms) != 0 {
		err = errors.Errorf("Resource pool is not empty: %q", rp.Name())
		return err
	}
	if _, err := tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return rp.Destroy(ctx)
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
