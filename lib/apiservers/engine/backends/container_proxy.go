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
	"io"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"github.com/mreiferson/go-httpclient"

	"github.com/go-swagger/go-swagger/httpkit"
	httptransport "github.com/go-swagger/go-swagger/httpkit/client"

	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	dnetwork "github.com/docker/engine-api/types/network"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/go-connections/nat"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	epoint "github.com/vmware/vic/lib/apiservers/engine/backends/endpoint"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/sys"
)

type VicContainerProxy interface {
	CreateContainerHandle(imageID string, config types.ContainerCreateConfig) (string, string, error)
	AddContainerToScope(handle string, config types.ContainerCreateConfig) (string, error)
	AddVolumesToContainer(handle string, config types.ContainerCreateConfig) (string, error)
	CommitContainerHandle(handle, imageID string) error
	StreamContainerLogs(name string, out io.Writer, started chan struct{}, showTimestamps bool, followLogs bool, since int64, tailLines int64) error
	ContainerRunning(vc *viccontainer.VicContainer) (bool, error)

	Client() *client.PortLayer
}

type ContainerProxy struct {
	client        *client.PortLayer
	portlayerAddr string
	portlayerName string
}

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
	forceLogType                         = "json-file" //Use in inspect to allow docker logs to work
)

// NewContainerProxy creates a new ContainerProxy
func NewContainerProxy(plClient *client.PortLayer, portlayerAddr string, portlayerName string) *ContainerProxy {
	return &ContainerProxy{client: plClient, portlayerAddr: portlayerAddr, portlayerName: portlayerName}
}

func (c *ContainerProxy) Client() *client.PortLayer {
	return c.client
}

// CreateContainerHandle creates a new VIC container by calling the portlayer
//
// returns:
//	(containerID, containerHandle, error)
func (c *ContainerProxy) CreateContainerHandle(imageID string, config types.ContainerCreateConfig) (string, string, error) {
	defer trace.End(trace.Begin(imageID))

	if c.client == nil {
		return "", "", InternalServerError("ContainerProxy.CreateContainerHandle failed to create a portlayer client")
	}

	if imageID == "" {
		return "", "", NotFoundError("No image specified")
	}

	// Call the Exec port layer to create the container
	host, err := sys.UUID()
	if err != nil {
		return "", "", InternalServerError("ContainerProxy.CreateContainerHandle got unexpected error getting VCH UUID")
	}

	plCreateParams := dockerContainerCreateParamsToPortlayer(config, imageID, host)
	createResults, err := c.client.Containers.Create(plCreateParams)
	if err != nil {
		if _, ok := err.(*containers.CreateNotFound); ok {
			cerr := fmt.Errorf("No such image: %s", imageID)
			log.Errorf("%s (%s)", cerr, err)
			return "", "", NotFoundError(cerr.Error())
		}

		// If we get here, most likely something went wrong with the port layer API server
		return "", "", InternalServerError(err.Error())
	}

	id := createResults.Payload.ID
	h := createResults.Payload.Handle

	return id, h, nil
}

// AddContainerToScope adds a container, referenced by handle, to a scope.
// If an error is return, the returned handle should not be used.
//
// returns:
//	modified handle
func (c *ContainerProxy) AddContainerToScope(handle string, config types.ContainerCreateConfig) (string, error) {
	defer trace.End(trace.Begin(handle))

	if c.client == nil {
		return "", InternalServerError("ContainerProxy.AddContainerToScope failed to create a portlayer client")
	}

	log.Debugf("Network Configuration Section - Container Create")
	// configure networking
	netConf := toModelsNetworkConfig(config)
	if netConf != nil {
		addContRes, err := c.client.Scopes.AddContainer(scopes.NewAddContainerParamsWithContext(ctx).
			WithScope(netConf.NetworkName).
			WithConfig(&models.ScopesAddContainerConfig{
				Handle:        handle,
				NetworkConfig: netConf,
			}))

		if err != nil {
			log.Errorf("ContainerProxy.AddContainerToScope: Scopes error: %s", err.Error())
			return handle, InternalServerError(err.Error())
		}

		defer func() {
			if err == nil {
				return
			}
			// roll back the AddContainer call
			if _, err2 := c.client.Scopes.RemoveContainer(scopes.NewRemoveContainerParamsWithContext(ctx).WithHandle(handle).WithScope(netConf.NetworkName)); err2 != nil {
				log.Warnf("could not roll back container add: %s", err2)
			}
		}()

		handle = addContRes.Payload
	}

	return handle, nil
}

