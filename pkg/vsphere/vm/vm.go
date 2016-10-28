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
	"container/list"
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
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
	return NewVirtualMachineFromVM(ctx, session, object.NewVirtualMachine(session.Vim25(), moref))
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
	u, err := vm.DSPath(ctx)
	if err != nil {
		return "", err
	}

	return path.Base(u.Path), nil
}

// DSPath returns the full datastore path of the VM as a url. The datastore name is in the host
// portion, the path is in the Path field, the scheme is set to "ds"
func (vm *VirtualMachine) DSPath(ctx context.Context) (url.URL, error) {
	var mvm mo.VirtualMachine

	if err := vm.Properties(ctx, vm.Reference(), []string{"runtime.host", "config"}, &mvm); err != nil {
		log.Errorf("Unable to get managed config for VM: %s", err)
		return url.URL{}, err
	}

	if mvm.Config == nil {
		return url.URL{}, errors.New("failed to get datastore path - config not found")
	}
	path := path.Dir(mvm.Config.Files.VmPathName)
	val := url.URL{
		Scheme: "ds",
	}

	// split the dsPath into the url components
	if ix := strings.Index(path, "] "); ix != -1 {
		val.Host = path[strings.Index(path, "[")+1 : ix]
		val.Path = path[ix+2:]
	}

	return val, nil
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
		pgref := object.NewDistributedVirtualPortgroup(vm.Session.Vim25(), types.ManagedObjectReference{
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

func (vm *VirtualMachine) FetchExtraConfigBaseOptions(ctx context.Context) ([]types.BaseOptionValue, error) {
	var err error

	var mvm mo.VirtualMachine

	if err = vm.Properties(ctx, vm.Reference(), []string{"config.extraConfig"}, &mvm); err != nil {
		log.Infof("Unable to get vm config: %s", err)
		return nil, err
	}

	return mvm.Config.ExtraConfig, nil
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
	err := property.Wait(ctx, p, vm.Reference(), []string{object.PropRuntimePowerState, "config.extraConfig"}, waitFunc)
	if err != nil {
		log.Errorf("Property collector error: %s", err)
		return err
	}
	return nil
}

func (vm *VirtualMachine) WaitForKeyInExtraConfig(ctx context.Context, key string) (string, error) {
	var detail string
	var poweredOff error

	waitFunc := func(pc []types.PropertyChange) bool {
		for _, c := range pc {
			if c.Op != types.PropertyChangeOpAssign {
				continue
			}

			switch v := c.Val.(type) {
			case types.ArrayOfOptionValue:
				for _, value := range v.OptionValue {
					// check the status of the key and return true if it's been set to non-nil
					if key == value.GetOptionValue().Key {
						detail = value.GetOptionValue().Value.(string)
						if detail != "" && detail != "<nil>" {
							return true
						}
						break // continue the outer loop as we may have a powerState change too
					}
				}
			case types.VirtualMachinePowerState:
				if v != types.VirtualMachinePowerStatePoweredOn {
					// Give up if the vm has powered off
					poweredOff = fmt.Errorf("%s=%s", c.Name, v)
					return true
				}
			}

		}
		return false
	}

	err := vm.WaitForExtraConfig(ctx, waitFunc)
	if err == nil && poweredOff != nil {
		err = poweredOff
	}

	if err != nil {
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

// DeleteExceptDisks destroys the VM after detaching all virtual disks
func (vm *VirtualMachine) DeleteExceptDisks(ctx context.Context) (*object.Task, error) {
	devices, err := vm.Device(ctx)
	if err != nil {
		return nil, err
	}

	disks := devices.SelectByType(&types.VirtualDisk{})

	err = vm.RemoveDevice(ctx, true, disks...)
	if err != nil {
		return nil, err
	}

	return vm.Destroy(ctx)
}

func (vm *VirtualMachine) VMPathName(ctx context.Context) (string, error) {
	var err error
	var mvm mo.VirtualMachine

	if err = vm.Properties(ctx, vm.Reference(), []string{"config.files"}, &mvm); err != nil {
		log.Errorf("Unable to get vm config.files property: %s", err)
		return "", err
	}

	return mvm.Config.Files.VmPathName, nil
}

// RemoveSnapshot delete one snapshot
func (vm *VirtualMachine) RemoveSnapshot(ctx context.Context, id types.ManagedObjectReference, removeChildren bool, consolidate bool) (*object.Task, error) {
	req := types.RemoveSnapshot_Task{
		This:           id,
		RemoveChildren: removeChildren,
		Consolidate:    &consolidate,
	}
	res, err := methods.RemoveSnapshot_Task(ctx, vm.Client.RoundTripper, &req)
	if err != nil {
		return nil, err
	}

	return object.NewTask(vm.Vim25(), res.Returnval), nil
}

// GetCurrentSnapshotTree returns current snapshot, with tree information
func (vm *VirtualMachine) GetCurrentSnapshotTree(ctx context.Context) (*types.VirtualMachineSnapshotTree, error) {
	var err error

	var mvm mo.VirtualMachine

	if err = vm.Properties(ctx, vm.Reference(), []string{"snapshot"}, &mvm); err != nil {
		log.Infof("Unable to get vm properties: %s", err)
		return nil, err
	}
	if mvm.Snapshot == nil {
		// no snapshot at all
		return nil, nil
	}

	current := mvm.Snapshot.CurrentSnapshot
	q := list.New()
	for _, c := range mvm.Snapshot.RootSnapshotList {
		q.PushBack(c)
	}

	compareID := func(node types.VirtualMachineSnapshotTree) bool {
		if node.Snapshot == *current {
			return true
		}
		return false
	}
	return vm.bfsSnapshotTree(q, compareID), nil
}

// GetCurrentSnapshotTree returns current snapshot, with tree information
func (vm *VirtualMachine) GetSnapshotTreeByName(ctx context.Context, name string) (*types.VirtualMachineSnapshotTree, error) {
	var err error

	var mvm mo.VirtualMachine

	if err = vm.Properties(ctx, vm.Reference(), []string{"snapshot"}, &mvm); err != nil {
		log.Infof("Unable to get vm properties: %s", err)
		return nil, err
	}
	if mvm.Snapshot == nil {
		// no snapshot at all
		return nil, nil
	}

	q := list.New()
	for _, c := range mvm.Snapshot.RootSnapshotList {
		q.PushBack(c)
	}

	compareName := func(node types.VirtualMachineSnapshotTree) bool {
		if node.Name == name {
			return true
		}
		return false
	}
	return vm.bfsSnapshotTree(q, compareName), nil
}

func (vm *VirtualMachine) bfsSnapshotTree(q *list.List, compare func(node types.VirtualMachineSnapshotTree) bool) *types.VirtualMachineSnapshotTree {
	if q.Len() == 0 {
		return nil
	}

	e := q.Front()
	tree := q.Remove(e).(types.VirtualMachineSnapshotTree)
	if compare(tree) {
		return &tree
	}
	for _, c := range tree.ChildSnapshotList {
		q.PushBack(c)
	}
	return vm.bfsSnapshotTree(q, compare)
}

// UpgradeInProgress tells if an upgrade has already been started based on snapshot name beginning with upgradePrefix
func (vm *VirtualMachine) UpgradeInProgress(ctx context.Context, upgradePrefix string) (bool, string, error) {
	node, err := vm.GetCurrentSnapshotTree(ctx)
	if err != nil {
		return false, "", fmt.Errorf("Failed to check upgrade snapshot status: %s", err)
	}

	if node != nil && strings.HasPrefix(node.Name, upgradePrefix) {
		return true, node.Name, nil
	}

	return false, "", nil
}

func (vm *VirtualMachine) registerVM(ctx context.Context, path, name string,
	vapp, pool, host *types.ManagedObjectReference, vmfolder *object.Folder) (*object.Task, error) {
	log.Debugf("Register VM %s", name)

	if vapp == nil {
		var hostObject *object.HostSystem
		if host != nil {
			hostObject = object.NewHostSystem(vm.Vim25(), *host)
		}
		poolObject := object.NewResourcePool(vm.Vim25(), *pool)
		return vmfolder.RegisterVM(ctx, path, name, false, poolObject, hostObject)
	}

	req := types.RegisterChildVM_Task{
		This: vapp.Reference(),
		Path: path,
		Host: host,
	}

	if name != "" {
		req.Name = name
	}

	res, err := methods.RegisterChildVM_Task(ctx, vm.Vim25(), &req)
	if err != nil {
		return nil, err
	}

	return object.NewTask(vm.Vim25(), res.Returnval), nil
}

// FixInvalidState fix vm invalid state issue through unregister & register
func (vm *VirtualMachine) FixInvalidState(ctx context.Context) error {
	log.Debugf("Fix invalid state VM: %s", vm.Reference())
	folders, err := vm.Session.Datacenter.Folders(ctx)
	if err != nil {
		log.Errorf("Unable to get vm folder: %s", err)
		return err
	}

	properties := []string{"config.files", "summary.config", "summary.runtime", "resourcePool", "parentVApp"}
	log.Debugf("Get vm properties %s", properties)
	var mvm mo.VirtualMachine
	if err = vm.Properties(ctx, vm.Reference(), properties, &mvm); err != nil {
		log.Errorf("Unable to get vm properties: %s", err)
		return err
	}

	name := mvm.Summary.Config.Name
	log.Debugf("Unregister VM %s", name)
	if err := vm.Unregister(ctx); err != nil {
		log.Errorf("Unable to unregister vm %q: %s", name, err)
		return err
	}

	task, err := vm.registerVM(ctx, mvm.Config.Files.VmPathName, name, mvm.ParentVApp, mvm.ResourcePool, mvm.Summary.Runtime.Host, folders.VmFolder)
	if err != nil {
		log.Errorf("Unable to register VM %q back: %s", name, err)
		return err
	}
	info, err := task.WaitForResult(ctx, nil)
	if err != nil {
		return err
	}
	// re-register vm will change vm reference, so reset the object reference here
	if info.Error != nil {
		return errors.New(info.Error.LocalizedMessage)
	}

	// set new registered vm attribute back
	newRef := info.Result.(types.ManagedObjectReference)
	common := object.NewCommon(vm.Vim25(), newRef)
	common.InventoryPath = vm.InventoryPath
	vm.Common = common
	return nil
}

func (vm *VirtualMachine) needsFix(err error) bool {
	f, ok := err.(types.HasFault)
	if !ok {
		return false
	}
	switch f.Fault().(type) {
	case *types.InvalidState:
		return true
	default:
		log.Debugf("Do not fix non invalid state error")
		return false
	}
}

// WaitForResult is designed to handle VM invalid state error for any VM operations.
// It will call tasks.WaitForResult to retry if there is task in progress error.
func (vm *VirtualMachine) WaitForResult(ctx context.Context, f func(context.Context) (tasks.Task, error)) (*types.TaskInfo, error) {
	info, err := tasks.WaitForResult(ctx, f)
	if err == nil || !vm.needsFix(err) {
		return info, err
	}
	log.Debugf("Try to fix task failure %s", err)
	if nerr := vm.FixInvalidState(ctx); nerr != nil {
		log.Errorf("Failed to fix task failure: %s", nerr)
		return info, err
	}
	log.Debugf("Fixed")
	return tasks.WaitForResult(ctx, f)
}
