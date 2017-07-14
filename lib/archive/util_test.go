// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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

package archive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// COPY TO TESTS

func TestComplexWriteSpec(t *testing.T) {
	copyTarget := "/mnt/vols"
	direction := CopyTo

	mounts := []testMount{

		{
			mount:              "/",
			CopyTarget:         copyTarget,
			primaryMountTarget: true,
			direction:          direction,
		},
		{
			mount:              "/mnt/vols/A",
			CopyTarget:         copyTarget,
			primaryMountTarget: false,
			direction:          direction,
		},
		{
			mount:              "/mnt/vols/B",
			CopyTarget:         copyTarget,
			primaryMountTarget: false,
			direction:          direction,
		},
		{
			mount:              "/mnt/vols/A/subvols/AB",
			CopyTarget:         copyTarget,
			primaryMountTarget: false,
			direction:          direction,
		},
	}

	expectedFilterSpecs := map[string]FilterSpec{
		"/": FilterSpec{
			RebasePath: "mnt/vols",
			StripPath:  "",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/A": FilterSpec{
			RebasePath: "",
			StripPath:  "A",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/B": FilterSpec{
			RebasePath: "",
			StripPath:  "B",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/A/subvols/AB": FilterSpec{
			RebasePath: "",
			StripPath:  "A/subvols/AB",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
	}

	for _, v := range mounts {
		actualFilterSpec := GenerateFilterSpec(v.CopyTarget, v.mount, v.primaryMountTarget, v.direction)
		expectedFilterSpec := expectedFilterSpecs[v.mount]

		if !assert.Equal(t, expectedFilterSpec.RebasePath, actualFilterSpec.RebasePath, "rebase path check failed (%s)", v.mount) {
			return
		}

		if !assert.Equal(t, expectedFilterSpec.StripPath, actualFilterSpec.StripPath, "strip path check failed (%s)", v.mount) {
			return
		}

	}

}

// COPY FROM TESTS

func TestComplexReadSpec(t *testing.T) {
	copyTarget := "/mnt/vols"
	direction := CopyFrom

	mounts := []testMount{

		{
			mount:              "/",
			CopyTarget:         copyTarget,
			primaryMountTarget: true,
			direction:          direction,
		},
		{
			mount:              "/mnt/vols/A",
			CopyTarget:         copyTarget,
			primaryMountTarget: false,
			direction:          direction,
		},
		{
			mount:              "/mnt/vols/B",
			CopyTarget:         copyTarget,
			primaryMountTarget: false,
			direction:          direction,
		},
		{
			mount:              "/mnt/vols/A/subvols/AB",
			CopyTarget:         copyTarget,
			primaryMountTarget: false,
			direction:          direction,
		},
	}

	expectedFilterSpecs := map[string]FilterSpec{
		"/": FilterSpec{
			RebasePath: "vols",
			StripPath:  "mnt/vols",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/A": FilterSpec{
			RebasePath: "vols/A",
			StripPath:  "",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/B": FilterSpec{
			RebasePath: "vols/B",
			StripPath:  "",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/A/subvols/AB": FilterSpec{
			RebasePath: "vols/A/subvols/AB",
			StripPath:  "",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
	}

	for _, v := range mounts {
		actualFilterSpec := GenerateFilterSpec(v.CopyTarget, v.mount, v.primaryMountTarget, v.direction)
		expectedFilterSpec := expectedFilterSpecs[v.mount]

		if !assert.Equal(t, expectedFilterSpec.RebasePath, actualFilterSpec.RebasePath, "rebase path check failed (%s)", v.mount) {
			return
		}

		if !assert.Equal(t, expectedFilterSpec.StripPath, actualFilterSpec.StripPath, "strip path check failed (%s)", v.mount) {
			return
		}

	}

}

// test utility functions and structs

type testMount struct {
	mount              string
	CopyTarget         string
	primaryMountTarget bool
	direction          bool
}
