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
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	derr "github.com/docker/docker/api/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	containertypes "github.com/docker/docker/api/types/container"
	eventtypes "github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	dnetwork "github.com/docker/docker/api/types/network"
	timetypes "github.com/docker/docker/api/types/time"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/reference"
	"github.com/docker/docker/utils"
	gonat "github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	"github.com/docker/libnetwork/iptables"
	"github.com/docker/libnetwork/portallocator"
	"github.com/vishvananda/netlink"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	"github.com/vmware/vic/lib/apiservers/engine/backends/convert"
	"github.com/vmware/vic/lib/apiservers/engine/backends/filter"
	"github.com/vmware/vic/lib/apiservers/engine/backends/portmap"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/interaction"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/tasks"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/retry"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
)

// valid filters as of docker commit 49bf474
var acceptedPsFilterTags = map[string]bool{
	"ancestor":  true,
	"before":    true,
	"exited":    true,
	"id":        true,
	"isolation": true,
	"label":     true,
	"name":      true,
	"status":    true,
	"health":    true,
	"since":     true,
	"volume":    true,
	"network":   true,
	"is-task":   true,
}

// currently not supported by vic
var unSupportedPsFilters = map[string]bool{
	"ancestor":  false,
	"health":    false,
	"isolation": false,
	"is-task":   false,
}

const (
	bridgeIfaceName = "bridge"

	// MemoryAlignMB is the value to which container VM memory must align in order for hotadd to work
	MemoryAlignMB = 128
	// MemoryMinMB - the minimum allowable container memory size
	MemoryMinMB = 512
	// MemoryDefaultMB - the default container VM memory size
	MemoryDefaultMB = 2048
	// MinCPUs - the minimum number of allowable CPUs the container can use
	MinCPUs = 1
	// DefaultCPUs - the default number of container VM CPUs
	DefaultCPUs = 2
	// Default timeout to stop a container if not specified in container config
	DefaultStopTimeout = 10
)

var (
	publicIfaceName = "public"

	defaultScope struct {
		sync.Mutex
		scope string
	}

	portMapper portmap.PortMapper

	// bridge-to-bridge rules, indexed by mapped port;
	// this map is used to delete the rule once
	// the container stops or is removed
	btbRules map[string][]string

	cbpLock         sync.Mutex
	containerByPort map[string]string // port:containerID

	ctx = context.TODO()
)

func init() {
	portMapper = portmap.NewPortMapper()
	btbRules = make(map[string][]string)
	containerByPort = make(map[string]string)

	l, err := netlink.LinkByName(publicIfaceName)
	if l == nil {
		l, err = netlink.LinkByAlias(publicIfaceName)
		if err != nil {
			log.Errorf("interface %s not found", publicIfaceName)
			return
		}
	}

	// don't use interface alias for iptables rules
	publicIfaceName = l.Attrs().Name

	// seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())
}

// type and funcs to provide sorting by created date
type containerByCreated []*types.Container

func (r containerByCreated) Len() int           { return len(r) }
func (r containerByCreated) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r containerByCreated) Less(i, j int) bool { return r[i].Created < r[j].Created }

// Container struct represents the Container
type Container struct {
	containerProxy VicContainerProxy
}

const (
	defaultEnvPath = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
)

