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
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/metadata"
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

func (d *Dispatcher) getNetworkDevices(conf *metadata.VirtualContainerHostConfigSpec) ([]types.BaseVirtualDeviceConfigSpec, error) {
	var devices []types.BaseVirtualDeviceConfigSpec
	var err error
	var backing types.BaseVirtualDeviceBackingInfo

	for _, info := range conf.Networks {
		network := info.PortGroup
		backing, err = network.EthernetCardBackingInfo(d.ctx)
		if err != nil {
			err = errors.Errorf("Failed to get network backing info for %s: %s", network, err)
			return nil, err
		}

		net := &types.VirtualDeviceConfigSpec{
			Operation: types.VirtualDeviceConfigSpecOperationAdd,
			Device: &types.VirtualVmxnet3{
				VirtualVmxnet: types.VirtualVmxnet{
					VirtualEthernetCard: types.VirtualEthernetCard{
						VirtualDevice: types.VirtualDevice{
							Key:     -3,
							Backing: backing,
						},
						AddressType: string(types.VirtualEthernetCardMacTypeGenerated),
					},
				},
			},
		}
		devices = append(devices, net)
	}
	return devices, nil
}

func (d *Dispatcher) createApplianceSpec(conf *metadata.VirtualContainerHostConfigSpec) (*types.VirtualMachineConfigSpec, error) {
	netDevices, err := d.getNetworkDevices(conf)
	if err != nil {
		return nil, err
	}

	unitNumber = -1

	spec := &types.VirtualMachineConfigSpec{
		Name:     conf.Name,
		GuestId:  "other3xLinux64Guest",
		Files:    &types.VirtualMachineFileInfo{VmPathName: fmt.Sprintf("[%s]", conf.ImageStoreName)},
		NumCPUs:  int32(conf.ApplianceSize.CPU.Limit),
		MemoryMB: conf.ApplianceSize.Memory.Limit,
		ExtraConfig: []types.BaseOptionValue{
			&types.OptionValue{Key: "guestinfo.vch/components", Value: "/sbin/docker-engine-server /sbin/port-layer-server /sbin/vicadmin"},
			&types.OptionValue{Key: "guestinfo.vch/sbin/imagec", Value: "-debug -logfile=/var/log/vic/imagec.log -insecure"},
			&types.OptionValue{Key: "guestinfo.vch/sbin/port-layer-server",
				Value: fmt.Sprintf("--host=localhost --port=8080 --insecure --sdk=%s --datacenter=%s --cluster=%s --pool=%s --datastore=%s --network=%s --vch=%s",
					conf.Target, conf.DatacenterName, conf.ClusterPath, d.vchPoolPath,
					conf.ImageStores[0], conf.Networks["client"].InventoryPath, conf.Name)},
		},
		DeviceChange: []types.BaseVirtualDeviceConfigSpec{
			&types.VirtualDeviceConfigSpec{
				Operation: types.VirtualDeviceConfigSpecOperationAdd,
				Device: &types.ParaVirtualSCSIController{
					VirtualSCSIController: types.VirtualSCSIController{
						SharedBus: types.VirtualSCSISharingNoSharing,
						VirtualController: types.VirtualController{
							BusNumber: 0,
							VirtualDevice: types.VirtualDevice{
								Key: 100,
							},
						},
					},
				},
			},
			// ide controller for cdrom
			&types.VirtualDeviceConfigSpec{
				Operation: types.VirtualDeviceConfigSpecOperationAdd,
				Device: &types.VirtualIDEController{
					VirtualController: types.VirtualController{
						VirtualDevice: types.VirtualDevice{
							Key: 200,
						},
					},
				},
			},
		},
	}
	cfg := make(map[string]string)
	extraconfig.EncodeWithPrefix(extraconfig.MapSink(cfg), conf, "vch")
	spec.ExtraConfig = append(spec.ExtraConfig, extraconfig.OptionValueFromMap(cfg)...)

	spec.DeviceChange = append(spec.DeviceChange, netDevices...)
	return spec, nil
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

	// get VM
	vm, err := d.findAppliance(conf)
	if err != nil || vm == nil {
		err = errors.Errorf("Failed query back appliance: %s", err)
		return err
	}

	// update the displayname to the actual folder name used
	if d.vmPathName, err = vm.FolderName(d.ctx); err != nil {
		log.Errorf("Failed to get canonical name for appliance: %s", err)
		return err
	}
	log.Debugf("vm folder name: %s", d.vmPathName)
	log.Debugf("vm inventory path: %s", vm.InventoryPath)

	unitNumber = -1
	// set component execution parameters into guestinfo
	spec.DeviceChange = []types.BaseVirtualDeviceConfigSpec{
		// currently we're hardcoded to boot from a base ISO so add that
		&types.VirtualDeviceConfigSpec{
			Operation: types.VirtualDeviceConfigSpecOperationAdd,
			Device: &types.VirtualCdrom{
				VirtualDevice: types.VirtualDevice{
					Key:           -2,
					ControllerKey: 200,
					UnitNumber:    &unitNumber,
					Backing: &types.VirtualCdromIsoBackingInfo{
						VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
							FileName: fmt.Sprintf("[%s] %s/appliance.iso", conf.ImageStoreName, d.vmPathName),
						},
					},
				},
			},
		},
	}

	files := "/var/tmp/images/ /var/log/vic/"

	if conf.CertPEM != "" && conf.KeyPEM != "" {
		d.VICAdminProto = "https"
		spec.ExtraConfig = append(
			spec.ExtraConfig,
			&types.OptionValue{Key: "guestinfo.vch/etc/pki/tls/certs/vic-host-cert.pem", Value: conf.CertPEM},
		)
		spec.ExtraConfig = append(
			spec.ExtraConfig,
			&types.OptionValue{Key: "guestinfo.vch/etc/pki/tls/certs/vic-host-key.pem", Value: conf.KeyPEM},
		)
		d.dockertlsargs = "-TLS -tls-certificate=/etc/pki/tls/certs/vic-host-cert.pem -tls-key=/etc/pki/tls/certs/vic-host-key.pem"
		vicadmintlsargs := " -hostcert=/etc/pki/tls/certs/vic-host-cert.pem -hostkey=/etc/pki/tls/certs/vic-host-key.pem"
		files = fmt.Sprintf("%s /etc/pki/tls/certs/vic-host-cert.pem /etc/pki/tls/certs/vic-host-key.pem", files)
		d.DockerPort = "2376"
		spec.ExtraConfig = append(spec.ExtraConfig, &types.OptionValue{Key: "guestinfo.vch/sbin/docker-engine-server",
			Value: fmt.Sprintf("-serveraddr=0.0.0.0 -port=%s -port-layer-port=8080 %s", d.DockerPort, d.dockertlsargs)})
		spec.ExtraConfig = append(spec.ExtraConfig, &types.OptionValue{Key: "guestinfo.vch/sbin/vicadmin",
			Value: fmt.Sprintf("-docker-host=unix:///var/run/docker.sock -insecure -sdk=%s -ds=%s -vm-path=%s -cluster=%s -pool=%s %s",
				conf.Target, conf.ImageStores[0], vm.InventoryPath, conf.ClusterPath, d.vchPoolPath, vicadmintlsargs)})
	} else {
		d.VICAdminProto = "http"
		d.DockerPort = "2375"
		spec.ExtraConfig = append(spec.ExtraConfig, &types.OptionValue{Key: "guestinfo.vch/sbin/docker-engine-server",
			Value: fmt.Sprintf("-serveraddr=0.0.0.0 -port=%s -port-layer-port=8080", d.DockerPort)})
		spec.ExtraConfig = append(spec.ExtraConfig, &types.OptionValue{Key: "guestinfo.vch/sbin/vicadmin",
			Value: fmt.Sprintf("-docker-host=unix:///var/run/docker.sock -insecure -sdk=%s -ds=%s -vm-path=%s -cluster=%s -pool=%s -tls=%t",
				conf.Target, conf.ImageStores[0], vm.InventoryPath, conf.ClusterPath, d.vchPoolPath, false)})
	}
	spec.ExtraConfig = append(spec.ExtraConfig, &types.OptionValue{Key: "guestinfo.vch/files", Value: files})

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

		spec.ExtraConfig = append(spec.ExtraConfig, &types.OptionValue{Key: fmt.Sprintf("guestinfo.vch/networks/%s/portgroup", nicName), Value: netInfo.PortGroupName})
		spec.ExtraConfig = append(spec.ExtraConfig, &types.OptionValue{Key: fmt.Sprintf("guestinfo.vch/networks/%s/mac", nicName), Value: mac})
		keys = append(keys, nicName)
		netInfo.Mac = mac
	}
	spec.ExtraConfig = append(spec.ExtraConfig, &types.OptionValue{Key: "guestinfo.vch/networks", Value: strings.Join(keys, " ")})
	cfg := make(map[string]string)
	extraconfig.EncodeWithPrefix(extraconfig.MapSink(cfg), conf, "vch")
	spec.ExtraConfig = append(spec.ExtraConfig, extraconfig.OptionValueFromMap(cfg)...)

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

func (d *Dispatcher) waitingForIP(dul time.Duration) (map[string]string, error) {
	timeout := time.After(dul)
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		case <-timeout:
			return nil, errors.New("Timeout")
		case <-tick.C:
			info, err := d.appliance.FetchExtraConfig(d.ctx)
			if err != nil {
				return nil, err
			}
			if value, ok := info["guestinfo.vch.clientip"]; ok {
				d.HostIP = value
				return info, nil
			}
		}
	}
}

func (d *Dispatcher) makeSureApplianceRuns() error {
	var err error

	if d.appliance == nil {
		return nil
	}
	log.Infof("Waiting for IP information")

	_, err = d.waitingForIP(d.timeout)
	if err != nil {
		err = fmt.Errorf("Timed out waiting for appliance to publish URI for docker API")
		log.Infof("Log files can be found on the appliance:")
		return err
	}
	return nil
}
