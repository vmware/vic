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
// container.go
//
// Rules for code to be in here:
// 1. No remote or swagger calls.  Move those code to container_portlayer.go
// 2. Always return docker engine-api compatible errors.
//		- Do NOT return fmt.Errorf()
//		- Do NOT return errors.New()
//		- DO USE the aliased docker error package 'derr'
//		- It is OK to return errors returned from functions in container_portlayer.go

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/backend"
	derr "github.com/docker/docker/errors"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/version"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
)

// Container struct represents the Container
type Container struct {
}

// docker's container.execBackend

// ContainerExecCreate sets up an exec in a running container.
func (c *Container) ContainerExecCreate(config *types.ExecConfig) (string, error) {
	return "", fmt.Errorf("%s does not implement container.ContainerExecCreate", ProductName())
}

// ContainerExecInspect returns low-level information about the exec
// command. An error is returned if the exec cannot be found.
func (c *Container) ContainerExecInspect(id string) (*backend.ExecInspect, error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerExecInspect", ProductName())
}

// ContainerExecResize changes the size of the TTY of the process
// running in the exec with the given name to the given height and
// width.
func (c *Container) ContainerExecResize(name string, height, width int) error {
	return fmt.Errorf("%s does not implement container.ContainerExecResize", ProductName())
}

// ContainerExecStart starts a previously set up exec instance. The
// std streams are set up.
func (c *Container) ContainerExecStart(name string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer) error {
	return fmt.Errorf("%s does not implement container.ContainerExecStart", ProductName())
}

// ExecExists looks up the exec instance and returns a bool if it exists or not.
// It will also return the error produced by `getConfig`
func (c *Container) ExecExists(name string) (bool, error) {
	return false, fmt.Errorf("%s does not implement container.ExecExists", ProductName())
}

// docker's container.copyBackend

// ContainerArchivePath creates an archive of the filesystem resource at the
// specified path in the container identified by the given name. Returns a
// tar archive of the resource and whether it was a directory or a single file.
func (c *Container) ContainerArchivePath(name string, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error) {
	return nil, nil, fmt.Errorf("%s does not implement container.ContainerArchivePath", ProductName())
}

// ContainerCopy performs a deprecated operation of archiving the resource at
// the specified path in the container identified by the given name.
func (c *Container) ContainerCopy(name string, res string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerCopy", ProductName())
}

// ContainerExport writes the contents of the container to the given
// writer. An error is returned if the container cannot be found.
func (c *Container) ContainerExport(name string, out io.Writer) error {
	return fmt.Errorf("%s does not implement container.ContainerExport", ProductName())
}

// ContainerExtractToDir extracts the given archive to the specified location
// in the filesystem of the container identified by the given name. The given
// path must be of a directory in the container. If it is not, the error will
// be ErrExtractPointNotDirectory. If noOverwriteDirNonDir is true then it will
// be an error if unpacking the given content would cause an existing directory
// to be replaced with a non-directory and vice versa.
func (c *Container) ContainerExtractToDir(name, path string, noOverwriteDirNonDir bool, content io.Reader) error {
	return fmt.Errorf("%s does not implement container.ContainerExtractToDir", ProductName())
}

// ContainerStatPath stats the filesystem resource at the specified path in the
// container identified by the given name.
func (c *Container) ContainerStatPath(name string, path string) (stat *types.ContainerPathStat, err error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerStatPath", ProductName())
}

// docker's container.stateBackend

