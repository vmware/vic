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
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	"golang.org/x/net/context"
)

func (d *Dispatcher) DeleteVCH(conf *metadata.VirtualContainerHostConfigSpec) error {
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

	if err = d.DeleteDataStores(vmm, conf); err != nil {
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
	folder, err := d.deleteVM(vmm, true)
	if err != nil {
		log.Debugf("Error deleting appliance VM %s", err)
		return err
	}
	if _, err = d.deleteDatastoreFiles(d.session.Datastore, folder, true); err != nil {
		log.Warnf("Appliance VM path %s is not removed, %s", folder, err)
	}
	if err = d.destroyResourcePoolIfEmpty(conf); err != nil {
		log.Warnf("VCH resource pool is not removed, %s", err)
	}
	return nil
}

func (d *Dispatcher) DeleteVCHInstances(vmm *vm.VirtualMachine, conf *metadata.VirtualContainerHostConfigSpec) error {
	log.Infof("Removing VMs")
	var errs []string

	var err error
	var children []*vm.VirtualMachine

	rpRef := conf.ComputeResources[len(conf.ComputeResources)-1]
	rp := compute.NewResourcePool(d.ctx, d.session, rpRef)
	if children, err = rp.GetChildrenVMs(d.ctx, d.session); err != nil {
		return err
	}

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
		if _, err = d.deleteVM(child, d.force); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		log.Debugf("Error deleting container VMs %s", errs)
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (d *Dispatcher) deleteNetworkDevices(vmm *vm.VirtualMachine, conf *metadata.VirtualContainerHostConfigSpec) error {
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

	spec, err := d.deleteNetworkSpec(vmm, conf)
	if err != nil {
		log.Errorf("Unable to create deleting spec: %s", err)
		return err
	}

	if spec == nil {
		log.Infof("No network device attached")
		return nil
	}
	// reconfig
	info, err := tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return vmm.Reconfigure(ctx, *spec)
	})

	if err != nil {
		log.Errorf("Error while removing network devices from appliance VM: %s", err)
		return err
	}
	if info.State != types.TaskInfoStateSuccess {
		log.Errorf("Removing network devices reported: %s", info.Error.LocalizedMessage)
		return err
	}
	return nil
}

func (d *Dispatcher) deleteNetworkSpec(vmm *vm.VirtualMachine, conf *metadata.VirtualContainerHostConfigSpec) (*types.VirtualMachineConfigSpec, error) {
	var err error

	spec := &types.VirtualMachineConfigSpec{
		Name:    conf.Name,
		GuestId: "other3xLinux64Guest",
		Files:   &types.VirtualMachineFileInfo{VmPathName: fmt.Sprintf("[%s]", conf.ImageStores[0].Host)},
	}

	vmDevices, err := vmm.Device(d.ctx)
	if err != nil {
		log.Errorf("Failed to get vm devices for appliance: %s", err)
		return nil, err
	}
	var found bool
	for _, device := range vmDevices {
		if _, ok := device.(types.BaseVirtualEthernetCard); ok {
			found = true
			spec.DeviceChange = append(spec.DeviceChange,
				&types.VirtualDeviceConfigSpec{
					Operation: types.VirtualDeviceConfigSpecOperationRemove,
					Device:    device,
				},
			)

		}
	}
	if found {
		return spec, nil
	}
	return nil, nil
}
