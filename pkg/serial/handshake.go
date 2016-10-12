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

	"github.com/vmware/vic/pkg/trace"
)

const (
	flagSyn      = 0x16
	flagAck      = 0x06
	flagDebugAck = 0x07
	flagNak      = 0x15
)

// PurgeIncoming is used to clear a channel of bytes prior to handshaking
func PurgeIncoming(conn net.Conn) {
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
	syn := make([]byte, 2)
	synack := make([]byte, 2)
	ack := make([]byte, 2)

	buf := make([]byte, 255)

	syn[0] = flagSyn
	synack[0] = flagAck
	ack[0] = flagAck

	// set the read deadline for timeout
	// this has no effect on windows as the deadline is set at port open time
	deadline, ok := ctx.Deadline()
	if ok {
		conn.SetReadDeadline(deadline)
	}

	rand.Read(syn[1:])

	conn.Write(syn)

	if n, err := io.ReadFull(conn, buf[:3]); n != 3 || err != nil {

		if n == 0 && err != nil {
			return err
		}

		msg := fmt.Sprintf("HandshakeClient: failed to read expected SYN-ACK: n=%d, err=%s buf=[%#x]", n, err, buf[:n])
		if err != nil {
			log.Error(msg)
		} else {
			log.Debug(msg)
		}
		return err
	}

	synack[1] = syn[1] + 1
	if bytes.Compare(synack, buf[:2]) != 0 {
		msg := fmt.Sprintf("HandshakeClient: did not receive synack: %#x != %#x", synack, buf[:2])
		log.Debugf(msg)
		conn.Write([]byte{flagNak})
		return errors.New(msg)
	}

	log.Infof("HandshakeClient: received synack: %#x == %#x\n", synack, buf[:2])
	log.Debug("client: writing ack")
	ack[1] = buf[2] + 1
	if debug {
		ack[0] = flagDebugAck
	}
	conn.Write(ack)

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
			log.Error(err)
			return err
		}

		if n != len(rxbuf) {
			err = fmt.Errorf("packet size mismatch (expected %d, received %d)", len(rxbuf), n)
			log.Error(err)
			return err
		}

		if bytes.Compare(rxbuf, txbuf) != 0 {
			conn.Write([]byte{flagNak})
			err = fmt.Errorf("client: lossiness check FAILED")
			log.Error(err)
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
	defer trace.End(trace.Begin(""))

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

	log.Debug("server: reading syn")
	// syn is 3 bytes as that will eventually syn us again if we're offset
	if n, err := readMultiple(conn, syn); n != 2 || err != nil || syn[0] != flagSyn {
		var msg string
		if err != nil {
			msg = fmt.Sprintf("server: failed to read expected SYN: n=%d, err=%s", n, err)
		} else if syn[0] != flagSyn {
			msg = fmt.Sprintf("server: did not receive SYN (read %v bytes): %#x != %#x", n, flagSyn, syn[0])
			conn.Write([]byte{flagNak})
		} else {
			msg = fmt.Sprintf("server: received SYN (read %d) bytes", n)
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
	synack[0] = flagAck
	synack[1] = syn[1] + 1
	rand.Read(synack[2:])

	conn.Write(synack)

	ack[0] = flagAck
	ack[1] = synack[2] + 1
	log.Debug("server: reading ack")
	readMultiple(conn, buf)
	if (buf[0] != flagAck && buf[0] != flagDebugAck) || bytes.Compare(ack[1:], buf[1:]) != 0 {
		msg := fmt.Sprintf("server: did not receive ack: %#x != %#x", ack, buf)
		log.Debug(msg)
		conn.Write([]byte{flagNak})
		return errors.New(msg)
	}
	log.Infof("server: received ack: %#x == %#x\n", ack, buf)

	// disable the read timeout
	// this has no effect on windows as the deadline is set at port open time
	conn.SetReadDeadline(time.Time{})

	if buf[0] == flagDebugAck {
		// Check for lossiness
		rxbuf := make([]byte, 23)
		log.Debugf("Checking for lossiness")
		n, err := io.ReadFull(conn, rxbuf)
		if err != nil {
			log.Error(err)
			return err
		}

		if n != len(rxbuf) {
			err = fmt.Errorf("packet size mismatch (expected %d, received %d)", len(rxbuf), n)
			log.Error(err)
			return err
		}

		// echo the data back
		_, err = conn.Write(rxbuf)
		if err != nil {
			log.Error(err)
			return err
		}

		// wait for the ack
		if _, err = conn.Read(ack[:1]); err != nil {
			log.Error(err)
			return err
		}

		if ack[0] != flagAck {
			err = fmt.Errorf("server: lossiness check FAILED")
			log.Error(err)
			return err
		}
		log.Infof("server: lossiness check PASSED")
	}

	return nil
}
