// Copyright 2017 VMware, Inc. All Rights Reserved.
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
	"net/http"

	enginetypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/reference"
	"golang.org/x/net/context"
)

type Plugin struct {
}

func NewPluginBackend() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Disable(name string, config *enginetypes.PluginDisableConfig) error {
	return fmt.Errorf("%s does not yet support plugins", ProductName())
}

func (p *Plugin) Enable(name string, config *enginetypes.PluginEnableConfig) error {
	return fmt.Errorf("%s does not yet support plugins", ProductName())
}

func (p *Plugin) List() ([]enginetypes.Plugin, error) {
	return nil, fmt.Errorf("%s does not yet support plugins", ProductName())
}

func (p *Plugin) Inspect(name string) (*enginetypes.Plugin, error) {
	return nil, PluginNotFoundError(name)
}

func (p *Plugin) Remove(name string, config *enginetypes.PluginRmConfig) error {
	return fmt.Errorf("%s does not yet support plugins", ProductName())
}

func (p *Plugin) Set(name string, args []string) error {
	return fmt.Errorf("%s does not yet support plugins", ProductName())
}

func (p *Plugin) Privileges(ctx context.Context, ref reference.Named, metaHeaders http.Header, authConfig *enginetypes.AuthConfig) (enginetypes.PluginPrivileges, error) {
	return nil, fmt.Errorf("%s does not yet support plugins", ProductName())
}

func (p *Plugin) Pull(ctx context.Context, ref reference.Named, name string, metaHeaders http.Header, authConfig *enginetypes.AuthConfig, privileges enginetypes.PluginPrivileges, outStream io.Writer) error {
	return fmt.Errorf("%s does not yet support plugins", ProductName())
}

func (p *Plugin) Push(ctx context.Context, name string, metaHeaders http.Header, authConfig *enginetypes.AuthConfig, outStream io.Writer) error {
	return fmt.Errorf("%s does not yet support plugins", ProductName())
}

func (p *Plugin) CreateFromContext(ctx context.Context, tarCtx io.ReadCloser, options *enginetypes.PluginCreateOptions) error {
	return fmt.Errorf("%s does not yet support plugins", ProductName())
}
