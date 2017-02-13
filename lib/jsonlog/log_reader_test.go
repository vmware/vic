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

package jsonlog

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
	"time"
)

type TestReadCloser struct {
	*bytes.Buffer
	io.Closer
}

func TestReadEntry(t *testing.T) {
	now := time.Now()
	entry := &LogEntry{
		Msg: "The quick brown fox jumped over the lazy dog\n",
		Ts:  now,
	}

	b, err := json.Marshal(&entry)
	if err != nil {
		t.Errorf(err.Error())
	}
	b = append(b, 0xa)

	in := TestReadCloser{
		Buffer: bytes.NewBuffer(b),
	}

	out := make([]byte, len(entry.Msg))

	r := NewLogReader(in)
	n, err := r.Read(out)
	if err != nil {
		t.Errorf(err.Error())
	}

	if n != len(out) {
		t.Errorf("Read %d bytes, expected to read %d", n, len(out))
	}

	expected := string(out)
	if expected != entry.Msg {
		t.Errorf("Read %s, expected %s", expected, entry.Msg)
	}
}

func TestReadSpecialChars(t *testing.T) {
	now := time.Now()
	entry := &LogEntry{
		Msg: "~!@#$%^&*()_+\n",
		Ts:  now,
	}

	b, err := json.Marshal(&entry)
	if err != nil {
		t.Errorf(err.Error())
	}
	b = append(b, 0xa)

	in := TestReadCloser{
		Buffer: bytes.NewBuffer(b),
	}

	out := make([]byte, len(entry.Msg))

	r := NewLogReader(in)
	n, err := r.Read(out)
	if err != nil {
		t.Errorf(err.Error())
	}

	if n != len(out) {
		t.Errorf("Read %d bytes, expected to read %d", n, len(out))
	}

	expected := string(out)
	if entry.Msg != expected {
		t.Errorf("Read %s, expected %s", expected, entry.Msg)
	}
}

func TestReadShortBuffer(t *testing.T) {
	now := time.Now()
	entry := &LogEntry{
		Msg: "The quick brown fox jumped over the lazy dog\n",
		Ts:  now,
	}

	b, err := json.Marshal(&entry)
	if err != nil {
		t.Errorf(err.Error())
	}
	b = append(b, 0xa)

	in := TestReadCloser{
		Buffer: bytes.NewBuffer(b),
	}

	out := make([]byte, len(entry.Msg)/2)

	r := NewLogReader(in)
	n, err := r.Read(out)
	if err != nil && err != io.ErrShortBuffer {
		t.Errorf(err.Error())
	}

	if n != 0 {
		t.Errorf("Read %d bytes, expected to read 0 bytes", n)
	}

	out = make([]byte, len(entry.Msg))
	n, err = r.Read(out)
	if err != nil {
		t.Errorf(err.Error())
	}

	if n != len(out) {
		t.Errorf("Read %d bytes, expected to read %d", n, len(out))
	}

	expected := string(out)
	if expected != entry.Msg {
		t.Errorf("Read %s, expected %s", expected, entry.Msg)
	}
}
