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

package dio

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

const (
	count = 32

	base    = "base functionality"
	dynamic = "dynamic add/remove functionality"
)

func write(t *testing.T, mwriter DynamicMultiWriter, p []byte) {
	n, err := mwriter.Write(p)
	if err != nil {
		t.Errorf("write: %v", err)
	}
	if n != len(p) {
		t.Errorf("short write: %d != %d", n, len(p))
	}
}

func read(t *testing.T, mreader DynamicMultiReader, limit int) []byte {
	total := 0

	var buf = make([]byte, 32*1024)
	for {
		n, err := mreader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("read: %v", err)
		}

		total += n
		if total >= limit {
			break
		}
	}

	return buf[:limit]
}

func control(t *testing.T, buffers []*bytes.Buffer, p []byte) {
	for i := range buffers {
		buffer := buffers[i]
		if buffer.String() != string(p) {
			t.Errorf("Expected: %q, actual: %q", string(p), buffer.String())
		}
	}
}

func equal(t *testing.T, mreader DynamicMultiReader, str string, count int) {
	expected := strings.Repeat(str, count)
	// read from multi reader
	buffer := string(read(t, mreader, len(expected)))
	// compare the data
	if buffer != expected {
		t.Errorf("Expected: %q, actual: %q", expected, buffer)
	}
}

// TestMultiWrite creates multi writers and writes to them, then removes some of them and writes again
func TestMultiWrite(t *testing.T) {
	var writers []io.Writer
	var buffers []*bytes.Buffer

	// create & initialize writers and buffers
	for i := 0; i < count; i++ {
		var buffer bytes.Buffer

		reader, writer := io.Pipe()

		writers = append(writers, writer)
		buffers = append(buffers, &buffer)

		// set up a goroutine so we don't block writes
		go io.Copy(&buffer, reader)
	}

	// create the multi writer
	mwriter := MultiWriter(writers...)

	write(t, mwriter, []byte(base))
	control(t, buffers, []byte(base))
}

// TestMultiWrite creates bunch of multi writers and writes to them, then adds more and writes again
func TestWriteAdd(t *testing.T) {
	var writers []io.Writer
	var buffers []*bytes.Buffer

	var writersAdded []io.Writer
	var buffersAdded []*bytes.Buffer

	// create & initialize writers and buffers into two categories
	// *Added ones will be added to multi writer later
	for i := 0; i < count; i++ {
		var buffer bytes.Buffer

		reader, writer := io.Pipe()

		// add the writer & buffer to skip list if i divisible by three
		if i%3 == 0 {
			writersAdded = append(writersAdded, writer)
			buffersAdded = append(buffersAdded, &buffer)
		} else {
			writers = append(writers, writer)
			buffers = append(buffers, &buffer)
		}

		// set up a goroutine so we don't block writes
		go io.Copy(&buffer, reader)
	}

	// create the multi writer
	mwriter := MultiWriter(writers...)

	write(t, mwriter, []byte(base))
	control(t, buffers, []byte(base))

	// add skipped writers to the writer
	mwriter.Add(writersAdded...)

	write(t, mwriter, []byte(dynamic))
	control(t, buffers, []byte(base+dynamic))
	control(t, buffersAdded, []byte(dynamic))
}

// TestMultiWrite creates multi writers and writes to them, then removes some of them and writes again
func TestWriteRemove(t *testing.T) {
	var writers []io.Writer
	var buffers []*bytes.Buffer

	var buffersLeft []*bytes.Buffer
	var buffersRemoved []*bytes.Buffer

	// create & initialize writers and buffers into two categories
	// *Removed ones will be filtered out from multi writer later
	for i := 0; i < count; i++ {
		var buffer bytes.Buffer

		reader, writer := io.Pipe()

		writers = append(writers, writer)
		buffers = append(buffers, &buffer)

		// set up a goroutine so we don't block writes
		go io.Copy(&buffer, reader)
	}

	// create the multi writer
	mwriter := MultiWriter(writers...)

	write(t, mwriter, []byte(base))
	control(t, buffers, []byte(base))

	// add the writer & buffer to skip list if i divisible by three
	for i := 0; i < count; i++ {
		if i%3 == 0 {
			mwriter.Remove(writers[i])

			buffersRemoved = append(buffersRemoved, buffers[i])
		} else {
			buffersLeft = append(buffersLeft, buffers[i])
		}
	}

	write(t, mwriter, []byte(dynamic))
	control(t, buffersLeft, []byte(base+dynamic))
	control(t, buffersRemoved, []byte(base))
}

// TestMultiRead creates multi readers and reads from them
func TestMultiRead(t *testing.T) {
	var readers []io.Reader

	// create & initialize writers and buffers
	for i := 0; i < count; i++ {
		reader, writer := io.Pipe()

		readers = append(readers, reader)

		// set up a  goroutine so we don't block reads
		go io.Copy(writer, bytes.NewReader([]byte(base)))
	}

	// create the multi writer
	mreader := MultiReader(readers...)

	equal(t, mreader, base, count)
}

// TestMultiRead creates multi readers and reads from them, then adds mores and reads again
func TestReadAdd(t *testing.T) {
	var readers []io.Reader
	var pipereaders []io.PipeReader

	var readersAdded []io.Reader

	skipped := 0
	// create & initialize writers and buffers
	for i := 0; i < count; i++ {
		reader, writer := io.Pipe()

		if i%3 == 0 {
			skipped++

			readersAdded = append(readersAdded, reader)
			// set up a  goroutine so we don't block reads
			go io.Copy(writer, bytes.NewReader([]byte(dynamic)))
		} else {
			pipereaders = append(pipereaders, *reader)
			readers = append(readers, reader)

			// set up a  goroutine so we don't block reads
			go io.Copy(writer, bytes.NewReader([]byte(base)))
		}
	}

	// create the multi writer
	mreader := MultiReader(readers...)

	equal(t, mreader, base, count-skipped)

	// add the rest of the readers
	mreader.Add(readersAdded...)

	// close the initial set otherwise they will block
	for i := range pipereaders {
		pipereaders[i].Close()
	}

	equal(t, mreader, dynamic, skipped)
}

// TestReadRemove creates multi readers and reads from them, then removes some and reads again
func TestReadRemove(t *testing.T) {
	var readers []io.Reader
	var writers []io.Writer

	// create & initialize writers and buffers
	for i := 0; i < count; i++ {
		reader, writer := io.Pipe()

		readers = append(readers, reader)
		writers = append(writers, writer)

		// set up a goroutine so we don't block reads
		go io.Copy(writer, bytes.NewReader([]byte(base)))
	}

	// create the multi writer
	mreader := MultiReader(readers...)

	equal(t, mreader, base, count)

	removed := 0
	for i := 0; i < count; i++ {
		if i%3 == 0 {
			removed++

			mreader.Remove(readers[i])
		} else {
			// set up another goroutine so we don't block reads
			go io.Copy(writers[i], bytes.NewReader([]byte(dynamic)))
		}
	}
	equal(t, mreader, dynamic, count-removed)
}
