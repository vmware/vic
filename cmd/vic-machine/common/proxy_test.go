// Copyright 2018 VMware, Inc. All Rights Reserved.
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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessProxies(t *testing.T) {
	gurls := [...]string{
		"https://fully.qualified.example.com",
		"https://fully.qualified.example.com:443",
		"http://fully.qualified.example.com",
		"http://fully.qualified.example.com:80",
		"http://203.0.113.123",
		"http://[2001:DB8:0123::]",
	}

	burls := [...]string{
		"example.com",
		"example.com:80",
		"localhost",
		"localhost:80",
		"ftp://example.com",
		"httpd://example.com",
	}

	for _, ghttp := range gurls {
		for _, ghttps := range gurls {
			gproxy := Proxies{HTTPProxy: &ghttp, HTTPSProxy: &ghttps}

			_, _, err := gproxy.ProcessProxies()
			assert.NoError(t, err, "Expected %s and %s to be accepted", ghttp, ghttps)
			assert.True(t, gproxy.IsSet, "Expected proxy to be marked as set")
		}
	}

	for _, ghttp := range gurls {
		for _, bhttps := range burls {
			bproxy := Proxies{HTTPProxy: &ghttp, HTTPSProxy: &bhttps}

			_, _, err := bproxy.ProcessProxies()
			assert.Error(t, err, "Expected %s to be rejected", bhttps)
		}
	}

	for _, bhttp := range burls {
		for _, ghttps := range gurls {
			bproxy := Proxies{HTTPProxy: &bhttp, HTTPSProxy: &ghttps}

			_, _, err := bproxy.ProcessProxies()
			assert.Error(t, err, "Expected %s to be rejected", bhttp)
		}
	}

	for _, bhttp := range burls {
		for _, bhttps := range burls {
			bproxy := Proxies{HTTPProxy: &bhttp, HTTPSProxy: &bhttps}

			_, _, err := bproxy.ProcessProxies()
			assert.Error(t, err, "Expected %s and %s to be rejected", bhttp, bhttps)
		}
	}
}
