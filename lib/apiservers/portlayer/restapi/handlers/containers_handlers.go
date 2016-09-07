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
	"time"

	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"net/http"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/containers"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/options"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
	"github.com/vmware/vic/pkg/version"
)

const containerWaitTimeout = 3 * time.Minute

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
	api.ContainersContainerSignalHandler = containers.ContainerSignalHandlerFunc(handler.ContainerSignalHandler)
	api.ContainersGetContainerLogsHandler = containers.GetContainerLogsHandlerFunc(handler.GetContainerLogsHandler)
	api.ContainersContainerWaitHandler = containers.ContainerWaitHandlerFunc(handler.ContainerWaitHandler)

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
	id := uid.New().String()

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

	m := executor.ExecutorConfig{
		Common: executor.Common{
			ID:   id,
			Name: *params.CreateConfig.Name,
		},
		Version: version.GetBuild(),
		Sessions: map[string]executor.SessionConfig{
			id: {
				Common: executor.Common{
					ID:   id,
					Name: *params.CreateConfig.Name,
				},
				Tty: *params.CreateConfig.Tty,
				// FIXME: default to true for now until we can have a more sophisticated approach
				Attach: true,
				Cmd: executor.Cmd{
					Env:  params.CreateConfig.Env,
					Dir:  *params.CreateConfig.WorkingDir,
					Path: *params.CreateConfig.Path,
					Args: append([]string{*params.CreateConfig.Path}, params.CreateConfig.Args...),
				},
			},
		},
		Key:        pem.EncodeToMemory(&privateKeyBlock),
		LayerID:    *params.CreateConfig.Image,
		RepoName:   *params.CreateConfig.RepoName,
		StopSignal: *params.CreateConfig.StopSignal,
	}
	log.Infof("CreateHandler Metadata: %#v", m)

	// Create new portlayer executor and call Create on it
	h := exec.NewContainer(uid.Parse(id))
	// Create the executor.ExecutorCreateConfig
	c := &exec.ContainerCreateConfig{
		Metadata:       m,
		ParentImageID:  *params.CreateConfig.Image,
		ImageStoreName: params.CreateConfig.ImageStore.Name,
		VCHName:        options.PortLayerOptions.VCHName,
	}

	err = h.Create(ctx, session, c)
	if err != nil {
		log.Errorf("ContainerCreate error: %s", err.Error())
		return containers.NewCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}

	//  send the container id back to the caller
	return containers.NewCreateOK().WithPayload(&models.ContainerCreatedInfo{ID: id, Handle: h.String()})
}

// StateChangeHandler changes the state of a container
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
	case "CREATED":
		state = exec.StateCreated
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

	case exec.StateCreated:
		state = "CREATED"

	default:
		return containers.NewGetStateDefault(http.StatusServiceUnavailable)
	}

	return containers.NewGetStateOK().WithPayload(&models.ContainerGetStateResponse{Handle: h.String(), State: state})
}

func (handler *ContainersHandlersImpl) GetHandler(params containers.GetParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.GetHandler"))

	h := exec.GetContainer(uid.Parse(params.ID))
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

	if err := h.Commit(context.Background(), handler.handlerCtx.Session, params.WaitTime); err != nil {
		log.Errorf("CommitHandler error (%s): %s", h.String(), err)
		return containers.NewCommitDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: err.Error()})
	}

	return containers.NewCommitOK()
}

func (handler *ContainersHandlersImpl) RemoveContainerHandler(params containers.ContainerRemoveParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.RemoveContainerHandler"))

	// get the indicated container for removal
	cID := uid.Parse(params.ID)
	h := exec.GetContainer(cID)
	if h == nil {
		return containers.NewContainerRemoveNotFound()
	}

	err := h.Container.Remove(context.Background(), handler.handlerCtx.Session)
	if err != nil {
		switch err := err.(type) {
		case exec.NotFoundError:
			return containers.NewContainerRemoveNotFound()
		case exec.RemovePowerError:
			return containers.NewContainerRemoveConflict().WithPayload(&models.Error{Message: err.Error()})
		default:
			return containers.NewContainerRemoveInternalServerError()
		}
	}

	return containers.NewContainerRemoveOK()
}

func (handler *ContainersHandlersImpl) GetContainerInfoHandler(params containers.GetContainerInfoParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.GetContainerInfoHandler"))

	container := exec.ContainerInfo(params.ID)
	if container == nil {
		info := fmt.Sprintf("GetContainerInfoHandler ContainerCache miss for container(%s)", params.ID)
		log.Error(info)
		return containers.NewGetContainerInfoNotFound().WithPayload(&models.Error{Message: info})
	}

	containerInfo := convertContainerToContainerInfo(container)
	return containers.NewGetContainerInfoOK().WithPayload(containerInfo)
}

func (handler *ContainersHandlersImpl) GetContainerListHandler(params containers.GetContainerListParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.GetContainerListHandler"))

	containerVMs := exec.Containers(*params.All)
	containerList := make([]models.ContainerListInfo, 0, len(containerVMs))

	for i := range containerVMs {
		// convert to return model
		container := containerVMs[i]
		info := models.ContainerListInfo{}
		info.ContainerID = &container.ExecConfig.ID
		info.LayerID = &container.ExecConfig.LayerID
		info.Created = &container.ExecConfig.Created
		state := container.State.String()
		info.State = &state
		processStatus := container.ExecConfig.Sessions[*info.ContainerID].Started
		info.Status = &processStatus

		exitCode := int32(container.ExecConfig.Sessions[*info.ContainerID].ExitStatus)
		info.ExitCode = &exitCode

		info.Names = []string{container.ExecConfig.Name}
		info.ExecArgs = container.ExecConfig.Sessions[*info.ContainerID].Cmd.Args
		info.StorageSize = &container.VMUnsharedDisk
		info.RepoName = &container.ExecConfig.RepoName
		containerList = append(containerList, info)
	}
	return containers.NewGetContainerListOK().WithPayload(containerList)
}

