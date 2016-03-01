package handlers

import (
	"net/http"
	"os"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/swag"

	"github.com/vmware/vic/apiservers/portlayer/models"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations/storage"

	linux "github.com/vmware/vic/portlayer/linux/storage"
	portlayer "github.com/vmware/vic/portlayer/storage"
	"github.com/vmware/vic/portlayer/util"
)

// StorageHandlersImpl
type StorageHandlersImpl struct{}

var ls = &linux.LocalStore{
	Path: "/var/lib/portlayer",
}

var cache = &portlayer.NameLookupCache{
	DataStore: ls,
}

// Configure assigns functions to all the storage api handlers
func (handler *StorageHandlersImpl) Configure(api *operations.PortLayerAPI) {
	api.StorageCreateImageStoreHandler = storage.CreateImageStoreHandlerFunc(handler.CreateImageStore)
	api.StorageGetImageHandler = storage.GetImageHandlerFunc(handler.GetImage)
	api.StorageGetImageTarHandler = storage.GetImageTarHandlerFunc(handler.GetImageTar)
	api.StorageListImagesHandler = storage.ListImagesHandlerFunc(handler.ListImages)
	api.StorageWriteImageHandler = storage.WriteImageHandlerFunc(handler.WriteImage)
}

// CreateImageStore creates a new image store
func (handler *StorageHandlersImpl) CreateImageStore(params storage.CreateImageStoreParams) middleware.Responder {
	url, err := cache.CreateImageStore(params.Body.Name)
	if err != nil {
		if os.IsExist(err) {
			return storage.NewCreateImageStoreConflict().WithPayload(
				&models.Error{
					Code:    swag.Int64(http.StatusConflict),
					Message: "An image store with that name already exists",
				})
		} else {
			return storage.NewCreateImageStoreDefault(http.StatusInternalServerError).WithPayload(
				&models.Error{
					Code:    swag.Int64(http.StatusInternalServerError),
					Message: err.Error(),
				})
		}
	}
	s := &models.StoreURL{Code: swag.Int64(http.StatusCreated), URL: url.String()}
	return storage.NewCreateImageStoreCreated().WithPayload(s)
}

// GetImage retrieves an image from a store
func (handler *StorageHandlersImpl) GetImage(params storage.GetImageParams) middleware.Responder {
	id := params.ID
	url, err := util.StoreNameToURL(params.StoreName)
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
}

// GetImageTar returns an image tar file
func (handler *StorageHandlersImpl) GetImageTar(params storage.GetImageTarParams) middleware.Responder {
	return middleware.NotImplemented("operation storage.GetImageTar has not yet been implemented")
}

// ListImages returns a list of images in a store
func (handler *StorageHandlersImpl) ListImages(params storage.ListImagesParams) middleware.Responder {
	// TODO(jzt): support multiple query args i.e.:  /storage/ListImages?id=1,2,3,4...
	u, err := util.StoreNameToURL(params.StoreName)
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
}

// WriteImage writes an image to an image store
func (handler *StorageHandlersImpl) WriteImage(params storage.WriteImageParams) middleware.Responder {
	u, err := util.StoreNameToURL(params.StoreName)
	if err != nil {
		return storage.NewWriteImageDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    swag.Int64(http.StatusInternalServerError),
				Message: err.Error(),
			})
	}

	parent := &portlayer.Image{
		Store: u,
		ID:    params.ParentID,
	}

	image, err := cache.WriteImage(parent, params.ImageID, params.ImageFile)
	if err != nil {
		return storage.NewWriteImageDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    swag.Int64(http.StatusInternalServerError),
				Message: err.Error(),
			})
	}
	i := convertImage(image)
	return storage.NewWriteImageCreated().WithPayload(i)
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
