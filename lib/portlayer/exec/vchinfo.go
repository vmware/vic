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
	"fmt"
	"context"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type StatKey int

const (
	StatNCPU StatKey = iota
	StatMemTotal
	StatCPUUsage
	StatMemUsage
)

var statKey map[StatKey][]string = map[StatKey][]string{
	StatNCPU:     {"parent", "config.cpuAllocation"},
	StatMemTotal: {"parent", "config.memoryAllocation"},
	StatCPUUsage: {"parent", "runtime.cpu"},
	StatMemUsage: {"parent", "runtime.memory"},
}

func getVCHstats(statID StatKey, ctx context.Context, moref ...types.ManagedObjectReference) int64 {

	if Config.ResourcePool == nil {
		return 0
	}

	resPool, err := getProperties(ctx, statID, moref...)
	if err != nil {
		return 0
	}

	v := getLimit(statID, resPool)

	if v == -1 {
		return getVCHstats(statID, ctx, *resPool.Parent)
	}
	return v
}

func getProperties(ctx context.Context, statID StatKey, moref ...types.ManagedObjectReference) (*mo.ResourcePool, error) {
	r := Config.ResourcePool.Reference()
	if len(moref) > 0 {
		r = moref[0]
	}

	v, ok := statKey[statID]
	p := &mo.ResourcePool{}

	if !ok {
		panic(fmt.Sprintf("Unexpected stat requested: %q", statID))
	}
	err := Config.ResourcePool.Properties(ctx, r, v, p)
	return p, err
}

func getLimit(statID StatKey, p *mo.ResourcePool) int64 {
	switch statID {
	case StatNCPU:
		return p.Config.CpuAllocation.GetResourceAllocationInfo().Limit
	case StatMemTotal:
		return p.Config.MemoryAllocation.GetResourceAllocationInfo().Limit
	case StatCPUUsage:
		return p.Runtime.Cpu.OverallUsage
	case StatMemUsage:
		return p.Runtime.Memory.OverallUsage
	default:
		panic(fmt.Sprintf("Unexpected stat requested: %q", statID))
	}
}