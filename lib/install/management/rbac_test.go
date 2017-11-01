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

	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/test/env"
)

func TestRolesSimulatorVPX(t *testing.T) {
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

	session, err := session.NewSession(config).Connect(ctx)
	require.NoError(t, err, "Cannot connect to VPX Simulator")

	mgr := NewAuthzManager(ctx, session.Vim25(), "ops-user")

	doTestRoles(ctx, t, mgr)
}

func TestRolesVCenter(t *testing.T) {
	ctx := context.Background()

	config := &session.Config{
		Service:   env.URL(t),
		Insecure:  true,
		Keepalive: time.Duration(5) * time.Minute,
	}

	session, err := session.NewSession(config).Connect(ctx)
	if err != nil {
		t.SkipNow()
	}

	mgr := NewAuthzManager(ctx, session.Vim25(), "ops-user")

	doTestRoles(ctx, t, mgr)
}

func initRoles(ctx context.Context, t *testing.T, am *AuthzManager, tRoles []types.AuthorizationRole) int {
	cleanup(ctx, t, am, false)

	count, err := am.createOrRepairRoles(ctx, tRoles)
	require.NoError(t, err, "Failed to initialize Roles")

	return count
}

func doTestRoles(ctx context.Context, t *testing.T, am *AuthzManager) {
	targetRoles := getTargetRoles()
	var roleCount = len(targetRoles)
	count := initRoles(ctx, t, am, targetRoles)

	defer cleanup(ctx, t, am, true)
	require.Equal(t, roleCount, count, "Incorrect number of roles: expected %d, actual %d", roleCount, count)

	// Test correct role validation, it should return 0
	roleCount = 0
	count, err := am.createOrRepairRoles(ctx, targetRoles)
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
	count, err = am.createOrRepairRoles(ctx, targetRoles)
	require.NoError(t, err, "Failed to repair roles 1")
	require.Equal(t, roleCount, count, "Incorrect number of roles: expected %d, actual %d", roleCount, count)

	// Test correct role validation, it should return 0
	roleCount = 0
	count, err = am.createOrRepairRoles(ctx, targetRoles)
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
	targetRoles := getTargetRoles()
	var roleCount = len(targetRoles)
	count, err := am.deleteRoles(ctx, targetRoles)
	require.NoError(t, err, "Failed to delete roles")

	if checkCount && count != roleCount {
		t.Fatalf("Incorrect number of roles: expcted %d, actual %d", roleCount, count)
	}
}
