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
	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
)

func (d *Dispatcher) createBridgeNetwork(conf *metadata.VirtualContainerHostConfigSpec) error {
	// if the bridge network is already extant there's nothing to do
	bnet := conf.ContainerNetworks[conf.DefaultNetwork]
	if bnet != nil && bnet.ID != "" {
		return nil
	}

	// network didn't exist during validation given we don't have a moref, so create it
	if d.session.Client.IsVC() {
		// double check
		return errors.New("bridge network must already exist for vCenter environments")
	}

	log.Infof("Creating VirtualSwitch")
	hostNetSystem, err := d.session.Host.ConfigManager().NetworkSystem(d.ctx)
	if err != nil {
		err = errors.Errorf("Failed to retrieve host network system: %s", err)
		return err
	}

	if err = hostNetSystem.AddVirtualSwitch(d.ctx, bnet.Name, &types.HostVirtualSwitchSpec{
		NumPorts: 1024,
	}); err != nil {
		err = errors.Errorf("Failed to add virtual switch (%s): %s", bnet.Name, err)
		return err
	}

	log.Infof("Creating Portgroup")
	if err = hostNetSystem.AddPortGroup(d.ctx, types.HostPortGroupSpec{
		Name:        bnet.Name,
		VlanId:      1, // TODO: expose this for finer grained grouping within the switch
		VswitchName: bnet.Name,
		Policy:      types.HostNetworkPolicy{},
	}); err != nil {
		err = errors.Errorf("Failed to add port group (%s): %s", bnet.Name, err)
		return err
	}

	net, err := d.session.Finder.Network(d.ctx, bnet.Name)
	if err != nil {
		_, ok := err.(*find.NotFoundError)
		if !ok {
			err = errors.Errorf("Failed to query virtual switch (%s): %s", bnet.Name, err)
			return err
		}
	}

	// assign the moref to the bridge network config on the appliance
	bnet.ID = net.Reference().String()
	return nil
}

func (d *Dispatcher) removeNetwork(name string) error {
	if d.session.Client.IsVC() {
		log.Warnf("Remove network is not supported for vCenter")
		return nil
	}
	log.Infof("Removing Portgroup")
	hostNetSystem, err := d.session.Host.ConfigManager().NetworkSystem(d.ctx)
	if err != nil {
		return err
	}
	err = hostNetSystem.RemovePortGroup(d.ctx, name)
	if err != nil {
		return err
	}

	log.Infof("Removing VirtualSwitch")
	err = hostNetSystem.RemoveVirtualSwitch(d.ctx, name)
	if err != nil {
		return err
	}
	return nil
}
