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

package handlers

import (
	"fmt"
	"net"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/docker/pkg/stringid"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/apiservers/portlayer/models"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations/exec"
	"github.com/vmware/vic/apiservers/portlayer/restapi/options"
	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/portlayer/network"

	executor "github.com/vmware/vic/portlayer/exec"
)

// ExecHandlersImpl is the receiver for all of the exec handler methods
type ExecHandlersImpl struct {
	netCtx *network.Context
}

var (
	execSession = &session.Session{}
)

// Configure assigns functions to all the exec api handlers
func (handler *ExecHandlersImpl) Configure(api *operations.PortLayerAPI, netCtx *network.Context) {
	var err error

	api.ExecContainerCreateHandler = exec.ContainerCreateHandlerFunc(handler.ContainerCreateHandler)
	api.ExecContainerStartHandler = exec.ContainerStartHandlerFunc(handler.ContainerStartHandler)

	ctx := context.Background()

	sessionconfig := &session.Config{
		Service:        options.PortLayerOptions.SDK,
		Insecure:       options.PortLayerOptions.Insecure,
		Keepalive:      options.PortLayerOptions.Keepalive,
		DatacenterPath: options.PortLayerOptions.DatacenterPath,
		ClusterPath:    options.PortLayerOptions.ClusterPath,
		PoolPath:       options.PortLayerOptions.PoolPath,
		DatastorePath:  options.PortLayerOptions.DatastorePath,
		NetworkPath:    options.PortLayerOptions.NetworkPath,
	}

	execSession, err = session.NewSession(sessionconfig).Create(ctx)
	if err != nil {
		log.Fatalf("ExecHandler ERROR: %s", err)
	}

	handler.netCtx = netCtx
}

func (handler *ExecHandlersImpl) addContainerToScope(name string, ns *models.NetworkConfig) (*metadata.NetworkEndpoint, *network.Scope, error) {
	if ns == nil {
		return nil, nil, nil
	}

	var err error
	var s *network.Scope

	switch ns.NetworkName {
	// docker's default network, usually maps to the default bridge network
	case "default":
		s = handler.netCtx.DefaultScope()

	default:
		var scopes []*network.Scope
		scopes, err = handler.netCtx.Scopes(&ns.NetworkName)
		if err != nil || len(scopes) != 1 {
			return nil, nil, err
		}

		// should have only one match at this point
		s = scopes[0]
	}

	var ip *net.IP
	if ns.Address != nil {
		i := net.ParseIP(*ns.Address)
		if i == nil {
			return nil, nil, fmt.Errorf("invalid ip address")
		}

		ip = &i
	}

	var e *network.Endpoint
	e, err = s.AddContainer(name, ip)
	if err != nil {
		return nil, nil, err
	}

	ne := &metadata.NetworkEndpoint{
		IP: net.IPNet{
			IP:   e.IP(),
			Mask: e.Subnet().Mask,
		},
		Network: metadata.ContainerNetwork{
			Name: e.Scope().Name(),
			Gateway: net.IPNet{
				IP:   e.Gateway(),
				Mask: e.Subnet().Mask,
			},
		},
	}

	return ne, s, nil
}

// ContainerCreateHandler creates a new container
func (handler *ExecHandlersImpl) ContainerCreateHandler(params exec.ContainerCreateParams) middleware.Responder {
	defer trace.End(trace.Begin("ContainerCreate"))

	ctx := context.Background()

	// FIXME: Move id/name creation into personality
	id := stringid.GenerateNonCryptoID()

	var name string
	// Autogenerate a name if client doesn't specify one
	if params.Name == nil {
		name = namesgenerator.GetRandomName(0)
	} else {
		name = *params.Name
	}

	m := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   id,
			Name: name,
		},
		Sessions: map[string]metadata.SessionConfig{
			id: metadata.SessionConfig{
				Common: metadata.Common{
					ID: id,
				},
				Tty: false,
				// FIXME: default to true for now until we can have a more sophisticated approach
				Attach: true,
				Cmd: metadata.Cmd{
					Env:  params.CreateConfig.Env,
					Dir:  *params.CreateConfig.WorkingDir,
					Path: *params.CreateConfig.Path,
					Args: append([]string{*params.CreateConfig.Path}, params.CreateConfig.Args...),
				},
			},
		},
		Networks: make(map[string]metadata.NetworkEndpoint),
	}
	log.Infof("Metadata: %#v", m)

	// network config
	ns := params.CreateConfig.NetworkSettings
	ne, s, err := handler.addContainerToScope(name, ns)
	defer func() {
		if err != nil {
			log.Errorf(err.Error())
			if s != nil {
				s.RemoveContainer(name)
			}
		}
	}()

	if ne != nil {
		m.Networks[ne.Network.Name] = *ne
	}

	// Create the executor.ExecutorCreateConfig
	c := &executor.ExecutorCreateConfig{
		Metadata: m,

		ParentImageID:  *params.CreateConfig.Image,
		ImageStoreName: params.CreateConfig.ImageStore.Name,
		VCHName:        options.PortLayerOptions.VCHName,
	}

	// Create new portlayer executor and call Create on it
	executor, err := executor.NewExecutor(ctx, execSession)
	if err != nil {
		return exec.NewContainerCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}
	err = executor.Create(ctx, c)
	if err != nil {
		return exec.NewContainerCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}

	//  send the container id back to the caller
	payload := &models.ContainerCreatedInfo{
		ContainerID: &id,
	}
	return exec.NewContainerCreateOK().WithPayload(payload)

}

// ContainerStartHandler starts the container
func (handler *ExecHandlersImpl) ContainerStartHandler(params exec.ContainerStartParams) middleware.Responder {
	defer trace.End(trace.Begin("ContainerStart"))

	ctx := context.Background()

	c := &executor.ExecutorStartConfig{
		ID: params.ID,
	}
	// Create new portlayer executor and call Start on it
	executor, err := executor.NewExecutor(ctx, execSession)
	if err != nil {
		return exec.NewContainerCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}
	err = executor.Start(ctx, c)
	if err != nil {
		return exec.NewContainerCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}

	return exec.NewContainerStartOK()
}
