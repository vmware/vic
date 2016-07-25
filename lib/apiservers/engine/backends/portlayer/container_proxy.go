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

package portlayer

//****
// container_proxy.go
//
// Contains all code that touches the portlayer for container operations and all
// code that converts swagger based returns to docker personality backend structs.
// The goal is to make the backend code that implements the docker engine-api
// interfaces be as simple as possible and contain no swagger or portlayer code.
//
// Rule for code to be in here:
// 1. touches VIC portlayer
// 2. converts swagger to docker engine-api structs
// 3. errors MUST be docker engine-api compatible errors.  DO NOT return arbitrary errors!
//		- Do NOT return portlayer errors
//		- Do NOT return fmt.Errorf()
//		- Do NOT return errors.New()
//		- Please USE the aliased docker error package 'derr'

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/google/uuid"

	derr "github.com/docker/docker/errors"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/strslice"

	"github.com/vmware/vic/lib/apiservers/engine/backends/endpoint"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/sys"
)

type VicContainerProxy interface {
	CreateContainerHandle(imageID string, config types.ContainerCreateConfig) (string, string, error)
	AddContainerToScope(handle string, config types.ContainerCreateConfig) (string, error)
	AddVolumesToContainer(handle string, config types.ContainerCreateConfig) (string, error)
	CommitContainerHandle(handle, imageID string) error
}

type ContainerProxy struct {
	client *client.PortLayer
}

type volumeFields struct {
	ID    string
	Dest  string
	Flags string
}

func NewContainerProxy(plClient *client.PortLayer) *ContainerProxy {
	return &ContainerProxy{client: plClient}
}

// containerCreateHandle() creates a new VIC container by calling the portlayer
//
// returns:
// 	(containerID, containerHandle, error)
func (c *ContainerProxy) CreateContainerHandle(imageID string, config types.ContainerCreateConfig) (string, string, error) {
	defer trace.End(trace.Begin("ContainerProxy.containerCreateHandle"))

	if c.client == nil {
		return "", "",
			derr.NewErrorWithStatusCode(fmt.Errorf("ContainerProxy.CreateContainerHandle failed to create a portlayer client"),
				http.StatusInternalServerError)
	}

	if imageID == "" {
		return "", "",
			derr.NewRequestNotFoundError(fmt.Errorf("No image specified"))
	}

	// Call the Exec port layer to create the container
	host, err := sys.UUID()
	if err != nil {
		return "", "",
			derr.NewErrorWithStatusCode(fmt.Errorf("ContainerProxy.CreateContainerHandle got unexpected error getting VCH UUID"),
				http.StatusInternalServerError)
	}

	plCreateParams := dockerContainerCreateParamsToPortlayer(config, imageID, host)
	createResults, err := c.client.Containers.Create(plCreateParams)
	if err != nil {
		if _, ok := err.(*containers.CreateNotFound); ok {
			err = fmt.Errorf("No such image: %s", imageID)
			log.Errorf(err.Error())
			return "", "", derr.NewRequestNotFoundError(err)
		}

		// If we get here, most likely something went wrong with the port layer API server
		return "", "",
			derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
	}

	id := createResults.Payload.ID
	h := createResults.Payload.Handle

	return id, h, nil
}

// AddContainerToScope() adds a container, referenced by handle, to a scope.
// If an error is return, the returned handle should not be used.
//
// returns:
//	modified handle
func (c *ContainerProxy) AddContainerToScope(handle string, config types.ContainerCreateConfig) (string, error) {
	defer trace.End(trace.Begin("ContainerProxy.AddContainerToScope"))

	if c.client == nil {
		return "",
			derr.NewErrorWithStatusCode(fmt.Errorf("ContainerProxy.AddContainerToScope failed to create a portlayer client"),
				http.StatusInternalServerError)
	}

	log.Debugf("Network Configuration Section - Container Create")
	// configure networking
	netConf := toModelsNetworkConfig(config)
	if netConf != nil {
		addContRes, err := c.client.Scopes.AddContainer(scopes.NewAddContainerParams().
			WithScope(netConf.NetworkName).
			WithConfig(&models.ScopesAddContainerConfig{
				Handle:        handle,
				NetworkConfig: netConf,
			}))

		if err != nil {
			log.Errorf("ContainerProxy.AddContainerToScope: Scopes error: %s", err.Error())
			return handle, derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
		}

		defer func() {
			if err == nil {
				return
			}
			// roll back the AddContainer call
			if _, err2 := c.client.Scopes.RemoveContainer(scopes.NewRemoveContainerParams().WithHandle(handle).WithScope(netConf.NetworkName)); err2 != nil {
				log.Warnf("could not roll back container add: %s", err2)
			}
		}()

		handle = addContRes.Payload
	}

	return handle, nil
}

