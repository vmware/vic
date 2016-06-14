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

package validate

import (
	"net/url"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/cmd/vic-machine/data"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/simulator"

	"golang.org/x/net/context"
)

func TestMain(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.Background()

	//	for i, model := range []*simulator.Model{simulator.ESX(), simulator.VPX()} {
	model := simulator.ESX()
	i := 0
	{
		t.Logf("%d", i)
		defer model.Remove()
		err := model.Create()
		if err != nil {
			t.Fatal(err)
		}

		s := model.Service.NewServer()
		defer s.Close()

		s.URL.User = url.UserPassword("user", "pass")
		s.URL.Path = ""
		t.Logf("server URL: %s", s.URL)

		var sess *session.Session
		var input *data.Data
		if i == 0 {
			sess, err = getESXSession(ctx, s.URL.String())
			input = getESXData(s.URL)
		} else {
			sess, err = getVPXSession(ctx, s.URL.String())
			input = getVPXData(s.URL)
		}
		if err != nil {
			t.Fatal(err)
		}
		defer sess.Logout(ctx)
		testGetResourcePool(ctx, input, t)
		//		testValidate(ctx, input, t)
		testGetVCHConfigFailed(ctx, sess, input, t)
		testGetVCHFailed(ctx, sess, input, t)
	}
}

func getESXData(url *url.URL) *data.Data {
	result := data.NewData()
	result.URL = url
	result.DisplayName = "test001"
	result.ComputeResourcePath = "/ha-datacenter/host/localhost.localdomain/Resources"
	result.ImageDatastoreName = "LocalDS_0"
	result.BridgeNetworkName = "VM Network"
	result.ManagementNetworkName = "VM Network"
	result.ExternalNetworkName = "VM Network"
	return result
}

func getVPXData(url *url.URL) *data.Data {
	result := data.NewData()
	result.URL = url
	result.DisplayName = "test001"
	result.ComputeResourcePath = "/DC0/host/DC0_C0/Resources"
	result.ImageDatastoreName = "LocalDS_0"
	result.ExternalNetworkName = "VM Network"
	result.BridgeNetworkName = "VM Network"
	return result
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
	return s, nil
}

func getVPXSession(ctx context.Context, service string) (*session.Session, error) {
	config := &session.Config{
		Service:        service,
		Insecure:       true,
		Keepalive:      time.Duration(5) * time.Minute,
		DatacenterPath: "/DC0",
		DatastorePath:  "/DC0/datastore/LocalDS_0",
		PoolPath:       "/DC0/host/DC0_C0/Resources",
	}

	s, err := session.NewSession(config).Connect(ctx)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func testGetResourcePool(ctx context.Context, input *data.Data, t *testing.T) {
	v, _ := NewValidator(ctx, input)

	if rp, err := v.GetResourcePool(input); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Got rp: %s", rp)
	}
}

func testValidate(ctx context.Context, input *data.Data, t *testing.T) {
	v, _ := NewValidator(ctx, input)
	if vchConfig, err := v.Validate(ctx, input); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Got vch: %s", vchConfig.Name)
	}
}

func testGetVCHConfigFailed(ctx context.Context, session *session.Session, input *data.Data, t *testing.T) {
	orp, err := session.Finder.ResourcePool(ctx, session.PoolPath)
	if err != nil {
		t.Errorf("test failed for cannot get resource pool: %s", err)
	}
	v, err := NewValidator(ctx, input)
	v.Session = session
	_, err = v.GetVCHConfig(orp, "test001")
	if err == nil {
		t.Errorf("Should not get VCH VM")
	}
}

func testGetVCHFailed(ctx context.Context, session *session.Session, input *data.Data, t *testing.T) {
	v, _ := NewValidator(ctx, input)
	v.Session = session
	_, err := v.GetVCH(input)
	if err == nil {
		t.Errorf("Should not get VCH VM")
	}
}

type TestValidator struct {
	Validator
}

func TestPathConversionESX(t *testing.T) {
	v := &TestValidator{}

	v.DatacenterPath = "/ha-datacenter"

	ipath := v.computePathToInventoryPath("/")
	assert.Equal(t, "/ha-datacenter/host/*/Resources", ipath, "Expected top level resource pool")
	cpath := v.inventoryPathToComputePath(ipath)
	assert.Equal(t, "*", cpath, "Expected root resource specifier")

	ipath = v.computePathToInventoryPath("*")
	assert.Equal(t, "/ha-datacenter/host/*/Resources/*", ipath, "Expected top level resource pool")
	cpath = v.inventoryPathToComputePath(ipath)
	assert.Equal(t, "*/*", cpath, "Expected top level wildcard specifier")
}

func TestSampleConversion(t *testing.T) {
	v := &TestValidator{}
	v.DatacenterPath = "/ha-datacenter"

	translations := map[string]string{
		"/":              "/ha-datacenter/host/*/Resources",
		"*":              "/ha-datacenter/host/*/Resources/*",
		"testpool":       "/ha-datacenter/host/*/Resources/testpool",
		"test/deep/path": "/ha-datacenter/host/*/Resources/test/deep/path",
	}

	for in, expected := range translations {
		ipath := v.computePathToInventoryPath(in)
		assert.Equal(t, expected, ipath, "Translation did not match")
	}
}
