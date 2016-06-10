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
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/go-swagger/go-swagger/httpkit"
	httptransport "github.com/go-swagger/go-swagger/httpkit/client"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/backend"
	derr "github.com/docker/docker/errors"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/version"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"

	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/interaction"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/pkg/trace"
)

type Container struct {
	ProductName string
}

const attachRequestTimeout time.Duration = 2 * time.Hour

// docker's container.execBackend

// ContainerExecCreate sets up an exec in a running container.
func (c *Container) ContainerExecCreate(config *types.ExecConfig) (string, error) {
	return "", fmt.Errorf("%s does not implement container.ContainerExecCreate", c.ProductName)
}

// ContainerExecInspect returns low-level information about the exec
// command. An error is returned if the exec cannot be found.
func (c *Container) ContainerExecInspect(id string) (*backend.ExecInspect, error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerExecInspect", c.ProductName)
}

// ContainerExecResize changes the size of the TTY of the process
// running in the exec with the given name to the given height and
// width.
func (c *Container) ContainerExecResize(name string, height, width int) error {
	return fmt.Errorf("%s does not implement container.ContainerExecResize", c.ProductName)
}

// ContainerExecStart starts a previously set up exec instance. The
// std streams are set up.
func (c *Container) ContainerExecStart(name string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer) error {
	return fmt.Errorf("%s does not implement container.ContainerExecStart", c.ProductName)
}

// ExecExists looks up the exec instance and returns a bool if it exists or not.
// It will also return the error produced by `getConfig`
func (c *Container) ExecExists(name string) (bool, error) {
	return false, fmt.Errorf("%s does not implement container.ExecExists", c.ProductName)
}

// docker's container.copyBackend

// ContainerArchivePath creates an archive of the filesystem resource at the
// specified path in the container identified by the given name. Returns a
// tar archive of the resource and whether it was a directory or a single file.
func (c *Container) ContainerArchivePath(name string, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error) {
	return nil, nil, fmt.Errorf("%s does not implement container.ContainerArchivePath", c.ProductName)
}

// ContainerCopy performs a deprecated operation of archiving the resource at
// the specified path in the container identified by the given name.
func (c *Container) ContainerCopy(name string, res string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerCopy", c.ProductName)
}

// ContainerExport writes the contents of the container to the given
// writer. An error is returned if the container cannot be found.
func (c *Container) ContainerExport(name string, out io.Writer) error {
	return fmt.Errorf("%s does not implement container.ContainerExport", c.ProductName)
}

// ContainerExtractToDir extracts the given archive to the specified location
// in the filesystem of the container identified by the given name. The given
// path must be of a directory in the container. If it is not, the error will
// be ErrExtractPointNotDirectory. If noOverwriteDirNonDir is true then it will
// be an error if unpacking the given content would cause an existing directory
// to be replaced with a non-directory and vice versa.
func (c *Container) ContainerExtractToDir(name, path string, noOverwriteDirNonDir bool, content io.Reader) error {
	return fmt.Errorf("%s does not implement container.ContainerExtractToDir", c.ProductName)
}

// ContainerStatPath stats the filesystem resource at the specified path in the
// container identified by the given name.
func (c *Container) ContainerStatPath(name string, path string) (stat *types.ContainerPathStat, err error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerStatPath", c.ProductName)
}

// docker's container.stateBackend

