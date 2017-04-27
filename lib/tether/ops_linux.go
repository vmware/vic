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

package tether

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/d2g/dhcp4"
	"github.com/docker/docker/pkg/archive"
	// need to use libcontainer for user validation, for os/user package cannot find user here if container image is busybox
	"github.com/opencontainers/runc/libcontainer/user"
	"github.com/vishvananda/netlink"

	"github.com/vmware/vic/lib/dhcp"
	"github.com/vmware/vic/lib/dhcp/client"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vmw-guestinfo/rpcout"
)

var (
	hostnameFile = "/etc/hostname"
	byLabelDir   = "/dev/disk/by-label"

	defaultExecUser = &user.ExecUser{
		Uid:  syscall.Getuid(),
		Gid:  syscall.Getgid(),
		Home: "/",
	}
)

const (
	pciDevPath         = "/sys/bus/pci/devices"
	nfsFileSystemType  = "nfs"
	ext4FileSystemType = "ext4"
	bridgeTableNumber  = 201
)

type BaseOperations struct {
	dhcpLoops    map[string]chan struct{}
	dynEndpoints map[string][]*NetworkEndpoint
	config       Config
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
	AddrDel(netlink.Link, *netlink.Addr) error
	RouteAdd(*netlink.Route) error
	RouteDel(*netlink.Route) error
	RuleList(family int) ([]netlink.Rule, error)
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

func (t *BaseOperations) AddrDel(link netlink.Link, addr *netlink.Addr) error {
	return netlink.AddrDel(link, addr)
}

func (t *BaseOperations) RouteAdd(route *netlink.Route) error {
	return netlink.RouteAdd(route)
}

func (t *BaseOperations) RouteDel(route *netlink.Route) error {
	return netlink.RouteDel(route)
}

func (t *BaseOperations) RuleList(family int) ([]netlink.Rule, error) {
	return netlink.RuleList(family)
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

	err = Sys.Syscall.Sethostname([]byte(hostname))
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
			err2 := Sys.Syscall.Sethostname([]byte(old))
			if err2 != nil {
				log.Errorf("Unable to revert kernel hostname - kernel and hostname file are out of sync! Error: %s", err2)
			}
		}

		return err
	}

	// add entry to hosts for resolution without nameservers
	lo4 := net.IPv4(127, 0, 1, 1)
	for _, a := range append(aliases, hostname) {
		Sys.Hosts.SetHost(a, lo4)
	}
	if err = Sys.Hosts.Save(); err != nil {
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
		return "", fmt.Errorf("%d eth interfaces match %s (%v)", len(matches), p, matches)
	}

	return path.Base(matches[0]), nil
}

