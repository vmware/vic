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
	"net/url"
	"strings"
	"sync"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/fs"
)

type FilesystemType uint8

const (
	Ext4 FilesystemType = iota + 1
	Xfs
	Ntfs
)

type Filesystem interface {
	Mkfs(devPath, label string) error
	SetLabel(devPath, labelName string) error
	Mount(devPath, targetPath string, options []string) error
	Unmount(path string) error
}

func FilesystemTypeToFilesystem(fstype FilesystemType) Filesystem {
	switch fstype {
	case Xfs:
		return fs.NewXFS()
	default:
		return fs.NewExt4()
	}
}

// VirtualDisk represents a VMDK in the datastore, the device node it may be
// attached at (if it's attached), the mountpoint it is mounted at (if
// mounted), and other configuration.
type VirtualDisk struct {
	// The URI in the datastore this disk can be found with
	DatastoreURI *object.DatastorePath

	// The URI in the datastore to the parent of this disk
	ParentDatastoreURI *object.DatastorePath

	// The device node the disk is attached to
	DevicePath string

	// The path on the filesystem this device is attached to.
	mountPath string

	// To avoid attach/detach races, this lock serializes operations to the disk.
	l sync.Mutex

	fs Filesystem
}

func NewVirtualDisk(DatastoreURI *object.DatastorePath, fst FilesystemType) (*VirtualDisk, error) {
	if err := VerifyDatastoreDiskURI(DatastoreURI.String()); err != nil {
		return nil, err
	}

	d := &VirtualDisk{
		DatastoreURI: DatastoreURI,
		// We only support ext4 for now
		fs: FilesystemTypeToFilesystem(fst),
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
	if d.Attached() {
		return fmt.Errorf("%s is already attached (%s)", d.DatastoreURI, devicePath)
	}

	if devicePath == "" {
		return fmt.Errorf("no device path specified")
	}

	d.DevicePath = devicePath
	return nil
}

func (d *VirtualDisk) canBeDetached() error {
	if !d.Attached() {
		return fmt.Errorf("%s is already detached", d.DatastoreURI)
	}

	if d.Mounted() {
		return fmt.Errorf("%s is mounted (%s)", d.DatastoreURI, d.mountPath)
	}

	return nil
}

func (d *VirtualDisk) setDetached() error {
	if !d.Attached() {
		return fmt.Errorf("%s is already dettached", d.DatastoreURI)
	}

	if d.Mounted() {
		return fmt.Errorf("%s is still mounted (%s)", d.DatastoreURI, d.mountPath)
	}

	d.DevicePath = ""
	return nil
}

func (d *VirtualDisk) Mkfs(labelName string) error {
	d.lock()
	defer d.unlock()

	if !d.Attached() {
		return fmt.Errorf("%s isn't attached", d.DatastoreURI)
	}

	if d.Mounted() {
		return fmt.Errorf("%s is mounted mounted", d.DatastoreURI)
	}

	return d.fs.Mkfs(d.DevicePath, labelName)
}

func (d *VirtualDisk) SetLabel(labelName string) error {
	d.lock()
	defer d.unlock()

	if !d.Attached() {
		return fmt.Errorf("%s isn't attached", d.DatastoreURI)
	}

	return d.fs.SetLabel(d.DevicePath, labelName)
}

func (d *VirtualDisk) Attached() bool {
	return d.DevicePath != ""
}

func (d *VirtualDisk) Mount(mountPath string, options []string) error {
	d.lock()
	defer d.unlock()

	if !d.Attached() {
		return fmt.Errorf("%s isn't attached", d.DatastoreURI)
	}

	if d.Mounted() {
		return fmt.Errorf("%s already mounted", d.DatastoreURI)
	}

	if err := d.fs.Mount(d.DevicePath, mountPath, options); err != nil {
		return err
	}

	d.mountPath = mountPath
	return nil
}

func (d *VirtualDisk) Unmount() error {
	d.lock()
	defer d.unlock()

	if !d.Mounted() {
		return fmt.Errorf("%s already unmounted", d.DatastoreURI)
	}

	if err := d.fs.Unmount(d.mountPath); err != nil {
		return err
	}

	d.mountPath = ""
	return nil
}

func (d *VirtualDisk) MountPath() (string, error) {
	if !d.Mounted() {
		return "", fmt.Errorf("%s isn't mounted", d.DatastoreURI)
	}

	return d.mountPath, nil
}

func (d *VirtualDisk) DiskPath() url.URL {

	return url.URL{
		Scheme: "ds",
		Path:   d.DatastoreURI.String(),
	}
}

func (d *VirtualDisk) Mounted() bool {
	return d.mountPath != ""
}

func (d *VirtualDisk) canBeUnmounted() error {
	if !d.Attached() {
		return fmt.Errorf("%s is detached", d.DatastoreURI)
	}

	if !d.Mounted() {
		return fmt.Errorf("%s is unmounted", d.DatastoreURI)
	}

	return nil
}

func (d *VirtualDisk) setUmounted() error {
	if !d.Mounted() {
		return fmt.Errorf("%s already unmounted", d.DatastoreURI)
	}

	d.mountPath = ""
	return nil
}

func VerifyDatastoreDiskURI(name string) error {
	if !strings.HasSuffix(name, ".vmdk") {
		return fmt.Errorf("%s isn't a vmdk", name)
	}
	return nil
}
