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

package decode

import (
	"net/http"
	"fmt"

	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/client"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/lib/constants"
	"github.com/vmware/vic/cmd/vic-machine/common"
)

func ProcessCompute(op trace.Operation, d *data.Data, vch *models.VCH, finder client.Finder) error {
	if vch.Compute != nil {
		if vch.Compute.CPU != nil {
			d.ResourceLimits.VCHCPULimitsMHz = MHzFromValueHertz(vch.Compute.CPU.Limit)
			d.ResourceLimits.VCHCPUReservationsMHz = MHzFromValueHertz(vch.Compute.CPU.Reservation)
			d.ResourceLimits.VCHCPUShares = FromShares(vch.Compute.CPU.Shares)
		}

		if vch.Compute.Memory != nil {
			d.ResourceLimits.VCHMemoryLimitsMB = MBFromValueBytes(vch.Compute.Memory.Limit)
			d.ResourceLimits.VCHMemoryReservationsMB = MBFromValueBytes(vch.Compute.Memory.Reservation)
			d.ResourceLimits.VCHMemoryShares = FromShares(vch.Compute.Memory.Shares)
		}

		// TODO (#6711): Do we need to handle clusters differently?
		resourcePath, err := FromManagedObject(op, finder, "ResourcePool", vch.Compute.Resource)
		if err != nil {
			errors.NewError(http.StatusBadRequest, "error finding resource pool: %s", err)
		}
		if resourcePath == "" {
			return errors.NewError(http.StatusBadRequest, "resource pool must be specified (by name or id)")
		}
		d.ComputeResourcePath = resourcePath

		if vch.Compute.Affinity != nil {
			d.UseVMGroup = vch.Compute.Affinity.UseVMGroup
		}
	}

	return nil
}

func ProcessEndpoint(op trace.Operation, d *data.Data, vch *models.VCH) error {
	d.MemoryMB = constants.DefaultEndpointMemoryMB

	if vch.Endpoint != nil {
		if vch.Endpoint.Memory != nil {
			d.MemoryMB = *MBFromValueBytes(vch.Endpoint.Memory)
		}
		if vch.Endpoint.CPU != nil {
			d.NumCPUs = int(vch.Endpoint.CPU.Sockets)
		}

		if vch.Endpoint.OperationsCredentials != nil {
			opsPassword := string(vch.Endpoint.OperationsCredentials.Password)
			d.OpsCredentials = common.OpsCredentials{
				OpsUser:     &vch.Endpoint.OperationsCredentials.User,
				OpsPassword: &opsPassword,
				GrantPerms:  &vch.Endpoint.OperationsCredentials.GrantPermissions,
			}
		}
	}

	err := processOpsCredentials(op, d.OpsCredentials, d.Target.User, d.Target.Password)
	if err != nil {
		return errors.WrapError(http.StatusBadRequest, err)
	}

	return nil
}

// processOpsCredentials check if the ops credentials given is valid. If not, use administrative user as ops user
func processOpsCredentials(op trace.Operation, opsCreds common.OpsCredentials, adminUser string, adminPassword *string) error {
	if opsCreds.OpsUser == nil {
		if opsCreds.GrantPerms != nil {
			if *opsCreds.GrantPerms { // grant perms set but no ops user
				return fmt.Errorf("no ops user to grant permissions to")
			}
			opsCreds.GrantPerms = nil
		}
		op.Warn("Ops credentials not specified. Using administrative user for VCH operation")
		opsCreds.OpsUser = &adminUser
		// TODO [AngieCris]: is there the need to check if adminPassword is nil? (passed from target, shouldn't be nil
		opsCreds.OpsPassword = adminPassword
	}

	return nil
}
