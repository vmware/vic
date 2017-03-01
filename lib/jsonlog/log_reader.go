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
	"bufio"
	"io"

	"github.com/mailru/easyjson/jlexer"
)

// LogReader unwraps entires into LogEntry structs.
type LogReader struct {
	io.ReadCloser

	s    *bufio.Scanner
	prev string
}

// NewLogReader wraps an io.ReadCloser in a LogReader.
func NewLogReader(r io.ReadCloser) *LogReader {
	return &LogReader{
		ReadCloser: r,
		s:          bufio.NewScanner(r),
	}
}

// Read unmarshals a line from the underlying stream of log entries into
// a LogEntry struct, and then writes that struct's log message to the
// supplied []byte slice.
func (lr *LogReader) Read(p []byte) (int, error) {

	// use unwritten previous message if we happen to have one
	msg := lr.prev

	if msg == "" {
		// scanning one Read call at a time allows us to stay current when tailing
		// output from the underlying stream
		if lr.s.Scan() {
			l := &jlexer.Lexer{Data: lr.s.Bytes()}
			entry := &LogEntry{}
			entry.UnmarshalEasyJSON(l)
			msg = entry.Msg
		} else {
			// If we get here, there are no more bytes to scan from the underlying stream.
			err := lr.s.Err()
			// return non-EOF error if we encountered one while scanning
			if err != nil {
				return 0, err
			}
			return 0, io.EOF
		}
	}

	// ensure the supplied buffer size is adequate
	if len(p) < len(msg) {
		// if not, save it and hold out for a larger buffer
		lr.prev = msg
		return 0, io.ErrShortBuffer
	}

	// write the log message
	w := copy(p, []byte(msg))

	// reset previous message to nil
	lr.prev = ""

	return w, nil
}
