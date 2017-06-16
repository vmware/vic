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

package disk

import (
	"context"
	"fmt"
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
	// You can assign the device to (1:z ), where 1 is SCSI controller 1 and z is a virtual device node from 0 to 15.
	// https://pubs.vmware.com/vsphere-65/index.jsp#com.vmware.vsphere.vm_admin.doc/GUID-5872D173-A076-42FE-8D0B-9DB0EB0E7362.html
	MaxAttachedDisks = 16
)

// Manager manages disks for the vm it runs on.  The expectation is this is run
// from a VM on a vsphere instance.  This VM creates disks on ESX, attaches
// them to itself, writes to them, then detaches them.
type Manager struct {
	// We can't have more than this number of disks attached.
	maxAttached chan bool

	// reference to the vm this is running on.
	vm *vm.VirtualMachine

	// VirtualDiskManager that is used to create vmdks directly on datastore
	// from https://pubs.vmware.com/vsphere-65/index.jsp?topic=%2Fcom.vmware.vspsdk.apiref.doc%2Fvim.VirtualDiskManager.html
	// Most VirtualDiskManager APIs will be DEPRECATED as of vSphere 6.5. Please use VStorageObjectManager APIs to manage Virtual disks.
	vdm *object.VirtualDiskManager

	// The controller on this vm.
	controller *types.ParaVirtualSCSIController

	// The PCI + SCSI device /dev node string format the disks can be attached with
	byPathFormat string

	// serialize reconfigure operations
	mu sync.Mutex

	// map of URIs to VirtualDisk structs so that we can return the same instance to the caller, required for ref counting
	Disks map[uint64]*VirtualDisk
}

// NewDiskManager creates a new Manager instance associated with the caller VM
func NewDiskManager(op trace.Operation, session *session.Session) (*Manager, error) {
	defer trace.End(trace.Begin(""))

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

	return &Manager{
		maxAttached:  make(chan bool, MaxAttachedDisks),
		vm:           vm,
		vdm:          object.NewVirtualDiskManager(vm.Vim25()),
		controller:   controller,
		byPathFormat: byPathFormat,
		Disks:        make(map[uint64]*VirtualDisk),
	}, nil
}

// toSpec converts the given config to VirtualDisk spec
func (m *Manager) toSpec(config *VirtualDiskConfig) *types.VirtualDisk {
	backing := &types.VirtualDiskFlatVer2BackingInfo{
		DiskMode:        string(config.DiskMode),
		ThinProvisioned: types.NewBool(true),
		VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
			FileName: config.DatastoreURI.String(),
		},
	}

	disk := &types.VirtualDisk{
		VirtualDevice: types.VirtualDevice{
			Key:           -1,
			ControllerKey: m.controller.Key,
			UnitNumber:    new(int32),
			Backing:       backing,
		},
		// As of vSphere API 5.5 capacityInKB is deprecated. Documentation suggest using capacityInBytes but we can't unset CapacityInKB and its default value 0 causes problems
		// ... Exception thrown during reconfigure: (vim.vm.ConfigSpec) {
		// ...
		// -->             unitNumber = -1,
		// -->             capacityInKB = 0,
		// -->             capacityInBytes = 8192000000,
		// -->             shares = (vim.SharesInfo) null,
		// ...
		CapacityInBytes: config.CapacityInKB * 1024,
		CapacityInKB:    config.CapacityInKB,
	}

	if config.ParentDatastoreURI != nil {
		backing.Parent = &types.VirtualDiskFlatVer2BackingInfo{
			VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
				FileName: config.ParentDatastoreURI.String(),
			},
		}

		// Capacity needs to be 0 as we inherit it from the parent
		disk.CapacityInBytes = 0
		disk.CapacityInKB = 0
	}

	// It's possible the VCH has a disk already attached.
	*disk.VirtualDevice.UnitNumber = -1

	return disk
}

// CreateAndAttach creates a new VMDK, attaches it and ensures that the device becomes visible to the caller.
// Returns a VirtualDisk corresponding to the created and attached disk.
func (m *Manager) CreateAndAttach(op trace.Operation, config *VirtualDiskConfig) (*VirtualDisk, error) {
	defer trace.End(trace.Begin(config.DatastoreURI.String()))

	// ensure we abide by max attached disks limits
	m.maxAttached <- true

	op.Infof("Create/attach vmdk %s from parent %s", config.DatastoreURI, config.ParentDatastoreURI)

	if err := m.attach(op, config); err != nil {
		return nil, errors.Trace(err)
	}

	op.Debugf("Mapping vmdk to pci device %s", config.DatastoreURI)
	devicePath, err := m.devicePathByURI(op, config.DatastoreURI)
	if err != nil {
		return nil, errors.Trace(err)
	}

	d, err := NewVirtualDisk(config, m.Disks)
	if err != nil {
		return nil, errors.Trace(err)
	}

	blockDev, err := waitForDevice(op, devicePath)
	if err != nil {
		op.Errorf("waitForDevice failed for %s with %s", config.DatastoreURI, errors.ErrorStack(err))
		// ensure that the disk is detached if it's the publish that's failed

		if detachErr := m.Detach(op, config); detachErr != nil {
			op.Debugf("detach(%s) failed with %s", config.DatastoreURI, errors.ErrorStack(detachErr))
		}

		return nil, errors.Trace(err)
	}

	d.ParentDatastoreURI = config.ParentDatastoreURI
	d.setAttached(blockDev)

	return d, nil
}

