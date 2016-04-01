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
	"math/rand"
	"strings"

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
	"github.com/vmware/vic/pkg/vsphere/guest"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/spec"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// ExecHandlersImpl is the receiver for all of the exec handler methods
type ExecHandlersImpl struct{}

var (
	execSession = &session.Session{}
)

const (
	serialOverLANPort = 2377
)

// Configure assigns functions to all the exec api handlers
func (handler *ExecHandlersImpl) Configure(api *operations.PortLayerAPI) {
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
		DatastorePath:  options.PortLayerOptions.DatastorePath,
		NetworkPath:    options.PortLayerOptions.NetworkPath,
	}

	execSession, err = session.NewSession(sessionconfig).Create(ctx)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}

// ContainerCreateHandler creates a new container
func (handler *ExecHandlersImpl) ContainerCreateHandler(params exec.ContainerCreateParams) middleware.Responder {
	defer trace.End(trace.Begin("ContainerCreate"))

	var name string
	session := execSession

	ctx := context.Background()

	log.Debugf("Cmd: %#v", params.CreateConfig.Cmd)
	log.Debugf("EntryPoint: %#v", params.CreateConfig.EntryPoint)
	log.Debugf("Env: %#v", params.CreateConfig.Env)
	log.Debugf("WorkingDir: %#v", params.CreateConfig.WorkingDir)

	id := stringid.GenerateNonCryptoID()
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
				Cmd: metadata.Cmd{
					Path: params.CreateConfig.Cmd[0],
					Args: params.CreateConfig.Cmd,
					Env:  params.CreateConfig.Env,
					Dir:  *params.CreateConfig.WorkingDir,
				},
			},
		},
	}
	log.Infof("Metadata: %#v", m)

	specconfig := &spec.VirtualMachineConfigSpecConfig{
		// FIXME: hardcoded values
		NumCPUs:  2,
		MemoryMB: 2048,
		// FIXME: hardcoded value
		ConnectorURI: fmt.Sprintf("tcp://%s:%d", "127.0.0.1", serialOverLANPort),

		// They will be redundant with the Metadata
		ID:   id,
		Name: name,

		ParentImageID: *params.CreateConfig.Image,

		// FIXME: hardcoded value
		BootMediaPath: session.Datastore.Path(fmt.Sprintf("%s/bootstrap.iso", options.PortLayerOptions.VCHName)),
		VMPathName:    fmt.Sprintf("[%s]", session.Datastore.Name()),
		NetworkName:   strings.Split(session.Network.Reference().Value, "-")[1],

		ImageStoreName: params.CreateConfig.ImageStore.Name,

		Metadata: m,
	}
	log.Debugf("Config: %#v", specconfig)

	// Create a linux guest
	linux, err := guest.NewLinuxGuest(ctx, session, specconfig)
	if err != nil {
		return exec.NewContainerCreateNotFound().WithPayload(&models.Error{Message: fmt.Sprintf("Error constructing container vm specification: %s", err)})
	}

	// Find the Virtual Machine folder that we use
	folders, err := session.Datacenter.Folders(ctx)
	if err != nil {
		return exec.NewContainerCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}
	parent := folders.VmFolder

	// FIXME: Replace this simple logic with DRS placement
	// Pick a random host
	hosts, err := session.Datastore.AttachedClusterHosts(ctx, session.Cluster)
	if err != nil {
		return exec.NewContainerCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}
	host := hosts[rand.Intn(len(hosts))]

	// Create the vm
	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return parent.CreateVM(ctx, *linux.Spec(), session.Pool, host)
	})
	if err != nil {
		return exec.NewContainerCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}

	//  send the container id back
	payload := &models.ContainerCreatedInfo{
		ContainerID: &specconfig.ID,
	}
	return exec.NewContainerCreateOK().WithPayload(payload)

}

// ContainerStartHandler starts the container
func (handler *ExecHandlersImpl) ContainerStartHandler(params exec.ContainerStartParams) middleware.Responder {
	defer trace.End(trace.Begin("ContainerStart"))

	session := execSession
	ctx := context.Background()

	foundvm, err := session.Finder.VirtualMachine(ctx, params.ID)
	if err != nil {
		return exec.NewContainerCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}

	// Wrap the result with our version of VirtualMachine
	vm := vm.NewVirtualMachine(ctx, session, foundvm.Reference())

	// Power on
	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return vm.PowerOn(ctx)
	})
	if err != nil {
		return exec.NewContainerCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}

	return exec.NewContainerStartOK()
}
