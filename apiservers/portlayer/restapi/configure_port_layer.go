package restapi

import (
	"net/http"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/go-swagger/go-swagger/swag"

	"github.com/vmware/vic/apiservers/portlayer/models"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations/storage"

	linux "github.com/vmware/vic/portlayer/linux/storage"
	portlayer "github.com/vmware/vic/portlayer/storage"
	"github.com/vmware/vic/portlayer/util"
)

// This file is safe to edit. Once it exists it will not be overwritten

var ls = &linux.LocalStore{
	Path: "/var/lib/portlayer",
}

var cache = &portlayer.NameLookupCache{
	DataStore: ls,
}

// convert an SPL Image to a swagger-defined Image
func convertImage(image *portlayer.Image) *models.Image {
	var parent *string

	// scratch image
	if image.Parent != nil {
		s := image.Parent.String()
		parent = &s
	}

	s := image.SelfLink.String()
	selflink := &s

	return &models.Image{
		ID:       image.ID,
		SelfLink: selflink,
		Parent:   parent,
		Store:    image.Store.String(),
	}
}

func configureAPI(api *operations.PortLayerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.StorageCreateImageStoreHandler = storage.CreateImageStoreHandlerFunc(func(params storage.CreateImageStoreParams) middleware.Responder {
		url, err := cache.CreateImageStore(params.Body.Name)
		if err != nil {
			return storage.NewCreateImageStoreDefault(http.StatusInternalServerError).WithPayload(
				&models.Error{
					Code:    swag.Int64(http.StatusInternalServerError),
					Message: err.Error(),
				})
		}
		s := &models.StoreURL{Code: swag.Int64(http.StatusCreated), URL: url.String()}
		return storage.NewCreateImageStoreCreated().WithPayload(s)
	})

	api.StorageGetImageHandler = storage.GetImageHandlerFunc(func(params storage.GetImageParams) middleware.Responder {
		id := params.ID
		url, err := util.StoreNameToUrl(params.StoreName)
		if err != nil {
			return storage.NewGetImageDefault(http.StatusInternalServerError).WithPayload(
				&models.Error{
					Code:    swag.Int64(http.StatusInternalServerError),
					Message: err.Error(),
				})
		}

		image, err := cache.GetImage(url, id)
		if err != nil {
			e := &models.Error{Code: swag.Int64(http.StatusNotFound), Message: err.Error()}
			return storage.NewGetImageDefault(http.StatusNotFound).WithPayload(e)
		}
		result := convertImage(image)
		return storage.NewGetImageOK().WithPayload(result)
	})

	api.StorageGetImageTarHandler = storage.GetImageTarHandlerFunc(func(params storage.GetImageTarParams) middleware.Responder {
		return middleware.NotImplemented("operation storage.GetImageTar has not yet been implemented")
	})

	api.StorageListImagesHandler = storage.ListImagesHandlerFunc(func(params storage.ListImagesParams) middleware.Responder {
		// TODO(jzt): support multiple query args i.e.:  /storage/ListImages?id=1,2,3,4...
		u, err := util.StoreNameToUrl(params.StoreName)
		if err != nil {
			return storage.NewListImagesDefault(http.StatusInternalServerError).WithPayload(
				&models.Error{
					Code:    swag.Int64(http.StatusInternalServerError),
					Message: err.Error(),
				})
		}

		// FIXME(jzt): not populating the cache at startup will result in 404's
		images, err := cache.ListImages(u, nil)
		if err != nil {
			return storage.NewListImagesDefault(http.StatusNotFound).WithPayload(
				&models.Error{
					Code:    swag.Int64(http.StatusNotFound),
					Message: err.Error(),
				})
		}

		result := make([]*models.Image, 0, len(images))

		for _, image := range images {
			result = append(result, convertImage(image))
		}
		return storage.NewListImagesOK().WithPayload(result)
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
