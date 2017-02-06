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
	"context"
	"strconv"

	"github.com/vmware/vic/lib/migration/errors"
	"github.com/vmware/vic/lib/migration/manager"
	_ "github.com/vmware/vic/lib/migration/plugins"
	"github.com/vmware/vic/pkg/vsphere/session"
)

// MigrateApplianceConfigure migrate VCH appliance configuration, including guestinfo, keyvaluestore, or any other configuration related change.
//
// Note: Input map conf is VCH appliance guestinfo map, and returned map is the new guestinfo. Any other changes should be made in plugin.
// If there is error returned, returned map might have half-migrated value
func MigrateApplianceConfig(ctx context.Context, s *session.Session, conf map[string]string) (map[string]string, bool, error) {
	return migrateConfig(ctx, s, conf, manager.ApplianceConfigure, manager.ApplianceVersionKey)
}

// MigrateContainerConfigure migrate container configuration
//
// Note: Migrated data will be returned in map, and input object is not changed.
// If there is error returned, returned map might have half-migrated value.
func MigrateContainerConfig(conf map[string]string) (map[string]string, bool, error) {
	return migrateConfig(nil, nil, conf, manager.ContainerConfigure, manager.ContainerVersionKey)
}

// IsContainerDataOlder returns true if input container config is older than latest version. If error happens, always returns false
func IsContainerDataOlder(conf map[string]string) (bool, error) {
	return isDataOlder(conf, manager.ContainerConfigure, manager.ContainerVersionKey)
}

// IsApplianceDataOlder returns true if input appliance config is older than latest version. If error happens, always returns false
func IsApplianceDataOlder(conf map[string]string) (bool, error) {
	return isDataOlder(conf, manager.ApplianceConfigure, manager.ApplianceVersionKey)
}

// isDataOlder returns true if data is older than latest. If error happens, always returns false
func isDataOlder(data map[string]string, target string, verKey string) (bool, error) {
	var currentID int
	var err error

	if currentID, err = getCurrentID(data, verKey); err != nil {
		return false, err
	}
	latestVer := manager.Migrator.LatestVersion(target)
	return latestVer > currentID, nil
}

func migrateConfig(ctx context.Context, s *session.Session, data map[string]string, target string, verKey string) (map[string]string, bool, error) {
	dst := make(map[string]string)
	if len(data) == 0 {
		return dst, false, nil
	}

	var currentID int
	var err error

	if currentID, err = getCurrentID(data, verKey); err != nil {
		return dst, false, err
	}

	for k, v := range data {
		dst[k] = v
	}

	latestID, err := manager.Migrator.Migrate(ctx, s, target, currentID, dst)
	if latestID == currentID {
		return dst, false, err
	}
	dst[verKey] = strconv.Itoa(latestID)
	return dst, true, err
}

func getCurrentID(data map[string]string, verKey string) (int, error) {
	var currentID int
	var err error
	strID := data[verKey]

	if strID == "" {
		return 0, nil
	}
	if currentID, err = strconv.Atoi(strID); err != nil {
		return 0, &errors.InvalidMigrationVersion{
			Version: strID,
			Err:     err,
		}
	}
	return currentID, nil
}
