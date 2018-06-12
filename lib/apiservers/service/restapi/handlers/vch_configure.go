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
	"sort"
	"time"

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

// TODO [AngieCris]: move it to some util package
const DefaultTimeout = 3 * time.Minute

type configureStrategy int

const (
	patch configureStrategy = iota
	put
)

type VCHPatch struct {
	vchConfigure
}

type VCHDatacenterPatch struct {
	vchConfigure
}

type VCHPut struct {
	vchConfigure
}

type VCHDatacenterPut struct {
	vchConfigure
}

// vchConfigure allows for VCHPatch and VCHDatacenterPatch to share common code without polluting the package
type vchConfigure struct{}

func (h *VCHPatch) Handle(params operations.PatchTargetTargetVchVchIDParams, principal interface{}) middleware.Responder {
	op := trace.FromContext(params.HTTPRequest.Context(), "VCHPatch: %s", params.VchID)

	b := target.Params{
		Target:     params.Target,
		Thumbprint: params.Thumbprint,
		VCHID:      &params.VchID,
	}

	task, err := h.handle(op, b, principal, params.Vch, patch)
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

	task, err := h.handle(op, b, principal, params.Vch, patch)
	if err != nil {
		return operations.NewPatchTargetTargetDatacenterDatacenterVchVchIDDefault(errors.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewPatchTargetTargetDatacenterDatacenterVchVchIDAccepted().WithPayload(operations.PatchTargetTargetDatacenterDatacenterVchVchIDAcceptedBody{Task: task})
}

func (h *VCHPut) Handle(params operations.PutTargetTargetVchVchIDParams, principal interface{}) middleware.Responder {
	op := trace.FromContext(params.HTTPRequest.Context(), "VCHPut: %s", params.VchID)

	b := target.Params{
		Target:     params.Target,
		Thumbprint: params.Thumbprint,
		VCHID:      &params.VchID,
	}

	task, err := h.handle(op, b, principal, params.Vch, put)
	if err != nil {
		return operations.NewPutTargetTargetVchVchIDDefault(errors.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewPutTargetTargetVchVchIDAccepted().WithPayload(operations.PutTargetTargetVchVchIDAcceptedBody{Task: task})
}

func (h *VCHDatacenterPut) Handle(params operations.PutTargetTargetDatacenterDatacenterVchVchIDParams, principal interface{}) middleware.Responder {
	op := trace.FromContext(params.HTTPRequest.Context(), "VCHDatacenterPut: %s", params.VchID)

	b := target.Params{
		Target:     params.Target,
		Thumbprint: params.Thumbprint,
		Datacenter: &params.Datacenter,
		VCHID:      &params.VchID,
	}

	task, err := h.handle(op, b, principal, params.Vch, put)
	if err != nil {
		return operations.NewPutTargetTargetDatacenterDatacenterVchVchIDDefault(errors.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewPutTargetTargetDatacenterDatacenterVchVchIDAccepted().WithPayload(operations.PutTargetTargetDatacenterDatacenterVchVchIDAcceptedBody{Task: task})
}

func (h *vchConfigure) handle(op trace.Operation, params target.Params, principal interface{}, vch *models.VCH, strategy configureStrategy) (*strfmt.URI, error) {
	// validate target
	d, hc, err := target.Validate(op, management.ActionConfigure, params, principal)
	if err != nil {
		return nil, err
	}

	switch strategy {
	case patch:
		// check for immutable fields, bail out if there is any
		err = h.checkImmutableVCHProperties(op, vch)
		if err != nil {
			return nil, err
		}
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

	switch strategy {
	case put:
		err := h.checkImmutableFieldsIfSet(op, oldData, c.Data)
		if err != nil {
			return nil, err
		}
	}

	// merge old and new data
	err = oldData.CopyNonEmpty(c.Data)
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "error copying new VCH data: %s", err)
	}

	c.Data = oldData

	// perform vch configure
	return nil, h.handleConfigure(op, c, vch, vchConfig, hc)
}

func (h *vchConfigure) mergeData(op trace.Operation, oldData *data.Data, newData *data.Data) (*data.Data, error) {
	err := h.checkImmutableFieldsIfSet(op, oldData, newData)
	if err != nil {
		return nil, err
	}

	err = oldData.CopyNonEmpty(newData)
	if err != nil {
		return nil, errors.NewError(http.StatusInternalServerError, "error copying new VCH data: %s", err)
	}

	// TODO [AngieCris]: there're also a bunch of other fields not copied by CopyNonEmpty, and treated separately (they're still mutable tho). Check with CLI code and figure out a plan

	return oldData, nil
}

// Check immutable fields in data struct, and output error whenever there is a mismatch
// only check the fields if it's set in source data (newData)
// TODO [AngieCris]: log specifics about what fields mismatch and what are the values (expected and actual) (also other helpful logs)
// TODO [AngieCris]: this manipulates data. Not sure if it should be moved out of handler. Not re-usable code for PATCH
func (h *vchConfigure) checkImmutableFieldsIfSet(op trace.Operation, oldData *data.Data, newData *data.Data) error {
	// vch name
	// TODO [AngieCris]: can't differentiate between unset and "". Now just treat "" as unset, even if empty string is provided
	if newData.DisplayName != "" && newData.DisplayName != oldData.DisplayName {
		return errors.NewError(http.StatusConflict, "Provided VCH name does not match with VCH configuration")
	}

	// networks: bridge, management, client, public
	if newData.BridgeNetworkName != "" {
		if newData.BridgeNetworkName != oldData.BridgeNetworkName || newData.BridgeIPRange.String() != oldData.BridgeIPRange.String() {
			return errors.NewError(http.StatusConflict, "Provided bridge network does not match with VCH configuration")
		}
	}

	if newData.ClientNetwork.IsSet() {
		if !checkNetworkConfig(oldData.ClientNetwork, newData.ClientNetwork) {
			return errors.NewError(http.StatusConflict, "provided client network does not match with VCH configuration")
		}
	}

	if newData.ManagementNetwork.IsSet() {
		if !checkNetworkConfig(oldData.ManagementNetwork, newData.ManagementNetwork) {
			return errors.NewError(http.StatusConflict, "provided management network does not match with VCH configuration")
		}
	}

	if newData.PublicNetwork.IsSet() {
		if !checkNetworkConfig(oldData.PublicNetwork, newData.PublicNetwork) {
			return errors.NewError(http.StatusConflict, "provided public network does not match with VCH configuration")
		}
	}

	// registry info
	if len(newData.WhitelistRegistries) != 0 && !checkTwoStringSlicesEqual(oldData.WhitelistRegistries, newData.WhitelistRegistries) {
		return errors.NewError(http.StatusConflict, "provided whitelist registries do not match with VCH configuration")
	}

	if len(newData.InsecureRegistries) != 0 && !checkTwoStringSlicesEqual(oldData.InsecureRegistries, newData.InsecureRegistries) {
		return errors.NewError(http.StatusConflict, "provided insecure registries do not match with VCH configuration")
	}

	// tls and ca
	// TODO [AngieCris]: skip for now because certFactory is stored configure.Configure instead of data, and it's processed later (according to CLI) (need to figure out why)

	// storage: datastore path
	if newData.ImageDatastorePath != "" && newData.ImageDatastorePath != oldData.ImageDatastorePath {
		return errors.NewError(http.StatusConflict, "provided image store path does not match with VCH configuration")
	}

	// endpoint VM resource limits
	// TODO [AngieCris]: check if NumCPUs and MemoryMB have default values if non provided from API (checking 0 works but..)
	if newData.NumCPUs != 0 && newData.NumCPUs != oldData.NumCPUs {
		return errors.NewError(http.StatusConflict, "provided number of CPUs reserved for VCH endpoint VM does not match with VCH configuration")
	}

	if newData.NumCPUs != 0 && newData.MemoryMB != oldData.MemoryMB {
		return errors.NewError(http.StatusConflict, "provided memory limit (in MB) reserved for VCH endpoint VM does not match with VCH configuration")
	}

	// syslog config
	if newData.SyslogConfig.IsSet() {
		if newData.SyslogConfig.Addr.String() != oldData.SyslogConfig.Addr.String() {
			return errors.NewError(http.StatusConflict, "provided syslog server address does not match with VCH configuration")
		}

		if newData.SyslogConfig.Tag != "" && newData.SyslogConfig.Tag != oldData.SyslogConfig.Tag {
			return errors.NewError(http.StatusConflict, "provided syslog server tag does not match with VCH configuration")
		}
	}

	// base image size
	if newData.ScratchSize != "" && newData.ScratchSize != oldData.ScratchSize {
		return errors.NewError(http.StatusConflict, "provided base image size does not match up with VCH configuration")
	}

	return nil
}

// TODO [AngieCris]: better log messages (more detailized)
func (h *vchConfigure) checkImmutableVCHProperties(op trace.Operation, vch *models.VCH) error {
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

func (h *vchConfigure) buildConfigure(op trace.Operation, d *data.Data, finder client.Finder, vch *models.VCH) (*configure.Configure, error) {
	c := &configure.Configure{Data: d}

	// parse input vch
	if vch != nil {
		if vch.Version != "" && version.String() != string(vch.Version) {
			return nil, errors.NewError(http.StatusBadRequest, "invalid version: %s", vch.Version)
		}

		c.DisplayName = vch.Name

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
func (h *vchConfigure) handleConfigure(op trace.Operation, c *configure.Configure, vch *models.VCH, config *config.VirtualContainerHostConfigSpec, hc *client.HandlerClient) error {
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

	// TODO [AngieCris]: patch affinity VM group thing separately

	// TODO [AngieCris]: set UpgradeInProgress flag, and RollBack logic

	err = hc.Executor().Configure(config, vConfig)
	if err != nil {
		return errors.NewError(http.StatusInternalServerError, "failed to configure VCH: %s", err)
	}

	// TODO [AngieCris]: Question: what happens if configure failed and left a broken VCH? Is it the API's job to manually rollback to the original state? (or that's what rollback does)
	return nil
}

// TODO [AngieCris]: this is a util function that needs to be moved to a util package
// check if two network configs match up. If not, returns false, if they are the same, return true
func checkNetworkConfig(oldConfig data.NetworkConfig, newConfig data.NetworkConfig) bool {

	if oldConfig.Name != newConfig.Name || oldConfig.IP.String() != newConfig.IP.String() || oldConfig.Gateway.String() != newConfig.Gateway.String() {
		return false
	}

	// check routing destinations
	// TODO [AngieCris]: this look up is n*2 time, needs to optimize
	for _, dest := range newConfig.Destinations {
		// check if this routing destination is included in the old config
		contains := false
		for _, oldDest := range oldConfig.Destinations {
			if oldDest.String() == dest.String() {
				contains = true
			}
		}
		if !contains {
			return false
		}
	}

	return true
}

// TODO [AngieCris]: this is a util function that needs to be moved to a util package
func checkTwoStringSlicesEqual(list1 []string, list2 []string) bool {
	if len(list1) != len(list2) {
		return false
	}

	sort.Strings(list1)
	sort.Strings(list2)

	for i, item := range list1 {
		if item != list2[i] {
			return false
		}
	}

	return true
}
