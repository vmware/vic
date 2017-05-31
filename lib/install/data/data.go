// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"net/url"
	"time"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/ip"
)

// Data wraps all parameters required by value validation
type Data struct {
	*common.Target
	common.Debug
	common.Compute
	common.VCHID

	OpsUser     string `cmd:"ops-user"`
	OpsPassword *string

	CertPEM     []byte
	KeyPEM      []byte
	ClientCAs   []byte
	RegistryCAs []byte
	common.Images

	ImageDatastorePath     string              `cmd:"image-store"`
	VolumeLocations        map[string]*url.URL `cmd:"volume-store" label:"value-key"`
	ContainerDatastoreName string

	BridgeNetworkName string        `cmd:"bridge-network"`
	ClientNetwork     NetworkConfig `cmd:"client-network"`
	PublicNetwork     NetworkConfig `cmd:"public-network"`
	ManagementNetwork NetworkConfig `cmd:"management-network"`
	DNS               []net.IP      `cmd:"dns-server"`

	ContainerNetworks `cmd:"container-network"`

	VCHCPULimitsMHz       int               `cmd:"cpu"`
	VCHCPUReservationsMHz int               `cmd:"cpu-reservation"`
	VCHCPUShares          *types.SharesInfo `cmd:"cpu-shares"`

	VCHMemoryLimitsMB       int               `cmd:"memory"`
	VCHMemoryReservationsMB int               `cmd:"memory-reservation"`
	VCHMemoryShares         *types.SharesInfo `cmd:"memory-shares"`

	BridgeIPRange *net.IPNet `cmd:"bridge-network-range"`

	InsecureRegistries  []url.URL `cmd:"insecure-registry"`
	WhitelistRegistries []url.URL `cmd:"whitelist-registry"`

	HTTPSProxy *url.URL `cmd:"https-proxy"`
	HTTPProxy  *url.URL `cmd:"http-proxy"`

	NumCPUs  int `cmd:"endpoint-cpu"`
	MemoryMB int `cmd:"endpoint-memory"`

	Timeout time.Duration

	Force               bool
	UseRP               bool `cmd:"use-rp"`
	ResetInProgressFlag bool

	AsymmetricRouting bool `cmd:"asymmetric-routes"`

	ScratchSize string `cmd:"base-image-size"`
	Rollback    bool

	SyslogConfig SyslogConfig `cmd:"syslog"`
}

type SyslogConfig struct {
	Addr *url.URL `cmd:"address"`
	Tag  string
}

type ContainerNetworks struct {
	MappedNetworks         map[string]string     `cmd:"parent" label:"key-value"`
	MappedNetworksGateways map[string]net.IPNet  `cmd:"gateway" label:"key-value"`
	MappedNetworksIPRanges map[string][]ip.Range `cmd:"ip-range" label:"key-value"`
	MappedNetworksDNS      map[string][]net.IP   `cmd:"dns" label:"key-value"`
}

func (c *ContainerNetworks) IsSet() bool {
	return len(c.MappedNetworks) > 0 ||
		len(c.MappedNetworksGateways) > 0 ||
		len(c.MappedNetworksIPRanges) > 0 ||
		len(c.MappedNetworksDNS) > 0
}

// NetworkConfig is used to set IP addr for each network
type NetworkConfig struct {
	Name         string `cmd:"parent"`
	Destinations []net.IPNet
	Gateway      net.IPNet
	IP           net.IPNet `cmd:"ip"`
}

// Empty determines if ip and gateway are unset
func (n *NetworkConfig) Empty() bool {
	return ip.Empty(n.Gateway) && ip.Empty(n.IP)
}

func (n *NetworkConfig) IsSet() bool {
	return n.Name != "" || !n.Empty()
}

// InstallerData is used to hold the transient installation configuration that shouldn't be serialized
type InstallerData struct {
	// Virtual Container Host capacity
	VCHSize config.Resources
	// Appliance capacity
	ApplianceSize config.Resources

	KeyPEM  string
	CertPEM string

	//FIXME: remove following attributes after port-layer-server read config from guestinfo
	DatacenterName         string
	ClusterPath            string
	ResourcePoolPath       string
	ApplianceInventoryPath string

	Datacenter types.ManagedObjectReference
	Cluster    types.ManagedObjectReference

	ImageFiles map[string]string

	ApplianceISO      string
	BootstrapISO      string
	ISOVersion        string
	PreUpgradeVersion string
	Timeout           time.Duration

	UseRP bool

	HTTPSProxy *url.URL
	HTTPProxy  *url.URL
}

func NewData() *Data {
	d := &Data{
		Target: common.NewTarget(),
		ContainerNetworks: ContainerNetworks{
			MappedNetworks:         make(map[string]string),
			MappedNetworksGateways: make(map[string]net.IPNet),
			MappedNetworksIPRanges: make(map[string][]ip.Range),
			MappedNetworksDNS:      make(map[string][]net.IP),
		},
		Timeout: 3 * time.Minute,
	}
	return d
}

// CopyNonEmpty will shallow copy src value to override existing value if the value is set
// This copy will take care of relationship between variables, that means if any variable in ContainerNetwork or
// NetworkConfig is not empty, the whole ContainerNetwork or NetworkConfig will be copied
func (d *Data) CopyNonEmpty(src *Data) {
	// TODO: Add data copy here for each reconfigure items, to make sure specified variables present in the Data object.

	if src.ClientNetwork.IsSet() {
		d.ClientNetwork = src.ClientNetwork
	}
	if src.PublicNetwork.IsSet() {
		d.PublicNetwork = src.PublicNetwork
	}
	if src.ManagementNetwork.IsSet() {
		d.ManagementNetwork = src.ManagementNetwork
	}

	// copy container networks
	if src.ContainerNetworks.IsSet() {
		d.ContainerNetworks = src.ContainerNetworks
	}
	d.Timeout = src.Timeout
	if src.HTTPProxy != nil || src.HTTPSProxy != nil {
		d.HTTPProxy = src.HTTPProxy
		d.HTTPSProxy = src.HTTPSProxy
	}
}
