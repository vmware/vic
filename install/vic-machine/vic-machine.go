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
	"flag"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
)

// InstallOptions holds the CLI arguments for reference
type InstallOptions struct {
	logfile string
	debug   bool
	stdout  bool

	from     string
	target   string
	username string
	password string
}

// Default values for flag arguments
const (
	DefaultLogfile = "vic-machine.log"
	// CEIPGatewayBase is where we access the phone home reporting
	// The full string looks like:
	//  "https://vcsa.vmware.com/ph-stg/api/hyper/send?_v=1.0&_c=vic.1_0&_i="+vc.uuid
	// TODO: this is currently the staging server and should be switched over once the
	// schema has stabilized
	CEIPGatewayBase = "https://vcsa.vmware.com/ph-stg/api/hyper/send"

	// FromUsage describes how to use the from command arguments
	FromUsage = "URI of configuration file"

	// TargetUsage describes the possible target variations
	TargetUsage = `Target is a URI of the forms:
                    vmomi://vsphere-host/moid=<moid> where ACTION=[update|delete]
                    vmomi://vsphere-host/compute/resource/path/to/vch`
)

var (
	options InstallOptions
)

func init() {
	flag.StringVar(&options.from, "from", "", "URI of configuration file")
	flag.StringVar(&options.target, "target", "", "URI of the target compute resource")
	flag.StringVar(&options.target, "username", "", "Username for the target")
	flag.StringVar(&options.target, "passwd", "", "Password for the target")

	flag.StringVar(&options.logfile, "logfile", DefaultLogfile, "Path of the installer log file")

	flag.Parse()
}

func main() {
	// Open log file
	f, err := os.OpenFile(options.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening logfile %s: %v", options.logfile, err)
	}
	defer f.Close()

	if options.debug {
		log.SetLevel(log.DebugLevel)
	}

	// Initiliaze logger with default TextFormatter
	log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})

	// SetOutput to log file and/or stdout
	log.SetOutput(f)
	if options.stdout {
		log.SetOutput(io.MultiWriter(os.Stdout, f))
	}

}
