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

type TelnetConn struct {
	conn              net.Conn
	readCh            chan []byte
	writeCh           chan []byte
	unackedServerOpts map[byte]bool
	unackedClientOpts map[byte]bool
	//server            *TelnetServer
	serverOpts         map[byte]bool
	clientOpts         map[byte]bool
	dataRW             io.ReadWriter
	cmdBuffer          bytes.Buffer
	fsm                *telnetFSM
	fsmInputCh         chan byte
	handlerWriter      io.Writer
	cmdHandler         CmdHandlerFunc
	dataHandler        DataHandlerFunc
	dataHandlerCloseCh chan chan struct{}
	dataWrittenCh      chan bool
	optionCallback     func(byte, byte)
	connReadDoneCh     chan chan struct{}
	connWriteDoneCh    chan chan struct{}
	negotiationDone    chan struct{}
	closedMutex        sync.Mutex
	closed             bool
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
	cw.ch <- b
	return len(b), nil
}

func newTelnetConn(opts connOpts) *TelnetConn {
	tc := &TelnetConn{
		conn:               opts.conn,
		readCh:             make(chan []byte),
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
		fsmInputCh:         make(chan byte),
		connReadDoneCh:     make(chan chan struct{}),
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

func (c *TelnetConn) connectionLoop() {
	log.Debugf("Entered connectionLoop")
	// this is the reading thread
	go func() {
		for {
			select {
			case readBytes := <-c.readCh:
				for _, ch := range readBytes {
					c.fsmInputCh <- ch
				}

			case ch := <-c.connReadDoneCh:
				ch <- struct{}{}
				return
			}
		}
	}()
	// this is the writing thread
	go func() {
		for {
			select {
			case writeBytes := <-c.writeCh:
				c.conn.Write(writeBytes)
			case ch := <-c.connWriteDoneCh:
				ch <- struct{}{}
				return
			}
		}
	}()
}

// reads from the connection and dumps into the connection read channel
func (c *TelnetConn) readLoop() {
	defer func() {
		log.Debugf("read loop closed")
	}()
	for {
		buf := make([]byte, 4096)
		n, err := c.conn.Read(buf)
		if n > 0 {
			log.Debug("read %d bytes from the TCP Connection %v", n, buf[:n])
			c.readCh <- buf[:n]
		}
		if err != nil {
			log.Debugf("connection read: %v", err)
			c.Close()
			break
		}
	}
}

func (c *TelnetConn) startNegotiation() {
	for k := range c.serverOpts {
		log.Infof("sending WILL %d", k)
		c.unackedServerOpts[k] = true
		c.sendCmd(WILL, k)
	}
	for k := range c.clientOpts {
		log.Infof("sending DO %d", k)
		c.unackedClientOpts[k] = true
		c.sendCmd(DO, k)
	}
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
	c.closeConnLoopRead()
	c.closeConnLoopWrite()
	c.closeFSM()
	c.closeDatahandler()
	log.Infof("telnet connection closed")
	c.closedMutex.Lock()
	defer c.closedMutex.Unlock()
	c.closed = true
}

func (c *TelnetConn) closeConnLoopRead() {
	connLoopReadCh := make(chan struct{})
	c.connReadDoneCh <- connLoopReadCh
	<-connLoopReadCh
	log.Infof("connection loop read-side closed")
}

func (c *TelnetConn) closeConnLoopWrite() {
	connLoopWriteCh := make(chan struct{})
	c.connWriteDoneCh <- connLoopWriteCh
	<-connLoopWriteCh
	log.Infof("connection loop write-side closed")
}

func (c *TelnetConn) closeFSM() {
	fsmCh := make(chan struct{})
	c.fsm.doneCh <- fsmCh
	<-fsmCh
}

func (c *TelnetConn) closeDatahandler() {
	dataCh := make(chan struct{})
	c.dataHandlerCloseCh <- dataCh
	<-dataCh
}

func (c *TelnetConn) sendCmd(cmd byte, opt byte) {
	b := []byte{IAC, cmd, opt}
	log.Infof("Sending command: %v %v", cmd, opt)
	c.writeCh <- b
}

func (c *TelnetConn) handleOptionCommand(cmd byte, opt byte) {
	if cmd == WILL || cmd == WONT {
		if _, ok := c.clientOpts[opt]; !ok {
			c.sendCmd(DONT, opt)
			return
		}

		if _, ok := c.unackedClientOpts[opt]; ok {
			delete(c.unackedClientOpts, opt)
			if len(c.unackedClientOpts) == 0 && len(c.unackedServerOpts) == 0 {
				close(c.negotiationDone)
			}
		} else {
			c.sendCmd(DO, opt)
		}
	}

	if cmd == DO || cmd == DONT {
		if _, ok := c.serverOpts[opt]; !ok {
			c.sendCmd(WONT, opt)
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
			c.sendCmd(WILL, opt)
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
