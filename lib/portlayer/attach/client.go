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
	"sync/atomic"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/vmware/vic/cmd/tether/msgs"
	"github.com/vmware/vic/pkg/trace"
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
	CloseWrite() error
	IsClosed() bool

	// Resize the terminal
	Resize(cols, rows, widthpx, heightpx uint32) error

	CloseStdin() error
}

type attachSSH struct {
	client *ssh.Client

	channel  ssh.Channel
	requests <-chan *ssh.Request
	closed   uint32
}

func SSHls(client *ssh.Client) ([]string, error) {
	defer trace.End(trace.Begin(""))

	ok, reply, err := client.SendRequest(msgs.ContainersReq, true, nil)
	if !ok || err != nil {
		return nil, fmt.Errorf("failed to get container IDs from remote: %s", err)
	}

	ids := msgs.ContainersMsg{}

	if err = ids.Unmarshal(reply); err != nil {
		log.Debugf("raw IDs response: %+v", reply)
		return nil, fmt.Errorf("failed to unmarshal ids from remote: %s", err)
	}

	return ids.IDs, nil
}

// SSHAttach returns a stream connection to the requested session
// The ssh client is assumed to be connected to the Executor hosting the session
func SSHAttach(client *ssh.Client, id string) (SessionInteraction, error) {
	defer trace.End(trace.Begin(""))

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
			// default, preserving OpenSSH behavior
			req.Reply(false, nil)
		}
	}()

	return sessionSSH, nil
}

func (t *attachSSH) Signal(signal ssh.Signal) error {
	defer trace.End(trace.Begin(""))

	msg := msgs.SignalMsg{Signal: signal}
	ok, err := t.channel.SendRequest(msgs.SignalReq, true, msg.Marshal())
	if err == nil && !ok {
		return fmt.Errorf("unknown error")
	}

	if err != nil {
		return fmt.Errorf("signal error: %s", err)
	}

	return nil
}

func (t *attachSSH) CloseStdin() error {
	defer trace.End(trace.Begin(""))

	ok, err := t.channel.SendRequest(msgs.CloseStdinReq, true, nil)
	if err == nil && !ok {
		return fmt.Errorf("unknown error closing stdin")
	}

	if err != nil {
		return fmt.Errorf("close stdin error: %s", err)
	}

	return nil
}

func (t *attachSSH) Stdout() io.Reader {
	defer trace.End(trace.Begin(""))

	return t.channel
}

func (t *attachSSH) Stderr() io.Reader {
	defer trace.End(trace.Begin(""))

	return t.channel.Stderr()
}

func (t *attachSSH) Stdin() io.WriteCloser {
	defer trace.End(trace.Begin(""))

	return t.channel
}

func (t *attachSSH) Close() error {
	defer trace.End(trace.Begin(""))

	if !atomic.CompareAndSwapUint32(&t.closed, 0, 1) {
		return nil // another routine beat us to it
	}

	err := t.channel.Close()

	return err
}

func (t *attachSSH) CloseWrite() error {
	defer trace.End(trace.Begin(""))

	return t.channel.CloseWrite()
}

func (t *attachSSH) IsClosed() bool {
	return atomic.LoadUint32(&t.closed) == 1
}

// Resize resizes the terminal.
func (t *attachSSH) Resize(cols, rows, widthpx, heightpx uint32) error {
	defer trace.End(trace.Begin(""))

	msg := msgs.WindowChangeMsg{
		Columns:  cols,
		Rows:     rows,
		WidthPx:  widthpx,
		HeightPx: heightpx,
	}
	ok, err := t.channel.SendRequest(msgs.WindowChangeReq, true, msg.Marshal())
	if err == nil && !ok {
		return fmt.Errorf("unknown error")
	}

	if err != nil {
		return fmt.Errorf("resize error: %s", err)
	}
	return nil
}
