package handlers

import (
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/apiservers/docker/restapi/operations"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/volume"
)

type VolumeHandlersImpl struct{}

func (handlers *VolumeHandlersImpl) Configure(api *operations.DockerAPI) {
	api.VolumeDeleteVolumesNameHandler = volume.DeleteVolumesNameHandlerFunc(handlers.DeleteVolumesNames)
	api.VolumeGetVolumesHandler = volume.GetVolumesHandlerFunc(handlers.GetVolumes)
	api.VolumeGetVolumesNameHandler = volume.GetVolumesNameHandlerFunc(handlers.GetVolumesName)
	api.VolumePostVolumesCreateHandler = volume.PostVolumesCreateHandlerFunc(handlers.PostVolumesCreate)
}

func (handlers *VolumeHandlersImpl) DeleteVolumesNames(params volume.DeleteVolumesNameParams) middleware.Responder {
	return middleware.NotImplemented("operation volume.DeleteVolumesName has not yet been implemented")
}

func (handlers *VolumeHandlersImpl) GetVolumes(params volume.GetVolumesParams) middleware.Responder {
	return middleware.NotImplemented("operation volume.GetVolumes has not yet been implemented")
}

func (handlers *VolumeHandlersImpl) GetVolumesName(params volume.GetVolumesNameParams) middleware.Responder {
	return middleware.NotImplemented("operation volume.GetVolumesName has not yet been implemented")
}

func (handlers *VolumeHandlersImpl) PostVolumesCreate(params volume.PostVolumesCreateParams) middleware.Responder {
	return middleware.NotImplemented("operation volume.PostVolumesCreate has not yet been implemented")
}
