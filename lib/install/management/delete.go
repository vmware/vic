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

package management

import (
	"bytes"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	"golang.org/x/net/context"
)

func (d *Dispatcher) DeleteVCH(conf *metadata.VirtualContainerHostConfigSpec) error {
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

	if err = d.DeleteVCHInstances(vmm, conf); err != nil {
		// if container delete failed, do not remove anything else
		log.Infof("Specify --force to force delete")
		return err
	}

	if err = d.DeleteStores(vmm, conf); err != nil {
		errs = append(errs, err.Error())
	}

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

	if d.isVC {
		log.Infoln("Removing VCH vSphere extension")
		if err = d.GenerateExtensionName(conf); err != nil {
			log.Warnf("Failed to get extension name during VCH deletion: %s", err)
		}
		if err = d.UnregisterExtension(conf.ExtensionName); err != nil {
			log.Warnf("Failed to remove extension %s: %s", conf.ExtensionName, err)
		}
	}

	err = d.deleteVM(vmm, true)
	if err != nil {
		log.Debugf("Error deleting appliance VM %s", err)
		return err
	}
	if err = d.destroyResourcePoolIfEmpty(conf); err != nil {
		log.Warnf("VCH resource pool is not removed, %s", err)
	}

	log.Infoln("Removing volume stores...")
	if err = d.deleteVolumeStoreIfForced(conf); err != nil {
		log.Warnf("Error while deleting volume store: %s", err)
		return err
	}
	return nil
}

func (d *Dispatcher) DeleteVCHInstances(vmm *vm.VirtualMachine, conf *metadata.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin(""))

	log.Infof("Removing VMs")
	var errs []string

	var err error
	var children []*vm.VirtualMachine

	rpRef := conf.ComputeResources[len(conf.ComputeResources)-1]
	ref, err := d.session.Finder.ObjectReference(d.ctx, rpRef)
	if err != nil {
		err = errors.Errorf("Failed to get VCH resource pool (%s): %s", rpRef, err)
		return err
	}
	_, ok := ref.(*object.ResourcePool)
	if !ok {
		log.Errorf("Failed to find resource pool %s, %s", rpRef, err)
		return err
	}
	rp := compute.NewResourcePool(d.ctx, d.session, ref.Reference())
	if children, err = rp.GetChildrenVMs(d.ctx, d.session); err != nil {
		return err
	}

	ds, err := d.session.Finder.Datastore(d.ctx, conf.ImageStores[0].Host)
	if err != nil {
		err = errors.Errorf("Failed to find image datastore %s", conf.ImageStores[0].Host)
		return err
	}
	d.session.Datastore = ds

	for _, child := range children {
		name, err := child.Name(d.ctx)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		//Leave VCH appliance there until everything else is removed, cause it has VCH configuration. Then user could retry delete in case of any failure.
		if name == conf.Name {
			continue
		}
		if err = d.deleteVM(child, d.force); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		log.Debugf("Error deleting container VMs %s", errs)
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (d *Dispatcher) deleteVolumeStoreIfForced(conf *metadata.VirtualContainerHostConfigSpec) error {
	if d.force {

		for label, url := range conf.VolumeLocations {

			// separate the host from the path in the provided path URL
			pathAry := strings.SplitN(url.Path, " ", 2)
			if len(pathAry) != 2 {
				return errors.New("Didn't receive an expected volume store path format")
			}

			// convert the URL path to vSphere style pathing, and omit everything from the path except the root of the volume store
			vSpherePath := fmt.Sprintf("%s %s", pathAry[0], strings.Split(pathAry[1], "/")[0])

			// connect to vSphere and do the actual deletion
			m := object.NewFileManager(d.session.Vim25())
			log.Infof("Deleting volume store %s at path %s", label, vSpherePath)
			task, err := m.DeleteDatastoreFile(d.ctx, vSpherePath, d.session.Datacenter)
			if err != nil {
				return errors.Errorf("Failed to start delete of %s due to error: %s", vSpherePath, err)
			}

			if err = task.Wait(d.ctx); err != nil {
				return errors.Errorf("Failed to finish delete of %s due to error: %s", vSpherePath, err)
			}

		}

	} else {
		volumeStores := new(bytes.Buffer)
		for label, url := range conf.VolumeLocations {
			if _, err := volumeStores.WriteString(fmt.Sprintf("\t%s: %s\n", label, url.Path)); err != nil {
				return err
			}
		}
		log.Warnf("Since --force was not specified, the following volume stores will not be removed. Use the vSphere UI to delete content you do not wish to keep.\n%s", volumeStores.String())
	}
	return nil
}

func (d *Dispatcher) deleteNetworkDevices(vmm *vm.VirtualMachine, conf *metadata.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin(""))

	log.Infof("Removing appliance VM network devices")

	power, err := vmm.PowerState(d.ctx)
	if err != nil {
		log.Errorf("Failed to get vm power status %s: %s", vmm.Reference(), err)
		return err

	}
	if power != types.VirtualMachinePowerStatePoweredOff {
		if _, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
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

func (d *Dispatcher) UnregisterExtension(name string) error {
	defer trace.End(trace.Begin(name))

	extensionManager := object.NewExtensionManager(d.session.Vim25())
	if err := extensionManager.Unregister(d.ctx, name); err != nil {
		return errors.Errorf("Failed to remove extension w/ name %s due to error: %s", name, err)
	}
	return nil
}
