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
	"context"
	"fmt"
	"sort"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

var (
	cache = &vsanDSDomCache{
		ds: &object.Datastore{
			Common: object.NewCommon(nil, types.ManagedObjectReference{
				Type:  "test",
				Value: "1",
			}),
		},
		hvis: &testHostSystem{},
	}
)

type testHostSystem struct {
	deleteRes []types.HostVsanInternalSystemDeleteVsanObjectsResult
	queryRes  []string
	extAttrs  map[string]VSANExtAttrs
	err       error
}

func (h *testHostSystem) DeleteVsanObjects(ctx context.Context, uuids []string, force *bool) ([]types.HostVsanInternalSystemDeleteVsanObjectsResult, error) {
	return h.deleteRes, h.err
}

func (h *testHostSystem) QueryVsanObjectUuidsByFilter(ctx context.Context, uuids []string, limit int32, version int32) ([]string, error) {
	return h.queryRes, h.err
}

func (h *testHostSystem) GetVsanObjExtAttrs(ctx context.Context, uuids []string) (map[string]VSANExtAttrs, error) {
	return h.extAttrs, h.err
}

func TestTruncateFilePath(t *testing.T) {
	data := []struct {
		in  string
		out string
		err error
	}{
		{"/vmfs/volumes/vsan:521865c0389b0348-0dd6bfe38a2f193e/4eb8a458-8086-3055-d715-02000857aaba/VIC/4239f04c-bc6a-6653-988e-ad10e048d372/images/c40e708042c6a113e5827705cf9ace0377fae60cf9ae7a80ea8502d6460d87c6/c40e708042c6a113e5827705cf9ace0377fae60cf9ae7a80ea8502d6460d87c6.vmdk", "4eb8a458-8086-3055-d715-02000857aaba/VIC/4239f04c-bc6a-6653-988e-ad10e048d372/images/c40e708042c6a113e5827705cf9ace0377fae60cf9ae7a80ea8502d6460d87c6/c40e708042c6a113e5827705cf9ace0377fae60cf9ae7a80ea8502d6460d87c6.vmdk", nil},
		{"/vmfs/volumes/vsan:4239f04c-bc6a-6653-988e-ad10e048d372/", "", nil},
		{"/vmfs/volumes/vsan:521865c0389b0348-0dd6bfe38a2f193e/4eb8a458-8086-3055-d715-02000857aaba/VIC/4239f04c-bc6a-6653-988e-ad10e048d372/", "", nil},
		{"/vmfs/volumes/vsan:521865c0389b0348-0dd6bfe38a2f193e/4eb8a458-8086-3055-d715-02000857aaba/VIC/4239f04c-bc6a-6653-988e-ad10e048d372", "", nil},
		{"/vmfs/volumes/vsan:521865c0389b0348-0dd6bfe38a2f193e", "", nil},
		{"/vmfs/volumes/vsan:5218[65c0389b0348]-0dd6bfe38a2f193e", "", nil},
		{"[vsanDatastore] /VIC/images/4239f04c-bc6a-6653-988e-ad10e048d372/scratch.vmdk", "", errors.Errorf("not nil")},
	}
	for _, tt := range data {
		out, err := cache.truncateFilePath(tt.in)
		assert.Equal(t, tt.out, out, "wrong path")
		if tt.err != nil {
			assert.True(t, err != nil, "should have err returned")
			t.Logf("Got error: %s", err)
		}
	}
}

