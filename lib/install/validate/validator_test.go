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
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/simulator"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
)

type TestValidator struct {
	Validator
}

func TestParseURL(t *testing.T) {
	var hosts = []string{
		"host.domain.com",
		"host.domain.com:123",
		"1.2.3.4",
		"1.2.3.4:10",
		"[2001:4860:0:2001::68]",
		"[2001:db8:1f70::999:de8:7648:6e8]:123",
	}

	for _, urlString := range hosts {
		u, err := ParseURL(urlString)
		assert.Nil(t, err)
		assert.Equal(t, u.String(), "https://"+urlString)
		// Null the scheme
		u.Scheme = ""
		assert.Equal(t, u.String(), "//"+urlString)
		assert.Equal(t, u.Host, urlString)
	}

	// Add path to create a more significant URL
	var urls = []string{}

	for i, h := range hosts {
		url := fmt.Sprintf("%s/path%d/test", h, i)
		urls = append(urls, url)
	}

	for i, urlString := range urls {
		u, err := ParseURL(urlString)
		assert.Nil(t, err)
		assert.Equal(t, u.String(), "https://"+urlString)

		// Null the scheme
		u.Scheme = ""
		assert.Equal(t, u.String(), "//"+urlString)

		// Check host
		assert.Equal(t, u.Host, hosts[i])
		// Check path
		path := fmt.Sprintf("/path%d/test", i)
		assert.Equal(t, u.Path, path)
		// Check concatenation
		assert.Equal(t, u.Host+u.Path, urlString)
	}

	// Add an HTTP scheme to verify that it is preserved
	var urlsWithHTTPScheme = []string{}

	for _, u := range urls {
		uws := fmt.Sprintf("http://%s", u)
		urlsWithHTTPScheme = append(urlsWithHTTPScheme, uws)
	}

	for _, urlString := range urlsWithHTTPScheme {
		u, err := ParseURL(urlString)
		fmt.Printf("UrlString: %s\n", u.String())
		assert.Nil(t, err)
		assert.Equal(t, u.String(), urlString)
	}

	var invalidUrls = []string{
		"[2001:db8/path",
		"1.2.3.4\\path",
	}

	for _, urlString := range invalidUrls {
		_, err := ParseURL(urlString)
		assert.NotNil(t, err)
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
		t.Logf("session pool: %s", validator.Session.Pool)
		if err = createPool(ctx, validator.Session, input.ComputeResourcePath, "validator", t); err != nil {
			t.Errorf("Unable to create resource pool: %s", err)
		}

		conf := testCompute(validator, input, t)
		testTargets(validator, input, conf, t)
		testStorage(validator, input, conf, t)
	}
}

func getESXData(testURL *url.URL) *data.Data {
	result := data.NewData()
	testURL.Path = testURL.Path + "/ha-datacenter"
	result.OpsUser = testURL.User.Username()
	passwd, _ := testURL.User.Password()
	result.OpsPassword = &passwd
	result.URL = testURL
	result.DisplayName = "test001"
	result.ComputeResourcePath = "/ha-datacenter/host/localhost.localdomain/Resources"
	result.ImageDatastorePath = "LocalDS_0"
	result.BridgeNetworkName = "bridge"
	_, result.BridgeIPRange, _ = net.ParseCIDR("172.16.0.0/12")
	result.ManagementNetwork.Name = "VM Network"
	result.PublicNetwork.Name = "VM Network"
	result.VolumeLocations = make(map[string]*url.URL)
	testVolumeStoreURL := &url.URL{
		Host: "LocalDS_0",
		Path: "volumes/test",
	}
	result.VolumeLocations["volume-store"] = testVolumeStoreURL
	return result
}

func getVPXData(testURL *url.URL) *data.Data {
	result := data.NewData()
	testURL.Path = testURL.Path + "/DC0"
	result.OpsUser = testURL.User.Username()
	passwd, _ := testURL.User.Password()
	result.OpsPassword = &passwd
	result.URL = testURL
	result.DisplayName = "test001"
	result.ComputeResourcePath = "/DC0/host/DC0_C0/Resources"
	result.ImageDatastorePath = "LocalDS_0"
	result.PublicNetwork.Name = "VM Network"
	result.BridgeNetworkName = "bridge"
	_, result.BridgeIPRange, _ = net.ParseCIDR("172.16.0.0/12")

	result.VolumeLocations = make(map[string]*url.URL)
	testVolumeStoreURL := &url.URL{
		Host: "LocalDS_0",
		Path: "volumes/test",
	}
	result.VolumeLocations["volume-store"] = testVolumeStoreURL
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
		{"DC0_C0/validator", true, true},
		{"validator", true, false},
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
	v.credentials(v.Context, input, conf)

	u, err := url.Parse(conf.Target)
	assert.NoError(t, err)
	assert.Nil(t, u.User)
	assert.NotEmpty(t, conf.Token)
	assert.NotEmpty(t, conf.Username)
}

