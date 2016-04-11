// Copyright 2016 VMware, Inc. All Rights Reserved.
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

// Package io adds dynamic behaviour to the standard io package mutliX types
package dio

import (
	"io"
	"sync"
)

// DynamicMultiWriter adds dynamic add/remove to the base multiwriter behaviour
type DynamicMultiWriter interface {
	io.Writer
	Add(...io.Writer)
	Remove(io.Writer)
}

type multiWriter struct {
	mutex sync.Mutex

	writers []io.Writer
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	// stash a local copy of the slice
	wTmp := t.writers

	// possibly want to add buffering or parallelize this
	for _, w := range wTmp {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

func (t *multiWriter) Add(writer ...io.Writer) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.writers = append(t.writers, writer...)
}

// Remove doesn't return an error if element isn't found as the end result is
// identical
func (t *multiWriter) Remove(writer io.Writer) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for i, w := range t.writers {
		if w == writer {
			t.writers = append(t.writers[:i], t.writers[i+1:]...)
			// using range directly means that we're looping up, so indexes are now
			// invalid
			return
		}
	}
}

// MultiWriter extends io.MultiWriter to allow add/remove of writers dynamically
// without disrupting existing writing
func MultiWriter(writers ...io.Writer) DynamicMultiWriter {
	w := make([]io.Writer, len(writers))
	copy(w, writers)
	return &multiWriter{writers: w}
}

// DynamicMultiReader adds dynamic add/remove to the base multireader behaviour
type DynamicMultiReader interface {
	io.Reader
	Add(...io.Reader)
	Remove(io.Reader)
}

type multiReader struct {
	mutex sync.Mutex

	readers []io.Reader
}

func (t *multiReader) Read(p []byte) (int, error) {
	eof := io.EOF
	n := 0
	for _, r := range t.readers {
		slice := p[n:]
		if len(slice) == 0 {
			// we've run out of target space and don't know what
			// the remaining readers have, so not EOF
			return n, nil
		}

		x, err := r.Read(slice)
		n += x
		// if any of the readers don't return EOF, it's not EOF
		if err == nil {
			eof = nil
		} else if err != io.EOF {
			// if the was an actual error, return that
			// we cannot handle multiple not EOF errors, so return now
			return n, err
		}
	}

	// we'd have returned anything other than EOF/nil inline
	return n, eof
}

func (t *multiReader) Add(reader ...io.Reader) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.readers = append(t.readers, reader...)
}

// Remove doesn't return an error if element isn't found as the end result is
// identical
func (t *multiReader) Remove(reader io.Reader) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for i, r := range t.readers {
		if r == reader {
			t.readers = append(t.readers[:i], t.readers[i+1:]...)
			// using range directly means that we're looping up, so indexes are now
			// invalid
			return
		}
	}
}

// MultiReader returns a Reader that's the logical concatenation of
// the provided input readers.  They're read sequentially.  Once all
// inputs have returned EOF, Read will return EOF.  If any of the readers
// return a non-nil, non-EOF error, Read will return that error.
func MultiReader(readers ...io.Reader) DynamicMultiReader {
	r := make([]io.Reader, len(readers))
	copy(r, readers)
	return &multiReader{readers: r}
}