// AddVolumesToContainer adds volumes to a container, referenced by handle.
// If an error is returned, the returned handle should not be used.
//
// returns:
//	modified handle
func (c *ContainerProxy) AddVolumesToContainer(handle string, config types.ContainerCreateConfig) (string, error) {
	defer trace.End(trace.Begin(handle))

	if c.client == nil {
		return "", InternalServerError("ContainerProxy.AddVolumesToContainer failed to create a portlayer client")
	}

	//Volume Attachment Section
	log.Debugf("ContainerProxy.AddVolumesToContainer - VolumeSection")
	log.Debugf("Raw Volume arguments : binds:  %#v : volumes : %#v", config.HostConfig.Binds, config.Config.Volumes)

	// Collect all volume mappings. In a docker create/run, they
	// can be anonymous (-v /dir) or specific (-v vol-name:/dir).
	volumes := config.HostConfig.Binds
	for anonVol := range config.Config.Volumes {
		volumes = append(volumes, anonVol)
	}

	volList, err := processVolumeFields(volumes)
	if err != nil {
		return handle, BadRequestError(err.Error())
	}

	// Create and join volumes.
	for _, fields := range volList {

		driverArgs := make(map[string]string)
		driverArgs["flags"] = fields.Flags

		// NOTE: calling volumeCreate regardless of whether the volume is already
		// present can be avoided by adding an extra optional param to VolumeJoin,
		// which would then call volumeCreate if the volume does not exist.
		vol := &Volume{}
		_, err := vol.volumeCreate(fields.ID, "vsphere", driverArgs, nil)
		if err != nil {
			switch err := err.(type) {
			case *storage.CreateVolumeConflict:
				// Implicitly ignore the error where a volume with the same name
				// already exists. We can just join the said volume to the container.
				log.Infof("a volume with the name %s already exists", fields.ID)
			case *storage.CreateVolumeNotFound:
				return handle, VolumeCreateNotFoundError(volumeStore(driverArgs))
			default:
				return handle, InternalServerError(err.Error())
			}
		} else {
			log.Infof("volumeCreate succeeded. Volume mount section ID: %s", fields.ID)
		}

		flags := make(map[string]string)
		//NOTE: for now we are passing the flags directly through. This is NOT SAFE and only a stop gap.
		flags["Mode"] = fields.Flags
		joinParams := storage.NewVolumeJoinParamsWithContext(ctx).WithJoinArgs(&models.VolumeJoinConfig{
			Flags:     flags,
			Handle:    handle,
			MountPath: fields.Dest,
		}).WithName(fields.ID)

		res, err := c.client.Storage.VolumeJoin(joinParams)
		if err != nil {
			switch err := err.(type) {
			case *storage.VolumeJoinInternalServerError:
				return handle, InternalServerError(err.Payload.Message)
			case *storage.VolumeJoinDefault:
				return handle, InternalServerError(err.Payload.Message)
			case *storage.VolumeJoinNotFound:
				return handle, VolumeJoinNotFoundError(err.Payload.Message)
			default:
				return handle, InternalServerError(err.Error())
			}
		}

		handle = res.Payload
	}

	return handle, nil
}

// CommitContainerHandle() commits any changes to container handle.
func (c *ContainerProxy) CommitContainerHandle(handle, imageID string) error {
	defer trace.End(trace.Begin(handle))

	if c.client == nil {
		return InternalServerError("ContainerProxy.CommitContainerHandle failed to create a portlayer client")
	}

	_, err := c.client.Containers.Commit(containers.NewCommitParamsWithContext(ctx).WithHandle(handle))
	if err != nil {
		cerr := fmt.Errorf("No such image: %s", imageID)
		log.Errorf("%s (%s)", cerr, err)
		// FIXME: Containers.Commit returns more errors than it's swagger spec says.
		// When no image exist, it also sends back non swagger errors.  We should fix
		// this once Commit returns correct error codes.
		return NotFoundError(cerr.Error())
	}

	return nil
}

