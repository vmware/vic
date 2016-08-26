// Copyright 2016 VMware, Inc. All Rights Reserved.
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

package vm

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/sys"

	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/test"

	"golang.org/x/net/context"
)

func CreateVM(ctx context.Context, session *session.Session, host *object.HostSystem, name string) (*types.ManagedObjectReference, error) {
	// Create the spec config
	specconfig := test.SpecConfig(session, name)

	// Create a linux guest
	linux, err := guest.NewLinuxGuest(ctx, session, specconfig)
	if err != nil {
		return nil, err
	}

	// Find the Virtual Machine folder that we use
	folders, err := session.Datacenter.Folders(ctx)
	if err != nil {
		return nil, err
	}
	parent := folders.VmFolder

	// Create the vm
	info, err := tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return parent.CreateVM(ctx, *linux.Spec().Spec(), session.Pool, host)
	})
	if err != nil {
		return nil, err
	}

	moref := info.Result.(types.ManagedObjectReference)
	// Return the moRef
	return &moref, nil
}

func TestDeleteExceptDisk(t *testing.T) {
	s := os.Getenv("DRONE")
	if s != "" {
		t.Skip("Skipping: test must be run in a VM")
	}

	ctx := context.Background()

	session := test.Session(ctx, t)
	defer session.Logout(ctx)

	host := test.PickRandomHost(ctx, session, t)

	uuid, err := sys.UUID()
	if err != nil {
		t.Fatalf("unable to get UUID for guest - used for VM name: %s", err)
	}
	name := fmt.Sprintf("%s-%d", uuid, rand.Intn(math.MaxInt32))

	moref, err := CreateVM(ctx, session, host, name)
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}
	// Wrap the result with our version of VirtualMachine
	vm := NewVirtualMachine(ctx, session, *moref)

	folder, err := vm.FolderName(ctx)
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}

	// generate the disk name
	diskName := fmt.Sprintf("%s/%s.vmdk", folder, folder)

	// Delete the VM but not it's disk
	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return vm.DeleteExceptDisks(ctx)
	})
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}

	// check that the disk still exists
	session.Datastore.Stat(ctx, diskName)
	if err != nil {
		t.Fatalf("Disk does not exist")
	}

	// clean up
	dm := object.NewVirtualDiskManager(session.Client.Client)

	task, err := dm.DeleteVirtualDisk(context.TODO(), diskName, nil)
	if err != nil {
		t.Fatalf("Unable to locate orphan vmdk: %s", err)
	}

	if err = task.Wait(context.TODO()); err != nil {
		t.Fatalf("Unable to remove orphan vmdk: %s", err)
	}
}

func TestVM(t *testing.T) {

	s := os.Getenv("DRONE")
	if s != "" {
		t.Skip("Skipping: test must be run in a VM")
	}

	ctx := context.Background()

	session := test.Session(ctx, t)
	defer session.Logout(ctx)

	host := test.PickRandomHost(ctx, session, t)

	uuid, err := sys.UUID()
	if err != nil {
		t.Fatalf("unable to get UUID for guest - used for VM name: %s", err)
		return
	}
	name := fmt.Sprintf("%s-%d", uuid, rand.Intn(math.MaxInt32))

	moref, err := CreateVM(ctx, session, host, name)
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}
	// Wrap the result with our version of VirtualMachine
	vm := NewVirtualMachine(ctx, session, *moref)

	// Check the state
	state, err := vm.PowerState(ctx)
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}

	assert.Equal(t, types.VirtualMachinePowerStatePoweredOff, state)

	// Check VM name
	rname, err := vm.Name(ctx)
	if err != nil {
		t.Errorf("Failed to load VM name: %s", err)
	}
	assert.Equal(t, name, rname)

	// Get VM UUID
	ruuid, err := vm.UUID(ctx)
	if err != nil {
		t.Errorf("Failed to load VM UUID: %s", err)
	}
	t.Logf("Got UUID: %s", ruuid)

	// Destroy the vm
	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return vm.Destroy(ctx)
	})
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}
}

func TestVMFailureWithTimeout(t *testing.T) {
	ctx := context.Background()

	session := test.Session(ctx, t)
	defer session.Logout(ctx)

	host := test.PickRandomHost(ctx, session, t)

	ctx, cancel := context.WithTimeout(ctx, 1*time.Microsecond)
	defer cancel()

	uuid, err := sys.UUID()
	if err != nil {
		t.Fatalf("unable to get UUID for guest - used for VM name: %s", err)
		return
	}
	name := fmt.Sprintf("%s-%d", uuid, rand.Intn(math.MaxInt32))

	_, err = CreateVM(ctx, session, host, name)
	if err != nil && err != context.DeadlineExceeded {
		t.Fatalf("ERROR: %s", err)
	}
}

func TestVMAttributes(t *testing.T) {

	ctx := context.Background()

	session := test.Session(ctx, t)
	defer session.Logout(ctx)

	host := test.PickRandomHost(ctx, session, t)

	uuid, err := sys.UUID()
	if err != nil {
		t.Fatalf("unable to get UUID for guest - used for VM name: %s", err)
		return
	}
	ID := fmt.Sprintf("%s-%d", uuid, rand.Intn(math.MaxInt32))

	moref, err := CreateVM(ctx, session, host, ID)
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}
	// Wrap the result with our version of VirtualMachine
	vm := NewVirtualMachine(ctx, session, *moref)

	folder, err := vm.FolderName(ctx)
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}

	name, err := vm.Name(ctx)
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}
	assert.Equal(t, name, folder)

	_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return vm.PowerOn(ctx)
	})
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}

	if guest, err := vm.FetchExtraConfig(ctx); err != nil {
		t.Fatalf("ERROR: %s", err)
	} else {
		assert.NotEmpty(t, guest)
	}
	defer func() {
		// Destroy the vm
		_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
			return vm.PowerOff(ctx)
		})
		if err != nil {
			t.Fatalf("ERROR: %s", err)
		}
		_, err = tasks.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
			return vm.Destroy(ctx)
		})
		if err != nil {
			t.Fatalf("ERROR: %s", err)
		}
	}()
}
