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
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"sync"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/index"
)

var Scratch = Image{
	ID: "scratch",
}

var ErrCorruptImageStore = errors.New("Corrupt image store")

// NameLookupCache the global view of all of the image stores.  To avoid unnecessary
// lookups, the image cache keeps an in memory map of the store URI to the map
// of images on disk.
type NameLookupCache struct {

	// The individual store locations -> Index
	storeCache map[url.URL]*index.Index
	// Guard against concurrent writes to the storeCache map
	storeCacheLock sync.Mutex

	// The image store implementation.  This mutates the actual disk images.
	DataStore ImageStorer
}

func NewLookupCache(ds ImageStorer) *NameLookupCache {
	return &NameLookupCache{
		DataStore:  ds,
		storeCache: make(map[url.URL]*index.Index),
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
		infof("Refreshing image cache from datastore.")
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

		indx := index.NewIndex()

		c.storeCache[*store] = indx

		// Add Scratch
		scratch, err := c.DataStore.GetImage(ctx, store, Scratch.ID)
		if err != nil {
			log.Errorf("ImageCache Error: looking up scratch on %s: %s", store.String(), err)
			return nil, ErrCorruptImageStore
		}

		if err = indx.Insert(scratch); err != nil {
			return nil, err
		}

		// XXX after creating the indx and populating the map, we can put the rest in a go routine

		// Fall out here if there are no images.  We should at least have a scratch.
		images, err := c.DataStore.ListImages(ctx, store, nil)
		if err != nil {
			return nil, err
		}

		// add the images we retrieved to the cache.
		for _, image := range images {
			if image.ID == Scratch.ID {
				continue
			}

			infof("Found image %s on datastore.", image.ID)

			if err := indx.Insert(image); err != nil {
				return nil, err
			}
		}
	}

	return store, nil
}

func (c *NameLookupCache) CreateImageStore(ctx context.Context, storeName string) (*url.URL, error) {
	store, err := util.ImageStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	// Check for existance and rehydrate the cache if it exists on disk.
	_, err = c.GetImageStore(ctx, storeName)
	// we expect this not to exist.
	if err == nil {
		return nil, os.ErrExist
	}

	c.storeCacheLock.Lock()
	defer c.storeCacheLock.Unlock()

	store, err = c.DataStore.CreateImageStore(ctx, storeName)
	if err != nil {
		return nil, err
	}

	// Create the root image
	scratch, err := c.DataStore.WriteImage(ctx, &Image{Store: store}, Scratch.ID, nil, "", nil)
	if err != nil {
		return nil, err
	}

	indx := index.NewIndex()
	c.storeCache[*store] = indx
	if err = indx.Insert(scratch); err != nil {
		return nil, err
	}

	return store, nil
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
		return nil, fmt.Errorf("parent (%s) doesn't exist in %s: %s", parent.ID, parent.Store.String(), err)
	}

	// Check the image doesn't already exist in the cache.  A miss in this will trigger a datastore lookup.
	i, err := c.GetImage(ctx, p.Store, ID)
	if err == nil && i != nil {
		// TODO(FA) check sums to make sure this is the right image

		return i, nil
	}

	// Definitely not in cache or image store, create image.
	i, err = c.DataStore.WriteImage(ctx, p, ID, meta, sum, r)
	if err != nil {
		errorf("WriteImage of %s failed with: %s", ID, err)
		return nil, err
	}

	c.storeCacheLock.Lock()
	indx := c.storeCache[*parent.Store]
	c.storeCacheLock.Unlock()

	// Add the new image to the cache
	if err = indx.Insert(i); err != nil {
		return nil, err
	}

	return i, nil
}

// GetImage gets the specified image from the given store by retreiving it from the cache.
func (c *NameLookupCache) GetImage(ctx context.Context, store *url.URL, ID string) (*Image, error) {

	debugf("Getting image %s from %s", ID, store.String())

	storeName, err := util.ImageStoreName(store)
	if err != nil {
		return nil, err
	}

	// Check the store exists
	if _, err = c.GetImageStore(ctx, storeName); err != nil {
		return nil, err
	}

	c.storeCacheLock.Lock()
	indx := c.storeCache[*store]
	c.storeCacheLock.Unlock()

	imgUrl, err := util.ImageURL(storeName, ID)
	if err != nil {
		return nil, err
	}
	node, err := c.storeCache[*store].Get(imgUrl.String())

	var img *Image
	if err != nil {
		if err == index.ErrNodeNotFound {
			debugf("Image %s not in cache, retreiving from datastore", ID)
			// Not in the cache.  Try to load it.
			img, err = c.DataStore.GetImage(ctx, store, ID)
			if err != nil {
				return nil, err
			}

			if err = indx.Insert(img); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		img, _ = node.(*Image)
	}

	return img, nil
}

// ListImages returns a list of Images for a list of IDs, or all if no IDs are passed
func (c *NameLookupCache) ListImages(ctx context.Context, store *url.URL, IDs []string) ([]*Image, error) {
	// Filter the results
	imageList := make([]*Image, 0, len(IDs))

	if len(IDs) > 0 {
		for _, id := range IDs {
			i, err := c.GetImage(ctx, store, id)
			if err == nil {
				imageList = append(imageList, i)
			}
		}

	} else {

		storeName, err := util.ImageStoreName(store)
		if err != nil {
			return nil, err
		}
		// Check the store exists before we start iterating it.  This will populate the cache if it's empty.
		if _, err := c.GetImageStore(ctx, storeName); err != nil {
			return nil, err
		}

		// get the relevant cache
		c.storeCacheLock.Lock()
		indx := c.storeCache[*store]
		c.storeCacheLock.Unlock()

		images, err := indx.List()
		if err != nil {
			return nil, err
		}

		for _, v := range images {
			img, _ := v.(*Image)
			// filter out scratch
			if img.ID == Scratch.ID {
				continue
			}

			imageList = append(imageList, img)
		}
	}

	return imageList, nil
}

// DeleteImage deletes an image from the image store.  If it is in use or is being inheritted from, then this will return an error.
func (c *NameLookupCache) DeleteImage(ctx context.Context, image *Image) error {
	infof("DeleteImage: deleting %s", image.Self())

	// Check the image exists.  This will rehydrate the cache if necessary.
	img, err := c.GetImage(ctx, image.Store, image.ID)
	if err != nil {
		errorf("DeleteImage: %s", err)
		return err
	}

	// get the relevant cache
	c.storeCacheLock.Lock()
	indx := c.storeCache[*img.Store]
	c.storeCacheLock.Unlock()

	hasChildren, err := indx.HasChildren(img.Self())
	if err != nil {
		errorf("DeleteImage: %s", err)
		return err
	}

	if hasChildren {
		return &ErrImageInUse{img.Self() + " in use by child images"}
	}

	// The datastore will tell us if the image is attached
	if err = c.DataStore.DeleteImage(ctx, img); err != nil {
		errorf("%s", err)
		return err
	}

	// Remove the image from the cache
	if _, err = indx.Delete(img.Self()); err != nil {
		errorf("%s", err)
		return err
	}

	return nil
}

func infof(format string, args ...interface{}) {
	log.Infof("ImageCache: "+format, args...)
}

func errorf(format string, args ...interface{}) {
	err := fmt.Errorf("ImageCache error: "+format, args...)
	log.Errorf(err.Error())
}

func debugf(format string, args ...interface{}) {
	log.Debugf("ImageCache: "+format, args...)
}
