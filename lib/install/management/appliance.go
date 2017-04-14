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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/opts"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/diag"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/extraconfig/vmomi"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

const (
	portLayerPort = constants.SerialOverLANPort

	// this is generated in the crypto/tls.alert code
	badTLSCertificate = "tls: bad certificate"
)

var (
	lastSeenProgressMessage string
	unitNumber              int32
)

func (d *Dispatcher) isVCH(vm *vm.VirtualMachine) (bool, error) {
	if vm == nil {
		return false, errors.New("nil parameter")
	}
	defer trace.End(trace.Begin(vm.InventoryPath))

	info, err := vm.FetchExtraConfig(d.ctx)
	if err != nil {
		err = errors.Errorf("Failed to fetch guest info of appliance vm: %s", err)
		return false, err
	}

	var remoteConf config.VirtualContainerHostConfigSpec
	extraconfig.Decode(extraconfig.MapSource(info), &remoteConf)

	// if the moref of the target matches where we expect to find it for a VCH, run with it
	if remoteConf.ExecutorConfig.ID == vm.Reference().String() || remoteConf.IsCreating() {
		return true, nil
	}

	return false, nil
}

func (d *Dispatcher) checkExistence(conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) error {
	defer trace.End(trace.Begin(""))

	var err error
	d.vchPoolPath = path.Join(settings.ResourcePoolPath, conf.Name)
	var orp *object.ResourcePool
	var vapp *object.VirtualApp
	if d.isVC {
		vapp, err = d.findVirtualApp(d.vchPoolPath)
		if err != nil {
			return err
		}
		if vapp != nil {
			orp = vapp.ResourcePool
		}
	}
	if orp == nil {
		if orp, err = d.findResourcePool(d.vchPoolPath); err != nil {
			return err
		}
	}
	if orp == nil {
		return nil
	}

	rp := compute.NewResourcePool(d.ctx, d.session, orp.Reference())
	vm, err := rp.GetChildVM(d.ctx, d.session, conf.Name)
	if err != nil {
		return err
	}
	if vm == nil {
		if vapp != nil {
			err = errors.Errorf("virtual app %q is found, but is not VCH, please choose different name", d.vchPoolPath)
			log.Error(err)
			return err
		}
		return nil
	}

	log.Debugf("Appliance is found")
	if ok, verr := d.isVCH(vm); !ok {
		verr = errors.Errorf("VM %q is found, but is not VCH appliance, please choose different name", conf.Name)
		return verr
	}
	err = errors.Errorf("Appliance %q exists, to install with same name, please delete it first.", conf.Name)
	return err
}

func (d *Dispatcher) getName(vm *vm.VirtualMachine) string {
	name, err := vm.Name(d.ctx)
	if err != nil {
		log.Errorf("VM name not found: %s", err)
		return ""
	}
	return name
}

func (d *Dispatcher) deleteVM(vm *vm.VirtualMachine, force bool) error {
	defer trace.End(trace.Begin(fmt.Sprintf("vm %q, force %t", vm.String(), force)))

	var err error
	power, err := vm.PowerState(d.ctx)
	if err != nil || power != types.VirtualMachinePowerStatePoweredOff {
		if err != nil {
			log.Warnf("Failed to get vm power status %q: %s", vm.Reference(), err)
		}
		if !force {
			if err != nil {
				return err
			}
			name := d.getName(vm)
			if name != "" {
				err = errors.Errorf("VM %q is powered on", name)
			} else {
				err = errors.Errorf("VM %q is powered on", vm.Reference())
			}
			return err
		}
		if _, err = vm.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
			return vm.PowerOff(ctx)
		}); err != nil {
			log.Debugf("Failed to power off existing appliance for %s, try to remove anyway", err)
		}
	}
	// get the actual folder name before we delete it
	folder, err := vm.FolderName(d.ctx)
	if err != nil {
		// failed to get folder name, might not be able to remove files for this VM
		name := d.getName(vm)
		if name == "" {
			log.Errorf("Unable to automatically remove all files in datastore for VM %q", vm.Reference())
		} else {
			// try to use the vm name in place of folder
			log.Infof("Delete will attempt to remove datastore files for VM %q", name)
			folder = name
		}
	}

	_, err = vm.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
		return vm.DeleteExceptDisks(ctx)
	})
	if err != nil {
		err = errors.Errorf("Failed to destroy VM %q: %s", vm.Reference(), err)
		err2 := vm.Unregister(d.ctx)
		if err2 != nil {
			return errors.Errorf("%s then failed to unregister VM: %s", err, err2)
		}
		log.Infof("Unregistered VM to cleanup after failed destroy: %q", vm.Reference())
	}
	if _, err = d.deleteDatastoreFiles(d.session.Datastore, folder, true); err != nil {
		log.Warnf("Failed to remove datastore files for VM path %q: %s", folder, err)
	}

	return nil
}

