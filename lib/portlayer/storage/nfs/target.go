// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package nfs

import (
	"io"
	"net/url"
	"os"
)

// MountServer is an interface used to communicate with network attached storage.
type MountServer interface {
	// Mount initiates the NAS Target and returns a Target interface.
	Mount(target *url.URL) (Target, error)

	// Unmount terminates the Mount on the Target.
	Unmount(target Target) error
}

// Target is the filesystem interface for performing actions against attached storage.
type Target interface {
	// Open opens a file on the Target in RD_ONLY
	Open(path string) (io.ReadCloser, error)

	// OpenFile opens a file on the Target with the given mode
	OpenFile(path string, perm os.FileMode) (io.ReadWriteCloser, error)

	// Create creates a file, errors out if file already exists
	Create(path string, perm os.FileMode) (io.ReadWriteCloser, error)

	// Mkdir creates a directory at the given path
	Mkdir(path string, perm os.FileMode) ([]byte, error)

	// RemoveAll deletes Directory recursively
	RemoveAll(Path string) error

	// ReadDir reads the dirents of the given directory
	ReadDir(path string) ([]os.FileInfo, error)

	// Lookup reads os.FileInfo for the given path
	Lookup(path string) (os.FileInfo, error)
}
