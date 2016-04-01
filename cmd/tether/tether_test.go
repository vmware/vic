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
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/vmware/vic/cmd/tether/utils"
	"github.com/vmware/vic/metadata"
)

// createFakeDevices creates regular files or pipes in place of the char devices used
// in a full VM
func createFakeDevices() error {
	// create serial devices
	for i := 0; i < 3; i++ {
		path := fmt.Sprintf("%s/ttyS%d", pathPrefix, i)
		_, err := os.Create(path)
		if err != nil {
			detail := fmt.Sprintf("failed to create %s for com%d: %s", path, i+1, err)
			return errors.New(detail)
		}
	}

	// make an access to urandom
	path := fmt.Sprintf("%s/urandom", pathPrefix)
	err := os.Symlink("/dev/urandom", path)
	if err != nil {
		detail := fmt.Sprintf("failed to create urandom access: %s", err)
		return errors.New(detail)
	}

	return nil
}

type TestMissingBinaryConfig struct{}

func (c *TestMissingBinaryConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestMissingBinaryConfig) LoadConfig(blobl string) (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "deadbeef"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"feebdaed": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "feebdaed",
				Name: "tether_test_session",
			},
			Tty: false,
			Cmd: metadata.Cmd{
				// test relative path
				Path: "/not/there",
				Args: []string{"/not/there"},
				Env:  []string{"PATH=/not"},
				Dir:  "/",
			},
		},
	}

	return &config, nil
}

type TestRelativePathConfig struct{}

func (c *TestRelativePathConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestRelativePathConfig) LoadConfig(blobl string) (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "deadbeef"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"feebdaed": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "feebdaed",
				Name: "tether_test_session",
			},
			Tty: false,
			Cmd: metadata.Cmd{
				// test relative path
				Path: "date",
				Args: []string{"date", "--reference=/"},
				Env:  []string{"PATH=/bin"},
				Dir:  "/bin",
			},
		},
	}

	return &config, nil
}

type TestAbsPathConfig struct{}

func (c *TestAbsPathConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestAbsPathConfig) LoadConfig(blobl string) (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "deadbeef"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"feebdaed": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "feebdaed",
				Name: "tether_test_session",
			},
			Tty: false,
			Cmd: metadata.Cmd{
				// test relative path
				Path: "/bin/date",
				Args: []string{"date", "--reference=/"},
				Env:  []string{},
				Dir:  "/",
			},
		},
	}

	return &config, nil
}

func testSetup(t *testing.T) {
	var err error

	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	pathPrefix, err = ioutil.TempDir("", path.Base(name))
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}
	utils.SetPathPrefix(pathPrefix)

	err = os.MkdirAll(pathPrefix, 0777)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	err = createFakeDevices()
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}
}

func testTeardown(t *testing.T) {
	// cleanup
	os.RemoveAll(pathPrefix)
	log.SetOutput(os.Stdout)
}

func TestRelativePath(t *testing.T) {
	t.Skip("Relative path resolution not yet implemented")

	testSetup(t)

	if err := run(&TestRelativePathConfig{}, ""); err != nil {
		t.Error(err)
	}

	testTeardown(t)
}

func TestAbsPath(t *testing.T) {
	testSetup(t)

	if err := run(&TestAbsPathConfig{}, ""); err != nil {
		t.Error(err)
	}

	// read the output from the session
	log, err := ioutil.ReadFile(pathPrefix + "/ttyS2")
	if err != nil {
		fmt.Printf("Failed to open log file for command: %s", err)
		t.Error(err)
	}

	// run the command directly
	out, err := exec.Command("/bin/date", "--reference=/").Output()
	if err != nil {
		fmt.Printf("Failed to run date for comparison data: %s", err)
		t.Error(err)
	}

	if !bytes.Equal(out, log) {
		err := fmt.Errorf("Actual and expected output did not match\nExpected: %s\nActual:   %s\n", out, log)
		t.Error(err)
	}

	testTeardown(t)
}

func TestMissingBinary(t *testing.T) {
	testSetup(t)

	err := run(&TestMissingBinaryConfig{}, "")
	if err == nil {
		t.Error("Expected error from missing binary")
	}

	// read the output from the session
	log, err := ioutil.ReadFile(pathPrefix + "/ttyS2")
	if err != nil {
		fmt.Printf("Failed to open log file for command: %s", err)
	}
	if len(log) > 0 {
		fmt.Printf("Command output: %s", string(log))
	}

	testTeardown(t)
}

func TestSetIpAddress(t *testing.T) {
	testSetup(t)

	// err := run(&TestIPConfig{})

	testTeardown(t)
}
