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

package management

import (
	"context"
	"fmt"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/opts"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/retry"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

const (
	UpgradePrefix = "upgrade for"
)

// Upgrade will try to upgrade vch appliance to new version. If failed will try to roll back to original status.
func (d *Dispatcher) Upgrade(vch *vm.VirtualMachine, conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) (err error) {
	defer trace.End(trace.Begin(conf.Name))

	d.appliance = vch

	// update the displayname to the actual folder name used
	if d.vmPathName, err = d.appliance.FolderName(d.ctx); err != nil {
		log.Errorf("Failed to get canonical name for appliance: %s", err)
		return err
	}

	ds, err := d.session.Finder.Datastore(d.ctx, conf.ImageStores[0].Host)
	if err != nil {
		err = errors.Errorf("Failed to find image datastore %q", conf.ImageStores[0].Host)
		return err
	}
	d.session.Datastore = ds
	if !conf.HostCertificate.IsNil() {
		d.DockerPort = fmt.Sprintf("%d", opts.DefaultTLSHTTPPort)
	} else {
		d.DockerPort = fmt.Sprintf("%d", opts.DefaultHTTPPort)
	}

	if err = d.uploadImages(settings.ImageFiles); err != nil {
		return errors.Errorf("Uploading images failed with %s. Exiting...", err)
	}

	conf.BootstrapImagePath = fmt.Sprintf("[%s] %s/%s", conf.ImageStores[0].Host, d.vmPathName, settings.BootstrapISO)

	// ensure that we wait for components to come up
	for _, s := range conf.ExecutorConfig.Sessions {
		s.Started = ""
	}

	snapshotName := fmt.Sprintf("%s %s", UpgradePrefix, conf.Version.BuildNumber)
	snapshotName = strings.TrimSpace(snapshotName)

	// check for old snapshot
	oldSnapshot, _ := d.appliance.GetCurrentSnapshotTree(d.ctx)

	if err = d.tryCreateSnapshot(snapshotName, "upgrade snapshot"); err != nil {
		d.deleteUpgradeImages(ds, settings)
		return err
	}

	if err = d.update(conf, settings); err == nil {
		if oldSnapshot != nil && vm.IsUpgradeSnapshot(oldSnapshot, UpgradePrefix) {
			d.retryDeleteSnapshot(oldSnapshot.Name, conf.Name)
		}
		return nil
	}
	log.Errorf("Failed to upgrade: %s", err)
	log.Infof("Rolling back upgrade")

	if rerr := d.rollback(conf, snapshotName, settings); rerr != nil {
		log.Errorf("Failed to revert appliance to snapshot: %s", rerr)
		return err
	}
	log.Infof("Appliance is rolled back to old version")

	d.deleteUpgradeImages(ds, settings)
	d.retryDeleteSnapshot(snapshotName, conf.Name)

	// return the error message for upgrade
	return err
}

func (d *Dispatcher) Rollback(vch *vm.VirtualMachine, conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) error {

	// some setup that is only necessary because we didn't just create a VCH in this case
	d.appliance = vch
	d.setDockerPort(conf, settings)

	notfound := "A VCH version available from before the last upgrade could not be found."
	snapshot, err := d.appliance.GetCurrentSnapshotTree(d.ctx)
	if err != nil {
		return errors.Errorf("%s An error was reported while trying to discover it: %s", notfound, err)
	}

	if snapshot == nil {
		return errors.Errorf("%s No error was reported, so it's possible that this VCH has never been upgraded or the saved previous version was removed out-of-band.", notfound)
	}

	err = d.rollback(conf, snapshot.Name, settings)
	if err != nil {
		return errors.Errorf("could not complete manual rollback: %s", err)
	}

	return d.retryDeleteSnapshot(snapshot.Name, conf.Name)
}

// retryDeleteSnapshot will retry to delete snpashot if there is GenericVmConfigFault returned. This is a workaround for vSAN delete snapshot
func (d *Dispatcher) retryDeleteSnapshot(snapshotName string, applianceName string) error {
	// delete snapshot immediately after snapshot rollback usually fail in vSAN, so have to retry several times
	operation := func() error {
		return d.deleteSnapshot(snapshotName, applianceName)
	}
	var err error
	if err = retry.Do(operation, isSystemError); err != nil {
		log.Errorf("Failed to clean up appliance upgrade snapshot %q: %s.", snapshotName, err)
		log.Errorf("Snapshot %q of appliance virtual machine %q MUST be removed manually before upgrade again", snapshotName, applianceName)
	}
	return err
}

func isSystemError(err error) bool {
	if soap.IsSoapFault(err) {
		if _, ok := soap.ToSoapFault(err).VimFault().(*types.SystemError); ok {
			return true
		}
	}

	if soap.IsVimFault(err) {
		if _, ok := soap.ToVimFault(err).(*types.SystemError); ok {
			return true
		}
	}

	if terr, ok := err.(task.Error); ok {
		if _, ok := terr.Fault().(*types.SystemError); ok {
			return true
		}
	}

	return false
}

func (d *Dispatcher) deleteSnapshot(snapshotName string, applianceName string) error {
	defer trace.End(trace.Begin(snapshotName))
	log.Infof("Deleting upgrade snapshot %q", snapshotName)
	// do clean up aggressively, even the previous operation failed with context deadline exceeded.
	ctx := context.Background()
	if _, err := d.appliance.WaitForResult(ctx, func(ctx context.Context) (tasks.Task, error) {
		consolidate := true
		return d.appliance.RemoveSnapshot(ctx, snapshotName, true, &consolidate)
	}); err != nil {
		return err
	}
	return nil
}

// tryCreateSnapshot try to create upgrade snapshot. It will check if upgrade snapshot already exists. If exists, return error.
// if succeed, return snapshot refID
func (d *Dispatcher) tryCreateSnapshot(name, desc string) error {
	defer trace.End(trace.Begin(name))

	// TODO detect whether another upgrade is in progress & bail if it is.
	// Use solution from https://github.com/vmware/vic/issues/4069 to do this either as part of 4069 or once it's closed

	if _, err := d.appliance.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
		return d.appliance.CreateSnapshot(d.ctx, name, desc, true, false)
	}); err != nil {
		return errors.Errorf("Failed to create upgrade snapshot %q: %s.", name, err)
	}
	return nil
}

