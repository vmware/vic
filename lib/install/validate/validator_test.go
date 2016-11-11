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
	"net"
	"net/url"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/simulator"

	"context"
)

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

func TestMain(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	trace.Logger.Level = log.DebugLevel
	ctx := context.Background()

	for i, model := range []*simulator.Model{simulator.ESX(), simulator.VPX()} {
		t.Logf("%d", i)
		model.Datastore = 3
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

		var input *data.Data
		if i == 0 {
			input = getESXData(s.URL)
		} else {
			input = getVPXData(s.URL)
		}
		if err != nil {
			t.Fatal(err)
		}

		validator, err := NewValidator(ctx, input)
		if err != nil {
			t.Errorf("Failed to new validator: %s", err)
		}
		ds, _ := validator.Session.Finder.Datastore(validator.Context, "LocalDS_1")
		simulator.Map.Get(ds.Reference()).(mo.Entity).Entity().Name = "Local DS_0"

		ds, _ = validator.Session.Finder.Datastore(validator.Context, "LocalDS_2")
		simulator.Map.Get(ds.Reference()).(mo.Entity).Entity().Name = `😗`

		validator.DisableFirewallCheck = true
		validator.DisableDRSCheck = true
		t.Logf("session pool: %s", validator.Session.Pool)
		if err = createPool(ctx, validator.Session, input.ComputeResourcePath, "validator", t); err != nil {
			t.Errorf("Unable to create resource pool: %s", err)
		}

		conf := testCompute(validator, input, t)
		testTargets(validator, input, conf, t)
		testStorage(validator, input, conf, t)
		//		testNetwork() need dvs support
	}
}

func getESXData(url *url.URL) *data.Data {
	result := data.NewData()
	url.Path = url.Path + "/ha-datacenter"
	result.URL = url
	result.DisplayName = "test001"
	result.ComputeResourcePath = "/ha-datacenter/host/localhost.localdomain/Resources"
	result.ImageDatastorePath = "LocalDS_0"
	result.BridgeNetworkName = "bridge"
	_, result.BridgeIPRange, _ = net.ParseCIDR("172.16.0.0/12")
	result.ManagementNetwork.Name = "VM Network"
	result.ExternalNetwork.Name = "VM Network"
	result.VolumeLocations = make(map[string]string)
	result.VolumeLocations["volume-store"] = "LocalDS_0/volumes/test"

	return result
}

func getVPXData(url *url.URL) *data.Data {
	result := data.NewData()
	url.Path = url.Path + "/DC0"
	result.URL = url
	result.DisplayName = "test001"
	result.ComputeResourcePath = "/DC0/host/DC0_C0/Resources"
	result.ImageDatastorePath = "LocalDS_0"
	result.ExternalNetwork.Name = "VM Network"
	result.BridgeNetworkName = "bridge"
	_, result.BridgeIPRange, _ = net.ParseCIDR("172.16.0.0/12")
	result.VolumeLocations = make(map[string]string)
	result.VolumeLocations["volume-store"] = "LocalDS_0/volumes/test"

	return result
}

func createPool(ctx context.Context, sess *session.Session, poolPath string, name string, t *testing.T) error {
	rp, err := sess.Finder.ResourcePool(ctx, poolPath)
	if err != nil {
		t.Logf("Failed to get parent pool: %s", err)
		return err
	}
	t.Logf("Creating Resource Pool %s", name)
	resSpec := types.ResourceConfigSpec{
		CpuAllocation: &types.ResourceAllocationInfo{
			Shares: &types.SharesInfo{
				Level: types.SharesLevelNormal,
			},
			ExpandableReservation: types.NewBool(true),
			Limit:       -1,
			Reservation: 1,
		},
		MemoryAllocation: &types.ResourceAllocationInfo{
			Shares: &types.SharesInfo{
				Level: types.SharesLevelNormal,
			},
			ExpandableReservation: types.NewBool(true),
			Limit:       -1,
			Reservation: 1,
		},
	}

	_, err = rp.Create(ctx, name, resSpec)
	if err != nil {
		t.Logf("Failed to create resource pool %s: %s", name, err)
		return err
	}
	return nil
}

