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
	"testing"
)

func TestWriteJsonEntry(t *testing.T) {
	var buf bytes.Buffer
	w := NewLogWriter(&buf)
	msg := "The quick brown fox jumped over the lazy dog\n"

	n, err := w.Write([]byte(msg))
	if err != nil {
		t.Errorf(err.Error())
	}

	if n != len(msg) {
		t.Errorf("Wrote %d bytes, expected to write %d", n, len(msg))
	}

	entry := &LogEntry{}
	if err := json.Unmarshal(buf.Bytes(), entry); err != nil {
		t.Errorf(err.Error())
	}

	if entry.Msg != msg {
		t.Errorf("Unmarshaled \"%s\", expected \"%s\"", entry.Msg, msg)
	}
}

func TestWriteSpecialChars(t *testing.T) {
	var buf bytes.Buffer
	w := NewLogWriter(&buf)
	msg := "~!@#$%^&*()_+\n"

	n, err := w.Write([]byte(msg))
	if err != nil {
		t.Errorf(err.Error())
	}

	if n != len(msg) {
		t.Errorf("Wrote %d bytes, expected to write %d", n, len(msg))
	}

	entry := &LogEntry{}
	if err := json.Unmarshal(buf.Bytes(), entry); err != nil {
		t.Errorf(err.Error())
	}

	if entry.Msg != msg {
		t.Errorf("Unmarshaled \"%s\", expected \"%s\"", entry.Msg, msg)
	}
}

func TestWritePartialEntry(t *testing.T) {
	var buf bytes.Buffer
	w := NewLogWriter(&buf)
	msg1 := "The quick brown fox"
	msg2 := " jumped over "
	msg3 := "the lazy dog\n"

	n, err := w.Write([]byte(msg1))
	if err != nil {
		t.Errorf(err.Error())
	}

	if n != len(msg1) {
		t.Errorf("Wrote %d bytes, expected to write %d", n, len(msg1))
	}

	n, err = w.Write([]byte(msg2))
	if err != nil {
		t.Errorf(err.Error())
	}

	if n != len(msg2) {
		t.Errorf("Wrote %d bytes, expected to write %d", n, len(msg2))
	}

	n, err = w.Write([]byte(msg3))
	if err != nil {
		t.Errorf(err.Error())
	}

	if n != len(msg3) {
		t.Errorf("Wrote %d bytes, expected to write %d", n, len(msg3))
	}

	entry := &LogEntry{}
	if err := json.Unmarshal(buf.Bytes(), entry); err != nil {
		t.Errorf(err.Error())
	}

	expected := msg1 + msg2 + msg3
	if entry.Msg != expected {
		t.Errorf("Unmarshaled \"%s\", expected \"%s\"", entry.Msg, expected)
	}
}
