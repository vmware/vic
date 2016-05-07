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
	return &VirtualMachine{
		VirtualMachine: object.NewVirtualMachine(
			session.Vim25(),
			moref,
		),
		Session: session,
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

func (vm VirtualMachine) WaitForMAC(ctx context.Context) (map[string]string, error) {
	macs := make(map[string]string)

	p := property.DefaultCollector(vm.Session.Vim25())

	// Wait for all NICs to have a MacAddress, which may not be generated yet.
	err := property.Wait(ctx, p, vm.Reference(), []string{"config.hardware.device"}, func(pc []types.PropertyChange) bool {
		for _, c := range pc {
			if c.Op != types.PropertyChangeOpAssign {
				continue
			}

			devices := c.Val.(types.ArrayOfVirtualDevice).VirtualDevice
			for _, device := range devices {
				if nic, ok := device.(types.BaseVirtualEthernetCard); ok {
					mac := nic.GetVirtualEthernetCard().MacAddress
					if mac == "" {
						return false
					}
					summary := nic.GetVirtualEthernetCard().DeviceInfo.GetDescription().Summary
					macs[summary] = mac
				}
			}
		}

		return true
	})
	return macs, err
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