// ContainerCreate creates a container.
func (c *Container) ContainerCreate(config types.ContainerCreateConfig) (types.ContainerCreateResponse, error) {
	defer trace.End(trace.Begin("ContainerCreate"))

	log.Debugf("config.Config = %+v", config.Config)

	// bail early if container name already exists
	if exists := cache.ContainerCache().GetContainer(config.Name); exists != nil {
		err := fmt.Errorf("Conflict. The name \"%s\" is already in use by container %s. You have to remove (or rename) that container to be able to re use that name.", config.Name, exists.ContainerID)
		log.Errorf("%s", err.Error())
		return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(err, http.StatusConflict)
	}

	// get the image from the personality server's cache
	image, err := cache.ImageCache().GetImage(config.Config.Image)
	if err != nil {
		// if no image found then error thrown and a pull will be initiated by the docker client
		log.Errorf("ContainerCreate: image %s error: %s", config.Config.Image, err.Error())
		return types.ContainerCreateResponse{}, err
	}

	// Create a container representation in the personality server.  This representation
	// will be stored in the cache if create succeeds in the port layer.
	container, err := createInternalVicContainer(image, &config)
	if err != nil {
		return types.ContainerCreateResponse{}, err
	}

	// Create a container in the port layer
	id, h, err := VicCreateContainer(container, config)
	if err != nil {
		return types.ContainerCreateResponse{}, err
	}

	// Container created ok, save the container id and save the container internal
	// represenation in our personality server's cache
	container.ContainerID = id
	cache.ContainerCache().AddContainer(container)

	log.Debugf("Container create: %#v", container)

	// Success!
	log.Debugf("container.ContainerCreate succeeded.  Returning container handle %s", h)
	return types.ContainerCreateResponse{ID: id}, nil
}

// ContainerKill sends signal to the container
// If no signal is given (sig 0), then Kill with SIGKILL and wait
// for the container to exit.
// If a signal is given, then just send it to the container and return.
func (c *Container) ContainerKill(name string, sig uint64) error {
	return fmt.Errorf("%s does not implement container.ContainerKill", ProductName())
}

// ContainerPause pauses a container
func (c *Container) ContainerPause(name string) error {
	return fmt.Errorf("%s does not implement container.ContainerPause", ProductName())
}

// ContainerRename changes the name of a container, using the oldName
// to find the container. An error is returned if newName is already
// reserved.
func (c *Container) ContainerRename(oldName, newName string) error {
	return fmt.Errorf("%s does not implement container.ContainerRename", ProductName())
}

// ContainerResize changes the size of the TTY of the process running
// in the container with the given name to the given height and width.
func (c *Container) ContainerResize(name string, height, width int) error {
	defer trace.End(trace.Begin("ContainerResize"))

	// Look up the container name in the metadata cache to get long ID
	if vc := cache.ContainerCache().GetContainer(name); vc != nil {
		log.Debugf("Found %q in cache as %q", name, vc.ContainerID)
		name = vc.ContainerID
	}

	resp := VicResizeContainer(name, height, width)

	return resp
}

// ContainerRestart stops and starts a container. It attempts to
// gracefully stop the container within the given timeout, forcefully
// stopping it if the timeout is exceeded. If given a negative
// timeout, ContainerRestart will wait forever until a graceful
// stop. Returns an error if the container cannot be found, or if
// there is an underlying error at any stage of the restart.
func (c *Container) ContainerRestart(name string, seconds int) error {
	defer trace.End(trace.Begin("ContainerRestart"))

	err := VicContainerStop(name, seconds, false)
	if err != nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("Stop failed with: %s", err), http.StatusInternalServerError)
	}

	err = VicContainerStart(name, nil, false)
	if err != nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("Start failed with: %s", err), http.StatusInternalServerError)
	}

	return nil
}

// ContainerRm removes the container id from the filesystem. An error
// is returned if the container is not found, or if the remove
// fails. If the remove succeeds, the container name is released, and
// network links are removed.
func (c *Container) ContainerRm(name string, config *types.ContainerRmConfig) error {
	defer trace.End(trace.Begin("ContainerRm"))

	// Look up the container name in the metadata cache to get long ID
	if vc := cache.ContainerCache().GetContainer(name); vc != nil {
		log.Debugf("Found %q in cache as %q", name, vc.ContainerID)
		name = vc.ContainerID
	}

	// Call the port layer to remove the container
	err := VicContainerRemove(name)
	if err != nil {
		return err
	}

	// delete container from the cache
	cache.ContainerCache().DeleteContainer(name)

	return nil
}

