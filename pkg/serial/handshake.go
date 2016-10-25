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
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/trace"
)

const (
	flagSyn      = 0x16
	flagAck      = 0x06
	flagDebugAck = 0x07
	flagNak      = 0x15
)

var connEstablised = []byte("established")

// PurgeIncoming is used to clear a channel of bytes prior to handshaking
func PurgeIncoming(conn net.Conn) {
	if tracing {
		defer trace.End(trace.Begin(""))
	}
	buf := make([]byte, 255)

	// read until the incoming channel is empty
	conn.SetReadDeadline(time.Now().Add(time.Duration(10 * time.Millisecond)))
	for n, err := conn.Read(buf); n != 0 || err == nil; n, err = conn.Read(buf) {
		log.Debugf("discarding following %d bytes from input channel\n", n)
		log.Debugf("%+v\n", buf[0:n])
	}

	// disable the read timeout
	conn.SetReadDeadline(time.Time{})
}

func HandshakeClient(ctx context.Context, conn net.Conn, debug bool) error {
	if tracing {
		defer trace.End(trace.Begin(""))
	}
	synAck := make([]byte, 3)
	buf := make([]byte, 255)
	syn := []byte{flagSyn, 0}
	rand.Read(syn[1:])

	// set the read deadline for timeout
	// this has no effect on windows as the deadline is set at port open time
	deadline, ok := ctx.Deadline()
	if ok {
		log.Debugf("Setting deadline to receive data from server: %s", deadline)
		conn.SetReadDeadline(deadline)
	}

	log.Debugf("HandshakeClient: Writing SYN: %#v", syn)
	conn.Write(syn)

	if n, err := io.ReadFull(conn, synAck); err != nil {
		log.Errorf("HandshakeClient: failed to read expected SYN-ACK: n=%d, err=%s buf=[%#x]",
			n, err, buf[:n])
		if err != nil {
			return err
		}
	}

	log.Debugf("Received SYNACK: %#v", synAck)
	if synAck[0] != flagAck || synAck[1] != syn[1]+1 {
		conn.Write([]byte{flagNak})
		return fmt.Errorf("Unexpected sequence received for SYNACK: %#v", synAck)
	}

	ack := []byte{flagAck, synAck[2] + 1}
	log.Debug("client: writing ack")
	if debug {
		ack[0] = flagDebugAck
	}
	conn.Write(ack)

	recvData := make([]byte, len(connEstablised))

	if _, err := io.ReadFull(conn, recvData); err != nil {
		return fmt.Errorf("Failed to read confirmation line: %s", err)
	}

	if bytes.Compare(connEstablised, recvData) != 0 {
		return fmt.Errorf("Received data should be %s, not %s",
			connEstablised, recvData)
	}

	conn.Write(connEstablised)

	// disable the read timeout
	conn.SetReadDeadline(time.Time{})

	if debug {
		// Verify packet length handling works.  We're going to send a known stream
		// of data to the container and it will echo it back.  Verify the sent and
		// received bufs are the same and we know the channel is lossless.

		log.Debugf("Checking for lossiness")
		txbuf := []byte("\x1b[32mhello world\x1b[39m!\n")
		rxbuf := make([]byte, len(txbuf))

		_, err := conn.Write(txbuf)
		if err != nil {
			log.Error(err)
			return err
		}

		var n int
		n, err = io.ReadFull(conn, rxbuf)
		if err != nil {
			return err
		}

		if n != len(rxbuf) {
			err = fmt.Errorf("packet size mismatch (expected %d, received %d)", len(rxbuf), n)
			return err
		}

		if bytes.Compare(rxbuf, txbuf) != 0 {
			conn.Write([]byte{flagNak})
			err = fmt.Errorf("client: lossiness check FAILED")
			return err
		}

		// Tell the server we're good.
		if _, err = conn.Write([]byte{flagAck}); err != nil {
			return err
		}

		log.Infof("client: lossiness check PASSED")

	}

	return nil
}

func HandshakeServer(ctx context.Context, conn net.Conn) error {
	if tracing {
		defer trace.End(trace.Begin(""))
	}
	synBuf := make([]byte, 2)

	// set the read deadline for timeout
	// this has no effect on windows as the deadline is set at port open time
	deadline, ok := ctx.Deadline()
	if ok {
		conn.SetReadDeadline(deadline)
	}

	log.Debug("server: reading syn")
	// syn is 2 bytes as that will eventually syn us again if we're offset
	n, err := io.ReadAtLeast(conn, synBuf, 2)
	if err != nil {
		return fmt.Errorf("server: failed to read expected SYN: n=%d, err=%s", n, err)
	}

	if synBuf[0] != flagSyn {
		conn.Write([]byte{flagNak})
		// to aid in debug we always dump the full handshake
		log.Debugf("server: read %v bytes: %#v", n, synBuf[:n])
		return fmt.Errorf("server: did not receive SYN (read %v bytes): %#x != %#x", n, flagSyn, synBuf[0])
	}

	log.Debugf("server: received syn: %#v", synBuf)

	log.Debug("server: writing synack")

	// last byte is random.
	synackBuf := []byte{flagAck, synBuf[1] + 1, 0}
	rand.Read(synackBuf[2:])

	conn.Write(synackBuf)

	ack := []byte{flagAck, synackBuf[2] + 1}
	log.Debug("server: reading ack")

	buf := make([]byte, 2)
	io.ReadFull(conn, buf)

	if (buf[0] != flagAck && buf[0] != flagDebugAck) || ack[1] != buf[1] {
		conn.Write([]byte{flagNak})
		return fmt.Errorf("server: did not receive ack: %#x != %#x", ack, buf)
	}

	log.Debugf("server: received ack: %#x == %#x", ack, buf)

	recvData := make([]byte, len(connEstablised))

	conn.Write(connEstablised)

	if _, err := io.ReadFull(conn, recvData); err != nil {
		return fmt.Errorf("Failed to read confirmation line: %s", err)
	}

	if bytes.Compare(connEstablised, recvData) != 0 {
		return fmt.Errorf("Received data should be %s, not %s",
			connEstablised, recvData)
	}

	// disable the read timeout
	// this has no effect on windows as the deadline is set at port open time
	conn.SetReadDeadline(time.Time{})

	if buf[0] == flagDebugAck {
		// Check for lossiness
		rxbuf := make([]byte, 23)
		log.Debugf("Checking for lossiness")

		n, err := io.ReadFull(conn, rxbuf)
		if err != nil {
			return err
		}

		if n != len(rxbuf) {
			return fmt.Errorf("packet size mismatch (expected %d, received %d)", len(rxbuf), n)
		}

		// echo the data back
		_, err = conn.Write(rxbuf)
		if err != nil {
			return err
		}

		// wait for the ack
		if _, err = conn.Read(ack[:1]); err != nil {
			return err
		}

		if ack[0] != flagAck {
			return fmt.Errorf("server: lossiness check FAILED")
		}
		log.Infof("server: lossiness check PASSED")
	}

	return nil
}
