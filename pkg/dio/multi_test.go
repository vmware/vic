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
	"bytes"
	"io"
	"sync"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func init() {
	// enable verbose logging during tests
	log.SetLevel(log.DebugLevel)
	verbose = true
}

func TestMultiWrite(t *testing.T) {
	pipeAR, pipeAW := io.Pipe()
	pipeBR, pipeBW := io.Pipe()

	mwriter := MultiWriter(pipeAW, pipeBW)

	var bufA bytes.Buffer
	var bufB bytes.Buffer

	// set up a copy so we don't block writes
	go io.Copy(&bufA, pipeAR)
	go io.Copy(&bufB, pipeBR)

	// send the test string
	data := "verify base multiwriter function"
	_, err := mwriter.Write([]byte(data))
	if err != nil {
		t.Error(err)
		return
	}

	// compare the data
	if bufA.String() != data {
		t.Errorf("A: expected: %s, actual: %s", data, bufA.String())
		return
	}

	if bufB.String() != data {
		t.Errorf("B: expected: %s, actual: %s", data, bufB.String())
		return
	}
}

func TestWriteAdd(t *testing.T) {
	pipeAR, pipeAW := io.Pipe()
	pipeBR, pipeBW := io.Pipe()

	mwriter := MultiWriter(pipeAW)

	var bufA bytes.Buffer
	var bufB bytes.Buffer

	// set up a copy so we don't block writes
	go io.Copy(&bufA, pipeAR)
	go io.Copy(&bufB, pipeBR)

	// send the test string
	data := "verify base multiwriter function"
	_, err := mwriter.Write([]byte(data))
	if err != nil {
		t.Error(err)
		return
	}

	// compare the data - shouldn't be present
	if bufA.String() != data {
		t.Errorf("A: expected: %s, actual: %s", data, bufA.String())
		return
	}

	if bufB.String() != "" {
		t.Errorf("B: expected: %s, actual: %s", "", bufB.String())
		return
	}

	// Add writer to existing MultiWriter
	mwriter.Add(pipeBW)

	data2 := "verify dynamic add"
	_, err = mwriter.Write([]byte(data2))
	if err != nil {
		t.Error(err)
		return
	}

	// compare the data - shouldn't be present
	if bufA.String() != data+data2 {
		t.Errorf("A: expected: %s, actual: %s", data+data2, bufA.String())
		return
	}

	if bufB.String() != data2 {
		t.Errorf("B: expected: %s, actual: %s", data2, bufB.String())
		return
	}

}

func TestWriteRemove(t *testing.T) {
	pipeAR, pipeAW := io.Pipe()
	pipeBR, pipeBW := io.Pipe()

	mwriter := MultiWriter(pipeAW, pipeBW)

	var bufA bytes.Buffer
	var bufB bytes.Buffer

	// set up a copy so we don't block writes
	go io.Copy(&bufA, pipeAR)
	go io.Copy(&bufB, pipeBR)

	// send the test string
	data := "verify base multiwriter function"
	_, err := mwriter.Write([]byte(data))
	if err != nil {
		t.Error(err)
		return
	}

	// compare the data
	if bufA.String() != data {
		t.Errorf("A: expected: %s, actual: %s", data, bufA.String())
		return
	}

	if bufB.String() != data {
		t.Errorf("B: expected: %s, actual: %s", data, bufB.String())
		return
	}

	// Add writer to existing MultiWriter
	mwriter.Remove(pipeBW)

	data2 := "verify dynamic remove"
	_, err = mwriter.Write([]byte(data2))
	if err != nil {
		t.Error(err)
		return
	}

	// compare the data - shouldn't be present
	if bufA.String() != data+data2 {
		t.Errorf("A: expected: %s, actual: %s", data+data2, bufA.String())
		return
	}

	if bufB.String() != data {
		t.Errorf("B: expected: %s, actual: %s", data, bufB.String())
		return
	}
}

func TestWriteConcurrentRemove(t *testing.T) {
	t.Skip("not sure how to test concurrency in this case")
}

