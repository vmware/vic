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

package network

import (
	"net"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/ip"
)

var Config Configuration

type Configuration struct {
	// The default bridge network supplied for the Virtual Container Host
	BridgeNetwork string `vic:"0.1" scope:"read-only" key:"bridge_network"`
	// Published networks available for containers to join, keyed by consumption name
	ContainerNetworks map[string]*ContainerNetwork `vic:"0.1" scope:"read-only" key:"container_networks"`
}

type ContainerNetwork struct {
	// Common.Name - the symbolic name for the network, e.g. web or backend
	// Common.ID - identifier of the underlay for the network
	metadata.Common

	// The network scope the IP belongs to.
	// The IP address is the default gateway
	Gateway net.IPNet `vic:"0.1" scope:"read-only" key:"gateway"`

	// The set of nameservers associated with this network - may be empty
	Nameservers []net.IP `vic:"0.1" scope:"read-only" key:"dns"`

	// The IP ranges for this network
	Pools []ip.Range `vic:"0.1" scope:"read-only" key:"pools"`

	PortGroup object.NetworkReference
}
