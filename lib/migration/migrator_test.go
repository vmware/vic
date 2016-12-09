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

package migration

import (
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/migration/manager"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

func TestMigrateConfigure(t *testing.T) {
	log.SetLevel(log.DebugLevel)

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
	newData, migrated, err := MigrateApplianceConfigure(nil, nil, mapData)
	if err != nil {
		t.Errorf("migration failed: %s", err)
		assert.Fail(t, "migration failed")
	}
	assert.True(t, migrated, "should be migrated")

	latestID := newData[manager.ConfigureVersionKey]
	assert.Equal(t, "1", latestID, "upgrade version mismatch")

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
}
