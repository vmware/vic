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

package vsphere

import (
	"os"
	"os/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

func URL(t *testing.T) string {
	s := os.Getenv("TEST_URL")
	if s == "" {
		t.SkipNow()
	}
	return s
}

func TestAttach(t *testing.T) {
	// need to be root and on esx to run this test
	u, err := user.Current()
	if !assert.NoError(t, err) {
		return
	}

	if u.Uid != "0" {
		t.SkipNow()
		return
	}

	config := &session.Config{
		Service:        URL(t),
		Insecure:       true,
		Keepalive:      time.Duration(5) * time.Minute,
		DatacenterPath: "",
		DatastorePath:  "/ha-datacenter/datastore/*",
		HostPath:       "/ha-datacenter/host/*/*",
		NetworkPath:    "/ha-datacenter/network/*",
		PoolPath:       "/ha-datacenter/host/*/Resources",
	}
	client, err := session.NewSession(config).Create(context.Background())
	if !assert.NoError(t, err) {
		return
	}

	s, err := GetSelf(client)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.NotNil(t, s) {
		return
	}
}
