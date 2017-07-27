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

package archive

import (
	"archive/tar"
	"io"

	"github.com/vmware/vic/pkg/trace"
)

// Stripper strips the end-of-archive entries from a tar stream
type Stripper struct {
	// op allows threaded tracing
	op trace.Operation

	// be opinionated about the type of the source
	source *tar.Reader
}

// NewStripper returns a WriterTo that will strip the trailing end-of-archive bytes
// from the supplied tar stream.
// It implements io.Reader only so that it can be passed to io.Copy
func NewStripper(op trace.Operation, reader *tar.Reader) *Stripper {
	return &Stripper{
		op:     op,
		source: reader,
	}
}

// Read is implemented solely so this can be provided to io.Copy as an io.Reader.
// This works on the assumption of io.Copy making use of the WriterTo implementation.
func (s *Stripper) Read(b []byte) (int, error) {
	panic("io.Reader usage not supported - intended use is as io.WriterTo")
}

// WriteTo is the primary function, allowing easy use of the underlying tar stream without
// requiring chunking and assocated tracking to another buffer size.
func (s *Stripper) WriteTo(w io.Writer) (sum int64, err error) {
	// TODO: should we nil s.source on error then handle a post-error call? What's the expected
	// semantic?

	tw := tar.NewWriter(w)

	for {
		var header *tar.Header
		header, err = s.source.Next()
		if err == io.EOF {
			// do NOT call tarwriter.Close()
			s.op.Debugf("Stripper dropping end of archive")
			return
		}

		if err != nil {
			s.op.Errorf("Error reading archive header: %s", err)
			return
		}

		err = tw.WriteHeader(header)
		if err != nil {
			s.op.Errorf("Error writing tar header: %s", err)
			return
		}

		var n int64
		n, err = io.Copy(tw, s.source)
		sum += n
		if err != nil {
			s.op.Errorf("Error copying file data: %s", err)
			return
		}

		tw.Flush()
	}
}

// eofReader copied from io package to support MultiWriterTo variant
type eofReader struct{}

func (eofReader) WriteTo(w io.Writer) (int64, error) {
	return 0, io.EOF
}

func (eofReader) Read(b []byte) (int, error) {
	return 0, io.EOF
}

func (eofReader) ReadFrom(r io.Reader) (int64, error) {
	return 0, io.EOF
}

// multiStripper based off io.MultiReader but delegating to io.WriterTo
// instead of performing buffer copy
type multiReader struct {
	readers []io.Reader
}

func (mr *multiReader) Read(p []byte) (n int, err error) {
	panic("io.Reader usage not supported - intended use is as io.WriterTo")
}

func (mr *multiReader) WriteTo(w io.Writer) (sum int64, err error) {
	for len(mr.readers) > 0 {
		var n int64

		n, err = io.Copy(w, mr.readers[0])
		sum += n

		// io.Copy strips EOF
		if err == io.EOF || err == nil {
			mr.readers[0] = eofReader{} // permit earlier GC
			mr.readers = mr.readers[1:]
		}

		if n > 0 || err != io.EOF {
			if err == io.EOF && len(mr.readers) > 0 {
				// Don't return EOF yet. More readers remain.
				err = nil
				continue
			}

			break
		}
	}

	if err == nil {
		err = io.EOF
	}
	return
}

// Close allows this to be a Closer as well - specific to expected usage but necessary.
func (mr *multiReader) Close() error {
	return nil
}

// MultiReader is based off the io.MultiReader but will make use of WriteTo or
// ReadFrom delegation and ONLY supports usage via the WriteTo method on itself.
// It is specifically intended to be passed to io.Copy
func MultiReader(readers ...io.Reader) io.ReadCloser {
	r := make([]io.Reader, len(readers))
	copy(r, readers)
	return &multiReader{r}
}
