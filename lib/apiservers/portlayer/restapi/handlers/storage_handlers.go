// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/storage"
	vicarchive "github.com/vmware/vic/lib/archive"
	epl "github.com/vmware/vic/lib/portlayer/exec"
	spl "github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/lib/portlayer/storage/nfs"
	"github.com/vmware/vic/lib/portlayer/storage/vsphere"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
)

// StorageHandlersImpl is the receiver for all of the storage handler methods
type StorageHandlersImpl struct {
	imageCache     *spl.NameLookupCache
	volumeCache    *spl.VolumeLookupCache
	containerStore *spl.ContainerStore
}

const (
	nfsScheme = "nfs"
	dsScheme  = "ds"

	uidQueryKey = "uid"
	gidQueryKey = "gid"
)

// Configure assigns functions to all the storage api handlers
func (h *StorageHandlersImpl) Configure(api *operations.PortLayerAPI, handlerCtx *HandlerContext) {
	var err error

	ctx := context.Background()
	op := trace.NewOperation(ctx, "configure")

	if len(spl.Config.ImageStores) == 0 {
		op.Panicf("No image stores provided; unable to instantiate storage layer")
	}

	imageStoreURL := spl.Config.ImageStores[0]
	// TODO: support multiple image stores. Right now we only support the first one
	if len(spl.Config.ImageStores) > 1 {
		op.Warnf("Multiple image stores found. Multiple image stores are not yet supported. Using [%s] %s", imageStoreURL.Host, imageStoreURL.Path)
	}

	ds, err := vsphere.NewImageStore(op, handlerCtx.Session, &imageStoreURL)
	if err != nil {
		op.Panicf("Cannot instantiate storage layer: %s", err)
	}

	// The imagestore is implemented via a cache which is backed via an
	// implementation that writes to disks.  The cache is used to avoid
	// expensive metadata lookups.
	h.imageCache = spl.NewLookupCache(ds)

	spl.RegisterImporter(op, imageStoreURL.String(), ds)
	spl.RegisterExporter(op, imageStoreURL.String(), ds)

	c, err := spl.NewContainerStore(op, handlerCtx.Session, h.imageCache)
	if err != nil {
		op.Panicf("Couldn't create containerStore: %s", err.Error())
	}
	h.containerStore = c

	spl.RegisterImporter(op, "container", h.containerStore)
	spl.RegisterExporter(op, "container", h.containerStore)

	// add the volume stores, errors are logged within this function.
	h.configureVolumeStores(op, handlerCtx)

	api.StorageCreateImageStoreHandler = storage.CreateImageStoreHandlerFunc(h.CreateImageStore)
	api.StorageGetImageHandler = storage.GetImageHandlerFunc(h.GetImage)
	api.StorageGetImageTarHandler = storage.GetImageTarHandlerFunc(h.GetImageTar)
	api.StorageListImagesHandler = storage.ListImagesHandlerFunc(h.ListImages)
	api.StorageWriteImageHandler = storage.WriteImageHandlerFunc(h.WriteImage)
	api.StorageDeleteImageHandler = storage.DeleteImageHandlerFunc(h.DeleteImage)

	api.StorageVolumeStoresListHandler = storage.VolumeStoresListHandlerFunc(h.VolumeStoresList)
	api.StorageCreateVolumeHandler = storage.CreateVolumeHandlerFunc(h.CreateVolume)
	api.StorageRemoveVolumeHandler = storage.RemoveVolumeHandlerFunc(h.RemoveVolume)
	api.StorageVolumeJoinHandler = storage.VolumeJoinHandlerFunc(h.VolumeJoin)
	api.StorageListVolumesHandler = storage.ListVolumesHandlerFunc(h.VolumesList)
	api.StorageGetVolumeHandler = storage.GetVolumeHandlerFunc(h.GetVolume)

	api.StorageExportArchiveHandler = storage.ExportArchiveHandlerFunc(h.ExportArchive)
	api.StorageImportArchiveHandler = storage.ImportArchiveHandlerFunc(h.ImportArchive)
}

