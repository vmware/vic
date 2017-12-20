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
	"fmt"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	gvsession "github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/session"
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

type NameToRef map[string]types.ManagedObjectReference

type AuthzManager struct {
	authzManager *object.AuthorizationManager
	client       *vim25.Client
	configSpec   *config.VirtualContainerHostConfigSpec
	principal    string
	rbacConfig   *RBACConfig
	resources    map[int8]*RBACResource
	targetRoles  []types.AuthorizationRole
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

type RBACResourcePermission struct {
	rType      int8
	reference  types.ManagedObjectReference
	permission types.Permission
}

func NewAuthzManager(ctx context.Context, client *vim25.Client, configSpec *config.VirtualContainerHostConfigSpec) *AuthzManager {
	authManager := object.NewAuthorizationManager(client)
	mgr := &AuthzManager{
		client:       client,
		authzManager: authManager,
		configSpec:   configSpec,
	}
	return mgr
}

func GrantOpsUserPerms(ctx context.Context, client *vim25.Client, configSpec *config.VirtualContainerHostConfigSpec) error {
	am := NewAuthzManager(ctx, client, configSpec)
	am.configSpec = configSpec
	am.InitRBACConfig(configSpec.Connection.Username, &OpsUserRBACConf)
	_, err := am.SetupRolesAndPermissions(ctx)
	return err
}

func (am *AuthzManager) InitRBACConfig(principal string, config *RBACConfig) {
	am.principal = principal
	am.rbacConfig = config
	am.initTargetRoles()
	am.initResourceMap()
}

func (am *AuthzManager) CreateRoles(ctx context.Context) (int, error) {
	return am.createOrRepairRoles(ctx)
}

func (am *AuthzManager) DeleteRoles(ctx context.Context) (int, error) {
	return am.deleteRoles(ctx)
}

func (am *AuthzManager) RoleList(ctx context.Context) (object.AuthorizationRoleList, error) {
	return am.getRoleList(ctx)
}

func (am *AuthzManager) SetupRolesAndPermissions(ctx context.Context) ([]RBACResourcePermission, error) {
	res, err := am.isPrincipalAnAdministrator(ctx)
	if err != nil {
		return nil, err
	}
	if res {
		log.Warnf("Skipping ops-user Role/Permissions initialization. The current ops-user (%s) has administrative privileges.", am.principal)
		log.Warnf("This occurs when \"%s\" is a member of the \"Administrators\" group or has been granted \"Admin\" role to any of the resources in the system.", am.principal)
		return nil, nil
	}
	if _, err = am.CreateRoles(ctx); err != nil {
		return nil, err
	}
	return am.SetupPermissions(ctx)
}

func (am *AuthzManager) SetupPermissions(ctx context.Context) ([]RBACResourcePermission, error) {
	return am.setupPermissions(ctx)
}

func (am *AuthzManager) setupPermissions(ctx context.Context) ([]RBACResourcePermission, error) {
	type ResourceDesc struct {
		rType int8
		ref   types.ManagedObjectReference
	}

	resourceDescs := make([]ResourceDesc, 0, len(am.rbacConfig.Resources))

	// Get a reference to the top object
	finder := find.NewFinder(am.client, false)

	root, err := finder.Folder(ctx, "/")
	if err != nil {
		return nil, errors.Errorf("Ops-User: AuthzManager, Unable to find top object: %s", err.Error())
	}

	resourceDescs = append(resourceDescs, ResourceDesc{VCenter, root.Reference()})

	session := session.NewSession(&session.Config{})
	// Set client
	session.Client = &govmomi.Client{
		Client:         am.client,
		SessionManager: gvsession.NewManager(am.client),
	}

	// Use the VirtualContainerHostConfigSpec to find the various resources
	// Start with Resource Pool, Cluster and Datacenter
	rpRef := am.configSpec.ComputeResources[0]
	rp := compute.NewResourcePool(ctx, session, rpRef)

	datacenter, err := rp.GetDatacenter(ctx)
	if err != nil {
		return nil, errors.Errorf("Ops-User: AuthzManager, Unable to find Datacenter: %s", err.Error())
	}
	resourceDescs = append(resourceDescs, ResourceDesc{Datacenter, datacenter.Reference()})

	finder.SetDatacenter(datacenter)

	cluster, err := rp.GetCluster(ctx)
	if err != nil {
		return nil, errors.Errorf("Ops-User: AuthzManager, Unable to find Cluster: %s", err.Error())
	}
	resourceDescs = append(resourceDescs, ResourceDesc{Cluster, cluster.Reference()})

	// Find image datastore
	dsNameToRef := make(NameToRef)
	err = am.collectDatastores(ctx, finder, dsNameToRef)
	if err != nil {
		return nil, errors.Errorf("Ops-User: AuthzManager, Unable to find Datastores: %s", err.Error())
	}

	// Loop over Datastores
	for _, ref := range dsNameToRef {
		resourceDescs = append(resourceDescs, ResourceDesc{Datastore, ref})
	}

	// Loop over Networks
	for _, network := range am.configSpec.Network.ContainerNetworks {
		netRef := &types.ManagedObjectReference{}
		netRef.FromString(network.ID)
		if netRef.Type == "" || netRef.Value == "" {
			return nil, errors.Errorf("Ops-User: AuthzManager, Unable to build Bridged Network MoRef: %s", network.ID)
		}
		resourceDescs = append(resourceDescs, ResourceDesc{Network, *netRef})
	}

	// Loop over Resource Pools
	for _, rPoolRef := range am.configSpec.ComputeResources {
		resourceDescs = append(resourceDescs, ResourceDesc{Endpoint, rPoolRef})
	}

	resourcePermissions := make([]RBACResourcePermission, 0, len(am.rbacConfig.Resources))
	// Apply permissions
	for _, desc := range resourceDescs {
		resourcePermission, err := am.addPermission(ctx, desc.ref, desc.rType, false)
		if err != nil {
			return nil, errors.Errorf("Ops-User: AuthzManager, Unable to set permissions on %s, error: %s",
				desc.ref.String(), err.Error())
		}
		if resourcePermission != nil {
			resourcePermissions = append(resourcePermissions, *resourcePermission)
		}
	}

	return resourcePermissions, nil
}

func (am *AuthzManager) isPrincipalAnAdministrator(ctx context.Context) (bool, error) {
	// Check if the principal belongs to the Administrators group
	res, err := am.principalBelongsToGroup(ctx, "Administrators")
	if err != nil {
		return false, err
	}

	if res {
		return res, nil
	}

	// Check if the principal has an Admin Role
	res, err = am.principalHasRole(ctx, "Admin")
	if err != nil {
		return false, err
	}

	return res, nil
}

func (am *AuthzManager) principalBelongsToGroup(ctx context.Context, group string) (bool, error) {
	ref := *am.client.ServiceContent.UserDirectory

	components := strings.Split(am.principal, "@")
	var domain string
	name := components[0]
	if len(components) < 2 {
		domain = ""
	} else {
		domain = components[1]
	}

	req := types.RetrieveUserGroups{
		This:           ref,
		Domain:         domain,
		SearchStr:      name,
		ExactMatch:     true,
		BelongsToGroup: group,
		FindUsers:      true,
		FindGroups:     false,
	}

	results, err := methods.RetrieveUserGroups(ctx, am.client, &req)
	if err != nil {
		return false, err
	}

	if len(results.Returnval) > 0 {
		return true, nil
	}

	return false, nil
}

func (am *AuthzManager) principalHasRole(ctx context.Context, roleName string) (bool, error) {
	// Build expected representation of the ops-user
	principal := strings.ToLower(am.principal)

	// Get role id for admin Role
	roleList, err := am.RoleList(ctx)
	if err != nil {
		return false, err
	}

	role := roleList.ByName(roleName)

	allPerms, err := am.authzManager.RetrieveAllPermissions(ctx)
	if err != nil {
		return false, err
	}

	for _, perm := range allPerms {
		if perm.RoleId != role.RoleId {
			continue
		}

		fPrincipal := am.formatPrincipal(perm.Principal)
		if fPrincipal == principal {
			return true, nil
		}
	}

	return false, nil
}

func (am *AuthzManager) getPermissions(ctx context.Context,
	ref types.ManagedObjectReference) ([]types.Permission, error) {
	// Get current Permissions
	return am.authzManager.RetrieveEntityPermissions(ctx, ref, false)
}

func (am *AuthzManager) addPermission(ctx context.Context, ref types.ManagedObjectReference,
	resourceType int8, isGroup bool) (*RBACResourcePermission, error) {

	resource := am.getResource(resourceType)
	if resource == nil {
		return nil, fmt.Errorf("cannot find resource of type %d", resourceType)
	}

	// Collect the new roles, possibly cache the result in the Authz manager
	roleList, err := am.getRoleList(ctx)
	if err != nil {
		return nil, err
	}

	// Locate target role
	role := roleList.ByName(rolePrefix + resource.Role.Name)
	if role == nil {
		return nil, fmt.Errorf("cannot find role: %s", resource.Role.Name)
	}

	// Get current Permissions
	permissions, err := am.authzManager.RetrieveEntityPermissions(ctx, ref, false)
	if err != nil {
		return nil, err
	}

	for _, permission := range permissions {
		if permission.Principal == am.principal &&
			permission.RoleId == role.RoleId &&
			permission.Propagate == resource.Propagate {
			return nil, nil
		}
	}

	// No match found, create new permission
	permission := types.Permission{
		Principal: am.principal,
		RoleId:    role.RoleId,
		Propagate: resource.Propagate,
		Group:     isGroup,
	}

	permissions = append(permissions, permission)

	if err = am.authzManager.SetEntityPermissions(ctx, ref, permissions); err != nil {
		return nil, err
	}

	resourcePermission := &RBACResourcePermission{
		permission: permission,
		reference:  ref,
		rType:      resourceType,
	}

	return resourcePermission, nil
}

func (am *AuthzManager) collectDatastores(ctx context.Context, finder *find.Finder, dsNameToRef NameToRef) error {
	err := am.findDatastores(ctx, finder, am.configSpec.Storage.ImageStores, dsNameToRef)
	if err != nil {
		return err
	}
	volumeLocations := make([]url.URL, 0, len(am.configSpec.Storage.VolumeLocations))
	for _, volumeLocation := range am.configSpec.Storage.VolumeLocations {
		volumeLocations = append(volumeLocations, *volumeLocation)
	}
	if err = am.findDatastores(ctx, finder, volumeLocations, dsNameToRef); err != nil {
		return err
	}
	return nil
}

func (am *AuthzManager) findDatastores(ctx context.Context, finder *find.Finder,
	storeURLs []url.URL, dsNameToRef NameToRef) error {
	for _, storeURL := range storeURLs {
		dsName := storeURL.Host
		// Skip if we already have one
		if _, ok := dsNameToRef[dsName]; ok {
			continue
		}
		ds, err := finder.Datastore(ctx, dsName)
		if err != nil {
			return err
		}
		dsNameToRef[dsName] = ds.Reference()
	}
	return nil
}

func (am *AuthzManager) createOrRepairRoles(ctx context.Context) (int, error) {
	// Get all the existing roles
	mgr := am.authzManager
	roleList, err := mgr.RoleList(ctx)
	if err != nil {
		return 0, err
	}

	var count int
	for _, targetRole := range am.targetRoles {
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

func (am *AuthzManager) deleteRoles(ctx context.Context) (int, error) {
	mgr := am.authzManager
	// Get all the existing roles
	roleList, err := mgr.RoleList(ctx)
	if err != nil {
		return 0, err
	}

	var count int
	for _, targetRole := range am.targetRoles {
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

func (am *AuthzManager) initTargetRoles() {
	count := len(am.rbacConfig.Resources)
	roles := make([]types.AuthorizationRole, 0, count)
	dSet := make(map[string]bool)
	for _, resource := range am.rbacConfig.Resources {
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
	am.targetRoles = roles
}

func (am *AuthzManager) initResourceMap() {
	am.resources = make(map[int8]*RBACResource)
	for i, resource := range am.rbacConfig.Resources {
		am.resources[resource.Type] = &am.rbacConfig.Resources[i]
	}
}

func (am *AuthzManager) getResource(resourceType int8) *RBACResource {
	resource, ok := am.resources[resourceType]
	if !ok {
		panic(errors.Errorf("Cannot find RBAC resource type: %d", resourceType))
	}
	return resource
}

func (am *AuthzManager) formatPrincipal(principal string) string {
	components := strings.Split(principal, "\\")
	if len(components) != 2 {
		return strings.ToLower(principal)
	}
	ret := strings.ToLower(components[1]) + "@" + strings.ToLower(components[0])
	return ret
}
