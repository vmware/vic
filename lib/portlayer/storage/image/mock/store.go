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

package mock

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/vic/lib/portlayer/storage/image"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

type MockDataStore struct {
	object.DatastorePath

	// id -> image
	db map[url.URL]map[string]*image.Image

	createImageStoreError error
	writeImageError       error
}

func NewMockDataStore(images map[url.URL]map[string]*image.Image) *MockDataStore {
	m := &MockDataStore{
		db: images,
	}

	if m.db == nil {
		m.db = make(map[url.URL]map[string]*image.Image)
	}

	return m
}

// SetWriteImageError injects an error to be returned on next call to WriteImage
func (c *MockDataStore) SetWriteImageError(err error) {
	c.writeImageError = err
}

// SetCreateImageStoreError injects an error to be returned on next call to CreateImageStore
func (c *MockDataStore) SetCreateImageStoreError(err error) {
	c.createImageStoreError = err
}

// GetImageStore checks to see if a named image store exists and returls the
// URL to it if so or error.
func (c *MockDataStore) GetImageStore(op trace.Operation, storeName string) (*url.URL, error) {
	u, err := util.ImageStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	if _, ok := c.db[*u]; !ok {
		return nil, os.ErrNotExist
	}

	return u, nil
}

func (c *MockDataStore) CreateImageStore(op trace.Operation, storeName string) (*url.URL, error) {
	if c.createImageStoreError != nil {
		return nil, c.createImageStoreError
	}

	u, err := util.ImageStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	c.db[*u] = make(map[string]*image.Image)
	return u, nil
}

func (c *MockDataStore) DeleteImageStore(op trace.Operation, storeName string) error {
	u, err := util.ImageStoreNameToURL(storeName)
	if err != nil {
		return err
	}

	c.db[*u] = nil
	return nil
}

func (c *MockDataStore) ListImageStores(op trace.Operation) ([]*url.URL, error) {
	stores := make([]*url.URL, len(c.db))
	i := 0
	for k := range c.db {
		stores[i] = &k
		i++
	}

	return stores, nil
}

func (c *MockDataStore) URL(op trace.Operation, id string) (*url.URL, error) {
	return nil, nil
}

func (c *MockDataStore) Owners(op trace.Operation, url *url.URL, filter func(vm *mo.VirtualMachine) bool) ([]*vm.VirtualMachine, error) {
	return nil, nil
}

func (c *MockDataStore) WriteImage(op trace.Operation, parent *image.Image, ID string, meta map[string][]byte, sum string, r io.Reader) (*image.Image, error) {
	if c.writeImageError != nil {
		op.Infof("WriteImage: returning error")
		return nil, c.writeImageError
	}

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

	i := &image.Image{
		ID:         ID,
		Store:      parent.Store,
		ParentLink: parentLink,
		SelfLink:   selflink,
		Metadata:   meta,
		DatastorePath: &object.DatastorePath{
			Datastore: c.Datastore,
			Path:      c.Path,
		},
	}

	c.db[*parent.Store][ID] = i

	return i, nil
}

func (c *MockDataStore) WriteMetadata(op trace.Operation, storeName string, ID string, meta map[string][]byte) error {
	return nil
}

// GetImage gets the specified image from the given store by retreiving it from the cache.
func (c *MockDataStore) GetImage(op trace.Operation, store *url.URL, ID string) (*image.Image, error) {
	i, ok := c.db[*store][ID]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return i, nil
}

// ListImages resturns a list of Images for a list of IDs, or all if no IDs are passed
func (c *MockDataStore) ListImages(op trace.Operation, store *url.URL, IDs []string) ([]*image.Image, error) {
	var imageList []*image.Image
	for _, i := range c.db[*store] {
		imageList = append(imageList, i)
	}
	return imageList, nil
}

// DeleteImage removes an image from the image store
func (c *MockDataStore) DeleteImage(op trace.Operation, image *image.Image) (*image.Image, error) {
	delete(c.db[*image.Store], image.ID)
	return image, nil
}
