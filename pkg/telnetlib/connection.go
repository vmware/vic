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
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type connOpts struct {
	conn        net.Conn
	fsm         *telnetFSM
	cmdHandler  CmdHandlerFunc
	dataHandler DataHandlerFunc
	serverOpts  map[byte]bool
	clientOpts  map[byte]bool
	optCallback func(byte, byte)
}

// TelnetConn is the struct representing the telnet connection
type TelnetConn struct {
	// conn is the underlying connection
	conn net.Conn
	// the connection write channel. Everything required to be written to the connection goes to this channel
	writeCh           chan []byte
	srvrOptsLock      sync.Mutex
	unackedServerOpts map[byte]bool
	clientOptsLock    sync.Mutex
	unackedClientOpts map[byte]bool
	//server            *TelnetServer
	serverOpts map[byte]bool
	clientOpts map[byte]bool

	// dataRW is the data buffer. It is written to by the FSM and read from by the data handler
	dataRW        io.ReadWriter
	cmdBuffer     bytes.Buffer
	fsm           *telnetFSM
	handlerWriter io.WriteCloser
	// cmdHandler is the command handler for the telnet server
	// it is a callback function when receiving commands issued by the telnet client
	// it is wrapped by the cmdHandlerWrapper
	cmdHandler CmdHandlerFunc

	// dataHandler is the data handler of the telnet server
	// it is a call back function when receiving data from the telnet client
	// it is wrapped by the dataHandlerWrapper
	dataHandler DataHandlerFunc

	// used in the dataHandlerWrapper to notify that the telnet connection is closed
	dataHandlerCloseCh chan chan struct{}
	// used in the dataHandlerWrapper to notify that data has been writeen to the dataRW buffer
	dataWrittenCh chan bool

	// callBack function for receiving an option during negotiation
	optionCallback func(byte, byte)

	// connWriteDoneCh closes the write loop when the telnet connection is closed
	connWriteDoneCh chan chan struct{}
	// negotiationDone notifies that telnet negotiation has been done
	negotiationDone chan struct{}

	closedMutex sync.Mutex
	closed      bool
}

// Safely read/write concurrently to the data Buffer
// databuffer is written to by the FSM and it is read from by the dataHandler
type dataReadWriter struct {
	dataBuffer bytes.Buffer
	dataMux    *sync.Mutex
}

func (drw *dataReadWriter) Read(p []byte) (int, error) {
	drw.dataMux.Lock()
	defer drw.dataMux.Unlock()
	return drw.dataBuffer.Read(p)
}

func (drw *dataReadWriter) Write(p []byte) (int, error) {
	drw.dataMux.Lock()
	defer drw.dataMux.Unlock()
	return drw.dataBuffer.Write(p)
}

type connectionWriter struct {
	ch chan []byte
}

func (cw *connectionWriter) Write(b []byte) (int, error) {
	if cw.ch != nil {
		cw.ch <- b
	}
	return len(b), nil
}

func (cw *connectionWriter) Close() error {
	close(cw.ch)
	cw.ch = nil
	return nil
}

func newTelnetConn(opts connOpts) *TelnetConn {
	tc := &TelnetConn{
		conn:               opts.conn,
		writeCh:            make(chan []byte),
		unackedServerOpts:  make(map[byte]bool),
		unackedClientOpts:  make(map[byte]bool),
		cmdHandler:         opts.cmdHandler,
		dataHandler:        opts.dataHandler,
		dataHandlerCloseCh: make(chan chan struct{}),
		dataWrittenCh:      make(chan bool),
		serverOpts:         opts.serverOpts,
		clientOpts:         opts.clientOpts,
		optionCallback:     opts.optCallback,
		connWriteDoneCh:    make(chan chan struct{}),
		negotiationDone:    make(chan struct{}),
		closed:             false,
	}
	if tc.optionCallback == nil {
		tc.optionCallback = tc.handleOptionCommand
	}
	tc.handlerWriter = &connectionWriter{
		ch: tc.writeCh,
	}
	tc.dataRW = &dataReadWriter{
		dataMux: new(sync.Mutex),
	}
	fsm := opts.fsm
	fsm.tc = tc
	tc.fsm = fsm
	for k := range tc.serverOpts {
		tc.unackedServerOpts[k] = true
	}
	for k := range tc.clientOpts {
		tc.unackedClientOpts[k] = true
	}
	return tc
}