func (c *Container) Handle(id, name string) (string, error) {
	resp, err := c.containerProxy.Client().Containers.Get(containers.NewGetParamsWithContext(ctx).WithID(id))
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

// NewContainerBackend returns a new Container
func NewContainerBackend() *Container {
	return &Container{
		containerProxy: NewContainerProxy(PortLayerClient(), PortLayerServer(), PortLayerName()),
	}
}

// docker's container.execBackend

func (c *Container) TaskInspect(cid, cname, eid string) (*models.TaskInspectResponse, error) {
	// obtain a portlayer client
	client := c.containerProxy.Client()

	handle, err := c.Handle(cid, cname)
	if err != nil {
		return nil, err
	}

	// inspect the Task to obtain ProcessConfig
	config := &models.TaskInspectConfig{
		Handle: handle,
		ID:     eid,
	}
	params := tasks.NewInspectParamsWithContext(ctx).WithConfig(config)
	resp, err := client.Tasks.Inspect(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil

}

// ContainerExecCreate sets up an exec in a running container.
func (c *Container) ContainerExecCreate(name string, config *types.ExecConfig) (string, error) {
	defer trace.End(trace.Begin(name))

	if !config.Detach {
		return "", fmt.Errorf("%s only supports detached exec commands at this time", ProductName())
	}

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return "", NotFoundError(name)
	}
	id := vc.ContainerID

	handle, err := c.Handle(id, name)
	if err != nil {
		return "", InternalServerError(err.Error())
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
	log.Debugf("JoinConfig: %#v", joinconfig)

	// obtain a portlayer client
	client := c.containerProxy.Client()

	// call Join with JoinParams
	joinparams := tasks.NewJoinParamsWithContext(ctx).WithConfig(joinconfig)
	resp, err := client.Tasks.Join(joinparams)
	if err != nil {
		return "", InternalServerError(err.Error())
	}
	eid := resp.Payload.ID

	handle = resp.Payload.Handle.(string)
	_, err = client.Containers.Commit(containers.NewCommitParamsWithContext(ctx).WithHandle(handle))
	if err != nil {
		return "", InternalServerError(err.Error())
	}

	// associate newly created exec task with container
	cache.ContainerCache().AddExecToContainer(vc, eid)

	ec, err := c.TaskInspect(id, name, eid)
	if err != nil {
		return "", InternalServerError(err.Error())
	}

	// exec_create event
	event := "exec_create: " + ec.ProcessConfig.ExecPath + " " + strings.Join(ec.ProcessConfig.ExecArgs[1:], " ")
	actor := CreateContainerEventActorWithAttributes(vc, map[string]string{})
	EventService().Log(event, eventtypes.ContainerEventType, actor)

	return eid, nil
}

// ContainerExecInspect returns low-level information about the exec
// command. An error is returned if the exec cannot be found.
func (c *Container) ContainerExecInspect(eid string) (*backend.ExecInspect, error) {
	defer trace.End(trace.Begin(eid))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainerFromExec(eid)
	if vc == nil {
		return nil, NotFoundError(eid)
	}
	id := vc.ContainerID
	name := vc.Name

	ec, err := c.TaskInspect(id, name, eid)
	if err != nil {
		return nil, InternalServerError(err.Error())
	}

	exit := int(ec.ExitCode)
	return &backend.ExecInspect{
		ID:       ec.ID,
		Running:  ec.Running,
		ExitCode: &exit,
		ProcessConfig: &backend.ExecProcessConfig{
			Tty:        ec.Tty,
			Entrypoint: ec.ProcessConfig.ExecPath,
			Arguments:  ec.ProcessConfig.ExecArgs,
			User:       ec.User,
		},
		OpenStdin:   ec.OpenStdin,
		OpenStdout:  ec.OpenStdout,
		OpenStderr:  ec.OpenStderr,
		ContainerID: vc.ContainerID,
		Pid:         int(ec.Pid),
	}, nil
}

// ContainerExecResize changes the size of the TTY of the process
// running in the exec with the given name to the given height and
// width.
func (c *Container) ContainerExecResize(name string, height, width int) error {
	defer trace.End(trace.Begin(name))

	// FIXME(caglar10ur)

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainerFromExec(name)
	if vc == nil {
		return NotFoundError(name)
	}
	/*
		// portlayer client
		client := c.containerProxy.Client()

		resizeParams := containers.NewExecResizeParamsWithContext(ctx).WithHeight(int32(height)).WithWidth(int32(width)).WithID(vc.ContainerID).WithEid(name)
		_, err := client.Containers.ExecResize(resizeParams)
		if err != nil {
			return InternalServerError(err.Error())
		}
	*/
	return nil
}

// ContainerExecStart starts a previously set up exec instance. The
// std streams are set up.
func (c *Container) ContainerExecStart(ctx context.Context, eid string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer) error {
	defer trace.End(trace.Begin(eid))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainerFromExec(eid)
	if vc == nil {
		return NotFoundError(eid)
	}
	id := vc.ContainerID
	name := vc.Name

	handle, err := c.Handle(id, name)
	if err != nil {
		return InternalServerError(err.Error())
	}

	bindconfig := &models.TaskBindConfig{
		Handle: handle,
		ID:     eid,
	}

	// obtain a portlayer client
	client := c.containerProxy.Client()

	// call Bind with bindparams
	bindparams := tasks.NewBindParamsWithContext(ctx).WithConfig(bindconfig)
	resp, err := client.Tasks.Bind(bindparams)
	if err != nil {
		return InternalServerError(err.Error())
	}
	handle = resp.Payload.Handle.(string)

	_, err = client.Containers.Commit(containers.NewCommitParamsWithContext(ctx).WithHandle(handle))
	if err != nil {
		return InternalServerError(err.Error())
	}

	ec, err := c.TaskInspect(id, name, eid)
	if err != nil {
		return InternalServerError(err.Error())
	}

	// exec_start event
	event := "exec_start: " + ec.ProcessConfig.ExecPath + " " + strings.Join(ec.ProcessConfig.ExecArgs[1:], " ")
	actor := CreateContainerEventActorWithAttributes(vc, map[string]string{})
	EventService().Log(event, eventtypes.ContainerEventType, actor)

	return nil
}

// ExecExists looks up the exec instance and returns a bool if it exists or not.
// It will also return the error produced by `getConfig`
func (c *Container) ExecExists(name string) (bool, error) {
	defer trace.End(trace.Begin(name))

	vc := cache.ContainerCache().GetContainerFromExec(name)
	if vc == nil {
		return false, NotFoundError(name)
	}
	return true, nil
}

// docker's container.copyBackend

// ContainerArchivePath creates an archive of the filesystem resource at the
// specified path in the container identified by the given name. Returns a
// tar archive of the resource and whether it was a directory or a single file.
func (c *Container) ContainerArchivePath(name string, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error) {
	return nil, nil, fmt.Errorf("%s does not yet implement ContainerArchivePath", ProductName())
}

// ContainerCopy performs a deprecated operation of archiving the resource at
// the specified path in the container identified by the given name.
func (c *Container) ContainerCopy(name string, res string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("%s does not yet implement ContainerCopy", ProductName())
}

// ContainerExport writes the contents of the container to the given
// writer. An error is returned if the container cannot be found.
func (c *Container) ContainerExport(name string, out io.Writer) error {
	return fmt.Errorf("%s does not yet implement ContainerExport", ProductName())
}

// ContainerExtractToDir extracts the given archive to the specified location
// in the filesystem of the container identified by the given name. The given
// path must be of a directory in the container. If it is not, the error will
// be ErrExtractPointNotDirectory. If noOverwriteDirNonDir is true then it will
// be an error if unpacking the given content would cause an existing directory
// to be replaced with a non-directory and vice versa.
func (c *Container) ContainerExtractToDir(name, path string, noOverwriteDirNonDir bool, content io.Reader) error {
	return fmt.Errorf("%s does not yet implement ContainerExtractToDir", ProductName())
}

// ContainerStatPath stats the filesystem resource at the specified path in the
// container identified by the given name.
func (c *Container) ContainerStatPath(name string, path string) (stat *types.ContainerPathStat, err error) {
	return nil, fmt.Errorf("%s does not yet implement ContainerStatPath", ProductName())
}

// docker's container.stateBackend

// ContainerCreate creates a container.
func (c *Container) ContainerCreate(config types.ContainerCreateConfig) (containertypes.ContainerCreateCreatedBody, error) {
	defer trace.End(trace.Begin(""))

	var err error

	// bail early if container name already exists
	if exists := cache.ContainerCache().GetContainerByName(config.Name); exists != nil {
		err := fmt.Errorf("Conflict. The name %q is already in use by container %s. You have to remove (or rename) that container to be able to re use that name.", config.Name, exists.ContainerID)
		log.Errorf("%s", err.Error())
		return containertypes.ContainerCreateCreatedBody{}, derr.NewRequestConflictError(err)
	}

	// get the image from the cache
	image, err := cache.ImageCache().Get(config.Config.Image)
	if err != nil {
		// if no image found then error thrown and a pull
		// will be initiated by the docker client
		log.Errorf("ContainerCreate: image %s error: %s", config.Config.Image, err.Error())
		return containertypes.ContainerCreateCreatedBody{}, derr.NewRequestNotFoundError(err)
	}

	setCreateConfigOptions(config.Config, image.Config)

	log.Debugf("config.Config = %+v", config.Config)
	if err = validateCreateConfig(&config); err != nil {
		return containertypes.ContainerCreateCreatedBody{}, err
	}

	// Create a container representation in the personality server.  This representation
	// will be stored in the cache if create succeeds in the port layer.
	container, err := createInternalVicContainer(image, &config)
	if err != nil {
		return containertypes.ContainerCreateCreatedBody{}, err
	}

	// Create an actualized container in the VIC port layer
	id, err := c.containerCreate(container, config)
	if err != nil {
		return containertypes.ContainerCreateCreatedBody{}, err
	}

	// Container created ok, save the container id and save the config override from the API
	// caller and save this container internal representation in our personality server's cache
	copyConfigOverrides(container, config)
	container.ContainerID = id
	cache.ContainerCache().AddContainer(container)

	log.Debugf("Container create - name(%s), containerID(%s), config(%#v), host(%#v)",
		container.Name, container.ContainerID, container.Config, container.HostConfig)

	// Add create event
	actor := CreateContainerEventActorWithAttributes(container, map[string]string{})
	EventService().Log(containerCreateEvent, eventtypes.ContainerEventType, actor)

	return containertypes.ContainerCreateCreatedBody{ID: id}, nil
}

// createContainer() makes calls to the container proxy to actually create the backing
// VIC container.  All remoting code is in the proxy.
//
// returns:
//	(container id, error)
func (c *Container) containerCreate(vc *viccontainer.VicContainer, config types.ContainerCreateConfig) (string, error) {
	defer trace.End(trace.Begin("Container.containerCreate"))

	if vc == nil {
		return "", InternalServerError("Failed to create container")
	}

	id, h, err := c.containerProxy.CreateContainerHandle(vc, config)
	if err != nil {
		return "", err
	}

	h, err = c.containerProxy.CreateContainerTask(h, id, config)
	if err != nil {
		return "", err
	}

	h, err = c.containerProxy.AddContainerToScope(h, config)
	if err != nil {
		return id, err
	}

	h, err = c.containerProxy.AddInteractionToContainer(h, config)
	if err != nil {
		return id, err
	}

	h, err = c.containerProxy.AddLoggingToContainer(h, config)
	if err != nil {
		return id, err
	}

	h, err = c.containerProxy.AddVolumesToContainer(h, config)
	if err != nil {
		return id, err
	}

	err = c.containerProxy.CommitContainerHandle(h, id, -1)
	if err != nil {
		return id, err
	}

	return id, nil
}

// ContainerKill sends signal to the container
// If no signal is given (sig 0), then Kill with SIGKILL and wait
// for the container to exit.
// If a signal is given, then just send it to the container and return.
func (c *Container) ContainerKill(name string, sig uint64) error {
	defer trace.End(trace.Begin(fmt.Sprintf("%s, %d", name, sig)))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}

	err := c.containerProxy.Signal(vc, sig)
	if err == nil {
		actor := CreateContainerEventActorWithAttributes(vc, map[string]string{"signal": fmt.Sprintf("%d", sig)})

		EventService().Log(containerKillEvent, eventtypes.ContainerEventType, actor)

	}

	return err
}

// ContainerPause pauses a container
func (c *Container) ContainerPause(name string) error {
	return fmt.Errorf("%s does not yet implement ContainerPause", ProductName())
}

// ContainerResize changes the size of the TTY of the process running
// in the container with the given name to the given height and width.
func (c *Container) ContainerResize(name string, height, width int) error {
	defer trace.End(trace.Begin(name))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}

	// Call the port layer to resize
	plHeight := int32(height)
	plWidth := int32(width)

	var err error
	if err = c.containerProxy.Resize(vc, plHeight, plWidth); err == nil {
		actor := CreateContainerEventActorWithAttributes(vc, map[string]string{
			"height": fmt.Sprintf("%d", height),
			"width":  fmt.Sprintf("%d", width),
		})

		EventService().Log(containerResizeEvent, eventtypes.ContainerEventType, actor)
	}

	return err
}

