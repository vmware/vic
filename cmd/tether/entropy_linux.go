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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/tether"
)

// TODO(morris-jason): Not sure this is the best file path. Open to suggestions.
const entropyConfigPath = "/.tether/opt/config/entropy.txt"

// Entropy is a tether extension that wraps command
type Entropy struct {
	p    *os.Process
	exec func() (*os.Process, error)
}

// NewEntropy returns a tether.Extension that wraps entropy
func NewEntropy() *Entropy {
	return &Entropy{
		exec: func() (*os.Process, error) {
			entropyConfig, err := ioutil.ReadFile(entropyConfigPath)
			if err != nil {
				log.Errorf("Cannot read entropy configuration %q", err.Error())
				return nil, err
			}
			args := strings.Split(string(entropyConfig), " ")
			// #nosec: Subprocess launching with variable
			cmd := exec.Command(args[0], args[1:]...)
			log.Infof("Starting entropy daemon")
			if err := cmd.Start(); err != nil {
				log.Errorf("Starting entropy failed with %q", err.Error())
				return nil, err
			}
			return cmd.Process, nil
		},
	}
}

// Start implementation of the tether.Extension interface
func (e *Entropy) Start(system tether.System) error {
	log.Infof("Starting entropy")

	var err error
	e.p, err = e.exec()
	return err
}

// Stop implementation of the tether.Extension interface
func (e *Entropy) Stop() error {
	log.Infof("Stopping entropy")

	if e.p != nil {
		return e.p.Kill()
	}
	return fmt.Errorf("Entropy process is missing")
}

// Reload implementation of the tether.Extension interface
func (e *Entropy) Reload(config *tether.ExecutorConfig) error {
	// entropy doesn't support reloading
	return nil
}
