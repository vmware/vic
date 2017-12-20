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

package backends

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	derr "github.com/docker/docker/api/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/go-units"
	"github.com/google/uuid"

	vicfilter "github.com/vmware/vic/lib/apiservers/engine/backends/filter"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/pkg/trace"
)

// NOTE: FIXME: These might be moved to a utility package once there are multiple personalities
const (
	OptsVolumeStoreKey     string = "volumestore"
	OptsCapacityKey        string = "capacity"
	dockerMetadataModelKey string = "DockerMetaData"
	DefaultVolumeDriver    string = "vsphere"
)

// define a set (whitelist) of valid driver opts keys for command line argument validation
var validDriverOptsKeys = map[string]struct{}{
	OptsVolumeStoreKey:    {},
	OptsCapacityKey:       {},
	DriverArgFlagKey:      {},
	DriverArgContainerKey: {},
	DriverArgImageKey:     {},
}

// Volume drivers currently supported. "local" is the default driver supplied by the client
// and is equivalent to "vsphere" for our implementation.
var supportedVolDrivers = map[string]struct{}{
	"vsphere": {},
	"local":   {},
}

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

// acceptedVolumeFilters are volume filters that are supported by VIC
var acceptedVolumeFilters = map[string]bool{
	"dangling": true,
	"name":     true,
	"driver":   true,
	"label":    true,
}

var errPortlayerClient = fmt.Errorf("failed to get a portlayer client")

func NewVolumeBackend() *Volume {
	return &Volume{}
}

// Volumes docker personality implementation for VIC
func (v *Volume) Volumes(filter string) ([]*types.Volume, []string, error) {
	defer trace.End(trace.Begin("Volume.Volumes"))
	var volumes []*types.Volume

	client := PortLayerClient()
	if client == nil {
		return nil, nil, VolumeInternalServerError(errPortlayerClient)
	}

	res, err := client.Storage.ListVolumes(storage.NewListVolumesParamsWithContext(ctx).WithFilterString(&filter))
	if err != nil {
		switch err := err.(type) {
		case *storage.ListVolumesInternalServerError:
			return nil, nil, VolumeInternalServerError(fmt.Errorf("error from portlayer server: %s", err.Payload.Message))
		case *storage.ListVolumesDefault:
			return nil, nil, VolumeInternalServerError(fmt.Errorf("error from portlayer server: %s", err.Payload.Message))
		default:
			return nil, nil, VolumeInternalServerError(fmt.Errorf("error from portlayer server: %s", err.Error()))
		}
	}

	volumeResponses := res.Payload

	// Parse and validate filters
	volumeFilters, err := filters.FromParam(filter)
	if err != nil {
		return nil, nil, VolumeInternalServerError(err)
	}
	volFilterContext, err := vicfilter.ValidateVolumeFilters(volumeFilters, acceptedVolumeFilters, nil)
	if err != nil {
		return nil, nil, VolumeInternalServerError(err)
	}

	// joinedVolumes stores names of volumes that are joined to a container
	// and is used while filtering the output by dangling (dangling=true should
	// return volumes that are not attached to a container)
	joinedVolumes := make(map[string]struct{})
	if volumeFilters.Include("dangling") {
		// If the dangling filter is specified, gather required items beforehand
		joinedVolumes, err = fetchJoinedVolumes()
		if err != nil {
			return nil, nil, VolumeInternalServerError(err)
		}
	}

	log.Infoln("volumes found:")
	for _, vol := range volumeResponses {
		log.Infof("%s", vol.Name)

		volumeMetadata, err := extractDockerMetadata(vol.Metadata)
		if err != nil {
			return nil, nil, VolumeInternalServerError(fmt.Errorf("error unmarshalling docker metadata: %s", err))
		}

		// Set fields needed for filtering the output
		volFilterContext.Name = vol.Name
		volFilterContext.Driver = vol.Driver
		_, volFilterContext.Joined = joinedVolumes[vol.Name]
		volFilterContext.Labels = volumeMetadata.Labels

		// Include the volume in the output if it meets the filtering criteria
		filterAction := vicfilter.IncludeVolume(volumeFilters, volFilterContext)
		if filterAction == vicfilter.IncludeAction {
			volume := NewVolumeModel(vol, volumeMetadata.Labels)
			volumes = append(volumes, volume)
		}
	}

	return volumes, nil, nil
}

// fetchJoinedVolumes obtains all containers from the portlayer and returns a map with all
// volumes that are joined to at least one container.
func fetchJoinedVolumes() (map[string]struct{}, error) {
	conts, err := allContainers()
	if err != nil {
		return nil, VolumeInternalServerError(err)
	}

	joinedVolumes := make(map[string]struct{})
	var s struct{}
	for i := range conts {
		for _, vol := range conts[i].VolumeConfig {
			joinedVolumes[vol.Name] = s
		}
	}

	return joinedVolumes, nil
}

// allContainers obtains all containers from the portlayer, akin to `docker ps -a`.
func allContainers() ([]*models.ContainerInfo, error) {
	client := PortLayerClient()
	if client == nil {
		return nil, errPortlayerClient
	}

	all := true
	cons, err := client.Containers.GetContainerList(containers.NewGetContainerListParamsWithContext(ctx).WithAll(&all))
	if err != nil {
		return nil, err
	}

	return cons.Payload, nil
}

