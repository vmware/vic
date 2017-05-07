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

package data

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	d := NewData()
	s := NewData()

	ipAddr, mask, _ := net.ParseCIDR("1.1.1.1/32")
	s.ClientNetwork.IP.IP = ipAddr
	s.ClientNetwork.IP.Mask = mask.Mask

	s.ContainerNetworks.MappedNetworks["container"] = "external"
	d.CopyNonEmpty(s)
	assert.Equal(t, s.ClientNetwork.IP.IP, d.ClientNetwork.IP.IP, "ip is not right")
	assert.Equal(t, s.ContainerNetworks, d.ContainerNetworks, "container network is not copied")
}
