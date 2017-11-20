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

// +build !windows

package archive

import (
	"archive/tar"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"errors"

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

// unpack will unpack the given tarstream(if it is a tar stream) on the local filesystem based on the specified root
// combined with any rebase from the path spec
//
// the pathSpec will include the following elements
// - include : any tar entry that has a path below(after stripping) the include path will be written
// - strip : The strip string will indicate the
// - exlude : marks paths that are to be excluded from the write operation
// - rebase : marks the the write path that will be tacked onto (appended or prepended? TODO improve this comment) the "root". e.g /tmp/unpack + /my/target/path = /tmp/unpack/my/target/path
func InvokeUnpack(op trace.Operation, tarStream io.Reader, filter *FilterSpec) error {
	// op.Debugf("unpacking archive to root: %s, filter: %+v", root, filter)

	// Online datasource is sending a tar reader instead of an io reader.
	// Type check here to see if we actually need to create a tar reader.
	var tr *tar.Reader
	if trCheck, ok := tarStream.(*tar.Reader); ok {
		tr = trCheck
	} else {
		tr = tar.NewReader(tarStream)
	}

	// fi, err := os.Stat(root)
	// if err != nil {
	// 	// the target unpack path does not exist. We should not get here.
	// 	op.Errorf("tar unpack target does not exist: %s", root)
	// 	return err
	// }

	// if !fi.IsDir() {
	// 	err := fmt.Errorf("unpack root target is not a directory: %s", root)
	// 	op.Error(err)
	// 	return err
	// }

	op.Debugf("using FilterSpec : (%#v)", *filter)
	// process the tarball onto the filesystem
	for {
		header, err := tr.Next()
		if err == io.EOF {
			// This indicates the end of the archive
			break
		}

		if err != nil {
			op.Errorf("Error reading tar header: %s", err)
			return err
		}

		op.Debugf("processing tar header: asset(%s), size(%d)", header.Name, header.Size)
		// skip excluded elements unless explicitly included
		if filter.Excludes(op, header.Name) {
			continue
		}

		// fix up path
		stripped := strings.TrimPrefix(header.Name, filter.StripPath)
		rebased := filepath.Join(filter.RebasePath, stripped)
		// absPath := filepath.Join(root, rebased)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(rebased, header.FileInfo().Mode())
			if err != nil {
				op.Errorf("Failed to create directory%s: %s", rebased, err)
				return err
			}
		case tar.TypeSymlink:
			err := os.Symlink(header.Linkname, rebased)
			if err != nil {
				op.Errorf("Failed to create symlink %s->%s: %s", rebased, header.Linkname, err)
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(rebased, fileWriteFlags, header.FileInfo().Mode())
			if err != nil {
				op.Errorf("Failed to open file %s: %s", rebased, err)
				return err
			}
			_, err = io.Copy(f, tr)
			// TODO: add ctx.Done cancellation
			f.Close()
			if err != nil {
				return err
			}
		default:
			// TODO: add support for special file types - otherwise we will do absurd things such as read infinitely from /dev/random
		}
		op.Debugf("Finished writing to: %s", rebased)
	}
	return nil
}

func simpleCopy(src, dst string) error {

	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := out.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	err = out.Sync()
	return err

}

func Unpack(op trace.Operation, tarStream io.Reader, filter *FilterSpec, root string) error {
	// p := fmt.Sprintf("%s%sunpack", filter.RebasePath, string(os.PathSeparator))
	// op.Infof("XXX %s", p)

	// op.Infof("XXX copying unpack to %s", root)
	// dest := fmt.Sprintf("%s%sunpack", root, string(os.PathSeparator))

	// op.Infof("XXX dest is %s", dest)

	// if err := simpleCopy("/bin/unpack", dest); err != nil {
	// 	op.Errorf("XXX Couldn't copy because %s", err.Error())
	// 	return err
	// }

	// defer os.Remove(dest)

	// op.Infof("XXX has root changed for some reason?? %s", root)
	// files, err := ioutil.ReadDir("/bin")
	// if err != nil {
	// 	op.Error(err)
	// }
	// op.Infof("contents of /bin")
	// for _, f := range files {
	// 	op.Infof("XXX file %s", f.Name())
	// }
	// files, err = ioutil.ReadDir("/sbin")
	// if err != nil {
	// 	op.Error(err)
	// }

	// op.Infof("contents of /sbin")
	// for _, f := range files {
	// 	op.Infof("XXX file %s", f.Name())
	// }

	// #nosec: executable is allowed
	// if err := os.Chmod(dest, 0755); err != nil {
	// 	op.Errorf("XXX Couldn't chmod because %s", err.Error())
	// 	return err
	// }

	// execute the unpack binary
	cmd := exec.Cmd{
		Path: "/bin/unpack",
		Dir:  root,
		Args: []string{"/bin/unpack", op.ID(), root},
		// SysProcAttr: &syscall.SysProcAttr{
		// 	Chroot: root,
		// },
	}

	encFilter, err := EncodeFilterSpec(op, filter)
	if err != nil {
		return err
	}

	op.Infof("XXX Creating stdinpipe")
	stdin, err := cmd.StdinPipe()

	if err != nil {
		op.Infof("XXX err non-nill")
		op.Error(err)
		return err
	}

	if stdin == nil {
		err = errors.New("stdin was nil")
		op.Error(err)
		return err
	}

	op.Infof("XXX Creating stdoutpipe")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		op.Error(err)
		return err
	}
	if stdout == nil {
		err = errors.New("stdout was nil")
		op.Error(err)
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(stdout io.ReadCloser, wg *sync.WaitGroup) {
		defer wg.Done()
		op.Infof("XXX Running stdout reader")
		output := []byte{}
		for {
			if stdout == nil {
				op.Infof("XXX done")
				break
			}
			op.Infof("XXX repeat read")
			n, err := stdout.Read(output)
			if n > 0 && err != nil {
				op.Infof("XXX %s", string(output))
				output = []byte{}
			}
			if err != nil {
				op.Errorf("XXX %s", err.Error())
				break
			}
			if n <= 0 {
				op.Infof("XXX End of stream")
				break
			}
		}
	}(stdout, &wg)

	// op.Infof("stdin ok but need to defer closing it")
	// defer func() {
	// 	if stdin == nil {
	// 		op.Errorf("stdin was nil at callback time")
	// 		return
	// 	}
	// 	if err := stdin.Close(); err != nil {
	// 		op.Error(err)
	// 	}
	// }()

	op.Infof("XXX running command %+v", cmd)
	if err := cmd.Start(); err != nil {
		op.Errorf("Couldn't start archive binary: %s", err.Error())
		return err
	}
	op.Infof("XXX filterspec stuff, write that to stdin")

	bencFilter := []byte(*encFilter)
	o, err := stdin.Write(bencFilter)

	if o != len(bencFilter) {
		return errors.New("XXX Fucked up trying to send the filterspec over")
	}
	if err != nil {
		return err
	}

	op.Infof("XXX Insert the field sep")
	stdin.Write([]byte("\n"))

	op.Infof("XXX Write the tarstream to the binary")
	oneByte := *new([]byte)
	for n, err := tarStream.Read(oneByte); err == nil; {
		if n > 0 {
			o, er := stdin.Write(oneByte)
			if err != nil {
				return er
			}
			if o != n {
				op.Errorf("didnt read and write the same stuff")
				return errors.New("didn't read and write the same stuff")
			}
			op.Infof("XXX wrote %d bytes", o)
		} else {
			break
		}
	}

	if err != nil {
		return err
	}

	op.Infof("XXX process completed, waiting on gofunc")
	wg.Wait()

	op.Infof("XXX Wait for process")
	// go func read from stdout for logging
	if err := cmd.Wait(); err != nil {
		op.Errorf("XXX Got %s while waiting for process", err.Error())
		return err
	}

	return nil
}
