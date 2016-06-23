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
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	"golang.org/x/net/context"
)

var (
	lastSeenProgressMessage string
	unitNumber              int32
)

func (d *Dispatcher) isVCH(vm *vm.VirtualMachine) (bool, error) {
	if vm == nil {
		return false, errors.New("nil parameter")
	}

	info, err := vm.FetchExtraConfig(d.ctx)
	if err != nil {
		err = errors.Errorf("Failed to fetch guest info of appliance vm, %s", err)
		return false, err
	}

	var remoteConf metadata.VirtualContainerHostConfigSpec
	extraconfig.Decode(extraconfig.MapSource(info), &remoteConf)

	// if the moref of the target matches where we expect to find it for a VCH, run with it
	if remoteConf.ExecutorConfig.ID == vm.Reference().String() {
		return true, nil
	}

	return false, nil
}

func (d *Dispatcher) checkExistence(conf *metadata.VirtualContainerHostConfigSpec) error {
	vm, err := d.findAppliance(conf)
	if err != nil {
		return err
	}
	if vm == nil {
		return nil
	}

	log.Debugf("Appliance is found")
	if ok, verr := d.isVCH(vm); !ok {
		verr = errors.Errorf("VM %s is found, but is not VCH appliance, please choose different name", conf.Name)
		return verr
	}
	err = errors.Errorf("Appliance %s exists, to install with same name, please delete it first.", conf.Name)
	return err
}

func (d *Dispatcher) deleteVM(vm *vm.VirtualMachine, force bool) (string, error) {
	var err error
	power, err := vm.PowerState(d.ctx)
	if err != nil {
		err = errors.Errorf("Failed to get vm power status %s: %s", vm.Reference(), err)
		return "", err

	}
	if power != types.VirtualMachinePowerStatePoweredOff {
		if !force {
			name, err := vm.Name(d.ctx)
			if err != nil {
				log.Errorf("VM name is not found, %s", err)
			}
			if name != "" {
				err = errors.Errorf("VM %s is powered on", name)
			} else {
				err = errors.Errorf("VM %s is powered on", vm.Reference())
			}
			return "", err
		}
		if _, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
			return vm.PowerOff(ctx)
		}); err != nil {
			log.Debugf("Failed to power off existing appliance for %s, try to remove anyway", err)
		}
	}
	// get the actual folder name before we delete it
	folder, err := vm.FolderName(d.ctx)
	if err != nil {
		log.Warnf("Failed to get actual folder name for VM. Will not attempt to delete additional data files in VM directory: %s", err)
	}

	_, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return vm.Destroy(ctx)
	})
	if err != nil {
		err = errors.Errorf("Failed to destroy vm %s: %s", vm.Reference(), err)
		return "", err
	}
	return folder, nil
}

func (d *Dispatcher) addNetworkDevices(conf *metadata.VirtualContainerHostConfigSpec, cspec *spec.VirtualMachineConfigSpec, devices object.VirtualDeviceList) (object.VirtualDeviceList, error) {
	// network name:alias, to avoid create multiple devices for same network
	slots := make(map[int32]bool)
	nets := make(map[string]*metadata.NetworkEndpoint)

	for name, endpoint := range conf.ExecutorConfig.Networks {
		if pnic, ok := nets[endpoint.Network.Common.ID]; ok {
			// there's already a NIC on this network
			endpoint.Common.ID = pnic.Common.ID
			log.Infof("Network role %s is sharing NIC with %s", name, pnic.Network.Common.Name)
			continue
		}

		moref := new(types.ManagedObjectReference)
		if ok := moref.FromString(endpoint.Network.ID); !ok {
			return nil, fmt.Errorf("serialized managed object reference in unexpected format: %s", endpoint.Network.ID)
		}
		obj, err := d.session.Finder.ObjectReference(d.ctx, *moref)
		if err != nil {
			return nil, fmt.Errorf("unable to reacquire reference for network %s from serialized form: %s", endpoint.Network.Name, endpoint.Network.ID)
		}
		network, ok := obj.(object.NetworkReference)
		if !ok {
			return nil, fmt.Errorf("reacquired reference for network %s, from serialized form %s, was not a network: %T", endpoint.Network.Name, endpoint.Network.ID, obj)
		}

		backing, err := network.EthernetCardBackingInfo(d.ctx)
		if err != nil {
			err = errors.Errorf("Failed to get network backing info for %s: %s", network, err)
			return nil, err
		}

		nic, err := devices.CreateEthernetCard("vmxnet3", backing)
		if err != nil {
			err = errors.Errorf("Failed to create Ethernet Card spec for %s", err)
			return nil, err
		}

		slot := cspec.AssignSlotNumber(nic, slots)
		if slot == spec.NilSlot {
			err = errors.Errorf("Failed to assign stable PCI slot for %s network card", name)
		}

		endpoint.Common.ID = strconv.Itoa(int(slot))
		slots[slot] = true
		log.Debugf("Setting %s to slot %d", name, slot)

		devices = append(devices, nic)

		nets[endpoint.Network.Common.ID] = endpoint
	}
	return devices, nil
}

