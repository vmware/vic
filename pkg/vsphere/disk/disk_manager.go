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
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/guest"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
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
	vm *vm.VirtualMachine

	// The controller on this vm.
	controller *types.ParaVirtualSCSIController

	// The PCI + SCSI device /dev node string format the disks can be attached with
	byPathFormat string

	reconfig sync.Mutex
}

func NewDiskManager(op trace.Operation, session *session.Session) (*Manager, error) {
	vm, err := guest.GetSelf(op, session)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// create handle to the docker daemon VM as we need to mount disks on it
	controller, byPathFormat, err := verifyParavirtualScsiController(op, vm)
	if err != nil {
		op.Errorf("scsi controller verification failed: %s", err.Error())
		return nil, err
	}

	d := &Manager{
		maxAttached:  make(chan bool, MaxAttachedDisks),
		vm:           vm,
		controller:   controller,
		byPathFormat: byPathFormat,
	}

	// Remove any attached disks
	if err = d.detachAll(op); err != nil {
		return nil, err
	}

	return d, nil
}

// CreateAndAttach creates a new vmdk child from parent of the given size.
// Returns a VirtualDisk corresponding to the created and attached disk.  The
// newDiskURI and parentURI are both Datastore URI paths in the form of
// [datastoreN] /path/to/disk.vmdk.
func (m *Manager) CreateAndAttach(op trace.Operation, newDiskURI,
	parentURI string,
	capacity int64, flags int) (*VirtualDisk, error) {
	defer trace.End(trace.Begin(newDiskURI))

	// ensure we abide by max attached disks limits
	m.maxAttached <- true

	spec := m.createDiskSpec(newDiskURI, parentURI, capacity, flags)

	op.Infof("Create/attach vmdk %s from parent %s", newDiskURI, parentURI)

	if err := m.Attach(op, spec); err != nil {
		return nil, errors.Trace(err)
	}

	op.Debugf("Mapping vmdk to pci device %s", newDiskURI)
	devicePath, err := m.devicePathByURI(op, newDiskURI)
	if err != nil {
		return nil, errors.Trace(err)
	}

	d, err := NewVirtualDisk(newDiskURI)
	if err != nil {
		return nil, errors.Trace(err)
	}

	blockDev, err := waitForDevice(op, devicePath)
	if err != nil {
		op.Errorf("waitForDevice failed for %s with %s", newDiskURI, errors.ErrorStack(err))
		// ensure that the disk is detached if it's the publish that's failed

		if detachErr := m.Detach(op, d); detachErr != nil {
			op.Debugf("detach(%s) failed with %s", newDiskURI, errors.ErrorStack(detachErr))
		}

		return nil, errors.Trace(err)
	}

	var ppth *object.DatastorePath
	if parentURI != "" {
		ppth, err = datastore.DatastorePathFromString(parentURI)
		if err != nil {
			return nil, err
		}
	}

	d.ParentDatastoreURI = ppth
	d.setAttached(blockDev)

	return d, nil
}

func (m *Manager) createDiskSpec(childURI, parentURI string, capacity int64, flags int) *types.VirtualDisk {
	// TODO: migrate this method to govmomi CreateDisk method
	backing := &types.VirtualDiskFlatVer2BackingInfo{
		DiskMode:        string(types.VirtualDiskModeIndependent_persistent),
		ThinProvisioned: types.NewBool(true),
		VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
			FileName: childURI,
		},
	}

	if flags == os.O_RDONLY {
		backing.DiskMode = string(types.VirtualDiskModeIndependent_nonpersistent)
		capacity = 0
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
			UnitNumber:    new(int32),
			Backing:       backing,
		},
		CapacityInKB: capacity,
	}

	// It's possible the VCH has a disk already attached.
	*disk.VirtualDevice.UnitNumber = -1

	return disk
}

// Create creates a disk without a parent (and doesn't attach it).
func (m *Manager) Create(op trace.Operation, newDiskURI string,
	capacityKB int64) (*VirtualDisk, error) {

	defer trace.End(trace.Begin(newDiskURI))

	vdm := object.NewVirtualDiskManager(m.vm.Vim25())

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

	op.Infof("Creating vmdk for layer or volume %s", d.DatastoreURI)
	err = tasks.Wait(op, func(ctx context.Context) (tasks.Task, error) {
		return vdm.CreateVirtualDisk(ctx, d.DatastoreURI.String(), nil, spec)
	})

	if err != nil {
		return nil, errors.Trace(err)
	}

	return d, nil
}

