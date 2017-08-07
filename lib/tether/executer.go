// Copyright 2017 VMware, Inc. All Rights Reserved.
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
	"fmt"
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
)

// Executor is a tether extension that wraps command
type Executor struct {
	p    *os.Process
	exec func() (*os.Process, error)
}

// NewExecutor returns a tether.Extension that wraps the executor service
func NewExecutor() *Executor {
	return &Executor{}
}

// NewHaveged returns a tether.Extension that wraps haveged
func NewHaveged() *Executor {
	return &Executor{
		exec: func() (*os.Process, error) {
			args := []string{"/.tether/lib/ld-linux-x86-64.so.2", "--library-path", "/.tether/lib", "/.tether/haveged", "-w", "1024", "-v", "1", "-F"}
			// #nosec: Subprocess launching with variable
			cmd := exec.Command(args[0], args[1:]...)

			logrus.Infof("Starting haveged with args: %q", args)
			if err := cmd.Start(); err != nil {
				logrus.Errorf("Starting haveged failed with %q", err.Error())
				return nil, err
			}
			return cmd.Process, nil
		},
	}
}

// Start implementation of the tether.Extension interface
func (e *Executor) Start() error {
	logrus.Infof("Starting haveged")

	var err error
	e.p, err = e.exec()
	return err
}

// Stop implementation of the tether.Extension interface
func (e *Executor) Stop() error {
	logrus.Infof("Stopping haveged")

	if e.p != nil {
		return e.p.Kill()
	}
	return fmt.Errorf("haveged process is missing")
}

// Reload implementation of the tether.Extension interface
func (e *Executor) Reload(config *ExecutorConfig) error {
	// haveged doesn't support reloading
	return nil
}
