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

package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/vishvananda/netlink"
)

func (ch *GlobalHandler) DynamicIPAddress() (string, error) {
	cmd := exec.Command("/.tether/udhcpc", "-R", "-b", "-i", "eth0", "-p", "/var/run/udhcpc.eth0.pid", "-s", "/.tether/udhcpc.script")
	if out, err := cmd.CombinedOutput(); err != nil {
		detail := fmt.Sprintf("failed to set IP for eth0: %s\n%s", err, string(out))
		log.Print(detail)
		return "", errors.New(detail)
	}

	link, err := netlink.LinkByName("eth0")
	if err != nil {
		detail := fmt.Sprintf("failed to get link for eth0: %s", err)
		log.Print(detail)
		return "", errors.New(detail)
	}

	existingAddr, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		detail := fmt.Sprintf("failed to list existing address for eth0: %s", err)
		log.Print(detail)
		return "", errors.New(detail)
	}

	if len(existingAddr) > 0 {
		return existingAddr[0].String(), nil
	}

	// no IPv4 address on the interface
	return "", nil
}

func (ch *GlobalHandler) StaticIPAddress(cidr, gateway string) error {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName("eth0")
	if err != nil {
		detail := fmt.Sprintf("failed to get link for eth0: %s", err)
		log.Print(detail)
		return errors.New(detail)
	}

	// set the ip address for our interface
	addr, err := netlink.ParseAddr(cidr)
	if err != nil {
		detail := fmt.Sprintf("failed to parse address for eth0: %s", err)
		log.Print(detail)
		return errors.New(detail)
	}

	existingAddr, err := netlink.AddrList(link, syscall.AF_UNSPEC)
	if err != nil {
		detail := fmt.Sprintf("failed to list existing address for eth0: %s", err)
		log.Print(detail)
		return errors.New(detail)
	}

	for _, oldAddr := range existingAddr {
		err = netlink.AddrDel(link, &oldAddr)
		if err != nil {
			detail := fmt.Sprintf("failed to del existing address for eth0: %s", err)
			log.Print(detail)
		}
	}

	if err := netlink.AddrAdd(link, addr); err != nil {
		detail := fmt.Sprintf("failed to add address to eth0: %s", err)
		log.Print(detail)
		return errors.New(detail)
	}

	// bring the interface up
	if err = netlink.LinkSetUp(link); err != nil {
		detail := fmt.Sprintf("failed to bring up eth0: %s", err)
		log.Print(detail)
		return errors.New(detail)
	}

	// add a default route
	//route := netlink.Route{LinkIndex: link.Attrs().Index, Dst: ipnet, Src: ip, Gw: net.ParseIP(gateway)}
	_, defaultNet, _ := net.ParseCIDR("0.0.0.0/0")
	route := netlink.Route{LinkIndex: link.Attrs().Index, Dst: defaultNet, Gw: net.ParseIP(gateway)}
	err = netlink.RouteAdd(&route)
	if err != nil {
		detail := fmt.Sprintf("failed to add route for gateway: %s", err)
		log.Print(detail)
		return errors.New(detail)
	}
	/*
		cmd := exec.Command("/sbin/ip", "address", "add", cidr, "dev", "eth0")
		if out, err := cmd.CombinedOutput(); err != nil {
			detail := fmt.Sprintf("failed to set IP for eth0: %s\n%s", err, string(out))
			log.Print(detail)
			return errors.New(detail)
		}

		cmd = exec.Command("/sbin/ip", "link", "set", "dev", "eth0", "up")
		if out, err := cmd.CombinedOutput(); err != nil {
			detail := fmt.Sprintf("failed to bring up eth0: %s\n%s", err, string(out))
			log.Print(detail)
			return errors.New(detail)
		}

		cmd = exec.Command("/sbin/ip", "route", "add", "default", "dev", "eth0", "via", gateway)
		if out, err := cmd.CombinedOutput(); err != nil {
			detail := fmt.Sprintf("failed to add route for gateway: %s\n%s", err, string(out))
			log.Print(detail)
			return errors.New(detail)
		}
	*/

	// TODO: even if we want to keep an entry in the local /etc/hosts to speed resolution
	// we should do it in a maner that doesn't result in it being present in a diff with the parent
	f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	defer f.Close()

	hostname, err := os.Hostname()
	if err != nil {
		log.Print(err.Error())
		return err
	}

	if _, err = f.WriteString(fmt.Sprintf("%s %s\n", ip.String(), hostname)); err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func (c *GlobalHandler) MountLabel(label, target string) error {
	if err := os.MkdirAll(target, 0600); err != nil {
		return fmt.Errorf("Unable to create mount point %s: %s", target, err)
	}

	//volumes := "/.tether/volumes"
	volumes := "/dev/disk/by-label"

	source := volumes + "/" + label

	// wait for mount source to appear or timeout
	for start := time.Now(); time.Since(start) < (5 * time.Second); {
		_, err := os.Stat(source)
		if err == nil || !os.IsNotExist(err) {
			break
		}
	}

	if err := syscall.Mount(source, target, "ext4", syscall.MS_NOATIME, ""); err != nil {
		detail := fmt.Sprintf("Unable to mount %s: %s", source, err)
		log.Print(detail)
		// for debug purposes, dump the directory listing of volumes and /dev/disk/by-label
		for _, dir := range []string{volumes, "/dev/disk/by-label"} {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Printf("unable to read listing for %s: %s\n", dir, err)
				continue
			}

			log.Printf("%s/\n", dir)
			for _, fi := range files {
				log.Printf("\t%s\n", fi.Name())
			}
		}

		return errors.New(detail)
	}

	return nil
}

func (c *GlobalHandler) Sync() {
	syscall.Sync()
}
