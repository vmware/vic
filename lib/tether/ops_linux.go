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
	"github.com/vmware/vic/lib/dhcp"
	"github.com/vmware/vic/lib/etcconf"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vmw-guestinfo/rpcout"
)

var hostnameFile = "/etc/hostname"
var byLabelDir = "/dev/disk/by-label"

const pciDevPath = "/sys/bus/pci/devices"

type BaseOperations struct {
	dhcpClient   dhcp.Client
	dhcpLoops    []chan bool
	hosts        etcconf.Hosts
	resolvConf   etcconf.ResolvConf
	dynEndpoints map[string][]*metadata.NetworkEndpoint
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
	RouteDel(*netlink.Route) error
	// Not quite netlink, but tightly associated

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

func (t *BaseOperations) RouteDel(route *netlink.Route) error {
	return netlink.RouteDel(route)
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
	lo4 := net.IPv4(127, 0, 0, 1)
	for _, a := range aliases {
		t.hosts.SetHost(a, lo4)
	}
	if err = t.hosts.Save(); err != nil {
		return err
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
			detail := fmt.Sprintf("failed to set link %s down for rename: %s", link.Attrs().Name, err)
			return nil, errors.New(detail)
		}

		err = t.LinkSetName(link, endpoint.Name)
		if err != nil {
			return nil, err
		}

		err = t.LinkSetUp(link)
		if err != nil {
			detail := fmt.Sprintf("failed to bring link %s up after rename: %s", link.Attrs().Name, err)
			return nil, errors.New(detail)
		}

		// reacquire link with updated attributes
		link, err := t.LinkBySlot(slot)
		if err != nil {
			detail := fmt.Sprintf("unable to reacquire link %s after rename pass: %s", link.Attrs().Name, err)
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
			detail := fmt.Sprintf("unable to reacquire link %s after rename pass: %s", link.Attrs().Name, err)
			return nil, errors.New(detail)
		}

		return link, nil
	}

	log.Warnf("Unable to add additional alias on link %s for %s", link.Attrs().Name, endpoint.Name)
	return link, nil
}

func assignStaticIP(t Netlink, link netlink.Link, endpoint *metadata.NetworkEndpoint) error {
	if endpoint.IsDynamic() {
		return nil
	}

	if err := linkAddrAdd(endpoint.Static, t, link); err != nil {
		return err
	}

	updateEndpoint(endpoint.Static, nil, endpoint)
	return nil
}

func assignDynamicIP(t Netlink, link netlink.Link, dc dhcp.Client, endpoint *metadata.NetworkEndpoint) (*dhcp.Packet, error) {
	var ack *dhcp.Packet
	var newIP *net.IPNet
	var err error

	addAddr := true
	timeout := time.After(30 * time.Second)

	for newIP == nil {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timed out")

		default:
			if dc != nil {
				// use dhcp to acquire address
				ack, err = dc.Request(link.Attrs().Index, link.Attrs().HardwareAddr)
				if err != nil {
					log.Errorf("error sending dhcp request: %s", err)
					return nil, err
				}

				if ack.YourIP() == nil || ack.SubnetMask() == nil {
					err = fmt.Errorf("dhcp assigned nil ip or subnet mask")
					log.Error(err)
					return nil, err
				}

				log.Infof("DHCP response: IP=%s, SubnetMask=%s, Gateway=%s, DNS=%s, Lease Time=%s", ack.YourIP(), ack.SubnetMask(), ack.Gateway(), ack.DNS(), ack.LeaseTime())

				newIP = &net.IPNet{IP: ack.YourIP(), Mask: ack.SubnetMask()}

				defer func() {
					if err != nil && ack != nil {
						dc.Release(ack)
					}
				}()
				break
			} else {
				// we do not have a dhcp client, just use the first ip
				// on the interface
				var addrs []netlink.Addr
				addrs, err = t.AddrList(link, netlink.FAMILY_V4)
				if err != nil {
					return nil, err
				}

				if len(addrs) > 0 {
					newIP = addrs[0].IPNet
					addAddr = false
					break
				}
			}

		}

		time.Sleep(1 * time.Second)
	}

	if addAddr {
		if err = linkAddrAdd(newIP, t, link); err != nil {
			return nil, err
		}
	}

	updateEndpoint(newIP, ack, endpoint)
	return ack, nil
}

