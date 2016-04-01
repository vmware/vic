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

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/vishvananda/netlink"
	"github.com/vmware/vic/metadata"
)

func linkByAddress(address string) (netlink.Link, error) {
	nis, err := net.Interfaces()
	if err != nil {
		detail := fmt.Sprintf("unable to iterate interfaces for LinkByAddress: %s", err)
		return nil, errors.New(detail)
	}

	for _, iface := range nis {
		if bytes.Equal([]byte(address), iface.HardwareAddr) {
			return netlink.LinkByName(iface.Name)
		}
	}

	return nil, fmt.Errorf("unable to locate interface for address %s", address)
}

// Apply takes the network endpoint configuration and applies it to the system
func Apply(endpoint *metadata.NetworkEndpoint) error {
	// Locate interface
	link, err := linkByAddress(endpoint.MAC)
	if err != nil {
		return err
	}

	// Take interface down
	err = netlink.LinkSetDown(link)
	if err != nil {
		detail := fmt.Sprintf("unable to take interface down for setup: %s", err)
		return errors.New(detail)
	}

	// Rename interface
	err = netlink.LinkSetName(link, endpoint.Network.Name)
	if err != nil {
		detail := fmt.Sprintf("unable to set interface name: %s", err)
		return errors.New(detail)
	}

	// Remove any existing addresses
	existingAddr, err := netlink.AddrList(link, syscall.AF_UNSPEC)
	if err != nil {
		detail := fmt.Sprintf("failed to list existing address for %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}

	for _, oldAddr := range existingAddr {
		err = netlink.AddrDel(link, &oldAddr)
		if err != nil {
			detail := fmt.Sprintf("failed to del existing address for %s: %s", endpoint.Network.Name, err)
			return errors.New(detail)
		}
	}

	// Set IP address
	addr, err := netlink.ParseAddr(endpoint.IP.String())
	if err != nil {
		detail := fmt.Sprintf("failed to parse address for %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}

	if err := netlink.AddrAdd(link, addr); err != nil {
		detail := fmt.Sprintf("failed to add address to %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}

	// Bring up interface
	if err = netlink.LinkSetUp(link); err != nil {
		detail := fmt.Sprintf("failed to bring up %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}

	// Add routes
	_, defaultNet, _ := net.ParseCIDR("0.0.0.0/0")
	route := netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet, Gw: endpoint.Network.Gateway.IP}
	err = netlink.RouteAdd(&route)
	if err != nil {
		detail := fmt.Sprintf("failed to add gateway route for endpoint %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}

	// Add /etc/hosts entry
	// TODO - figure out how to name us for each network
	hosts, err := os.OpenFile(pathPrefix+"/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		detail := fmt.Sprintf("failed to update hosts for endpoint %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}
	defer hosts.Close()

	_, err = hosts.WriteString(fmt.Sprintf("localhost.%s %s", endpoint.Network.Name, endpoint.IP.IP))
	if err != nil {
		detail := fmt.Sprintf("failed to add nameserver for endpoint %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}

	// Add nameservers
	// This is incredibly trivial for now - should be updated to a less messy approach
	resolv, err := os.OpenFile(pathPrefix+"/etc/resolv.conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		detail := fmt.Sprintf("failed to update resolv.confg for endpoint %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}
	defer resolv.Close()

	for _, server := range endpoint.Network.Nameservers {
		_, err = resolv.WriteString("nameserver " + server.String())
		if err != nil {
			detail := fmt.Sprintf("failed to add nameserver for endpoint %s: %s", endpoint.Network.Name, err)
			return errors.New(detail)
		}
	}

	return nil
}