func (d *Dispatcher) addNetworkDevices(conf *config.VirtualContainerHostConfigSpec, cspec *spec.VirtualMachineConfigSpec, devices object.VirtualDeviceList) (object.VirtualDeviceList, error) {
	defer trace.End(trace.Begin(""))

	// network name:alias, to avoid create multiple devices for same network
	slots := make(map[int32]bool)
	nets := make(map[string]*executor.NetworkEndpoint)

	for name, endpoint := range conf.ExecutorConfig.Networks {
		if pnic, ok := nets[endpoint.Network.Common.ID]; ok {
			// there's already a NIC on this network
			endpoint.Common.ID = pnic.Common.ID
			log.Infof("Network role %q is sharing NIC with %q", name, pnic.Network.Common.Name)
			continue
		}

		moref := new(types.ManagedObjectReference)
		if ok := moref.FromString(endpoint.Network.ID); !ok {
			return nil, fmt.Errorf("serialized managed object reference in unexpected format: %q", endpoint.Network.ID)
		}
		obj, err := d.session.Finder.ObjectReference(d.ctx, *moref)
		if err != nil {
			return nil, fmt.Errorf("unable to reacquire reference for network %q from serialized form: %q", endpoint.Network.Name, endpoint.Network.ID)
		}
		network, ok := obj.(object.NetworkReference)
		if !ok {
			return nil, fmt.Errorf("reacquired reference for network %q, from serialized form %q, was not a network: %T", endpoint.Network.Name, endpoint.Network.ID, obj)
		}

		backing, err := network.EthernetCardBackingInfo(d.ctx)
		if err != nil {
			err = errors.Errorf("Failed to get network backing info for %q: %s", network, err)
			return nil, err
		}

		nic, err := devices.CreateEthernetCard("vmxnet3", backing)
		if err != nil {
			err = errors.Errorf("Failed to create Ethernet Card spec for %s", err)
			return nil, err
		}

		slot := cspec.AssignSlotNumber(nic, slots)
		if slot == spec.NilSlot {
			err = errors.Errorf("Failed to assign stable PCI slot for %q network card", name)
		}

		endpoint.Common.ID = strconv.Itoa(int(slot))
		slots[slot] = true
		log.Debugf("Setting %q to slot %d", name, slot)

		devices = append(devices, nic)

		nets[endpoint.Network.Common.ID] = endpoint
	}
	return devices, nil
}

func (d *Dispatcher) addIDEController(devices object.VirtualDeviceList) (object.VirtualDeviceList, error) {
	defer trace.End(trace.Begin(""))

	// IDE controller
	scsi, err := devices.CreateIDEController()
	if err != nil {
		return nil, err
	}
	devices = append(devices, scsi)
	return devices, nil
}

func (d *Dispatcher) addParaVirtualSCSIController(devices object.VirtualDeviceList) (object.VirtualDeviceList, error) {
	defer trace.End(trace.Begin(""))

	// para virtual SCSI controller
	scsi, err := devices.CreateSCSIController("pvscsi")
	if err != nil {
		return nil, err
	}
	devices = append(devices, scsi)
	return devices, nil
}

