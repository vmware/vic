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

package archive

import (
	"archive/tar"
	"bytes"
)

// Below are some structs and utility functions for unit testing the untar functionality

type tarFile struct {
	Name string
	Type byte
	Body string
}

func tarFiles(files []tarFile) (*bytes.Buffer, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new tar archive.
	tw := tar.NewWriter(buf)

	// Write data to the tar as if it came from the hub
	for _, file := range files {

		var hdr *tar.Header

		switch file.Type {
		case tar.TypeDir:
			hdr = &tar.Header{
				Name:     file.Name,
				Mode:     0777,
				Typeflag: file.Type,
				Size:     0,
			}
		case tar.TypeSymlink:
			hdr = &tar.Header{
				Name:     file.Name,
				Mode:     0777,
				Typeflag: file.Type,
				Size:     0,
				Linkname: file.Body,
			}
		default:
			hdr = &tar.Header{
				Name:     file.Name,
				Mode:     0777,
				Typeflag: file.Type,
				Size:     int64(len(file.Body)),
			}
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}

		if file.Type == tar.TypeDir || file.Type == tar.TypeSymlink {
			continue
		}

		if _, err := tw.Write([]byte(file.Body)); err != nil {
			return nil, err
		}
	}

	// Make sure to check the error on Close.
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