func testStorage(v *Validator, input *data.Data, conf *config.VirtualContainerHostConfigSpec, t *testing.T) {
	// specifically ignoring err here because we do not care about the parse result.
	testURL1, _ := url.Parse("LocalDS_0/volumes/volume1")
	testURL2, _ := url.Parse("ds://LocalDS_0/volumes/volume2")
	testURL3, _ := url.Parse("ds://LocalDS_0/volumes/volume1")

	// These two should report errors due to bad characters in the url. These should test how DatastoreHelper handles a nil or malformed url.
	testURL4, _ := url.Parse("😗/volumes/volume1")
	testURL5, _ := url.Parse("ds://😗/volumes/volume2")

	testURL6, _ := url.Parse("LocalDS_1/volumes/volume1")
	testURL7, _ := url.Parse("ds://LocalDS_1/volumes/volume2")
	testURL8, _ := url.Parse("")
	testURL9, _ := url.Parse("ds://")

	// positive nfs case
	nfsTestURL1, _ := url.Parse("nfs://prod.shared.storage/vchprod/volumes")

	// the two current negative nfs cases for validation
	nfsTestURL2, _ := url.Parse("nfs:///no/host/here")
	nfsTestURL3, _ := url.Parse("nfs://no.actual.path")

	tests := []struct {
		image         string
		volumes       map[string]*url.URL
		hasErr        bool
		expectImage   string
		expectVolumes map[string]*url.URL
	}{
		{"LocalDS_0",
			map[string]*url.URL{"volume1": testURL1,
				"volume2": testURL2},
			false,
			"ds://LocalDS_0/test001",
			map[string]*url.URL{"volume1": testURL3,
				"volume2": testURL2}},
		{"LocalDS_0/images",
			map[string]*url.URL{"volume1": testURL1,
				"volume2": testURL2},
			false,
			"ds://LocalDS_0/images",
			map[string]*url.URL{"volume1": testURL3,
				"volume2": testURL2}},

		{"ds://LocalDS_0/images",
			map[string]*url.URL{"volume1": testURL1,
				"volume2": testURL2},
			false,
			"ds://LocalDS_0/images",
			map[string]*url.URL{"volume1": testURL3,
				"volume2": testURL2}},

		{"ds://LocalDS_0/images/xyz",
			map[string]*url.URL{"volume1": testURL1,
				"volume2": testURL2},
			false,
			"ds://LocalDS_0/images/xyz",
			map[string]*url.URL{"volume1": testURL3,
				"volume2": testURL2}},

		{"ds://😗",
			map[string]*url.URL{"volume1": testURL4,
				"volume2": testURL5},
			true,
			"ds://😗/test001",
			nil},

		{"ds://LocalDS_0",
			map[string]*url.URL{"volume1": testURL6,
				"volume2": testURL7},
			true,
			"ds://LocalDS_0/test001",
			nil},

		{"LocalDS_0",
			map[string]*url.URL{"volume1": testURL6,
				"volume2": testURL7},
			true,
			"ds://LocalDS_0/test001",
			nil},

		{"LocalDS_0",
			map[string]*url.URL{"volume1": testURL6,
				"volume2": testURL7},
			true,
			"ds://LocalDS_0/test001",
			nil},

		{"",
			map[string]*url.URL{"volume1": testURL8,
				"volume2": testURL9},
			true,
			"",
			nil},

		{"ds://",
			map[string]*url.URL{"volume1": testURL8,
				"volume2": testURL9},
			true,
			"",
			nil},
		// below here lies the setup for nfs validation checks

		{"LocalDS_0",
			map[string]*url.URL{"volume1": nfsTestURL1},
			false,
			"ds://LocalDS_0/test001",
			map[string]*url.URL{"volume1": nfsTestURL1}},

		{"LocalDS_0",
			map[string]*url.URL{"volume1": nfsTestURL1,
				"volume2": nfsTestURL2},
			true,
			"ds://LocalDS_0/test001",
			map[string]*url.URL{"volume1": nfsTestURL1}},
		{"LocalDS_0",
			map[string]*url.URL{"volume1": nfsTestURL1,
				"volume2": nfsTestURL3},
			true,
			"ds://LocalDS_0/test001",
			map[string]*url.URL{"volume1": nfsTestURL1}},
		{"LocalDS_0",
			map[string]*url.URL{"volume1": nfsTestURL3,
				"volume2": nfsTestURL2},
			true,
			"ds://LocalDS_0/test001",
			nil},
		// below here lies the mixed store validation checks
		{"LocalDS_0",
			map[string]*url.URL{"volume1": testURL1,
				"volume2": nfsTestURL1,
				"volume3": nfsTestURL2,
				"volume4": testURL4,
			},
			true,
			"ds://LocalDS_0/test001",
			map[string]*url.URL{"volume1": testURL3,
				"volume2": nfsTestURL1}},
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
				if _, ok := test.expectVolumes[key]; !ok {
					assert.Fail(t, "Could not find volume store that was expected to present", "volume : %s", volume.String())
				} else {
					assert.Equal(t, test.expectVolumes[key].String(), volume.String())
				}
			}
		} else {
			assert.True(t, len(v.issues) > 0, "Should have errors")
		}
		v.issues = nil
		conf.VolumeLocations = nil
	}
}

