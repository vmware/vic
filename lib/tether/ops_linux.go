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

package tether

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vmw-guestinfo/rpcout"
)

var hostnameFile = "/etc/hostname"
var hostsFile = "/etc/hosts"
var resolvFile = "/etc/resolv.conf"
var byLabelDir = "/dev/disk/by-label"

const pciDevPath = "/sys/bus/pci/devices"

type BaseOperations struct {
}

// NetLink gives us an interface to the netlink calls used so that
// we can test the calling code.
type Netlink interface {
	LinkByName(string) (netlink.Link, error)
	LinkSetName(netlink.Link, string) error
	LinkSetDown(netlink.Link) error
	LinkSetUp(netlink.Link) error
	LinkSetAlias(netlink.Link, string) error
	AddrList(netlink.Link, int) ([]netlink.Addr, error)
	AddrAdd(netlink.Link, *netlink.Addr) error
	RouteAdd(*netlink.Route) error

	// Not quite netlink, but tightly assocaited
	LinkBySlot(slot int32) (netlink.Link, error)
}

func (t *BaseOperations) LinkByName(name string) (netlink.Link, error) {
	return netlink.LinkByName(name)
}

func (t *BaseOperations) LinkSetName(link netlink.Link, name string) error {
	return netlink.LinkSetName(link, name)
}

func (t *BaseOperations) LinkSetDown(link netlink.Link) error {
	return netlink.LinkSetDown(link)
}

func (t *BaseOperations) LinkSetUp(link netlink.Link) error {
	return netlink.LinkSetUp(link)
}

func (t *BaseOperations) LinkSetAlias(link netlink.Link, alias string) error {
	return netlink.LinkSetAlias(link, alias)
}

func (t *BaseOperations) AddrList(link netlink.Link, family int) ([]netlink.Addr, error) {
	return netlink.AddrList(link, family)
}

func (t *BaseOperations) AddrAdd(link netlink.Link, addr *netlink.Addr) error {
	return netlink.AddrAdd(link, addr)
}

func (t *BaseOperations) RouteAdd(route *netlink.Route) error {
	return netlink.RouteAdd(route)
}

func (t *BaseOperations) LinkBySlot(slot int32) (netlink.Link, error) {
	pciPath, err := slotToPCIPath(slot)
	if err != nil {
		return nil, err
	}

	name, err := pciToLinkName(pciPath)
	if err != nil {
		return nil, err
	}

	log.Debugf("got link name: %#v", name)
	return t.LinkByName(name)
}

// SetHostname sets both the kernel hostname and /etc/hostname to the specified string
func (t *BaseOperations) SetHostname(hostname string, aliases ...string) error {
	defer trace.End(trace.Begin("setting hostname to " + hostname))

	old, err := os.Hostname()
	if err != nil {
		log.Warnf("Unable to get current hostname - will not be able to revert on failure: %s", err)
	}

	err = syscall.Sethostname([]byte(hostname))
	if err != nil {
		log.Errorf("Unable to set hostname: %s", err)
		return err
	}
	log.Debugf("Updated kernel hostname")

	// update /etc/hostname to match
	err = ioutil.WriteFile(hostnameFile, []byte(hostname), 0644)
	if err != nil {
		log.Errorf("Failed to update hostname in %s", hostnameFile)

		// revert the hostname
		if old != "" {
			log.Warnf("Reverting kernel hostname to %s", old)
			err2 := syscall.Sethostname([]byte(old))
			if err2 != nil {
				log.Errorf("Unable to revert kernel hostname - kernel and hostname file are out of sync! Error: %s", err2)
			}
		}

		return err
	}

	// add entry to hosts for resolution without nameservers
	hosts, err := os.OpenFile(hostsFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		detail := fmt.Sprintf("failed to update hosts with hostname: %s", err)
		return errors.New(detail)
	}
	defer hosts.Close()

	names := strings.Join(aliases, " ")
	_, err = hosts.WriteString(fmt.Sprintf("\n127.0.0.1 %s %s\n", hostname, names))
	if err != nil {
		detail := fmt.Sprintf("failed to add hosts entry for hostname %s: %s", hostname, err)
		return errors.New(detail)
	}

	return nil
}

