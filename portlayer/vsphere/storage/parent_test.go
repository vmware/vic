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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/test"
	"golang.org/x/net/context"
)

func parentSetup(t *testing.T) *session.Session {
	datastoreParentPath = "testingParentDirectory"

	return test.Session(context.TODO(), t)
}

func TestParentEmptyRestore(t *testing.T) {
	client := parentSetup(t)
	if client == nil {
		return
	}

	par, err := restoreParentMap(context.TODO(), client)
	if !assert.NoError(t, err) && !assert.NotNil(t, par) {
		return
	}
}

func TestParentEmptySaveRestore(t *testing.T) {
	client := parentSetup(t)
	if client == nil {
		return
	}
	// Nuke the parent image store directory
	defer rm(t, client, "")

	par, err := restoreParentMap(context.TODO(), client)
	if !assert.NoError(t, err) && !assert.NotNil(t, par) {
		return
	}

	err = par.Save(context.TODO())
	if !assert.NoError(t, err) {
		return
	}

	p, err := restoreParentMap(context.TODO(), client)
	if !assert.NoError(t, err) && !assert.NotNil(t, p) {
		return
	}
}

// Write some child -> parent mappings and see if we can read them.
func TestParentSaveRestore(t *testing.T) {
	client := parentSetup(t)
	if client == nil {
		return
	}
	// Nuke the parent image store directory
	defer rm(t, client, "")

	par, err := restoreParentMap(context.TODO(), client)
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
	err = par.Save(context.TODO())
	if !assert.NoError(t, err) {
		return
	}

	// load into a different map
	p, err := restoreParentMap(context.TODO(), client)
	if !assert.NoError(t, err) && !assert.NotNil(t, p) {
		return
	}

	// check if the 2nd map loaded everything correctly
	if !assert.Equal(t, expected, p.db) {
		return
	}

	// Now save it to be extra paranoid
	err = p.Save(context.TODO())
	if !assert.NoError(t, err) {
		return
	}
}
