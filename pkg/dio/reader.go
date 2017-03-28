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

// Package dio adds dynamic behaviour to the standard io package mutliX types
package dio

import (
	"io"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// DynamicMultiReader adds dynamic add/remove to the base multireader behaviour
type DynamicMultiReader interface {
	io.Reader
	Add(...io.Reader)
	Remove(io.Reader)
	Close() error
}

type multiReader struct {
	mutex     sync.Mutex
	readersMu sync.Mutex
	cond      *sync.Cond
	readers   []io.Reader
	pr        *io.PipeReader
	pw        *io.PipeWriter
}

func (t *multiReader) Read(p []byte) (int, error) {
	t.readersMu.Lock()
	for len(t.readers) == 0 {
		t.cond.Wait()
	}
	t.readersMu.Unlock()

	t.mutex.Lock()
	defer t.mutex.Unlock()

	n, err := t.pr.Read(p)
	if err == io.ErrClosedPipe {
		t.pr.Close()
		err = io.EOF
	}

	return n, err
}

func (t *multiReader) Add(readers ...io.Reader) {
	t.readersMu.Lock()
	defer t.readersMu.Unlock()

	if len(t.readers) == 0 {
		t.mutex.Lock()
		t.pr, t.pw = io.Pipe()
		t.mutex.Unlock()
	}

	t.readers = append(t.readers, readers...)
	for i := range readers {
		go func(r io.Reader) {
			var err error

			defer func() {
				log.Debugf("reader %p exited with error: %s", r, err)
				t.Remove(r)
			}()

			buf := make([]byte, 512)
			n := 0
			for {
				n, err = r.Read(buf)
				if n > 0 {
					if _, ew := t.pw.Write(buf[:n]); ew != nil {
						return
					}
				}

				if err != nil {
					if err != io.EOF && err != io.ErrClosedPipe {
						t.pw.CloseWithError(err)
					}
					return
				}
			}
		}(readers[i])
	}

	t.cond.Broadcast()
}

// TODO: add a WriteTo for more efficient copy

func (t *multiReader) Close() error {
	t.readersMu.Lock()
	defer t.readersMu.Unlock()

	log.Debugf("[%p] Close on readers", t)
	for _, r := range t.readers {
		if c, ok := r.(io.Closer); ok {
			log.Debugf("[%p] Closing reader %+v", t, r)
			c.Close()
		}
	}

	return nil
}

// Remove doesn't return an error if element isn't found as the end result is
// identical
func (t *multiReader) Remove(reader io.Reader) {
	t.readersMu.Lock()
	defer t.readersMu.Unlock()

	if verbose {
		log.Debugf("[%p] removing reader - currently %d readers", t, len(t.readers))
	}

	for i, r := range t.readers {
		if r == reader {
			t.readers = append(t.readers[:i], t.readers[i+1:]...)
			// using range directly means that we're looping up, so indexes are now invalid
			if verbose {
				log.Debugf("[%p] removed reader - now %d readers", t, len(t.readers))

				for i, r := range t.readers {
					log.Debugf("[%p] Reader %d [%p]", t, i, r)
				}
			}
			break
		}
	}

	if len(t.readers) == 0 {
		t.pw.Close()
	}
}

// MultiReader returns a Reader that's the logical concatenation of
// the provided input readers.  They're read sequentially.  Once all
// inputs have returned EOF, Read will return EOF.  If any of the readers
// return a non-nil, non-EOF error, Read will return that error.
func MultiReader(readers ...io.Reader) DynamicMultiReader {
	t := &multiReader{}
	t.cond = sync.NewCond(&t.readersMu)
	t.Add(readers...)

	if verbose {
		log.Debugf("[%p] created multireader", t)
	}
	return t
}
