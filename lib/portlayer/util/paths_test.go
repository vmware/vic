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
	"net/url"
	"testing"
)

func TestStoreName(t *testing.T) {
	u, _ := url.Parse("/storage/imgstore/image")
	store, err := StoreName(u)
	if err != nil {
		t.Errorf("StoreName failed %v", err)
	}
	expectedStore := "imgstore"
	if store != expectedStore {
		t.Errorf("Got: %s Expected: %s", store, expectedStore)
	}
}

func TestStoreNameErrors(t *testing.T) {
	u, _ := url.Parse("fail")
	_, err := StoreName(u)
	expectedError := "invalid uri path"
	if err.Error() != expectedError {
		t.Errorf("Got: %s Expected: %s", err, expectedError)
	}

	u, _ = url.Parse("/storage:123")
	_, err = StoreName(u)
	expectedError = "not a storage path"
	if err.Error() != expectedError {
		t.Errorf("Got: %s Expected: %s", err, expectedError)
	}

	u, _ = url.Parse("/storage")
	_, err = StoreName(u)
	expectedError = "uri path mismatch"
	if err.Error() != expectedError {
		t.Errorf("Got: %s Expected: %s", err, expectedError)
	}
}

func TestImageURL(t *testing.T) {
	DefaultHost, _ = url.Parse("http://foo.com/")
	storeName := "storage"
	imageName := "image"

	u, err := ImageURL(storeName, imageName)
	if err != nil {
		t.Errorf("ImageURL failed %v", err)
	}
	expectedURL := "http://foo.com/storage/image"
	if u.String() != expectedURL {
		t.Errorf("Got: %s Expected: %s", u, expectedURL)
	}
}
