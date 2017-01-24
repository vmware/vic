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

package scp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

// the following code is a modified version of https://github.com/gnicod/goscplib
// which follows https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works

type opType string

//Constants
const (
	// C0644 6511 test_utils.go
	BEGINFILE   = opType("C")
	BEGINFOLDER = opType("D")

	BEGINENDFOLDER = opType("0")
	ENDFOLDER      = opType("E")

	// Not yet supported
	// T1449462739 0 1449507489 0
	UPDATETIME = opType("T")

	END = opType("\x00")
)

// ([A-Z])(\d+)\s(\d+)\s(\w*)
var regex = regexp.MustCompile("([A-Z])(\\d+)\\s(\\d+)\\s(.*)")

type Operation struct {
	File *os.File
	Mode os.FileMode
	Size int64
	Name string
}

// Unmarshal an scp header
func Unmarshal(h string) (*Operation, error) {

	s := &Operation{}
	op := opType(h[0:1])
	if op == BEGINENDFOLDER ||
		op == ENDFOLDER ||
		op == END ||
		op == UPDATETIME {

		return s, nil
	}

	res := regex.FindStringSubmatch(h)
	if len(res) != regex.NumSubexp()+1 {
		return nil, fmt.Errorf("malformed header")
	}

	switch op {
	case BEGINFILE:
		fallthrough
	case BEGINFOLDER:
		mode, err := strconv.ParseUint(res[2], 8, 32)
		if err != nil {
			return nil, err
		}
		s.Mode = os.FileMode(mode)

		if op == BEGINFOLDER {
			s.Mode = s.Mode | os.ModeDir
		} else {
			s.Size, err = strconv.ParseInt(res[3], 10, 64)
			if err != nil {
				return nil, err
			}
		}

		s.Name = res[4]

	default:
		return nil, fmt.Errorf("unknown header type")
	}

	return s, nil
}

func (s *Operation) String() string {

	size := s.Size

	typ := BEGINFILE
	if s.Mode.IsDir() {
		typ = BEGINFOLDER
		// folders don't have a size.
		size = 0
	}

	// strip the upper bits of mode
	mode := s.Mode & ^os.ModeType

	name := filepath.Base(s.Name)

	return fmt.Sprintf("%s0%o %d %s\n", typ, mode, size, name)
}

// update the stat fields
func (s *Operation) stat() error {
	stat, err := os.Stat(s.File.Name())
	if err != nil {
		return fmt.Errorf("Error stat'ing %s: %s", s.Name, err)
	}

	s.Mode = stat.Mode()
	s.Size = stat.Size()

	// we just want the basepath
	s.Name = filepath.Base(s.File.Name())

	return nil
}

// Write will unmarshal an scp header, open the dirent for writing, write the
// file contents, and return the relevant Operation
func Write(r io.ReadCloser, path string) (*Operation, error) {
	rdr := bufio.NewReader(r)
	header, err := rdr.ReadString('\n')
	if err != nil {
		return nil, err
	}

	s, err := Unmarshal(header)
	if err != nil {
		return nil, err
	}

	s.File, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, s.Mode)
	if err != nil {
		return nil, err
	}

	n, err := rdr.WriteTo(s.File)
	if err != nil {
		return nil, err
	}

	r.Close()

	// it's unclear why we have to do this since we pass the mode bits in
	// openfile, but setting a file to 777 yields 775 without it.
	if err := os.Chmod(s.Name, s.Mode); err != nil {
		return nil, err
	}

	log.Printf("Created %s (%d bytes)", s.File.Name(), n)
	return s, nil
}

// OpenFile is like os.OpenFile except it populates the Operation.  It opens an existing file.
// Useful for visiting an existing directory.
func OpenFile(path string, flag int, perm os.FileMode) (*Operation, error) {
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	s := &Operation{}
	s.File = f
	s.Name = path

	if err = s.stat(); err != nil {
		return nil, err
	}

	return s, nil
}

// Once a file is open (using OpenFile), a header will be generated, and the
// file's contents will be written to the writer.
func (s *Operation) Read(r io.Writer) (int64, error) {

	hdr := s.String()

	// write the header
	if _, err := r.Write([]byte(hdr)); err != nil {
		return 0, err
	}

	n, err := io.CopyN(r, s.File, s.Size)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// Close closes the file and updates the header information by stat'ing the dirent
func (s *Operation) Close() error {
	if err := s.File.Close(); err != nil {
		return err
	}

	if err := s.stat(); err != nil {
		return err
	}

	return nil
}
