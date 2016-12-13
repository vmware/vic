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

package tether

import (
	"errors"
	"strings"

	"github.com/vmware/vic/pkg/trace"
)

const (
	pidFilePath = "var/run"
)

// Mkdev will hopefully get rolled into go.sys at some point
func Mkdev(majorNumber int, minorNumber int) int {
	return (majorNumber << 8) | (minorNumber & 0xff) | ((minorNumber & 0xfff00) << 12)
}

// childReaper is used to handle events from child processes, including child exit.
// If running as pid=1 then this means it handles zombie process reaping for orphaned children
// as well as direct child processes.
func (t *tether) childReaper() error {
	return errors.New("Child reaping unimplemented on OSX")
}

func (t *tether) stopReaper() {
	defer trace.End(trace.Begin("Shutting down child reaping"))
}

func (t *tether) triggerReaper() error {
	return errors.New("Child reaping unimplemented on OSX")
}

// processEnvOS does OS specific checking and munging on the process environment prior to launch
func (t *tether) processEnvOS(env []string) []string {
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

// lookPath searches for an executable binary named file in the directories
// specified by the path argument.
// This is a direct modification of the unix os/exec core library impl
func lookPath(file string, env []string, dir string) (string, error) {
	return "", errors.New("unimplemented on OSX")
}

func establishPty(session *SessionConfig) error {
	return errors.New("unimplemented on OSX")
}

func establishNonPty(session *SessionConfig) error {
	return errors.New("unimplemented on OSX")
}
