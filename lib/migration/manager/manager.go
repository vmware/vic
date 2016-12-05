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
	"fmt"
	"sort"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/migration/errors"
	"github.com/vmware/vic/pkg/trace"
)

const (
	ApplianceConfigure = "ApplianceConfigure"
	ContainerConfigure = "ContainerConfigure"
	KeyValueStore      = "KeyValueStore"

	ConfigureVersionKey     = "guestinfo.vice.migration.version"
	KeyValueStoreVersionKey = "vice.migration.version"
)

var (
	// MaxPluginID must be increased to add new plugin and make sure the new plugin id is same to this value
	MaxPluginID = 2

	Migrator = NewDataMigrator()
)

type Plugin interface {
	Migrate(data interface{}) (bool, error)
}

type DataMigration interface {
	// Register plugin to data migration system
	Register(id int, target string, plugin Plugin) error
	// Migrate data with current version ID, return true if has any plugin executed
	Migrate(target string, currentID int, data interface{}) (int, error)
}

type DataMigrator struct {
	targetIDs map[string][]int
	idPlugins map[int]Plugin
}

func NewDataMigrator() DataMigration {
	return &DataMigrator{
		targetIDs: make(map[string][]int),
		idPlugins: make(map[int]Plugin),
	}
}

func (m *DataMigrator) Register(id int, target string, plugin Plugin) error {
	defer trace.End(trace.Begin(fmt.Sprintf("plugin %s:%d", target, id)))
	// assert if plugin id less than mast plugin id, which is forcing deveoper to change MaxPluginID variable everytime new plugin is added
	if plugin == nil {
		return &errors.InternalError{
			fmt.Sprintf("Empty Plugin object is not allowed"),
		}
	}
	if id > MaxPluginID {
		return &errors.InternalError{
			fmt.Sprintf("Plugin %d is bigger than Max Plugin ID %d", id, MaxPluginID),
		}
	}

	if m.idPlugins[id] != nil {
		return &errors.InternalError{
			fmt.Sprintf("Plugin %d is conflict with another plugin, please make sure the plugin ID is unique and ascending", id),
		}
	}

	m.insertID(id, target)
	m.idPlugins[id] = plugin
	return nil
}

func (m *DataMigrator) insertID(id int, target string) {
	defer trace.End(trace.Begin(fmt.Sprintf("id array: %s, insert %s:%d", m.targetIDs[target], target, id)))

	s := m.targetIDs[target]
	if len(s) == 0 {
		m.targetIDs[target] = append(s, id)
		return
	}
	i := sort.SearchInts(s, id)
	m.targetIDs[target] = append(s[:i], append([]int{id}, s[i:]...)...)
	log.Debugf("id array: %d", m.targetIDs[target])
}

func (m *DataMigrator) Migrate(target string, currentID int, data interface{}) (int, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("migrate %s from %d", target, currentID)))

	pluginIDs := m.targetIDs[target]
	if len(pluginIDs) == 0 {
		log.Debugf("No plugin registered for %s", target)
		return currentID, nil
	}

	i := sort.SearchInts(pluginIDs, currentID)
	if i >= len(pluginIDs) {
		log.Debugf("No plugin bigger than %d", currentID)
		return currentID, nil
	}

	latestID := currentID
	j := i
	if pluginIDs[i] == currentID {
		j = i + 1
	}
	for ; j < len(pluginIDs); j++ {
		id := pluginIDs[j]
		p := m.idPlugins[id]
		c, err := p.Migrate(data)
		if err != nil {
			return latestID, err
		}
		if c {
			latestID = id
		}
	}
	return latestID, nil
}
