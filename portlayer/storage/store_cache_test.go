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
	"fmt"
	"io"
	"net/url"
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/portlayer/util"
)

type MockDataStore struct {
	// id -> image
	images map[string]*Image
}

func NewMockDataStore() *MockDataStore {
	return &MockDataStore{
		images: make(map[string]*Image),
	}
}

// GetImageStore checks to see if a named image store exists and returls the
// URL to it if so or error.
func (c *MockDataStore) GetImageStore(ctx context.Context, storeName string) (*url.URL, error) {
	return nil, nil
}

func (c *MockDataStore) CreateImageStore(ctx context.Context, storeName string) (*url.URL, error) {
	u, err := util.StoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (c *MockDataStore) ListImageStores(ctx context.Context) ([]*url.URL, error) {
	return nil, nil
}

func (c *MockDataStore) WriteImage(ctx context.Context, parent *Image, ID string, meta map[string][]byte, r io.Reader) (*Image, error) {
	i := &Image{
		ID:       ID,
		Store:    parent.Store,
		Parent:   parent.SelfLink,
		Metadata: meta,
	}

	c.images[ID] = i

	return i, nil
}

// GetImage gets the specified image from the given store by retreiving it from the cache.
func (c *MockDataStore) GetImage(ctx context.Context, store *url.URL, ID string) (*Image, error) {
	i, ok := c.images[ID]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return i, nil
}

// ListImages resturns a list of Images for a list of IDs, or all if no IDs are passed
func (c *MockDataStore) ListImages(ctx context.Context, store *url.URL, IDs []string) ([]*Image, error) {
	return nil, nil
}

func TestListImages(t *testing.T) {
	s := NewLookupCache(NewMockDataStore())

	storeURL, err := s.CreateImageStore(context.TODO(), "testStore")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, storeURL) {
		return
	}

	// Create a set of images
	images := make(map[string]*Image)
	images[Scratch.ID] = &Scratch
	parent := Scratch
	parent.Store = storeURL
	testSum := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	for i := 1; i < 50; i++ {
		id := fmt.Sprintf("ID-%d", i)

		img, err := s.WriteImage(context.TODO(), &parent, id, nil, testSum, nil)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, img) {
			return
		}

		images[id] = img
	}

	// List all images
	outImages, err := s.ListImages(context.TODO(), storeURL, nil)
	if !assert.NoError(t, err) {
		return
	}

	// check we retrieve all of the iamges
	assert.Equal(t, len(outImages), len(images))
	for _, img := range outImages {
		_, ok := images[img.ID]
		if !assert.True(t, ok) {
			return
		}
	}

	// Check we can retrieve a subset
	inIDs := []string{"ID-1", "ID-2", "ID-3"}
	outImages, err = s.ListImages(context.TODO(), storeURL, inIDs)
	if !assert.NoError(t, err) {
		return
	}

	for _, img := range outImages {
		reference, ok := images[img.ID]
		if !assert.True(t, ok) {
			return
		}

		if !assert.Equal(t, reference, img) {
			return
		}
	}
}

// Create 2 store caches but use the same backing datastore.  Create images
// with the first cache, then get the image with the second.  This simulates
// restart since the second cache is empty and has to go to the backing store.
func TestImageStoreRestart(t *testing.T) {
	ds := NewMockDataStore()

	firstCache := NewLookupCache(ds)
	secondCache := NewLookupCache(ds)

	storeURL, err := firstCache.CreateImageStore(context.TODO(), "testStore")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, storeURL) {
		return
	}

	// Create a set of images
	expectedImages := make(map[string]*Image)

	parent := Scratch
	parent.Store = storeURL

	testSum := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	for i := 1; i < 50; i++ {
		id := fmt.Sprintf("ID-%d", i)

		img, err := firstCache.WriteImage(context.TODO(), &parent, id, nil, testSum, nil)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, img) {
			return
		}

		expectedImages[id] = img
	}

	// get the images from the second cache to ensure it goes to the ds
	for id, expectedImg := range expectedImages {
		img, err := secondCache.GetImage(context.TODO(), storeURL, id)
		if !assert.NoError(t, err) || !assert.Equal(t, expectedImg, img) {
			return
		}
	}

	// Nuke the second cache's datastore.  All data should come from the cache.
	secondCache.DataStore = nil
	for id, expectedImg := range expectedImages {
		img, err := secondCache.GetImage(context.TODO(), storeURL, id)
		if !assert.NoError(t, err) || !assert.Equal(t, expectedImg, img) {
			return
		}
	}
}
