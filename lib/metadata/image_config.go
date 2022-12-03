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

package metadata

import (
	docker "github.com/docker/docker/image"
)

// ImageConfig defines the docker format for representing an image. When marshaled to JSON, the sha256 sum
// of the resulting bytes is the image ID.
type ImageConfig struct {
	docker.V1Image

	// image specific data
	ImageID   string            `json:"image_id"`
	Digests   []string          `json:"digests,omitempty"`
	Tags      []string          `json:"tags,omitempty"`
	Name      string            `json:"name,omitempty"`
	DiffIDs   map[string]string `json:"diff_ids,omitempty"`
	History   []docker.History  `json:"history,omitempty"`
	Reference string            `json:"registry"`

	// VMDK is the ID of the VMDK to be used as the R/W layer's parent disk when
	// creating a container from the image. This is the type that is stored in our
	// image cache.
	//
	// This field is ignored when marshalling into JSON to preserve the image ID.
	VMDK string `json:"vmdk,omitempty"`
}
