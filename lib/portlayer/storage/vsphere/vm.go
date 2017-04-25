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

package vsphere

import (
	"fmt"
	"net/url"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/trace"
)

func VolumeJoin(op trace.Operation, handle *exec.Handle, volume *storage.Volume, mountPath string, diskOpts map[string]string) (*exec.Handle, error) {
	defer trace.End(trace.Begin("vsphere.VolumeJoin"))

	if _, ok := handle.ExecConfig.Mounts[volume.ID]; ok {
		return nil, fmt.Errorf("Volume with ID %s is already in container %s's mountspec config", volume.ID, handle.ExecConfig.ID)
	}

	//constuct MountSpec for the tether
	mountSpec := createMountSpec(volume, mountPath, diskOpts)
	//append a device addition spec change to the container config
	diskDevice := createVolumeVirtualDisk(volume)
	config := createDeviceConfigSpec(diskDevice)
	handle.Spec.DeviceChange = append(handle.Spec.DeviceChange, config)

	if handle.ExecConfig.Mounts == nil {
		handle.ExecConfig.Mounts = make(map[string]executor.MountSpec)
	}
	handle.ExecConfig.Mounts[volume.ID] = mountSpec

	return handle, nil
}

func createVolumeVirtualDisk(volume *storage.Volume) *types.VirtualDisk {
	unitNumber := int32(-1)
	diskDevice := &types.VirtualDisk{
		CapacityInKB: 0,
		VirtualDevice: types.VirtualDevice{
			Key:           -1,
			ControllerKey: 100, //FIXME: This is hardcoded for now and should be located from the config spec in the future.
			UnitNumber:    &unitNumber,
			Backing: &types.VirtualDiskFlatVer2BackingInfo{
				DiskMode: string(types.VirtualDiskModeIndependent_persistent),
				VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
					FileName: volume.Device.DiskPath().Path,
				},
			},
		},
	}
	return diskDevice
}

func createDeviceConfigSpec(diskDevice *types.VirtualDisk) *types.VirtualDeviceConfigSpec {
	config := &types.VirtualDeviceConfigSpec{
		Device:        diskDevice,
		Operation:     types.VirtualDeviceConfigSpecOperationAdd,
		FileOperation: "", //blank for existing disk
	}
	return config
}

func createMountSpec(volume *storage.Volume, mountPath string, diskOpts map[string]string) executor.MountSpec {
	deviceMode := diskOpts["Mode"]
	newMountSpec := executor.MountSpec{
		Source: url.URL{
			Scheme: "label",
			Path:   volume.Label,
		},
		Path:     mountPath,
		Mode:     deviceMode,
		CopyMode: volume.CopyMode,
	}
	return newMountSpec
}