func (d *Dispatcher) addIDEController(devices object.VirtualDeviceList) (object.VirtualDeviceList, error) {
	// IDE controller
	scsi, err := devices.CreateIDEController()
	if err != nil {
		return nil, err
	}
	devices = append(devices, scsi)
	return devices, nil
}

func (d *Dispatcher) addParaVirtualSCSIController(devices object.VirtualDeviceList) (object.VirtualDeviceList, error) {
	// para virtual SCSI controller
	scsi, err := devices.CreateSCSIController("pvscsi")
	if err != nil {
		return nil, err
	}
	devices = append(devices, scsi)
	return devices, nil
}

func (d *Dispatcher) createApplianceSpec(conf *metadata.VirtualContainerHostConfigSpec, vConf *data.InstallerData) (*types.VirtualMachineConfigSpec, error) {
	var devices object.VirtualDeviceList
	var err error

	cfg := make(map[string]string)
	extraconfig.Encode(extraconfig.MapSink(cfg), conf)

	spec := &spec.VirtualMachineConfigSpec{
		VirtualMachineConfigSpec: &types.VirtualMachineConfigSpec{
			Name:     conf.Name,
			GuestId:  "other3xLinux64Guest",
			Files:    &types.VirtualMachineFileInfo{VmPathName: fmt.Sprintf("[%s]", conf.ImageStores[0].Host)},
			NumCPUs:  int32(vConf.ApplianceSize.CPU.Limit),
			MemoryMB: vConf.ApplianceSize.Memory.Limit,
			// Encode the config both here and after the VMs created so that it can be identified as a VCH appliance as soon as
			// creation is complete.
			ExtraConfig: extraconfig.OptionValueFromMap(cfg),
		},
	}

	if devices, err = d.addIDEController(devices); err != nil {
		return nil, err
	}

	if devices, err = d.addParaVirtualSCSIController(devices); err != nil {
		return nil, err
	}

	if devices, err = d.addNetworkDevices(conf, spec, devices); err != nil {
		return nil, err
	}

	deviceChange, err := devices.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	if err != nil {
		return nil, err
	}

	spec.DeviceChange = deviceChange
	return spec.VirtualMachineConfigSpec, nil
}

func (d *Dispatcher) findApplianceByID(conf *metadata.VirtualContainerHostConfigSpec) (*vm.VirtualMachine, error) {
	var err error
	var vmm *vm.VirtualMachine

	moref := new(types.ManagedObjectReference)
	if ok := moref.FromString(conf.ID); !ok {
		message := "Failed to get appliance VM mob reference"
		log.Errorf(message)
		return nil, errors.New(message)
	}
	ref, err := d.session.Finder.ObjectReference(d.ctx, *moref)
	if err != nil {
		if _, ok := err.(*find.NotFoundError); !ok {
			err = errors.Errorf("Failed to query appliance (%s): %s", moref, err)
			return nil, err
		}
		log.Debugf("Appliance is not found")
		return nil, nil

	}
	ovm, ok := ref.(*object.VirtualMachine)
	if !ok {
		log.Errorf("Failed to find VM %s, %s", moref, err)
		return nil, err
	}
	vmm = vm.NewVirtualMachine(d.ctx, d.session, ovm.Reference())
	return vmm, nil
}