func TestRemoveFromCache(t *testing.T) {
	tests := []struct {
		uuids       map[string]string
		remove      []string
		expectUuids map[string]string
		expectPaths map[string]string
	}{
		{
			uuids:       map[string]string{},
			remove:      []string{"2", "3"},
			expectUuids: map[string]string{},
			expectPaths: map[string]string{},
		},
		{
			uuids: map[string]string{
				"1": "1.vmdk",
				"2": "2.vmdk",
				"3": "3.vmdk",
				"4": "4.vmdk",
				"5": "5.vmdk",
				"6": "6.vmdk",
			},
			remove: []string{"2", "3"},
			expectUuids: map[string]string{
				"1": "1.vmdk",
				"4": "4.vmdk",
				"5": "5.vmdk",
				"6": "6.vmdk",
			},
			expectPaths: map[string]string{
				"1.vmdk": "1",
				"4.vmdk": "4",
				"5.vmdk": "5",
				"6.vmdk": "6",
			},
		},
		{
			uuids: map[string]string{
				"1": "1.vmdk",
				"2": "2.vmdk",
				"3": "3.vmdk",
				"4": "4.vmdk",
				"5": "5.vmdk",
				"6": "6.vmdk",
			},
			remove: []string{"1", "8"},
			expectUuids: map[string]string{
				"2": "2.vmdk",
				"3": "3.vmdk",
				"4": "4.vmdk",
				"5": "5.vmdk",
				"6": "6.vmdk",
			},
			expectPaths: map[string]string{
				"2.vmdk": "2",
				"3.vmdk": "3",
				"4.vmdk": "4",
				"5.vmdk": "5",
				"6.vmdk": "6",
			},
		},
	}

	for _, data := range tests {
		cache.uuids = data.uuids
		cache.paths = make(map[string]string)
		// create reverse index
		for key, path := range cache.uuids {
			cache.paths[path] = key
		}
		cache.removeFromCache(data.remove)
		checkMaps(t, data.expectUuids, cache.uuids)
		checkMaps(t, data.expectPaths, cache.paths)
	}
}

func checkMaps(t *testing.T, expects map[string]string, actual map[string]string) {
	for k, v := range expects {
		val, ok := actual[k]
		assert.True(t, ok, fmt.Sprintf("key %s should exist in uuids", k))
		assert.Equal(t, v, val, "map value of key %s is inconsistent", k)
	}
	for k, v := range actual {
		val, ok := expects[k]
		assert.True(t, ok, fmt.Sprintf("key %s should be removed", k))
		assert.Equal(t, v, val, "map value of key %s is inconsistent", k)
	}
}

func TestRefresh(t *testing.T) {
	ctx := context.Background()
	thvic := cache.hvis.(*testHostSystem)
	cache.uuids = make(map[string]string)
	cache.paths = make(map[string]string)

	tests := []struct {
		queryUuids  []string
		extAttrs    map[string]VSANExtAttrs
		expectUuids map[string]string
		expectPaths map[string]string
	}{
		{
			queryUuids:  []string{},
			extAttrs:    map[string]VSANExtAttrs{},
			expectUuids: map[string]string{},
			expectPaths: map[string]string{},
		},
		{
			queryUuids: []string{"1", "2", "3", "4"},
			extAttrs: map[string]VSANExtAttrs{
				"1": VSANExtAttrs{"/vmfs/volumes/vsan:abc/VIC/1.vmdk"},
				"2": VSANExtAttrs{"/vmfs/volumes/vsan:abc/VIC/2.vmdk"},
				"3": VSANExtAttrs{"/vmfs/volumes/vsan:abc/VIC/3.vmdk"},
				"4": VSANExtAttrs{"/vmfs/volumes/vsan:abc/VIC/4.vmdk"},
			},
			expectUuids: map[string]string{
				"1": "VIC/1.vmdk",
				"2": "VIC/2.vmdk",
				"3": "VIC/3.vmdk",
				"4": "VIC/4.vmdk",
			},
			expectPaths: map[string]string{
				"VIC/1.vmdk": "1",
				"VIC/2.vmdk": "2",
				"VIC/3.vmdk": "3",
				"VIC/4.vmdk": "4",
			},
		},
		{
			queryUuids: []string{"1", "2", "5", "6"},
			extAttrs: map[string]VSANExtAttrs{
				"1": VSANExtAttrs{"/vmfs/volumes/vsan:abc/VIC/1.vmdk"},
				"2": VSANExtAttrs{"/vmfs/volumes/vsan:abc/VIC/2.vmdk"},
				"5": VSANExtAttrs{"/vmfs/volumes/vsan:abc/VIC/5.vmdk"},
				"6": VSANExtAttrs{"/vmfs/volumes/vsan:abc/VIC/6.vmdk"},
			},
			expectUuids: map[string]string{
				"1": "VIC/1.vmdk",
				"2": "VIC/2.vmdk",
				"5": "VIC/5.vmdk",
				"6": "VIC/6.vmdk",
			},
			expectPaths: map[string]string{
				"VIC/1.vmdk": "1",
				"VIC/2.vmdk": "2",
				"VIC/5.vmdk": "5",
				"VIC/6.vmdk": "6",
			},
		},
	}

	for _, data := range tests {
		thvic.queryRes = data.queryUuids
		thvic.extAttrs = data.extAttrs
		cache.Refresh(ctx)
		checkMaps(t, data.expectUuids, cache.uuids)
		checkMaps(t, data.expectPaths, cache.paths)
	}
}

