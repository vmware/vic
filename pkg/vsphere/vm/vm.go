// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"strconv"
	"strings"
	"sync/atomic"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
)

const UpdateStatus = "UpdateInProgress"

type InvalidState struct {
	r types.ManagedObjectReference
}

func (i *InvalidState) Error() string {
	return fmt.Sprintf("vm %s is invalid", i.r.String())
}

// VirtualMachine struct defines the VirtualMachine which provides additional
// VIC specific methods over object.VirtualMachine as well as keeps some state
type VirtualMachine struct {
	// TODO: Wrap Internal VirtualMachine struct when we have it
	// *internal.VirtualMachine

	*object.VirtualMachine

	*session.Session

	// fxing is 1 means this vm is fixing for it's in invalid status. 0 means not in fixing status
	fixing int32
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

	if err := vm.Properties(ctx, vm.Reference(), []string{"config.files.vmPathName"}, &mvm); err != nil {
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
		log.Errorf("Unable to get vm config: %s", err)
		return nil, err
	}

	return mvm.Config.ExtraConfig, nil
}

func (vm *VirtualMachine) FetchExtraConfig(ctx context.Context) (map[string]string, error) {
	var err error

	var mvm mo.VirtualMachine
	info := make(map[string]string)

	if err = vm.Properties(ctx, vm.Reference(), []string{"config.extraConfig"}, &mvm); err != nil {
		log.Errorf("Unable to get vm config: %s", err)
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
					msg := "powered off"
					if v == types.VirtualMachinePowerStateSuspended {
						// Unlikely, but possible if the VM was suspended out-of-band
						msg = string(v)
					}
					poweredOff = fmt.Errorf("container VM has unexpectedly %s", msg)
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

// GetCurrentSnapshotTreeByName returns current snapshot, with tree information
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

// Finds a snapshot tree based on comparator function 'compare' via a breadth first search of the snapshot tree attached to the VM
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

// helper func that returns true if node is an upgrade snapshot image
func IsUpgradeSnapshot(node *types.VirtualMachineSnapshotTree, upgradePrefix string) bool {
	return node != nil && strings.HasPrefix(node.Name, upgradePrefix)
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

func (vm *VirtualMachine) IsFixing() bool {
	return vm.fixing > 0
}

func (vm *VirtualMachine) EnterFixingState() {
	atomic.AddInt32(&vm.fixing, 1)
}

func (vm *VirtualMachine) LeaveFixingState() {
	atomic.StoreInt32(&vm.fixing, 0)
}

// FixInvalidState fix vm invalid state issue through unregister & register
func (vm *VirtualMachine) fixVM(ctx context.Context) error {
	log.Debugf("Fix invalid state VM: %s", vm.Reference())
	folders, err := vm.Session.Datacenter.Folders(ctx)
	if err != nil {
		log.Errorf("Unable to get vm folder: %s", err)
		return err
	}

	properties := []string{"summary.config", "summary.runtime.host", "resourcePool", "parentVApp"}
	log.Debugf("Get vm properties %s", properties)
	var mvm mo.VirtualMachine
	if err = vm.VirtualMachine.Properties(ctx, vm.Reference(), properties, &mvm); err != nil {
		log.Errorf("Unable to get vm properties: %s", err)
		return err
	}

	name := mvm.Summary.Config.Name
	log.Debugf("Unregister VM %s", name)
	vm.EnterFixingState()
	if err := vm.Unregister(ctx); err != nil {
		log.Errorf("Unable to unregister vm %q: %s", name, err)

		// Leave fixing state since it will not be reset in the remove event handler
		vm.LeaveFixingState()
		return err
	}

	task, err := vm.registerVM(ctx, mvm.Summary.Config.VmPathName, name, mvm.ParentVApp, mvm.ResourcePool, mvm.Summary.Runtime.Host, folders.VmFolder)
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

func (vm *VirtualMachine) needsFix(ctx context.Context, err error) bool {
	if err == nil {
		return false
	}
	if vm.IsInvalidState(ctx) {
		log.Debugf("vm %s is invalid", vm.Reference())
		return true
	}
	log.Debugf("Do not fix non invalid state error")
	return false
}

func (vm *VirtualMachine) IsInvalidState(ctx context.Context) bool {
	var o mo.VirtualMachine
	if err := vm.VirtualMachine.Properties(ctx, vm.Reference(), []string{"summary.runtime.connectionState"}, &o); err != nil {
		log.Debugf("Failed to get vm properties: %s", err)
		return false
	}
	if o.Summary.Runtime.ConnectionState == types.VirtualMachineConnectionStateInvalid {
		return true
	}
	return false
}

// WaitForResult is designed to handle VM invalid state error for any VM operations.
// It will call tasks.WaitForResult to retry if there is task in progress error.
func (vm *VirtualMachine) WaitForResult(ctx context.Context, f func(context.Context) (tasks.Task, error)) (*types.TaskInfo, error) {
	info, err := tasks.WaitForResult(ctx, f)
	if err == nil || !vm.needsFix(ctx, err) {
		return info, err
	}

	log.Debugf("Try to fix task failure %s", err)
	if nerr := vm.fixVM(ctx); nerr != nil {
		log.Errorf("Failed to fix task failure: %s", nerr)
		return info, err
	}
	log.Debugf("Fixed")
	return tasks.WaitForResult(ctx, f)
}

func (vm *VirtualMachine) Properties(ctx context.Context, r types.ManagedObjectReference, ps []string, o *mo.VirtualMachine) error {
	log.Debugf("get vm properties %s of vm %s", ps, r)
	contains := false
	for i := range ps {
		if ps[i] == "summary" || ps[i] == "summary.runtime" {
			contains = true
			break
		}
	}
	var newps []string
	if !contains {
		newps = append(ps, "summary.runtime.connectionState")
	} else {
		newps = append(newps, ps...)
	}
	log.Debugf("properties: %s", newps)
	if err := vm.VirtualMachine.Properties(ctx, r, newps, o); err != nil {
		return err
	}
	if o.Summary.Runtime.ConnectionState != types.VirtualMachineConnectionStateInvalid {
		return nil
	}
	log.Infof("vm %s is in invalid state", r)
	if err := vm.fixVM(ctx); err != nil {
		log.Errorf("Failed to fix vm %s: %s", vm.Reference(), err)
		return &InvalidState{r: vm.Reference()}
	}
	log.Debugf("Retry properties query %s of vm %s", ps, vm.Reference())
	return vm.VirtualMachine.Properties(ctx, vm.Reference(), ps, o)
}

func (vm *VirtualMachine) Parent(ctx context.Context) (*types.ManagedObjectReference, error) {
	var mvm mo.VirtualMachine

	if err := vm.Properties(ctx, vm.Reference(), []string{"parentVApp", "resourcePool"}, &mvm); err != nil {
		log.Errorf("Unable to get VM parent: %s", err)
		return nil, err
	}
	if mvm.ParentVApp != nil {
		return mvm.ParentVApp, nil
	}
	return mvm.ResourcePool, nil
}

func (vm *VirtualMachine) DatastoreReference(ctx context.Context) ([]types.ManagedObjectReference, error) {
	var mvm mo.VirtualMachine

	if err := vm.Properties(ctx, vm.Reference(), []string{"datastore"}, &mvm); err != nil {
		log.Errorf("Unable to get VM datastore: %s", err)
		return nil, err
	}
	return mvm.Datastore, nil
}

// VCHUpdateStatus tells if an upgrade/configure has already been started based on the UpdateInProgress flag in ExtraConfig
// It returns the error if the vm operation does not succeed
func (vm *VirtualMachine) VCHUpdateStatus(ctx context.Context) (bool, error) {
	info, err := vm.FetchExtraConfig(ctx)
	if err != nil {
		log.Errorf("Unable to get vm ExtraConfig: %s", err)
		return false, err
	}

	if v, ok := info[UpdateStatus]; ok {
		status, err := strconv.ParseBool(v)
		if err != nil {
			//  If error occurs, the bool return value does not matter for the caller.
			return false, fmt.Errorf("failed to parse %s to bool: %s", v, err)
		}
		return status, nil
	}

	// If UpdateStatus is not found, it might be the case that no upgrade/configure has been done to this VCH before
	return false, nil
}

// SetVCHUpdateStatus sets the VCH update status in ExtraConfig
func (vm *VirtualMachine) SetVCHUpdateStatus(ctx context.Context, status bool) error {
	info := make(map[string]string)
	info[UpdateStatus] = strconv.FormatBool(status)

	s := &types.VirtualMachineConfigSpec{
		ExtraConfig: vmomi.OptionValueFromMap(info),
	}

	_, err := vm.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		return vm.Reconfigure(ctx, *s)
	})

	return err
}
