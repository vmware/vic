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

package exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerCache(t *testing.T) {
	NewContainerCache()
	containerID := "1234"
	id := ParseID(containerID)
	// create a new container
	NewContainer(id)
	// should only have 1 in cache
	assert.Equal(t, len(containers.cache), 1)
	container := containers.Container(containerID)
	// do we have this one in the cache?
	assert.Equal(t, container.ExecConfig.ID, containerID)
	// remove the container
	containers.Remove(containerID)
	assert.Equal(t, len(containers.cache), 0)
}