// AddVolumesToContainer() adds volumes to a container, referenced by handle.
// If an error is return, the returned handle should not be used.
//
// returns:
//	modified handle
func (c *ContainerProxy) AddVolumesToContainer(handle string, config types.ContainerCreateConfig) (string, error) {
	defer trace.End(trace.Begin("ContainerProxy.AddVolumesToContainer"))

	if c.client == nil {
		return "",
			derr.NewErrorWithStatusCode(fmt.Errorf("ContainerProxy.AddVolumesToContainer failed to create a portlayer client"),
				http.StatusInternalServerError)
	}

	//Volume Attachment Section
	log.Debugf("ContainerProxy.AddVolumesToContainer - VolumeSection")
	log.Debugf("Raw Volume arguments : binds:  %#v : volumes : %#v", config.HostConfig.Binds, config.Config.Volumes)
	var joinList []volumeFields
	var err error

	joinList, err = processAnonymousVolumes(&handle, config.Config.Volumes, c.client)
	if err != nil {
		return handle, derr.NewErrorWithStatusCode(fmt.Errorf("%s", err), http.StatusBadRequest)
	}

	volumeSubset, err := processSpecifiedVolumes(config.HostConfig.Binds)
	if err != nil {
		return handle, derr.NewErrorWithStatusCode(fmt.Errorf("%s", err), http.StatusBadRequest)
	}
	joinList = append(joinList, volumeSubset...)

	for _, fields := range joinList {
		flags := make(map[string]string)
		//NOTE: for now we are passing the flags directly through. This is NOT SAFE and only a stop gap.
		flags["Mode"] = fields.Flags
		joinParams := storage.NewVolumeJoinParams().WithJoinArgs(&models.VolumeJoinConfig{
			Flags:     flags,
			Handle:    handle,
			MountPath: fields.Dest,
		}).WithName(fields.ID)

		res, err := c.client.Storage.VolumeJoin(joinParams)
		if err != nil {
			switch err := err.(type) {
			case *storage.VolumeJoinInternalServerError:
				return handle, derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)
			case *storage.VolumeJoinDefault:
				return handle, derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)
			default:
				return handle, derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
			}
		}

		handle = res.Payload
	}

	return handle, nil
}

