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
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/test/env"
)

func newOpsUserAuthzManager(ctx context.Context, client *vim25.Client, configSpec *config.VirtualContainerHostConfigSpec) *AuthzManager {
	am := NewAuthzManager(ctx, client, configSpec)
	am.InitRBACConfig("ops-user@vsphere.local", &OpsUserRBACConf)
	return am
}

func TestOpsUserRolesSimulatorVPX(t *testing.T) {
	ctx := context.Background()
	m := simulator.VPX()
	defer m.Remove()

	err := m.Create()
	require.NoError(t, err, "Cannot create VPX Simulator")

	s := m.Service.NewServer()
	defer s.Close()

	config := &session.Config{
		Service:   s.URL.String(),
		Insecure:  true,
		Keepalive: time.Duration(5) * time.Minute,
	}

	sess, err := session.NewSession(config).Connect(ctx)
	require.NoError(t, err, "Cannot connect to VPX Simulator")

	am := newOpsUserAuthzManager(ctx, sess.Vim25(), nil)

	doTestRoles(ctx, t, am)
}

func TestOpsUserRolesVCenter(t *testing.T) {
	ctx := context.Background()

	config := &session.Config{
		Service:   env.URL(t),
		Insecure:  true,
		Keepalive: time.Duration(5) * time.Minute,
	}

	sess, err := session.NewSession(config).Connect(ctx)
	if err != nil {
		t.SkipNow()
	}

	am := newOpsUserAuthzManager(ctx, sess.Vim25(), nil)

	doTestRoles(ctx, t, am)
}

func TestOpsUserPermsSimulatorVPX(t *testing.T) {
	ctx := context.Background()
	m := simulator.VPX()

	defer m.Remove()

	err := m.Create()
	require.NoError(t, err)

	s := m.Service.NewServer()
	defer s.Close()

	fmt.Println(s.URL.String())

	config := &session.Config{
		Service:   s.URL.String(),
		Insecure:  true,
		Keepalive: time.Duration(5) * time.Minute,
	}

	sess, err := session.NewSession(config).Connect(ctx)
	require.NoError(t, err)

	am := newOpsUserAuthzManager(ctx, sess.Vim25(), nil)

	var roleCount = len(am.targetRoles)
	count := initRoles(ctx, t, am)

	defer cleanup(ctx, t, am, true)
	require.Equal(t, roleCount, count, "Incorrect number of roles: expected %d, actual %d", roleCount, count)

	c := sess.Client
	// Find the Datacenter
	finder := find.NewFinder(c.Client, false)

	dcList, err := finder.DatacenterList(ctx, "/*")
	require.NoError(t, err)
	require.NotEqual(t, 0, len(dcList))

	dc := dcList[0]
	finder.SetDatacenter(dc)

	resourcePermission, err := am.addPermission(ctx, dc.Reference(), Datacenter, false)
	require.NoError(t, err)
	require.NotNil(t, resourcePermission)

	// Get permission back
	permissions, err := am.getPermissions(ctx, dc.Reference())

	if err != nil || len(permissions) == 0 {
		t.Fatalf("Failed to get permissions for Datacenter")
	}

	foundPermission := permissions[0]

	permission := &resourcePermission.permission

	if foundPermission.Principal != permission.Principal ||
		foundPermission.RoleId != permission.RoleId ||
		foundPermission.Propagate != permission.Propagate ||
		foundPermission.Group != permission.Group {
		t.Fatalf("Permission mismatch, exp: %v, found: %v", permission, foundPermission)
	}
}

func TestOpsUserPermsFromConfigSimulatorVPX(t *testing.T) {
	ctx := context.Background()
	m := simulator.VPX()

	m.Datacenter = 3
	m.Folder = 2
	m.Pool = 1
	m.App = 1
	m.Pod = 1

	defer m.Remove()

	err := m.Create()
	require.NoError(t, err)

	s := m.Service.NewServer()
	defer s.Close()

	fmt.Println(s.URL.String())

	config := &session.Config{
		Service:   s.URL.String(),
		Insecure:  true,
		Keepalive: time.Duration(5) * time.Minute,
	}
	sess, err := session.NewSession(config).Connect(ctx)
	require.NoError(t, err)

	configSpec, err := validate.TestValidateForSim(ctx, s.URL)
	require.NoError(t, err)

	// Set up the Authz Manager
	am := newOpsUserAuthzManager(ctx, sess.Vim25(), configSpec)

	resourcePermissions, err := am.SetupRolesAndPermissions(ctx)
	require.NoError(t, err)
	defer cleanup(ctx, t, am, true)
	require.True(t, len(am.rbacConfig.Resources) >= len(resourcePermissions))

	verifyResourcePermissions(ctx, t, am, resourcePermissions, am.rbacConfig.Resources)
}