func slotToPCIPath(pciSlot int32) (string, error) {
	// see https://kb.vmware.com/kb/2047927
	dev := pciSlot & 0x1f
	bus := (pciSlot >> 5) & 0x1f
	fun := (pciSlot >> 10) & 0x7
	if bus == 0 {
		return path.Join(pciDevPath, fmt.Sprintf("0000:%02x:%02x.%d", bus, dev, fun)), nil
	}

	// device on secondary bus, prepend pci bridge address
	bridgeSlot := 0x11 + (bus - 1)
	bridgeAddr, err := slotToPCIPath(bridgeSlot)
	if err != nil {
		return "", err
	}

	return path.Join(bridgeAddr, fmt.Sprintf("0000:*:%02x.%d", dev, fun)), nil
}

func pciToLinkName(pciPath string) (string, error) {
	p := filepath.Join(pciPath, "net", "*")
	matches, err := filepath.Glob(p)
	if err != nil {
		return "", err
	}

	if len(matches) != 1 {
		return "", fmt.Errorf("more than one eth interface matches %s", p)
	}

	return path.Base(matches[0]), nil
}

func renameLink(t Netlink, link netlink.Link, slot int32, endpoint *metadata.NetworkEndpoint) (netlink.Link, error) {
	if link.Attrs().Name == endpoint.Name || link.Attrs().Alias == endpoint.Name || endpoint.Name == "" {
		// if the network is already identified, whether as primary name or alias it doesn't need repeating.
		// if the name is empty then it should not be aliases or named directly. IPAM data should still be applied.
		return link, nil
	}

	if strings.HasPrefix(link.Attrs().Name, "eno") {
		log.Infof("Renaming link %s to %s", link.Attrs().Name, endpoint.Name)

		err := t.LinkSetDown(link)
		if err != nil {
			detail := fmt.Sprintf("failed to set link %s down for rename: %s", endpoint.Name, err)
			return nil, errors.New(detail)
		}

		err = t.LinkSetName(link, endpoint.Name)
		if err != nil {
			return nil, err
		}

		err = t.LinkSetUp(link)
		if err != nil {
			detail := fmt.Sprintf("failed to bring link %s up after rename: %s", endpoint.Name, err)
			return nil, errors.New(detail)
		}

		// reacquire link with updated attributes
		link, err := t.LinkBySlot(slot)
		if err != nil {
			detail := fmt.Sprintf("unable to reacquire link %s after rename pass: %s", endpoint.ID, err)
			return nil, errors.New(detail)
		}

		return link, nil
	}

	if link.Attrs().Alias == "" {
		log.Infof("Aliasing link %s to %s", link.Attrs().Name, endpoint.Name)
		err := t.LinkSetAlias(link, endpoint.Name)
		if err != nil {
			return nil, err
		}

		// reacquire link with updated attributes
		link, err := t.LinkBySlot(slot)
		if err != nil {
			detail := fmt.Sprintf("unable to reacquire link %s after rename pass: %s", endpoint.ID, err)
			return nil, errors.New(detail)
		}

		return link, nil
	}

	log.Warnf("Unable to add additional alias on link %s for %s", link.Attrs().Name, endpoint.Name)
	return link, nil
}

