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
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/vmware/vic/pkg/nfs/xdr"
)

var xid uint32

func init() {
	// seed the XID (which is set by the client)
	xid = rand.New(rand.NewSource(time.Now().UnixNano())).Uint32()
}

type Client struct {
	transport
}

func DialTCP(network string, ldr *net.TCPAddr, addr string) (*Client, error) {
	a, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP(a.Network(), ldr, a)
	if err != nil {
		return nil, err
	}

	t := &tcpTransport{
		Reader:      bufio.NewReader(conn),
		WriteCloser: conn,
	}
	return &Client{t}, nil
}

func (c *Client) Call(call interface{}) ([]byte, error) {
	msg := &message{
		Xid:     atomic.AddUint32(&xid, 1),
		Msgtype: 0,
		Body:    call,
	}

	w := new(bytes.Buffer)
	if err := xdr.Write(w, msg); err != nil {
		return nil, err
	}

	if err := c.send(w.Bytes()); err != nil {
		return nil, err
	}

	buf, err := c.recv()
	if err != nil {
		return nil, err
	}

	xid, buf := xdr.Uint32(buf)
	if xid != msg.Xid {
		return nil, fmt.Errorf("xid did not match, expected: %x, received: %x", msg.Xid, xid)
	}

	mtype, buf := xdr.Uint32(buf)
	if mtype != 1 {
		return nil, fmt.Errorf("message as not a reply: %d", mtype)
	}

	reply_stat, buf := xdr.Uint32(buf)
	switch reply_stat {
	case MSG_ACCEPTED:
		_, buf = xdr.Uint32(buf)
		opaque_len, buf := xdr.Uint32(buf)
		_ = buf[0:int(opaque_len)]
		buf = buf[opaque_len:]
		accept_stat, buf := xdr.Uint32(buf)

		switch accept_stat {
		case SUCCESS:
			return buf, nil
		case PROG_UNAVAIL:
			return nil, fmt.Errorf("PROG_UNAVAIL")
		case PROG_MISMATCH:
			// TODO(dfc) decode mismatch_info
			return nil, fmt.Errorf("rpc: PROG_MISMATCH")
		default:
			return nil, fmt.Errorf("rpc: %d", accept_stat)
		}

	case MSG_DENIED:
		rejected_stat, _ := xdr.Uint32(buf)
		switch rejected_stat {
		case RPC_MISMATCH:

		default:
			return nil, fmt.Errorf("rejected_stat was not valid: %d", rejected_stat)
		}

	default:
		return nil, fmt.Errorf("reply_stat was not valid: %d", reply_stat)
	}

	panic("unreachable")
}
