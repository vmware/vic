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

	"github.com/vmware/vic/portlayer/util"
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

// GetImageStore checks to see if a named image store exists and returls the
// URL to it if so or error.
func (c *NameLookupCache) GetImageStore(ctx context.Context, storeName string) (*url.URL, error) {
	u, err := util.StoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()
	_, ok := c.storeCache[*u]
	if !ok {
		return nil, os.ErrNotExist
	}

	return u, nil
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

	if c.storeCache == nil {
		c.storeCache = make(map[url.URL]map[string]Image)
	}

	c.storeCache[*u] = make(map[string]Image)

	// Create the root image
	scratch, err := c.DataStore.WriteImage(ctx, &Image{Store: u}, Scratch.ID, nil, nil)
	if err != nil {
		return nil, err
	}

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

func (c *NameLookupCache) WriteImage(ctx context.Context, parent *Image, ID string, meta map[string][]byte, sum string, r io.Reader) (*Image, error) {
	// Check the parent exists (at least in the cache).
	p, err := c.GetImage(ctx, parent.Store, parent.ID)
	if err != nil {
		return nil, fmt.Errorf("parent (%s) doesn't exist in %s", parent.ID, parent.Store.String())
	}

	h := sha256.New()
	t := io.TeeReader(r, h)

	i, err := c.DataStore.WriteImage(ctx, p, ID, meta, t)
	if err != nil {
		return nil, err
	}

	actualSum := fmt.Sprintf("sha256:%x", h.Sum(nil))
	if actualSum != sum {
		// TODO(jzt): cleanup?
		return nil, fmt.Errorf("Failed to validate image checksum. Expected %s, got %s", sum, actualSum)
	}

	// Add the new image to the cache
	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()
	c.storeCache[*p.Store][i.ID] = *i

	return i, nil
}

// GetImage gets the specified image from the given store by retreiving it from the cache.
func (c *NameLookupCache) GetImage(ctx context.Context, store *url.URL, ID string) (*Image, error) {
	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()

	s, ok := c.storeCache[*store]
	if !ok {
		return nil, fmt.Errorf("store (%s) doesn't exist", store.String())
	}

	i, ok := s[ID]
	if !ok {
		return nil, fmt.Errorf("store (%s) doesn't have image %s", store.String(), ID)
	}

	return &i, nil
}

// ListImages resturns a list of Images for a list of IDs, or all if no IDs are passed
func (c *NameLookupCache) ListImages(ctx context.Context, store *url.URL, IDs []string) ([]*Image, error) {
	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()

	// check the store exists
	_, ok := c.storeCache[*store]
	if !ok {
		return nil, fmt.Errorf("store (%s) doesn't exist", store.String())
	}

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
			newImage := v
			imageList = append(imageList, &newImage)
		}
	}

	return imageList, nil
}
