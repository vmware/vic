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

package simulator

import (
	"strings"
	"time"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
)

type HostSystem struct {
	mo.HostSystem
}

func NewHostSystem(host mo.HostSystem) *HostSystem {
	now := time.Now()

	host.Name = host.Summary.Config.Name
	host.Summary.Runtime = &host.Runtime
	host.Summary.Runtime.BootTime = &now

	return &HostSystem{
		HostSystem: host,
	}
}

// CreateDefaultESX creates a standalone ESX
// Adds objects of type: Datacenter, Network, ComputeResource, ResourcePool and HostSystem
func CreateDefaultESX(f *Folder) {
	dc := &esx.Datacenter
	createDatacenterFolders(dc, false)
	f.putChild(dc)

	host := NewHostSystem(esx.HostSystem)

	for _, ref := range host.Network {
		network := &mo.Network{}
		network.Self = ref
		network.Name = strings.Split(ref.Value, "-")[1]
		Map.Get(dc.NetworkFolder).(*Folder).putChild(network)
	}

	cr := &mo.ComputeResource{}
	cr.Self = *host.Parent
	cr.Name = host.Name
	cr.Host = append(cr.Host, host.Reference())
	Map.PutEntity(cr, host)

	pool := esx.ResourcePool
	cr.ResourcePool = &pool.Self
	Map.PutEntity(cr, &pool)

	Map.Get(dc.HostFolder).(*Folder).putChild(cr)
}
