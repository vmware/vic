// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package common

import (
	"encoding"
	"fmt"
	"net"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/vic/pkg/ip"
)

// Container networks - mapped from vSphere
type CNetworks struct {
	ContainerNetworks         cli.StringSlice `arg:"container-network"`
	ContainerNetworksGateway  cli.StringSlice `arg:"container-network-gateway"`
	ContainerNetworksIPRanges cli.StringSlice `arg:"container-network-ip-range"`
	ContainerNetworksDNS      cli.StringSlice `arg:"container-network-dns"`
	IsSet                     bool
}

// CNetworkFields holds items returned by ProcessContainerNetworks for each
// container network.
type CNetworkFields struct {
	// NetAlias is the network alias for use by docker
	NetAlias string
	// VNet is the vSphere network name
	VNet     string
	Gateways net.IPNet
	IPRanges []ip.Range
	DNS      []net.IP
}

func (c *CNetworks) CNetworkFlags(hidden bool) []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:  "container-network, cn",
			Value: &c.ContainerNetworks,
			Usage: "vSphere network list that containers can use directly with labels, e.g. vsphere-net:backend. Defaults to DCHP - see advanced help (-x).",
		},
		cli.StringSliceFlag{
			Name:   "container-network-gateway, cng",
			Value:  &c.ContainerNetworksGateway,
			Usage:  "Gateway for the container network's subnet in CONTAINER-NETWORK:SUBNET format, e.g. vsphere-net:172.16.0.1/16",
			Hidden: hidden,
		},
		cli.StringSliceFlag{
			Name:   "container-network-ip-range, cnr",
			Value:  &c.ContainerNetworksIPRanges,
			Usage:  "IP range for the container network in CONTAINER-NETWORK:IP-RANGE format, e.g. vsphere-net:172.16.0.0/24, vsphere-net:172.16.0.10-172.16.0.20",
			Hidden: hidden,
		},
		cli.StringSliceFlag{
			Name:   "container-network-dns, cnd",
			Value:  &c.ContainerNetworksDNS,
			Usage:  "DNS servers for the container network in CONTAINER-NETWORK:DNS format, e.g. vsphere-net:8.8.8.8. Ignored if no static IP assigned.",
			Hidden: hidden,
		},
	}
}

func parseContainerNetworkGateways(cgs []string) (map[string]net.IPNet, error) {
	gws := make(map[string]net.IPNet)
	for _, cg := range cgs {
		m := &ipNetUnmarshaler{}
		vnet, err := parseVnetParam(cg, m)
		if err != nil {
			return nil, err
		}

		if _, ok := gws[vnet]; ok {
			return nil, fmt.Errorf("Duplicate gateway specified for container network %s", vnet)
		}

		gws[vnet] = net.IPNet{IP: m.ip, Mask: m.ipnet.Mask}
	}

	return gws, nil
}

func parseContainerNetworkIPRanges(cps []string) (map[string][]ip.Range, error) {
	pools := make(map[string][]ip.Range)
	for _, cp := range cps {
		ipr := &ip.Range{}
		vnet, err := parseVnetParam(cp, ipr)
		if err != nil {
			return nil, err
		}

		pools[vnet] = append(pools[vnet], *ipr)
	}

	return pools, nil
}

func parseContainerNetworkDNS(cds []string) (map[string][]net.IP, error) {
	dns := make(map[string][]net.IP)
	for _, cd := range cds {
		var ip net.IP
		vnet, err := parseVnetParam(cd, &ip)
		if err != nil {
			return nil, err
		}

		if ip == nil {
			return nil, fmt.Errorf("DNS IP not specified for container network %s", vnet)
		}

		dns[vnet] = append(dns[vnet], ip)
	}

	return dns, nil
}

