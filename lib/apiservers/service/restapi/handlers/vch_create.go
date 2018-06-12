// Copyright 2017-2018 VMware, Inc. All Rights Reserved.
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
	"net/http"
	"path"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/client"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/decode"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/errors"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/logging"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/target"
	"github.com/vmware/vic/lib/apiservers/service/restapi/operations"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/vchlog"
	"github.com/vmware/vic/pkg/trace"
)

const (
	logFile = "vic-machine.log" // name of local log file
)

// VCHCreate is the handler for creating a VCH without specifying a datacenter
type VCHCreate struct {
	vchCreate
}

// VCHDatacenterCreate is the handler for creating a VCH within a specified datacenter
type VCHDatacenterCreate struct {
	vchCreate
}

// vchCreate allows for VCHCreate and VCHDatacenterCreate to share common code without polluting the package
type vchCreate struct{}

// Handle is the handler implementation for creating a VCH without specifying a datacenter
func (h *VCHCreate) Handle(params operations.PostTargetTargetVchParams, principal interface{}) middleware.Responder {
	op := trace.FromContext(params.HTTPRequest.Context(), "VCHCreate")

	b := target.Params{
		Target:     params.Target,
		Thumbprint: params.Thumbprint,
	}

	task, err := h.handle(op, b, principal, params.Vch)
	if err != nil {
		return operations.NewPostTargetTargetVchDefault(errors.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewPostTargetTargetVchCreated().WithPayload(operations.PostTargetTargetVchCreatedBody{Task: task})
}

// Handle is the handler implementation for creating a VCH within a specified datacenter
func (h *VCHDatacenterCreate) Handle(params operations.PostTargetTargetDatacenterDatacenterVchParams, principal interface{}) middleware.Responder {
	op := trace.FromContext(params.HTTPRequest.Context(), "VCHDatacenterCreate")

	b := target.Params{
		Target:     params.Target,
		Thumbprint: params.Thumbprint,
		Datacenter: &params.Datacenter,
	}

	task, err := h.handle(op, b, principal, params.Vch)
	if err != nil {
		return operations.NewPostTargetTargetDatacenterDatacenterVchDefault(errors.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewPostTargetTargetDatacenterDatacenterVchCreated().WithPayload(operations.PostTargetTargetDatacenterDatacenterVchCreatedBody{Task: task})
}

// handle creates a VCH with the settings from vch at the target described by params, using the credentials from
// principal. If an error occurs validating the requested settings, a 400 is returned. If an error occurs during
// creation, a 500 is returned. Currently, no task is ever returned.
func (h *vchCreate) handle(op trace.Operation, params target.Params, principal interface{}, vch *models.VCH) (*strfmt.URI, error) {
	datastoreLogger := logging.SetUpLogger(&op)
	defer datastoreLogger.Close()

	d, hc, err := target.Validate(op, management.ActionCreate, params, principal)
	if err != nil {
		return nil, err
	}

	err = decode.ProcessVCH(op, d, vch, hc.Finder())
	if err != nil {
		return nil, err
	}

	return h.handleCreate(op, d, hc, datastoreLogger)
}

func (h *vchCreate) handleCreate(op trace.Operation, d *data.Data, hc *client.HandlerClient, receiver vchlog.Receiver) (*strfmt.URI, error) {
	validator := hc.Validator() // TODO (#6032): Move some of the logic that uses this into methods on hc

	vchConfig, err := validator.Validate(op, d, false)
	if err != nil {
		issues := validator.GetIssues()
		messages := make([]string, 0, len(issues))
		for _, issue := range issues {
			messages = append(messages, issue.Error())
		}

		return nil, errors.NewError(http.StatusBadRequest, "failed to validate VCH: %s", strings.Join(messages, ", "))
	}

	vConfig := validator.AddDeprecatedFields(op, vchConfig, d)

	// TODO (#6714): make this configurable
	images := common.Images{}
	vConfig.ImageFiles, err = images.CheckImagesFiles(op, true)
	vConfig.ApplianceISO = path.Base(images.ApplianceISO)
	vConfig.BootstrapISO = path.Base(images.BootstrapISO)

	vConfig.HTTPProxy = d.HTTPProxy
	vConfig.HTTPSProxy = d.HTTPSProxy

	err = hc.Executor().CreateVCH(vchConfig, vConfig, receiver)
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "failed to create VCH: %s", err)
	}

	return nil, nil
}
