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
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

func (d *Dispatcher) NewVCHFromID(id string) (*vm.VirtualMachine, error) {
	defer trace.End(trace.Begin(id))

	var err error
	var vmm *vm.VirtualMachine

	moref := new(types.ManagedObjectReference)
	if ok := moref.FromString(id); !ok {
		message := "Failed to get appliance VM mob reference"
		log.Errorf(message)
		return nil, errors.New(message)
	}
	ref, err := d.session.Finder.ObjectReference(d.ctx, *moref)
	if err != nil {
		if _, ok := err.(*find.NotFoundError); !ok {
			err = errors.Errorf("Failed to query appliance (%q): %s", moref, err)
			return nil, err
		}
		log.Debugf("Appliance is not found")
		return nil, nil
	}
	ovm, ok := ref.(*object.VirtualMachine)
	if !ok {
		log.Errorf("Failed to find VM %q: %s", moref, err)
		return nil, err
	}
	vmm = vm.NewVirtualMachine(d.ctx, d.session, ovm.Reference())

	// check if it's VCH
	if ok, err = d.isVCH(vmm); err != nil {
		log.Error(err)
		return nil, err
	}
	if !ok {
		err = errors.Errorf("Not a VCH")
		log.Error(err)
		return nil, err
	}
	return vmm, nil
}

func (d *Dispatcher) NewVCHFromComputePath(computePath string, name string, v *validate.Validator) (*vm.VirtualMachine, string, error) {
	defer trace.End(trace.Begin(fmt.Sprintf("path %q, name %q", computePath, name)))

	var err error

	parent, err := v.ResourcePoolHelper(d.ctx, computePath)
	if err != nil {
		return nil, "", err
	}
	d.vchPoolPath = path.Join(parent.InventoryPath, name)
	var vchPool *object.ResourcePool
	if d.isVC {
		vapp, err := d.findVirtualApp(d.vchPoolPath)
		if err != nil {
			log.Errorf("Failed to get VCH virtual app %q: %s", d.vchPoolPath, err)
			return nil, d.vchPoolPath, err
		}
		if vapp != nil {
			vchPool = vapp.ResourcePool
		}
	}
	if vchPool == nil {
		vchPool, err = d.session.Finder.ResourcePool(d.ctx, d.vchPoolPath)
		if err != nil {
			log.Errorf("Failed to get VCH resource pool %q: %s", d.vchPoolPath, err)
			return nil, d.vchPoolPath, err
		}
	}

	rp := compute.NewResourcePool(d.ctx, d.session, vchPool.Reference())
	var vmm *vm.VirtualMachine
	if vmm, err = rp.GetChildVM(d.ctx, d.session, name); err != nil {
		log.Errorf("Failed to get VCH VM: %s", err)
		return nil, d.vchPoolPath, err
	}
	if vmm == nil {
		err = errors.Errorf("Didn't find VM %q in resource pool %q", name, rp.Name())
		log.Error(err)
		return nil, d.vchPoolPath, err
	}

	// check if it's VCH
	var ok bool
	if ok, err = d.isVCH(vmm); err != nil {
		log.Error(err)
		return nil, d.vchPoolPath, err
	}
	if !ok {
		err = errors.Errorf("Not a VCH")
		log.Error(err)
		return nil, d.vchPoolPath, err
	}
	return vmm, d.vchPoolPath, nil
}

func (d *Dispatcher) GetVCHConfig(vm *vm.VirtualMachine) (*config.VirtualContainerHostConfigSpec, error) {
	defer trace.End(trace.Begin(""))

	//this is the appliance vm
	mapConfig, err := vm.FetchExtraConfig(d.ctx)
	if err != nil {
		err = errors.Errorf("Failed to get VM extra config of %q: %s", vm.Reference(), err)
		log.Error(err)
		return nil, err
	}
	data := extraconfig.MapSource(mapConfig)
	vchConfig := &config.VirtualContainerHostConfigSpec{}
	result := extraconfig.Decode(data, vchConfig)
	if result == nil {
		err = errors.Errorf("Failed to decode VM configuration %q: %s", vm.Reference(), err)
		log.Error(err)
		return nil, err
	}

	//	vchConfig.ID
	return vchConfig, nil
}
