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

package disk

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/guest"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"golang.org/x/net/context"
)

const (
	MaxAttachedDisks = 8
)

// Manager manages disks for the vm it runs on.  The expectation is this is run
// from a VM on a vsphere instance.  This VM creates disks on ESX, attaches
// them to itself, writes to them, then detaches them.
type Manager struct {
	// We can't have more than this number of disks attached.
	maxAttached chan bool

	// reference to the vm this is running on.
	vm *object.VirtualMachine

	// The controller on this vm.
	controller *types.ParaVirtualSCSIController

	// The PCI + SCSI device /dev node string format the disks can be attached with
	byPathFormat string
}

func NewDiskManager(ctx context.Context, session *session.Session) (*Manager, error) {
	vm, err := guest.GetSelf(ctx, session)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create handle to the docker daemon VM as we need to mount disks on it
	controller, byPathFormat, err := verifyParavirtualScsiController(ctx, vm)
	if err != nil {
		return nil, errors.Trace(err)
	}

	d := &Manager{
		maxAttached:  make(chan bool, MaxAttachedDisks),
		vm:           vm,
		controller:   controller,
		byPathFormat: byPathFormat,
	}
	return d, nil
}

// CreateAndAttach creates a new vmdk child from parent of the given size.
// Returns a VirtualDisk corresponding to the created and attached disk.  The
// newDiskURI and parentURI are both Datastore URI paths in the form of
// [datastoreN] /path/to/disk.vmdk.
func (m *Manager) CreateAndAttach(ctx context.Context, newDiskURI,
	parentURI string,
	capacity int64) (*VirtualDisk, error) {
	defer trace.End(trace.Begin(newDiskURI))

	// ensure we abide by max attached disks limits
	m.maxAttached <- true

	d, err := NewVirtualDisk(newDiskURI)
	if err != nil {
		return nil, errors.Trace(err)
	}

	spec := m.createDiskSpec(newDiskURI, parentURI, capacity)

	log.Infof("Attempting to create/attach vmdk for layer or volume %s at %s, with parent %s", newDiskURI, newDiskURI, parentURI)

	if err := m.vm.AddDevice(ctx, spec); err != nil {
		log.Errorf("vmdk storage driver failed to attach disk: %s", errors.ErrorStack(err))
		return nil, errors.Trace(err)
	}

	log.Debugf("Mapping vmdk to pci device %s", newDiskURI)
	devicePath, err := m.devicePathByURI(ctx, newDiskURI)
	if err != nil {
		return nil, errors.Trace(err)
	}

	d.setAttached(devicePath)

	if err := waitForPath(ctx, devicePath); err != nil {
		log.Infof("waitForPath failed for %s with %s", newDiskURI, errors.ErrorStack(err))
		// ensure that the disk is detached if it's the publish that's failed

		if detachErr := m.Detach(ctx, d); detachErr != nil {
			log.Debugf("detach(%s) failed with %s", newDiskURI, errors.ErrorStack(detachErr))
		}

		return nil, errors.Trace(err)
	}

	return d, nil
}

func (m *Manager) createDiskSpec(childURI, parentURI string, capacity int64) *types.VirtualDisk {
	// TODO: migrate this method to govmomi CreateDisk method
	backing := &types.VirtualDiskFlatVer2BackingInfo{
		DiskMode:        string(types.VirtualDiskModeIndependent_persistent),
		ThinProvisioned: types.NewBool(true),
		VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
			FileName: childURI,
		},
	}

	if parentURI != "" {
		backing.Parent = &types.VirtualDiskFlatVer2BackingInfo{
			VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
				FileName: parentURI,
			},
		}
	}

	disk := &types.VirtualDisk{
		VirtualDevice: types.VirtualDevice{
			Key:           -1,
			ControllerKey: m.controller.Key,
			UnitNumber:    -1,
			Backing:       backing,
		},
		CapacityInKB: capacity,
	}

	return disk
}

// Create creates a disk without a parent (and doesn't attach it).
func (m *Manager) Create(ctx context.Context, newDiskURI string,
	capacityKB int64) (*VirtualDisk, error) {

	defer trace.End(trace.Begin(newDiskURI))

	vdm := object.NewVirtualDiskManager(m.vm.Client())

	d, err := NewVirtualDisk(newDiskURI)
	if err != nil {
		return nil, errors.Trace(err)
	}

	spec := &types.FileBackedVirtualDiskSpec{
		VirtualDiskSpec: types.VirtualDiskSpec{
			DiskType:    string(types.VirtualDiskTypeThin),
			AdapterType: string(types.VirtualDiskAdapterTypeLsiLogic),
		},

		CapacityKb: capacityKB,
	}

	log.Infof("Attempting to create vmdk for layer or volume %s", d.DatastoreURI)

	err = tasks.Wait(ctx, func(ctx context.Context) (tasks.Waiter, error) {
		return vdm.CreateVirtualDisk(ctx, d.DatastoreURI, nil, spec)
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	return d, nil
}

// TODO(FA) this doesn't work since delta disks get set with `deletable =
// false` when they become parents.  This needs some thought and will require
// some answers from a larger context.
//func (m *DiskManager) Delete(ctx context.Context, d *VirtualDisk) error {
//	defer trace.End(trace.Begin(d.DatastoreURI))
//
//	log.Infof("Deleting %s", d.DatastoreURI)
//
//	d.lock()
//	defer d.unlock()
//
//	if d.isAttached() {
//		return fmt.Errorf("cannot delete %s, still attached (%s)", d.DatastoreURI, d.devicePath)
//	}
//
//	// TODO(FA) Check if disk is a parent.
//
//	vdm := object.NewVirtualDiskManager(m.vm.Client())
//	task, err := vdm.DeleteVirtualDisk(ctx, d.DatastoreURI, nil)
//	if err != nil {
//		return err
//	}
//
//	err = task.Wait(ctx)
//	if err != nil {
//		return errors.Trace(err)
//	}
//
//	return nil
// }

func (m *Manager) Detach(ctx context.Context, d *VirtualDisk) error {
	defer trace.End(trace.Begin(d.DevicePath))
	log.Infof("starting detaching disk %s", d.DevicePath)

	d.lock()
	defer d.unlock()

	if err := d.canBeDetached(); err != nil {
		return errors.Trace(err)
	}

	spec := types.VirtualMachineConfigSpec{}

	disk, err := findDisk(ctx, m.vm, d.DatastoreURI)
	if err != nil {
		return errors.Trace(err)
	}

	config := []types.BaseVirtualDeviceConfigSpec{
		&types.VirtualDeviceConfigSpec{
			Device:    disk,
			Operation: types.VirtualDeviceConfigSpecOperationRemove,
		},
	}

	spec.DeviceChange = config

	err = tasks.Wait(ctx, func(ctx context.Context) (tasks.Waiter, error) {
		return m.vm.Reconfigure(ctx, spec)
	})
	if err != nil {
		log.Warnf("detach for %s failed with %s", d.DevicePath, errors.ErrorStack(err))
		return errors.Trace(err)
	}

	func() {
		select {
		case <-m.maxAttached:
		default:
		}
	}()

	return d.setDetached()
}

func (m *Manager) devicePathByURI(ctx context.Context, datastoreURI string) (string, error) {
	disk, err := findDisk(ctx, m.vm, datastoreURI)
	if err != nil {
		log.Debugf("findDisk failed for %s with %s", datastoreURI, errors.ErrorStack(err))
		return "", errors.Trace(err)
	}

	return fmt.Sprintf(m.byPathFormat, disk.UnitNumber), nil
}
