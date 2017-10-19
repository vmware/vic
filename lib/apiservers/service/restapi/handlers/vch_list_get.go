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
	"fmt"
	"net/url"
	"path"

	"github.com/docker/docker/opts"
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/util"
	"github.com/vmware/vic/lib/apiservers/service/restapi/operations"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// VCHListGet is the handler for listing VCHs
type VCHListGet struct {
}

// VCHListGet is the handler for listing VCHs within a Datacenter
type VCHDatacenterListGet struct {
}

func (h *VCHListGet) Handle(params operations.GetTargetTargetVchParams, principal interface{}) middleware.Responder {
	d, err := buildData(params.HTTPRequest.Context(),
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		nil,
		params.ComputeResource)
	if err != nil {
		return operations.NewGetTargetTargetVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	op := trace.NewOperation(params.HTTPRequest.Context(), "vch_list get handler")
	vchs, err := listVCHs(op, d)
	if err != nil {
		return operations.NewGetTargetTargetVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewGetTargetTargetVchOK().WithPayload(operations.GetTargetTargetVchOKBody{Vchs: vchs})
}

func (h *VCHDatacenterListGet) Handle(params operations.GetTargetTargetDatacenterDatacenterVchParams, principal interface{}) middleware.Responder {
	d, err := buildData(params.HTTPRequest.Context(),
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		&params.Datacenter,
		params.ComputeResource)
	if err != nil {
		return operations.NewGetTargetTargetVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	op := trace.NewOperation(params.HTTPRequest.Context(), "vch_list get handler")
	vchs, err := listVCHs(op, d)
	if err != nil {
		return operations.NewGetTargetTargetVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewGetTargetTargetVchOK().WithPayload(operations.GetTargetTargetVchOKBody{Vchs: vchs})
}

func listVCHs(op trace.Operation, d *data.Data) ([]*models.VCHListItem, error) {
	validator, err := validateTarget(op, d)
	if err != nil {
		return nil, util.WrapError(400, err)
	}

	executor := management.NewDispatcher(validator.Context, validator.Session, nil, false)
	vchs, err := executor.SearchVCHs(validator.ClusterPath)
	if err != nil {
		return nil, util.NewError(500, fmt.Sprintf("Failed to search VCHs in %s: %s", validator.ResourcePoolPath, err))
	}

	return vchsToModels(op, vchs, executor), nil
}

func vchsToModels(op trace.Operation, vchs []*vm.VirtualMachine, executor *management.Dispatcher) []*models.VCHListItem {
	installerVer := version.GetBuild()
	payload := make([]*models.VCHListItem, 0)
	for _, vch := range vchs {
		var version *version.Build
		var dockerHost string
		var adminPortal string
		if vchConfig, err := executor.GetNoSecretVCHConfig(vch); err == nil {
			version = vchConfig.Version

			if public := vchConfig.ExecutorConfig.Networks["public"]; public != nil {
				if publicIP := public.Assigned.IP; publicIP != nil {
					var dockerPort = opts.DefaultTLSHTTPPort
					if vchConfig.HostCertificate.IsNil() {
						dockerPort = opts.DefaultHTTPPort
					}

					dockerHost = fmt.Sprintf("%s:%d", publicIP, dockerPort)
					adminPortal = fmt.Sprintf("https://%s:2378", publicIP)
				}
			}
		}

		name := path.Base(vch.InventoryPath)

		model := &models.VCHListItem{ID: vch.Reference().Value, Name: name, AdminPortal: adminPortal, DockerHost: dockerHost}

		if version != nil {
			model.Version = version.ShortVersion()
			model.UpgradeStatus = upgradeStatusMessage(op, vch, installerVer, version)
		}

		payload = append(payload, model)
	}

	return payload
}
