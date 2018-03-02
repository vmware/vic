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

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/performance"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// RandomHostPolicy chooses a random host on which to power-on a VM.
type RandomHostPolicy struct{}

// NewRandomHostPolicy returns a RandomHostPolicy instance.
func NewRandomHostPolicy(s performance.MetricsProvider) *RandomHostPolicy {
	return &RandomHostPolicy{}
}

// CheckHost always returns false in a RandomHostPolicy.
func (p *RandomHostPolicy) CheckHost(op trace.Operation, vm *vm.VirtualMachine) bool {
	return false
}

// RecommendHost recommends a random host on which to place a newly created VM.
func (p *RandomHostPolicy) RecommendHost(op trace.Operation, vm *vm.VirtualMachine) (*object.HostSystem, error) {
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