func renameLink(t Netlink, link netlink.Link, slot int32, endpoint *NetworkEndpoint) (rlink netlink.Link, err error) {
	rlink = link
	defer func() {
		// we still need to ensure that the link is up irrespective of path
		err = t.LinkSetUp(link)
		if err != nil {
			err = fmt.Errorf("failed to bring link %s up: %s", link.Attrs().Name, err)
		}
	}()

	if link.Attrs().Name == endpoint.Name || link.Attrs().Alias == endpoint.Name || endpoint.Name == "" {
		// if the network is already identified, whether as primary name or alias it doesn't need repeating.
		// if the name is empty then it should not be aliases or named directly. IPAM data should still be applied.
		return link, nil
	}

	if strings.HasPrefix(link.Attrs().Name, "eth") {
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

func getDynamicIP(t Netlink, link netlink.Link, endpoint *NetworkEndpoint) (client.Client, error) {
	var ack *dhcp.Packet
	var err error

	// use dhcp to acquire address
	dc, err := client.NewClient(link.Attrs().Index, link.Attrs().HardwareAddr)
	if err != nil {
		return nil, err
	}

	params := []byte{byte(dhcp4.OptionSubnetMask)}
	if ip.IsUnspecifiedIP(endpoint.Network.Gateway.IP) {
		params = append(params, byte(dhcp4.OptionRouter))
	}
	if len(endpoint.Network.Nameservers) == 0 {
		params = append(params, byte(dhcp4.OptionDomainNameServer))
	}

	dc.SetParameterRequestList(params...)

	err = dc.Request()
	if err != nil {
		log.Errorf("error sending dhcp request: %s", err)
		return nil, err
	}

	ack = dc.LastAck()
	if ack.YourIP() == nil || ack.SubnetMask() == nil {
		err = fmt.Errorf("dhcp assigned nil ip or subnet mask")
		log.Error(err)
		return nil, err
	}

	log.Infof("DHCP response: IP=%s, SubnetMask=%s, Gateway=%s, DNS=%s, Lease Time=%s", ack.YourIP(), ack.SubnetMask(), ack.Gateway(), ack.DNS(), ack.LeaseTime())
	defer func() {
		if err != nil && ack != nil {
			dc.Release()
		}
	}()

	return dc, nil
}

func updateEndpoint(newIP *net.IPNet, endpoint *NetworkEndpoint) {
	log.Debugf("updateEndpoint(%s, %+v)", newIP, endpoint)

	dhcp := endpoint.DHCP
	if dhcp == nil {
		endpoint.Assigned = *newIP
		endpoint.Network.Assigned.Gateway = endpoint.Network.Gateway
		endpoint.Network.Assigned.Nameservers = endpoint.Network.Nameservers
		return
	}

	endpoint.Assigned = dhcp.Assigned
	endpoint.Network.Assigned.Gateway = dhcp.Gateway
	if len(dhcp.Nameservers) > 0 {
		endpoint.Network.Assigned.Nameservers = dhcp.Nameservers
	}
}

func linkAddrUpdate(old, new *net.IPNet, t Netlink, link netlink.Link) error {
	log.Infof("setting ip address %s for link %s", new, link.Attrs().Name)

	if old != nil && !old.IP.Equal(new.IP) {
		log.Debugf("removing old address %s", old)
		if err := t.AddrDel(link, &netlink.Addr{IPNet: old}); err != nil {
			if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EADDRNOTAVAIL {
				log.Errorf("failed to remove existing address %s: %s", old, err)
				return err
			}
		}

		log.Debugf("removed old address %s for link %s", old, link.Attrs().Name)
	}

	// assign IP to NIC
	if err := t.AddrAdd(link, &netlink.Addr{IPNet: new}); err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EEXIST {
			log.Errorf("failed to assign ip %s for link %s", new, link.Attrs().Name)
			return err

		}

		log.Warnf("address %s already set on interface %s", new, link.Attrs().Name)
	}

	log.Debugf("added address %s to link %s", new, link.Attrs().Name)
	return nil
}

func updateRoutes(t Netlink, link netlink.Link, endpoint *NetworkEndpoint) error {
	gw := endpoint.Network.Assigned.Gateway
	if ip.IsUnspecifiedIP(gw.IP) {
		return nil
	}

	if endpoint.Network.Default {
		return updateDefaultRoute(t, link, endpoint)
	}

	for _, d := range endpoint.Network.Destinations {
		r := &netlink.Route{
			LinkIndex: link.Attrs().Index,
			Dst:       &d,
			Gw:        gw.IP,
		}

		if err := t.RouteAdd(r); err != nil && !os.IsNotExist(err) {
			log.Errorf("failed to add route for destination %s via gateway %s", d, gw.IP)
		}
	}

	return nil
}

func bridgeTableExists(t Netlink) bool {
	rules, err := t.RuleList(syscall.AF_INET)
	if err != nil {
		return false
	}

	for _, r := range rules {
		if r.Table == bridgeTableNumber {
			return true
		}
	}

	return false
}

func updateDefaultRoute(t Netlink, link netlink.Link, endpoint *NetworkEndpoint) error {
	gw := endpoint.Network.Assigned.Gateway
	// Add routes
	if !endpoint.Network.Default || ip.IsUnspecifiedIP(gw.IP) {
		log.Debugf("not setting route for network: default=%v gateway=%s", endpoint.Network.Default, gw.IP)
		return nil
	}

	_, defaultNet, _ := net.ParseCIDR("0.0.0.0/0")
	// delete default route first
	if err := t.RouteDel(&netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet}); err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.ESRCH {
			return fmt.Errorf("could not update default route: %s", err)
		}
	}

	// delete the default route for the bridge.out table, if it exists
	bTablePresent := bridgeTableExists(t)
	if bTablePresent {
		if err := t.RouteDel(&netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet, Table: bridgeTableNumber}); err != nil {
			if errno, ok := err.(syscall.Errno); !ok || errno != syscall.ESRCH {
				return fmt.Errorf("could not update default route for bridge.out table: %s", err)
			}
		}
	}

	log.Infof("Setting default gateway to %s", gw.IP)
	route := &netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet, Gw: gw.IP}
	if err := t.RouteAdd(route); err != nil {
		return fmt.Errorf("failed to add gateway route for endpoint %s: %s", endpoint.Network.Name, err)
	}

	if bTablePresent {
		route = &netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet, Gw: gw.IP, Table: bridgeTableNumber}
		if err := t.RouteAdd(route); err != nil {
			return fmt.Errorf("failed to add gateway route for table bridge.out for endpoint %s: %s", endpoint.Network.Name, err)
		}
	}

	log.Infof("updated default route to %s interface, gateway: %s", endpoint.Network.Name, gw.IP)
	return nil
}

