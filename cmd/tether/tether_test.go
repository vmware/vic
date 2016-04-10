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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/docker/docker/pkg/stringid"
	"github.com/vmware/vic/metadata"
)

var mockOps osopsMock

func TestMain(m *testing.M) {
	// use the mock ops
	mockOps = osopsMock{
		updated: make(chan bool, 1),
	}
	ops = &mockOps

	retCode := m.Run()

	// call with result of m.Run()
	os.Exit(retCode)
}

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

type testAttachServer struct {
	attachServerSSH

	updated chan bool
}

func (t *testAttachServer) start() error {
	err := t.attachServerSSH.start()

	t.updated <- true
	return err
}

func (t *testAttachServer) stop() {
	t.updated <- true
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

	if server != nil {
		server = &attachServerSSH{}
	}
}

func testTeardown(t *testing.T) {
	// let the main tether loop exit
	close(reload)
	// cleanup
	os.RemoveAll(pathPrefix)
	log.SetOutput(os.Stdout)
}

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
	testConfig, _ := cfg.LoadConfig()

	if err := run(cfg); err != nil {
		t.Error(err)
	}

	// check the exit code was sets
	if !testConfig.Sessions["feeddead"].Cmd.Cmd.ProcessState.Success() {
		t.Error("reference process 'data --reference=/' did not exit cleanly")
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

	if !bytes.Equal(out, log) {
		err := fmt.Errorf("Actual and expected output did not match\nExpected: %s\nActual:   %s\n", out, log)
		t.Error(err)
		return
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

/////////////////////////////////////////////////////////////////////////////////////
// TestHostnameConfig constructs the spec with no sesisons specifically for testing
// hostname setting - this checks the value passed to the mocked SetHostname
//
type TestSetHostnameConfig struct{}

func (c *TestSetHostnameConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestSetHostnameConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "deadbeef"
	config.Name = "tether_test_executor"

	return &config, nil
}

func TestSetHostname(t *testing.T) {
	testSetup(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestSetHostnameConfig{}
	testConfig, _ := cfg.LoadConfig()
	go func() {
		err := run(cfg)
		if err != nil {
			t.Error(err)
		}
	}()

	// wait for updates to occur
	<-mockOps.updated

	expected := stringid.TruncateID(testConfig.ID)
	if mockOps.hostname != expected {
		t.Errorf("expected: %s, actual: %s", expected, mockOps.hostname)
	}

	testTeardown(t)
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestSetIpAddressConfig constructs the spec for setting IP addresses - this checks
// the values passed to the Apply mock match those from the test config
//
type TestIPConfig struct{}

func (c *TestIPConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestIPConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "deadbeef"
	config.Name = "tether_test_executor"

	return &config, nil
}

func TestSetIpAddress(t *testing.T) {
	t.Skip("Network configuration processing not yet implemented")

	testSetup(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestIPConfig{}
	testConfig, _ := cfg.LoadConfig()
	go func() {
		err := run(cfg)
		if err != nil {
			t.Error(err)
		}
	}()

	// wait for updates to occur
	<-mockOps.updated

	// TEST LOGIC GOES HERE
	_ = testConfig

	testTeardown(t)
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestAttachConfig sets up the config for attach testing - the grep will echo anything
// sent and adds colour which is useful for tty testing
//
type TestAttachConfig struct{}

func (c *TestAttachConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestAttachConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "deadbeef"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"deadbeef": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "deadbeef",
				Name: "tether_test_session",
			},
			Tty:    false,
			Attach: true,
			Cmd: metadata.Cmd{
				Path: "/bin/grep",
				// grep, matching everything, reading from stdin
				Args: []string{"/bin/grep", ".", "-"},
				Env:  []string{},
				Dir:  "/",
			},
		},
	}

	// generate a host key for the tether
	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		return nil, err
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}

	config.Key = pem.EncodeToMemory(&privateKeyBlock)

	return &config, nil
}

func TestAttach(t *testing.T) {
	// supply custom attach server so we can inspect its state
	testServer := &testAttachServer{
		updated: make(chan bool, 10),
	}
	server = testServer

	testSetup(t)
	defer testTeardown(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestAttachConfig{}
	testConfig, err := cfg.LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		err := run(cfg)
		if err != nil {
			t.Error(err)
		}
	}()

	// wait for updates to occur
	<-testServer.updated

	if !testServer.enabled {
		t.Error("attach server was not enabled")
		return
	}

	// TEST LOGIC GOES HERE
	_ = testConfig

}

//
/////////////////////////////////////////////////////////////////////////////////////
