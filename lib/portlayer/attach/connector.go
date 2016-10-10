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
	"github.com/vmware/vic/pkg/trace"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

// Connection represents a communication channel initiated by the client TO the
// client.  The client connects (via TCP) to the server, then the server
// initiates an SSH connection over the same sock to the client.
type Connection struct {
	spty SessionInteraction

	// the container's ID
	id string
}

type Connector struct {
	mutex       sync.RWMutex
	cond        *sync.Cond
	connections map[string]*Connection

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
		connections:  make(map[string]*Connection),
		listener:     listener,
		listenerQuit: make(chan bool),
		debug:        debug,
	}
	connector.cond = sync.NewCond(connector.mutex.RLocker())

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

	c.mutex.RLock()
	conn := c.connections[id]
	c.mutex.RUnlock()
	if conn != nil {
		return conn.spty, nil
	} else if timeout == 0 {
		return nil, fmt.Errorf("no such connection")
	}

	result := make(chan *Connection, 1)

	go func() {
		ok := false
		var conn *Connection

		c.mutex.RLock()
		defer c.mutex.RUnlock()

		for !ok && ctx.Err() == nil {
			conn, ok = c.connections[id]
			if ok {
				result <- conn
				return
			}

			// block until cond is updated
			log.Infof("attach connector:  Connection not found yet for %s", id)
			c.cond.Wait()
		}
		log.Debugf("attach connector:  Giving up on connection for %s", id)
	}()

	select {
	case client := <-result:
		log.Debugf("attach connector: Found connection for %s: %p", id, client)
		return client.spty, nil
	case <-ctx.Done():
		err := fmt.Errorf("attach connector: Connection not found error for id:%s: %s", id, ctx.Err())
		log.Error(err)
		// wake up the result gofunc before returning
		c.cond.Broadcast()
		return nil, err
	}
}

func (c *Connector) Remove(id string) {
	defer trace.End(trace.Begin(id))

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.connections[id] != nil {
		if c.connections[id].id == id {
			c.connections[id].spty.Close()
		}
		delete(c.connections, id)
	}
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

		// This needs to timeout with a *longer* wait than the ticker set on
		// the tether side (in tether_linux.go) or alignment may not happen.
		// The PL sends the first SYN in the handshake and if the tether is not
		// waiting, the handshake may never succeed.
		ctx, cancel := context.WithTimeout(context.TODO(), 50*time.Millisecond)
		if err = serial.HandshakeClient(ctx, conn, c.debug); err == nil {
			log.Debugf("attach connector: New connection")
			cancel()
			break
		} else if err == io.EOF {
			log.Debugf("caught EOF")
			conn.Close()
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
		connection := &Connection{
			spty: si,
			id:   id,
		}

		c.connections[connection.id] = connection

		c.cond.Broadcast()
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

	close(c.listenerQuit)
	c.wg.Wait()
}

func (c *Connector) URL() string {
	defer trace.End(trace.Begin(""))

	addr := c.listener.Addr()
	return fmt.Sprintf("%s://%s", addr.Network(), addr.String())
}
