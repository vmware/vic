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
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/portlayer/attach"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

const (
	attachChannelType = "attach"
)

// server is the singleton attachServer for the tether - there can be only one
// as the backchannel line protocol may not provide multiplexing of connections
var server attachServer

type attachServer interface {
	start() error
	stop()
}

// conn is held directly as it's how we stop the attach server
type attachServerSSH struct {
	conn   *net.Conn
	config *ssh.ServerConfig

	enabled bool
}

// start is not thread safe with stop
func (t *attachServerSSH) start() error {
	defer trace.End(trace.Begin("start attach server"))

	if t == nil {
		return errors.New("attach server is not configured")
	}

	if t.enabled {
		return nil
	}

	// don't assume that the key hasn't changed
	pkey, err := ssh.ParsePrivateKey([]byte(config.Key))
	if err != nil {
		detail := fmt.Sprintf("failed to load key for attach: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	// TODO: update this with generated credentials for the appliance
	t.config = &ssh.ServerConfig{
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
	t.config.AddHostKey(pkey)

	t.enabled = true
	go t.run()

	return nil
}

// stop is not thread safe with start
func (t *attachServerSSH) stop() {
	defer trace.End(trace.Begin("stop attach server"))

	if t == nil || !t.enabled {
		return
	}

	t.enabled = false
	if t.conn != nil {
		(*t.conn).Close()
		t.conn = nil
	}
}

// run should not be called directly, but via start
// run will establish an ssh server listening on the backchannel
func (t *attachServerSSH) run() error {
	defer trace.End(trace.Begin("main attach server loop"))

	var sConn *ssh.ServerConn
	var chans <-chan ssh.NewChannel
	var reqs <-chan *ssh.Request
	var err error

	// keep waiting for the connection to establish
	for t.enabled && sConn == nil {
		// wait for backchannel to establish
		conn, errb := utils.backchannel(context.Background())
		if errb != nil {
			err = errb
			detail := fmt.Sprintf("failed to establish backchannel: %s", err)
			log.Error(detail)
			continue
		}
		t.conn = &conn

		// create the SSH server
		sConn, chans, reqs, err = ssh.NewServerConn(*t.conn, t.config)
		if err != nil {
			detail := fmt.Sprintf("failed to establish ssh handshake: %s", err)
			log.Error(detail)
			continue
		}
	}
	if err != nil {
		detail := fmt.Sprintf("abandoning attempt to start attach server: %s", err)
		log.Error(detail)
		return err
	}

	defer sConn.Close()

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
		session, ok := config.Sessions[sessionid]

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
		session.outwriter.Add(channel)
		session.reader.Add(channel)

		// cleanup on detach from the session
		detach := func() {
			session.outwriter.Remove(channel)
			session.reader.Remove(channel)
		}

		// tty's merge stdout and stderr so we don't bind an additional reader in that case
		// but we need to do so for non-tty
		if session.pty == nil {
			session.errwriter.Add(channel.Stderr())

			// no good way to function chain, so reimplement appropriately
			detach = func() {
				session.outwriter.Remove(channel)
				session.reader.Remove(channel)
				session.errwriter.Remove(channel)
			}
		}
		log.Debugf("reader/writers bound for channel for %s", sessionid)

		go t.channelMux(requests, session.Cmd.Process, session.pty, detach)
	}

	log.Info("incoming attach channel closed")

	return nil
}

func (t *attachServerSSH) globalMux(reqchan <-chan *ssh.Request) {
	defer trace.End(trace.Begin("start attach server global request handler"))

	for req := range reqchan {
		var pendingFn func()
		var payload []byte
		ok := true

		log.Infof("received global request type %v", req.Type)

		switch req.Type {
		case attach.ContainersReq:
			keys := make([]string, len(config.Sessions))
			i := 0
			for k := range config.Sessions {
				keys[i] = k
				i++
			}
			msg := attach.ContainersMsg{IDs: keys}
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

func (t *attachServerSSH) channelMux(in <-chan *ssh.Request, process *os.Process, pty *os.File, detach func()) {
	defer trace.End(trace.Begin("start attach server channel request handler"))

	var err error
	for req := range in {
		var pendingFn func()
		ok := true

		switch req.Type {
		case attach.WindowChangeReq:
			msg := attach.WindowChangeMsg{}
			if pty == nil {
				ok = false
				log.Errorf("illegal window-change request for non-tty")
			} else if err = msg.Unmarshal(req.Payload); err != nil {
				ok = false
				log.Errorf(err.Error())
			} else if err = utils.resizePty(pty.Fd(), &msg); err != nil {
				ok = false
				log.Errorf(err.Error())
			}
		case attach.SignalReq:
			msg := attach.SignalMsg{}
			if err = msg.Unmarshal(req.Payload); err != nil {
				ok = false
				log.Errorf(err.Error())
			} else {
				log.Infof("Sending signal %s to container process, pid=%d\n", string(msg.Signal), process.Pid)
				err = utils.signalProcess(process, msg.Signal)
				if err != nil {
					log.Errorf("Failed to dispatch signal to process: %s\n", err)
				}
			}
		default:
			ok = false
			err = fmt.Errorf("ssh request type %s is not supported", req.Type)
			log.Error(err.Error())
		}

		// payload is ignored on channel specific replies.  The ok is passed, however.
		if req.WantReply {
			req.Reply(ok, nil)
		}

		// run any pending work now that a reply has been sent
		if pendingFn != nil {
			log.Debug("Invoking pending work")
			go pendingFn()
			pendingFn = nil
		}
	}

	detach()
}
