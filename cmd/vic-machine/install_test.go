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
	"io/ioutil"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/install/configuration"
)

func TestImageNotFound(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	tmpfile, err := ioutil.TempFile("", "appIso")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	os.Args = []string{"cmd", fmt.Sprintf("-appIso=%s", tmpfile.Name())}
	flag.Parse()
	data.osType = "linux"
	data.conf = configuration.NewConfig()
	if err = checkImagesFiles(); err == nil {
		t.Errorf("Error is expected for boot iso file is not found.")
	}
}

func TestImageChecks(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	tmpfile, err := ioutil.TempFile("", "bootIso")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	_, err = os.Create("appliance.iso")
	if err != nil {
		t.Errorf("Failed to create default appliance iso file")
	}
	defer os.Remove("appliance.iso")

	os.Args = []string{"cmd", fmt.Sprintf("-bootIso=%s", tmpfile.Name())}
	flag.Parse()
	data.applianceISO = ""
	data.osType = "linux"
	data.conf = configuration.NewConfig()
	if err = checkImagesFiles(); err != nil {
		t.Errorf("Error returned: %s", err)
	}
	found := false
	for _, file := range data.conf.ImageFiles {
		if file == tmpfile.Name() {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Image file list does not contain input, %s", data.conf.ImageFiles)
	}
}

func TestLoadKey(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	os.Args = []string{"cmd"}
	flag.Parse()
	if err := loadCertificate(); err != nil {
		t.Errorf("Error returned: %s", err)
	}
}

func TestGenKey(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	os.Args = []string{"cmd"}
	flag.Parse()
	data.tlsGenerate = true
	if err := loadCertificate(); err != nil {
		t.Errorf("Error returned: %s", err)
	}
}
