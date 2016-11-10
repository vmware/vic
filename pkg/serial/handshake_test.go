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

package serial

import (
	"errors"
	"sync"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
)

type readErr struct {
	err error
	n   int
}

type BlockingSendReceiver struct {
	c        chan byte
	deadline chan struct{}
}

func NewBlockingSendReceiver() *BlockingSendReceiver {
	return &BlockingSendReceiver{
		c:        make(chan byte, 10240),
		deadline: make(chan struct{}, 1),
	}
}

func (f *BlockingSendReceiver) Send(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		f.c <- b[i]
	}
	return 0, nil
}

func (f *BlockingSendReceiver) Timeout(d time.Duration) *BlockingSendReceiver {
	go func() {
		time.Sleep(d)
		f.deadline <- struct{}{}
	}()
	return f
}

func (f *BlockingSendReceiver) Receive(b []byte) (int, error) {
	select {
	case <-f.deadline:
		return 0, errors.New("Timeout error")
	default:
	}
	count := 0
	for count < len(b) {
		select {
		case v := <-f.c:
			b[count] = v
			count++
		case _, ok := <-f.deadline:
			if ok {
				close(f.deadline)
			}
			return 0, errors.New("Timeout error")
		default:
			if count == 0 {
				time.Sleep(time.Millisecond)
			} else {
				return count, nil
			}
		}
	}
	return count, nil
}

type BiChannel struct {
	L *BlockingSendReceiver
	R *BlockingSendReceiver
}

func (bc *BiChannel) Write(b []byte) (int, error) {
	return bc.L.Send(b)
}

func (bc *BiChannel) Read(b []byte) (int, error) {
	return bc.R.Receive(b)
}

func NewFakeConnection(t time.Duration) (*BiChannel, *BiChannel) {
	l := NewBlockingSendReceiver().Timeout(t)
	r := NewBlockingSendReceiver().Timeout(t)
	return &BiChannel{L: l, R: r}, &BiChannel{L: r, R: l}
}

func TestHandshakeServerNormalCaseScenario(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	clientConn, serverConn := NewFakeConnection(time.Second)

	go func() {
		buf := make([]byte, 10)
		clientConn.Write([]byte{flagSyn, 200})

		if n, e := clientConn.Read(buf); e != nil || n != 3 {
			t.Errorf("Only 3 bytes are expected: %x, received: %d", buf[:n], n)
			return
		}

		if buf[0] != flagAck || buf[1] != 201 || buf[2] < 128 {
			t.Errorf("Error, unexpected data: %x", buf[:3])
			return
		}
		clientConn.Write([]byte{flagAck, incrementByte(buf[2])})
	}()

	if e := HandshakeServer(serverConn); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}
}

func TestHandshakeServerLotsOfTrashOnTheLine(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	clientConn, serverConn := NewFakeConnection(time.Second * 10)

	go func() {
		buf := make([]byte, 10)

		// Do not send too many bytes, otherwise "write" will block on server side due too many flagNak.
		x := "sdfkgn sdflkjsfdfdgis dfgs"
		for i := 0; i < 5; i++ {
			x += x
		}
		clientConn.Write([]byte(x))
		clientConn.Write([]byte{flagSyn, 200})

		for {
			n, e := clientConn.Read(buf)
			if e != nil {
				t.Errorf("Unexpected server error: %v", e)
				return
			}

			if n < 3 {
				continue
			}

			data := buf[n-3:]
			if data[0] == flagNak {
				continue
			}

			if data[0] != flagAck || data[1] != 201 || data[2] < 128 {
				t.Errorf("Error, unexpected data: %x", data[:3])
				return
			}
			clientConn.Write([]byte{flagAck, incrementByte(data[2])})
		}
	}()

	for {
		if e := HandshakeServer(serverConn); e == nil {
			break
		} else {
			if _, ok := e.(*HandshakeError); !ok {
				t.Errorf("Unexpected error: %v", e)
				return
			}
		}
	}
}

func TestHandshakeServerComportSync(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	clientConn, serverConn := NewFakeConnection(time.Second * 3)

	go func() {
		buf := make([]byte, 10)
		clientConn.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22})

		clientConn.Write([]byte{flagSyn, 200})

		for {
			n, e := clientConn.Read(buf)
			if e != nil {
				t.Errorf("Unexpected server error: %v", e)
				return
			}

			if n < 3 {
				continue
			}

			data := buf[n-3:]
			if data[0] == flagNak {
				continue
			}

			if data[0] != flagAck || data[1] != 201 || data[2] < 128 {
				t.Errorf("Error, unexpected data: %x", data[:3])
				return
			}
			clientConn.Write([]byte{flagAck, incrementByte(data[2])})
		}
	}()

	for {
		if e := HandshakeServer(serverConn); e == nil {
			break
		} else {
			if _, ok := e.(*HandshakeError); !ok {
				t.Errorf("Unexpected error: %v", e)
				return
			}
		}
	}
}

