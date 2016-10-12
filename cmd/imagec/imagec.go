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
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"runtime/trace"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/pkg/streamformatter"

	"github.com/vmware/vic/lib/imagec"
	"github.com/vmware/vic/pkg/i18n"
	"github.com/vmware/vic/pkg/version"

	"github.com/pkg/profile"
)

var (
	options = imagec.Options{
		Outstream: os.Stdout,
	}

	// https://raw.githubusercontent.com/docker/docker/master/distribution/pull_v2.go
	sf = streamformatter.NewJSONStreamFormatter()
)

func init() {
	// TODO: get language from host OS
	lang := i18n.DefaultLang
	languageFile := fmt.Sprintf("messages/%s", lang.String())
	data, err := Asset(languageFile)
	if err != nil {
		panic(fmt.Sprintf("Invalid language asset: %v", languageFile))
	}
	i18n.LoadLanguageBytes(lang, data)

	flag.StringVar(&options.Reference, "reference", "", i18n.T("Name of the reference"))

	flag.StringVar(&options.Destination, "destination", imagec.DefaultDestination, i18n.T("Destination directory"))

	flag.StringVar(&options.Host, "host", imagec.DefaultPortLayerHost, i18n.T("Host that runs portlayer API (FQDN:port format)"))

	flag.StringVar(&options.Logfile, "logfile", imagec.DefaultLogfile, i18n.T("Path of the imagec log file"))

	flag.StringVar(&options.Username, "username", "", i18n.T("Username"))
	flag.StringVar(&options.Password, "password", "", i18n.T("Password"))

	flag.DurationVar(&options.Timeout, "timeout", imagec.DefaultHTTPTimeout, i18n.T("HTTP timeout"))

	flag.BoolVar(&options.Stdout, "stdout", false, i18n.T("Enable writing to stdout"))
	flag.BoolVar(&options.Debug, "debug", false, i18n.T("Show debug logging"))
	flag.BoolVar(&options.InsecureSkipVerify, "insecure-skip-verify", false, i18n.T("Don't verify certificates when fetching images"))
	flag.BoolVar(&options.InsecureAllowHTTP, "insecure-allow-http", false, i18n.T("Uses unencrypted connections when fetching images"))
	flag.BoolVar(&options.Standalone, "standalone", false, i18n.T("Disable port-layer integration"))

	flag.StringVar(&options.Profiling, "profile.mode", "", i18n.T("Enable profiling mode, one of [cpu, mem, block]"))
	flag.BoolVar(&options.Tracing, "tracing", false, i18n.T("Enable runtime tracing"))

	flag.Parse()
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
	switch options.Profiling {
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
	if options.Tracing {
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
	f, err := os.OpenFile(options.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open the logfile %s: %s", options.Logfile, err)
	}
	defer f.Close()

	// Initiliaze logger with default TextFormatter
	log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})

	// Set the log level
	if options.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// SetOutput to log file and/or stdout
	log.SetOutput(f)
	if options.Stdout {
		log.SetOutput(io.MultiWriter(os.Stdout, f))
	}

	ic := imagec.NewImageC(options, sf)
	err = ic.PullImage()
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Infof("Successfully pulled image %s", options.Reference)

}
