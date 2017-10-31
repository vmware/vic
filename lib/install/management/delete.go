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
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

type DeleteContainers int

const (
	AllContainers DeleteContainers = iota
	PoweredOffContainers
)

type DeleteVolumeStores int

const (
	AllVolumeStores DeleteVolumeStores = iota
	NoVolumeStores
)

func (d *Dispatcher) DeleteVCH(conf *config.VirtualContainerHostConfigSpec, containers *DeleteContainers, volumeStores *DeleteVolumeStores) error {
	defer trace.End(trace.Begin(conf.Name))

	var errs []string

	var err error
	var vmm *vm.VirtualMachine

	if vmm, err = d.findApplianceByID(conf); err != nil {
		return err
	}
	if vmm == nil {
		return nil
	}

	if err = d.DeleteVCHInstances(vmm, conf, containers); err != nil {
		// if container delete failed, do not remove anything else
		log.Infof("Specify --force to force delete")
		return err
	}

	if err = d.deleteImages(conf); err != nil {
		errs = append(errs, err.Error())
	}

	d.deleteVolumeStoreIfForced(conf, volumeStores) // logs errors but doesn't ever bail out if it has an issue

	if err = d.deleteNetworkDevices(vmm, conf); err != nil {
		errs = append(errs, err.Error())
	}
	if err = d.removeNetwork(conf); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		// stop here, leave vch appliance there for next time delete
		return errors.New(strings.Join(errs, "\n"))
	}

	err = d.deleteVM(vmm, true)
	if err != nil {
		log.Debugf("Error deleting appliance VM %s", err)
		return err
	}
	if err = d.destroyResourcePoolIfEmpty(conf); err != nil {
		log.Warnf("VCH resource pool is not removed: %s", err)
	}
	return nil
}

func (d *Dispatcher) getComputeResource(vmm *vm.VirtualMachine, conf *config.VirtualContainerHostConfigSpec) (*compute.ResourcePool, error) {
	var rpRef types.ManagedObjectReference
	var err error

	ignoreFailureToFindComputeResources := d.force

	if len(conf.ComputeResources) == 0 {
		if !ignoreFailureToFindComputeResources {
			err = errors.Errorf("Cannot find compute resources from configuration")
			return nil, err
		}
		log.Warnf("Cannot find compute resources from configuration, attempting to delete under parent resource pool")
		parent, err := vmm.Parent(d.ctx)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			err = errors.Errorf("Cannot find VCH parent resource pool")
			return nil, err
		}
		rpRef = *parent
	} else {
		rpRef = conf.ComputeResources[len(conf.ComputeResources)-1]
	}

	ref, err := d.session.Finder.ObjectReference(d.ctx, rpRef)
	if err != nil {
		err = errors.Errorf("Failed to get VCH resource pool %q: %s", rpRef, err)
		return nil, err
	}
	switch ref.(type) {
	case *object.VirtualApp:
	case *object.ResourcePool:
		//		ok
	default:
		err = errors.Errorf("Unsupported compute resource %q", rpRef)
		return nil, err
	}

	rp := compute.NewResourcePool(d.ctx, d.session, ref.Reference())
	return rp, nil
}

func (d *Dispatcher) getImageDatastore(vmm *vm.VirtualMachine, conf *config.VirtualContainerHostConfigSpec, force bool) (*object.Datastore, error) {
	var err error
	if conf == nil || len(conf.ImageStores) == 0 {
		if !force {
			err = errors.Errorf("Cannot find image stores from configuration")
			return nil, err
		}
		log.Debugf("Cannot find image stores from configuration; attempting to find from vm datastore")
		dss, err := vmm.DatastoreReference(d.ctx)
		if err != nil {
			return nil, errors.Errorf("Failed to query vm datastore: %s", err)
		}
		if len(dss) == 0 {
			return nil, errors.New("No VCH datastore found, cannot continue")
		}
		ds, err := d.session.Finder.ObjectReference(d.ctx, dss[0])
		if err != nil {
			return nil, errors.Errorf("Failed to search vm datastore %s: %s", dss[0], err)
		}
		return ds.(*object.Datastore), nil
	}
	ds, err := d.session.Finder.Datastore(d.ctx, conf.ImageStores[0].Host)
	if err != nil {
		err = errors.Errorf("Failed to find image datastore %q", conf.ImageStores[0].Host)
		return nil, err
	}
	return ds, nil
}