func (d *Dispatcher) findAppliance(conf *metadata.VirtualContainerHostConfigSpec) (*vm.VirtualMachine, error) {
	ovm, err := d.session.Finder.VirtualMachine(d.ctx, conf.Name)
	if err != nil {
		_, ok := err.(*find.NotFoundError)
		if !ok {
			err = errors.Errorf("Failed to query appliance (%s): %s", conf.Name, err)
			return nil, err
		}
		log.Debugf("Appliance is not found")
		return nil, nil
	}
	newVM := vm.NewVirtualMachine(d.ctx, d.session, ovm.Reference())
	// workaround here. We lost the value set in ovm cause we wrap the object to another type
	newVM.InventoryPath = ovm.InventoryPath
	return newVM, nil
}

// retrieves the uuid of the appliance vm to create a unique vsphere extension name
func (d *Dispatcher) GenerateExtensionName(conf *metadata.VirtualContainerHostConfigSpec) error {
	// must be called after appliance VM is created
	vm, err := d.findAppliance(conf)

	if err != nil {
		return errors.Errorf("Could not find appliance at extension creation time; failed with error: %s", err)
	}

	var o mo.VirtualMachine
	err = vm.Properties(d.ctx, vm.Reference(), []string{"config.uuid"}, &o)
	if err != nil {
		return errors.Errorf("Could not get VM UUID from appliance VM due to error: %s", err)
	}

	conf.ExtensionName = "com.vmware.vic." + o.Config.Uuid
	return nil
}

func (d *Dispatcher) configIso(conf *metadata.VirtualContainerHostConfigSpec, vm *vm.VirtualMachine) (object.VirtualDeviceList, error) {
	var devices object.VirtualDeviceList
	var err error

	vmDevices, err := vm.Device(d.ctx)
	if err != nil {
		log.Errorf("Failed to get vm devices for appliance: %s", err)
		return nil, err
	}
	ide, err := vmDevices.FindIDEController("")
	if err != nil {
		log.Errorf("Failed to find IDE controller for appliance: %s", err)
		return nil, err
	}
	cdrom, err := devices.CreateCdrom(ide)
	if err != nil {
		log.Errorf("Failed to create Cdrom device for appliance: %s", err)
		return nil, err
	}
	cdrom = devices.InsertIso(cdrom, fmt.Sprintf("[%s] %s/appliance.iso", conf.ImageStores[0].Host, d.vmPathName))
	devices = append(devices, cdrom)
	return devices, nil
}

