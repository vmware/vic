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

// Package dio adds dynamic behaviour to the standard io package mutliX types
package dio

import (
	"io"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

var verbose = true

// DynamicMultiWriter adds dynamic add/remove to the base multiwriter behaviour
type DynamicMultiWriter interface {
	io.Writer
	Add(...io.Writer)
	Remove(io.Writer)
	Close() error
}

type multiWriter struct {
	mutex sync.Mutex

	writers []io.Writer
}

func (t *multiWriter) Write(p []byte) (int, error) {
	var n int
	var err error

	// stash a local copy of the slice as we never want to write twice to a single writer
	// if remove is called during this flow
	wTmp := t.writers

	if verbose {
		defer func() {
			log.Debugf("[%p] write %q to %d writers (err: %#+v)", t, string(p[:n]), len(t.writers), err)
		}()
	}

	// possibly want to add buffering or parallelize this
	for _, w := range wTmp {
		n, err = w.Write(p)
		if err != nil {
			if err != io.EOF {
				return n, err
			}

			// remove the writer
			log.Debugf("[%p] removing writer due to EOF", t)
			t.Remove(w)
		}

		// FIXME: figure out what semantics we need here - currently we may not write to
		// everything as we abort
		if n != len(p) {
			err = io.ErrShortWrite
			return n, err
		}
	}
	return len(p), nil
}

func (t *multiWriter) Add(writer ...io.Writer) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.writers = append(t.writers, writer...)
	if verbose {
		log.Debugf("[%p] added writer - now %d writers", t, len(t.writers))
	}
}

// FIXME: provide a mechanism for selectively closing writers
//  - currently this closes /dev/stdout and logging as well if present
func (t *multiWriter) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	log.Debugf("Close on writers")
	for _, w := range t.writers {
		// squash closing of stdout/err if bound
		if c, ok := w.(io.Closer); ok && c != os.Stdout && c != os.Stderr {
			log.Debugf("Closing writer %+v", w)
			c.Close()
		}
	}

	return nil
}

// TODO: add a ReadFrom for more efficient copy

// Remove doesn't return an error if element isn't found as the end result is
// identical
func (t *multiWriter) Remove(writer io.Writer) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if verbose {
		log.Debugf("[%p] removing writer - currently %d writers", t, len(t.writers))
	}
	for i, w := range t.writers {
		if w == writer {
			t.writers = append(t.writers[:i], t.writers[i+1:]...)
			// using range directly means that we're looping up, so indexes are now
			// invalid
			if verbose {
				log.Debugf("[%p] removed writer - now %d writers", t, len(t.writers))
			}
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
	Close() error
}

type multiReader struct {
	mutex sync.Mutex
	cond  *sync.Cond

	readers []io.Reader
	err     error
}

func (t *multiReader) Read(p []byte) (int, error) {
	n := 0
	if verbose {
		defer func() {
			log.Debugf("[%p] read %q from %d readers (err: %#+v)", t, string(p[:n]), len(t.readers), t.err)
		}()
	}

	if t.err == io.EOF {
		if verbose {
			log.Debugf("[%p] read from closed multi-reader, returning EOF", t)
		}
		return 0, io.EOF
	}

	// if there's no readers we are steady state - has to be after t.err check to
	// get correct Close behaviour.
	// Blocking behaviour!
	t.mutex.Lock()
	for len(t.readers) == 0 && t.err == nil {
		t.cond.Wait()
	}
	t.mutex.Unlock()

	if t.err != nil {
		return 0, t.err
	}

	eof := io.EOF
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
			// if there was an actual error, return that
			// we cannot handle multiple not EOF errors, so return now
			t.err = err
			return n, err
		}
	}

	// we'd have returned anything other than EOF/nil inline
	t.err = eof
	return n, eof
}

// TODO: add a WriteTo for more efficient copy

func (t *multiReader) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	log.Debugf("Close on readers")
	for _, r := range t.readers {
		if c, ok := r.(io.Closer); ok {
			log.Debugf("Closing reader %+v", r)
			c.Close()
		}
	}

	t.err = io.EOF
	t.cond.Broadcast()
	return nil
}

func (t *multiReader) Add(reader ...io.Reader) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.readers = append(t.readers, reader...)
	// if we've got a new reader, we're not EOF any more until that reader EOFs
	t.err = nil
	t.cond.Broadcast()

	if verbose {
		log.Debugf("[%p] adding reader - now %d readers", t, len(t.readers))
	}
}

// Remove doesn't return an error if element isn't found as the end result is
// identical
func (t *multiReader) Remove(reader io.Reader) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if verbose {
		log.Debugf("[%p] removing reader - currently %d readers", t, len(t.readers))
	}

	for i, r := range t.readers {
		if r == reader {
			t.readers = append(t.readers[:i], t.readers[i+1:]...)
			// using range directly means that we're looping up, so indexes are now invalid
			if verbose {
				log.Debugf("[%p] removed reader - currently %d readers", t, len(t.readers))
			}
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
	t := &multiReader{readers: r}
	t.cond = sync.NewCond(&t.mutex)

	if verbose {
		log.Debugf("[%p] created multireader", t)
	}
	return t
}
