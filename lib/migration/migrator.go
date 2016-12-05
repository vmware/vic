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
	"strconv"

	"github.com/vmware/vic/lib/migration/errors"
	"github.com/vmware/vic/lib/migration/manager"
	_ "github.com/vmware/vic/lib/migration/plugins"
)

// MigrateApplianceConfigure migrate VCH appliance configuration.
// Migrated data will be returned in map, and input object is not changed.
// If there is error returned, returned map might have half-migrated value.
func MigrateApplianceConfigure(data map[string]string) (map[string]string, bool, error) {
	return migrateConfigure(data, manager.ApplianceConfigure)
}

// MigrateContainerConfigure migrate container configuration.
// Migrated data will be returned in map, and input object is not changed.
// If there is error returned, returned map might have half-migrated value.
func MigrateContainerConfigure(data map[string]string) (map[string]string, bool, error) {
	return migrateConfigure(data, manager.ContainerConfigure)
}

func migrateConfigure(data map[string]string, target string) (map[string]string, bool, error) {
	dst := make(map[string]string)
	if len(data) == 0 {
		return dst, false, nil
	}

	var currentID int
	var err error
	strID := data[manager.ConfigureVersionKey]

	if strID == "" {
		currentID = 0
	} else {
		if currentID, err = strconv.Atoi(strID); err != nil {
			return dst, false, &errors.InvalidMigrationID{
				ID:  strID,
				Err: err,
			}
		}
	}

	for k, v := range data {
		dst[k] = v
	}

	latestID, err := manager.Migrator.Migrate(target, currentID, dst)
	if latestID == currentID {
		return dst, false, err
	}
	dst[manager.ConfigureVersionKey] = strconv.Itoa(latestID)
	return dst, true, err
}

// MigrateKeyValueStore migrate vic keyvalue store.
// Migrated data will be returned in map, and input object is not changed.
// If there is error returned, returned map might have half-migrated value.
func MigrateKeyValueStore(data map[string][]byte) (map[string][]byte, bool, error) {
	dst := make(map[string][]byte)
	if len(data) == 0 {
		return dst, false, nil
	}

	var currentID int
	var err error
	byteID := data[manager.KeyValueStoreVersionKey]

	if len(byteID) == 0 {
		currentID = 0
	} else {
		strID := string(byteID)
		if currentID, err = strconv.Atoi(strID); err != nil {
			return dst, false, &errors.InvalidMigrationID{
				ID:  strID,
				Err: err,
			}
		}
	}

	for k, v := range data {
		dst[k] = v
	}

	latestID, err := manager.Migrator.Migrate(manager.KeyValueStore, currentID, dst)
	if latestID == currentID {
		return dst, false, err
	}
	dst[manager.KeyValueStoreVersionKey] = []byte(strconv.Itoa(latestID))
	return dst, true, err
}
