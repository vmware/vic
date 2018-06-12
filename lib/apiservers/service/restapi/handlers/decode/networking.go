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

package decode

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/client"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/errors"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"
)

func ProcessNetworks(op trace.Operation, d *data.Data, vch *models.VCH, finder client.Finder) error {
	if vch.Network != nil {
		if vch.Network.Bridge != nil {
			path, err := FromManagedObject(op, finder, "Network", vch.Network.Bridge.PortGroup)
			if err != nil {
				return errors.NewError(http.StatusBadRequest, "error finding bridge network portgroup: %s", err)
			}
			if path == "" {
				return errors.NewError(http.StatusBadRequest, "bridge network portgroup must be specified (by name or id)")
			}

			d.BridgeNetworkName = path
			bridgeIPRange := FromCIDR(&vch.Network.Bridge.IPRange)
			_, d.BridgeIPRange, err = net.ParseCIDR(bridgeIPRange)

			if err != nil {
				return errors.WrapError(http.StatusBadRequest, err)
			}
		}

		if vch.Network.Client != nil {
			err := processNetwork(op, finder, vch.Network.Client, &d.ClientNetwork, "client")
			if err != nil {
				return errors.WrapError(http.StatusBadRequest, err)
			}
		}

		if vch.Network.Management != nil {
			err := processNetwork(op, finder, vch.Network.Management, &d.ManagementNetwork, "management")
			if err != nil {
				return errors.WrapError(http.StatusBadRequest, err)
			}
		}

		if vch.Network.Public != nil {
			err := processNetwork(op, finder, vch.Network.Public, &d.PublicNetwork, "public")
			if err != nil {
				return errors.WrapError(http.StatusBadRequest, err)
			}

			// Process DNS server to be applied to all public, management and client network
			d.DNS = fromIPAddresses(vch.Network.Public.Nameservers)

			if len(d.DNS) > 3 {
				op.Warn("Maximum of 3 DNS servers allowed. Additional servers specified will be ignored.")
			}
		}

		if vch.Network.Container != nil {
			d.ContainerNetworks = common.ContainerNetworks{
				MappedNetworks:          make(map[string]string),
				MappedNetworksGateways:  make(map[string]net.IPNet),
				MappedNetworksIPRanges:  make(map[string][]ip.Range),
				MappedNetworksDNS:       make(map[string][]net.IP),
				MappedNetworksFirewalls: make(map[string]executor.TrustLevel),
			}

			for _, cNetwork := range vch.Network.Container {
				err := processContainerNetwork(op, finder, cNetwork, &d.ContainerNetworks)
				if err != nil {
					return errors.WrapError(http.StatusBadRequest, err)
				}
			}
		}
	}
	return nil
}

func processNetwork(op trace.Operation, finder client.Finder, network *models.Network, networkConfig *data.NetworkConfig, netType string) error {
	name, err := FromManagedObject(op, finder, "Network", network.PortGroup)
	if err != nil {
		return fmt.Errorf("error finding %s network portgroup: %s", netType, err)
	}

	gateway := fromGateway(network.Gateway)
	ip := FromCIDR(&network.Static)

	err = processNetworkConfig(op, networkConfig, netType, name, ip, gateway)
	if err != nil {
		return err
	}

	return nil
}

func processContainerNetwork(op trace.Operation, finder client.Finder, cNetwork *models.ContainerNetwork, cNetworkConfig *common.ContainerNetworks) error {
	alias := cNetwork.Alias

	// TODO [AngieCris]; figure out what fields are required and what are not
	// portgroup
	path, err := FromManagedObject(op, finder, "Network", cNetwork.PortGroup)
	if err != nil {
		return fmt.Errorf("error finding container network %s portgroup: %s", alias, err)
	}
	if path == "" {
		return fmt.Errorf("container network %s portgroup must be specified (by name or id)", alias)
	}
	cNetworkConfig.MappedNetworks[alias] = path

	// firewall
	// TODO [AngieCris]: firewall is "" -> trust level "unspecified"
	trustLevel, err := executor.ParseTrustLevel(cNetwork.Firewall)
	if err != nil {
		return fmt.Errorf("error parsing trust level for container network %s: %s", alias, err)
	}
	cNetworkConfig.MappedNetworksFirewalls[alias] = trustLevel

	// gateway
	if cNetwork.Gateway != nil {
		ip, mask, err := processGateway(op, cNetwork.Gateway)
		if err != nil {
			return fmt.Errorf("error parsing container network %s: %s", alias, err)
		}
		cNetworkConfig.MappedNetworksGateways[alias] = net.IPNet{
			IP:   ip,
			Mask: mask.Mask,
		}
	}

	// ip ranges
	cNetworkConfig.MappedNetworksIPRanges[alias] = fromIPRanges(cNetwork.IPRanges)

	// nameservers
	cNetworkConfig.MappedNetworksDNS[alias] = fromIPAddresses(cNetwork.Nameservers)

	return nil
}

