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
	"net"
	"testing"

	"github.com/docker/docker/pkg/stringid"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/lib/metadata"
)

func TestSetHostname(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	cfg := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "sethostname",
			Name: "tether_test_executor",
		},
	}

	tthr, _ := StartTether(t, &cfg)

	<-Mocked.Started

	// prevent indefinite wait in tether - normally session exit would trigger this
	tthr.Stop()

	// wait for tether to exit
	<-Mocked.Cleaned

	expected := stringid.TruncateID(cfg.ID)
	if Mocked.Hostname != expected {
		t.Errorf("expected: %s, actual: %s", expected, Mocked.Hostname)
	}
}

func TestNoNetwork(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	cfg := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "ipconfig",
			Name: "tether_test_executor",
		},
	}

	tthr, _ := StartTether(t, &cfg)

	<-Mocked.Started

	// prevent indefinite wait in tether - normally session exit would trigger this
	tthr.Stop()

	// wait for tether to exit
	<-Mocked.Cleaned
}

func TestSetIpAddress(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	cfg := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "ipconfig",
			Name: "tether_test_executor",
		},
		Networks: map[string]*metadata.NetworkEndpoint{
			"netA": &metadata.NetworkEndpoint{
				Network: metadata.ContainerNetwork{
					Common: metadata.Common{
						Name: "netA",
					},
				},
				Static: &net.IPNet{
					IP:   localhost,
					Mask: lmask.Mask,
				},
			},
			"netB": &metadata.NetworkEndpoint{
				Network: metadata.ContainerNetwork{
					Common: metadata.Common{
						Name: "netB",
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

	<-Mocked.Started

	// check addresses
	checkA, ok := Mocked.IPs["netA"]
	assert.True(t, ok, "Expected entry in IP map for netA")
	assert.Equal(t, localhost, checkA.IP, "netA has incorrect IP address")

	checkB, ok := Mocked.IPs["netB"]
	assert.True(t, ok, "Expected entry in IP map for netB")
	assert.Equal(t, gateway, checkB.IP, "netB has incorrect IP address")

	// prevent indefinite wait in tether - normally session exit would trigger this
	tthr.Stop()

	// wait for tether to exit
	<-Mocked.Cleaned
}
