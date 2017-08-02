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
	"strings"
	"sync"
	"time"

	docker "github.com/docker/docker/pkg/archive"

	"github.com/vmware/vic/pkg/trace"
)

const (
	// ChangeTypeKey defines the key for the type of diff change stored in the tar Xattrs header
	ChangeTypeKey = "change_type"
)

// CancelNotifyKey allows for a notification when cancelation is complete
type CancelNotifyKey struct{}

// for sort.Sort
type changesByPath []docker.Change

func (c changesByPath) Less(i, j int) bool { return c[i].Path < c[j].Path }
func (c changesByPath) Len() int           { return len(c) }
func (c changesByPath) Swap(i, j int)      { c[j], c[i] = c[i], c[j] }

// Diff produces a tar archive containing the differences between two filesystems
func Diff(op trace.Operation, newDir, oldDir string, spec *FilterSpec, data bool, xattr bool) (io.ReadCloser, error) {
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

	return Tar(op, newDir, changes, spec, data, xattr)
}

func Tar(op trace.Operation, dir string, changes []docker.Change, spec *FilterSpec, data bool, xattr bool) (io.ReadCloser, error) {
	var (
		err error
		hdr *tar.Header
	)

	// I don't like this as we've not assurance that it's appropriately buffered, but cannot think of a go-idomatic way to
	// do better right now
	var notify *sync.WaitGroup
	n := op.Value(CancelNotifyKey{})
	if n != nil {
		var ok bool
		if notify, ok = n.(*sync.WaitGroup); ok {
			// this lets us block the creator of the cancel-notifier until we've cleaned up
			notify.Add(1)
		}
	}

	// Note: it is up to the caller to handle errors and the closing of the read side of the pipe
	r, w := io.Pipe()
	go func() {
		tw := tar.NewWriter(w)

		defer func() {
			var cerr error
			defer w.Close()

			if notify != nil {
				// inform waiters that we're done with cleanup
				notify.Done()
			}

			if oerr := op.Err(); oerr != nil {
				// don't close the archive if we're truncating the copy - it's misleading
				_ = w.CloseWithError(oerr)
				return
			}

			if cerr = tw.Close(); cerr != nil {
				op.Errorf("Error closing tar writer: %s", cerr.Error())
			}
			if err == nil {
				if cerr != nil {
					op.Errorf("Closing down tar writer with clean exit: %s", cerr)
				} else {
					op.Debugf("Closing down tar writer with pristine exit")
				}
				_ = w.CloseWithError(cerr)
			} else {
				op.Errorf("Closing down tar writer with error during tar: %s", err)
				_ = w.CloseWithError(err)
			}
		}()

		for _, change := range changes {
			if cerr := op.Err(); cerr != nil {
				// this will still trigger the defer to close the archive neatly
				op.Warnf("Aborting tar due to cancellation: %s", cerr)
				break
			}

			if spec.Excludes(op, strings.TrimPrefix(change.Path, "/")) {
				continue
			}

			hdr, err = createHeader(op, dir, change, spec, xattr)

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
			if (hdr.Typeflag == tar.TypeReg || hdr.Typeflag == tar.TypeRegA) && hdr.Size != 0 {
				f, err = os.Open(p)
				if err != nil {
					if os.IsPermission(err) {
						err = nil
					}
					return
				}

				if f != nil {
					// make sure we get out of io.Copy if context is canceled
					done := make(chan struct{})
					go func() {
						select {
						case <-op.Done():
							f.Close()
							// force the io.Copy to exit whether it's in the read or write portion of the copy.
							// If this causes problems with inflight data we can try moving the cancellation notifier
							// Done call here and see if the other end of w/tw will be shut down.
							w.Close()
						case <-done:
						}
					}()

					_, err = io.Copy(tw, f)
					close(done)
					_ = f.Close()
					if err != nil {
						op.Errorf("Error writing archive data: %s", err.Error())
						if err == io.EOF || err == io.ErrClosedPipe {
							// no point in continuing
							break
						}
					}
				}
			}
		}
	}()
	return r, err
}

func createHeader(op trace.Operation, dir string, change docker.Change, spec *FilterSpec, xattr bool) (*tar.Header, error) {
	var hdr *tar.Header
	timestamp := time.Now()

	switch change.Kind {
	case docker.ChangeDelete:
		whiteOutDir := filepath.Dir(change.Path)
		whiteOutBase := filepath.Base(change.Path)
		whiteOut := filepath.Join(whiteOutDir, docker.WhiteoutPrefix+whiteOutBase)

		// FIXME: consult @jzt about how strip should work here
		hdr = &tar.Header{
			ModTime:    timestamp,
			AccessTime: timestamp,
			ChangeTime: timestamp,
		}

		whiteOut = strings.TrimPrefix(whiteOut, "/")
		strippedName := strings.TrimPrefix(whiteOut, spec.StripPath)
		hdr.Name = filepath.Join(spec.RebasePath, strippedName)
	default:
		fi, err := os.Lstat(filepath.Join(dir, change.Path))
		if err != nil {
			op.Errorf("Error getting file info: %s", err.Error())
			return nil, err
		}

		link, err := os.Readlink(filepath.Join(dir, change.Path))
		if err != nil {
			link = change.Path
		}
		hdr, err = tar.FileInfoHeader(fi, link)
		if err != nil {
			op.Errorf("Error getting file info header: %s", err.Error())
			return nil, err
		}

		change.Path = strings.TrimPrefix(change.Path, "/")
		strippedName := strings.TrimPrefix(change.Path, spec.StripPath)
		change.Path = strings.TrimPrefix(change.Path, "/")
		hdr.Name = filepath.Join(spec.RebasePath, strippedName)

		if fi.IsDir() && !strings.HasSuffix(change.Path, "/") {
			hdr.Name += "/"
		}
	}

	if xattr {
		hdr.Xattrs = make(map[string]string)
		hdr.Xattrs[ChangeTypeKey] = change.Kind.String()
	}

	return hdr, nil
}
