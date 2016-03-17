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

package guest

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/spec"
	"golang.org/x/net/context"
)

func URL(t *testing.T) string {
	s := os.Getenv("TEST_URL")
	if s == "" {
		t.SkipNow()
	}
	return s
}

func TestNewLinuxGuest(t *testing.T) {

	ctx := context.Background()

	sessionconfig := &session.Config{
		Service:        URL(t),
		Insecure:       true,
		Keepalive:      time.Duration(5) * time.Minute,
		DatacenterPath: "",
		DatastorePath:  "/ha-datacenter/datastore/*",
		HostPath:       "/ha-datacenter/host/*/*",
		NetworkPath:    "/ha-datacenter/network/*",
		PoolPath:       "/ha-datacenter/host/*/Resources",
	}

	session, err := session.NewSession(sessionconfig).Create(ctx)
	if err != nil {
		t.Logf("%+v", err.Error())
		if _, ok := err.(*find.MultipleFoundError); !ok {
			t.Errorf(err.Error())
		} else {
			t.SkipNow()
		}
	}

	specconfig := &spec.VirtualMachineConfigSpecConfig{
		NumCPUs:       2,
		MemoryMB:      2048,
		VMForkEnabled: true,

		ConnectorURI: "tcp://1.2.3.4:9876",

		ID:            "zombie_attack",
		BootMediaPath: session.Datastore.Path("brainz.iso"),
		VMPathName:    fmt.Sprintf("[%s]", session.Datastore.Name()),
		NetworkName:   strings.Split(session.Network.Reference().Value, "-")[0],
	}
	t.Logf("%+v", specconfig)

	root := NewLinuxGuest(session, specconfig)
	assert.Equal(t, "other3xLinux64Guest", root.GuestId)
}
