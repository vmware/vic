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

// The MountServer is an interface used to communicate with network attached storage.
type MountServer interface {
	// Mount initiates the NAS Target and returns a Target interface.
	Mount(target *url.URL) (target, error)

	// Unmount terminates the Mount on the Target.
	Unmount(target target) error
}

// Target is the filesystem interface for performing actions against attached storage.:w

type target interface {
	// Opens a target path in a READONLY context
	Open(path string) (io.ReadCloser, error)

	// Opens targeted file with the supplied attr.
	OpenFile(path string, perm os.FileMode) (io.ReadWriteCloser, error)

	// Creates a file, errors out if file already exists. assumes write permissions.
	Create(path string, perm os.FileMode) (io.ReadWriteCloser, error)

	// Create directory path
	Mkdir(path string, perm os.FileMode) ([]byte, error)

	// Delete Directory Path, and children
	RemoveAll(Path string) error

	// Reads the contents of the targeted directory
	ReadDir(path string) ([]os.FileInfo, error)

	// Looks up the file information for a target entry
	Lookup(path string) (os.FileInfo, error)
}
