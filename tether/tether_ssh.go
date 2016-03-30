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

package tether

import (
	"fmt"
	"log"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
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

type setenvRequest struct {
	Name  string
	Value string
}

type execMsg struct {
	Command string
}

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

// This is the RFC4254 struct
type WindowChangeMsg struct {
	Columns   uint32
	Rows      uint32
	Width_px  uint32
	Height_px uint32
}

type subsystemRequestMsg struct {
	Subsystem string
}

type signalMsg struct {
	Signal string
}

type stringArrayMsg struct {
	Args []string
}

type nameValueMsg struct {
	Name  string
	Value string
}

func globalMux(reqchan <-chan *ssh.Request, handler GlobalContext) {
	for req := range reqchan {
		var pendingFn func()
		var payload []byte
		ok := true

		log.Printf("received global request type %v", req.Type)

		switch req.Type {
		case "container-id":
			payload = []byte(handler.ContainerId())
		case "ip-address":
			// set up arguments for a future exec call
			msg := stringArrayMsg{}
			if err := ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				args := msg.Args
				if len(args) == 2 {
					// this is synchronous as it should be non-interactive pass/fail with no data returned otherwise
					if err = handler.StaticIPAddress(args[0], args[1]); err != nil {
						ok = false
						payload = []byte(err.Error())
					}
				} else {
					ok = false
					payload = []byte("ip-address request requires CIDR address and gateway")
				}
			}
		case "dynamic-ip-address":
			if v4addr, err := handler.DynamicIPAddress(); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				payload = []byte(v4addr)
			}
		case "mount-label":
			msg := nameValueMsg{}
			if err := ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				// this is synchronous as it should be non-interactive pass/fail with no data returned otherwise
				if err = handler.MountLabel(msg.Name, msg.Value); err != nil {
					ok = false
					payload = []byte(err.Error())
				}
			}
		default:
			ok = false
			payload = []byte("unknown global request type: " + req.Type)
		}

		log.Printf("Returning payload: %s", string(payload))

		// make sure that errors get send back if we failed
		if req.WantReply {
			req.Reply(ok, payload)
		}

		// run any pending work now that a reply has been sent
		if pendingFn != nil {
			go pendingFn()
			pendingFn = nil
		}
	}
}

func channelMux(in <-chan *ssh.Request, handler SessionContext) {
	var err error
	args := []string{}
	config := make(map[string]string)
	for req := range in {
		var payload []byte
		ok := true

		switch req.Type {
		case "shell":
			// shell can't have a payload
			ok, payload = handler.Shell()
		case "pty-req":
			msg := ptyRequestMsg{}
			if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				handler.AssignPty()
			}
		case "window-change":
			msg := WindowChangeMsg{}
			if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else if err := handler.ResizePty(&msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			}
		case "env":
			msg := setenvRequest{}
			if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				ok, payload = handler.Setenv(msg.Name, msg.Value)
			}
		case "exec":
			msg := execMsg{}
			if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				split := strings.Split(msg.Command, " ")
				if len(split) < 1 {
					ok = false
					payload = []byte("empty command")
				}
				ok, payload = handler.Exec(split[0], append(split[1:], args...), config)
			}
		case "exec-args":
			// set up arguments for a future exec call
			msg := stringArrayMsg{}
			if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				args = msg.Args
				payload = []byte{}
			}
		case "exec-config":
			msg := nameValueMsg{}
			if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				config[msg.Name] = msg.Value
				payload = []byte{}
			}
		case "signal":
			msg := signalMsg{}
			if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				handler.Signal(ssh.Signal(msg.Signal))
				payload = []byte{}
			}
		case "kill":
			handler.Kill()
			payload = []byte{}
		case "sync":
			// ensure that filesystem is flushed
			handler.Sync()
			payload = []byte{}
		default:
			ok = false
			err = fmt.Errorf("ssh request type %s is not supported", req.Type)
			log.Println(err.Error())
		}

		// make sure that errors get send back if we failed
		if req.WantReply {
			req.Reply(ok, payload)
		}

		// run any pending work now that a reply has been sent
		pendingFn := handler.GetPendingWork()
		if pendingFn != nil {
			log.Println("Invoking pending work")
			go pendingFn()
			handler.ClearPendingWork()
		}
	}
	log.Println("incoming request channel closed")
}

// this will block until the tether closes the connection
func StartTether(conn net.Conn, privateKey ssh.Signer, handler GlobalContext) error {
	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
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
	config.AddHostKey(privateKey)

	sConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		err := fmt.Errorf("failed to handshake: %s", err)
		log.Println(err.Error())
		return err
	}

	// start the connection manager if extant
	go handler.StartConnectionManager(sConn)

	// Global requests
	go globalMux(reqs, handler)

	log.Println("ready to service ssh requests")
	// Service the incoming session channels
	for sessionchan := range chans {
		// each will need their own channel
		shandler := handler.NewSessionContext()

		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if sessionchan.ChannelType() != "session" {
			err := fmt.Errorf("unknown session type %s", sessionchan.ChannelType())
			log.Println(err.Error())
			sessionchan.Reject(ssh.UnknownChannelType, "unknown channel type")
			return err
		}

		channel, requests, err := sessionchan.Accept()
		if err != nil {
			err := fmt.Errorf("could not accept channel: %s", err)
			log.Println(err.Error())
			return err
		}

		shandler.SetChannel(&channel)
		go channelMux(requests, shandler)
	}

	log.Println("incoming channel channel closed")
	sConn.Close()

	return nil
}
