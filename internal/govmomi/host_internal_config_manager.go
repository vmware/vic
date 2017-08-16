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

package llpm

import (
	"context"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
)

// HostInternalConfigManager struct
type HostInternalConfigManager struct {
	object.Common
	c *vim25.Client

	Llpm *LowLevelProvisioningManager
}

// NewHostInternalConfigManager returns a new internal config manager for given host
func NewHostInternalConfigManager(ctx context.Context, client *vim25.Client, host *object.HostSystem) (*HostInternalConfigManager, error) {
	req := RetrieveInternalConfigManagerRequest{
		This: host.Reference(),
	}

	res, err := RetrieveInternalConfigManager(ctx, client, &req)
	if err != nil {
		return nil, err
	}

	return &HostInternalConfigManager{
		Common: object.NewCommon(client, res.Returnval.Self),
		c:      client,
		Llpm:   NewLowLevelProvisioningManager(client, res.Returnval.LlProvisioningManager),
	}, nil
}
