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
	"github.com/vmware/vic/pkg/dio"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

const (
	attachChannelType = "attach"
)

var (
	Signals = map[ssh.Signal]int{
		ssh.SIGABRT: 6,
		ssh.SIGALRM: 14,
		ssh.SIGFPE:  8,
		ssh.SIGHUP:  1,
		ssh.SIGILL:  4,
		ssh.SIGINT:  2,
		ssh.SIGKILL: 9,
		ssh.SIGPIPE: 13,
		ssh.SIGQUIT: 3,
		ssh.SIGSEGV: 11,
		ssh.SIGTERM: 15,
		ssh.SIGUSR1: 10,
		ssh.SIGUSR2: 12,
	}
)

// WindowChangeMsg the RFC4254 struct
type WindowChangeMsg struct {
	Columns  uint32
	Rows     uint32
	WidthPx  uint32
	HeightPx uint32
}

type stringMsg struct {
	Signal string
}

type stringArrayMsg struct {
	Strings []string
}

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
	if t == nil {
		return errors.New("attach server is not configured")
	}

	if t.enabled {
		return nil
	}

	// don't assume that the key hasn't changed
	pkey, err := ssh.ParsePrivateKey(config.Key)
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
		live, ok := sessions[sessionid]

		reason := ""
		if !ok || live.cmd == nil {
			reason = "is unknown"
		} else if live.cmd.Process == nil {
			reason = "process has not been launched"
		} else if live.cmd.Process.Signal(syscall.Signal(0)) != nil {
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
		if live.pty == nil {
			// if it's not a TTY then bind the channel directly to the multiwriter that's already associated with the process
			dmwStdout, okA := live.cmd.Stdout.(dio.DynamicMultiWriter)
			dmwStderr, okB := live.cmd.Stderr.(dio.DynamicMultiWriter)
			dmrStdin, okC := live.cmd.Stdin.(dio.DynamicMultiReader)
			if !okA || !okB || !okC {
				detail := fmt.Sprintf("target session IO cannot be duplicated to attach streams: %s", string(bytes))
				log.Error(detail)
				attachchan.Reject(ssh.ConnectionFailed, detail)
				continue
			}

			log.Debugf("binding reader/writers for channel for %s", sessionid)
			dmwStdout.Add(channel)
			dmwStderr.Add(channel.Stderr())
			dmrStdin.Add(channel)
			log.Debugf("reader/writers bound for channel for %s", sessionid)

			// cleanup on detach from the session
			detach := func() {
				dmwStdout.Remove(channel)
				dmwStderr.Remove(channel.Stderr())
				dmrStdin.Remove(channel)
			}
			go t.channelMux(requests, live.cmd.Process, nil, detach)
			continue
		}

		// if it's a TTY bind the channel to the multiwriter that's on the far side of the PTY from the process
		// this is done so the logging is done with processed output
		live.outwriter.Add(channel)
		// PTY merges stdout & stderr so the two are the same
		live.reader.Add(channel)

		// cleanup on detach from the session
		detach := func() {
			live.outwriter.Remove(channel)
			live.reader.Remove(channel)
		}
		go t.channelMux(requests, live.cmd.Process, live.pty, detach)
	}

	log.Info("incoming attach channel closed")

	return nil
}

func (t *attachServerSSH) globalMux(reqchan <-chan *ssh.Request) {
	for req := range reqchan {
		var pendingFn func()
		var payload []byte
		ok := true

		log.Infof("received global request type %v", req.Type)

		switch req.Type {
		case "container-ids":
			keys := make([]string, len(config.Sessions))
			i := 0
			for k := range config.Sessions {
				keys[i] = k
				i++
			}
			msg := stringArrayMsg{Strings: keys}

			payload = []byte(ssh.Marshal(msg))
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
	var err error
	for req := range in {
		var pendingFn func()
		var payload []byte
		ok := true

		switch req.Type {
		case "window-change":
			msg := WindowChangeMsg{}
			if pty == nil {
				ok = false
				payload = []byte("illegal window-change request for non-tty")
			} else if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else if err = utils.resizePty(pty.Fd(), &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			}
		case "signal":
			msg := stringMsg{}
			if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				log.Infof("Sending signal %s to container process, pid=%d\n", string(msg.Signal), process.Pid)
				err = utils.signalProcess(process, ssh.Signal(msg.Signal))
				if err != nil {
					log.Errorf("Failed to dispatch signal to process: %s\n", err)
				}
				payload = []byte{}
			}
		default:
			ok = false
			err = fmt.Errorf("ssh request type %s is not supported", req.Type)
			log.Error(err.Error())
		}

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

	detach()

	log.Info("incoming attach request channel closed")
}
