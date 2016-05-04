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
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
)

const (
	SYN = 0x16
	ACK = 0x06
	NAK = 0x15
)

// PurgeIncoming is used to clear a channel of bytes prior to handshaking
func PurgeIncoming(ctx context.Context, conn net.Conn) {
	buf := make([]byte, 255)

	// read until the incoming channel is empty
	log.Debug("purging incoming channel")
	conn.SetReadDeadline(time.Now().Add(time.Duration(10 * time.Millisecond)))
	for n, err := conn.Read(buf); n != 0 || err == nil; n, err = conn.Read(buf) {
		log.Debugf("discarding following %d bytes from input channel\n", n)
		log.Debugf("%+v\n", buf[0:n])
	}

	log.Debug("Incoming channel is purged of content")

	// disable the read timeout
	conn.SetReadDeadline(time.Time{})
}

func HandshakeClient(ctx context.Context, conn net.Conn) error {
	syn := make([]byte, 2)
	synack := make([]byte, 2)
	buf := make([]byte, 255)
	ack := make([]byte, 2)

	syn[0] = SYN
	synack[0] = ACK
	ack[0] = ACK

	// set the read deadline for timeout
	// this has no effect on windows as the deadline is set at port open time
	deadline, ok := ctx.Deadline()
	if ok {
		conn.SetReadDeadline(deadline)
	}

	rand.Read(syn[1:])

	log.Debug("client: writing syn")
	conn.Write(syn)
	log.Debug("client: reading synack")
	if n, err := io.ReadFull(conn, buf[:3]); n != 3 || err != nil {
		msg := fmt.Sprintf("HandshakeClient: failed to read expected SYN-ACK: n=%d, err=%s buf=[%#x]", n, err, buf[:n])
		if err != nil {
			log.Error(msg)
		} else {
			log.Debug(msg)
		}
		return errors.New(msg)
	}

	synack[1] = syn[1] + 1
	if bytes.Compare(synack, buf[:2]) != 0 {
		msg := fmt.Sprintf("HandshakeClient: did not receive synack: %#x != %#x", synack, buf[:2])
		log.Debugf(msg)
		conn.Write([]byte{NAK})
		return errors.New(msg)
	}

	log.Debugf("HandshakeClient: received synack: %#x == %#x\n", synack, buf[:2])
	log.Debug("client: writing ack")
	ack[1] = buf[2] + 1
	conn.Write(ack)

	// disable the read timeout
	conn.SetReadDeadline(time.Time{})

	return nil
}

func readMultiple(conn net.Conn, b []byte) (int, error) {
	if runtime.GOOS != "windows" {
		return conn.Read(b)
	}

	// we want a blocking read, but the behaviour described in the remarks section
	// here is a problem as we never get the syn in the same read as the rest.
	// https://msdn.microsoft.com/en-us/library/aa363190(v=VS.85).aspx

	// we know we're never reading single bytes for the handshake, so we hack this
	// multi read path into place
	// ideally we'd allow fine tweaking of the ReadIntervalTimeout so that we could
	// toggle between this behaviour and blocking read at a windows level
	n, err := conn.Read(b)
	if err == nil && n == 1 {
		cn, cerr := conn.Read(b[1:])
		return cn + 1, cerr
	}
	return n, err
}

func HandshakeServer(ctx context.Context, conn net.Conn) error {
	syn := make([]byte, 3)
	synack := make([]byte, 3)
	buf := make([]byte, 2)
	ack := make([]byte, 2)

	// set the read deadline for timeout
	// this has no effect on windows as the deadline is set at port open time
	deadline, ok := ctx.Deadline()
	if ok {
		conn.SetReadDeadline(deadline)
	}

	fmt.Println("server: reading syn")
	// loop here until we get a valid syn opening. syn is 3 bytes as that will eventually
	// syn us again if we're offset
	if n, err := readMultiple(conn, syn); n != 2 || err != nil || syn[0] != SYN {
		var msg string
		if err != nil {
			msg = fmt.Sprintf("server: failed to read expected SYN: n=%d, err=%s", n, err)
		} else if syn[0] != SYN {
			msg = fmt.Sprintf("server: did not receive syn (read %v bytes): %#x != %#x", n, SYN, syn[0])
			conn.Write([]byte{NAK})
		} else {
			msg = fmt.Sprintf("server: received syn but expected single sequence byte")
		}

		// to aid in debug we always dump the full handhsake
		log.Debug(msg)
		log.Debugf("server: read %v bytes: ", n)
		for i := 0; i < n; i++ {
			log.Debugf("%#x ", syn[i])
		}
		log.Debug("")

		return errors.New(msg)
	}
	log.Debugf("server: received syn: %#x\n", syn)

	log.Debug("server: writing synack")
	synack[0] = ACK
	synack[1] = syn[1] + 1
	rand.Read(synack[2:])

	conn.Write(synack)

	ack[0] = ACK
	ack[1] = synack[2] + 1
	log.Debug("server: reading ack")
	readMultiple(conn, buf)
	if bytes.Compare(ack, buf) != 0 {
		msg := fmt.Sprintf("server: did not receive ack: %#x != %#x", ack, buf)
		log.Debug(msg)
		conn.Write([]byte{NAK})
		return errors.New(msg)
	}

	log.Debugf("server: received ack: %#x == %#x\n", ack, buf)

	// disable the read timeout
	// this has no effect on windows as the deadline is set at port open time
	conn.SetReadDeadline(time.Time{})

	return nil
}
