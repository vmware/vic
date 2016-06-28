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
	"os"
	"runtime/debug"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

var tthr tether.Tether

func init() {
	trace.Logger.Level = log.DebugLevel
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("run time panic: %s : %s", r, debug.Stack())
		}

		reboot()
	}()

	if strings.HasSuffix(os.Args[0], "-debug") {
		extraconfig.DecodeLogLevel = log.DebugLevel
		extraconfig.EncodeLogLevel = log.DebugLevel
		log.SetLevel(log.DebugLevel)
	}

	src, err := extraconfig.GuestInfoSourceWithPrefix("init")
	if err != nil {
		log.Error(err)
		return
	}

	sink, err := extraconfig.GuestInfoSinkWithPrefix("init")
	if err != nil {
		log.Error(err)
		return
	}

	// create the tether and register the attach extension
	tthr = tether.New(src, sink, &operations{})

	err = tthr.Start()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Clean exit from init")
}

// exit cleanly shuts down the system
func halt() {
	log.Infof("Powering off the system")
	if strings.HasSuffix(os.Args[0], "-debug") {
		log.Info("Squashing power off for debug init")
		return
	}

	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}

func reboot() {
	log.Infof("Rebooting the system")
	if strings.HasSuffix(os.Args[0], "-debug") {
		log.Info("Squashing reboot for debug init")
		return
	}

	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}
