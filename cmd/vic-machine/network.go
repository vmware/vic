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

package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/errors"
)

func (v *Validator) createBridgeNetwork() error {
	net, err := v.Session.Finder.Network(v.Context, v.BridgeNetworkPath)
	if err != nil {
		_, ok := err.(*find.NotFoundError)
		if !ok {
			err = errors.Errorf("Failed to query virtual switch (%s): %s", v.BridgeNetworkName, err)
			return err
		}
	}
	if net != nil {
		log.Infof("Network %s exists", v.BridgeNetworkName)
		return nil
	}

	if v.Session.Host == nil {
		err = errors.Errorf("Network creation is not supported in VC, please specify bridge network")
		return err
	}
	log.Infof("Creating VirtualSwitch")
	hostNetSystem, err := v.Session.Host.ConfigManager().NetworkSystem(v.Context)
	if err != nil {
		err = errors.Errorf("Failed to retrieve host network system: %s", err)
		return err
	}

	if err = hostNetSystem.AddVirtualSwitch(v.Context, v.BridgeNetworkName, &types.HostVirtualSwitchSpec{
		NumPorts: 1024,
	}); err != nil {
		err = errors.Errorf("Failed to add virtual switch (%s): %s", v.BridgeNetworkName, err)
		return err
	}

	log.Infof("Creating Portgroup")
	if err = hostNetSystem.AddPortGroup(v.Context, types.HostPortGroupSpec{
		Name:        v.BridgeNetworkName,
		VlanId:      1, // TODO: expose this for finer grained grouping within the switch
		VswitchName: v.BridgeNetworkName,
		Policy:      types.HostNetworkPolicy{},
	}); err != nil {
		err = errors.Errorf("Failed to add port group (%s): %s", v.BridgeNetworkName, err)
		return err
	}

	return nil
}

func (v *Validator) removeNetwork() error {
	if v.Session.Host == nil {
		log.Warnf("Remove network is not supported in VC")
		return nil
	}
	log.Infof("Removing Portgroup")
	hostNetSystem, err := v.Session.Host.ConfigManager().NetworkSystem(v.Context)
	if err != nil {
		return err
	}
	err = hostNetSystem.RemovePortGroup(v.Context, v.BridgeNetworkName)
	if err != nil {
		return err
	}

	log.Infof("Removing VirtualSwitch")
	err = hostNetSystem.RemoveVirtualSwitch(v.Context, v.BridgeNetworkName)
	if err != nil {
		return err
	}
	return nil
}
