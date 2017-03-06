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

package plugin1

import (
	"context"
	"fmt"

	"github.com/vmware/vic/lib/migration/manager"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/lib/migration/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
)

const (
	version = 1
	target  = manager.ContainerConfigure
)

func init() {
	defer trace.End(trace.Begin(fmt.Sprintf("Registering plugin %s:%d", target, version)))
	if err := manager.Migrator.Register(version, target, &AddCommonSpecForVM{}); err != nil {
		log.Errorf("Failed to register plugin %s:%d, %s", target, version, err)
	}
}

// AddCommonSpecForVM is plugin for vic 0.8.0-GA version upgrade
type AddCommonSpecForVM struct {
}

type ExecutorConfig struct {
	Common `vic:"0.1" scope:"read-only" key:"common"`
}

type Common struct {
	// A reference to the components hosting execution environment, if any
	ExecutionEnvironment string

	// Unambiguous ID with meaning in the context of its hosting execution environment
	ID string `vic:"0.1" scope:"read-only" key:"id"`

	// Convenience field to record a human readable name
	Name string `vic:"0.1" scope:"read-only" key:"name"`

	// Freeform notes related to the entity
	Notes string `vic:"0.1" scope:"hidden" key:"notes"`
}

type NewExecutorConfig struct {
	CommonSpecForVM `vic:"0.1" scope:"read-only" key:"common"`
}

type CommonSpecForVM struct {
	// A reference to the components hosting execution environment, if any
	ExecutionEnvironment string

	// Unambiguous ID with meaning in the context of its hosting execution environment
	ID string `vic:"0.1" scope:"read-only" key:"id"`

	// Convenience field to record a human readable name
	Name string `vic:"0.1" scope:"hidden" key:"name"`

	// Freeform notes related to the entity
	Notes string `vic:"0.1" scope:"hidden" key:"notes"`
}

func (p *AddCommonSpecForVM) Migrate(ctx context.Context, s *session.Session, data interface{}) error {
	defer trace.End(trace.Begin(fmt.Sprintf("%d", version)))
	if data == nil {
		return nil
	}
	mapData := data.(map[string]string)
	oldStruct := &ExecutorConfig{}
	result := extraconfig.Decode(extraconfig.MapSource(mapData), oldStruct)
	log.Infof("The mapData is %+v", mapData)
	log.Infof("The oldStruct is %+v", oldStruct)
	if result == nil {
		return &errors.DecodeError{}
	}

	newStruct := &NewExecutorConfig{
		CommonSpecForVM: CommonSpecForVM{
			Name: oldStruct.Name,
		}}
	log.Infof("The newStruct is %+v", newStruct)
	cfg := make(map[string]string)
	extraconfig.Encode(extraconfig.MapSink(cfg), newStruct)
	log.Infof("The newStruct is %+v", newStruct)

	for k, v := range cfg {
		log.Debugf("New data: %s:%s", k, v)
		mapData[k] = v
	}
	return nil
}
