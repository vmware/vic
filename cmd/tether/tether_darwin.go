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
	"os"
	"strings"

	"github.com/vmware/vic/pkg/dio"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

// allow us to pick up some of the osops implementations when mocking
// allowing it to be less all or nothing
func init() {
	ops = &osopsOSX{}
	utils = &osopsOSX{}
}

var backchannelMode = os.ModePerm

func (t *osopsOSX) setup() error {
}

func (t *osopsOSX) cleanup() {
}

// sessionLogWriter returns a writer that will persist the session output
func (t *osopsOSX) sessionLogWriter() (dio.DynamicMultiWriter, error) {
	return nil, errors.New("unimplemented on OSX")
}

// processEnvOS does OS specific checking and munging on the process environment prior to launch
func (t *osopsOSX) processEnvOS(env []string) []string {
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

func (t *osopsOSX) establishPty(live *liveSession) error {
	return errors.New("unimplemented on OSX")
}

func (t *osopsOSX) resizePty(pty uintptr, winSize *WindowChangeMsg) error {
	return errors.New("unimplemented on OSX")
}

func (t *osopsOSX) signalProcess(process *os.Process, sig ssh.Signal) error {
	return errors.New("unimplemented on OSX")
}

func (t *osopsOSX) backchannel(ctx context.Context) (net.Conn, error) {
	return nil, errors.New("unimplemented on OSX")
}