// CommitContainerHandle() commits any changes to container handle.
//
func (c *ContainerProxy) CommitContainerHandle(handle, imageID string) error {
	defer trace.End(trace.Begin("ContainerProxy.CommitContainerHandle"))

	if c.client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("ContainerProxy.CommitContainerHandle failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	_, err := c.client.Containers.Commit(containers.NewCommitParams().WithHandle(handle))
	if err != nil {
		err = fmt.Errorf("No such image: %s", imageID)
		log.Errorf("%s", err.Error())
		// FIXME: Containers.Commit returns more errors than it's swagger spec says.
		// When no image exist, it also sends back non swagger errors.  We should fix
		// this once Commit returns correct error codes.
		return derr.NewRequestNotFoundError(err)
	}

	return nil
}

//----------
// Utility Functions
//----------

func dockerContainerCreateParamsToPortlayer(cc types.ContainerCreateConfig, layerID string, imageStore string) *containers.CreateParams {
	config := &models.ContainerCreateConfig{}

	// Image
	config.Image = new(string)
	*config.Image = layerID

	// Repo Requested
	config.RepoName = new(string)
	*config.RepoName = cc.Config.Image

	var path string
	var args []string

	// Expand cmd into entrypoint and args
	cmd := strslice.StrSlice(cc.Config.Cmd)
	if len(cc.Config.Entrypoint) != 0 {
		path, args = cc.Config.Entrypoint[0], append(cc.Config.Entrypoint[1:], cmd...)
	} else {
		path, args = cmd[0], cmd[1:]
	}

	//copy friendly name
	config.Name = new(string)
	*config.Name = cc.Name

	// copy the path
	config.Path = new(string)
	*config.Path = path

	// copy the args
	config.Args = make([]string, len(args))
	copy(config.Args, args)

	// copy the env array
	config.Env = make([]string, len(cc.Config.Env))
	copy(config.Env, cc.Config.Env)

	// image store
	config.ImageStore = &models.ImageStore{Name: imageStore}

	// network
	config.NetworkDisabled = new(bool)
	*config.NetworkDisabled = cc.Config.NetworkDisabled

	// working dir
	config.WorkingDir = new(string)
	*config.WorkingDir = cc.Config.WorkingDir

	// tty
	config.Tty = new(bool)
	*config.Tty = cc.Config.Tty

	log.Debugf("dockerContainerCreateParamsToPortlayer = %+v", config)

	return containers.NewCreateParams().WithCreateConfig(config)
}

func toModelsNetworkConfig(cc types.ContainerCreateConfig) *models.NetworkConfig {
	if cc.Config.NetworkDisabled {
		return nil
	}

	nc := &models.NetworkConfig{
		NetworkName: cc.HostConfig.NetworkMode.NetworkName(),
	}
	if cc.NetworkingConfig != nil {
		log.Debugf("EndpointsConfig: %#v", cc.NetworkingConfig)

		es, ok := cc.NetworkingConfig.EndpointsConfig[nc.NetworkName]
		if ok {
			if es.IPAMConfig != nil {
				nc.Address = &es.IPAMConfig.IPv4Address
			}

			// Docker copies Links to NetworkConfig only if it is a UserDefined network, handle that
			// https://github.com/docker/docker/blame/master/runconfig/opts/parse.go#L598
			if !cc.HostConfig.NetworkMode.IsUserDefined() && len(cc.HostConfig.Links) > 0 {
				es.Links = make([]string, len(cc.HostConfig.Links))
				copy(es.Links, cc.HostConfig.Links)
			}
			// Pass Links and Aliases to PL
			nc.Aliases = endpoint.Alias(es)

		}
	}

	nc.Ports = make([]string, len(cc.HostConfig.PortBindings))
	i := 0
	for p := range cc.HostConfig.PortBindings {
		nc.Ports[i] = string(p)
		i++
	}

	return nc
}

//This function is used to turn any call from docker create -v <stuff> into a volumeFields object.
//the -v has 3 forms. 1: -v <anonymouse mount path>, -v <Volume Name>:<Destination Mount Path>, and -v <Volume Name>:<Destination Mount Path>:<mount flags>
func processVolumeParam(volString string) (volumeFields, error) {
	volumeStrings := strings.Split(volString, ":")
	fields := volumeFields{}

	//This switch determines which type of -v was invoked.
	switch len(volumeStrings) {
	case 1:
		VolumeID, err := uuid.NewUUID()
		if err != nil {
			return volumeFields{}, nil
		}
		fields.ID = VolumeID.String()
		fields.Dest = volumeStrings[0]
		fields.Flags = "rw"
	case 2:
		fields.ID = volumeStrings[0]
		fields.Dest = volumeStrings[1]
		fields.Flags = "rw"
	case 3:
		fields.ID = volumeStrings[0]
		fields.Dest = volumeStrings[1]
		fields.Flags = volumeStrings[2]
	default:
		//NOTE: the docker cli should cover this case. This is here for posterity.
		return volumeFields{}, fmt.Errorf("Volume bind input is invalid : -v %s", volString)
	}
	return fields, nil
}

func processAnonymousVolumes(h *string, volumes map[string]struct{}, client *client.PortLayer) ([]volumeFields, error) {
	var volumeFields []volumeFields

	for v := range volumes {
		fields, err := processVolumeParam(v)
		log.Infof("Processed Volume arguments : %#v", fields)
		if err != nil {
			return nil, err
		}
		//NOTE: This should be the guard for the case of an anonymous volume.
		//NOTE: we should not expect any driver args if the drive is anonymous.
		log.Infof("anonymous volume being created - Container Create - volume mount section ID: %s ", fields.ID)
		metadata := make(map[string]string)
		metadata["flags"] = fields.Flags
		volumeRequest := models.VolumeRequest{
			Capacity: -1,
			Driver:   "vsphere",
			Store:    "default",
			Name:     fields.ID,
			Metadata: metadata,
		}
		_, err = client.Storage.CreateVolume(storage.NewCreateVolumeParams().WithVolumeRequest(&volumeRequest))
		if err != nil {
			return nil, err
		}
		volumeFields = append(volumeFields, fields)
	}
	return volumeFields, nil
}

func processSpecifiedVolumes(volumes []string) ([]volumeFields, error) {
	var volumeFields []volumeFields
	for _, v := range volumes {
		fields, err := processVolumeParam(v)
		if err != nil {
			return volumeFields, err
		}
		volumeFields = append(volumeFields, fields)
	}
	return volumeFields, nil
}