// StreamContainerLogs reads the log stream from the portlayer rest server and writes
// it directly to the io.Writer that is passed in.
func (c *ContainerProxy) StreamContainerLogs(name string, out io.Writer, started chan struct{}, showTimestamps bool, followLogs bool, since int64, tailLines int64) error {
	defer trace.End(trace.Begin(""))

	plClient, transport := c.createNewAttachClientWithTimeouts(attachConnectTimeout, 0, attachAttemptTimeout)
	defer transport.Close()
	close(started)

	params := containers.NewGetContainerLogsParamsWithContext(ctx).
		WithID(name).
		WithFollow(&followLogs).
		WithTimestamp(&showTimestamps).
		WithSince(&since).
		WithTaillines(&tailLines)
	_, err := plClient.Containers.GetContainerLogs(params, out)
	if err != nil {
		switch err := err.(type) {
		case *containers.GetContainerLogsNotFound:
			return NotFoundError(fmt.Sprintf("No such container: %s", name))
		case *containers.GetContainerLogsInternalServerError:
			return InternalServerError("Server error from the interaction port layer")
		default:
			//Check for EOF.  Since the connection, transport, and data handling are
			//encapsulated inisde of Swagger, we can only detect EOF by checking the
			//error string
			if strings.Contains(err.Error(), swaggerSubstringEOF) {
				return nil
			}
			return InternalServerError(fmt.Sprintf("Unknown error from the interaction port layer: %s", err))
		}
	}

	return nil
}

// ContainerRunning returns true if the given container is running
func (c *ContainerProxy) ContainerRunning(vc *viccontainer.VicContainer) (bool, error) {
	defer trace.End(trace.Begin(""))

	if c.client == nil {
		return false, InternalServerError("ContainerProxy.CommitContainerHandle failed to create a portlayer client")
	}

	results, err := c.client.Containers.GetContainerInfo(containers.NewGetContainerInfoParamsWithContext(ctx).WithID(vc.ContainerID))
	if err != nil {
		switch err := err.(type) {
		case *containers.GetContainerInfoNotFound:
			return false, NotFoundError(fmt.Sprintf("No such container: %s", vc.ContainerID))
		case *containers.GetContainerInfoInternalServerError:
			return false, InternalServerError(err.Payload.Message)
		default:
			return false, InternalServerError(fmt.Sprintf("Unknown error from the interaction port layer: %s", err))
		}
	}

	inspectJSON, err := ContainerInfoToDockerContainerInspect(vc, results.Payload, c.portlayerName)
	if err != nil {
		log.Errorf("containerInfoToDockerContainerInspect failed with %s", err)
		return false, err
	}

	return inspectJSON.State.Running, nil
}

