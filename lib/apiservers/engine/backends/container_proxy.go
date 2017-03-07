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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime"
	rc "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	httpclient "github.com/mreiferson/go-httpclient"

	derr "github.com/docker/docker/api/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/container"
	dnetwork "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/go-connections/nat"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	epoint "github.com/vmware/vic/lib/apiservers/engine/backends/endpoint"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/interaction"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/logging"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/sys"
)

// VicContainerProxy interface
type VicContainerProxy interface {
	CreateContainerHandle(imageID string, config types.ContainerCreateConfig) (string, string, error)
	AddContainerToScope(handle string, config types.ContainerCreateConfig) (string, error)
	AddVolumesToContainer(handle string, config types.ContainerCreateConfig) (string, error)
	AddLoggingToContainer(handle string, config types.ContainerCreateConfig) (string, error)
	AddInteractionToContainer(handle string, config types.ContainerCreateConfig) (string, error)
	CommitContainerHandle(handle, containerID string, waitTime int32) error
	StreamContainerLogs(name string, out io.Writer, started chan struct{}, showTimestamps bool, followLogs bool, since int64, tailLines int64) error

	Stop(vc *viccontainer.VicContainer, name string, seconds *int, unbound bool) error
	IsRunning(vc *viccontainer.VicContainer) (bool, error)
	IsBroken(vc *viccontainer.VicContainer) (bool, error)
	Wait(vc *viccontainer.VicContainer, timeout time.Duration) (exitCode int32, processStatus string, containerState string, reterr error)
	Signal(vc *viccontainer.VicContainer, sig uint64) error
	Resize(vc *viccontainer.VicContainer, height, width int32) error
	AttachStreams(ctx context.Context, vc *viccontainer.VicContainer, clStdin io.ReadCloser, clStdout, clStderr io.Writer, ca *backend.ContainerAttachConfig) error

	Handle(id, name string) (string, error)
	Client() *client.PortLayer
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

const (
	attachConnectTimeout  time.Duration = 15 * time.Second //timeout for the connection
	attachAttemptTimeout  time.Duration = 40 * time.Second //timeout before we ditch an attach attempt
	attachPLAttemptDiff   time.Duration = 10 * time.Second
	attachStdinInitString               = "v1c#>"
	swaggerSubstringEOF                 = "EOF"
	forceLogType                        = "json-file" //Use in inspect to allow docker logs to work
	annotationKeyLabels                 = "docker.labels"
	ShortIDLen                          = 12

	DriverArgFlagKey      = "flags"
	DriverArgContainerKey = "container"
	DriverArgImageKey     = "image"
)

// NewContainerProxy creates a new ContainerProxy
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
	// anonymous volumes can also come from Image Metadata

	rawAnonVolumes := make([]string, 0, len(config.Config.Volumes))
	for k := range config.Config.Volumes {
		rawAnonVolumes = append(rawAnonVolumes, k)
	}

	volList, err := finalizeVolumeList(config.HostConfig.Binds, rawAnonVolumes)
	if err != nil {
		return handle, BadRequestError(err.Error())
	}

	log.Infof("Finalized Volume list : %#v", volList)

	// Create and join volumes.
	for _, fields := range volList {

		//we only set these here for volumes made on a docker create
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

// AddInteractionToContainer adds interaction capabilies to a container, referenced by handle.
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
	inspectJSON, err := c.dockerInfo(vc)
	if err != nil && IsNotFoundError(err) {
		cache.ContainerCache().DeleteContainer(vc.ContainerID)
		return err
	}
	// attempt to stop container if status is running or broken
	if !inspectJSON.State.Running && inspectJSON.State.Status != "broken" {
		return nil
	}

	if unbound {
		unbindParams := scopes.NewUnbindContainerParamsWithContext(ctx).WithHandle(handle)
		ub, err := c.client.Scopes.UnbindContainer(unbindParams)
		if err != nil {
			switch err := err.(type) {
			case *scopes.UnbindContainerNotFound:
				// ignore error
				log.Warnf("Container %s not found by network unbind", vc.ContainerID)
			case *scopes.UnbindContainerInternalServerError:
				return InternalServerError(err.Payload.Message)
			default:
				return InternalServerError(err.Error())
			}
		} else {
			handle = ub.Payload.Handle
		}

		// unmap ports
		if err = UnmapPorts(vc.HostConfig); err != nil {
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
	wait := int32(*seconds)
	err = c.CommitContainerHandle(handle, vc.ContainerID, wait)
	if err != nil {
		if IsNotFoundError(err) {
			cache.ContainerCache().DeleteContainer(vc.ContainerID)
		}
		return err
	}

	return nil
}

// IsRunning returns true if the given container is running
func (c *ContainerProxy) IsRunning(vc *viccontainer.VicContainer) (bool, error) {
	defer trace.End(trace.Begin(""))

	inspectJSON, err := c.dockerInfo(vc)
	if err != nil {
		return false, err
	}
	return inspectJSON.State.Running, nil
}

func (c *ContainerProxy) dockerInfo(vc *viccontainer.VicContainer) (*types.ContainerJSON, error) {
	defer trace.End(trace.Begin(""))

	if c.client == nil {
		return nil, InternalServerError("ContainerProxy.dockerInfo failed to get a portlayer client")
	}

	results, err := c.client.Containers.GetContainerInfo(containers.NewGetContainerInfoParamsWithContext(ctx).WithID(vc.ContainerID))
	if err != nil {
		switch err := err.(type) {
		case *containers.GetContainerInfoNotFound:
			return nil, NotFoundError(fmt.Sprintf("No such container: %s", vc.ContainerID))
		case *containers.GetContainerInfoInternalServerError:
			return nil, InternalServerError(err.Payload.Message)
		default:
			return nil, InternalServerError(fmt.Sprintf("Unknown error from the interaction port layer: %s", err))
		}
	}

	return ContainerInfoToDockerContainerInspect(vc, results.Payload, c.portlayerName)
}

// IsBroken returns true if the given container is broken, e.g. migration failed
func (c *ContainerProxy) IsBroken(vc *viccontainer.VicContainer) (bool, error) {
	defer trace.End(trace.Begin(""))

	inspectJSON, err := c.dockerInfo(vc)
	if err != nil {
		return false, err
	}
	return inspectJSON.State.Status == "broken", nil
}

func (c *ContainerProxy) Wait(vc *viccontainer.VicContainer, timeout time.Duration) (
	exitCode int32, processStatus string, containerState string, reterr error) {

	defer trace.End(trace.Begin(vc.ContainerID))

	if vc == nil {
		reterr = InternalServerError("Wait bad arguments")
		return
	}

	// Get an API client to the portlayer
	client := PortLayerClient()
	if client == nil {
		reterr = InternalServerError("Wait failed to create a portlayer client")
		return
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
			reterr = NotFoundError(vc.ContainerID)
			return
		case *containers.ContainerWaitInternalServerError:
			reterr = InternalServerError(err.Payload.Message)
			return
		default:
			reterr = InternalServerError(err.Error())
			return
		}
	}

	if results == nil || results.Payload == nil {
		reterr = InternalServerError("Unexpected swagger error")
	}

	ci := results.Payload

	return ci.ProcessConfig.ExitCode, ci.ProcessConfig.Status, ci.ContainerConfig.State, nil
}

func (c *ContainerProxy) Signal(vc *viccontainer.VicContainer, sig uint64) error {
	defer trace.End(trace.Begin(vc.ContainerID))

	if vc == nil {
		return InternalServerError("Signal bad arguments")
	}

	// Get an API client to the portlayer
	client := PortLayerClient()
	if client == nil {
		return InternalServerError("Signal failed to create a portlayer client")
	}

	if running, err := c.IsRunning(vc); !running && err == nil {
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

	if running, err := c.IsRunning(vc); !running && err == nil {
		// unmap ports
		if err = UnmapPorts(vc.HostConfig); err != nil {
			return err
		}
	}

	return nil
}

func (c *ContainerProxy) createNewAttachClientWithTimeouts(connectTimeout, responseTimeout, responseHeaderTimeout time.Duration) (*client.PortLayer, *httpclient.Transport) {

	r := rc.New(c.portlayerAddr, "/", []string{"http"})
	transport := &httpclient.Transport{
		ConnectTimeout:        connectTimeout,
		ResponseHeaderTimeout: responseHeaderTimeout,
		RequestTimeout:        responseTimeout,
	}

	r.Transport = transport

	plClient := client.New(r, nil)
	r.Consumers["application/octet-stream"] = runtime.ByteStreamConsumer()
	r.Producers["application/octet-stream"] = runtime.ByteStreamProducer()

	return plClient, transport
}

func (c *ContainerProxy) Resize(vc *viccontainer.VicContainer, height, width int32) error {
	defer trace.End(trace.Begin(vc.ContainerID))

	if c.client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("ContainerProxy failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	plResizeParam := interaction.NewContainerResizeParamsWithContext(ctx).
		WithID(vc.ContainerID).
		WithHeight(height).
		WithWidth(width)

	_, err := c.client.Interaction.ContainerResize(plResizeParam)
	if err != nil {
		if _, isa := err.(*interaction.ContainerResizeNotFound); isa {
			return ResourceNotFoundError(vc.ContainerID, "interaction connection")
		}

		// If we get here, most likely something went wrong with the port layer API server
		return InternalServerError(err.Error())
	}

	return nil
}

// AttachStreams takes the the hijacked connections from the calling client and attaches
// them to the 3 streams from the portlayer's rest server.
// clStdin, clStdout, clStderr are the hijacked connection
func (c *ContainerProxy) AttachStreams(ctx context.Context, vc *viccontainer.VicContainer, clStdin io.ReadCloser, clStdout, clStderr io.Writer, ca *backend.ContainerAttachConfig) error {
	// Cancel will close the child connections.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errors := make(chan error, 3)

	// For stdin, we only have a timeout for connection.  There can be a long duration before
	// the first entry so there is no timeout for response.
	plClient, transport := c.createNewAttachClientWithTimeouts(attachConnectTimeout, 0, attachAttemptTimeout)
	defer transport.Close()

	var keys []byte
	var err error
	if ca.DetachKeys != "" {
		keys, err = term.ToBytes(ca.DetachKeys)
		if err != nil {
			return fmt.Errorf("Invalid escape keys (%s) provided", ca.DetachKeys)
		}
	}

	if ca.UseStdin {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := copyStdIn(ctx, plClient, vc, clStdin, keys)
			if err != nil {
				log.Errorf("container attach: stdin (%s): %s", vc.ContainerID, err.Error())
			} else {
				log.Infof("container attach: stdin (%s) done: %s", vc.ContainerID)
			}

			// no need to take action if we are canceled
			// as that means error happened somewhere else
			if ctx.Err() == context.Canceled {
				log.Infof("returning from stdin as context canceled somewhere else")
				return
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

			// no need to take action if we are canceled
			// as that means error happened somewhere else
			if ctx.Err() == context.Canceled {
				log.Infof("returning from stdin as context canceled somewhere else")
				return
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

			// no need to take action if we are canceled
			// as that means error happened somewhere else
			if ctx.Err() == context.Canceled {
				log.Infof("returning from stdin as context canceled somewhere else")
				return
			}
			cancel()

			errors <- err
		}()
	}

	// Wait for all stream copy to exit
	wg.Wait()

	// close the channel so that we don't leak (if there is an error)/or get blocked (if there are no errors)
	close(errors)

	log.Infof("cleaned up connections to %s. Checking errors", vc.ContainerID)
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

//----------
// Utility Functions
//----------

func dockerContainerCreateParamsToPortlayer(cc types.ContainerCreateConfig, layerID string, imageStore string) *containers.CreateParams {
	config := &models.ContainerCreateConfig{}

	config.NumCpus = cc.HostConfig.CPUCount
	config.MemoryMB = cc.HostConfig.Memory

	// Image
	config.Image = layerID

	// Repo Requested
	config.RepoName = cc.Config.Image

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
	config.Name = cc.Name

	// copy the path
	config.Path = path

	// copy the args
	config.Args = make([]string, len(args))
	copy(config.Args, args)

	// copy the env array
	config.Env = make([]string, len(cc.Config.Env))
	copy(config.Env, cc.Config.Env)

	// image store
	config.ImageStore = &models.ImageStore{Name: imageStore}

	// network
	config.NetworkDisabled = cc.Config.NetworkDisabled

	// working dir
	config.WorkingDir = cc.Config.WorkingDir

	// attach
	config.Attach = cc.Config.AttachStdin || cc.Config.AttachStdout || cc.Config.AttachStderr

	// openstdin
	config.OpenStdin = cc.Config.OpenStdin

	// tty
	config.Tty = cc.Config.Tty

	// container stop signal
	config.StopSignal = cc.Config.StopSignal

	// Stuff the Docker labels into VIC container annotations
	annotationsFromLabels(config, cc.Config.Labels)

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

	// Docker copies Links to NetworkConfig only if it is a UserDefined network, handle that
	if !cc.HostConfig.NetworkMode.IsUserDefined() && len(cc.HostConfig.Links) > 0 {
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

	// Set default container state attributes
	containerState := &types.ContainerState{}

	if info.ProcessConfig != nil {
		containerState.Pid = int(info.ProcessConfig.Pid)
		containerState.ExitCode = int(info.ProcessConfig.ExitCode)
		containerState.Error = info.ProcessConfig.ErrorMsg
		if info.ProcessConfig.StartTime > 0 {
			containerState.StartedAt = time.Unix(info.ProcessConfig.StartTime, 0).Format(time.RFC3339Nano)
		}

		if info.ProcessConfig.StopTime > 0 {
			containerState.FinishedAt = time.Unix(info.ProcessConfig.StopTime, 0).Format(time.RFC3339Nano)
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
		inspectJSON.Path = info.ProcessConfig.ExecPath
		if len(info.ProcessConfig.ExecArgs) > 0 {
			// args[0] is the command and should not appear in the args list here
			inspectJSON.Args = info.ProcessConfig.ExecArgs[1:]
		}
	}

	if info.ContainerConfig != nil {
		containerState.Status = strings.ToLower(info.ContainerConfig.State)

		// https://github.com/docker/docker/blob/master/container/state.go#L77
		if containerState.Status == "stopped" {
			containerState.Status = "exited"
		}
		if containerState.Status == "running" {
			containerState.Running = true
		}

		inspectJSON.Image = info.ContainerConfig.LayerID
		inspectJSON.LogPath = info.ContainerConfig.LogPath
		inspectJSON.RestartCount = int(info.ContainerConfig.RestartCount)
		inspectJSON.ID = info.ContainerConfig.ContainerID
		inspectJSON.Created = time.Unix(info.ContainerConfig.CreateTime, 0).Format(time.RFC3339Nano)
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

	// Set this to json-file to force the docker CLI to allow us to use docker logs
	hostConfig.LogConfig.Type = forceLogType

	var err error
	_, hostConfig.PortBindings, err = nat.ParsePortSpecs(info.HostConfig.Ports)
	if err != nil {
		log.Errorf("Failed to parse port mapping %s: %s", info.HostConfig.Ports, err)
	}

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
		mountConfig.Name = vConfig.MountPoint
		mountConfig.Source = vConfig.MountPoint
		mountConfig.RW = vConfig.ReadWrite

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
		imageConfig, _ = cache.ImageCache().Get(info.ContainerConfig.LayerID)
	}

	// Fill in the values with defaults from the original image's container config
	// structure
	if imageConfig != nil {
		container.StopSignal = imageConfig.ContainerConfig.StopSignal // Signal to stop a container

		container.OnBuild = imageConfig.ContainerConfig.OnBuild // ONBUILD metadata that were defined on the image Dockerfile
	}

	// Pull labels from the annotation
	labelsFromAnnotations(&container, info.ContainerConfig.Annotations)

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

func ContainerInfoToVicContainer(info models.ContainerInfo) *viccontainer.VicContainer {
	log.Debugf("Convert container info to vic container")

	vc := viccontainer.NewVicContainer()

	var name string
	if len(info.ContainerConfig.Names) > 0 {
		vc.Name = info.ContainerConfig.Names[0]
	}
	log.Debugf("Container %q", name)

	if info.ContainerConfig.LayerID != "" {
		vc.ImageID = info.ContainerConfig.LayerID
	}

	if info.ContainerConfig.ContainerID != "" {
		vc.ContainerID = info.ContainerConfig.ContainerID
	}

	tempVC := viccontainer.NewVicContainer()
	tempVC.HostConfig = &container.HostConfig{}
	vc.Config = containerConfigFromContainerInfo(tempVC, &info)
	vc.HostConfig = hostConfigFromContainerInfo(tempVC, &info, PortLayerName())
	return vc
}

// annotationsFromLabels() encodes labels into annotations within the swagger
// create config.  The difference between labels and annotations is that labels
// is specific to Docker.  Annotations is a generic per VIC container k,v struct.
// We store the labels in an annotation key.
func annotationsFromLabels(config *models.ContainerCreateConfig, labels map[string]string) error {
	var err error

	if config == nil || len(labels) == 0 {
		return nil
	}

	if config.Annotations == nil {
		config.Annotations = make(map[string]string)
	}

	// Encoding the labels map into a blob that can be stored as ansi regardless
	// of what encoding the input labels are.  We do this by first marshaling to
	// to a json byte array to get a self describing encoding and then encoding
	// to base64.  We could use another encoding for the self describing part,
	// such as Golang GOB, but this data will be pushed over to a standard REST
	// server so we use standard web standards instead.
	if labelsBytes, merr := json.Marshal(labels); merr == nil {
		labelsBlob := base64.StdEncoding.EncodeToString(labelsBytes)
		config.Annotations[annotationKeyLabels] = labelsBlob
	} else {
		err = merr
		log.Errorf("Unable to marshal docker labels to json: %s", err)
	}

	return err
}

// labelsFromAnnotations() decodes the Docker label value from the VIC annotations.
func labelsFromAnnotations(config *container.Config, annotations map[string]string) error {
	var err error

	if config == nil || len(annotations) == 0 {
		return nil
	}

	if config.Labels == nil {
		config.Labels = make(map[string]string)
	}

	if labelsBlob, ok := annotations[annotationKeyLabels]; ok {
		if labelsBytes, decodeErr := base64.StdEncoding.DecodeString(labelsBlob); decodeErr == nil {
			if err = json.Unmarshal(labelsBytes, &config.Labels); err != nil {
				log.Errorf("Unable to unmarshal docker labels: %s", err)
			}
		} else {
			err = decodeErr
			log.Errorf("Unable to decode container annotations: %s", err)
		}
	}

	return err
}

//------------------------------------
// ContainerAttach() Utility Functions
//------------------------------------

func copyStdIn(ctx context.Context, pl *client.PortLayer, vc *viccontainer.VicContainer, clStdin io.ReadCloser, keys []byte) error {
	// Pipe for stdin so we can interject and watch the input streams for detach keys.
	stdinReader, stdinWriter := io.Pipe()
	defer stdinWriter.Close()

	var detach bool

	go func() {
		defer stdinReader.Close()

		// Copy the stdin from the CLI and write to a pipe.  We need to do this so we can
		// watch the stdin stream for the detach keys.
		var err error

		// Write some init bytes into the pipe to force Swagger to make the initial
		// call to the portlayer, prior to any user input in whatever attach client
		// he/she is using.
		log.Debugf("copyStdIn writing primer bytes")
		stdinWriter.Write([]byte(attachStdinInitString))
		if vc.Config.Tty {
			_, err = copyEscapable(stdinWriter, clStdin, keys)
		} else {
			_, err = io.Copy(stdinWriter, clStdin)
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

	// Swagger wants an io.reader so give it the reader pipe.  Also, the swagger call
	// to set the stdin is synchronous so we need to run in a goroutine
	setStdinParams := interaction.NewContainerSetStdinParamsWithContext(ctx).WithID(vc.ContainerID)
	setStdinParams = setStdinParams.WithRawStream(stdinReader)

	_, err := pl.Interaction.ContainerSetStdin(setStdinParams)
	if vc.Config.StdinOnce && !vc.Config.Tty {
		// Close the stdin connection.  Mimicing Docker's behavior.
		log.Errorf("Attach stream has stdinOnce set.  Closing the stdin.")
		params := interaction.NewContainerCloseStdinParamsWithContext(ctx).WithID(vc.ContainerID)
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

func copyStdOut(ctx context.Context, pl *client.PortLayer, attemptTimeout time.Duration, vc *viccontainer.VicContainer, clStdout io.Writer) error {
	id := vc.ContainerID
	//Calculate how much time to let portlayer attempt
	plAttemptTimeout := attemptTimeout - attachPLAttemptDiff //assumes personality deadline longer than portlayer's deadline
	plAttemptDeadline := time.Now().Add(plAttemptTimeout)
	swaggerDeadline := strfmt.DateTime(plAttemptDeadline)
	log.Debugf("* stdout portlayer deadline: %s", plAttemptDeadline.Format(time.UnixDate))
	log.Debugf("* stdout personality deadline: %s", time.Now().Add(attemptTimeout).Format(time.UnixDate))

	log.Debugf("* stdout attach start %s", time.Now().Format(time.UnixDate))
	getStdoutParams := interaction.NewContainerGetStdoutParamsWithContext(ctx).WithID(id).WithDeadline(&swaggerDeadline)
	_, err := pl.Interaction.ContainerGetStdout(getStdoutParams, clStdout)
	log.Debugf("* stdout attach end %s", time.Now().Format(time.UnixDate))
	if err != nil {
		if _, ok := err.(*interaction.ContainerGetStdoutNotFound); ok {
			return ResourceNotFoundError(id, "interaction connection")
		}

		return InternalServerError(err.Error())
	}

	return nil
}

func copyStdErr(ctx context.Context, pl *client.PortLayer, vc *viccontainer.VicContainer, clStderr io.Writer) error {
	name := vc.ContainerID
	getStderrParams := interaction.NewContainerGetStderrParamsWithContext(ctx).WithID(name)

	_, err := pl.Interaction.ContainerGetStderr(getStderrParams, clStderr)
	if err != nil {
		if _, ok := err.(*interaction.ContainerGetStderrNotFound); ok {
			ResourceNotFoundError(name, "interaction connection")
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
