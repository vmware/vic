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

package cache

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"

	derr "github.com/docker/docker/errors"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/metadata"
)

// ImageCache is an in-memory cache of image metadata. It is refreshed at startup
// by a call to the portlayer. It is updated when new images are pulled or
// images are deleted.
type ImageCache struct {
	m sync.RWMutex

	// cache maps image ID to image metadata
	cache map[string]*metadata.ImageConfig
}

func NewImageCache() *ImageCache {
	return &ImageCache{
		cache: make(map[string]*metadata.ImageConfig),
	}
}

// Update adds new layers to the cache
func (c *ImageCache) Update(client *client.PortLayer) error {
	c.m.Lock()
	defer c.m.Unlock()

	log.Debugf("Updating image cache...")

	host, err := guest.UUID()
	if host == "" {
		host, err = os.Hostname()
	}
	if err != nil {
		return fmt.Errorf("Unexpected error getting hostname: %s", err)
	}

	// attempt to create the image store if it doesn't exist
	store := &models.ImageStore{Name: host}
	_, err = client.Storage.CreateImageStore(
		storage.NewCreateImageStoreParams().WithBody(store),
	)

	if err != nil {
		if _, ok := err.(*storage.CreateImageStoreConflict); ok {
			log.Debugf("Store already exists")
		} else {
			// TODO(jzt): add some retry/backoff logic for this?
			// Or just let the watchdog handle restarting the docker engine api?
			log.Debugf("Creating a store failed: %#v", err)
			return err
		}
	}

	params := storage.NewListImagesParams().WithStoreName(host)

	layers, err := client.Storage.ListImages(params)
	if err != nil {
		return fmt.Errorf("Failed to retrieve image list from portlayer: %s", err)
	}

	for _, layer := range layers.Payload {
		imageConfig := &metadata.ImageConfig{}
		if err := json.Unmarshal([]byte(layer.Metadata["metaData"]), imageConfig); err != nil {
			derr.NewErrorWithStatusCode(fmt.Errorf("Failed to unmarshal image config: %s", err),
				http.StatusInternalServerError)
		}

		if imageConfig.ImageID != "" {
			c.cache[imageConfig.ImageID] = imageConfig
		}
	}

	return nil
}

// GetImages returns a slice containing metadata for all cached images
func (c *ImageCache) GetImages() []*metadata.ImageConfig {
	c.m.RLock()
	defer c.m.RUnlock()

	result := make([]*metadata.ImageConfig, 0, len(c.cache))
	for _, image := range c.cache {
		newImage := new(metadata.ImageConfig)
		*newImage = *image
		result = append(result, newImage)
	}

	return result
}
