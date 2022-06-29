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

package plugina

import (
	"context"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/migration/errors"
	"github.com/vmware/vic/lib/migration/feature"
	"github.com/vmware/vic/lib/migration/manager"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
)

const (
	ctarget = manager.ContainerConfigure
)

func init() {
	defer trace.End(trace.Begin(fmt.Sprintf("Registering plugins %s:%d", ctarget, feature.ContainerGuestVisibleName)))

	if err := manager.Migrator.Register(feature.ContainerGuestVisibleName, ctarget, &ContainerGuestVisibleName{}); err != nil {
		log.Errorf("Failed to register plugin %s:%d, %s", ctarget, feature.ContainerGuestVisibleName, err)
		panic(err)
	}
}

// ContainerGuestVisibleName is plugin for vic 1.5.6 version upgrade
type ContainerGuestVisibleName struct {
}

// Migrate deals with the container config and is distinct from the Migrate which is for VCH config.
func (p *ContainerGuestVisibleName) Migrate(ctx context.Context, s *session.Session, data interface{}) error {
	defer trace.End(trace.Begin(fmt.Sprintf("ContainerGuestVisibleName version %d", feature.ContainerGuestVisibleName)))
	if data == nil {
		return nil
	}
	mapData, ok := data.(map[string]string)
	if !ok {
		// Log the error here and return nil so that other plugins can proceed
		log.Errorf("Migration data format is not map: %+v", data)
		return nil
	}
	oldStruct := &ExecutorConfig{}
	result := extraconfig.Decode(extraconfig.MapSource(mapData), oldStruct)
	log.Debugf("The oldStruct is %+v", oldStruct)
	if result == nil {
		return &errors.DecodeError{Err: fmt.Errorf("decode oldStruct %+v failed", oldStruct)}
	}

	newStruct := &UpdatedExecutorConfig{
		UpdatedCommon: UpdatedCommon{
			Name:                 oldStruct.Name,
			ID:                   oldStruct.ID,
			Notes:                oldStruct.Notes,
			ExecutionEnvironment: oldStruct.ExecutionEnvironment,
		}}

	cfg := make(map[string]string)
	extraconfig.Encode(extraconfig.MapSink(cfg), newStruct)

	for k, v := range cfg {
		log.Debugf("New data: %s:%s", k, v)
		mapData[k] = v
	}
	return nil
}
