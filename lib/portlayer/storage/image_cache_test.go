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
	"strconv"
	"testing"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/trace"
)

type MockDataStore struct {
	// id -> image
	db map[url.URL]map[string]*Image
}

func NewMockDataStore() *MockDataStore {
	m := &MockDataStore{
		db: make(map[url.URL]map[string]*Image),
	}

	return m
}

// GetImageStore checks to see if a named image store exists and returls the
// URL to it if so or error.
func (c *MockDataStore) GetImageStore(op trace.Operation, storeName string) (*url.URL, error) {
	return nil, nil
}

func (c *MockDataStore) CreateImageStore(op trace.Operation, storeName string) (*url.URL, error) {
	u, err := util.ImageStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	c.db[*u] = make(map[string]*Image)
	return u, nil
}

func (c *MockDataStore) ListImageStores(op trace.Operation) ([]*url.URL, error) {
	return nil, nil
}

func (c *MockDataStore) WriteImage(op trace.Operation, parent *Image, ID string, meta map[string][]byte, sum string, r io.Reader) (*Image, error) {
	storeName, err := util.ImageStoreName(parent.Store)
	if err != nil {
		return nil, err
	}

	selflink, err := util.ImageURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	var parentLink *url.URL
	if parent.ID != "" {
		parentLink, err = util.ImageURL(storeName, parent.ID)
		if err != nil {
			return nil, err
		}
	}

	i := &Image{
		ID:         ID,
		Store:      parent.Store,
		ParentLink: parentLink,
		SelfLink:   selflink,
		Metadata:   meta,
	}

	c.db[*parent.Store][ID] = i

	return i, nil
}

// GetImage gets the specified image from the given store by retreiving it from the cache.
func (c *MockDataStore) GetImage(op trace.Operation, store *url.URL, ID string) (*Image, error) {
	i, ok := c.db[*store][ID]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return i, nil
}

// ListImages resturns a list of Images for a list of IDs, or all if no IDs are passed
func (c *MockDataStore) ListImages(op trace.Operation, store *url.URL, IDs []string) ([]*Image, error) {
	var imageList []*Image
	for _, i := range c.db[*store] {
		imageList = append(imageList, i)
	}
	return imageList, nil
}

// DeleteImage removes an image from the image store
func (c *MockDataStore) DeleteImage(op trace.Operation, image *Image) error {
	delete(c.db[*image.Store], image.ID)
	return nil
}

func TestListImages(t *testing.T) {
	s := NewLookupCache(NewMockDataStore())

	op := trace.NewOperation(context.Background(), "test")
	storeURL, err := s.CreateImageStore(op, "testStore")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, storeURL) {
		return
	}

	// Create a set of images
	images := make(map[string]*Image)
	parent := Scratch
	parent.Store = storeURL
	testSum := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	for i := 1; i < 50; i++ {
		id := fmt.Sprintf("ID-%d", i)

		img, werr := s.WriteImage(op, &parent, id, nil, testSum, nil)
		if !assert.NoError(t, werr) {
			return
		}
		if !assert.NotNil(t, img) {
			return
		}

		images[id] = img
	}

	// List all images
	outImages, err := s.ListImages(op, storeURL, nil)
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
	outImages, err = s.ListImages(op, storeURL, inIDs)

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

// Create an image on the datastore directly and try to WriteImage via the
// cache.  The datastore should reflect the image already exists and bale out
// without an error.
func TestOutsideCacheWriteImage(t *testing.T) {
	s := NewLookupCache(NewMockDataStore())
	op := trace.NewOperation(context.Background(), "test")

	storeURL, err := s.CreateImageStore(op, "testStore")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, storeURL) {
		return
	}

	// Create a set of images
	images := make(map[string]*Image)
	parent := Scratch
	parent.Store = storeURL
	for i := 1; i < 50; i++ {
		id := fmt.Sprintf("ID-%d", i)

		// Write to the datastore creating images
		img, werr := s.DataStore.WriteImage(op, &parent, id, nil, "", nil)
		if !assert.NoError(t, werr) {
			return
		}
		if !assert.NotNil(t, img) {
			return
		}

		images[id] = img
	}

	testSum := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	// Try to write the same images as above, but this time via the cache.  WriteImage should return right away without any data written.
	for i := 1; i < 50; i++ {
		id := fmt.Sprintf("ID-%d", i)

		// Write to the datastore creating images
		img, werr := s.WriteImage(op, &parent, id, nil, testSum, nil)
		if !assert.NoError(t, werr) {
			return
		}
		if !assert.NotNil(t, img) {
			return
		}

		// assert it's the same image
		if !assert.Equal(t, images[img.ID], img) {
			return
		}
	}
}

