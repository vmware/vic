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
	io.WriteCloser
	Add(...io.WriteCloser)
	Remove(io.WriteCloser)
}

type multiWriter struct {
	mutex sync.Mutex

	writers []io.WriteCloser
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

func (t *multiWriter) Add(writer ...io.WriteCloser) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.writers = append(t.writers, writer...)
}

// Remove doesn't return an error if element isn't found as the end result is
// identical
func (t *multiWriter) Remove(writer io.WriteCloser) {
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

func (t *multiWriter) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var err error
	for _, w := range t.writers {
		err2 := w.Close()
		if err2 != nil && err == nil {
			err = err2
		}
	}

	return err
}

// MultiWriter extends io.MultiWriter to allow add/remove of writers dynamically
// without disrupting existing writing
func MultiWriter(writers ...io.WriteCloser) DynamicMultiWriter {
	w := make([]io.WriteCloser, len(writers))
	copy(w, writers)
	return &multiWriter{writers: w}
}