func (h *StorageHandlersImpl) configureVolumeStores(op trace.Operation, handlerCtx *HandlerContext) {
	var (
		vs  spl.VolumeStorer
		err error
	)

	h.volumeCache = spl.NewVolumeLookupCache(op)

	// Configure the datastores
	// Each volume store name maps to a datastore + path, which can be referred to by the name.
	for name, dsurl := range spl.Config.VolumeLocations {
		switch dsurl.Scheme {
		case nfsScheme:
			vs, err = createNFSVolumeStore(op, dsurl, name)
		case dsScheme:
			vs, err = createVsphereVolumeStore(op, dsurl, name, handlerCtx)
		default:
			err = fmt.Errorf("unknown scheme for %s", dsurl.String())
			log.Error(err.Error())
		}

		// if an error has been logged skip volume store cache addition
		if err != nil {
			continue
		}

		op.Infof("Adding volume store %s (%s)", name, dsurl.String())
		if _, err = h.volumeCache.AddStore(op, name, vs); err != nil {
			log.Errorf("volume addition error %s", err)
		}

		spl.RegisterImporter(op, dsurl.String(), vs)
		spl.RegisterExporter(op, dsurl.String(), vs)
	}
}

// CreateImageStore creates a new image store
func (h *StorageHandlersImpl) CreateImageStore(params storage.CreateImageStoreParams) middleware.Responder {
	op := trace.NewOperation(context.Background(), fmt.Sprintf("CreateImageStore(%s)", params.Body.Name))
	url, err := h.imageCache.CreateImageStore(op, params.Body.Name)
	if err != nil {
		if os.IsExist(err) {
			return storage.NewCreateImageStoreConflict().WithPayload(
				&models.Error{
					Code:    http.StatusConflict,
					Message: "An image store with that name already exists",
				})
		}

		return storage.NewCreateImageStoreDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
	}
	s := &models.StoreURL{
		Code: http.StatusCreated,
		URL:  url.String(),
	}
	return storage.NewCreateImageStoreCreated().WithPayload(s)
}

