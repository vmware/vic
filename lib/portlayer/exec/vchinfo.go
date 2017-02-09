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
)

func GetVCHstats(ctx context.Context, moref ...types.ManagedObjectReference) (*mo.ResourcePool, error) {

	p := &mo.ResourcePool{}

	if Config.ResourcePool == nil {
		log.Debugf("Config.ResourcePool is nil")
		return p, nil
	}

	r := Config.ResourcePool.Reference()
	if len(moref) > 0 {
		r = moref[0]
	}

	ps := []string{"parent", "runtime", "config"}

	if err := Config.ResourcePool.Properties(ctx, r, ps, p); err != nil {
		log.Errorf("Error while obtaining VCH stats: %s", err)
		return p, err
	}
	stats := []int64{p.Config.CpuAllocation.GetResourceAllocationInfo().Limit,
		p.Config.MemoryAllocation.GetResourceAllocationInfo().Limit,
		p.Runtime.Cpu.OverallUsage,
		p.Runtime.Memory.OverallUsage}

	// If any of the stats is -1 (s is true), we need to get the vch stats from the parent resource pool
	s := false
	for v := range stats {
		if v == -1 {
			s = true
			break
		}
	}

	if s {
		return GetVCHstats(ctx, *p.Parent)
	}

	return p, nil
}