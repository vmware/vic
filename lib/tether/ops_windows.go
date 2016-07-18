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
	"fmt"
	"os"

	"golang.org/x/net/context"

	"github.com/vmware/vic/pkg/trace"
)

const pciDevPath = ""

type BaseOperations struct{}

// SetHostname sets the system hostname
func (t *BaseOperations) SetHostname(hostname string, aliases ...string) error {
	defer trace.End(trace.Begin("setting hostname to " + hostname))

	return errors.New("not implemented on windows")
}

// Apply takes the network endpoint configuration and applies it to the system
func (t *BaseOperations) Apply(endpoint *NetworkEndpoint) error {
	defer trace.End(trace.Begin("applying endpoint configuration for " + endpoint.Network.Name))

	return errors.New("not implemented on windows")
}

// MountLabel performs a mount with the source treated as a disk label
// This assumes that /dev/disk/by-label is being populated, probably by udev
func (t *BaseOperations) MountLabel(label, target string, ctx context.Context) error {
	defer trace.End(trace.Begin(fmt.Sprintf("Mounting %s on %s", label, target)))

	return errors.New("not implemented on windows")
}

// processEnvOS does OS specific checking and munging on the process environment prior to launch
func (t *BaseOperations) ProcessEnv(env []string) []string {
	return env
}

// Fork triggers a vmfork, address the pre and post-fork operations necessary at an OS level
func (t *BaseOperations) Fork() error {
	return errors.New("not implemented on windows")
}

func MkNamedPipe(path string, mode os.FileMode) error {
	return errors.New("not implemented on windows")
}

func (t *BaseOperations) Setup(_ Config) error {
	return nil
}

func (t *BaseOperations) Cleanup() error {
	return nil
}
