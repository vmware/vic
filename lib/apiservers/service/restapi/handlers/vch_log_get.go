// Copyright 2017 VMware, Inc. All Rights Reserved.
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
	"io/ioutil"
	"net/url"

	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/util"
	"github.com/vmware/vic/lib/apiservers/service/restapi/operations"
	"github.com/vmware/vic/lib/install/vchlog"
)

// VCHLogGet is the handler for getting the log messages for a VCH
type VCHLogGet struct {
}

// VCHDatacenterLogGet is the handler for getting the log messages for a VCH within a Datacenter
type VCHDatacenterLogGet struct {
}

func (h *VCHLogGet) Handle(params operations.GetTargetTargetVchVchIDLogParams, principal interface{}) middleware.Responder {
	d, err := buildData(params.HTTPRequest.Context(),
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		nil,
		nil)
	if err != nil {
		return operations.NewGetTargetTargetVchVchIDLogDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	d.ID = params.VchID

	response, err := downloadAndReadLogFileFromDatastore(params.HTTPRequest.Context(), d.ID)

	if err != nil {
		return operations.NewGetTargetTargetVchVchIDLogDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewGetTargetTargetVchVchIDLogOK().WithPayload(response)
}

func (h *VCHDatacenterLogGet) Handle(params operations.GetTargetTargetDatacenterDatacenterVchVchIDLogParams, principal interface{}) middleware.Responder {
	d, err := buildData(params.HTTPRequest.Context(),
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		&params.Datacenter,
		nil)
	if err != nil {
		return operations.NewGetTargetTargetDatacenterDatacenterVchVchIDLogDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	d.ID = params.VchID

	response, err := downloadAndReadLogFileFromDatastore(params.HTTPRequest.Context(), d.ID)

	if err != nil {
		return operations.NewGetTargetTargetDatacenterDatacenterVchVchIDLogDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewGetTargetTargetDatacenterDatacenterVchVchIDLogOK().WithPayload(response)
}

// downloadAndReadLogFileFromDatastore downloads log file from VCH datastore folder and returns the content of log
// file upon success, and returns error if log file is not found or read fails.
func downloadAndReadLogFileFromDatastore(ctx context.Context, ID string) (string, error) {
	reader, err := vchlog.DownloadLogFile(ctx)
	if err != nil {
		return "", util.NewError(404, fmt.Sprintf("No log file found for VCH %s", ID))
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", util.NewError(404, fmt.Sprintf("Error reading log file from datastore for VCH %s", ID))
	}

	return string(bytes), nil
}