func (c *ContainerProxy) createNewAttachClientWithTimeouts(connectTimeout, responseTimeout, responseHeaderTimeout time.Duration) (*client.PortLayer, *httpclient.Transport) {
	runtime := httptransport.New(c.portlayerAddr, "/", []string{"http"})
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

	// container stop signal
	config.StopSignal = new(string)
	*config.StopSignal = cc.Config.StopSignal

	log.Debugf("dockerContainerCreateParamsToPortlayer = %+v", config)

	return containers.NewCreateParamsWithContext(ctx).WithCreateConfig(config)
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
			nc.Aliases = epoint.Alias(es)

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

// processVolumeParam is used to turn any call from docker create -v <stuff> into a volumeFields object.
// The -v has 3 forms. -v <anonymous mount path>, -v <Volume Name>:<Destination Mount Path> and
// -v <Volume Name>:<Destination Mount Path>:<mount flags>
func processVolumeParam(volString string) (volumeFields, error) {
	volumeStrings := strings.Split(volString, ":")
	fields := volumeFields{}

	// Error out if the intended volume is a directory on the client filesystem.
	numVolParams := len(volumeStrings)
	if numVolParams > 1 && strings.HasPrefix(volumeStrings[0], "/") {
		return volumeFields{}, InvalidVolumeError{}
	}

	// This switch determines which type of -v was invoked.
	switch numVolParams {
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
		// NOTE: the docker cli should cover this case. This is here for posterity.
		return volumeFields{}, InvalidBindError{volume: volString}
	}
	return fields, nil
}

// processVolumeFields parses fields for volume mappings specified in a create/run -v.
func processVolumeFields(volumes []string) ([]volumeFields, error) {
	var volumeFields []volumeFields

	for _, v := range volumes {
		fields, err := processVolumeParam(v)
		log.Infof("Processed volume arguments: %#v", fields)
		if err != nil {
			return nil, err
		}
		volumeFields = append(volumeFields, fields)
	}
	return volumeFields, nil
}

//-------------------------------------
// Inspect Utility Functions
//-------------------------------------

// ContainerInfoToDockerContainerInspect takes a ContainerInfo swagger-based struct
// returned from VIC's port layer and creates an engine-api based container inspect struct.
// There maybe other asset gathering if ContainerInfo does not have all the information
func ContainerInfoToDockerContainerInspect(vc *viccontainer.VicContainer, info *models.ContainerInfo, portlayerName string) (*types.ContainerJSON, error) {
	if vc == nil || info == nil || info.ContainerConfig == nil {
		return nil, NotFoundError(fmt.Sprintf("No such container: %s", vc.ContainerID))
	}

	// Set default container state attributes
	containerState := &types.ContainerState{}

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
		if info.ProcessConfig.StartTime != nil && *info.ProcessConfig.StartTime != 0 {
			containerState.StartedAt = time.Unix(*info.ProcessConfig.StartTime, 0).Format(time.RFC3339Nano)
		}

		if info.ProcessConfig.StopTime != nil && *info.ProcessConfig.StopTime != 0 {
			containerState.FinishedAt = time.Unix(*info.ProcessConfig.StopTime, 0).Format(time.RFC3339Nano)
		}
	}

	inspectJSON := &types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State:           containerState,
			ResolvConfPath:  "",
			HostnamePath:    "",
			HostsPath:       "",
			Driver:          portlayerName,
			MountLabel:      "",
			ProcessLabel:    "",
			AppArmorProfile: "",
			ExecIDs:         nil,
			HostConfig:      hostConfigFromContainerInfo(vc, info, portlayerName),
			GraphDriver:     types.GraphDriverData{Name: portlayerName},
			SizeRw:          nil,
			SizeRootFs:      nil,
		},
		Mounts:          mountsFromContainerInfo(vc, info, portlayerName),
		Config:          containerConfigFromContainerInfo(vc, info),
		NetworkSettings: networkFromContainerInfo(vc, info),
	}

	if inspectJSON.NetworkSettings != nil {
		log.Debugf("Docker inspect - network settings = %#v", inspectJSON.NetworkSettings)
	} else {
		log.Debug("Docker inspect - network settings = nil")
	}

	if info.ProcessConfig != nil {
		if info.ProcessConfig.ExecPath != nil {
			inspectJSON.Path = *info.ProcessConfig.ExecPath
		}
		if info.ProcessConfig.ExecArgs != nil {
			// args[0] is the command and should not appear in the args list here
			inspectJSON.Args = info.ProcessConfig.ExecArgs[1:]
		}
	}

	if info.ContainerConfig != nil {
		if info.ContainerConfig.State != nil {
			containerState.Status = strings.ToLower(*info.ContainerConfig.State)

			// https://github.com/docker/docker/blob/master/container/state.go#L77
			if containerState.Status == "stopped" {
				containerState.Status = "exited"
			}
			if containerState.Status == "running" {
				containerState.Running = true
			}
		}
		if info.ContainerConfig.LayerID != nil {
			inspectJSON.Image = *info.ContainerConfig.LayerID
		}
		if info.ContainerConfig.LogPath != nil {
			inspectJSON.LogPath = *info.ContainerConfig.LogPath
		}
		if info.ContainerConfig.RestartCount != nil {
			inspectJSON.RestartCount = int(*info.ContainerConfig.RestartCount)
		}
		if info.ContainerConfig.ContainerID != nil {
			inspectJSON.ID = *info.ContainerConfig.ContainerID
		}
		if info.ContainerConfig.CreateTime != nil {
			inspectJSON.Created = time.Unix(*info.ContainerConfig.CreateTime, 0).Format(time.RFC3339Nano)
		}
		if len(info.ContainerConfig.Names) > 0 {
			inspectJSON.Name = info.ContainerConfig.Names[0]
		}
	}

	return inspectJSON, nil
}

// hostConfigFromContainerInfo() gets the hostconfig that is passed to the backend during
// docker create and updates any needed info
func hostConfigFromContainerInfo(vc *viccontainer.VicContainer, info *models.ContainerInfo, portlayerName string) *container.HostConfig {
	if vc == nil || vc.HostConfig == nil || info == nil {
		return nil
	}

	// Create a copy of the created container's hostconfig.  This is passed in during
	// container create
	hostConfig := *vc.HostConfig

	// Resources don't really map well to VIC so we leave most of them empty. If we look
	// at the struct in engine-api/types/container/host_config.go, Microsoft added
	// additional attributes to the struct that are applicable to Windows containers.
	// If understanding VIC's host resources are desirable, we should go down this
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

	hostConfig.VolumeDriver = portlayerName
	hostConfig.Resources = resourceConfig

	if len(info.ScopeConfig) > 0 {
		if info.ScopeConfig[0].DNS != nil {
			hostConfig.DNS = info.ScopeConfig[0].DNS
		}

		hostConfig.NetworkMode = container.NetworkMode(info.ScopeConfig[0].ScopeType)
	}

	// Set this to json-file to force the docker CLI to allow us to use docker logs
	hostConfig.LogConfig.Type = forceLogType

	return &hostConfig
}

