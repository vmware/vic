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

package vsphere

import (
	"fmt"
	"net/url"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/trace"

	"golang.org/x/net/context"
)

func VolumeJoin(ctx context.Context, handle *exec.Handle, volume *storage.Volume, diskOpts map[string]string) (*exec.Handle, error) {
	defer trace.End(trace.Begin("vsphere.VolumeJoin"))

	//TODO: populate the mode field of the MountSpec from the diskOpts(rw/ro)
	mountPath, err := volume.Device.MountPath()
	if err != nil {
		return nil, err
	}
	newMountSpec := metadata.MountSpec{
		Source: url.URL{
			Scheme: "label",
			Host:   volume.Label,
			Mode:   DiskOpts["Mode"],
		},
		Path: mountPath,
	}

	unitNumber := int32(-1)
	diskDevice := &types.VirtualDisk{
		CapacityInKB: 0,
		VirtualDevice: types.VirtualDevice{
			Key:           -1,
			ControllerKey: 100, //FIXME: This is hardcoded for now and should be located from the config spec in the future.y
			UnitNumber:    &unitNumber,
			Backing: &types.VirtualDiskFlatVer2BackingInfo{
				DiskMode: string(types.VirtualDiskModeIndependent_persistent),
				VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
					FileName: volume.Device.DiskPath(),
				},
			},
		},
	}

	config := &types.VirtualDeviceConfigSpec{
		Device:        diskDevice,
		Operation:     types.VirtualDeviceConfigSpecOperationAdd,
		FileOperation: "", //blank for existing disk
	}

	handle.Spec.DeviceChange = append(handle.Spec.DeviceChange, config)

	if _, ok := handle.ExecConfig.Mounts[volume.ID]; !ok {
		return nil, fmt.Errorf("Volume with ID %s is already in container %s's mountspec'", volume.ID, handle.Container.ID)
	}
	handle.ExecConfig.Mounts[volume.ID] = newMountSpec

	return handle, nil
}
