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
)

// LogReader reads containerVM entries from a stream and decodes them into
// their original form
type LogReader struct {
	io.ReadCloser
	prev []byte
}

// NewLogReader wraps an io.ReadCloser in a LogReader.
func NewLogReader(r io.ReadCloser) *LogReader {
	return &LogReader{ReadCloser: r}
}

// Read reads a 10 byte header and decodes it into the timestamp, stream and
// size of an entry. It uses the size to read the next set of bytes as the
// message, and then copies the message into the supplied buffer, saving
// what will not fit in the buffer for the next call to Read.
func (lr *LogReader) Read(p []byte) (int, error) {
	var (
		err  error
		n, w int
	)

	ehdr := make([]byte, encodedHeaderLengthBytes)
	msg := lr.prev

	if msg == nil {
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
		// ts := time.Unix(0, int64(binary.LittleEndian.Uint64(h[:8])))
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

	// write the log message
	return copy(p, msg), err
}
