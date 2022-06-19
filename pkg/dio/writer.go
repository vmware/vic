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

// DynamicMultiWriter adds dynamic add/remove to the base multiwriter behaviour
type DynamicMultiWriter interface {
	io.Writer
	Add(...io.Writer)
	Remove(io.Writer)
	Close() error
}

type multiWriter struct {
	mutex   sync.Mutex
	writers []io.Writer
}

func (t *multiWriter) Write(p []byte) (int, error) {
	var n int
	var err error

	t.mutex.Lock()
	defer t.mutex.Unlock()

	eof := 0
	// possibly want to add buffering or parallelize this
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			// remove the writer
			log.Debugf("[%p] removing writer %p due to %s", t, w, err.Error())
			// Remove grabs the lock
			t.remove(w)
			if err == io.EOF {
				eof++
			}
			continue
		}

		// FIXME: figure out what semantics we need here - currently we may not write to
		// everything as we abort
		if n != len(p) {
			// remove the writer
			log.Debugf("[%p] removing writer %p due to short write: %d != %d", t, w, n, len(p))
			t.remove(w)
		}
	}

	// This means writers closed/removed while we iterate
	if eof != 0 && n == 0 && err == nil && eof == len(t.writers) {
		log.Debugf("[%p] All of the writers returned EOF (%d)", t, len(t.writers))
	}
	if verbose {
		if err != nil {
			log.Debugf("[%p] write %q to %d writers (err: %#+v)", t, string(p[:n]), len(t.writers), err)
		} else {
			log.Debugf("[%p] write %q to %d writers", t, string(p[:n]), len(t.writers))
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

// CloseWriter is an interface that implements structs
// that close input streams to prevent from writing.
type CloseWriter interface {
	CloseWrite() error
}

// FIXME: provide a mechanism for selectively closing writers
//  - currently this closes /dev/stdout and logging as well if present
func (t *multiWriter) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	log.Debugf("[%p] Close on writers", t)
	for _, w := range t.writers {
		log.Debugf("[%p] Closing writer %+v", t, w)

		if c, ok := w.(CloseWriter); ok {
			log.Debugf("[%p] is a CloseWriter", t, w)
			c.CloseWrite()
		} else if c, ok := w.(io.Closer); ok && c != os.Stdout && c != os.Stderr {
			log.Debugf("[%p] is a Closer", t, w)
			// squash closing of stdout/err if bound
			c.Close()
		}
	}

	return nil
}

func (t *multiWriter) remove(writer io.Writer) {
	wTmp := make([]io.Writer, 0, len(t.writers))

	if verbose {
		log.Debugf("[%p] removing writer %p - currently %d writers", t, writer, len(t.writers))
	}
	for _, w := range t.writers {
		if w != writer {
			wTmp = append(wTmp, w)
		}
	}
	if len(t.writers) != len(wTmp) {
		log.Debugf("[%p] removed writer - now %d writers", t, len(wTmp))
	}
	t.writers = wTmp
}

// Remove doesn't return an error if element isn't found as the end result is
// identical
func (t *multiWriter) Remove(writer io.Writer) {
	t.mutex.Lock()
	t.remove(writer)
	t.mutex.Unlock()
}

// MultiWriter extends io.MultiWriter to allow add/remove of writers dynamically
// without disrupting existing writing
func MultiWriter(writers ...io.Writer) DynamicMultiWriter {
	t := &multiWriter{writers: writers}
	if verbose {
		log.Debugf("[%p] created multiwriter", t)
	}
	return t
}
