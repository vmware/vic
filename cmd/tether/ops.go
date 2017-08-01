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
	"fmt"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/etcconf"
	"github.com/vmware/vic/lib/tether"
)

// pathPrefix is present to allow the various files referenced by tether to be placed
// in specific directories, primarily for testing.
var pathPrefix string

func (t *operations) Cleanup() error {
	return t.BaseOperations.Cleanup()
}

func (t *operations) Apply(endpoint *tether.NetworkEndpoint) error {
	err := t.BaseOperations.Apply(endpoint)
	if err != nil {
		return err
	}

	bindMountMap := map[string]string{
		hostsPathBindSrc:      etcconf.HostsPath,
		resolvConfPathBindSrc: etcconf.ResolvConfPath,
		hostnameFileBindSrc:   hostnameFile,
	}

	for src, target := range bindMountMap {
		err = bindMount(src, target)
		if err != nil {
			return err
		}
	}

	return nil
}

// HandleSessionExit controls the behaviour on session exit - for the tether if the session exiting
// is the primary session (i.e. SessionID matches ExecutorID) then we exit everything.
func (t *operations) HandleSessionExit(config *tether.ExecutorConfig, session *tether.SessionConfig) func() {
	// if the session that's exiting is the primary session, stop the tether
	return func() {
		if session.ID == config.ID {
			tthr.Stop()
		}
	}
}

func bindMount(src, target string) error {
	// no need to return if unmount fails; it's possible that the target is not mounted previously
	log.Infof("unmounting %s", target)
	if err := tether.Sys.Syscall.Unmount(target, syscall.MNT_DETACH); err != nil {
		log.Errorf("failed to unmount %s: %s", target, err)
	}

	// bind mount src to target
	log.Infof("bind-mounting %s on %s", src, target)
	if err := tether.Sys.Syscall.Mount(src, target, "bind", syscall.MS_BIND, ""); err != nil {
		return fmt.Errorf("faild to mount %s to %s: %s", src, target, err)
	}

	// make sure the file is readable
	// #nosec: Expect file permissions to be 0600 or less
	if err := os.Chmod(target, 0644); err != nil {
		return err
	}

	return nil
}
