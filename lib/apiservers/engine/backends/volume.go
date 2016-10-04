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

package backends

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	derr "github.com/docker/docker/errors"

	"regexp"
	"strconv"

	"github.com/docker/engine-api/types"
	"github.com/docker/go-units"
	"github.com/google/uuid"

	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/pkg/trace"
)

// NOTE: FIXME: These might be moved to a utility package once there are multiple personalities
const (
	OptsVolumeStoreKey     string = "VolumeStore"
	OptsCapacityKey        string = "Capacity"
	dockerMetadataModelKey string = "DockerMetaData"
)

//Validation pattern for Volume Names
var volumeNameRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_.-]*$")

func NewVolumeModel(volume *models.VolumeResponse, labels map[string]string) *types.Volume {
	return &types.Volume{
		Driver:     volume.Driver,
		Name:       volume.Name,
		Labels:     labels,
		Mountpoint: volume.Label,
	}
}

// Volume which defines the docker personalities view of a Volume
type Volume struct {
}

// Volumes docker personality implementation for VIC
func (v *Volume) Volumes(filter string) ([]*types.Volume, []string, error) {
	defer trace.End(trace.Begin("Volume.Volumes"))
	var volumes []*types.Volume

	client := PortLayerClient()
	if client == nil {
		return nil, nil, derr.NewErrorWithStatusCode(fmt.Errorf("Failed to get a portlayer client"), http.StatusInternalServerError)
	}

	res, err := client.Storage.ListVolumes(storage.NewListVolumesParamsWithContext(ctx).WithFilterString(&filter))
	if err != nil {
		switch err := err.(type) {
		case *storage.ListVolumesInternalServerError:
			return nil, nil, derr.NewErrorWithStatusCode(fmt.Errorf("error from portlayer server: %s", err.Payload.Message), http.StatusInternalServerError)
		case *storage.ListVolumesDefault:
			return nil, nil, derr.NewErrorWithStatusCode(fmt.Errorf("error from portlayer server: %s", err.Payload.Message), http.StatusInternalServerError)
		default:
			return nil, nil, derr.NewErrorWithStatusCode(fmt.Errorf("error from portlayer server: %s", err.Error()), http.StatusInternalServerError)
		}
	}

	volumeResponses := res.Payload

	log.Infoln("volumes found: ")
	for _, vol := range volumeResponses {
		log.Infof("%s", vol.Name)
		volumeMetadata, err := extractDockerMetadata(vol.Metadata)
		if err != nil {
			return nil, nil, fmt.Errorf("error unmarshalling docker metadata: %s", err)
		}
		volume := NewVolumeModel(vol, volumeMetadata.Labels)
		volumes = append(volumes, volume)
	}
	return volumes, nil, nil
}

// VolumeInspect : docker personality implementation for VIC
func (v *Volume) VolumeInspect(name string) (*types.Volume, error) {
	defer trace.End(trace.Begin(name))

	client := PortLayerClient()
	if client == nil {
		return nil, fmt.Errorf("failed to get a portlayer client")
	}

	if name == "" {
		return nil, nil
	}

	param := storage.NewGetVolumeParamsWithContext(ctx).WithName(name)
	res, err := client.Storage.GetVolume(param)
	if err != nil {
		switch err := err.(type) {
		case *storage.GetVolumeNotFound:
			return nil, VolumeNotFoundError(name)
		default:
			return nil, derr.NewErrorWithStatusCode(fmt.Errorf("error from portlayer server: %s", err.Error()), http.StatusInternalServerError)
		}
	}

	volumeMetadata, err := extractDockerMetadata(res.Payload.Metadata)
	if err != nil {
		return nil, derr.NewErrorWithStatusCode(fmt.Errorf("error unmarshalling docker metadata: %s", err), http.StatusInternalServerError)
	}
	volume := NewVolumeModel(res.Payload, volumeMetadata.Labels)

	return volume, nil
}

// volumeCreate issues a CreateVolume request to the portlayer
func (v *Volume) volumeCreate(name, driverName string, driverArgs, labels map[string]string) (*types.Volume, error) {
	defer trace.End(trace.Begin(""))
	result := &types.Volume{}

	client := PortLayerClient()
	if client == nil {
		return nil, fmt.Errorf("failed to get a portlayer client")
	}

	if name == "" {
		name = uuid.New().String()
	}

	// TODO: support having another driver besides vsphere.
	// assign the values of the model to be passed to the portlayer handler
	req, varErr := newVolumeCreateReq(name, driverName, driverArgs, labels)
	if varErr != nil {
		return result, varErr
	}
	log.Infof("Finalized model for volume create request to portlayer: %#v", req)

	res, err := client.Storage.CreateVolume(storage.NewCreateVolumeParamsWithContext(ctx).WithVolumeRequest(req))
	if err != nil {
		return result, err
	}
	result = NewVolumeModel(res.Payload, labels)
	return result, nil
}

