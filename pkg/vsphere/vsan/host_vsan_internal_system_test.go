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

package vsan

import (
	"context"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/test/env"
)

func Session(ctx context.Context, t *testing.T) *session.Session {
	config := &session.Config{
		Service: env.URL(t),

		DatastorePath: env.DS(t),

		Insecure:  true,
		Keepalive: time.Duration(5) * time.Minute,
	}

	s, err := session.NewSession(config).Connect(context.Background())
	s.Populate(context.Background())
	if err != nil {
		t.Error(err)
	}

	return s
}

// TODO: This test should be enabled with vcSim support
func testDelete(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.Background()

	s := Session(context.Background(), t)
	if s == nil {
		t.Errorf("Client is not created")
	}

	t.Logf("session: %#v", s.Vim25())

	//	testVHIS(t, ctx, s)
	testCache(ctx, s, t)
}

func testCache(ctx context.Context, s *session.Session, t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	SyncedDomCache.AddDomCache(ctx, s.Datastore)
	t.Logf("Wait cache refresh")
	SyncedDomCache.waitRefresh(s.Datastore.Reference().String())
	t.Logf("Cache refresh finished")
	dsCache := SyncedDomCache.dsMap[s.Datastore.Reference().String()].(*vsanDSDomCache)
	t.Logf("uuid cache: %#v", dsCache.uuids)
	t.Logf("path cache: %#v", dsCache.paths)

	SyncedDomCache.SyncCleanOrphanDoms(ctx, s.Datastore, false)
	dsCache = SyncedDomCache.dsMap[s.Datastore.Reference().String()].(*vsanDSDomCache)
	t.Logf("uuid cache: %#v", dsCache.uuids)
	t.Logf("path cache: %#v", dsCache.paths)

	SyncedDomCache.SyncCleanOrphanDoms(ctx, s.Datastore, true)
	dsCache = SyncedDomCache.dsMap[s.Datastore.Reference().String()].(*vsanDSDomCache)
	t.Logf("uuid cache: %#v", dsCache.uuids)
	t.Logf("path cache: %#v", dsCache.paths)

	SyncedDomCache.SyncDeleteVMDKDoms(ctx, s.Datastore, []string{"ef13a258-af82-6c3a-739f-020014517637/volumes/volume1/volume1.vmdk"}, true)
	t.Logf("uuid cache: %#v", dsCache.uuids)
	t.Logf("path cache: %#v", dsCache.paths)
}
