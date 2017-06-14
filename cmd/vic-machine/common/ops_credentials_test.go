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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessOpsCredentials(t *testing.T) {
	createOps := &OpsCredentials{}
	isCreateOp := true
	adminUser := "admin"
	adminPassword := ""

	// There should be an error if the admin password is not specified for a create operation.
	err := createOps.ProcessOpsCredentials(isCreateOp, adminUser, nil)
	assert.NotNil(t, err)

	err = createOps.ProcessOpsCredentials(isCreateOp, adminUser, &adminPassword)
	assert.NoError(t, err)
	assert.Equal(t, *createOps.OpsUser, adminUser)
	assert.Equal(t, *createOps.OpsPassword, adminPassword)

	opsUser := "op"
	opsPassword := "opPass"
	createOps.OpsUser = &opsUser
	createOps.OpsPassword = &opsPassword
	err = createOps.ProcessOpsCredentials(isCreateOp, adminUser, &adminPassword)
	assert.NoError(t, err)
	assert.Equal(t, *createOps.OpsUser, opsUser)
	assert.Equal(t, *createOps.OpsPassword, opsPassword)

	// Ensure that fields are set correctly for a configure operation.
	configureOps := &OpsCredentials{
		OpsUser:     &opsUser,
		OpsPassword: &opsPassword,
	}
	isCreateOp = false
	err = configureOps.ProcessOpsCredentials(isCreateOp, "", nil)
	assert.NoError(t, err)
	assert.True(t, configureOps.IsSet)
	assert.Equal(t, *createOps.OpsUser, opsUser)
	assert.Equal(t, *createOps.OpsPassword, opsPassword)

	// There should be an error if the ops-password is specified without ops-user.
	configureOps.OpsUser = nil
	err = configureOps.ProcessOpsCredentials(isCreateOp, "", nil)
	assert.NotNil(t, err)
}