// ContainerStart starts a container.
func (c *Container) ContainerStart(name string, hostConfig *container.HostConfig) error {
	defer trace.End(trace.Begin("ContainerStart"))

	// Look up the container name in the metadata cache to get long ID
	if vc := cache.ContainerCache().GetContainer(name); vc != nil {
		log.Debugf("Found %q in cache as %q", name, vc.ContainerID)
		name = vc.ContainerID
	}

	return VicContainerStart(name, hostConfig, true)
}

// ContainerStop looks for the given container and terminates it,
// waiting the given number of seconds before forcefully killing the
// container. If a negative number of seconds is given, ContainerStop
// will wait for a graceful termination. An error is returned if the
// container is not found, is already stopped, or if there is a
// problem stopping the container.
func (c *Container) ContainerStop(name string, seconds int) error {
	defer trace.End(trace.Begin("ContainerStop"))

	// Look up the container name in the metadata cache to get long ID
	if vc := cache.ContainerCache().GetContainer(name); vc != nil {
		log.Debugf("Found %q in cache as %q", name, vc.ContainerID)
		name = vc.ContainerID
	}

	return VicContainerStop(name, seconds, true)
}

// ContainerUnpause unpauses a container
func (c *Container) ContainerUnpause(name string) error {
	return fmt.Errorf("%s does not implement container.ContainerUnpause", ProductName())
}

// ContainerUpdate updates configuration of the container
func (c *Container) ContainerUpdate(name string, hostConfig *container.HostConfig) ([]string, error) {
	return make([]string, 0, 0), fmt.Errorf("%s does not implement container.ContainerUpdate", ProductName())
}

// ContainerWait stops processing until the given container is
// stopped. If the container is not found, an error is returned. On a
// successful stop, the exit code of the container is returned. On a
// timeout, an error is returned. If you want to wait forever, supply
// a negative duration for the timeout.
func (c *Container) ContainerWait(name string, timeout time.Duration) (int, error) {
	return 0, fmt.Errorf("%s does not implement container.ContainerWait", ProductName())
}

// docker's container.monitorBackend

// ContainerChanges returns a list of container fs changes
func (c *Container) ContainerChanges(name string) ([]archive.Change, error) {
	return make([]archive.Change, 0, 0), fmt.Errorf("%s does not implement container.ContainerChanges", ProductName())
}

// ContainerInspect returns low-level information about a
// container. Returns an error if the container cannot be found, or if
// there is an error getting the data.
func (c *Container) ContainerInspect(name string, size bool, version version.Version) (interface{}, error) {
	// Ignore version.  We're supporting post-1.20 version.
	defer trace.End(trace.Begin("ContainerInspect"))

	// Look up the container name in the metadata cache to get long ID
	if vc := cache.ContainerCache().GetContainer(name); vc != nil {
		log.Debugf("Found %q in cache as %q", name, vc.ContainerID)
		name = vc.ContainerID
	}

	plInspectData, err := VicContainerInspect(name)
	if err != nil {
		return nil, err
	}

	inspectJSON, err := plContainerInfoToDockerContainerInspect(name, plInspectData)
	if err != nil {
		log.Errorf("containerInfoToDockerContainerInspect failed with %s", err)
		return nil, err
	}

	log.Debugf("ContainerInspect json config = %+v\n", inspectJSON.Config)

	return inspectJSON, nil
}

// ContainerLogs hooks up a container's stdout and stderr streams
// configured with the given struct.
func (c *Container) ContainerLogs(name string, config *backend.ContainerLogsConfig, started chan struct{}) error {
	return fmt.Errorf("%s does not implement container.ContainerLogs", ProductName())
}

// ContainerStats writes information about the container to the stream
// given in the config object.
func (c *Container) ContainerStats(name string, config *backend.ContainerStatsConfig) error {
	return fmt.Errorf("%s does not implement container.ContainerStats", ProductName())
}

