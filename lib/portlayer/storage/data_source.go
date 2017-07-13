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
	Path    *os.File
	Clean   func()
	cleanOp trace.Operation
}

// NewMountDataSource creates a new data source assocaited with a specific mount, with the mount
// point being the path argument.
// The cleanup function is invoked with the Close of the ReadCloser from Export, or explicitly
func NewMountDataSource(op trace.Operation, path *os.File, cleanup func()) *MountDataSource {
	if path == nil {
		return nil
	}

	op.Debugf("Created mount data source at %s", path.Name())

	return &MountDataSource{
		Path:    path,
		Clean:   cleanup,
		cleanOp: trace.FromOperation(op, "clean up from new mount source"),
	}
}

// Source returns the data source associated with the DataSource
func (m *MountDataSource) Source() interface{} {
	return m.Path
}

// Export reads data from the associated data source and returns it as a tar archive
func (m *MountDataSource) Export(op trace.Operation, spec *archive.FilterSpec, data bool) (io.ReadCloser, error) {
	// reparent cleanup to Export operation
	m.cleanOp = trace.FromOperation(op, "clean up from export")

	name := m.Path.Name()
	fi, err := m.Path.Stat()
	if err != nil {
		op.Errorf("Unable to stat mount path %s for data source: %s", name, err)
		return nil, err
	}

	if !fi.IsDir() {
		return nil, errors.New("path must be a directory")
	}

	// NOTE: this isn't actually diffing - it's just creating a tar. @jzt to explain why
	op.Infof("Exporting data from %s", name)
	rc, err := archive.Diff(op, name, "", spec, data, m.Clean)
	if err != nil {
		return nil, err
	}

	return &ProxyReadCloser{
		rc,
		m.Close,
	}, nil
}

func (m *MountDataSource) Close() error {
	m.cleanOp.Infof("cleaning up after export")

	m.Path.Close()
	if m.Clean != nil {
		m.cleanOp.Debugf("calling specified cleaner function")
		m.Clean()
	}

	return nil
}
