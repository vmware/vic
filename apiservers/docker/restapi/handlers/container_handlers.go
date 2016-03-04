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
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/apiservers/docker/restapi/operations"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/container"
)

type ContainerHandlersImpl struct{}

func (handlers *ContainerHandlersImpl) Configure(api *operations.DockerAPI) {
	api.ContainerAttachHandler = container.AttachHandlerFunc(handlers.Attach)
	api.ContainerAttachWebsocketHandler = container.AttachWebsocketHandlerFunc(handlers.AttachWebsocket)
	api.ContainerChangesHandler = container.ChangesHandlerFunc(handlers.Changes)
	api.ContainerCreateHandler = container.CreateHandlerFunc(handlers.Create)
	api.ContainerExportHandler = container.ExportHandlerFunc(handlers.Export)
	api.ContainerFindHandler = container.FindHandlerFunc(handlers.Find)
	api.ContainerFindAllHandler = container.FindAllHandlerFunc(handlers.FindAll)
	api.ContainerGetArchiveHandler = container.GetArchiveHandlerFunc(handlers.GetArchive)
	api.ContainerGetArchiveInformationHandler = container.GetArchiveInformationHandlerFunc(handlers.GetArchiveInformation)
	api.ContainerKillHandler = container.KillHandlerFunc(handlers.Kill)
	api.ContainerListProcessesHandler = container.ListProcessesHandlerFunc(handlers.ListProcesses)
	api.ContainerLogsHandler = container.LogsHandlerFunc(handlers.Logs)
	api.ContainerPauseHandler = container.PauseHandlerFunc(handlers.Pause)
	api.ContainerPutArchiveHandler = container.PutArchiveHandlerFunc(handlers.PutArchive)
	api.ContainerRemoveHandler = container.RemoveHandlerFunc(handlers.Remove)
	api.ContainerRenameHandler = container.RenameHandlerFunc(handlers.Rename)
	api.ContainerResizeHandler = container.ResizeHandlerFunc(handlers.Resize)
	api.ContainerRestartHandler = container.RestartHandlerFunc(handlers.Restart)
	api.ContainerStartHandler = container.StartHandlerFunc(handlers.Start)
	api.ContainerStatsHandler = container.StatsHandlerFunc(handlers.Stats)
	api.ContainerStopHandler = container.StopHandlerFunc(handlers.Stop)
	api.ContainerUnpauseHandler = container.UnpauseHandlerFunc(handlers.Unpause)
	api.ContainerWaitHandler = container.WaitHandlerFunc(handlers.Wait)

}

func (handlers *ContainerHandlersImpl) Attach(params container.AttachParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Attach has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) AttachWebsocket(params container.AttachWebsocketParams) middleware.Responder {
	return middleware.NotImplemented("operation container.AttachWebsocket has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Changes(params container.ChangesParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Changes has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Create(params container.CreateParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Create has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Export(params container.ExportParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Export has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Find(params container.FindParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Find has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) FindAll(params container.FindAllParams) middleware.Responder {
	return middleware.NotImplemented("operation container.FindAll has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) GetArchive(params container.GetArchiveParams) middleware.Responder {
	return middleware.NotImplemented("operation container.GetArchive has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) GetArchiveInformation(params container.GetArchiveInformationParams) middleware.Responder {
	return middleware.NotImplemented("operation container.GetArchiveInformation has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Kill(params container.KillParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Kill has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) ListProcesses(params container.ListProcessesParams) middleware.Responder {
	return middleware.NotImplemented("operation container.ListProcesses has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Logs(params container.LogsParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Logs has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Pause(params container.PauseParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Pause has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) PutArchive(params container.PutArchiveParams) middleware.Responder {
	return middleware.NotImplemented("operation container.PutArchive has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Remove(params container.RemoveParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Remove has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Rename(params container.RenameParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Rename has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Resize(params container.ResizeParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Resize has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Restart(params container.RestartParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Restart has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Start(params container.StartParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Start has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Stats(params container.StatsParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Stats has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Stop(params container.StopParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Stop has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Unpause(params container.UnpauseParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Unpause has not yet been implemented")
}

func (handlers *ContainerHandlersImpl) Wait(params container.WaitParams) middleware.Responder {
	return middleware.NotImplemented("operation container.Wait has not yet been implemented")
}
