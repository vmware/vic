package restapi

import (
	"net/http"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations/storage"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.PortLayerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.StorageCreateImageStoreHandler = storage.CreateImageStoreHandlerFunc(func(params storage.CreateImageStoreParams) middleware.Responder {
		return middleware.NotImplemented("operation storage.CreateImageStore has not yet been implemented")
	})
	api.StorageGetImageInfoHandler = storage.GetImageInfoHandlerFunc(func(params storage.GetImageInfoParams) middleware.Responder {
		return middleware.NotImplemented("operation storage.GetImageInfo has not yet been implemented")
	})
	api.StorageGetImageTarHandler = storage.GetImageTarHandlerFunc(func(params storage.GetImageTarParams) middleware.Responder {
		return middleware.NotImplemented("operation storage.GetImageTar has not yet been implemented")
	})
	api.StorageListImageStoresHandler = storage.ListImageStoresHandlerFunc(func() middleware.Responder {
		return middleware.NotImplemented("operation storage.ListImageStores has not yet been implemented")
	})
	api.StorageListImagesHandler = storage.ListImagesHandlerFunc(func(params storage.ListImagesParams) middleware.Responder {
		return middleware.NotImplemented("operation storage.ListImages has not yet been implemented")
	})
	api.StorageWriteImageHandler = storage.WriteImageHandlerFunc(func(params storage.WriteImageParams) middleware.Responder {
		return middleware.NotImplemented("operation storage.WriteImage has not yet been implemented")
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
