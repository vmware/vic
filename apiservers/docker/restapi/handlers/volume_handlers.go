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
