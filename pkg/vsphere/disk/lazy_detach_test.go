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

package disk

import (
	"context"
	"path"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/tasks"
)

// TestLazyDetach tests lazy detach functionality to make sure that every ESXi version shows this behaviour
// https://github.com/vmware/vic/issues/5565
func TestLazyDetach(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	client := Session(context.Background(), t)
	if client == nil {
		return
	}

	imagestore := &object.DatastorePath{
		Datastore: client.Datastore.Name(),
		Path:      datastore.TestName("lazyReconfigure"),
	}

	fm := object.NewFileManager(client.Vim25())

	// create a directory in the datastore
	// eat the error because we dont care if it exists
	fm.MakeDirectory(context.TODO(), imagestore.String(), nil, true)
	// Nuke the image store
	defer func() {
		task, err := fm.DeleteDatastoreFile(context.TODO(), imagestore.String(), nil)
		if err != nil && err.Error() == "can't find the hosting vm" {
			t.Skip("Skipping: test must be run in a VM")
		}
		if !assert.NoError(t, err) {
			return
		}
		_, err = task.WaitForResult(context.TODO(), nil)
		if !assert.NoError(t, err) {
			return
		}
	}()

	op := trace.NewOperation(context.TODO(), "lazyReconfigure")

	// create a diskmanager
	vdm, err := NewDiskManager(op, client)
	if err != nil && err.Error() == "can't find the hosting vm" {
		t.Skip("Skipping: test must be run in a VM")
	}
	if !assert.NoError(t, err) || !assert.NotNil(t, vdm) {
		return
	}

	// helper fn
	reconfigure := func(changes []types.BaseVirtualDeviceConfigSpec) error {
		t.Logf("Calling reconfigure")

		machineSpec := types.VirtualMachineConfigSpec{}
		machineSpec.DeviceChange = changes

		_, err := vdm.vm.WaitForResult(op, func(ctx context.Context) (tasks.Task, error) {
			t, er := vdm.vm.Reconfigure(ctx, machineSpec)

			if t != nil {
				op.Debugf("reconfigure task=%s", t.Reference())
			}

			return t, er
		})
		return err
	}

	// 1MB
	diskSize := int64(1 << 10)
	scratch := &object.DatastorePath{
		Datastore: client.Datastore.Name(),
		Path:      path.Join(imagestore.Path, "scratch.vmdk"),
	}
	// config
	config := NewPersistentDisk(scratch).WithCapacity(diskSize)

	// attach + create spec (scratch)
	spec := []types.BaseVirtualDeviceConfigSpec{
		&types.VirtualDeviceConfigSpec{
			Device:        vdm.toSpec(config),
			Operation:     types.VirtualDeviceConfigSpecOperationAdd,
			FileOperation: types.VirtualDeviceConfigSpecFileOperationCreate,
		},
	}
	// reconfigure
	err = reconfigure(spec)
	if !assert.NoError(t, err) {
		return
	}
	t.Logf("scratch created and attached")

	// ref to scratch (needed for detach as initial spec's Key and UnitNumber was unset)
	disk, err := findDiskByFilename(op, vdm.vm, scratch.Path)
	if !assert.NoError(t, err) {
		return
	}

	// DO NOT DETACH AND START WORKING ON THE CHILD

	// child
	child := &object.DatastorePath{
		Datastore: client.Datastore.Name(),
		Path:      path.Join(imagestore.Path, "child.vmdk"),
	}
	// config
	config = NewPersistentDisk(child).WithParent(scratch)

	// detach (scratch) AND attach + create (child) spec
	spec = []types.BaseVirtualDeviceConfigSpec{
		&types.VirtualDeviceConfigSpec{
			Device:    disk,
			Operation: types.VirtualDeviceConfigSpecOperationRemove,
		},
		&types.VirtualDeviceConfigSpec{
			Device:        vdm.toSpec(config),
			Operation:     types.VirtualDeviceConfigSpecOperationAdd,
			FileOperation: types.VirtualDeviceConfigSpecFileOperationCreate,
		},
	}
	// reconfigure
	err = reconfigure(spec)
	if !assert.NoError(t, err) {
		return
	}
	t.Logf("scratch detached, child created and attached")

	// ref to child (needed for detach as initial spec's Key and UnitNumber was unset)
	disk, err = findDiskByFilename(op, vdm.vm, child.Path)
	if !assert.NoError(t, err) {
		return
	}

	// detach  spec (child)
	spec = []types.BaseVirtualDeviceConfigSpec{
		&types.VirtualDeviceConfigSpec{
			Device:    disk,
			Operation: types.VirtualDeviceConfigSpecOperationRemove,
		},
	}
	// reconfigure
	err = reconfigure(spec)
	if !assert.NoError(t, err) {
		return
	}
	t.Logf("child detached")

	// Expected outcome; 3 reconfigure operations
	//	--- PASS: TestLazyDetach (0.68s)
	//        lazy_detach_test.go:80: Calling reconfigure
	//        lazy_detach_test.go:119: scratch created and attached
	//        lazy_detach_test.go:80: Calling reconfigure
	//        lazy_detach_test.go:154: scratch detached, child created and attached
	//        lazy_detach_test.go:80: Calling reconfigure
	//        lazy_detach_test.go:174: child detached
	//	PASS
}
