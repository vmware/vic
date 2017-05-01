// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/containers"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/iolog"
	"github.com/vmware/vic/lib/migration/feature"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/metrics"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
	"github.com/vmware/vic/pkg/version"
)

const (
	containerWaitTimeout = 3 * time.Minute
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
	api.ContainersContainerSignalHandler = containers.ContainerSignalHandlerFunc(handler.ContainerSignalHandler)
	api.ContainersGetContainerLogsHandler = containers.GetContainerLogsHandlerFunc(handler.GetContainerLogsHandler)
	api.ContainersContainerWaitHandler = containers.ContainerWaitHandlerFunc(handler.ContainerWaitHandler)
	api.ContainersContainerRenameHandler = containers.ContainerRenameHandlerFunc(handler.RenameContainerHandler)
	api.ContainersGetContainerStatsHandler = containers.GetContainerStatsHandlerFunc(handler.GetContainerStatsHandler)

	handler.handlerCtx = handlerCtx
}

// CreateHandler creates a new container
func (handler *ContainersHandlersImpl) CreateHandler(params containers.CreateParams) middleware.Responder {
	defer trace.End(trace.Begin(""))

	var err error

	session := handler.handlerCtx.Session

	ctx := context.Background()

	id := uid.New().String()

	// Init key for tether
	// #nosec: RSA keys should be at least 2048 bits
	// Size is 512 because key validation is not performed - see GitHub #2849
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return containers.NewCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}

	m := &executor.ExecutorConfig{
		ExecutorConfigCommon: executor.ExecutorConfigCommon{
			ID:   id,
			Name: params.CreateConfig.Name,
		},
		CreateTime: time.Now().UTC().Unix(),
		Version:    version.GetBuild(),
		Key:        pem.EncodeToMemory(&privateKeyBlock),
		LayerID:    params.CreateConfig.Layer,
		ImageID:    params.CreateConfig.Image,
		RepoName:   params.CreateConfig.RepoName,
	}

	if params.CreateConfig.Annotations != nil && len(params.CreateConfig.Annotations) > 0 {
		m.Annotations = make(map[string]string)
		for k, v := range params.CreateConfig.Annotations {
			m.Annotations[k] = v
		}
	}

	// Create the executor.ExecutorCreateConfig
	c := &exec.ContainerCreateConfig{
		Metadata:       m,
		ParentImageID:  params.CreateConfig.Layer,
		ImageStoreName: params.CreateConfig.ImageStore.Name,
		Resources: exec.Resources{
			NumCPUs:  params.CreateConfig.NumCpus,
			MemoryMB: params.CreateConfig.MemoryMB,
		},
	}

	h, err := exec.Create(ctx, session, c)
	if err != nil {
		log.Errorf("ContainerCreate error: %s", err.Error())
		return containers.NewCreateNotFound().WithPayload(&models.Error{Message: err.Error()})
	}

	//  send the container id back to the caller
	return containers.NewCreateOK().WithPayload(&models.ContainerCreatedInfo{ID: id, Handle: h.String()})
}

// StateChangeHandler changes the state of a container
func (handler *ContainersHandlersImpl) StateChangeHandler(params containers.StateChangeParams) middleware.Responder {
	defer trace.End(trace.Begin(fmt.Sprintf("handle(%s)", params.Handle)))

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

	h.SetTargetState(state)
	return containers.NewStateChangeOK().WithPayload(h.String())
}

func (handler *ContainersHandlersImpl) GetStateHandler(params containers.GetStateParams) middleware.Responder {
	defer trace.End(trace.Begin(fmt.Sprintf("handle(%s)", params.Handle)))

	// NOTE: I've no idea why GetStateHandler takes a handle instead of an ID - hopefully there was a reason for an inspection
	// operation to take this path
	h := exec.GetHandle(params.Handle)
	if h == nil || h.ExecConfig == nil {
		return containers.NewGetStateNotFound()
	}

	container := exec.Containers.Container(h.ExecConfig.ID)
	if container == nil {
		return containers.NewGetStateNotFound()
	}

	var state string
	switch container.CurrentState() {
	case exec.StateRunning:
		state = "RUNNING"

	case exec.StateStopped:
		state = "STOPPED"

	case exec.StateCreated:
		state = "CREATED"

	default:
		return containers.NewGetStateDefault(http.StatusServiceUnavailable)
	}

	return containers.NewGetStateOK().WithPayload(
		&models.ContainerGetStateResponse{
			Handle: h.String(),
			State:  state,
		})
}

func (handler *ContainersHandlersImpl) GetHandler(params containers.GetParams) middleware.Responder {
	defer trace.End(trace.Begin(params.ID))

	h := exec.GetContainer(context.Background(), uid.Parse(params.ID))
	if h == nil {
		return containers.NewGetNotFound().WithPayload(&models.Error{Message: fmt.Sprintf("container %s not found", params.ID)})
	}

	return containers.NewGetOK().WithPayload(h.String())
}

