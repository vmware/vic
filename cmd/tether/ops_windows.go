// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	winserial "github.com/tarm/serial"

	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/trace"
)

type operations struct {
	tether.BaseOperations

	logging bool
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
		return os.OpenFile(parts[1], os.O_RDWR|os.O_SYNC, 0)
	default:
		return nil, errors.New("unrecognised destination scheme: " + scheme)
	}
}

func (t *operations) Log() (io.Writer, error) {
	defer trace.End(trace.Begin("operations.Log"))

	com := "COM2"

	// redirect logging to the serial log
	log.Infof("opening %s%s for debug log", pathPrefix, com)
	out, err := OpenPort(fmt.Sprintf("%s%s", pathPrefix, com))
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for debug log: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	return out, nil
}

// sessionLogWriter returns a writer that will persist the session output
func (t *operations) SessionLog(session *tether.SessionConfig) (dio.DynamicMultiWriter, dio.DynamicMultiWriter, error) {
	com := "COM3"

	defer trace.End(trace.Begin("configure session log writer"))

	if t.logging {
		detail := "unable to log more than one session concurrently"
		log.Error(detail)
		return nil, nil, errors.New(detail)
	}

	t.logging = true

	// redirect backchannel to the serial connection
	log.Infof("opening %s%s for session logging", pathPrefix, com)
	f, err := OpenPort(fmt.Sprintf("%s%s", pathPrefix, com))
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for session log: %s", err)
		log.Error(detail)
		return nil, nil, errors.New(detail)
	}

	// use multi-writer so it goes to both screen and session log
	return dio.MultiWriter(f, os.Stdout), dio.MultiWriter(f, os.Stderr), nil
}

func (t *operations) Setup(sink tether.Config) error {

	if err := t.BaseOperations.Setup(sink); err != nil {
		return err
	}

	return nil
}

func (t *operations) SetupFirewall(config *tether.ExecutorConfig) error {
	return errors.New("Not implemented on windows")
}
