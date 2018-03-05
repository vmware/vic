// Copyright 2018 VMware, Inc. All Rights Reserved.
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

package placement

import (
	"sort"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/performance"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

const (
	// TODO(jzt): move these values into a Configuration struct so that consumers can provide
	// different weights.
	memUnconsumedWeight = 0.7 // available memory (total - consumed)
	memInactiveWeight   = 0.3 // active memory on the host
)

type rankedHost struct {
	HostReference string
	*performance.HostMetricsInfo
	score float64
}

type rankedHosts []rankedHost

func (r rankedHosts) Len() int           { return len(r) }
func (r rankedHosts) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r rankedHosts) Less(i, j int) bool { return r[i].score > r[j].score }

// RankedHostPolicy uses data from a MetricsProvider to decide on which host to power-on a VM.
type RankedHostPolicy struct {
	source performance.MetricsProvider
}

// NewRankedHostPolicy returns a RandomHostPolicy instance using the supplied MetricsProvider.
func NewRankedHostPolicy(s performance.MetricsProvider) *RankedHostPolicy {
	return &RankedHostPolicy{source: s}
}

// CheckHost returns true if the host has adequate capacity to power on the VM, false otherwise.
func (r *RankedHostPolicy) CheckHost(op trace.Operation, vm *vm.VirtualMachine) bool {
	// TODO(jzt): return false until we have host checking logic decided
	return false
}

// RecommendHost recommends an ideal host on which to place a newly created VM.
func (r *RankedHostPolicy) RecommendHost(op trace.Operation, vm *vm.VirtualMachine) (*object.HostSystem, error) {
	return nil, nil
}

func (r *RankedHostPolicy) rankHosts(op trace.Operation, hm map[string]*performance.HostMetricsInfo) []rankedHost {
	ranking := []rankedHost{}
	for h, m := range hm {
		rh := rankedHost{
			HostReference:   h,
			HostMetricsInfo: m,
			score:           r.rankMemory(m) * (1 - m.CPU.UsagePercent),
		}
		ranking = append(ranking, rh)
	}
	sort.Sort(rankedHosts(ranking))
	return ranking
}

func (r *RankedHostPolicy) rankMemory(hm *performance.HostMetricsInfo) float64 {
	free := float64(hm.Memory.TotalKB-hm.Memory.ConsumedKB) / 1024.0
	return free * memUnconsumedWeight
}
