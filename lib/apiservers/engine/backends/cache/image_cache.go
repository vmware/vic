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
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/distribution/digest"
	derr "github.com/docker/docker/errors"
	"github.com/docker/docker/reference"

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
	cacheByID   map[string]*metadata.ImageConfig
	cacheByName map[string]*metadata.ImageConfig
}

// NewImageCache creates and returns a new ImageCache
func NewImageCache() *ImageCache {
	return &ImageCache{
		cacheByID:   make(map[string]*metadata.ImageConfig),
		cacheByName: make(map[string]*metadata.ImageConfig),
	}
}

// Update adds new images to the cache
func (c *ImageCache) Update(client *client.PortLayer) error {

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
			c.AddImage(imageConfig)
		}
	}

	return nil
}

// GetImages returns a slice containing metadata for all cached images
func (c *ImageCache) GetImages() []*metadata.ImageConfig {
	c.m.RLock()
	defer c.m.RUnlock()

	result := make([]*metadata.ImageConfig, 0, len(c.cacheByID))
	for _, image := range c.cacheByID {
		newImage := new(metadata.ImageConfig)
		*newImage = *image
		result = append(result, newImage)
	}

	return result
}

func (c *ImageCache) GetImageByDigest(digest digest.Digest) *metadata.ImageConfig {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cacheByID[string(digest)]
}

// Looks up image by reference.Named
func (c *ImageCache) GetImageByNamed(named reference.Named) *metadata.ImageConfig {
	c.m.RLock()
	defer c.m.RUnlock()

	var config *metadata.ImageConfig

	if tagged, ok := named.(reference.NamedTagged); ok {
		taggedName := tagged.Name() + ":" + tagged.Tag()
		config = c.cacheByName[taggedName]
	} else {
		// First try just the name.
		config = c.cacheByName[named.Name()]
		if config == nil {
			// try with the default docker tag
			taggedName := named.Name() + ":" + reference.DefaultTag
			config = c.cacheByName[taggedName]
		}
	}

	return config
}

// AddImage adds an image to the image cache
func (c *ImageCache) AddImage(imageConfig *metadata.ImageConfig) {
	c.m.Lock()
	defer c.m.Unlock()
	var imageID string

	// Don't assume the image id in image has "sha256:<id> as format.  We store it in
	// this fomat to make it easier to lookup by digest
	if strings.HasPrefix(imageConfig.ImageID, "sha") {
		imageID = imageConfig.ImageID
	} else {
		imageID = "sha256:" + imageConfig.ImageID
	}

	c.cacheByID[imageID] = imageConfig

	// Normalize the name stored in imageConfig using Docker's reference code
	ref, err := reference.WithName(imageConfig.Name)
	if err != nil {
		log.Errorf("Tried to create reference from %s: %s", imageConfig.Name, err.Error())
		return
	}

	for id := range imageConfig.Tags {
		tag := imageConfig.Tags[id]
		ref, err = reference.WithTag(ref, tag)
		if err != nil {
			log.Errorf("Tried to create tagged reference from %s and tag %s: %s", imageConfig.Name, tag, err.Error())
			return
		}

		if tagged, ok := ref.(reference.NamedTagged); ok {
			taggedName := fmt.Sprintf("%s:%s", tagged.Name(), tagged.Tag())
			c.cacheByName[taggedName] = imageConfig
		} else {
			c.cacheByName[ref.Name()] = imageConfig
		}
	}
}
