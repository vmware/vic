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
	"io"
	"sync"
	"time"

	"github.com/vmware/vic/pkg/errors"
)

// BufferedPipe struct implements a pipe readwriter with buffer
type BufferedPipe struct {
	// buffer: the internal buffer to hold data
	buffer *bytes.Buffer
	// c: the sync locker to manage concurrent reads and writes
	c *sync.Cond
	// readClosed: boolean indicating if read end of the pipe is closed
	readClosed bool
	// writeClosed: boolean indicating if write end of the pipe is closed
	writeClosed bool
	// readerReady: boolean indicating if the reader is ready
	readerReady bool
}

// NewBufferedPipe returns a new buffered pipe instance
// the internal buffer is initialized to default size.
// Since internal memory is used, need to make sure that buffered data is bounded.
func NewBufferedPipe() *BufferedPipe {
	var m sync.Mutex
	c := sync.NewCond(&m)
	return &BufferedPipe{
		buffer:      bytes.NewBuffer(nil),
		c:           c,
		readClosed:  false,
		writeClosed: false,
		readerReady: false,
	}
}

// Read is blocked until a writer in the queue is done writing (until data is available)
func (bp *BufferedPipe) Read(data []byte) (n int, err error) {
	bp.c.L.Lock()
	defer bp.c.L.Unlock()
	defer bp.c.Broadcast()

	bp.readerReady = true // now we have a valid consumer to switch to flushing in Close

	for bp.buffer.Len() == 0 && !bp.readClosed {
		bp.c.Wait()
	}
	if bp.readClosed {
		return 0, io.EOF
	}

	return bp.buffer.Read(data)
}

// Write writes to the internal buffer, and signals one of the reader in queue to start reading.
func (bp *BufferedPipe) Write(data []byte) (n int, err error) {
	bp.c.L.Lock()
	defer bp.c.L.Unlock()
	defer bp.c.Broadcast()

	if bp.writeClosed {
		return 0, io.ErrUnexpectedEOF
	}

	return bp.buffer.Write(data)
}

// Close closes the pipe, flushes remaining data or drops if consumer times out (30 minutes).
// timeout implementation reference: https://github.com/golang/go/issues/9578
func (bp *BufferedPipe) Close() (err error) {
	bp.c.L.Lock()
	defer bp.c.L.Unlock()
	defer bp.c.Broadcast()

	bp.writeClosed = true

	if bp.readerReady { // flush only if there is a valid consumer
		done := make(chan bool)
		for bp.buffer.Len() > 0 {
			go func() {
				bp.c.Wait()
				done <- true
			}()
			select {
			case <-time.After(time.Minute * 30): // timeout if consumer inactive for 30 minutes
				err = errors.New("buffered data left in pipe; consumer inactive for 30 minutes")
				break
			case <-done:
			}
		}
	}

	bp.readClosed = true
	return err
}