// ContainerCreate creates a container.
func (c *Container) ContainerCreate(config types.ContainerCreateConfig) (types.ContainerCreateResponse, error) {
	defer trace.End(trace.Begin("ContainerCreate"))

	var err error

	//TODO: validate the config parameters
	log.Printf("config.Config = %+v", config.Config)

	// Get an API client to the portlayer
	client := PortLayerClient()
	if client == nil {
		return types.ContainerCreateResponse{},
			derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate failed to create a portlayer client"),
				http.StatusInternalServerError)
	}

	var layer *viccontainer.VicContainer
	container := viccontainer.GetCache().GetContainerByName(config.Config.Image)
	if container == nil {
		layer, err = c.getImageMetadataFromStoragePL(config.Config.Image)

		if err != nil {
			return types.ContainerCreateResponse{}, err
		}
	}

	// Overwrite or append the image's config from the CLI with the metadata from the image's
	// layer metadata where appropriate
	if len(config.Config.Cmd) == 0 {
		config.Config.Cmd = layer.Config.Cmd
	}
	if config.Config.WorkingDir == "" {
		config.Config.WorkingDir = layer.Config.WorkingDir
	}
	if len(config.Config.Entrypoint) == 0 {
		config.Config.Entrypoint = layer.Config.Entrypoint
	}
	// Set TERM to xterm if tty is set
	if config.Config.Tty {
		config.Config.Env = append(config.Config.Env, "TERM=xterm")
	}
	config.Config.Env = append(config.Config.Env, layer.Config.Env...)

	// Was a name provided - if not create a friendly name
	if config.Name == "" {
		//TODO: Assume we could have a name collison here : need to
		// provide validation / retry CDG June 9th 2016
		config.Name = namesgenerator.GetRandomName(0)
	}

	log.Printf("ContainerCreate config' = %+v", config)
	// Call the Exec port layer to create the container
	host, err := guest.UUID()
	if err != nil {
		return types.ContainerCreateResponse{},
			derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate got unexpected error getting VCH UUID"),
				http.StatusInternalServerError)
	}

	plCreateParams := c.dockerContainerCreateParamsToPortlayer(config, layer.ID, host)
	createResults, err := client.Containers.Create(plCreateParams)
	// transfer port layer swagger based response to Docker backend data structs and return to the REST front-end
	if err != nil {
		if _, ok := err.(*containers.CreateNotFound); ok {
			return types.ContainerCreateResponse{}, derr.NewRequestNotFoundError(fmt.Errorf("No such image: %s", layer.ID))
		}

		// If we get here, most likely something went wrong with the port layer API server
		return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
	}

	id := createResults.Payload.ID
	h := createResults.Payload.Handle

	// configure networking
	netConf := toModelsNetworkConfig(config)
	if netConf != nil {
		addContRes, err := client.Scopes.AddContainer(
			scopes.NewAddContainerParams().WithHandle(h).WithNetworkConfig(netConf))
		if err != nil {
			return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
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

	// commit the create op
	_, err = client.Containers.Commit(containers.NewCommitParams().WithHandle(h))
	if err != nil {
		// FIXME: Containers.Commit returns more errors than it's swagger spec says.
		// When no image exist, it also sends back non swagger errors.  We should fix
		// this once Commit returns correct error codes.
		return types.ContainerCreateResponse{}, derr.NewRequestNotFoundError(fmt.Errorf("No such image: %s", layer.ID))
	}

	// Container created ok, overwrite the container params in the container store as
	// these are the parameters that the containers were actually created with
	layer.Config.Cmd = config.Config.Cmd
	layer.Config.WorkingDir = config.Config.WorkingDir
	layer.Config.Entrypoint = config.Config.Entrypoint
	layer.Config.Env = config.Config.Env
	layer.Config.AttachStdin = config.Config.AttachStdin
	layer.Config.AttachStdout = config.Config.AttachStdout
	layer.Config.AttachStderr = config.Config.AttachStderr
	layer.Config.Tty = config.Config.Tty
	layer.Config.OpenStdin = config.Config.OpenStdin
	layer.Config.StdinOnce = config.Config.StdinOnce
	layer.ContainerID = createResults.Payload.ID

	log.Debugf("Container create: %#v", layer)
	viccontainer.GetCache().SaveContainer(createResults.Payload.ID, layer)

	// Success!
	log.Printf("container.ContainerCreate succeeded.  Returning container handle %s", *createResults.Payload)
	return types.ContainerCreateResponse{ID: id}, nil
}

// ContainerKill sends signal to the container
// If no signal is given (sig 0), then Kill with SIGKILL and wait
// for the container to exit.
// If a signal is given, then just send it to the container and return.
func (c *Container) ContainerKill(name string, sig uint64) error {
	return fmt.Errorf("%s does not implement container.ContainerKill", c.ProductName)
}

// ContainerPause pauses a container
func (c *Container) ContainerPause(name string) error {
	return fmt.Errorf("%s does not implement container.ContainerPause", c.ProductName)
}

// ContainerRename changes the name of a container, using the oldName
// to find the container. An error is returned if newName is already
// reserved.
func (c *Container) ContainerRename(oldName, newName string) error {
	return fmt.Errorf("%s does not implement container.ContainerRename", c.ProductName)
}

// ContainerResize changes the size of the TTY of the process running
// in the container with the given name to the given height and width.
func (c *Container) ContainerResize(name string, height, width int) error {
	defer trace.End(trace.Begin("ContainerResize"))

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

// ContainerRestart stops and starts a container. It attempts to
// gracefully stop the container within the given timeout, forcefully
// stopping it if the timeout is exceeded. If given a negative
// timeout, ContainerRestart will wait forever until a graceful
// stop. Returns an error if the container cannot be found, or if
// there is an underlying error at any stage of the restart.
func (c *Container) ContainerRestart(name string, seconds int) error {
	return fmt.Errorf("%s does not implement container.ContainerRestart", c.ProductName)
}

// ContainerRm removes the container id from the filesystem. An error
// is returned if the container is not found, or if the remove
// fails. If the remove succeeds, the container name is released, and
// network links are removed.
func (c *Container) ContainerRm(name string, config *types.ContainerRmConfig) error {
	defer trace.End(trace.Begin("ContainerRm"))

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

// ContainerStart starts a container.
func (c *Container) ContainerStart(name string, hostConfig *container.HostConfig) error {
	defer trace.End(trace.Begin("ContainerStart"))

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

	// bind network
	bindRes, err := client.Scopes.BindContainer(scopes.NewBindContainerParams().WithHandle(h))
	if err != nil {
		if _, ok := err.(*scopes.BindContainerNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer"))
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer"), http.StatusInternalServerError)
	}

	defer func() {
		if err != nil {
			// roll back the BindContainer call
			if _, err = client.Scopes.UnbindContainer(scopes.NewUnbindContainerParams().WithHandle(h)); err != nil {
				log.Warnf("failed to roll back container bind: %s", err.Error())
			}
		}
	}()

	h = bindRes.Payload

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

// ContainerStop looks for the given container and terminates it,
// waiting the given number of seconds before forcefully killing the
// container. If a negative number of seconds is given, ContainerStop
// will wait for a graceful termination. An error is returned if the
// container is not found, is already stopped, or if there is a
// problem stopping the container.
func (c *Container) ContainerStop(name string, seconds int) error {
	defer trace.End(trace.Begin("ContainerStop"))
	//retrieve client to portlayer
	client := PortLayerClient()
	if client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	getResponse, err := client.Containers.Get(containers.NewGetParams().WithID(name))
	if err != nil {
		if _, ok := err.(*containers.GetNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer"), http.StatusInternalServerError)
	}

	handle := getResponse.Payload

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

// ContainerUnpause unpauses a container
func (c *Container) ContainerUnpause(name string) error {
	return fmt.Errorf("%s does not implement container.ContainerUnpause", c.ProductName)
}

// ContainerUpdate updates configuration of the container
func (c *Container) ContainerUpdate(name string, hostConfig *container.HostConfig) ([]string, error) {
	return make([]string, 0, 0), fmt.Errorf("%s does not implement container.ContainerUpdate", c.ProductName)
}

// ContainerWait stops processing until the given container is
// stopped. If the container is not found, an error is returned. On a
// successful stop, the exit code of the container is returned. On a
// timeout, an error is returned. If you want to wait forever, supply
// a negative duration for the timeout.
func (c *Container) ContainerWait(name string, timeout time.Duration) (int, error) {
	return 0, fmt.Errorf("%s does not implement container.ContainerWait", c.ProductName)
}

// docker's container.monitorBackend

// ContainerChanges returns a list of container fs changes
func (c *Container) ContainerChanges(name string) ([]archive.Change, error) {
	return make([]archive.Change, 0, 0), fmt.Errorf("%s does not implement container.ContainerChanges", c.ProductName)
}

// ContainerInspect returns low-level information about a
// container. Returns an error if the container cannot be found, or if
// there is an error getting the data.
func (c *Container) ContainerInspect(name string, size bool, version version.Version) (interface{}, error) {
	//Ignore version.  We're supporting post-1.20 version.

	defer trace.End(trace.Begin("ContainerInspect"))

	// Look up the container info in the metadata cache
	vc := viccontainer.GetCache().GetContainerByName(name)
	if vc == nil {
		return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
	}

	// FIXME: Just set a bunch of dummy info needed to get docker attach working.
	// Once we have event stream working (part of docker ps work), we need to go
	// back and replace this.
	base := &types.ContainerJSONBase{
		State: &types.ContainerState{Status: "running",
			Running: true,
			Paused:  false,
		},
	}

	conJSON := &types.ContainerJSON{ContainerJSONBase: base, Config: vc.Config}

	log.Debugf("ContainerInspect json config = %+v\n", conJSON.Config)

	return conJSON, nil
}

// ContainerLogs hooks up a container's stdout and stderr streams
// configured with the given struct.
func (c *Container) ContainerLogs(name string, config *backend.ContainerLogsConfig, started chan struct{}) error {
	return fmt.Errorf("%s does not implement container.ContainerLogs", c.ProductName)
}

// ContainerStats writes information about the container to the stream
// given in the config object.
func (c *Container) ContainerStats(name string, config *backend.ContainerStatsConfig) error {
	return fmt.Errorf("%s does not implement container.ContainerStats", c.ProductName)
}

// ContainerTop lists the processes running inside of the given
// container by calling ps with the given args, or with the flags
// "-ef" if no args are given.  An error is returned if the container
// is not found, or is not running, or if there are any problems
// running ps, or parsing the output.
func (c *Container) ContainerTop(name string, psArgs string) (*types.ContainerProcessList, error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerTop", c.ProductName)
}

// Containers returns the list of containers to show given the user's filtering.
func (c *Container) Containers(config *types.ContainerListOptions) ([]*types.Container, error) {
	return nil, fmt.Errorf("%s does not implement container.Containers", c.ProductName)
}

// docker's container.attachBackend

// ContainerAttach attaches to logs according to the config passed in. See ContainerAttachConfig.
func (c *Container) ContainerAttach(prefixOrName string, ca *backend.ContainerAttachConfig) error {
	defer trace.End(trace.Begin("ContainerAttach"))

	vc := viccontainer.GetCache().GetContainerByName(prefixOrName)

	if vc == nil {
		//FIXME: If we didn't find in the cache, we should goto the port layer and
		//see if it exists there.  API server might have been bounced.  For now,
		//just return error

		return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", prefixOrName))
	}

	clStdin, clStdout, clStderr, err := ca.GetStreams()

	if err != nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("Unable to get stdio streams for calling client"), http.StatusInternalServerError)
	}

	if !ca.UseStdin {
		clStdin = nil
	}
	if !ca.UseStdout {
		clStdout = nil
	}
	if !ca.UseStderr {
		clStderr = nil
	}

	if !vc.Config.Tty && ca.MuxStreams {
		// replace the stdout/stderr with Docker's multiplex stream
		if ca.UseStdout {
			clStderr = stdcopy.NewStdWriter(clStderr, stdcopy.Stderr)
		}
		if ca.UseStderr {
			clStdout = stdcopy.NewStdWriter(clStdout, stdcopy.Stdout)
		}
	}

	err = attachStreams(vc.ContainerID, vc.Config.Tty, vc.Config.StdinOnce, clStdin, clStdout, clStderr, ca.DetachKeys)

	return err
}

//----------
// Utility Functions
//----------

func (c *Container) dockerContainerCreateParamsToPortlayer(cc types.ContainerCreateConfig, layerID string, imageStore string) *containers.CreateParams {
	config := &models.ContainerCreateConfig{}

	// Image
	config.Image = new(string)
	*config.Image = layerID

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

	log.Printf("dockerContainerCreateParamsToPortlayer = %+v", config)

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
		if es, ok := cc.NetworkingConfig.EndpointsConfig[nc.NetworkName]; ok {
			if es.IPAMConfig != nil {
				nc.Address = &es.IPAMConfig.IPv4Address
			}
		}
	}

	return nc
}

func (c *Container) imageExist(imageID string) (storeName string, err error) {
	// Call the storage port layer to determine if the image currently exist
	host, err := guest.UUID()
	if err != nil {
		return "", derr.NewBadRequestError(fmt.Errorf("container.ContainerCreate got unexpected error getting VCH UUID"))
	}

	getParams := storage.NewGetImageParams().WithID(imageID).WithStoreName(host)
	if _, err := PortLayerClient().Storage.GetImage(getParams); err != nil {
		// If the image does not exist
		if _, ok := err.(*storage.GetImageNotFound); ok {
			// return error and "No such image" which the client looks for to determine if the image didn't exist
			return "", derr.NewRequestNotFoundError(fmt.Errorf("No such image: %s", imageID))
		}

		// If we get here, most likely something went wrong with the port layer API server
		return "", derr.NewErrorWithStatusCode(fmt.Errorf("Unknown error from the storage portlayer"),
			http.StatusInternalServerError)
	}

	return host, nil
}

func (c *Container) getImageMetadataFromStoragePL(image string) (*viccontainer.VicContainer, error) {
	// FIXME: This is a temporary workaround until we have a name resolution story.
	// Call imagec with -resolv parameter to learn the name of the vmdk and put it into in-memory map
	cmdArgs := []string{"-reference", image, "-resolv", "-standalone", "-destination", os.TempDir()}

	out, err := exec.Command(Imagec, cmdArgs...).Output()
	if err != nil {
		log.Printf("%s exit code: %s", Imagec, err)
		return nil,
			derr.NewErrorWithStatusCode(fmt.Errorf("Container look up failed"),
				http.StatusInternalServerError)
	}
	var v1 struct {
		ID string `json:"id"`
		// https://github.com/docker/engine-api/blob/master/types/container/config.go
		Config container.Config `json:"config"`
	}
	if err := json.Unmarshal(out, &v1); err != nil {
		return nil,
			derr.NewErrorWithStatusCode(fmt.Errorf("Failed to unmarshall image history: %s", err),
				http.StatusInternalServerError)
	}
	log.Printf("v1 = %+v", v1)

	imageMetadata := &viccontainer.VicContainer{
		ID:     v1.ID,
		Config: &v1.Config,
	}

	return imageMetadata, nil
}

// attacheStreams takes the the hijacked connections from the calling client and attaches
// them to the 3 streams from the portlayer's rest server.
// name is the container id
// tty indicates whether we should use Docker's copyescapable for stdin
// stdinOnce indicates whether to close the stdin at the container when we finish (we don't handle this but still need to check)
// clStdin, clStdout, clStderr are the hijacked connection
// keys are the keys that are used for detaching streams without closing them
func attachStreams(name string, tty, stdinOnce bool, clStdin io.ReadCloser, clStdout, clStderr io.Writer, keys []byte) error {
	var wg sync.WaitGroup

	errors := make(chan error, 3)
	//FIXME: Swagger will timeout on us.  We need to either have an infinite timeout or the timeout should
	// start after some inactivity?
	inContext, inCancel := context.WithTimeout(context.Background(), attachRequestTimeout)
	outContext, outCancel := context.WithTimeout(context.Background(), attachRequestTimeout)
	errContext, errCancel := context.WithTimeout(context.Background(), attachRequestTimeout)

	if clStdin != nil {
		wg.Add(1)
		// Pipe for stdin so we can interject and watch the input streams for detach keys.
		stdinReader, stdinWriter := io.Pipe()

		defer clStdin.Close()
		defer stdinReader.Close()

		transport := httptransport.New(PortLayerServer(), "/", []string{"http"})
		plClient := client.New(transport, nil)
		transport.Consumers["application/octet-stream"] = httpkit.ByteStreamConsumer()
		transport.Producers["application/octet-stream"] = httpkit.ByteStreamProducer()

		// Swagger wants an io.reader so give it the reader pipe.  Also, the swagger call
		// to set the stdin is synchronous so we need to run in a goroutine
		go func() {
			setStdinParams := interaction.NewContainerSetStdinParamsWithContext(inContext).WithID(name)
			setStdinParams = setStdinParams.WithRawStream(stdinReader)
			plClient.Interaction.ContainerSetStdin(setStdinParams)
		}()

		// Copy the stdin from the CLI and write to a pipe.  We need to do this so we can
		// watch the stdin stream for the detach keys.
		go func() {
			var err error
			defer wg.Done()

			if tty {
				_, err = copyEscapable(stdinWriter, clStdin, keys)
			} else {
				_, err = io.Copy(stdinWriter, clStdin)
			}

			if err != nil {
				log.Errorf(err.Error())
			}
			errors <- err

			if stdinOnce && !tty {
				// Close the stdin connection.  Mimicing Docker's behavior.
				// FIXME: If we close this stdin connection.  The portlayer does
				// not really close stdin.  This is diff from current Docker
				// behavior.  However, we're not sure why Docker even has this
				// behavior where you connect to stdin on the first time only.
				// If we really want to add this behavior, we need to add support
				// in the ssh tether in the portlayer.
				log.Printf("Attach stream has stdinOnce set.  VIC does not yet support this.")
			} else {
				// Shutdown the client's request for stdout, stderr
				errCancel()
				outCancel()
			}
		}()
	}

	if clStdout != nil {
		wg.Add(1)
		transport := httptransport.New(PortLayerServer(), "/", []string{"http"})
		plClient := client.New(transport, nil)
		transport.Consumers["application/octet-stream"] = httpkit.ByteStreamConsumer()
		transport.Producers["application/octet-stream"] = httpkit.ByteStreamProducer()

		// Swagger -> pipewriter.  Synchronous, blocking call
		go func() {
			defer wg.Done()
			getStdoutParams := interaction.NewContainerGetStdoutParamsWithContext(outContext).WithID(name)
			_, err := plClient.Interaction.ContainerGetStdout(getStdoutParams, clStdout)
			if clStdin != nil {
				// Close the client stdin connection (e.g. CLI)
				clStdin.Close()
				inCancel()
			}
			if err != nil {
				if _, ok := err.(*interaction.ContainerGetStdoutNotFound); ok {
					errors <- derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
					return
				}

				// If we get here, most likely something went wrong with the port layer API server
				errors <- derr.NewErrorWithStatusCode(fmt.Errorf("Unknown error from the interaction port layer: %s", err),
					http.StatusInternalServerError)
				return
			}

			errors <- nil
		}()
	}

	if clStderr != nil {
		wg.Add(1)
		transport := httptransport.New(PortLayerServer(), "/", []string{"http"})
		plClient := client.New(transport, nil)
		transport.Consumers["application/octet-stream"] = httpkit.ByteStreamConsumer()
		transport.Producers["application/octet-stream"] = httpkit.ByteStreamProducer()

		// Swagger -> pipewriter.  Synchronous, blocking call
		go func() {
			defer wg.Done()
			getStderrParams := interaction.NewContainerGetStderrParamsWithContext(errContext).WithID(name)
			_, err := plClient.Interaction.ContainerGetStderr(getStderrParams, clStderr)
			if clStdin != nil {
				// Close the client stdin connection (e.g. CLI)
				clStdin.Close()
				inCancel()
			}
			if err != nil {
				if _, ok := err.(*interaction.ContainerGetStderrNotFound); ok {
					errors <- derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
					return
				}

				// If we get here, most likely something went wrong with the port layer API server
				errors <- derr.NewErrorWithStatusCode(fmt.Errorf("Unknown error from the interaction port layer: %s", err),
					http.StatusInternalServerError)
				return
			}

			errors <- nil
		}()
	}

	// Wait for all stream copy to exit
	wg.Wait()
	log.Printf("Attach stream closed")
	defer close(errors)
	for err := range errors {
		if err != nil {
			return err
		}
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
