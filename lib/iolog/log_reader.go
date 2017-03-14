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

package iolog

import (
	"encoding/base64"
	"encoding/binary"
	"io"
	"time"
)

const (
	// RFC3339NanoFixed is Docker's version of RFC339Nano that pads the
	// nanoseconds with zeros to ensure that the timestamps are aligned in the logs.
	RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"
)

// LogReader reads containerVM entries from a stream and decodes them into
// their original form
type LogReader struct {
	io.ReadCloser
	prev []byte
	ts   bool
}

// NewLogReader wraps an io.ReadCloser in a LogReader.
func NewLogReader(r io.ReadCloser, ts bool) *LogReader {
	return &LogReader{
		ReadCloser: r,
		ts:         ts}
}

// Read reads a 10 byte header and decodes it into the timestamp, stream and
// size of an entry. It uses the size to read the next set of bytes as the
// message, and then copies the message into the supplied buffer, saving
// what will not fit in the buffer for the next call to Read.
func (lr *LogReader) Read(p []byte) (int, error) {
	var (
		err  error
		n, w int
		ts   string
	)

	ehdr := make([]byte, encodedHeaderLengthBytes)
	msg := lr.prev
	partial := true // treat msg as a partial entry until we verify otherwise

	if msg == nil {
		// we know msg is not a partial entry as we had no bytes left from the previous call
		partial = false

		// read a header
		n = 0
		for n < encodedHeaderLengthBytes {
			w, err = lr.ReadCloser.Read(ehdr[n:])
			n += w
			if err != nil {
				return 0, err
			}
		}

		// decode base64 header
		hdr, err := base64.StdEncoding.DecodeString(string(ehdr))
		if err != nil {
			return 0, err
		}

		// parse header
		ts = time.Unix(0, int64(binary.LittleEndian.Uint64(hdr[:8]))).Format(RFC3339NanoFixed)
		s := binary.LittleEndian.Uint16(hdr[8:10])
		// stream := int((s&streamFlag) >> 3)
		size := int(s >> 4)

		// read the associated entry
		msg = make([]byte, size)
		n = 0
		for n < size {
			w, err = lr.ReadCloser.Read(msg[n:])
			n += w
			if err != nil {
				if err != io.EOF {
					// only return if not EOF as we may actually have some bytes to copy
					return 0, err
				}
				break
			}
		}
	}

	lr.prev = nil
	if len(p) < len(msg) {
		// copy what we can and save the rest for the next call
		lr.prev = msg[len(p):]
		msg = msg[:len(p)]
	}

	// add timestamp if enabled
	if lr.ts && !partial {
		msg = append([]byte(ts+" "), msg...)
	}

	// write the log message
	return copy(p, msg), err
}