// Create creates a disk without a parent (and doesn't attach it).
func (m *Manager) Create(op trace.Operation, config *VirtualDiskConfig) (*VirtualDisk, error) {
	defer trace.End(trace.Begin(config.DatastoreURI.String()))

	var err error

	d, err := NewVirtualDisk(config, m.Disks)
	if err != nil {
		return nil, errors.Trace(err)
	}

	spec := &types.FileBackedVirtualDiskSpec{
		VirtualDiskSpec: types.VirtualDiskSpec{
			DiskType:    string(types.VirtualDiskTypeThin),
			AdapterType: string(types.VirtualDiskAdapterTypeLsiLogic),
		},
		CapacityKb: config.CapacityInKB,
	}

	op.Infof("Creating vmdk for layer or volume %s", d.DatastoreURI)
	err = tasks.Wait(op, func(ctx context.Context) (tasks.Task, error) {
		return m.vdm.CreateVirtualDisk(ctx, d.DatastoreURI.String(), nil, spec)
	})

	if err != nil {
		return nil, errors.Trace(err)
	}

	return d, nil
}

// Gets a disk given a datastore path URI to the vmdk
func (m *Manager) Get(op trace.Operation, config *VirtualDiskConfig) (*VirtualDisk, error) {
	defer trace.End(trace.Begin(config.DatastoreURI.String()))

	d, err := NewVirtualDisk(config, m.Disks)
	if err != nil {
		return nil, errors.Trace(err)
	}

	info, err := m.vdm.QueryVirtualDiskInfo(op, config.DatastoreURI.String(), m.vm.Datacenter, true)
	if err != nil {
		op.Errorf("error querying parents (%s): %s", config.DatastoreURI, err.Error())
		return nil, err
	}

	// the last elem in the info list is the disk we just looked up.
	p := info[len(info)-1]

	if p.Parent != "" {
		ppth, err := datastore.DatastorePathFromString(p.Parent)
		if err != nil {
			return nil, err
		}
		d.ParentDatastoreURI = ppth
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

// Attach attempts to attach a virtual disk
func (m *Manager) attach(op trace.Operation, config *VirtualDiskConfig) error {
	defer trace.End(trace.Begin(""))

	disk := m.toSpec(config)

	deviceList := object.VirtualDeviceList{}
	deviceList = append(deviceList, disk)

	changeSpec, err := deviceList.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	if err != nil {
		return err
	}

	machineSpec := types.VirtualMachineConfigSpec{}
	machineSpec.DeviceChange = append(machineSpec.DeviceChange, changeSpec...)

	m.mu.Lock()
	defer m.mu.Unlock()

	_, err = m.vm.WaitForResult(op, func(ctx context.Context) (tasks.Task, error) {
		t, er := m.vm.Reconfigure(ctx, machineSpec)

		if t != nil {
			op.Debugf("Attach reconfigure task=%s", t.Reference())
		}

		return t, er
	})

	if err != nil {
		op.Errorf("vmdk storage driver failed to attach disk: %s", errors.ErrorStack(err))
		return errors.Trace(err)
	}

	return nil
}

// Detach attempts to detach a virtual disk
func (m *Manager) Detach(op trace.Operation, config *VirtualDiskConfig) error {
	defer trace.End(trace.Begin(""))

	d, err := NewVirtualDisk(config, m.Disks)
	if err != nil {
		return errors.Trace(err)
	}

	op.Infof("Detaching disk %s", d.DevicePath)

	d.lock()
	defer d.unlock()

	if !d.Attached() {
		op.Infof("Disk %s is already detached", d.DevicePath)
		return nil
	}

	if err := d.canBeDetached(); err != nil {
		// even though canBeDetached() is called here and nowhere else, it does not imply
		// an attempt to detach, so decrease the ref count from here for that reason
		d.attachedRefs--
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

	// don't decrement reference count here as setDetached() does it already
	return d.setDetached()
}

func (m *Manager) DetachAll(op trace.Operation) error {
	defer trace.End(trace.Begin(""))

	disks, err := findAllDisks(op, m.vm)
	if err != nil {
		return err
	}

	for _, disk := range disks {
		if err2 := m.detach(op, disk); err != nil {
			op.Errorf("error detaching disk: %s", err2.Error())
			// return the last error on the return of this function
			err = err2
			// if we failed here that means we have a disk attached, ensure we abide by max attached disks limits
			m.maxAttached <- true
		}
	}

	return err
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

	m.mu.Lock()
	defer m.mu.Unlock()

	_, err := m.vm.WaitForResult(op, func(ctx context.Context) (tasks.Task, error) {
		t, er := m.vm.Reconfigure(ctx, spec)

		if t != nil {
			op.Debugf("Detach reconfigure task=%s", t.Reference())
		}

		return t, er
	})

	return err
}

func (m *Manager) devicePathByURI(op trace.Operation, datastoreURI *object.DatastorePath) (string, error) {
	disk, err := findDiskByFilename(op, m.vm, datastoreURI.String())
	if err != nil {
		op.Errorf("findDisk failed for %s with %s", datastoreURI.String(), errors.ErrorStack(err))
		return "", errors.Trace(err)
	}

	sysPath := fmt.Sprintf(m.byPathFormat, *disk.UnitNumber)

	return sysPath, nil
}
