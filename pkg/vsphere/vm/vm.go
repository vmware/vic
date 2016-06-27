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

package vm

import (
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"

	"github.com/vmware/vic/pkg/vsphere/session"
)

// VirtualMachine struct defines the VirtualMachine which provides additional
// VIC specific methods over object.VirtualMachine as well as keeps some state
type VirtualMachine struct {
	// TODO: Wrap Internal VirtualMachine struct when we have it
	// *internal.VirtualMachine

	*object.VirtualMachine

	*session.Session
}

// NewVirtualMachine returns a NewVirtualMachine object
func NewVirtualMachine(ctx context.Context, session *session.Session, moref types.ManagedObjectReference) *VirtualMachine {
	return NewVirtualMachineFromVM(ctx, session, object.NewVirtualMachine(session.Client.Client, moref))
}

func NewVirtualMachineFromVM(ctx context.Context, session *session.Session, vm *object.VirtualMachine) *VirtualMachine {
	return &VirtualMachine{
		VirtualMachine: vm,
		Session:        session,
	}
}

// FolderName returns the name of the namespace(vsan) or directory(vmfs) that holds the VM
// this equates to the normal directory that contains the vmx file, stripped of any parent path
func (vm *VirtualMachine) FolderName(ctx context.Context) (string, error) {
	var mvm mo.VirtualMachine

	if err := vm.Properties(ctx, vm.Reference(), []string{"runtime.host", "config"}, &mvm); err != nil {
		log.Errorf("Unable to get managed config for VM: %s", err)
		return "", err
	}

	path := path.Base(path.Dir(mvm.Config.Files.VmPathName))
	if path[0] == '[' {
		path = strings.Split(path, "] ")[1]
	}
	return path, nil
}

// WaitForMac will wait until VM get mac for all attached nics.
// Returns map "Virtual Network Name": "nic MAC address"
func (vm VirtualMachine) WaitForMAC(ctx context.Context) (map[string]string, error) {
	devices, err := vm.Device(ctx)
	if err != nil {
		log.Errorf("Unable to get device listing for VM")
		return nil, err
	}

	nics := devices.SelectByType(&types.VirtualEthernetCard{})
	macs := make(map[string]string)
	// device name:network name
	nicMappings := make(map[string]string)
	for _, nic := range nics {
		if n, ok := nic.(types.BaseVirtualEthernetCard); ok {
			netName, err := vm.getNetworkName(ctx, n)
			if err != nil {
				log.Errorf("failed to get network name: %s", err)
				return nil, err
			}
			macs[netName] = ""

			nicMappings[devices.Name(nic)] = netName
		} else {
			log.Errorf("Failed to get network name of vNIC: %v", nic)
			return nil, err
		}
	}

	p := property.DefaultCollector(vm.Session.Vim25())

	// Wait for all NICs to have a MacAddress, which may not be generated yet.
	err = property.Wait(ctx, p, vm.Reference(), []string{"config.hardware.device"}, func(pc []types.PropertyChange) bool {
		for _, c := range pc {
			if c.Op != types.PropertyChangeOpAssign {
				continue
			}

			changedDevices := c.Val.(types.ArrayOfVirtualDevice).VirtualDevice
			for _, device := range changedDevices {
				if nic, ok := device.(types.BaseVirtualEthernetCard); ok {
					mac := nic.GetVirtualEthernetCard().MacAddress
					if mac == "" {
						continue
					}
					netName := nicMappings[devices.Name(device)]
					macs[netName] = mac
				}
			}
		}
		for key, value := range macs {
			if value == "" {
				log.Debugf("Didn't get mac address for nic on %s, continue", key)
				return false
			}
		}
		return true
	})
	return macs, err
}

func (vm *VirtualMachine) getNetworkName(ctx context.Context, nic types.BaseVirtualEthernetCard) (string, error) {
	if card, ok := nic.GetVirtualEthernetCard().Backing.(*types.VirtualEthernetCardDistributedVirtualPortBackingInfo); ok {
		pg := card.Port.PortgroupKey
		pgref := object.NewDistributedVirtualPortgroup(vm.Session.Client.Client, types.ManagedObjectReference{
			Type:  "DistributedVirtualPortgroup",
			Value: pg,
		})

		var pgo mo.DistributedVirtualPortgroup
		err := pgref.Properties(ctx, pgref.Reference(), []string{"config"}, &pgo)
		if err != nil {
			log.Errorf("Failed to query portgroup %s for %s", pg, err)
			return "", err
		}
		return pgo.Config.Name, nil
	}
	return nic.GetVirtualEthernetCard().DeviceInfo.GetDescription().Summary, nil
}

func (vm *VirtualMachine) FetchExtraConfig(ctx context.Context) (map[string]string, error) {
	var err error

	var mvm mo.VirtualMachine
	info := make(map[string]string)

	if err = vm.Properties(ctx, vm.Reference(), []string{"config.extraConfig"}, &mvm); err != nil {
		log.Infof("Unable to get vm config: %s", err)
		return info, err
	}

	for _, bov := range mvm.Config.ExtraConfig {
		ov := bov.GetOptionValue()
		value, _ := ov.Value.(string)
		info[ov.Key] = value
	}
	return info, nil
}

// WaitForExtraConfig waits until key shows up with the expected value inside the ExtraConfig
func (vm *VirtualMachine) WaitForExtraConfig(ctx context.Context, waitFunc func(pc []types.PropertyChange) bool) error {
	// Get the default collector
	p := property.DefaultCollector(vm.Vim25())

	// Wait on config.extraConfig
	// https://www.vmware.com/support/developer/vc-sdk/visdk2xpubs/ReferenceGuide/vim.vm.ConfigInfo.html
	err := property.Wait(ctx, p, vm.Reference(), []string{"config.extraConfig"}, waitFunc)
	if err != nil {
		log.Errorf("Property collector error: %s", err)
		return err
	}
	return nil
}

func (vm *VirtualMachine) WaitForKeyInExtraConfig(ctx context.Context, key string) (string, error) {
	var detail string
	waitFunc := func(pc []types.PropertyChange) bool {
		for _, c := range pc {
			if c.Op != types.PropertyChangeOpAssign {
				continue
			}

			values := c.Val.(types.ArrayOfOptionValue).OptionValue
			for _, value := range values {
				// check the status of the key and return true if it's been set to non-nil
				if key == value.GetOptionValue().Key {
					detail = value.GetOptionValue().Value.(string)
					return detail != "" && detail != "<nil>"
				}
			}
		}
		return false
	}

	if err := vm.WaitForExtraConfig(ctx, waitFunc); err != nil {
		log.Errorf("Unable to wait for extra config property %s: %s", key, err.Error())
		return "", err
	}
	return detail, nil
}

func (vm *VirtualMachine) Name(ctx context.Context) (string, error) {
	var err error
	var mvm mo.VirtualMachine

	if err = vm.Properties(ctx, vm.Reference(), []string{"summary.config"}, &mvm); err != nil {
		log.Errorf("Unable to get vm summary.config property: %s", err)
		return "", err
	}

	return mvm.Summary.Config.Name, nil
}

func (vm *VirtualMachine) UUID(ctx context.Context) (string, error) {
	var err error
	var mvm mo.VirtualMachine

	if err = vm.Properties(ctx, vm.Reference(), []string{"summary.config"}, &mvm); err != nil {
		log.Errorf("Unable to get vm summary.config property: %s", err)
		return "", err
	}

	return mvm.Summary.Config.Uuid, nil
}