func (c *TelnetConn) writeLoop() {
	log.Debugf("entered write loop")
	for {
		select {
		case writeBytes := <-c.writeCh:
			c.conn.Write(writeBytes)
		case ch := <-c.connWriteDoneCh:
			ch <- struct{}{}
			return
		}
	}
}

func (c *TelnetConn) startNegotiation() {
	c.srvrOptsLock.Lock()
	for k := range c.serverOpts {
		log.Infof("sending WILL %d", k)
		c.unackedServerOpts[k] = true
		c.sendCmd(Will, k)
	}
	c.srvrOptsLock.Unlock()
	c.clientOptsLock.Lock()
	for k := range c.clientOpts {
		log.Infof("sending DO %d", k)
		c.unackedClientOpts[k] = true
		c.sendCmd(Do, k)
	}
	c.clientOptsLock.Unlock()
	select {
	case <-c.negotiationDone:
		log.Infof("Negotiation finished")
		return
	case <-time.After(10 * time.Second):
		log.Infof("Negotiation failed. Exiting")
		c.Close()
		c.closed = true
		return
	}
}

// Close closes the telnet connection
func (c *TelnetConn) Close() {
	log.Infof("Closing the connection")
	c.conn.Close()
	c.closeConnLoopWrite()
	c.closeDatahandler()
	c.handlerWriter.Close()
	log.Infof("telnet connection closed")
	c.closedMutex.Lock()
	defer c.closedMutex.Unlock()
	c.closed = true
}

func (c *TelnetConn) closeConnLoopWrite() {
	connLoopWriteCh := make(chan struct{})
	c.connWriteDoneCh <- connLoopWriteCh
	<-connLoopWriteCh
	log.Infof("connection loop write-side closed")
}

func (c *TelnetConn) closeDatahandler() {
	dataCh := make(chan struct{})
	c.dataHandlerCloseCh <- dataCh
	<-dataCh
}

func (c *TelnetConn) sendCmd(cmd byte, opt byte) {
	b := []byte{Iac, cmd, opt}
	log.Infof("Sending command: %v %v", cmd, opt)
	c.writeCh <- b
}

func (c *TelnetConn) handleOptionCommand(cmd byte, opt byte) {
	c.clientOptsLock.Lock()
	defer c.clientOptsLock.Unlock()
	c.srvrOptsLock.Lock()
	defer c.srvrOptsLock.Unlock()
	if cmd == Will || cmd == Wont {
		if _, ok := c.clientOpts[opt]; !ok {
			c.sendCmd(Dont, opt)
			return
		}

		if _, ok := c.unackedClientOpts[opt]; ok {
			delete(c.unackedClientOpts, opt)
			if len(c.unackedClientOpts) == 0 && len(c.unackedServerOpts) == 0 {
				close(c.negotiationDone)
			}
		} else {
			c.sendCmd(Do, opt)
		}
	}

	if cmd == Do || cmd == Dont {
		if _, ok := c.serverOpts[opt]; !ok {
			c.sendCmd(Wont, opt)
			return
		}
		if _, ok := c.unackedServerOpts[opt]; ok {
			log.Infof("removing from the unack list")
			delete(c.unackedServerOpts, opt)
			if len(c.unackedClientOpts) == 0 && len(c.unackedServerOpts) == 0 {
				close(c.negotiationDone)
			}
		} else {
			log.Infof("Sending WILL command")
			c.sendCmd(Will, opt)
		}
	}
}

func (c *TelnetConn) dataHandlerWrapper(w io.Writer, r io.Reader) {
	defer func() {
		log.Infof("data handler closed")
	}()
	for {
		select {
		case ch := <-c.dataHandlerCloseCh:
			ch <- struct{}{}
			return
		case <-c.dataWrittenCh:
			if b, err := ioutil.ReadAll(r); err == nil {
				c.dataHandler(w, b, c)
			}
		}
	}
}

func (c *TelnetConn) cmdHandlerWrapper(w io.Writer, r io.Reader) {
	if cmd, err := ioutil.ReadAll(r); err == nil {
		c.cmdHandler(w, cmd, c)
	}
}

// IsClosed returns true if the connection is already closed
func (c *TelnetConn) IsClosed() bool {
	c.closedMutex.Lock()
	defer c.closedMutex.Unlock()
	return c.closed
}
