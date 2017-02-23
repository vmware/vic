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

package nfs

import (
	"context"
	"io"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/pkg/trace"
)

const (
	nfsTestDir = "NFSVolumeStoreTests"
)

// MOCK TARGET STRUCT AND IMPL
type MockTarget struct {
	dirPath string
}

func NewMocktarget(path string) MockTarget {
	return MockTarget{dirPath: path}
}

func (v MockTarget) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (v MockTarget) OpenFile(path string, mode os.FileMode) (io.ReadWriteCloser, error) {
	return os.OpenFile(path, os.O_RDWR, mode)
}

func (v MockTarget) Create(path string, perm os.FileMode) (io.ReadWriteCloser, error) {
	return os.Create(path)
}

func (v MockTarget) Mkdir(path string, perm os.FileMode) ([]byte, error) {
	return nil, os.Mkdir(path, perm)
}

func (v MockTarget) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (v MockTarget) ReadDir(path string) ([]os.FileInfo, error) {
	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return nil, err
	}

	return dir.Readdir(0)
}

func (v MockTarget) Lookup(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// MOCK MOUNT STRUCT AND IMPL

type MockMount struct {
}

func (m MockMount) Mount(target *url.URL) (Target, error) {
	return NewMocktarget(target.Path), nil
}

func (m MockMount) Unmount(target Target) error {
	return nil
}

func TestMain(m *testing.M) {
	testPath := path.Join(os.TempDir(), nfsTestDir)
	os.Mkdir(testPath, 0755)
	result := m.Run()
	os.RemoveAll(testPath)
	os.Exit(result)
}

func TestSimpleVolumeStoreOperations(t *testing.T) {
	mockMount := MockMount{}
	testVolName := "testVolume"
	dirpath := path.Join(os.TempDir(), nfsTestDir)

	targetURL := url.URL{Path: dirpath}
	op := trace.NewOperation(context.TODO(), "TestOp")

	//Create a Volume Store
	vs, err := NewVolumeStore(op, "testStore", &targetURL, mockMount)
	if !assert.NoError(t, err, "Failed during call to NewVolumeStore with err (%s)", err) {
		return
	}

	_, err = os.Stat(dirpath)
	if !assert.NoError(t, err, "Could not find the initial volume store directory after creation of volume store. err (%s)", err) {
		return
	}

	if !assert.NotNil(t, vs, "Volume Store created with nil err, but return is also nil") {
		return
	}

	info := make(map[string][]byte)
	testInfoKey := "junk"
	info[testInfoKey] = make([]byte, 20)

	file, err := os.Open(dirpath)
	defer file.Close()

	//Create a Volume
	vol, err := vs.VolumeCreate(op, testVolName, vs.Target, 0 /*we do not use this*/, info)
	if !assert.NoError(t, err, "Failed during call to VolumeCreate with err (%s)", err) {
		return
	}

	if !assert.Equal(t, testVolName, vol.ID, "expected volume ID (%s) got ID (%s)", testVolName, vol.ID) {
		return
	}

	_, ok := vol.Info[testInfoKey]
	if !assert.True(t, ok, "TestInfoKey did not exist in the return metadata map") {
		return
	}

	//Check Metadata Pathing
	metadataPath := path.Join(dirpath, metadataDir)
	volumePath := path.Join(dirpath, VolumesDir)

	metaFilesDir, err := os.Open(metadataPath)
	defer metaFilesDir.Close()
	if !assert.NoError(t, err, "opening the metadata directory failed with err (%s)", err) {
		return
	}

	volumeDir, err := os.Open(volumePath)
	defer volumeDir.Close()
	if !assert.NoError(t, err, "Opening the volume directory failed with err (%s)", err) {
		return
	}

	metaDirEntries, err := metaFilesDir.Readdir(0)
	if !assert.NoError(t, err, "Failed to read the metadata directory with err (%s)", err) {
		return
	}

	volumeDirEntries, err := volumeDir.Readdir(0)
	if !assert.NoError(t, err, "Failed to read the volume data directory with err (%s)", err) {
		return
	}

	t.Logf("meta directory (%s)", metaDirEntries)
	if !assert.Equal(t, len(metaDirEntries), 1, "expected metadata directory to have 1 entry and it had (%s)", len(metaDirEntries)) {
		return
	}
	t.Logf("vol directory (%s)", volumeDirEntries)
	if !assert.Equal(t, len(volumeDirEntries), 1, "expected metadata directory to have 1 entry and it had (%s)", len(volumeDirEntries)) {
		return
	}

	//Remove the Volume
	err = vs.VolumeDestroy(op, vol)
	if !assert.NoError(t, err, "Failed during a call to VolumeDestroy with err (%s)", err) {
		return
	}

	volumeMetadatapath := path.Join(metadataPath, vol.ID)
	volDirCheck, err := os.Open(volumeMetadatapath)
	defer volDirCheck.Close()
	if !assert.Error(t, err, "expected the path (%s) to be deleted after the deletion of vol (%s)", volumeMetadatapath, vol.ID) {
		return
	}

	metaFilesDir, err = os.Open(metadataPath)
	defer metaFilesDir.Close()
	if !assert.NoError(t, err, "opening the metadata directory failed with err (%s)", err) {
		return
	}

	volumeDir, err = os.Open(volumePath)
	defer volumeDir.Close()
	if !assert.NoError(t, err, "Opening the volume directory failed with err (%s)", err) {
		return
	}

	metaDirEntries, err = metaFilesDir.Readdir(0)
	if !assert.NoError(t, err, "Failed to read the metadata directory with err (%s)", err) {
		return
	}

	volumeDirEntries, err = volumeDir.Readdir(0)
	if !assert.NoError(t, err, "Failed to read the volume data directory with err (%s)", err) {
		return
	}

	t.Logf("meta directory (%s)", metaDirEntries)
	if !assert.Equal(t, len(metaDirEntries), 0, "expected metadata directory to have 1 entry and it had (%s)", len(metaDirEntries)) {
		return
	}
	t.Logf("vol directory (%s)", volumeDirEntries)
	if !assert.Equal(t, len(volumeDirEntries), 0, "expected metadata directory to have 1 entry and it had (%s)", len(volumeDirEntries)) {
		return
	}

	volToCheck, err := vs.VolumeCreate(op, testVolName, vs.Target, 0, info)
	if !assert.NoError(t, err, "Failed during call to VolumeCreate with err (%s)", err) {
		return
	}

	volumeList, err := vs.VolumesList(op)
	if !assert.NoError(t, err, "Failed during call to VolumesList with err (%s)", err) {
		return
	}

	if !assert.Equal(t, len(volumeList), 1, "Expected 1 entry in volumeList, but it had (%s)", len(volumeList)) {
		return
	}

	if !assert.Equal(t, volumeList[0].ID, volToCheck.ID, "Failed due to VolumeList returning an unexpected volume %#v when volume %#v was expected.", volumeList[0], volToCheck) {
		return
	}

	RetrievedInfo := volumeList[0].Info
	CreatedInfo := volToCheck.Info

	if !assert.Equal(t, len(RetrievedInfo), len(CreatedInfo), "Length mismatch between the created volume(%s) and the volume returned from VolumeList(%s)", len(CreatedInfo), len(RetrievedInfo)) {
		return
	}

	if !assert.Equal(t, RetrievedInfo[testInfoKey], CreatedInfo[testInfoKey], "Failed due to mismatch in metadata between the content of the Created volume(%s) and the volume return from VolumesList", CreatedInfo[testInfoKey], RetrievedInfo[testInfoKey]) {
		return
	}

	err = vs.VolumeDestroy(op, volToCheck)
	if !assert.NoError(t, err, "Failed during a call to VolumeDestroy with err (%s)", err) {
		return
	}

	volToCheckMetaDataPath := path.Join(metadataPath, volToCheck.ID)
	volToCheckDirCheck, err := os.Open(volToCheckMetaDataPath)
	defer volToCheckDirCheck.Close()
	if !assert.Error(t, err, "expected path (%s) to no longer exist after the deletion of volume (%s)", volToCheckMetaDataPath, volToCheck.ID) {
		return
	}

	volumeList, err = vs.VolumesList(op)
	if !assert.NoError(t, err, "Failed during a call to VolumesListwith err (%s)", err) {
		return
	}

	if !assert.Equal(t, len(volumeList), 0, "Expected %s volumes, VolumesList returned %s", 0, len(volumeList)) {
		return
	}

	os.Remove(volumePath)
	os.Remove(metadataDir)
}

func TestMultipleVolumes(t *testing.T) {
	mockMount := MockMount{}
	dirpath := path.Join(os.TempDir(), nfsTestDir)

	targetURL := url.URL{Path: dirpath}
	op := trace.NewOperation(context.TODO(), "TestOp")

	//Create a Volume Store
	vs, err := NewVolumeStore(op, "testStore", &targetURL, mockMount)
	if !assert.NoError(t, err, "Failed during call to NewVolumeStore with err (%s)", err) {
		return
	}

	_, err = os.Stat(dirpath)
	if !assert.NoError(t, err, "Could not find the initial volume store directory after creation of volume store. err (%s)", err) {
		return
	}

	if !assert.NotNil(t, vs, "Volume Store created with nil err, but return is also nil") {
		return
	}

	// setup volume inputs
	testVolNameOne := "test1"
	infoOne := make(map[string][]byte)
	testOneInfoKey := "junk"
	infoOne[testOneInfoKey] = make([]byte, 20)

	testVolNameTwo := "test2"
	testTwoInfoKey := "important"
	infoTwo := make(map[string][]byte)
	infoTwo[testTwoInfoKey] = []byte("42")
	testTwoInfoKeyTwo := "lessImportant"
	infoTwo[testTwoInfoKeyTwo] = []byte("41")

	testVolNameThree := "test3"
	infoThree := make(map[string][]byte)
	testThreeInfoKey := "lotsOfStuff"
	infoThree[testThreeInfoKey] = []byte("importantData")
	testThreeInfoKeyTwo := "someMoreStuff"
	infoThree[testThreeInfoKeyTwo] = []byte("maybeSomeLabels")

	//make volume one
	volOne, err := vs.VolumeCreate(op, testVolNameOne, vs.Target, 0 /*we do not use this*/, infoOne)

	if !assert.NoError(t, err, "Failed during call to VolumeCreate with err (%s)", err) {
		return
	}

	if !assert.Equal(t, testVolNameOne, volOne.ID, "expected volume ID (%s) got ID (%s)", testVolNameOne, volOne.ID) {
		return
	}

	valOne, ok := volOne.Info[testOneInfoKey]
	if !assert.True(t, ok, "TestInfoKey did not exist in the return metadata map") {
		return
	}

	if !assert.Equal(t, valOne, volOne.Info[testOneInfoKey], "TestVolOne expected to have data (%s) and (%s) was found", infoOne, valOne) {
		return
	}

	//make volume two
	volTwo, err := vs.VolumeCreate(op, testVolNameTwo, vs.Target, 0 /*we do not use this*/, infoTwo)

	if !assert.NoError(t, err, "Failed during call to VolumeCreate with err (%s)", err) {
		return
	}

	if !assert.Equal(t, testVolNameTwo, volTwo.ID, "expected volume ID (%s) got ID (%s)", testVolNameTwo, volTwo.ID) {
		return
	}

	valOne, ok = volTwo.Info[testTwoInfoKey]
	if !assert.True(t, ok, "TestInfoKey did not exist in the return metadata map") {
		return
	}

	if !assert.Equal(t, []byte("42"), valOne, "TestVolTwo expected to have data (%s) and (%s) was found", []byte("42"), valOne) {
		return
	}

	valTwo, ok := volTwo.Info[testTwoInfoKeyTwo]
	if !assert.True(t, ok, "TestTwoInfoKeyTwo did not exist in the return metadata map for volTwo") {
		return
	}

	if !assert.Equal(t, []byte("41"), valTwo, "TestVolTwo expected to have data (%s) and (%s) was found", []byte("41"), valTwo) {
		return
	}

	//make volume three
	volThree, err := vs.VolumeCreate(op, testVolNameThree, vs.Target, 0 /*we do not use this*/, infoThree)

	if !assert.NoError(t, err, "Failed during call to VolumeCreate with err (%s)", err) {
		return
	}

	if !assert.Equal(t, testVolNameThree, volThree.ID, "expected volume ID (%s) got ID (%s)", testVolNameThree, volThree.ID) {
		return
	}

	valOne, ok = volThree.Info[testThreeInfoKey]
	if !assert.True(t, ok, "TestInfoKey did not exist in the return metadata map") {
		return
	}

	if !assert.Equal(t, []byte("importantData"), valOne, "TestVolThree expected to have data (%s) and (%s) was found", []byte("importantData"), valOne) {
		return
	}

	valTwo, ok = volThree.Info[testThreeInfoKeyTwo]
	if !assert.True(t, ok, "TestThreeInfoKeyTwo did not exist in the return metadata map for volThree") {
		return
	}

	if !assert.Equal(t, []byte("maybeSomeLabels"), valTwo, "TestVolThree expected to have data (%s) and (%s) was found", []byte("maybeSomeLabels"), valTwo) {
		return
	}

	//list volumes
	volumes, err := vs.VolumesList(op)
	if !assert.NoError(t, err, "Failed during a call to VolumesList with err (%s)", err) {
		return
	}

	volCount := len(volumes)
	if !assert.Equal(t, volCount, 3, "VolumesList returned unexpected volume count. expected (%s), but received (%s) ", 3, volCount) {
		return
	}

	//check metadatas
	metadataPath := path.Join(dirpath, metadataDir)
	volumePath := path.Join(dirpath, VolumesDir)

	metaFilesDir, err := os.Open(metadataPath)
	defer metaFilesDir.Close()
	if !assert.NoError(t, err, "opening the metadata directory failed with err (%s)", err) {
		return
	}

	volumeDir, err := os.Open(volumePath)
	defer volumeDir.Close()
	if !assert.NoError(t, err, "Opening the volume directory failed with err (%s)", err) {
		return
	}

	metaDirEntries, err := metaFilesDir.Readdir(0)
	if !assert.NoError(t, err, "Failed to read the metadata directory with err (%s)", err) {
		return
	}

	volumeDirEntries, err := volumeDir.Readdir(0)
	if !assert.NoError(t, err, "Failed to read the volume data directory with err (%s)", err) {
		return
	}

	if !assert.Equal(t, len(metaDirEntries), 3, "expected metadata directory to have 1 entry and it had (%s)", len(metaDirEntries)) {
		return
	}
	if !assert.Equal(t, len(volumeDirEntries), 3, "expected metadata directory to have 1 entry and it had (%s)", len(volumeDirEntries)) {
		return
	}

	//check and individual metadata dir
	volThreeMetadataPath := path.Join(metadataPath, testVolNameThree)
	volThreeMetadataDir, err := os.Open(volThreeMetadataPath)
	defer volThreeMetadataDir.Close()
	if !assert.NoError(t, err, "Expected for path (%s) to exist, but received this error instead (%s)", volThreeMetadataPath, err) {
		return
	}

	metadataFiles, err := volThreeMetadataDir.Readdir(0)
	if !assert.Len(t, metadataFiles, 2, "Expected %s files of metadata, instead %s were found", 2, len(metadataFiles)) {
		return
	}

	volTwoMetadataPath := path.Join(metadataPath, testVolNameTwo)
	volTwoMetadataDir, err := os.Open(volTwoMetadataPath)
	defer volTwoMetadataDir.Close()
	if !assert.NoError(t, err, "Expected for path (%s) to exist, but received this error instead (%s)", volTwoMetadataPath, err) {
		return
	}

	metadataFiles, err = volTwoMetadataDir.Readdir(0)
	if !assert.Len(t, metadataFiles, 2, "Expected %s files of metadata, instead %s were found", 2, len(metadataFiles)) {
		return
	}

	volOneMetadataPath := path.Join(metadataPath, testVolNameOne)
	volOneMetadataDir, err := os.Open(volOneMetadataPath)
	defer volOneMetadataDir.Close()
	if !assert.NoError(t, err, "Expected for path (%s) to exist, but received this error instead (%s)", volOneMetadataPath, err) {
		return
	}

	metadataFiles, err = volOneMetadataDir.Readdir(0)
	if !assert.Len(t, metadataFiles, 1, "Expected %s files of metadata, instead %s were found", 2, len(metadataFiles)) {
		return
	}

	//remove volume one
	err = vs.VolumeDestroy(op, volOne)
	if !assert.NoError(t, err, "Failed during a call to VolumeDestroy with error (%s)", err) {
		return
	}

	volOneMetaDataPath := path.Join(metadataPath, volOne.ID)
	volOneDirCheck, err := os.Open(volOneMetaDataPath)
	defer volOneDirCheck.Close()
	if !assert.Error(t, err, "expected path (%s) to no longer exist after the deletion of volume (%s)", volOneMetaDataPath, volOne.ID) {
		return
	}

	//check that volume two and three exist with appropriate metadata
	volThreeMetadataPath = path.Join(metadataPath, testVolNameThree)
	volThreeMetadataDir, err = os.Open(volThreeMetadataPath)
	defer volThreeMetadataDir.Close()
	if !assert.NoError(t, err, "Expected for path (%s) to exist, but received this error instead (%s)", volThreeMetadataPath, err) {
		return
	}

	metadataFiles, err = volThreeMetadataDir.Readdir(0)
	if !assert.Len(t, metadataFiles, 2, "Expected %s files of metadata, instead %s were found", 2, len(metadataFiles)) {
		return
	}

	volTwoMetadataPath = path.Join(metadataPath, testVolNameTwo)
	volTwoMetadataDir, err = os.Open(volTwoMetadataPath)
	defer volTwoMetadataDir.Close()
	if !assert.NoError(t, err, "Expected for path (%s) to exist, but received this error instead (%s)", volTwoMetadataPath, err) {
		return
	}

	metadataFiles, err = volTwoMetadataDir.Readdir(0)
	if !assert.Len(t, metadataFiles, 2, "Expected %s files of metadata, instead %s were found", 2, len(metadataFiles)) {
		return
	}

	//remove the rest of the volumes
	err = vs.VolumeDestroy(op, volTwo)
	if !assert.NoError(t, err, "Failed during a call to VolumeDestroy with error (%s)", err) {
		return
	}

	volTwoMetaDataPath := path.Join(metadataPath, volTwo.ID)
	volTwoDirCheck, err := os.Open(volTwoMetaDataPath)
	defer volTwoDirCheck.Close()
	if !assert.Error(t, err, "expected path (%s) to no longer exist after the deletion of volume (%s)", volTwoMetaDataPath, volTwo.ID) {
		return
	}

	err = vs.VolumeDestroy(op, volThree)
	if !assert.NoError(t, err, "Failed during a call to VolumeDestroy with error (%s)", err) {
		return
	}

	volThreeMetaDataPath := path.Join(metadataPath, volThree.ID)
	volThreeDirCheck, err := os.Open(volThreeMetaDataPath)
	defer volThreeDirCheck.Close()
	if !assert.Error(t, err, "expected path (%s) to no longer exist after the deletion of volume (%s)", volThreeMetaDataPath, volThree.ID) {
		return
	}

	os.Remove(volumePath)
	os.Remove(metadataDir)
	return
}
