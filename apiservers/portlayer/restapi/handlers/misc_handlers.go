package handlers

import (
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations/misc"
)

// MiscHandlersImpl
type MiscHandlersImpl struct{}

// Configure assigns functions to all the miscellaneous api handlers
func (handler *MiscHandlersImpl) Configure(api *operations.PortLayerAPI) {
	api.MiscPingHandler = misc.PingHandlerFunc(handler.Ping)
}

// Ping sends an OK response to let the client know the server is up
func (handler *MiscHandlersImpl) Ping() middleware.Responder {
	return misc.NewPingOK().WithPayload("OK")
}