// ContainerRestart stops and starts a container. It attempts to
// gracefully stop the container within the given timeout, forcefully
// stopping it if the timeout is exceeded. If given a negative
// timeout, ContainerRestart will wait forever until a graceful
// stop. Returns an error if the container cannot be found, or if
// there is an underlying error at any stage of the restart.
func (c *Container) ContainerRestart(name string, seconds *int) error {
	defer trace.End(trace.Begin(name))

	// Look up the container name in the metadata cache ot get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}

	operation := func() error {
		return c.containerProxy.Stop(vc, name, seconds, false)
	}
	if err := retry.Do(operation, IsConflictError); err != nil {
		return InternalServerError(fmt.Sprintf("Stop failed with: %s", err))
	}

	operation = func() error {
		return c.containerStart(name, nil, false)
	}
	if err := retry.Do(operation, IsConflictError); err != nil {
		return InternalServerError(fmt.Sprintf("Start failed with: %s", err))
	}

	actor := CreateContainerEventActorWithAttributes(vc, map[string]string{})
	EventService().Log(containerRestartEvent, eventtypes.ContainerEventType, actor)

	return nil
}

// ContainerRm removes the container id from the filesystem. An error
// is returned if the container is not found, or if the remove
// fails. If the remove succeeds, the container name is released, and
// network links are removed.
func (c *Container) ContainerRm(name string, config *types.ContainerRmConfig) error {
	defer trace.End(trace.Begin(name))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}
	id := vc.ContainerID

	// Get the portlayer Client API
	client := c.containerProxy.Client()

	// TODO: Pass this RemoveVolume flag to somewhere
	_ = &config.RemoveVolume

	// Use the force and stop the container first
	secs := 0

	if config.ForceRemove {
		c.containerProxy.Stop(vc, name, &secs, true)
	} else {
		state, err := c.containerProxy.State(vc)
		if err != nil {
			if IsNotFoundError(err) {
				cache.ContainerCache().DeleteContainer(id)
				return NotFoundError(name)
			}
			return InternalServerError(err.Error())
		}
		// force stop if container state is error to make sure container is deletable later
		if state.Status == ContainerError {
			c.containerProxy.Stop(vc, name, &secs, true)
		}
	}

	//call the remove directly on the name. No need for using a handle.
	_, err := client.Containers.ContainerRemove(containers.NewContainerRemoveParamsWithContext(ctx).WithID(id))
	if err != nil {
		switch err := err.(type) {
		case *containers.ContainerRemoveNotFound:
			cache.ContainerCache().DeleteContainer(id)
			return NotFoundError(name)
		case *containers.ContainerRemoveDefault:
			return InternalServerError(err.Payload.Message)
		case *containers.ContainerRemoveConflict:
			return derr.NewRequestConflictError(fmt.Errorf("You cannot remove a running container. Stop the container before attempting removal or use -f"))
		default:
			return InternalServerError(err.Error())
		}
	}

	return nil
}

// cleanupPortBindings gets port bindings for the container and
// unmaps ports if the cVM that previously bound them isn't powered on
func (c *Container) cleanupPortBindings(vc *viccontainer.VicContainer) error {
	defer trace.End(trace.Begin(vc.ContainerID))
	for ctrPort, hostPorts := range vc.HostConfig.PortBindings {
		for _, hostPort := range hostPorts {
			hPort := hostPort.HostPort

			cbpLock.Lock()
			mappedCtr, mapped := containerByPort[hPort]
			cbpLock.Unlock()
			if !mapped {
				continue
			}

			log.Debugf("Container %q maps host port %s to container port %s", mappedCtr, hPort, ctrPort)
			// check state of the previously bound container with PL
			cc := cache.ContainerCache().GetContainer(mappedCtr)
			if cc == nil {
				// The container was removed from the cache and
				// port bindings were cleaned up by another operation.
				continue
			}
			state, err := c.containerProxy.State(cc)
			if err != nil {
				if IsNotFoundError(err) {
					log.Debugf("container(%s) not found in portLayer, removing from persona cache", cc.ContainerID)
					// we have a container in the persona cache, but it's been removed from the portLayer
					// which is the source of truth -- so remove from the persona cache after this func
					// completes
					defer cache.ContainerCache().DeleteContainer(cc.ContainerID)
				} else {
					// we have issues of an unknown variety...return..
					return InternalServerError(err.Error())
				}
			}

			if state != nil && state.Running {
				log.Debugf("Running container %q still holds port %s", mappedCtr, hPort)
				continue
			}

			log.Debugf("Unmapping ports for powered off / removed container %q", mappedCtr)
			err = UnmapPorts(cc.HostConfig)
			if err != nil {
				return fmt.Errorf("Failed to unmap host port %s for container %q: %s",
					hPort, mappedCtr, err)
			}
		}
	}
	return nil
}

// ContainerStart starts a container.
func (c *Container) ContainerStart(name string, hostConfig *containertypes.HostConfig, checkpoint string, checkpointDir string) error {
	defer trace.End(trace.Begin(name))

	operation := func() error {
		return c.containerStart(name, hostConfig, true)
	}
	if err := retry.Do(operation, IsConflictError); err != nil {
		return err
	}
	return nil
}

