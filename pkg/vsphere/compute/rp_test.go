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

package compute

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/simulator"

	"golang.org/x/net/context"
)

func TestMain(t *testing.T) {

	ctx := context.Background()

	for i, model := range []*simulator.Model{simulator.ESX(), simulator.VPX()} {
		t.Logf("%d", i)
		defer model.Remove()
		err := model.Create()
		if err != nil {
			t.Fatal(err)
		}

		s := model.Service.NewServer()
		defer s.Close()

		s.URL.User = url.UserPassword("user", "pass")
		t.Logf("server URL: %s", s.URL)

		var sess *session.Session
		if i == 0 {
			sess, err = getESXSession(ctx, s.URL.String())
		} else {
			sess, err = getVPXSession(ctx, s.URL.String())
		}
		if err != nil {
			t.Fatal(err)
		}
		defer sess.Logout(ctx)
		testGetChildrenVMs(ctx, sess, t)
		testGetChildVM(ctx, sess, t)
		testFindResourcePool(ctx, sess, t)
		testGetCluster(ctx, sess, t)
	}
}

func getESXSession(ctx context.Context, service string) (*session.Session, error) {
	config := &session.Config{
		Service:        service,
		Insecure:       true,
		Keepalive:      time.Duration(5) * time.Minute,
		DatacenterPath: "/ha-datacenter",
		ClusterPath:    "*",
		DatastorePath:  "/ha-datacenter/datastore/LocalDS_0",
		NetworkPath:    "/ha-datacenter/network/VM Network",
		PoolPath:       "/ha-datacenter/host/localhost.localdomain/Resources",
	}

	s, err := session.NewSession(config).Connect(ctx)
	if err != nil {
		return nil, err
	}
	s.Finder = find.NewFinder(s.Client.Client, false)
	if s, err = s.Populate(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func getVPXSession(ctx context.Context, service string) (*session.Session, error) {
	config := &session.Config{
		Service:        service,
		Insecure:       true,
		Keepalive:      time.Duration(5) * time.Minute,
		DatacenterPath: "/DC0",
		ClusterPath:    "/DC0/host/DC0_C0",
		DatastorePath:  "/DC0/datastore/LocalDS_0",
		PoolPath:       "/DC0/host/DC0_C0/Resources",
	}

	s, err := session.NewSession(config).Connect(ctx)
	if err != nil {
		return nil, err
	}
	s.Finder = find.NewFinder(s.Client.Client, false)
	if s, err = s.Populate(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func testGetChildrenVMs(ctx context.Context, sess *session.Session, t *testing.T) {
	rp := NewResourcePool(ctx, sess, sess.Pool.Reference())
	vms, err := rp.GetChildrenVMs(ctx, sess)
	if err != nil {
		t.Errorf("Failed to get children vm of resource pool %s, %s", rp.Name(), err)
	}
	//	if vms == nil || len(vms) == 0 {
	//		t.Error("Didn't get children VM")
	//	}
	for _, vm := range vms {
		t.Logf("vm: %s", vm)
	}
}

func testGetChildVM(ctx context.Context, sess *session.Session, t *testing.T) {
	rp := NewResourcePool(ctx, sess, sess.Pool.Reference())
	vm, err := rp.GetChildVM(ctx, sess, "random")
	if err == nil && vm != nil {
		t.Logf("vm: %s", vm.Reference())
		t.Errorf("Should not find VM random")
	}
}

func testFindResourcePool(ctx context.Context, sess *session.Session, t *testing.T) {
	tests := []struct {
		name   string
		hasErr bool
	}{
		{"", false},
		{"random123", true},
	}

	for _, test := range tests {
		_, err := FindResourcePool(ctx, sess, test.name)
		assert.Equal(t, test.hasErr, err != nil)
	}
}

func testGetCluster(ctx context.Context, sess *session.Session, t *testing.T) {
	rp := NewResourcePool(ctx, sess, sess.Pool.Reference())
	cluster, err := rp.GetCluster(ctx)
	if err != nil {
		t.Logf("Failed to owner cluster: %s", err)
		t.Errorf("Should get owner")
	}
	t.Logf("Cluster: %s", cluster)
}
