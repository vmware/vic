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

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"syscall"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vmw-guestinfo/rpcout"
)

var hostnameFile = "/etc/hostname"
var hostsFile = "/etc/hosts"
var resolvFile = "/etc/resolv.conf"
var byLabelDir = "/dev/disk/by-label"

type osopsLinux struct{}

// SetHostname sets both the kernel hostname and /etc/hostname to the specified string
func (t *osopsLinux) SetHostname(hostname string) error {
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
	err = ioutil.WriteFile(pathPrefix+hostnameFile, []byte(hostname), 0644)
	if err != nil {
		log.Errorf("Failed to update hostname in %s%s", pathPrefix, hostnameFile)

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
	hosts, err := os.OpenFile(pathPrefix+"/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		detail := fmt.Sprintf("failed to update hosts with hostname: %s", err)
		return errors.New(detail)
	}
	defer hosts.Close()

	_, err = hosts.WriteString(fmt.Sprintf("127.0.0.1 %s", hostname))
	if err != nil {
		detail := fmt.Sprintf("failed to add hosts entry for hostname %s: %s", hostname, err)
		return errors.New(detail)
	}

	return nil
}

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
func (t *osopsLinux) Apply(endpoint *metadata.NetworkEndpoint) error {
	defer trace.End(trace.Begin("applying endpoint configuration for " + endpoint.Network.Name))

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

	if err = netlink.AddrAdd(link, addr); err != nil {
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
	hosts, err := os.OpenFile(pathPrefix+hostsFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		detail := fmt.Sprintf("failed to update hosts for endpoint %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}
	defer hosts.Close()

	_, err = hosts.WriteString(fmt.Sprintf("localhost.%s %s", endpoint.Network.Name, endpoint.IP.IP))
	if err != nil {
		detail := fmt.Sprintf("failed to add hosts entry for endpoint %s: %s", endpoint.Network.Name, err)
		return errors.New(detail)
	}

	// Add nameservers
	// This is incredibly trivial for now - should be updated to a less messy approach
	resolv, err := os.OpenFile(pathPrefix+resolvFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		detail := fmt.Sprintf("failed to update %s for endpoint %s: %s", resolvFile, endpoint.Network.Name, err)
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

// MountLabel performs a mount with the source treated as a disk label
// This assumes that /dev/disk/by-label is being populated, probably by udev
func (t *osopsLinux) MountLabel(label, target string, ctx context.Context) error {
	defer trace.End(trace.Begin(fmt.Sprintf("Mounting %s on %s", label, target)))

	if err := os.MkdirAll(target, 0600); err != nil {
		return fmt.Errorf("unable to create mount point %s: %s", target, err)
	}

	volumes := pathPrefix + byLabelDir
	source := volumes + "/" + label

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

// Fork triggers vmfork and handles the necessary pre/post OS level operations
func (t *osopsLinux) Fork(config *metadata.ExecutorConfig) error {
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

func MkNamedPipe(path string, mode os.FileMode) error {
	return syscall.Mkfifo(path, (uint32(mode)))
}
