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
	"runtime"
	"time"
)

const (
	SYN = 0x16
	ACK = 0x06
	NAK = 0x15
)

func HandshakeClient(conn net.Conn, interval time.Duration) error {
	syn := make([]byte, 2)
	synack := make([]byte, 2)
	buf := make([]byte, 255)
	ack := make([]byte, 2)

	syn[0] = SYN
	synack[0] = ACK
	ack[0] = ACK

	for {
		// read until the incoming channel is empty
		// TODO: raw socket doesn't yet have read timeout
		fmt.Printf("purging incoming channel\n")
		conn.SetReadDeadline(time.Now().Add(time.Duration(10 * time.Millisecond)))
		for n, err := conn.Read(buf); n != 0 || err == nil; n, err = conn.Read(buf) {
			fmt.Printf("discarding following %d bytes from input channel\n", n)
			fmt.Printf("%+v\n", buf[0:n])
		}

		fmt.Println("Incoming channel is purged of content")

		// set the read deadline for timeout
		conn.SetReadDeadline(time.Now().Add(interval))

		rand.Read(syn[1:])

		fmt.Println("writing syn")
		conn.Write(syn)
		fmt.Println("reading synack")
		if n, err := io.ReadFull(conn, buf[:3]); n != 3 || err != nil {
			fmt.Printf("Failed to read expected SYN-ACK: n=%d, err=%s buf=[%#x]\n", n, err, buf[:n])
			if err != nil {
				return err
			}

			continue
		}

		synack[1] = syn[1] + 1
		if bytes.Compare(synack, buf[:2]) != 0 {
			fmt.Printf("Did not receive synack: %#x != %#x\n", synack, buf[:2])
			conn.Write([]byte{NAK})
			//time.Sleep(time.Duration(10 * time.Second))
			continue
		} else {
			fmt.Printf("Received synack: %#x == %#x\n", synack, buf[:2])
		}

		fmt.Println("writing ack")
		ack[1] = buf[2] + 1
		conn.Write(ack)
		break
	}

	// disable the read timeout
	conn.SetReadDeadline(time.Time{})

	return nil
}

func readMultiple(conn net.Conn, b []byte) (int, error) {
	if runtime.GOOS != "windows" {
		return conn.Read(b)
	} else {
		// we want a blocking read, but the behaviour described in the remarks section
		// here is a problem as we never get the syn in the same read as the rest.
		// https://msdn.microsoft.com/en-us/library/aa363190(v=VS.85).aspx

		// we know we're never reading single bytes for the handshake, so we hack this
		// multi read path into place
		// ideally we'd allow fine tweaking of the ReadIntervalTimeout so that we could
		// toggle between this behaviour and blocking read at a windows level
		n, err := conn.Read(b)
		if err == nil && n == 1 {
			n, err := conn.Read(b[1:])
			return n + 1, err
		}
		return n, err
	}
}

func HandshakeServer(conn net.Conn, interval time.Duration) {
	syn := make([]byte, 3)
	synack := make([]byte, 3)
	buf := make([]byte, 2)
	ack := make([]byte, 2)

	for {
		// set the read deadline for timeout
		// this has no effect on windows as the deadline is set at port open time
		conn.SetReadDeadline(time.Now().Add(interval))

		fmt.Println("reading syn")
		// loop here until we get a valid syn opening. syn is 3 bytes as that will eventually
		// syn us again if we're offset
		if n, err := readMultiple(conn, syn); n != 2 || err != nil || syn[0] != SYN {
			if err != nil {
				fmt.Printf("Failed to read expected SYN: n=%d, err=%s\n", n, err)
			} else if syn[0] != SYN {
				fmt.Printf("Did not receive syn (read %v bytes): %#x != %#x\n", n, SYN, syn[0])
				conn.Write([]byte{NAK})
			} else {
				fmt.Printf("Received syn but expected single sequence byte\n")
			}

			// to aid in debug we always dump the full handhsake
			fmt.Printf("read %v bytes: ", n)
			for i := 0; i < n; i++ {
				fmt.Printf("%#x ", syn[i])
			}
			fmt.Println()

			continue
		}
		fmt.Printf("Received syn: %#x\n", syn)

		fmt.Println("writing synack")
		synack[0] = ACK
		synack[1] = syn[1] + 1
		rand.Read(synack[2:])

		conn.Write(synack)

		ack[0] = ACK
		ack[1] = synack[2] + 1
		fmt.Println("reading ack")
		readMultiple(conn, buf)
		if bytes.Compare(ack, buf) != 0 {
			fmt.Printf("Did not receive ack: %#x != %#x\n", ack, buf)
			conn.Write([]byte{NAK})
			//time.Sleep(time.Duration(10 * time.Second))
			continue
		} else {
			fmt.Printf("Received ack: %#x == %#x\n", ack, buf)
		}

		break
	}

	// disable the read timeout
	// this has no effect on windows as the deadline is set at port open time
	conn.SetReadDeadline(time.Time{})
}
