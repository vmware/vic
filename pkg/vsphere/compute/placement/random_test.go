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

package placement

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/test"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// vpxModelSetup creates a VPX model, starts its server and populates the session. The caller must
// clean up the model and the server once it is done using them.
func vpxModelSetup(ctx context.Context, t *testing.T) (*simulator.Model, *simulator.Server, *session.Session) {
	model := simulator.VPX()
	if err := model.Create(); err != nil {
		t.Fatal(err)
	}

	server := model.Service.NewServer()
	sess, err := test.SessionWithVPX(ctx, server.URL.String())
	if err != nil {
		t.Fatal(err)
	}

	return model, server, sess
}

func TestRandomRecommendHost(t *testing.T) {
	op := trace.NewOperation(context.Background(), "TestRandomRecommendHost")

	model, server, sess := vpxModelSetup(op, t)
	defer func() {
		model.Remove()
		server.Close()
	}()

	cls := sess.Cluster

	hosts, err := cls.Hosts(op)
	assert.NoError(t, err)

	name := "test-vm"
	vmx := fmt.Sprintf("%s/%s.vmx", name, name)
	ds := sess.Datastore
	secretKey, err := extraconfig.NewSecretKey()
	if err != nil {
		t.Fatal(err)
	}

	spec := types.VirtualMachineConfigSpec{
		Name:    name,
		GuestId: string(types.VirtualMachineGuestOsIdentifierOtherGuest),
		Files: &types.VirtualMachineFileInfo{
			VmPathName: fmt.Sprintf("[%s] %s", ds.Name(), vmx),
		},
		ExtraConfig: []types.BaseOptionValue{
			&types.OptionValue{
				Key:   extraconfig.GuestInfoSecretKey,
				Value: secretKey.String(),
			},
		},
	}

	res, err := tasks.WaitForResult(op, func(op context.Context) (tasks.Task, error) {
		return sess.VMFolder.CreateVM(op, spec, sess.Pool, nil)
	})
	assert.NoError(t, err)

	v := vm.NewVirtualMachine(op, sess, res.Result.(types.ManagedObjectReference))

	rhp := NewRandomHostPolicy()
	assert.False(t, rhp.CheckHost(op, v))
	h, err := rhp.RecommendHost(op, v)
	assert.NoError(t, err)

	// TODO(jzt): come up with a better way to verify this than using a loop/moref comparison
	found := false
	for _, host := range hosts {
		if h.Reference().String() == host.Reference().String() {
			found = true
			break
		}
	}
	assert.True(t, found)
}
