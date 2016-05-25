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
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/tasks"

	"golang.org/x/net/context"
)

func (d *Dispatcher) createResourcePool(conf *metadata.VirtualContainerHostConfigSpec) (*compute.ResourcePool, error) {
	d.vchPoolPath = path.Join(d.session.Pool.InventoryPath, conf.Name)
	conf.ResourcePoolPath = d.vchPoolPath

	rp, err := d.session.Finder.ResourcePool(d.ctx, d.vchPoolPath)
	if err != nil {
		_, ok := err.(*find.NotFoundError)
		if !ok {
			err = errors.Errorf("Failed to query compute resource (%s): %s", d.vchPoolPath, err)
			return nil, err
		}
	} else {
		conf.ComputeResources = append(conf.ComputeResources, rp)
		return compute.NewResourcePool(d.ctx, d.session, rp.Reference()), nil
	}

	log.Infof("Creating a Resource Pool")
	resSpec := types.ResourceConfigSpec{
		CpuAllocation: &types.ResourceAllocationInfo{
			Shares: &types.SharesInfo{
				Level: types.SharesLevelNormal,
			},
			ExpandableReservation: types.NewBool(true),
			Limit: -1,
			// FIXME: govmomi omitempty
			Reservation: 42,
		},
		MemoryAllocation: &types.ResourceAllocationInfo{
			Shares: &types.SharesInfo{
				Level: types.SharesLevelNormal,
			},
			ExpandableReservation: types.NewBool(true),
			Limit: -1,
			// FIXME: govmomi omitempty
			Reservation: 42,
		},
	}

	_, err = d.session.Pool.Create(d.ctx, conf.Name, resSpec)
	if err != nil {
		return nil, err
	}
	vrp, err := compute.FindResourcePool(d.ctx, d.session, d.vchPoolPath)
	if err != nil {
		return nil, err
	}
	conf.ComputeResources = append(conf.ComputeResources, rp)
	return vrp, nil
}

func (d *Dispatcher) destroyResourcePool(conf *metadata.VirtualContainerHostConfigSpec) error {
	log.Infof("Destroying the Resource Pool")

	vrp, err := compute.FindResourcePool(d.ctx, d.session, d.vchPoolPath)
	if err != nil {
		return err
	}

	_, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return vrp.Destroy(ctx)
	})
	if err != nil {
		return err
	}
	return nil
}