func updateEndpoint(newIP *net.IPNet, ack *dhcp.Packet, endpoint *metadata.NetworkEndpoint) {
	if newIP == nil {
		return
	}

	endpoint.Assigned = newIP.IP
	if ack != nil {
		endpoint.Network.Gateway = net.IPNet{IP: ack.Gateway(), Mask: ack.SubnetMask()}
		dns := ack.DNS()
		if len(dns) > 0 {
			endpoint.Network.Nameservers = ack.DNS()
		}
	}
}

func linkAddrAdd(addr *net.IPNet, t Netlink, link netlink.Link) error {
	log.Infof("setting ip address %s for link %s", addr, link.Attrs().Name)

	var err error
	// assign IP to NIC
	if err = t.AddrAdd(link, &netlink.Addr{IPNet: addr}); err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EEXIST {
			log.Errorf("failed to assign dhcp ip %s for link %s", addr, link.Attrs().Name)
			return err
		}

		log.Warnf("address %s already set on interface %s", addr, link.Attrs().Name)
		err = nil
	}

	return err
}

func updateDefaultRoute(t Netlink, link netlink.Link, endpoint *metadata.NetworkEndpoint) error {
	// Add routes
	if !endpoint.Network.Default || ip.IsUnspecifiedIP(endpoint.Network.Gateway.IP) {
		log.Debugf("not setting route for network: default=%v gateway=%s", endpoint.Network.Default, endpoint.Network.Gateway.IP)
		return nil
	}

	_, defaultNet, _ := net.ParseCIDR("0.0.0.0/0")
	// delete default route first
	if err := t.RouteDel(&netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet}); err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.ESRCH {
			return fmt.Errorf("could not update default route: %s", err)
		}
	}

	route := &netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet, Gw: endpoint.Network.Gateway.IP}
	if err := t.RouteAdd(route); err != nil {
		detail := fmt.Sprintf("failed to add gateway route for endpoint %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}

	log.Infof("updated default route to %s interface: %s", endpoint.Network.Name, defaultNet.String())
	return nil
}

func (t *BaseOperations) updateHosts(endpoint *metadata.NetworkEndpoint) error {
	log.Debugf("%+v", endpoint)
	// Add /etc/hosts entry
	if endpoint.Network.Name == "" {
		return nil
	}

	t.hosts.SetHost(fmt.Sprintf("%s.localhost", endpoint.Network.Name), endpoint.Assigned)

	if err := t.hosts.Save(); err != nil {
		return err
	}

	return nil
}

func (t *BaseOperations) updateNameservers(endpoint *metadata.NetworkEndpoint) error {
	// Add nameservers
	// This is incredibly trivial for now - should be updated to a less messy approach
	if len(endpoint.Network.Nameservers) > 0 {
		t.resolvConf.AddNameservers(endpoint.Network.Nameservers...)
		log.Infof("Added nameservers: %+v", endpoint.Network.Nameservers)
	} else if !ip.IsUnspecifiedIP(endpoint.Network.Gateway.IP) {
		t.resolvConf.AddNameservers(endpoint.Network.Gateway.IP)
		log.Infof("Added nameserver: %s", endpoint.Network.Gateway.IP)
	}

	if err := t.resolvConf.Save(); err != nil {
		return err
	}

	return nil
}

func (t *BaseOperations) Apply(endpoint *metadata.NetworkEndpoint) error {
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

	var ack *dhcp.Packet
	defer func() {
		if err != nil && ack != nil {
			t.dhcpClient.Release(ack)
		}
	}()

	if ip.IsUnspecifiedIP(endpoint.Assigned) {
		// assign IP address as needed
		if endpoint.IsDynamic() {
			ack, err = assignDynamicIP(t, link, t.dhcpClient, endpoint)
		} else {
			err = assignStaticIP(t, link, endpoint)
		}

		if err != nil {
			detail := fmt.Sprintf("unable to assign IP for net %s: %s", endpoint.Network.Name, err)
			return errors.New(detail)
		}
	}

	if err = updateDefaultRoute(t, link, endpoint); err != nil {
		return err
	}

	if err = t.updateHosts(endpoint); err != nil {
		return err
	}

	if err = t.updateNameservers(endpoint); err != nil {
		return err
	}

	if endpoint.IsDynamic() {
		t.dynEndpoints[endpoint.ID] = append(t.dynEndpoints[endpoint.ID], endpoint)
	}

	// add renew/release loop if necessary
	if ack != nil {
		stop := make(chan bool)
		go t.dhcpLoop(stop, endpoint, ack)
		t.dhcpLoops = append(t.dhcpLoops, stop)
	}

	return nil
}

func (t *BaseOperations) dhcpUpdate(ack *dhcp.Packet, e *metadata.NetworkEndpoint) {
	old := &metadata.NetworkEndpoint{}
	old.Assigned = make(net.IP, len(e.Assigned))
	copy(old.Assigned, e.Assigned)
	old.Network.Nameservers = make([]net.IP, len(e.Network.Nameservers))
	copy(old.Network.Nameservers, e.Network.Nameservers)

	updateEndpoint(&net.IPNet{IP: ack.YourIP(), Mask: ack.SubnetMask()}, ack, e)
	id, err := strconv.Atoi(e.ID)
	if err != nil {
		log.Errorf("error updating endpoint: %s", err)
		return
	}

	link, err := t.LinkBySlot(int32(id))
	if err != nil {
		log.Errorf("error updating endpoint: %s", err)
		return
	}

	if err = updateDefaultRoute(t, link, e); err != nil {
		log.Errorf("error updating endpoint: %s", err)
	}

	if !old.Assigned.Equal(e.Assigned) {
		if err = t.updateHosts(e); err != nil {
			log.Errorf("could not add hosts for network %s: %s", e.Network.Name, err)
		}
	}

	updated := false
	for _, n := range old.Network.Nameservers {
		found := false
		for _, n2 := range e.Network.Nameservers {
			if n.Equal(n2) {
				found = true
				break
			}
		}

		if !found {
			updated = true
			break
		}
	}

	if updated {
		t.resolvConf.RemoveNameservers(old.Network.Nameservers...)
		if err = t.updateNameservers(e); err != nil {
			log.Errorf("could not add nameservers for network %s: %s", e.Network.Name, err)
		}
	}
}

func (t *BaseOperations) dhcpLoop(stop chan bool, e *metadata.NetworkEndpoint, ack *dhcp.Packet) {
	exp := time.After(ack.LeaseTime() / 2)
	for {
		select {
		case <-stop:
			// release the ip
			log.Infof("releasing IP address for network %s", e.Name)
			t.dhcpClient.Release(ack)
			return

		case <-exp:
			log.Infof("renewing IP address for network %s", e.Name)
			newack, err := t.dhcpClient.Renew(ack)
			if err != nil {
				log.Errorf("failed to renew ip address for network %s", e.Name)
				continue
			}

			log.Infof("successfully renewed ip address: %s", newack.YourIP())
			ack = newack

			t.dhcpUpdate(ack, e)

			// update any endpoints that share this NIC
			for _, d := range t.dynEndpoints[e.Network.ID] {
				t.dhcpUpdate(ack, d)
			}

			exp = time.After(ack.LeaseTime() / 2)
		}
	}
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

func (t *BaseOperations) Setup() error {
	c, err := dhcp.NewClient()
	if err != nil {
		return err
	}

	h := etcconf.NewHosts("")
	if err = h.Load(); err != nil {
		return err
	}

	// start with empty resolv.conf
	os.Remove("/etc/resolv.conf")

	rc := etcconf.NewResolvConf("")
	if err = rc.Load(); err != nil {
		return err
	}

	t.dynEndpoints = make(map[string][]*metadata.NetworkEndpoint)

	t.dhcpClient = c
	t.hosts = h
	t.resolvConf = rc
	return nil
}

func (t *BaseOperations) Cleanup() error {
	for _, stop := range t.dhcpLoops {
		stop <- true
	}

	return nil
}
