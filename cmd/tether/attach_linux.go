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
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/cmd/tether/msgs"
	"github.com/vmware/vic/pkg/serial"
	"github.com/vmware/vic/pkg/trace"
)

var backchannelMode = os.ModePerm

func backchannel(ctx context.Context, conn *net.Conn) error {
	defer trace.End(trace.Begin("establish tether backchannel"))

	// HACK: currently RawConn dosn't implement timeout so throttle the spinning
	// it does implement the Timeout methods so the intermediary code can be written
	// to support it, but they are stub implementation in rawconn impl.

	// This needs to tick *faster* than the ticker in connection.go on the
	// portlayer side.  The PL sends the first syn and if this isn't waiting,
	// alignment will take a few rounds (or it may never happen).
	ticker := time.NewTicker(10 * time.Millisecond)
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
	defer trace.End(trace.Begin(""))

	t.m.Lock()
	defer t.m.Unlock()

	var err error

	rand.Reader, err = os.Open(pathPrefix + "/urandom")
	if err != nil {
		detail := fmt.Sprintf("failed to open new urandom device: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}

	log.Info("opening ttyS0 for backchannel")
	f, err := os.OpenFile(pathPrefix+"/ttyS0", os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, backchannelMode)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for backchannel: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}

	// set the provided FDs to raw if it's a termial
	// 0 is the uninitialized value for Fd
	if f.Fd() != 0 && terminal.IsTerminal(int(f.Fd())) {
		log.Debug("setting terminal to raw mode")
		s, err := terminal.MakeRaw(int(f.Fd()))
		if err != nil {
			return err
		}

		log.Infof("s = %#v", s)
	}

	log.Infof("creating raw connection from ttyS0 (fd=%d)\n", f.Fd())
	var conn net.Conn
	conn, err = serial.NewFileConn(f)
	if err != nil {
		detail := fmt.Sprintf("failed to create raw connection from ttyS0 file handle: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}

	t.conn = &conn
	return nil
}

func resizePty(pty uintptr, winSize *msgs.WindowChangeMsg) error {
	defer trace.End(trace.Begin(""))

	ws := &winsize{uint16(winSize.Rows), uint16(winSize.Columns), uint16(winSize.WidthPx), uint16(winSize.HeightPx)}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		pty,
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func signalProcess(process *os.Process, sig ssh.Signal) error {
	signal := msgs.Signals[sig]
	defer trace.End(trace.Begin(fmt.Sprintf("signal process %d: %s", process.Pid, sig)))

	s := syscall.Signal(signal)
	return process.Signal(s)
}