// detach all VMDKs attached to vm
func (d *Dispatcher) detachAttachedDisks(v *vm.VirtualMachine) error {
	devices, err := v.Device(d.ctx)
	if err != nil {
		log.Debugf("Couldn't find any devices to detach: %s", err.Error())
		return nil
	}

	disks := devices.SelectByType(&types.VirtualDisk{})
	if disks == nil {
		// nothing attached
		log.Debugf("No disks found attached to VM")
		return nil
	}

	config := []types.BaseVirtualDeviceConfigSpec{}
	for _, disk := range disks {
		config = append(config,
			&types.VirtualDeviceConfigSpec{
				Device:    disk,
				Operation: types.VirtualDeviceConfigSpecOperationRemove,
			})
	}

	op := trace.NewOperation(d.ctx, "detach disks before delete")
	_, err = v.WaitForResult(op,
		func(ctx context.Context) (tasks.Task, error) {
			t, er := v.Reconfigure(ctx,
				types.VirtualMachineConfigSpec{DeviceChange: config})
			if t != nil {
				op.Debugf("Detach reconfigure task=%s", t.Reference())
			}
			return t, er
		})

	return err
}

func (d *Dispatcher) DeleteVCHInstances(vmm *vm.VirtualMachine, conf *config.VirtualContainerHostConfigSpec, containers *DeleteContainers) error {
	defer trace.End(trace.Begin(conf.Name))

	deletePoweredOnContainers := d.force || (containers != nil && *containers == AllContainers)
	ignoreFailureToFindImageStores := d.force

	log.Infof("Removing VMs")

	// serializes access to errs
	var mu sync.Mutex
	var errs []string

	var err error
	var children []*vm.VirtualMachine
	d.parentResourcepool, err = d.getComputeResource(vmm, conf)
	if err != nil {
		return err
	}

	if children, err = d.parentResourcepool.GetChildrenVMs(d.ctx, d.session); err != nil {
		return err
	}

	if d.session.Datastore, err = d.getImageDatastore(vmm, conf, ignoreFailureToFindImageStores); err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, child := range children {
		//Leave VCH appliance there until everything else is removed, cause it has VCH configuration. Then user could retry delete in case of any failure.
		ok, err := d.isVCH(child)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		if ok {
			// child is vch; detach all attached disks so later removal of images is successful
			if err = d.detachAttachedDisks(child); err != nil {
				errs = append(errs, err.Error())
			}
			continue
		}

		wg.Add(1)
		go func(child *vm.VirtualMachine) {
			defer wg.Done()
			if err = d.deleteVM(child, deletePoweredOnContainers); err != nil {
				mu.Lock()
				errs = append(errs, err.Error())
				mu.Unlock()
			}
		}(child)
	}
	wg.Wait()

	if len(errs) > 0 {
		log.Debugf("Error deleting container VMs %s", errs)
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (d *Dispatcher) deleteNetworkDevices(vmm *vm.VirtualMachine, conf *config.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin(conf.Name))

	log.Infof("Removing appliance VM network devices")

	power, err := vmm.PowerState(d.ctx)
	if err != nil {
		log.Errorf("Failed to get vm power status %q: %s", vmm.Reference(), err)
		return err

	}
	if power != types.VirtualMachinePowerStatePoweredOff {
		if _, err = vmm.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
			return vmm.PowerOff(ctx)
		}); err != nil {
			log.Errorf("Failed to power off existing appliance for %s", err)
			return err
		}
	}

	devices, err := d.networkDevices(vmm)
	if err != nil {
		log.Errorf("Unable to get network devices: %s", err)
		return err
	}

	if len(devices) == 0 {
		log.Infof("No network device attached")
		return nil
	}
	// remove devices
	return vmm.RemoveDevice(d.ctx, false, devices...)
}

func (d *Dispatcher) networkDevices(vmm *vm.VirtualMachine) ([]types.BaseVirtualDevice, error) {
	defer trace.End(trace.Begin(""))

	var err error
	vmDevices, err := vmm.Device(d.ctx)
	if err != nil {
		log.Errorf("Failed to get vm devices for appliance: %s", err)
		return nil, err
	}
	var devices []types.BaseVirtualDevice
	for _, device := range vmDevices {
		if _, ok := device.(types.BaseVirtualEthernetCard); ok {
			devices = append(devices, device)
		}
	}
	return devices, nil
}
