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
	"context"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/migration/errors"
	"github.com/vmware/vic/lib/migration/manager"
	"github.com/vmware/vic/pkg/kvstore"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
)

// Sample plugin to migrate data in keyvalue store
// Migrate keyvalue plugin should read configuration from input VirtualContainerHost configuration, and then read from keyvalue store file directly
// After migration, write back to datastore file
// Note: Currently allowed change in keyvalue store is only "ADD". Old value could be migrated to new key, but old key could not be removed. This is to make sure if eventually
// VCH upgrade failed, no datastore file revert available in current migration framework, so this limitation can make sure old binary still works well with half migrated keyvalue store.
// While there is rollback plugin and framework support, this limitation could be removed.
const (
	id     = 2
	target = manager.ApplianceConfigure

	KVStoreFolder = "kvStores"
	APIKV         = "apiKV"

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

func (p *NewImageMeta) Migrate(ctx context.Context, s *session.Session, data interface{}) (bool, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("%d", id)))
	if data == nil {
		return false, nil
	}
	vchConfMap := data.(map[string]string)
	// No plugin query keyvalue store yet, load from datastore file
	// get a ds helper for this ds url
	vchConf := &VirtualContainerHostConfigSpec{}
	extraconfig.Decode(extraconfig.MapSource(vchConfMap), vchConf)

	imageURL := vchConf.ImageStores[0]
	// TODO: sample code, should get datastore from imageURL
	dsHelper, err := datastore.NewHelper(trace.NewOperation(ctx, "datastore helper creation"), s,
		s.Datastore, fmt.Sprintf("%s/%s", imageURL.Path, KVStoreFolder))
	if err != nil {
		return false, &errors.InternalError{
			fmt.Sprintf("unable to get datastore helper for %s store creation: %s", APIKV, err.Error()),
		}
	}

	// restore the modified K/V store
	keyValStore, err := kvstore.NewKeyValueStore(ctx, kvstore.NewDatastoreBackend(dsHelper), APIKV)
	if err != nil && !os.IsExist(err) {
		return false, &errors.InternalError{
			fmt.Sprintf("unable to create %s datastore backed store: %s", APIKV, err.Error()),
		}
	}
	val, err := keyValStore.Get(oldKey)
	if err != nil {
		return false, &errors.InternalError{
			fmt.Sprintf("failed to get %s from store %s: %s", oldKey, APIKV, err.Error()),
		}
	}
	if val == nil {
		log.Debugf("Nothing to migrate")
		return false, nil
	}
	// put the new key/value to store, and leave the old key/value there, in case upgrade failed, old binary still works well with half-changed store
	keyValStore.Put(ctx, newKey, []byte(fmt.Sprintf("%s:%s", val, "latest")))
	// persist new data back to vsphere, framework does not take of it
	keyValStore.Save(ctx)
	return false, nil
}
