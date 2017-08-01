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

package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBindSrcTarget(t *testing.T) {
	dir, err := ioutil.TempDir("", "testCreateBindSrcTarget")
	assert.NoError(t, err)

	defer os.RemoveAll(dir)

	// Create a file under an existing directory
	file := dir + "/test1"
	err = createBindSrcTarget(map[string]os.FileMode{file: 0644})
	assert.NoError(t, err, "createBindSrcTarget failed: %s", err)
	_, err = os.Stat(dir)
	assert.NoError(t, err)

	// Create a file under non-existent directory
	file = dir + "/testDir/test1"
	err = createBindSrcTarget(map[string]os.FileMode{file: 0644})
	assert.NoError(t, err, "createBindSrcTarget failed: %s", err)
	_, err = os.Stat(dir)
	assert.NoError(t, err)

	// Create an existing file
	err = createBindSrcTarget(map[string]os.FileMode{file: 0644})
	assert.NoError(t, err, "createBindSrcTarget failed: %s", err)
}
