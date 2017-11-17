// Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"

	derr "github.com/docker/docker/api/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dnetwork "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/go-connections/nat"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	"github.com/vmware/vic/lib/apiservers/engine/backends/convert"
	epoint "github.com/vmware/vic/lib/apiservers/engine/backends/endpoint"
	"github.com/vmware/vic/lib/apiservers/engine/backends/filter"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/interaction"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/logging"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/tasks"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/lib/constants"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/sys"
)

// VicContainerProxy interface
type VicContainerProxy interface {
	CreateContainerHandle(vc *viccontainer.VicContainer, config types.ContainerCreateConfig) (string, string, error)
	AddContainerToScope(handle string, config types.ContainerCreateConfig) (string, error)
	AddVolumesToContainer(handle string, config types.ContainerCreateConfig) (string, error)
	AddLoggingToContainer(handle string, config types.ContainerCreateConfig) (string, error)
	AddInteractionToContainer(handle string, config types.ContainerCreateConfig) (string, error)

	CreateContainerTask(handle string, id string, config types.ContainerCreateConfig) (string, error)
	CreateExecTask(handle string, config *types.ExecConfig) (string, string, error)
	InspectTask(op trace.Operation, handle string, eid string, cid string) (*models.TaskInspectResponse, error)
	BindTask(op trace.Operation, handle string, eid string) (*models.TaskBindResponse, error)

	BindInteraction(handle string, name string, id string) (string, error)
	UnbindInteraction(handle string, name string, id string) (string, error)

	CommitContainerHandle(handle, containerID string, waitTime int32) error
	AttachStreams(ctx context.Context, ac *AttachConfig, stdin io.ReadCloser, stdout, stderr io.Writer) error
	StreamContainerLogs(ctx context.Context, name string, out io.Writer, started chan struct{}, showTimestamps bool, followLogs bool, since int64, tailLines int64) error
	StreamContainerStats(ctx context.Context, config *convert.ContainerStatsConfig) error

	StatPath(op trace.Operation, sotre, deviceID string, filterSpec archive.FilterSpec) (*types.ContainerPathStat, error)

	Stop(vc *viccontainer.VicContainer, name string, seconds *int, unbound bool) error
	State(vc *viccontainer.VicContainer) (*types.ContainerState, error)
	Wait(vc *viccontainer.VicContainer, timeout time.Duration) (*types.ContainerState, error)
	Signal(vc *viccontainer.VicContainer, sig uint64) error
	Resize(id string, height, width int32) error
	Rename(vc *viccontainer.VicContainer, newName string) error
	Remove(vc *viccontainer.VicContainer, config *types.ContainerRmConfig) error

	GetContainerChanges(op trace.Operation, vc *viccontainer.VicContainer, data bool) (io.ReadCloser, error)

	UnbindContainerFromNetwork(vc *viccontainer.VicContainer, handle string) (string, error)

	Handle(id, name string) (string, error)
	Client() *client.PortLayer
	exitCode(vc *viccontainer.VicContainer) (string, error)
}

// ContainerProxy struct
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

// AttachConfig wraps backend.ContainerAttachConfig and adds other required fields
// Similar to https://github.com/docker/docker/blob/master/container/stream/attach.go
type AttachConfig struct {
	*backend.ContainerAttachConfig

	// ID of the session
	ID string
	// Tells the attach copier that the stream's stdin is a TTY and to look for
	// escape sequences in stdin to detach from the stream.
	// When true the escape sequence is not passed to the underlying stream
	UseTty bool
	// CloseStdin signals that once done, stdin for the attached stream should be closed
	// For example, this would close the attached container's stdin.
	CloseStdin bool
}

const (
	attachConnectTimeout  time.Duration = 15 * time.Second //timeout for the connection
	attachAttemptTimeout  time.Duration = 60 * time.Second //timeout before we ditch an attach attempt
	attachPLAttemptDiff   time.Duration = 10 * time.Second
	attachStdinInitString               = "v1c#>"
	swaggerSubstringEOF                 = "EOF"
	forceLogType                        = "json-file" //Use in inspect to allow docker logs to work
	ShortIDLen                          = 12
	archiveStreamBufSize                = 64 * 1024

	DriverArgFlagKey      = "flags"
	DriverArgContainerKey = "container"
	DriverArgImageKey     = "image"

	ContainerRunning = "running"
	ContainerError   = "error"
	ContainerStopped = "stopped"
	ContainerExited  = "exited"
	ContainerCreated = "created"
)

// NewContainerProxy will create a new proxy
func NewContainerProxy(plClient *client.PortLayer, portlayerAddr string, portlayerName string) *ContainerProxy {
	return &ContainerProxy{client: plClient, portlayerAddr: portlayerAddr, portlayerName: portlayerName}
}

// Handle retrieves a handle to a VIC container.  Handles should be treated as opaque strings.
//
// returns:
//	(handle string, error)
func (c *ContainerProxy) Handle(id, name string) (string, error) {
	if c.client == nil {
		return "", InternalServerError("ContainerProxy.Handle failed to get a portlayer client")
	}

	resp, err := c.client.Containers.Get(containers.NewGetParamsWithContext(ctx).WithID(id))
	if err != nil {
		switch err := err.(type) {
		case *containers.GetNotFound:
			cache.ContainerCache().DeleteContainer(id)
			return "", NotFoundError(name)
		case *containers.GetDefault:
			return "", InternalServerError(err.Payload.Message)
		default:
			return "", InternalServerError(err.Error())
		}
	}
	return resp.Payload, nil
}

func (c *ContainerProxy) Client() *client.PortLayer {
	return c.client
}