func TestValidateWithFolders(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	ctx := context.Background()

	m := simulator.VPX()
	m.Datacenter = 3
	m.Folder = 2
	m.Datastore = 2
	m.ClusterHost = 3
	m.Pool = 1

	defer m.Remove()

	err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	m.Service.TLS = new(tls.Config)
	s := m.Service.NewServer()
	defer s.Close()

	input := data.NewData()
	input.URL = &url.URL{
		Scheme: s.URL.Scheme,
		Host:   s.URL.Host,
	}

	newShouldFail := true
	license := simulator.EvalLicense
	simulator.EvalLicense.Properties = nil // erase features

	var validator *Validator
	var dc string

	// Cover various failure paths while we're at it
	steps := []func(){
		func() {},
		func() {
			input.Thumbprint = "nope"
		},
		func() {
			input.Force = true
			input.Thumbprint = ""
		},
		func() {
			input.Force = false
			input.Thumbprint = s.CertificateInfo().ThumbprintSHA1
		},
		func() {
			input.URL.Path = "/"
			input.URL.User = s.URL.User
			newShouldFail = false
			if _, err = validator.ValidateCompute(ctx, input, false); err != nil {
				t.Error(err)
			}
		},
		func() {
			input.URL.Path = "/enoent" // Datacenter "enoent" in --target is not found
		},
		func() {
			input.URL.Path = "/DC1/sorry" // --target should only specify datacenter in the path
		},
		func() {
			input.URL.Path = "/DC1" // ok
			dc = input.URL.Path
		},
		func() {
			if _, err = validator.ValidateCompute(ctx, input, true); err == nil {
				t.Error("expected error")
			}
			input.ComputeResourcePath = "enoent"
		},
		func() {
			input.ComputeResourcePath = "DC1_C0"
		},
		func() {
			input.PublicNetwork.Name = "enoent"
		},
		func() {
			input.PublicNetwork.Name = "VM Network"
			input.ManagementNetwork.Name = input.PublicNetwork.Name
			input.ClientNetwork.Name = input.PublicNetwork.Name
			input.BridgeNetworkName = "DC1_DVPG0"
		},
		func() {
			input.ScratchSize = "10GB"
			p, _ := s.URL.User.Password()
			input.OpsPassword = &p
		},
		func() {
			input.ImageDatastorePath = "enoent"
		},
		func() {
			input.ImageDatastorePath = "LocalDS_*" // > 1
		},
		func() {
			input.ImageDatastorePath = "LocalDS_0"
		},
		func() {
			// TODO: volume
		},
		func() {
			input.OpsUser = s.URL.User.Username()
		},
		func() {
			simulator.EvalLicense.Properties = license.Properties // restore license features
		},
	}

	for i, step := range steps {
		if testing.Verbose() {
			fmt.Fprintf(os.Stderr, "TestValidateVPX(%d)%s\n", i, strings.Repeat(".", 30))
		}
		step()

		validator, err = NewValidator(ctx, input)
		if err != nil {
			continue
		}

		if newShouldFail {
			t.Fatalf("%d: expected error", i)
		}

		validator.DisableFirewallCheck = true

		_, err = validator.Validate(ctx, input)
		if i == len(steps)-1 {
			if err != nil {
				t.Fatal(err)
			}
		} else {
			if err == nil {
				t.Fatal("expected error")
			}
		}

		if dc != "" {
			// NewValidator has the side-effect of setting input.URL.Path=""
			input.URL.Path = dc
		}
	}

	// we have valid input at this point, test various compute-resource suggestions
	vs := validator.Session
	crs := []struct {
		flag    string
		pool    string
		cluster string
	}{
		{"*", "", ""},         // MultipleFoundError
		{"Resources", "", ""}, // MultipleFoundError
		{"DC1_[CH]0", "", ""}, // MultipleFoundError
		{"DC1_C0_RP1", "/F0/DC1/host/F0/DC1_C0/Resources/DC1_C0_RP1", "/F0/DC1/host/F0/DC1_C0"}, // ResourcePool (nested)
		{"DC1_H0", "/F0/DC1/host/F0/DC1_H0/Resources", "/F0/DC1/host/F0/DC1_H0"},                // Host (standalone)
		{"DC1_C0", "/F0/DC1/host/F0/DC1_C0/Resources", "/F0/DC1/host/F0/DC1_C0"},                // Cluster
	}

	for _, cr := range crs {
		vs.Pool = nil
		vs.PoolPath = ""

		vs.Cluster = nil
		vs.ClusterPath = ""

		_, err = validator.ResourcePoolHelper(ctx, cr.flag)

		if vs.ClusterPath != cr.cluster {
			t.Errorf("%s ClusterPath=%s", cr.flag, vs.ClusterPath)
		}

		if vs.PoolPath != cr.pool {
			t.Errorf("%s PoolPath=%s", cr.flag, vs.PoolPath)
		}

		if err == nil {
			continue
		}

		switch err.(type) {
		case *find.MultipleFoundError:
			// expected
		default:
			t.Errorf("ResourcePoolHelper(%s): %s", cr.flag, err)
		}
	}

	// cover some other paths now that we have a valid config
	spec, err := validator.ValidateTarget(ctx, input)
	if err != nil {
		t.Fatal(err)
	}

	validator.AddDeprecatedFields(ctx, spec, input)

	_, err = CreateFromVCHConfig(ctx, spec, vs)
	if err != nil {
		t.Fatal(err)
	}

	// force vim25.NewClient to fail
	simulator.Map.Remove(methods.ServiceInstance)
	validator.credentials(ctx, input, spec)

	// cover some of the return+error paths for certs (TODO: move this elsewhere and include valid data)
	validator.certificate(ctx, input, spec)
	validator.certificateAuthorities(ctx, input, spec)

	input.CertPEM = s.Certificate().Raw
	validator.certificate(ctx, input, spec)
	validator.certificateAuthorities(ctx, input, spec)

	input.ClientCAs = []byte{1}
	validator.certificateAuthorities(ctx, input, spec)

	validator.registries(ctx, input, spec)
	input.RegistryCAs = input.ClientCAs
	validator.registries(ctx, input, spec)
}

