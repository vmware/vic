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

package attach

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/serial"
	"github.com/vmware/vic/pkg/trace"

	"golang.org/x/crypto/ssh"
)

type Connector struct {
	mutex       sync.Mutex
	connections map[string]SessionInteraction
	gets        []chan struct{}

	listener net.Listener
	// Quit channel for listener routine
	listenerQuit chan bool
	wg           sync.WaitGroup

	// enable extra debug on the line
	debug bool
}

// On connect from a client (over TCP), attempt to SSH (over the same sock) to the client.
func NewConnector(listener net.Listener, debug bool) *Connector {
	defer trace.End(trace.Begin(""))

	connector := &Connector{
		connections:  make(map[string]SessionInteraction),
		listener:     listener,
		listenerQuit: make(chan bool),
		debug:        debug,
	}

	connector.wg.Add(1)
	go connector.serve()

	return connector
}

// Returns a connection corresponding to the specified ID. If the connection doesn't exist
// the method will wait for the specified timeout, returning when the connection is created
// or the timeout expires, whichever occurs first
func (c *Connector) Get(ctx context.Context, id string, timeout time.Duration) (SessionInteraction, error) {
	defer trace.End(trace.Begin(id))

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ch := make(chan struct{}, 1)
	defer func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		// remove ch from waiters
		for i := range c.gets {
			if c.gets[i] == ch {
				c.gets = append(c.gets[:i], c.gets[i+1:]...)
				return
			}
		}
	}()

	for {
		c.mutex.Lock()

		if conn, ok := c.connections[id]; ok && !conn.IsClosed() {
			c.mutex.Unlock()
			return conn, nil
		}

		c.gets = append(c.gets, ch)
		c.mutex.Unlock()
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out waiting for connection for container %s", id)
		case <-ch:
		}
	}
}

func (c *Connector) Remove(id string) error {
	defer trace.End(trace.Begin(id))

	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error

	if conn := c.connections[id]; conn != nil {
		err = conn.Close()
		delete(c.connections, id)
	}
	return err
}

// takes the base connection, determines the ID of the source and stashes it in the map
func (c *Connector) processIncoming(conn net.Conn) {
	var err error
	defer func() {
		if err != nil && conn != nil {
			conn.Close()
		}
	}()

	for {
		if conn == nil {
			log.Infof("attach connector: connection closed")
			return
		}

		serial.PurgeIncoming(conn)

		// TODO needs timeout handling.  This could take 30s.

		// Timeout for client handshake should be reasonably small.
		// Server will try to drain a buffer and if the buffer doesn't contain
		// 2 or more bytes it will just wait, so client should timeout.
		// However, if timeout is too short, client will flood server with Syn requests.
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()

		deadline, ok := ctx.Deadline()
		if ok {
			conn.SetReadDeadline(deadline)
		}
		if err = serial.HandshakeClient(conn, c.debug); err == nil {
			conn.SetReadDeadline(time.Time{})
			log.Debugf("attach connector: New connection")
			cancel()
			break
		} else if err == io.EOF {
			log.Debugf("caught EOF")
			conn.Close()
			return
		} else if _, ok := err.(*serial.HandshakeError); ok {
			log.Debugf("HandshakeClient: %v", err)
		} else {
			log.Errorf("HandshakeClient: %v", err)
		}
	}

	callback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	config := &ssh.ClientConfig{
		User:            "daemon",
		HostKeyCallback: callback,
	}

	log.Debugf("Initiating ssh handshake with new connection attempt")
	var (
		ccon    ssh.Conn
		newchan <-chan ssh.NewChannel
		request <-chan *ssh.Request
	)

	ccon, newchan, request, err = ssh.NewClientConn(conn, "", config)
	if err != nil {
		log.Errorf("SSH connection could not be established: %s", errors.ErrorStack(err))
		return
	}

	client := ssh.NewClient(ccon, newchan, request)

	var ids []string
	ids, err = SSHls(client)
	if err != nil {
		log.Errorf("SSH connection could not be established: %s", errors.ErrorStack(err))
		return
	}

	var si SessionInteraction
	for _, id := range ids {
		si, err = SSHAttach(client, id)
		if err != nil {
			log.Errorf("SSH connection could not be established (id=%s): %s", id, errors.ErrorStack(err))
			return
		}

		log.Infof("Established connection with container VM: %s", id)

		c.mutex.Lock()
		c.connections[id] = si
		for i := range c.gets {
			c.gets[i] <- struct{}{}
		}
		c.mutex.Unlock()
	}

	return
}

// Starts the connector listening on the specified source
// TODO: should have mechanism for stopping this, and probably handing off the connections to another
// routine to insert into the map
func (c *Connector) serve() {
	defer c.wg.Done()
	for {
		if c.listener == nil {
			log.Debugf("attach connector: listener closed")
			break
		}

		conn, err := c.listener.Accept()

		select {
		case <-c.listenerQuit:
			log.Debugf("attach connector: serve exitting")
			return
		default:
		}

		if err != nil {
			log.Errorf("Error waiting for incoming connection: %s", errors.ErrorStack(err))
			continue
		}

		log.Info("attach connector: Received incoming connection")
		go c.processIncoming(conn)
	}
}

func (c *Connector) Stop() {
	defer trace.End(trace.Begin(""))

	c.listener.Close()
	close(c.listenerQuit)
	c.wg.Wait()
}

func (c *Connector) URL() string {
	defer trace.End(trace.Begin(""))

	addr := c.listener.Addr()
	return fmt.Sprintf("%s://%s", addr.Network(), addr.String())
}
