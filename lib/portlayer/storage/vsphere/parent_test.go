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

package vsphere

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/pkg/vsphere/datastore"
)

const testStore = "testStore"

func TestParentEmptyRestore(t *testing.T) {
	ctx, ds, cleanupfunc := datastore.DSsetup(t)
	if t.Failed() {
		return
	}
	defer cleanupfunc()

	par, err := restoreParentMap(ctx, ds, testStore)
	if !assert.NoError(t, err) && !assert.NotNil(t, par) {
		return
	}
}

func TestParentEmptySaveRestore(t *testing.T) {
	ctx, ds, cleanupfunc := datastore.DSsetup(t)
	if t.Failed() {
		return
	}
	defer cleanupfunc()

	par, err := restoreParentMap(ctx, ds, testStore)
	if !assert.NoError(t, err) && !assert.NotNil(t, par) {
		return
	}

	err = par.Save(ctx)
	if !assert.NoError(t, err) {
		return
	}

	p, err := restoreParentMap(ctx, ds, testStore)
	if !assert.NoError(t, err) && !assert.NotNil(t, p) {
		return
	}
}

// Write some child -> parent mappings and see if we can read them.
func TestParentSaveRestore(t *testing.T) {
	ctx, ds, cleanupfunc := datastore.DSsetup(t)
	if t.Failed() {
		return
	}
	defer cleanupfunc()

	par, err := restoreParentMap(ctx, ds, testStore)
	if !assert.NoError(t, err) && !assert.NotNil(t, par) {
		return
	}

	expected := make(map[string]string)
	for i := 0; i < 10; i++ {
		child := fmt.Sprintf("c%d", i)
		parent := fmt.Sprintf("p%d", i)
		expected[child] = parent
		par.Add(child, parent)
	}
	err = par.Save(ctx)
	if !assert.NoError(t, err) {
		return
	}

	// load into a different map
	p, err := restoreParentMap(ctx, ds, testStore)
	if !assert.NoError(t, err) && !assert.NotNil(t, p) {
		return
	}

	// check if the 2nd map loaded everything correctly
	if !assert.Equal(t, expected, p.db) {
		return
	}

	// Now save it to be extra paranoid
	err = p.Save(ctx)
	if !assert.NoError(t, err) {
		return
	}
}