func (t *BaseOperations) updateHosts(endpoint *NetworkEndpoint) error {
	log.Debugf("%+v", endpoint)
	// Add /etc/hosts entry
	if endpoint.Network.Name == "" {
		return nil
	}

	Sys.Hosts.SetHost(fmt.Sprintf("%s.localhost", endpoint.Network.Name), endpoint.Assigned.IP)

	if err := Sys.Hosts.Save(); err != nil {
		return err
	}

	return nil
}

func (t *BaseOperations) updateNameservers(endpoint *NetworkEndpoint) error {
	ns := endpoint.Network.Assigned.Nameservers
	gw := endpoint.Network.Assigned.Gateway
	// Add nameservers
	// This is incredibly trivial for now - should be updated to a less messy approach
	if len(ns) > 0 {
		Sys.ResolvConf.AddNameservers(ns...)
		log.Infof("Added nameservers: %+v", ns)
	} else if !ip.IsUnspecifiedIP(gw.IP) {
		Sys.ResolvConf.AddNameservers(gw.IP)
		log.Infof("Added nameserver: %s", gw.IP)
	}

	if err := Sys.ResolvConf.Save(); err != nil {
		return err
	}

	return nil
}

func (t *BaseOperations) Apply(endpoint *NetworkEndpoint) error {
	defer trace.End(trace.Begin("applying endpoint configuration for " + endpoint.Network.Name))

	return apply(t, t, endpoint)
}

