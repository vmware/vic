// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"errors"
	"io"
	"net/url"

	"github.com/vmware/vic/pkg/trace"
	"os"
)

type MountHandler interface {
	Mount(target *url.URL) (NFSTarget, error)
	Unmount(target NFSTarget) error
}

type NFSTarget interface {
	//opens a target path in a READONLY context
	Open(path string) (io.ReadCloser, error)

	//Opens targeted file with the supplied attr.
	OpenFile(path string, perm os.FileMode) (io.ReadWriteCloser, error)

	//Creates a file, errors out if file already exists. assumes write permissions.
	Create(path string, perm os.FileMode) (io.ReadWriteCloser, error)

	// Create directory path
	MkDir(path string, perm os.FileMode) ([]byte, error)

	// Delete Directory Path, and children
	RemoveAll(Path string) error

	//Reads the contents of the targeted directory
	ReadDir(path string) ([]os.FileInfo, error)

	//Looks up the file information for a target entry
	Lookup(path string) (os.FileInfo, error)
}

type NFSv3Target struct {
	//nfs endpoint
	host *url.URL

	//Path to the the volume store
	directoryPath string
}

func (t NFSv3Target) MkDir(path string, perm os.FileMode) ([]byte, error) {
	return nil, nil
}

func (t NFSv3Target) RemoveAll(path string) error {
	return nil
}

func (t NFSv3Target) ReadDir(path string) ([]os.FileInfo, error) {
	return nil, nil
}

func (t NFSv3Target) Open(path string) (io.ReadCloser, error) {
	return nil, nil
}

func (t NFSv3Target) OpenFile(path string, perm os.FileMode) (io.ReadWriteCloser, error) {
	return nil, nil
}

func (t NFSv3Target) Create(path string, perm os.FileMode) (io.ReadWriteCloser, error) {
	return nil, nil
}

func (t NFSv3Target) Lookup(path string) (os.FileInfo, error) {
	return nil, nil
}

func Mount(target *url.URL) (NFSTarget, error) {
	return nil, nil
}

func Unmount(target NFSTarget) error {
	return nil
}

// TODO: client implementation here
func NewNFSv3Target(op trace.Operation, fqdn *url.URL, volumeDirectory string) (NFSTarget, error) {
	return NFSv3Target{}, errors.New("FUNCTION NOT IMPLEMENTED")
}

//MAKE MOCK TARGET
//USER TEMP DIR
