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
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/serial"
)

func (t *operations) backchannel(ctx context.Context) (net.Conn, error) {
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
