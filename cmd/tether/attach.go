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

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/vmware/vic/cmd/tether/msgs"
	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/pkg/serial"
	"github.com/vmware/vic/pkg/trace"
)

const (
	attachChannelType = "attach"
)

// server is the singleton attachServer for the tether - there can be only one
// as the backchannel line protocol may not provide multiplexing of connections
var server AttachServer
var once sync.Once

type AttachServer interface {
	tether.Extension

	start() error
	stop() error
}

type attachServerSSH struct {
	// serializes data access for exported functions
	m sync.Mutex

	// conn is held directly as it is how we stop the attach server
	conn struct {
		// serializes data access for the underlying conn
		sync.Mutex
		conn net.Conn
	}

	config    *tether.ExecutorConfig
	sshConfig *ssh.ServerConfig

	enabled int32

	// Cancelable context and its cancel func. Used for resolving the deadlock
	// between run() and stop()
	ctx    context.Context
	cancel context.CancelFunc

	// INTERNAL: must set by testAttachServer only
	testing bool
}

// NewAttachServerSSH either creates a new instance or returns the initialized one
func NewAttachServerSSH() AttachServer {
	once.Do(func() {
		// create a cancelable context and assign it to the CancelFunc
		// it isused for resolving the deadlock between run() and stop()
		// it has a Background parent as we don't want timeouts here,
		// otherwise we may start leaking goroutines in the handshake code
		ctx, cancel := context.WithCancel(context.Background())
		server = &attachServerSSH{
			ctx:    ctx,
			cancel: cancel,
		}
	})
	return server
}

// Reload - tether.Extension implementation
func (t *attachServerSSH) Reload(config *tether.ExecutorConfig) error {
	defer trace.End(trace.Begin("attach reload"))

	t.m.Lock()
	defer t.m.Unlock()

	t.config = config

	err := server.start()
	if err != nil {
		detail := fmt.Sprintf("unable to start attach server: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}
	return nil
}

// Enabled sets the enabled to true
func (t *attachServerSSH) Enabled() {
	atomic.StoreInt32(&t.enabled, 1)
}

// Disabled sets the enabled to false
func (t *attachServerSSH) Disabled() {
	atomic.StoreInt32(&t.enabled, 0)
}

// IsEnabled returns whether the enabled is true
func (t *attachServerSSH) IsEnabled() bool {
	return atomic.LoadInt32(&t.enabled) == 1
}

// Start is implemented at _ARCH.go files

// Stop needed for tether.Extensions interface
func (t *attachServerSSH) Stop() error {
	defer trace.End(trace.Begin("stop attach server"))

	t.m.Lock()
	defer t.m.Unlock()

	// calling server.start not t.start so that test impl gets invoked
	return server.stop()
}

func (t *attachServerSSH) start() error {
	defer trace.End(trace.Begin("start attach server"))

	if t == nil {
		err := fmt.Errorf("attach server is not configured")
		log.Error(err)
		return err
	}

	if t.IsEnabled() {
		err := fmt.Errorf("attach server is already enabled")
		log.Error(err)
		return err
	}

	// don't assume that the key hasn't changed
	pkey, err := ssh.ParsePrivateKey([]byte(t.config.Key))
	if err != nil {
		detail := fmt.Sprintf("failed to load key for attach: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	// TODO: update this with generated credentials for the appliance
	t.sshConfig = &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if c.User() == "daemon" {
				return &ssh.Permissions{}, nil
			}
			return nil, fmt.Errorf("expected daemon user")
		},
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == "daemon" {
				return &ssh.Permissions{}, nil
			}
			return nil, fmt.Errorf("expected daemon user")
		},
		NoClientAuth: true,
	}
	t.sshConfig.AddHostKey(pkey)

	t.Enabled()
	go t.run()

	return nil
}

// stop is not thread safe with start
func (t *attachServerSSH) stop() error {
	defer trace.End(trace.Begin("stop attach server"))

	if t == nil {
		err := fmt.Errorf("attach server is not configured")
		log.Error(err)
		return err
	}

	if !t.IsEnabled() {
		err := fmt.Errorf("attach server is not enabled")
		log.Error(err)
		return err
	}

	log.Debugf("Setting enabled to false")
	t.Disabled()

	// This context is used by backchannel only. We need to cancel it before
	// trying to obtain the following lock so that backchannel interrupts the
	// underlying Read call by calling Close on it.
	// The lock is  held by backchannel's caller and not released until it returns
	log.Debugf("Canceling the context")
	t.cancel()

	log.Debugf("Acquiring the connection lock")
	t.conn.Lock()
	if t.conn.conn != nil {
		log.Debugf("Close called again on rawconn - squashing")
		t.conn.conn.Close()
		t.conn.conn = nil
	}
	t.conn.Unlock()
	log.Debugf("Released the connection lock")

	return nil
}

func backchannel(ctx context.Context, conn net.Conn) error {
	defer trace.End(trace.Begin("establish tether backchannel"))

	// HACK: currently RawConn dosn't implement timeout so throttle the spinning
	// it does implement the Timeout methods so the intermediary code can be written
	// to support it, but they are stub implementation in rawconn impl.

	// This needs to tick *faster* than the ticker in connection.go on the
	// portlayer side.  The PL sends the first syn and if this isn't waiting,
	// alignment will take a few rounds (or it may never happen).
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	// We run this in a separate goroutine because HandshakeServer
	// calls a Read on rawconn which is a blocking call which causes
	// the caller to block as well so this is the only way to cancel.
	// Calling Close() will unblock us and on the next tick we will
	// return ctx.Err()
	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		}
	}()

	for {
		select {
		case <-ticker.C:
			if ctx.Err() != nil {
				return ctx.Err()
			}
			deadline, ok := ctx.Deadline()
			if ok {
				conn.SetReadDeadline(deadline)
			}

			err := serial.HandshakeServer(conn)
			if err == nil {
				conn.SetReadDeadline(time.Time{})
				return nil
			}

			switch et := err.(type) {
			case *serial.HandshakeError:
				log.Debugf("HandshakeServer: %v", et)
			default:
				log.Errorf("HandshakeServer: %v", err)
			}
		}
	}
}

