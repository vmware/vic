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
	"os"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/cmd/tether/serial"
	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/trace"
	"golang.org/x/net/context"
)

type osopsMock struct {
	// allow tests to tell when the struct has been updated
	updated chan bool

	// the hostname of the system
	hostname string
	// the ip configuration for mac indexed interfaces
	ips map[string]net.IPNet
	// filesystem mounts, indexed by disk label
	mounts map[string]string
}

// SetHostname sets both the kernel hostname and /etc/hostname to the specified string
func (t *osopsMock) SetHostname(hostname string) error {
	defer trace.End(trace.Begin("mocking hostname to " + hostname))

	// TODO: we could mock at a much finer granularity, only extracting the syscall
	// that would exercise the file modification paths, however it's much less generalizable
	t.hostname = hostname

	t.updated <- true
	return nil
}

// Apply takes the network endpoint configuration and applies it to the system
func (t *osopsMock) Apply(endpoint *metadata.NetworkEndpoint) error {
	defer trace.End(trace.Begin("mocking endpoint configuration for " + endpoint.Network.Name))

	t.updated <- true
	return errors.New("Apply test not implemented")
}

// MountLabel performs a mount with the source treated as a disk label
// This assumes that /dev/disk/by-label is being populated, probably by udev
func (t *osopsMock) MountLabel(label, target string, ctx context.Context) error {
	defer trace.End(trace.Begin(fmt.Sprintf("mocking mounting %s on %s", label, target)))

	if t.mounts == nil {
		t.mounts = make(map[string]string)
	}

	t.mounts[label] = target

	t.updated <- true
	return nil
}

// Fork triggers vmfork and handles the necessary pre/post OS level operations
func (t *osopsMock) Fork(config *metadata.ExecutorConfig) error {
	defer trace.End(trace.Begin("mocking fork"))

	t.updated <- true
	return errors.New("Fork test not implemented")
}

func (t *osopsMock) backchannel(ctx context.Context) (net.Conn, error) {
	log.Info("opening ttyS0 pipe pair for backchannel")
	c, err := os.OpenFile(pathPrefix+"/ttyS0c", os.O_WRONLY|syscall.O_NOCTTY, 0777)
	if err != nil {
		detail := fmt.Sprintf("failed to open cpipe for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	s, err := os.OpenFile(pathPrefix+"/ttyS0s", os.O_RDONLY|syscall.O_NOCTTY, 0777)
	if err != nil {
		detail := fmt.Sprintf("failed to open spipe for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	log.Infof("creating raw connection from ttyS0 pipe pair (c=%d, s=%d)\n", c.Fd(), s.Fd())
	conn, err := serial.NewHalfDuplixFileConn(s, c, pathPrefix+"/ttyS0", "file")

	if err != nil {
		detail := fmt.Sprintf("failed to create raw connection from ttyS0 pipe pair: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// still run handshake over it to test that
	ticker := time.NewTicker(1000 * time.Millisecond)
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
