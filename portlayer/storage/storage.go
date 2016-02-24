package storage

import (
	"io"
	"net/url"
)

// Interface to store images in the Image Store
type ImageStorer interface {

	// CreateImageStore creates a location to store images and creates a root
	// disk which serves as the parent of all layers.
	//
	// storeName - The name of the image store to be created.  This must be
	// unique.
	//
	// Returns the URL of the created store
	CreateImageStore(storeName string) (*url.URL, error)

	// Gets the url to an image store via name
	GetImageStore(storeName string) (*url.URL, error)

	// WriteImage creates a new image layer from the given parent.  Eg
	// parentImage + newLayer = new Image built from parent
	//
	// parent - The parent image to create the new image from.  ID - textual ID
	// for the image to be written Tag - the tag of the image to be written
	WriteImage(parent *Image, ID string, r io.Reader) (*Image,
		error)

	// GetImage queries the image store for the specified image.
	//
	// store - The image store to query name - The name of the image (optional)
	// tag - The tagged version of the image (optional)
	GetImage(store *url.URL, ID string) (*Image, error)

	// ListImages returns a list of Images given a list of image IDs, or all
	// images in the image store if no param is passed.
	ListImages(store *url.URL, IDs []string) ([]*Image, error)
}
