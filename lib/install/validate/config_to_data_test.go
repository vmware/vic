// Copyright 2017 VMware, Inc. All Rights Reserved.
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
	"fmt"
	"net/url"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

type testFinder struct {
}

func setUp() {

}

func TestTest(t *testing.T) {
	urls := make([]url.URL, 1)
	u, _ := url.Parse("ds://vsanDatastore/path")
	urls[0] = *u
	t.Logf("urls[0]: %s", urls[0].String())
	t.Logf("urls[0] path: %s", urls[0].Path)
	url := urls[0]
	url.Path = ""
	t.Logf("urls[0]: %s", urls[0].String())
	t.Logf("url: %s", url.String())
}

func TestConvert(t *testing.T) {
	//	ConverterLogLevel = log.DebugLevel
	//	trace.Logger.Level = log.DebugLevel

	ctx := context.Background()
	target := "administrator@vsphere.local:Admin!23@10.192.171.116"
	sessionconfig := &session.Config{
		Service:        target,
		Insecure:       true,
		DatacenterPath: "/vcqaDC",
		ClusterPath:    "/vcqaDC/host/cls",
		DatastorePath:  "/vcqaDC/datastore/vsanDatastore",
		PoolPath:       "/vcqaDC/host/cls/Resources",
	}

	//	u, err := url.Parse(sessionconfig.Service)
	//	if err != nil {
	//		t.Errorf("Failed to parse url %s: %s", sessionconfig.Service, err)
	//	} else {
	//		t.Logf("Got url: %s", u)
	//	}
	installSession := session.NewSession(sessionconfig)

	_, err := installSession.Connect(ctx)
	if err != nil {
		log.Errorf("Failed to connect session: %s", err)
	}
	if _, err = installSession.Populate(ctx); err != nil {
		log.Errorf("Failed to get resources: %s", err)
	}
	log.Infof("cluster: %s", installSession.Cluster)
	log.Infof("datacenter: %s", installSession.Datacenter)
	log.Infof("datastore: %s", installSession.Datastore)
	log.Infof("host: %s", installSession.Host)
	log.Infof("pool: %s", installSession.Pool)

	test1, err := installSession.Finder.VirtualMachine(ctx, "test1/test1")
	if err != nil {
		log.Infof("Failed to find VM: %s", err)
	} else {
		log.Infof("vm path: %s", test1.InventoryPath)
	}

	vm := vm.NewVirtualMachine(ctx, installSession, test1.Reference())
	//this is the appliance vm
	mapConfig, err := vm.FetchExtraConfigBaseOptions(ctx)
	if err != nil {
		err = fmt.Errorf("Failed to get VM extra config of %q: %s", vm.Reference(), err)
		log.Error(err)
	}

	kv := vmomi.OptionValueMap(mapConfig)

	conf := &config.VirtualContainerHostConfigSpec{}
	extraconfig.DecodeWithPrefix(extraconfig.MapSource(kv), conf, "")
	t.Logf("Command: %s", conf.VicMachineCreateOptions)
	t.Logf("networks: %#v", conf.Network.ContainerNetworks)

	data, err := NewDataFromConfig(ctx, installSession.Finder, conf)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("bridge ip range: %s", data.BridgeIPRange.String())
	t.Logf("dest: %#v", data)
}
