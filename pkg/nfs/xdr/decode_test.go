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

package xdr

import (
	"bytes"
	"encoding/binary"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestRead(t *testing.T) {
	type X struct {
		A, B, C, D uint32
	}
	x := new(X)
	b := []byte{
		0, 0, 0, 1,
		0, 0, 0, 2,
		0, 0, 0, 3,
		0, 0, 0, 4,
		1,
	}
	buf := bytes.NewBuffer(b)
	Read(buf, x)
}

func TestByteSlice(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// byte slices have a length field up front, followed by the data.  The
	// data is aligned to 4B.
	type ByteSlice struct {
		Length uint32
		Data   []byte
		Pad    []byte
	}

	in := &ByteSlice{
		Length: 6,
		Data:   []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5},
		Pad:    []byte{0x0, 0x0},
	}

	b := &bytes.Buffer{}
	binary.Write(b, binary.BigEndian, uint32(in.Length))
	b.Write(in.Data)
	b.Write(in.Pad)

	var out []byte
	if err := Read(b, &out); err != nil {
		t.Log("fail in read")
		t.Fail()
		return
	}

	if len(out) != int(in.Length) {
		t.Logf("legth mismatch, expected %d, actual %d", in.Length, len(out))
		t.Fail()
		return
	}

	if bytes.Compare(in.Data, out) != 0 {
		t.FailNow()
	}
}
