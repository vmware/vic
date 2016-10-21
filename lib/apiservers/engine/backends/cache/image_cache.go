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
	"strings"
	"sync"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/distribution/digest"
	derr "github.com/docker/docker/errors"
	"github.com/docker/docker/pkg/truncindex"
	"github.com/docker/docker/reference"

	"github.com/vmware/vic/lib/apiservers/engine/backends/kv"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
)

// ICache is an in-memory cache of image metadata. It is refreshed at startup
// by a call to the portlayer. It is updated when new images are pulled or
// images are deleted.
type ICache struct {
	m sync.RWMutex
	// IDIndex allows image ID lookup by short prefix
	IDIndex *truncindex.TruncIndex
	// CacheByID maps imageID to image metadata
	CacheByID map[string]*metadata.ImageConfig
	// CacheByName maps image name to image metadata
	CacheByName map[string]*metadata.ImageConfig

	client *client.PortLayer
}

const (
	imageCacheKey = "images"
)

var (
	imageCache *ICache
	ctx        = context.TODO()
)

func init() {
	imageCache = &ICache{
		IDIndex:     truncindex.NewTruncIndex([]string{}),
		CacheByID:   make(map[string]*metadata.ImageConfig),
		CacheByName: make(map[string]*metadata.ImageConfig),
	}
}

// ImageCache returns a reference to the image cache
func ImageCache() *ICache {
	return imageCache
}

// NewImageCache will create a new image cache or rehydrate an
// existing image cache from the portlayer k/v store
func NewImageCache(client *client.PortLayer) error {
	defer trace.End(trace.Begin(""))

	imageCache.client = client

	log.Debugf("Updating image cache")

	val, err := kv.Get(client, imageCacheKey)
	if err != nil && err != kv.ErrKeyNotFound {
		return err
	}

	if val != "" {
		if err = json.Unmarshal([]byte(val), imageCache); err != nil {
			return fmt.Errorf("Failed to unmarshal image cache: %s", err)
		}
	}

	return nil
}

// GetImages returns a slice containing metadata for all cached images
func (ic *ICache) GetImages() []*metadata.ImageConfig {
	defer trace.End(trace.Begin(""))
	ic.m.RLock()
	defer ic.m.RUnlock()

	result := make([]*metadata.ImageConfig, 0, len(ic.CacheByID))
	for _, image := range ic.CacheByID {
		result = append(result, copyImageConfig(image))
	}
	return result
}

// IsImageID will check that a full or partial imageID
// exists in the cache
func (ic *ICache) IsImageID(id string) bool {
	ic.m.RLock()
	defer ic.m.RUnlock()
	if _, err := ic.IDIndex.Get(id); err == nil {
		return true
	}
	return false
}

// Get parses input to retrieve a cached image
func (ic *ICache) Get(idOrRef string) (*metadata.ImageConfig, error) {
	defer trace.End(trace.Begin(""))
	ic.m.RLock()
	defer ic.m.RUnlock()

	// cover the case of creating by a full reference
	if config, ok := ic.CacheByName[idOrRef]; ok {
		return config, nil
	}

	// get the full image ID if supplied a prefix
	if id, err := ic.IDIndex.Get(idOrRef); err == nil {
		idOrRef = id
	}

	imgDigest, named, err := reference.ParseIDOrReference(idOrRef)
	if err != nil {
		return nil, err
	}

	var config *metadata.ImageConfig
	if imgDigest != "" {
		config = ic.getImageByDigest(imgDigest)
	} else {
		config = ic.getImageByNamed(named)
	}

	if config == nil {
		// docker automatically prints out ":latest" tag if not specified in case if image is not found.
		postfixLatest := ""
		if !strings.Contains(idOrRef, ":") {
			postfixLatest += ":" + reference.DefaultTag
		}
		return nil, derr.NewRequestNotFoundError(fmt.Errorf(
			"No such image: %s%s", idOrRef, postfixLatest))
	}

	return copyImageConfig(config), nil
}

