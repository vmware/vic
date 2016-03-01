package storage

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/vmware/vic/portlayer/util"
)

type Image struct {
	// Identifer for this layer.  Usually a SHA
	ID string

	// location of the layer.  Filled in by the runtime.
	SelfLink *url.URL

	// Parent's location.  It's the VMDK this snapshot inerits from.
	Parent *url.URL

	Store *url.URL
}

func Parse(u *url.URL) (*Image, error) {
	// Check the path isn't malformed.
	if !filepath.IsAbs(u.Path) {
		return nil, errors.New("invalid uri path")
	}

	segments := strings.Split(filepath.Clean(u.Path), "/")

	if segments[0] != util.StoragePath {
		return nil, errors.New("not a storage path")
	}

	if len(segments) < 3 {
		return nil, errors.New("uri path mismatch")
	}

	store, err := util.StoreNameToURL(segments[2])
	if err != nil {
		return nil, err
	}

	id := segments[3]

	var SelfLink url.URL
	SelfLink = *u

	i := &Image{
		ID:       id,
		SelfLink: &SelfLink,
		Store:    store,
	}

	return i, nil
}
