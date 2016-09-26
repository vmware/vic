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
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/go-swagger/go-swagger/httpkit"
	httptransport "github.com/go-swagger/go-swagger/httpkit/client"
	strfmt "github.com/go-swagger/go-swagger/strfmt"
	"github.com/mreiferson/go-httpclient"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/backend"
	derr "github.com/docker/docker/errors"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/version"
	"github.com/docker/engine-api/types"
	containertypes "github.com/docker/engine-api/types/container"
	dnetwork "github.com/docker/engine-api/types/network"
	timetypes "github.com/docker/engine-api/types/time"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	"github.com/docker/libnetwork/portallocator"

	"github.com/vishvananda/netlink"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	"github.com/vmware/vic/lib/apiservers/engine/backends/portmap"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/interaction"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
)

const (
	bridgeIfaceName = "bridge"
)

var (
	clientIfaceName = "client"

	defaultScope struct {
		sync.Mutex
		scope string
	}

	portMapper portmap.PortMapper

	ctx = context.TODO()
)

func init() {
	portMapper = portmap.NewPortMapper()

	l, err := netlink.LinkByName(clientIfaceName)
	if l == nil {
		l, err = netlink.LinkByAlias(clientIfaceName)
		if err != nil {
			log.Errorf("interface %s not found", clientIfaceName)
			return
		}
	}

	// don't use interface alias for iptables rules
	clientIfaceName = l.Attrs().Name
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
	defer trace.End(trace.Begin(""))

	var err error

	// bail early if container name already exists
	if exists := cache.ContainerCache().GetContainer(config.Name); exists != nil {
		err := fmt.Errorf("Conflict. The name %q is already in use by container %s. You have to remove (or rename) that container to be able to re use that name.", config.Name, exists.ContainerID)
		log.Errorf("%s", err.Error())
		return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(err, http.StatusConflict)
	}

	// get the image from the cache
	image, err := cache.ImageCache().GetImage(config.Config.Image)
	if err != nil {
		// if no image found then error thrown and a pull
		// will be initiated by the docker client
		log.Errorf("ContainerCreate: image %s error: %s", config.Config.Image, err.Error())
		return types.ContainerCreateResponse{}, derr.NewRequestNotFoundError(err)
	}

	setCreateConfigOptions(config.Config, image.Config)

	log.Debugf("config.Config = %+v", config.Config)
	if err = validateCreateConfig(&config); err != nil {
		return types.ContainerCreateResponse{}, err
	}

	// Create a container representation in the personality server.  This representation
	// will be stored in the cache if create succeeds in the port layer.
	container, err := createInternalVicContainer(image, &config)
	if err != nil {
		return types.ContainerCreateResponse{}, err
	}

	// Create an actualized container in the VIC port layer
	id, err := c.containerCreate(container, config)
	if err != nil {
		return types.ContainerCreateResponse{}, err
	}

	// Container created ok, save the container id and save the config override from the API
	// caller and save this container internal representation in our personality server's cache
	copyConfigOverrides(container, config)
	container.ContainerID = id
	cache.ContainerCache().AddContainer(container)

	log.Debugf("Container create - name(%s), containerID(%s), config(%#v), host(%#v)",
		container.Name, container.ContainerID, container.Config, container.HostConfig)

	return types.ContainerCreateResponse{ID: id}, nil
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

	imageID := vc.ImageID

	id, h, err := c.containerProxy.CreateContainerHandle(imageID, config)
	if err != nil {
		return "", err
	}

	h, err = c.containerProxy.AddContainerToScope(h, config)
	if err != nil {
		return id, err
	}

	h, err = c.containerProxy.AddVolumesToContainer(h, config)
	if err != nil {
		return id, err
	}

	err = c.containerProxy.CommitContainerHandle(h, imageID)
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
	defer trace.End(trace.Begin(name))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}
	name = vc.ContainerID

	if err := ContainerSignal(name, sig); err != nil {
		return err
	}

	if sig == 0 {
		// Use ContainerWait infrastructure to wait for container to exit
	}

	return nil
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
	defer trace.End(trace.Begin(name))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)

	}
	id := vc.ContainerID

	// Get an API client to the portlayer
	client := c.containerProxy.Client()

	// Call the port layer to resize
	plHeight := int32(height)
	plWidth := int32(width)
	plResizeParam := interaction.NewContainerResizeParamsWithContext(ctx).WithID(id).WithHeight(plHeight).WithWidth(plWidth)

	_, err := client.Interaction.ContainerResize(plResizeParam)
	if err != nil {
		if _, isa := err.(*interaction.ContainerResizeNotFound); isa {
			cache.ContainerCache().DeleteContainer(id)
			return NotFoundError(name)
		}

		// If we get here, most likely something went wrong with the port layer API server
		return InternalServerError(err.Error())
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
	defer trace.End(trace.Begin(name))

	err := c.containerStop(name, seconds, false)
	if err != nil {
		return InternalServerError(fmt.Sprintf("Stop failed with: %s", err))
	}

	err = c.containerStart(name, nil, false)
	if err != nil {
		return InternalServerError(fmt.Sprintf("Start failed with: %s", err))
	}

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
	if config.ForceRemove {
		c.containerStop(id, 0, true)
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
			return derr.NewErrorWithStatusCode(fmt.Errorf("You cannot remove a running container. Stop the container before attempting removal or use -f"), http.StatusConflict)
		default:
			return InternalServerError(err.Error())
		}
	}
	// delete container from the cache
	cache.ContainerCache().DeleteContainer(id)
	return nil
}

