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
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
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

const pciDevPath = "/sys/bus/pci/devices"

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
	hosts, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
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

func linkBySlot(slot int32) (netlink.Link, error) {
	pciPath, err := slotToPCIPath(slot)
	if err != nil {
		return nil, err
	}

	name, err := pciToLinkName(pciPath)
	if err != nil {
		return nil, err
	}

	log.Debugf("got link name: %#v", name)
	return netlink.LinkByName(name)
}

// Apply takes the network endpoint configuration and applies it to the system
func (t *osopsLinux) Apply(endpoint *metadata.NetworkEndpoint) error {
	defer trace.End(trace.Begin("applying endpoint configuration for " + endpoint.Network.Name))

	// Locate interface
	link, err := linkBySlot(endpoint.PCISlot)
	if err != nil {
		return err
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

	// Add routes
	_, defaultNet, _ := net.ParseCIDR("0.0.0.0/0")
	route := netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet, Gw: endpoint.Network.Gateway.IP}
	err = netlink.RouteAdd(&route)
	if err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EEXIST {
			detail := fmt.Sprintf("failed to add gateway route for endpoint %s: %s", endpoint.Network.Name, err)
			return errors.New(detail)
		}
	}

	// TODO update /etc/hosts and /etc/resolv.conf

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
func (t *osopsLinux) Fork(config *ExecutorConfig) error {
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
