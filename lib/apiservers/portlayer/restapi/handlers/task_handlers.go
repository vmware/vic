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
	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/tasks"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/task"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
)

// TaskHandlersImpl is the receiver for all of the task handler methods
type TaskHandlersImpl struct {
}

func (handler *TaskHandlersImpl) Configure(api *operations.PortLayerAPI, _ *HandlerContext) {
	api.TasksJoinHandler = tasks.JoinHandlerFunc(handler.JoinHandler)
	api.TasksBindHandler = tasks.BindHandlerFunc(handler.BindHandler)
	api.TasksUnbindHandler = tasks.UnbindHandlerFunc(handler.UnbindHandler)
	api.TasksRemoveHandler = tasks.RemoveHandlerFunc(handler.RemoveHandler)
}

// JoinHandler calls the Join
func (handler *TaskHandlersImpl) JoinHandler(params tasks.JoinParams) middleware.Responder {
	defer trace.End(trace.Begin(""))
	op := trace.NewOperation(context.Background(), "task.Join(%s, %s)", params.Config.Handle, params.Config.ID)

	handle := exec.HandleFromInterface(params.Config.Handle)
	if handle == nil {
		err := &models.Error{Message: "Failed to get the Handle"}
		return tasks.NewJoinInternalServerError().WithPayload(err)
	}

	// TODO: ensure uniqueness of ID - this is already an issue with containercreate now we're not using it as
	// the VM name and cannot rely on vSphere for uniqueness guarantee
	id := params.Config.ID
	if id == "" {
		id = uid.New().String()
	}

	op.Debugf("ID: %#v", id)
	op.Debugf("Path: %#v", params.Config.Path)
	op.Debugf("WorkingDir: %#v", params.Config.WorkingDir)
	op.Debugf("OpenStdin: %#v", params.Config.OpenStdin)

	sessionConfig := &executor.SessionConfig{
		Common: executor.Common{
			ID: id,
		},
		Tty:       params.Config.Tty,
		Attach:    params.Config.Attach,
		OpenStdin: params.Config.OpenStdin,
		Cmd: executor.Cmd{
			Env:  params.Config.Env,
			Dir:  params.Config.WorkingDir,
			Path: params.Config.Path,
			Args: append([]string{params.Config.Path}, params.Config.Args...),
		},
		StopSignal: params.Config.StopSignal,
	}

	handleprime, err := task.Join(&op, handle, sessionConfig)
	if err != nil {
		log.Errorf("%s", err.Error())

		return tasks.NewJoinInternalServerError().WithPayload(
			&models.Error{Message: err.Error()},
		)
	}
	res := &models.TaskJoinResponse{
		ID:     id,
		Handle: exec.ReferenceFromHandle(handleprime),
	}
	return tasks.NewJoinOK().WithPayload(res)
}

// BindHandler calls the Bind
func (handler *TaskHandlersImpl) BindHandler(params tasks.BindParams) middleware.Responder {
	defer trace.End(trace.Begin(""))
	op := trace.NewOperation(context.Background(), "task.Bind(%s, %s)", params.Config.Handle, params.Config.ID)

	handle := exec.HandleFromInterface(params.Config.Handle)
	if handle == nil {
		err := &models.Error{Message: "Failed to get the Handle"}
		return tasks.NewBindInternalServerError().WithPayload(err)
	}

	handleprime, err := task.Bind(&op, handle, params.Config.ID)
	if err != nil {
		log.Errorf("%s", err.Error())

		return tasks.NewBindInternalServerError().WithPayload(
			&models.Error{Message: err.Error()},
		)
	}

	res := &models.TaskBindResponse{
		Handle: exec.ReferenceFromHandle(handleprime),
	}
	return tasks.NewBindOK().WithPayload(res)
}

// UnbindHandler calls the Unbind
func (handler *TaskHandlersImpl) UnbindHandler(params tasks.UnbindParams) middleware.Responder {
	defer trace.End(trace.Begin(""))
	op := trace.NewOperation(context.Background(), "task.Unbind(%s, %s)", params.Config.Handle, params.Config.ID)

	handle := exec.HandleFromInterface(params.Config.Handle)
	if handle == nil {
		err := &models.Error{Message: "Failed to get the Handle"}
		return tasks.NewUnbindInternalServerError().WithPayload(err)
	}

	handleprime, err := task.Unbind(&op, handle, params.Config.ID)
	if err != nil {
		log.Errorf("%s", err.Error())

		return tasks.NewUnbindInternalServerError().WithPayload(
			&models.Error{Message: err.Error()},
		)
	}

	res := &models.TaskUnbindResponse{
		Handle: exec.ReferenceFromHandle(handleprime),
	}
	return tasks.NewUnbindOK().WithPayload(res)
}

// RemoveHandler calls remove
func (handler *TaskHandlersImpl) RemoveHandler(params tasks.RemoveParams) middleware.Responder {
	defer trace.End(trace.Begin(""))
	op := trace.NewOperation(context.Background(), "task.Remove(%s, %s)", params.Config.Handle, params.Config.ID)

	handle := exec.HandleFromInterface(params.Config.Handle)
	if handle == nil {
		err := &models.Error{Message: "Failed to get the Handle"}
		return tasks.NewRemoveInternalServerError().WithPayload(err)
	}

	handleprime, err := task.Remove(&op, handle, params.Config.ID)
	if err != nil {
		log.Errorf("%s", err.Error())

		return tasks.NewRemoveInternalServerError().WithPayload(
			&models.Error{Message: err.Error()},
		)
	}

	res := &models.TaskRemoveResponse{
		Handle: exec.ReferenceFromHandle(handleprime),
	}
	return tasks.NewRemoveOK().WithPayload(res)
}
