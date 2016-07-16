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

//****
// container_portlayer.go
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
//		- DO USE the aliased docker error package 'derr'

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/go-swagger/go-swagger/httpkit"
	httptransport "github.com/go-swagger/go-swagger/httpkit/client"
	strfmt "github.com/go-swagger/go-swagger/strfmt"
	"github.com/mreiferson/go-httpclient"

	"github.com/docker/docker/api/types/backend"
	derr "github.com/docker/docker/errors"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"

	"github.com/google/uuid"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/interaction"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/metadata"
)

type volumeFields struct {
	ID    string
	Dest  string
	Flags string
}

const (
	attachConnectTimeout   time.Duration = 15 * time.Second //timeout for the connection
	attachAttemptTimeout   time.Duration = 40 * time.Second //timeout before we ditch an attach attempt
	attachPLAttemptDiff    time.Duration = 10 * time.Second
	attachPLAttemptTimeout time.Duration = attachAttemptTimeout - attachPLAttemptDiff //timeout for the portlayer before ditching an attempt
	attachRequestTimeout   time.Duration = 2 * time.Hour                              //timeout to hold onto the attach connection
	swaggerSubstringEOF                  = "EOF"
)

// VicCreateContainer() creates a VIC container by making remote calls to the
// VIC portlayer.
//
// returns:
// 	(containerID, containerHandle, error)
func VicCreateContainer(vc *viccontainer.VicContainer, config types.ContainerCreateConfig) (string, string, error) {
	// Get an API client to the portlayer
	client := PortLayerClient()
	if client == nil {
		return "", "", derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	// Call the Exec port layer to create the container
	host, err := guest.UUID()
	if err != nil {
		return "", "", derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate got unexpected error getting VCH UUID"),
			http.StatusInternalServerError)
	}

	plCreateParams := dockerContainerCreateParamsToPortlayer(config, vc.ID, host)
	createResults, err := client.Containers.Create(plCreateParams)
	// transfer port layer swagger based response to Docker backend data structs and return to the REST front-end
	if err != nil {
		if _, ok := err.(*containers.CreateNotFound); ok {
			err = fmt.Errorf("No such image: %s", vc.ID)
			log.Errorf(err.Error())
			return "", "", derr.NewRequestNotFoundError(err)
		}

		// If we get here, most likely something went wrong with the port layer API server
		return "", "", derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
	}

	id := createResults.Payload.ID
	h := createResults.Payload.Handle

	log.Debugf("Network Configuration Section - Container Create")
	// configure networking
	netConf := toModelsNetworkConfig(config)
	if netConf != nil {
		addContRes, err := client.Scopes.AddContainer(scopes.NewAddContainerParams().
			WithScope(netConf.NetworkName).
			WithConfig(&models.ScopesAddContainerConfig{
				Handle:        h,
				NetworkConfig: netConf,
			}))

		if err != nil {
			log.Errorf("ContainerCreate: Scopes error: %s", err.Error())
			return "", "", derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
		}

		defer func() {
			if err == nil {
				return
			}
			// roll back the AddContainer call
			if _, err2 := client.Scopes.RemoveContainer(scopes.NewRemoveContainerParams().WithHandle(h).WithScope(netConf.NetworkName)); err2 != nil {
				log.Warnf("could not roll back container add: %s", err2)
			}
		}()

		h = addContRes.Payload
	}

	//Volume Attachment Section
	log.Debugf("Container.ContainerCreate - VolumeSection")
	log.Debugf("Raw Volume arguments : binds:  %#v : volumes : %#v", config.HostConfig.Binds, config.Config.Volumes)
	var joinList []volumeFields

	joinList, err = processAnonymousVolumes(&h, config.Config.Volumes, client)
	if err != nil {
		return "", "", derr.NewErrorWithStatusCode(fmt.Errorf("Server error from Portlayer: %s", err), http.StatusBadRequest)
	}

	volumeSubset, err := processSpecifiedVolumes(config.HostConfig.Binds)
	if err != nil {
		return "", "", derr.NewErrorWithStatusCode(fmt.Errorf("Server error from Portlayer: %s", err), http.StatusBadRequest)
	}
	joinList = append(joinList, volumeSubset...)

	for _, fields := range joinList {
		flags := make(map[string]string)
		//NOTE: for now we are passing the flags directly through. This is NOT SAFE and only a stop gap.
		flags["Mode"] = fields.Flags
		joinParams := storage.NewVolumeJoinParams().WithJoinArgs(&models.VolumeJoinConfig{
			Flags:     flags,
			Handle:    h,
			MountPath: fields.Dest,
		}).WithName(fields.ID)

		res, err := client.Storage.VolumeJoin(joinParams)
		if err != nil {
			return "", "", derr.NewErrorWithStatusCode(fmt.Errorf("Server error from Portlayer: %s", err), http.StatusInternalServerError)
		}
		h = res.Payload
	}
	// commit the create op
	_, err = client.Containers.Commit(containers.NewCommitParams().WithHandle(h))
	if err != nil {
		err = fmt.Errorf("No such image: %s", vc.ID)
		log.Errorf("%s", err.Error())
		// FIXME: Containers.Commit returns more errors than it's swagger spec says.
		// When no image exist, it also sends back non swagger errors.  We should fix
		// this once Commit returns correct error codes.
		return "", "", derr.NewRequestNotFoundError(err)
	}

	return id, h, nil
}

