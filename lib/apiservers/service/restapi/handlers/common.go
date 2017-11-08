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
	"net/http"
	"net/url"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/util"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

type buildDataParams struct {
	target          string
	thumbprint      *string
	datacenter      *string
	computeResource *string
}

func buildData(ctx context.Context, params buildDataParams, principal interface{}) (*data.Data, error) {
	d := data.Data{
		Target: &common.Target{
			URL: &url.URL{Host: params.target},
		},
	}

	if c, ok := principal.(Credentials); ok {
		d.Target.User = c.user
		d.Target.Password = &c.pass
	} else {
		d.Target.CloneTicket = principal.(Session).ticket
	}

	if params.thumbprint != nil {
		d.Thumbprint = *params.thumbprint
	}

	if params.datacenter != nil {
		validator, err := validateTarget(ctx, &d)
		if err != nil {
			return nil, util.WrapError(http.StatusInternalServerError, err)
		}

		datacenterManagedObjectReference := types.ManagedObjectReference{Type: "Datacenter", Value: *params.datacenter}

		datacenterObject, err := validator.Session.Finder.ObjectReference(ctx, datacenterManagedObjectReference)
		if err != nil {
			return nil, util.WrapError(http.StatusNotFound, err)
		}

		d.Target.URL.Path = datacenterObject.(*object.Datacenter).InventoryPath
	}

	if params.computeResource != nil {
		d.ComputeResourcePath = *params.computeResource
	}

	return &d, nil
}

func validateTarget(ctx context.Context, d *data.Data) (*validate.Validator, error) {
	if err := d.HasCredentials(); err != nil {
		return nil, fmt.Errorf("Invalid Credentials: %s", err)
	}

	validator, err := validate.NewValidator(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("Validation Error: %s", err)
	}

	// If dc is not set, and multiple datacenters are available, operate on all datacenters.
	validator.AllowEmptyDC()

	if _, err = validator.ValidateTarget(ctx, d); err != nil {
		return nil, fmt.Errorf("Target validation failed: %s", err)
	}

	if _, err = validator.ValidateCompute(ctx, d, false); err != nil {
		return nil, fmt.Errorf("Compute resource validation failed: %s", err)
	}

	return validator, nil
}

// Copied from list.go, and appears to be present other places. TODO (#6032): deduplicate
func upgradeStatusMessage(op trace.Operation, vch *vm.VirtualMachine, installerVer *version.Build, vchVer *version.Build) string {
	if sameVer := installerVer.Equal(vchVer); sameVer {
		return "Up to date"
	}

	upgrading, err := vch.VCHUpdateStatus(op)
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

func getVCHConfig(op trace.Operation, d *data.Data) (*config.VirtualContainerHostConfigSpec, error) {
	// TODO (#6032): abstract some of this boilerplate into helpers
	validator, err := validateTarget(op.Context, d)
	if err != nil {
		return nil, util.WrapError(http.StatusBadRequest, err)
	}

	executor := management.NewDispatcher(validator.Context, validator.Session, nil, false)
	vch, err := executor.NewVCHFromID(d.ID)
	if err != nil {
		return nil, util.NewError(http.StatusNotFound, fmt.Sprintf("Unable to find VCH %s: %s", d.ID, err))
	}

	err = validate.SetDataFromVM(validator.Context, validator.Session.Finder, vch, d)
	if err != nil {
		return nil, util.NewError(http.StatusInternalServerError, fmt.Sprintf("Failed to load VCH data: %s", err))
	}

	vchConfig, err := executor.GetNoSecretVCHConfig(vch)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve VCH information: %s", err)
	}

	return vchConfig, nil
}