// mountsFromContainerInfo()
func mountsFromContainerInfo(vc *viccontainer.VicContainer, info *models.ContainerInfo, portlayerName string) []types.MountPoint {
	if vc == nil || info == nil {
		return nil
	}

	var mounts []types.MountPoint

	for _, vConfig := range info.VolumeConfig {
		// Fill with defaults
		mountConfig := types.MountPoint{
			Destination: "",
			Driver:      portlayerName,
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
// container config that is part of the image metadata AND the overridden container
// config.  The user can override these via the remote API or the docker CLI.
func containerConfigFromContainerInfo(vc *viccontainer.VicContainer, info *models.ContainerInfo) *container.Config {
	if vc == nil || vc.Config == nil || info == nil || info.ContainerConfig == nil || info.ProcessConfig == nil {
		return nil
	}

	// Copy the working copy of our container's config
	container := *vc.Config

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

	if info.ContainerConfig.RepoName != nil {
		container.Image = *info.ContainerConfig.RepoName // Name of the image as it was passed by the operator (eg. could be symbolic)
	}

	// Fill in information about the process
	if info.ProcessConfig.Env != nil {
		container.Env = info.ProcessConfig.Env // List of environment variable to set in the container
	}

	if info.ProcessConfig.WorkingDir != nil {
		container.WorkingDir = *info.ProcessConfig.WorkingDir // Current directory (PWD) in the command will be launched
	}

	// Fill in information about the container network
	if info.ScopeConfig == nil {
		container.NetworkDisabled = true
	} else {
		container.NetworkDisabled = false
		container.MacAddress = ""
		container.ExposedPorts = vc.Config.ExposedPorts
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

	return &container
}

func networkFromContainerInfo(vc *viccontainer.VicContainer, info *models.ContainerInfo) *types.NetworkSettings {
	networks := &types.NetworkSettings{
		NetworkSettingsBase: types.NetworkSettingsBase{
			Bridge:                 "",
			SandboxID:              "",
			HairpinMode:            false,
			LinkLocalIPv6Address:   "",
			LinkLocalIPv6PrefixLen: 0,
			Ports:                  portMapFromVicContainer(vc),
			SandboxKey:             "",
			SecondaryIPAddresses:   nil,
			SecondaryIPv6Addresses: nil,
		},
		Networks: make(map[string]*dnetwork.EndpointSettings),
	}

	return networks
}

// portMapFromVicContainer() constructs a docker portmap from both the container's
// hostconfig and config (both stored in VicContainer).  They are added and modified
// during docker create.  This function creates a new map that is adhere's to docker's
// structure for types.NetworkSettings.Ports.
func portMapFromVicContainer(vc *viccontainer.VicContainer) nat.PortMap {
	var portMap nat.PortMap

	if vc == nil {
		return portMap
	}

	portMap = make(nat.PortMap)

	// Iterate over the hostconfig that was set in docker create.  Get non-nil
	// bindings and fix up ip addr and add to networks
	if vc.HostConfig != nil && vc.HostConfig.PortBindings != nil {
		//		networks.Ports = vc.HostConfig.PortBindings
		for port, portbindings := range vc.HostConfig.PortBindings {

			var newbindings []nat.PortBinding

			for _, binding := range portbindings {
				nb := binding

				// Check host IP.  VIC only support 0.0.0.0
				if nb.HostIP == "" {
					nb.HostIP = "0.0.0.0"
				}

				newbindings = append(newbindings, nb)
			}

			portMap[port] = newbindings
		}
	}

	// Iterate over the container's original image config.  This is the set of
	// exposed ports.  For ports that were not in hostConfig, we assign value of
	// nil.  This appears to be the behavior of regular docker.
	if vc.Config != nil {
		for port := range vc.Config.ExposedPorts {
			if _, ok := portMap[port]; ok {
				continue
			}

			portMap[port] = nil
		}
	}

	return portMap
}