func VicContainerRemove(name string) error {
	// Get the portlayer Client API
	client := PortLayerClient()
	if client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerRm failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	//TODO: verify params that are passed in(like -f) and then pass them along.

	//FIXME: currently the portlayer yaml does not support more params than the simple name param

	//call the remove directly on the name. No need for using a handle.
	_, err := client.Containers.ContainerRemove(containers.NewContainerRemoveParams().WithID(name))
	if err != nil {
		if _, ok := err.(*containers.ContainerRemoveNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer"), http.StatusInternalServerError)
	}

	return nil
}

func VicContainerStart(name string, hostConfig *container.HostConfig, bind bool) error {
	var err error

	// Get an API client to the portlayer
	client := PortLayerClient()
	if client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	// handle legacy hostConfig
	if hostConfig != nil {
		// hostConfig exist for backwards compatibility.  TODO: Figure out which parameters we
		// need to look at in hostConfig
	}

	// get a handle to the container
	getRes, err := client.Containers.Get(containers.NewGetParams().WithID(name))
	if err != nil {
		if _, ok := err.(*containers.GetNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer"), http.StatusInternalServerError)
	}

	h := getRes.Payload

	// error handling just in case bind fails
	defer func() {
		if err != nil {
			// roll back the BindContainer call
			if _, err = client.Scopes.UnbindContainer(scopes.NewUnbindContainerParams().WithHandle(h)); err != nil {
				log.Warnf("failed to roll back container bind: %s", err.Error())
			}
		}
	}()

	// bind network
	if bind {
		bindRes, err := client.Scopes.BindContainer(scopes.NewBindContainerParams().WithHandle(h))
		if err != nil {
			switch err := err.(type) {
			case *scopes.BindContainerNotFound:
				return derr.NewRequestNotFoundError(fmt.Errorf(err.Payload.Message))

			case *scopes.BindContainerInternalServerError:
				return derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

			default:
				return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
			}
		}

		h = bindRes.Payload
	}

	// change the state of the container
	// TODO: We need a resolved ID from the name
	stateChangeRes, err := client.Containers.StateChange(containers.NewStateChangeParams().WithHandle(h).WithState("RUNNING"))
	if err != nil {
		if _, ok := err.(*containers.StateChangeNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer"))
		}

		// If we get here, most likely something went wrong with the port layer API server

		return derr.NewErrorWithStatusCode(fmt.Errorf("Unknown error from the exec port layer"), http.StatusInternalServerError)
	}

	h = stateChangeRes.Payload

	// commit the handle; this will reconfigure and start the vm
	_, err = client.Containers.Commit(containers.NewCommitParams().WithHandle(h))
	if err != nil {
		if _, ok := err.(*containers.CommitNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer"))
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer"), http.StatusInternalServerError)
	}

	return nil
}

func VicContainerStop(containerID string, seconds int, unbound bool) error {
	//retrieve client to portlayer
	client := PortLayerClient()
	if client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	getResponse, err := client.Containers.Get(containers.NewGetParams().WithID(containerID))
	if err != nil {
		if _, ok := err.(*containers.GetNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", containerID))
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer"), http.StatusInternalServerError)
	}

	handle := getResponse.Payload

	if unbound {
		ub, err := client.Scopes.UnbindContainer(scopes.NewUnbindContainerParams().WithHandle(handle))
		if err != nil {
			switch err := err.(type) {
			case *scopes.UnbindContainerNotFound:
				return derr.NewRequestNotFoundError(fmt.Errorf("container %s not found", containerID))

			case *scopes.UnbindContainerInternalServerError:
				return derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

			default:
				return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
			}
		}

		handle = ub.Payload
	}

	// change the state of the container
	// TODO: We need a resolved ID from the name
	stateChangeResponse, err := client.Containers.StateChange(containers.NewStateChangeParams().WithHandle(handle).WithState("STOPPED"))
	if err != nil {
		if _, ok := err.(*containers.StateChangeNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer"))
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer"), http.StatusInternalServerError)
	}

	handle = stateChangeResponse.Payload

	_, err = client.Containers.Commit(containers.NewCommitParams().WithHandle(handle))
	if err != nil {
		if _, ok := err.(*containers.CommitNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer"))
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer"), http.StatusInternalServerError)
	}

	return nil
}

func VicResizeContainer(name string, height, width int) error {
	// Get an API client to the portlayer
	client := PortLayerClient()
	if client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerResize failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	// Call the port layer to resize
	plHeight := int32(height)
	plWidth := int32(width)
	plResizeParam := interaction.NewContainerResizeParams().WithID(name).WithHeight(plHeight).WithWidth(plWidth)

	_, err := client.Interaction.ContainerResize(plResizeParam)
	if err != nil {
		if _, isa := err.(*interaction.ContainerResizeNotFound); isa {
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
		}

		// If we get here, most likely something went wrong with the port layer API server
		return derr.NewErrorWithStatusCode(fmt.Errorf("Unknown error from the interaction port layer: %s", err),
			http.StatusInternalServerError)
	}

	return nil
}

func VicContainerInspect(containerID string) (*models.ContainerInfo, error) {
	client := PortLayerClient()
	if client == nil {
		return nil, derr.NewErrorWithStatusCode(fmt.Errorf("Failed to get portlayer client"), http.StatusInternalServerError)
	}

	results, err := client.Containers.GetContainerInfo(containers.NewGetContainerInfoParams().WithID(containerID))
	if err != nil {
		switch err := err.(type) {
		case *containers.GetContainerInfoNotFound:
			return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", containerID))
		case *containers.GetContainerInfoInternalServerError:
			return nil, derr.NewErrorWithStatusCode(fmt.Errorf("Error from portlayer: %#v", err.Payload), http.StatusInternalServerError)
		default:
			return nil, derr.NewErrorWithStatusCode(fmt.Errorf("Unknown error from the container portlayer"), http.StatusInternalServerError)
		}
	}

	return results.Payload, nil
}

func VicContainerList(listAll bool) ([]models.ContainerListInfo, error) {
	portLayerClient := PortLayerClient()
	if portLayerClient == nil {
		return nil, derr.NewErrorWithStatusCode(fmt.Errorf("container.Containers failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	containme, err := portLayerClient.Containers.GetContainerList(containers.NewGetContainerListParams().WithAll(&listAll))
	if err != nil {
		return nil, fmt.Errorf("Error invoking GetContainerList: %s", err.Error())
	}

	return containme.Payload, nil
}

// VicAttachStreams takes the the hijacked connections from the calling client and attaches
// them to the 3 streams from the portlayer's rest server.
// clStdin, clStdout, clStderr are the hijacked connection
func VicAttachStreams(ctx context.Context, vc *viccontainer.VicContainer, clStdin io.ReadCloser, clStdout, clStderr io.Writer, ca *backend.ContainerAttachConfig) error {
	defer clStdin.Close()

	// Cancel will close the child connections.
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	errors := make(chan error, 3)

	// For stdin, we only have a timeout for connection.  There can be a long duration before
	// the first entry so there is no timeout for response.
	plClient, transport := createNewAttachClientWithTimeouts(attachConnectTimeout, attachAttemptTimeout, 0)
	defer transport.Close()

	if ca.UseStdin {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := copyStdIn(ctx, plClient, vc, clStdin, ca.DetachKeys)
			if err != nil {
				log.Errorf("container attach: stdin (%s): %s", vc.ContainerID, err.Error())
			} else {
				log.Infof("container attach: stdin (%s) done: %s", vc.ContainerID)
			}

			cancel()
			errors <- err
		}()
	}

	if ca.UseStdout {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := copyStdOut(ctx, plClient, attachAttemptTimeout, vc, clStdout)
			if err != nil {
				log.Errorf("container attach: stdout (%s): %s", vc.ContainerID, err.Error())
			} else {
				log.Infof("container attach: stdout (%s) done: %s", vc.ContainerID)
			}

			cancel()
			errors <- err
		}()
	}

	if ca.UseStderr {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := copyStdErr(ctx, plClient, vc, clStderr)
			if err != nil {
				log.Errorf("container attach: stderr (%s): %s", vc.ContainerID, err.Error())
			} else {
				log.Infof("container attach: stderr (%s) done: %s", vc.ContainerID)
			}

			cancel()
			errors <- err
		}()
	}

	// Wait for all stream copy to exit
	wg.Wait()

	log.Infof("container attach:  cleaned up connections to %s.", vc.ContainerID)
	for err := range errors {
		if err != nil {
			// If we get here, most likely something went wrong with the port layer API server
			// These errors originate within the go-swagger client itself.
			// Go-swagger returns untyped errors to us if the error is not one that we define
			// in the swagger spec.  Even EOF.  Therefore, we must scan the error string (if there
			// is an error string in the untyped error) for the term EOF.

			log.Errorf("container attach error: %s", err.Error())

			return err
		}
	}

	return nil
}

//----------
// Conversion utility functions - convert swagger to docker engine-api structs
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
			nc.Aliases = EP2Alias(es)
		}
	}
	return nc
}

// plContainerInfoToDockerContainerInspect() takes a ContainerInfo swagger-based struct
// returned from VIC's port layer and creates an engine-api based container inspect struct.
// There maybe other asset gathering if ContainerInfo does not have all the information
func plContainerInfoToDockerContainerInspect(id string, info *models.ContainerInfo) (*types.ContainerJSON, error) {
	if info == nil || info.ContainerConfig == nil {
		return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", id))
	}

	// Set default container state attributes
	containerState := &types.ContainerState{
		Restarting: false,
		OOMKilled:  false,
	}

	if info.ProcessConfig != nil {
		if info.ProcessConfig.Pid != nil {
			containerState.Pid = int(*info.ProcessConfig.Pid)
		}
		if info.ProcessConfig.ExitCode != nil {
			containerState.ExitCode = int(*info.ProcessConfig.ExitCode)
		}
		if info.ProcessConfig.ErrorMsg != nil {
			containerState.Error = *info.ProcessConfig.ErrorMsg
		}
		if info.ProcessConfig.Started != nil {
			swaggerTime := time.Time(*info.ProcessConfig.Started)
			containerState.StartedAt = swaggerTime.Format(time.RFC3339Nano)
		}
		if info.ProcessConfig.Finished != nil {
			swaggerTime := time.Time(*info.ProcessConfig.Finished)
			containerState.FinishedAt = swaggerTime.Format(time.RFC3339Nano)
		}
	}

	inpsectJSON := &types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State:           containerState,
			ResolvConfPath:  "",
			HostnamePath:    "",
			HostsPath:       "",
			Driver:          PortLayerName(),
			MountLabel:      "",
			ProcessLabel:    "",
			AppArmorProfile: "",
			ExecIDs:         nil,
			HostConfig:      hostConfigFromContainerInfo(id, info),
			GraphDriver:     types.GraphDriverData{Name: PortLayerName()},
			SizeRw:          nil,
			SizeRootFs:      nil,
		},
		Mounts:          mountsFromContainerInfo(id, info),
		Config:          containerConfigFromContainerInfo(id, info),
		NetworkSettings: networkFromContainerInfo(id, info),
	}

	if info.ProcessConfig != nil {
		if info.ProcessConfig.ExecPath != nil {
			inpsectJSON.Path = *info.ProcessConfig.ExecPath
		}
		if info.ProcessConfig.ExecArgs != nil {
			inpsectJSON.Args = info.ProcessConfig.ExecArgs
		}
	}

	if info.ContainerConfig != nil {
		if info.ContainerConfig.State != nil {
			containerState.Status = strings.ToLower(*info.ContainerConfig.State)

			containerState.Running = false
			containerState.Paused = false //We do not yet support paused/resumed
			containerState.Dead = false
			if containerState.Status == "running" {
				containerState.Running = true
				containerState.Dead = false //This is only true during docker rm
			}
		}
		if info.ContainerConfig.LayerID != nil {
			inpsectJSON.Image = *info.ContainerConfig.LayerID
		}
		if info.ContainerConfig.LogPath != nil {
			inpsectJSON.LogPath = *info.ContainerConfig.LogPath
		}
		if info.ContainerConfig.RestartCount != nil {
			inpsectJSON.RestartCount = int(*info.ContainerConfig.RestartCount)
		}
		if info.ContainerConfig.ContainerID != nil {
			inpsectJSON.ID = *info.ContainerConfig.ContainerID
		}
		if info.ContainerConfig.Created != nil {
			inpsectJSON.Created = time.Unix(*info.ContainerConfig.Created, 0).String()
		}
		if len(info.ContainerConfig.Names) > 0 {
			inpsectJSON.Name = info.ContainerConfig.Names[0]
		}
	}

	return inpsectJSON, nil
}

// hostConfigFromContainerInfo() extracts docker compatible hostconfig data from the
// Swagger-based ContainerInfo object.
func hostConfigFromContainerInfo(id string, info *models.ContainerInfo) *container.HostConfig {
	if info == nil {
		return nil
	}

	// Resources don't really map well to VIC so we leave mose of them empty. If we look
	// at the struct in engine-api/types/container/host_config.go, Microsoft added
	// additional attributes to the struct that are applicable to Windows containers.
	// If understanding VIC's host resources are desireable, we should go down this
	// same route.
	//
	// The values we fill out below is an abridged list of the original struct.
	resourceConfig := container.Resources{
	// Applicable to all platforms
	//			CPUShares int64 `json:"CpuShares"` // CPU shares (relative weight vs. other containers)
	//			Memory    int64 // Memory limit (in bytes)

	//			// Applicable to UNIX platforms
	//			DiskQuota            int64           // Disk limit (in bytes)
	}

	hostConfig := &container.HostConfig{
		Binds:           nil,
		ContainerIDFile: "",
		LogConfig: container.LogConfig{
			Type:   "",
			Config: nil,
		},
		PortBindings: nil, // Port mapping between the exposed port (container) and the host
		RestartPolicy: container.RestartPolicy{
			Name:              "",
			MaximumRetryCount: 0,
		}, // Restart policy to be used for the container
		AutoRemove:   false,           // Automatically remove container when it exits
		VolumeDriver: PortLayerName(), // Name of the volume driver used to mount volumes
		VolumesFrom:  nil,             // List of volumes to take from other container

		// Applicable to UNIX platforms
		CapAdd:          nil,   // List of kernel capabilities to add to the container
		CapDrop:         nil,   // List of kernel capabilities to remove from the container
		DNSOptions:      nil,   // List of DNSOption to look for
		DNSSearch:       nil,   // List of DNSSearch to look for
		ExtraHosts:      nil,   // List of extra hosts
		GroupAdd:        nil,   // List of additional groups that the container process will run as
		IpcMode:         "",    // IPC namespace to use for the container
		Cgroup:          "",    // Cgroup to use for the container
		Links:           nil,   // List of links (in the name:alias form)
		OomScoreAdj:     0,     // Container preference for OOM-killing
		PidMode:         "",    // PID namespace to use for the container
		Privileged:      false, // Is the container in privileged mode
		PublishAllPorts: false, // Should docker publish all exposed port for the container
		ReadonlyRootfs:  false, // Is the container root filesystem in read-only
		SecurityOpt:     nil,   // List of string values to customize labels for MLS systems, such as SELinux.
		StorageOpt:      nil,   // Storage driver options per container.
		Tmpfs:           nil,   // List of tmpfs (mounts) used for the container
		UTSMode:         "",    // UTS namespace to use for the container
		UsernsMode:      "",    // The user namespace to use for the container
		ShmSize:         0,     // Total shm memory usage
		Sysctls:         nil,   // List of Namespaced sysctls used for the container

		// Applicable to Windows
		Isolation: "", // Isolation technology of the container (eg default, hyperv)

		// Contains container's resources (cgroups, ulimits)
		Resources: resourceConfig,
	}

	if len(info.ScopeConfig) > 0 {
		if info.ScopeConfig[0].DNS != nil {
			hostConfig.DNS = info.ScopeConfig[0].DNS
		}

		hostConfig.NetworkMode = container.NetworkMode(info.ScopeConfig[0].ScopeType)
	}

	return hostConfig
}

// mountsFromContainerInfo()
func mountsFromContainerInfo(id string, info *models.ContainerInfo) []types.MountPoint {
	if info == nil {
		return nil
	}

	var mounts []types.MountPoint

	for _, vConfig := range info.VolumeConfig {
		// Fill with defaults
		mountConfig := types.MountPoint{
			Destination: "",
			Driver:      PortLayerName(),
			Mode:        "",
			Propagation: "",
		}

		// Fill with info from portlayer
		if vConfig.MountPoint != nil {
			mountConfig.Name = *vConfig.MountPoint
		}
		if vConfig.MountPoint != nil {
			mountConfig.Source = *vConfig.MountPoint
		}
		if vConfig.ReadWrite != nil {
			mountConfig.RW = *vConfig.ReadWrite
		}

		mounts = append(mounts, mountConfig)
	}

	return mounts
}

// containerConfigFromContainerInfo() returns a container.Config that has attributes
// overridden at create or start time.  This is important.  This function is called
// to help build the Container Inspect struct.  That struct contains the original
// container config that is part of the image metadata AND the overriden container
// config.  The user can override these via the remote API or the docker CLI.
func containerConfigFromContainerInfo(id string, info *models.ContainerInfo) *container.Config {
	if info == nil || info.ContainerConfig == nil || info.ProcessConfig == nil {
		return nil
	}

	container := &container.Config{
		Domainname: "",    // Domainname
		User:       "",    // User that will run the command(s) inside the container
		StdinOnce:  false, // If true, close stdin after the 1 attached client disconnects.
	}

	if info.ContainerConfig.ContainerID != nil {
		container.Hostname = stringid.TruncateID(*info.ContainerConfig.ContainerID) // Hostname
	}
	if info.ContainerConfig.AttachStdin != nil {
		container.AttachStdin = *info.ContainerConfig.AttachStdin // Attach the standard input, makes possible user interaction
	}
	if info.ContainerConfig.AttachStdout != nil {
		container.AttachStdout = *info.ContainerConfig.AttachStdout // Attach the standard output
	}
	if info.ContainerConfig.AttachStderr != nil {
		container.AttachStderr = *info.ContainerConfig.AttachStderr // Attach the standard error
	}
	if info.ContainerConfig.Tty != nil {
		container.Tty = *info.ContainerConfig.Tty // Attach standard streams to a tty, including stdin if it is not closed.
	}
	// They are not coming from PL so set them to true unconditionally
	container.OpenStdin = true // Open stdin
	container.StdinOnce = true

	if info.ContainerConfig.LayerID != nil {
		container.Image = *info.ContainerConfig.LayerID // Name of the image as it was passed by the operator (eg. could be symbolic)
	}
	if info.ContainerConfig.Labels != nil {
		container.Labels = info.ContainerConfig.Labels // List of labels set to this container
	}

	// Fill in information about the process
	if info.ProcessConfig.ExecArgs != nil {
		container.Env = info.ProcessConfig.ExecArgs // List of environment variable to set in the container
	}
	if info.ProcessConfig.ExecPath != nil {
		container.Cmd = append(container.Cmd, *info.ProcessConfig.ExecPath) // Command to run when starting the container
	}
	if info.ProcessConfig.WorkingDir != nil {
		container.WorkingDir = *info.ProcessConfig.WorkingDir // Current directory (PWD) in the command will be launched
	}
	//		container.Entrypoint      strslice.StrSlice     				// Entrypoint to run when starting the container

	// Fill in information about the container network
	if info.ScopeConfig == nil {
		container.NetworkDisabled = true
	} else {
		container.NetworkDisabled = false
		container.MacAddress = ""
		container.ExposedPorts = nil  //FIXME:  Add once port mapping is implemented
		container.PublishService = "" // Name of the network service exposed by the container
	}

	// Get the original container config from the image's metadata in our image cache.
	var imageConfig *metadata.ImageConfig

	if info.ContainerConfig.LayerID != nil {
		imageConfig, _ = cache.ImageCache().GetImage(*info.ContainerConfig.LayerID)
	}

	// Fill in the values with defaults from the original image's container config
	// structure
	if imageConfig != nil {
		container.StopSignal = imageConfig.ContainerConfig.StopSignal // Signal to stop a container

		container.OnBuild = imageConfig.ContainerConfig.OnBuild // ONBUILD metadata that were defined on the image Dockerfile

		// Fill in information about the container's volumes
		// FIXME:  Why does types.ContainerJSON have Mounts and also ContainerConfig,
		// which also has Volumes?  Assuming this is a copy from image's container
		// config till we figure this out.
		container.Volumes = imageConfig.ContainerConfig.Volumes
	}

	return container
}

func networkFromContainerInfo(id string, info *models.ContainerInfo) *types.NetworkSettings {
	if info == nil || info.ScopeConfig == nil {
		return nil
	}

	networks := &types.NetworkSettings{
		NetworkSettingsBase: types.NetworkSettingsBase{
			Bridge:                 "",
			SandboxID:              "",
			HairpinMode:            false,
			LinkLocalIPv6Address:   "",
			LinkLocalIPv6PrefixLen: 0,
			Ports:                  nil, //FIXME:  Fill in once port mapping is implemented.
			SandboxKey:             "",
			SecondaryIPAddresses:   nil,
			SecondaryIPv6Addresses: nil,
		},
		Networks: nil,
	}

	return networks
}

func plContainerListToDockerContainerList(plContainerList []models.ContainerListInfo) []*types.Container {
	containers := make([]*types.Container, 0, len(plContainerList))
	for _, t := range plContainerList {
		cmd := strings.Join(t.ExecArgs, " ")
		// the docker client expects the friendly name to be prefixed
		// with a forward slash -- create a new slice and add here
		names := make([]string, 0, len(t.Names))
		for i := range t.Names {
			names = append(names, clientFriendlyContainerName(t.Names[i]))
		}
		c := &types.Container{
			ID:      *t.ContainerID,
			Image:   *t.RepoName,
			Created: *t.Created,
			Status:  *t.Status,
			Names:   names,
			Command: cmd,
			SizeRw:  *t.StorageSize,
		}
		containers = append(containers, c)
	}

	return containers
}

// helper function to format the container name
// to the docker client approved format
func clientFriendlyContainerName(name string) string {
	return fmt.Sprintf("/%s", name)
}

//----------
// Utility Functions
//----------

func createNewAttachClientWithTimeouts(connectTimeout, responseHeaderTimeout, responseTimeout time.Duration) (*client.PortLayer, *httpclient.Transport) {
	runtime := httptransport.New(PortLayerServer(), "/", []string{"http"})
	transport := &httpclient.Transport{
		ConnectTimeout:        connectTimeout,
		ResponseHeaderTimeout: responseHeaderTimeout,
		RequestTimeout:        responseTimeout,
	}
	runtime.Transport = transport

	plClient := client.New(runtime, nil)
	runtime.Consumers["application/octet-stream"] = httpkit.ByteStreamConsumer()
	runtime.Producers["application/octet-stream"] = httpkit.ByteStreamProducer()

	return plClient, transport
}

func copyStdIn(ctx context.Context, pl *client.PortLayer, vc *viccontainer.VicContainer, clStdin io.ReadCloser, keys []byte) error {
	// Pipe for stdin so we can interject and watch the input streams for detach keys.
	stdinReader, stdinWriter := io.Pipe()

	defer stdinWriter.Close()

	go func() {
		defer stdinReader.Close()
		// Copy the stdin from the CLI and write to a pipe.  We need to do this so we can
		// watch the stdin stream for the detach keys.
		var err error
		if vc.Config.Tty {
			_, err = copyEscapable(stdinWriter, clStdin, keys)
		} else {
			_, err = io.Copy(stdinWriter, clStdin)
		}

		if err != nil {
			log.Errorf("container attach: stdin err: %s", err.Error())
		}
	}()

	// Swagger wants an io.reader so give it the reader pipe.  Also, the swagger call
	// to set the stdin is synchronous so we need to run in a goroutine
	setStdinParams := interaction.NewContainerSetStdinParamsWithContext(ctx).WithID(vc.ContainerID)
	setStdinParams = setStdinParams.WithRawStream(stdinReader)
	_, err := pl.Interaction.ContainerSetStdin(setStdinParams)

	if vc.Config.StdinOnce && !vc.Config.Tty {
		// Close the stdin connection.  Mimicing Docker's behavior.
		// FIXME: If we close this stdin connection.  The portlayer does
		// not really close stdin.  This is diff from current Docker
		// behavior.  However, we're not sure why Docker even has this
		// behavior where you connect to stdin on the first time only.
		// If we really want to add this behavior, we need to add support
		// in the tether in the portlayer.
		log.Errorf("Attach stream has stdinOnce set.  VIC does not yet support this.")
	}
	return err
}

func copyStdOut(ctx context.Context, pl *client.PortLayer, attemptTimeout time.Duration, vc *viccontainer.VicContainer, clStdout io.Writer) error {
	name := vc.ContainerID
	//Calculate how much time to let portlayer attempt
	plAttemptTimeout := attemptTimeout - attachPLAttemptDiff //assumes personality deadline longer than portlayer's deadline
	plAttemptDeadline := time.Now().Add(plAttemptTimeout)
	swaggerDeadline := strfmt.DateTime(plAttemptDeadline)
	log.Debugf("* stdout portlayer deadline: %s", plAttemptDeadline.Format(time.UnixDate))
	log.Debugf("* stdout personality deadline: %s", time.Now().Add(attemptTimeout).Format(time.UnixDate))

	log.Debugf("* stdout attach start %s", time.Now().Format(time.UnixDate))
	getStdoutParams := interaction.NewContainerGetStdoutParamsWithContext(ctx).WithID(name).WithDeadline(&swaggerDeadline)
	_, err := pl.Interaction.ContainerGetStdout(getStdoutParams, clStdout)
	log.Debugf("* stdout attach end %s", time.Now().Format(time.UnixDate))
	if err != nil {
		if _, ok := err.(*interaction.ContainerGetStdoutNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
		}

		if _, ok := err.(*interaction.ContainerGetStdoutInternalServerError); ok {
			return derr.NewErrorWithStatusCode(fmt.Errorf("Server error from the interaction port layer"),
				http.StatusInternalServerError)
		}

		unknownErrMsg := fmt.Errorf("Unknown error from the interaction port layer: %s", err)
		return derr.NewErrorWithStatusCode(unknownErrMsg, http.StatusInternalServerError)
	}

	return nil
}

func copyStdErr(ctx context.Context, pl *client.PortLayer, vc *viccontainer.VicContainer, clStderr io.Writer) error {
	name := vc.ContainerID
	getStderrParams := interaction.NewContainerGetStderrParamsWithContext(ctx).WithID(name)
	_, err := pl.Interaction.ContainerGetStderr(getStderrParams, clStderr)

	if err != nil {
		if _, ok := err.(*interaction.ContainerGetStderrNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
		}

		if _, ok := err.(*interaction.ContainerGetStderrInternalServerError); ok {
			return derr.NewErrorWithStatusCode(fmt.Errorf("Server error from the interaction port layer"),
				http.StatusInternalServerError)
		}

		unknownErrMsg := fmt.Errorf("Unknown error from the interaction port layer: %s", err)
		return derr.NewErrorWithStatusCode(unknownErrMsg, http.StatusInternalServerError)
	}

	return nil
}

// FIXME: Move this function to a pkg to show it's origination from Docker once
// we have ignore capabilities in our header-check.sh that checks for copyright
// header.
// Code c/c from io.Copy() modified by Docker to handle escape sequence
func copyEscapable(dst io.Writer, src io.ReadCloser, keys []byte) (written int64, err error) {
	if len(keys) == 0 {
		// Default keys : ctrl-p ctrl-q
		keys = []byte{16, 17}
	}
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			// ---- Docker addition
			for i, key := range keys {
				if nr != 1 || buf[0] != key {
					break
				}
				if i == len(keys)-1 {
					if err := src.Close(); err != nil {
						return 0, err
					}
					return 0, nil
				}
				nr, er = src.Read(buf)
			}
			// ---- End of docker
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
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
