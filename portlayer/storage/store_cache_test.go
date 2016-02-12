package storage

import (
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/portlayer/util"
)

type MockDataStore struct {
}

// GetImageStore checks to see if a named image store exists and returls the
// URL to it if so or error.
func (c *MockDataStore) GetImageStore(storeName string) (*url.URL, error) {
	return nil, nil
}

func (c *MockDataStore) CreateImageStore(storeName string) (*url.URL, error) {
	u, err := util.StoreNameToUrl(storeName)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (c *MockDataStore) WriteImage(parent *Image, ID string, r io.Reader) (*Image, error) {
	i := Image{
		ID:     ID,
		Store:  parent.Store,
		Parent: parent.SelfLink,
	}

	return &i, nil
}

// GetImage gets the specified image from the given store by retreiving it from the cache.
func (c *MockDataStore) GetImage(store *url.URL, ID string) (*Image, error) {
	return nil, nil
}

// ListImages resturns a list of Images for a list of IDs, or all if no IDs are passed
func (c *MockDataStore) ListImages(store *url.URL, IDs []string) ([]*Image, error) {
	return nil, nil
}

func TestListImages(t *testing.T) {
	s := &NameLookupCache{
		DataStore: &MockDataStore{},
	}

	storeUrl, err := s.CreateImageStore("testStore")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, storeUrl) {
		return
	}

	// Create a set of images
	images := make(map[string]*Image)
	images[Scratch.ID] = &Scratch
	parent := Scratch
	parent.Store = storeUrl
	for i := 1; i < 50; i++ {
		id := fmt.Sprintf("ID-%d", i)

		img, err := s.WriteImage(&parent, id, nil)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, img) {
			return
		}

		images[id] = img
	}

	// List all images
	outImages, err := s.ListImages(storeUrl, nil)
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
	outImages, err = s.ListImages(storeUrl, inIDs)
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
