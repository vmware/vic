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
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
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
	if _, ok := info["guestinfo.vch/components"]; ok {
		return true, nil
	}
	return false, nil
}

func (d *Dispatcher) removeApplianceIfForced(conf *metadata.VirtualContainerHostConfigSpec) error {
	vm, err := d.findAppliance(conf)
	if err != nil {
		return err
	}
	log.Debugf("Appliance is found")
	if vm != nil && d.force {
		if ok, verr := d.isVCH(vm); !ok {
			verr = errors.Errorf("VM %s is found, but is not VCH appliance, please choose different name", conf.Name)
			return verr
		}
		log.Infof("Appliance exists, remove it...")
		_, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
			return vm.PowerOff(ctx)
		})
		if err != nil {
			log.Warnf("Failed to power off existing appliance for %s, try to remove anyway", err)
		}
		_, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
			return vm.Destroy(ctx)
		})
		if err != nil {
			err = errors.Errorf("Failed to destroy existing appliance: %s", err)
			return err
		}
		m := object.NewFileManager(d.session.Client.Client)
		if _, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
			return m.DeleteDatastoreFile(ctx, d.session.Datastore.Path(conf.Name), d.session.Datacenter)
		}); err != nil {
			err = errors.Errorf("Failed to remove existing VCH data files, %s", err)
			return err
		}
	} else if vm != nil {
		err = errors.Errorf("VM already exists with display name %s. Name must be unique. Exiting...", conf.Name)
		return err
	}

	return nil
}