// TODO [AngieCris]: duplicate of process network logic in cmd/create
func processNetworkConfig(op trace.Operation, network *data.NetworkConfig, netType string, name string, staticIP string, gateway string) error {
	network.Name = name

	if gateway != "" && staticIP == "" {
		return fmt.Errorf("Gateway provided without static IP for %s network", netType)
	}

	var ipNet *net.IPNet

	if staticIP != "" {
		var ipAddr net.IP
		ipAddr, ipNet, err := net.ParseCIDR(staticIP)
		if err != nil {
			return fmt.Errorf("Failed to parse the provided %s network IP address %s: %s", netType, staticIP, err)
		}

		network.IP.IP = ipAddr
		network.IP.Mask = ipNet.Mask
	}

	var err error
	if gateway != "" {
		network.Destinations, network.Gateway, err = parseGateway(gateway)

		if err != nil {
			return fmt.Errorf("Invalid %s network gateway: %s", netType, err)
		}

		if !network.IP.Contains(network.Gateway.IP) {
			return fmt.Errorf("%s gateway with IP %s is not reachable from %s", netType, network.Gateway.IP, ipNet.String())
		}

		// TODO(vburenin): this seems ugly, and it actually is. The reason is that a gateway required to specify
		// a network mask for it, which is just not how network configuration should be done. Network mask has to
		// be provided separately or with the IP address. It is hard to change all dependencies to keep mask
		// with IP address, so it will be stored with network gateway as it was previously.
		network.Gateway.Mask = network.IP.Mask
	}

	return nil
}

// process gateway model to gateway IP and routing destinations
func processGateway(op trace.Operation, gateway *models.Gateway) (net.IP, *net.IPNet, error) {
	addr := net.ParseIP(string(gateway.Address))
	if addr == nil {
		return nil, nil, fmt.Errorf("error parsing gateway IP %s", gateway.Address)
	}

	// only one routing destination is needed for our usage
	if gateway.RoutingDestinations == nil || len(gateway.RoutingDestinations) != 1 {
		return nil, nil, fmt.Errorf("error parsing network mask: exactly one subnet must be specified")
	}

	_, mask, err := net.ParseCIDR(string(gateway.RoutingDestinations[0]))
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing network mask: %s", err)
	}

	return addr, mask, nil
}

// parse gateway string to gateway IP and routing destinations
// TODO [AngieCris]: complete duplicate of an util function in cmd/create (maybe no need to de-duplicate?)
func parseGateway(gw string) (cidrs []net.IPNet, gwIP net.IPNet, err error) {
	ss := strings.Split(gw, ":")
	if len(ss) > 2 {
		err = fmt.Errorf("gateway %s specified incorrectly", gw)
		return
	}

	gwStr := ss[0]
	cidrsStr := ""
	if len(ss) > 1 {
		gwStr = ss[1]
		cidrsStr = ss[0]
	}

	if gwIP.IP = net.ParseIP(gwStr); gwIP.IP == nil {
		err = fmt.Errorf("Provided gateway IP address is not valid: %s", gwStr)
	}

	if err != nil {
		return
	}

	if cidrsStr != "" {
		for _, c := range strings.Split(cidrsStr, ",") {
			var ipnet *net.IPNet
			_, ipnet, err = net.ParseCIDR(c)
			if err != nil {
				err = fmt.Errorf("invalid CIDR in gateway specification: %s", err)
				return
			}
			cidrs = append(cidrs, *ipnet)
		}
	}

	return
}

func FromCIDR(m *models.CIDR) string {
	if m == nil {
		return ""
	}

	return string(*m)
}

func fromCIDRs(m *[]models.CIDR) *[]string {
	s := make([]string, 0, len(*m))
	for _, d := range *m {
		s = append(s, FromCIDR(&d))
	}

	return &s
}

func fromIPAddress(m *models.IPAddress) net.IP {
	if m == nil {
		return nil
	}

	return net.ParseIP(string(*m))
}

func fromIPAddresses(ipAddresses []models.IPAddress) []net.IP {
	ips := make([]net.IP, len(ipAddresses))

	for i, n := range ipAddresses {
		ips[i] = fromIPAddress(&n)
	}

	return ips
}

func fromIPRanges(ipRanges []models.IPRange) []ip.Range {
	ranges := make([]ip.Range, len(ipRanges))
	for i, r := range ipRanges {
		parsedR := ip.ParseRange(string(r))
		ranges[i] = *parsedR
	}

	return ranges
}

func fromGateway(m *models.Gateway) string {
	if m == nil {
		return ""
	}

	if m.RoutingDestinations == nil {
		return fmt.Sprintf("%s",
			m.Address,
		)
	}

	return fmt.Sprintf("%s:%s",
		strings.Join(*fromCIDRs(&m.RoutingDestinations), ","),
		m.Address,
	)
}
