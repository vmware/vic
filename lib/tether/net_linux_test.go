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

// +build linux

package tether

import (
	"io/ioutil"
	"net"
	"strconv"
	"testing"

	"github.com/vishvananda/netlink"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/lib/metadata"
)

// Utility method to add an interface to Mocked
// This assigns the interface name and returns the "slot" as a string
func AddInterface(name string) string {
	Mocked.maxSlot++

	Mocked.Interfaces[name] = &Interface{
		LinkAttrs: netlink.LinkAttrs{
			Name:  name,
			Index: Mocked.maxSlot,
		},
		Up: true,
	}

	return strconv.Itoa(Mocked.maxSlot)
}

func TestSetIpAddress(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	hFile, err := ioutil.TempFile("", "vic_set_ip_test_hosts")
	if err != nil {
		t.Errorf("Failed to create tmp hosts file: %s", err)
	}
	rFile, err := ioutil.TempFile("", "vic_set_ip_test_resolv")
	if err != nil {
		t.Errorf("Failed to create tmp resolv file: %s", err)
	}

	// give us a hosts file we can modify
	defer func(hosts, resolv string) {
		hostsFile = hosts
		resolvFile = resolv
	}(hostsFile, resolvFile)

	hostsFile = hFile.Name()
	resolvFile = rFile.Name()

	bridge := AddInterface("eth1")
	external := AddInterface("eth2")

	secondIP, _ := netlink.ParseIPNet("172.16.0.10/24")
	gwIP, _ := netlink.ParseIPNet("172.16.0.1/24")
	cfg := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "ipconfig",
			Name: "tether_test_executor",
		},
		Networks: map[string]*metadata.NetworkEndpoint{
			"bridge": &metadata.NetworkEndpoint{
				Common: metadata.Common{
					ID: bridge,
					// interface rename
					Name: "bridge",
				},
				Network: metadata.ContainerNetwork{
					Common: metadata.Common{
						Name: "bridge",
					},
					Default: true,
					Gateway: *gwIP,
				},
				Static: &net.IPNet{
					IP:   localhost,
					Mask: lmask.Mask,
				},
			},
			"cnet": &metadata.NetworkEndpoint{
				Common: metadata.Common{
					ID: bridge,
					// no interface rename
				},
				Network: metadata.ContainerNetwork{
					Common: metadata.Common{
						Name: "cnet",
					},
				},
				Static: secondIP,
			},
			"external": &metadata.NetworkEndpoint{
				Common: metadata.Common{
					ID: external,
					// interface rename
					Name: "external",
				},
				Network: metadata.ContainerNetwork{
					Common: metadata.Common{
						Name: "external",
					},
				},
				Static: &net.IPNet{
					IP:   gateway,
					Mask: gmask.Mask,
				},
			},
		},
	}

	tthr, _ := StartTether(t, &cfg)

	defer func() {
		// prevent indefinite wait in tether - normally session exit would trigger this
		tthr.Stop()

		// wait for tether to exit
		<-Mocked.Cleaned
	}()

	<-Mocked.Started

	assert.NotNil(t, Mocked.Interfaces["bridge"], "Expected bridge network if endpoints applied correctly")
	// check addresses
	bIface, _ := Mocked.Interfaces["bridge"].(*Interface)
	assert.NotNil(t, bIface)

	assert.Equal(t, 2, len(bIface.Addrs), "Expected two addresses on bridge interface")

	eIface, _ := Mocked.Interfaces["external"].(*Interface)
	assert.NotNil(t, eIface)

	assert.Equal(t, 1, len(eIface.Addrs), "Expected one address on external interface")
}
