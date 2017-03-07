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

package telnetlib

import (
	"fmt"
	"io"
	"net"

	log "github.com/Sirupsen/logrus"
)

// DataHandlerFunc is the callback function in the event of receiving data from the telnet client
type DataHandlerFunc func(w io.Writer, data []byte, tc *TelnetConn)

// CmdHandlerFunc is the callback function in the event of receiving a command from the telnet client
type CmdHandlerFunc func(w io.Writer, cmd []byte, tc *TelnetConn)

var defaultDataHandlerFunc = func(w io.Writer, data []byte, tc *TelnetConn) {}

var defaultCmdHandlerFunc = func(w io.Writer, cmd []byte, tc *TelnetConn) {}

// TelnetOpts is the telnet server constructor options
type TelnetOpts struct {
	Addr        string
	ServerOpts  []byte
	ClientOpts  []byte
	DataHandler DataHandlerFunc
	CmdHandler  CmdHandlerFunc
}

// TelnetServer is the struct representing the telnet server
type TelnetServer struct {
	ServerOptions map[byte]bool
	ClientOptions map[byte]bool
	DataHandler   DataHandlerFunc
	CmdHandler    CmdHandlerFunc
	ln            net.Listener
}

// NewTelnetServer is the constructor of the telnet server
func NewTelnetServer(opts TelnetOpts) *TelnetServer {
	ts := new(TelnetServer)
	ts.ClientOptions = make(map[byte]bool)
	ts.ServerOptions = make(map[byte]bool)
	for _, v := range opts.ServerOpts {
		ts.ServerOptions[v] = true
	}
	for _, v := range opts.ClientOpts {
		ts.ClientOptions[v] = true
	}
	ts.DataHandler = opts.DataHandler
	if ts.DataHandler == nil {
		ts.DataHandler = defaultDataHandlerFunc
	}

	ts.CmdHandler = opts.CmdHandler
	if ts.CmdHandler == nil {
		ts.CmdHandler = defaultCmdHandlerFunc
	}
	ln, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		panic(fmt.Sprintf("cannot start telnet server: %v", err))
	}
	ts.ln = ln
	return ts
}

// Accept accepts a connection and returns the Telnet connection
func (ts *TelnetServer) Accept() (*TelnetConn, error) {
	conn, _ := ts.ln.Accept()
	log.Info("connection received")
	opts := connOpts{
		conn:        conn,
		cmdHandler:  ts.CmdHandler,
		dataHandler: ts.DataHandler,
		serverOpts:  ts.ServerOptions,
		clientOpts:  ts.ClientOptions,
		fsm:         newTelnetFSM(),
	}
	tc := newTelnetConn(opts)
	go tc.writeLoop()
	go tc.dataHandlerWrapper(tc.handlerWriter, tc.dataRW)
	go tc.fsm.start()
	go tc.startNegotiation()
	return tc, nil
}
