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

package vicbackends

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessVolumeParams(t *testing.T) {
	rawTestVolumes := []string{"/blah", "testVolume:/mount", "testVolume:/mount/path:r"}
	var processedTestVolumes []volumeFields

	for _, testString := range rawTestVolumes {
		processedFields, err := processVolumeParam(testString)
		assert.Nil(t, err)
		processedTestVolumes = append(processedTestVolumes, processedFields)
	}
	assert.Equal(t, 3, len(processedTestVolumes))

	assert.NotEmpty(t, processedTestVolumes[0].ID)
	assert.Equal(t, "/blah", processedTestVolumes[0].Dest)
	assert.Equal(t, "rw", processedTestVolumes[0].Flags)

	assert.Equal(t, "testVolume", processedTestVolumes[1].ID)
	assert.Equal(t, "/mount", processedTestVolumes[1].Dest)
	assert.Equal(t, "rw", processedTestVolumes[1].Flags)

	assert.Equal(t, "testVolume", processedTestVolumes[2].ID)
	assert.Equal(t, "/mount/path", processedTestVolumes[2].Dest)
	assert.Equal(t, "r", processedTestVolumes[2].Flags)
}

func TestProcessSpecifiedVolumes(t *testing.T) {
	rawTestVolumes := []string{"masterVolume:/blah", "testVolume:/mount:r", "specVol:/mount/path:r"}
	var processedTestVolumes []volumeFields

	processedFields, err := processSpecifiedVolumes(rawTestVolumes)
	assert.Nil(t, err)
	processedTestVolumes = append(processedTestVolumes, processedFields...)

	assert.Len(t, processedFields, 3)

	assert.Equal(t, "masterVolume", processedTestVolumes[0].ID)
	assert.Equal(t, "/blah", processedTestVolumes[0].Dest)
	assert.Equal(t, "rw", processedTestVolumes[0].Flags)

	assert.Equal(t, "testVolume", processedTestVolumes[1].ID)
	assert.Equal(t, "/mount", processedTestVolumes[1].Dest)
	assert.Equal(t, "r", processedTestVolumes[1].Flags)

	assert.Equal(t, "testVolume", processedTestVolumes[2].ID)
	assert.Equal(t, "/mount/path", processedTestVolumes[2].Dest)
	assert.Equal(t, "r", processedTestVolumes[2].Flags)
}
