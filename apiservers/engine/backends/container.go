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
	goexec "os/exec"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/backend"
	derr "github.com/docker/docker/errors"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/version"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"

	"github.com/vmware/vic/apiservers/portlayer/client/exec"
	"github.com/vmware/vic/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/apiservers/portlayer/models"
	"github.com/vmware/vic/pkg/trace"
)

type V1Compatibility struct {
	ID     string
	Config container.Config
}

type Container struct {
	ProductName string

	// Protects following map
	m sync.Mutex

	// FIXME: in-memory map to keep image name to vmdk name relationship
	HackMap map[string]V1Compatibility
}

// docker's container.execBackend

func (c *Container) ContainerExecCreate(config *types.ExecConfig) (string, error) {
	return "", fmt.Errorf("%s does not implement container.ContainerExecCreate", c.ProductName)
}

func (c *Container) ContainerExecInspect(id string) (*backend.ExecInspect, error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerExecInspect", c.ProductName)
}

func (c *Container) ContainerExecResize(name string, height, width int) error {
	return fmt.Errorf("%s does not implement container.ContainerExecResize", c.ProductName)
}

func (c *Container) ContainerExecStart(name string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer) error {
	return fmt.Errorf("%s does not implement container.ContainerExecStart", c.ProductName)
}

func (c *Container) ExecExists(name string) (bool, error) {
	return false, fmt.Errorf("%s does not implement container.ExecExists", c.ProductName)
}

// docker's container.copyBackend

func (c *Container) ContainerArchivePath(name string, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error) {
	return nil, nil, fmt.Errorf("%s does not implement container.ContainerArchivePath", c.ProductName)
}

func (c *Container) ContainerCopy(name string, res string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerCopy", c.ProductName)
}

func (c *Container) ContainerExport(name string, out io.Writer) error {
	return fmt.Errorf("%s does not implement container.ContainerExport", c.ProductName)
}

func (c *Container) ContainerExtractToDir(name, path string, noOverwriteDirNonDir bool, content io.Reader) error {
	return fmt.Errorf("%s does not implement container.ContainerExtractToDir", c.ProductName)
}

func (c *Container) ContainerStatPath(name string, path string) (stat *types.ContainerPathStat, err error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerStatPath", c.ProductName)
}

// docker's container.stateBackend

func (c *Container) ContainerCreate(config types.ContainerCreateConfig) (types.ContainerCreateResponse, error) {
	defer trace.End(trace.Begin("ContainerCreate"))

	//TODO: validate the config parameters
	log.Printf("config.Config = %+v", config.Config)

	// Get an API client to the portlayer
	client := PortLayerClient()
	if client == nil {
		return types.ContainerCreateResponse{},
			derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate failed to create a portlayer client"),
				http.StatusInternalServerError)
	}

	c.m.Lock()
	defer c.m.Unlock()

	layer, found := c.HackMap[config.Config.Image]
	if !found {
		// FIXME: This is a temporary workaround until we have a name resolution story.
		// Call imagec with -resolv parameter to learn the name of the vmdk and put it into in-memory map
		cmdArgs := []string{"-reference", config.Config.Image, "-resolv", "-standalone"}

		out, err := goexec.Command("/sbin/imagec", cmdArgs...).Output()
		if err != nil {
			log.Printf("imagec exit code:", err)
			return types.ContainerCreateResponse{},
				derr.NewErrorWithStatusCode(fmt.Errorf("Container look up failed"),
					http.StatusInternalServerError)
		}
		var v1 struct {
			ID string `json:"id"`
			// https://github.com/docker/engine-api/blob/master/types/container/config.go
			Config container.Config `json:"config"`
		}
		if err := json.Unmarshal(out, &v1); err != nil {
			return types.ContainerCreateResponse{},
				derr.NewErrorWithStatusCode(fmt.Errorf("Failed to unmarshall image history: %s", err),
					http.StatusInternalServerError)
		}
		log.Printf("v1 = %+v", v1)

		c.HackMap[config.Config.Image] = V1Compatibility{
			v1.ID,
			v1.Config,
		}

		layer = c.HackMap[config.Config.Image]
	}

	// Overwrite the config struct
	if len(config.Config.Cmd) == 0 {
		config.Config.Cmd = layer.Config.Cmd
	}
	if config.Config.WorkingDir == "" {
		config.Config.WorkingDir = layer.Config.WorkingDir
	}
	if len(config.Config.Entrypoint) == 0 {
		config.Config.Entrypoint = layer.Config.Entrypoint
	}
	config.Config.Env = append(config.Config.Env, layer.Config.Env...)

	log.Printf("config.Config' = %+v", config.Config)

	// Call the Exec port layer to create the container
	host, err := os.Hostname()
	if err != nil {
		return types.ContainerCreateResponse{},
			derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate got unexpected error getting hostname"),
				http.StatusInternalServerError)
	}

	plCreateParams := c.dockerContainerCreateParamsToPortlayer(config, layer.ID, host)
	createResults, err := client.Exec.ContainerCreate(plCreateParams)

	// transfer port layer swagger based response to Docker backend data structs and return to the REST front-end
	if err != nil {
		if _, isa := err.(*exec.ContainerCreateNotFound); isa {
			return types.ContainerCreateResponse{}, derr.NewRequestNotFoundError(fmt.Errorf("No such image: %s", layer.ID))
		}

		// If we get here, most likely something went wrong with the port layer API server
		return types.ContainerCreateResponse{},
			derr.NewErrorWithStatusCode(fmt.Errorf("Unknown error from the exec port layer"), http.StatusInternalServerError)
	}

	// Success!
	log.Printf("container.ContainerCreate succeeded.  Returning container id %s", *createResults.Payload.ContainerID)
	return types.ContainerCreateResponse{ID: *createResults.Payload.ContainerID}, nil
}

