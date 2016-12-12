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

	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/misc"
	"github.com/vmware/vic/lib/portlayer/exec"
)

// MiscHandlersImpl is the receiver for all the misc handler methods
type MiscHandlersImpl struct{}

// Configure assigns functions to all the miscellaneous api handlers
func (handler *MiscHandlersImpl) Configure(api *operations.PortLayerAPI, handlerCtx *HandlerContext) {
	api.MiscPingHandler = misc.PingHandlerFunc(handler.Ping)
	api.MiscGetVCHInfoHandler = misc.GetVCHInfoHandlerFunc(handler.GetVCHInfo)
}

// Ping sends an OK response to let the client know the server is up
func (handler *MiscHandlersImpl) Ping() middleware.Responder {
	return misc.NewPingOK().WithPayload("OK")
}

// GetVCHInfo returns VCH-related info for a `docker info` call
func (handler *MiscHandlersImpl) GetVCHInfo() middleware.Responder {
	ctx := context.Background()
	vchCPUMhz := exec.NCPU(ctx)
	vchMemLimit := exec.MemTotal(ctx)

	vchInfo := &models.VCHInfo{
		CPUMhz:          &vchCPUMhz,
		Memory:          &vchMemLimit,
		HostOS:          &exec.Config.HostOS,
		HostOSVersion:   &exec.Config.HostOSVersion,
		HostProductName: &exec.Config.HostProductName,
	}

	return misc.NewGetVCHInfoOK().WithPayload(vchInfo)
}
