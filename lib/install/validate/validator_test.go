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

package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestValidator struct {
	Validator
}

func TestPathConversionESX(t *testing.T) {
	v := &TestValidator{}

	v.DatacenterPath = "/ha-datacenter"

	ipath := v.computePathToInventoryPath("/")
	assert.Equal(t, "/ha-datacenter/host/*/Resources", ipath, "Expected top level resource pool")
	cpath := v.inventoryPathToComputePath(ipath)
	assert.Equal(t, "*", cpath, "Expected root resource specifier")

	ipath = v.computePathToInventoryPath("*")
	assert.Equal(t, "/ha-datacenter/host/*/Resources/*", ipath, "Expected top level resource pool")
	cpath = v.inventoryPathToComputePath(ipath)
	assert.Equal(t, "*/*", cpath, "Expected top level wildcard specifier")
}

func TestSampleConversion(t *testing.T) {
	v := &TestValidator{}
	v.DatacenterPath = "/ha-datacenter"

	translations := map[string]string{
		"/":              "/ha-datacenter/host/*/Resources",
		"*":              "/ha-datacenter/host/*/Resources/*",
		"testpool":       "/ha-datacenter/host/*/Resources/testpool",
		"test/deep/path": "/ha-datacenter/host/*/Resources/test/deep/path",
	}

	for in, expected := range translations {
		ipath := v.computePathToInventoryPath(in)
		assert.Equal(t, expected, ipath, "Translation did not match")
	}
}
