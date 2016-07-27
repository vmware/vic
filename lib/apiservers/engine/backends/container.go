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
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
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
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/version"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	dnetwork "github.com/docker/engine-api/types/network"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/go-connections/nat"

	"strconv"

	"github.com/google/uuid"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"

	"github.com/vishvananda/netlink"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	"github.com/vmware/vic/lib/apiservers/engine/backends/portmap"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/interaction"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/scopes"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
)

const bridgeIfaceName = "bridge"

var (
	clientIfaceName = "client"

	defaultScope struct {
		sync.Mutex
		scope string
	}

	portMapper portmap.PortMapper
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

//TODO: gotta be a better way...
// type and funcs to provide sorting by created date
type containerByCreated []*types.Container

func (r containerByCreated) Len() int           { return len(r) }
func (r containerByCreated) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r containerByCreated) Less(i, j int) bool { return r[i].Created < r[j].Created }

// Container struct represents the Container
type Container struct {
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
)

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

	var err error

	// validate port bindings
	if config.HostConfig != nil {
		for _, pbs := range config.HostConfig.PortBindings {
			for _, pb := range pbs {
				if pb.HostIP != "" {
					return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(fmt.Errorf("host IP is not supported for port bindings"), http.StatusInternalServerError)
				}

				start, end, _ := nat.ParsePortRangeToInt(pb.HostPort)
				if start != end {
					return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(fmt.Errorf("host port ranges are not supported for port bindings"), http.StatusInternalServerError)
				}
			}
		}
	}

	//TODO: validate the config parameters
	log.Debugf("Image fetch section - Container Create")
	log.Debugf("config.Config = %+v", config.Config)
	// Get an API client to the portlayer
	client := PortLayerClient()
	if client == nil {
		return types.ContainerCreateResponse{},
			derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate failed to create a portlayer client"),
				http.StatusInternalServerError)
	}

	// bail early if container name already exists
	if exists := cache.ContainerCache().GetContainer(config.Name); exists != nil {
		err := fmt.Errorf("Conflict. The name \"%s\" is already in use by container %s. You have to remove (or rename) that container to be able to re use that name.", config.Name, exists.ContainerID)
		log.Errorf("%s", err.Error())
		return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(err, http.StatusConflict)
	}

	// get the image from the cache
	image, err := cache.ImageCache().GetImage(config.Config.Image)
	if err != nil {
		// if no image found then error thrown and a pull
		// will be initiated by the docker client
		log.Errorf("ContainerCreate: image %s error: %s", config.Config.Image, err.Error())
		return types.ContainerCreateResponse{}, err
	}

	// provide basic container config via the image
	container := viccontainer.NewVicContainer()
	container.ID = image.ID
	container.Config = image.Config

	// set default configuration values and those supplied by image where needed
	container.SetConfigOptions(config.Config)

	// TODO(jzt): users other than root are not currently supported
	// We should check for USER in config.Config.Env once we support Dockerfiles.
	if config.Config.User != "" && config.Config.User != "root" {
		return types.ContainerCreateResponse{},
			derr.NewErrorWithStatusCode(fmt.Errorf("Failed to create container - users other than root are not currently supported"),
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
		return types.ContainerCreateResponse{}, derr.NewRequestNotFoundError(fmt.Errorf("No command specified"))
	}

	// Call the Exec port layer to create the container
	host, err := guest.UUID()
	if err != nil {
		return types.ContainerCreateResponse{},
			derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate got unexpected error getting VCH UUID"),
				http.StatusInternalServerError)
	}

	plCreateParams := c.dockerContainerCreateParamsToPortlayer(config, container.ID, host)
	createResults, err := client.Containers.Create(plCreateParams)
	// transfer port layer swagger based response to Docker backend data structs and return to the REST front-end
	if err != nil {
		if _, ok := err.(*containers.CreateNotFound); ok {
			err = fmt.Errorf("No such image: %s", container.ID)
			log.Errorf(err.Error())
			return types.ContainerCreateResponse{}, derr.NewRequestNotFoundError(err)
		}

		// If we get here, most likely something went wrong with the port layer API server
		return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
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
	//Volume Attachment Section
	log.Debugf("Container.ContainerCreate - VolumeSection")
	log.Debugf("Raw Volume arguments : binds:  %#v : volumes : %#v", config.HostConfig.Binds, config.Config.Volumes)
	var joinList []volumeFields

	joinList, err = processAnonymousVolumes(&h, config.Config.Volumes, client)
	if err != nil {
		return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(fmt.Errorf("%s", err), http.StatusBadRequest)
	}

	volumeSubset, err := processSpecifiedVolumes(config.HostConfig.Binds)
	if err != nil {
		return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(fmt.Errorf("%s", err), http.StatusBadRequest)
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
			switch err := err.(type) {
			case *storage.VolumeJoinInternalServerError:
				return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

			case *storage.VolumeJoinDefault:
				return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)
			default:
				return types.ContainerCreateResponse{}, derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
			}
		}

		h = res.Payload
	}
	// commit the create op
	_, err = client.Containers.Commit(containers.NewCommitParams().WithHandle(h))
	if err != nil {
		err = fmt.Errorf("No such image: %s", container.ID)
		log.Errorf("%s", err.Error())
		// FIXME: Containers.Commit returns more errors than it's swagger spec says.
		// When no image exist, it also sends back non swagger errors.  We should fix
		// this once Commit returns correct error codes.
		return types.ContainerCreateResponse{}, derr.NewRequestNotFoundError(err)
	}

	// Container created ok, overwrite the container params in the container store as
	// these are the parameters that the containers were actually created with
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
	container.HostConfig = config.HostConfig
	container.ContainerID = createResults.Payload.ID
	container.Name = config.Name

	log.Debugf("Container create: %#v", container)
	cache.ContainerCache().AddContainer(container)

	// Success!
	log.Debugf("container.ContainerCreate succeeded.  Returning container handle %s", *createResults.Payload)
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
	defer trace.End(trace.Begin("ContainerRestart"))

	err := c.containerStop(name, seconds, false)
	if err != nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("Stop failed with: %s", err), http.StatusInternalServerError)
	}

	err = c.containerStart(name, nil, false)
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

	// Get the portlayer Client API
	client := PortLayerClient()
	if client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerRm failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	// TODO: Pass this RemoveVolume flag to somewhere
	_ = &config.RemoveVolume

	// Use the force and stop the container first
	if config.ForceRemove {
		c.containerStop(name, 0, true)
	}

	//call the remove directly on the name. No need for using a handle.
	_, err := client.Containers.ContainerRemove(containers.NewContainerRemoveParams().WithID(name))
	if err != nil {
		if _, ok := err.(*containers.ContainerRemoveNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
		}
		if errModel, ok := err.(*containers.ContainerRemoveDefault); ok {
			return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer : %s", errModel.Payload.Message), http.StatusInternalServerError)
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer : %s", err), http.StatusInternalServerError)
	}
	// delete container from the cache
	cache.ContainerCache().DeleteContainer(name)

	return nil
}

