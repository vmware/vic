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
	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

func (d *Dispatcher) createVApp(conf *metadata.VirtualContainerHostConfigSpec, settings *data.InstallerData) (*object.VirtualApp, error) {
	defer trace.End(trace.Begin(""))
	var err error

	log.Infof("Creating virtual app %s", conf.Name)

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

	prodSpec := types.VAppProductSpec{
		Info: &types.VAppProductInfo{
			Name:      "vSphere Integrated Containers",
			Vendor:    "VMware",
			VendorUrl: "http://www.vmware.com/",
			Version:   "0.0.1",
		},
		ArrayUpdateSpec: types.ArrayUpdateSpec{
			Operation: types.ArrayUpdateOperationAdd,
		},
	}

	configSpec := types.VAppConfigSpec{
		Annotation: "vSphere Integrated Containers",
		VmConfigSpec: types.VmConfigSpec{
			Product: []types.VAppProductSpec{prodSpec},
		},
	}

	app, err := d.session.Pool.CreateVApp(d.ctx, conf.Name, resSpec, configSpec, d.session.Folders(d.ctx).VmFolder)
	if err != nil {
		log.Debugf("Failed to create virtual app %s: %s", conf.Name, err)
		return nil, err
	}
	conf.ComputeResources = append(conf.ComputeResources, app.Reference())
	return app, nil
}

func (d *Dispatcher) findVirtualApp(path string) (*object.VirtualApp, error) {
	defer trace.End(trace.Begin(path))
	vapp, err := d.session.Finder.VirtualApp(d.ctx, path)
	if err != nil {
		_, ok := err.(*find.NotFoundError)
		if !ok {
			err = errors.Errorf("Failed to query virtual app (%s): %s", path, err)
			return nil, err
		}
		return nil, nil
	}
	return vapp, nil
}
