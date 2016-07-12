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
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/swag"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/storage"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/options"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"

	epl "github.com/vmware/vic/lib/portlayer/exec"
	spl "github.com/vmware/vic/lib/portlayer/storage"
	vsphereSpl "github.com/vmware/vic/lib/portlayer/storage/vsphere"
	"github.com/vmware/vic/lib/portlayer/util"

	"golang.org/x/net/context"
)

// StorageHandlersImpl is the receiver for all of the storage handler methods
type StorageHandlersImpl struct{}

var (
	storageSession     = &session.Session{}
	storageImageLayer  = &spl.NameLookupCache{}
	storageVolumeLayer = &spl.VolumeLookupCache{}
)

// Configure assigns functions to all the storage api handlers
func (handler *StorageHandlersImpl) Configure(api *operations.PortLayerAPI, handlerCtx *HandlerContext) {
	var err error

	ctx := context.Background()

	sessionconfig := &session.Config{
		Service:        options.PortLayerOptions.SDK,
		Insecure:       options.PortLayerOptions.Insecure,
		Keepalive:      options.PortLayerOptions.Keepalive,
		DatacenterPath: options.PortLayerOptions.DatacenterPath,
		ClusterPath:    options.PortLayerOptions.ClusterPath,
		PoolPath:       options.PortLayerOptions.PoolPath,
		DatastorePath:  options.PortLayerOptions.DatastorePath,
	}

	storageSession, err = session.NewSession(sessionconfig).Create(ctx)
	if err != nil {
		log.Fatalf("StorageHandler ERROR: %s", err)
	}

	ds, err := vsphereSpl.NewImageStore(ctx, storageSession)
	if err != nil {
		log.Panicf("Cannot instantiate storage layer: %s", err)
	}

	// The imagestore is implemented via a cache which is backed via an
	// implementation that writes to disks.  The cache is used to avoid
	// expensive metadata lookups.
	storageImageLayer = spl.NewLookupCache(ds)
	//FIXME: this may need another viewing after ian/faiyaz's changes
	vsVolumeStore, err := vsphereSpl.NewVolumeStore(context.TODO(), storageSession)
	if err != nil {
		log.Panicf("Cannot instantiate the volume store: %s", err)
	}

	// Get the datastores for volumes.
	// Each volume store name maps to a datastore + path, which can be referred to by the name.
	dstores, err := vsphereSpl.GetDatastores(context.TODO(), storageSession, spl.Config.VolumeLocations)
	if err != nil {
		log.Panicf("Cannot find datastores: %s", err)
	}

	// Add datastores to the vsphere volume store impl
	for volStoreName, volDatastore := range dstores {
		log.Infof("Adding volume store %s (%s)", volStoreName, volDatastore.RootURL)
		_, err := vsVolumeStore.AddStore(context.TODO(), volDatastore, volStoreName)
		if err != nil {
			log.Errorf("volume addition error %s", err)
		}
	}

	storageVolumeLayer, err = spl.NewVolumeLookupCache(context.TODO(), vsVolumeStore)
	if err != nil {
		log.Panicf("Cannot instantiate the Volume Lookup cache: %s", err)
	}

	api.StorageCreateImageStoreHandler = storage.CreateImageStoreHandlerFunc(handler.CreateImageStore)
	api.StorageGetImageHandler = storage.GetImageHandlerFunc(handler.GetImage)
	api.StorageGetImageTarHandler = storage.GetImageTarHandlerFunc(handler.GetImageTar)
	api.StorageListImagesHandler = storage.ListImagesHandlerFunc(handler.ListImages)
	api.StorageWriteImageHandler = storage.WriteImageHandlerFunc(handler.WriteImage)
	api.StorageRemoveVolumeHandler = storage.RemoveVolumeHandlerFunc(handler.RemoveVolume)
	api.StorageCreateVolumeHandler = storage.CreateVolumeHandlerFunc(handler.CreateVolume)
	api.StorageVolumeJoinHandler = storage.VolumeJoinHandlerFunc(handler.VolumeJoin)
}

// CreateImageStore creates a new image store
func (handler *StorageHandlersImpl) CreateImageStore(params storage.CreateImageStoreParams) middleware.Responder {
	url, err := storageImageLayer.CreateImageStore(context.TODO(), params.Body.Name)
	if err != nil {
		if os.IsExist(err) {
			return storage.NewCreateImageStoreConflict().WithPayload(
				&models.Error{
					Code:    swag.Int64(http.StatusConflict),
					Message: "An image store with that name already exists",
				})
		}

		return storage.NewCreateImageStoreDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    swag.Int64(http.StatusInternalServerError),
				Message: err.Error(),
			})
	}
	s := &models.StoreURL{Code: swag.Int64(http.StatusCreated), URL: url.String()}
	return storage.NewCreateImageStoreCreated().WithPayload(s)
}