func apply(nl Netlink, t *BaseOperations, endpoint *NetworkEndpoint) error {
	if endpoint.configured {
		log.Infof("skipping applying config for network %s as it has been applied already", endpoint.Network.Name)
		return nil // already applied
	}

	// Locate interface
	slot, err := strconv.Atoi(endpoint.ID)
	if err != nil {
		return fmt.Errorf("endpoint ID must be a base10 numeric pci slot identifier: %s", err)
	}

	defer func() {
		if err == nil {
			log.Infof("successfully applied config for network %s", endpoint.Network.Name)
			endpoint.configured = true
		}
	}()

	var link netlink.Link
	link, err = nl.LinkBySlot(int32(slot))
	if err != nil {
		return fmt.Errorf("unable to acquire reference to link %s: %s", endpoint.ID, err)
	}

	// rename the link if needed
	link, err = renameLink(nl, link, int32(slot), endpoint)
	if err != nil {
		return fmt.Errorf("unable to reacquire link %s after rename pass: %s", endpoint.ID, err)
	}

	var dc client.Client
	defer func() {
		if err != nil && dc != nil {
			dc.Release()
		}
	}()

	var newIP *net.IPNet

	if endpoint.IsDynamic() && endpoint.DHCP == nil {
		if e, ok := t.dynEndpoints[endpoint.ID]; ok {
			// endpoint shares NIC, copy over DHCP
			endpoint.DHCP = e[0].DHCP
		}
	}

	log.Debugf("%+v", endpoint)
	if endpoint.IsDynamic() {
		if endpoint.DHCP == nil {
			dc, err = getDynamicIP(nl, link, endpoint)
			if err != nil {
				return err
			}

			ack := dc.LastAck()
			endpoint.DHCP = &DHCPInfo{
				Assigned:    net.IPNet{IP: ack.YourIP(), Mask: ack.SubnetMask()},
				Nameservers: ack.DNS(),
				Gateway:     net.IPNet{IP: ack.Gateway(), Mask: ack.SubnetMask()},
			}
		}
		newIP = &endpoint.DHCP.Assigned
	} else {
		newIP = endpoint.IP
		if newIP.IP.Equal(net.IPv4zero) {
			// managed externally
			return nil
		}
	}

	var old *net.IPNet
	if !ip.IsUnspecifiedIP(endpoint.Assigned.IP) {
		old = &endpoint.Assigned
	}

	if err = linkAddrUpdate(old, newIP, nl, link); err != nil {
		return err
	}

	updateEndpoint(newIP, endpoint)

	if err = updateRoutes(nl, link, endpoint); err != nil {
		return err
	}

	if err = t.updateHosts(endpoint); err != nil {
		return err
	}

	Sys.ResolvConf.RemoveNameservers(endpoint.Network.Assigned.Nameservers...)
	if err = t.updateNameservers(endpoint); err != nil {
		return err
	}

	if endpoint.IsDynamic() {
		eps := t.dynEndpoints[endpoint.ID]
		found := false
		for _, e := range eps {
			if e == endpoint {
				found = true
				break
			}
		}

		if !found {
			eps = append(eps, endpoint)
			t.dynEndpoints[endpoint.ID] = eps
		}
	}

	// add renew/release loop if necessary
	if dc != nil {
		if _, ok := t.dhcpLoops[endpoint.ID]; !ok {
			stop := make(chan struct{})
			if err != nil {
				log.Errorf("could not make DHCP client id for link %s: %s", link.Attrs().Name, err)
			} else {
				t.dhcpLoops[endpoint.ID] = stop
				go t.dhcpLoop(stop, endpoint, dc)
			}
		}
	}

	return nil
}

func (t *BaseOperations) dhcpLoop(stop chan struct{}, e *NetworkEndpoint, dc client.Client) {
	exp := time.After(dc.LastAck().LeaseTime() / 2)
	for {
		select {
		case <-stop:
			// release the ip
			log.Infof("releasing IP address for network %s", e.Name)
			dc.Release()
			return

		case <-exp:
			log.Infof("renewing IP address for network %s", e.Name)
			err := dc.Renew()
			if err != nil {
				log.Errorf("failed to renew ip address for network %s: %s", e.Name, err)
				continue
			}

			ack := dc.LastAck()
			log.Infof("successfully renewed ip address: IP=%s, SubnetMask=%s, Gateway=%s, DNS=%s, Lease Time=%s", ack.YourIP(), ack.SubnetMask(), ack.Gateway(), ack.DNS(), ack.LeaseTime())

			e.DHCP = &DHCPInfo{
				Assigned:    net.IPNet{IP: ack.YourIP(), Mask: ack.SubnetMask()},
				Gateway:     net.IPNet{IP: ack.Gateway(), Mask: ack.SubnetMask()},
				Nameservers: ack.DNS(),
			}

			e.configured = false
			t.Apply(e)
			if err = t.config.UpdateNetworkEndpoint(e); err != nil {
				log.Error(err)
			}
			// update any endpoints that share this NIC
			for _, d := range t.dynEndpoints[e.ID] {
				if e == d {
					// already applied above
					continue
				}

				d.DHCP = e.DHCP
				d.configured = false
				t.Apply(d)
				if err = t.config.UpdateNetworkEndpoint(d); err != nil {
					log.Error(err)
				}
			}

			t.config.Flush()

			exp = time.After(ack.LeaseTime() / 2)
		}
	}
}