func (d *Dispatcher) deleteUpgradeImages(ds *object.Datastore, settings *data.InstallerData) {
	defer trace.End(trace.Begin(""))

	log.Infof("Deleting upgrade images")

	// do clean up aggressively, even the previous operation failed with context deadline exceeded.
	d.ctx = context.Background()

	m := ds.NewFileManager(d.session.Datacenter, true)

	file := ds.Path(path.Join(d.vmPathName, settings.ApplianceISO))
	if err := d.deleteVMFSFiles(m, ds, file); err != nil {
		log.Warnf("Image file %q is not removed for %s. Use the vSphere UI to delete content", file, err)
	}

	file = ds.Path(path.Join(d.vmPathName, settings.BootstrapISO))
	if err := d.deleteVMFSFiles(m, ds, file); err != nil {
		log.Warnf("Image file %q is not removed for %s. Use the vSphere UI to delete content", file, err)
	}
}

func (d *Dispatcher) update(conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) error {
	defer trace.End(trace.Begin(conf.Name))

	power, err := d.appliance.PowerState(d.ctx)
	if err != nil {
		log.Errorf("Failed to get vm power status %q: %s", d.appliance.Reference(), err)
		return err
	}
	if power != types.VirtualMachinePowerStatePoweredOff {
		if _, err = d.appliance.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
			return d.appliance.PowerOff(ctx)
		}); err != nil {
			log.Errorf("Failed to power off appliance: %s", err)
			return err
		}
	}

	if err = d.reconfigVCH(conf, fmt.Sprintf("[%s] %s/%s", conf.ImageStores[0].Host, d.vmPathName, settings.ApplianceISO)); err != nil {
		return err
	}

	if err = d.startAppliance(conf); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(d.ctx, settings.Timeout)
	defer cancel()
	if err = d.CheckServiceReady(ctx, conf, nil); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			//context deadline exceeded, replace returned error message
			err = errors.Errorf("Upgrading VCH exceeded time limit of %s. Please increase the timeout using --timeout to accommodate for a busy vSphere target", settings.Timeout)
		}

		log.Info("\tAPI may be slow to start - please retry with increased timeout using --timeout: %s", err)
		return err
	}
	return nil
}