func verifyResourcePermissions(ctx context.Context, t *testing.T, am *AuthzManager, retPerms []RBACResourcePermission, configResources []RBACResource) {
	for _, retPerm := range retPerms {

		// Validate returned permission against the configured permission
		configPerm := am.getResource(retPerm.rType)
		require.Equal(t, am.principal, retPerm.permission.Principal)
		require.Equal(t, configPerm.Propagate, retPerm.permission.Propagate)

		actPerms, err := am.getPermissions(ctx, retPerm.reference)
		require.NoError(t, err)

		for _, actPerm := range actPerms {
			if actPerm.Principal != am.principal {
				continue
			}
			// RoleId must be the same
			require.Equal(t, retPerm.permission.RoleId, actPerm.RoleId)
		}
	}
}

func initRoles(ctx context.Context, t *testing.T, am *AuthzManager) int {
	cleanup(ctx, t, am, false)

	count, err := am.createOrRepairRoles(ctx)
	require.NoError(t, err, "Failed to initialize Roles")

	return count
}

func doTestRoles(ctx context.Context, t *testing.T, am *AuthzManager) {
	var roleCount = len(am.targetRoles)
	count := initRoles(ctx, t, am)

	defer cleanup(ctx, t, am, true)
	require.Equal(t, roleCount, count, "Incorrect number of roles: expected %d, actual %d", roleCount, count)

	// Test correct role validation, it should return 0
	roleCount = 0
	count, err := am.createOrRepairRoles(ctx)
	require.NoError(t, err, "Failed to create roles")
	require.Equal(t, roleCount, count, "Incorrect number of roles: expected %d, actual %d", roleCount, count)

	// Remove two Privileges from two roles
	roles, err := am.getRoleList(ctx)
	fmt.Println(err)
	fmt.Println(roles)

	targetRoleName1 := rolePrefix + "datastore"
	targetRoleName2 := rolePrefix + "endpoint"

	for _, role := range roles {
		if role.Name == targetRoleName1 {
			removePrivilege(&role, "Datastore.DeleteFile")
			am.authzManager.UpdateRole(ctx, role.RoleId, role.Name, role.Privilege)
		}
		if role.Name == targetRoleName2 {
			removePrivilege(&role, "VirtualMachine.Config.AddNewDisk")
			am.authzManager.UpdateRole(ctx, role.RoleId, role.Name, role.Privilege)
		}
	}

	// Test
	roleCount = 2
	count, err = am.createOrRepairRoles(ctx)
	require.NoError(t, err, "Failed to repair roles 1")
	require.Equal(t, roleCount, count, "Incorrect number of roles: expected %d, actual %d", roleCount, count)

	// Test correct role validation, it should return 0
	roleCount = 0
	count, err = am.createOrRepairRoles(ctx)
	require.NoError(t, err, "Failed to repair roles 2")
	require.Equal(t, roleCount, count, "Incorrect number of roles: expected %d, actual %d", roleCount, count)
}

func removePrivilege(role *types.AuthorizationRole, privilege string) {
	for i, priv := range role.Privilege {
		if priv == privilege {
			role.Privilege = append(role.Privilege[:i], role.Privilege[i+1:]...)
			return
		}
	}
}

func cleanup(ctx context.Context, t *testing.T, am *AuthzManager, checkCount bool) {
	var roleCount = len(am.targetRoles)
	count, err := am.deleteRoles(ctx)
	require.NoError(t, err, "Failed to delete roles")

	if checkCount && count != roleCount {
		t.Fatalf("Incorrect number of roles: expcted %d, actual %d", roleCount, count)
	}
}
