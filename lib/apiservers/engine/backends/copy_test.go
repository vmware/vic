// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package backends

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
)

type MockCopyToData struct {
	containerDestPath string
	tarAssetName      string
	expectedPrefix    string
}

type MockCopyFromData struct {
	containerSourcePath string
	expectedPrefices    []string
}

func TestFindArchiveWriter(t *testing.T) {
	mounts := []types.MountPoint{
		{Name: "volA", Destination: "/mnt/A"},
		{Name: "volAB", Destination: "/mnt/A/AB"},
		{Name: "volB", Destination: "/mnt/B"},
		{Name: "R/W", Destination: "/"},
	}

	mockData := []MockCopyToData{
		// mock data for tar asset as a file and container dest path including a mount point
		{
			containerDestPath: "/mnt/A/",
			tarAssetName:      "file.txt",
			expectedPrefix:    "/mnt/A",
		},
		{
			containerDestPath: "/mnt/A/AB",
			tarAssetName:      "file.txt",
			expectedPrefix:    "/mnt/A/AB",
		},
		// mock data for tar asset containing a mount point and the container dest path as /
		{
			containerDestPath: "/",
			tarAssetName:      "mnt/A/file.txt",
			expectedPrefix:    "/mnt/A",
		},
		{
			containerDestPath: "/",
			tarAssetName:      "mnt/A/AB/file.txt",
			expectedPrefix:    "/mnt/A/AB",
		},
		// mock data for cases that do not involve mount points
		{
			containerDestPath: "/",
			tarAssetName:      "test/file.txt",
			expectedPrefix:    "/",
		},
	}

	for _, data := range mockData {
		writerMap := NewArchiveStreamWriterMap(mounts, data.containerDestPath)
		aw, err := writerMap.FindArchiveWriter(data.containerDestPath, data.tarAssetName)
		assert.Nil(t, err, "Expected success from finding archive writer for container dest %s and tar asset path %s", data.containerDestPath, data.tarAssetName)
		assert.NotNil(t, aw, "Expected non-nil archive writer")
		if aw != nil {
			assert.Contains(t, aw.mountPoint.Destination, data.expectedPrefix,
				"Expected to find prefix %s for container dest %s and tar asset path %s",
				data.expectedPrefix, data.containerDestPath, data.tarAssetName)
		}
	}
}

func TestFindArchiveReaders(t *testing.T) {
	mounts := []types.MountPoint{
		{Name: "volA", Destination: "/mnt/A"},     //mount point
		{Name: "volAB", Destination: "/mnt/A/AB"}, //mount point
		{Name: "volB", Destination: "/mnt/B"},     //mount point
		{Name: "R/W", Destination: "/"},           //container base volume
	}

	mockData := []MockCopyFromData{
		// case 1: Get all mount prefix
		{
			containerSourcePath: "/",
			expectedPrefices:    []string{"/", "/mnt/A", "/mnt/B", "/mnt/A/AB"},
		},
		{
			containerSourcePath: "/mnt",
			expectedPrefices:    []string{"/", "/mnt/A", "/mnt/B", "/mnt/A/AB"},
		},
		{
			containerSourcePath: "/mnt/",
			expectedPrefices:    []string{"/", "/mnt/A", "/mnt/B", "/mnt/A/AB"},
		},
		// case 3: Do not include /mnt/B
		{
			containerSourcePath: "/mnt/A",
			expectedPrefices:    []string{"/mnt/A", "/mnt/A/AB"},
		},
		// case 4: Return the container base "/"
		{
			containerSourcePath: "/mnt/not-a-mount",
			expectedPrefices:    []string{"/"},
		},
		{
			containerSourcePath: "/etc/",
			expectedPrefices:    []string{"/"},
		},
	}

	readerMap := NewArchiveStreamReaderMap(mounts)

	for _, data := range mockData {
		archiveReaders, err := readerMap.FindArchiveReaders(data.containerSourcePath)
		assert.Nil(t, err, "Expected success from finding archive readers for container source %s", data.containerSourcePath)
		assert.NotNil(t, archiveReaders, "Expected an array of archive readers but got nil for container source path %s", data.containerSourcePath)
		assert.NotEmpty(t, archiveReaders, "Expected an array of archive readers %s with more than one items", data.containerSourcePath)

		log.Debugf("Data = %#v", data)
		pa := PrefixArray(archiveReaders)
		nonOverlap := UnionMinusIntersection(pa, data.expectedPrefices)
		assert.Empty(t, nonOverlap, "Found mismatch in the prefix array and expected array for source path %s.  Non-overlapped result = %#v", data.containerSourcePath, nonOverlap)
	}
}

func PrefixArray(readers []ArchiveReader) (pa []string) {
	for _, reader := range readers {
		pa = append(pa, reader.mountPoint.Destination)
	}

	log.Debugf("prefix array - %#v", pa)
	return
}

func UnionMinusIntersection(A, B []string) (res []string) {
	test := make(map[string]bool)

	log.Debugf("Looking for non overlapping in array A-%#v and array B-%#v", A, B)

	for _, data := range A {
		test[data] = true
	}

	for _, data := range B {
		if _, ok := test[data]; ok {
			delete(test, data)
		} else {
			res = append(res, data)
		}
	}

	for key := range test {
		res = append(res, key)
	}

	log.Debugf("Resulting non overlapped array - %#v", res)

	return
}