// Gets a disk given a datastore path URI to the vmdk
func (m *Manager) Get(op trace.Operation, diskURI string) (*VirtualDisk, error) {
	defer trace.End(trace.Begin(diskURI))

	dsk, err := NewVirtualDisk(diskURI)
	if err != nil {
		return nil, err
	}

	vdm := object.NewVirtualDiskManager(m.vm.Vim25())

	info, err := vdm.QueryVirtualDiskInfo(op, diskURI, nil, true)
	if err != nil {
		op.Errorf("error querying parents (%s): %s", diskURI, err.Error())
		return nil, err
	}

	// the last elem in the info list is the disk we just looked up.
	p := info[len(info)-1]

	if p.Parent != "" {
		ppth, err := datastore.DatastorePathFromString(p.Parent)
		if err != nil {
			return nil, err
		}
		dsk.ParentDatastoreURI = ppth
	}

	return dsk, nil
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

func (m *Manager) Attach(op trace.Operation, disk *types.VirtualDisk) error {
	deviceList := object.VirtualDeviceList{}
	deviceList = append(deviceList, disk)

	changeSpec, err := deviceList.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	if err != nil {
		return err
	}

	machineSpec := types.VirtualMachineConfigSpec{}
	machineSpec.DeviceChange = append(machineSpec.DeviceChange, changeSpec...)

	m.reconfig.Lock()
	_, err = m.vm.WaitForResult(op, func(ctx context.Context) (tasks.Task, error) {
		t, er := m.vm.Reconfigure(ctx, machineSpec)

		if t != nil {
			op.Debugf("Attach reconfigure task=%s", t.Reference())
		}

		return t, er
	})
	m.reconfig.Unlock()

	if err != nil {
		op.Errorf("vmdk storage driver failed to attach disk: %s", errors.ErrorStack(err))
		return errors.Trace(err)
	}
	return nil
}

func (m *Manager) Detach(op trace.Operation, d *VirtualDisk) error {
	defer trace.End(trace.Begin(d.DevicePath))
	op.Infof("Detaching disk %s", d.DevicePath)

	d.lock()
	defer d.unlock()

	if !d.Attached() {
		op.Infof("Disk %s is already detached", d.DevicePath)
		return nil
	}

	if err := d.canBeDetached(); err != nil {
		return errors.Trace(err)
	}

	disk, err := findDiskByFilename(op, m.vm, d.DatastoreURI.String())
	if err != nil {
		return errors.Trace(err)
	}

	if err = m.detach(op, disk); err != nil {
		op.Errorf("detach for %s failed with %s", d.DevicePath, errors.ErrorStack(err))
		return errors.Trace(err)
	}

	select {
	case <-m.maxAttached:
	default:
	}

	return d.setDetached()
}

func (m *Manager) detach(op trace.Operation, disk *types.VirtualDisk) error {

	config := []types.BaseVirtualDeviceConfigSpec{
		&types.VirtualDeviceConfigSpec{
			Device:    disk,
			Operation: types.VirtualDeviceConfigSpecOperationRemove,
		},
	}

	spec := types.VirtualMachineConfigSpec{}
	spec.DeviceChange = config

	m.reconfig.Lock()
	_, err := m.vm.WaitForResult(op, func(ctx context.Context) (tasks.Task, error) {
		t, er := m.vm.Reconfigure(ctx, spec)

		if t != nil {
			op.Debugf("Detach reconfigure task=%s", t.Reference())
		}

		return t, er
	})
	m.reconfig.Unlock()

	return err
}

// detachAll detaches all disks from this vm
func (m *Manager) detachAll(op trace.Operation) error {
	disks, err := findAllDisks(op, m.vm)
	if err != nil {
		return err
	}

	for _, disk := range disks {
		if er := m.detach(op, disk); err != nil {
			// late exit on error
			op.Errorf("error detaching disk: %s", er.Error())
			err = er
		}
	}

	return err
}

func (m *Manager) devicePathByURI(op trace.Operation, datastoreURI string) (string, error) {
	disk, err := findDiskByFilename(op, m.vm, datastoreURI)
	if err != nil {
		op.Errorf("findDisk failed for %s with %s", datastoreURI, errors.ErrorStack(err))
		return "", errors.Trace(err)
	}

	sysPath := fmt.Sprintf(m.byPathFormat, *disk.UnitNumber)

	return sysPath, nil
}