func (d *Dispatcher) createApplianceSpec(conf *config.VirtualContainerHostConfigSpec, vConf *data.InstallerData) (*types.VirtualMachineConfigSpec, error) {
	defer trace.End(trace.Begin(""))

	var devices object.VirtualDeviceList
	var err error

	// set to creating VCH
	conf.SetIsCreating(true)

	cfg, err := d.encodeConfig(conf)
	if err != nil {
		return nil, err
	}

	spec := &spec.VirtualMachineConfigSpec{
		VirtualMachineConfigSpec: &types.VirtualMachineConfigSpec{
			Name:               conf.Name,
			GuestId:            string(types.VirtualMachineGuestOsIdentifierOtherGuest64),
			AlternateGuestName: constants.DefaultAltVCHGuestName(),
			Files:              &types.VirtualMachineFileInfo{VmPathName: fmt.Sprintf("[%s]", conf.ImageStores[0].Host)},
			NumCPUs:            int32(vConf.ApplianceSize.CPU.Limit),
			MemoryMB:           vConf.ApplianceSize.Memory.Limit,
			// Encode the config both here and after the VMs created so that it can be identified as a VCH appliance as soon as
			// creation is complete.
			ExtraConfig: append(vmomi.OptionValueFromMap(cfg), &types.OptionValue{Key: "answer.msg.serial.file.open", Value: "Append"}),
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

func isManagedObjectNotFoundError(err error) bool {
	if soap.IsSoapFault(err) {
		_, ok := soap.ToSoapFault(err).VimFault().(types.ManagedObjectNotFound)
		return ok
	}

	return false
}

func (d *Dispatcher) findApplianceByID(conf *config.VirtualContainerHostConfigSpec) (*vm.VirtualMachine, error) {
	defer trace.End(trace.Begin(""))

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
		if !isManagedObjectNotFoundError(err) {
			err = errors.Errorf("Failed to query appliance (%q): %s", moref, err)
			return nil, err
		}
		log.Debugf("Appliance is not found")
		return nil, nil

	}
	ovm, ok := ref.(*object.VirtualMachine)
	if !ok {
		log.Errorf("Failed to find VM %q: %s", moref, err)
		return nil, err
	}
	vmm = vm.NewVirtualMachine(d.ctx, d.session, ovm.Reference())
	return vmm, nil
}

func (d *Dispatcher) configIso(conf *config.VirtualContainerHostConfigSpec, vm *vm.VirtualMachine, settings *data.InstallerData) (object.VirtualDeviceList, error) {
	defer trace.End(trace.Begin(""))

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
	cdrom = devices.InsertIso(cdrom, fmt.Sprintf("[%s] %s/%s", conf.ImageStores[0].Host, d.vmPathName, settings.ApplianceISO))
	devices = append(devices, cdrom)
	return devices, nil
}

func (d *Dispatcher) configLogging(conf *config.VirtualContainerHostConfigSpec, vm *vm.VirtualMachine, settings *data.InstallerData) (object.VirtualDeviceList, error) {
	defer trace.End(trace.Begin(""))

	devices, err := vm.Device(d.ctx)
	if err != nil {
		log.Errorf("Failed to get vm devices for appliance: %s", err)
		return nil, err
	}

	p, err := devices.CreateSerialPort()
	if err != nil {
		return nil, err
	}

	err = vm.AddDevice(d.ctx, p)
	if err != nil {
		return nil, err
	}

	devices, err = vm.Device(d.ctx)
	if err != nil {
		log.Errorf("Failed to get vm devices for appliance: %s", err)
		return nil, err
	}

	serial, err := devices.FindSerialPort("")
	if err != nil {
		log.Errorf("Failed to locate serial port for persistent log configuration: %s", err)
		return nil, err
	}

	// TODO: we need to add an accessor for generating paths within the VM directory
	vmx, err := vm.VMPathName(d.ctx)
	if err != nil {
		log.Errorf("Unable to determine path of appliance VM: %s", err)
		return nil, err
	}

	// TODO: move this construction into the spec package and update portlayer/logging to use it as well
	serial.Backing = &types.VirtualSerialPortFileBackingInfo{
		VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
			// name consistency with containerVM
			FileName: fmt.Sprintf("%s/tether.debug", path.Dir(vmx)),
		},
	}

	return []types.BaseVirtualDevice{serial}, nil
}

func (d *Dispatcher) setDockerPort(conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) {
	if conf.HostCertificate != nil {
		d.DockerPort = fmt.Sprintf("%d", opts.DefaultTLSHTTPPort)
	} else {
		d.DockerPort = fmt.Sprintf("%d", opts.DefaultHTTPPort)
	}
}

func (d *Dispatcher) createAppliance(conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) error {
	defer trace.End(trace.Begin(""))

	log.Infof("Creating appliance on target")

	spec, err := d.createApplianceSpec(conf, settings)
	if err != nil {
		log.Errorf("Unable to create appliance spec: %s", err)
		return err
	}

	var info *types.TaskInfo
	// create appliance VM
	if d.isVC && d.vchVapp != nil {
		info, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
			return d.vchVapp.CreateChildVM(ctx, *spec, d.session.Host)
		})
	} else {
		// if vapp is not created, fall back to create VM under default resource pool
		folder := d.session.Folders(d.ctx).VmFolder
		info, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
			return folder.CreateVM(ctx, *spec, d.vchPool, d.session.Host)
		})
	}

	if err != nil {
		log.Errorf("Unable to create appliance VM: %s", err)
		return err
	}
	if info.Error != nil || info.State != types.TaskInfoStateSuccess {
		log.Errorf("Create appliance reported: %s", info.Error.LocalizedMessage)
	}

	if err = d.createVolumeStores(conf); err != nil {
		return errors.Errorf("Exiting because we could not create volume stores due to error: %s", err)
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
	log.Debugf("vm folder name: %q", d.vmPathName)
	log.Debugf("vm inventory path: %q", vm2.InventoryPath)

	vicadmin := executor.Cmd{
		Path: "/sbin/vicadmin",
		Args: []string{
			"/sbin/vicadmin",
			"--dc=" + settings.DatacenterName,
			"--pool=" + settings.ResourcePoolPath,
			"--cluster=" + settings.ClusterPath,
		},
		Env: []string{
			"PATH=/sbin:/bin",
			"GOTRACEBACK=all",
		},
		Dir: "/home/vicadmin",
	}
	if settings.HTTPProxy != nil {
		vicadmin.Env = append(vicadmin.Env, fmt.Sprintf("HTTP_PROXY=%s", settings.HTTPProxy.String()))
	}
	if settings.HTTPSProxy != nil {
		vicadmin.Env = append(vicadmin.Env, fmt.Sprintf("HTTPS_PROXY=%s", settings.HTTPSProxy.String()))
	}

	conf.AddComponent("vicadmin", &executor.SessionConfig{
		User:    "vicadmin",
		Group:   "vicadmin",
		Cmd:     vicadmin,
		Restart: true,
		Active:  true,
	},
	)

	d.setDockerPort(conf, settings)

	personality := executor.Cmd{
		Path: "/sbin/docker-engine-server",
		Args: []string{
			"/sbin/docker-engine-server",
			//FIXME: hack during config migration
			"-port=" + d.DockerPort,
			fmt.Sprintf("-port-layer-port=%d", portLayerPort),
		},
		Env: []string{
			"PATH=/sbin",
			"GOTRACEBACK=all",
		},
	}
	if settings.HTTPProxy != nil {
		personality.Env = append(personality.Env, fmt.Sprintf("HTTP_PROXY=%s", settings.HTTPProxy.String()))
	}
	if settings.HTTPSProxy != nil {
		personality.Env = append(personality.Env, fmt.Sprintf("HTTPS_PROXY=%s", settings.HTTPSProxy.String()))
	}

	conf.AddComponent("docker-personality", &executor.SessionConfig{
		// currently needed for iptables interaction
		// User:  "nobody",
		// Group: "nobody",
		Cmd:     personality,
		Restart: true,
		Active:  true,
	},
	)

	conf.AddComponent("port-layer", &executor.SessionConfig{
		Cmd: executor.Cmd{
			Path: "/sbin/port-layer-server",
			Args: []string{
				"/sbin/port-layer-server",
				"--host=localhost",
				fmt.Sprintf("--port=%d", portLayerPort),
			},
			Env: []string{
				//FIXME: hack during config migration
				"VC_URL=" + conf.Target,
				"DC_PATH=" + settings.DatacenterName,
				"CS_PATH=" + settings.ClusterPath,
				"POOL_PATH=" + settings.ResourcePoolPath,
				"DS_PATH=" + conf.ImageStores[0].Host,
			},
		},
		Restart: true,
		Active:  true,
	},
	)

	conf.BootstrapImagePath = fmt.Sprintf("[%s] %s/%s", conf.ImageStores[0].Host, d.vmPathName, settings.BootstrapISO)

	spec, err = d.reconfigureApplianceSpec(vm2, conf, settings)
	if err != nil {
		log.Errorf("Error while getting appliance reconfig spec: %s", err)
		return err
	}

	// reconfig
	info, err = vm2.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
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

func (d *Dispatcher) encodeConfig(conf *config.VirtualContainerHostConfigSpec) (map[string]string, error) {
	if d.secret == nil {
		log.Debug("generating new config secret key")

		s, err := extraconfig.NewSecretKey()
		if err != nil {
			return nil, err
		}

		d.secret = s
	}

	cfg := make(map[string]string)
	extraconfig.Encode(d.secret.Sink(extraconfig.MapSink(cfg)), conf)

	return cfg, nil
}

func (d *Dispatcher) reconfigureApplianceSpec(vm *vm.VirtualMachine, conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) (*types.VirtualMachineConfigSpec, error) {
	defer trace.End(trace.Begin(""))

	var devices object.VirtualDeviceList
	var err error

	spec := &types.VirtualMachineConfigSpec{
		Name:               conf.Name,
		GuestId:            string(types.VirtualMachineGuestOsIdentifierOtherGuest64),
		AlternateGuestName: constants.DefaultAltVCHGuestName(),
		Files:              &types.VirtualMachineFileInfo{VmPathName: fmt.Sprintf("[%s]", conf.ImageStores[0].Host)},
	}

	// create new devices
	if devices, err = d.configIso(conf, vm, settings); err != nil {
		return nil, err
	}

	newDevices, err := devices.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	if err != nil {
		log.Errorf("Failed to create config spec for appliance: %s", err)
		return nil, err
	}

	spec.DeviceChange = newDevices

	// update existing devices
	if devices, err = d.configLogging(conf, vm, settings); err != nil {
		return nil, err
	}

	updateDevices, err := devices.ConfigSpec(types.VirtualDeviceConfigSpecOperationEdit)
	if err != nil {
		log.Errorf("Failed to create config spec for logging update: %s", err)
		return nil, err
	}

	spec.DeviceChange = append(spec.DeviceChange, updateDevices...)

	cfg, err := d.encodeConfig(conf)
	if err != nil {
		return nil, err
	}

	spec.ExtraConfig = append(spec.ExtraConfig, vmomi.OptionValueFromMap(cfg)...)
	return spec, nil
}

// applianceConfiguration updates the configuration passed in with the latest from the appliance VM.
// there's no guarantee of consistency within the configuration at this time
func (d *Dispatcher) applianceConfiguration(conf *config.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin(""))

	extraConfig, err := d.appliance.FetchExtraConfig(d.ctx)
	if err != nil {
		return err
	}

	extraconfig.Decode(extraconfig.MapSource(extraConfig), conf)
	return nil
}

