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
	"time"

	v1 "github.com/docker/docker/image"
	"github.com/docker/engine-api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/apiservers/portlayer/models"
)

func TestGetV1ImageMapEmpty(t *testing.T) {
	images := []*models.Image{
		&models.Image{
			ID: "scratch",
		},
	}
	v1Map := getV1MapFromImages(images)
	assert.Empty(t, v1Map, "v1Map should be empty")
}

func TestGetV1ImageMap(t *testing.T) {
	p := "scratch"
	s := "test_selflink"
	images := []*models.Image{
		&models.Image{
			ID: "deadbeef",
			Metadata: map[string]string{
				"v1Compatibility": "{\"id\": \"deadbeef\", \"Size\":1024, \"parent\":\"\"}",
			},
			Parent:   &p,
			SelfLink: &s,
			Store:    "teststore",
		},
	}
	v1Map := getV1MapFromImages(images)
	assert.NotEmpty(t, v1Map, "v1Map should not be empty")
	if !assert.Equal(t, 1, len(v1Map), "len(v1Map) should be 1") {
		return
	}
}

func TestResolveImageSize(t *testing.T) {
	p := ""
	pp := "deadbeef"
	s := "test_selflink"
	images := []*models.Image{
		&models.Image{
			ID:       "scratch",
			Metadata: map[string]string{},
			Parent:   nil,
			SelfLink: &s,
			Store:    "teststore",
		},
		&models.Image{
			ID: "deadbeef",
			Metadata: map[string]string{
				"v1Compatibility": "{\"id\": \"deadbeef\", \"Size\":1024, \"parent\":\"\"}",
			},
			Parent:   &p,
			SelfLink: &s,
			Store:    "teststore",
		},
		&models.Image{
			ID: "test_id",
			Metadata: map[string]string{
				"v1Compatibility": "{\"id\": \"test_id\", \"Size\":1024, \"parent\":\"deadbeef\"}",
			},
			Parent:   &pp,
			SelfLink: &s,
			Store:    "teststore",
		},
	}
	v1Map := getV1MapFromImages(images)
	assert.NotEmpty(t, v1Map, "v1Map should not be empty")
	// minus the scratch image which gets skipped
	if !assert.Equal(t, len(images)-1, len(v1Map), "len(v1Map) should be 1") {
		return
	}

	// "test_id" image size should be the size of its parent ("deadbeef") plus its own size before resolution
	expected := v1Map["deadbeef"].Size + v1Map["test_id"].Size

	// perform the size resolution
	resolveImageSizes(v1Map)

	actual := v1Map["test_id"].Size
	assert.Equal(t, expected, actual, "Error: expected size %d, got %d", expected, actual)
}

func TestConvertV1ImageToDockerImage(t *testing.T) {
	now := time.Now()

	v1Image := &v1.V1Image{
		ID:      "deadbeef",
		Size:    1024,
		Created: now,
		Parent:  "",
		Config: &container.Config{
			Labels: map[string]string{},
		},
	}
	dockerImage := convertV1ImageToDockerImage(v1Image)
	assert.Equal(t, v1Image.ID, dockerImage.ID, "Error: expected id %s, got %s", v1Image.ID, dockerImage.ID)
	assert.Equal(t, v1Image.Size, dockerImage.VirtualSize, "Error: expected size %s, got %s", v1Image.Size, dockerImage.VirtualSize)
	assert.Equal(t, v1Image.Created.Unix(), dockerImage.Created, "Error: expected created %s, got %s", v1Image.Created, dockerImage.Created)
	assert.Equal(t, v1Image.Parent, dockerImage.ParentID, "Error: expected parent %s, got %s", v1Image.Parent, dockerImage.ParentID)
	assert.Equal(t, v1Image.Config.Labels, dockerImage.Labels, "Error: expected labels %s, got %s", v1Image.Config.Labels, dockerImage.Labels)
}