// ContainerStart starts a container.
func (c *Container) ContainerStart(name string, hostConfig *containertypes.HostConfig) error {
	defer trace.End(trace.Begin(name))
	return c.containerStart(name, hostConfig, true)
}

func (c *Container) containerStart(name string, hostConfig *containertypes.HostConfig, bind bool) error {
	var err error

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}
	id := vc.ContainerID

	// Get an API client to the portlayer
	client := c.containerProxy.Client()

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

	}

	// change the state of the container
	// TODO: We need a resolved ID from the name
	var stateChangeRes *containers.StateChangeOK
	stateChangeRes, err = client.Containers.StateChange(containers.NewStateChangeParamsWithContext(ctx).WithHandle(handle).WithState("RUNNING"))
	if err != nil {
		cache.ContainerCache().DeleteContainer(id)
		switch err := err.(type) {
		case *containers.StateChangeNotFound:
			return NotFoundError(name)
		case *containers.StateChangeDefault:
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	handle = stateChangeRes.Payload

	if bind {
		e := c.findPortBoundNetworkEndpoint(hostConfig, endpoints)
		if err = c.mapPorts(portmap.Map, hostConfig, e); err != nil {
			return InternalServerError(fmt.Sprintf("error mapping ports: %s", err))
		}

		defer func() {
			if err != nil {
				c.mapPorts(portmap.Unmap, hostConfig, e)
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
		case *containers.CommitDefault:
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	return nil
}

// requestHostPort finds a free port on the host
func requestHostPort(proto string) (int, error) {
	pa := portallocator.Get()
	return pa.RequestPortInRange(nil, proto, 0, 0)
}

func (c *Container) mapPorts(op portmap.Operation, hostconfig *containertypes.HostConfig, endpoint *models.EndpointConfig) error {
	if len(hostconfig.PortBindings) == 0 || endpoint == nil {
		return nil
	}

	var containerIP net.IP
	containerIP = net.ParseIP(endpoint.Address)
	if containerIP == nil {
		return fmt.Errorf("invalid endpoint address %s", endpoint.Address)
	}

	for _, p := range endpoint.Ports {
		proto, port := nat.SplitProtoPort(p)
		var nport nat.Port
		nport, err := nat.NewPort(proto, port)
		if err != nil {
			return err
		}

		pbs, ok := hostconfig.PortBindings[nport]
		if !ok {
			continue
		}

		for i := range pbs {
			var hostPort int
			if pbs[i].HostPort == "" {
				// use a random port since no host port is specified
				hostPort, err = requestHostPort(proto)
				if err != nil {
					log.Errorf("could not find available port on host")
					return err
				}
				// update the hostconfig
				pbs[i].HostPort = strconv.Itoa(hostPort)

			} else {
				hostPort, err = strconv.Atoi(pbs[i].HostPort)
				if err != nil {
					return err
				}
			}

			if err = portMapper.MapPort(op, nil, hostPort, nport.Proto(), containerIP.String(), nport.Int(), clientIfaceName, bridgeIfaceName); err != nil {
				return err
			}
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
func (c *Container) ContainerStop(name string, seconds int) error {
	defer trace.End(trace.Begin(name))
	return c.containerStop(name, seconds, true)
}

func (c *Container) containerStop(name string, seconds int, unbound bool) error {

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}
	id := vc.ContainerID

	//retrieve client to portlayer
	client := c.containerProxy.Client()
	handle, err := c.Handle(id, name)
	if err != nil {
		return err
	}

	// we have a container on the PL side lets check the state before proceeding
	// ignore the error  since others will be checking below..this is an attempt to short circuit the op
	// TODO: can be replaced with simple cache check once power events are propagated to persona
	infoResponse, err := client.Containers.GetContainerInfo(containers.NewGetContainerInfoParamsWithContext(ctx).WithID(id))
	if err != nil {
		cache.ContainerCache().DeleteContainer(id)
		return NotFoundError(name)
	}
	if *infoResponse.Payload.ContainerConfig.State == "Stopped" || *infoResponse.Payload.ContainerConfig.State == "Created" {
		return nil
	}

	if unbound {
		var endpoints []*models.EndpointConfig
		ub, err := client.Scopes.UnbindContainer(scopes.NewUnbindContainerParamsWithContext(ctx).WithHandle(handle))
		if err != nil {
			switch err := err.(type) {
			case *scopes.UnbindContainerNotFound:
				// ignore error
				log.Warnf("Container %s not found by network unbind", id)
			case *scopes.UnbindContainerInternalServerError:
				return InternalServerError(err.Payload.Message)
			default:
				return InternalServerError(err.Error())
			}
		} else {
			handle = ub.Payload.Handle
			endpoints = ub.Payload.Endpoints
		}

		// unmap ports
		if err = c.mapPorts(portmap.Unmap, vc.HostConfig, c.findPortBoundNetworkEndpoint(vc.HostConfig, endpoints)); err != nil {
			return err
		}
	}

	// change the state of the container
	// TODO: We need a resolved ID from the name
	stateChangeResponse, err := client.Containers.StateChange(containers.NewStateChangeParamsWithContext(ctx).WithHandle(handle).WithState("STOPPED"))
	if err != nil {
		cache.ContainerCache().DeleteContainer(id)
		switch err := err.(type) {
		case *containers.StateChangeNotFound:
			return NotFoundError(name)
		case *containers.StateChangeDefault:
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	handle = stateChangeResponse.Payload
	wait := int32(seconds)

	_, err = client.Containers.Commit(containers.NewCommitParamsWithContext(ctx).WithHandle(handle).WithWaitTime(&wait))
	if err != nil {
		// delete from cache since all cases are 404's
		cache.ContainerCache().DeleteContainer(id)
		switch err := err.(type) {
		case *containers.CommitNotFound:
			return NotFoundError(name)
		case *containers.CommitDefault:
			return InternalServerError(err.Payload.Message)
		default:
			return InternalServerError(err.Error())
		}
	}

	return nil
}

// ContainerUnpause unpauses a container
func (c *Container) ContainerUnpause(name string) error {
	return fmt.Errorf("%s does not implement container.ContainerUnpause", ProductName())
}

// ContainerUpdate updates configuration of the container
func (c *Container) ContainerUpdate(name string, hostConfig *containertypes.HostConfig) ([]string, error) {
	return make([]string, 0, 0), fmt.Errorf("%s does not implement container.ContainerUpdate", ProductName())
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
	id := vc.ContainerID

	//retrieve client to portlayer
	client := c.containerProxy.Client()
	results, err := client.Containers.ContainerWait(containers.NewContainerWaitParamsWithContext(ctx).WithTimeout(int64(timeout.Seconds())).WithID(id))
	if err != nil {
		switch err := err.(type) {
		case *containers.ContainerWaitNotFound:
			// since the container wasn't found on the PL lets remove from the local
			// cache
			cache.ContainerCache().DeleteContainer(id)
			return -1, NotFoundError(name)
		case *containers.ContainerWaitInternalServerError:
			return -1, InternalServerError(err.Payload.Message)
		default:
			return -1, InternalServerError(err.Error())
		}
	}

	containerInfo := results.Payload
	// call to the dockerStatus function to retrieve the docker friendly exitCode
	// we are only calling for the exitCode, thus the zero time
	exitCode, _ := dockerStatus(int(*containerInfo.ProcessConfig.ExitCode),
		*containerInfo.ProcessConfig.Status,
		*containerInfo.ContainerConfig.State,
		time.Time{}, time.Time{})

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
				// or recieved a sigkill
				exitCode = 137
			} else if status == "" && exitCode == 0 {
				// the process was stopped via the cli
				// or recieved a sigterm
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
	return make([]archive.Change, 0, 0), fmt.Errorf("%s does not implement container.ContainerChanges", ProductName())
}

// ContainerInspect returns low-level information about a
// container. Returns an error if the container cannot be found, or if
// there is an error getting the data.
func (c *Container) ContainerInspect(name string, size bool, version version.Version) (interface{}, error) {
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
	if results.Payload.ProcessConfig.StartTime != nil && *results.Payload.ProcessConfig.StartTime > 0 {
		started = time.Unix(*results.Payload.ProcessConfig.StartTime, 0)
	}
	if results.Payload.ProcessConfig.StopTime != nil && *results.Payload.ProcessConfig.StopTime > 0 {
		stopped = time.Unix(*results.Payload.ProcessConfig.StopTime, 0)
	}
	// call to the dockerStatus function to retrieve the docker friendly exitCode
	exitCode, status := dockerStatus(int(*results.Payload.ProcessConfig.ExitCode),
		*results.Payload.ProcessConfig.Status,
		*results.Payload.ContainerConfig.State,
		started, stopped)

	// set the payload values
	exit := int32(exitCode)
	results.Payload.ProcessConfig.ExitCode = &exit
	results.Payload.ProcessConfig.Status = &status

	inspectJSON, err := ContainerInfoToDockerContainerInspect(vc, results.Payload, PortLayerName())
	if err != nil {
		log.Errorf("containerInfoToDockerContainerInspect failed with %s", err)
		return nil, err
	}

	log.Debugf("ContainerInspect json config = %+v\n", inspectJSON.Config)
	if inspectJSON.NetworkSettings != nil {
		log.Debugf("Docker inspect - network settings = %#v", inspectJSON.NetworkSettings)
	} else {
		log.Debugf("Docker inspect - network settings = null")
	}

	return inspectJSON, nil
}

// ContainerLogs hooks up a container's stdout and stderr streams
// configured with the given struct.
func (c *Container) ContainerLogs(name string, config *backend.ContainerLogsConfig, started chan struct{}) error {
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

	var outStream io.Writer = wf
	if !vc.Config.Tty {
		outStream = stdcopy.NewStdWriter(outStream, stdcopy.Stdout)
	}

	// Make a call to our proxy to handle the remoting
	err = c.containerProxy.StreamContainerLogs(name, outStream, started, config.Timestamps, config.Follow, since, tailLines)

	return err
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

	// Get an API client to the portlayer
	client := c.containerProxy.Client()

	containme, err := client.Containers.GetContainerList(containers.NewGetContainerListParamsWithContext(ctx).WithAll(&config.All))
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
	for _, t := range containme.Payload {
		cmd := strings.Join(t.ProcessConfig.ExecArgs, " ")
		// the docker client expects the friendly name to be prefixed
		// with a forward slash -- create a new slice and add here
		names := make([]string, 0, len(t.ContainerConfig.Names))
		for i := range t.ContainerConfig.Names {
			names = append(names, clientFriendlyContainerName(t.ContainerConfig.Names[i]))
		}
		var started time.Time
		var stopped time.Time
		if t.ProcessConfig.StartTime != nil && *t.ProcessConfig.StartTime > 0 {
			started = time.Unix(*t.ProcessConfig.StartTime, 0)
		}
		if t.ProcessConfig.StopTime != nil && *t.ProcessConfig.StopTime > 0 {
			stopped = time.Unix(*t.ProcessConfig.StopTime, 0)
		}
		// get the docker friendly status
		_, status := dockerStatus(int(*t.ProcessConfig.ExitCode), *t.ProcessConfig.Status, *t.ContainerConfig.State, started, stopped)

		c := &types.Container{
			ID:      *t.ContainerConfig.ContainerID,
			Image:   *t.ContainerConfig.RepoName,
			Created: *t.ContainerConfig.CreateTime,
			Status:  status,
			Names:   names,
			Command: cmd,
			SizeRw:  *t.ContainerConfig.StorageSize,
		}
		containers = append(containers, c)
	}
	// sort on creation time
	sort.Sort(sort.Reverse(containerByCreated(containers)))
	return containers, nil
}

// docker's container.attachBackend

// ContainerAttach attaches to logs according to the config passed in. See ContainerAttachConfig.
func (c *Container) ContainerAttach(name string, ca *backend.ContainerAttachConfig) error {
	defer trace.End(trace.Begin(name))

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)

	}

	clStdin, clStdout, clStderr, err := ca.GetStreams()
	if err != nil {
		return InternalServerError("Unable to get stdio streams for calling client")
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

	err = attachStreams(context.Background(), vc, clStdin, clStdout, clStderr, ca)

	return err
}

//------------------------------------
// ContainerAttach() Utility Functions
//------------------------------------

// attacheStreams takes the the hijacked connections from the calling client and attaches
// them to the 3 streams from the portlayer's rest server.
// clStdin, clStdout, clStderr are the hijacked connection
func attachStreams(ctx context.Context, vc *viccontainer.VicContainer, clStdin io.ReadCloser, clStdout, clStderr io.Writer, ca *backend.ContainerAttachConfig) error {
	defer clStdin.Close()

	// Cancel will close the child connections.
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	errors := make(chan error, 3)

	// For stdin, we only have a timeout for connection.  There can be a long duration before
	// the first entry so there is no timeout for response.
	plClient, transport := createNewAttachClientWithTimeouts(attachConnectTimeout, 0, attachAttemptTimeout)
	defer transport.Close()

	if ca.UseStdin {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := copyStdIn(ctx, plClient, vc, clStdin, ca.DetachKeys)
			if err != nil {
				log.Errorf("stdin (%s) error: %s", vc.ContainerID, err.Error())
			} else {
				log.Infof("stdin (%s) done", vc.ContainerID)
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
				log.Errorf("stdout (%s) error: %s", vc.ContainerID, err.Error())
			} else {
				log.Infof("stdout (%s) done", vc.ContainerID)
			}

			// no need to take action if we are canceled
			// as that means error happened somewhere else
			if ctx.Err() == context.Canceled {
				log.Infof("returning from stdout as context canceled somewhere else")
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
				log.Errorf("stderr (%s) error: %s", vc.ContainerID, err.Error())
			} else {
				log.Infof("stderr (%s) done", vc.ContainerID)
			}

			// no need to take action if we are canceled
			// as that means error happened somewhere else
			if ctx.Err() == context.Canceled {
				log.Infof("returning from stderr as context canceled somewhere else")
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

func createNewAttachClientWithTimeouts(connectTimeout, responseTimeout, responseHeaderTimeout time.Duration) (*client.PortLayer, *httpclient.Transport) {
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

	var detach bool

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
		// FIXME: If we close this stdin connection.  The portlayer does
		// not really close stdin.  This is diff from current Docker
		// behavior.  However, we're not sure why Docker even has this
		// behavior where you connect to stdin on the first time only.
		// If we really want to add this behavior, we need to add support
		// in the tether in the portlayer.
		log.Errorf("Attach stream has stdinOnce set.  VIC does not yet support this.")
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
			return NotFoundError(id)
		}
		if _, ok := err.(*interaction.ContainerGetStdoutInternalServerError); ok {
			return InternalServerError(err.Error())
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
			return NotFoundError(name)
		}

		if _, ok := err.(*interaction.ContainerGetStderrInternalServerError); ok {
			return NotFoundError(name)
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
	container.ImageID = image.ID
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

	if config.NetworkingConfig == nil {
		config.NetworkingConfig = &dnetwork.NetworkingConfig{}
	}

	if config.HostConfig == nil || config.Config == nil {
		return BadRequestError("invalid config")
	}

	// validate port bindings
	if config.HostConfig != nil {
		for _, pbs := range config.HostConfig.PortBindings {
			for _, pb := range pbs {
				if pb.HostIP != "" && pb.HostIP != "0.0.0.0" {
					return InternalServerError("host IP is not supported for port bindings")
				}

				start, end, _ := nat.ParsePortRangeToInt(pb.HostPort)
				if start != end {
					return InternalServerError("host port ranges are not supported for port bindings")
				}
			}
		}
	}

	// TODO(jzt): users other than root are not currently supported
	// We should check for USER in config.Config.Env once we support Dockerfiles.
	if config.Config.User != "" && config.Config.User != "root" {
		return InternalServerError("Failed to create container - users other than root are not currently supported")
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
	vc.Config.Labels = config.Config.Labels
	vc.HostConfig = config.HostConfig
}

func ContainerSignal(containerID string, sig uint64) error {
	// Get an API client to the portlayer
	client := PortLayerClient()
	if client == nil {
		return InternalServerError("ContainerSignal failed to create a portlayer client")
	}

	params := containers.NewContainerSignalParamsWithContext(ctx).WithID(containerID).WithSignal(int64(sig))
	if _, err := client.Containers.ContainerSignal(params); err != nil {
		return InternalServerError(err.Error())
	}

	return nil
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
		return 0, 0, fmt.Errorf("%s does not yet support '--%s'", ProductName(), opt)
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

	if config.Timestamps {
		return unsupported("timestamps")
	}

	if config.Since != "" {
		return unsupported("since")
	}

	return tailLines, since.Unix(), nil
}