// VolumeCreate : docker personality implementation for VIC
func (v *Volume) VolumeCreate(name, driverName string, driverArgs, labels map[string]string) (*types.Volume, error) {
	defer trace.End(trace.Begin("Volume.VolumeCreate"))

	result, err := v.volumeCreate(name, driverName, driverArgs, labels)
	if err != nil {
		switch err := err.(type) {
		case *storage.CreateVolumeConflict:
			return result, derr.NewErrorWithStatusCode(fmt.Errorf("A volume named %s already exists. Choose a different volume name.", name), http.StatusInternalServerError)

		case *storage.CreateVolumeNotFound:
			return result, derr.NewErrorWithStatusCode(fmt.Errorf("No volume store named (%s) exists", volumeStore(driverArgs)), http.StatusInternalServerError)

		case *storage.CreateVolumeInternalServerError:
			// FIXME: right now this does not return an error model...
			return result, derr.NewErrorWithStatusCode(fmt.Errorf("%s", err.Error()), http.StatusInternalServerError)

		case *storage.CreateVolumeDefault:
			return result, derr.NewErrorWithStatusCode(fmt.Errorf("%s", err.Payload.Message), http.StatusInternalServerError)

		default:
			return result, derr.NewErrorWithStatusCode(fmt.Errorf("%s", err), http.StatusInternalServerError)
		}
	}

	return result, nil
}

// VolumeRm : docker personality for VIC
func (v *Volume) VolumeRm(name string) error {
	defer trace.End(trace.Begin("Volume.VolumeRm"))

	client := PortLayerClient()
	if client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("Failed to get a portlayer client"), http.StatusInternalServerError)
	}

	// FIXME: check whether this is a name or a UUID. UUID expected for now.
	_, err := client.Storage.RemoveVolume(storage.NewRemoveVolumeParamsWithContext(ctx).WithName(name))
	if err != nil {

		switch err := err.(type) {
		case *storage.RemoveVolumeNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf("Get %s: no such volume", name))

		case *storage.RemoveVolumeConflict:
			return derr.NewRequestConflictError(fmt.Errorf(err.Payload.Message))

		case *storage.RemoveVolumeInternalServerError:
			return derr.NewErrorWithStatusCode(fmt.Errorf("Server error from portlayer: %s", err.Payload.Message), http.StatusInternalServerError)
		default:
			return derr.NewErrorWithStatusCode(fmt.Errorf("Server error from portlayer: %s", err), http.StatusInternalServerError)
		}
	}
	return nil
}

type volumeMetadata struct {
	Driver        string
	DriverOpts    map[string]string
	Name          string
	Labels        map[string]string
	AttachHistory []string
	Image         string
}

func createVolumeMetadata(req *models.VolumeRequest, labels map[string]string, Container, Image string) (string, error) {
	metadata := volumeMetadata{
		Driver:     req.Driver,
		DriverOpts: req.DriverArgs,
		Name:       req.Name,
		Labels:     labels,
	}
	result, err := json.Marshal(metadata)
	return string(result), err
}

// Unmarshal the docker metadata using the docker metadata key.  The docker
// metadatakey.  We stash the vals we know about in that map with that key.
func extractDockerMetadata(metadataMap map[string]string) (*volumeMetadata, error) {
	v, ok := metadataMap[dockerMetadataModelKey]
	if !ok {
		return nil, fmt.Errorf("metadata %s missing", dockerMetadataModelKey)
	}

	result := &volumeMetadata{}
	err := json.Unmarshal([]byte(v), result)
	return result, err
}

// Utility Functions

func newVolumeCreateReq(name, driverName string, driverArgs, labels map[string]string) (*models.VolumeRequest, error) {
	defaultDriver := driverName == "local"
	vsphereDriver := driverName == "vsphere"

	if !defaultDriver && !vsphereDriver {
		return nil, fmt.Errorf("Error looking up volume plugin %s: plugin not found", driverName)
	}

	if !volumeNameRegex.Match([]byte(name)) && name != "" {
		return nil, fmt.Errorf("volume name %q includes invalid characters, only \"[a-zA-Z0-9][a-zA-Z0-9_.-]\" are allowed", name)
	}

	req := &models.VolumeRequest{
		Driver:     driverName,
		DriverArgs: driverArgs,
		Name:       name,
		Metadata:   make(map[string]string),
	}

	metadata, err := createVolumeMetadata(req, labels)
	if err != nil {
		return nil, err
	}

	req.Metadata[dockerMetadataModelKey] = metadata

	if err := validateDriverArgs(driverArgs, req); err != nil {
		return nil, fmt.Errorf("bad driver value - %s", err)
	}

	return req, nil
}

// volumeStore returns the value of the optional volume store param specified in the CLI.
func volumeStore(args map[string]string) string {
	storeName, ok := args[OptsVolumeStoreKey]
	if !ok {
		return "default"
	}
	return storeName
}

func validateDriverArgs(args map[string]string, req *models.VolumeRequest) error {
	// volumestore name validation
	req.Store = volumeStore(args)

	// capacity validation
	capstr, ok := args[OptsCapacityKey]
	if !ok {
		req.Capacity = -1
		return nil
	}

	//check if it is just a numerical value
	capacity, err := strconv.ParseInt(capstr, 10, 64)
	if err == nil {
		//input has no units in this case.
		if capacity < 1 {
			return fmt.Errorf("Invalid size: %s", capstr)
		}
		req.Capacity = capacity
		return nil
	}

	capacity, err = units.FromHumanSize(capstr)
	if err != nil {
		return err
	}

	if capacity < 1 {
		return fmt.Errorf("Capacity value too large: %s", capstr)
	}

	req.Capacity = int64(capacity) / int64(units.MB)
	return nil
}
