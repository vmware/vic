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

package migration

import (
	"strconv"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/migration/manager"
	"github.com/vmware/vic/lib/migration/samples/plugins/plugin1"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

func setUp() {
	// register sample plugin into test
	log.SetLevel(log.DebugLevel)
	trace.Logger.Level = log.DebugLevel
	version.MaxPluginVersion = version.MaxPluginVersion + 1

	if err := manager.Migrator.Register(1, manager.ApplianceConfigure, &plugin1.ApplianceStopSignalRename{}); err != nil {
		log.Errorf("Failed to register plugin %s:%d, %s", manager.ApplianceConfigure, version.MaxPluginVersion, err)
	}
}

func TestMigrateConfigure(t *testing.T) {
	setUp()

	conf := &config.VirtualContainerHostConfigSpec{
		ExecutorConfig: executor.ExecutorConfig{
			Sessions: map[string]*executor.SessionConfig{
				"abc": &executor.SessionConfig{
					Attach:     true,
					StopSignal: "2",
				},
				"def": &executor.SessionConfig{
					Attach:     false,
					StopSignal: "10",
				},
			},
		},
		Network: config.Network{
			BridgeNetwork: "VM Network",
		},
	}
	mapData := make(map[string]string)
	extraconfig.Encode(extraconfig.MapSink(mapData), conf)
	t.Logf("Old data: %#v", mapData)
	newData, migrated, err := MigrateApplianceConfig(nil, nil, mapData)
	if err != nil {
		t.Errorf("migration failed: %s", err)
		assert.Fail(t, "migration failed")
	}
	assert.True(t, migrated, "should be migrated")

	latestVer := newData[manager.ApplianceVersionKey]
	assert.Equal(t, strconv.Itoa(version.MaxPluginVersion), latestVer, "upgrade version mismatch")

	// check new data
	var found bool
	for k, _ := range newData {
		if strings.Contains(k, "stopSignal") {
			assert.Fail(t, "key %s still exists in migrated data", k)
		}
		if strings.Contains(k, "forceStopSignal") {
			found = true
		}
	}
	assert.True(t, found, "Should found migrated data")

	// verify old data
	found = false
	for k, _ := range mapData {
		if strings.Contains(k, "stopSignal") {
			found = true
		}
		if strings.Contains(k, "forceStopSignal") {
			assert.Fail(t, "key %s is found in old data", k)
		}
	}
	assert.True(t, found, "Should found old data")

	t.Logf("New data: %#v", newData)
	newConf := &config.VirtualContainerHostConfigSpec{}
	extraconfig.Decode(extraconfig.MapSource(newData), newConf)

	assert.Equal(t, version.MaxPluginVersion, newConf.Version.PluginVersion, "should not be migrated")
	t.Logf("other version fields: %s", newConf.Version.String())
}

func TestIsDataOlder(t *testing.T) {
	setUp()

	conf := &config.VirtualContainerHostConfigSpec{
		ExecutorConfig: executor.ExecutorConfig{
			Sessions: map[string]*executor.SessionConfig{
				"abc": &executor.SessionConfig{
					Attach:     true,
					StopSignal: "2",
				},
				"def": &executor.SessionConfig{
					Attach:     false,
					StopSignal: "10",
				},
			},
		},
		Network: config.Network{
			BridgeNetwork: "VM Network",
		},
	}
	mapData := make(map[string]string)
	extraconfig.Encode(extraconfig.MapSink(mapData), conf)
	t.Logf("Old appliance data: %#v", mapData)
	older, err := ApplianceDataIsOlder(mapData)
	assert.Equal(t, nil, err, "should not have error")
	assert.True(t, older, "Test data should be older than latest")

	mapData = make(map[string]string)
	extraconfig.Encode(extraconfig.MapSink(mapData), conf.ExecutorConfig)
	t.Logf("Old container data: %#v", mapData)

	older, err = ContainerDataIsOlder(mapData)
	assert.Equal(t, nil, err, "should not have error")
	assert.True(t, older, "Test data should be older than latest since a container update plugin has been registered")
}
