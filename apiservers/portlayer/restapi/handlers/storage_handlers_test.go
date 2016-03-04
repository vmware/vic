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

package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	//"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/swag"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/apiservers/portlayer/models"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations/storage"
	spl "github.com/vmware/vic/portlayer/storage"
	"github.com/vmware/vic/portlayer/util"
)

var (
	testImageID     = "testImage"
	testImageSum    = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	testHostName, _ = os.Hostname()
	testStoreName   = "testStore"
	testStoreURL    = url.URL{
		Scheme: "http",
		Host:   testHostName,
		Path:   "/storage/" + testStoreName,
	}
)

type MockDataStore struct {
}

// GetImageStore checks to see if a named image store exists and returls the
// URL to it if so or error.
func (c *MockDataStore) GetImageStore(storeName string) (*url.URL, error) {
	return nil, nil
}

func (c *MockDataStore) CreateImageStore(storeName string) (*url.URL, error) {
	u, err := util.StoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (c *MockDataStore) ListImageStores() ([]*url.URL, error) {
	return nil, nil
}

func (c *MockDataStore) WriteImage(parent *spl.Image, ID string, r io.Reader) (*spl.Image, error) {
	i := spl.Image{
		ID:     ID,
		Store:  parent.Store,
		Parent: parent.SelfLink,
	}

	return &i, nil
}

// GetImage gets the specified image from the given store by retreiving it from the cache.
func (c *MockDataStore) GetImage(store *url.URL, ID string) (*spl.Image, error) {
	return nil, nil
}

// ListImages resturns a list of Images for a list of IDs, or all if no IDs are passed
func (c *MockDataStore) ListImages(store *url.URL, IDs []string) ([]*spl.Image, error) {
	return nil, nil
}

func TestCreateImageStore(t *testing.T) {
	cache = &spl.NameLookupCache{
		DataStore: &MockDataStore{},
	}

	s := &StorageHandlersImpl{}
	store := &models.ImageStore{
		Name: "testStore",
	}

	params := &storage.CreateImageStoreParams{
		Body: store,
	}

	result := s.CreateImageStore(*params)
	if !assert.NotNil(t, result) {
		return
	}

	// try to recreate the same image store
	result = s.CreateImageStore(*params)
	if !assert.NotNil(t, result) {
		return
	}

	// expect 409 since it already exists
	conflict := &storage.CreateImageStoreConflict{
		Payload: &models.Error{
			Code:    swag.Int64(http.StatusConflict),
			Message: "An image store with that name already exists",
		},
	}
	if !assert.Equal(t, conflict, result) {
		return
	}
}

func TestGetImage(t *testing.T) {
	cache = &spl.NameLookupCache{
		DataStore: &MockDataStore{},
	}

	s := &StorageHandlersImpl{}

	params := &storage.GetImageParams{
		ID:        testImageID,
		StoreName: testStoreName,
	}

	// expect 404 since no image store exists by that name
	storeNotFound := &storage.GetImageNotFound{
		Payload: &models.Error{
			Code:    swag.Int64(http.StatusNotFound),
			Message: fmt.Sprintf("store (%s) doesn't exist", testStoreURL.String()),
		},
	}

	result := s.GetImage(*params)
	if !assert.NotNil(t, result) {
		return
	}
	if !assert.Equal(t, storeNotFound, result) {
		return
	}

	// create the image store
	url, err := cache.CreateImageStore(testStoreName)
	// TODO(jzt): these are testing NameLookupCache, do we need them here?
	if !assert.Nil(t, err, "Error while creating image store") {
		return
	}
	if !assert.Equal(t, &testStoreURL, url) {
		return
	}

	// expect 404 since no image exists by that name in that store
	imageNotFound := &storage.GetImageNotFound{
		Payload: &models.Error{
			Code:    swag.Int64(http.StatusNotFound),
			Message: fmt.Sprintf("store (%s) doesn't have image %s", testStoreURL.String(), testImageID),
		},
	}
	// try GetImage again
	result = s.GetImage(*params)
	if !assert.NotNil(t, result) {
		return
	}
	if !assert.Equal(t, imageNotFound, result) {
		return
	}

	// add image to store
	parent := spl.Image{
		ID:       "scratch",
		SelfLink: nil,
		Parent:   nil,
		Store:    &testStoreURL,
	}

	// add the image to the store
	image, err := cache.WriteImage(&parent, testImageID, testImageSum, nil)
	if !assert.NotNil(t, image) {
		return
	}

	// expect our image back now that we've created it
	expected := &storage.GetImageOK{
		Payload: &models.Image{
			ID:       image.ID,
			SelfLink: nil,
			Parent:   nil,
			Store:    testStoreURL.String(),
		},
	}

	result = s.GetImage(*params)
	if !assert.NotNil(t, result) {
		return
	}
	if !assert.Equal(t, expected, result) {
		return
	}
}

func TestListImages(t *testing.T) {
	cache = &spl.NameLookupCache{
		DataStore: &MockDataStore{},
	}

	s := &StorageHandlersImpl{}

	params := &storage.ListImagesParams{
		StoreName: testStoreName,
	}

	// expect 404 if image store doesn't exist
	notFound := &storage.ListImagesNotFound{
		Payload: &models.Error{
			Code:    swag.Int64(http.StatusNotFound),
			Message: fmt.Sprintf("store (%s) doesn't exist", testStoreURL.String()),
		},
	}

	outImages := s.ListImages(*params)
	if !assert.NotNil(t, outImages) {
		return
	}
	if !assert.Equal(t, notFound, outImages) {
		return
	}

	// create the image store
	url, err := cache.CreateImageStore(testStoreName)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, url) {
		return
	}

	// create a set of images
	images := make(map[string]*spl.Image)
	images[spl.Scratch.ID] = &spl.Scratch
	parent := spl.Scratch
	parent.Store = &testStoreURL
	for i := 1; i < 50; i++ {
		id := fmt.Sprintf("id-%d", i)
		img, err := cache.WriteImage(&parent, id, testImageSum, nil)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, img) {
			return
		}
		images[id] = img
	}

	// List all images
	outImages = s.ListImages(*params)
	assert.IsType(t, &storage.ListImagesOK{}, outImages)
	assert.Equal(t, len(outImages.(*storage.ListImagesOK).Payload), len(images))

	for _, img := range outImages.(*storage.ListImagesOK).Payload {
		_, ok := images[img.ID]
		if !assert.True(t, ok) {
			return
		}
	}
}

func TestWriteImage(t *testing.T) {

	cache = &spl.NameLookupCache{
		DataStore: &MockDataStore{},
	}

	// create image store
	_, err := cache.CreateImageStore(testStoreName)
	if err != nil {
		return
	}

	s := &StorageHandlersImpl{}

	params := &storage.WriteImageParams{
		StoreName: testStoreName,
		ImageID:   testImageID,
		ParentID:  "scratch",
		Sum:       testImageSum,
		ImageFile: nil,
	}

	expected := &storage.WriteImageCreated{
		Payload: &models.Image{
			ID:       testImageID,
			Parent:   nil,
			Store:    testStoreURL.String(),
			SelfLink: nil,
		},
	}

	result := s.WriteImage(*params)
	if !assert.NotNil(t, result) {
		return
	}
	if !assert.Equal(t, expected, result) {
		return
	}
}
