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
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/serial"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

// Connection represents a communication channel initiated by the client TO the
// client.  The client connects (via TCP) to the server, then the server
// initiates an SSH connection over the same sock to the client.
type Connection struct {
	*ssh.Client

	// the container's ID
	id string

	requests <-chan *ssh.Request
}

type Connector struct {
	mutex       sync.RWMutex
	cond        *sync.Cond
	connections map[string]*Connection

	listener net.Listener
	// Quit channel for listener routine
	listenerQuit chan bool
	wg           sync.WaitGroup
}

// On connect from a client (over TCP), attempt to SSH (over the same sock) to the client.
func NewConnector(listener net.Listener) *Connector {
	connector := &Connector{
		connections:  make(map[string]*Connection),
		listener:     listener,
		listenerQuit: make(chan bool),
	}
	connector.cond = sync.NewCond(connector.mutex.RLocker())

	connector.wg.Add(1)
	go connector.serve()

	return connector
}

// Returns a connection corresponding to the specified ID. If the connection doesn't exist
// the method will wait for the specified timeout, returning when the connection is created
// or the timeout expires, whichever occurs first
func (c *Connector) Get(ctx context.Context, id string, timeout time.Duration) (*Connection, error) {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	c.mutex.RLock()
	conn := c.connections[id]
	c.mutex.RUnlock()
	if conn != nil {
		return conn, nil
	} else if timeout == 0 {
		return nil, fmt.Errorf("no such connection")
	}

	result := make(chan *Connection, 1)

	go func() {
		ok := false
		var conn *Connection

		for !ok && ctx.Err() == nil {
			c.mutex.RLock()
			defer c.mutex.RUnlock()

			conn, ok = c.connections[id]
			if ok {
				log.Debugf("Found connection for %s: %p", id, conn)
				result <- conn
				return
			}

			// block until cond is updated
			c.cond.Wait()
		}
	}()

	select {
	case client := <-result:
		log.Debugf("Found connection for %s: %p", id, client)
		return client, nil
	case <-ctx.Done():
		err := fmt.Errorf("id:%s: %s", id, ctx.Err())
		log.Error(err)
		// wake up the result gofunc before returning
		c.cond.Broadcast()
		return nil, err
	}
}

func (c *Connector) Remove(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.connections[id] != nil {
		c.connections[id].Close()
		delete(c.connections, id)
	}
}

// takes the base connection, determines the ID of the source and stashes it in the map
func (c *Connector) processIncoming(conn net.Conn) {

	for {
		if conn == nil {
			log.Infof("connection closed")
			return
		}
		defer conn.Close()

		// TODO needs timeout handling.  This could take 30s.
		ctx, cancel := context.WithTimeout(context.TODO(), 50*time.Millisecond)
		if err := serial.HandshakeClient(ctx, conn); err == nil {
			log.Infof("New connection")
			cancel()
			break
		} else if err == io.EOF {
			log.Debugf("caught EOF")
			return
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
	ccon, newchan, request, err := ssh.NewClientConn(conn, "", config)

	if err != nil {
		log.Errorf("SSH connection could not be established: %s", errors.ErrorStack(err))
		return
	}

	// get the container id
	ok, reply, err := ccon.SendRequest(containerID, true, nil)
	if ok || err != nil {
		log.Errorf("Error retrieving container ID via global ssh request: %s", errors.ErrorStack(err))
		return
	}

	log.Debugf("Established connection with container VM: %s", string(reply))

	c.mutex.Lock()
	connection := &Connection{
		Client:   ssh.NewClient(ccon, newchan, nil),
		id:       string(reply),
		requests: request,
	}

	c.connections[connection.id] = connection

	c.cond.Broadcast()
	c.mutex.Unlock()

	return
}

// Starts the connector listening on the specified source
// TODO: should have mechanism for stopping this, and probably handing off the connections to another
// routine to insert into the map
func (c *Connector) serve() {
	defer c.wg.Done()
	for {
		if c.listener == nil {
			log.Debugf("listener closed")
			break
		}

		conn, err := c.listener.Accept()

		select {
		case <-c.listenerQuit:
			log.Debugf("serve exitting")
			return
		default:
		}

		if err != nil {
			log.Errorf("Error waiting for incoming connection: %s", errors.ErrorStack(err))
			continue
		}

		go c.processIncoming(conn)
	}
}

func (c *Connector) Stop() {
	close(c.listenerQuit)
	c.wg.Wait()
}

func (c *Connector) URL() string {
	addr := c.listener.Addr()
	return fmt.Sprintf("%s://%s", addr.Network(), addr.String())
}
