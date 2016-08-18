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
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/pkg/uid"
)

func TestContainerCache(t *testing.T) {
	NewContainerCache()
	containerID := "1234"
	id := uid.Parse(containerID)

	// create a new container
	NewContainer(id)
	// shouldn't have a container as it's not commited
	assert.Equal(t, len(containers.cache), 0)

	// create a new container
	container := &Container{ExecConfig: &executor.ExecutorConfig{}}
	container.ExecConfig.ID = containerID
	// put it in the cache
	containers.Put(container)
	// Get it
	cachedContainer := containers.Container(containerID)
	// did we find it?
	assert.NotNil(t, cachedContainer)
	// do we have this one in the cache?
	assert.Equal(t, cachedContainer.ExecConfig.ID, containerID)
	// remove the container
	containers.Remove(containerID)
	assert.Equal(t, len(containers.cache), 0)
	// remove non-existent container
	containers.Remove("blahblah")
}
