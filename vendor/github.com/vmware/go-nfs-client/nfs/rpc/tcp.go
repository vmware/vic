// Copyright © 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause
//
package rpc

import (
	"bytes"
	"encoding/binary"
	"io"
	"sync"
)

type tcpTransport struct {
	r  io.Reader
	wc io.WriteCloser

	rlock, wlock sync.Mutex
}

// Get the response from the conn, buffer the contents, and return a reader to
// it.
func (t *tcpTransport) recv() (io.ReadSeeker, error) {
	t.rlock.Lock()
	defer t.rlock.Unlock()

	var hdr uint32
	if err := binary.Read(t.r, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}

	buf := make([]byte, hdr&0x7fffffff)
	if _, err := io.ReadFull(t.r, buf); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf), nil
}

func (t *tcpTransport) Write(buf []byte) (int, error) {
	t.wlock.Lock()
	defer t.wlock.Unlock()

	var hdr uint32 = uint32(len(buf)) | 0x80000000
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, hdr)
	n, err := t.wc.Write(append(b, buf...))

	return n, err
}

func (t *tcpTransport) Close() error {
	return t.wc.Close()
}
