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

	"golang.org/x/net/context"

	basictypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	types "github.com/docker/docker/api/types/swarm"
)

type Swarm struct {
}

func NewSwarmBackend() *Swarm {
	return &Swarm{}
}

func (s *Swarm) Init(req types.InitRequest) (string, error) {
	return "", fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) Join(req types.JoinRequest) error {
	return fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) Leave(force bool) error {
	return fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) Inspect() (types.Swarm, error) {
	return types.Swarm{}, SwarmNotSupportedError()
}

func (s *Swarm) Update(uint64, types.Spec, types.UpdateFlags) error {
	return fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) GetUnlockKey() (string, error) {
	return "", fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) UnlockSwarm(req types.UnlockRequest) error {
	return fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) GetServices(basictypes.ServiceListOptions) ([]types.Service, error) {
	return nil, SwarmNotSupportedError()
}

func (s *Swarm) GetService(string) (types.Service, error) {
	return types.Service{}, SwarmNotSupportedError()
}

func (s *Swarm) CreateService(types.ServiceSpec, string) (*basictypes.ServiceCreateResponse, error) {
	return nil, fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) UpdateService(string, uint64, types.ServiceSpec, string, string) (*basictypes.ServiceUpdateResponse, error) {
	return nil, fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) RemoveService(string) error {
	return fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) ServiceLogs(context.Context, string, *backend.ContainerLogsConfig, chan struct{}) error {
	return fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) GetNodes(basictypes.NodeListOptions) ([]types.Node, error) {
	return nil, SwarmNotSupportedError()
}

func (s *Swarm) GetNode(string) (types.Node, error) {
	return types.Node{}, SwarmNotSupportedError()
}

func (s *Swarm) UpdateNode(string, uint64, types.NodeSpec) error {
	return fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) RemoveNode(string, bool) error {
	return fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) GetTasks(basictypes.TaskListOptions) ([]types.Task, error) {
	return nil, SwarmNotSupportedError()
}

func (s *Swarm) GetTask(string) (types.Task, error) {
	return types.Task{}, SwarmNotSupportedError()
}

func (s *Swarm) GetSecrets(opts basictypes.SecretListOptions) ([]types.Secret, error) {
	return nil, fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) CreateSecret(sp types.SecretSpec) (string, error) {
	return "", fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) RemoveSecret(id string) error {
	return fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) GetSecret(id string) (types.Secret, error) {
	return types.Secret{}, fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}

func (s *Swarm) UpdateSecret(id string, version uint64, spec types.SecretSpec) error {
	return fmt.Errorf("%s does not yet support Docker Swarm", ProductName())
}