// ContainerTop lists the processes running inside of the given
// container by calling ps with the given args, or with the flags
// "-ef" if no args are given.  An error is returned if the container
// is not found, or is not running, or if there are any problems
// running ps, or parsing the output.
func (c *Container) ContainerTop(name string, psArgs string) (*types.ContainerProcessList, error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerTop", ProductName())
}

// Containers returns the list of containers to show given the user's filtering.
func (c *Container) Containers(config *types.ContainerListOptions) ([]*types.Container, error) {
	plContainerListData, err := VicContainerList(config.All)
	if err != nil {
		return nil, err
	}

	containers := plContainerListToDockerContainerList(plContainerListData)

	return containers, nil
}

// docker's container.attachBackend

// ContainerAttach attaches to logs according to the config passed in. See ContainerAttachConfig.
func (c *Container) ContainerAttach(name string, ca *backend.ContainerAttachConfig) error {
	defer trace.End(trace.Begin("ContainerAttach"))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
	}

	clStdin, clStdout, clStderr, err := ca.GetStreams()

	if err != nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("Unable to get stdio streams for calling client"), http.StatusInternalServerError)
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

	err = VicAttachStreams(context.Background(), vc, clStdin, clStdout, clStderr, ca)

	return err
}

//----------
// Utility Functions
//----------

func createInternalVicContainer(image *metadata.ImageConfig, config *types.ContainerCreateConfig) (*viccontainer.VicContainer, error) {
	// provide basic container config via the image
	container := viccontainer.NewVicContainer()
	container.ID = image.ID
	container.Config = image.Config

	// Overwrite or append the image's config from the CLI with the metadata from the image's
	// layer metadata where appropriate
	if len(config.Config.Cmd) == 0 {
		config.Config.Cmd = container.Config.Cmd
	}
	if config.Config.WorkingDir == "" {
		config.Config.WorkingDir = container.Config.WorkingDir
	}
	if len(config.Config.Entrypoint) == 0 {
		config.Config.Entrypoint = container.Config.Entrypoint
	}
	// Set TERM to xterm if tty is set
	if config.Config.Tty {
		config.Config.Env = append(config.Config.Env, "TERM=xterm")
	}
	config.Config.Env = append(config.Config.Env, container.Config.Env...)

	// TODO(jzt): users other than root are not currently supported
	// We should check for USER in config.Config.Env once we support Dockerfiles.
	if config.Config.User != "" && config.Config.User != "root" {
		return nil, derr.NewErrorWithStatusCode(fmt.Errorf("Failed to create container - users other than root are not currently supported"),
			http.StatusInternalServerError)
	}

	// Was a name provided - if not create a friendly name
	if config.Name == "" {
		//TODO: Assume we could have a name collison here : need to
		// provide validation / retry CDG June 9th 2016
		config.Name = namesgenerator.GetRandomName(0)
	}
	log.Debugf("ContainerCreate config' = %+v", config)

	// https://github.com/vmware/vic/issues/1378
	if len(config.Config.Entrypoint) == 0 && len(config.Config.Cmd) == 0 {
		return nil, derr.NewRequestNotFoundError(fmt.Errorf("No command specified"))
	}

	// Copy the create overrides
	container.Name = config.Name
	container.Config.Cmd = config.Config.Cmd
	container.Config.WorkingDir = config.Config.WorkingDir
	container.Config.Entrypoint = config.Config.Entrypoint
	container.Config.Env = config.Config.Env
	container.Config.AttachStdin = config.Config.AttachStdin
	container.Config.AttachStdout = config.Config.AttachStdout
	container.Config.AttachStderr = config.Config.AttachStderr
	container.Config.Tty = config.Config.Tty
	container.Config.OpenStdin = config.Config.OpenStdin
	container.Config.StdinOnce = config.Config.StdinOnce

	return container, nil
}
