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
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/docker/docker/pkg/term"
)

type NamedChannel interface {
	io.ReadWriteCloser
	Name() string
	Fd() uintptr
}

type RawConn struct {
	channel    NamedChannel
	localAddr  net.Addr
	remoteAddr net.Addr
	state      *term.State
	err        chan error
	mutex      sync.Mutex
	closed     bool
}

func NewTypedConn(channel NamedChannel, net string) (*RawConn, error) {
	conn := &RawConn{
		channel:    channel,
		localAddr:  *NewRawAddr(net, channel.Name()),
		remoteAddr: *NewRawAddr("void", "void"),
		err:        make(chan error, 1),
		closed:     false,
	}

	// set the provided FD to raw if it's a termial
	// 0 is the uninitialized value for Fd
	if channel.Fd() != 0 && term.IsTerminal(channel.Fd()) {
		log.Println("setting terminal into raw mode")
		terminal, err := term.SetRawTerminal(channel.Fd())
		if err != nil {
			return nil, err
		}
		conn.state = terminal
	}

	return conn, nil
}

func NewFileConn(file *os.File) (*RawConn, error) {
	return NewTypedConn(file, "file")
}

func NewRawConn(fd uintptr, name string, net string) (*RawConn, error) {
	return NewTypedConn(os.NewFile(fd, name), net)
}

func (conn *RawConn) Read(b []byte) (n int, err error) {
	// TODO: this is horrific from a performance perspective - really need a better
	// way to interrupt that file.Read call
	bytes := make(chan int, 1)

	go func() {
		n, err = conn.channel.Read(b)

		// if we've got any bytes we need to pass them back so we cannot return
		// the error via conn.err
		bytes <- n
	}()

	select {
	case n = <-bytes:
		if err != nil {
			log.Printf("Returning error and bytes from read: %d, %s\n", n, err)
		}
		return n, err
	case e := <-conn.err:
		log.Printf("Returning error from read: %s\n", e)
		// only one close will send an error and we have that, so this won't block
		// we do need to interrupt all reads
		conn.err <- e
		return n, e
	}
}

func (conn *RawConn) Write(b []byte) (n int, err error) {
	n, err = conn.channel.Write(b)
	return
}

func (conn *RawConn) Close() error {
	var closed bool

	conn.mutex.Lock()
	closed = conn.closed
	conn.closed = true
	conn.mutex.Unlock()

	if closed {
		log.Printf("Close called again on RawConn\n")
		return nil
	}

	// process the close
	err := conn.channel.Close()

	buf := make([]byte, 4096)
	bytes := runtime.Stack(buf, false)
	log.Printf("Close called on RawConn:\n%s\n", string(buf[:bytes]))

	log.Println("Pushing EOF to any blocked readers on the raw connection")
	conn.err <- io.EOF

	return err
}

func (conn *RawConn) LocalAddr() net.Addr {
	return conn.localAddr
}

func (conn *RawConn) RemoteAddr() net.Addr {
	return conn.remoteAddr
}

func (conn *RawConn) SetDeadline(t time.Time) error {
	// https://golang.org/src/net/fd_poll_runtime.go#L133
	// consider implementing this by making RawConn a netFD
	// if we can find a way around the lack of export
	return nil
}

func (conn *RawConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (conn *RawConn) SetWriteDeadline(t time.Time) error {
	return nil
}
