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
	"testing"

	"github.com/docker/docker/pkg/stringid"
	"github.com/vmware/vic/lib/metadata"
)

/////////////////////////////////////////////////////////////////////////////////////
// TestHostnameConfig constructs the spec with no sesisons specifically for testing
// hostname setting - this checks the value passed to the mocked SetHostname
//
func TestSetHostname(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	cfg := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "sethostname",
			Name: "tether_test_executor",
		},
	}

	startTether(t, &cfg)

	<-mocked.started

	// prevent indefinite wait in tether - normally session exit would trigger this
	close(reload)

	// wait for tether to exit
	<-mocked.cleaned

	expected := stringid.TruncateID(cfg.ID)
	if mocked.hostname != expected {
		t.Errorf("expected: %s, actual: %s", expected, mocked.hostname)
	}
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestSetIpAddressConfig constructs the spec for setting IP addresses - this checks
// the values passed to the Apply mock match those from the test config
//
func TestSetIpAddress(t *testing.T) {
	t.Skip("Not yet testing network config")

	testSetup(t)
	defer testTeardown(t)

	cfg := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   "ipconfig",
			Name: "tether_test_executor",
		},
	}

	startTether(t, &cfg)

	<-mocked.started

	// prevent indefinite wait in tether - normally session exit would trigger this
	close(reload)

	// wait for tether to exit
	<-mocked.cleaned

	testTeardown(t)
}

//
/////////////////////////////////////////////////////////////////////////////////////
