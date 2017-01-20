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

package manager

import (
	"context"
	"fmt"
	"sort"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/migration/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/session"
)

const (
	ApplianceConfigure = "ApplianceConfigure"
	ContainerConfigure = "ContainerConfigure"

	ApplianceVersionKey = "guestinfo.vice./init/version/PluginVersion"
	ContainerVersionKey = "guestinfo.vice./version/PluginVersion"
)

var (
	Migrator = NewDataMigrator()
)

type Plugin interface {
	Migrate(ctx context.Context, s *session.Session, data interface{}) error
}

type DataMigration interface {
	// Register plugin to data migration system
	Register(version int, target string, plugin Plugin) error
	// Migrate data with current version ID, return true if has any plugin executed
	Migrate(ctx context.Context, s *session.Session, target string, currentVersion int, data interface{}) (int, error)
	// GetLatestVersion return the latest plugin id for specified target
	GetLatestVersion(target string) int
}

type DataMigrator struct {
	targetVers map[string][]int
	verPlugins map[int]Plugin
}

func NewDataMigrator() DataMigration {
	return &DataMigrator{
		targetVers: make(map[string][]int),
		verPlugins: make(map[int]Plugin),
	}
}

func (m *DataMigrator) Register(ver int, target string, plugin Plugin) error {
	defer trace.End(trace.Begin(fmt.Sprintf("plugin %s:%d", target, ver)))
	// assert if plugin version less than max plugin version, which is forcing deveoper to change MaxPluginVersion variable everytime new plugin is added
	if plugin == nil {
		return &errors.InternalError{
			Message: "Empty Plugin object is not allowed",
		}
	}
	if ver > version.MaxPluginVersion {
		return &errors.InternalError{
			Message: fmt.Sprintf("Plugin %d is bigger than Max Plugin Version %d", ver, version.MaxPluginVersion),
		}
	}

	if m.verPlugins[ver] != nil {
		return &errors.InternalError{
			Message: fmt.Sprintf("Plugin %d is conflict with another plugin, please make sure the plugin Version is unique and ascending", ver),
		}
	}

	m.insertVersion(ver, target)
	m.verPlugins[ver] = plugin
	return nil
}

func (m *DataMigrator) insertVersion(version int, target string) {
	defer trace.End(trace.Begin(fmt.Sprintf("insert %s:%d", target, version)))

	s := m.targetVers[target]
	if len(s) == 0 {
		m.targetVers[target] = append(s, version)
		return
	}
	i := sort.SearchInts(s, version)
	m.targetVers[target] = append(s[:i], append([]int{version}, s[i:]...)...)
	log.Debugf("version array: %d", m.targetVers[target])
}

func (m *DataMigrator) Migrate(ctx context.Context, s *session.Session, target string, currentVersion int, data interface{}) (int, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("migrate %s from %d", target, currentVersion)))

	pluginVers := m.targetVers[target]
	if len(pluginVers) == 0 {
		log.Debugf("No plugin registered for %s", target)
		return currentVersion, nil
	}

	i := sort.SearchInts(pluginVers, currentVersion)
	if i >= len(pluginVers) {
		log.Debugf("No plugin bigger than %d", currentVersion)
		return currentVersion, nil
	}

	latestVer := currentVersion
	j := i
	if pluginVers[i] == currentVersion {
		j = i + 1
	}
	for ; j < len(pluginVers); j++ {
		ver := pluginVers[j]
		p := m.verPlugins[ver]
		err := p.Migrate(ctx, s, data)
		if err != nil {
			return latestVer, err
		}
		latestVer = ver
	}
	return latestVer, nil
}

func (m *DataMigrator) GetLatestVersion(target string) int {
	pluginVers := m.targetVers[target]
	l := len(pluginVers)
	if l == 0 {
		log.Debugf("No plugin registered for %s", target)
		return 0
	}
	return pluginVers[l-1]
}