// MountLabel performs a mount with the label and target being absolute paths
func (t *BaseOperations) MountLabel(ctx context.Context, label, target string) error {
	defer trace.End(trace.Begin(fmt.Sprintf("Mounting %s on %s", label, target)))

	if err := os.MkdirAll(target, 0644); err != nil {
		// same as MountFileSystem error for consistency
		return fmt.Errorf("unable to create mount point %s: %s", target, err)
	}

	// convert the label to a filesystem path
	label = "/dev/disk/by-label/" + label

	// do..while ! timedout
	var timeout bool
	for timeout = false; !timeout; {
		_, err := os.Stat(label)
		if err == nil || !os.IsNotExist(err) {
			break
		}

		deadline, ok := ctx.Deadline()
		timeout = ok && time.Now().After(deadline)
	}

	if timeout {
		detail := fmt.Sprintf("timed out waiting for %s to appear", label)
		return errors.New(detail)
	}

	if err := Sys.Syscall.Mount(label, target, ext4FileSystemType, syscall.MS_NOATIME, ""); err != nil {
		// consistent with MountFileSystem
		detail := fmt.Sprintf("mounting %s on %s failed: %s", label, target, err)
		return errors.New(detail)
	}

	return nil
}

// MountTarget performs a mount based on the target path from the source url
// This assumes that the source url is valid and available.
func (t *BaseOperations) MountTarget(ctx context.Context, source url.URL, target string, mountOptions string) error {
	defer trace.End(trace.Begin(fmt.Sprintf("Mounting %s on %s", source.String(), target)))

	if err := os.MkdirAll(target, 0644); err != nil {
		// same as MountLabel error for consistency
		return fmt.Errorf("unable to create mount point %s: %s", target, err)
	}

	rawSource := source.Hostname() + ":/" + source.Path
	// NOTE: by default we are supporting "NOATIME" and it can be configurable later. this must be specfied as a flag.
	// Additionally, we must parse out the "ro" option and supply it as a flag as well for this flavor of the mount call.
	if err := Sys.Syscall.Mount(rawSource, target, nfsFileSystemType, syscall.MS_NOATIME, mountOptions); err != nil {
		log.Errorf("mounting %s on %s failed: %s", source.String(), target, err)
		return err
	}

	return nil
}