func (c *Container) containerStart(name string, hostConfig *containertypes.HostConfig, bind bool) error {
	var err error

	// Get an API client to the portlayer
	client := c.containerProxy.Client()

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}
	id := vc.ContainerID

	// handle legacy hostConfig
	if hostConfig != nil {
		// hostConfig exist for backwards compatibility.  TODO: Figure out which parameters we
		// need to look at in hostConfig
	} else if vc != nil {
		hostConfig = vc.HostConfig
	}

	if vc != nil && hostConfig.NetworkMode.NetworkName() == "" {
		hostConfig.NetworkMode = vc.HostConfig.NetworkMode
	}

	// get a handle to the container
	handle, err := c.Handle(id, name)
	if err != nil {
		return err
	}

	var endpoints []*models.EndpointConfig
	// bind network
	if bind {
		var bindRes *scopes.BindContainerOK
		bindRes, err = client.Scopes.BindContainer(scopes.NewBindContainerParamsWithContext(ctx).WithHandle(handle))
		if err != nil {
			switch err := err.(type) {
			case *scopes.BindContainerNotFound:
				cache.ContainerCache().DeleteContainer(id)
				return NotFoundError(name)
			case *scopes.BindContainerInternalServerError:
				return InternalServerError(err.Payload.Message)
			default:
				return InternalServerError(err.Error())
			}
		}

		handle = bindRes.Payload.Handle
		endpoints = bindRes.Payload.Endpoints

		// unbind in case we fail later
		defer func() {
			if err != nil {
				client.Scopes.UnbindContainer(scopes.NewUnbindContainerParamsWithContext(ctx).WithHandle(handle))
			}
		}()

		// unmap ports that vc needs if they're not being used by previously mapped container
		err = c.cleanupPortBindings(vc)
		if err != nil {
			return err
		}
	}

	// change the state of the container
	// TODO: We need a resolved ID from the name
	var stateChangeRes *containers.StateChangeOK
	stateChangeRes, err = client.Containers.StateChange(containers.NewStateChangeParamsWithContext(ctx).WithHandle(handle).WithState("RUNNING"))
	if err != nil {
		switch err := err.(type) {
		case *containers.StateChangeNotFound:
			cache.ContainerCache().DeleteContainer(id)
			return NotFoundError(name)
		case *containers.StateChangeDefault:
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	handle = stateChangeRes.Payload

	// map ports
	if bind {
		e := c.findPortBoundNetworkEndpoint(hostConfig, endpoints)
		if err = MapPorts(hostConfig, e, id); err != nil {
			return InternalServerError(fmt.Sprintf("error mapping ports: %s", err))
		}

		defer func() {
			if err != nil {
				UnmapPorts(hostConfig)
			}
		}()
	}

	// commit the handle; this will reconfigure and start the vm
	_, err = client.Containers.Commit(containers.NewCommitParamsWithContext(ctx).WithHandle(handle))
	if err != nil {
		switch err := err.(type) {
		case *containers.CommitNotFound:
			cache.ContainerCache().DeleteContainer(id)
			return NotFoundError(name)
		case *containers.CommitConflict:
			return ConflictError(err.Error())
		case *containers.CommitDefault:
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	actor := CreateContainerEventActorWithAttributes(vc, map[string]string{})
	EventService().Log(containerStartEvent, eventtypes.ContainerEventType, actor)
	return nil
}

// requestHostPort finds a free port on the host
func requestHostPort(proto string) (int, error) {
	pa := portallocator.Get()
	return pa.RequestPortInRange(nil, proto, 0, 0)
}

type portMapping struct {
	intHostPort int
	strHostPort string
	portProto   gonat.Port
}

// unrollPortMap processes config for mapping/unmapping ports e.g. from hostconfig.PortBindings
func unrollPortMap(portMap gonat.PortMap) ([]*portMapping, error) {
	var portMaps []*portMapping
	for i, pb := range portMap {

		proto, port := gonat.SplitProtoPort(string(i))
		nport, err := gonat.NewPort(proto, port)
		if err != nil {
			return nil, err
		}

		// iterate over all the ports in pb []nat.PortBinding
		for i := range pb {
			var hostPort int
			var hPort string
			if pb[i].HostPort == "" {
				// use a random port since no host port is specified
				hostPort, err = requestHostPort(proto)
				if err != nil {
					log.Errorf("could not find available port on host")
					return nil, err
				}
				log.Infof("using port %d on the host for port mapping", hostPort)

				// update the hostconfig
				pb[i].HostPort = strconv.Itoa(hostPort)

			} else {
				hostPort, err = strconv.Atoi(pb[i].HostPort)
				if err != nil {
					return nil, err
				}
			}
			hPort = strconv.Itoa(hostPort)
			portMaps = append(portMaps, &portMapping{
				intHostPort: hostPort,
				strHostPort: hPort,
				portProto:   nport,
			})
		}
	}
	return portMaps, nil
}

// MapPorts maps ports defined in hostconfig for containerID
func MapPorts(hostconfig *containertypes.HostConfig, endpoint *models.EndpointConfig, containerID string) error {
	log.Debugf("mapPorts for %q: %v", containerID, hostconfig.PortBindings)

	if len(hostconfig.PortBindings) == 0 {
		return nil
	}
	if endpoint == nil {
		return fmt.Errorf("invalid endpoint")
	}

	var containerIP net.IP
	containerIP = net.ParseIP(endpoint.Address)
	if containerIP == nil {
		return fmt.Errorf("invalid endpoint address %s", endpoint.Address)
	}

	portMap, err := unrollPortMap(hostconfig.PortBindings)
	if err != nil {
		return err
	}

	cbpLock.Lock()
	defer cbpLock.Unlock()
	for _, p := range portMap {
		if err = portMapper.MapPort(nil, p.intHostPort, p.portProto.Proto(), containerIP.String(), p.portProto.Int(), publicIfaceName, bridgeIfaceName); err != nil {
			return err
		}

		// bridge-to-bridge pin hole for traffic from containers for exposed port
		if err = interBridgeTraffic(portmap.Map, p.strHostPort, p.portProto.Proto(), containerIP.String(), p.portProto.Port()); err != nil {
			return err
		}

		// update mapped ports
		containerByPort[p.strHostPort] = containerID
		log.Debugf("mapped port %s for container %s", p.strHostPort, containerID)
	}
	return nil
}

// UnmapPorts unmaps ports defined in hostconfig
func UnmapPorts(hostconfig *containertypes.HostConfig) error {
	log.Debugf("UnmapPorts: %v", hostconfig.PortBindings)

	if len(hostconfig.PortBindings) == 0 {
		return nil
	}

	portMap, err := unrollPortMap(hostconfig.PortBindings)
	if err != nil {
		return err
	}

	cbpLock.Lock()
	defer cbpLock.Unlock()
	for _, p := range portMap {
		// check if we should actually unmap based on current mappings
		_, mapped := containerByPort[p.strHostPort]
		if !mapped {
			log.Debugf("skipping already unmapped %s", p.strHostPort)
			continue
		}

		if err = portMapper.UnmapPort(nil, p.intHostPort, p.portProto.Proto(), p.portProto.Int(), publicIfaceName, bridgeIfaceName); err != nil {
			return err
		}

		// bridge-to-bridge pin hole for traffic from containers for exposed port
		if err = interBridgeTraffic(portmap.Unmap, p.strHostPort, "", "", ""); err != nil {
			return err
		}

		// update mapped ports
		delete(containerByPort, p.strHostPort)
		log.Debugf("unmapped port %s", p.strHostPort)
	}
	return nil
}

// interBridgeTraffic enables traffic for exposed port from one bridge network to another
func interBridgeTraffic(op portmap.Operation, hostPort, proto, containerAddr, containerPort string) error {
	switch op {
	case portmap.Map:
		switch proto {
		case "udp", "tcp":
		default:
			return fmt.Errorf("unknown protocol: %s", proto)
		}

		// rule to allow connections from bridge interface for the
		// specific mapped port. has to inserted at the top of the
		// chain rather than appended to supersede bridge-to-bridge
		// traffic blocking
		baseArgs := []string{"-t", string(iptables.Filter),
			"-i", bridgeIfaceName,
			"-o", bridgeIfaceName,
			"-p", proto,
			"-d", containerAddr,
			"--dport", containerPort,
			"-j", "ACCEPT",
		}

		args := append([]string{string(iptables.Insert), "VIC", "1"}, baseArgs...)
		if _, err := iptables.Raw(args...); err != nil && !os.IsExist(err) {
			return err
		}

		btbRules[hostPort] = baseArgs
	case portmap.Unmap:
		if args, ok := btbRules[hostPort]; ok {
			args = append([]string{string(iptables.Delete), "VIC"}, args...)
			if _, err := iptables.Raw(args...); err != nil && !os.IsNotExist(err) {
				return err
			}

			delete(btbRules, hostPort)
		}
	}

	return nil
}

func (c *Container) defaultScope() string {
	defaultScope.Lock()
	defer defaultScope.Unlock()

	if defaultScope.scope != "" {
		return defaultScope.scope
	}

	client := c.containerProxy.Client()
	listRes, err := client.Scopes.List(scopes.NewListParamsWithContext(ctx).WithIDName("default"))
	if err != nil {
		log.Error(err)
		return ""
	}

	if len(listRes.Payload) != 1 {
		log.Errorf("could not get default scope name")
		return ""
	}

	defaultScope.scope = listRes.Payload[0].Name
	return defaultScope.scope
}

func (c *Container) findPortBoundNetworkEndpoint(hostconfig *containertypes.HostConfig, endpoints []*models.EndpointConfig) *models.EndpointConfig {
	if len(hostconfig.PortBindings) == 0 {
		return nil
	}

	// check if the port binding network is a bridge type
	listRes, err := PortLayerClient().Scopes.List(scopes.NewListParamsWithContext(ctx).WithIDName(hostconfig.NetworkMode.NetworkName()))
	if err != nil {
		log.Error(err)
		return nil
	}

	if len(listRes.Payload) != 1 || listRes.Payload[0].ScopeType != "bridge" {
		log.Warnf("port binding for network %s is not bridge type", hostconfig.NetworkMode.NetworkName())
		return nil
	}

	// look through endpoints to find the container's IP on the network that has the port binding
	for _, e := range endpoints {
		if hostconfig.NetworkMode.NetworkName() == e.Scope || (hostconfig.NetworkMode.IsDefault() && e.Scope == c.defaultScope()) {
			return e
		}
	}

	return nil
}

// ContainerStop looks for the given container and terminates it,
// waiting the given number of seconds before forcefully killing the
// container. If a negative number of seconds is given, ContainerStop
// will wait for a graceful termination. An error is returned if the
// container is not found, is already stopped, or if there is a
// problem stopping the container.
func (c *Container) ContainerStop(name string, seconds *int) error {
	defer trace.End(trace.Begin(name))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}

	if seconds == nil {
		timeout := DefaultStopTimeout
		if vc.Config.StopTimeout != nil {
			timeout = *vc.Config.StopTimeout
		}
		seconds = &timeout
	}

	operation := func() error {
		return c.containerProxy.Stop(vc, name, seconds, true)
	}
	if err := retry.Do(operation, IsConflictError); err != nil {
		return err
	}

	actor := CreateContainerEventActorWithAttributes(vc, map[string]string{})
	EventService().Log(containerStopEvent, eventtypes.ContainerEventType, actor)

	return nil
}

// ContainerUnpause unpauses a container
func (c *Container) ContainerUnpause(name string) error {
	return fmt.Errorf("%s does not yet implement ContainerUnpause", ProductName())
}

// ContainerUpdate updates configuration of the container
func (c *Container) ContainerUpdate(name string, hostConfig *containertypes.HostConfig) (containertypes.ContainerUpdateOKBody, error) {
	return containertypes.ContainerUpdateOKBody{}, fmt.Errorf("%s does not yet implement ontainerUpdate", ProductName())
}

// ContainerWait stops processing until the given container is
// stopped. If the container is not found, an error is returned. On a
// successful stop, the exit code of the container is returned. On a
// timeout, an error is returned. If you want to wait forever, supply
// a negative duration for the timeout.
func (c *Container) ContainerWait(name string, timeout time.Duration) (int, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("name(%s):timeout(%s)", name, timeout)))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return -1, NotFoundError(name)
	}

	processExitCode, processStatus, containerState, err := c.containerProxy.Wait(vc, timeout)
	if err != nil {
		return -1, err
	}

	// call to the dockerStatus function to retrieve the docker friendly exitCode
	exitCode, _ := dockerStatus(int(processExitCode), processStatus, containerState, time.Time{}, time.Time{})

	return exitCode, nil
}

