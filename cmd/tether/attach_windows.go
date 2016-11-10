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

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/cmd/tether/msgs"
	"github.com/vmware/vic/pkg/trace"
)

const (
	com = "COM1"
)

func rawConnectionFromSerial() (net.Conn, error) {
	var err error

	// redirect backchannel to the serial connection
	log.Infof("opening %s%s for backchannel", pathPrefix, com)
	// TODO: set read timeout on port during open
	_, err = OpenPort(fmt.Sprintf("%s%s", pathPrefix, com))
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	log.Errorf("creating raw connection from %s\n", com)

	// TODO: sort out the named port impl so that we can transparently switch from that to/from
	// regular files for testing
	// t.conn, err := serial.NewTypedConn(f, "file")
	return nil, nil
}

func (t *attachServerSSH) Start() error {
	defer trace.End(trace.Begin(""))

	return nil
}

func resizePty(pty uintptr, winSize *msgs.WindowChangeMsg) error {
	return errors.New("not supported on windows")
}
