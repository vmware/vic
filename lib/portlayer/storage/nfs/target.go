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

type NFSTarget interface {
	//returns the NFSTarget Endpoint
	EndPoint() *url.URL

	//Binds the target NFS server to the VCH(Possible concurrency/performance issue?)
	Open() (NFSTarget, error)

	//Unbinds NFS server from the VCH(ditto to open)
	Close() error

	// write data to target nfs filesystem
	Write(op trace.Operation, path string, r io.Reader) error

	// read data from target nfs filesystem
	Read(path string) (io.Reader, error)

	// Create directory path
	MkDir(path string, makeParents bool) (string, error)

	// Delete Directory Path, and children
	RemoveAll(Path string) error

	//Reads the contents of the targeted directory
	ReadDir(path string) ([]os.FileInfo, error)
}

type NFSv3Target struct {
	//nfs enpoint
	host *url.URL

	//Path to the the volume store
	directoryPath string
}

func (t NFSv3Target) EndPoint() *url.URL {
	return nil
}

func (t NFSv3Target) Write(op trace.Operation, path string, r io.Reader) error {
	return nil
}

func (t NFSv3Target) Read(path string) (io.Reader, error) {
	return nil, nil
}

func (t NFSv3Target) MkDir(path string, makeParents bool) (string, error) {
	return "", nil
}

func (t NFSv3Target) RemoveAll(path string) error {
	return nil
}

func (t NFSv3Target) ReadDir(path string) ([]os.FileInfo, error) {
	return nil, nil
}

func (t NFSv3Target) Open() (NFSTarget, error) {
	return nil, nil
}

func (t NFSv3Target) Close() error {
	return nil
}

// TODO: client implementation here
func NewNFSv3Target(op trace.Operation, fqdn *url.URL, volumeDirectory string) (NFSTarget, error) {
	return NFSv3Target{}, errors.New("FUNCTION NOT IMPLEMENTED")
}

//MAKE MOCK TARGET
//USER TEMP DIR