func TestHandshakeServerAckNakResponse(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	clientConn, serverConn := NewFakeConnection(time.Second * 3)

	go func() {
		buf := make([]byte, 10)

		// Do not send too many bytes, otherwise "write" will block on server side due too many flagNak.
		clientConn.Write([]byte{flagSyn, 200})

		n, e := clientConn.Read(buf)
		if e != nil {
			t.Errorf("Unexpected server error: %v", e)
			return
		}

		data := buf[n-3:]
		if data[0] != flagAck || data[1] != 201 || data[2] < 128 {
			t.Errorf("Error, unexpected data: %x", data[:3])
			return
		}
		// intentional error. data[2] has to be incremented.
		clientConn.Write([]byte{flagAck, data[2]})
		if n, err := clientConn.Read(buf); n != 1 || err != nil || buf[0] != flagNak {
			t.Errorf("Unexpected data or error %d, %v", n, err)
		}

		clientConn.Write([]byte{flagSyn, 200})

		n, e = clientConn.Read(buf)
		if e != nil {
			t.Errorf("Unexpected server error: %v", e)
			return
		}

		data = buf[n-3:]
		if data[0] != flagAck || data[1] != 201 || data[2] < 128 {
			t.Errorf("Error, unexpected data: %x", data[:3])
			return
		}

		// intentional error. 99 in a wrong code.
		clientConn.Write([]byte{99, incrementByte(data[2])})
		if n, err := clientConn.Read(buf); n != 1 || err != nil || buf[0] != flagNak {
			t.Errorf("Unexpected data or error %d, %v", n, err)
		}

	}()

	e := HandshakeServer(serverConn)
	if e != nil {
		if _, ok := e.(*HandshakeError); !ok {
			t.Errorf("Unexpected error: %v", e)
			return
		}
	}
	if e != nil {
		e = HandshakeServer(serverConn)
		if _, ok := e.(*HandshakeError); !ok {
			t.Errorf("Unexpected error: %v", e)
		}
	}
}

func TestHandshakeClientNormalConnection(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	clientConn, serverConn := NewFakeConnection(time.Second * 3)

	go func() {
		pos := byte(200)
		buf := make([]byte, 1024)
		if n, err := serverConn.Read(buf); n != 2 || err != nil || buf[0] != flagSyn || buf[1] < 128 {
			t.Errorf("Unexpected data or error %d, %v", n, err)
		}
		serverConn.Write([]byte{flagAck, incrementByte(buf[1]), pos})

		if n, err := serverConn.Read(buf); n != 2 || err != nil || buf[0] != flagAck || buf[1] != incrementByte(pos) {
			t.Errorf("Unexpected data or error %d, %v", n, err)
		}
	}()

	e := HandshakeClient(clientConn, false)
	if e != nil {
		t.Errorf("Unexpected error: %v", e)
	}
}

func TestHandshakeClientWrongServerAckPos(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	clientConn, serverConn := NewFakeConnection(time.Second * 3)

	go func() {
		pos := byte(200)
		buf := make([]byte, 1024)
		if n, err := serverConn.Read(buf); n != 2 || err != nil || buf[0] != flagSyn || buf[1] < 128 {
			t.Errorf("Unexpected data or error %d, %v", n, err)
		}

		// writing the wrong buf[1] that supposed to be incremented.
		serverConn.Write([]byte{flagAck, buf[1], pos})

		if n, err := serverConn.Read(buf); n != 2 || err != nil || buf[0] != flagAck || buf[1] != incrementByte(pos) {
			t.Errorf("Unexpected data or error %d, %v", n, err)
		}
	}()

	err, ok := HandshakeClient(clientConn, false).(*HandshakeError)
	if !ok {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestHandshakeClientWrongServerAck(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	clientConn, serverConn := NewFakeConnection(time.Second * 3)

	go func() {
		pos := byte(200)
		buf := make([]byte, 1024)
		if n, err := serverConn.Read(buf); n != 2 || err != nil || buf[0] != flagSyn || buf[1] < 128 {
			t.Errorf("Unexpected data or error %d, %v", n, err)
		}

		// writing 90 instead of flagAck
		serverConn.Write([]byte{90, incrementByte(buf[1]), pos})

		if n, err := serverConn.Read(buf); n != 2 || err != nil || buf[0] != flagAck || buf[1] != incrementByte(pos) {
			t.Errorf("Unexpected data or error %d, %v", n, err)
		}
	}()

	err, ok := HandshakeClient(clientConn, false).(*HandshakeError)
	if !ok {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestHandshakeServerVsClient(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	clientConn, serverConn := NewFakeConnection(time.Second * 3)
	w := sync.WaitGroup{}
	w.Add(2)

	go func() {
		defer w.Done()
		err := HandshakeClient(clientConn, false)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}()

	go func() {
		defer w.Done()
		err := HandshakeServer(serverConn)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}()

	w.Wait()
}
