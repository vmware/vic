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
	"testing"
)

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

func TestAdd(t *testing.T) {
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

func TestRemove(t *testing.T) {
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

func TestConcurrentRemove(t *testing.T) {
	t.Skip("not sure how to test concurrency in this case")
}
