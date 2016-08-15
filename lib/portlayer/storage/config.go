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

import "net/url"

var Config Configuration

// Configuration is a slice of the VCH config that is relevant to the exec part of the port layer
type Configuration struct {
	// Turn on debug logging
	DebugLevel int `vic:"0.1" scope:"read-only" key:"init/common/debug"`

	////////////// Port Layer - storage
	// Datastore URLs for image stores - the top layer is [0], the bottom layer is [len-1]
	ImageStores []url.URL `vic:"0.1" scope:"read-only" key:"image_stores"`

	// Permitted datastore URL roots for volumes
	// Keyed by the volume store name (which is used by the docker user to
	// refer to the datstore + path), valued by the datastores and the path.
	VolumeLocations map[string]url.URL `vic:"0.1" scope:"read-only"`

	ScratchSize int64 `vic:"0.1" scope:"read-only" key:"scratch_size"`
}