// CopyExistingContent copies the underlying files shadowed by a mount on a directory
// to the volume mounted on the directory
// see bug https://github.com/vmware/vic/issues/3482
func (t *BaseOperations) CopyExistingContent(source string) error {
	defer trace.End(trace.Begin(fmt.Sprintf("copyExistingContent from %s", source)))

	source = filepath.Clean(source)

	// if mounted volume is not empty skip the copy task
	if empty, err := isEmpty(source); err != nil || !empty {
		if err != nil {
			log.Errorf("error checking directory for contents %s: %+v", source, err)
			return err
		}
		log.Debugf("Skipping copy as volume %s is not empty", source)
		return nil
	}

	log.Debugf("creating directory %s", bindDir)
	if err := os.MkdirAll(bindDir, 0644); err != nil {
		log.Errorf("error creating directory %s: %+v", bindDir, err)
		return err
	}

	// remove dir
	defer func() {
		log.Debugf("removing %s", bindDir)
		if err := os.Remove(bindDir); err != nil {
			log.Errorf("error removing directory %s: %+v", bindDir, err)
		}
	}()

	parentDir := filepath.Dir(source)
	// mount the parent directory of the source to bindDir
	// e.g if source is /foo/bar, mount /foo to ./bindDir
	log.Debugf("mounting %s on %s", parentDir, bindDir)
	if err := Sys.Syscall.Mount(parentDir, bindDir, ext4FileSystemType, syscall.MS_BIND, ""); err != nil {
		log.Errorf("error mounting to %s: %+v", bindDir, err)
		return err
	}

	// unmount
	defer func() {
		log.Debugf("unmounting %s", bindDir)
		if err := Sys.Syscall.Unmount(bindDir, syscall.MNT_DETACH); err != nil {
			log.Errorf("error unmounting %+v", err)
		}
	}()

	mountedSource := filepath.Join(bindDir, filepath.Base(source))
	// copy data from the bindDir to the source
	// e.g if source is /foo/bar, copy ./bindDir/bar to /foo/bar
	log.Debugf("copying contents from to %s to %s", mountedSource, source)
	if err := archive.CopyWithTar(mountedSource, source); err != nil {
		log.Errorf("err copying %s to %s: %+v", mountedSource, source, err)
		return err
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

func (t *BaseOperations) Setup(config Config) error {
	err := Sys.Hosts.Load()
	if err != nil {
		return err
	}

	// make sure localhost entries are present
	entries := []struct {
		hostname string
		addr     net.IP
	}{
		{"localhost", net.ParseIP("127.0.0.1")},
		{"ip6-localhost", net.ParseIP("::1")},
		{"ip6-loopback", net.ParseIP("::1")},
		{"ip6-localnet", net.ParseIP("fe00::0")},
		{"ip6-mcastprefix", net.ParseIP("ff00::0")},
		{"ip6-allnodes", net.ParseIP("ff02::1")},
		{"ip6-allrouters", net.ParseIP("ff02::2")},
	}

	for _, e := range entries {
		Sys.Hosts.SetHost(e.hostname, e.addr)
	}

	if err = Sys.Hosts.Save(); err != nil {
		return err
	}

	t.dynEndpoints = make(map[string][]*NetworkEndpoint)
	t.dhcpLoops = make(map[string]chan struct{})
	t.config = config

	return nil
}

func (t *BaseOperations) Cleanup() error {
	for _, stop := range t.dhcpLoops {
		close(stop)
	}

	return nil
}

// Need to put this here because Windows does not support SysProcAttr.Credential
// getUserSysProcAttr relies on docker user package to verify user specification
// Examples of valid user specifications are:
//     * ""
//     * "user"
//     * "uid"
//     * "user:group"
//     * "uid:gid
//     * "user:gid"
//     * "uid:group"
func getUserSysProcAttr(uid, gid string) (*syscall.SysProcAttr, error) {
	if uid == "" && gid == "" {
		log.Debugf("no user id or group id specified")
		return nil, nil
	}

	userGroup := uid
	if gid != "" {
		userGroup = fmt.Sprintf("%s:%s", uid, gid)
	}
	passwdPath, err := user.GetPasswdPath()
	if err != nil {
		return nil, err
	}
	groupPath, err := user.GetGroupPath()
	if err != nil {
		return nil, err
	}
	execUser, err := user.GetExecUserPath(userGroup, defaultExecUser, passwdPath, groupPath)
	if err != nil {
		return nil, err
	}

	sysProc := &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(execUser.Uid),
			Gid: uint32(execUser.Gid),
		},
		Setsid: true,
	}
	for _, sgid := range execUser.Sgids {
		sysProc.Credential.Groups = append(sysProc.Credential.Groups, uint32(sgid))
	}
	return sysProc, nil
}

// isEmpty returns true if the directory is empty or contains a lost+found folder
func isEmpty(name string) (bool, error) {
	files, err := readDir(name)
	if err != nil || len(files) > 0 {
		return false, err
	}
	return true, nil
}

// readDir reads a directory and hides a specific dir "lost+found"
func readDir(dir string) ([]os.FileInfo, error) {
	lostnfound := "lost+found"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := files[:0]
	for _, f := range files {
		if f.Name() != lostnfound {
			result = append(result, f)
		}
	}

	return result, nil
}
