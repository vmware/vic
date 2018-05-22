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

package fs

import (
	"fmt"

	"github.com/docker/docker/pkg/mount"

	"github.com/vmware/vic/pkg/trace"
)

// Raw satisfies the Filesystem interface
type Raw struct{}

func NewRaw() *Raw {
	return &Raw{}
}

// Mkfs creates an ext4 fs on the given device and applices the given label
func (e *Raw) Mkfs(op trace.Operation, devPath, label string) error {
	return nil
}

// Mount mounts an ext4 formatted device at the given path.  From the Docker
// mount pkg, args must in the form arg=val.
func (e *Raw) Mount(op trace.Operation, devPath, targetPath string, options []string) error {
	defer trace.End(trace.Begin(devPath))
	return fmt.Errorf("mount not available for raw filesystem type")
}

// Unmount unmounts the disk.
// path can be a device path or a mount point
func (e *Raw) Unmount(op trace.Operation, path string) error {
	defer trace.End(trace.Begin(path))
	op.Infof("Unmounting %s", path)
	return mount.Unmount(path)
}

// SetLabel sets the label of an ext4 formated device
func (e *Raw) SetLabel(op trace.Operation, devPath, labelName string) error {
	return fmt.Errorf("setlabel not available for raw filesystem type")
}
