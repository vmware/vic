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
	"net"
	_ "net/http/pprof"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"

	"github.com/vmware/vic/lib/portlayer/attach"
	"github.com/vmware/vic/pkg/serial"
)

func backchannel(ctx context.Context, conn *net.Conn) error {

	// HACK: currently RawConn dosn't implement timeout so throttle the spinning
	ticker := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			err := serial.HandshakeServer(ctx, *conn)
			if err == nil {
				return nil
			}
		case <-ctx.Done():
			(*conn).Close()
			ticker.Stop()
			return ctx.Err()
		}
	}
}

func (t *attachServerSSH) Start() error {
	return errors.New("not supported on OSX")
}

func resizePty(pty uintptr, winSize *attach.WindowChangeMsg) error {
	return errors.New("not supported on OSX")
}

func signalProcess(process *os.Process, sig ssh.Signal) error {
	return errors.New("not supported on OSX")
}
