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
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"runtime/trace"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/reference"
	"github.com/pkg/profile"

	"github.com/vmware/vic/lib/imagec"
	"github.com/vmware/vic/pkg/version"
)

var (
	imageCOptions = ImageCOptions{}

	// https://raw.githubusercontent.com/docker/docker/master/distribution/pull_v2.go
	sf = streamformatter.NewJSONStreamFormatter()
)

const (
	PullImage = "pull"
	PushImage = "push"
)

// ImageCOptions wraps the cli arguments
type ImageCOptions struct {
	imageName string

	options imagec.Options

	logfile string

	stdout bool
	debug  bool

	profiling string
	tracing   bool

	operation string
}

func init() {

	flag.StringVar(&imageCOptions.imageName, "reference", "", "Name of the reference")

	flag.StringVar(&imageCOptions.options.Destination, "destination", imagec.DefaultDestination, "Destination directory")

	flag.StringVar(&imageCOptions.options.Host, "host", imagec.DefaultPortLayerHost, "Host that runs portlayer API (FQDN:port format)")

	flag.StringVar(&imageCOptions.logfile, "logfile", imagec.DefaultLogfile, "Path of the imagec log file")

	flag.StringVar(&imageCOptions.options.Username, "username", "", "Username")
	flag.StringVar(&imageCOptions.options.Password, "password", "", "Password")

	flag.DurationVar(&imageCOptions.options.Timeout, "timeout", imagec.DefaultHTTPTimeout, "HTTP timeout")

	flag.BoolVar(&imageCOptions.stdout, "stdout", false, "Enable writing to stdout")
	flag.BoolVar(&imageCOptions.debug, "debug", false, "Show debug logging")
	flag.BoolVar(&imageCOptions.options.InsecureSkipVerify, "insecure-skip-verify", false, "Don't verify certificates when fetching images")
	flag.BoolVar(&imageCOptions.options.InsecureAllowHTTP, "insecure-allow-http", false, "Uses unencrypted connections when fetching images")
	flag.BoolVar(&imageCOptions.options.Standalone, "standalone", false, "Disable port-layer integration")

	flag.StringVar(&imageCOptions.profiling, "profile.mode", "", "Enable profiling mode, one of [cpu, mem, block]")
	flag.BoolVar(&imageCOptions.tracing, "tracing", false, "Enable runtime tracing")

	flag.StringVar(&imageCOptions.operation, "operation", "pull", "Pull/push an image")

	flag.StringVar(&imageCOptions.options.Registry, "registry", imagec.DefaultDockerURL, "Registry to pull/push images (default: registry-1.docker.io)")

	flag.Parse()

	var err error
	if imageCOptions.options.Reference, err = reference.ParseNamed(imageCOptions.imageName); err != nil {
		panic(err)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, string(sf.FormatError(fmt.Errorf("%s : %s", r, debug.Stack()))))
		}
	}()

	if version.Show() {
		fmt.Fprintf(os.Stdout, "%s\n", version.String())
		return
	}

	// Enable profiling if mode is set
	switch imageCOptions.profiling {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.Quiet).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.Quiet).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile, profile.ProfilePath("."), profile.Quiet).Stop()
	default:
		// do nothing
	}

	// Register our custom Error hook
	log.AddHook(NewErrorHook(os.Stderr))

	// Enable runtime tracing if tracing is true
	if imageCOptions.tracing {
		tracing, err := os.Create(time.Now().Format("2006-01-02T150405.pprof"))
		if err != nil {
			log.Fatalf("Failed to create tracing logfile: %s", err)
		}
		defer tracing.Close()

		if err := trace.Start(tracing); err != nil {
			log.Fatalf("Failed to start tracing: %s", err)
		}
		defer trace.Stop()
	}

	// Open the log file
	// #nosec: Expect file permissions to be 0600 or less
	f, err := os.OpenFile(imageCOptions.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open the logfile %s: %s", imageCOptions.logfile, err)
	}
	defer f.Close()

	// Initiliaze logger with default TextFormatter
	log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})

	// Set the log level
	if imageCOptions.debug {
		log.SetLevel(log.DebugLevel)
	}

	// SetOutput to log file and/or stdout
	log.SetOutput(f)
	if imageCOptions.stdout {
		log.SetOutput(io.MultiWriter(os.Stdout, f))
	}

	if imageCOptions.operation == PullImage {
		options := imageCOptions.options

		options.Outstream = os.Stdout

		ic := imagec.NewImageC(options, streamformatter.NewJSONStreamFormatter())
		if err := ic.PullImage(); err != nil {
			log.Fatalf("Pulling image failed due to %s\n", err)
		}
	} else if imageCOptions.operation == PushImage {
		log.Errorf("The operation '%s' is not implemented\n", PushImage)
	} else {
		log.Errorf("The operation '%s' is not valid\n", imageCOptions.operation)
	}
}
