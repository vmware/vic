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
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/storage"
	spl "github.com/vmware/vic/lib/portlayer/storage"
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

type MockDataStore struct {
}

type MockVolumeStore struct {
	// id -> volume
	db map[string]*spl.Volume
}

func NewMockVolumeStore() *MockVolumeStore {
	m := &MockVolumeStore{
		db: make(map[string]*spl.Volume),
	}

	return m
}

// Creates a volume on the given volume store, of the given size, with the given metadata.
func (m *MockVolumeStore) VolumeCreate(op trace.Operation, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*spl.Volume, error) {
	storeName, err := util.VolumeStoreName(store)
	if err != nil {
		return nil, err
	}

	selfLink, err := util.VolumeURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	vol := &spl.Volume{
		ID:       ID,
		Store:    store,
		SelfLink: selfLink,
	}

	m.db[ID] = vol

	return vol, nil
}

// Get an existing volume via it's ID and volume store.
func (m *MockVolumeStore) VolumeGet(op trace.Operation, ID string) (*spl.Volume, error) {
	vol, ok := m.db[ID]
	if !ok {
		return nil, os.ErrNotExist
	}

	return vol, nil
}

// Destroys a volume
func (m *MockVolumeStore) VolumeDestroy(op trace.Operation, vol *spl.Volume) error {
	if _, ok := m.db[vol.ID]; !ok {
		return os.ErrNotExist
	}

	delete(m.db, vol.ID)

	return nil
}

func (m *MockVolumeStore) VolumeStoresList(op trace.Operation) (map[string]url.URL, error) {
	return nil, fmt.Errorf("not implemented")
}

// Lists all volumes on the given volume store`
func (m *MockVolumeStore) VolumesList(op trace.Operation) ([]*spl.Volume, error) {
	var i int
	list := make([]*spl.Volume, len(m.db))
	for _, v := range m.db {
		t := *v
		list[i] = &t
		i++
	}

	return list, nil
}

// GetImageStore checks to see if a named image store exists and returls the
// URL to it if so or error.
func (c *MockDataStore) GetImageStore(op trace.Operation, storeName string) (*url.URL, error) {
	u, err := util.ImageStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("store (%s) doesn't exist", u.String())
}

func (c *MockDataStore) CreateImageStore(op trace.Operation, storeName string) (*url.URL, error) {
	u, err := util.ImageStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (c *MockDataStore) DeleteImageStore(op trace.Operation, storeName string) error {
	return nil
}

func (c *MockDataStore) ListImageStores(op trace.Operation) ([]*url.URL, error) {
	return nil, nil
}

func (c *MockDataStore) WriteImage(op trace.Operation, parent *spl.Image, ID string, meta map[string][]byte, sum string, r io.Reader) (*spl.Image, error) {
	storeName, err := util.ImageStoreName(parent.Store)
	if err != nil {
		return nil, err
	}

	selflink, err := util.ImageURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	i := spl.Image{
		ID:         ID,
		Store:      parent.Store,
		ParentLink: parent.SelfLink,
		SelfLink:   selflink,
		Metadata:   meta,
	}

	return &i, nil
}
func (c *MockDataStore) WriteMetadata(op trace.Operation, storeName string, ID string, meta map[string][]byte) error {
	return nil
}

// GetImage gets the specified image from the given store by retreiving it from the cache.
func (c *MockDataStore) GetImage(op trace.Operation, store *url.URL, ID string) (*spl.Image, error) {
	if ID == spl.Scratch.ID {
		return &spl.Image{Store: store}, nil
	}

	return nil, fmt.Errorf("store (%s) doesn't have image %s", store.String(), ID)
}

// ListImages resturns a list of Images for a list of IDs, or all if no IDs are passed
func (c *MockDataStore) ListImages(op trace.Operation, store *url.URL, IDs []string) ([]*spl.Image, error) {
	return nil, fmt.Errorf("store (%s) doesn't exist", store.String())
}

func (c *MockDataStore) DeleteImage(op trace.Operation, image *spl.Image) (*spl.Image, error) {
	return nil, nil
}

func TestCreateImageStore(t *testing.T) {
	s := &StorageHandlersImpl{
		imageCache: spl.NewLookupCache(&MockDataStore{}),
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
		imageCache: spl.NewLookupCache(&MockDataStore{}),
	}

	params := &storage.GetImageParams{
		ID:        testImageID,
		StoreName: testStoreName,
	}

	// expect 404 since no image store exists by that name
	storeNotFound := &storage.GetImageNotFound{
		Payload: &models.Error{
			Code:    http.StatusNotFound,
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

	// expect 404 since no image exists by that name in that store
	imageNotFound := &storage.GetImageNotFound{
		Payload: &models.Error{
			Code:    http.StatusNotFound,
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
		imageCache: spl.NewLookupCache(&MockDataStore{}),
	}

	params := &storage.ListImagesParams{
		StoreName: testStoreName,
	}

	// expect 404 if image store doesn't exist
	notFound := &storage.ListImagesNotFound{
		Payload: &models.Error{
			Code:    http.StatusNotFound,
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
	images := make(map[string]*spl.Image)
	parent := spl.Scratch
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
	ic := spl.NewLookupCache(&MockDataStore{})

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

	testStore := NewMockVolumeStore()
	op := trace.NewOperation(context.Background(), "test")
	volCache, err := spl.NewVolumeLookupCache(op, testStore)
	if !assert.NoError(t, err) {
		return
	}

	handler := StorageHandlersImpl{
		volumeCache: volCache,
	}

	model := models.VolumeRequest{}
	model.Store = "blah"
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
	if !assert.Equal(t, "blah", testVolumeStoreName) {
		return
	}
	if !assert.Equal(t, "testVolume", testVolume.ID) {
		return
	}
}