// VolumeInspect : docker personality implementation for VIC
func (v *Volume) VolumeInspect(name string) (*types.Volume, error) {
	defer trace.End(trace.Begin(name))

	client := PortLayerClient()
	if client == nil {
		return nil, VolumeInternalServerError(errPortlayerClient)
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
			return nil, VolumeInternalServerError(fmt.Errorf("error from portlayer server: %s", err.Error()))
		}
	}

	volumeMetadata, err := extractDockerMetadata(res.Payload.Metadata)
	if err != nil {
		return nil, VolumeInternalServerError(fmt.Errorf("error unmarshalling docker metadata: %s", err))
	}
	volume := NewVolumeModel(res.Payload, volumeMetadata.Labels)

	return volume, nil
}

// volumeCreate issues a CreateVolume request to the portlayer
func (v *Volume) volumeCreate(name, driverName string, volumeData, labels map[string]string) (*types.Volume, error) {
	defer trace.End(trace.Begin(""))
	result := &types.Volume{}

	client := PortLayerClient()
	if client == nil {
		return nil, errPortlayerClient
	}

	if name == "" {
		name = uuid.New().String()
	}

	// TODO: support having another driver besides vsphere.
	// assign the values of the model to be passed to the portlayer handler
	req, varErr := newVolumeCreateReq(name, driverName, volumeData, labels)
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
func (v *Volume) VolumeCreate(name, driverName string, volumeData, labels map[string]string) (*types.Volume, error) {
	defer trace.End(trace.Begin("Volume.VolumeCreate"))

	result, err := v.volumeCreate(name, driverName, volumeData, labels)
	if err != nil {
		switch err := err.(type) {
		case *storage.CreateVolumeConflict:
			return result, VolumeInternalServerError(fmt.Errorf("A volume named %s already exists. Choose a different volume name.", name))

		case *storage.CreateVolumeNotFound:
			return result, VolumeInternalServerError(fmt.Errorf("No volume store named (%s) exists", volumeStore(volumeData)))

		case *storage.CreateVolumeInternalServerError:
			// FIXME: right now this does not return an error model...
			return result, VolumeInternalServerError(fmt.Errorf("%s", err.Error()))

		case *storage.CreateVolumeDefault:
			return result, VolumeInternalServerError(fmt.Errorf("%s", err.Payload.Message))

		default:
			return result, VolumeInternalServerError(fmt.Errorf("%s", err))
		}
	}

	return result, nil
}

// VolumeRm : docker personality for VIC
func (v *Volume) VolumeRm(name string, force bool) error {
	defer trace.End(trace.Begin(name))

	client := PortLayerClient()
	if client == nil {
		return VolumeInternalServerError(errPortlayerClient)
	}

	_, err := client.Storage.RemoveVolume(storage.NewRemoveVolumeParamsWithContext(ctx).WithName(name))
	if err != nil {

		switch err := err.(type) {
		case *storage.RemoveVolumeNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf("Get %s: no such volume", name))

		case *storage.RemoveVolumeConflict:
			return derr.NewRequestConflictError(fmt.Errorf(err.Payload.Message))

		case *storage.RemoveVolumeInternalServerError:
			return VolumeInternalServerError(fmt.Errorf("Server error from portlayer: %s", err.Payload.Message))
		default:
			return VolumeInternalServerError(fmt.Errorf("Server error from portlayer: %s", err))
		}
	}
	return nil
}

func (v *Volume) VolumesPrune(pruneFilters filters.Args) (*types.VolumesPruneReport, error) {
	return nil, fmt.Errorf("%s does not yet implement VolumesPrune", ProductName())
}

type volumeMetadata struct {
	Driver        string
	DriverOpts    map[string]string
	Name          string
	Labels        map[string]string
	AttachHistory []string
	Image         string
}

func createVolumeMetadata(req *models.VolumeRequest, driverargs, labels map[string]string) (string, error) {
	metadata := volumeMetadata{
		Driver:        req.Driver,
		DriverOpts:    req.DriverArgs,
		Name:          req.Name,
		Labels:        labels,
		AttachHistory: []string{driverargs[DriverArgContainerKey]},
		Image:         driverargs[DriverArgImageKey],
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

func newVolumeCreateReq(name, driverName string, volumeData, labels map[string]string) (*models.VolumeRequest, error) {
	if _, ok := supportedVolDrivers[driverName]; !ok {
		return nil, fmt.Errorf("error looking up volume plugin %s: plugin not found", driverName)
	}

	if !volumeNameRegex.Match([]byte(name)) && name != "" {
		return nil, fmt.Errorf("volume name %q includes invalid characters, only \"[a-zA-Z0-9][a-zA-Z0-9_.-]\" are allowed", name)
	}

	req := &models.VolumeRequest{
		Driver:     driverName,
		DriverArgs: volumeData,
		Name:       name,
		Metadata:   make(map[string]string),
	}

	metadata, err := createVolumeMetadata(req, volumeData, labels)
	if err != nil {
		return nil, err
	}

	req.Metadata[dockerMetadataModelKey] = metadata

	if err := validateDriverArgs(volumeData, req); err != nil {
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

func normalizeDriverArgs(args map[string]string) error {
	// normalize keys to lowercase & validate them
	for k, val := range args {
		lowercase := strings.ToLower(k)

		if _, ok := validDriverOptsKeys[lowercase]; !ok {
			return fmt.Errorf("%s is not a supported option", k)
		}

		if strings.Compare(lowercase, k) != 0 {
			delete(args, k)
			args[lowercase] = val
		}
	}
	return nil
}

func validateDriverArgs(args map[string]string, req *models.VolumeRequest) error {
	if err := normalizeDriverArgs(args); err != nil {
		return err
	}

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