func (d *Dispatcher) rollback(conf *config.VirtualContainerHostConfigSpec, snapshot string, settings *data.InstallerData) error {
	defer trace.End(trace.Begin(fmt.Sprintf("old appliance iso: %q, snapshot: %q", d.oldApplianceISO, snapshot)))

	// do not power on appliance in this snapshot revert
	log.Infof("Reverting to snapshot %s", snapshot)
	if _, err := d.appliance.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
		return d.appliance.RevertToSnapshot(d.ctx, snapshot, true)
	}); err != nil {
		return errors.Errorf("Failed to roll back upgrade: %s.", err)
	}
	return d.ensureRollbackReady(conf, settings)
}

func (d *Dispatcher) ensureRollbackReady(conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) error {
	defer trace.End(trace.Begin(conf.Name))

	power, err := d.appliance.PowerState(d.ctx)
	if err != nil {
		log.Errorf("Failed to get vm power status %q after rollback: %s", d.appliance.Reference(), err)
		return err
	}
	if power == types.VirtualMachinePowerStatePoweredOff {
		log.Infof("Roll back finished - Appliance is kept in powered off status")
		return nil
	}
	if err = d.startAppliance(conf); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(d.ctx, settings.Timeout)
	defer cancel()
	if err = d.CheckServiceReady(ctx, conf, nil); err != nil {
		// do not return error in this case, to make sure clean up continues
		log.Info("\tAPI may be slow to start - try to connect to API after a few minutes")
	}
	return nil
}

func (d *Dispatcher) reconfigVCH(conf *config.VirtualContainerHostConfigSpec, isoFile string) error {
	defer trace.End(trace.Begin(isoFile))

	spec := &types.VirtualMachineConfigSpec{}

	deviceChange, err := d.switchISO(isoFile)
	if err != nil {
		return err
	}

	spec.DeviceChange = deviceChange

	if conf != nil {
		// reset service started attribute
		for _, sess := range conf.ExecutorConfig.Sessions {
			sess.Started = ""
		}
		cfg := make(map[string]string)
		extraconfig.Encode(extraconfig.MapSink(cfg), conf)
		spec.ExtraConfig = append(spec.ExtraConfig, vmomi.OptionValueFromMap(cfg)...)
	}

	if spec.DeviceChange == nil && spec.ExtraConfig == nil {
		// nothing need to do
		log.Debugf("Nothing changed, no need to reconfigure appliance")
		return nil
	}

	// reconfig
	log.Infof("Setting VM configuration")
	info, err := d.appliance.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
		return d.appliance.Reconfigure(ctx, *spec)
	})

	if err != nil {
		log.Errorf("Error while reconfiguring appliance: %s", err)
		return err
	}
	if info.State != types.TaskInfoStateSuccess {
		log.Errorf("Reconfiguring appliance reported: %s", info.Error.LocalizedMessage)
		return err
	}
	return nil
}

func (d *Dispatcher) switchISO(filePath string) ([]types.BaseVirtualDeviceConfigSpec, error) {
	defer trace.End(trace.Begin(filePath))

	var devices object.VirtualDeviceList
	var err error

	log.Infof("Switching appliance iso to %s", filePath)
	devices, err = d.appliance.Device(d.ctx)
	if err != nil {
		log.Errorf("Failed to get vm devices for appliance: %s", err)
		return nil, err
	}
	// find the single cdrom
	cd, err := devices.FindCdrom("")
	if err != nil {
		log.Errorf("Failed to get CD rom device from appliance: %s", err)
		return nil, err
	}

	oldApplianceISO := cd.Backing.(*types.VirtualCdromIsoBackingInfo).FileName
	if oldApplianceISO == filePath {
		log.Debugf("Target file name %q is same to old one, no need to change.")
		return nil, nil
	}
	cd = devices.InsertIso(cd, filePath)
	changedDevices := object.VirtualDeviceList([]types.BaseVirtualDevice{cd})

	deviceChange, err := changedDevices.ConfigSpec(types.VirtualDeviceConfigSpecOperationEdit)
	if err != nil {
		log.Errorf("Failed to create config spec for appliance: %s", err)
		return nil, err
	}

	d.oldApplianceISO = oldApplianceISO
	return deviceChange, nil
}
