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

package communication

import (
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"

	"github.com/vmware/vic/cmd/tether/msgs"
	"github.com/vmware/vic/pkg/trace"
)

const (
	attachChannelType = "attach"
)

// SessionInteractor defines the interaction interface
type SessionInteractor interface {
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

	CloseStdin() error

	Ping() error
	Unblock() error
}

// interaction implements SessionInteractor using SSH
type interaction struct {
	channel ssh.Channel
}

// ContainerIDs asks the ids of the containers on the other hand and return them to the caller
func ContainerIDs(conn ssh.Conn) ([]string, error) {
	defer trace.End(trace.Begin(""))

	ok, reply, err := conn.SendRequest(msgs.ContainersReq, true, nil)
	if !ok || err != nil {
		return nil, fmt.Errorf("failed to get container IDs from remote: %s", err)
	}

	ids := msgs.ContainersMsg{}
	if err = ids.Unmarshal(reply); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ids from remote: %s", err)
	}

	return ids.IDs, nil
}

// NewSSHInteraction returns a stream connection to the requested session
// The ssh conn is assumed to be connected to the Executor hosting the session
func NewSSHInteraction(conn ssh.Conn, id string) (SessionInteractor, error) {
	defer trace.End(trace.Begin(id))

	channel, _, err := conn.OpenChannel(attachChannelType, []byte(id))
	if err != nil {
		return nil, err
	}

	return &interaction{channel: channel}, nil
}

func (t *interaction) Signal(signal ssh.Signal) error {
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

func (t *interaction) CloseStdin() error {
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

func (t *interaction) Stdout() io.Reader {
	defer trace.End(trace.Begin(""))

	return t.channel
}

func (t *interaction) Stderr() io.Reader {
	defer trace.End(trace.Begin(""))

	return t.channel.Stderr()
}

func (t *interaction) Stdin() io.WriteCloser {
	defer trace.End(trace.Begin(""))

	return t.channel
}

func (t *interaction) Close() error {
	defer trace.End(trace.Begin(""))

	return t.channel.Close()
}

// Resize resizes the terminal.
func (t *interaction) Resize(cols, rows, widthpx, heightpx uint32) error {
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

// Ping checks the liveleness of the connection
func (t *interaction) Ping() error {
	defer trace.End(trace.Begin(""))

	ok, err := t.channel.SendRequest(msgs.PingReq, true, []byte(msgs.PingMsg))
	if !ok || err != nil {
		return fmt.Errorf("failed to ping the other side: %s", err)
	}

	return nil
}

// Unblock sends an unblock msg
func (t *interaction) Unblock() error {
	defer trace.End(trace.Begin(""))

	ok, err := t.channel.SendRequest(msgs.UnblockReq, true, []byte(msgs.UnblockMsg))
	if !ok || err != nil {
		return fmt.Errorf("failed to unblock the other side: %s", err)
	}

	return nil
}
