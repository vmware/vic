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
	"net/url"
	"path"

	"github.com/docker/docker/opts"
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/util"
	"github.com/vmware/vic/lib/apiservers/service/restapi/operations"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
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
	d := buildData(
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		nil,
		params.ComputeResource)

	vchs, err := handle(params.HTTPRequest.Context(), d)
	if err != nil {
		return operations.NewGetTargetTargetVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewGetTargetTargetVchOK().WithPayload(operations.GetTargetTargetVchOKBody{Vchs: vchs})
}

func (h *VCHDatacenterListGet) Handle(params operations.GetTargetTargetDatacenterDatacenterVchParams, principal interface{}) middleware.Responder {
	d := buildData(
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		&params.Datacenter,
		params.ComputeResource)

	vchs, err := handle(params.HTTPRequest.Context(), d)
	if err != nil {
		return operations.NewGetTargetTargetVchDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	return operations.NewGetTargetTargetVchOK().WithPayload(operations.GetTargetTargetVchOKBody{Vchs: vchs})
}

func handle(ctx context.Context, d *data.Data) ([]*models.VCHListItem, error) {
	validator, err := validateTarget(ctx, d)
	if err != nil {
		return nil, util.WrapError(400, err)
	}

	executor := management.NewDispatcher(validator.Context, validator.Session, nil, false)
	vchs, err := executor.SearchVCHs(validator.ClusterPath)
	if err != nil {
		return nil, util.NewError(500, fmt.Sprintf("Failed to search VCHs in %s: %s", validator.ResourcePoolPath, err))
	}

	return vchsToModels(ctx, vchs, executor), nil
}

func buildData(url url.URL, user string, pass string, thumbprint *string, datacenter *string, computeResource *string) *data.Data {
	d := data.Data{
		Target: &common.Target{
			URL:      &url,
			User:     user,
			Password: &pass,
		},
	}

	if datacenter != nil {
		// TODO: Convert ID to Name or update underlying code to accept ID
		d.Target.URL.Path = *datacenter
	}

	if thumbprint != nil {
		d.Thumbprint = *thumbprint
	}

	if computeResource != nil {
		d.ComputeResourcePath = *computeResource
	}

	return &d
}

func validateTarget(ctx context.Context, d *data.Data) (*validate.Validator, error) {
	if err := d.HasCredentials(); err != nil {
		return nil, fmt.Errorf("Invalid Credentials: %s", err)
	}

	validator, err := validate.NewValidator(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("Validation Error: %s", err)
	}
	// If dc is not set, and multiple datacenter is available, vic-machine ls will list VCHs under all datacenters.
	validator.AllowEmptyDC()

	_, err = validator.ValidateTarget(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("Target validation failed: %s", err)
	}
	_, err = validator.ValidateCompute(ctx, d, false)
	if err != nil {
		return nil, fmt.Errorf("Compute resource validation failed: %s", err)
	}

	return validator, nil
}

// Copied from list.go, and appears to be present other places. TODO: deduplicate
func upgradeStatusMessage(ctx context.Context, vch *vm.VirtualMachine, installerVer *version.Build, vchVer *version.Build) string {
	if sameVer := installerVer.Equal(vchVer); sameVer {
		return "Up to date"
	}

	upgrading, err := vch.VCHUpdateStatus(ctx)
	if err != nil {
		return fmt.Sprintf("Unknown: %s", err)
	}
	if upgrading {
		return "Upgrade in progress"
	}

	canUpgrade, err := installerVer.IsNewer(vchVer)
	if err != nil {
		return fmt.Sprintf("Unknown: %s", err)
	}
	if canUpgrade {
		return fmt.Sprintf("Upgradeable to %s", installerVer.ShortVersion())
	}

	oldInstaller, err := installerVer.IsOlder(vchVer)
	if err != nil {
		return fmt.Sprintf("Unknown: %s", err)
	}
	if oldInstaller {
		return fmt.Sprintf("VCH has newer version")
	}

	// can't get here
	return "Invalid upgrade status"
}

func vchsToModels(ctx context.Context, vchs []*vm.VirtualMachine, executor *management.Dispatcher) []*models.VCHListItem {
	installerVer := version.GetBuild()
	payload := make([]*models.VCHListItem, 0)
	for _, vch := range vchs {
		var version *version.Build
		var dockerHost string
		var adminPortal string
		if vchConfig, err := executor.GetNoSecretVCHConfig(vch); err == nil {
			version = vchConfig.Version

			if public := vchConfig.ExecutorConfig.Networks["public"]; public != nil {
				if public_ip := public.Assigned.IP; public_ip != nil {
					var docker_port = opts.DefaultTLSHTTPPort
					if vchConfig.HostCertificate.IsNil() {
						docker_port = opts.DefaultHTTPPort
					}

					dockerHost = fmt.Sprintf("%s:%d", public_ip, docker_port)
					adminPortal = fmt.Sprintf("https://%s:2378", public_ip)
				}
			}
		}

		name := path.Base(vch.InventoryPath)

		model := &models.VCHListItem{ID: vch.Reference().Value, Name: name, AdminPortal: adminPortal, DockerHost: dockerHost}

		if version != nil {
			model.Version = version.ShortVersion()
			model.UpgradeStatus = upgradeStatusMessage(ctx, vch, installerVer, version)
		}

		payload = append(payload, model)
	}

	return payload
}
