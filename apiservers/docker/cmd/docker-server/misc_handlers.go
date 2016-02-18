package main

import (
	//	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/apiservers/docker/restapi/operations"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/misc"
)

type MiscHandlersImpl struct{}

func (handlers *MiscHandlersImpl) Configure(api *operations.DockerAPI) {
	api.MiscCheckAuthenticationHandler = misc.CheckAuthenticationHandlerFunc(handlers.CheckAuthentication)
	api.MiscGetEventsHandler = misc.GetEventsHandlerFunc(handlers.GetEvents)
	api.MiscGetSystemInformationHandler = misc.GetSystemInformationHandlerFunc(handlers.GetSystemInfo)
	api.MiscGetVersionHandler = misc.GetVersionHandlerFunc(handlers.GetVersion)
	api.MiscPingHandler = misc.PingHandlerFunc(handlers.Ping)
}

func (handlers *MiscHandlersImpl) CheckAuthentication(params misc.CheckAuthenticationParams) middleware.Responder {
	return middleware.NotImplemented("operation misc.CheckAuthentication has not yet been implemented")
}

func (handlers *MiscHandlersImpl) GetEvents(params misc.GetEventsParams) middleware.Responder {
	return middleware.NotImplemented("operation misc.GetEvents has not yet been implemented")
}

func (handlers *MiscHandlersImpl) GetSystemInfo() middleware.Responder {
	return middleware.NotImplemented("operation misc.GetSystemInformation has not yet been implemented")
}

func (handlers *MiscHandlersImpl) GetVersion() middleware.Responder {
	return middleware.NotImplemented("operation misc.GetVersion has not yet been implemented")
}

func (handlers *MiscHandlersImpl) Ping() middleware.Responder {
	return misc.NewPingOK().WithPayload("OK")
}