// assignIP assigns an IP to a NIC, using a label to provide an associated between address and network role.
// returns true if an address has been updated so that /etc/hosts can be updated.
func assignIP(t Netlink, link netlink.Link, endpoint *metadata.NetworkEndpoint) (bool, error) {

	// get the current ip addresses on the link
	active, err := t.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		detail := fmt.Sprintf("unable to confirm assigned IP address for net %s: %s", endpoint.Network.Name, err)
		return false, errors.New(detail)
	}

	// Set IP address if it's specified - this is now named for for the network role rather than the nic
	if !ip.IsUnspecifiedIP(endpoint.Static.IP) {
		addr, err := netlink.ParseAddr(endpoint.Static.String())
		if err != nil {
			detail := fmt.Sprintf("failed to parse address for %s ednpoint: %s", endpoint.Network.Name, err)
			return false, errors.New(detail)
		}

		// see if there's a need to set an address
		for _, ipaddr := range active {
			log.Debugf("checking existing IPs on link: %s vs %s", ipaddr.IP.String(), addr.IP.String())
			// couldn't get bytes.Equal to match this - trailing data in the array maybe?
			if ipaddr.IP.String() == addr.IP.String() {
				log.Infof("address is already assigned to link, skipping assignment")
				// ensure the assigned field is set no matter what
				endpoint.Assigned = addr.IP

				return false, nil
			}
		}

		// add a label to identify the network
		addr.Label = fmt.Sprintf("%s:%s", link.Attrs().Name, endpoint.Network.Name)
		if err = t.AddrAdd(link, addr); err != nil {
			detail := fmt.Sprintf("failed to add address to %s: %s", endpoint.Network.Name, err)
			return false, errors.New(detail)
		}

		// report it
		endpoint.Assigned = endpoint.Static.IP
		log.Infof("Added IP address for %s: %s", endpoint.Network.Name, endpoint.Assigned.String())

		return true, nil
	}

	// TODO: split the entire network management out into an extension.
	// move extensions to Pre/Post so we can ensure setup prior to session launch
	// this will allow us to release leases on shutdown or disconnect

	// if there's already an address assigned, obtain it otherwise wait for one
	for {
		// update the current ip addresses on the link
		active, err = t.AddrList(link, netlink.FAMILY_V4)
		if err != nil {
			detail := fmt.Sprintf("unable to confirm assigned IP address for net %s: %s", endpoint.Network.Name, err)
			return false, errors.New(detail)
		}

		if !ip.IsUnspecifiedIP(endpoint.Network.Gateway.IP) {
			// if gateway is supplied filter with it
			for _, ipaddr := range active {
				if ipaddr.IP == nil {
					continue
				}

				log.Debugf("filtering ip %s with gateway %s", ipaddr.IP, endpoint.Network.Gateway.IP)
				if endpoint.Network.Gateway.Contains(ipaddr.IP) {
					// couldn't get bytes.Equal to match even when the string addresses do
					updated := endpoint.Assigned.String() != ipaddr.IP.String()
					if updated {
						endpoint.Assigned = ipaddr.IP
					}
					log.Infof("Using dynamic IP for network %s: %s", endpoint.Network.Name, endpoint.Assigned.String())
					return updated, nil
				}

				log.Debugf("rejecting IP %s on link %s due to mismatch with endpoint gateway", ipaddr.IP.String(), endpoint.Network.Name)
			}
		} else if len(active) > 0 {
			// if no gateway is specified then just take the first non-nil address
			for _, ipaddr := range active {
				if ipaddr.IP == nil {
					continue
				}

				updated := endpoint.Assigned.String() != ipaddr.IP.String()
				if updated {
					endpoint.Assigned = ipaddr.IP
				}

				log.Infof("Using dynamic IP (unvetted) for network %s: %s", endpoint.Network.Name, endpoint.Assigned.String())
				return updated, nil
			}
		}

		// we don't want to busy wait but I don't currently know how to wait for interface updates
		time.Sleep(100 * time.Millisecond)
	}
}

