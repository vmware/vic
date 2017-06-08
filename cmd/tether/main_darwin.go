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

package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/tether"
	viclog "github.com/vmware/vic/pkg/log"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

var tthr tether.Tether

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("run time panic: %s : %s", r, debug.Stack())
		}
		halt()
	}()

	// where to look for the various devices and files related to tether
	pathPrefix = "/.tether"

	if strings.HasSuffix(os.Args[0], "-debug") {
		extraconfig.DecodeLogLevel = log.DebugLevel
		extraconfig.EncodeLogLevel = log.DebugLevel
	}
	trace.Logger.Level = log.DebugLevel
	log.SetLevel(log.DebugLevel)

	// Initiliaze logger with default TextFormatter
	log.SetFormatter(viclog.NewTextFormatter())

	// TODO: hard code executor initialization status reporting via guestinfo here
	err := createDevices()
	if err != nil {
		log.Error(err)
		// return gives us good behaviour in the case of "-debug" binary
		return
	}

	sshserver := NewAttachServerSSH()
	src, err := extraconfig.GuestInfoSource()
	if err != nil {
		log.Error(err)
		return
	}

	sink, err := extraconfig.GuestInfoSink()
	if err != nil {
		log.Error(err)
		return
	}

	// create the tether
	tthr = tether.New(src, sink, &operations{})

	// register the attach extension
	tthr.Register("Attach", sshserver)

	err = tthr.Start()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Clean exit from tether")
}

func createDevices() error {
	// #nosec: Expect directory permissions to be 0700 or less
	err := os.MkdirAll(pathPrefix, 644)
	if err != nil {
		log.Warnf("Failed to ensure presence of tether device directory: %s", err)
	}

	// create serial devices
	for i := 0; i < 3; i++ {
		path := fmt.Sprintf("%s/ttyS%d", pathPrefix, i)
		minor := 64 + i
		err = syscall.Mknod(path, syscall.S_IFCHR|uint32(os.FileMode(0660)), tether.Mkdev(4, minor))
		if err != nil {
			return fmt.Errorf("failed to create %s for com%d: %s", path, i+1, err)
		}
	}

	// make an access to urandom
	path := fmt.Sprintf("%s/urandom", pathPrefix)
	err = syscall.Mknod(path, syscall.S_IFCHR|uint32(os.FileMode(0444)), tether.Mkdev(1, 9))
	if err != nil {
		return fmt.Errorf("failed to create urandom access %s: %s", path, err)
	}

	return nil
}

// exit cleanly shuts down the system
func halt() {
	log.Infof("Powering off the system")
	// if strings.HasSuffix(os.Args[0], "-debug") {
	log.Info("Squashing power off for WIP tether")
	return
	// }
}