// GetImage retrieves an image from a store
func (handler *StorageHandlersImpl) GetImage(params storage.GetImageParams) middleware.Responder {
	id := params.ID

	url, err := util.ImageStoreNameToURL(params.StoreName)
	if err != nil {
		return storage.NewGetImageDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    swag.Int64(http.StatusInternalServerError),
				Message: err.Error(),
			})
	}

	image, err := storageImageLayer.GetImage(context.TODO(), url, id)
	if err != nil {
		e := &models.Error{Code: swag.Int64(http.StatusNotFound), Message: err.Error()}
		return storage.NewGetImageNotFound().WithPayload(e)
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
	u, err := util.ImageStoreNameToURL(params.StoreName)
	if err != nil {
		return storage.NewListImagesDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    swag.Int64(http.StatusInternalServerError),
				Message: err.Error(),
			})
	}

	images, err := storageImageLayer.ListImages(context.TODO(), u, params.Ids)
	if err != nil {
		return storage.NewListImagesNotFound().WithPayload(
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
	u, err := util.ImageStoreNameToURL(params.StoreName)
	if err != nil {
		return storage.NewWriteImageDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    swag.Int64(http.StatusInternalServerError),
				Message: err.Error(),
			})
	}

	parent := &spl.Image{
		Store: u,
		ID:    params.ParentID,
	}

	var meta map[string][]byte

	if params.Metadatakey != nil && params.Metadataval != nil {
		meta = map[string][]byte{*params.Metadatakey: []byte(*params.Metadataval)}
	}

	image, err := storageImageLayer.WriteImage(context.TODO(), parent, params.ImageID, meta, params.Sum, params.ImageFile)
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

//CreateVolume : Create a Volume
func (handler *StorageHandlersImpl) CreateVolume(params storage.CreateVolumeParams) middleware.Responder {
	defer trace.End(trace.Begin("storage_handlers.CreateVolume"))

	//TODO: FIXME: add more errorcodes as we identify error scenarios.
	storeURL, err := util.VolumeStoreNameToURL(params.VolumeRequest.Store)
	if err != nil {
		log.Errorf("storagehandler: VolumeStoreName error: %s", err)
		return storage.NewCreateVolumeDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}

	byteMap := make(map[string][]byte)
	for key, value := range params.VolumeRequest.Metadata {
		byteMap[key] = []byte(value)
	}

	capacity := uint64(0)
	if params.VolumeRequest.Capacity < 0 {
		capacity = uint64(1024) //FIXME: this should look for a default cap and set or fail here.
	} else {
		capacity = uint64(params.VolumeRequest.Capacity)
	}

	volume, err := storageVolumeLayer.VolumeCreate(context.TODO(), params.VolumeRequest.Name, storeURL, capacity*1024, byteMap)
	if err != nil {
		log.Errorf("storagehandler: VolumeCreate error: %s", err)
		return storage.NewCreateVolumeDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}

	response := volumeToCreateResponse(volume, params.VolumeRequest)
	return storage.NewCreateVolumeCreated().WithPayload(&response)
}

//RemoveVolume : Remove a Volume from existence
func (handler *StorageHandlersImpl) RemoveVolume(storage.RemoveVolumeParams) middleware.Responder {
	defer trace.End(trace.Begin("storage_handlers.RemoveVolume"))
	return storage.NewRemoveVolumeOK() //TODO: this is just a stub for now.
}

// convert an SPL Image to a swagger-defined Image
func convertImage(image *spl.Image) *models.Image {
	var parent, selfLink *string

	// scratch image
	if image.Parent != nil {
		s := image.Parent.String()
		parent = &s
	}

	if image.SelfLink != nil {
		l := image.SelfLink.String()
		selfLink = &l
	}

	meta := make(map[string]string)
	if image.Metadata != nil {
		for k, v := range image.Metadata {
			meta[k] = string(v)
		}
	}

	return &models.Image{
		ID:       image.ID,
		SelfLink: selfLink,
		Parent:   parent,
		Metadata: meta,
		Store:    image.Store.String(),
	}
}

func (handler *StorageHandlersImpl) VolumeJoin(params storage.VolumeJoinParams) middleware.Responder {
	defer trace.End(trace.Begin("storage_handlers.RemoveVolume"))
	actualHandle := epl.GetHandle(params.JoinArgs.Handle)

	//Note: Name should already be populated by now.
	volume, err := storageVolumeLayer.VolumeGet(context.Background(), params.Name)
	if err != nil {
		log.Errorf("Volumes: StorageHandler : %#v", err)
		return storage.NewVolumeJoinInternalServerError().WithPayload(&models.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	log.Infof("found volume %s for volume join", volume.ID)
	actualHandle, err = vsphereSpl.VolumeJoin(context.Background(), actualHandle, volume, params.JoinArgs.MountPath, params.JoinArgs.Flags)
	if err != nil {
		log.Errorf("Volumes: StorageHandler : %#v", err)
		return storage.NewVolumeJoinInternalServerError().WithPayload(&models.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	log.Infof("volume %s has been joined to a container", volume.ID)
	return storage.NewVolumeJoinOK().WithPayload(actualHandle.String())
}

//utility functions

func volumeToCreateResponse(volume *spl.Volume, model *models.VolumeRequest) models.VolumeResponse {
	response := models.VolumeResponse{
		Driver:   model.Driver,
		Name:     volume.ID,
		Label:    volume.Label,
		Store:    model.Store,
		Metadata: model.Metadata,
	}
	return response
}

func findVolume(volumeList []*spl.Volume, ID string) (*spl.Volume, error) {
	for _, v := range volumeList {
		if v.ID == ID {
			return v, nil
		}
	}
	return &spl.Volume{}, fmt.Errorf("The volume with ID '%s' does not exist", ID)
}
