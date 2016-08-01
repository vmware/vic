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
	"github.com/docker/docker/pkg/truncindex"
	"github.com/docker/docker/reference"
	"github.com/docker/engine-api/types/container"

	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/vsphere/sys"
)

// ICache is an in-memory cache of image metadata. It is refreshed at startup
// by a call to the portlayer. It is updated when new images are pulled or
// images are deleted.
type ICache struct {
	m sync.RWMutex

	// cache maps image ID to image metadata
	idIndex     *truncindex.TruncIndex
	cacheByID   map[string]*metadata.ImageConfig
	cacheByName map[string]*metadata.ImageConfig
}

var imageCache *ICache

func init() {
	imageCache = &ICache{
		idIndex:     truncindex.NewTruncIndex([]string{}),
		cacheByID:   make(map[string]*metadata.ImageConfig),
		cacheByName: make(map[string]*metadata.ImageConfig),
	}
}

// ImageCache returns a reference to the image cache
func ImageCache() *ICache {
	return imageCache
}

// Update runs only once at startup to hydrate the image cache
func (ic *ICache) Update(client *client.PortLayer) error {

	log.Debugf("Updating image cache...")

	host, err := sys.UUID()
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
			ic.AddImage(imageConfig)
		}
	}

	return nil
}

// GetImages returns a slice containing metadata for all cached images
func (ic *ICache) GetImages() []*metadata.ImageConfig {
	ic.m.RLock()
	defer ic.m.RUnlock()

	result := make([]*metadata.ImageConfig, 0, len(ic.cacheByID))
	for _, image := range ic.cacheByID {
		result = append(result, copyImageConfig(image))
	}
	return result
}

// GetImage parses input to retrieve a cached image
func (ic *ICache) GetImage(idOrRef string) (*metadata.ImageConfig, error) {
	ic.m.RLock()
	defer ic.m.RUnlock()

	// get the full image ID if supplied a prefix
	if id, err := ic.idIndex.Get(idOrRef); err == nil {
		idOrRef = id
	}

	digest, named, err := reference.ParseIDOrReference(idOrRef)
	if err != nil {
		return nil, err
	}

	var config *metadata.ImageConfig
	if digest != "" {
		config, err = ic.getImageByDigest(digest)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = ic.getImageByNamed(named)
		if err != nil {
			return nil, err
		}
	}

	if config == nil {
		return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such image: %s", idOrRef))
	}

	return copyImageConfig(config), nil
}

func (ic *ICache) getImageByDigest(digest digest.Digest) (*metadata.ImageConfig, error) {
	var config *metadata.ImageConfig
	config, ok := ic.cacheByID[string(digest)]
	if !ok {
		return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such image: %s", digest))
	}
	return copyImageConfig(config), nil
}

// Looks up image by reference.Named
func (ic *ICache) getImageByNamed(named reference.Named) (*metadata.ImageConfig, error) {

	var config *metadata.ImageConfig

	if tagged, ok := named.(reference.NamedTagged); ok {
		taggedName := tagged.Name() + ":" + tagged.Tag()
		config = ic.cacheByName[taggedName]
	} else {
		// First try just the name.
		if config, ok = ic.cacheByName[named.Name()]; !ok {
			// try with the default docker tag
			taggedName := named.Name() + ":" + reference.DefaultTag
			config = ic.cacheByName[taggedName]
		}
	}

	return copyImageConfig(config), nil
}

// AddImage adds an image to the image cache
func (ic *ICache) AddImage(imageConfig *metadata.ImageConfig) {

	ic.m.Lock()
	defer ic.m.Unlock()

	var imageID string

	// Don't assume the image id in image has "sha256:<id> as format.  We store it in
	// this fomat to make it easier to lookup by digest
	if strings.HasPrefix(imageConfig.ImageID, "sha") {
		imageID = imageConfig.ImageID
	} else {
		imageID = "sha256:" + imageConfig.ImageID
	}

	// add image id to truncindex
	ic.idIndex.Add(imageConfig.ImageID)

	ic.cacheByID[imageID] = imageConfig

	// Normalize the name stored in imageConfig using Docker's reference code
	ref, err := reference.WithName(imageConfig.Name)
	if err != nil {
		log.Errorf("Tried to create reference from %s: %s", imageConfig.Name, err.Error())
		return
	}

	for _, tag := range imageConfig.Tags {
		ref, err = reference.WithTag(ref, tag)
		if err != nil {
			log.Errorf("Tried to create tagged reference from %s and tag %s: %s", imageConfig.Name, tag, err.Error())
			return
		}

		if tagged, ok := ref.(reference.NamedTagged); ok {
			taggedName := fmt.Sprintf("%s:%s", tagged.Name(), tagged.Tag())
			ic.cacheByName[taggedName] = imageConfig
		} else {
			ic.cacheByName[ref.Name()] = imageConfig
		}
	}
}

// copyImageConfig performs and returns deep copy of an ImageConfig struct
func copyImageConfig(image *metadata.ImageConfig) *metadata.ImageConfig {

	if image == nil {
		return nil
	}

	newImage := new(metadata.ImageConfig)

	// copy everything
	*newImage = *image

	// replace the pointer to metadata.ImageConfig.Config and copy the contents
	newConfig := new(container.Config)
	*newConfig = *image.Config
	newImage.Config = newConfig

	return newImage
}