func (d *Dispatcher) addNetworkDevices(conf *metadata.VirtualContainerHostConfigSpec, cspec *spec.VirtualMachineConfigSpec, devices object.VirtualDeviceList) (object.VirtualDeviceList, error) {
	var err error
	var backing types.BaseVirtualDeviceBackingInfo
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

		refparts := strings.SplitN(endpoint.Network.Common.ID, "-", 2)
		if len(refparts) != 2 {
			err = errors.Errorf("Network identifier for %s role was not in expected format: %s", name, endpoint.Network.Common.ID)
			return nil, err
		}

		log.Infof("Rehydrating network from %s and %s", refparts[0], refparts[1])
		network := object.NewReference(d.session.Client.Client, types.ManagedObjectReference{Type: refparts[0], Value: refparts[1]}).(object.NetworkReference)

		// FIXME: hack to work around rehydration issues - must not assume the net name like this
		backing, err = network.EthernetCardBackingInfo(d.ctx)
		if err != nil {
			err = errors.Errorf("Failed to get network backing info for %s: %s", network, err)
			return nil, err
		}
		backing2 := backing.(*types.VirtualEthernetCardNetworkBackingInfo)
		netname := strings.SplitN(refparts[1], "-", 2)
		backing2.DeviceName = netname[1]

		// END FIXME

		log.Infof("Using network %s for network role %s: %s", netname[1], name, err)

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
		log.Infof("Setting %s to slot %d", name, slot)

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

func (d *Dispatcher) createApplianceSpec(conf *metadata.VirtualContainerHostConfigSpec) (*types.VirtualMachineConfigSpec, error) {
	var devices object.VirtualDeviceList
	var err error

	cfg := make(map[string]string)
	extraconfig.Encode(extraconfig.MapSink(cfg), conf)

	spec := &spec.VirtualMachineConfigSpec{
		VirtualMachineConfigSpec: &types.VirtualMachineConfigSpec{
			Name:        conf.Name,
			GuestId:     "other3xLinux64Guest",
			Files:       &types.VirtualMachineFileInfo{VmPathName: fmt.Sprintf("[%s]", conf.ImageStoreName)},
			NumCPUs:     int32(conf.ApplianceSize.CPU.Limit),
			MemoryMB:    conf.ApplianceSize.Memory.Limit,
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

func (d *Dispatcher) getPresetExtraconfig(conf *metadata.VirtualContainerHostConfigSpec) []types.BaseOptionValue {
	extraConfig :=
		[]types.BaseOptionValue{
			&types.OptionValue{
				Key:   "guestinfo.vch/components",
				Value: "/sbin/docker-engine-server /sbin/port-layer-server /sbin/vicadmin",
			},
			&types.OptionValue{
				Key:   "guestinfo.vch/sbin/imagec",
				Value: "-debug -logfile=/var/log/vic/imagec.log -insecure-skip-verify",
			},
			&types.OptionValue{
				Key: "guestinfo.vch/sbin/port-layer-server",
				Value: fmt.Sprintf("--host=localhost --port=8080 --insecure --sdk=%s --datacenter=%s --cluster=%s --pool=%s --datastore=%s --network=%s --vch=%s",
					conf.Target.String(), conf.DatacenterName, conf.ClusterPath, d.vchPoolPath,
					conf.ImageStores[0], conf.Networks["client"].InventoryPath, conf.Name)},
		}

	files := "/var/tmp/images/ /var/log/vic/"

	if conf.CertPEM != "" && conf.KeyPEM != "" {
		d.VICAdminProto = "https"
		extraConfig = append(
			extraConfig,
			&types.OptionValue{
				Key:   "guestinfo.vch/etc/pki/tls/certs/vic-host-cert.pem",
				Value: conf.CertPEM,
			},
		)
		extraConfig = append(
			extraConfig,
			&types.OptionValue{
				Key:   "guestinfo.vch/etc/pki/tls/certs/vic-host-key.pem",
				Value: conf.KeyPEM,
			},
		)
		d.dockertlsargs = "-TLS -tls-certificate=/etc/pki/tls/certs/vic-host-cert.pem -tls-key=/etc/pki/tls/certs/vic-host-key.pem"
		vicadmintlsargs := " -hostcert=/etc/pki/tls/certs/vic-host-cert.pem -hostkey=/etc/pki/tls/certs/vic-host-key.pem"
		files = fmt.Sprintf("%s /etc/pki/tls/certs/vic-host-cert.pem /etc/pki/tls/certs/vic-host-key.pem", files)
		d.DockerPort = "2376"
		extraConfig = append(extraConfig,
			&types.OptionValue{
				Key:   "guestinfo.vch/sbin/docker-engine-server",
				Value: fmt.Sprintf("-serveraddr=0.0.0.0 -port=%s -port-layer-port=8080 %s", d.DockerPort, d.dockertlsargs),
			})
		extraConfig = append(extraConfig,
			&types.OptionValue{
				Key: "guestinfo.vch/sbin/vicadmin",
				Value: fmt.Sprintf("-docker-host=unix:///var/run/docker.sock -insecure -sdk=%s -ds=%s -vm-path=%s -cluster=%s -pool=%s %s",
					conf.Target.String(), conf.ImageStores[0].String(), conf.ApplianceInventoryPath, conf.ClusterPath, d.vchPoolPath, vicadmintlsargs),
			})
	} else {
		d.VICAdminProto = "http"
		d.DockerPort = "2375"
		extraConfig = append(extraConfig,
			&types.OptionValue{
				Key:   "guestinfo.vch/sbin/docker-engine-server",
				Value: fmt.Sprintf("-serveraddr=0.0.0.0 -port=%s -port-layer-port=8080", d.DockerPort),
			})
		extraConfig = append(extraConfig,
			&types.OptionValue{Key: "guestinfo.vch/sbin/vicadmin",
				Value: fmt.Sprintf("-docker-host=unix:///var/run/docker.sock -insecure -sdk=%s -ds=%s -vm-path=%s -cluster=%s -pool=%s -tls=%t",
					conf.Target.String(), conf.ImageStores[0].String(), conf.ApplianceInventoryPath, conf.ClusterPath, d.vchPoolPath, false),
			})
	}
	extraConfig = append(extraConfig,
		&types.OptionValue{
			Key:   "guestinfo.vch/files",
			Value: files,
		})
	// Set network info into guestinfo before VM is powered on, although the mac address is not availalbe yet.
	// This is to make sure the related attrs are persisted
	for nicName, netInfo := range conf.Networks {
		extraConfig = append(extraConfig,
			&types.OptionValue{
				Key:   fmt.Sprintf("guestinfo.vch/networks/%s/portgroup", nicName),
				Value: netInfo.PortGroupName},
		)
		extraConfig = append(extraConfig,
			&types.OptionValue{
				Key:   fmt.Sprintf("guestinfo.vch/networks/%s/mac", nicName),
				Value: " ",
			})
	}
	extraConfig = append(extraConfig,
		&types.OptionValue{
			Key:   "guestinfo.vch/networks",
			Value: " ",
		})
	return extraConfig
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
	cdrom = devices.InsertIso(cdrom, fmt.Sprintf("[%s] %s/appliance.iso", conf.ImageStoreName, d.vmPathName))
	devices = append(devices, cdrom)
	return devices, nil
}

func (d *Dispatcher) createAppliance(conf *metadata.VirtualContainerHostConfigSpec) error {
	log.Infof("Creating appliance on target")

	spec, err := d.createApplianceSpec(conf)
	if err != nil {
		log.Errorf("Unable to create appliance spec: %s", err)
		return err
	}

	// create test VM
	info, err := tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return d.session.Folders(ctx).VmFolder.CreateVM(ctx, *spec, d.vchPool.ResourcePool, d.session.Host)
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
	vm := vm.NewVirtualMachine(d.ctx, d.session, moref)

	// update the displayname to the actual folder name used
	if d.vmPathName, err = vm.FolderName(d.ctx); err != nil {
		log.Errorf("Failed to get canonical name for appliance: %s", err)
		return err
	}
	log.Debugf("vm folder name: %s", d.vmPathName)
	log.Debugf("vm inventory path: %s", vm.InventoryPath)

	conf.ApplianceInventoryPath = vm.InventoryPath

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
				"-cluster=" + conf.ClusterPath,
				"-pool=" + conf.ResourcePoolPath,
				// FIXME: tls is hardcoded false until vicadmin is migrated to extraconfig
				// this is to avoid having to put in code to push files into the appliance
				"-tls=false",
			},
		},
	},
	)

	conf.AddComponent("docker-personality", &metadata.SessionConfig{
		Cmd: metadata.Cmd{
			Path: "/sbin/docker-engine-server",
			Args: []string{
				"/sbin/docker-engine-server",
				//FIXME: hack during config migration
				"-serveraddr=0.0.0.0",
				"-port=2375",
				"-port-layer-port=8080",
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
				"--datacenter=" + conf.DatacenterName,
				"--cluster=" + conf.ClusterPath,
				"--pool=" + conf.ResourcePoolPath,
				"--datastore=" + conf.ImageStores[0].Host,
				"--network=" + conf.Networks["client"].InventoryPath,
				"--vch=" + conf.Common.Name,
			},
		},
	},
	)

	spec, err = d.reconfigureApplianceSpec(vm, conf)

	// reconfig
	info, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return vm.Reconfigure(ctx, *spec)
	})

	if err != nil {
		log.Errorf("Error while setting component parameters to appliance: %s", err)
		return err
	}
	if info.State != types.TaskInfoStateSuccess {
		log.Errorf("Setting parameters to appliance reported: %s", info.Error.LocalizedMessage)
		return err
	}

	d.appliance = vm
	return nil
}