func checkSlices(t *testing.T, exp []string, act []string) {
	assert.Equal(t, len(exp), len(act), "different length")
	sort.Strings(exp)
	sort.Strings(act)
	for i := range exp {
		assert.Equal(t, exp[i], act[i], fmt.Sprintf("array element %d is different", i))
	}
}

func TestDeleteVMDKDoms(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	trace.Logger.Level = log.DebugLevel
	ctx := context.Background()
	thvic := cache.hvis.(*testHostSystem)
	cache.uuids = map[string]string{
		"1": "1.vmdk",
		"2": "2.vmdk",
		"3": "3.vmdk",
		"4": "4.vmdk",
		"5": "5.vmdk",
		"6": "6.vmdk",
	}
	cache.paths = map[string]string{
		"1.vmdk": "1",
		"2.vmdk": "2",
		"3.vmdk": "3",
		"4.vmdk": "4",
		"5.vmdk": "5",
		"6.vmdk": "6",
	}

	tests := []struct {
		delRes      []types.HostVsanInternalSystemDeleteVsanObjectsResult
		in          []string
		out         []string
		failed      []string
		expectUuids map[string]string
		expectPaths map[string]string
	}{
		{
			expectUuids: map[string]string{
				"1": "1.vmdk",
				"2": "2.vmdk",
				"3": "3.vmdk",
				"4": "4.vmdk",
				"5": "5.vmdk",
				"6": "6.vmdk",
			},
			expectPaths: map[string]string{
				"1.vmdk": "1",
				"2.vmdk": "2",
				"3.vmdk": "3",
				"4.vmdk": "4",
				"5.vmdk": "5",
				"6.vmdk": "6",
			},
		},
		{
			in: []string{"5.vmdk", "2.vmdk"},
			delRes: []types.HostVsanInternalSystemDeleteVsanObjectsResult{
				{
					Uuid:    "2",
					Success: true,
				},
				{
					Uuid:    "5",
					Success: true,
				},
			},
			expectUuids: map[string]string{
				"1": "1.vmdk",
				"3": "3.vmdk",
				"4": "4.vmdk",
				"6": "6.vmdk",
			},
			expectPaths: map[string]string{
				"1.vmdk": "1",
				"3.vmdk": "3",
				"4.vmdk": "4",
				"6.vmdk": "6",
			},
		},
		{
			in: []string{"1.vmdk", "4.vmdk"},
			delRes: []types.HostVsanInternalSystemDeleteVsanObjectsResult{
				{
					Uuid:    "1",
					Success: true,
				},
				{
					Uuid:    "4",
					Success: false,
				},
			},
			expectUuids: map[string]string{
				"3": "3.vmdk",
				"4": "4.vmdk",
				"6": "6.vmdk",
			},
			expectPaths: map[string]string{
				"3.vmdk": "3",
				"4.vmdk": "4",
				"6.vmdk": "6",
			},
			failed: []string{"4"},
		},
		{
			in: []string{"7.vmdk", "4.vmdk"},
			delRes: []types.HostVsanInternalSystemDeleteVsanObjectsResult{
				{
					Uuid:    "4",
					Success: true,
				},
			},
			expectUuids: map[string]string{
				"3": "3.vmdk",
				"6": "6.vmdk",
			},
			expectPaths: map[string]string{
				"3.vmdk": "3",
				"6.vmdk": "6",
			},
			out: []string{"7.vmdk"},
		},
	}
	for _, data := range tests {
		thvic.deleteRes = data.delRes
		out, err := cache.DeleteVMDKDoms(ctx, data.in)
		checkSlices(t, data.out, out)
		if len(data.failed) > 0 {
			assert.True(t, err != nil, "Should get error for test %s", data)
			domErr := err.(DomDeleteError)
			checkSlices(t, data.failed, domErr.FailedUuids)
		} else {
			assert.True(t, err == nil, "Should not return error")
		}
		checkMaps(t, data.expectUuids, cache.uuids)
		checkMaps(t, data.expectPaths, cache.paths)
	}
}
