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

// MountDataSink implements the DataSink interface for mounted devices
// This is a single use mechanism and will be tidied up on exit from MountDataSource.Import
type MountDataSink struct {
	Path  *os.File
	Clean func()
}

// Sink returns the data source associated with the DataSink
func (m *MountDataSink) Sink() interface{} {
	return m.Path
}

// Import writes `data` to the data source associated with this DataSource
// This will call MountDataSink.Close on exit, irrespective of success or error
func (m *MountDataSink) Import(op trace.Operation, spec *archive.FilterSpec, data io.ReadCloser) error {
	// ensure that mounts are tidied up - a data sink is a single use mechanism.
	defer m.Close()

	fi, err := m.Path.Stat()
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New("Path must be a directory")
	}

	// This assumes that m.Path was opened with a useful path (i.e. absolute) as that argument is what's
	// returned by Name.
	return archive.Unpack(op, data, spec, m.Path.Name())
}

func (m *MountDataSink) Close() error {
	m.Path.Close()
	m.Clean()

	return nil
}
