package restapi

import (
	"io"
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/vmware/vic/apiservers/docker/restapi/handlers"
	"github.com/vmware/vic/apiservers/docker/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

type dockerhandlers struct {
	imageHandlers     handlers.ImageHandlersImpl
	containerHandlers handlers.ContainerHandlersImpl
	networkHandlers   handlers.NetworkHandlersImpl
	volumeHandlers    handlers.VolumeHandlersImpl
	execHandlers      handlers.ExecHandlersImpl
	miscHandlers      handlers.MiscHandlersImpl
}

func configureAPI(api *operations.DockerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.TxtConsumer = httpkit.ConsumerFunc(func(r io.Reader, target interface{}) error {
		return nil
	})

	api.JSONProducer = httpkit.JSONProducer()
	api.BinConsumer = httpkit.ConsumerFunc(func(r io.Reader, target interface{}) error {
		return nil
	})

	allhandlers := dockerhandlers{}

	allhandlers.containerHandlers.Configure(api)
	allhandlers.imageHandlers.Configure(api)
	allhandlers.networkHandlers.Configure(api)
	allhandlers.volumeHandlers.Configure(api)
	allhandlers.execHandlers.Configure(api)
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