// dockerStatus will evaluate the container state, exit code and
// process status to return a docker friendly status
//
// exitCode is the container process exit code
// status is the container process status -- stored in the vmx file as "started"
// started & finished are the process start / finish times
func dockerStatus(exitCode int, status string, state string, started time.Time, finished time.Time) (int, string) {

	// set docker status to state and we'll change if needed
	dockStatus := state

	switch state {
	case "Running":
		// if we don't have a start date leave the status as the state
		if !started.IsZero() {
			dockStatus = fmt.Sprintf("Up %s", units.HumanDuration(time.Now().UTC().Sub(started)))
		}
	case "Stopped":
		// if we don't have a finished date then don't process exitCode and return "Stopped" for the status
		if !finished.IsZero() {
			// interrogate the process status returned from the portlayer
			// and based on status text and exit codes set the appropriate
			// docker exit code
			if strings.Contains(status, "permission denied") {
				exitCode = 126
			} else if strings.Contains(status, "no such") {
				exitCode = 127
			} else if status == "true" && exitCode == -1 {
				// most likely the process was killed via the cli
				// or received a sigkill
				exitCode = 137
			} else if status == "" && exitCode == 0 {
				// the process was stopped via the cli
				// or received a sigterm
				exitCode = 143
			}

			dockStatus = fmt.Sprintf("Exited (%d) %s ago", exitCode, units.HumanDuration(time.Now().UTC().Sub(finished)))
		}
	}

	return exitCode, dockStatus
}

// docker's container.monitorBackend

// ContainerChanges returns a list of container fs changes
func (c *Container) ContainerChanges(name string) ([]archive.Change, error) {
	return make([]archive.Change, 0, 0), fmt.Errorf("%s does not yet implement ontainerChanges", ProductName())
}

// ContainerInspect returns low-level information about a
// container. Returns an error if the container cannot be found, or if
// there is an error getting the data.
func (c *Container) ContainerInspect(name string, size bool, version string) (interface{}, error) {
	// Ignore version.  We're supporting post-1.20 version.
	defer trace.End(trace.Begin(name))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return nil, NotFoundError(name)
	}
	id := vc.ContainerID
	log.Debugf("Found %q in cache as %q", id, vc.ContainerID)

	client := c.containerProxy.Client()

	results, err := client.Containers.GetContainerInfo(containers.NewGetContainerInfoParamsWithContext(ctx).WithID(id))
	if err != nil {
		switch err := err.(type) {
		case *containers.GetContainerInfoNotFound:
			cache.ContainerCache().DeleteContainer(id)
			return nil, NotFoundError(name)
		case *containers.GetContainerInfoInternalServerError:
			return nil, InternalServerError(err.Payload.Message)
		default:
			return nil, InternalServerError(err.Error())
		}
	}
	var started time.Time
	var stopped time.Time

	if results.Payload.ProcessConfig.StartTime > 0 {
		started = time.Unix(results.Payload.ProcessConfig.StartTime, 0)
	}
	if results.Payload.ProcessConfig.StopTime > 0 {
		stopped = time.Unix(results.Payload.ProcessConfig.StopTime, 0)
	}

	// call to the dockerStatus function to retrieve the docker friendly exitCode
	exitCode, status := dockerStatus(
		int(results.Payload.ProcessConfig.ExitCode),
		results.Payload.ProcessConfig.Status,
		results.Payload.ContainerConfig.State,
		started, stopped)

	// set the payload values
	exit := int32(exitCode)
	results.Payload.ProcessConfig.ExitCode = exit
	results.Payload.ProcessConfig.Status = status

	inspectJSON, err := ContainerInfoToDockerContainerInspect(vc, results.Payload, PortLayerName())
	if err != nil {
		log.Errorf("containerInfoToDockerContainerInspect failed with %s", err)
		return nil, err
	}

	log.Debugf("ContainerInspect json config = %+v\n", inspectJSON.Config)

	return inspectJSON, nil
}

// ContainerLogs hooks up a container's stdout and stderr streams
// configured with the given struct.
func (c *Container) ContainerLogs(ctx context.Context, name string, config *backend.ContainerLogsConfig, started chan struct{}) error {
	defer trace.End(trace.Begin(""))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}
	name = vc.ContainerID

	tailLines, since, err := c.validateContainerLogsConfig(vc, config)
	if err != nil {
		return err
	}

	// Outstream modification (from Docker's code) so the stream is streamed with the
	// necessary headers that the CLI expects.  This is Docker's scheme.
	wf := ioutils.NewWriteFlusher(config.OutStream)
	defer wf.Close()

	wf.Flush()

	outStream := io.Writer(wf)
	if !vc.Config.Tty {
		outStream = stdcopy.NewStdWriter(outStream, stdcopy.Stdout)
	}

	// Make a call to our proxy to handle the remoting
	err = c.containerProxy.StreamContainerLogs(name, outStream, started, config.Timestamps, config.Follow, since, tailLines)
	if err != nil {
		// Don't return an error encountered while streaming logs.
		// Once we've started streaming logs, the Docker client doesn't expect
		// an error to be returned as it leads to a malformed response body.
		log.Errorf("error while streaming logs: %#v", err)
	}

	return nil
}

// ContainerStats writes information about the container to the stream
// given in the config object.
func (c *Container) ContainerStats(ctx context.Context, name string, config *backend.ContainerStatsConfig) error {
	defer trace.End(trace.Begin(name))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}

	// get the configured CPUMhz for this VCH so that we can calculate docker CPU stats
	cpuMhz, err := systemBackend.SystemCPUMhzLimit()
	if err != nil {
		// wrap error to provide a bit more detail
		sysErr := fmt.Errorf("unable to gather system CPUMhz for container(%s): %s", vc.ContainerID, err)
		log.Error(sysErr)
		return InternalServerError(sysErr.Error())
	}

	out := config.OutStream
	if config.Stream {
		// Outstream modification (from Docker's code) so the stream is streamed with the
		// necessary headers that the CLI expects.  This is Docker's scheme.
		wf := ioutils.NewWriteFlusher(config.OutStream)
		defer wf.Close()
		wf.Flush()
		out = io.Writer(wf)
	}

	// stats configuration
	statsConfig := &convert.ContainerStatsConfig{
		VchMhz:      cpuMhz,
		Stream:      config.Stream,
		ContainerID: vc.ContainerID,
		Out:         out,
		Memory:      vc.HostConfig.Memory,
	}

	// if we are not streaming then we need to get the container state
	if !config.Stream {
		statsConfig.ContainerState, err = c.containerProxy.State(vc)
		if err != nil {
			return InternalServerError(err.Error())
		}

	}

	err = c.containerProxy.StreamContainerStats(ctx, statsConfig)
	if err != nil {
		log.Errorf("error while streaming container (%s) stats: %s", vc.ContainerID, err)
	}
	return nil
}

