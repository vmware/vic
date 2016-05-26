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
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/vmware/vic/lib/portlayer/util"
)

type Image struct {
	// Identifer for this layer.  Usually a SHA
	ID string

	// location of the layer.  Filled in by the runtime.
	SelfLink *url.URL

	// Parent's location.  It's the VMDK this snapshot inerits from.
	Parent *url.URL

	Store *url.URL

	// Metadata associated with the image.
	Metadata map[string][]byte
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
