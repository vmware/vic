// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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

package management

import (
	"context"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
)

const (
	VCenter = iota
	Datacenter
	Cluster
	DatastoreFolder
	Datastore
	VSANDatastore
	Network
	Endpoint
)

type AuthzManager struct {
	authzManager *object.AuthorizationManager
	client       *vim25.Client
	vchConfig    *config.VirtualContainerHostConfigSpec
	principal    string
}

type RBACResource struct {
	Type      int8
	Propagate bool
	Role      types.AuthorizationRole
}

type RBACConfig struct {
	Resources []RBACResource
}

type PermissionList []types.Permission

var rolePrefix = "vic-vch-"

var RoleVCenter = types.AuthorizationRole{
	Name: "vcenter",
	Privilege: []string{
		"Datastore.Config",
	},
}

var RoleDataCenter = types.AuthorizationRole{
	Name: "datacenter",
	Privilege: []string{
		"Datastore.Config",
		"Datastore.FileManagement",
	},
}

var RoleDataStore = types.AuthorizationRole{
	Name: "datastore",
	Privilege: []string{
		"Datastore.AllocateSpace",
		"Datastore.Browse",
		"Datastore.Config",
		"Datastore.DeleteFile",
		"Datastore.FileManagement",
		"Host.Config.SystemManagement",
	},
}

var RoleNetwork = types.AuthorizationRole{
	Name: "network",
	Privilege: []string{
		"Network.Assign",
	},
}

var RoleEndpoint = types.AuthorizationRole{
	Name: "endpoint",
	Privilege: []string{
		"DVPortgroup.Create",
		"DVPortgroup.Delete",
		"DVPortgroup.Modify",
		"DVPortgroup.PolicyOp",
		"DVPortgroup.ScopeOp",
		"Resource.AssignVMToPool",
		"VirtualMachine.Config.AddNewDisk",
		"VirtualMachine.Config.AdvancedConfig",
		"VirtualMachine.Config.EditDevice",
		"VirtualMachine.Config.RemoveDisk",
		"VirtualMachine.GuestOperations.Execute",
		"VirtualMachine.Interact.DeviceConnection",
		"VirtualMachine.Interact.PowerOff",
		"VirtualMachine.Interact.PowerOn",
		"VirtualMachine.Inventory.Create",
		"VirtualMachine.Inventory.Delete",
		"VirtualMachine.Inventory.Register",
		"VirtualMachine.Inventory.Unregister",
	},
}

var RBACConf = RBACConfig{
	Resources: []RBACResource{
		{
			Type:      VCenter,
			Propagate: false,
			Role:      RoleVCenter,
		},
		{
			Type:      Datacenter,
			Propagate: false,
			Role:      RoleDataCenter,
		},
		{
			Type:      Cluster,
			Propagate: true,
			Role:      RoleDataStore,
		},
		{
			Type:      DatastoreFolder,
			Propagate: true,
			Role:      RoleDataStore,
		},
		{
			Type:      Datastore,
			Propagate: true,
			Role:      RoleDataStore,
		},
		{
			Type:      VSANDatastore,
			Propagate: false,
			Role:      RoleDataStore,
		},
		{
			Type:      Network,
			Propagate: false,
			Role:      RoleNetwork,
		},
		{
			Type:      Endpoint,
			Propagate: true,
			Role:      RoleEndpoint,
		},
	},
}

func NewAuthzManager(ctx context.Context, client *vim25.Client, principal string) *AuthzManager {
	authManager := object.NewAuthorizationManager(client)
	mgr := &AuthzManager{
		client:       client,
		authzManager: authManager,
		principal:    principal,
	}
	return mgr
}

func (am *AuthzManager) CreateRoles(ctx context.Context) (int, error) {
	return am.createOrRepairRoles(ctx, getTargetRoles())
}

func (am *AuthzManager) DeleteRoles(ctx context.Context) (int, error) {
	return am.deleteRoles(ctx, getTargetRoles())
}

func (am *AuthzManager) RoleList(ctx context.Context) (object.AuthorizationRoleList, error) {
	return am.getRoleList(ctx)
}

func (am *AuthzManager) createOrRepairRoles(ctx context.Context, targetRoles []types.AuthorizationRole) (int, error) {
	// Get all the existing roles
	mgr := am.authzManager
	roleList, err := mgr.RoleList(ctx)
	if err != nil {
		return 0, err
	}

	var count int
	for _, targetRole := range targetRoles {
		foundRole := roleList.ByName(targetRole.Name)
		if foundRole != nil {
			isMod, err := am.checkAndRepairRole(ctx, &targetRole, foundRole)
			if isMod && err == nil {
				count++
			}
		} else {
			_, err = mgr.AddRole(ctx, targetRole.Name, targetRole.Privilege)
			if err == nil {
				count++
			}
		}
		if err != nil {
			return count, err
		}
	}
	return count, nil
}

func (am *AuthzManager) deleteRoles(ctx context.Context, targetRoles []types.AuthorizationRole) (int, error) {
	mgr := am.authzManager
	// Get all the existing roles
	roleList, err := mgr.RoleList(ctx)
	if err != nil {
		return 0, err
	}

	var count int
	for _, targetRole := range targetRoles {
		foundRole := roleList.ByName(targetRole.Name)
		if foundRole != nil {
			err = mgr.RemoveRole(ctx, foundRole.RoleId, true)
			if err == nil {
				count++
			}
		}
	}
	return count, nil
}

func (am *AuthzManager) getRoleList(ctx context.Context) (object.AuthorizationRoleList, error) {
	return am.authzManager.RoleList(ctx)
}

func (am *AuthzManager) checkAndRepairRole(ctx context.Context, tRole *types.AuthorizationRole, fRole *types.AuthorizationRole) (bool, error) {
	mgr := am.authzManager
	// Check that the privileges list in Target Role is a subset of the list in Found role
	fSet := make(map[string]bool)
	for _, p := range fRole.Privilege {
		fSet[p] = true
	}

	var isModified bool
	for _, p := range tRole.Privilege {
		if _, found := fSet[p]; !found {
			// Privilege not found
			// Add it to the found Role
			fRole.Privilege = append(fRole.Privilege, p)
			isModified = true
		}
	}

	if !isModified {
		return false, nil
	}

	// Not a subset need to call go-vmomi to set the new privileges
	err := mgr.UpdateRole(ctx, fRole.RoleId, fRole.Name, fRole.Privilege)

	return true, err
}

func getTargetRoles() []types.AuthorizationRole {
	count := len(RBACConf.Resources)
	roles := make([]types.AuthorizationRole, 0, count)
	dSet := make(map[string]bool)
	for _, resource := range RBACConf.Resources {
		name := rolePrefix + resource.Role.Name
		// Discard duplicates
		if _, found := dSet[name]; !found {
			role := new(types.AuthorizationRole)
			*role = resource.Role
			role.Name = name
			dSet[name] = true
			roles = append(roles, *role)
		}
	}
	return roles
}

func getResource(resourceType int8) *RBACResource {
	for _, resource := range RBACConf.Resources {
		if resource.Type == resourceType {
			return &resource
		}
	}
	return nil
}