// ContainerTop lists the processes running inside of the given
// container by calling ps with the given args, or with the flags
// "-ef" if no args are given.  An error is returned if the container
// is not found, or is not running, or if there are any problems
// running ps, or parsing the output.
func (c *Container) ContainerTop(name string, psArgs string) (*types.ContainerProcessList, error) {
	return nil, fmt.Errorf("%s does not yet implement ContainerTop", ProductName())
}

// Containers returns the list of containers to show given the user's filtering.
func (c *Container) Containers(config *types.ContainerListOptions) ([]*types.Container, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("ListOptions %#v", config)))

	// validate filters for support and validity
	listContext, err := filter.ValidateContainerFilters(config, acceptedPsFilterTags, unSupportedPsFilters)
	if err != nil {
		return nil, err
	}

	// Get an API client to the portlayer
	client := c.containerProxy.Client()

	containme, err := client.Containers.GetContainerList(containers.NewGetContainerListParamsWithContext(ctx).WithAll(&listContext.All))
	if err != nil {
		switch err := err.(type) {

		case *containers.GetContainerListInternalServerError:
			return nil, fmt.Errorf("Error invoking GetContainerList: %s", err.Payload.Message)

		default:
			return nil, fmt.Errorf("Error invoking GetContainerList: %s", err.Error())
		}
	}
	// TODO: move to conversion function
	containers := make([]*types.Container, 0, len(containme.Payload))

payloadLoop:
	for _, t := range containme.Payload {
		var started time.Time
		var stopped time.Time

		if t.ProcessConfig.StartTime > 0 {
			started = time.Unix(t.ProcessConfig.StartTime, 0)
		}

		if t.ProcessConfig.StopTime > 0 {
			stopped = time.Unix(t.ProcessConfig.StopTime, 0)
		}

		// get the docker friendly status
		exitCode, status := dockerStatus(
			int(t.ProcessConfig.ExitCode),
			t.ProcessConfig.Status,
			t.ContainerConfig.State,
			started,
			stopped)

		// labels func requires a config be passed
		// TODO: refactor labelsFromAnnotations func for broader use
		tempConfig := &containertypes.Config{}
		err = labelsFromAnnotations(tempConfig, t.ContainerConfig.Annotations)
		if err != nil && config.Filters.Include("label") {
			return nil, fmt.Errorf("unable to convert vic annotations to docker labels (%s)", t.ContainerConfig.ContainerID)
		}

		listContext.Labels = tempConfig.Labels
		listContext.ExitCode = exitCode
		listContext.ID = t.ContainerConfig.ContainerID

		// prior to further conversion lets determine if this container
		// is needed or if the list is complete -- if the container is
		// needed conversion will continue and the container will be added to the
		// return array
		action := filter.IncludeContainer(listContext, t)
		switch action {
		case filter.ExcludeAction:
			// skip to next container
			continue payloadLoop
		case filter.StopAction:
			// we're done
			break payloadLoop
		}

		cmd := strings.Join(t.ProcessConfig.ExecArgs, " ")
		// the docker client expects the friendly name to be prefixed
		// with a forward slash -- create a new slice and add here
		names := make([]string, 0, len(t.ContainerConfig.Names))
		for i := range t.ContainerConfig.Names {
			names = append(names, clientFriendlyContainerName(t.ContainerConfig.Names[i]))
		}

		ips, err := publicIPv4Addrs()
		var ports []types.Port
		if err != nil {
			log.Errorf("Could not get IP information for reporting port bindings.")
		} else {
			ports = portInformation(t, ips)
		}

		// verify that the repo:tag exists for the container -- if it doesn't then we should present the
		// truncated imageID -- if we have a failure determining then we'll show the data we have
		repo := *t.ContainerConfig.RepoName
		ref, _ := reference.ParseNamed(*t.ContainerConfig.RepoName)
		if ref != nil {
			imageID, err := cache.RepositoryCache().Get(ref)
			if err != nil && err == cache.ErrDoesNotExist {
				// the tag has been removed, so we need to show the truncated imageID
				imageID = cache.RepositoryCache().GetImageID(t.ContainerConfig.LayerID)
				if imageID != "" {
					id := uid.Parse(imageID)
					repo = id.Truncate().String()
				}
			}
		}

		c := &types.Container{
			ID:      t.ContainerConfig.ContainerID,
			Image:   repo,
			Created: t.ContainerConfig.CreateTime,
			Status:  status,
			Names:   names,
			Command: cmd,
			SizeRw:  t.ContainerConfig.StorageSize,
			Ports:   ports,
			State:   filter.DockerState(t.ContainerConfig.State),
		}

		// The container should be included in the list
		containers = append(containers, c)
		listContext.Counter++

	}

	return containers, nil
}

func (c *Container) ContainersPrune(pruneFilters filters.Args) (*types.ContainersPruneReport, error) {
	return nil, fmt.Errorf("%s does not yet implement ContainersPrune", ProductName())
}

// docker's container.attachBackend

// ContainerAttach attaches to logs according to the config passed in. See ContainerAttachConfig.
func (c *Container) ContainerAttach(name string, ca *backend.ContainerAttachConfig) error {
	defer trace.End(trace.Begin(name))

	operation := func() error {
		return c.containerAttach(name, ca)
	}
	if err := retry.Do(operation, IsConflictError); err != nil {
		return err
	}
	return nil
}

