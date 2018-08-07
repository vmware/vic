// Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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
	"context"
	"fmt"
	"net/url"
	"strconv"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/vic/lib/constants"
	"github.com/vmware/vic/lib/portlayer/storage/image"
	"github.com/vmware/vic/lib/portlayer/storage/image/mock"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/trace"
)

func TestListImages(t *testing.T) {
	s := NewLookupCache(mock.NewMockDataStore(nil))

	op := trace.NewOperation(context.Background(), "test")
	storeURL, err := s.CreateImageStore(op, "testStore")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, storeURL) {
		return
	}

	// Create a set of images
	images := make(map[string]*image.Image)
	parent := image.Image{
		ID: constants.ScratchLayerID,
	}
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
	s := NewLookupCache(mock.NewMockDataStore(nil))
	op := trace.NewOperation(context.Background(), "test")

	storeURL, err := s.CreateImageStore(op, "testStore")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, storeURL) {
		return
	}

	// Create a set of images
	images := make(map[string]*image.Image)
	parent := image.Image{
		ID: constants.ScratchLayerID,
	}
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
	ds := mock.NewMockDataStore(nil)
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
	expectedImages := make(map[string]*image.Image)

	parent, err := firstCache.GetImage(op, storeURL, constants.ScratchLayerID)
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
	imageCache := NewLookupCache(mock.NewMockDataStore(nil))
	op := trace.NewOperation(context.Background(), "test")

	storeURL, err := imageCache.CreateImageStore(op, "testStore")
	if !assert.NoError(t, err) || !assert.NotNil(t, storeURL) {
		return
	}

	scratch, err := imageCache.GetImage(op, storeURL, constants.ScratchLayerID)
	if !assert.NoError(t, err) {
		return
	}

	// create a 3 level tree with 4 branches
	branches := 4
	images := make(map[int]*image.Image)
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
	_, err = imageCache.DeleteImage(op, images[1])
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
		for _, img := range []*image.Image{images[branch*100], images[branch*10], images[branch]} {

			_, err = imageCache.DeleteImage(op, img)
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

func TestDeleteBranch(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	trace.Logger.Level = logrus.DebugLevel

	imageCache := NewLookupCache(mock.NewMockDataStore(nil))
	op := trace.NewOperation(context.Background(), "test")

	storeURL, err := imageCache.CreateImageStore(op, "testStore")
	if !assert.NoError(t, err) || !assert.NotNil(t, storeURL) {
		return
	}

	scratch, err := imageCache.GetImage(op, storeURL, constants.ScratchLayerID)
	if !assert.NoError(t, err) {
		return
	}

	// create a 3 level tree with 3 branches.  The third branch will have an extra node.
	//             scratch
	//        1    2      3
	//       10   20      30
	//       100  200   300 301
	branches := 4
	images := make(map[int]*image.Image)
	for branch := 1; branch < branches; branch++ {
		// level 1
		img, err := imageCache.WriteImage(op, scratch, strconv.Itoa(branch), nil, "", nil)
		if !assert.NoError(t, err) || !assert.NotNil(t, img) {
			return
		}
		images[branch] = img

		// level 2
		img, err = imageCache.WriteImage(op, img, strconv.Itoa(branch*10), nil, "", nil)
		if !assert.NoError(t, err) || !assert.NotNil(t, img) {
			return
		}
		images[branch*10] = img

		// level 3
		img, err = imageCache.WriteImage(op, img, strconv.Itoa(branch*100), nil, "", nil)
		if !assert.NoError(t, err) || !assert.NotNil(t, img) {
			return
		}
		images[branch*100] = img
	}

	// Add an extra node to the last branch
	img, err := imageCache.WriteImage(op, images[30], "301", nil, "", nil)
	if !assert.NoError(t, err) || !assert.NotNil(t, img) {
		return
	}
	images[301] = img

	//
	// Everything above here is just setup.  Everything from here on is the test.
	//

	// Deletion of an intermediate node should fail
	imagesDeleted, err := imageCache.DeleteBranch(op, images[1], nil)
	if !assert.Error(t, err) && assert.Nil(t, imagesDeleted) {
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

	//
	// Deletion of a branch
	//
	imagesDeleted, err = imageCache.DeleteBranch(op, images[100], nil)
	if !assert.NoError(t, err) {
		return
	}

	// List images should be missing a branch
	imageList, err = imageCache.ListImages(op, storeURL, nil)
	if !assert.NoError(t, err) || !assert.NotNil(t, imageList) {
		return
	}

	if !assert.Equal(t, 7, len(imageList)) || !assert.Equal(t, 3, len(imagesDeleted)) {
		return
	}

	//
	// Deletion of the split branch should only allow deletion of a single image
	//
	imagesDeleted, err = imageCache.DeleteBranch(op, images[300], nil)
	if !assert.NoError(t, err) {
		return
	}

	imageList, err = imageCache.ListImages(op, storeURL, nil)
	if !assert.NoError(t, err) || !assert.NotNil(t, imageList) {
		return
	}

	// only 300 should have been deleted
	if !assert.Equal(t, 6, len(imageList)) || !assert.Equal(t, images[300], imagesDeleted[0]) {
		return
	}

	//
	// Test keep with our 1 remaining branch
	//

	imagesDeleted, err = imageCache.DeleteBranch(op, images[200], []*url.URL{images[2].SelfLink})
	if !assert.NoError(t, err) {
		return
	}

	imageList, err = imageCache.ListImages(op, storeURL, nil)
	if !assert.NoError(t, err) || !assert.NotNil(t, imageList) {
		return
	}

	// only 20 and 200 should have been deleted
	if !assert.Equal(t, 4, len(imageList)) || !assert.Equal(t, images[200], imagesDeleted[0]) || !assert.Equal(t, images[20], imagesDeleted[1]) {
		for _, img = range imageList {
			t.Logf("image = %#v", img)
		}
		return
	}

}

func TestCreateImageStoreFailureCleanup(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	trace.Logger.Level = logrus.DebugLevel

	mds := mock.NewMockDataStore(nil)
	imageCache := NewLookupCache(mds)
	op := trace.NewOperation(context.Background(), "create image store error")
	mds.SetCreateImageStoreError(fmt.Errorf("foo error"))

	storeURL, err := imageCache.CreateImageStore(op, "testStore")
	if !assert.Error(t, err) || !assert.Nil(t, storeURL) {
		return
	}

	mds.SetCreateImageStoreError(nil)
	storeURL, err = imageCache.CreateImageStore(op, "testStore")
	if !assert.NoError(t, err) || !assert.NotNil(t, storeURL) {
		return
	}

	op = trace.NewOperation(context.Background(), "write image error")
	mds = mock.NewMockDataStore(nil)
	mds.SetWriteImageError(fmt.Errorf("foo error"))
	imageCache = NewLookupCache(mds)

	storeURL, err = imageCache.CreateImageStore(op, "testStore2")
	if !assert.Error(t, err) || !assert.Nil(t, storeURL) {
		return
	}

	mds.SetWriteImageError(nil)
	storeURL, err = imageCache.CreateImageStore(op, "testStore2")
	if !assert.NoError(t, err) || !assert.NotNil(t, storeURL) {
		return
	}
}

// Cache population should be happening in order starting from parent(id1) to children(id4)
func TestPopulateCacheInExpectedOrder(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	trace.EnableTracing()
	trace.Logger.Level = logrus.DebugLevel

	op := trace.NewOperation(context.Background(), "test")

	storeURL, _ := util.ImageStoreNameToURL("testStore")

	storageURLStr := storeURL.String()

	url1, _ := url.Parse(storageURLStr + "/id1")
	url2, _ := url.Parse(storageURLStr + "/id2")
	url3, _ := url.Parse(storageURLStr + "/id3")
	url4, _ := url.Parse(storageURLStr + "/id4")
	scratchURL, _ := url.Parse(storageURLStr + constants.ScratchLayerID)

	img1 := &image.Image{ID: "id1", SelfLink: url1, ParentLink: scratchURL, Store: storeURL}
	img2 := &image.Image{ID: "id2", SelfLink: url2, ParentLink: url1, Store: storeURL}
	img3 := &image.Image{ID: "id3", SelfLink: url3, ParentLink: url2, Store: storeURL}
	img4 := &image.Image{ID: "id4", SelfLink: url4, ParentLink: url3, Store: storeURL}
	scratchImg := &image.Image{
		ID:         constants.ScratchLayerID,
		SelfLink:   scratchURL,
		ParentLink: scratchURL,
		Store:      storeURL,
	}

	// Order does matter for some reason.
	imageMap := map[string]*image.Image{
		img1.ID:       img1,
		img4.ID:       img4,
		img2.ID:       img2,
		img3.ID:       img3,
		scratchImg.ID: scratchImg,
	}

	storeImages := map[url.URL]map[string]*image.Image{
		*storeURL: imageMap,
	}
	st := mock.NewMockDataStore(storeImages)

	imageCache := NewLookupCache(st)
	sURL, err := imageCache.GetImageStore(op, "testStore")
	require.NoError(t, err, "Received error while getting image store from cache")

	// Check if all images are available.
	imageIds := []string{"id1", "id2", "id3", "id4"}
	for _, imageID := range imageIds {
		v, err := imageCache.GetImage(op, sURL, imageID)
		require.NoError(t, err, "Received error instead of image from cache")
		assert.NotNil(t, v)
	}
}
