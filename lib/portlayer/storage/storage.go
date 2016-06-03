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
	"io"
	"net/url"

	"golang.org/x/net/context"
)

// ImageStorer is an interface to store images in the Image Store
type ImageStorer interface {

	// CreateImageStore creates a location to store images and creates a root
	// disk which serves as the parent of all layers.
	//
	// storeName - The name of the image store to be created.  This must be
	// unique.
	//
	// Returns the URL of the created store
	CreateImageStore(ctx context.Context, storeName string) (*url.URL, error)

	// Gets the url to an image store via name
	GetImageStore(ctx context.Context, storeName string) (*url.URL, error)

	// ListImageStores lists the available image stores
	ListImageStores(ctx context.Context) ([]*url.URL, error)

	// WriteImage creates a new image layer from the given parent.  Eg
	// parentImage + newLayer = new Image built from parent
	//
	// parent - The parent image to create the new image from.
	// ID - textual ID for the image to be written
	// meta - metadata associated with the image
	// r - the image tar to be written
	WriteImage(ctx context.Context, parent *Image, ID string, meta map[string][]byte, r io.Reader) (*Image,
		error)

	// GetImage queries the image store for the specified image.
	//
	// store - The image store to query name - The name of the image (optional)
	// ID - textual ID for the image to be retrieved
	GetImage(ctx context.Context, store *url.URL, ID string) (*Image, error)

	// ListImages returns a list of Images given a list of image IDs, or all
	// images in the image store if no param is passed.
	ListImages(ctx context.Context, store *url.URL, IDs []string) ([]*Image, error)
}