func (c *Container) containerAttach(name string, ca *backend.ContainerAttachConfig) error {
	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)

	}
	id := vc.ContainerID

	client := c.containerProxy.Client()
	handle, err := c.Handle(id, name)
	if err != nil {
		return err
	}

	bind, err := client.Interaction.InteractionBind(interaction.NewInteractionBindParamsWithContext(ctx).
		WithConfig(&models.InteractionBindConfig{
			Handle: handle,
		}))
	if err != nil {
		return InternalServerError(err.Error())
	}
	handle, ok := bind.Payload.Handle.(string)
	if !ok {
		return InternalServerError(fmt.Sprintf("Type assertion failed for %#+v", handle))
	}

	// commit the handle; this will reconfigure the vm
	_, err = client.Containers.Commit(containers.NewCommitParamsWithContext(ctx).WithHandle(handle))
	if err != nil {
		switch err := err.(type) {
		case *containers.CommitNotFound:
			return NotFoundError(name)
		case *containers.CommitConflict:
			return ConflictError(err.Error())
		case *containers.CommitDefault:
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	clStdin, clStdout, clStderr, err := ca.GetStreams()
	if err != nil {
		return InternalServerError("Unable to get stdio streams for calling client")
	}
	defer clStdin.Close()

	if !vc.Config.Tty && ca.MuxStreams {
		// replace the stdout/stderr with Docker's multiplex stream
		if ca.UseStdout {
			clStderr = stdcopy.NewStdWriter(clStderr, stdcopy.Stderr)
		}
		if ca.UseStderr {
			clStdout = stdcopy.NewStdWriter(clStdout, stdcopy.Stdout)
		}
	}

	actor := CreateContainerEventActorWithAttributes(vc, map[string]string{})
	EventService().Log(containerAttachEvent, eventtypes.ContainerEventType, actor)
	err = c.containerProxy.AttachStreams(context.Background(), vc, clStdin, clStdout, clStderr, ca)
	if err != nil {
		if _, ok := err.(DetachError); ok {
			// fire detach event
			actor := CreateContainerEventActorWithAttributes(vc, map[string]string{})
			EventService().Log(containerDetachEvent, eventtypes.ContainerEventType, actor)

			log.Infof("Detach detected, tearing down connection")
			client = c.containerProxy.Client()
			handle, err = c.Handle(id, name)
			if err != nil {
				return err
			}

			unbind, err := client.Interaction.InteractionUnbind(interaction.NewInteractionUnbindParamsWithContext(ctx).
				WithConfig(&models.InteractionUnbindConfig{
					Handle: handle,
				}))
			if err != nil {
				return InternalServerError(err.Error())
			}

			handle, ok = unbind.Payload.Handle.(string)
			if !ok {
				return InternalServerError("type assertion failed")
			}

			// commit the handle; this will reconfigure the vm
			_, err = client.Containers.Commit(containers.NewCommitParamsWithContext(ctx).WithHandle(handle))
			if err != nil {
				switch err := err.(type) {
				case *containers.CommitNotFound:
					return NotFoundError(name)
				case *containers.CommitConflict:
					return ConflictError(err.Error())
				case *containers.CommitDefault:
					return InternalServerError(err.Payload.Message)
				default:
					return InternalServerError(err.Error())
				}
			}
		}
		return err
	}

	return nil
}

// ContainerRename changes the name of a container, using the oldName
// to find the container. An error is returned if newName is already
// reserved.
func (c *Container) ContainerRename(oldName, newName string) error {
	defer trace.End(trace.Begin(newName))

	if oldName == "" || newName == "" {
		err := fmt.Errorf("neither old nor new names may be empty")
		log.Errorf("%s", err.Error())
		return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
	}

	if !utils.RestrictedNamePattern.MatchString(newName) {
		err := fmt.Errorf("invalid container name (%s), only %s are allowed", newName, utils.RestrictedNameChars)
		log.Errorf("%s", err.Error())
		return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
	}

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(oldName)
	if vc == nil {
		log.Errorf("Container %s not found", oldName)
		return NotFoundError(oldName)
	}

	oldName = vc.Name
	if oldName == newName {
		err := fmt.Errorf("renaming a container with the same name as its current name")
		log.Errorf("%s", err.Error())
		return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
	}

	// reserve the new name in containerCache
	if err := cache.ContainerCache().ReserveName(vc, newName); err != nil {
		log.Errorf("%s", err.Error())
		return derr.NewRequestConflictError(err)
	}

	if err := c.containerProxy.Rename(vc, newName); err != nil {
		log.Errorf("Rename error: %s", err)
		cache.ContainerCache().ReleaseName(newName)
		return err
	}

	// update containerCache
	if err := cache.ContainerCache().UpdateContainerName(oldName, newName); err != nil {
		log.Errorf("Failed to update container cache: %s", err)
		cache.ContainerCache().ReleaseName(newName)
		return err
	}

	log.Infof("Container %s renamed to %s", oldName, newName)

	actor := CreateContainerEventActorWithAttributes(vc, map[string]string{"newName": fmt.Sprintf("%s", newName)})

	EventService().Log("Rename", eventtypes.ContainerEventType, actor)

	return nil
}

// helper function to format the container name
// to the docker client approved format
func clientFriendlyContainerName(name string) string {
	return fmt.Sprintf("/%s", name)
}

//------------------------------------
// ContainerCreate() Utility Functions
//------------------------------------

// createInternalVicContainer() creates an container representation (for docker personality)
// This is called by ContainerCreate()
func createInternalVicContainer(image *metadata.ImageConfig, config *types.ContainerCreateConfig) (*viccontainer.VicContainer, error) {
	// provide basic container config via the image
	container := viccontainer.NewVicContainer()
	container.LayerID = image.V1Image.ID // store childmost layer ID to map to the proper vmdk
	container.ImageID = image.ImageID
	container.Config = image.Config //Set defaults.  Overrides will get copied below.

	return container, nil
}

// SetConfigOptions is a place to add necessary container configuration
// values that were not explicitly supplied by the user
func setCreateConfigOptions(config, imageConfig *containertypes.Config) {
	// Overwrite or append the image's config from the CLI with the metadata from the image's
	// layer metadata where appropriate
	if len(config.Cmd) == 0 {
		config.Cmd = imageConfig.Cmd
	}
	if config.WorkingDir == "" {
		config.WorkingDir = imageConfig.WorkingDir
	}
	if len(config.Entrypoint) == 0 {
		config.Entrypoint = imageConfig.Entrypoint
	}

	if config.Volumes == nil {
		config.Volumes = imageConfig.Volumes
	} else {
		for k, v := range imageConfig.Volumes {
			//NOTE: the value of the map is an empty struct.
			//      we also do not care about duplicates.
			//      This Volumes map is really a Set.
			config.Volumes[k] = v
		}
	}

	if config.User == "" {
		config.User = imageConfig.User
	}
	// set up environment
	setEnvFromImageConfig(config, imageConfig)
}

func setEnvFromImageConfig(config, imageConfig *containertypes.Config) {
	// Set PATH in ENV if needed
	setPathFromImageConfig(config, imageConfig)

	containerEnv := make(map[string]string, len(config.Env))
	for _, env := range config.Env {
		kv := strings.SplitN(env, "=", 2)
		var val string
		if len(kv) == 2 {
			val = kv[1]
		}
		containerEnv[kv[0]] = val
	}

	// Set TERM to xterm if tty is set, unless user supplied a different TERM
	if config.Tty {
		if _, ok := containerEnv["TERM"]; !ok {
			config.Env = append(config.Env, "TERM=xterm")
		}
	}

	// add remaining environment variables from the image config to the container
	// config, taking care not to overwrite anything
	for _, imageEnv := range imageConfig.Env {
		key := strings.SplitN(imageEnv, "=", 2)[0]
		// is environment variable already set in container config?
		if _, ok := containerEnv[key]; !ok {
			// no? let's copy it from the image config
			config.Env = append(config.Env, imageEnv)
		}
	}
}

func setPathFromImageConfig(config, imageConfig *containertypes.Config) {
	// check if user supplied PATH environment variable at creation time
	for _, v := range config.Env {
		if strings.HasPrefix(v, "PATH=") {
			// a PATH is set, bail
			return
		}
	}

	// check to see if the image this container is created from supplies a PATH
	for _, v := range imageConfig.Env {
		if strings.HasPrefix(v, "PATH=") {
			// a PATH was found, add it to the config
			config.Env = append(config.Env, v)
			return
		}
	}

	// no PATH set, use the default
	config.Env = append(config.Env, fmt.Sprintf("PATH=%s", defaultEnvPath))
}

// validateCreateConfig() checks the parameters for ContainerCreate().
// It may "fix up" the config param passed into ConntainerCreate() if needed.
func validateCreateConfig(config *types.ContainerCreateConfig) error {
	defer trace.End(trace.Begin("Container.validateCreateConfig"))

	if config.Config == nil {
		return BadRequestError("invalid config")
	}

	if config.HostConfig == nil {
		config.HostConfig = &containertypes.HostConfig{}
	}

	// process cpucount here
	var cpuCount int64 = DefaultCPUs

	// support windows client
	if config.HostConfig.CPUCount > 0 {
		cpuCount = config.HostConfig.CPUCount
	} else {
		// we hijack --cpuset-cpus in the non-windows case
		if config.HostConfig.CpusetCpus != "" {
			cpus := strings.Split(config.HostConfig.CpusetCpus, ",")
			if c, err := strconv.Atoi(cpus[0]); err == nil {
				cpuCount = int64(c)
			} else {
				return fmt.Errorf("Error parsing CPU count: %s", err)
			}
		}
	}
	config.HostConfig.CPUCount = cpuCount

	// fix-up cpu/memory settings here
	if cpuCount < MinCPUs {
		config.HostConfig.CPUCount = MinCPUs
	}
	log.Infof("Container CPU count: %d", config.HostConfig.CPUCount)

	// convert from bytes to MiB for vsphere
	memoryMB := config.HostConfig.Memory / units.MiB
	if memoryMB == 0 {
		memoryMB = MemoryDefaultMB
	} else if memoryMB < MemoryMinMB {
		memoryMB = MemoryMinMB
	}

	// check that memory is aligned
	if remainder := memoryMB % MemoryAlignMB; remainder != 0 {
		log.Warnf("Default container VM memory must be %d aligned for hotadd, rounding up.", MemoryAlignMB)
		memoryMB += MemoryAlignMB - remainder
	}

	config.HostConfig.Memory = memoryMB
	log.Infof("Container memory: %d MB", config.HostConfig.Memory)

	if config.NetworkingConfig == nil {
		config.NetworkingConfig = &dnetwork.NetworkingConfig{}
	} else {
		if l := len(config.NetworkingConfig.EndpointsConfig); l > 1 {
			return fmt.Errorf("NetworkMode error: Container can be connected to one network endpoint only")
		}
		// If NetworkConfig exists, set NetworkMode to the default endpoint network, assuming only one endpoint network as the default network during container create
		for networkName := range config.NetworkingConfig.EndpointsConfig {
			config.HostConfig.NetworkMode = containertypes.NetworkMode(networkName)
		}
	}

	// validate port bindings
	var ips []string
	if addrs, err := publicIPv4Addrs(); err != nil {
		log.Warnf("could not get address for public interface: %s", err)
	} else {
		ips = make([]string, len(addrs))
		for i := range addrs {
			ips[i] = addrs[i].IP.String()
		}
	}

	for _, pbs := range config.HostConfig.PortBindings {
		for _, pb := range pbs {
			if pb.HostIP != "" && pb.HostIP != "0.0.0.0" {
				// check if specified host ip equals any of the addresses on the "client" interface
				found := false
				for _, i := range ips {
					if i == pb.HostIP {
						found = true
						break
					}
				}
				if !found {
					return InternalServerError("host IP for port bindings is only supported for 0.0.0.0 and the public interface IP address")
				}
			}

			start, end, _ := gonat.ParsePortRangeToInt(pb.HostPort)
			if start != end {
				return InternalServerError("host port ranges are not supported for port bindings")
			}
		}
	}

	// https://github.com/vmware/vic/issues/1378
	if len(config.Config.Entrypoint) == 0 && len(config.Config.Cmd) == 0 {
		return derr.NewRequestNotFoundError(fmt.Errorf("No command specified"))
	}

	// Was a name provided - if not create a friendly name
	if config.Name == "" {
		//TODO: Assume we could have a name collison here : need to
		// provide validation / retry CDG June 9th 2016
		config.Name = namesgenerator.GetRandomName(0)
	}

	return nil
}

func copyConfigOverrides(vc *viccontainer.VicContainer, config types.ContainerCreateConfig) {
	// Copy the create overrides to our new container
	vc.Name = config.Name
	vc.Config.Cmd = config.Config.Cmd
	vc.Config.WorkingDir = config.Config.WorkingDir
	vc.Config.Entrypoint = config.Config.Entrypoint
	vc.Config.Env = config.Config.Env
	vc.Config.AttachStdin = config.Config.AttachStdin
	vc.Config.AttachStdout = config.Config.AttachStdout
	vc.Config.AttachStderr = config.Config.AttachStderr
	vc.Config.Tty = config.Config.Tty
	vc.Config.OpenStdin = config.Config.OpenStdin
	vc.Config.StdinOnce = config.Config.StdinOnce
	vc.Config.StopSignal = config.Config.StopSignal
	vc.Config.Volumes = config.Config.Volumes
	vc.HostConfig = config.HostConfig
}

func publicIPv4Addrs() ([]netlink.Addr, error) {
	l, err := netlink.LinkByName(publicIfaceName)
	if err != nil {
		return nil, fmt.Errorf("Could not look up link from public interface name %s due to error %s",
			publicIfaceName, err.Error())
	}
	ips, err := netlink.AddrList(l, netlink.FAMILY_V4)
	if err != nil {
		return nil, fmt.Errorf("Could not get IP addresses of link due to error %s", err.Error())
	}

	return ips, nil
}

// returns port bindings as a slice of Docker Ports for return to the client
// returns empty slice on error
func portInformation(t *models.ContainerInfo, ips []netlink.Addr) []types.Port {
	// create a port for each IP on the interface (usually only 1, but could be more)
	// (works with both IPv4 and IPv6 addresses)
	var ports []types.Port

	cid := t.ContainerConfig.ContainerID
	c := cache.ContainerCache().GetContainer(cid)

	if c == nil {
		log.Errorf("Could not find container with ID %s", cid)
		return ports
	}

	for _, ip := range ips {
		ports = append(ports, types.Port{IP: ip.IP.String()})
	}

	portBindings := c.HostConfig.PortBindings
	var resultPorts []types.Port

	for _, port := range ports {
		for portBindingPrivatePort, hostPortBindings := range portBindings {
			portAndType := strings.SplitN(string(portBindingPrivatePort), "/", 2)
			portNum, err := strconv.Atoi(portAndType[0])
			if err != nil {
				log.Infof("Got an error trying to convert private port number to an int")
				continue
			}
			port.PrivatePort = uint16(portNum)
			port.Type = portAndType[1]

			for i := 0; i < len(hostPortBindings); i++ {
				newport := port
				publicPort, err := strconv.Atoi(hostPortBindings[i].HostPort)
				if err != nil {
					log.Infof("Got an error trying to convert public port number to an int")
					continue
				}
				newport.PublicPort = uint16(publicPort)
				// sanity check -- sometimes these come back as 0 when no binding actually exists
				// that doesn't make sense, so in that case we don't want to report these bindings
				if newport.PublicPort != 0 && newport.PrivatePort != 0 {
					resultPorts = append(resultPorts, newport)
				}
			}
		}
	}
	return resultPorts
}

//----------------------------------
// ContainerLogs() utility functions
//----------------------------------

// validateContainerLogsConfig() validates and extracts options for logging from the
// backend.ContainerLogsConfig object we're given.
//
// returns:
//	tail lines, since (in unix time), error
func (c *Container) validateContainerLogsConfig(vc *viccontainer.VicContainer, config *backend.ContainerLogsConfig) (int64, int64, error) {
	if !(config.ShowStdout || config.ShowStderr) {
		return 0, 0, fmt.Errorf("You must choose at least one stream")
	}

	unsupported := func(opt string) (int64, int64, error) {
		return 0, 0, fmt.Errorf("container %s does not support '--%s'", vc.ContainerID, opt)
	}

	tailLines := int64(-1)
	if config.Tail != "" && config.Tail != "all" {
		n, err := strconv.ParseInt(config.Tail, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("error parsing tail option: %s", err)
		}
		tailLines = n
	}

	var since time.Time
	if config.Since != "" {
		s, n, err := timetypes.ParseTimestamps(config.Since, 0)
		if err != nil {
			return 0, 0, err
		}
		since = time.Unix(s, n)
	}

	// TODO(jzt): this should not require an extra call to the portlayer. We should
	// update container.DataVersion when we hydrate the container cache at VCH startup
	// see https://github.com/vmware/vic/issues/4194
	if config.Timestamps {
		// check container DataVersion to make sure it's supported
		params := containers.NewGetContainerInfoParams()
		params.SetID(vc.ContainerID)
		info, err := PortLayerClient().Containers.GetContainerInfo(params)
		if err != nil {
			return 0, 0, err
		}
		if info.Payload.DataVersion == 0 {
			return unsupported("timestamps")
		}
	}

	if config.Since != "" {
		return unsupported("since")
	}

	return tailLines, since.Unix(), nil
}

func CreateContainerEventActorWithAttributes(vc *viccontainer.VicContainer, attributes map[string]string) eventtypes.Actor {
	if vc.Config != nil {
		for k, v := range vc.Config.Labels {
			attributes[k] = v
		}
	}
	if vc.Config.Image != "" {
		attributes["image"] = vc.Config.Image
	}
	attributes["name"] = strings.TrimLeft(vc.Name, "/")

	return eventtypes.Actor{
		ID:         vc.ContainerID,
		Attributes: attributes,
	}
}
