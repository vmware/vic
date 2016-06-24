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

package data

import (
	"net"
	"time"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/ip"
)

// Data wrapps all parameters required by value validation
type Data struct {
	*common.Target
	common.Debug

	Insecure bool

	CertPEM []byte
	KeyPEM  []byte

	ComputeResourcePath string
	ImageDatastoreName  string
	DisplayName         string

	ContainerDatastoreName string
	ExternalNetworkName    string
	ManagementNetworkName  string
	BridgeNetworkName      string
	ClientNetworkName      string

	MappedNetworks         map[string]string
	MappedNetworksGateways map[string]net.IPNet
	MappedNetworksIPRanges map[string][]ip.Range
	MappedNetworksDNS      map[string][]net.IP

	NumCPUs  int
	MemoryMB int

	Timeout time.Duration

	Force bool
}

// InstallerData is used to hold the transient installation configuration that shouldn't be serialized
type InstallerData struct {
	// Virtual Container Host capacity
	VCHSize metadata.Resources
	// Appliance capacity
	ApplianceSize metadata.Resources

	KeyPEM  string
	CertPEM string

	//FIXME: remove following attributes after port-layer-server read config from guestinfo
	DatacenterName         string
	ClusterPath            string
	ResourcePoolPath       string
	ApplianceInventoryPath string

	Datacenter types.ManagedObjectReference
	Cluster    types.ManagedObjectReference

	ImageFiles []string

	Extension types.Extension
}

func NewData() *Data {
	d := &Data{
		Target:                 common.NewTarget(),
		MappedNetworks:         make(map[string]string),
		MappedNetworksGateways: make(map[string]net.IPNet),
		MappedNetworksIPRanges: make(map[string][]ip.Range),
		MappedNetworksDNS:      make(map[string][]net.IP),
	}
	return d
}
