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
	log "github.com/Sirupsen/logrus"
	"fmt"
)

func GetVCHstats(ctx context.Context, moref ...types.ManagedObjectReference) (*mo.ResourcePool, error) {
	log.Infof("In GetVCHstats")

	if Config.ResourcePool == nil {
		return &mo.ResourcePool{}, fmt.Errorf("Config.ResourcePool is nil")
	}

	var p mo.ResourcePool

	r := Config.ResourcePool.Reference()
	if len(moref) > 0 {
		r = moref[0]
	}

	ps := []string{"config.cpuAllocation", "config.memoryAllocation", "runtime.cpu", "runtime.memory", "parent"}

	if err := Config.ResourcePool.Properties(ctx, r, ps, &p); err != nil {
		return &p, fmt.Errorf("VCH stats error: %s", err)
	}
	stats := []int64{p.Config.CpuAllocation.GetResourceAllocationInfo().Limit,
		p.Config.MemoryAllocation.GetResourceAllocationInfo().Limit,
		p.Runtime.Cpu.OverallUsage,
		p.Runtime.Memory.OverallUsage}

	log.Debugf("The VCH stats are: %+v", stats)

	// If any of the stats is -1 (s is true), we need to get the vch stats from the parent resource pool
	s := false
	for _, v := range stats {
		if v == -1 {
			s = true
			break
		}
	}
	if s {
		return GetVCHstats(ctx, *p.Parent)
	}

	return &p, nil
}