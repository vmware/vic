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

package rpc

import (
	"encoding/binary"
	"io"
	"sync"
)

type tcpTransport struct {
	io.Reader
	io.WriteCloser
	rlock, wlock sync.Mutex
}

func (t *tcpTransport) recv() ([]byte, error) {
	t.rlock.Lock()
	defer t.rlock.Unlock()
	var hdr uint32
	if err := binary.Read(t, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}
	buf := make([]byte, hdr&0x7fffffff)
	if _, err := io.ReadFull(t, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (t *tcpTransport) send(buf []byte) error {
	t.wlock.Lock()
	defer t.wlock.Unlock()
	var hdr uint32 = uint32(len(buf)) | 0x80000000
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, hdr)
	_, err := t.WriteCloser.Write(append(b, buf...))
	return err
}
