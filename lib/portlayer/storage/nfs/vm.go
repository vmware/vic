// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package nfs

import (
	"fmt"
	"net/url"

	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/trace"
)

const (
	nfsMountOptions = "rw,noatime,vers=3,rsize=131072,wsize=131072,hard,proto=tcp,timeo=600,sec=sys,mountvers=3,mountproto=tcp,local_lock=none"
)

func VolumeJoin(op trace.Operation, handle *exec.Handle, volume *storage.Volume, mountPath string, diskOpts map[string]string) (*exec.Handle, error) {
	defer trace.End(trace.Begin("nfs.VolumeJoin"))

	if _, ok := handle.ExecConfig.Mounts[volume.ID]; ok {
		return nil, fmt.Errorf("Volume with ID %s is already in container %s's mountspec config", volume.ID, handle.ExecConfig.ID)
	}

	// construct MountSpec for the tether
	mountSpec := createMountSpec(volume.Device.DiskPath(), mountPath, diskOpts)

	if handle.ExecConfig.Mounts == nil {
		handle.ExecConfig.Mounts = make(map[string]executor.MountSpec)
	}
	handle.ExecConfig.Mounts[volume.ID] = mountSpec

	return handle, nil
}

func createMountSpec(host url.URL, mountPath string, diskOpts map[string]string) executor.MountSpec {
	deviceMode := nfsMountOptions + ",addr=" + host.Host
	newMountSpec := executor.MountSpec{
		Source: host,
		Path:   mountPath,
		Mode:   deviceMode,
	}
	return newMountSpec
}