// waitForKey squashes the return values and simpy blocks until the key is updated or there is an error
func (d *Dispatcher) waitForKey(key string) {
	defer trace.End(trace.Begin(key))

	d.appliance.WaitForKeyInExtraConfig(d.ctx, key)
	return
}

// isPortLayerRunning decodes the `docker info` response to check if the portlayer is running
func isPortLayerRunning(res *http.Response) bool {
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Debugf("error while reading res body: %s", err.Error())
		return false
	}

	var sysInfo dockertypes.Info
	if err = json.Unmarshal(resBody, &sysInfo); err != nil {
		log.Debugf("error while unmarshalling res body: %s", err.Error())
		return false
	}

	for _, status := range sysInfo.SystemStatus {
		if status[0] == sysInfo.Driver {
			return status[1] == "RUNNING"
		}
	}
	return false
}

// CheckDockerAPI checks if the appliance components are initialized by issuing
// `docker info` to the appliance
func (d *Dispatcher) CheckDockerAPI(conf *config.VirtualContainerHostConfigSpec, clientCert *tls.Certificate) error {
	defer trace.End(trace.Begin(""))

	var (
		proto          string
		client         *http.Client
		res            *http.Response
		err            error
		req            *http.Request
		tlsErrExpected bool
	)

	if conf.HostCertificate.IsNil() {
		// TLS disabled
		proto = "http"
		client = &http.Client{}
	} else {
		// TLS enabled
		proto = "https"

		// #nosec: TLS InsecureSkipVerify set true
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		// appliance is configured for tlsverify, but we don't have a client certificate
		if len(conf.CertificateAuthorities) > 0 {
			// if tlsverify was configured at all then we must verify the remote
			tr.TLSClientConfig.InsecureSkipVerify = false

			func() {
				log.Debug("Loading CAs for client auth")
				pool, err := x509.SystemCertPool()
				if err != nil {
					log.Warnf("Unable to load system root certificates - continuing with only the provided CA")
					pool = x509.NewCertPool()
				}

				if !pool.AppendCertsFromPEM(conf.CertificateAuthorities) {
					log.Warn("Unable add CAs from config to validation pool")
				}

				// tr.TLSClientConfig.ClientCAs = pool
				tr.TLSClientConfig.RootCAs = pool

				if clientCert == nil {
					// we know this will fail, but we can try to distinguish the expected error vs
					// unresponsive endpoint
					tlsErrExpected = true
					log.Debugf("CA configured on appliance but no client certificate available")
					return
				}

				cert, err := conf.HostCertificate.X509Certificate()
				if err != nil {
					log.Debugf("Unable to extract host certificate: %s", err)
					tlsErrExpected = true
					return
				}

				cip := net.ParseIP(d.HostIP)
				if err != nil {
					log.Debugf("Unable to process Docker API host address from %q: %s", d.HostIP, err)
					tlsErrExpected = true
					return
				}

				// find the name to use and override the IP if found
				addr, err := addrToUse([]net.IP{cip}, cert, conf.CertificateAuthorities)
				if err != nil {
					log.Debugf("Unable to determine address to use with remote certificate, checking SANs")
					addr, _ = viableHostAddress([]net.IP{cip}, cert, conf.CertificateAuthorities)
					log.Debugf("Using host address: %s", addr)
				}
				if addr != "" {
					d.HostIP = addr
				} else {
					log.Debug("Failed to find a viable address for Docker API from certificates")
					// Server certificate won't validate since we don't have a hostname
					tlsErrExpected = true
				}
				log.Debugf("Host address set to: %q", d.HostIP)
			}()
		}

		if clientCert != nil {
			log.Debug("Assigning certificates for client auth")
			tr.TLSClientConfig.Certificates = []tls.Certificate{*clientCert}
		}

		client = &http.Client{Transport: tr}
	}

	dockerInfoURL := fmt.Sprintf("%s://%s:%s/info", proto, d.HostIP, d.DockerPort)
	log.Debugf("Docker API endpoint: %s", dockerInfoURL)
	req, err = http.NewRequest("GET", dockerInfoURL, nil)
	if err != nil {
		return errors.New("invalid HTTP request for docker info")
	}
	req = req.WithContext(d.ctx)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		res, err = client.Do(req)
		if err == nil {
			if res.StatusCode == http.StatusOK {
				if isPortLayerRunning(res) {
					log.Debug("Confirmed port layer is operational")
					break
				}
			}

			log.Debugf("Received HTTP status %d: %s", res.StatusCode, res.Status)
		} else {
			// DEBU[2016-10-11T22:22:38Z] Error received from endpoint: Get https://192.168.78.127:2376/info: dial tcp 192.168.78.127:2376: getsockopt: connection refused &{%!t(string=Get) %!t(string=https://192.168.78.127:2376/info) %!t(*net.OpError=&{dial tcp <nil> 0xc4204505a0 0xc4203a5e00})}
			// DEBU[2016-10-11T22:22:39Z] Components not yet initialized, retrying
			// ERR=&url.Error{
			//     Op:  "Get",
			//     URL: "https://192.168.78.127:2376/info",
			//     Err: &net.OpError{
			//         Op:     "dial",
			//         Net:    "tcp",
			//         Source: nil,
			//         Addr:   &net.TCPAddr{
			//             IP:   {0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xc0, 0xa8, 0x4e, 0x7f},
			//             Port: 2376,
			//             Zone: "",
			//         },
			//         Err: &os.SyscallError{
			//             Syscall: "getsockopt",
			//             Err:     syscall.Errno(0x6f),
			//         },
			//     },
			// }
			// DEBU[2016-10-11T22:22:41Z] Error received from endpoint: Get https://192.168.78.127:2376/info: remote error: tls: bad certificate &{%!t(string=Get) %!t(string=https://192.168.78.127:2376/info) %!t(*net.OpError=&{remote error  <nil> <nil> 42})}
			// DEBU[2016-10-11T22:22:42Z] Components not yet initialized, retrying
			// ERR=&url.Error{
			//     Op:  "Get",
			//     URL: "https://192.168.78.127:2376/info",
			//     Err: &net.OpError{
			//         Op:     "remote error",
			//         Net:    "",
			//         Source: nil,
			//         Addr:   nil,
			//         Err:    tls.alert(0x2a),
			//     },
			// }

			// ECONNREFUSED: 111, 0x6f

			uerr, ok := err.(*url.Error)
			if ok {
				switch neterr := uerr.Err.(type) {
				case *net.OpError:
					switch root := neterr.Err.(type) {
					case *os.SyscallError:
						if root.Err == syscall.Errno(syscall.ECONNREFUSED) {
							// waiting for API server to start
							log.Debug("connection refused")
						} else {
							log.Debugf("Error was expected to be ECONNREFUSED: %#v", root.Err)
						}
					default:
						errmsg := root.Error()

						if tlsErrExpected {
							log.Warnf("Expected TLS error without access to client certificate, received error: %s", errmsg)
							return nil
						}

						// the TLS package doesn't expose the raw reason codes
						// but we're actually looking for alertBadCertificate (42)
						if errmsg == badTLSCertificate {
							// TODO: programmatic check for clock skew on host
							log.Errorf("Connection failed with TLS error \"bad certificate\" - check for clock skew on the host")
						} else {
							log.Errorf("Connection failed with error: %s", root)
						}

						return fmt.Errorf("failed to connect to %s: %s", dockerInfoURL, root)
					}

				case x509.UnknownAuthorityError:
					// This will occur if the server certificate was signed by a CA that is not the one used for client authentication
					// and does not have a trusted root registered on the system running vic-machine
					msg := fmt.Sprintf("Unable to validate server certificate with configured CAs (unknown CA): %s", neterr.Error())
					if tlsErrExpected {
						// Legitimate deployment so no error, but definitely requires a warning.
						log.Warn(msg)
						return nil
					}
					// TLS error not expected, the validation failure is a problem
					log.Error(msg)
					return neterr

				case x509.HostnameError:
					// e.g. "doesn't contain any IP SANs"
					msg := fmt.Sprintf("Server certificate hostname doesn't match: %s", neterr.Error())
					if tlsErrExpected {
						log.Warn(msg)
						return nil
					}
					log.Error(msg)
					return neterr

				default:
					log.Debugf("Unhandled net error type: %#v", neterr)
					return neterr
				}
			} else {
				log.Debugf("Error type was expected to be url.Error: %#v", err)
			}
		}

		select {
		case <-ticker.C:
		case <-d.ctx.Done():
			return d.ctx.Err()
		}

		log.Debug("Components not yet initialized, retrying")
	}

	return nil
}

