// Copyright 2018 VMware, Inc. All Rights Reserved.
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

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/vmware/vic/cmd/vic-machine/configure"
	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/client"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/errors"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/target"
	"github.com/vmware/vic/lib/apiservers/service/restapi/operations"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
)

type VCHPatch struct {
	vchPatch
}

type VCHDatacenterPatch struct {
	vchPatch
}

// vchPatch allows for VCHPatch and VCHDatacenterPatch to share common code without polluting the package
type vchPatch struct{}

func (h *VCHPatch) Handle(params operations.PatchTargetTargetVchVchIDParams, principal interface{}) middleware.Responder {
	op := trace.FromContext(params.HTTPRequest.Context(), "VCHPatch: %s", params.VchID)

	b := target.Params{
		Target:     params.Target,
		Thumbprint: params.Thumbprint,
		VCHID:      &params.VchID,
	}

	task, err := h.handle(op, b, principal, params.Vch)
	if err != nil {
		return operations.NewPatchTargetTargetVchVchIDDefault(errors.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewPatchTargetTargetVchVchIDAccepted().WithPayload(operations.PatchTargetTargetVchVchIDAcceptedBody{Task: task})
}

func (h *VCHDatacenterPatch) Handle(params operations.PatchTargetTargetDatacenterDatacenterVchVchIDParams, principal interface{}) middleware.Responder {
	op := trace.FromContext(params.HTTPRequest.Context(), "VCHDatacenterPatch: %s", params.VchID)

	b := target.Params{
		Target:     params.Target,
		Thumbprint: params.Thumbprint,
		Datacenter: &params.Datacenter,
		VCHID:      &params.VchID,
	}

	task, err := h.handle(op, b, principal, params.Vch)
	if err != nil {
		return operations.NewPatchTargetTargetDatacenterDatacenterVchVchIDDefault(errors.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewPatchTargetTargetDatacenterDatacenterVchVchIDAccepted().WithPayload(operations.PatchTargetTargetDatacenterDatacenterVchVchIDAcceptedBody{Task: task})
}

func (h *vchPatch) handle(op trace.Operation, params target.Params, principal interface{}, vch *models.VCH) (*strfmt.URI, error) {
	// validate target
	d, hc, err := target.Validate(op, management.ActionConfigure, params, principal)
	if err != nil {
		return nil, err
	}

	// check for immutable fields, bail out if there is any
	err = h.checkImmutableVCHProperties(op, vch)
	if err != nil {
		return nil, err
	}

	// build configure object
	// TODO [AngieCris]: duplicate code with API PUT
	c, err := h.buildConfigure(op, d, hc.Finder(), vch)
	if err != nil {
		return nil, err
	}

	// get old vch config and data
	vchConfig, oldData, err := hc.GetDataAndVCHSecretConfig(op, c.Data)
	if err != nil {
		return nil, err
	}

	// merge old and new data
	err = oldData.CopyNonEmpty(c.Data)
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "error copying new VCH data: %s", err)
	}
	c.Data = oldData

	// perform vch configure
	// TODO [AngieCris]: duplicate code with API PUT
	err = h.handleConfigure(op, c, vch, vchConfig, hc)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// TODO [AngieCris]: better log messages (more detailized)
func (h *vchPatch) checkImmutableVCHProperties(op trace.Operation, vch *models.VCH) error {
	// vch name
	if vch.Name != "" { // TODO [AngieCris]: can't differentiate between nil (not provided) and "" (empty string provided)
		return errors.NewError(http.StatusBadRequest, "VCH name is not re-configurable")
	}

	// networks: bridge, management, client, public
	if vch.Network != nil {
		// TODO [AngieCris]: what if the network fields are non-nil, but portgroup or nameservers is empty? Is it considered BadRequest?
		if vch.Network.Bridge != nil || vch.Network.Client != nil || vch.Network.Management != nil || vch.Network.Public != nil {
			return errors.NewError(http.StatusBadRequest, "VCH network is not re-configurable")
		}
	}

	// registry info
	if vch.Registry != nil {
		if vch.Registry.Insecure != nil || vch.Registry.Whitelist != nil { // TODO [AngieCris]: here empty list is not allowed (double check)
			return errors.NewError(http.StatusBadRequest, "VCH registries are not reconfigurable")
		}
	}

	// tls and ca
	// TODO [AngieCris]: skip for now because certFactory is stored configure.Configure instead of data, and it's processed later (according to CLI) (need to figure out why)

	// storage: datastore path, base image size
	if vch.Storage != nil {
		if vch.Storage.ImageStores != nil { // TODO [AngieCris]: here empty list is not allowed (double check)
			return errors.NewError(http.StatusBadRequest, "VCH image store path is not reconfigurable")
		}
		if vch.Storage.BaseImageSize != nil {
			return errors.NewError(http.StatusBadRequest, "VCH base image size is not reconfigurable")
		}
	}

	// endpoint VM resource limits
	if vch.Endpoint != nil {
		if vch.Endpoint.CPU != nil || vch.Endpoint.Memory != nil {
			return errors.NewError(http.StatusBadRequest, "VCH endpoint resource limits are not reconfigurable")
		}
	}

	// syslog config
	if vch.SyslogAddr != "" {
		return errors.NewError(http.StatusBadRequest, "VCH syslog address is not reconfigurable")
	}

	// TODO [AngieCris]: what to do with vch.Runtime (that includes docker host, admin portal and power status)? Error out?

	return nil
}

func (h *vchPatch) buildConfigure(op trace.Operation, d *data.Data, finder client.Finder, vch *models.VCH) (*configure.Configure, error) {
	c := &configure.Configure{Data: d}

	// parse input vch
	if vch != nil {
		if vch.Version != "" && version.String() != string(vch.Version) {
			return nil, errors.NewError(http.StatusBadRequest, "invalid version: %s", vch.Version)
		}

		// TODO [AngieCris]: prototype stage, only debug configurable
		debug := int(vch.Debug)
		c.Debug.Debug = &debug

		// TODO [AngieCris]: timeout is not configurable from API. Make it less hacky
		// Set default timeout to 3 minutes
		c.Timeout = DefaultTimeout
	}

	return c, nil
}

// TODO [AngieCris]: complete duplicate from API PUT implementation. Potential plan: have a configure package that includes common code for both PUT and PATCH
func (h *vchPatch) handleConfigure(op trace.Operation, c *configure.Configure, vch *models.VCH, config *config.VirtualContainerHostConfigSpec, hc *client.HandlerClient) error {
	validator := hc.Validator()

	grantPerms := false
	if vch.Endpoint != nil && vch.Endpoint.OperationsCredentials != nil {
		grantPerms = vch.Endpoint.OperationsCredentials.GrantPermissions
	}

	newConfig, err := validator.Validate(op, c.Data, false)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "cannot validate configuration: %s", err)
	}

	// copy changed config from new to old
	// TODO: [AngieCris] refactor this code to not depend on CLI (option: maybe move this changeConfig logic to vchConfig package)
	c.CopyChangedConf(config, newConfig, grantPerms)

	// add deprecated fields
	// TODO [AngieCris]: there're a bunch of other extra fields set separately
	vConfig := validator.AddDeprecatedFields(op, config, c.Data)
	vConfig.Timeout = c.Timeout
	vConfig.VCHSizeIsSet = c.ResourceLimits.IsSet

	// TODO [AngieCris]: handle affinity VM group thing separately

	// TODO [AngieCris]: set UpgradeInProgress flag, and RollBack logic

	err = hc.Executor().Configure(config, vConfig)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "failed to configure VCH: %s", err)
	}

	// TODO [AngieCris]: Question: what happens if configure failed and left a broken VCH? Is it the API's job to manually rollback to the original state? (or that's what rollback does)
	return nil
}
