// Copyright 2017-2018 VMware, Inc. All Rights Reserved.
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
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/storage"
	"github.com/vmware/vic/lib/constants"
	"github.com/vmware/vic/lib/portlayer/storage/image"
	icache "github.com/vmware/vic/lib/portlayer/storage/image/cache"
	imagemock "github.com/vmware/vic/lib/portlayer/storage/image/mock"
	vcache "github.com/vmware/vic/lib/portlayer/storage/volume/cache"
	volumemock "github.com/vmware/vic/lib/portlayer/storage/volume/mock"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/trace"
)

var (
	testImageID     = "testImage"
	testImageSum    = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	testHostName, _ = os.Hostname()
	testStoreName   = "testStore"
	testStoreURL    = url.URL{
		Scheme: "http",
		Host:   testHostName,
		Path:   "/" + util.ImageURLPath + "/" + testStoreName,
	}
)

func TestCreateImageStore(t *testing.T) {
	s := &StorageHandlersImpl{
		imageCache: icache.NewLookupCache(imagemock.NewMockDataStore(nil)),
	}

	store := &models.ImageStore{
		Name: "testStore",
	}

	params := storage.CreateImageStoreParams{
		Body: store,
	}

	result := s.CreateImageStore(params)
	if !assert.NotNil(t, result) {
		return
	}

	// try to recreate the same image store
	result = s.CreateImageStore(params)
	if !assert.NotNil(t, result) {
		return
	}
	// expect 409 since it already exists
	conflict := &storage.CreateImageStoreConflict{
		Payload: &models.Error{
			Code:    http.StatusConflict,
			Message: "An image store with that name already exists",
		},
	}
	if !assert.Equal(t, conflict, result) {
		return
	}
}

