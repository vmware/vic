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
	"strings"
	"time"

	"github.com/mailru/easyjson/jwriter"
)

// LogWriter wraps entries in a LogEntry struct with the
// timestamp of the entry.
type LogWriter struct {
	w    io.Writer
	prev []byte
}

// NewLogWriter wraps an io.WriteCloser in a LogWriter
func NewLogWriter(w io.Writer) *LogWriter {
	return &LogWriter{w: w}
}

// Write takes an incoming byte slice and wraps it in a LogEntry,
// adding a timestamp. It then serializes the LogEntry and writes the
// serialized bytes to the underlying stream, adding a newline at the end for
// human readability
func (lw *LogWriter) Write(p []byte) (n int, err error) {

	i := 0
	l := len(p)

	s := bufio.NewScanner(strings.NewReader(string(p)))
	for s.Scan() {
		b := s.Bytes()

		// This is used to keep count of how many bytes we have scanned from the supplied
		// buffer. We add one to account for the '\n' byte stripped by the scanner.
		i += len(b) + 1

		// did we store the beginning of a log entry on the previous call?
		if lw.prev != nil {
			// if so, append to it what we just read to continue or complete that entry
			b = append(lw.prev, b...)
		}

		// One of three cases now must be true:
		// 1. i < l:  We have not hit the end of the buffer yet, but have a complete entry
		// 2. i == l: We hit the end and the last byte in the buffer was '\n', and we have a
		//    complete log entry.
		// 3. i > l:  We hit the end and the last byte was not '\n', indicating the beginning
		//    (or middle) of a log entry that we must store so that we can get the end of the
		//    entry from the new buffer supplied in a subsequent call to Write.
		//
		// Cases 1 and 2 require no action here as the next scan will either yield another
		// token, or return false and signal that we're done with this buffer.
		if i > l {
			// store partial log entry
			lw.prev = b
			continue // next call to Scan() should return false
		}

		// we have a complete entry, time to write it
		entry := LogEntry{
			Msg: string(b) + "\n",
			Ts:  time.Now(),
		}

		jw := &jwriter.Writer{}
		entry.MarshalEasyJSON(jw)

		n, err = jw.DumpTo(lw.w)
		if err != nil {
			return n, err
		}
		// add a newline to the end of a JSON log entry
		lw.w.Write([]byte{'\n'})

		// we wrote an entry successfully, so reset prev to nil
		lw.prev = nil
	}

	if err := s.Err(); err != nil {
		return 0, err
	}

	// inform the caller that all received bytes were successfully written
	return len(p), err
}
