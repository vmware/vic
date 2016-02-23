package main

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