func (d *Dispatcher) createAppliance(conf *metadata.VirtualContainerHostConfigSpec, settings *data.InstallerData) error {
	log.Infof("Creating appliance on target")

	spec, err := d.createApplianceSpec(conf, settings)
	if err != nil {
		log.Errorf("Unable to create appliance spec: %s", err)
		return err
	}

	// create test VM
	info, err := tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return d.session.Folders(ctx).VmFolder.CreateVM(ctx, *spec, d.vchPool, d.session.Host)
	})

	if err != nil {
		log.Errorf("Unable to create appliance VM: %s", err)
		return err
	}
	if info.Error != nil || info.State != types.TaskInfoStateSuccess {
		log.Errorf("Create appliance reported: %s", info.Error.LocalizedMessage)
	}

	// get VM reference and save it
	moref := info.Result.(types.ManagedObjectReference)
	conf.SetMoref(&moref)
	obj, err := d.session.Finder.ObjectReference(d.ctx, moref)
	if err != nil {
		log.Errorf("Failed to reacquire reference to appliance VM after creation: %s", err)
		return err
	}
	gvm, ok := obj.(*object.VirtualMachine)
	if !ok {
		return fmt.Errorf("Required reference after appliance creation was not for a VM: %T", obj)
	}
	vm2 := vm.NewVirtualMachineFromVM(d.ctx, d.session, gvm)

	// update the displayname to the actual folder name used
	if d.vmPathName, err = vm2.FolderName(d.ctx); err != nil {
		log.Errorf("Failed to get canonical name for appliance: %s", err)
		return err
	}
	log.Debugf("vm folder name: %s", d.vmPathName)
	log.Debugf("vm inventory path: %s", vm2.InventoryPath)

	// create an extension to register the appliance as
	if err = d.GenerateExtensionName(conf); err != nil {
		return errors.Errorf("Could not generate extension name during appliance creation due to error: %s", err)
	}

	settings.Extension = types.Extension{
		Description: &types.Description{
			Label:   "VIC",
			Summary: "vSphere Integrated Containers Virtual Container Host",
		},
		Company: "VMware, Inc.",
		Version: "0.0",
		Key:     conf.ExtensionName,
	}

	conf.AddComponent("vicadmin", &metadata.SessionConfig{
		Cmd: metadata.Cmd{
			Path: "/sbin/vicadmin",
			Args: []string{
				"/sbin/vicadmin",
				"-docker-host=unix:///var/run/docker.sock",
				// FIXME: hack during config migration
				"-insecure",
				"-sdk=" + conf.Target.String(),
				"-ds=" + conf.ImageStores[0].Host,
				"-cluster=" + settings.ClusterPath,
				"-pool=" + settings.ResourcePoolPath,
				"-vm-path=" + vm2.InventoryPath,
			},
			Env: []string{
				"PATH=/sbin:/bin",
			},
		},
	},
	)

	if conf.HostCertificate != nil {
		d.VICAdminProto = "https"
		d.DockerPort = "2376"
	} else {
		d.VICAdminProto = "http"
		d.DockerPort = "2375"
	}

	conf.AddComponent("docker-personality", &metadata.SessionConfig{
		Cmd: metadata.Cmd{
			Path: "/sbin/docker-engine-server",
			Args: []string{
				"/sbin/docker-engine-server",
				//FIXME: hack during config migration
				"-serveraddr=0.0.0.0",
				"-port=" + d.DockerPort,
				"-port-layer-port=8080",
			},
			Env: []string{
				"PATH=/sbin",
			},
		},
	},
	)

	conf.AddComponent("port-layer", &metadata.SessionConfig{
		Cmd: metadata.Cmd{
			Path: "/sbin/port-layer-server",
			Args: []string{
				"/sbin/port-layer-server",
				//FIXME: hack during config migration
				"--host=localhost",
				"--port=8080",
				"--insecure",
				"--sdk=" + conf.Target.String(),
				"--datacenter=" + settings.DatacenterName,
				"--cluster=" + settings.ClusterPath,
				"--pool=" + settings.ResourcePoolPath,
				"--datastore=" + conf.ImageStores[0].Host,
				"--vch=" + conf.ExecutorConfig.Name,
			},
		},
	},
	)

	spec, err = d.reconfigureApplianceSpec(vm2, conf)

	// reconfig
	info, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return vm2.Reconfigure(ctx, *spec)
	})

	if err != nil {
		log.Errorf("Error while setting component parameters to appliance: %s", err)
		return err
	}
	if info.State != types.TaskInfoStateSuccess {
		log.Errorf("Setting parameters to appliance reported: %s", info.Error.LocalizedMessage)
		return err
	}

	d.appliance = vm2
	return nil
}

