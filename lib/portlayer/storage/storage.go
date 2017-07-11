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

package storage

import (
	"context"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
)

var once sync.Once

func create(ctx context.Context, session *session.Session, pool *object.ResourcePool) error {
	var err error

	mngr := view.NewManager(session.Vim25())

	// Create view of VirtualMachine objects under the VCH's resource pool
	Config.ContainerView, err = mngr.CreateContainerView(ctx, pool.Reference(), []string{"VirtualMachine"}, true)
	if err != nil {
		return err
	}
	return nil
}

func Init(ctx context.Context, session *session.Session, pool *object.ResourcePool, source extraconfig.DataSource, _ extraconfig.DataSink) error {
	defer trace.End(trace.Begin(""))

	var err error

	once.Do(func() {
		// Grab the storage layer config blobs from extra config
		extraconfig.Decode(source, &Config)
		log.Debugf("Decoded VCH config for storage: %#v", Config)

		err = create(ctx, session, pool)
	})
	return err
}
