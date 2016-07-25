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

package container

import (
	"fmt"
	"strings"

	containertypes "github.com/docker/engine-api/types/container"
)

// DefaultEnvPath defines the default PATH environment variable to use in a container
const DefaultEnvPath = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

// VicContainer is VIC's abridged version of Docker's container object.
type VicContainer struct {
	Name        string
	ID          string
	ContainerID string
	Config      *containertypes.Config
	HostConfig  *containertypes.HostConfig
}

// NewVicContainer returns a reference to a new VicContainer
func NewVicContainer() *VicContainer {
	return &VicContainer{
		Config: &containertypes.Config{},
	}
}

// SetConfigOptions is a place to add necessary container configuration
// values that were not explicitly supplied by the user
func (c *VicContainer) SetConfigOptions(config *containertypes.Config) {
	// Overwrite or append the image's config from the CLI with the metadata from the image's
	// layer metadata where appropriate
	if len(config.Cmd) == 0 {
		config.Cmd = c.Config.Cmd
	}
	if config.WorkingDir == "" {
		config.WorkingDir = c.Config.WorkingDir
	}
	if len(config.Entrypoint) == 0 {
		config.Entrypoint = c.Config.Entrypoint
	}

	// set up environment
	setEnv(config, c.Config)
}

func setEnv(config, imageConfig *containertypes.Config) {
	// Set PATH in ENV if needed
	setPath(config, imageConfig)

	containerEnv := make(map[string]string, len(config.Env))
	for _, env := range config.Env {
		kv := strings.Split(env, "=")
		containerEnv[kv[0]] = kv[1]
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
		key := strings.Split(imageEnv, "=")[0]
		// is environment variable already set in container config?
		if _, ok := containerEnv[key]; !ok {
			// no? let's copy it from the image config
			config.Env = append(config.Env, imageEnv)
		}
	}
}

func setPath(config, imageConfig *containertypes.Config) {
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
	config.Env = append(config.Env, fmt.Sprintf("PATH=%s", DefaultEnvPath))
}
