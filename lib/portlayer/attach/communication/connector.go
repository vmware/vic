// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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

package communication

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/cmd/tether/msgs"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/serial"
	"github.com/vmware/vic/pkg/trace"

	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/singleflight"
)

// Connector defines the connection and interactions
type Connector struct {
	mutex        sync.RWMutex
	cond         *sync.Cond
	interactions map[string]SessionInteractor

	listener net.Listener
	// Quit channel for serve
	done chan struct{}

	// deduplication of incoming calls
	fg singleflight.Group

	// graceful shutdown
	wg sync.WaitGroup

	// enable extra debug on the line
	debug bool
}

// NewConnector returns a new Connector
func NewConnector(listener net.Listener, debug bool) *Connector {
	defer trace.End(trace.Begin(""))

	connector := &Connector{
		interactions: make(map[string]SessionInteractor),
		listener:     listener,
		done:         make(chan struct{}),
		debug:        debug,
	}
	connector.cond = sync.NewCond(connector.mutex.RLocker())

	return connector
}

func (c *Connector) aliveAndKicking(ctx context.Context, id string) SessionInteractor {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conn := c.interactions[id]
	// we already established this connection, let's check it status
	if conn != nil {
		log.Infof("attach connector: Pinging for %s", id)

		if err := conn.Ping(); err == nil {
			log.Infof("attach connector: Pinged %s, returning", id)

			if err := conn.Unblock(); err == nil {
				log.Infof("attach connector: Unblocked %s, returning", id)
			}
			return conn
		}
		// ping failed so we need to remove it from the map
		log.Infof("attach connector: Ping test failed, removing %s from connection map", id)
		delete(c.interactions, id)
	}

	return nil
}

// Interaction returns the interactor corresponding to the specified ID. If the connection doesn't exist
// the method will wait for the specified timeout, returning when the connection is created
// or the timeout expires, whichever occurs first
func (c *Connector) Interaction(ctx context.Context, id string) (SessionInteractor, error) {
	defer trace.End(trace.Begin(id))

	// make sure that we have only one call in-flight for each ID at any given time
	si, err, shared := c.fg.Do(id, func() (interface{}, error) {
		return c.interaction(ctx, id)
	})
	if err != nil {
		c.fg.Forget(id)
		return nil, err
	}
	if shared {
		log.Debugf("Eliminated duplicated calls to Interaction for %s", id)
	}
	return si.(SessionInteractor), nil
}

