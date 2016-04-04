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
	"os"
	"os/exec"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func setup(t *testing.T) {
	log.SetLevel(log.DebugLevel)
}

func TestSet(t *testing.T) {
	os.Args = []string{"cmd", "-set", "foo", "bar"}
	flag.Parse()

	set = true
	main()
}

func TestSetNoArgs(t *testing.T) {

	cmd := exec.Command(os.Args[0], "-test.run=TestSetNoArgs", os.Args[1], "-set=true")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		log.Debugf("Excepted error")
	}
}

func TestGetNoArgs(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	cmd := exec.Command(os.Args[0], "-test.run=TestGetNoArgs", os.Args[1], "-get=true")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		log.Debugf("Excepted error")
	}
}

func TestGet(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	os.Args = []string{"cmd", "-get", "foo"}
	flag.Parse()

	set = false
	get = true
	main()
}

func TestForkNoArgs(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	cmd := exec.Command(os.Args[0], "-test.run=TestForkNoArgs", os.Args[1], "-fork=true")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		log.Debugf("Excepted error")
	}
}
