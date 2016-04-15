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
	"testing"

	"github.com/docker/docker/pkg/stringid"
	"github.com/vmware/vic/metadata"
)

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

	config.ID = "sethostname"
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

	<-mocked.started

	// prevent indefinite wait in tether - normally session exit would trigger this
	close(reload)

	// wait for tether to exit
	<-mocked.cleaned

	expected := stringid.TruncateID(testConfig.ID)
	if mocked.hostname != expected {
		t.Errorf("expected: %s, actual: %s", expected, mocked.hostname)
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

	config.ID = "ipconfig"
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

	<-mocked.started

	// prevent indefinite wait in tether - normally session exit would trigger this
	close(reload)

	// wait for tether to exit
	<-mocked.cleaned

	// mocked state still exists and can be verified
	_ = testConfig

	testTeardown(t)
}

//
/////////////////////////////////////////////////////////////////////////////////////