func (c *Connector) interaction(ctx context.Context, id string) (SessionInteractor, error) {
	defer trace.End(trace.Begin(id))

	conn := c.aliveAndKicking(ctx, id)
	if conn != nil {
		return conn, nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("attach connector: no such connection")
	}

	result := make(chan SessionInteractor, 1)
	go func() {
		ok := false
		var conn SessionInteractor

		c.mutex.RLock()
		defer c.mutex.RUnlock()

		for !ok && ctx.Err() == nil {
			conn, ok = c.interactions[id]
			if ok {
				// no need to test this as we just created it
				if err := conn.Unblock(); err == nil {
					log.Infof("attach connector: Unblocked %s, returning", id)
				}
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
		return client, nil
	case <-ctx.Done():
		err := fmt.Errorf("attach connector: Connection not found error for id:%s: %s", id, ctx.Err())
		log.Error(err)
		// wake up the result gofunc before returning
		c.mutex.RLock()
		c.cond.Broadcast()
		c.mutex.RUnlock()

		return nil, err
	}
}

// RemoveInteraction removes the session the inteactions map
func (c *Connector) RemoveInteraction(id string) error {
	defer trace.End(trace.Begin(id))

	c.mutex.Lock()
	defer c.mutex.Unlock()

	var err error

	if c.interactions[id] != nil {
		log.Debugf("attach connector: Removing %s from the connection map", id)
		err = c.interactions[id].Close()
		delete(c.interactions, id)
		c.fg.Forget(id)
	}

	return err
}

// Start starts the connector
func (c *Connector) Start() {
	defer trace.End(trace.Begin(""))

	c.wg.Add(1)
	go c.serve()
}

// Stop stops the connector
func (c *Connector) Stop() {
	defer trace.End(trace.Begin(""))

	c.listener.Close()
	close(c.done)
	c.wg.Wait()
}

// URL returns the listener's URL
func (c *Connector) URL() string {
	defer trace.End(trace.Begin(""))

	addr := c.listener.Addr()
	return fmt.Sprintf("%s://%s", addr.Network(), addr.String())
}

// Starts the connector listening on the specified source
// TODO: should have mechanism for stopping this, and probably handing off the interactions to another
// routine to insert into the map
func (c *Connector) serve() {
	defer c.wg.Done()
	for {
		if c.listener == nil {
			log.Debugf("attach connector: listener closed")
			break
		}

		// check to see whether we should stop accepting new connections and exit
		select {
		case <-c.done:
			log.Debugf("attach connector: done closed")
			return
		default:
		}

		conn, err := c.listener.Accept()
		if err != nil {
			log.Errorf("Error waiting for incoming connection: %s", errors.ErrorStack(err))
			continue
		}
		log.Debugf("attach connector: Received incoming connection")

		go c.processIncoming(conn)
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

	log.Debugf("Initiating ssh handshake with new connection attempt")
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

	// create the SSH connection
	clientConn, chans, reqs, err := ssh.NewClientConn(conn, "", config)
	if err != nil {
		log.Errorf("SSH connection could not be established: %s", errors.ErrorStack(err))
		return
	}

	// ask the IDs
	ids, err := ContainerIDs(clientConn)
	if err != nil {
		log.Errorf("SSH connection could not be established: %s", errors.ErrorStack(err))
		return
	}

	// Handle global requests
	go c.reqs(reqs, clientConn, ids)
	// Handle channel open messages
	go c.chans(chans)

	// create the connections
	c.ids(clientConn, ids)

	return
}

// ids iterates over the gived ids and
// - calls Ping for existing connections
// - calls NewSSHInteraction for new connections and fills the connection map
func (c *Connector) ids(conn ssh.Conn, ids []string) {
	for _, id := range ids {
		c.mutex.RLock()
		si, ok := c.interactions[id]
		c.mutex.RUnlock()

		if ok {
			if err := si.Ping(); err == nil {
				log.Debugf("Connection %s found and alive", id)

				continue
			}
			log.Warnf("Connection found but it wasn't alive. Creating a new one")
		}

		// this is a new connection so learn the version
		version, err := ContainerVersion(conn)
		if err != nil {
			log.Errorf("SSH version could not be learned (id=%s): %s", id, errors.ErrorStack(err))
			return
		}

		si, err = NewSSHInteraction(conn, id, version)
		if err != nil {
			log.Errorf("SSH connection could not be established (id=%s): %s", id, errors.ErrorStack(err))
			return
		}

		log.Infof("Established connection with container VM: %s", id)

		c.mutex.Lock()

		c.interactions[id] = si

		c.cond.Broadcast()
		c.mutex.Unlock()
	}
}

// reqs is the global request channel of the portlayer side of the connection
// we keep a list of  sessions assosiacated with this connection and drop them from the map when the global mux exits
func (c *Connector) reqs(reqs <-chan *ssh.Request, conn ssh.Conn, ids []string) {
	defer trace.End(trace.Begin(""))

	var pending func()

	// list of session ids mux'ed on this connection
	droplist := make(map[string]struct{})

	// fill the map with the initial ids
	for _, id := range ids {
		droplist[id] = struct{}{}
	}

	for req := range reqs {
		ok := true

		log.Infof("received global request type %v", req.Type)
		switch req.Type {
		case msgs.ContainersReq:
			pending = func() {
				ids := msgs.ContainersMsg{}
				if err := ids.Unmarshal(req.Payload); err != nil {
					log.Errorf("Unmarshal failed with %s", err)
					return
				}
				c.ids(conn, ids.IDs)

				// fill the droplist with the latest info
				for _, id := range ids.IDs {
					droplist[id] = struct{}{}
				}
			}
		default:
			ok = false
		}

		// make sure that errors get send back if we failed
		if req.WantReply {
			log.Infof("Sending global request reply %t", ok)
			if err := req.Reply(ok, nil); err != nil {
				log.Warnf("Failed to reply a request back")
			}
		}

		// run any pending work now that a reply has been sent
		if pending != nil {
			log.Debug("Invoking pending work for global mux")
			go pending()
			pending = nil
		}
	}

	// global mux closed so it is time to do cleanup
	for id := range droplist {
		log.Infof("Droping %s from connection map", id)
		c.RemoveInteraction(id)
	}
}

// this is the channel mux for the ssh channel . It is configured to reject everything (required)
func (c *Connector) chans(chans <-chan ssh.NewChannel) {
	defer trace.End(trace.Begin(""))

	for ch := range chans {
		ch.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %v", ch.ChannelType()))
	}
}
