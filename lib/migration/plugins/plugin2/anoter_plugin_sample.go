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

package plugin2

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/migration/errors"
	"github.com/vmware/vic/lib/migration/manager"
	"github.com/vmware/vic/pkg/trace"
)

// Sample plugin to migrate data in keyvalue store
const (
	id     = 2
	target = manager.KeyValueStore

	oldKey = "image.name"
	newKey = "image.tag"
)

func init() {
	log.Debugf("Registering plugin %s:%d", target, id)
	if err := manager.Migrator.Register(id, target, &NewImageMeta{}); err != nil {
		log.Errorf("Failed to register plugin %s:%d", target, id, err)
	}
}

// NewImageMeta is plugin for vic 0.8.0-GA version upgrade
type NewImageMeta struct {
}

func (p *NewImageMeta) Migrate(data interface{}) (bool, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("%d", id)))
	mapData, ok := data.(map[string][]byte)
	if !ok {
		return false, &errors.DataTypeError{
			"map[string][]byte",
		}
	}

	val := mapData[oldKey]
	if val != nil {
		mapData[newKey] = []byte(fmt.Sprintf("%s:%s", val, "latest"))
		delete(mapData, oldKey)
		return true, nil
	}
	log.Debugf("Nothing to migrate")
	return false, nil
}
