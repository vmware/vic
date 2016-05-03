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
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

/////////////////////////////////////////////////////////////////////////////////////
// TestPathLookup constructs the spec for a Session where the binary path must be
// resolved from the PATH environment variable - this is a variation from normal
// Cmd handling where that is done during creation of Cmd
//

func TestPathLookup(t *testing.T) {
	t.Skip("PATH based lookup not yet implemented")

	testSetup(t)
	defer testTeardown(t)

	cfg := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "pathlookup",
			Name: "tether_test_executor",
		},

		Sessions: map[string]metadata.SessionConfig{
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
		},
	}

	sink := map[string]string{}
	extraconfig.Encode(extraconfig.MapSink(sink), cfg)
	src := extraconfig.MapSource(sink)
	log.Debugf("Test configuration: %#v", sink)

	if err := run(src); err != nil {
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

func TestRelativePath(t *testing.T) {
	t.Skip("Relative path resolution not yet implemented")

	testSetup(t)
	defer testTeardown(t)

	cfg := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "relpath",
			Name: "tether_test_executor",
		},

		Sessions: map[string]metadata.SessionConfig{
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
		},
	}

	sink := map[string]string{}
	extraconfig.Encode(extraconfig.MapSink(sink), cfg)
	src := extraconfig.MapSource(sink)
	log.Debugf("Test configuration: %#v", sink)

	if err := run(src); err != nil {
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

func TestAbsPath(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	cfg := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "abspath",
			Name: "tether_test_executor",
		},

		Sessions: map[string]metadata.SessionConfig{
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
		},
	}

	sink := map[string]string{}
	extraconfig.Encode(extraconfig.MapSink(sink), cfg)
	src := extraconfig.MapSource(sink)
	log.Debugf("Test configuration: %#v", sink)

	if err := run(src); err != nil {
		t.Error(err)
	}

	// check the exit code was set
	status := config.Sessions["abspath"].ExitStatus
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

func TestMissingBinary(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	cfg := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "missing",
			Name: "tether_test_executor",
		},

		Sessions: map[string]metadata.SessionConfig{
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
		},
	}

	sink := map[string]string{}
	extraconfig.Encode(extraconfig.MapSink(sink), cfg)
	src := extraconfig.MapSource(sink)
	log.Debugf("Test configuration: %#v", sink)

	err := run(src)
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
}

//
/////////////////////////////////////////////////////////////////////////////////////
