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
	"math/rand"
	"sort"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/performance"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

const (
	memUnconsumedWeight = 0.7 // available memory (total - consumed)
	memInactiveWeight   = 0.3 // active memory on the host
)

// RandomHostPolicy uses data from a MetricsProvider to decide on which host to powerOn a VM.
type RandomHostPolicy struct {
	source performance.MetricsProvider
}

// NewRandomHostPolicy returns a RandomHostPolicy instance using the supplied MetricsProvider.
func NewRandomHostPolicy(s performance.MetricsProvider) *RandomHostPolicy {
	return &RandomHostPolicy{source: s}
}

// CheckHost returns true if the host has adequate capacity to power on the VM, false otherwise.
func (e *RandomHostPolicy) CheckHost(op trace.Operation, vm *vm.VirtualMachine) bool {
	// TODO(jzt): return false until we have host checking logic decided
	return false
}

// RecommendHost recommends an ideal host on which to place a newly created VM.
func (e *RandomHostPolicy) RecommendHost(op trace.Operation, vm *vm.VirtualMachine) (*object.HostSystem, error) {
	// TODO(jzt): randomize placement initially to allow usage of this
	// interface for development towards other ROBO-related issues.
	r, err := vm.ResourcePool(op)
	if err != nil {
		return nil, err
	}

	rp := compute.NewResourcePool(op, vm.Session, r.Reference())

	cls, err := rp.GetCluster(op)
	if err != nil {
		return nil, err
	}

	hosts, err := cls.Hosts(op)
	if err != nil {
		return nil, err
	}

	return hosts[rand.Intn(len(hosts))], nil
}

func rankHosts(op trace.Operation, hm map[string]*performance.HostMetricsInfo) []rankedHost {
	ranking := []rankedHost{}
	for h, m := range hm {
		rh := rankedHost{
			HostReference:   h,
			HostMetricsInfo: m,
			score:           rankMemory(m) * (1 - m.CPU.UsagePercent),
		}
		ranking = append(ranking, rh)
	}
	sort.Sort(rankedHosts(ranking))
	return ranking
}

func rankMemory(hm *performance.HostMetricsInfo) float64 {
	free := float64(hm.Memory.TotalKB-hm.Memory.ConsumedKB) / 1024.0
	return free * memUnconsumedWeight
}
