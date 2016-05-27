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
	_ "net/http/pprof"
	"os"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/serial"
	"github.com/vmware/vic/pkg/trace"
)

var backchannelMode = os.ModePerm

func backchannel(ctx context.Context) (net.Conn, error) {
	defer trace.End(trace.Begin("establish tether backchannel"))

	log.Info("opening ttyS0 for backchannel")
	f, err := os.OpenFile(pathPrefix+"/ttyS0", os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, backchannelMode)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// set the provided FDs to raw if it's a termial
	// 0 is the uninitialized value for Fd
	if f.Fd() != 0 && terminal.IsTerminal(int(f.Fd())) {
		log.Debug("setting terminal to raw mode")
		s, err := terminal.MakeRaw(int(f.Fd()))
		if err != nil {
			return nil, err
		}

		log.Infof("s = %#v", s)
	}

	log.Infof("creating raw connection from ttyS0 (fd=%d)\n", f.Fd())
	conn, err := serial.NewFileConn(f)

	if err != nil {
		detail := fmt.Sprintf("failed to create raw connection from ttyS0 file handle: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// HACK: currently RawConn dosn't implement timeout so throttle the spinning

	// This needs to tick *faster* than the ticker in connection.go on the
	// portlayer side.  The PL sends the first syn and if this isn't waiting,
	// alignment will take a few rounds (or it may never happen).
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			err := serial.HandshakeServer(ctx, conn)
			if err == nil {
				return conn, nil
			}
		case <-ctx.Done():
			conn.Close()
			ticker.Stop()
			return nil, ctx.Err()
		}
	}
}
