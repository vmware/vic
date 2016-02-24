package handlers

import (
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/apiservers/docker/restapi/operations"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/exec"
)

type ExecHandlersImpl struct{}

func (handlers *ExecHandlersImpl) Configure(api *operations.DockerAPI) {
	api.ExecPostContainersIDExecHandler = exec.PostContainersIDExecHandlerFunc(handlers.PostContainersIDExec)
	api.ExecPostExecIDJSONHandler = exec.PostExecIDJSONHandlerFunc(handlers.PostExecIDJSON)
	api.ExecPostExecIDResizeHandler = exec.PostExecIDResizeHandlerFunc(handlers.PostExecIDResize)
	api.ExecPostExecIDStartHandler = exec.PostExecIDStartHandlerFunc(handlers.PostExecIDStart)
}

func (handlers *ExecHandlersImpl) PostContainersIDExec(params exec.PostContainersIDExecParams) middleware.Responder {
	return middleware.NotImplemented("operation exec.PostContainersIDExec has not yet been implemented")
}

func (handlers *ExecHandlersImpl) PostExecIDJSON(params exec.PostExecIDJSONParams) middleware.Responder {
	return middleware.NotImplemented("operation exec.PostExecIDJSON has not yet been implemented")
}

func (handlers *ExecHandlersImpl) PostExecIDResize(params exec.PostExecIDResizeParams) middleware.Responder {
	return middleware.NotImplemented("operation exec.PostExecIDResize has not yet been implemented")
}

func (handlers *ExecHandlersImpl) PostExecIDStart(params exec.PostExecIDStartParams) middleware.Responder {
	return middleware.NotImplemented("operation exec.PostExecIDStart has not yet been implemented")
}