// CreateContainerHandle creates a new VIC container by calling the portlayer
//
// returns:
//	(containerID, containerHandle, error)
func (c *ContainerProxy) CreateContainerHandle(vc *viccontainer.VicContainer, config types.ContainerCreateConfig) (string, string, error) {
	defer trace.End(trace.Begin(vc.ImageID))

	if c.client == nil {
		return "", "", InternalServerError("ContainerProxy.CreateContainerHandle failed to create a portlayer client")
	}

	if vc.ImageID == "" {
		return "", "", NotFoundError("No image specified")
	}

	if vc.LayerID == "" {
		return "", "", NotFoundError("No layer specified")
	}

	// Call the Exec port layer to create the container
	host, err := sys.UUID()
	if err != nil {
		return "", "", InternalServerError("ContainerProxy.CreateContainerHandle got unexpected error getting VCH UUID")
	}

	plCreateParams := dockerContainerCreateParamsToPortlayer(config, vc, host)
	createResults, err := c.client.Containers.Create(plCreateParams)
	if err != nil {
		if _, ok := err.(*containers.CreateNotFound); ok {
			cerr := fmt.Errorf("No such image: %s", vc.ImageID)
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

// CreateContainerTask sets the primary command to run in the container
//
// returns:
//	(containerHandle, error)
func (c *ContainerProxy) CreateContainerTask(handle, id string, config types.ContainerCreateConfig) (string, error) {
	defer trace.End(trace.Begin(""))

	if c.client == nil {
		return "", InternalServerError("ContainerProxy.CreateContainerTask failed to create a portlayer client")
	}

	plTaskParams := dockerContainerCreateParamsToTask(id, config)
	plTaskParams.Config.Handle = handle

	responseJoin, err := c.client.Tasks.Join(plTaskParams)
	if err != nil {
		log.Errorf("Unable to join primary task to container: %+v", err)
		return "", InternalServerError(err.Error())
	}

	handle, ok := responseJoin.Payload.Handle.(string)
	if !ok {
		return "", InternalServerError(fmt.Sprintf("Type assertion failed on handle from task join: %#+v", handle))
	}

	plBindParams := tasks.NewBindParamsWithContext(ctx).WithConfig(&models.TaskBindConfig{Handle: handle, ID: id})
	responseBind, err := c.client.Tasks.Bind(plBindParams)
	if err != nil {
		log.Errorf("Unable to bind primary task to container: %+v", err)
		return "", InternalServerError(err.Error())
	}

	handle, ok = responseBind.Payload.Handle.(string)
	if !ok {
		return "", InternalServerError(fmt.Sprintf("Type assertion failed on handle from task bind %#+v", handle))
	}

	return handle, nil
}

func (c *ContainerProxy) CreateExecTask(handle string, config *types.ExecConfig) (string, string, error) {
	defer trace.End(trace.Begin(""))

	if c.client == nil {
		return "", "", InternalServerError("ContainerProxy.CreateExecTask failed to create a portlayer client")
	}

	joinconfig := &models.TaskJoinConfig{
		Handle:    handle,
		Path:      config.Cmd[0],
		Args:      config.Cmd[1:],
		Env:       config.Env,
		User:      config.User,
		Attach:    config.AttachStdin || config.AttachStdout || config.AttachStderr,
		OpenStdin: config.AttachStdin,
		Tty:       config.Tty,
	}

	// call Join with JoinParams
	joinparams := tasks.NewJoinParamsWithContext(ctx).WithConfig(joinconfig)
	resp, err := c.client.Tasks.Join(joinparams)
	if err != nil {
		return "", "", InternalServerError(err.Error())
	}
	eid := resp.Payload.ID

	handleprime, ok := resp.Payload.Handle.(string)
	if !ok {
		return "", "", InternalServerError(fmt.Sprintf("Type assertion failed on handle from task bind %#+v", handleprime))
	}

	return handleprime, eid, nil
}

func (c *ContainerProxy) InspectTask(op trace.Operation, handle string, eid string, cid string) (*models.TaskInspectResponse, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("handle(%s), eid(%s), cid(%s)", handle, eid, cid)))

	// inspect the Task to obtain ProcessConfig
	config := &models.TaskInspectConfig{
		Handle: handle,
		ID:     eid,
	}

	// FIXME: right now we are only using this path for exec targets. But later the error messages may need to be changed
	// to be more accurate.
	params := tasks.NewInspectParamsWithContext(ctx).WithConfig(config)
	resp, err := c.client.Tasks.Inspect(params)
	if err != nil {
		switch err := err.(type) {
		case *tasks.InspectNotFound:
			// These error types may need to be expanded. NotFoundError does not fit here.
			op.Errorf("received a TaskNotFound error during task inspect: %s", err.Payload.Message)
			return nil, ConflictError("container (%s) has been poweredoff")
		case *tasks.InspectInternalServerError:
			op.Errorf("received an internal server error during task inspect: %s", err.Payload.Message)
			return nil, InternalServerError(err.Payload.Message)
		case *tasks.InspectConflict:
			op.Errorf("received a conflict error during task inspect: %s", err.Payload.Message)
			return nil, ConflictError(fmt.Sprintf("Cannot complete the operation, container %s has been powered off during execution", cid))
		default:
			return nil, InternalServerError(err.Error())
		}
	}
	return resp.Payload, nil
}

func (c *ContainerProxy) BindTask(op trace.Operation, handle string, eid string) (*models.TaskBindResponse, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("handle(%s), eid(%s)", handle, eid)))

	bindconfig := &models.TaskBindConfig{
		Handle: handle,
		ID:     eid,
	}
	bindparams := tasks.NewBindParamsWithContext(ctx).WithConfig(bindconfig)

	// call Bind with bindparams
	resp, err := c.client.Tasks.Bind(bindparams)
	if err != nil {
		switch err := err.(type) {
		case *tasks.BindNotFound:
			op.Errorf("received TaskNotFound error during task bind: %s", err.Payload.Message)
			return nil, NotFoundError("container (%s) has been poweredoff")
		case *tasks.BindInternalServerError:

			op.Errorf("received unexpected error attempting to bind task(%s) for handle(%s): %s", eid, handle, err.Payload.Message)
			return nil, InternalServerError(err.Payload.Message)
		default:
			op.Errorf("received unexpected error attempting to bind task(%s) for handle(%s): %s", eid, handle, err.Error())
			return nil, InternalServerError(err.Error())
		}

	}

	return resp.Payload, nil
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

	// Volume Attachment Section
	log.Debugf("ContainerProxy.AddVolumesToContainer - VolumeSection")
	log.Debugf("Raw volume arguments: binds:  %#v, volumes: %#v", config.HostConfig.Binds, config.Config.Volumes)

	// Collect all volume mappings. In a docker create/run, they
	// can be anonymous (-v /dir) or specific (-v vol-name:/dir).
	// anonymous volumes can also come from Image Metadata

	rawAnonVolumes := make([]string, 0, len(config.Config.Volumes))
	for k := range config.Config.Volumes {
		rawAnonVolumes = append(rawAnonVolumes, k)
	}

	volList, err := finalizeVolumeList(config.HostConfig.Binds, rawAnonVolumes)
	if err != nil {
		return handle, BadRequestError(err.Error())
	}
	log.Infof("Finalized volume list: %#v", volList)

	if len(config.Config.Volumes) > 0 {
		// override anonymous volume list with generated volume id
		for _, vol := range volList {
			if _, ok := config.Config.Volumes[vol.Dest]; ok {
				delete(config.Config.Volumes, vol.Dest)
				mount := getMountString(vol.ID, vol.Dest, vol.Flags)
				config.Config.Volumes[mount] = struct{}{}
				log.Debugf("Replace anonymous volume config %s with %s", vol.Dest, mount)
			}
		}
	}

	// Create and join volumes.
	for _, fields := range volList {
		// We only set these here for volumes made on a docker create
		volumeData := make(map[string]string)
		volumeData[DriverArgFlagKey] = fields.Flags
		volumeData[DriverArgContainerKey] = config.Name
		volumeData[DriverArgImageKey] = config.Config.Image

		// NOTE: calling volumeCreate regardless of whether the volume is already
		// present can be avoided by adding an extra optional param to VolumeJoin,
		// which would then call volumeCreate if the volume does not exist.
		vol := &Volume{}
		_, err := vol.volumeCreate(fields.ID, "vsphere", volumeData, nil)
		if err != nil {
			switch err := err.(type) {
			case *storage.CreateVolumeConflict:
				// Implicitly ignore the error where a volume with the same name
				// already exists. We can just join the said volume to the container.
				log.Infof("a volume with the name %s already exists", fields.ID)
			case *storage.CreateVolumeNotFound:
				return handle, VolumeCreateNotFoundError(volumeStore(volumeData))
			default:
				return handle, InternalServerError(err.Error())
			}
		} else {
			log.Infof("volumeCreate succeeded. Volume mount section ID: %s", fields.ID)
		}

		flags := make(map[string]string)
		//NOTE: for now we are passing the flags directly through. This is NOT SAFE and only a stop gap.
		flags[constants.Mode] = fields.Flags
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

// AddLoggingToContainer adds logging capability to a container, referenced by handle.
// If an error is return, the returned handle should not be used.
//
// returns:
//	modified handle
func (c *ContainerProxy) AddLoggingToContainer(handle string, config types.ContainerCreateConfig) (string, error) {
	defer trace.End(trace.Begin(handle))

	if c.client == nil {
		return "", InternalServerError("ContainerProxy.AddLoggingToContainer failed to get the portlayer client")
	}

	response, err := c.client.Logging.LoggingJoin(logging.NewLoggingJoinParamsWithContext(ctx).
		WithConfig(&models.LoggingJoinConfig{
			Handle: handle,
		}))
	if err != nil {
		return "", InternalServerError(err.Error())
	}
	handle, ok := response.Payload.Handle.(string)
	if !ok {
		return "", InternalServerError(fmt.Sprintf("Type assertion failed for %#+v", handle))
	}

	return handle, nil
}

// AddInteractionToContainer adds interaction capabilities to a container, referenced by handle.
// If an error is return, the returned handle should not be used.
//
// returns:
//	modified handle
func (c *ContainerProxy) AddInteractionToContainer(handle string, config types.ContainerCreateConfig) (string, error) {
	defer trace.End(trace.Begin(handle))

	if c.client == nil {
		return "", InternalServerError("ContainerProxy.AddInteractionToContainer failed to get the portlayer client")
	}

	response, err := c.client.Interaction.InteractionJoin(interaction.NewInteractionJoinParamsWithContext(ctx).
		WithConfig(&models.InteractionJoinConfig{
			Handle: handle,
		}))
	if err != nil {
		return "", InternalServerError(err.Error())
	}
	handle, ok := response.Payload.Handle.(string)
	if !ok {
		return "", InternalServerError(fmt.Sprintf("Type assertion failed for %#+v", handle))
	}

	return handle, nil
}

// BindInteraction enables interaction capabilities
func (c *ContainerProxy) BindInteraction(handle string, name string, id string) (string, error) {
	defer trace.End(trace.Begin(handle))

	if c.client == nil {
		return "", InternalServerError("ContainerProxy.AddInteractionToContainer failed to get the portlayer client")
	}

	bind, err := c.client.Interaction.InteractionBind(
		interaction.NewInteractionBindParamsWithContext(ctx).
			WithConfig(&models.InteractionBindConfig{
				Handle: handle,
				ID:     id,
			}))
	if err != nil {
		switch err := err.(type) {
		case *interaction.InteractionBindInternalServerError:
			return "", InternalServerError(err.Payload.Message)
		default:
			return "", InternalServerError(err.Error())
		}
	}
	handle, ok := bind.Payload.Handle.(string)
	if !ok {
		return "", InternalServerError(fmt.Sprintf("Type assertion failed for %#+v", handle))
	}
	return handle, nil
}

// UnbindInteraction disables interaction capabilities
func (c *ContainerProxy) UnbindInteraction(handle string, name string, id string) (string, error) {
	defer trace.End(trace.Begin(handle))

	if c.client == nil {
		return "", InternalServerError("ContainerProxy.AddInteractionToContainer failed to get the portlayer client")
	}

	unbind, err := c.client.Interaction.InteractionUnbind(
		interaction.NewInteractionUnbindParamsWithContext(ctx).
			WithConfig(&models.InteractionUnbindConfig{
				Handle: handle,
				ID:     id,
			}))
	if err != nil {
		return "", InternalServerError(err.Error())
	}
	handle, ok := unbind.Payload.Handle.(string)
	if !ok {
		return "", InternalServerError("type assertion failed")
	}

	return handle, nil
}

// CommitContainerHandle commits any changes to container handle.
//
// Args:
//	waitTime <= 0 means no wait time
func (c *ContainerProxy) CommitContainerHandle(handle, containerID string, waitTime int32) error {
	defer trace.End(trace.Begin(handle))

	if c.client == nil {
		return InternalServerError("ContainerProxy.CommitContainerHandle failed to get a portlayer client")
	}

	var commitParams *containers.CommitParams
	if waitTime > 0 {
		commitParams = containers.NewCommitParamsWithContext(ctx).WithHandle(handle).WithWaitTime(&waitTime)
	} else {
		commitParams = containers.NewCommitParamsWithContext(ctx).WithHandle(handle)
	}

	_, err := c.client.Containers.Commit(commitParams)
	if err != nil {
		switch err := err.(type) {
		case *containers.CommitNotFound:
			return NotFoundError(containerID)
		case *containers.CommitConflict:
			return ConflictError(err.Error())
		case *containers.CommitDefault:
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	return nil
}

// StreamContainerLogs reads the log stream from the portlayer rest server and writes
// it directly to the io.Writer that is passed in.
func (c *ContainerProxy) StreamContainerLogs(ctx context.Context, name string, out io.Writer, started chan struct{}, showTimestamps bool, followLogs bool, since int64, tailLines int64) error {
	defer trace.End(trace.Begin(""))

	close(started)

	params := containers.NewGetContainerLogsParamsWithContext(ctx).
		WithID(name).
		WithFollow(&followLogs).
		WithTimestamp(&showTimestamps).
		WithSince(&since).
		WithTaillines(&tailLines)
	_, err := c.client.Containers.GetContainerLogs(params, out)
	if err != nil {
		switch err := err.(type) {
		case *containers.GetContainerLogsNotFound:
			return NotFoundError(name)
		case *containers.GetContainerLogsInternalServerError:
			return InternalServerError("Server error from the interaction port layer")
		default:
			//Check for EOF.  Since the connection, transport, and data handling are
			//encapsulated inside of Swagger, we can only detect EOF by checking the
			//error string
			if strings.Contains(err.Error(), swaggerSubstringEOF) {
				return nil
			}
			return InternalServerError(fmt.Sprintf("Unknown error from the interaction port layer: %s", err))
		}
	}

	return nil
}

// StreamContainerStats will provide a stream of container stats written to the provided
// io.Writer.  Prior to writing to the provided io.Writer there will be a transformation
// from the portLayer representation of stats to the docker format
func (c *ContainerProxy) StreamContainerStats(ctx context.Context, config *convert.ContainerStatsConfig) error {
	defer trace.End(trace.Begin(config.ContainerID))

	// create a child context that we control
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	params := containers.NewGetContainerStatsParamsWithContext(ctx)
	params.ID = config.ContainerID
	params.Stream = config.Stream

	config.Ctx = ctx
	config.Cancel = cancel

	// create our converter
	containerConverter := convert.NewContainerStats(config)
	// provide the writer for the portLayer and start listening for metrics
	writer := containerConverter.Listen()
	if writer == nil {
		// problem with the listener
		return InternalServerError(fmt.Sprintf("unable to gather container(%s) statistics", config.ContainerID))
	}

	_, err := c.client.Containers.GetContainerStats(params, writer)
	if err != nil {
		switch err := err.(type) {
		case *containers.GetContainerStatsNotFound:
			return NotFoundError(config.ContainerID)
		case *containers.GetContainerStatsInternalServerError:
			return InternalServerError("Server error from the interaction port layer")
		default:
			if ctx.Err() == context.Canceled {
				return nil
			}
			//Check for EOF.  Since the connection, transport, and data handling are
			//encapsulated inside of Swagger, we can only detect EOF by checking the
			//error string
			if strings.Contains(err.Error(), swaggerSubstringEOF) {
				return nil
			}
			return InternalServerError(fmt.Sprintf("Unknown error from the interaction port layer: %s", err))
		}
	}
	return nil
}

// GetContainerChanges returns container changes from portlayer.
// Set data to true will return file data, otherwise, only return file headers with change type.
func (c *ContainerProxy) GetContainerChanges(op trace.Operation, vc *viccontainer.VicContainer, data bool) (io.ReadCloser, error) {
	host, err := sys.UUID()
	if err != nil {
		return nil, InternalServerError("Failed to determine host UUID")
	}

	parent := vc.LayerID
	spec := archive.FilterSpec{
		Inclusions: map[string]struct{}{},
		Exclusions: map[string]struct{}{},
	}

	r, err := archiveProxy.ArchiveExportReader(op, constants.ContainerStoreName, host, vc.ContainerID, parent, data, spec)
	if err != nil {
		return nil, InternalServerError(err.Error())
	}

	return r, nil
}

// StatPath requests the portlayer to stat the filesystem resource at the
// specified path in the container vc.
func (c *ContainerProxy) StatPath(op trace.Operation, store, deviceID string, filterSpec archive.FilterSpec) (*types.ContainerPathStat, error) {
	defer trace.End(trace.Begin(deviceID))

	statPathParams := storage.
		NewStatPathParamsWithContext(op).
		WithStore(store).
		WithDeviceID(deviceID)

	spec, err := archive.EncodeFilterSpec(op, &filterSpec)
	if err != nil {
		op.Errorf(err.Error())
		return nil, InternalServerError(err.Error())
	}
	statPathParams = statPathParams.WithFilterSpec(spec)

	statPathOk, err := c.client.Storage.StatPath(statPathParams)
	if err != nil {
		op.Errorf(err.Error())
		return nil, err
	}

	stat := &types.ContainerPathStat{
		Name:       statPathOk.Name,
		Mode:       os.FileMode(statPathOk.Mode),
		Size:       statPathOk.Size,
		LinkTarget: statPathOk.LinkTarget,
	}

	var modTime time.Time
	if err := modTime.GobDecode([]byte(statPathOk.ModTime)); err != nil {
		op.Debugf("error getting mod time from statpath: %s", err.Error())
	} else {
		stat.Mtime = modTime
	}

	return stat, nil
}

// Stop will stop (shutdown) a VIC container.
//
// returns
//	error
func (c *ContainerProxy) Stop(vc *viccontainer.VicContainer, name string, seconds *int, unbound bool) error {
	defer trace.End(trace.Begin(vc.ContainerID))

	if c.client == nil {
		return InternalServerError("ContainerProxy.Stop failed to get a portlayer client")
	}

	//retrieve client to portlayer
	handle, err := c.Handle(vc.ContainerID, name)
	if err != nil {
		return err
	}

	// we have a container on the PL side lets check the state before proceeding
	// ignore the error  since others will be checking below..this is an attempt to short circuit the op
	// TODO: can be replaced with simple cache check once power events are propagated to persona
	state, err := c.State(vc)
	if err != nil && IsNotFoundError(err) {
		cache.ContainerCache().DeleteContainer(vc.ContainerID)
		return err
	}
	// attempt to stop container only if container state is not stopped, exited or created.
	// we should allow user to stop and remove the container that is in unexpected status, e.g. starting, because of serial port connection issue
	if state.Status == ContainerStopped || state.Status == ContainerExited || state.Status == ContainerCreated {
		return nil
	}

	if unbound {
		handle, err = c.UnbindContainerFromNetwork(vc, handle)
		if err != nil {
			return err
		}

		// unmap ports
		if err = UnmapPorts(vc.ContainerID, vc); err != nil {
			return err
		}
	}

	// change the state of the container
	changeParams := containers.NewStateChangeParamsWithContext(ctx).WithHandle(handle).WithState("STOPPED")
	stateChangeResponse, err := c.client.Containers.StateChange(changeParams)
	if err != nil {
		switch err := err.(type) {
		case *containers.StateChangeNotFound:
			cache.ContainerCache().DeleteContainer(vc.ContainerID)
			return NotFoundError(name)
		case *containers.StateChangeDefault:
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	handle = stateChangeResponse.Payload

	// if no timeout in seconds provided then set to default of 10
	if seconds == nil {
		s := 10
		seconds = &s
	}

	err = c.CommitContainerHandle(handle, vc.ContainerID, int32(*seconds))
	if err != nil {
		if IsNotFoundError(err) {
			cache.ContainerCache().DeleteContainer(vc.ContainerID)
		}
		return err
	}

	return nil
}

// UnbindContainerFromNetwork unbinds a container from the networks that it connects to
func (c *ContainerProxy) UnbindContainerFromNetwork(vc *viccontainer.VicContainer, handle string) (string, error) {
	defer trace.End(trace.Begin(vc.ContainerID))

	unbindParams := scopes.NewUnbindContainerParamsWithContext(ctx).WithHandle(handle)
	ub, err := c.client.Scopes.UnbindContainer(unbindParams)
	if err != nil {
		switch err := err.(type) {
		case *scopes.UnbindContainerNotFound:
			// ignore error
			log.Warnf("Container %s not found by network unbind", vc.ContainerID)
		case *scopes.UnbindContainerInternalServerError:
			return "", InternalServerError(err.Payload.Message)
		default:
			return "", InternalServerError(err.Error())
		}
	}

	return ub.Payload.Handle, nil
}

// State returns container state
func (c *ContainerProxy) State(vc *viccontainer.VicContainer) (*types.ContainerState, error) {
	defer trace.End(trace.Begin(""))

	if c.client == nil {
		return nil, InternalServerError("ContainerProxy.State failed to get a portlayer client")
	}

	results, err := c.client.Containers.GetContainerInfo(containers.NewGetContainerInfoParamsWithContext(ctx).WithID(vc.ContainerID))
	if err != nil {
		switch err := err.(type) {
		case *containers.GetContainerInfoNotFound:
			return nil, NotFoundError(vc.Name)
		case *containers.GetContainerInfoInternalServerError:
			return nil, InternalServerError(err.Payload.Message)
		default:
			return nil, InternalServerError(fmt.Sprintf("Unknown error from the interaction port layer: %s", err))
		}
	}

	inspectJSON, err := ContainerInfoToDockerContainerInspect(vc, results.Payload, c.portlayerName)
	if err != nil {
		return nil, err
	}
	return inspectJSON.State, nil
}

// exitCode returns container exitCode
func (c *ContainerProxy) exitCode(vc *viccontainer.VicContainer) (string, error) {
	defer trace.End(trace.Begin(""))

	if c.client == nil {
		return "", InternalServerError("ContainerProxy.exitCode failed to get a portlayer client")
	}

	results, err := c.client.Containers.GetContainerInfo(containers.NewGetContainerInfoParamsWithContext(ctx).WithID(vc.ContainerID))
	if err != nil {
		switch err := err.(type) {
		case *containers.GetContainerInfoNotFound:
			return "", NotFoundError(vc.Name)
		case *containers.GetContainerInfoInternalServerError:
			return "", InternalServerError(err.Payload.Message)
		default:
			return "", InternalServerError(fmt.Sprintf("Unknown error from the interaction port layer: %s", err))
		}
	}
	// get the container state
	dockerState := convert.State(results.Payload)
	if dockerState == nil {
		return "", InternalServerError("Unable to determine container state")
	}

	return strconv.Itoa(dockerState.ExitCode), nil
}

func (c *ContainerProxy) Wait(vc *viccontainer.VicContainer, timeout time.Duration) (
	*types.ContainerState, error) {

	defer trace.End(trace.Begin(vc.ContainerID))

	if vc == nil {
		return nil, InternalServerError("Wait bad arguments")
	}

	// Get an API client to the portlayer
	client := c.client
	if client == nil {
		return nil, InternalServerError("Wait failed to create a portlayer client")
	}

	params := containers.NewContainerWaitParamsWithContext(ctx).
		WithTimeout(int64(timeout.Seconds())).
		WithID(vc.ContainerID)
	results, err := client.Containers.ContainerWait(params)
	if err != nil {
		switch err := err.(type) {
		case *containers.ContainerWaitNotFound:
			// since the container wasn't found on the PL lets remove from the local
			// cache
			cache.ContainerCache().DeleteContainer(vc.ContainerID)
			return nil, NotFoundError(vc.ContainerID)
		case *containers.ContainerWaitInternalServerError:
			return nil, InternalServerError(err.Payload.Message)
		default:
			return nil, InternalServerError(err.Error())
		}
	}

	if results == nil || results.Payload == nil {
		return nil, InternalServerError("Unexpected swagger error")
	}

	dockerState := convert.State(results.Payload)
	if dockerState == nil {
		return nil, InternalServerError("Unable to determine container state")
	}
	return dockerState, nil
}

func (c *ContainerProxy) Signal(vc *viccontainer.VicContainer, sig uint64) error {
	defer trace.End(trace.Begin(vc.ContainerID))

	if vc == nil {
		return InternalServerError("Signal bad arguments")
	}

	// Get an API client to the portlayer
	client := c.client
	if client == nil {
		return InternalServerError("Signal failed to create a portlayer client")
	}

	if state, err := c.State(vc); !state.Running && err == nil {
		return fmt.Errorf("%s is not running", vc.ContainerID)
	}

	// If Docker CLI sends sig == 0, we use sigkill
	if sig == 0 {
		sig = uint64(syscall.SIGKILL)
	}
	params := containers.NewContainerSignalParamsWithContext(ctx).WithID(vc.ContainerID).WithSignal(int64(sig))
	if _, err := client.Containers.ContainerSignal(params); err != nil {
		switch err := err.(type) {
		case *containers.ContainerSignalNotFound:
			return NotFoundError(vc.ContainerID)
		case *containers.ContainerSignalInternalServerError:
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	if state, err := c.State(vc); !state.Running && err == nil {
		// unmap ports
		if err = UnmapPorts(vc.ContainerID, vc); err != nil {
			return err
		}
	}

	return nil
}

func (c *ContainerProxy) Resize(id string, height, width int32) error {
	defer trace.End(trace.Begin(id))

	if c.client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("ContainerProxy failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	plResizeParam := interaction.NewContainerResizeParamsWithContext(ctx).
		WithID(id).
		WithHeight(height).
		WithWidth(width)

	_, err := c.client.Interaction.ContainerResize(plResizeParam)
	if err != nil {
		if _, isa := err.(*interaction.ContainerResizeNotFound); isa {
			return ResourceNotFoundError(id, "interaction connection")
		}

		// If we get here, most likely something went wrong with the port layer API server
		return InternalServerError(err.Error())
	}

	return nil
}

// AttachStreams takes the the hijacked connections from the calling client and attaches
// them to the 3 streams from the portlayer's rest server.
// stdin, stdout, stderr are the hijacked connection
func (c *ContainerProxy) AttachStreams(ctx context.Context, ac *AttachConfig, stdin io.ReadCloser, stdout, stderr io.Writer) error {
	// Cancel will close the child connections.
	var wg, outWg sync.WaitGroup
	errors := make(chan error, 3)

	var keys []byte
	var err error
	if ac.DetachKeys != "" {
		keys, err = term.ToBytes(ac.DetachKeys)
		if err != nil {
			return fmt.Errorf("Invalid escape keys (%s) provided", ac.DetachKeys)
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if ac.UseStdin {
		wg.Add(1)
	}

	if ac.UseStdout {
		wg.Add(1)
		outWg.Add(1)
	}

	if ac.UseStderr {
		wg.Add(1)
		outWg.Add(1)
	}

	// cancel stdin if all output streams are complete
	go func() {
		outWg.Wait()
		cancel()
	}()

	EOForCanceled := func(err error) bool {
		return err != nil && ctx.Err() != context.Canceled && !strings.HasSuffix(err.Error(), swaggerSubstringEOF)
	}

	if ac.UseStdin {
		go func() {
			defer wg.Done()
			err := copyStdIn(ctx, c.client, ac, stdin, keys)
			if err != nil {
				log.Errorf("container attach: stdin (%s): %s", ac.ID, err)
			} else {
				log.Infof("container attach: stdin (%s) done", ac.ID)
			}

			if !ac.CloseStdin || ac.UseTty {
				cancel()
			}

			// Check for EOF or canceled context. We can only detect EOF by checking the error string returned by swagger :/
			if EOForCanceled(err) {
				errors <- err
			}
		}()
	}

	if ac.UseStdout {
		go func() {
			defer outWg.Done()
			defer wg.Done()

			err := copyStdOut(ctx, c.client, ac, stdout, attachAttemptTimeout)
			if err != nil {
				log.Errorf("container attach: stdout (%s): %s", ac.ID, err)
			} else {
				log.Infof("container attach: stdout (%s) done", ac.ID)
			}

			// Check for EOF or canceled context. We can only detect EOF by checking the error string returned by swagger :/
			if EOForCanceled(err) {
				errors <- err
			}
		}()
	}

	if ac.UseStderr {
		go func() {
			defer outWg.Done()
			defer wg.Done()

			err := copyStdErr(ctx, c.client, ac, stderr)
			if err != nil {
				log.Errorf("container attach: stderr (%s): %s", ac.ID, err)
			} else {
				log.Infof("container attach: stderr (%s) done", ac.ID)
			}

			// Check for EOF or canceled context. We can only detect EOF by checking the error string returned by swagger :/
			if EOForCanceled(err) {
				errors <- err
			}
		}()
	}

	// Wait for all stream copy to exit
	wg.Wait()

	// close the channel so that we don't leak (if there is an error)/or get blocked (if there are no errors)
	close(errors)

	log.Infof("cleaned up connections to %s. Checking errors", ac.ID)
	for err := range errors {
		if err != nil {
			// check if we got DetachError
			if _, ok := err.(DetachError); ok {
				log.Infof("Detached from container detected")
				return err
			}

			// If we get here, most likely something went wrong with the port layer API server
			// These errors originate within the go-swagger client itself.
			// Go-swagger returns untyped errors to us if the error is not one that we define
			// in the swagger spec.  Even EOF.  Therefore, we must scan the error string (if there
			// is an error string in the untyped error) for the term EOF.
			log.Errorf("container attach error: %s", err)

			return err
		}
	}

	log.Infof("No error found. Returning nil...")
	return nil
}

// Rename calls the portlayer's RenameContainerHandler to update the container name in the handle,
// and then commit the new name to vSphere
func (c *ContainerProxy) Rename(vc *viccontainer.VicContainer, newName string) error {
	defer trace.End(trace.Begin(vc.ContainerID))

	//retrieve client to portlayer
	handle, err := c.Handle(vc.ContainerID, vc.Name)
	if err != nil {
		return InternalServerError(err.Error())
	}

	if c.client == nil {
		return InternalServerError("ContainerProxy.Rename failed to create a portlayer client")
	}

	// Call the rename functionality in the portlayer.
	renameParams := containers.NewContainerRenameParamsWithContext(ctx).WithName(newName).WithHandle(handle)
	result, err := c.client.Containers.ContainerRename(renameParams)
	if err != nil {
		switch err := err.(type) {
		// Here we don't check the portlayer error type for *containers.ContainerRenameConflict since
		// (1) we already check that in persona cache for ConflictError and
		// (2) the container name in portlayer cache will be updated when committing the handle in the next step
		case *containers.ContainerRenameNotFound:
			return NotFoundError(vc.Name)
		default:
			return InternalServerError(err.Error())
		}
	}

	h := result.Payload

	// commit handle
	_, err = c.client.Containers.Commit(containers.NewCommitParamsWithContext(ctx).WithHandle(h))
	if err != nil {
		switch err := err.(type) {
		case *containers.CommitNotFound:
			return NotFoundError(err.Payload.Message)
		case *containers.CommitConflict:
			return ConflictError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	return nil
}

// Remove calls the portlayer's ContainerRemove handler to remove the container and its
// anonymous volumes if the remove flag is set.
func (c *ContainerProxy) Remove(vc *viccontainer.VicContainer, config *types.ContainerRmConfig) error {
	if c.client == nil {
		return InternalServerError("ContainerProxy.Remove failed to get a portlayer client")
	}

	id := vc.ContainerID
	_, err := c.client.Containers.ContainerRemove(containers.NewContainerRemoveParamsWithContext(ctx).WithID(id))
	if err != nil {
		switch err := err.(type) {
		case *containers.ContainerRemoveNotFound:
			// Remove container from persona cache, but don't return error to the user.
			cache.ContainerCache().DeleteContainer(id)
			return nil
		case *containers.ContainerRemoveDefault:
			return InternalServerError(err.Payload.Message)
		case *containers.ContainerRemoveConflict:
			return derr.NewRequestConflictError(fmt.Errorf("You cannot remove a running container. Stop the container before attempting removal or use -f"))
		case *containers.ContainerRemoveInternalServerError:
			if err.Payload == nil || err.Payload.Message == "" {
				return InternalServerError(err.Error())
			}
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	// Once the container is removed, remove anonymous volumes (vc.Config.Volumes) if
	// the remove flag is set.
	if config.RemoveVolume && len(vc.Config.Volumes) > 0 {
		removeAnonContainerVols(c.client, id, vc)
	}

	return nil
}

//----------
// Utility Functions
//----------

// removeAnonContainerVols removes anonymous volumes joined to a container. It is invoked
// once the said container has been removed. It fetches a list of volumes that are joined
// to at least one other container, and calls the portlayer to remove this container's
// anonymous volumes if they are dangling. Errors, if any, are only logged.
func removeAnonContainerVols(pl *client.PortLayer, cID string, vc *viccontainer.VicContainer) {
	// NOTE: these strings come in the form of <volume id>:<destination>:<volume options>
	volumes := vc.Config.Volumes
	// NOTE: these strings come in the form of <volume id>:<destination path>
	namedVolumes := vc.HostConfig.Binds

	// assemble a mask of volume paths before processing binds. MUST be paths, as we want to move to honoring the proper metadata in the "volumes" section in the future.
	namedMaskList := make(map[string]struct{}, 0)
	for _, entry := range namedVolumes {
		fields := strings.SplitN(entry, ":", 2)
		if len(fields) != 2 {
			log.Errorf("Invalid entry in the HostConfig.Binds metadata section for container %s: %s", cID, entry)
			continue
		}
		destPath := fields[1]
		namedMaskList[destPath] = struct{}{}
	}

	joinedVols, err := fetchJoinedVolumes()
	if err != nil {
		log.Errorf("Unable to obtain joined volumes from portlayer, skipping removal of anonymous volumes for %s: %s", cID, err.Error())
		return
	}

	for vol := range volumes {
		// Extract the volume ID from the full mount path, which is of form "id:mountpath:flags" - see getMountString().
		volFields := strings.SplitN(vol, ":", 3)

		// NOTE(mavery): this check will start to fail when we fix our metadata correctness issues
		if len(volFields) != 3 {
			log.Debugf("Invalid entry in the volumes metadata section for container %s: %s", cID, vol)
			continue
		}
		volName := volFields[0]
		volPath := volFields[1]

		_, isNamed := namedMaskList[volPath]
		_, joined := joinedVols[volName]
		if !joined && !isNamed {
			_, err := pl.Storage.RemoveVolume(storage.NewRemoveVolumeParamsWithContext(ctx).WithName(volName))
			if err != nil {
				log.Debugf("Unable to remove anonymous volume %s in container %s: %s", volName, cID, err.Error())
				continue
			}
			log.Debugf("Successfully removed anonymous volume %s during remove operation against container(%s)", volName, cID)
		}
	}
}

func dockerContainerCreateParamsToTask(id string, cc types.ContainerCreateConfig) *tasks.JoinParams {
	config := &models.TaskJoinConfig{}

	var path string
	var args []string

	// we explicitly specify the ID for the primary task so that it's the same as the containerID
	config.ID = id

	// Expand cmd into entrypoint and args
	cmd := strslice.StrSlice(cc.Config.Cmd)
	if len(cc.Config.Entrypoint) != 0 {
		path, args = cc.Config.Entrypoint[0], append(cc.Config.Entrypoint[1:], cmd...)
	} else {
		path, args = cmd[0], cmd[1:]
	}

	// copy the path
	config.Path = path

	// copy the args
	config.Args = make([]string, len(args))
	copy(config.Args, args)

	// copy the env array
	config.Env = make([]string, len(cc.Config.Env))
	copy(config.Env, cc.Config.Env)

	// working dir
	config.WorkingDir = cc.Config.WorkingDir

	// user
	config.User = cc.Config.User

	// attach.  Always set to true otherwise we cannot attach later.
	// this tells portlayer container is attachable.
	config.Attach = true

	// openstdin
	config.OpenStdin = cc.Config.OpenStdin

	// tty
	config.Tty = cc.Config.Tty

	// container stop signal
	config.StopSignal = cc.Config.StopSignal

	log.Debugf("dockerContainerCreateParamsToTask = %+v", config)

	return tasks.NewJoinParamsWithContext(ctx).WithConfig(config)
}

func dockerContainerCreateParamsToPortlayer(cc types.ContainerCreateConfig, vc *viccontainer.VicContainer, imageStore string) *containers.CreateParams {
	config := &models.ContainerCreateConfig{}

	config.NumCpus = cc.HostConfig.CPUCount
	config.MemoryMB = cc.HostConfig.Memory

	// Layer/vmdk to use
	config.Layer = vc.LayerID

	// Image ID
	config.Image = vc.ImageID

	// Repo Requested
	config.RepoName = cc.Config.Image

	//copy friendly name
	config.Name = cc.Name

	// image store
	config.ImageStore = &models.ImageStore{Name: imageStore}

	// network
	config.NetworkDisabled = cc.Config.NetworkDisabled

	// Stuff the Docker labels into VIC container annotations
	if len(cc.Config.Labels) > 0 {
		convert.SetContainerAnnotation(config, convert.AnnotationKeyLabels, cc.Config.Labels)
	}
	// if autoremove then add to annotation
	if cc.HostConfig.AutoRemove {
		convert.SetContainerAnnotation(config, convert.AnnotationKeyAutoRemove, cc.HostConfig.AutoRemove)
	}

	// hostname
	config.Hostname = cc.Config.Hostname
	// domainname - https://github.com/moby/moby/issues/27067
	config.Domainname = cc.Config.Domainname

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

	// Docker supports link for bridge network and user defined network, we should handle that
	if len(cc.HostConfig.Links) > 0 {
		nc.Aliases = append(nc.Aliases, cc.HostConfig.Links...)
	}

	if cc.NetworkingConfig != nil {
		log.Debugf("EndpointsConfig: %#v", cc.NetworkingConfig)

		es, ok := cc.NetworkingConfig.EndpointsConfig[nc.NetworkName]
		if ok {
			if es.IPAMConfig != nil {
				nc.Address = es.IPAMConfig.IPv4Address
			}

			// Pass Links and Aliases to PL
			nc.Aliases = append(nc.Aliases, epoint.Alias(es)...)
		}
	}

	for p := range cc.HostConfig.PortBindings {
		nc.Ports = append(nc.Ports, fromPortbinding(p, cc.HostConfig.PortBindings[p])...)
	}

	return nc
}

// fromPortbinding translate Port/PortBinding pair to string array with format "hostPort:containerPort/protocol" or
// "containerPort/protocol" if hostPort is empty
// HostIP is ignored here, cause VCH ip address might change. Will query back real interface address in docker ps
func fromPortbinding(port nat.Port, binding []nat.PortBinding) []string {
	var portMappings []string
	if len(binding) == 0 {
		portMappings = append(portMappings, string(port))
		return portMappings
	}

	proto, privatePort := nat.SplitProtoPort(string(port))
	for _, bind := range binding {
		var portMap string
		if bind.HostPort != "" {
			portMap = fmt.Sprintf("%s:%s/%s", bind.HostPort, privatePort, proto)
		} else {
			portMap = string(port)
		}
		portMappings = append(portMappings, portMap)
	}
	return portMappings
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
			return fields, err
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
// It returns a map of unique mountable volumes. This means that it removes dupes favoring
// specified volumes over anonymous volumes.
func processVolumeFields(volumes []string) (map[string]volumeFields, error) {
	volumeFields := make(map[string]volumeFields)

	for _, v := range volumes {
		fields, err := processVolumeParam(v)
		log.Infof("Processed volume arguments: %#v", fields)
		if err != nil {
			return nil, err
		}
		volumeFields[fields.Dest] = fields
	}
	return volumeFields, nil
}

func finalizeVolumeList(specifiedVolumes, anonymousVolumes []string) ([]volumeFields, error) {
	log.Infof("Specified Volumes : %#v", specifiedVolumes)
	processedVolumes, err := processVolumeFields(specifiedVolumes)
	if err != nil {
		return nil, err
	}

	log.Infof("anonymous Volumes : %#v", anonymousVolumes)
	processedAnonVolumes, err := processVolumeFields(anonymousVolumes)
	if err != nil {
		return nil, err
	}

	//combine all volumes, specified volumes are taken over anonymous volumes
	for k, v := range processedVolumes {
		processedAnonVolumes[k] = v
	}

	finalizedVolumes := make([]volumeFields, 0, len(processedAnonVolumes))
	for _, v := range processedAnonVolumes {
		finalizedVolumes = append(finalizedVolumes, v)
	}
	return finalizedVolumes, nil
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
	// get the docker state
	containerState := convert.State(info)

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
			ExecIDs:         vc.List(),
			HostConfig:      hostConfigFromContainerInfo(vc, info, portlayerName),
			GraphDriver:     types.GraphDriverData{Name: portlayerName},
			SizeRw:          nil,
			SizeRootFs:      nil,
		},
		Mounts:          mountsFromContainer(vc),
		Config:          containerConfigFromContainerInfo(vc, info),
		NetworkSettings: networkFromContainerInfo(vc, info),
	}

	if inspectJSON.NetworkSettings != nil {
		log.Debugf("Docker inspect - network settings = %#v", inspectJSON.NetworkSettings)
	} else {
		log.Debug("Docker inspect - network settings = nil")
	}

	if info.ProcessConfig != nil {
		inspectJSON.Path = info.ProcessConfig.ExecPath
		if len(info.ProcessConfig.ExecArgs) > 0 {
			// args[0] is the command and should not appear in the args list here
			inspectJSON.Args = info.ProcessConfig.ExecArgs[1:]
		}
	}

	if info.ContainerConfig != nil {
		// set the status to the inspect expected values
		containerState.Status = filter.DockerState(info.ContainerConfig.State)

		// https://github.com/docker/docker/blob/master/container/state.go#L77
		if containerState.Status == ContainerStopped {
			containerState.Status = ContainerExited
		}

		inspectJSON.Image = info.ContainerConfig.ImageID
		inspectJSON.LogPath = info.ContainerConfig.LogPath
		inspectJSON.RestartCount = int(info.ContainerConfig.RestartCount)
		inspectJSON.ID = info.ContainerConfig.ContainerID
		inspectJSON.Created = time.Unix(0, info.ContainerConfig.CreateTime).Format(time.RFC3339Nano)
		if len(info.ContainerConfig.Names) > 0 {
			inspectJSON.Name = fmt.Sprintf("/%s", info.ContainerConfig.Names[0])
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
	hostConfig.DNS = make([]string, 0)

	if len(info.Endpoints) > 0 {
		for _, ep := range info.Endpoints {
			for _, dns := range ep.Nameservers {
				if dns != "" {
					hostConfig.DNS = append(hostConfig.DNS, dns)
				}
			}
		}

		hostConfig.NetworkMode = container.NetworkMode(info.Endpoints[0].Scope)
	}

	hostConfig.PortBindings = portMapFromContainer(vc, info)

	// Set this to json-file to force the docker CLI to allow us to use docker logs
	hostConfig.LogConfig.Type = forceLogType

	// get the autoremove annotation from the container annotations
	convert.ContainerAnnotation(info.ContainerConfig.Annotations, convert.AnnotationKeyAutoRemove, &hostConfig.AutoRemove)

	return &hostConfig
}

// mountsFromContainer derives []types.MountPoint (used in inspect) from the cached container
// data.
func mountsFromContainer(vc *viccontainer.VicContainer) []types.MountPoint {
	if vc == nil {
		return nil
	}

	var mounts []types.MountPoint

	rawAnonVolumes := make([]string, 0, len(vc.Config.Volumes))
	for k := range vc.Config.Volumes {
		rawAnonVolumes = append(rawAnonVolumes, k)
	}

	volList, err := finalizeVolumeList(vc.HostConfig.Binds, rawAnonVolumes)
	if err != nil {
		return mounts
	}

	for _, vol := range volList {
		mountConfig := types.MountPoint{
			Type:        mount.TypeVolume,
			Driver:      DefaultVolumeDriver,
			Name:        vol.ID,
			Source:      vol.ID,
			Destination: vol.Dest,
			RW:          false,
			Mode:        vol.Flags,
		}

		if strings.Contains(strings.ToLower(vol.Flags), "rw") {
			mountConfig.RW = true
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

	if info.ContainerConfig.ContainerID != "" {
		container.Hostname = stringid.TruncateID(info.ContainerConfig.ContainerID) // Hostname
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
	if info.ContainerConfig.OpenStdin != nil {
		container.OpenStdin = *info.ContainerConfig.OpenStdin
	}
	// They are not coming from PL so set them to true unconditionally
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

	container.User = info.ProcessConfig.User

	// Fill in information about the container network
	if info.Endpoints == nil {
		container.NetworkDisabled = true
	} else {
		container.NetworkDisabled = false
		container.MacAddress = ""
		container.ExposedPorts = vc.Config.ExposedPorts
	}

	// Get the original container config from the image's metadata in our image cache.
	var imageConfig *metadata.ImageConfig

	if info.ContainerConfig.LayerID != "" {
		// #nosec: Errors unhandled.
		imageConfig, _ = cache.ImageCache().Get(info.ContainerConfig.LayerID)
	}

	// Fill in the values with defaults from the original image's container config
	// structure
	if imageConfig != nil {
		container.StopSignal = imageConfig.ContainerConfig.StopSignal // Signal to stop a container

		container.OnBuild = imageConfig.ContainerConfig.OnBuild // ONBUILD metadata that were defined on the image Dockerfile
	}

	// Pull labels from the annotation
	convert.ContainerAnnotation(info.ContainerConfig.Annotations, convert.AnnotationKeyLabels, &container.Labels)
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
			Ports:                  portMapFromContainer(vc, info),
			SandboxKey:             "",
			SecondaryIPAddresses:   nil,
			SecondaryIPv6Addresses: nil,
		},
		Networks: make(map[string]*dnetwork.EndpointSettings),
	}

	shortCID := vc.ContainerID[0:ShortIDLen]

	// Fill in as much info from the endpoint struct inside of the ContainerInfo.
	// The rest of the data must be obtained from the Scopes portlayer.
	for _, ep := range info.Endpoints {
		netEp := &dnetwork.EndpointSettings{
			IPAMConfig:          nil, //Get from Scope PL
			Links:               nil,
			Aliases:             nil,
			NetworkID:           "", //Get from Scope PL
			EndpointID:          ep.ID,
			Gateway:             ep.Gateway,
			IPAddress:           "",
			IPPrefixLen:         0,  //Get from Scope PL
			IPv6Gateway:         "", //Get from Scope PL
			GlobalIPv6Address:   "", //Get from Scope PL
			GlobalIPv6PrefixLen: 0,  //Get from Scope PL
			MacAddress:          "", //Container endpoints currently do not have mac addr yet
		}

		if ep.Address != "" {
			ip, ipnet, err := net.ParseCIDR(ep.Address)
			if err == nil {
				netEp.IPAddress = ip.String()
				netEp.IPPrefixLen, _ = ipnet.Mask.Size()
			}
		}

		if len(ep.Aliases) > 0 {
			netEp.Aliases = make([]string, len(ep.Aliases))
			found := false
			for i, alias := range ep.Aliases {
				netEp.Aliases[i] = alias
				if alias == shortCID {
					found = true
				}
			}

			if !found {
				netEp.Aliases = append(netEp.Aliases, vc.ContainerID[0:ShortIDLen])
			}
		}

		networks.Networks[ep.Scope] = netEp
	}

	return networks
}

// portMapFromContainer constructs a docker portmap from the container's
// info as returned by the portlayer and adds nil entries for any exposed ports
// that are unmapped
func portMapFromContainer(vc *viccontainer.VicContainer, t *models.ContainerInfo) nat.PortMap {
	var mappings nat.PortMap

	if t != nil {
		mappings = addDirectEndpointsToPortMap(t.Endpoints, mappings)
	}
	if vc != nil && vc.Config != nil {
		if vc.NATMap != nil {
			// if there's a NAT map for the container then just use that for the indirect port set
			mappings = mergePortMaps(vc.NATMap, mappings)
		} else {
			// if there's no NAT map then we use the backend data every time
			mappings = addIndirectEndpointsToPortMap(t.Endpoints, mappings)
		}
		mappings = addExposedToPortMap(vc.Config, mappings)
	}

	return mappings
}

// mergePortMaps creates a new map containing the union of the two arguments
func mergePortMaps(map1, map2 nat.PortMap) nat.PortMap {
	resultMap := make(map[nat.Port][]nat.PortBinding)
	for k, v := range map1 {
		resultMap[k] = v
	}

	for k, v := range map2 {
		vr := resultMap[k]
		resultMap[k] = append(vr, v...)
	}

	return resultMap
}

// addIndirectEndpointToPortMap constructs a docker portmap from the container's info as returned by the portlayer for those ports that
// require NAT forward on the endpointVM.
// The portMap provided is modified and returned - the return value should always be used.
func addIndirectEndpointsToPortMap(endpoints []*models.EndpointConfig, portMap nat.PortMap) nat.PortMap {
	if len(endpoints) == 0 {
		return portMap
	}

	// will contain a combined set of port mappings
	if portMap == nil {
		portMap = make(nat.PortMap)
	}

	// add IP address into port spec to allow direct usage of data returned by calls such as docker port
	var ip string
	ips, _ := publicIPv4Addrs()
	if len(ips) > 0 {
		ip = ips[0]
	}

	// Preserve the existing behaviour if we do not have an IP for some reason.
	if ip == "" {
		ip = "0.0.0.0"
	}

	for _, ep := range endpoints {
		if ep.Direct {
			continue
		}

		for _, port := range ep.Ports {
			mappings, err := nat.ParsePortSpec(port)
			if err != nil {
				log.Error(err)
				// just continue if we do have partial port data
			}

			for i := range mappings {
				p := mappings[i].Port
				b := mappings[i].Binding

				if b.HostIP == "" {
					b.HostIP = ip
				}

				if mappings[i].Binding.HostPort == "" {
					// leave this undefined for dynamic assignment
					// TODO: for port stability over VCH restart we would expect to set the dynamically assigned port
					// recorded in containerVM annotations here, so that the old host->port mapping is preserved.
				}

				log.Debugf("Adding indirect mapping for port %v: %v (%s)", p, b, port)

				current, _ := portMap[p]
				portMap[p] = append(current, b)
			}
		}
	}

	return portMap
}

// addDirectEndpointsToPortMap constructs a docker portmap from the container's info as returned by the portlayer for those
// ports exposed directly from the containerVM via container network
// The portMap provided is modified and returned - the return value should always be used.
func addDirectEndpointsToPortMap(endpoints []*models.EndpointConfig, portMap nat.PortMap) nat.PortMap {
	if len(endpoints) == 0 {
		return portMap
	}

	if portMap == nil {
		portMap = make(nat.PortMap)
	}

	for _, ep := range endpoints {
		if !ep.Direct {
			continue
		}

		// add IP address into the port spec to allow direct usage of data returned by calls such as docker port
		var ip string
		rawIP, _, _ := net.ParseCIDR(ep.Address)
		if rawIP != nil {
			ip = rawIP.String()
		}

		if ip == "" {
			ip = "0.0.0.0"
		}

		for _, port := range ep.Ports {
			mappings, err := nat.ParsePortSpec(port)
			if err != nil {
				log.Error(err)
				// just continue if we do have partial port data
			}

			for i := range mappings {
				if mappings[i].Binding.HostIP == "" {
					mappings[i].Binding.HostIP = ip
				}

				if mappings[i].Binding.HostPort == "" {
					// If there's no explicit host port and it's a direct endpoint, then
					// mirror the actual port. It's a bit misleading but we're trying to
					// pack extended function into an existing structure.
					_, p := nat.SplitProtoPort(string(mappings[i].Port))
					mappings[i].Binding.HostPort = p
				}
			}

			for _, mapping := range mappings {
				p := mapping.Port
				current, _ := portMap[p]
				portMap[p] = append(current, mapping.Binding)
			}
		}
	}

	return portMap
}

// addExposedToPortMap ensures that exposed ports are all present in the port map.
// This means nil entries for any exposed ports that are not mapped.
// The portMap provided is modified and returned - the return value should always be used.
func addExposedToPortMap(config *container.Config, portMap nat.PortMap) nat.PortMap {
	if config == nil || len(config.ExposedPorts) == 0 {
		return portMap
	}

	if portMap == nil {
		portMap = make(nat.PortMap)
	}

	for p := range config.ExposedPorts {
		if _, ok := portMap[p]; ok {
			continue
		}

		portMap[p] = nil
	}

	return portMap
}

func ContainerInfoToVicContainer(info models.ContainerInfo) *viccontainer.VicContainer {
	vc := viccontainer.NewVicContainer()

	if info.ContainerConfig.ContainerID != "" {
		vc.ContainerID = info.ContainerConfig.ContainerID
	}

	log.Debugf("Convert container info to vic container: %s", vc.ContainerID)

	if len(info.ContainerConfig.Names) > 0 {
		vc.Name = info.ContainerConfig.Names[0]
		log.Debugf("Container %q", vc.Name)
	}

	if info.ContainerConfig.LayerID != "" {
		vc.LayerID = info.ContainerConfig.LayerID
	}

	if info.ContainerConfig.ImageID != "" {
		vc.ImageID = info.ContainerConfig.ImageID
	}

	tempVC := viccontainer.NewVicContainer()
	tempVC.HostConfig = &container.HostConfig{}
	vc.Config = containerConfigFromContainerInfo(tempVC, &info)
	vc.HostConfig = hostConfigFromContainerInfo(tempVC, &info, PortLayerName())

	// FIXME: duplicate Config.Volumes and HostConfig.Binds here for can not derive them from persisted value right now.
	// get volumes from volume config
	vc.Config.Volumes = make(map[string]struct{}, len(info.VolumeConfig))
	vc.HostConfig.Binds = []string{}
	for _, volume := range info.VolumeConfig {
		mount := getMountString(volume.Name, volume.MountPoint, volume.Flags[constants.Mode])
		vc.Config.Volumes[mount] = struct{}{}
		vc.HostConfig.Binds = append(vc.HostConfig.Binds, mount)
		log.Debugf("add volume mount %s to config.volumes and hostconfig.binds", mount)
	}

	vc.Config.Cmd = info.ProcessConfig.ExecArgs

	return vc
}

// getMountString returns a colon-delimited string containing a volume's name/ID, mount
// point and flags.
func getMountString(mounts ...string) string {
	return strings.Join(mounts, ":")
}

//------------------------------------
// ContainerAttach() Utility Functions
//------------------------------------

func copyStdIn(ctx context.Context, pl *client.PortLayer, ac *AttachConfig, stdin io.ReadCloser, keys []byte) error {
	// Pipe for stdin so we can interject and watch the input streams for detach keys.
	stdinReader, stdinWriter := io.Pipe()
	defer stdinReader.Close()

	var detach bool

	done := make(chan struct{})
	go func() {
		// make sure we get out of io.Copy if context is canceled
		select {
		case <-ctx.Done():
			// This will cause the transport to the API client to be shut down, so all output
			// streams will get closed as well.
			// See the closer in container_routes.go:postContainersAttach

			// We're closing this here to disrupt the io.Copy below
			// TODO: seems like we should be providing an io.Copy impl with ctx argument that honors
			// cancelation with the amount of code dedicated to working around it

			// TODO: I think this still leaves a race between closing of the API client transport and
			// copying of the output streams, it's just likely the error will be dropped as the transport is
			// closed when it occurs.
			// We should move away from needing to close transports to interrupt reads.
			stdin.Close()
		case <-done:
		}
	}()

	go func() {
		defer close(done)
		defer stdinWriter.Close()

		// Copy the stdin from the CLI and write to a pipe.  We need to do this so we can
		// watch the stdin stream for the detach keys.
		var err error

		// Write some init bytes into the pipe to force Swagger to make the initial
		// call to the portlayer, prior to any user input in whatever attach client
		// he/she is using.
		log.Debugf("copyStdIn writing primer bytes")
		stdinWriter.Write([]byte(attachStdinInitString))
		if ac.UseTty {
			_, err = copyEscapable(stdinWriter, stdin, keys)
		} else {
			_, err = io.Copy(stdinWriter, stdin)
		}

		if err != nil {
			if _, ok := err.(DetachError); ok {
				log.Infof("stdin detach detected")
				detach = true
			} else {
				log.Errorf("stdin err: %s", err)
			}
		}
	}()

	id := ac.ID

	// Swagger wants an io.reader so give it the reader pipe.  Also, the swagger call
	// to set the stdin is synchronous so we need to run in a goroutine
	setStdinParams := interaction.NewContainerSetStdinParamsWithContext(ctx).WithID(id)
	setStdinParams = setStdinParams.WithRawStream(stdinReader)

	_, err := pl.Interaction.ContainerSetStdin(setStdinParams)
	<-done

	if ac.CloseStdin && !ac.UseTty {
		// Close the stdin connection.  Mimicing Docker's behavior.
		log.Errorf("Attach stream has stdinOnce set.  Closing the stdin.")
		params := interaction.NewContainerCloseStdinParamsWithContext(ctx).WithID(id)
		_, err := pl.Interaction.ContainerCloseStdin(params)
		if err != nil {
			log.Errorf("CloseStdin failed with %s", err)
		}
	}

	// ignore the portlayer error when it is DetachError as that is what we should return to the caller when we detach
	if detach {
		return DetachError{}
	}

	return err
}

func copyStdOut(ctx context.Context, pl *client.PortLayer, ac *AttachConfig, stdout io.Writer, attemptTimeout time.Duration) error {
	id := ac.ID

	//Calculate how much time to let portlayer attempt
	plAttemptTimeout := attemptTimeout - attachPLAttemptDiff //assumes personality deadline longer than portlayer's deadline
	plAttemptDeadline := time.Now().Add(plAttemptTimeout)
	swaggerDeadline := strfmt.DateTime(plAttemptDeadline)
	log.Debugf("* stdout portlayer deadline: %s", plAttemptDeadline.Format(time.UnixDate))
	log.Debugf("* stdout personality deadline: %s", time.Now().Add(attemptTimeout).Format(time.UnixDate))

	log.Debugf("* stdout attach start %s", time.Now().Format(time.UnixDate))
	getStdoutParams := interaction.NewContainerGetStdoutParamsWithContext(ctx).WithID(id).WithDeadline(&swaggerDeadline)
	_, err := pl.Interaction.ContainerGetStdout(getStdoutParams, stdout)
	log.Debugf("* stdout attach end %s", time.Now().Format(time.UnixDate))
	if err != nil {
		if _, ok := err.(*interaction.ContainerGetStdoutNotFound); ok {
			return ResourceNotFoundError(id, "interaction connection")
		}

		return InternalServerError(err.Error())
	}

	return nil
}

func copyStdErr(ctx context.Context, pl *client.PortLayer, ac *AttachConfig, stderr io.Writer) error {
	id := ac.ID

	getStderrParams := interaction.NewContainerGetStderrParamsWithContext(ctx).WithID(id)
	_, err := pl.Interaction.ContainerGetStderr(getStderrParams, stderr)
	if err != nil {
		if _, ok := err.(*interaction.ContainerGetStderrNotFound); ok {
			ResourceNotFoundError(id, "interaction connection")
		}

		return InternalServerError(err.Error())
	}

	return nil
}

// FIXME: Move this function to a pkg to show it's origination from Docker once
// we have ignore capabilities in our header-check.sh that checks for copyright
// header.
// Code c/c from io.Copy() modified by Docker to handle escape sequence
// Begin

// DetachError is special error which returned in case of container detach.
type DetachError struct{}

func (DetachError) Error() string {
	return "detached from container"
}

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
			preservBuf := []byte{}
			for i, key := range keys {
				preservBuf = append(preservBuf, buf[0:nr]...)
				if nr != 1 || buf[0] != key {
					break
				}
				if i == len(keys)-1 {
					src.Close()
					return 0, DetachError{}
				}
				nr, er = src.Read(buf)
			}
			var nw int
			var ew error
			if len(preservBuf) > 0 {
				nw, ew = dst.Write(preservBuf)
				nr = len(preservBuf)
			} else {
				// ---- End of docker
				nw, ew = dst.Write(buf[0:nr])
			}
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

// End
