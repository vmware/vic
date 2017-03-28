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
	"io"
	"net"
	"sync"
	"sync/atomic"
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

	// Used for ordering events between global mux and channel mux
	askedAndAnswered chan struct{}

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

// Enable sets the enabled to true
func (t *attachServerSSH) Enable() {
	atomic.StoreInt32(&t.enabled, 1)
}

// Disable sets the enabled to false
func (t *attachServerSSH) Disable() {
	atomic.StoreInt32(&t.enabled, 0)
}

// Enabled returns whether the enabled is true
func (t *attachServerSSH) Enabled() bool {
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

	if t.Enabled() {
		err := fmt.Errorf("attach server is already enabled")
		log.Warn(err)
		return nil
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

	// enable the server and start it
	t.Enable()
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

	if !t.Enabled() {
		err := fmt.Errorf("attach server is not enabled")
		log.Error(err)
		return err
	}

	// disable the server
	t.Disable()

	// This context is used by backchannel only. We need to cancel it before
	// trying to obtain the following lock so that backchannel interrupts the
	// underlying Read call by calling Close on it.
	// The lock is  held by backchannel's caller and not released until it returns
	log.Debugf("Canceling AttachServer's context")
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

	for t.Enabled() {
		serverConn.Lock()
		established = serverConn.ServerConn != nil
		serverConn.Unlock()

		// keep waiting for the connection to establish
		for !established && t.Enabled() {
			log.Infof("Trying to establish a connection")

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
						detail := fmt.Errorf("failed to create raw connection: %s", err)
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

		log.Infof("Ready to service attach requests")
		// Service the incoming channels
		for attachchan := range chans {
			// The only channel type we'll support is attach
			if attachchan.ChannelType() != attachChannelType {
				detail := fmt.Sprintf("unknown channel type %s", attachchan.ChannelType())
				attachchan.Reject(ssh.UnknownChannelType, detail)
				log.Error(detail)
				continue
			}

			// check we have a Session matching the requested ID
			bytes := attachchan.ExtraData()
			if bytes == nil {
				detail := "attach channel requires ID in ExtraData"
				attachchan.Reject(ssh.Prohibited, detail)
				log.Error(detail)
				continue
			}

			sessionid := string(bytes)
			session, ok := t.config.Sessions[sessionid]
			if !ok {
				detail := fmt.Sprintf("session %s is invalid", sessionid)
				attachchan.Reject(ssh.Prohibited, detail)
				log.Error(detail)
				continue
			}

			channel, requests, err := attachchan.Accept()
			if err != nil {
				detail := fmt.Sprintf("could not accept channel: %s", err)
				log.Errorf(detail)
				continue
			}

			// bind the channel to the Session
			log.Debugf("binding reader/writers for channel for %s", sessionid)

			log.Debugf("Adding [%p] to Outwriter", channel)
			session.Outwriter.Add(channel)
			log.Debugf("Adding [%p] to Reader", channel)
			session.Reader.Add(channel)

			// cleanup on detach from the session
			cleanup := func() {
				log.Debugf("Cleanup on detach from the session")

				log.Debugf("Removing [%p] from Outwriter", channel)
				session.Outwriter.Remove(channel)

				log.Debugf("Removing [%p] from Reader", channel)
				session.Reader.Remove(channel)

				channel.Close()

				serverConn.Lock()
				serverConn.ServerConn = nil
				serverConn.Unlock()
			}

			detach := cleanup
			// tty's merge stdout and stderr so we don't bind an additional reader in that case but we need to do so for non-tty
			if !session.Tty {
				// persist the value as we end up with different values each time we access it
				stderr := channel.Stderr()

				log.Debugf("Adding [%p] to Errwriter", stderr)
				session.Errwriter.Add(stderr)

				detach = func() {
					log.Debugf("Cleanup on detach from the session (non-tty)")

					log.Debugf("Removing [%p] from Errwriter", stderr)
					session.Errwriter.Remove(stderr)

					cleanup()
				}
			}
			log.Debugf("reader/writers bound for channel for %s", sessionid)

			go t.channelMux(requests, session, detach)

			if session.RunBlock && session.ClearToLaunch != nil && session.Started != "true" {
				log.Debugf("Unblocking the launch of %s", sessionid)
				// make sure that portlayer received the container id back
				session.ClearToLaunch <- <-t.askedAndAnswered
				log.Debugf("Unblocked the launch of %s", sessionid)
			}
		}
		log.Info("Incoming attach channel closed")
	}
	return nil
}

func (t *attachServerSSH) globalMux(reqchan <-chan *ssh.Request) {
	defer trace.End(trace.Begin("attach server global request handler"))

	// to make sure we close the channel once
	var once sync.Once

	// ContainersReq will close this channel
	t.askedAndAnswered = make(chan struct{})

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

			// unblock ^ (above)
			pendingFn = func() {
				once.Do(func() {
					close(t.askedAndAnswered)
				})
			}
		default:
			ok = false
			payload = []byte("unknown global request type: " + req.Type)
		}

		log.Debugf("Returning payload: %s", string(payload))

		// make sure that errors get send back if we failed
		if req.WantReply {
			log.Debugf("Sending global request reply %t back with %#v", ok, payload)
			if err := req.Reply(ok, payload); err != nil {
				log.Warnf("Failed to reply a global request back")
			}
		}

		// run any pending work now that a reply has been sent
		if pendingFn != nil {
			log.Debug("Invoking pending work for global mux")
			go pendingFn()
			pendingFn = nil
		}
	}
}