func (handler *ContainersHandlersImpl) CommitHandler(params containers.CommitParams) middleware.Responder {
	defer trace.End(trace.Begin(fmt.Sprintf("handle(%s)", params.Handle)))

	h := exec.GetHandle(params.Handle)
	if h == nil {
		return containers.NewCommitNotFound().WithPayload(&models.Error{Message: "container not found"})
	}

	if err := h.Commit(context.Background(), handler.handlerCtx.Session, params.WaitTime); err != nil {
		log.Errorf("CommitHandler error on handle(%s) for %s: %#v", h.String(), h.ExecConfig.ID, err)
		switch err := err.(type) {
		case exec.ConcurrentAccessError:
			return containers.NewCommitConflict().WithPayload(&models.Error{Message: err.Error()})
		default:
			return containers.NewCommitDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: err.Error()})
		}
	}

	return containers.NewCommitOK()
}

func (handler *ContainersHandlersImpl) RemoveContainerHandler(params containers.ContainerRemoveParams) middleware.Responder {
	defer trace.End(trace.Begin(params.ID))

	// get the indicated container for removal
	container := exec.Containers.Container(params.ID)
	if container == nil {
		return containers.NewContainerRemoveNotFound()
	}

	// NOTE: this should allowing batching of operations, as with Create, Start, Stop, et al
	err := container.Remove(context.Background(), handler.handlerCtx.Session)
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
	defer trace.End(trace.Begin(params.ID))

	container := exec.Containers.Container(params.ID)
	if container == nil {
		info := fmt.Sprintf("GetContainerInfoHandler ContainerCache miss for container(%s)", params.ID)
		log.Error(info)
		return containers.NewGetContainerInfoNotFound().WithPayload(&models.Error{Message: info})
	}

	// Refresh to get up to date network info
	container.Refresh(context.Background())
	containerInfo := convertContainerToContainerInfo(container.Info())
	return containers.NewGetContainerInfoOK().WithPayload(containerInfo)
}

// type and funcs to provide sorting by created date
type containerByCreated []*models.ContainerInfo

func (r containerByCreated) Len() int      { return len(r) }
func (r containerByCreated) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r containerByCreated) Less(i, j int) bool {
	return r[i].ContainerConfig.CreateTime < r[j].ContainerConfig.CreateTime
}

func (handler *ContainersHandlersImpl) GetContainerListHandler(params containers.GetContainerListParams) middleware.Responder {
	defer trace.End(trace.Begin(""))

	var state *exec.State
	if params.All != nil && !*params.All {
		state = new(exec.State)
		*state = exec.StateRunning
	}

	containerVMs := exec.Containers.Containers(state)
	containerList := make([]*models.ContainerInfo, 0, len(containerVMs))

	for _, container := range containerVMs {
		// convert to return model
		info := convertContainerToContainerInfo(container.Info())
		containerList = append(containerList, info)
	}

	sort.Sort(sort.Reverse(containerByCreated(containerList)))
	return containers.NewGetContainerListOK().WithPayload(containerList)
}

func (handler *ContainersHandlersImpl) ContainerSignalHandler(params containers.ContainerSignalParams) middleware.Responder {
	defer trace.End(trace.Begin(params.ID))

	// NOTE: I feel that this should be in a Commit path for consistency
	// it would allow phrasings such as:
	// 1. join Volume to container
	// 2. send HUP to primary process
	// Only really relevant when we can connect networks or join volumes live
	container := exec.Containers.Container(params.ID)
	if container == nil {
		return containers.NewContainerSignalNotFound().WithPayload(&models.Error{Message: fmt.Sprintf("container %s not found", params.ID)})
	}

	err := container.Signal(context.Background(), params.Signal)
	if err != nil {
		return containers.NewContainerSignalInternalServerError().WithPayload(&models.Error{Message: err.Error()})
	}

	return containers.NewContainerSignalOK()
}