func testCompute(v *Validator, input *data.Data, t *testing.T) *config.VirtualContainerHostConfigSpec {
	tests := []struct {
		path   string
		vc     bool
		hasErr bool
	}{
		{"DC0_C0/Resources/validator", true, false},
		{"DC0_C0/validator", true, false},
		{"validator", true, true},
		{"DC0_C0/test", true, true},
		{"/DC0_C1/test", true, true},
		{"/DC0_C1/test", true, true},
		{"/DC0/host/DC0_C1/Resources/validator", true, true},
		{"/DC1/host/DC0_C1/Resources/validator", true, true},
		{"DC0_H0/Resources", true, false},
		{"DC0_H0", true, false},
		{"/DC0/host/DC0_C0/Resources/validator", true, false},
		{"localhost.localdomain/Resources/validator", false, false},
		{"validator", false, false},
		{"test", false, true},
		{"/ha-datacenter/host/localhost.localdomain/Resources/validator", false, false},
	}
	conf := &config.VirtualContainerHostConfigSpec{}

	for _, test := range tests {
		if v.isVC && !test.vc {
			continue
		}
		if !v.isVC && test.vc {
			continue
		}
		t.Logf("%+v", test)
		input.ComputeResourcePath = test.path
		v.compute(v.Context, input, conf)
		v.ListIssues()
		if !test.hasErr {
			assert.Equal(t, 0, len(v.issues))
		} else {
			assert.True(t, len(v.issues) > 0, "Should have errors")
		}
		v.issues = nil
	}
	return conf
}

func testTargets(v *Validator, input *data.Data, conf *config.VirtualContainerHostConfigSpec, t *testing.T) {
	v.target(v.Context, input, conf)
	u, err := url.Parse(conf.Target)
	assert.NoError(t, err)
	assert.Nil(t, u.User)
	assert.NotEmpty(t, conf.Token)
	assert.NotEmpty(t, conf.Username)
}

func testStorage(v *Validator, input *data.Data, conf *config.VirtualContainerHostConfigSpec, t *testing.T) {
	tests := []struct {
		image         string
		volumes       map[string]string
		hasErr        bool
		expectImage   string
		expectVolumes map[string]string
	}{
		{"LocalDS_0",
			map[string]string{"volume1": "LocalDS_0/volumes/volume1",
				"volume2": "ds://LocalDS_0/volumes/volume2"},
			false,
			"ds://LocalDS_0/test001",
			map[string]string{"volume1": "ds://LocalDS_0/volumes/volume1",
				"volume2": "ds://LocalDS_0/volumes/volume2"}},

		{"LocalDS_0/images",
			map[string]string{"volume1": "LocalDS_0/volumes/volume1",
				"volume2": "ds://LocalDS_0/volumes/volume2"},
			false,
			"ds://LocalDS_0/images",
			map[string]string{"volume1": "ds://LocalDS_0/volumes/volume1",
				"volume2": "ds://LocalDS_0/volumes/volume2"}},

		{"ds://LocalDS_0/images",
			map[string]string{"volume1": "LocalDS_0/volumes/volume1",
				"volume2": "ds://LocalDS_0/volumes/volume2"},
			false,
			"ds://LocalDS_0/images",
			map[string]string{"volume1": "ds://LocalDS_0/volumes/volume1",
				"volume2": "ds://LocalDS_0/volumes/volume2"}},

		{"ds://LocalDS_0/images/xyz",
			map[string]string{"volume1": "LocalDS_0/volumes/volume1",
				"volume2": "ds://LocalDS_0/volumes/volume2"},
			false,
			"ds://LocalDS_0/images/xyz",
			map[string]string{"volume1": "ds://LocalDS_0/volumes/volume1",
				"volume2": "ds://LocalDS_0/volumes/volume2"}},

		{"ds://😗",
			map[string]string{"volume1": "😗/volumes/volume1",
				"volume2": "ds://😗/volumes/volume2"},
			true,
			"ds://😗/test001",
			nil},

		{"ds://LocalDS_0",
			map[string]string{"volume1": "LocalDS_1/volumes/volume1",
				"volume2": "ds://LocalDS_1/volumes/volume2"},
			true,
			"ds://LocalDS_0/test001",
			nil},

		{"LocalDS_0",
			map[string]string{"volume1": "LocalDS_1/volumes/volume1",
				"volume2": "ds://LocalDS_1/volumes/volume2"},
			true,
			"ds://LocalDS_0/test001",
			nil},

		{"LocalDS_0",
			map[string]string{"volume1": "LocalDS_1/volumes/volume1",
				"volume2": "ds://LocalDS_1/volumes/volume2"},
			true,
			"ds://LocalDS_0/test001",
			nil},

		{"",
			map[string]string{"volume1": "",
				"volume2": "ds://"},
			true,
			"",
			nil},

		{"ds://",
			map[string]string{"volume1": "",
				"volume2": "ds://"},
			true,
			"",
			nil},
	}

	for _, test := range tests {
		t.Logf("%+v", test)
		input.ImageDatastorePath = test.image
		input.VolumeLocations = test.volumes
		v.storage(v.Context, input, conf)
		v.ListIssues()
		if !test.hasErr {
			assert.Equal(t, 0, len(v.issues))
			assert.Equal(t, test.expectImage, conf.ImageStores[0].String())
			conf.ImageStores = conf.ImageStores[1:]
			for key, volume := range conf.VolumeLocations {
				assert.Equal(t, test.expectVolumes[key], volume.String())
			}
		} else {
			assert.True(t, len(v.issues) > 0, "Should have errors")
		}
		v.issues = nil
	}
}