func (c *Container) ContainerKill(name string, sig uint64) error {
	return fmt.Errorf("%s does not implement container.ContainerKill", c.ProductName)
}

func (c *Container) ContainerPause(name string) error {
	return fmt.Errorf("%s does not implement container.ContainerPause", c.ProductName)
}

func (c *Container) ContainerRename(oldName, newName string) error {
	return fmt.Errorf("%s does not implement container.ContainerRename", c.ProductName)
}

func (c *Container) ContainerResize(name string, height, width int) error {
	return fmt.Errorf("%s does not implement container.ContainerResize", c.ProductName)
}

func (c *Container) ContainerRestart(name string, seconds int) error {
	return fmt.Errorf("%s does not implement container.ContainerRestart", c.ProductName)
}

func (c *Container) ContainerRm(name string, config *types.ContainerRmConfig) error {
	return fmt.Errorf("%s does not implement container.ContainerRm", c.ProductName)
}

func (c *Container) ContainerStart(name string, hostConfig *container.HostConfig) error {
	defer trace.End(trace.Begin("ContainerStart"))

	// Get an API client to the portlayer
	client := PortLayerClient()
	if client == nil {
		return derr.NewErrorWithStatusCode(fmt.Errorf("container.ContainerCreate failed to create a portlayer client"),
			http.StatusInternalServerError)
	}

	// handle legancy hostConfig
	if hostConfig != nil {
		// hostConfig exist for backwards compatibility.  TODO: Figure out which parameters we
		// need to look at in hostConfig
	}

	// Start the container
	// TODO: We need a resolved ID from the name
	plStartParams := &exec.ContainerStartParams{ID: name}
	_, err := client.Exec.ContainerStart(plStartParams)
	if err != nil {
		if _, isa := err.(*exec.ContainerStartNotFound); isa {
			return derr.NewRequestNotFoundError(fmt.Errorf("No such container: %s", name))
		}

		// If we get here, most likely something went wrong with the port layer API server
		return derr.NewErrorWithStatusCode(fmt.Errorf("Unknown error from the exec port layer"),
			http.StatusInternalServerError)
	}

	return nil
}

func (c *Container) ContainerStop(name string, seconds int) error {
	return fmt.Errorf("%s does not implement container.ContainerStop", c.ProductName)
}

func (c *Container) ContainerUnpause(name string) error {
	return fmt.Errorf("%s does not implement container.ContainerUnpause", c.ProductName)
}

func (c *Container) ContainerUpdate(name string, hostConfig *container.HostConfig) ([]string, error) {
	return make([]string, 0, 0), fmt.Errorf("%s does not implement container.ContainerUpdate", c.ProductName)
}

func (c *Container) ContainerWait(name string, timeout time.Duration) (int, error) {
	return 0, fmt.Errorf("%s does not implement container.ContainerWait", c.ProductName)
}

