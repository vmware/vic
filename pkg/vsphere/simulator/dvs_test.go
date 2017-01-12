// Copyright 2017 VMware, Inc. All Rights Reserved.
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
	"context"
	"testing"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/types"
)

func TestDVS(t *testing.T) {
	m := VPX()
	m.Portgroup = 0 // disabled DVS creation

	defer m.Remove()

	err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	c := m.Service.client

	finder := find.NewFinder(c, false)
	dc, _ := finder.DatacenterList(ctx, "*")
	finder.SetDatacenter(dc[0])
	folders, _ := dc[0].Folders(ctx)
	hosts, _ := finder.HostSystemList(ctx, "*/*")

	var spec types.DVSCreateSpec
	spec.ConfigSpec = &types.VMwareDVSConfigSpec{}
	spec.ConfigSpec.GetDVSConfigSpec().Name = "DVS0"

	dtask, err := folders.NetworkFolder.CreateDVS(ctx, spec)
	if err != nil {
		t.Fatal(err)
	}

	info, err := dtask.WaitForResult(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	dvs := object.NewDistributedVirtualSwitch(c, info.Result.(types.ManagedObjectReference))

	config := &types.DVSConfigSpec{}

	for _, host := range hosts {
		config.Host = append(config.Host, types.DistributedVirtualSwitchHostMemberConfigSpec{
			Operation: string(types.ConfigSpecOperationAdd),
			Host:      host.Reference(),
		})
	}

	// Add == OK
	dtask, err = dvs.Reconfigure(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	err = dtask.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Add == fail (AlreadyExists)
	dtask, err = dvs.Reconfigure(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	err = dtask.Wait(ctx)
	if _, ok := err.(task.Error).Fault().(*types.AlreadyExists); !ok {
		t.Fatalf("err=%v", err)
	}

	// Edit == fail (NotSupported)
	for i := range config.Host {
		config.Host[i].Operation = string(types.ConfigSpecOperationEdit)
	}

	dtask, err = dvs.Reconfigure(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	err = dtask.Wait(ctx)
	if _, ok := err.(task.Error).Fault().(*types.NotSupported); !ok {
		t.Fatalf("err=%v", err)
	}

	// Remove == OK
	for i := range config.Host {
		config.Host[i].Operation = string(types.ConfigSpecOperationRemove)
	}

	dtask, err = dvs.Reconfigure(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	err = dtask.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Add == OK
	for i := range config.Host {
		config.Host[i].Operation = string(types.ConfigSpecOperationAdd)
	}

	dtask, err = dvs.Reconfigure(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	err = dtask.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Add (PG) == OK
	dtask, err = dvs.AddPortgroup(ctx, []types.DVPortgroupConfigSpec{{Name: "DVPG0"}})
	if err != nil {
		t.Fatal(err)
	}

	err = dtask.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Remove == fail (ResourceInUse)
	for i := range config.Host {
		config.Host[i].Operation = string(types.ConfigSpecOperationRemove)
	}

	dtask, err = dvs.Reconfigure(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	err = dtask.Wait(ctx)
	if _, ok := err.(task.Error).Fault().(*types.ResourceInUse); !ok {
		t.Fatalf("err=%v", err)
	}

	// Remove == fail (ManagedObjectNotFound)
	for i := range config.Host {
		config.Host[i].Host.Value = "enoent"
	}

	dtask, err = dvs.Reconfigure(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	err = dtask.Wait(ctx)
	if _, ok := err.(task.Error).Fault().(*types.ManagedObjectNotFound); !ok {
		t.Fatalf("err=%v", err)
	}
}
