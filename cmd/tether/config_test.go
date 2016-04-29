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

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
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
				IP:  net.IPNet{IP: net.IP("127.0.0.2"), Mask: net.IPMask("255.255.255.0")},
				MAC: "a-mac-address",
				Network: metadata.ContainerNetwork{
					Name:        "notsure",
					Gateway:     net.IPNet{IP: net.IP("127.0.0.1"), Mask: net.IPMask("255.255.255.0")},
					Nameservers: []net.IP{},
				},
			},
		},
	}

	// encode metadata package's ExecutorConfig
	encoded := extraconfig.Encode(exec)
	// decode into this package's ExecutorConfig
	var decoded ExecutorConfig
	extraconfig.Decode(extraconfig.OptionValueSource(encoded), &decoded)

	assert.Equal(t, exec.Sessions["deadbeef"], *(decoded.Sessions["deadbeef"]))
	assert.Equal(t, exec.Sessions["beefed"], *(decoded.Sessions["beefed"]))
	assert.Equal(t, exec.Networks["eth0"], *(decoded.Networks["eth0"]))
}
