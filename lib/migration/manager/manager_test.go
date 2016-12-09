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

package manager

import (
	"context"
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
)

var testMap map[int]bool

type TestPlugin struct {
	ID int
}

func NewTestPlugin(id int) *TestPlugin {
	return &TestPlugin{id}
}

func (p *TestPlugin) Migrate(ctx context.Context, s *session.Session, data interface{}) (bool, error) {
	testMap[p.ID] = true
	return true, nil
}

func setUp() {
	log.SetLevel(log.DebugLevel)
	trace.Logger.Level = log.DebugLevel
	MaxPluginID = 100
	testMap = make(map[int]bool)
}

func TestInsertID(t *testing.T) {
	setUp()

	tester := &DataMigrator{
		targetIDs: make(map[string][]int),
		idPlugins: make(map[int]Plugin),
	}

	tester.insertID(1, ApplianceConfigure)
	tester.insertID(11, ApplianceConfigure)
	tester.insertID(9, ApplianceConfigure)
	assert.Equal(t, []int{1, 9, 11}, tester.targetIDs[ApplianceConfigure], "Should have expected array")
	tester.insertID(5, ApplianceConfigure)
	tester.insertID(8, ApplianceConfigure)
	tester.insertID(2, ApplianceConfigure)
	tester.insertID(4, ApplianceConfigure)
	assert.Equal(t, []int{1, 2, 4, 5, 8, 9, 11}, tester.targetIDs[ApplianceConfigure], "Should have expected array")
	tester.insertID(20, ApplianceConfigure)
	tester.insertID(15, ApplianceConfigure)
	assert.Equal(t, []int{1, 2, 4, 5, 8, 9, 11, 15, 20}, tester.targetIDs[ApplianceConfigure], "Should have expected array")
}

func TestMigratePluginExecution(t *testing.T) {
	setUp()

	tester := &DataMigrator{
		targetIDs: make(map[string][]int),
		idPlugins: make(map[int]Plugin),
	}

	ids := []int{1, 2, 4, 5, 8, 9, 11, 20, 15}
	var err error
	for _, id := range ids {
		if err = tester.Register(id, ApplianceConfigure, NewTestPlugin(id)); err != nil {
			t.Errorf("Failed to register plugin %d: %s", id, err)
		}
	}

	dataID, err := tester.Migrate(nil, nil, ApplianceConfigure, 0, nil)
	assert.Equal(t, 20, dataID, "migrated id mismatch")
	for _, id := range ids {
		assert.True(t, testMap[id], fmt.Sprintf("plugin %d should be executed", id))
	}
	testMap = make(map[int]bool)
	dataID, err = tester.Migrate(nil, nil, ApplianceConfigure, 5, nil)
	assert.Equal(t, 20, dataID, "migrated id mismatch")
	for _, id := range ids[:4] {
		assert.False(t, testMap[id], fmt.Sprintf("plugin %d should not be executed", id))
	}
	for _, id := range ids[4:] {
		assert.True(t, testMap[id], fmt.Sprintf("plugin %d should be executed", id))
	}

	testMap = make(map[int]bool)
	dataID, err = tester.Migrate(nil, nil, ApplianceConfigure, 20, nil)
	assert.Equal(t, 20, dataID, "migrated id mismatch")
	for _, id := range ids {
		assert.False(t, testMap[id], fmt.Sprintf("plugin %d should not be executed", id))
	}

	testMap = make(map[int]bool)
	dataID, err = tester.Migrate(nil, nil, ApplianceConfigure, 30, nil)
	assert.Equal(t, 30, dataID, "migrated id mismatch")
	for _, id := range ids {
		assert.False(t, testMap[id], fmt.Sprintf("plugin %d should not be executed", id))
	}
}