// ensureApplianceInitializes checks if the appliance component processes are launched correctly
func (d *Dispatcher) ensureApplianceInitializes(conf *config.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin(""))

	if d.appliance == nil {
		return errors.New("cannot validate appliance due to missing VM reference")
	}

	log.Infof("Waiting for IP information")
	d.waitForKey(extraconfig.CalculateKeys(conf, "ExecutorConfig.Networks.client.Assigned.IP", "")[0])
	ctxerr := d.ctx.Err()

	if ctxerr == nil {
		log.Info("Waiting for major appliance components to launch")
		for _, k := range extraconfig.CalculateKeys(conf, "ExecutorConfig.Sessions.*.Started", "") {
			d.waitForKey(k)
		}
	}

	// at this point either everything has succeeded or we're going into diagnostics, ignore error
	// as we're only using it for IP in the success case
	updateErr := d.applianceConfiguration(conf)

	// TODO: we should call to the general vic-machine inspect implementation here for more detail
	// but instead...
	if !ip.IsUnspecifiedIP(conf.ExecutorConfig.Networks["client"].Assigned.IP) {
		d.HostIP = conf.ExecutorConfig.Networks["client"].Assigned.IP.String()
		log.Infof("Obtained IP address for client interface: %q", d.HostIP)
		return nil
	}

	// it's possible we timed out... get updated info having adjusted context to allow it
	// keeping it short
	ctxerr = d.ctx.Err()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	d.ctx = ctx
	err := d.applianceConfiguration(conf)
	if err != nil {
		return fmt.Errorf("unable to retrieve updated configuration from appliance for diagnostics: %s", err)
	}

	if ctxerr == context.DeadlineExceeded {
		log.Info("")
		log.Error("Failed to obtain IP address for client interface")
		log.Info("Use vic-machine inspect to see if VCH has received an IP address at a later time")
		log.Info("  State of all interfaces:")

		// if we timed out, then report status - if cancelled this doesn't need reporting
		for name, net := range conf.ExecutorConfig.Networks {
			addr := net.Assigned.String()
			if ip.IsUnspecifiedIP(net.Assigned.IP) {
				addr = "waiting for IP"
			}
			log.Infof("    %q IP: %q", name, addr)
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
			log.Infof("    %q: %q", name, status)
		}

		return errors.New("Failed to obtain IP address for client interface (timed out)")
	}

	return fmt.Errorf("Failed to get IP address information from appliance: %s", updateErr)
}

