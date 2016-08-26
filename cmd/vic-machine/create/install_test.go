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

	"github.com/urfave/cli"
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
	create.ApplianceISO = tmpfile.Name()
	create.osType = "linux"
	if _, err = create.checkImagesFiles(nil); err == nil {
		t.Errorf("Error is expected for boot iso file is not found.")
	}
}

func writeImageVersion(fileName string, version string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Truncate(int64(0x10*2048) + 318); err != nil {
		return err
	}
	if _, err := f.WriteAt([]byte(version), int64(0x10*2048)+318); err != nil {
		return err
	}
	return nil
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
	create.ApplianceISO = ""
	create.BootstrapISO = tmpfile.Name()
	create.osType = "linux"
	var imageFiles map[string]string
	if _, err = create.checkImagesFiles(nil); err == nil {
		t.Errorf("Error is expected")
	}

	if err = writeImageVersion("appliance.iso", "Inc. 0.1-000-abcd"); err != nil {
		t.Error(err)
	}
	if err = writeImageVersion(tmpfile.Name(), "Inc. 0.1-000-abcd"); err != nil {
		t.Error(err)
	}

	cliContext := &cli.Context{
		App: &cli.App{
			Version: "Inconsistent",
		},
	}
	if _, err = create.checkImagesFiles(cliContext); err == nil {
		t.Errorf("Error is expected")
	}

	cliContext.App.Version = "0.1-000-abcd"
	if imageFiles, err = create.checkImagesFiles(cliContext); err != nil {
		t.Errorf("Error is returned: %s", err)
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
	create.noTLS = false
	if _, err := create.loadCertificate(); err != nil {
		t.Errorf("Error returned: %s", err)
	}
}
