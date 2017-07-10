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

package storage

import (
	"errors"
	"io"
	"os"

	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/pkg/trace"
)

// MountDataSource implements the DataSource interface for mounted devices
type MountDataSource struct {
	Path  *os.File
	Clean func()
}

// Source returns the data source associated with the DataSource
func (m *MountDataSource) Source() interface{} {
	return m.Path
}

// Export reads data from the associated data source and returns it as a tar archive
func (m *MountDataSource) Export(op trace.Operation, spec *archive.FilterSpec, data bool) (io.ReadCloser, error) {
	fi, err := m.Path.Stat()
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, errors.New("path must be a directory")
	}

	return archive.Diff(op, m.Path.Name(), "", spec, data)
}

func (m *MountDataSource) Close() error {
	m.Clean()

	return nil
}
