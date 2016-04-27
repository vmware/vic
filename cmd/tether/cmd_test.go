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
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/metadata"
)

/////////////////////////////////////////////////////////////////////////////////////
// TestPathLookup constructs the spec for a Session where the binary path must be
// resolved from the PATH environment variable - this is a variation from normal
// Cmd handling where that is done during creation of Cmd
//

type TestPathLookupConfig struct{}

func (c *TestPathLookupConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestPathLookupConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "pathlookup"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"pathlookup": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "pathlookup",
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

func TestPathLookup(t *testing.T) {
	t.Skip("PATH based lookup not yet implemented")

	testSetup(t)
	defer testTeardown(t)

	if err := run(&TestPathLookupConfig{}); err != nil {
		t.Error(err)
		return
	}
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestRelativePathConfig constructs the spec for a Session with relative path to existing
// binary - this tests resolution vs the working directory
//

type TestRelativePathConfig struct{}

func (c *TestRelativePathConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestRelativePathConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "relpath"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"relpath": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "relpath",
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

func TestRelativePath(t *testing.T) {
	t.Skip("Relative path resolution not yet implemented")

	testSetup(t)
	defer testTeardown(t)

	if err := run(&TestRelativePathConfig{}); err != nil {
		t.Error(err)
		return
	}
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestAbsPathConfig constructs the spec for a Session with absolute path to existing
// binary
//
type TestAbsPathConfig struct{}

func (c *TestAbsPathConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestAbsPathConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "abspath"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"abspath": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "abspath",
				Name: "tether_test_session",
			},
			Tty: false,
			Cmd: metadata.Cmd{
				// test abs path
				Path: "/bin/date",
				Args: []string{"date", "--reference=/"},
				Env:  []string{},
				Dir:  "/",
			},
		},
	}

	return &config, nil
}

func TestAbsPath(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestAbsPathConfig{}

	if err := run(cfg); err != nil {
		t.Error(err)
	}

	// check the exit code was set
	status := sessions["abspath"].exitStatus
	if status != 0 {
		t.Errorf("reference process 'data --reference=/' did not exit cleanly: %d", status)
		return
	}

	// read the output from the session
	log, err := ioutil.ReadFile(pathPrefix + "/ttyS2")
	if err != nil {
		fmt.Printf("Failed to open log file for command: %s", err)
		t.Error(err)
		return
	}

	// run the command directly
	out, err := exec.Command("/bin/date", "--reference=/").Output()
	if err != nil {
		fmt.Printf("Failed to run date for comparison data: %s", err)
		t.Error(err)
		return
	}

	if !assert.Equal(t, out, log) {
		return
	}
}

func TestAbsPathRepeat(t *testing.T) {
	t.Skip("Occasional issues with output not being flushed to log - #577")

	for i := 0; i < 1000; i++ {
		TestAbsPath(t)
	}
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestMissingBinaryConfig constructs the spec for a Session with invalid binary path
//
type TestMissingBinaryConfig struct{}

func (c *TestMissingBinaryConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestMissingBinaryConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "missing"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"missing": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "missing",
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

func TestMissingBinary(t *testing.T) {
	testSetup(t)

	err := run(&TestMissingBinaryConfig{})
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

//
/////////////////////////////////////////////////////////////////////////////////////
