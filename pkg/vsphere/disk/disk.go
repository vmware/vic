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

package disk

import (
	"fmt"
	"strings"
	"sync"
)

// VirtualDisk represents a VMDK in the datastore, the device node it may be
// attached at (if it's attached), the mountpoint it is mounted at (if
// mounted), and other configuration.
type VirtualDisk struct {
	// The URI in the datastore this disk can be found with
	DatastoreURI string

	// The device node the disk is attached to
	DevicePath string

	// The path on the filesystem this device is attached to.
	MountPath string

	// To avoid attach/detach races, this lock serializes operations to the disk.
	l sync.Mutex
}

func NewVirtualDisk(DatastoreURI string) (*VirtualDisk, error) {
	if err := VerifyDatastoreDiskURI(DatastoreURI); err != nil {
		return nil, err
	}

	d := &VirtualDisk{
		DatastoreURI: DatastoreURI,
	}

	return d, nil
}

func (d *VirtualDisk) lock() {
	d.l.Lock()
}

func (d *VirtualDisk) unlock() {
	d.l.Unlock()
}

func (d *VirtualDisk) setAttached(devicePath string) error {
	if d.isAttached() {
		return fmt.Errorf("%s is already attached (%s)", d.DatastoreURI, devicePath)
	}

	if devicePath == "" {
		return fmt.Errorf("no device path specified")
	}

	d.DevicePath = devicePath
	return nil
}

func (d *VirtualDisk) canBeDetached() error {
	if !d.isAttached() {
		return fmt.Errorf("%s is already detached", d.DatastoreURI)
	}

	if d.isMounted() {
		return fmt.Errorf("%s is mounted (%s)", d.DatastoreURI, d.MountPath)
	}

	return nil
}

func (d *VirtualDisk) setDetached() error {
	if !d.isAttached() {
		return fmt.Errorf("%s is already dettached", d.DatastoreURI)
	}

	if d.isMounted() {
		return fmt.Errorf("%s is still mounted (%s)", d.DatastoreURI, d.MountPath)
	}

	d.DevicePath = ""
	return nil
}

func (d *VirtualDisk) isAttached() bool {
	return d.DevicePath != ""
}

func (d *VirtualDisk) setMounted(mountPath string) error {
	if !d.isAttached() {
		return fmt.Errorf("%s isn't attached", d.DatastoreURI)
	}

	if d.isMounted() {
		return fmt.Errorf("%s already mounted", d.DatastoreURI)
	}

	if mountPath == "" {
		return fmt.Errorf("mountPath missing")
	}

	d.MountPath = mountPath
	return nil
}

func (d *VirtualDisk) isMounted() bool {
	return d.MountPath != ""
}

func (d *VirtualDisk) canBeUnmounted() error {
	if !d.isAttached() {
		return fmt.Errorf("%s is detached", d.DatastoreURI)
	}

	if !d.isMounted() {
		return fmt.Errorf("%s is unmounted", d.DatastoreURI)
	}

	return nil
}

func (d *VirtualDisk) setUmounted() error {
	if !d.isMounted() {
		return fmt.Errorf("%s already unmounted", d.DatastoreURI)
	}

	d.MountPath = ""
	return nil
}

func VerifyDatastoreDiskURI(name string) error {
	if !strings.HasSuffix(name, ".vmdk") {
		return fmt.Errorf("%s isn't a vmdk", name)
	}
	return nil
}