func splitVnetParam(p string) (vnet string, value string, err error) {
	mapped := strings.Split(p, ":")
	if len(mapped) == 0 || len(mapped) > 2 {
		err = fmt.Errorf("Invalid value for parameter %s", p)
		return
	}

	vnet = mapped[0]
	if vnet == "" {
		err = fmt.Errorf("Container network not specified in parameter %s", p)
		return
	}

	// If the supplied vSphere network contains spaces then the user must supply a network alias. Guest info won't receive a name with spaces.
	if strings.Contains(vnet, " ") && (len(mapped) == 1 || (len(mapped) == 2 && len(mapped[1]) == 0)) {
		err = fmt.Errorf("A network alias must be supplied when network name %q contains spaces.", p)
		return
	}

	if len(mapped) > 1 {
		// Make sure the alias does not contain spaces
		if strings.Contains(mapped[1], " ") {
			err = fmt.Errorf("The network alias supplied in %q cannot contain spaces.", p)
			return
		}
		value = mapped[1]
	}

	return
}

func parseVnetParam(p string, m encoding.TextUnmarshaler) (vnet string, err error) {
	vnet, v, err := splitVnetParam(p)
	if err != nil {
		return "", fmt.Errorf("Error parsing container network parameter %s: %s", p, err)
	}

	if err = m.UnmarshalText([]byte(v)); err != nil {
		return "", fmt.Errorf("Error parsing container network parameter %s: %s", p, err)
	}

	return vnet, nil
}

type ipNetUnmarshaler struct {
	ipnet *net.IPNet
	ip    net.IP
}

func (m *ipNetUnmarshaler) UnmarshalText(text []byte) error {
	s := string(text)
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}

	m.ipnet = ipnet
	m.ip = ip
	return nil
}

// ProcessContainerNetworks parses container network settings and returns a
// slice of container network fields on success.
func (c *CNetworks) ProcessContainerNetworks() ([]CNetworkFields, error) {
	var cNets []CNetworkFields

	if c.ContainerNetworks != nil || c.ContainerNetworksGateway != nil ||
		c.ContainerNetworksIPRanges != nil || c.ContainerNetworksDNS != nil {
		c.IsSet = true
	}

	gws, err := parseContainerNetworkGateways([]string(c.ContainerNetworksGateway))
	if err != nil {
		return nil, cli.NewExitError(err.Error(), 1)
	}

	pools, err := parseContainerNetworkIPRanges([]string(c.ContainerNetworksIPRanges))
	if err != nil {
		return nil, cli.NewExitError(err.Error(), 1)
	}

	dns, err := parseContainerNetworkDNS([]string(c.ContainerNetworksDNS))
	if err != nil {
		return nil, cli.NewExitError(err.Error(), 1)
	}

	// parse container networks
	for _, cn := range c.ContainerNetworks {
		vnet, v, err := splitVnetParam(cn)
		if err != nil {
			return nil, cli.NewExitError(err.Error(), 1)
		}

		alias := vnet
		if v != "" {
			alias = v
		}

		cNet := CNetworkFields{
			NetAlias: alias,
			VNet:     vnet,
			Gateways: gws[vnet],
			IPRanges: pools[vnet],
			DNS:      dns[vnet],
		}
		cNets = append(cNets, cNet)

		delete(gws, vnet)
		delete(pools, vnet)
		delete(dns, vnet)
	}

	var hasError bool
	fmtMsg := "The following container network %s is set, but CONTAINER-NETWORK cannot be found. Please check the --container-network and %s settings"
	if len(gws) > 0 {
		log.Error(fmt.Sprintf(fmtMsg, "gateway", "--container-network-gateway"))
		for key, value := range gws {
			mask, _ := value.Mask.Size()
			log.Errorf("\t%s:%s/%d, %q should be vSphere network name", key, value.IP, mask, key)
		}
		hasError = true
	}
	if len(pools) > 0 {
		log.Error(fmt.Sprintf(fmtMsg, "ip range", "--container-network-ip-range"))
		for key, value := range pools {
			log.Errorf("\t%s:%s, %q should be vSphere network name", key, value, key)
		}
		hasError = true
	}
	if len(dns) > 0 {
		log.Errorf(fmt.Sprintf(fmtMsg, "dns", "--container-network-dns"))
		for key, value := range dns {
			log.Errorf("\t%s:%s, %q should be vSphere network name", key, value, key)
		}
		hasError = true
	}
	if hasError {
		return nil, cli.NewExitError("Inconsistent container network configuration.", 1)
	}

	return cNets, nil
}