func (t *attachServerSSH) channelMux(in <-chan *ssh.Request, session *tether.SessionConfig, cleanup func()) {
	defer trace.End(trace.Begin("attach server channel request handler"))

	// for the actions after we process the request
	var pendingFn func()

	// cleanup function passed by the caller
	defer cleanup()

	for req := range in {
		ok := true

		switch req.Type {
		case msgs.WindowChangeReq:
			session.Lock()
			pty := session.Pty
			session.Unlock()

			msg := msgs.WindowChangeMsg{}
			if pty == nil {
				ok = false
				log.Errorf("illegal window-change request for non-tty")
			} else if err := msg.Unmarshal(req.Payload); err != nil {
				ok = false
				log.Errorf(err.Error())
			} else if err := resizePty(pty.Fd(), &msg); err != nil {
				ok = false
				log.Errorf(err.Error())
			}
		case msgs.CloseStdinReq:
			// call Close as the pendingFn so that we can send reply back before closing the channel
			pendingFn = func() {
				session.Lock()
				defer session.Unlock()

				log.Debugf("Closing stdin for %s", session.ID)
				session.Reader.Close()
			}
		default:
			ok = false
			err := fmt.Errorf("ssh request type %s is not supported", req.Type)
			log.Error(err.Error())
		}

		// payload is ignored on channel specific replies.  The ok is passed, however.
		if req.WantReply {
			log.Debugf("Sending channel request reply %t back", ok)
			if err := req.Reply(ok, nil); err != nil {
				log.Warnf("Failed to reply a channel request back")
			}
		}

		// run any pending work now that a reply has been sent
		if pendingFn != nil {
			log.Debug("Invoking pending work for channel mux")
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

type stdinReader struct {
	reader io.ReadCloser
}

func newStdinReader(channel ssh.Channel) *stdinReader {
	r, w := io.Pipe()
	s := &stdinReader{
		reader: r,
	}

	go func() {
		buf := make([]byte, 2*1024)
		for {
			r, err := channel.Read(buf)
			log.Debugf("stdin r=%d, err=%s buf=%q", r, err, buf[:r])
			if r > 0 {
				w.Write(buf[:r])
			}

			if err != nil {
				w.CloseWithError(err)
				return
			}
		}
	}()

	return s
}

func (r *stdinReader) Read(buf []byte) (int, error) {
	defer trace.End(trace.Begin(""))

	n, err := r.reader.Read(buf)
	if err == io.ErrClosedPipe {
		err = io.EOF
	}

	return n, err
}

func (r *stdinReader) Close() error {
	defer trace.End(trace.Begin(""))

	r.reader.Close()
	return nil
}
