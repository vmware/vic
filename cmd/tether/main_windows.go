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
	"runtime/debug"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/tether"
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
	pathPrefix = "com://"

	// use the same logger for trace and other logging
	trace.Logger = log.StandardLogger()
	log.SetLevel(log.DebugLevel)

	// Initiliaze logger with default TextFormatter
	log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})

	// get the windows service logic running so that we can play well in that mode
	runService("VMware Tether", false)

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

	// create the tether and register the attach extension
	tthr = tether.New(src, sink, &operations{})
	tthr.Register("Attach", sshserver)

	err = tthr.Start()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Clean exit from tether")
}

// exit reports completion detail in persistent fashion then cleanly
// shuts down the system
func halt() {
	log.Infof("Powering off the system")
	// TODO: windows fast halt command
}