func TestValidateWithESX(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	ctx := context.Background()

	m := simulator.ESX()
	defer m.Remove()

	err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	m.Service.TLS = new(tls.Config)
	s := m.Service.NewServer()
	defer s.Close()

	input := data.NewData()
	input.URL = &url.URL{
		Path: s.URL.Host,
		User: s.URL.User,
	}

	input.Thumbprint = s.CertificateInfo().ThumbprintSHA1

	steps := []func(){
		func() {
			input.ComputeResourcePath = "enoent"
		},
		func() {
			input.PublicNetwork.Name = "enoent"
		},
		func() {
			input.ImageDatastorePath = "enoent"
		},
		func() {
			input.ImageDatastorePath = "enoent"
		},
		func() {
			input = getESXData(s.URL)
			input.URL.Path = "/"
			input.ScratchSize = "10GB"
			input.Force = true
		},
	}

	var validator *Validator

	for i, step := range steps {
		if testing.Verbose() {
			fmt.Fprintf(os.Stderr, "TestValidateESX(%d)%s\n", i, strings.Repeat(".", 30))
		}

		step()

		validator, err = NewValidator(ctx, input)
		if err != nil {
			t.Fatal(err)
		}
		validator.DisableFirewallCheck = true
		validator.AllowEmptyDC()

		_, err = validator.Validate(ctx, input)
		if i == len(steps)-1 {
			if err != nil {
				t.Fatal(err)
			}
		} else {
			if err == nil {
				t.Fatal("expected error")
			}
		}
	}

	// cover some errors paths by destroying this ESX system
	ref := esx.HostSystem.Reference()
	host := simulator.Map.Get(ref).(*simulator.HostSystem)

	steps = []func(){
		func() {
			host.Summary.ManagementServerIp = "owned"
		},
		func() {
			host.Summary.ManagementServerIp = ""
			simulator.Map.Remove(ref) // remove the host, forcing Finder.DefaultHostSystem to fail
		},
	}

	for i, step := range steps {
		step()

		validator.managedbyVC(ctx)
		issues := validator.GetIssues()

		if len(issues) != 1 {
			t.Errorf("%d issues: %s", i, issues)
		}
		validator.ClearIssues()
	}

	simulator.Map.Remove(esx.Datacenter.Reference()) // goodnight now.
	validator.suggestDatacenter()
}