func (handler *ContainersHandlersImpl) GetContainerStatsHandler(params containers.GetContainerStatsParams) middleware.Responder {
	defer trace.End(trace.Begin(params.ID))

	c := exec.Containers.Container(params.ID)
	if c == nil {
		return containers.NewGetContainerStatsNotFound()
	}

	r, w := io.Pipe()
	enc := json.NewEncoder(w)
	flusher := NewFlushingReader(r)

	// subscribe to metrics
	// currently all stats requests will be a subscription and it will
	// be the responsibility of the caller to close the connection
	// and there by release the subscription
	ch, err := metrics.Supervisor.VMCollector().Subscribe(c)
	if err != nil {
		log.Errorf("unable to subscribe container(%s) to stats stream: %s", params.ID, err)
		return containers.NewGetContainerStatsInternalServerError()
	}
	log.Debugf("container(%s) stats stream subscribed @ %d", params.ID, &ch)

	// closer will be run when the http transport is closed
	cleaner := func() {
		log.Debugf("unsubscribing %s from stats %d", params.ID, &ch)
		metrics.Supervisor.VMCollector().Unsubscribe(c, ch)
		closePipe(r, w)
	}

	// routine that will listen for new metrics and encode to provided output stream
	// unsubscription or error will exit the routine
	go func() {
		for {
			select {
			case metric, ok := <-ch:
				if !ok {
					log.Debugf("container stats complete for %s @ %d", params.ID, &ch)
					return
				}
				err := enc.Encode(metric)
				if err != nil {
					log.Errorf("encoding error [%s] for container(%s) stats @ %d - stream(%t)", err, params.ID, &ch, params.Stream)
					return
				}
			}
		}
	}()

	return NewStreamOutputHandler("containerStats").WithPayload(flusher, params.ID, cleaner)
}

func (handler *ContainersHandlersImpl) GetContainerLogsHandler(params containers.GetContainerLogsParams) middleware.Responder {
	defer trace.End(trace.Begin(params.ID))

	container := exec.Containers.Container(params.ID)
	if container == nil {
		return containers.NewGetContainerLogsNotFound()
	}

	follow := false
	tail := -1

	if params.Follow != nil {
		follow = *params.Follow
	}

	if params.Taillines != nil {
		tail = int(*params.Taillines)
	}

	reader, err := container.LogReader(context.Background(), tail, follow)
	if err != nil {
		// Do not return an error here.  It's a workaround for a panic similar to #2594
		return containers.NewGetContainerLogsInternalServerError()
	}

	// containers with DataVersion > 0 will use updated output logging on the backend
	if container.DataVersion > 0 {
		ts := false
		if params.Timestamp != nil {
			ts = *params.Timestamp
		}

		// wrap the reader in a LogReader to deserialize persisted containerVM output
		reader = iolog.NewLogReader(reader, ts)
	}

	detachableOut := NewFlushingReader(reader)

	return NewStreamOutputHandler("containerLogs").WithPayload(detachableOut, params.ID, nil)
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

	c := exec.Containers.Container(uid.Parse(params.ID).String())
	if c == nil {
		return containers.NewContainerWaitNotFound().WithPayload(&models.Error{
			Message: fmt.Sprintf("container %s not found", params.ID),
		})
	}

	select {
	case <-c.WaitForState(exec.StateStopped):
		containerInfo := convertContainerToContainerInfo(c.Info())

		return containers.NewContainerWaitOK().WithPayload(containerInfo)
	case <-ctx.Done():
		return containers.NewContainerWaitInternalServerError().WithPayload(&models.Error{
			Message: fmt.Sprintf("ContainerWaitHandler(%s) Error: %s", params.ID, ctx.Err()),
		})
	}
}

func (handler *ContainersHandlersImpl) RenameContainerHandler(params containers.ContainerRenameParams) middleware.Responder {
	defer trace.End(trace.Begin(fmt.Sprintf("Rename container to %s", params.Name)))

	h := exec.GetHandle(params.Handle)
	if h == nil || h.ExecConfig == nil {
		return containers.NewContainerRenameNotFound()
	}

	// get the indicated container for rename
	container := exec.Containers.Container(h.ExecConfig.ID)
	if container == nil {
		return containers.NewContainerRenameNotFound()
	}

	if container.ExecConfig.Name == params.Name {
		err := &models.Error{
			Message: fmt.Sprintf("renaming a container with the same name as its current name: %s", params.Name),
		}
		return containers.NewContainerRenameInternalServerError().WithPayload(err)
	}

	// rename on container version < supportVersionForRename is not supported
	log.Debugf("The container DataVersion is: %d", h.DataVersion)
	if h.DataVersion < feature.RenameSupportedVersion {
		err := &models.Error{
			Message: fmt.Sprintf("container %s does not support rename", container.ExecConfig.Name),
		}
		return containers.NewContainerRenameInternalServerError().WithPayload(err)
	}

	h = h.Rename(params.Name)

	return containers.NewContainerRenameOK().WithPayload(h.String())
}