func TestMultiRead(t *testing.T) {
	var wg sync.WaitGroup

	dataA := "verify base multireader functionA"
	dataB := "verify base multireader functionB"

	readerA := bytes.NewReader([]byte(dataA))
	readerB := bytes.NewReader([]byte(dataB))

	mreader := MultiReader(readerA, readerB)

	var buf bytes.Buffer

	// do the read
	_, err := io.CopyN(&buf, mreader, int64(len(dataA)+len(dataB)))
	if err != nil {
		t.Error(err)
	}
	// compare the data
	if buf.String() != dataA+dataB {
		t.Errorf("A: expected: %s, actual: %s", dataA+dataB, buf.String())
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := mreader.Read(buf.Bytes())
		if err != io.EOF {
			t.Error(err)
		}
	}()

	mreader.Close()
	wg.Wait()
}

func TestReadAdd(t *testing.T) {
	var wg sync.WaitGroup

	dataA := "verify base multireader functionA"
	dataB := "verify base multireader functionB"

	readerA := bytes.NewReader([]byte(dataA))
	readerB := bytes.NewReader([]byte(dataB))

	mreader := MultiReader(readerA)

	var bufA bytes.Buffer
	var bufB bytes.Buffer

	// do the read - bytes.NewReader does not return data and EOF
	// from the same call, so this should have err==nil
	_, err := io.CopyN(&bufA, mreader, int64(len(dataA)))
	if err != nil {
		t.Error(err)
	}

	// compare the data
	if bufA.String() != dataA {
		t.Errorf("A: expected: %s, actual: %s", dataA, bufA.String())
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := mreader.Read(bufA.Bytes())
		if err != io.EOF {
			t.Error(err)
		}
	}()

	// Add reader to existing MultiReader
	// this should furnish new data to the copy without further action being
	// taken
	mreader.Add(readerB)

	// do the read - we expect mreader to now switch to the new source, which
	// has the standard bytes.NewReader behaviour
	_, err = io.CopyN(&bufB, mreader, int64(len(dataB)))
	if err != nil {
		t.Error(err)
	}

	// compare the data
	if bufB.String() != dataB {
		t.Errorf("A: expected: %s, actual: %s", dataB, bufB.String())
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := mreader.Read(bufB.Bytes())
		if err != io.EOF {
			t.Error(err)
		}
	}()

	mreader.Close()
	wg.Wait()
}

func TestReadRemove(t *testing.T) {
	var wg sync.WaitGroup

	pipeAR, pipeAW := io.Pipe()
	pipeBR, pipeBW := io.Pipe()

	mreader := MultiReader(pipeAR, pipeBR)

	var buf bytes.Buffer

	// send the test string
	data := "verify base multiwriter function"

	go func() {
		defer pipeAW.Close()
		_, err := pipeAW.Write([]byte(data))
		if err != nil {
			t.Error(err)
		}
	}()

	// do the read
	_, err := io.CopyN(&buf, mreader, int64(len(data)))
	if err != nil {
		t.Error(err)
	}

	// compare the data
	if buf.String() != data {
		t.Errorf("A: expected: %s, actual: %s", data, buf.String())
	}
	buf.Truncate(0)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := mreader.Read(buf.Bytes())
		if err != io.EOF {
			t.Error(err)
		}
	}()

	// Add writer to existing MultiWriter
	mreader.Remove(pipeAR)

	data = "verify dynamic remove"

	go func() {
		defer pipeAW.Close()
		_, err = pipeBW.Write([]byte(data))
		if err != nil {
			t.Error(err)
		}
	}()

	// do the read
	_, err = io.CopyN(&buf, mreader, int64(len(data)))
	if err != nil {
		t.Error(err)
	}

	// compare the data - shouldn't be present
	if buf.String() != data {
		t.Errorf("A: expected: %s, actual: %s", data, buf.String())
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := mreader.Read(buf.Bytes())
		if err != io.EOF {
			t.Error(err)
		}
	}()

	mreader.Close()
	wg.Wait()
}

func TestReadConcurrentRemove(t *testing.T) {
	t.Skip("not sure how to test concurrency in this case")
}