func TestGetImage(t *testing.T) {

	s := &StorageHandlersImpl{
		imageCache: icache.NewLookupCache(imagemock.NewMockDataStore(nil)),
	}

	params := &storage.GetImageParams{
		ID:        testImageID,
		StoreName: testStoreName,
	}

	result := s.GetImage(*params)
	if !assert.NotNil(t, result) {
		return
	}

	op := trace.NewOperation(context.Background(), "test")

	// create the image store
	url, err := s.imageCache.CreateImageStore(op, testStoreName)
	// TODO(jzt): these are testing NameLookupCache, do we need them here?
	if !assert.Nil(t, err, "Error while creating image store") {
		return
	}
	if !assert.Equal(t, testStoreURL.String(), url.String()) {
		return
	}

	// try GetImage again
	result = s.GetImage(*params)
	if !assert.NotNil(t, result) {
		return
	}

	// add image to store
	parent := image.Image{
		ID:         "scratch",
		SelfLink:   nil,
		ParentLink: nil,
		Store:      &testStoreURL,
	}

	expectedMeta := make(map[string][]byte)
	expectedMeta["foo"] = []byte("bar")
	// add the image to the store
	image, err := s.imageCache.WriteImage(op, &parent, testImageID, expectedMeta, testImageSum, nil)
	if !assert.NoError(t, err) || !assert.NotNil(t, image) {
		return
	}

	selflink, err := util.ImageURL(testStoreName, testImageID)
	if !assert.NoError(t, err) {
		return
	}
	sl := selflink.String()

	parentlink, err := util.ImageURL(testStoreName, parent.ID)
	if !assert.NoError(t, err) {
		return
	}
	p := parentlink.String()

	eMeta := make(map[string]string)
	eMeta["foo"] = "bar"
	// expect our image back now that we've created it
	expected := &storage.GetImageOK{
		Payload: &models.Image{
			ID:       image.ID,
			SelfLink: sl,
			Parent:   p,
			Store:    testStoreURL.String(),
			Metadata: eMeta,
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

	s := &StorageHandlersImpl{
		imageCache: icache.NewLookupCache(imagemock.NewMockDataStore(nil)),
	}

	params := &storage.ListImagesParams{
		StoreName: testStoreName,
	}

	outImages := s.ListImages(*params)
	if !assert.NotNil(t, outImages) {
		return
	}

	op := trace.NewOperation(context.Background(), "test")

	// create the image store
	url, err := s.imageCache.CreateImageStore(op, testStoreName)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, url) {
		return
	}

	// create a set of images
	images := make(map[string]*image.Image)
	parent := image.Image{
		ID: constants.ScratchLayerID,
	}
	parent.Store = &testStoreURL
	for i := 1; i < 50; i++ {
		id := fmt.Sprintf("id-%d", i)
		img, err := s.imageCache.WriteImage(op, &parent, id, nil, testImageSum, nil)
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

	// List specific images
	var ids []string

	// query for odd-numbered image ids
	for i := 1; i < 50; i += 2 {
		ids = append(ids, fmt.Sprintf("id-%d", i))
	}
	params.Ids = ids
	outImages = s.ListImages(*params)
	assert.IsType(t, &storage.ListImagesOK{}, outImages)
	assert.Equal(t, len(ids), len(outImages.(*storage.ListImagesOK).Payload))

	outmap := make(map[string]*models.Image)
	for _, image := range outImages.(*storage.ListImagesOK).Payload {
		outmap[image.ID] = image
	}

	// ensure no even-numbered image ids in our result
	for i := 2; i < 50; i += 2 {
		id := fmt.Sprintf("id-%d", i)
		_, ok := outmap[id]
		if !assert.False(t, ok) {
			return
		}
	}
}

func TestWriteImage(t *testing.T) {
	ic := icache.NewLookupCache(imagemock.NewMockDataStore(nil))

	// create image store
	op := trace.NewOperation(context.Background(), "test")
	_, err := ic.CreateImageStore(op, testStoreName)
	if err != nil {
		return
	}

	s := &StorageHandlersImpl{
		imageCache: ic,
	}

	eMeta := make(map[string]string)
	eMeta["foo"] = "bar"

	name := new(string)
	val := new(string)
	*name = "foo"
	*val = eMeta["foo"]

	params := &storage.WriteImageParams{
		StoreName:   testStoreName,
		ImageID:     testImageID,
		ParentID:    "scratch",
		Sum:         testImageSum,
		Metadatakey: name,
		Metadataval: val,
		ImageFile:   nil,
	}

	parentlink, err := util.ImageURL(testStoreName, params.ParentID)
	if !assert.NoError(t, err) {
		return
	}
	p := parentlink.String()

	selflink, err := util.ImageURL(testStoreName, testImageID)
	if !assert.NoError(t, err) {
		return
	}
	sl := selflink.String()

	expected := &storage.WriteImageCreated{
		Payload: &models.Image{
			ID:       testImageID,
			Parent:   p,
			SelfLink: sl,
			Store:    testStoreURL.String(),
			Metadata: eMeta,
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

func TestVolumeCreate(t *testing.T) {

	op := trace.NewOperation(context.Background(), "test")
	volCache := vcache.NewVolumeLookupCache(op)

	testStore := volumemock.NewMockVolumeStore()
	_, err := volCache.AddStore(op, "testStore", testStore)
	if !assert.NoError(t, err) {
		return
	}

	handler := StorageHandlersImpl{
		volumeCache: volCache,
	}

	model := models.VolumeRequest{}
	model.Store = "testStore"
	model.Name = "testVolume"
	model.Capacity = 1
	model.Driver = "vsphere"
	model.DriverArgs = make(map[string]string)
	model.DriverArgs["stuff"] = "things"
	model.Metadata = make(map[string]string)
	params := storage.NewCreateVolumeParams()
	params.VolumeRequest = &model

	handler.CreateVolume(params)
	testVolume, err := testStore.VolumeGet(op, model.Name)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.NotNil(t, testVolume) {
		return
	}
	testVolumeStoreName, err := util.VolumeStoreName(testVolume.Store)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, "testStore", testVolumeStoreName) {
		return
	}
	if !assert.Equal(t, "testVolume", testVolume.ID) {
		return
	}
}

func TestParseUIDAndGID(t *testing.T) {
	positiveCases := []url.URL{
		{
			Scheme:   "nfs",
			Host:     "testURL",
			RawQuery: "uid=1234&gid=5678",
			Path:     "/test/path",
		},
		{
			Scheme:   "nfs",
			Host:     "testURL",
			RawQuery: "uid=00000000000000&gid=00000000000000000000000000",
			Path:     "/test/path",
		},
		{
			Scheme:   "nfs",
			Host:     "testURL",
			RawQuery: "uid=321321321&gid=123123123",
			Path:     "/test/path",
		},
		{
			Scheme:   "nfs",
			Host:     "testURL",
			RawQuery: "uid=0&gid=0",
			Path:     "/test/path",
		},
		{
			Scheme:   "nfs",
			Host:     "testURL",
			RawQuery: "uid=&gid=",
			Path:     "/test/path",
		},
	}

	negativeCases := []url.URL{
		{
			Scheme:   "nfs",
			Host:     "testURL",
			RawQuery: "uid=Hello&gid=World",
			Path:     "/test/path",
		},
		{
			Scheme:   "nfs",
			Host:     "testURL",
			RawQuery: "uid=ASKJHG#!@#LJK$&gid=!@#$$%#@@!",
			Path:     "/test/path",
		},
		{
			Scheme:   "nfs",
			Host:     "testURL",
			RawQuery: "uid=9999999999999999999999999999999999999999999999999&gid=7777777777777777777777777777777777777777777777777777777",
			Path:     "/test/path",
		},
	}

	for _, v := range positiveCases {

		testUID, testGID, err := parseUIDAndGID(&v)
		assert.Nil(t, err, v.String())
		assert.NotEqual(t, -1, testUID, v.String())
		assert.NotEqual(t, -1, testGID, v.String())
	}

	for _, v := range negativeCases {
		testUID, testGID, err := parseUIDAndGID(&v)
		assert.NotNil(t, err, v.String())
		assert.Equal(t, -1, testUID, v.String())
		assert.Equal(t, -1, testGID, v.String())

	}

}