// run should not be called directly, but via start
// run will establish an ssh server listening on the backchannel
func (t *attachServerSSH) run() error {
	defer trace.End(trace.Begin("main attach server loop"))

	// we pass serverConn to the channelMux goroutine so we need to lock it
	var serverConn struct {
		sync.Mutex
		*ssh.ServerConn
	}
	var established bool

	var chans <-chan ssh.NewChannel
	var reqs <-chan *ssh.Request
	var err error

	for t.IsEnabled() {
		serverConn.Lock()
		established = serverConn.ServerConn != nil
		serverConn.Unlock()

		// keep waiting for the connection to establish
		for !established && t.IsEnabled() {
			log.Debugf("Trying to establish a connection")

			establishFn := func() error {
				// we hold the t.conn.Lock during the scope of this function
				t.conn.Lock()
				defer t.conn.Unlock()

				// tests are passing their own connections so do not create connections when testing is set
				if !t.testing {
					// close the connection if required
					if t.conn.conn != nil {
						t.conn.conn.Close()
						t.conn.conn = nil
					}
					t.conn.conn, err = rawConnectionFromSerial()
					if err != nil {
						detail := fmt.Errorf("failed to create raw connection raw connection: %s", err)
						log.Error(detail)
						return detail
					}
				}

				// wait for backchannel to establish
				err = backchannel(t.ctx, t.conn.conn)
				if err != nil {
					detail := fmt.Errorf("failed to establish backchannel: %s", err)
					log.Error(detail)
					return detail
				}

				// create the SSH server using underlying t.conn
				serverConn.Lock()
				defer serverConn.Unlock()

				serverConn.ServerConn, chans, reqs, err = ssh.NewServerConn(t.conn.conn, t.sshConfig)
				if err != nil {
					detail := fmt.Errorf("failed to establish ssh handshake: %s", err)
					log.Error(detail)
					return detail
				}

				return nil
			}
			established = establishFn() == nil
		}

		defer func() {
			log.Debugf("cleanup on connection")

			serverConn.Lock()
			defer serverConn.Unlock()

			if serverConn.ServerConn != nil {
				log.Debugf("closing underlying connection")
				serverConn.Close()
			}
		}()

		// Global requests
		go t.globalMux(reqs)

		log.Println("ready to service attach requests")
		// Service the incoming channels
		for attachchan := range chans {
			// The only channel type we'll support is attach
			if attachchan.ChannelType() != attachChannelType {
				detail := fmt.Sprintf("unknown channel type %s", attachchan.ChannelType())
				log.Error(detail)
				attachchan.Reject(ssh.UnknownChannelType, "unknown channel type")
				continue
			}

			// check we have a Session matching the requested ID
			bytes := attachchan.ExtraData()
			if bytes == nil {
				detail := "attach channel requires ID in ExtraData"
				log.Error(detail)
				attachchan.Reject(ssh.Prohibited, detail)
				continue
			}

			sessionid := string(bytes)
			session, ok := t.config.Sessions[sessionid]

			reason := ""
			if !ok {
				reason = "is unknown"
			} else if session.Cmd.Process == nil {
				reason = "process has not been launched"
			} else if session.Cmd.Process.Signal(syscall.Signal(0)) != nil {
				reason = "process has exited"
			}

			if reason != "" {
				detail := fmt.Sprintf("attach request: session %s %s", sessionid, reason)
				log.Error(detail)
				attachchan.Reject(ssh.Prohibited, detail)
				continue
			}

			log.Infof("accepting incoming channel for %s", sessionid)
			channel, requests, err := attachchan.Accept()
			log.Debugf("accepted incoming channel for %s", sessionid)
			if err != nil {
				detail := fmt.Sprintf("could not accept channel: %s", err)
				log.Errorf(detail)
				continue
			}

			// bind the channel to the Session
			log.Debugf("binding reader/writers for channel for %s", sessionid)
			session.Outwriter.Add(channel)
			session.Reader.Add(channel)

			// cleanup on detach from the session
			detach := func() {
				log.Debugf("cleanup on detach from the session")

				session.Outwriter.Remove(channel)
				session.Reader.Remove(channel)

				channel.Close()

				serverConn.Lock()
				serverConn.ServerConn = nil
				serverConn.Unlock()
			}

			// tty's merge stdout and stderr so we don't bind an additional reader in that case
			// but we need to do so for non-tty
			if session.Pty == nil {
				session.Errwriter.Add(channel.Stderr())

				// no good way to function chain, so reimplement appropriately
				detach = func() {
					log.Debugf("cleanup on detach from the session")

					session.Outwriter.Remove(channel)
					session.Reader.Remove(channel)
					session.Errwriter.Remove(channel)

					channel.Close()

					serverConn.Lock()
					serverConn.ServerConn = nil
					serverConn.Unlock()
				}
			}
			log.Debugf("reader/writers bound for channel for %s", sessionid)

			go t.channelMux(requests, session, detach)
		}

		log.Info("incoming attach channel closed")
	}
	return nil
}

