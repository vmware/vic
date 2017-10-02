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

package vchlog

import (
	"bytes"
	"sync"
	"io"
)

// BufferedPipe struct implements a pipe readwriter with buffer
// buffer: the internal buffer to hold data
// c: the sync locker to manage concurrent reads and writes
// closed: boolean indicating if the stream is closed
type BufferedPipe struct {
	buffer *bytes.Buffer
	c      *sync.Cond
	closed bool
}

// NewBufferedPipe returns a new buffered pipe instance
// the internal buffer is initialized to default size.
// Since internal memory is used, need to make sure that buffered data is bounded.
func NewBufferedPipe() *BufferedPipe {
	var m sync.Mutex
	c := sync.NewCond(&m)
	return &BufferedPipe{
		buffer: bytes.NewBuffer(nil),
		c:      c,
		closed: false,
	}
}

// Read is blocked until a writer in the queue is done writing (until data is available)
func (bp *BufferedPipe) Read(data []byte) (n int, err error) {
	bp.c.L.Lock()
	defer bp.c.L.Unlock()

	// pipe closed, drop all left-over data
	if bp.closed {
		return 0, io.EOF
	}
	for bp.buffer.Len() == 0 && !bp.closed {
		bp.c.Wait()
	}

	return bp.buffer.Read(data)
}

// Write writes to the internal buffer, and signals one of the reader in queue to start reading.
func (bp *BufferedPipe) Write(data []byte) (n int, err error) {
	bp.c.L.Lock()
	defer bp.c.L.Unlock()
	defer bp.c.Signal()

	if bp.closed {
		return 0, io.ErrUnexpectedEOF
	}

	return bp.buffer.Write(data)
}

// Close closes the pipe.
func (bp *BufferedPipe) Close() (err error) {
	bp.c.L.Lock()
	defer bp.c.L.Unlock()
	defer bp.c.Signal()
	bp.closed = true
	return nil
}
