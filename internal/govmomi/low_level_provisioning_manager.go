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

package llpm

import (
	"context"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
)

// LowLevelProvisioningManager struct
type LowLevelProvisioningManager struct {
	object.Common
	c *vim25.Client
}

// NewLowLevelProvisioningManager returns a new LLPM
func NewLowLevelProvisioningManager(client *vim25.Client, ref types.ManagedObjectReference) *LowLevelProvisioningManager {
	return &LowLevelProvisioningManager{
		Common: object.NewCommon(client, ref),
		c:      client,
	}
}

// Create calls LowLevelCreateVm
func (llpm *LowLevelProvisioningManager) Create(ctx context.Context, configSpec *types.VirtualMachineConfigSpec) (*object.Task, error) {
	req := LowLevelCreateVm_Task{
		This:       llpm.Reference(),
		ConfigSpec: *configSpec,
	}

	res, err := LowLevelCreateVm(ctx, llpm.c, &req)
	if err != nil {
		return nil, err
	}

	return object.NewTask(llpm.c, res.Returnval), nil
}

// Delete calls LowLevelDeleteVm
func (llpm *LowLevelProvisioningManager) Delete(ctx context.Context, configInfo *types.VirtualMachineConfigInfo) (*object.Task, error) {
	req := DeleteVm_Task{
		This:       llpm.Reference(),
		ConfigInfo: *configInfo,
	}

	res, err := LowLevelDeleteVm(ctx, llpm.c, &req)
	if err != nil {
		return nil, err
	}

	return object.NewTask(llpm.c, res.Returnval), nil
}

// DeleteExceptDisks calls LowLevelDeleteVmExceptDisks
func (llpm *LowLevelProvisioningManager) DeleteExceptDisks(ctx context.Context, vmPathName string) (*object.Task, error) {
	req := DeleteVmExceptDisks_Task{
		This:       llpm.Reference(),
		VmPathName: vmPathName,
	}

	res, err := LowLevelDeleteVmExceptDisks(ctx, llpm.c, &req)
	if err != nil {
		return nil, err
	}

	return object.NewTask(llpm.c, res.Returnval), nil
}

// Reconfigure calls LowLevelReconfigVM
func (llpm *LowLevelProvisioningManager) Reconfigure(ctx context.Context, configSpec *types.VirtualMachineConfigSpec) (*object.Task, error) {
	req := LLPMReconfigVM_Task{
		This:       llpm.Reference(),
		ConfigSpec: *configSpec,
	}

	res, err := LowLevelReconfigVM(ctx, llpm.c, &req)
	if err != nil {
		return nil, err
	}

	return object.NewTask(llpm.c, res.Returnval), nil
}
