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

package create

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

var (
	create = NewCreate()
)

func TestImageNotFound(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	tmpfile, err := ioutil.TempFile("", "appIso")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	os.Args = []string{"cmd", "create", fmt.Sprintf("-appliance-iso=%s", tmpfile.Name())}
	flag.Parse()
	create.applianceISO = tmpfile.Name()
	create.osType = "linux"
	if _, err = create.checkImagesFiles(); err == nil {
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

	os.Args = []string{"cmd", "create", fmt.Sprintf("-bootstrap-iso=%s", tmpfile.Name())}
	flag.Parse()
	create.applianceISO = ""
	create.bootstrapISO = tmpfile.Name()
	create.osType = "linux"
	var imageFiles []string
	if imageFiles, err = create.checkImagesFiles(); err != nil {
		t.Errorf("Error returned: %s", err)
	}
	found := false
	for _, file := range imageFiles {
		if file == tmpfile.Name() {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Image file list does not contain input, %s", imageFiles)
	}
}

func TestLoadKey(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	os.Args = []string{"cmd", "create"}
	flag.Parse()
	if _, err := create.loadCertificate(); err != nil {
		t.Errorf("Error returned: %s", err)
	}
}

func TestGenKey(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	os.Args = []string{"cmd", "create"}
	flag.Parse()
	create.tlsGenerate = true
	if _, err := create.loadCertificate(); err != nil {
		t.Errorf("Error returned: %s", err)
	}
}
