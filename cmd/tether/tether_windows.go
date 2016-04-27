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
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	winserial "github.com/tarm/serial"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/serial"
)

// allow us to pick up some of the osops implementations when mocking
// allowing it to be less all or nothing
func init() {
	ops = &osopsWin{}
	utils = &osopsWin{}
}

var backchannelMode = os.ModePerm

type NamedPort struct {
	*winserial.Port

	config winserial.Config
	fd     uintptr
}

func (p *NamedPort) Name() string {
	return p.config.Name
}

func (p *NamedPort) Fd() uintptr {
	return p.fd
}

// Writer interface for the named port
func (p *NamedPort) Write(b []byte) (int, error) {
	// TODO: glue in writer code
	return 0, errors.New("Write not yet implemented")
}

// Reader interface for the named port
func (p *NamedPort) Read(b []byte) (int, error) {
	// TODO: glue in reader code
	return 0, errors.New("Read not yet implemented")
}

func (p *NamedPort) Close() error {
	return errors.New("Read not yet implemented")
}

// OpenPort does exactly that
// TODO: this needs to be renamed updated to open a regular file if that's
// what is specified - use a URI scheme to disambiguate.
// this will let us test on windows without actually needing com ports.
func OpenPort(name string) (io.ReadWriteCloser, error) {
	parts := strings.Split(name, "://")
	if len(parts) != 2 {
		return nil, errors.New("expected name to have a scheme://<name> construction")
	}

	scheme := parts[0]
	switch scheme {
	case "com":
		cfg := &winserial.Config{Name: parts[1], Baud: 115200}
		port, err := winserial.OpenPort(cfg)
		if err != nil {
			return nil, err
		}

		// ensure we don't have significant obsolete data built up
		port.Flush()
		return &NamedPort{Port: port, config: *cfg, fd: 0}, nil
	case "file":
		return os.OpenFile(parts[1], os.O_RDWR|os.O_SYNC, 0777)
	default:
		return nil, errors.New("unrecognised destination scheme: " + scheme)
	}
}

func childReaper() {
	// TODO: windows child process notifications
}

func (t *osopsWin) setup() error {
	com := "COM2"

	// redirect logging to the serial log
	log.Infof("opening %s%s for debug log", pathPrefix, com)
	out, err := OpenPort(fmt.Sprintf("%s%s", pathPrefix, com))
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for debug log: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}
	log.SetOutput(out)

	// TODO: enabled for initial dev debugging only
	go func() {
		log.Info(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	go childReaper()

	return nil
}

func (t *osopsWin) cleanup() {
}

func (t *osopsWin) backchannel(ctx context.Context) (net.Conn, error) {
	com := "COM1"

	// redirect backchannel to the serial connection
	log.Infof("opening %s%s for backchannel", pathPrefix, com)
	// TODO: set read timeout on port during open
	_, err := OpenPort(fmt.Sprintf("%s%s", pathPrefix, com))
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	log.Errorf("creating raw connection from %s\n", com)

	// TODO: sort out the named port impl so that we can transparently switch from that to/from
	// regular files for testing
	// conn, err := serial.NewTypedConn(f, "file")
	var conn net.Conn

	if err != nil {
		detail := fmt.Sprintf("failed to create raw connection from %s file handle: %s", com, err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// HACK: currently RawConn dosn't implement timeout so throttle the spinning
	ticker := time.NewTicker(50 * time.Millisecond)
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

// sessionLogWriter returns a writer that will persist the session output
func (t *osopsWin) sessionLogWriter() (dio.DynamicMultiWriter, error) {
	com := "COM3"

	// redirect backchannel to the serial connection
	log.Infof("opening %s%s for session logging", pathPrefix, com)
	f, err := OpenPort(fmt.Sprintf("%s%s", pathPrefix, com))
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for session log: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// use multi-writer so it goes to both screen and session log
	return dio.MultiWriter(f, os.Stdout), nil
}

// processEnvOS does OS specific checking and munging on the process environment prior to launch
func (t *osopsWin) processEnvOS(env []string) []string {
	// TODO: figure out how we're going to specify user and pass all the settings along
	// in the meantime, hardcode HOME to /root
	homeIndex := -1
	for i, tuple := range env {
		if strings.HasPrefix(tuple, "HOME=") {
			homeIndex = i
			break
		}
	}
	if homeIndex == -1 {
		return append(env, "HOME=/root")
	}

	return env
}

func (t *osopsWin) signalProcess(process *os.Process, sig ssh.Signal) error {
	return errors.New("unimplemented on windows")
}

func (t *osopsWin) establishPty(live *liveSession) error {
	return errors.New("unimplemented on windows")
}

func (t *osopsWin) resizePty(pty uintptr, winSize *WindowChangeMsg) error {
	return errors.New("unimplemented on windows")
}
