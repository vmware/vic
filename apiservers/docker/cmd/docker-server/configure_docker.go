package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/apiservers/docker/restapi/operations"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/container"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/exec"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/image"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/misc"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/network"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/volume"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.DockerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.TxtConsumer = httpkit.ConsumerFunc(func(r io.Reader, target interface{}) error {
		fmt.Print("Inside httpkit.ConsumerFunc")
		return nil
//		return errors.NotImplemented("txt consumer has not yet been implemented")
	})

	api.JSONProducer = httpkit.JSONProducer()

	api.ImageDeleteImagesNameHandler = image.DeleteImagesNameHandlerFunc(func(params image.DeleteImagesNameParams) middleware.Responder {
		return middleware.NotImplemented("operation image.DeleteImagesName has not yet been implemented")
	})
	api.NetworkDeleteNetworksIDHandler = network.DeleteNetworksIDHandlerFunc(func(params network.DeleteNetworksIDParams) middleware.Responder {
		return middleware.NotImplemented("operation network.DeleteNetworksID has not yet been implemented")
	})
	api.VolumeDeleteVolumesNameHandler = volume.DeleteVolumesNameHandlerFunc(func(params volume.DeleteVolumesNameParams) middleware.Responder {
		return middleware.NotImplemented("operation volume.DeleteVolumesName has not yet been implemented")
	})
	api.ImageGetImagesJSONHandler = image.GetImagesJSONHandlerFunc(func(params image.GetImagesJSONParams) middleware.Responder {
		return middleware.NotImplemented("operation image.GetImagesJSON has not yet been implemented")
	})
	api.ImageGetImagesNameJSONHandler = image.GetImagesNameJSONHandlerFunc(func(params image.GetImagesNameJSONParams) middleware.Responder {
		return middleware.NotImplemented("operation image.GetImagesNameJSON has not yet been implemented")
	})
	api.NetworkGetNetworksHandler = network.GetNetworksHandlerFunc(func(params network.GetNetworksParams) middleware.Responder {
		return middleware.NotImplemented("operation network.GetNetworks has not yet been implemented")
	})
	api.NetworkGetNetworksIDHandler = network.GetNetworksIDHandlerFunc(func(params network.GetNetworksIDParams) middleware.Responder {
		return middleware.NotImplemented("operation network.GetNetworksID has not yet been implemented")
	})
	api.VolumeGetVolumesHandler = volume.GetVolumesHandlerFunc(func(params volume.GetVolumesParams) middleware.Responder {
		return middleware.NotImplemented("operation volume.GetVolumes has not yet been implemented")
	})
	api.VolumeGetVolumesNameHandler = volume.GetVolumesNameHandlerFunc(func(params volume.GetVolumesNameParams) middleware.Responder {
		return middleware.NotImplemented("operation volume.GetVolumesName has not yet been implemented")
	})
	api.ExecPostContainersIDExecHandler = exec.PostContainersIDExecHandlerFunc(func(params exec.PostContainersIDExecParams) middleware.Responder {
		return middleware.NotImplemented("operation exec.PostContainersIDExec has not yet been implemented")
	})
	api.ExecPostExecIDJSONHandler = exec.PostExecIDJSONHandlerFunc(func(params exec.PostExecIDJSONParams) middleware.Responder {
		return middleware.NotImplemented("operation exec.PostExecIDJSON has not yet been implemented")
	})
	api.ExecPostExecIDResizeHandler = exec.PostExecIDResizeHandlerFunc(func(params exec.PostExecIDResizeParams) middleware.Responder {
		return middleware.NotImplemented("operation exec.PostExecIDResize has not yet been implemented")
	})
	api.ExecPostExecIDStartHandler = exec.PostExecIDStartHandlerFunc(func(params exec.PostExecIDStartParams) middleware.Responder {
		return middleware.NotImplemented("operation exec.PostExecIDStart has not yet been implemented")
	})
	api.ImagePostImagesCreateHandler = image.PostImagesCreateHandlerFunc(func(params image.PostImagesCreateParams) middleware.Responder {
		return middleware.NotImplemented("operation image.PostImagesCreate has not yet been implemented")
	})
	api.NetworkPostNetworksCreateHandler = network.PostNetworksCreateHandlerFunc(func(params network.PostNetworksCreateParams) middleware.Responder {
		return middleware.NotImplemented("operation network.PostNetworksCreate has not yet been implemented")
	})
	api.VolumePostVolumesCreateHandler = volume.PostVolumesCreateHandlerFunc(func(params volume.PostVolumesCreateParams) middleware.Responder {
		return middleware.NotImplemented("operation volume.PostVolumesCreate has not yet been implemented")
	})
	api.ContainerAttachHandler = container.AttachHandlerFunc(func(params container.AttachParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Attach has not yet been implemented")
	})
	api.ContainerAttachWebsocketHandler = container.AttachWebsocketHandlerFunc(func(params container.AttachWebsocketParams) middleware.Responder {
		return middleware.NotImplemented("operation container.AttachWebsocket has not yet been implemented")
	})
	api.ImageBuildHandler = image.BuildHandlerFunc(func(params image.BuildParams) middleware.Responder {
		return middleware.NotImplemented("operation image.Build has not yet been implemented")
	})
	api.ContainerChangesHandler = container.ChangesHandlerFunc(func(params container.ChangesParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Changes has not yet been implemented")
	})
	api.MiscCheckAuthenticationHandler = misc.CheckAuthenticationHandlerFunc(func(params misc.CheckAuthenticationParams) middleware.Responder {
		return middleware.NotImplemented("operation misc.CheckAuthentication has not yet been implemented")
	})
	api.ImageCommitHandler = image.CommitHandlerFunc(func(params image.CommitParams) middleware.Responder {
		return middleware.NotImplemented("operation image.Commit has not yet been implemented")
	})
	api.NetworkConnectHandler = network.ConnectHandlerFunc(func(params network.ConnectParams) middleware.Responder {
		return middleware.NotImplemented("operation network.Connect has not yet been implemented")
	})
	api.ContainerCreateHandler = container.CreateHandlerFunc(func(params container.CreateParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Create has not yet been implemented")
	})
	api.NetworkDisconnectHandler = network.DisconnectHandlerFunc(func(params network.DisconnectParams) middleware.Responder {
		return middleware.NotImplemented("operation network.Disconnect has not yet been implemented")
	})
	api.ContainerExportHandler = container.ExportHandlerFunc(func(params container.ExportParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Export has not yet been implemented")
	})
	api.ContainerFindHandler = container.FindHandlerFunc(func(params container.FindParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Find has not yet been implemented")
	})
	api.ContainerFindAllHandler = container.FindAllHandlerFunc(func(params container.FindAllParams) middleware.Responder {
		return middleware.NotImplemented("operation container.FindAll has not yet been implemented")
	})
	api.ContainerGetArchiveHandler = container.GetArchiveHandlerFunc(func(params container.GetArchiveParams) middleware.Responder {
		return middleware.NotImplemented("operation container.GetArchive has not yet been implemented")
	})
	api.ContainerGetArchiveInformationHandler = container.GetArchiveInformationHandlerFunc(func(params container.GetArchiveInformationParams) middleware.Responder {
		return middleware.NotImplemented("operation container.GetArchiveInformation has not yet been implemented")
	})
	api.MiscGetEventsHandler = misc.GetEventsHandlerFunc(func(params misc.GetEventsParams) middleware.Responder {
		return middleware.NotImplemented("operation misc.GetEvents has not yet been implemented")
	})
	api.MiscGetSystemInformationHandler = misc.GetSystemInformationHandlerFunc(func() middleware.Responder {
		return middleware.NotImplemented("operation misc.GetSystemInformation has not yet been implemented")
	})
	api.MiscGetVersionHandler = misc.GetVersionHandlerFunc(func() middleware.Responder {
		return middleware.NotImplemented("operation misc.GetVersion has not yet been implemented")
	})
	api.ImageHistoryHandler = image.HistoryHandlerFunc(func(params image.HistoryParams) middleware.Responder {
		return middleware.NotImplemented("operation image.History has not yet been implemented")
	})
	api.ContainerKillHandler = container.KillHandlerFunc(func(params container.KillParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Kill has not yet been implemented")
	})
	api.ContainerListProcessesHandler = container.ListProcessesHandlerFunc(func(params container.ListProcessesParams) middleware.Responder {
		return middleware.NotImplemented("operation container.ListProcesses has not yet been implemented")
	})
	api.ImageLoadHandler = image.LoadHandlerFunc(func(params image.LoadParams) middleware.Responder {
		return middleware.NotImplemented("operation image.Load has not yet been implemented")
	})
	api.ContainerLogsHandler = container.LogsHandlerFunc(func(params container.LogsParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Logs has not yet been implemented")
	})
	api.ContainerPauseHandler = container.PauseHandlerFunc(func(params container.PauseParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Pause has not yet been implemented")
	})
	api.MiscPingHandler = misc.PingHandlerFunc(func() middleware.Responder {
		return middleware.NotImplemented("operation misc.Ping has not yet been implemented")
	})
	api.ImagePushHandler = image.PushHandlerFunc(func(params image.PushParams) middleware.Responder {
		return middleware.NotImplemented("operation image.Push has not yet been implemented")
	})
	api.ContainerPutArchiveHandler = container.PutArchiveHandlerFunc(func(params container.PutArchiveParams) middleware.Responder {
		return middleware.NotImplemented("operation container.PutArchive has not yet been implemented")
	})
	api.ContainerRemoveHandler = container.RemoveHandlerFunc(func(params container.RemoveParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Remove has not yet been implemented")
	})
	api.ContainerRenameHandler = container.RenameHandlerFunc(func(params container.RenameParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Rename has not yet been implemented")
	})
	api.ContainerResizeHandler = container.ResizeHandlerFunc(func(params container.ResizeParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Resize has not yet been implemented")
	})
	api.ContainerRestartHandler = container.RestartHandlerFunc(func(params container.RestartParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Restart has not yet been implemented")
	})
	api.ImageSaveHandler = image.SaveHandlerFunc(func(params image.SaveParams) middleware.Responder {
		return middleware.NotImplemented("operation image.Save has not yet been implemented")
	})
	api.ImageSaveAllHandler = image.SaveAllHandlerFunc(func(params image.SaveAllParams) middleware.Responder {
		return middleware.NotImplemented("operation image.SaveAll has not yet been implemented")
	})
	api.ImageSearchHandler = image.SearchHandlerFunc(func(params image.SearchParams) middleware.Responder {
		return middleware.NotImplemented("operation image.Search has not yet been implemented")
	})
	api.ContainerStartHandler = container.StartHandlerFunc(func(params container.StartParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Start has not yet been implemented")
	})
	api.ContainerStatsHandler = container.StatsHandlerFunc(func(params container.StatsParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Stats has not yet been implemented")
	})
	api.ContainerStopHandler = container.StopHandlerFunc(func(params container.StopParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Stop has not yet been implemented")
	})
	api.ImageTagHandler = image.TagHandlerFunc(func(params image.TagParams) middleware.Responder {
		return middleware.NotImplemented("operation image.Tag has not yet been implemented")
	})
	api.ContainerUnpauseHandler = container.UnpauseHandlerFunc(func(params container.UnpauseParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Unpause has not yet been implemented")
	})
	api.ContainerWaitHandler = container.WaitHandlerFunc(func(params container.WaitParams) middleware.Responder {
		return middleware.NotImplemented("operation container.Wait has not yet been implemented")
	})

	api.ServerShutdown = func() {}
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
