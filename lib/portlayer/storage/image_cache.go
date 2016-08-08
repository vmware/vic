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

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/url"
	"os"
	"sync"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/lib/portlayer/util"
)

var Scratch = Image{
	ID: "scratch",
}

// NameLookupCache the global view of all of the image stores.  To avoid unnecessary
// lookups, the image cache keeps an in memory map of the store URI to the map
// of images on disk.
type NameLookupCache struct {

	// The individual store locations
	//
	// We want to map a store to a list of images.  Images are resolveable by
	// ID (string) and resolve to an Image.  The keys/values in the map should
	// be immuteable, so we're passing by value here.  We don't want things
	// changing outside of the API calls.
	storeCache     map[url.URL]map[string]Image
	storeCacheLock sync.Mutex

	// The image store implementation.  This mutates the actual disk images.
	DataStore ImageStorer
}

func NewLookupCache(ds ImageStorer) *NameLookupCache {
	return &NameLookupCache{
		DataStore:  ds,
		storeCache: make(map[url.URL]map[string]Image),
	}
}

// GetImageStore checks to see if a named image store exists and returls the
// URL to it if so or error.
func (c *NameLookupCache) GetImageStore(ctx context.Context, storeName string) (*url.URL, error) {
	store, err := util.ImageStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()

	// check the cache
	_, ok := c.storeCache[*store]
	if !ok {
		log.Info("Refreshing image cache from datastore.")
		// Store isn't in the cache.  Look it up in the datastore.
		storeName, err := util.ImageStoreName(store)
		if err != nil {
			return nil, err
		}

		// If the store doesn't exist, we'll fall out here.
		_, err = c.DataStore.GetImageStore(ctx, storeName)
		if err != nil {
			return nil, err
		}

		c.storeCache[*store] = make(map[string]Image)

		// Fall out here if there are no images.  We should at least have a scratch.
		images, err := c.DataStore.ListImages(ctx, store, nil)
		if err != nil {
			return nil, err
		}

		// add the images we retrieved to the cache.
		for _, v := range images {
			log.Infof("Imagestore: Found image %s on datastore.", v.ID)
			c.storeCache[*store][v.ID] = *v
		}

		// Assert there's a scratch
		if _, ok = c.storeCache[*store][Scratch.ID]; !ok {
			return nil, fmt.Errorf("Scratch does not exist.  Imagestore is corrrupt.")
		}
	}

	return store, nil
}

func (c *NameLookupCache) CreateImageStore(ctx context.Context, storeName string) (*url.URL, error) {
	u, err := c.GetImageStore(ctx, storeName)
	// we expect this not to exist.
	if err == nil {
		return nil, os.ErrExist
	}

	u, err = c.DataStore.CreateImageStore(ctx, storeName)
	if err != nil {
		return nil, err
	}

	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()

	// Create the root image
	scratch, err := c.DataStore.WriteImage(ctx, &Image{Store: u}, Scratch.ID, nil, nil)
	if err != nil {
		return nil, err
	}

	c.storeCache[*u] = make(map[string]Image)
	c.storeCache[*u][scratch.ID] = *scratch
	return u, nil
}

// ListImageStores returns a list of strings representing all existing image stores
func (c *NameLookupCache) ListImageStores(ctx context.Context) ([]*url.URL, error) {
	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()

	stores := make([]*url.URL, 0, len(c.storeCache))
	for key := range c.storeCache {
		stores = append(stores, &key)
	}
	return stores, nil
}

// add to store cache
func (c *NameLookupCache) AddImageToStore(storeURL url.URL, imageID string, image Image) {
	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()
	c.storeCache[storeURL][imageID] = image
}

func (c *NameLookupCache) WriteImage(ctx context.Context, parent *Image, ID string, meta map[string][]byte, sum string, r io.Reader) (*Image, error) {
	// Check the parent exists (at least in the cache).
	p, err := c.GetImage(ctx, parent.Store, parent.ID)
	if err != nil {
		return nil, fmt.Errorf("parent (%s) doesn't exist in %s: %s", parent.ID, parent.Store.String(), err)
	}

	// Check the image doesn't already exist in the cache.  A miss in this will trigger a datastore lookup.
	i, err := c.GetImage(ctx, p.Store, ID)
	if err == nil && i != nil {
		// TODO(FA) check sums to make sure this is the right image

		return i, nil
	}

	// Definitely not in cache or image store, create image.
	h := sha256.New()
	t := io.TeeReader(r, h)

	i, err = c.DataStore.WriteImage(ctx, p, ID, meta, t)
	if err != nil {
		return nil, err
	}

	actualSum := fmt.Sprintf("sha256:%x", h.Sum(nil))
	if actualSum != sum {
		// TODO(jzt): cleanup?
		return nil, fmt.Errorf("Failed to validate image checksum. Expected %s, got %s", sum, actualSum)
	}

	// Add the new image to the cache
	c.AddImageToStore(*p.Store, i.ID, *i)

	return i, nil
}

// GetImage gets the specified image from the given store by retreiving it from the cache.
func (c *NameLookupCache) GetImage(ctx context.Context, store *url.URL, ID string) (*Image, error) {
	log.Debugf("Getting image from %#v", store)
	storeName, err := util.ImageStoreName(store)
	if err != nil {
		return nil, err
	}

	// Check the store exists
	if _, err = c.GetImageStore(ctx, storeName); err != nil {
		return nil, err
	}

	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()

	var ok bool
	i := &Image{}
	*i, ok = c.storeCache[*store][ID]
	if !ok {
		log.Infof("Image %s not in cache, retreiving from datastore", ID)
		// Not in the cache.  Try to load it.
		i, err = c.DataStore.GetImage(ctx, store, ID)
		if err != nil {
			return nil, err
		}

		c.storeCache[*store][ID] = *i
	}

	return i, nil
}

// ListImages resturns a list of Images for a list of IDs, or all if no IDs are passed
func (c *NameLookupCache) ListImages(ctx context.Context, store *url.URL, IDs []string) ([]*Image, error) {
	storeName, err := util.ImageStoreName(store)
	if err != nil {
		return nil, err
	}

	// Check the store exists before we start iterating it.  This will populate the cache if it's empty.
	if _, err := c.GetImageStore(ctx, storeName); err != nil {
		return nil, err
	}

	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()

	// Filter the results
	var imageList []*Image
	if len(IDs) > 0 {
		for _, id := range IDs {
			if i, ok := c.storeCache[*store][id]; ok {
				newImage := i
				imageList = append(imageList, &newImage)
			}
		}
	} else {
		for _, v := range c.storeCache[*store] {
			// filter out scratch
			if v.ID == Scratch.ID {
				continue
			}
			newImage := v
			imageList = append(imageList, &newImage)
		}
	}

	return imageList, nil
}
