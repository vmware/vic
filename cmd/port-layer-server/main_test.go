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
	"os"
	"testing"
)

var systemTest *bool

func init() {
	systemTest = flag.Bool("systemTest", false, "Run system test")
}

func TestMain(m *testing.M) {
	// collect the command line options
	// known to the standard flags package
	flagOpts := map[string]interface{}{}
	flag.VisitAll(func(f *flag.Flag) {
		flagOpts[f.Name] = nil
	})

	// parse the command line using the go-flags
	// package, getting a command line back
	// that contains any options only
	// known to the standard flags package
	flagArgs, err := parseArgs(flagOpts)
	if err != nil {
		os.Exit(1)
	}

	// now parse the "remaining" command line
	// with the standard flags package
	flag.CommandLine.Parse(flagArgs)

	// execute tests
	os.Exit(m.Run())
}

func TestSystem(t *testing.T) {
	if *systemTest {
		main()
	}
}
