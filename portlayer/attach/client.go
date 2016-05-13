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
	"errors"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

const (
	attachChannelType = "attach"
)

type SessionInteraction interface {
	// Send specific signal
	Signal(signal ssh.Signal) error
	// Stdout stream
	Stdout() io.Reader
	// Stderr stream
	Stderr() io.Reader
	// Stdin stream
	Stdin() io.WriteCloser
	Close() error

	// Resize the terminal
	Resize(cols, rows, widthpx, heightpx uint32) error
}

type attachSSH struct {
	client *ssh.Client

	channel  ssh.Channel
	requests <-chan *ssh.Request
}

func SSHls(client *ssh.Client) ([]string, error) {
	ok, reply, err := client.SendRequest(requestContainerIDs, true, nil)
	if !ok || err != nil {
		detail := fmt.Sprintf("failed to get container IDs from remote: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	ids := struct {
		Strings []string
	}{}

	err = ssh.Unmarshal(reply, &ids)
	if err != nil {
		detail := fmt.Sprintf("failed to unmarshall ids from remote: %s", err)
		log.Error(detail)
		log.Debugf("raw IDs response: %+v", reply)
		return nil, err
	}

	return ids.Strings, nil
}

// SSHAttach returns a stream connection to the requested session
// The ssh client is assumed to be connected to the Executor hosting the session
func SSHAttach(client *ssh.Client, id string) (SessionInteraction, error) {
	sessionSSH := &attachSSH{
		client: client,
	}

	var err error
	sessionSSH.channel, sessionSSH.requests, err = client.OpenChannel(attachChannelType, []byte(id))
	if err != nil {
		return nil, err
	}

	// we have to handle incoming requests to the client on this channel but we don't support any currently
	go func() {
		for req := range sessionSSH.requests {
			// default, preserving OpenSSH behaviour
			req.Reply(false, nil)
		}
	}()

	return sessionSSH, nil
}

func (t *attachSSH) Signal(signal ssh.Signal) error {
	return errors.New("signal is unimplemented")
}

func (t *attachSSH) Stdout() io.Reader {
	return t.channel
}

func (t *attachSSH) Stderr() io.Reader {
	return t.channel.Stderr()
}

func (t *attachSSH) Stdin() io.WriteCloser {
	return t.channel
}

func (t *attachSSH) Close() error {
	return t.channel.Close()
}

// Resize resizes the terminal.
func (t *attachSSH) Resize(cols, rows, widthpx, heightpx uint32) error {
	msg := WindowChangeMsg{cols, rows, widthpx, heightpx}
	ok, err := t.channel.SendRequest(WindowChangeReq, true, msg.Marshal())
	if err == nil && !ok {
		return fmt.Errorf("unknown error")
	}

	if err != nil {
		return fmt.Errorf("resize error: %s", err)
	}

	return nil
}