// Apply takes the network endpoint configuration and applies it to the system
func apply(t Netlink, endpoint *metadata.NetworkEndpoint) error {
	defer trace.End(trace.Begin("applying endpoint configuration for " + endpoint.Network.Name))

	// Locate interface
	slot, err := strconv.Atoi(endpoint.ID)
	if err != nil {
		detail := fmt.Sprintf("endpoint ID must be a base10 numeric pci slot identifier: %s", err)
		return errors.New(detail)
	}
	link, err := t.LinkBySlot(int32(slot))
	if err != nil {
		detail := fmt.Sprintf("unable to acquire reference to link %s: %s", endpoint.ID, err)
		return errors.New(detail)
	}

	// TODO: add dhcp client code

	// rename the link if needed
	link, err = renameLink(t, link, int32(slot), endpoint)
	if err != nil {
		detail := fmt.Sprintf("unable to reacquire link %s after rename pass: %s", endpoint.ID, err)
		return errors.New(detail)
	}

	// assign IP address as needed
	updated, err := assignIP(t, link, endpoint)
	if err != nil {
		detail := fmt.Sprintf("unable to assign IP for net %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}

	// Add routes
	if endpoint.Network.Default && len(endpoint.Network.Gateway.IP) > 0 {
		_, defaultNet, _ := net.ParseCIDR("0.0.0.0/0")
		route := netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet, Gw: endpoint.Network.Gateway.IP}
		err = t.RouteAdd(&route)
		if err != nil {
			if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EEXIST {
				detail := fmt.Sprintf("failed to add gateway route for endpoint %s: %s", endpoint.Network.Name, err)
				return errors.New(detail)
			}
		}

		log.Infof("Added route to %s interface: %s", endpoint.Network.Name, defaultNet.String())
	}

	// if there's not been any updates made then we don't want to edit hosts and nameservers
	if !updated {
		return nil
	}

	// Add /etc/hosts entry
	if endpoint.Network.Name != "" {
		hosts, err := os.OpenFile(hostsFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			detail := fmt.Sprintf("failed to update hosts for endpoint %s: %s", endpoint.Network.Name, err)
			return errors.New(detail)
		}
		defer hosts.Close()

		entry := fmt.Sprintf("%s %s.localhost", endpoint.Assigned, endpoint.Network.Name)
		_, err = hosts.WriteString(fmt.Sprintf("\n%s\n", entry))
		if err != nil {
			detail := fmt.Sprintf("failed to add hosts entry for endpoint %s: %s", endpoint.Network.Name, err)
			return errors.New(detail)
		}

		log.Infof("Added hosts entry: %s", entry)
	}

	// Add nameservers
	// This is incredibly trivial for now - should be updated to a less messy approach
	if len(endpoint.Network.Nameservers) > 0 {
		resolv, err := os.OpenFile(resolvFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			detail := fmt.Sprintf("failed to update %s for endpoint %s: %s", resolvFile, endpoint.Network.Name, err)
			return errors.New(detail)
		}
		defer resolv.Close()

		for _, server := range endpoint.Network.Nameservers {
			_, err = resolv.WriteString(fmt.Sprintf("\nnameserver %s\n", server.String()))
			if err != nil {
				detail := fmt.Sprintf("failed to add nameserver for endpoint %s: %s", endpoint.Network.Name, err)
				return errors.New(detail)
			}
			log.Infof("Added nameserver: %s", server.String())
		}
	}

	return nil
}

// Apply takes the network endpoint configuration and applies it to the system
func (t *BaseOperations) Apply(endpoint *metadata.NetworkEndpoint) error {
	return apply(t, endpoint)
}

// MountLabel performs a mount with the source and target being absolute paths
func (t *BaseOperations) MountLabel(source, target string, ctx context.Context) error {
	defer trace.End(trace.Begin(fmt.Sprintf("Mounting %s on %s", source, target)))

	if err := os.MkdirAll(target, 0600); err != nil {
		return fmt.Errorf("unable to create mount point %s: %s", target, err)
	}

	// do..while ! timedout
	var timeout bool
	for timeout = false; !timeout; {
		_, err := os.Stat(source)
		if err == nil || !os.IsNotExist(err) {
			break
		}

		deadline, ok := ctx.Deadline()
		timeout = ok && time.Now().After(deadline)
	}

	if timeout {
		detail := fmt.Sprintf("timed out waiting for %s to appear", source)
		return errors.New(detail)
	}

	if err := syscall.Mount(source, target, "ext4", syscall.MS_NOATIME, ""); err != nil {
		detail := fmt.Sprintf("mounting %s on %s failed: %s", source, target, err)
		return errors.New(detail)
	}

	return nil
}

// ProcessEnv does OS specific checking and munging on the process environment prior to launch
func (t *BaseOperations) ProcessEnv(env []string) []string {
	// TODO: figure out how we're going to specify user and pass all the settings along
	// in the meantime, hardcode HOME to /root
	homeIndex := -1
	for i, tuple := range env {
		if strings.HasPrefix(tuple, "HOME=") {
			homeIndex = i
			break
		}
	}
	if homeIndex == -1 {
		return append(env, "HOME=/root")
	}

	return env
}

// Fork triggers vmfork and handles the necessary pre/post OS level operations
func (t *BaseOperations) Fork() error {
	// unload vmxnet3 module

	// fork
	out, ok, err := rpcout.SendOne("vmfork-begin -1 -1")
	if err != nil {
		detail := fmt.Sprintf("error while calling vmfork: err=%s, out=%s, ok=%t", err, out, ok)
		log.Error(detail)
		return errors.New(detail)
	}

	if !ok {
		detail := fmt.Sprintf("failed to vmfork: %s", out)
		log.Error(detail)
		return errors.New(detail)
	}

	log.Infof("vmfork call succeeded: %s", out)

	// update system time

	// rescan scsi bus

	// reload vmxnet3 module

	// ensure memory and cores are brought online if not using udev

	return nil
}
