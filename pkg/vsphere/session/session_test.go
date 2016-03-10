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

package session

import (
	"os"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/vmware/govmomi/find"
)

func URL(t *testing.T) string {
	s := os.Getenv("TEST_URL")
	if s == "" {
		t.SkipNow()
	}
	return s
}

func TestSessionDefaults(t *testing.T) {
	ctx := context.Background()

	config := &Config{
		Service:  URL(t),
		Insecure: true,
	}

	session, err := NewSession(config).Create(ctx)
	if err != nil {
		t.Logf("%+v", err.Error())
		if _, ok := err.(*find.DefaultMultipleFoundError); !ok {
			t.Errorf(err.Error())
		} else {
			t.SkipNow()
		}
	}
	t.Logf("%+v", session)
}

func TestSession(t *testing.T) {
	ctx := context.Background()

	config := &Config{
		Service:        URL(t),
		Insecure:       true,
		Keepalive:      time.Duration(5) * time.Minute,
		DatacenterPath: "",
		DatastorePath:  "/ha-datacenter/datastore/*",
		HostPath:       "/ha-datacenter/host/*/*",
		NetworkPath:    "/ha-datacenter/network/*",
		PoolPath:       "/ha-datacenter/host/*/Resources",
	}

	session, err := NewSession(config).Create(ctx)
	if err != nil {
		t.Logf("%+v", err.Error())
		if _, ok := err.(*find.MultipleFoundError); !ok {
			t.Errorf(err.Error())
		} else {
			t.SkipNow()
		}
	}
	t.Logf("Session: %+v", session)

	t.Logf("IsVC: %t", session.IsVC())
	t.Logf("IsVSAN: %t", session.IsVSAN(ctx))
}