func (t *attachServerSSH) globalMux(reqchan <-chan *ssh.Request) {
	defer trace.End(trace.Begin("attach server global request handler"))

	for req := range reqchan {
		var pendingFn func()
		var payload []byte
		ok := true

		log.Infof("received global request type %v", req.Type)

		switch req.Type {
		case msgs.ContainersReq:
			keys := make([]string, len(t.config.Sessions))
			i := 0
			for k := range t.config.Sessions {
				keys[i] = k
				i++
			}
			msg := msgs.ContainersMsg{IDs: keys}
			payload = msg.Marshal()
		default:
			ok = false
			payload = []byte("unknown global request type: " + req.Type)
		}

		log.Debugf("Returning payload: %s", string(payload))

		// make sure that errors get send back if we failed
		if req.WantReply {
			req.Reply(ok, payload)
		}

		// run any pending work now that a reply has been sent
		if pendingFn != nil {
			log.Debug("Invoking pending work")
			go pendingFn()
			pendingFn = nil
		}
	}
}

func (t *attachServerSSH) channelMux(in <-chan *ssh.Request, session *tether.SessionConfig, detach func()) {
	defer trace.End(trace.Begin("attach server channel request handler"))

	// detach
	defer detach()

	process := session.Cmd.Process
	pty := session.Pty

	var err error
	for req := range in {
		var pendingFn func()
		ok := true

		switch req.Type {
		case msgs.WindowChangeReq:
			msg := msgs.WindowChangeMsg{}
			if pty == nil {
				ok = false
				log.Errorf("illegal window-change request for non-tty")
			} else if err = msg.Unmarshal(req.Payload); err != nil {
				ok = false
				log.Errorf(err.Error())
			} else if err = resizePty(pty.Fd(), &msg); err != nil {
				ok = false
				log.Errorf(err.Error())
			}
		case msgs.SignalReq:
			msg := msgs.SignalMsg{}
			if err = msg.Unmarshal(req.Payload); err != nil {
				ok = false
				log.Errorf(err.Error())
			} else {
				log.Infof("Sending signal %s to container process, pid=%d\n", string(msg.Signal), process.Pid)
				err = signalProcess(process, msg.Signal)
				if err != nil {
					log.Errorf("Failed to dispatch signal to process: %s\n", err)
				}
			}
		case msgs.CloseStdinReq:
			// call Close as the pendingFn so that we can send reply back before closing the channel
			pendingFn = func() {
				log.Debugf("Closing stdin for %s", session.ID)
				session.Reader.Close()
			}
		default:
			ok = false
			err = fmt.Errorf("ssh request type %s is not supported", req.Type)
			log.Error(err.Error())
		}

		// payload is ignored on channel specific replies.  The ok is passed, however.
		if req.WantReply {
			log.Debugf("Sending reply %t back", ok)
			req.Reply(ok, nil)
		}

		// run any pending work now that a reply has been sent
		if pendingFn != nil {
			log.Debug("Invoking pending work")
			go pendingFn()
			pendingFn = nil
		}
	}

}

// The syscall struct
type winsize struct {
	wsRow    uint16
	wsCol    uint16
	wsXpixel uint16
	wsYpixel uint16
}