// CheckServiceReady checks if service is launched correctly, including ip address, service initialization, VC connection and Docker API
// Should expand this method for any more VCH service checking
func (d *Dispatcher) CheckServiceReady(ctx context.Context, conf *config.VirtualContainerHostConfigSpec, clientCert *tls.Certificate) error {
	oldCtx := d.ctx
	d.ctx = ctx
	defer func() {
		d.ctx = oldCtx
	}()

	if err := d.ensureApplianceInitializes(conf); err != nil {
		return err
	}

	// vic-init will try to reach out to the vSphere target.
	log.Info("Checking VCH connectivity with vSphere target")
	// Checking access to vSphere API
	if cd, err := d.CheckAccessToVCAPI(d.ctx, d.appliance, conf.Target); err == nil {
		code := int(cd)
		if code > 0 {
			log.Warningf("vSphere API Test: %s %s", conf.Target, diag.UserReadableVCAPITestDescription(code))
		} else {
			log.Infof("vSphere API Test: %s %s", conf.Target, diag.UserReadableVCAPITestDescription(code))
		}
	} else {
		log.Warningf("Could not run VCH vSphere API target check due to %v but the VCH may still function normally", err)
	}

	if err := d.CheckDockerAPI(conf, clientCert); err != nil {
		err = errors.Errorf("Docker API endpoint check failed: %s", err)
		// log with info cause this might not be an error
		log.Info(err.Error())
		return err
	}
	return nil
}
