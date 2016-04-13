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
	"github.com/vmware/vic/install/configuration"
	"github.com/vmware/vic/pkg/errors"
)

func (d *Dispatcher) createBridgeNetwork(conf *configuration.Configuration) error {
	net, err := d.session.Finder.Network(d.ctx, conf.BridgeNetworkPath)
	if err != nil {
		_, ok := err.(*find.NotFoundError)
		if !ok {
			err = errors.Errorf("Failed to query virtual switch (%s): %s", conf.BridgeNetworkName, err)
			return err
		}
	}
	if net != nil {
		log.Infof("Network %s exists", conf.BridgeNetworkName)
		return nil
	}

	if d.session.Host == nil {
		err = errors.Errorf("Network creation is not supported in VC, please specify bridge network")
		return err
	}
	log.Infof("Creating VirtualSwitch")
	hostNetSystem, err := d.session.Host.ConfigManager().NetworkSystem(d.ctx)
	if err != nil {
		err = errors.Errorf("Failed to retrieve host network system: %s", err)
		return err
	}

	if err = hostNetSystem.AddVirtualSwitch(d.ctx, conf.BridgeNetworkPath, &types.HostVirtualSwitchSpec{
		NumPorts: 1024,
	}); err != nil {
		err = errors.Errorf("Failed to add virtual switch (%s): %s", conf.BridgeNetworkName, err)
		return err
	}

	log.Infof("Creating Portgroup")
	if err = hostNetSystem.AddPortGroup(d.ctx, types.HostPortGroupSpec{
		Name:        conf.BridgeNetworkName,
		VlanId:      1, // TODO: expose this for finer grained grouping within the switch
		VswitchName: conf.BridgeNetworkName,
		Policy:      types.HostNetworkPolicy{},
	}); err != nil {
		err = errors.Errorf("Failed to add port group (%s): %s", conf.BridgeNetworkName, err)
		return err
	}

	return nil
}

func (d *Dispatcher) removeNetwork(conf *configuration.Configuration) error {
	if d.session.Host == nil {
		log.Warnf("Remove network is not supported in VC")
		return nil
	}
	log.Infof("Removing Portgroup")
	hostNetSystem, err := d.session.Host.ConfigManager().NetworkSystem(d.ctx)
	if err != nil {
		return err
	}
	err = hostNetSystem.RemovePortGroup(d.ctx, conf.BridgeNetworkName)
	if err != nil {
		return err
	}

	log.Infof("Removing VirtualSwitch")
	err = hostNetSystem.RemoveVirtualSwitch(d.ctx, conf.BridgeNetworkName)
	if err != nil {
		return err
	}
	return nil
}
