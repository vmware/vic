package storage

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
	portlayer "github.com/vmware/vic/portlayer/storage"
	"github.com/vmware/vic/portlayer/util"
)

var (
	// Set NoLchown to true if our effective uid is not 0.
	// This way we support unprivileged usage (aka non-root users) so
	// portlayer cannot nuke the system, by mistake.
	//
	// Also support dev + test cycle that doesn't require sudo or sticky bits.
	tarOptions = &archive.TarOptions{NoLchown: os.Geteuid() != 0}
)

// LocalStore implements the storage API on linux.  Stores images to the local filesystem.
type LocalStore struct {
	// root to store images in relative to /
	Path string
}

func (s *LocalStore) CreateImageStore(storeName string) (*url.URL, error) {
	u, err := util.StoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	if err := os.Mkdir(path.Join(s.Path, storeName), os.ModeDir|0744); err != nil {
		return nil, err
	}

	return u, nil
}

// GetImageStore checks to see if the image store exists on disk and returns an
// error or the store's URL.
func (s *LocalStore) GetImageStore(storeName string) (*url.URL, error) {
	u, err := util.StoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path.Join(s.Path, storeName))
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, os.ErrExist
	}
	return u, nil
}

func (s *LocalStore) ListImageStores() ([]*url.URL, error) {
	return nil, fmt.Errorf("not implemented")
}

// WriteImage creates a new image layer from the given parent.
// Eg parentImage + newLayer = new Image built from parent
//
// parent - The parent image to create the new image from.
// ID - textual ID for the image to be written
// Tag - the tag of the image to be written
func (s *LocalStore) WriteImage(parent *portlayer.Image, ID string, r io.Reader) (*portlayer.Image, error) {

	storeName, err := util.StoreName(parent.Store)
	if err != nil {
		return nil, err
	}

	imageURL, err := util.ImageURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	// XXX need to copy the parent!!!

	// We know the parent exists in this store because the cache checked for us.
	// XXX TODO sanitize paths
	destPath := filepath.Join(s.Path, storeName, ID)

	if err = os.Mkdir(destPath, os.ModeDir|0744); err != nil {
		return nil, err
	}

	// check if this is scratch
	if ID != portlayer.Scratch.ID {
		err = archive.Untar(r, destPath, tarOptions)
		if err != nil {
			return nil, err
		}
	}

	newImage := &portlayer.Image{
		ID:       ID,
		SelfLink: imageURL,
		Parent:   parent.SelfLink,
		Store:    parent.Store,
	}

	return newImage, nil
}

// GetImage queries the image store for the specified image.
//
// store - The image store to query
// name - The name of the image (optional)
// tag - The tagged version of the image (optional)
func (s *LocalStore) GetImage(store *url.URL, ID string) (*portlayer.Image, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (s *LocalStore) ListImages(store *url.URL, IDs []string) ([]*portlayer.Image, error) {
	return nil, fmt.Errorf("not yet implemented")
}
