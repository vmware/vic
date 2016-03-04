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

package util

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	StoragePath = "storage/"
)

// StoreNameToURL parses the image URL in the form /storage/<image store>/<image name>
func StoreNameToURL(storeName string) (*url.URL, error) {
	return ServiceURL(StoragePath).Parse(storeName)
}

func StoreName(u *url.URL) (string, error) {
	// Check the path isn't malformed.
	if !filepath.IsAbs(u.Path) {
		return "", errors.New("invalid uri path")
	}

	segments := strings.Split(filepath.Clean(u.Path), "/")[1:]

	if segments[0] != filepath.Clean(StoragePath) {
		return "", errors.New("not a storage path")
	}

	if len(segments) < 2 {
		return "", errors.New("uri path mismatch")
	}

	return segments[1], nil
}

func ImageURL(storeName, imageName string) (*url.URL, error) {
	u, err := StoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	return u.Parse(imageName)
}
