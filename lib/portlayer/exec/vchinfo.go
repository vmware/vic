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

package exec

import (
	"context"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// NCPU returns the CPU limit (MHZ)
func NCPU(ctx context.Context, moref ...types.ManagedObjectReference) int64 {
	if Config.ResourcePool == nil {
		return 0
	}

	var p mo.ResourcePool

	r := Config.ResourcePool.Reference()
	if len(moref) > 0 {
		r = moref[0]
	}

	if err := Config.ResourcePool.Properties(ctx, r, []string{"parent", "config.cpuAllocation"}, &p); err != nil {
		return 0
	}

	limit := p.Config.CpuAllocation.GetResourceAllocationInfo().Limit
	if limit == -1 {
		return NCPU(ctx, *p.Parent)
	}
	return limit
}

// MemTotal returns the memory limit (GiB)
func MemTotal(ctx context.Context, moref ...types.ManagedObjectReference) int64 {
	if Config.ResourcePool == nil {
		return 0
	}

	var p mo.ResourcePool

	r := Config.ResourcePool.Reference()
	if len(moref) > 0 {
		r = moref[0]
	}

	if err := Config.ResourcePool.Properties(ctx, r, []string{"parent", "config.memoryAllocation"}, &p); err != nil {
		return 0
	}

	limit := p.Config.MemoryAllocation.GetResourceAllocationInfo().Limit
	if limit == -1 {
		return MemTotal(ctx, *p.Parent)
	}

	return limit
}