// Create 2 store caches but use the same backing datastore.  Create images
// with the first cache, then get the image with the second.  This simulates
// restart since the second cache is empty and has to go to the backing store.
func TestImageStoreRestart(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ds := NewMockDataStore()
	op := trace.NewOperation(context.Background(), "test")

	firstCache := NewLookupCache(ds)
	secondCache := NewLookupCache(ds)

	storeURL, err := firstCache.CreateImageStore(op, "testStore")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, storeURL) {
		return
	}

	// Create a set of images
	expectedImages := make(map[string]*Image)

	parent, err := firstCache.GetImage(op, storeURL, Scratch.ID)
	if !assert.NoError(t, err) {
		return
	}

	testSum := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	for i := 1; i < 50; i++ {
		id := fmt.Sprintf("ID-%d", i)

		img, werr := firstCache.WriteImage(op, parent, id, nil, testSum, nil)
		if !assert.NoError(t, werr) {
			return
		}
		if !assert.NotNil(t, img) {
			return
		}

		expectedImages[id] = img
	}

	// get the images from the second cache to ensure it goes to the ds
	for id, expectedImg := range expectedImages {
		img, werr := secondCache.GetImage(op, storeURL, id)
		if !assert.NoError(t, werr) || !assert.Equal(t, expectedImg, img) {
			return
		}
	}

	// Nuke the second cache's datastore.  All data should come from the cache.
	secondCache.DataStore = nil
	for id, expectedImg := range expectedImages {
		img, gerr := secondCache.GetImage(op, storeURL, id)
		if !assert.NoError(t, gerr) || !assert.Equal(t, expectedImg, img) {
			return
		}
	}

	// Same should happen with a third cache when image list is called
	thirdCache := NewLookupCache(ds)
	imageList, err := thirdCache.ListImages(op, storeURL, nil)
	if !assert.NoError(t, err) || !assert.NotNil(t, imageList) {
		return
	}

	if !assert.Equal(t, len(expectedImages), len(imageList)) {
		return
	}

	// check the image data is the same
	for id, expectedImg := range expectedImages {
		img, err := thirdCache.GetImage(op, storeURL, id)
		if !assert.NoError(t, err) || !assert.Equal(t, expectedImg, img) {
			return
		}
	}
}

func TestDeleteImage(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	imageCache := NewLookupCache(NewMockDataStore())
	op := trace.NewOperation(context.Background(), "test")

	storeURL, err := imageCache.CreateImageStore(op, "testStore")
	if !assert.NoError(t, err) || !assert.NotNil(t, storeURL) {
		return
	}

	scratch, err := imageCache.GetImage(op, storeURL, Scratch.ID)
	if !assert.NoError(t, err) {
		return
	}

	// create a 3 level tree with 4 branches
	branches := 4
	images := make(map[int]*Image)
	for branch := 1; branch < branches; branch++ {
		// level 1
		img, err := imageCache.WriteImage(op, scratch, strconv.Itoa(branch), nil, "", nil)
		if !assert.NoError(t, err) || !assert.NotNil(t, img) {
			return
		}
		images[branch] = img

		// level 2
		i, err := imageCache.WriteImage(op, img, strconv.Itoa(branch*10), nil, "", nil)
		if !assert.NoError(t, err) || !assert.NotNil(t, i) {
			return
		}
		images[branch*10] = i

		// level 3
		i, err = imageCache.WriteImage(op, img, strconv.Itoa(branch*100), nil, "", nil)
		if !assert.NoError(t, err) || !assert.NotNil(t, i) {
			return
		}
		images[branch*100] = i
	}

	// Deletion of an intermediate node should fail
	err = imageCache.DeleteImage(op, images[1])
	if !assert.Error(t, err) {
		return
	}

	imageList, err := imageCache.ListImages(op, storeURL, nil)
	if !assert.NoError(t, err) || !assert.NotNil(t, imageList) {
		return
	}

	// image list should be uneffected
	if !assert.Equal(t, len(images), len(imageList)) {
		return
	}

	// Deletion of leaves should be fine
	for branch := 1; branch < branches; branch++ {
		// range up the branch
		for _, img := range []*Image{images[branch*100], images[branch*10], images[branch]} {

			err = imageCache.DeleteImage(op, img)
			if !assert.NoError(t, err) {
				return
			}

			// the image should be gone
			i, err := imageCache.GetImage(op, storeURL, img.ID)
			if !assert.Error(t, err) || !assert.Nil(t, i) {
				return
			}
		}
	}

	// List images should be empty (because we filter out scratch)
	imageList, err = imageCache.ListImages(op, storeURL, nil)
	if !assert.NoError(t, err) || !assert.NotNil(t, imageList) {
		return
	}

	if !assert.True(t, len(imageList) == 0) {
		return
	}
}

// Cache population should be happening in order starting from parent(id1) to children(id4)
func TestPopulateCacheInExpectedOrder(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	st := NewMockDataStore()
	op := trace.NewOperation(context.Background(), "test")

	storeURL, _ := util.ImageStoreNameToURL("testStore")

	storageURLStr := storeURL.String()

	url1, _ := url.Parse(storageURLStr + "/id1")
	url2, _ := url.Parse(storageURLStr + "/id2")
	url3, _ := url.Parse(storageURLStr + "/id3")
	url4, _ := url.Parse(storageURLStr + "/id4")
	scratchURL, _ := url.Parse(storageURLStr + Scratch.ID)

	img1 := &Image{ID: "id1", SelfLink: url1, ParentLink: scratchURL, Store: storeURL}
	img2 := &Image{ID: "id2", SelfLink: url2, ParentLink: url1, Store: storeURL}
	img3 := &Image{ID: "id3", SelfLink: url3, ParentLink: url2, Store: storeURL}
	img4 := &Image{ID: "id4", SelfLink: url4, ParentLink: url3, Store: storeURL}
	scratchImg := &Image{
		ID:         Scratch.ID,
		SelfLink:   scratchURL,
		ParentLink: scratchURL,
		Store:      storeURL,
	}

	// Order does matter for some reason.
	imageMap := map[string]*Image{
		img1.ID:       img1,
		img4.ID:       img4,
		img2.ID:       img2,
		img3.ID:       img3,
		scratchImg.ID: scratchImg,
	}

	st.db[*storeURL] = imageMap

	imageCache := NewLookupCache(st)
	imageCache.GetImageStore(op, "testStore")

	// Check if all images are available.
	imageIds := []string{"id1", "id2", "id3", "id4"}
	for _, imageID := range imageIds {
		v, _ := imageCache.GetImage(op, storeURL, imageID)
		assert.NotNil(t, v)
	}
}
