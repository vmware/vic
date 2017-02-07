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

package ioutil

import (
	"bufio"
	"encoding/json"
	"io"

	log "github.com/Sirupsen/logrus"
)

// LogReader unwraps entires into LogEntry structs.
type LogReader struct {
	io.ReadCloser
	s *bufio.Scanner

	// if the buffer we want to read into is too short for our log entry,
	// we will need to return io.ErrShortBuffer and try again on the next call,
	// necessitating this pointer to the previous LogEntry we attempted to use.
	prev *LogEntry
}

// NewLogReader wraps an io.ReadCloser in a LogReader.
func NewLogReader(r io.ReadCloser) *LogReader {
	return &LogReader{
		ReadCloser: r,
		s:          bufio.NewScanner(r),
	}
}

// Read unmarhsals a line from the underlying stream of log entries into
// a LogEntry struct, and then writes that struct's log message to the
// supplied []byte slice.
func (lr *LogReader) Read(p []byte) (int, error) {

	entry := &LogEntry{}

	if lr.prev != nil {
		entry = lr.prev
	} else {
		if lr.s.Scan() {
			if err := json.Unmarshal(lr.s.Bytes(), entry); err != nil {
				log.Debugf("Error unmarshaling log entry: %s", err.Error())
				return 0, err
			}
		} else {
			err := lr.s.Err()
			if err == nil {
				return 0, io.EOF
			}
			return 0, err
		}
	}

	if len(p) < len(entry.Log) {
		lr.prev = entry
		return 0, io.ErrShortBuffer
	}

	w := copy(p, []byte(entry.Log))

	// reset previous entry to nil
	lr.prev = nil

	return w, nil
}
