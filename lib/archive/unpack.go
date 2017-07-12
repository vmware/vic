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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/vmware/vic/pkg/trace"
)

const (
	// fileWriteFlags is a collection of flags configuring our writes for general tar behavior
	//
	// O_CREATE = Create file if it does not exist
	// O_TRUNC = truncate file to 0 length if it does exist(overwrite the file)
	// O_WRONLY = We use this since we do not intend to read, we only need to write.
	fileWriteFlags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
)

// unpack will unpack the given tarstream(if it is a tar stream) on the local filesystem based on the unpackPath in the path spec
//
// the pathSpec will include the following elements
// - include : any tar entry that has a path below(after stripping) the include path will be written
// - strip : The strip string will indicate the
// - exlude : marks paths that are to be excluded from the write operation
// - rebase : marks the the write path that will be tacked onto the "unpackPath". e.g /tmp/unpack + /my/target/path = /tmp/unpack/my/target/path
func Unpack(op trace.Operation, tarStream io.Reader, filter *FilterSpec, unpackPath string) error {
	// the tar stream should be wrapped up at the end of this call
	tr := tar.NewReader(tarStream)

	strip := filter.StripPath
	target := filter.RebasePath

	if target == "" {
		op.Debugf("Bad target path in FilterSpec (%#v)", filter)
		return fmt.Errorf("Invalid write target specified")
	}

	if strip == "" {
		op.Debugf("Strip path was set to \"\"")
	}

	if _, err := os.Stat(unpackPath); err != nil {
		// the target unpack path does not exist. We should not get here.
		op.Errorf("tar unpack target does not exist (%s)", unpackPath)
		return err
	}

	finalTargetPath := filepath.Join(unpackPath, target)
	op.Debugf("finalized target path for Tar unpack operation at (%s)", finalTargetPath)

	// process the tarball onto the filesystem
	for {
		header, err := tr.Next()
		if err == io.EOF {
			// This indicates the end of the archive
			break
		}
		if err != nil {
			// it is likely in this case that we were not given a legitimate tar stream
			op.Debugf("Received error (%s) when attempting a tar read operation on target stream", err)
			return err
		}

		// skip excluded elements unless explicitly included
		if filter.Excludes(op, header.Name) {
			continue
		}

		// fix up path
		strippedTargetPath := strings.TrimPrefix(header.Name, strip)
		writePath := filepath.Join(finalTargetPath, strippedTargetPath)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(writePath, header.FileInfo().Mode())
			if err != nil {
				return err
			}
			continue
		case tar.TypeSymlink:

			// NOTE: symbolic links cannot span mounts.
			// This can cause us to create a regular file instead of
			// a symlink. this behavior should be evaluated for a proper
			// response, for now a regular file is made since that is the
			// target of a sym link operation
			err := os.Symlink(header.Linkname, writePath)
			if err == nil {
				continue
			}
			op.Infof("error from os symlink (%s)", err.Error())
			fallthrough
		default:
			// we will treat the default as a regular file
			targetDir, _ := filepath.Split(writePath)

			// FIXME: this is a hack we must include the directory before this instead of excluding it. since the permissions could be different.
			err = os.MkdirAll(targetDir, header.FileInfo().Mode())

			err = func(path string, flags int, perm os.FileMode, tr *tar.Reader) error {
				writtenTarFile, err := os.OpenFile(path, flags, perm)
				if err != nil {
					return nil
				}
				defer writtenTarFile.Close()

				_, err = io.Copy(writtenTarFile, tr)
				if err != nil {
					return err
				}
				return nil
			}(writePath, fileWriteFlags, header.FileInfo().Mode(), tr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
