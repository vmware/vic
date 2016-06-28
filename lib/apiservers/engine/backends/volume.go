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
	"fmt"
	"net/http"
	"strconv"

	derr "github.com/docker/docker/errors"

	"github.com/docker/engine-api/types"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/pkg/trace"
)

//NOTE: These might be moved to a utility package once there are multiple personalities
const optsVolumeStoreKey string = "VolumeStore"
const optsCapacityKey string = "Capacity"
const dockerMetadataModelKey string = "DockerMetaData"

//Volume : struct which defines the docker personalities view of a Volume
type Volume struct {
	ProductName string
}

type volumeMetadata struct {
	Driver     string
	DriverOpts map[string]string
	Name       string
	Labels     map[string]string
}

//Volumes : docker personality implementation for VIC
func (v *Volume) Volumes(filter string) ([]*types.Volume, []string, error) {
	return nil, make([]string, 0), fmt.Errorf("%s does not implement volume.Volumes", v.ProductName)
}

//VolumeInspect : docker personality implementation for VIC
func (v *Volume) VolumeInspect(name string) (*types.Volume, error) {
	return nil, fmt.Errorf("%s does not implement volume.VolumeInspect", v.ProductName)
}

//VolumeCreate : docker personality implementation for VIC
func (v *Volume) VolumeCreate(name, driverName string, opts, labels map[string]string) (*types.Volume, error) {
	defer trace.End(trace.Begin("Volume.VolumeCreate"))
	result := &types.Volume{}

	//TODO: design a way to have better error returns.

	client := PortLayerClient()
	if client == nil {
		return nil, derr.NewErrorWithStatusCode(fmt.Errorf("Failed to get a portlayer client"), http.StatusInternalServerError)
	}

	//TODO: support having another driver besides vsphere.
	//assign the values of the model to be passed to the portlayer handler
	model, varErr := translateInputsToPortlayerRequestModel(name, driverName, opts, labels)
	if varErr != nil {
		return result, derr.NewErrorWithStatusCode(fmt.Errorf("Bad Driver Arg: %s", varErr), http.StatusBadRequest)
	}

	//TODO: setup name randomization if name == nil

	res, err := client.Storage.CreateVolume(storage.NewCreateVolumeParams().WithVolumeRequest(&model))
	if err != nil {
		return result, derr.NewErrorWithStatusCode(fmt.Errorf("Server error from Portlayer: %s", err), http.StatusInternalServerError)
	}

	result = fillDockerVolumeModel(res.Payload, labels)
	return result, nil
}

//VolumeRm : docker personality for VIC
func (v *Volume) VolumeRm(name string) error {
	defer trace.End(trace.Begin("Volume.VolumeRm"))

	client := PortLayerClient()
	if client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("Failed to get a portlayer client"), http.StatusInternalServerError)
	}

	//FIXME: check whether this is a name or a UUID. UUID expected for now.
	_, err := client.Storage.RemoveVolume(storage.NewRemoveVolumeParams().WithName(name))
	if err != nil {
		if _, ok := err.(*storage.RemoveVolumeNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("Get %s: no such volume", name))
		}
		if _, ok := err.(*storage.RemoveVolumeConflict); ok {
			return derr.NewRequestConflictError(fmt.Errorf("Volume is in use"))
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("Server error from portlayer: %s", err), http.StatusInternalServerError)
	}
	return nil
}

//Utility Functions

func fillDockerVolumeModel(volume *models.VolumeResponse, labels map[string]string) *types.Volume {
	dockerVol := types.Volume{
		Driver:     volume.Driver,
		Name:       volume.Name,
		Labels:     labels,
		Mountpoint: volume.SelfLink,
	}
	return &dockerVol
}

func validateDriverArgs(args map[string]string, model *models.VolumeRequest) error {
	//volumestore name validation
	storeName, ok := args[optsVolumeStoreKey]
	if !ok {
		storeName = "default"
	}
	model.Store = storeName

	//capacity validation
	capstr, ok := args[optsCapacityKey]
	if !ok {
		model.Capacity = -1
		return nil
	}
	capacity, convErr := strconv.ParseInt(capstr, 10, 64)
	if convErr != nil {
		model.Capacity = -1
		return fmt.Errorf("Capacity must be an integer value. The unit is GB.: %s", convErr)
	}
	model.Capacity = int64(capacity)
	return nil
}

func translateInputsToPortlayerRequestModel(name, driverName string, opts, labels map[string]string) (models.VolumeRequest, error) {
	model := models.VolumeRequest{
		Driver:     driverName,
		DriverArgs: opts,
		Name:       name,
	}
	metadata := createVolumeMetadata(&model, labels)
	model.Metadata = make(map[string]string)
	model.Metadata["dockerMetaData"] = metadata
	if err := validateDriverArgs(opts, &model); err != nil {
		return model, err
	}
	return model, nil
}

func createVolumeMetadata(model *models.VolumeRequest, labels map[string]string) string {
	metadata := volumeMetadata{
		Driver:     model.Driver,
		DriverOpts: model.DriverArgs,
		Name:       model.Name,
		Labels:     labels,
	}
	result, _ := json.Marshal(metadata)
	return string(result)
}
