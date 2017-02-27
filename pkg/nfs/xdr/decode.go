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
	"encoding/binary"
	"io"

	xdr "github.com/davecgh/go-xdr/xdr2"
)

func Uint32(b []byte) (uint32, []byte) {
	return binary.BigEndian.Uint32(b[0:4]), b[4:]
}

func Opaque(b []byte) ([]byte, []byte) {
	l, b := Uint32(b)
	return b[:l], b[l:]
}

func Uint32List(b []byte) ([]uint32, []byte) {
	l, b := Uint32(b)
	v := make([]uint32, l)
	for i := 0; i < int(l); i++ {
		v[i], b = Uint32(b)
	}
	return v, b
}

func Read(r io.Reader, val interface{}) error {
	_, err := xdr.Unmarshal(r, val)
	return err
}