func (d *Dispatcher) reconfigureApplianceSpec(vm *vm.VirtualMachine, conf *metadata.VirtualContainerHostConfigSpec) (*types.VirtualMachineConfigSpec, error) {
	var devices object.VirtualDeviceList
	var err error

	spec := &types.VirtualMachineConfigSpec{
		Name:    conf.Name,
		GuestId: "other3xLinux64Guest",
		Files:   &types.VirtualMachineFileInfo{VmPathName: fmt.Sprintf("[%s]", conf.ImageStores[0].Host)},
	}

	if devices, err = d.configIso(conf, vm); err != nil {
		return nil, err
	}

	deviceChange, err := devices.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	if err != nil {
		log.Errorf("Failed to create config spec for appliance: %s", err)
		return nil, err
	}

	spec.DeviceChange = deviceChange

	cfg := make(map[string]string)
	extraconfig.Encode(extraconfig.MapSink(cfg), conf)
	spec.ExtraConfig = append(spec.ExtraConfig, extraconfig.OptionValueFromMap(cfg)...)
	return spec, nil
}

// applianceConfiguration updates the configuration passed in with the latest from the appliance VM.
// there's no guarantee of consistency within the configuration at this time
func (d *Dispatcher) applianceConfiguration(conf *metadata.VirtualContainerHostConfigSpec) error {
	extraConfig, err := d.appliance.FetchExtraConfig(d.ctx)
	if err != nil {
		return err
	}

	extraconfig.Decode(extraconfig.MapSource(extraConfig), conf)
	return nil
}

// waitForKey squashes the return values and simpy blocks until the key is updated or there is an error
func (d *Dispatcher) waitForKey(key string) {
	d.appliance.WaitForKeyInExtraConfig(d.ctx, key)
	return
}

func (d *Dispatcher) makeSureApplianceRuns(conf *metadata.VirtualContainerHostConfigSpec) error {
	if d.appliance == nil {
		return errors.New("cannot validate appliance due to missing VM reference")
	}

	log.Infof("Waiting for IP information")
	d.waitForKey("guestinfo..init.networks|client.ip")
	ctxerr := d.ctx.Err()

	if ctxerr == nil {
		log.Info("Waiting for major appliance components to launch")
		log.Debug("waiting for vicadmin to start")
		d.waitForKey("guestinfo..init.sessions|vicadmin.started")
		log.Debug("waiting for docker personality to start")
		d.waitForKey("guestinfo..init.sessions|docker-personality.started")
		log.Debug("waiting for port layer to start")
		d.waitForKey("guestinfo..init.sessions|port-layer.started")
	}

	// at this point either everything has succeeded or we're going into diagnostics, ignore error
	// as we're only using it for IP in the success case
	updateErr := d.applianceConfiguration(conf)

	// TODO: we should call to the general vic-machine inspect implementation here for more detail
	// but instead...
	if len(conf.ExecutorConfig.Networks["client"].Assigned) > 0 {
		d.HostIP = conf.ExecutorConfig.Networks["client"].Assigned.String()
		log.Debug("Obtained IP address for client interface: %s", d.HostIP)
		return nil
	}

	// it's possible we timed out... get updated info having adjusted context to allow it
	// keeping it short
	ctxerr = d.ctx.Err()

	d.ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err := d.applianceConfiguration(conf)
	if err != nil {
		return fmt.Errorf("unable to retrieve updated configuration from appliance for diagnostics: %s", err)
	}

	if ctxerr == context.DeadlineExceeded {
		log.Info("Failed to retrieve IP for client interface")
		log.Info("  State of all interfaces:")

		// if we timed out, then report status - if cancelled this doesn't need reporting
		for name, net := range conf.ExecutorConfig.Networks {
			addr := net.Assigned.String()
			if len(net.Assigned) == 0 {
				addr = "waiting for IP"
			}
			log.Infof("    %s IP: %s", name, addr)
		}

		// if we timed out, then report status - if cancelled this doesn't need reporting
		log.Info("  State of components:")
		for name, session := range conf.ExecutorConfig.Sessions {
			status := "waiting to launch"
			if session.Started == "true" {
				status = "started successfully"
			} else if session.Started != "" {
				status = session.Started
			}
			log.Infof("    %s: %s", name, status)
		}

		return errors.New("timed out waiting for IP address information from appliance")
	}

	return fmt.Errorf("could not obtain IP address information from appliance: %s", updateErr)
}
