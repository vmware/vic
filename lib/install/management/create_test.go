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

package management

import (
	"net/url"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/simulator"

	"golang.org/x/net/context"
)

func TestMain(t *testing.T) {
	log.SetLevel(log.DebugLevel)
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
		installSettings := &data.InstallerData{}
		installSettings.ApplianceSize.CPU.Limit = 1
		installSettings.ApplianceSize.Memory.Limit = 1024

		validator, err := validate.NewValidator(ctx, input)
		if err != nil {
			t.Errorf("Failed to validator, %s", err)
		}

		validator.DisableFirewallCheck = true
		validator.DisableDRSCheck = true

		conf, err := validator.Validate(ctx, input)
		if err != nil {
			log.Errorf("Failed to validate conf, %s", err)
			validator.ListIssues()
		}

		// FIXME: get host NetworkSystem failed in simulator
		//		testCreateNetwork(ctx, validator.Session, conf, t)

		testCreateVolumeStores(ctx, validator.Session, conf, false, t)
		testDeleteVolumeStores(ctx, validator.Session, conf, 1, t)
		errConf := &metadata.VirtualContainerHostConfigSpec{}
		*errConf = *conf
		errConf.VolumeLocations = make(map[string]*url.URL)
		errConf.VolumeLocations["volume-store"], _ = url.Parse("ds://store_not_exist/volumes/test")
		testCreateVolumeStores(ctx, validator.Session, errConf, true, t)
		testCreateAppliance(ctx, validator.Session, conf, installSettings, false, t)
	}
}

func getESXData(url *url.URL) *data.Data {
	result := data.NewData()
	result.URL = url
	result.DisplayName = "test001"
	result.ComputeResourcePath = "/ha-datacenter/host/localhost.localdomain/Resources"
	result.ImageDatastoreName = "LocalDS_0"
	result.BridgeNetworkName = "bridge"
	result.ManagementNetworkName = "VM Network"
	result.ExternalNetworkName = "VM Network"
	result.VolumeLocations = make(map[string]string)
	result.VolumeLocations["volume-store"] = "LocalDS_0/volumes/test"

	return result
}

func getVPXData(url *url.URL) *data.Data {
	result := data.NewData()
	result.URL = url
	result.DisplayName = "test001"
	result.ComputeResourcePath = "/DC0/host/DC0_C0/Resources"
	result.ImageDatastoreName = "LocalDS_0"
	result.ExternalNetworkName = "VM Network"
	result.BridgeNetworkName = "bridge"
	result.VolumeLocations = make(map[string]string)
	result.VolumeLocations["volume-store"] = "LocalDS_0/volumes/test"

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
		ClusterPath:    "/DC0/host/DC0_C0",
	}

	s, err := session.NewSession(config).Connect(ctx)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func testCreateNetwork(ctx context.Context, sess *session.Session, conf *metadata.VirtualContainerHostConfigSpec, t *testing.T) {
	d := &Dispatcher{
		session: sess,
		ctx:     ctx,
		isVC:    sess.IsVC(),
		force:   false,
	}

	err := d.createBridgeNetwork(conf)
	if d.isVC && err != nil {
		t.Logf("Got exepcted err: %s", err)
		return
	}
	if d.isVC {
		t.Errorf("Should not create network in VC")
		return
	}
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func testCreateVolumeStores(ctx context.Context, sess *session.Session, conf *metadata.VirtualContainerHostConfigSpec, hasErr bool, t *testing.T) {
	d := &Dispatcher{
		session: sess,
		ctx:     ctx,
		isVC:    sess.IsVC(),
		force:   false,
	}

	err := d.createVolumeStores(conf)
	if hasErr && err != nil {
		t.Logf("Got exepcted err: %s", err)
		return
	}
	if hasErr {
		t.Errorf("Should have error, but got success")
		return
	}
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func testDeleteVolumeStores(ctx context.Context, sess *session.Session, conf *metadata.VirtualContainerHostConfigSpec, numVols int, t *testing.T) {
	d := &Dispatcher{
		session: sess,
		ctx:     ctx,
		isVC:    sess.IsVC(),
		force:   true,
	}

	if removed := d.deleteVolumeStoreIfForced(conf); removed != numVols {
		t.Errorf("Did not successfully remove all specified volumes")
	}

}

// FIXME: Failed to find IDE controller in simulator, so create appliance failed
func testCreateAppliance(ctx context.Context, sess *session.Session, conf *metadata.VirtualContainerHostConfigSpec, vConf *data.InstallerData, hasErr bool, t *testing.T) {
	d := &Dispatcher{
		session: sess,
		ctx:     ctx,
		isVC:    sess.IsVC(),
		force:   false,
	}
	delete(conf.Networks, "bridge") // FIXME: cannot create bridge network right now
	d.vchPool = d.session.Pool
	err := d.createAppliance(conf, vConf)
	if err != nil {
		t.Logf("Expected error: %s", err)
	}
}
