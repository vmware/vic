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
	flagSyn      byte = 0x16
	flagAck      byte = 0x06
	flagDebugAck byte = 0x07
	flagNak      byte = 0x15
)

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

func nextPos(syncPos byte) byte { return (syncPos + 1) | 0x80 }

func HandshakeClient(ctx context.Context, conn net.Conn, debug bool) error {
	if tracing {
		defer trace.End(trace.Begin(""))
	}

	buf1byte := make([]byte, 1)
	rand.Read(buf1byte)
	pos := nextPos(buf1byte[0])

	// set the read deadline for timeout
	// this has no effect on windows as the deadline is set at port open time
	deadline, ok := ctx.Deadline()
	if ok {
		conn.SetReadDeadline(deadline)
	}

	log.Debugf("HandshakeClient: Establishing connection...")
	for {
		log.Debug("HandshakeClient: Sending syn.")
		conn.Write([]byte{flagSyn, pos})
		pos = nextPos(pos)

		if _, err := conn.Read(buf1byte); err != nil {
			return err
		}

		if buf1byte[0] != flagAck {
			if buf1byte[0] == flagNak {
				log.Debugf("HandshakeClient: Server didn't accept sync. Trying one more time.")
			} else {
				log.Debugf("HandshakeClient: Unexpected byte for SynAck: %x", buf1byte[0])
			}
			continue
		}

		// read response sync position.
		if _, err := conn.Read(buf1byte); err != nil {
			return err
		}

		if buf1byte[0] != pos {
			log.Debugf("HandshakeClient: Unexpected byte pos for SynAck: %x", buf1byte[0])
			continue
		}

		if _, err := conn.Read(buf1byte); err != nil {
			return err
		}

		log.Debug("HandshakeClient: Sending ack.")
		if !debug {
			conn.Write([]byte{flagAck, nextPos(buf1byte[0])})
		} else {
			conn.Write([]byte{flagDebugAck, nextPos(buf1byte[0])})
			// Verify packet length handling works.  We're going to send a known stream
			// of data to the container and it will echo it back.  Verify the sent and
			// received bufs are the same and we know the channel is lossless.

			log.Debugf("HandshakeClient: Checking for lossiness")
			txbuf := []byte("\x1b[32mhello world\x1b[39m!\n")
			rxbuf := make([]byte, len(txbuf))

			_, err := conn.Write(txbuf)
			if err != nil {
				return err
			}

			var n int
			log.Debugf("HandshakeClient: Reading response")
			n, err = io.ReadFull(conn, rxbuf)
			if err != nil {
				log.Error(err)
				return err
			}

			if n != len(rxbuf) {
				return fmt.Errorf("packet size mismatch (expected %d, received %d)", len(rxbuf), n)
			}

			if bytes.Compare(rxbuf, txbuf) != 0 {
				return fmt.Errorf("HandshakeClient: lossiness check FAILED")
			}

			// Tell the server we're good.
			if _, err = conn.Write([]byte{flagAck}); err != nil {
				return err
			}

			log.Infof("HandshakeClient: lossiness check PASSED")
		}
		break
	}
	log.Debug("HandshakeClient: Connection established.")
	conn.SetReadDeadline(time.Time{})

	return nil
}

func HandshakeServer(ctx context.Context, conn net.Conn) error {
	if tracing {
		defer trace.End(trace.Begin(""))
	}

	deadline, ok := ctx.Deadline()
	if ok {
		conn.SetReadDeadline(deadline)
	}

	buf1byte := make([]byte, 1)
	syncBuf := make([]byte, 4096)

	rand.Read(buf1byte)
	pos := nextPos(buf1byte[0])

	log.Debugf("HandshakeServer: Establishing connection...")
	for {
		// Sync packet is 2 bytes, however if we read more than 2
		// it means buffer is not empty and data is not trusted for this sync.

		n, err := io.ReadAtLeast(conn, syncBuf, 2)
		if err != nil {
			return err
		}

		if n > 2 {
			log.Debugf("HandshakeServer: More than 2 bytes were received for Syn: %x", syncBuf[:n])
			conn.Write([]byte{flagNak})
			continue
		}

		if syncBuf[0] != flagSyn {
			log.Debugf("HandshakeServer: Unexpected byte for sync: %x", syncBuf[0])
			conn.Write([]byte{flagNak})
			continue
		}
		log.Debugf("HandshakeServer: Received Syn. Writing SynAck.")

		// syncBuf[1] contains position token that needs to be incremented
		// by one to send it back.
		conn.Write([]byte{flagAck, nextPos(syncBuf[1]), pos})
		pos = nextPos(pos)

		if _, err := conn.Read(buf1byte); err != nil {
			return err
		}

		ackType := buf1byte[0]
		if ackType != flagAck && ackType != flagDebugAck {
			conn.Write([]byte{flagNak})
			continue
		}

		if _, err := conn.Read(buf1byte); err != nil {
			return err
		}

		if buf1byte[0] != pos {
			log.Debug("HandshakeServer: Unexpected position %x, expected: %x", buf1byte[0], pos)
			conn.Write([]byte{flagNak})
			continue
		}

		if ackType == flagDebugAck {
			log.Debugf("HandshakeServer: Debug ACK received")
			rxbuf := make([]byte, 23)
			log.Debugf("HandshakeServer: Checking for lossiness")

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
			if _, err = conn.Read(buf1byte); err != nil {
				return err
			}

			if buf1byte[0] != flagAck {
				return fmt.Errorf("HandshakeServer: lossiness check FAILED")
			}
			log.Infof("HandshakeServer: lossiness check PASSED")
		}
		break
	}
	log.Debug("HandshakeServer: Connection established.")
	// disable the read timeout
	// this has no effect on windows as the deadline is set at port open time
	conn.SetReadDeadline(time.Time{})

	return nil
}
