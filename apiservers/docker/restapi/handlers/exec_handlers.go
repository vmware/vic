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
