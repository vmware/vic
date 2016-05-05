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
	"net"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

var (
	localhost, lmask, _ = net.ParseCIDR("127.0.0.2/24")
	gateway, gmask, _   = net.ParseCIDR("127.0.0.1/24")
)

func TestToExtraConfig(t *testing.T) {
	exec := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "deadbeef",
			Name: "configtest",
		},
		Sessions: map[string]metadata.SessionConfig{
			"deadbeef": metadata.SessionConfig{
				Cmd: metadata.Cmd{
					Path: "/bin/bash",
					Args: []string{"/bin/bash", "-c", "echo hello"},
					Dir:  "/",
					Env:  []string{"HOME=/", "PATH=/bin"},
				},
			},
			"beefed": metadata.SessionConfig{
				Cmd: metadata.Cmd{
					Path: "/bin/bash",
					Args: []string{"/bin/bash", "-c", "echo goodbye"},
					Dir:  "/",
					Env:  []string{"HOME=/", "PATH=/bin"},
				},
			},
		},
		Networks: map[string]metadata.NetworkEndpoint{
			"eth0": metadata.NetworkEndpoint{
				IP:  net.IPNet{IP: localhost, Mask: lmask.Mask},
				MAC: "a-mac-address",
				Network: metadata.ContainerNetwork{
					Name:        "notsure",
					Gateway:     net.IPNet{IP: gateway, Mask: gmask.Mask},
					Nameservers: []net.IP{},
				},
			},
		},
	}

	// encode metadata package's ExecutorConfig
	encoded := map[string]string{}
	extraconfig.Encode(extraconfig.MapSink(encoded), exec)

	// decode into this package's ExecutorConfig
	var decoded ExecutorConfig
	extraconfig.DecodeLogLevel = log.DebugLevel
	extraconfig.Decode(extraconfig.MapSource(encoded), &decoded)

	// the networks should be identical
	assert.Equal(t, exec.Networks["eth0"], *(decoded.Networks["eth0"]))

	// the source and destination structs are different - we're doing a sparse comparison
	expected := exec.Sessions["deadbeef"]
	actual := *decoded.Sessions["deadbeef"]

	assert.Equal(t, expected.Cmd.Path, actual.Cmd.Path)
	assert.Equal(t, expected.Cmd.Args, actual.Cmd.Args)
	assert.Equal(t, expected.Cmd.Dir, actual.Cmd.Dir)
	assert.Equal(t, expected.Cmd.Env, actual.Cmd.Env)
}
