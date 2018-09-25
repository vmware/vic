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
	"context"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
)

// generateArchive takes a number of files and a file size and generates a tar archive
// based on that data. It returns:
// index of entry names in the archive
// archive byte stream
func generateArchive(t *testing.T, name string, num, size int) ([]string, io.Reader) {
	r, w := io.Pipe()
	tw := tar.NewWriter(w)

	// stable reference for expected archive entries
	index := make([]string, num)
	for i := 0; i < num; i++ {
		index[i] = string(uid.New())
	}

	// our file contents are zeros - this is the worst case for stripping trailing headers
	zeros := make([]byte, size)

	go func(t *testing.T) {
		for i := 0; i < num; i++ {
			// we only really care about two things in the header
			hdr := &tar.Header{
				Name: index[i],
				Size: int64(size),
			}

			t.Logf("Writing header for file %s:%d", name, i)
			err := tw.WriteHeader(hdr)
			if err != nil {
				panic(err)
			}

			t.Logf("Writing data for file %s:%d", name, i)
			n, err := tw.Write(zeros)
			if err != nil {
				panic(err)
			}
			if n != size {
				panic(fmt.Sprintf("Failed to write all bytes: %d != %d", n, size))
			}
		}

		// add the end-of-archive
		tw.Close()
		w.Close()
	}(t)

	return index, r
}

// TestSingleStripper ensures that basic function (stripping end-of-archive footer) works as
// expected. I found no real way, when using the TarReader to actually assert that the footer
// has been dropped so this is left to assert correct passthrough of archive data and the
// dropping of the footer is asserted by the multistream cases.
func TestSingleStripper(t *testing.T) {
	size := 2048
	count := 5
	index, reader := generateArchive(t, "single", count, size)

	source := tar.NewReader(reader)
	stripper := NewStripper(trace.NewOperation(context.Background(), "strip"), source, nil)

	pr, pw := io.Pipe()
	tr := tar.NewReader(pr)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		n, err := io.Copy(pw, stripper)
		pw.Close()

		wg.Done()
		t.Logf("Done copying from stripper: %d, %s", n, err)
		if !assert.NoError(t, err, "Expected nil error from io.Copy on end-of-file") {
			t.FailNow()
		}

	}()

	for i := 0; i <= len(index); i++ {
		t.Logf("Reading header for file %d", i)
		header, err := tr.Next()
		if err == io.EOF {
			fmt.Printf("End-of-file")
			// TODO: is this pass or fail?
			return
		}

		require.NoError(t, err, "Error from archive: %s", err)
		require.NotEqual(t, i, len(index), "Expected EOF after index exhausted")
		require.Equal(t, header.Name, index[i], "Expected header name to match index")
		require.Equal(t, header.Size, int64(size), "Expected file size in header to match target generated size")

		// make the buf just that little bit bigger to allow for errrors in the copy if they would occur
		buf := make([]byte, size+10)

		t.Logf("Reading data for file %d", i)
		n, err := tr.Read(buf)

		require.NoError(t, err, "No expected errors from file data copy")
		require.Equal(t, n, size, "Expected file data size to match target generated size")
	}

	wg.Wait()
}

// TestConjoinedStrippers ensures that the footer is correctly dropped from a stripped archive
// and that a TarReader continues to provide headers from the following conjoined streams.
func TestConjoinedStrippers(t *testing.T) {
	size := 2048
	count := 3
	multiplicity := 3

	indices := make([][]string, multiplicity)
	strippers := make([]io.Reader, multiplicity)

	for m := 0; m < multiplicity; m++ {
		var reader io.Reader
		indices[m], reader = generateArchive(t, fmt.Sprintf("archive-%d", m), count, size)
		source := tar.NewReader(reader)
		strippers[m] = NewStripper(trace.NewOperation(context.Background(), fmt.Sprintf("strip-%d", m)), source, nil)
	}

	conjoined := MultiReader(strippers...)

	pr, pw := io.Pipe()
	tr := tar.NewReader(pr)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		n, err := io.Copy(pw, conjoined)
		pw.Close()

		wg.Done()
		t.Logf("Done copying from strippers: %d, %s", n, err)
		require.NoError(t, err, "Expected nil error from io.Copy on end-of-file")
	}()

	expectedEntries := count * multiplicity
	for i := 0; i <= expectedEntries; i++ {
		t.Logf("Reading header for file %d", i)
		header, err := tr.Next()
		if err == io.EOF {
			t.Logf("End-of-file")
			// TODO: is this pass or fail?
			return
		}

		require.NoError(t, err, "Error from archive: %s", err)
		require.NotEqual(t, i, expectedEntries, "Expected EOF after index exhausted")

		member := i / count
		entry := i % count

		require.Equal(t, header.Name, indices[member][entry], "Expected header name to match index")
		require.Equal(t, header.Size, int64(size), "Expected file size in header to match target generated size")

		// make the buf just that little bit bigger to allow for errrors in the copy if they would occur
		buf := make([]byte, size+10)

		t.Logf("Reading data for file %d", i)
		n, err := tr.Read(buf)

		require.NoError(t, err, "No expected errors from file data copy")
		require.Equal(t, n, size, "Expected file data size to match target generated size")
	}

	wg.Wait()
}

// TestConjoinedStrippersWithCloser ensures that we can conjoin readers, multiple strippers and a regular, in order to get
// the end-of-archive footer as the finale. We have a wait group to ensure that all routines have finished before declaring
// success.
func TestConjoinedStrippersWithCloser(t *testing.T) {
	size := 2048
	count := 3
	multiplicity := 3

	indices := make([][]string, multiplicity)
	readers := make([]io.Reader, multiplicity)

	for m := 0; m < multiplicity; m++ {
		var reader io.Reader
		indices[m], reader = generateArchive(t, fmt.Sprintf("archive-%d", m), count, size)
		source := tar.NewReader(reader)

		if m < multiplicity-1 {
			t.Logf("added stripper")
			readers[m] = NewStripper(trace.NewOperation(context.Background(), fmt.Sprintf("strip-%d", m)), source, nil)
		} else {
			t.Logf("added raw")
			readers[m] = reader
		}
	}

	conjoined := MultiReader(readers...)

	pr, pw := io.Pipe()
	tr := tar.NewReader(pr)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		n, err := io.Copy(pw, conjoined)
		pw.Close()

		wg.Done()
		t.Logf("Done copying from all sources: %d, %s", n, err)
		require.NoError(t, err, "Expected nil error from io.Copy on end-of-file")
	}()

	expectedEntries := count * multiplicity
	for i := 0; i <= expectedEntries; i++ {
		t.Logf("Reading header for file %d", i)
		header, err := tr.Next()
		if err == io.EOF {
			t.Logf("End-of-file")

			wg.Wait()
			return
		}

		require.NoError(t, err, "Error from archive: %s", err)
		require.NotEqual(t, i, expectedEntries, "Expected EOF after index exhausted")

		member := i / count
		entry := i % count

		require.Equal(t, header.Name, indices[member][entry], "Expected header name to match index")
		require.Equal(t, header.Size, int64(size), "Expected file size in header to match target generated size")

		// make the buf just that little bit bigger to allow for errrors in the copy if they would occur
		buf := make([]byte, size+10)

		t.Logf("Reading data for file %d", i)
		n, err := tr.Read(buf)

		require.NoError(t, err, "No expected errors from file data copy")
		require.Equal(t, n, size, "Expected file data size to match target generated size")
	}

	wg.Wait()
}