// GetImage retrieves an image from a store
func (h *StorageHandlersImpl) GetImage(params storage.GetImageParams) middleware.Responder {
	id := params.ID

	url, err := util.ImageStoreNameToURL(params.StoreName)
	if err != nil {
		return storage.NewGetImageDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
	}

	op := trace.NewOperation(context.Background(), fmt.Sprintf("GetImage(%s)", id))
	image, err := h.imageCache.GetImage(op, url, id)
	if err != nil {
		e := &models.Error{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		return storage.NewGetImageNotFound().WithPayload(e)
	}

	result := convertImage(image)
	return storage.NewGetImageOK().WithPayload(result)
}

// DeleteImage deletes an image from a store
func (h *StorageHandlersImpl) DeleteImage(params storage.DeleteImageParams) middleware.Responder {

	ferr := func(err error, code int) middleware.Responder {
		log.Errorf("DeleteImage: error %s", err.Error())
		return storage.NewDeleteImageDefault(code).WithPayload(
			&models.Error{
				Code:    int64(code),
				Message: err.Error(),
			})
	}

	imageURL, err := util.ImageURL(params.StoreName, params.ID)
	if err != nil {
		return ferr(err, http.StatusInternalServerError)
	}

	image, err := spl.Parse(imageURL)
	if err != nil {
		return ferr(err, http.StatusInternalServerError)
	}

	keepNodes := make([]*url.URL, len(params.KeepNodes))
	for idx, kn := range params.KeepNodes {
		k, err := url.Parse(kn)
		if err != nil {
			return ferr(err, http.StatusInternalServerError)
		}

		keepNodes[idx] = k
	}

	op := trace.NewOperation(context.Background(), fmt.Sprintf("DeleteBranch(%s)", image.ID))
	deletedImages, err := h.imageCache.DeleteBranch(op, image, keepNodes)
	if err != nil {
		switch {
		case spl.IsErrImageInUse(err):
			return ferr(err, http.StatusLocked)

		case os.IsNotExist(err):
			return ferr(err, http.StatusNotFound)

		default:
			return ferr(err, http.StatusInternalServerError)
		}
	}

	result := make([]*models.Image, len(deletedImages))
	for idx, image := range deletedImages {
		result[idx] = convertImage(image)
	}

	return storage.NewDeleteImageOK().WithPayload(result)
}

// GetImageTar returns an image tar file
func (h *StorageHandlersImpl) GetImageTar(params storage.GetImageTarParams) middleware.Responder {
	return middleware.NotImplemented("operation storage.GetImageTar has not yet been implemented")
}

// ListImages returns a list of images in a store
func (h *StorageHandlersImpl) ListImages(params storage.ListImagesParams) middleware.Responder {
	u, err := util.ImageStoreNameToURL(params.StoreName)
	if err != nil {
		return storage.NewListImagesDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
	}

	op := trace.NewOperation(context.Background(), fmt.Sprintf("ListImages(%s, %q)", u.String(), params.Ids))
	images, err := h.imageCache.ListImages(op, u, params.Ids)
	if err != nil {
		return storage.NewListImagesNotFound().WithPayload(
			&models.Error{
				Code:    http.StatusNotFound,
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
func (h *StorageHandlersImpl) WriteImage(params storage.WriteImageParams) middleware.Responder {
	u, err := util.ImageStoreNameToURL(params.StoreName)
	if err != nil {
		return storage.NewWriteImageDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    http.StatusInternalServerError,
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

	op := trace.NewOperation(context.Background(), fmt.Sprintf("WriteImage(%s)", params.ImageID))
	image, err := h.imageCache.WriteImage(op, parent, params.ImageID, meta, params.Sum, params.ImageFile)
	if err != nil {
		return storage.NewWriteImageDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
	}
	i := convertImage(image)
	return storage.NewWriteImageCreated().WithPayload(i)
}

// VolumeStoresList lists the configured volume stores and their datastore path URIs.
func (h *StorageHandlersImpl) VolumeStoresList(params storage.VolumeStoresListParams) middleware.Responder {
	defer trace.End(trace.Begin("storage_handlers.VolumeStoresList"))

	op := trace.NewOperation(context.Background(), "VolumeStoresList")
	stores, err := h.volumeCache.VolumeStoresList(op)
	if err != nil {
		return storage.NewVolumeStoresListInternalServerError().WithPayload(
			&models.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
	}

	resp := &models.VolumeStoresListResponse{Stores: stores}

	return storage.NewVolumeStoresListOK().WithPayload(resp)
}

//CreateVolume : Create a Volume
func (h *StorageHandlersImpl) CreateVolume(params storage.CreateVolumeParams) middleware.Responder {
	defer trace.End(trace.Begin("storage_handlers.CreateVolume"))

	//TODO: FIXME: add more errorcodes as we identify error scenarios.
	storeURL, err := util.VolumeStoreNameToURL(params.VolumeRequest.Store)
	if err != nil {
		log.Errorf("storagehandler: VolumeStoreName error: %s", err)
		return storage.NewCreateVolumeInternalServerError().WithPayload(&models.Error{
			Code:    http.StatusInternalServerError,
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

	op := trace.NewOperation(context.Background(), fmt.Sprintf("VolumeCreate(%s)", params.VolumeRequest.Name))
	volume, err := h.volumeCache.VolumeCreate(op, params.VolumeRequest.Name, storeURL, capacity*1024, byteMap)
	if err != nil {
		log.Errorf("storagehandler: VolumeCreate error: %#v", err)

		if os.IsExist(err) {
			return storage.NewCreateVolumeConflict().WithPayload(&models.Error{
				Code:    http.StatusConflict,
				Message: err.Error(),
			})
		}

		if _, ok := err.(spl.VolumeStoreNotFoundError); ok {
			return storage.NewCreateVolumeNotFound().WithPayload(&models.Error{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			})
		}

		return storage.NewCreateVolumeInternalServerError().WithPayload(&models.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	response := volumeToCreateResponse(volume, params.VolumeRequest)
	return storage.NewCreateVolumeCreated().WithPayload(&response)
}

//GetVolume : Gets a handle to a volume
func (h *StorageHandlersImpl) GetVolume(params storage.GetVolumeParams) middleware.Responder {
	defer trace.End(trace.Begin(params.Name))

	op := trace.NewOperation(context.Background(), fmt.Sprintf("VolumeGet(%s)", params.Name))
	data, err := h.volumeCache.VolumeGet(op, params.Name)
	if err == os.ErrNotExist {
		return storage.NewGetVolumeNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	response, err := fillVolumeModel(data)
	if err != nil {
		return storage.NewListVolumesInternalServerError().WithPayload(&models.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	log.Debugf("VolumeGet returned : %#v", response)
	return storage.NewGetVolumeOK().WithPayload(&response)
}

//RemoveVolume : Remove a Volume from existence
func (h *StorageHandlersImpl) RemoveVolume(params storage.RemoveVolumeParams) middleware.Responder {
	defer trace.End(trace.Begin("storage_handlers.RemoveVolume"))

	op := trace.NewOperation(context.Background(), fmt.Sprintf("VolumeDestroy(%s)", params.Name))
	err := h.volumeCache.VolumeDestroy(op, params.Name)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return storage.NewRemoveVolumeNotFound().WithPayload(&models.Error{
				Message: err.Error(),
			})

		case spl.IsErrVolumeInUse(err):
			return storage.NewRemoveVolumeConflict().WithPayload(&models.Error{
				Message: err.Error(),
			})

		default:
			return storage.NewRemoveVolumeInternalServerError().WithPayload(&models.Error{
				Message: err.Error(),
			})
		}
	}
	return storage.NewRemoveVolumeOK()
}

//VolumesList : Lists available volumes for use
func (h *StorageHandlersImpl) VolumesList(params storage.ListVolumesParams) middleware.Responder {
	defer trace.End(trace.Begin(""))
	var result []*models.VolumeResponse

	op := trace.NewOperation(context.Background(), "VolumeList")
	portlayerVolumes, err := h.volumeCache.VolumesList(op)
	if err != nil {
		log.Error(err)
		return storage.NewListVolumesInternalServerError().WithPayload(&models.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	log.Debugf("volumes fetched from list call : %#v", portlayerVolumes)

	for i := range portlayerVolumes {
		model, err := fillVolumeModel(portlayerVolumes[i])
		if err != nil {
			log.Error(err)
			return storage.NewListVolumesInternalServerError().WithPayload(&models.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}

		result = append(result, &model)
	}

	log.Debugf("volumes returned from list call : %#v", result)
	return storage.NewListVolumesOK().WithPayload(result)
}

//VolumeJoin : modifies the config spec of a container to mount the specified container
func (h *StorageHandlersImpl) VolumeJoin(params storage.VolumeJoinParams) middleware.Responder {
	defer trace.End(trace.Begin(""))

	op := trace.NewOperation(context.Background(), fmt.Sprintf("VolumeJoin(%s)", params.Name))

	actualHandle := epl.GetHandle(params.JoinArgs.Handle)

	//Note: Name should already be populated by now.
	volume, err := h.volumeCache.VolumeGet(op, params.Name)
	if err != nil {
		op.Errorf("Volumes: StorageHandler : %#v", err)

		return storage.NewVolumeJoinInternalServerError().WithPayload(&models.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	switch volume.Device.DiskPath().Scheme {
	case nfsScheme:
		actualHandle, err = nfs.VolumeJoin(op, actualHandle, volume, params.JoinArgs.MountPath, params.JoinArgs.Flags)
	case dsScheme:
		actualHandle, err = vsphere.VolumeJoin(op, actualHandle, volume, params.JoinArgs.MountPath, params.JoinArgs.Flags)
	default:
		err = fmt.Errorf("unknown scheme (%s) for Volume (%#v)", volume.Device.DiskPath().Scheme, *volume)
	}

	if err != nil {
		op.Errorf("Volumes: StorageHandler : %#v", err)

		return storage.NewVolumeJoinInternalServerError().WithPayload(&models.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	op.Infof("volume %s has been joined to a container", volume.ID)
	return storage.NewVolumeJoinOK().WithPayload(actualHandle.String())
}

// ImportArchive takes an input tar archive and unpacks to destination
func (h *StorageHandlersImpl) ImportArchive(params storage.ImportArchiveParams) middleware.Responder {
	defer trace.End(trace.Begin(""))
	defer params.Archive.Close()

	id := params.DeviceID
	op := trace.NewOperation(context.Background(), "ImportArchive: %s", id)

	filterSpec, err := vicarchive.DecodeFilterSpec(op, params.FilterSpec)
	if err != nil {
		// hickeng: should be a 422 instead of 500
		return storage.NewImportArchiveInternalServerError()
	}

	store, ok := spl.GetImporter(params.Store)
	if !ok {
		return storage.NewImportArchiveNotFound()
	}

	err = store.Import(op, id, filterSpec, params.Archive)
	if err != nil {
		// hickeng: see if we can return usefully typed errors here
		return storage.NewExportArchiveInternalServerError()
	}

	return storage.NewImportArchiveOK()
}

// ExportArchive creates a tar archive and returns to caller
func (h *StorageHandlersImpl) ExportArchive(params storage.ExportArchiveParams) middleware.Responder {
	defer trace.End(trace.Begin(""))

	id := params.DeviceID
	ancestor := ""
	if params.Ancestor != nil {
		ancestor = *params.Ancestor
	}

	op := trace.NewOperation(context.Background(), "ExportArchive: %s:%s", id, ancestor)

	filterSpec, err := vicarchive.DecodeFilterSpec(op, params.FilterSpec)
	if err != nil {
		// hickeng: should be a 422 instead of 500
		return storage.NewExportArchiveInternalServerError()
	}

	store, ok := spl.GetExporter(params.Store)
	if !ok {
		// TODO: this should be a 404 but cannot seem to figure that out #shamed
		return storage.NewExportArchiveNotFound()
	}

	r, err := store.Export(op, id, ancestor, filterSpec, params.Data)
	if err != nil {
		// hickeng: we're in need of typed errors - should check for id not found for 404 return
		return storage.NewExportArchiveInternalServerError()
	}

	return NewStreamOutputHandler("ExportArchive").WithPayload(NewFlushingReader(r), params.DeviceID, nil)
}

//utility functions

// convert an SPL Image to a swagger-defined Image
func convertImage(image *spl.Image) *models.Image {
	var parent, selfLink string

	// scratch image
	if image.ParentLink != nil {
		parent = image.ParentLink.String()
	}

	if image.SelfLink != nil {
		selfLink = image.SelfLink.String()
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

func fillVolumeModel(volume *spl.Volume) (models.VolumeResponse, error) {
	storeName, err := util.VolumeStoreName(volume.Store)
	if err != nil {
		return models.VolumeResponse{}, err
	}

	metadata := createMetadataMap(volume)

	model := models.VolumeResponse{
		Name:     volume.ID,
		Driver:   "vsphere",
		Store:    storeName,
		Metadata: metadata,
		Label:    volume.Label,
	}

	return model, nil
}

func createMetadataMap(volume *spl.Volume) map[string]string {
	stringMap := make(map[string]string)
	for k, v := range volume.Info {
		stringMap[k] = string(v)
	}
	return stringMap
}

func createNFSVolumeStore(op trace.Operation, dsurl *url.URL, name string) (spl.VolumeStorer, error) {
	var err error
	uid, gid, err := parseUIDAndGID(dsurl)
	if err != nil {
		op.Errorf("%s", err.Error())
		return nil, err
	}

	// XXX replace with the vch name
	mnt := nfs.NewMount(dsurl, "vic", uint32(uid), uint32(gid))
	vs, err := nfs.NewVolumeStore(op, name, mnt)
	if err != nil {
		op.Errorf("%s", err.Error())
		return nil, err
	}

	return vs, nil
}

func parseUIDAndGID(queryURL *url.URL) (int, int, error) {
	var err error
	uid := nfs.DefaultUID
	gid := nfs.DefaultUID

	vsUID := queryURL.Query().Get(uidQueryKey)
	vsGID := queryURL.Query().Get(gidQueryKey)

	if vsGID == "" {
		vsGID = vsUID
	}

	if vsUID != "" {
		uid, err = strconv.Atoi(vsUID)
		if err != nil {
			return -1, -1, err
		}
	}

	if vsGID != "" {
		gid, err = strconv.Atoi(vsGID)
		if err != nil {
			return -1, -1, err
		}
	}

	if uid < 0 {
		return -1, -1, fmt.Errorf("supplied url (%s) for nfs volume store has invalid uid : (%d)", queryURL.String(), uid)
	}

	if gid < 0 {
		return -1, -1, fmt.Errorf("supplied url (%s) for nfs volume store has invalid gid : (%d)", queryURL.String(), gid)
	}

	return uid, gid, nil
}

func createVsphereVolumeStore(op trace.Operation, dsurl *url.URL, name string, handlerCtx *HandlerContext) (spl.VolumeStorer, error) {
	ds, err := datastore.NewHelperFromURL(op, handlerCtx.Session, dsurl)
	if err != nil {
		err = fmt.Errorf("cannot find datastores: %s", err)
		op.Errorf("%s", err.Error())
		return nil, err
	}

	vs, err := vsphere.NewVolumeStore(op, name, handlerCtx.Session, ds)
	if err != nil {
		err = fmt.Errorf("cannot instantiate the volume store: %s", err)
		op.Errorf("%s", err.Error())
		return nil, err
	}
	return vs, nil
}