func (ic *ICache) getImageByDigest(digest digest.Digest) *metadata.ImageConfig {
	defer trace.End(trace.Begin(""))
	var config *metadata.ImageConfig
	config, ok := ic.CacheByID[string(digest)]
	if !ok {
		return nil
	}
	return copyImageConfig(config)
}

// Looks up image by reference.Named
func (ic *ICache) getImageByNamed(named reference.Named) *metadata.ImageConfig {
	defer trace.End(trace.Begin(""))
	// get the imageID from the repoCache
	id, _ := RepositoryCache().Get(named)
	return copyImageConfig(ic.CacheByID[prefixImageID(id)])
}

// Add the "sha256:" prefix to the image ID if missing.
// Don't assume the image id in image has "sha256:<id> as format.  We store it in
// this format to make it easier to lookup by digest
func prefixImageID(imageID string) string {
	if strings.HasPrefix(imageID, "sha256:") {
		return imageID
	}
	return "sha256:" + imageID
}

// Add adds an image to the image cache
func (ic *ICache) Add(imageConfig *metadata.ImageConfig, save bool) {
	defer trace.End(trace.Begin(""))

	if save {
		defer ic.Save()
	}

	ic.m.Lock()
	defer ic.m.Unlock()

	// Normalize the name stored in imageConfig using Docker's reference code
	ref, err := reference.WithName(imageConfig.Name)
	if err != nil {
		log.Errorf("Tried to create reference from %s: %s", imageConfig.Name, err.Error())
		return
	}

	imageID := prefixImageID(imageConfig.ImageID)
	ic.IDIndex.Add(imageConfig.ImageID)
	ic.CacheByID[imageID] = imageConfig

	for _, tag := range imageConfig.Tags {
		ref, err = reference.WithTag(ref, tag)
		if err != nil {
			log.Errorf("Tried to create tagged reference from %s and tag %s: %s", imageConfig.Name, tag, err.Error())
			return
		}
		ic.CacheByName[imageConfig.Reference] = imageConfig
	}
}

// RemoveImageByConfig removes image from the cache.
func (ic *ICache) RemoveImageByConfig(imageConfig *metadata.ImageConfig) {
	defer trace.End(trace.Begin(""))
	ic.m.Lock()
	defer ic.m.Unlock()

	// If we get here we definitely want to remove image config from any data structure
	// where it can be present. So that, if there is something is wrong
	// it could be tracked on debug level.
	if err := ic.IDIndex.Delete(imageConfig.ImageID); err != nil {
		log.Debugf("Not found in image cache index: %v", err)
	}

	prefixedID := prefixImageID(imageConfig.ImageID)
	if _, ok := ic.CacheByID[prefixedID]; ok {
		delete(ic.CacheByID, prefixedID)
	} else {
		log.Debugf("Not found in cache by id: %s", prefixedID)
	}

	if _, ok := ic.CacheByName[imageConfig.Reference]; ok {
		delete(ic.CacheByName, imageConfig.Reference)
	} else {
		log.Debugf("Not found in cache by name: %s", imageConfig.Reference)
	}
}

// Save will persist the image cache to the portlayer k/v store
func (ic *ICache) Save() error {
	defer trace.End(trace.Begin(""))
	ic.m.Lock()
	defer ic.m.Unlock()

	bytes, err := json.Marshal(ic)
	if err != nil {
		log.Errorf("Unable to marshal image cache: %s", err.Error())
		return err
	}

	err = kv.Put(ic.client, imageCacheKey, string(bytes))
	if err != nil {
		log.Errorf("Unable to save image cache: %s", err.Error())
		return err
	}

	return nil
}

// copyImageConfig performs and returns deep copy of an ImageConfig struct
func copyImageConfig(image *metadata.ImageConfig) *metadata.ImageConfig {

	if image == nil {
		return nil
	}

	// copy everything
	newImage := *image

	// replace the pointer to metadata.ImageConfig.Config and copy the contents
	newConfig := *image.Config
	newImage.Config = &newConfig

	// add tags & digests from the repo cache
	newImage.Tags = RepositoryCache().Tags(newImage.ImageID)
	newImage.Digests = RepositoryCache().Digests(newImage.ImageID)

	return &newImage
}
