// Copyright 2016 VMware, Inc. All Rights Reserved.
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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func checkHeader(expected, actual *Operation) error {
	if expected.String() != actual.String() {
		return fmt.Errorf("String() doesn't match %s %s", expected.String(), actual.String())
	}

	if expected.Mode != actual.Mode {
		return fmt.Errorf("Mode doesn't match 0%o 0%o", expected.Mode, actual.Mode)
	}

	// ignore size if directory
	if !os.FileMode(expected.Mode).IsDir() &&
		expected.Size != actual.Size {
		return fmt.Errorf("Size doesn't match %d %d", expected.Size, actual.Size)
	}

	if expected.Name != actual.Name {
		return fmt.Errorf("Name doesn't match %s %s", expected.Name, actual.Name)
	}
	return nil
}

func TestHeaderParsing(t *testing.T) {
	headers := []*Operation{
		{
			Mode: 0644,
			Size: 6,
			Name: "test",
		},
		{
			// need to test this on M$
			Mode: 0644 | os.ModeDir,
			Size: 0,
			Name: "test",
		},
	}

	for _, actual := range headers {
		outHeader, err := Unmarshal(actual.String())
		if err != nil {
			t.Error(err)
			return
		}

		err = checkHeader(actual, outHeader)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestHeaderCreationFromExistingFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "empty")
	if err != nil {
		t.Error(err)
		return
	}
	defer tmpFile.Close()

	err = tmpFile.Chmod(os.FileMode(0644))
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(tmpFile.Name())

	expected := &Operation{
		Mode: os.FileMode(0644),
		Size: 0,
		Name: filepath.Base(tmpFile.Name()),
	}

	actual, err := OpenFile(tmpFile.Name(), os.O_RDONLY, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer actual.Close()

	err = checkHeader(expected, actual)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestHeaderCreationFromExistingDir(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "empty")
	if err != nil {
		t.Error(err)
		return
	}

	err = os.Chmod(tmpDir, os.FileMode(0644))
	if err != nil {
		t.Error(err)
		return
	}

	defer os.Remove(tmpDir)

	expected := &Operation{
		Mode: os.FileMode(0644) | os.ModeDir,
		Size: 0,
		Name: filepath.Base(tmpDir),
	}

	actual, err := OpenFile(tmpDir, os.O_RDONLY, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer actual.Close()

	err = checkHeader(expected, actual)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestFileCreationFromHeader(t *testing.T) {
	t.Skip("skipping test until functional with full path")

	tmpFile, err := ioutil.TempFile("", "empty")
	if err != nil {
		t.Error(err)
		return
	}
	fileToCreateName := tmpFile.Name()
	os.Remove(fileToCreateName)

	expected := &Operation{
		Mode: os.FileMode(0677),
		Size: 0,
		Name: fileToCreateName,
	}

	hdr := expected.String()

	dummyStream := strings.NewReader(hdr)
	// try creating the file from the header
	actual, err := Write(ioutil.NopCloser(dummyStream), fileToCreateName)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(actual.File.Name())

	if err = actual.Close(); err != nil {
		t.Error(err)
		return
	}

	err = checkHeader(expected, actual)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDirCreationFromHeader(t *testing.T) {
	t.Skip("skipping test until functional with full path")

	f, err := ioutil.TempDir("", "empty")
	if err != nil {
		t.Error(err)
		return
	}
	os.Remove(f)

	expected := &Operation{
		Mode: os.FileMode(0777),
		Size: 0,
		Name: f,
	}

	hdr := expected.String()

	dummyStream := strings.NewReader(hdr)
	// try creating the dir from the header
	actual, err := Write(ioutil.NopCloser(dummyStream), f)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(actual.File.Name())

	if err = actual.Close(); err != nil {
		t.Error(err)
		return
	}

	err = checkHeader(expected, actual)
	if err != nil {
		t.Error(err)
		return
	}
}
