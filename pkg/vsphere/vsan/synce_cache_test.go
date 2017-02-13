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

// This is a workaround for vsan DOM object leaking issue caused by FileManager.DeleteDatastoreFile, see github issue #3787 and bugzilla issue #1808703
// This file used draft vSphere API, which is subject to change in the future, so this workaround should be removed as soon as the DOM leaking issue is fixed by vSAN.
package vsan

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/trace"
)

type testDSCache struct {
	refreshed int32
	deleted   int32
}

func (tc *testDSCache) Refresh() error {
	atomic.AddInt32(&tc.refreshed, 1)
	return nil
}

func (tc *testDSCache) DeleteVMDKDoms(paths []string) ([]string, error) {
	atomic.AddInt32(&tc.deleted, 1)
	return []string{"test"}, nil
}

func (tc *testDSCache) CleanOrphanDoms() ([]string, error) {
	return nil, nil
}

func setUp(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	trace.Logger.Level = log.DebugLevel
}

func TestSync(t *testing.T) {
	setUp(t)

	refreshInterval = 50 * time.Millisecond
	tc := testDSCache{}
	SyncedDomCache.dsMap = map[string]DSDomCache{"test-1:": &tc}
	//	SyncedDomCache.dsMap["test-1:"] = &tc
	var wg sync.WaitGroup
	wg.Add(1)
	t.Log("Start refresh")
	SyncedDomCache.Refresh("test-1:", &wg)
	wg.Wait()
	assert.Equal(t, int32(1), atomic.LoadInt32(&tc.refreshed), "Should be refreshed once")
	<-time.After(60 * time.Millisecond)
	assert.True(t, atomic.LoadInt32(&tc.refreshed) >= 2, "Should be refreshed two times")
	wg.Add(1)
	t.Log("Trigger refresh again")
	cr := atomic.LoadInt32(&tc.refreshed)
	SyncedDomCache.Refresh("test-1:", &wg)
	SyncedDomCache.waitRefresh("test-1:")
	wg.Wait()
	assert.Equal(t, cr+1, atomic.LoadInt32(&tc.refreshed), "Should be only one more time")

	SyncedDomCache.waitRefresh("test-1:")
	assert.Equal(t, cr+1, atomic.LoadInt32(&tc.refreshed), "Should be no more refresh")

	ds := &object.Datastore{
		Common: object.NewCommon(nil, types.ManagedObjectReference{
			Type: "test-1",
		}),
	}
	SyncedDomCache.deleteVMDKDoms(ds, nil)
	assert.Equal(t, int32(1), atomic.LoadInt32(&tc.deleted), "Delete should be called once")
	d := atomic.LoadInt32(&tc.deleted)
	r := atomic.LoadInt32(&tc.refreshed)
	SyncedDomCache.SyncDeleteVMDKDoms(ds, nil, true)
	assert.Equal(t, r+1, atomic.LoadInt32(&tc.refreshed), "Refreshed once more")
	assert.Equal(t, d+2, atomic.LoadInt32(&tc.deleted), "Delete should be called twice more")

	d = atomic.LoadInt32(&tc.deleted)
	r = atomic.LoadInt32(&tc.refreshed)
	SyncedDomCache.SyncDeleteVMDKDoms(ds, nil, false)
	assert.Equal(t, r, atomic.LoadInt32(&tc.refreshed), "Refreshed once more")
	assert.Equal(t, d+2, atomic.LoadInt32(&tc.deleted), "Delete should be called twice more")
}