// ContainerStart starts a container.
func (c *Container) ContainerStart(name string, hostConfig *container.HostConfig) error {
	defer trace.End(trace.Begin("ContainerStart"))
	return c.containerStart(name, hostConfig, true)
}

func (c *Container) containerStart(name string, hostConfig *container.HostConfig, bind bool) error {
	var err error

	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc != nil {
		log.Debugf("Found %q in cache as %q", name, vc.ContainerID)
		name = vc.ContainerID
	}

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
	} else if vc != nil {
		hostConfig = vc.HostConfig
	}

	if vc != nil && hostConfig.NetworkMode.NetworkName() == "" {
		hostConfig.NetworkMode = vc.HostConfig.NetworkMode
	}

	// get a handle to the container
	var getRes *containers.GetOK
	getRes, err = client.Containers.Get(containers.NewGetParams().WithID(name))
	if err != nil {
		if _, ok := err.(*containers.GetNotFound); ok {
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
		}
		if errModel, ok := err.(*containers.GetDefault); ok {
			return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer : %s", errModel.Payload.Message), http.StatusInternalServerError)
		}
		return derr.NewErrorWithStatusCode(fmt.Errorf("server error from portlayer : %s", err), http.StatusInternalServerError)
	}

	h := getRes.Payload

	var endpoints []*models.EndpointConfig
	// bind network
	if bind {
		var bindRes *scopes.BindContainerOK
		bindRes, err = client.Scopes.BindContainer(scopes.NewBindContainerParams().WithHandle(h))
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

		h = bindRes.Payload.Handle
		endpoints = bindRes.Payload.Endpoints

		// unbind in case we fail later
		defer func() {
			if err != nil {
				client.Scopes.UnbindContainer(scopes.NewUnbindContainerParams().WithHandle(h))
			}
		}()

	}

	// change the state of the container
	// TODO: We need a resolved ID from the name
	var stateChangeRes *containers.StateChangeOK
	stateChangeRes, err = client.Containers.StateChange(containers.NewStateChangeParams().WithHandle(h).WithState("RUNNING"))
	if err != nil {
		switch err := err.(type) {
		case *containers.StateChangeNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s", err.Payload.Message))

		case *containers.StateChangeDefault:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s", err.Payload.Message))

		default:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s", err))
		}
	}

	h = stateChangeRes.Payload

	if bind {
		e := c.findPortBoundNetworkEndpoint(hostConfig, endpoints)
		if err = c.mapPorts(portmap.Map, hostConfig, e); err != nil {
			err = fmt.Errorf("error mapping ports: %s", err)
			log.Error(err)
			return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
		}

		defer func() {
			if err != nil {
				c.mapPorts(portmap.Unmap, hostConfig, e)
			}
		}()
	}

	// commit the handle; this will reconfigure and start the vm
	_, err = client.Containers.Commit(containers.NewCommitParams().WithHandle(h))
	if err != nil {
		switch err := err.(type) {
		case *containers.CommitNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s", err.Payload.Message))

		case *containers.CommitDefault:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s", err.Payload.Message))
		default:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s", err))
		}
	}

	return nil
}

