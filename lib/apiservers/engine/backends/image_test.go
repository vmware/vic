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
	"fmt"
	"testing"
	"time"

	v1 "github.com/docker/docker/image"
	"github.com/docker/engine-api/types/container"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/metadata"
)

func TestConvertV1ImageToDockerImage(t *testing.T) {
	now := time.Now()

	image := &metadata.ImageConfig{
		V1Image: v1.V1Image{
			ID:      "deadbeef",
			Size:    1024,
			Created: now,
			Parent:  "",
			Config: &container.Config{
				Labels: map[string]string{},
			},
		},
		ImageID: "test_id",
		Digests: []string{"12345"},
		Tags:    []string{"test_tag"},
		Name:    "test_name",
		DiffIDs: map[string]string{"test_diffid": "test_layerid"},
		History: []v1.History{},
	}
	digest := fmt.Sprintf("%s@sha:%s", image.Name, "12345")
	tag := fmt.Sprintf("%s:%s", image.Name, "test_tag")

	dockerImage := convertV1ImageToDockerImage(image)

	assert.Equal(t, image.ImageID, dockerImage.ID, "Error: expected id %s, got %s", image.ImageID, dockerImage.ID)
	assert.Equal(t, image.Size, dockerImage.VirtualSize, "Error: expected size %s, got %s", image.Size, dockerImage.VirtualSize)
	assert.Equal(t, image.Size, dockerImage.Size, "Error: expected size %s, got %s", image.Size, dockerImage.Size)
	assert.Equal(t, image.Created.Unix(), dockerImage.Created, "Error: expected created %s, got %s", image.Created, dockerImage.Created)
	assert.Equal(t, image.Parent, dockerImage.ParentID, "Error: expected parent %s, got %s", image.Parent, dockerImage.ParentID)
	assert.Equal(t, image.Config.Labels, dockerImage.Labels, "Error: expected labels %s, got %s", image.Config.Labels, dockerImage.Labels)
	assert.Equal(t, digest, dockerImage.RepoDigests[0], "Error: expected digest %s, got %s", digest, dockerImage.RepoDigests[0])
	assert.Equal(t, tag, dockerImage.RepoTags[0], "Error: expected tag %s, got %s", tag, dockerImage.RepoTags[0])
}

func TestClientFriendlyTags(t *testing.T) {
	imageName := "busybox"
	tags := []string{"1.24.2", "latest"}

	friendlyTags := clientFriendlyTags(imageName, tags)
	assert.Equal(t, len(friendlyTags), len(tags), "Error: expected %d tags, got %d", len(tags), len(friendlyTags))
	assert.Equal(t, friendlyTags[0], "busybox:1.24.2", "Error: expected %s, got %s", "busybox:1.24.2", friendlyTags[0])
	assert.Equal(t, friendlyTags[1], "busybox:latest", "Error: expected %s, got %s", "busybox:latest", friendlyTags[1])

	emptyTags := clientFriendlyTags(imageName, []string{})
	assert.Equal(t, len(emptyTags), 1, "Error: expected %d tags, got %d", 1, len(emptyTags))
	assert.Equal(t, emptyTags[0], "<none>:<none>", "Error: expected %s tags, got %s", "<none>:<none>", emptyTags[0])

}

func TestClientFriendlyDigests(t *testing.T) {
	imageName := "busybox"
	digests := []string{"12345", "6789"}

	friendlyDigests := clientFriendlyDigests(imageName, digests)
	assert.Equal(t, len(friendlyDigests), len(digests), "Error: expected %d digests, got %d", len(digests), len(friendlyDigests))
	assert.Equal(t, friendlyDigests[0], "busybox@sha:12345", "Error: expected %s, got %s", "busybox@sha:12345", friendlyDigests[0])
	assert.Equal(t, friendlyDigests[1], "busybox@sha:6789", "Error: expected %s, got %s", "busybox@sha:6789", friendlyDigests[1])

	emptyDigests := clientFriendlyDigests(imageName, []string{})
	assert.Equal(t, len(emptyDigests), 1, "Error: expected %d digests, got %d", 1, len(emptyDigests))
	assert.Equal(t, emptyDigests[0], "<none>@<none>", "Error: expected %s digests, got %s", "<none>@<none>", emptyDigests[0])

}
