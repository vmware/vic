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

package vicbackends

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
)

func TestFillDockerVolume(t *testing.T) {
	testResponse := &models.VolumeResponse{
		Driver:   "vsphere",
		Name:     "Test Volume",
		SelfLink: "Test SelfLink",
	}
	testLabels := make(map[string]string)
	testLabels["TestMeta"] = "custom info about my volume"

	dockerVolume := fillDockerVolumeModel(testResponse, testLabels)

	assert.Equal(t, "vsphere", dockerVolume.Driver)
	assert.Equal(t, "Test Volume", dockerVolume.Name)
	assert.Equal(t, "Test SelfLink", dockerVolume.Mountpoint)
	assert.Equal(t, "custom info about my volume", dockerVolume.Labels["TestMeta"])
}

func TestTranslatVolumeRequestModel(t *testing.T) {
	testLabels := make(map[string]string)
	testLabels["TestMeta"] = "custom info about my volume"

	testDriverArgs := make(map[string]string)
	testDriverArgs["testArg"] = "important driver stuff"
	testDriverArgs[optsVolumeStoreKey] = "testStore"
	testDriverArgs[optsCapacityKey] = "12"

	testRequest, err := translateInputsToPortlayerRequestModel("testName", "vsphere", testDriverArgs, testLabels)

	assert.Equal(t, "testName", testRequest.Name)
	assert.Equal(t, "important driver stuff", testRequest.DriverArgs["testArg"])
	assert.Equal(t, "testStore", testRequest.Store)
	assert.Equal(t, "vsphere", testRequest.Driver)
	assert.Equal(t, int64(12), testRequest.Capacity)
	assert.Equal(t, "important driver stuff", testRequest.DriverArgs["testArg"])
	testMetaDataString := createVolumeMetadata(&testRequest, testLabels)
	assert.Equal(t, testMetaDataString, testRequest.Metadata["dockerMetaData"])
	assert.Nil(t, err)
}

func TestCreateVolumeMetada(t *testing.T) {
	testDriverOpts := make(map[string]string)
	testDriverOpts["TestArg"] = "test"
	testModel := models.VolumeRequest{
		Driver:     "vsphere",
		DriverArgs: testDriverOpts,
		Name:       "testModel",
	}
	testLabels := make(map[string]string)
	testLabels["TestMeta"] = "custom info about my volume"

	testMetadataString := createVolumeMetadata(&testModel, testLabels)

	volumeMetadata := volumeMetadata{}
	json.Unmarshal([]byte(testMetadataString), &volumeMetadata)

	assert.Equal(t, testModel.Driver, volumeMetadata.Driver)
	assert.Equal(t, testModel.Name, volumeMetadata.Name)
	assert.Equal(t, testLabels["TestMeta"], volumeMetadata.Labels["TestMeta"])
	assert.Equal(t, testLabels["TestArg"], volumeMetadata.DriverOpts["testArg"])
}

func TestValidateDriverArgs(t *testing.T) {
	testMap := make(map[string]string)
	testStore := "Mystore"
	testCap := int64(12)
	testBadCap := "This is not valid!"
	testModel := models.VolumeRequest{
		Driver:     "vsphere",
		DriverArgs: testMap,
		Name:       "testModel",
	}

	err := validateDriverArgs(testMap, &testModel)
	assert.Equal(t, "default", testModel.Store)
	assert.Equal(t, int64(-1), testModel.Capacity)
	assert.Nil(t, err)

	testMap[optsVolumeStoreKey] = testStore
	testMap[optsCapacityKey] = strconv.FormatInt(testCap, 10)
	err = validateDriverArgs(testMap, &testModel)
	assert.Equal(t, testStore, testModel.Store)
	assert.Equal(t, testCap, testModel.Capacity)
	assert.Nil(t, err)

	testMap[optsCapacityKey] = testBadCap
	err = validateDriverArgs(testMap, &testModel)
	assert.Equal(t, testStore, testModel.Store)
	assert.Equal(t, int64(-1), testModel.Capacity)
	assert.NotNil(t, err)

	testMap[optsCapacityKey] = strconv.FormatInt(testCap, 10)
	delete(testMap, optsVolumeStoreKey)
	err = validateDriverArgs(testMap, &testModel)
	assert.Equal(t, "default", testModel.Store)
	assert.Equal(t, int64(12), testModel.Capacity)
	assert.NotNil(t, err)
}
