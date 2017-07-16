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
		"/": {
			RebasePath: "mnt/vols",
			StripPath:  "",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/A": {
			RebasePath: "",
			StripPath:  "A",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/B": {
			RebasePath: "",
			StripPath:  "B",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/A/subvols/AB": {
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

func TestWriteIntoMountSpec(t *testing.T) {
	copyTarget := "/mnt/vols/A/a/path"
	direction := CopyTo

	mounts := []testMount{
		// "/" will exist but since it is past the write path we do not care about the filterspec. It will be filled out as a non primary target in any case. and is a non primary target that comes before the primary target.
		{
			mount:              "/",
			CopyTarget:         copyTarget,
			primaryMountTarget: false,
			direction:          direction,
		},
		{
			mount:              "/mnt/vols/A",
			CopyTarget:         copyTarget,
			primaryMountTarget: true,
			direction:          direction,
		},
	}

	expectedFilterSpecs := map[string]FilterSpec{
		// not a valid spec since this is past the copy target. The filterspec will be completely bogus. The generateFilterSpec function assumes you have given it a target that lives along the CopyTarget. In our case here we have "/" as the mount point and the target as "/mnt/vols/A/a/path"
		"/": {
			RebasePath: "",
			StripPath:  "",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/A": {
			RebasePath: "a/path",
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

func TestWriteFileIntoMountSpec(t *testing.T) {
	copyTarget := "/mnt/vols/A/a/path/file.txt"
	direction := CopyTo

	mounts := []testMount{
		{
			mount:              "/mnt/vols/A",
			CopyTarget:         copyTarget,
			primaryMountTarget: true,
			direction:          direction,
		},
	}

	expectedFilterSpecs := map[string]FilterSpec{
		"/mnt/vols/A": {
			RebasePath: "a/path/file.txt",
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

func TestWriteFilespec(t *testing.T) {
	copyTarget := "/path/file.txt"
	direction := CopyTo

	mounts := []testMount{
		{
			mount:              "/",
			CopyTarget:         copyTarget,
			primaryMountTarget: true,
			direction:          direction,
		},
	}

	expectedFilterSpecs := map[string]FilterSpec{
		"/": {
			RebasePath: "path/file.txt",
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
		"/": {
			RebasePath: "vols",
			StripPath:  "mnt/vols",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/A": {
			RebasePath: "vols/A",
			StripPath:  "",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/B": {
			RebasePath: "vols/B",
			StripPath:  "",
			Exclusions: make(map[string]struct{}),
			Inclusions: make(map[string]struct{}),
		},
		"/mnt/vols/A/subvols/AB": {
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

func TestReadIntoMountSpec(t *testing.T) {
	copyTarget := "/mnt/vols/A/a/path"
	direction := CopyFrom

	mounts := []testMount{
		{
			mount:              "/mnt/vols/A",
			CopyTarget:         copyTarget,
			primaryMountTarget: true,
			direction:          direction,
		},
	}

	expectedFilterSpecs := map[string]FilterSpec{
		"/mnt/vols/A": {
			RebasePath: "path",
			StripPath:  "a/path",
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

func TestReadFileSpec(t *testing.T) {
	copyTarget := "/mnt/vols/A/a/path/file.txt"
	direction := CopyFrom

	mounts := []testMount{
		// "/" will exist but since it is past the write path we do not care about the filterspec. It will be filled out as a non primary target in any case.
		{
			mount:              "/",
			CopyTarget:         copyTarget,
			primaryMountTarget: true,
			direction:          direction,
		},
	}

	expectedFilterSpecs := map[string]FilterSpec{
		// not since this is past the copy target the filterspec will be completely bogus. The generateFilterSpec function assumes you have given it a target that lives along the CopyTarget. In our case here we have "/" as the mount point and the target as "/mnt/vols/A/a/path"
		"/": {
			RebasePath: "file.txt",
			StripPath:  "mnt/vols/A/a/path/file.txt",
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

func TestReadFileInMountSpec(t *testing.T) {
	copyTarget := "/mnt/vols/A/a/path/file.txt"
	direction := CopyFrom

	mounts := []testMount{
		{
			mount:              "/mnt/vols/A",
			CopyTarget:         copyTarget,
			primaryMountTarget: true,
			direction:          direction,
		},
	}

	expectedFilterSpecs := map[string]FilterSpec{
		"/mnt/vols/A": {
			RebasePath: "file.txt",
			StripPath:  "a/path/file.txt",
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

// ADD EXCLUSIONS TESTS

func TestAddExclusionsNonNested(t *testing.T) {

	testMounts := []string{
		"/",
		"/mnt/A",
		"/mnt/B",
		"/mnt/C",
	}

	expectedResults := map[string]FilterSpec{
		"/": FilterSpec{
			Exclusions: map[string]struct{}{
				"mnt/A": {},
				"mnt/B": {},
				"mnt/C": {},
			},
		},
		"/mnt/A": FilterSpec{
			Exclusions: make(map[string]struct{}),
		},
		"/mnt/B": FilterSpec{
			Exclusions: make(map[string]struct{}),
		},
		"/mnt/C": FilterSpec{
			Exclusions: make(map[string]struct{}),
		},
	}

	for _, mount := range testMounts {
		spec := FilterSpec{
			Exclusions: make(map[string]struct{}),
		}

		AddMountExclusions(mount, &spec, testMounts)

		expectedSpec := expectedResults[mount]
		expectedLength := len(expectedSpec.Exclusions)
		actualLength := len(spec.Exclusions)
		if !assert.Equal(t, expectedLength, actualLength, "there were %d entries instead of %d for the exclusions generated for mount (%s)", expectedLength, actualLength, mount) {
			return
		}

		for path := range expectedSpec.Exclusions {
			_, ok := spec.Exclusions[path]

			if !assert.True(t, ok, "Should have found (%s) in the exclusion map for %s", path, mount) {
				return
			}

		}

	}

}

func TestAddExclusionsNestedMounts(t *testing.T) {

	testMounts := []string{
		"/",
		"/mnt/A",
		"/mnt/B",
		"/mnt/C",
		"/mnt/A/dir/AB",
		"/mnt/A/dir/AC",
	}

	expectedResults := map[string]FilterSpec{
		"/": FilterSpec{
			Exclusions: map[string]struct{}{
				"mnt/A":        {},
				"mnt/B":        {},
				"mnt/C":        {},
				"mnt/A/dir/AB": {},
				"mnt/A/dir/AC": {},
			},
		},
		"/mnt/A": FilterSpec{
			Exclusions: map[string]struct{}{
				"dir/AB": {},
				"dir/AC": {},
			},
		},
		"/mnt/B": FilterSpec{
			Exclusions: make(map[string]struct{}),
		},
		"/mnt/C": FilterSpec{
			Exclusions: make(map[string]struct{}),
		},
	}

	for _, mount := range testMounts {
		spec := FilterSpec{
			Exclusions: make(map[string]struct{}),
		}

		AddMountExclusions(mount, &spec, testMounts)

		expectedSpec := expectedResults[mount]
		expectedLength := len(expectedSpec.Exclusions)
		actualLength := len(spec.Exclusions)
		if !assert.Equal(t, expectedLength, actualLength, "there were %d entries instead of %d for the exclusions generated for mount (%s)", expectedLength, actualLength, mount) {
			return
		}

		for path := range expectedSpec.Exclusions {
			_, ok := spec.Exclusions[path]

			if !assert.True(t, ok, "Should have found (%s) in the exclusion map for %s \n exlusion map: %s", path, mount, spec.Exclusions) {
				return
			}

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
