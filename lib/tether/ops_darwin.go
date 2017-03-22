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

// MountTarget performs a mount based on the target path from the source url
// This assumes that the source url is valid and available.
func (t *BaseOperations) MountTarget(ctx context.Context, source url.URL, target string, mountOptions string) error {
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
func getUserSysProcAttr(uid, gid string) (*syscall.SysProcAttr, error) {
	if len(uid) == 0 && len(gid) == 0 {
		log.Debugf("no user id or group id specified")
		return nil, nil
	}

	var suid, sgid int
	var suidErr, sgidErr error

	var fuid, fgid int
	var fuidErr, fgidErr error

	var uinfo *user.User
	var ginfo *user.Group

	if len(uid) > 0 {
		suid, suidErr = strconv.Atoi(uid)
		// lookup username
		uinfo, fuidErr = user.Lookup(uid)
		if fuidErr == nil {
			log.Debugf("User %s is found", uid)
			fuid, _ = strconv.Atoi(uinfo.Uid)
			fgid, _ = strconv.Atoi(uinfo.Gid)
		}
	}
	if len(gid) > 0 {
		sgid, sgidErr = strconv.Atoi(gid)

		// lookup groupname
		ginfo, fgidErr = user.LookupGroup(gid)
		// if found groupname, override user group
		if fgidErr == nil {
			log.Debugf("Group %s is found", gid)
			fgid, _ = strconv.Atoi(ginfo.Gid)
		}
	}

	// lookup user failed
	if fuidErr != nil {
		if suidErr != nil {
			// failed to loopup username, and user is not number
			detail := fmt.Sprintf("unable to find user %s: %s", uid, fuidErr)
			return nil, errors.New(detail)
		}
		// user set user id must be inside valid uid range.
		if suid < minID || suid > maxID {
			detail := fmt.Sprintf("user id %s is invalid", uid)
			return nil, errors.New(detail)
		}
		fuid = suid
	}

	// lookup group failed
	if fgidErr != nil {
		if sgidErr != nil {
			// failed to loopup groupname, and user is not number
			detail := fmt.Sprintf("unable to find group %s: %s", gid, fgidErr)
			return nil, errors.New(detail)
		}
		// user set group id must be inside valid uid range.
		if sgid < minID || sgid > maxID {
			detail := fmt.Sprintf("group id %s is invalid", gid)
			return nil, errors.New(detail)
		}
		fgid = sgid
	}
	log.Debugf("set user to %s:%s", fuid, fgid)
	return &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(fuid),
			Gid: uint32(fgid),
		},
		Setsid: true,
	}, nil
}