func (c *Container) mapPorts(op portmap.Operation, hostconfig *container.HostConfig, endpoint *models.EndpointConfig) error {
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

		for _, pb := range pbs {
			var hostPort int
			hostPort, err = strconv.Atoi(pb.HostPort)
			if err != nil {
				return err
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

	client := PortLayerClient()
	listRes, err := client.Scopes.List(scopes.NewListParams().WithIDName("default"))
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

func (c *Container) findPortBoundNetworkEndpoint(hostconfig *container.HostConfig, endpoints []*models.EndpointConfig) *models.EndpointConfig {
	if len(hostconfig.PortBindings) == 0 {
		return nil
	}

	// check if the port binding network is a bridge type
	listRes, err := PortLayerClient().Scopes.List(scopes.NewListParams().WithIDName(hostconfig.NetworkMode.NetworkName()))
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
	defer trace.End(trace.Begin("ContainerStop"))
	return c.containerStop(name, seconds, true)
}

func (c *Container) containerStop(name string, seconds int, unbound bool) error {
	// Look up the container name in the metadata cache to get long ID
	vc := cache.ContainerCache().GetContainer(name)
	if vc != nil {
		log.Debugf("Found %q in cache as %q", name, vc.ContainerID)
		name = vc.ContainerID
	}

	//retrieve client to portlayer
	client := PortLayerClient()
	if client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("container.containerStop failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	getResponse, err := client.Containers.Get(containers.NewGetParams().WithID(name))
	if err != nil {
		switch err := err.(type) {

		case *containers.GetNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))

		case *containers.GetDefault:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s", err.Payload.Message))

		default:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s", err))
		}
	}

	handle := getResponse.Payload

	// we have a container on the PL side lets check the state before proceeding
	// ignore the error  since others will be checking below..this is an attempt to short circuit the op
	// TODO: can be replaced with simple cache check once power events are propigated to persona
	infoResponse, _ := client.Containers.GetContainerInfo(containers.NewGetContainerInfoParams().WithID(name))
	if *infoResponse.Payload.ContainerConfig.State == "Stopped" || *infoResponse.Payload.ContainerConfig.State == "Created" {
		return nil
	}

	if unbound {
		// get the endpoints for the container
		conEpsRes, err := client.Scopes.GetContainerEndpoints(scopes.NewGetContainerEndpointsParams().WithHandleOrID(handle))
		if err != nil {
			switch err := err.(type) {
			case *scopes.GetContainerEndpointsNotFound:
				return derr.NewRequestNotFoundError(fmt.Errorf("container %s not found", name))

			case *scopes.GetContainerEndpointsInternalServerError:
				return derr.NewErrorWithStatusCode(fmt.Errorf("%s", err.Payload.Message), http.StatusInternalServerError)
			}
		}

		ub, err := client.Scopes.UnbindContainer(scopes.NewUnbindContainerParams().WithHandle(handle))
		if err != nil {
			switch err := err.(type) {
			case *scopes.UnbindContainerNotFound:
				return derr.NewRequestNotFoundError(fmt.Errorf("container %s not found", name))

			case *scopes.UnbindContainerInternalServerError:
				return derr.NewErrorWithStatusCode(fmt.Errorf(err.Payload.Message), http.StatusInternalServerError)

			default:
				return derr.NewErrorWithStatusCode(err, http.StatusInternalServerError)
			}
		}

		handle = ub.Payload

		// unmap ports
		if err = c.mapPorts(portmap.Unmap, vc.HostConfig, c.findPortBoundNetworkEndpoint(vc.HostConfig, conEpsRes.Payload)); err != nil {
			return err
		}
	}

	// change the state of the container
	// TODO: We need a resolved ID from the name
	stateChangeResponse, err := client.Containers.StateChange(containers.NewStateChangeParams().WithHandle(handle).WithState("STOPPED"))
	if err != nil {
		switch err := err.(type) {
		case *containers.StateChangeNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s ", err.Payload.Message))

		case *containers.StateChangeDefault:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s ", err.Payload.Message))

		default:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s ", err))
		}
	}

	handle = stateChangeResponse.Payload

	_, err = client.Containers.Commit(containers.NewCommitParams().WithHandle(handle))
	if err != nil {
		switch err := err.(type) {
		case *containers.CommitNotFound:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s ", err.Payload.Message))

		case *containers.CommitDefault:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s ", err.Payload.Message))

		default:
			return derr.NewRequestNotFoundError(fmt.Errorf("server error from portlayer : %s ", err))
		}

	}

	return nil
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
	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
	}

	log.Debugf("Found %q in cache as %q", name, vc.ContainerID)
	name = vc.ContainerID

	client := PortLayerClient()
	if client == nil {
		return nil, derr.NewErrorWithStatusCode(fmt.Errorf("Failed to get portlayer client"), http.StatusInternalServerError)
	}

	results, err := client.Containers.GetContainerInfo(containers.NewGetContainerInfoParams().WithID(name))
	if err != nil {
		switch err := err.(type) {
		case *containers.GetContainerInfoNotFound:
			return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
		case *containers.GetContainerInfoInternalServerError:
			return nil, derr.NewErrorWithStatusCode(fmt.Errorf("Error from portlayer: %#v", err.Payload), http.StatusInternalServerError)
		default:
			return nil, derr.NewErrorWithStatusCode(fmt.Errorf("Unknown error from the container portlayer"), http.StatusInternalServerError)
		}
	}

	inspectJSON, err := containerInfoToDockerContainerInspect(vc, results.Payload)
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

	// Get an API client to the portlayer
	portLayerClient := PortLayerClient()
	if portLayerClient == nil {
		return nil, derr.NewErrorWithStatusCode(fmt.Errorf("container.Containers failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	containme, err := portLayerClient.Containers.GetContainerList(containers.NewGetContainerListParams().WithAll(&config.All))
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
	// sort on creation time
	sort.Sort(sort.Reverse(containerByCreated(containers)))
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

	err = attachStreams(context.Background(), vc, clStdin, clStdout, clStderr, ca)

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

	nc.Ports = make([]string, len(cc.HostConfig.PortBindings))
	i := 0
	for p := range cc.HostConfig.PortBindings {
		nc.Ports[i] = string(p)
		i++
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

// containerInfoToDockerContainerInspect() takes a ContainerInfo swagger-based struct
// returned from VIC's port layer and creates an engine-api based container inspect struct.
// There maybe other asset gathering if ContainerInfo does not have all the information
func containerInfoToDockerContainerInspect(vc *viccontainer.VicContainer, info *models.ContainerInfo) (*types.ContainerJSON, error) {
	if vc == nil || info == nil || info.ContainerConfig == nil {
		return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", vc.ContainerID))
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
		if info.ProcessConfig.Started != nil {
			swaggerTime := time.Time(*info.ProcessConfig.Started)
			containerState.StartedAt = swaggerTime.Format(time.RFC3339Nano)
		}

		if info.ProcessConfig.Finished != nil {
			swaggerTime := time.Time(*info.ProcessConfig.Finished)
			containerState.FinishedAt = swaggerTime.Format(time.RFC3339Nano)
		}
	}

	inspectJSON := &types.ContainerJSON{
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
			HostConfig:      hostConfigFromContainerInfo(vc, info),
			GraphDriver:     types.GraphDriverData{Name: PortLayerName()},
			SizeRw:          nil,
			SizeRootFs:      nil,
		},
		Mounts:          mountsFromContainerInfo(vc, info),
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
		if info.ContainerConfig.Created != nil {
			inspectJSON.Created = time.Unix(*info.ContainerConfig.Created, 0).String()
		}
		if len(info.ContainerConfig.Names) > 0 {
			inspectJSON.Name = info.ContainerConfig.Names[0]
		}
	}

	return inspectJSON, nil
}

// hostConfigFromContainerInfo() gets the hostconfig that is passed to the backend during
// docker create and updates any needed info
func hostConfigFromContainerInfo(vc *viccontainer.VicContainer, info *models.ContainerInfo) *container.HostConfig {
	if vc == nil || vc.HostConfig == nil || info == nil {
		return nil
	}

	// Create a copy of the created container's hostconfig.  This is passed in during
	// container create
	hostConfig := *vc.HostConfig

	// Resources don't really map well to VIC so we leave most of them empty. If we look
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

	hostConfig.VolumeDriver = PortLayerName()
	hostConfig.Resources = resourceConfig

	if len(info.ScopeConfig) > 0 {
		if info.ScopeConfig[0].DNS != nil {
			hostConfig.DNS = info.ScopeConfig[0].DNS
		}

		hostConfig.NetworkMode = container.NetworkMode(info.ScopeConfig[0].ScopeType)
	}

	return &hostConfig
}

// mountsFromContainerInfo()
func mountsFromContainerInfo(vc *viccontainer.VicContainer, info *models.ContainerInfo) []types.MountPoint {
	if vc == nil || info == nil {
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

	if info.ContainerConfig.LayerID != nil {
		container.Image = *info.ContainerConfig.LayerID // Name of the image as it was passed by the operator (eg. could be symbolic)
	}
	if info.ContainerConfig.Labels != nil {
		container.Labels = info.ContainerConfig.Labels // List of labels set to this container
	}

	// Fill in information about the process
	if info.ProcessConfig.Env != nil {
		container.Env = info.ProcessConfig.Env // List of environment variable to set in the container
	}
	if info.ProcessConfig.ExecPath != nil {
		container.Cmd = append(container.Cmd, *info.ProcessConfig.ExecPath) // Command to run when starting the container
	}
	if info.ProcessConfig.ExecArgs != nil {
		container.Cmd = append(container.Cmd, info.ProcessConfig.ExecArgs[1:]...)
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

// helper function to format the container name
// to the docker client approved format
func clientFriendlyContainerName(name string) string {
	return fmt.Sprintf("/%s", name)
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
