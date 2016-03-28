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
	_ "net/http/pprof"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/metadata"
)

func main() {
	// get the windows service logic running so that we can play well in that mode
	runService("VMware Tether", false)

	// where to look for the various devices and files related to tether
	pathPrefix = "com://"

	err := run(metadata.New(), os.Args[1])

	if err != nil {
		log.Error(err)
	} else {
		log.Info("Clean exit from tether")
	}

	halt()
}

// exit reports completion detail in persistent fashion then cleanly
// shuts down the system
func halt() {
	log.Infof("Powering off the system")
	// TODO: windows fast halt command
}
