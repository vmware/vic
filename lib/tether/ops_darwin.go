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

package tether

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/trace"
)

type BaseOperations struct{}

// SetHostname sets the system hostname
func (t *BaseOperations) SetHostname(hostname string, aliases ...string) error {
	defer trace.End(trace.Begin("setting hostname to " + hostname))

	return errors.New("not implemented on OSX")
}

// Apply takes the network endpoint configuration and applies it to the system
func (t *BaseOperations) Apply(endpoint *NetworkEndpoint) error {
	defer trace.End(trace.Begin("applying endpoint configuration for " + endpoint.Network.Name))

	return errors.New("not implemented on OSX")
}

// MountLabel performs a mount with the source treated as a disk label
// This assumes that /dev/disk/by-label is being populated, probably by udev
func (t *BaseOperations) MountLabel(ctx context.Context, label, target string) error {
	defer trace.End(trace.Begin(fmt.Sprintf("Mounting %s on %s", label, target)))

	return errors.New("not implemented on OSX")
}

// MountFileSystem performs a mount based on the target path from the source url
// This assumes that the source url is valid and available.
func (t *BaseOperations) MountFileSystem(ctx context.Context, source url.URL, target string, mountOptions string) error {
	defer trace.End(trace.Begin(fmt.Sprintf("Mounting %s on %s", source.String(), target)))

	return errors.New("not implemented on OSX")
}

// ProcessEnv does OS specific checking and munging on the process environment prior to launch
func (t *BaseOperations) ProcessEnv(env []string) []string {
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

func (t *BaseOperations) Fork() error {
	return errors.New("not implemented on OSX")
}

func MkNamedPipe(path string, mode os.FileMode) error {
	return errors.New("not implemented on OSX")
}

func (t *BaseOperations) Setup(_ Config) error {
	return nil
}

func (t *BaseOperations) Cleanup() error {
	return nil
}

// Need to put this here because Windows does not
// support SysProcAttr.Credential
func getUserSysProcAttr(uname string) *syscall.SysProcAttr {
	uinfo, err := user.Lookup(uname)
	if err != nil {
		detail := fmt.Sprintf("Unable to find user %s: %s", uname, err)
		log.Error(detail)
		return nil
	} else {
		u, _ := strconv.Atoi(uinfo.Uid)
		g, _ := strconv.Atoi(uinfo.Gid)
		// Unfortunately lookup GID by name is currently unsupported in Go.
		return &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(u),
				Gid: uint32(g),
			},
			Setsid: true,
		}
	}
}