// docker's container.monitorBackend

func (c *Container) ContainerChanges(name string) ([]archive.Change, error) {
	return make([]archive.Change, 0, 0), fmt.Errorf("%s does not implement container.ContainerChanges", c.ProductName)
}

func (c *Container) ContainerInspect(name string, size bool, version version.Version) (interface{}, error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerInspect", c.ProductName)
}

func (c *Container) ContainerLogs(name string, config *backend.ContainerLogsConfig, started chan struct{}) error {
	return fmt.Errorf("%s does not implement container.ContainerLogs", c.ProductName)
}

func (c *Container) ContainerStats(name string, config *backend.ContainerStatsConfig) error {
	return fmt.Errorf("%s does not implement container.ContainerStats", c.ProductName)
}

func (c *Container) ContainerTop(name string, psArgs string) (*types.ContainerProcessList, error) {
	return nil, fmt.Errorf("%s does not implement container.ContainerTop", c.ProductName)
}

func (c *Container) Containers(config *types.ContainerListOptions) ([]*types.Container, error) {
	return nil, fmt.Errorf("%s does not implement container.Containers", c.ProductName)
}

// docker's container.attachBackend

func (c *Container) ContainerAttach(name string, cac *backend.ContainerAttachConfig) error {
	return fmt.Errorf("%s does not implement container.ContainerAttach", c.ProductName)
}

//----------
// Utility Functions
//----------

func (c *Container) dockerContainerCreateParamsToPortlayer(cc types.ContainerCreateConfig, layerID string, imageStore string) *exec.ContainerCreateParams {
	//TODO: Fill in the name
	portLayerConfig := &exec.ContainerCreateParams{Name: nil}

	portLayerConfig.CreateConfig = &models.ContainerCreateConfig{}

	// Image
	portLayerConfig.CreateConfig.Image = new(string)
	*portLayerConfig.CreateConfig.Image = layerID

	var path string
	var args []string

	// Expand cmd into entrypoint and args
	cmd := strslice.StrSlice(cc.Config.Cmd)
	if len(cc.Config.Entrypoint) != 0 {
		path, args = cc.Config.Entrypoint[0], append(cc.Config.Entrypoint[1:], cmd...)
	} else {
		path, args = cmd[0], cmd[1:]
	}

	// copy the path
	portLayerConfig.CreateConfig.Path = new(string)
	*portLayerConfig.CreateConfig.Path = path

	// copy the args
	portLayerConfig.CreateConfig.Args = make([]string, len(args))
	copy(portLayerConfig.CreateConfig.Args, args)

	// copy the env array
	portLayerConfig.CreateConfig.Env = make([]string, len(cc.Config.Env))
	copy(portLayerConfig.CreateConfig.Env, cc.Config.Env)

	// image store
	portLayerConfig.CreateConfig.ImageStore = &models.ImageStore{Name: imageStore}

	// network
	portLayerConfig.CreateConfig.NetworkDisabled = new(bool)
	*portLayerConfig.CreateConfig.NetworkDisabled = cc.Config.NetworkDisabled
	portLayerConfig.CreateConfig.NetworkSettings = toModelsNetworkConfig(cc)

	// working dir
	portLayerConfig.CreateConfig.WorkingDir = new(string)
	*portLayerConfig.CreateConfig.WorkingDir = cc.Config.WorkingDir

	log.Printf("dockerContainerCreateParamsToPortlayer = %+v", portLayerConfig.CreateConfig)
	return portLayerConfig
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
	host, err := os.Hostname()
	if err != nil {
		return "", derr.NewBadRequestError(fmt.Errorf("container.ContainerCreate got unexpected error getting hostname"))
	}

	getParams := &storage.GetImageParams{ID: imageID, StoreName: host}
	if _, err := PortLayerClient().Storage.GetImage(getParams); err != nil {
		// If the image does not exist
		if _, isa := err.(*storage.GetImageNotFound); isa {
			// return error and "No such image" which the client looks for to determine if the image didn't exist
			return "", derr.NewRequestNotFoundError(fmt.Errorf("No such image: %s", imageID))
		}

		// If we get here, most likely something went wrong with the port layer API server
		return "", derr.NewErrorWithStatusCode(fmt.Errorf("Unknown error from the storage portlayer"),
			http.StatusInternalServerError)
	}

	return host, nil
}
