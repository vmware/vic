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

package metadata

import (
	"net"

	"github.com/vmware/govmomi/object"
)

// NetworkEndpoint describes a network presence in the form a vNIC in sufficient detail that it can be:
// a. created - the vNIC added to a VM
// b. identified - the guestOS can determine which interface it corresponds to
// c. configured - the guestOS can configure the interface correctly
type NetworkEndpoint struct {
	// IP addresses to assign - may be empty if DHCP
	IP net.IP
	// The MAC address for the vNIC - this allows for interface idenitifcaton in the guest
	MAC string
	// The network in which this information should be interpreted. This is embedded directly rather than
	// as a pointer so that we can ensure the data is consistent
	Network ContainerNetwork
}

// ContainerNetwork is the data needed on a per container basis both for vSphere to ensure it's attached
// to the correct network, and in the guest to ensure the interface is correctly configured.
type ContainerNetwork struct {
	// The symbolic name of the network
	Name string

	// The network scope the IP belongs to
	Network net.IPNet
	// Default gateway if any
	Gateway net.IP
	// The set of nameservers associated with this network - may be empty
	Nameservers []net.IP
}

// NetworkMapping records which vSphere networks are mapped to a given symbolic
// network at a consumption level
type NetworkMapping map[string]object.Network
