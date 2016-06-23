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
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

func (d *Dispatcher) DeleteVCH(conf *metadata.VirtualContainerHostConfigSpec) error {
	log.Infof("Removing VMs")

	if err := d.DeleteVCHInstances(conf); err != nil {
		return err
	}
	if err := d.destroyResourcePoolIfEmpty(conf); err != nil {
		return err
	}
	return d.removeNetwork(conf)
}

func (d *Dispatcher) DeleteVCHInstances(conf *metadata.VirtualContainerHostConfigSpec) error {
	var errs []string

	var err error
	var vmm *vm.VirtualMachine
	var children []*vm.VirtualMachine

	if vmm, err = d.findApplianceByID(conf); err != nil {
		return err
	}
	if vmm == nil {
		return nil
	}
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
		if _, err = d.deleteVM(child); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		log.Debugf("Error deleting container VMs %s", errs)
		return errors.New(strings.Join(errs, "\n"))
	}

	if _, err = d.deleteVM(vmm); err != nil {
		log.Debugf("Error deleting appliance VM %s", err)
		return err
	}

	return nil
}
