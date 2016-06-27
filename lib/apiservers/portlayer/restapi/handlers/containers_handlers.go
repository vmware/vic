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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"net/http"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/options"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"
)

// ContainersHandlersImpl is the receiver for all of the exec handler methods
type ContainersHandlersImpl struct {
	handlerCtx *HandlerContext
}

// Configure assigns functions to all the exec api handlers
func (handler *ContainersHandlersImpl) Configure(api *operations.PortLayerAPI, handlerCtx *HandlerContext) {
	api.ContainersCreateHandler = containers.CreateHandlerFunc(handler.CreateHandler)
	api.ContainersStateChangeHandler = containers.StateChangeHandlerFunc(handler.StateChangeHandler)
	api.ContainersGetHandler = containers.GetHandlerFunc(handler.GetHandler)
	api.ContainersCommitHandler = containers.CommitHandlerFunc(handler.CommitHandler)
	api.ContainersGetStateHandler = containers.GetStateHandlerFunc(handler.GetStateHandler)
	api.ContainersContainerRemoveHandler = containers.ContainerRemoveHandlerFunc(handler.RemoveContainerHandler)
	api.ContainersGetContainerInfoHandler = containers.GetContainerInfoHandlerFunc(handler.GetContainerInfoHandler)
	api.ContainersGetContainerListHandler = containers.GetContainerListHandlerFunc(handler.GetContainerListHandler)

	handler.handlerCtx = handlerCtx
}

// CreateHandler creates a new container
func (handler *ContainersHandlersImpl) CreateHandler(params containers.CreateParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.CreateHandler"))

	var err error

	session := handler.handlerCtx.Session

	ctx := context.Background()

	log.Debugf("Path: %#v", params.CreateConfig.Path)
	log.Debugf("Args: %#v", params.CreateConfig.Args)
	log.Debugf("Env: %#v", params.CreateConfig.Env)
	log.Debugf("WorkingDir: %#v", params.CreateConfig.WorkingDir)

	id := exec.GenerateID().String()

	// Init key for tether
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return containers.NewCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}

	m := metadata.ExecutorConfig{
		Common: metadata.Common{
			ID:   id,
			Name: *params.CreateConfig.Name,
		},
		Sessions: map[string]metadata.SessionConfig{
			id: metadata.SessionConfig{
				Common: metadata.Common{
					ID:   id,
					Name: *params.CreateConfig.Name,
				},
				Tty: *params.CreateConfig.Tty,
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
		Key: pem.EncodeToMemory(&privateKeyBlock),
	}
	log.Infof("CreateHandler Metadata: %#v", m)

	// Create new portlayer executor and call Create on it
	h := exec.NewContainer(exec.ParseID(id))
	// Create the executor.ExecutorCreateConfig
	c := &exec.ContainerCreateConfig{
		Metadata: m,

		ParentImageID:  *params.CreateConfig.Image,
		ImageStoreName: params.CreateConfig.ImageStore.Name,
		VCHName:        options.PortLayerOptions.VCHName,
	}

	err = h.Create(ctx, session, c)
	if err != nil {
		return containers.NewCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}

	//  send the container id back to the caller
	return containers.NewCreateOK().WithPayload(&models.ContainerCreatedInfo{ID: id, Handle: h.String()})
}

// ContainerStartHandler starts the container
func (handler *ContainersHandlersImpl) StateChangeHandler(params containers.StateChangeParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.StateChangeHandler"))

	h := exec.GetHandle(params.Handle)
	if h == nil {
		return containers.NewStateChangeNotFound()
	}

	var state exec.State
	switch params.State {
	case "RUNNING":
		state = exec.StateRunning

	case "STOPPED":
		state = exec.StateStopped

	default:
		return containers.NewStateChangeDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: "unknown state"})
	}

	h.SetState(state)
	return containers.NewStateChangeOK().WithPayload(h.String())
}

func (handler *ContainersHandlersImpl) GetStateHandler(params containers.GetStateParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.GetStateHandler"))

	h := exec.GetHandle(params.Handle)
	if h == nil {
		return containers.NewGetStateNotFound()
	}

	var state string
	switch h.Container.State {
	case exec.StateRunning:
		state = "RUNNING"

	case exec.StateStopped:
		state = "STOPPED"

	default:
		return containers.NewGetStateDefault(http.StatusServiceUnavailable)
	}

	return containers.NewGetStateOK().WithPayload(&models.ContainerGetStateResponse{Handle: h.String(), State: state})
}

func (handler *ContainersHandlersImpl) GetHandler(params containers.GetParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.GetHandler"))

	h := exec.GetContainer(exec.ParseID(params.ID))
	if h == nil {
		return containers.NewGetNotFound().WithPayload(&models.Error{Message: fmt.Sprintf("container %s not found", params.ID)})
	}

	return containers.NewGetOK().WithPayload(h.String())
}

func (handler *ContainersHandlersImpl) CommitHandler(params containers.CommitParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.CommitHandler"))

	h := exec.GetHandle(params.Handle)
	if h == nil {
		return containers.NewCommitNotFound().WithPayload(&models.Error{Message: "container not found"})
	}

	if err := h.Commit(context.Background(), handler.handlerCtx.Session); err != nil {
		return containers.NewCommitDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: err.Error()})
	}

	return containers.NewCommitOK()
}

func (handler *ContainersHandlersImpl) RemoveContainerHandler(params containers.ContainerRemoveParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.RemoveContainerHandler"))

	//get the indicated container for removal
	cID := exec.ParseID(params.ID)
	h := exec.GetContainer(cID)
	if h == nil {
		return containers.NewContainerRemoveNotFound()
	}

	err := h.Container.Remove(context.Background())
	if err != nil {
		return containers.NewContainerRemoveInternalServerError()
	}

	return containers.NewContainerRemoveOK()
}

func (handler *ContainersHandlersImpl) GetContainerInfoHandler(params containers.GetContainerInfoParams) middleware.Responder {
	//FIXME:  Fill in with actual code!

	return containers.NewGetContainerInfoNotFound()
}

func (handler *ContainersHandlersImpl) GetContainerListHandler(params containers.GetContainerListParams) middleware.Responder {
	//FIXME:  Fill in with actual code!

	return containers.NewGetContainerListInternalServerError().WithPayload(&models.Error{Message: "Not implemented."})
}