func (handler *ContainersHandlersImpl) ContainerSignalHandler(params containers.ContainerSignalParams) middleware.Responder {
	defer trace.End(trace.Begin("Containers.ContainerSignal"))

	h := exec.GetContainer(uid.Parse(params.ID))
	if h == nil {
		return containers.NewContainerSignalNotFound().WithPayload(&models.Error{Message: fmt.Sprintf("container %s not found", params.ID)})
	}

	err := h.Container.Signal(context.Background(), params.Signal)
	if err != nil {
		return containers.NewContainerSignalInternalServerError().WithPayload(&models.Error{Message: err.Error()})
	}

	return containers.NewContainerSignalOK()
}

func (handler *ContainersHandlersImpl) GetContainerLogsHandler(params containers.GetContainerLogsParams) middleware.Responder {
	defer trace.End(trace.Begin(params.ID))

	h := exec.GetContainer(uid.Parse(params.ID))
	if h == nil {
		return containers.NewGetContainerLogsNotFound().WithPayload(&models.Error{
			Message: fmt.Sprintf("container %s not found", params.ID),
		})
	}

	follow := false
	tail := -1

	if params.Follow != nil {
		follow = *params.Follow
	}

	if params.Taillines != nil {
		tail = int(*params.Taillines)
	}

	reader, err := h.Container.LogReader(context.Background(), tail, follow)
	if err != nil {
		return containers.NewGetContainerLogsInternalServerError().WithPayload(&models.Error{Message: err.Error()})
	}

	detachableOut := NewFlushingReader(reader)

	return NewContainerOutputHandler("logs").WithPayload(detachableOut, params.ID)
}

func (handler *ContainersHandlersImpl) ContainerWaitHandler(params containers.ContainerWaitParams) middleware.Responder {
	defer trace.End(trace.Begin(fmt.Sprintf("%s:%d", params.ID, params.Timeout)))

	// default context timeout in seconds
	defaultTimeout := int64(containerWaitTimeout.Seconds())

	// if we have a positive timeout specified then use it
	if params.Timeout > 0 {
		defaultTimeout = params.Timeout
	}

	timeout := time.Duration(defaultTimeout) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c := exec.ContainerInfo(uid.Parse(params.ID).String())
	if c == nil {
		return containers.NewContainerWaitNotFound().WithPayload(&models.Error{
			Message: fmt.Sprintf("container %s not found", params.ID),
		})
	}

	// if the container is already stopped return
	if c.State != exec.StateRunning {
		containerInfo := convertContainerToContainerInfo(c)
		return containers.NewContainerWaitOK().WithPayload(containerInfo)
	}

	err := exec.WaitForContainerStop(ctx, params.ID)
	if err != nil {
		return containers.NewContainerWaitInternalServerError().WithPayload(&models.Error{Message: err.Error()})
	}

	c = exec.ContainerInfo(uid.Parse(params.ID).String())
	containerInfo := convertContainerToContainerInfo(c)
	return containers.NewContainerWaitOK().WithPayload(containerInfo)
}

// utility function to convert from a Container type to the API Model ContainerInfo (which should prob be called ContainerDetail)
func convertContainerToContainerInfo(container *exec.Container) *models.ContainerInfo {
	// convert the container type to the required model
	info := &models.ContainerInfo{ContainerConfig: &models.ContainerConfig{}, ProcessConfig: &models.ProcessConfig{}}

	ccid := container.ExecConfig.ID
	info.ContainerConfig.ContainerID = &ccid

	s := container.State.String()
	info.ContainerConfig.State = &s
	info.ContainerConfig.LayerID = &container.ExecConfig.LayerID
	info.ContainerConfig.RepoName = &container.ExecConfig.RepoName
	info.ContainerConfig.Created = &container.ExecConfig.Created
	info.ContainerConfig.Names = []string{container.ExecConfig.Name}

	restart := int32(container.ExecConfig.Diagnostics.ResurrectionCount)
	info.ContainerConfig.RestartCount = &restart

	tty := container.ExecConfig.Sessions[ccid].Tty
	info.ContainerConfig.Tty = &tty

	attach := container.ExecConfig.Sessions[ccid].Attach
	info.ContainerConfig.AttachStdin = &attach
	info.ContainerConfig.AttachStdout = &attach
	info.ContainerConfig.AttachStderr = &attach

	path := container.ExecConfig.Sessions[ccid].Cmd.Path
	info.ProcessConfig.ExecPath = &path

	dir := container.ExecConfig.Sessions[ccid].Cmd.Dir
	info.ProcessConfig.WorkingDir = &dir

	info.ProcessConfig.ExecArgs = container.ExecConfig.Sessions[ccid].Cmd.Args
	info.ProcessConfig.Env = container.ExecConfig.Sessions[ccid].Cmd.Env

	exitcode := int32(container.ExecConfig.Sessions[ccid].ExitStatus)
	info.ProcessConfig.ExitCode = &exitcode

	// started is a string in the vmx that is not to be confused
	// with started the datetime in the models.ContainerInfo
	status := container.ExecConfig.Sessions[ccid].Started
	info.ProcessConfig.Status = &status

	return info
}