// utility function to convert from a Container type to the API Model ContainerInfo (which should prob be called ContainerDetail)
func convertContainerToContainerInfo(container *exec.ContainerInfo) *models.ContainerInfo {
	defer trace.End(trace.Begin(container.ExecConfig.ID))
	// convert the container type to the required model
	info := &models.ContainerInfo{
		ContainerConfig: &models.ContainerConfig{},
		ProcessConfig:   &models.ProcessConfig{},
		VolumeConfig:    make([]*models.VolumeConfig, 0),
		Endpoints:       make([]*models.EndpointConfig, 0),
		DataVersion:     int64(container.DataVersion),
	}

	// Populate volume information
	for volName := range container.ExecConfig.Mounts {
		vol := &models.VolumeConfig{
			Name: volName,
		}
		info.VolumeConfig = append(info.VolumeConfig, vol)
	}

	ccid := container.ExecConfig.ID
	info.ContainerConfig.ContainerID = ccid

	var state string
	if container.MigrationError != nil {
		state = "error"
		info.ProcessConfig.ErrorMsg = fmt.Sprintf("Migration failed: %s", container.MigrationError.Error())
		info.ProcessConfig.Status = state
	} else {
		state = container.State().String()
	}
	info.ContainerConfig.State = state
	info.ContainerConfig.LayerID = container.ExecConfig.LayerID
	info.ContainerConfig.ImageID = container.ExecConfig.ImageID
	info.ContainerConfig.RepoName = &container.ExecConfig.RepoName
	info.ContainerConfig.CreateTime = container.ExecConfig.CreateTime
	info.ContainerConfig.Names = []string{container.ExecConfig.Name}
	info.ContainerConfig.RestartCount = int64(container.ExecConfig.Diagnostics.ResurrectionCount)
	info.ContainerConfig.StorageSize = container.VMUnsharedDisk

	if container.ExecConfig.Annotations != nil && len(container.ExecConfig.Annotations) > 0 {
		info.ContainerConfig.Annotations = make(map[string]string)

		for k, v := range container.ExecConfig.Annotations {
			info.ContainerConfig.Annotations[k] = v
		}
	}

	// in heavily loaded environments we were seeing a panic due to a missing
	// session id in execConfig -- this has only manifested itself in short lived containers
	// that were initilized via run
	if session, exists := container.ExecConfig.Sessions[ccid]; exists {
		info.ContainerConfig.Tty = &session.Tty
		info.ContainerConfig.AttachStdin = &session.Attach
		info.ContainerConfig.AttachStdout = &session.Attach
		info.ContainerConfig.AttachStderr = &session.Attach
		info.ContainerConfig.OpenStdin = &session.OpenStdin

		// started is a string in the vmx that is not to be confused
		// with started the datetime in the models.ContainerInfo
		info.ProcessConfig.Status = session.Started
		info.ProcessConfig.ExecPath = session.Cmd.Path
		info.ProcessConfig.WorkingDir = &session.Cmd.Dir
		info.ProcessConfig.ExecArgs = session.Cmd.Args
		info.ProcessConfig.Env = session.Cmd.Env
		info.ProcessConfig.ExitCode = int32(session.ExitStatus)
		info.ProcessConfig.StartTime = session.StartTime
		info.ProcessConfig.StopTime = session.StopTime

		info.ProcessConfig.User = session.User
		if session.Group != "" {
			info.ProcessConfig.User = fmt.Sprintf("%s:%s", session.User, session.Group)
		}
	} else {
		// log that sessionID is missing and print the ExecConfig
		log.Errorf("Session ID is missing from execConfig: %#v", container.ExecConfig)

		// panic if we are in debug / hopefully CI
		if log.DebugLevel > 0 {
			panic("nil session id")
		}

	}

	info.HostConfig = &models.HostConfig{}
	for _, endpoint := range container.ExecConfig.Networks {
		ep := &models.EndpointConfig{
			Address:     "",
			Container:   ccid,
			Gateway:     "",
			ID:          endpoint.ID,
			Name:        endpoint.Name,
			Ports:       make([]string, 0),
			Scope:       endpoint.Network.Name,
			Aliases:     make([]string, 0),
			Nameservers: make([]string, 0),
		}

		if len(endpoint.Network.Gateway.IP) > 0 {
			ep.Gateway = endpoint.Network.Gateway.String()
		}

		if len(endpoint.Assigned.IP) > 0 {
			ep.Address = endpoint.Assigned.String()
		}

		if len(endpoint.Ports) > 0 {
			ep.Ports = append(ep.Ports, endpoint.Ports...)
			info.HostConfig.Ports = append(info.HostConfig.Ports, endpoint.Ports...)
		}

		for _, alias := range endpoint.Network.Aliases {
			parts := strings.Split(alias, ":")
			if len(parts) > 1 {
				ep.Aliases = append(ep.Aliases, parts[1])
			} else {
				ep.Aliases = append(ep.Aliases, parts[0])
			}
		}

		for _, dns := range endpoint.Network.Nameservers {
			ep.Nameservers = append(ep.Nameservers, dns.String())
		}

		info.Endpoints = append(info.Endpoints, ep)
	}

	return info
}
