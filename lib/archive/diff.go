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

package archive

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	docker "github.com/docker/docker/pkg/archive"

	"github.com/vmware/vic/pkg/trace"
)

const (
	// ChangeTypeKey defines the key for the type of diff change stored in the tar Xattrs header
	ChangeTypeKey = "change_type"
)

// for sort.Sort
type changesByPath []docker.Change

func (c changesByPath) Less(i, j int) bool { return c[i].Path < c[j].Path }
func (c changesByPath) Len() int           { return len(c) }
func (c changesByPath) Swap(i, j int)      { c[j], c[i] = c[i], c[j] }

// Diff produces a tar archive containing the differences between two filesystems
func Diff(op trace.Operation, newDir, oldDir string, spec *FilterSpec, data bool) (io.ReadCloser, error) {

	var err error
	if spec == nil {
		spec, err = CreateFilterSpec(op, nil)
		if err != nil {
			return nil, err
		}
	}

	changes, err := docker.ChangesDirs(newDir, oldDir)
	if err != nil {
		return nil, err
	}

	sort.Sort(changesByPath(changes))

	return Tar(op, newDir, changes, spec, data)
}

func Tar(op trace.Operation, dir string, changes []docker.Change, spec *FilterSpec, data bool) (io.ReadCloser, error) {
	var (
		err error
		hdr *tar.Header
	)

	// Note: it is up to the caller to handle errors and the closing of the read side of the pipe
	r, w := io.Pipe()
	go func() {
		tw := tar.NewWriter(w)
		defer func() {
			var cerr error
			if cerr = tw.Close(); cerr != nil {
				op.Errorf("Error closing tar writer: %s", cerr.Error())
			}
			if err == nil {
				_ = w.CloseWithError(cerr)
			} else {
				_ = w.CloseWithError(err)
			}
			return
		}()

		for _, change := range changes {
			if spec.Excludes(op, change.Path) {
				continue
			}

			hdr, err = createHeader(op, dir, change, spec)
			if err != nil {
				op.Errorf("Error creating header from change: %s", err.Error())
				return
			}

			var f *os.File

			if !data {
				hdr.Size = 0
			}

			_ = tw.WriteHeader(hdr)

			p := filepath.Join(dir, change.Path)
			if hdr.Typeflag == tar.TypeReg && hdr.Size != 0 {
				f, err = os.Open(p)
				if err != nil {
					if os.IsPermission(err) {
						err = nil
					}
					return
				}
				if f != nil {
					_, err = io.Copy(tw, f)
					if err != nil {
						op.Errorf("Error writing archive data: %s", err.Error())
					}
					_ = f.Close()
				}
			}
		}
	}()
	return r, err
}

func createHeader(op trace.Operation, dir string, change docker.Change, spec *FilterSpec) (*tar.Header, error) {
	var hdr *tar.Header
	timestamp := time.Now()

	switch change.Kind {
	case docker.ChangeDelete:
		whiteOutDir := filepath.Dir(change.Path)
		whiteOutBase := filepath.Base(change.Path)
		whiteOut := filepath.Join(whiteOutDir, docker.WhiteoutPrefix+whiteOutBase)
		hdr = &tar.Header{
			Name:       filepath.Join(spec.RebasePath, whiteOut),
			ModTime:    timestamp,
			AccessTime: timestamp,
			ChangeTime: timestamp,
		}
	default:
		fi, err := os.Stat(filepath.Join(dir, change.Path))
		if err != nil {
			op.Errorf("Error getting file info: %s", err.Error())
			return nil, err
		}

		hdr, err = tar.FileInfoHeader(fi, change.Path)
		if err != nil {
			op.Errorf("Error getting file info header: %s", err.Error())
			return nil, err
		}

		hdr.Name = filepath.Join(spec.RebasePath, change.Path)

		if hdr.Typeflag == tar.TypeDir {
			hdr.Name += "/"
		}
	}

	hdr.Xattrs = make(map[string]string)
	hdr.Xattrs[ChangeTypeKey] = change.Kind.String()

	return hdr, nil
}