func (d *Dispatcher) reconfigureApplianceSpec(vm *vm.VirtualMachine, conf *metadata.VirtualContainerHostConfigSpec) (*types.VirtualMachineConfigSpec, error) {
	var devices object.VirtualDeviceList
	var err error

	spec := &types.VirtualMachineConfigSpec{
		Name:    conf.Name,
		GuestId: "other3xLinux64Guest",
		Files:   &types.VirtualMachineFileInfo{VmPathName: fmt.Sprintf("[%s]", conf.ImageStoreName)},
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
	// set component execution parameters into guestinfo
	spec.ExtraConfig = d.getPresetExtraconfig(conf)

	cfg := make(map[string]string)
	extraconfig.Encode(extraconfig.MapSink(cfg), conf)
	spec.ExtraConfig = append(spec.ExtraConfig, extraconfig.OptionValueFromMap(cfg)...)
	return spec, nil
}

func (d *Dispatcher) setMacToGuestInfo(conf *metadata.VirtualContainerHostConfigSpec) error {
	m, err := d.appliance.WaitForMAC(d.ctx)
	if err != nil {
		err = errors.Errorf("Failed to get VM mac address %s", err)
		return err
	}
	var spec types.VirtualMachineConfigSpec
	spec = types.VirtualMachineConfigSpec{
		ExtraConfig: []types.BaseOptionValue{},
	}

	var keys []string
	for nicName, netInfo := range conf.Networks {
		mac, ok := m[netInfo.PortGroupName]
		if !ok || mac == "" {
			// timeout to wait MAC address, so empty mac address is returned
			err = errors.Errorf("Timeout to get VM MAC address")
			return err
		}

		spec.ExtraConfig = append(spec.ExtraConfig,
			&types.OptionValue{
				Key:   fmt.Sprintf("guestinfo.vch/networks/%s/mac", nicName),
				Value: mac,
			})
		netInfo.Mac = mac
		keys = append(keys, nicName)
	}

	// Do not persist VirtualContainerHost now cause only MAC address is changed in this object.
	// guestinfo update has bug through SDK, so all values updated after VM is powered on, will be removed from vmx file, that means those values
	// will lose after VM is restarted.
	// Need to revisit this, while the above MAC address guestinfo update is removed.

	//	cfg := make(map[string]string)
	//	extraconfig.EncodeWithPrefix(extraconfig.MapSink(cfg), conf, "guestinfo.vch")
	//	spec.ExtraConfig = append(spec.ExtraConfig, extraconfig.OptionValueFromMap(cfg)...)
	spec.ExtraConfig = append(spec.ExtraConfig,
		&types.OptionValue{
			Key:   "guestinfo.vch/networks",
			Value: strings.Join(keys, " "),
		})

	// reconfig
	_, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return d.appliance.Reconfigure(ctx, spec)
	})

	if err != nil {
		log.Errorf("Error to set MacAddress into guestinfo: %s", err)
		return err
	}

	return nil
}

func (d *Dispatcher) waitingForIP() error {
	var err error
	if d.HostIP, err = d.appliance.WaitForKeyInExtraConfig(d.ctx, "guestinfo.vch.clientip"); err != nil {
		return err
	}
	return nil
}

func (d *Dispatcher) makeSureApplianceRuns() error {
	var err error

	if d.appliance == nil {
		return nil
	}
	log.Infof("Waiting for IP information")

	if err = d.waitingForIP(); err != nil {
		err = fmt.Errorf("Timed out waiting for appliance to publish URI for docker API: %s", err.Error())
		log.Infof("Log files can be found on the appliance:")
		return err
	}
	return nil
}
