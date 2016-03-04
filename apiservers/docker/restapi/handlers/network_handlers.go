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
	"github.com/vmware/vic/apiservers/docker/restapi/operations/network"
)

type NetworkHandlersImpl struct{}

func (handlers *NetworkHandlersImpl) Configure(api *operations.DockerAPI) {
	api.NetworkDeleteNetworksIDHandler = network.DeleteNetworksIDHandlerFunc(handlers.DeleteNetworks)
	api.NetworkGetNetworksHandler = network.GetNetworksHandlerFunc(handlers.GetNetworks)
	api.NetworkGetNetworksIDHandler = network.GetNetworksIDHandlerFunc(handlers.GetNetworksID)
	api.NetworkPostNetworksCreateHandler = network.PostNetworksCreateHandlerFunc(handlers.PostNetworksCreate)
	api.NetworkConnectHandler = network.ConnectHandlerFunc(handlers.Connect)
	api.NetworkDisconnectHandler = network.DisconnectHandlerFunc(handlers.Disconnect)
}

func (handlers *NetworkHandlersImpl) DeleteNetworks(params network.DeleteNetworksIDParams) middleware.Responder {
	return middleware.NotImplemented("operation network.DeleteNetworksID has not yet been implemented")
}

func (handlers *NetworkHandlersImpl) GetNetworks(params network.GetNetworksParams) middleware.Responder {
	return middleware.NotImplemented("operation network.GetNetworks has not yet been implemented")
}

func (handlers *NetworkHandlersImpl) GetNetworksID(params network.GetNetworksIDParams) middleware.Responder {
	return middleware.NotImplemented("operation network.GetNetworksID has not yet been implemented")
}

func (handlers *NetworkHandlersImpl) PostNetworksCreate(params network.PostNetworksCreateParams) middleware.Responder {
	return middleware.NotImplemented("operation network.PostNetworksCreate has not yet been implemented")
}

func (handlers *NetworkHandlersImpl) Connect(params network.ConnectParams) middleware.Responder {
	return middleware.NotImplemented("operation network.Connect has not yet been implemented")
}

func (handlers *NetworkHandlersImpl) Disconnect(params network.DisconnectParams) middleware.Responder {
	return middleware.NotImplemented("operation network.Disconnect has not yet been implemented")
}
