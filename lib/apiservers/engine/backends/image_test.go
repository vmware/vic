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

package backends

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/daemon/events"
	v1 "github.com/docker/docker/image"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware/vic/lib/imagec"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	"github.com/vmware/vic/lib/apiservers/engine/proxy/mocks"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/vsphere/sys"

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
		ImageID:   "test_id",
		Digests:   []string{fmt.Sprintf("%s@sha:%s", "test_name", "12345")},
		Tags:      []string{fmt.Sprintf("%s:%s", "test_name", "test_tag")},
		Name:      "test_name",
		DiffIDs:   map[string]string{"test_diffid": "test_layerid"},
		History:   []v1.History{},
		Reference: "test_name:test_tag",
	}

	dockerImage := convertV1ImageToDockerImage(image)

	assert.Equal(t, image.ImageID, dockerImage.ID, "Error: expected id %s, got %s", image.ImageID, dockerImage.ID)
	assert.Equal(t, image.Size, dockerImage.VirtualSize, "Error: expected size %s, got %s", image.Size, dockerImage.VirtualSize)
	assert.Equal(t, image.Size, dockerImage.Size, "Error: expected size %s, got %s", image.Size, dockerImage.Size)
	assert.Equal(t, image.Created.Unix(), dockerImage.Created, "Error: expected created %s, got %s", image.Created, dockerImage.Created)
	assert.Equal(t, image.Parent, dockerImage.ParentID, "Error: expected parent %s, got %s", image.Parent, dockerImage.ParentID)
	assert.Equal(t, image.Config.Labels, dockerImage.Labels, "Error: expected labels %s, got %s", image.Config.Labels, dockerImage.Labels)
	assert.Equal(t, image.Digests[0], dockerImage.RepoDigests[0], "Error: expected digest %s, got %s", image.Digests[0], dockerImage.RepoDigests[0])
	assert.Equal(t, image.Tags[0], dockerImage.RepoTags[0], "Error: expected tag %s, got %s", image.Tags[0], dockerImage.RepoTags[0])
}

func NewMockImageBackend(t *testing.T) (*gomock.Controller, *mocks.MockVicImageProxy, *Image) {
	// Get mock for the image proxy
	mockCtrl := gomock.NewController(t)
	mockProxy := mocks.NewMockVicImageProxy(mockCtrl)
	image := &Image{
		proxy: mockProxy,
	}
	return mockCtrl, mockProxy, image
}

func TestImageDelete(t *testing.T) {
	mockController, mockProxy, imageBackend := NewMockImageBackend(t)
	defer mockController.Finish()

	// Image ID
	imageID := "sha256:e27e9d7f7f28d67aa9e2d7540bdc2b33254b452ee8e60f388875e5b7d9b2b696"

	// Setup event service
	eventService = events.New()

	// Setup Image cache
	imageCache := cache.ImageCache()
	imageCache.SetEphemeral()
	image := &metadata.ImageConfig{
		V1Image: v1.V1Image{
			Config: &container.Config{},
		},
		ImageID: imageID,
		Name:    "my-image-name",
		Tags:    make([]string, 0),
		Digests: make([]string, 0),
	}
	imageCache.Add(image)

	image.Digests = []string{"<none>@<none>"}
	image.Tags = []string{"<none>:<none>"}

	// Setup Layer cache
	layerCache := imagec.LayerCache()
	layerCache.SetEphemeral()
	layer := &imagec.ImageWithMeta{
		Image: &models.Image{
			ID: "LayerOne",
		},
	}
	layerCache.Add(layer)
	layer = &imagec.ImageWithMeta{
		Image: &models.Image{
			ID: "LayerTwo",
		},
	}
	layerCache.Add(layer)

	// Setup Repo Cache
	repoCache := cache.RepositoryCache()
	repoCache.SetEphemeral()

	storeName, err := sys.UUID()
	require.Nil(t, err)
	keepNodes := make([]string, 1)
	imageURL, err := util.ImageURL(storeName, image.ImageID)
	require.Nil(t, err)
	keepNodes[0] = imageURL.String()

	// Build return value
	layers := []*models.Image{
		{
			ID: "LayerOne",
		},
		{
			ID: "LayerTwo",
		},
	}
	ret := &storage.DeleteImageOK{
		Payload: layers,
	}

	// Mock check parameters
	mockProxy.EXPECT().DeleteImage(gomock.Any(), gomock.Eq(storeName), gomock.Eq(image), keepNodes).Return(ret, nil)

	// Call the backend method
	imageBackend.ImageDelete(imageID, false, false)

	// Verify Image Cache
	img, err := imageCache.Get(imageID)
	require.Nil(t, img)
	errorString := fmt.Sprintf("No such image: %s", imageID)
	require.Equal(t, err.Error(), errorString)

	// Verify Layer Cache
	layer, err = layerCache.Get("LayerOne")
	require.Nil(t, layer)
	_, ok := err.(imagec.LayerNotFoundError)
	require.True(t, ok)
	layer, err = layerCache.Get("LayerTwo")
	_, ok = err.(imagec.LayerNotFoundError)
	require.True(t, ok)
}
