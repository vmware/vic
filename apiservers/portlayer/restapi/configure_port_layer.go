package restapi

import (
	//"io"
	"net/http"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"

	"github.com/vmware/vic/apiservers/portlayer/restapi/handlers"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

type portlayerhandlers struct {
	storageHandlers handlers.StorageHandlersImpl
	miscHandlers    handlers.MiscHandlersImpl
}

func configureAPI(api *operations.PortLayerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.BinConsumer = httpkit.ByteStreamConsumer()

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	allhandlers := portlayerhandlers{}

	allhandlers.storageHandlers.Configure(api)
	allhandlers.miscHandlers.Configure(api)

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
