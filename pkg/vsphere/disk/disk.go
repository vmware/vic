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

package disk

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/fs"
)

// FilesystemType represents the filesystem in use by a virtual disk
type FilesystemType uint8

const (
	// Ext4 represents the ext4 file system
	Ext4 FilesystemType = iota + 1

	// Xfs represents the XFS file system
	Xfs

	// Ntfs represents the NTFS file system
	Ntfs
)

// Filesystem defines the interface for handling an attached virtual disk
type Filesystem interface {
	Mkfs(devPath, label string) error
	SetLabel(devPath, labelName string) error
	Mount(devPath, targetPath string, options []string) error
	Unmount(path string) error
}

// FilesystemTypeToFilesystem returns a new Filesystem implementation
// associated with the supplied FilesystemType
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

	mountedRefs int

	attachedRefs int

	fs Filesystem
}

// NewVirtualDisk creates and returns a new VirtualDisk object associated with the
// given datastore formatted with the specified FilesystemType
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

func (d *VirtualDisk) setAttached(devicePath string) (err error) {
	defer func() {
		if err == nil {
			// bump the attached reference count
			d.attachedRefs++
			log.Debugf("Increased attach references for %s to %d", d.DatastoreURI, d.attachedRefs)
		}
	}()

	if d.Attached() {
		log.Warnf("%s is already attached (%s)", d.DatastoreURI, devicePath)
		return nil
	}

	if devicePath == "" {
		err = fmt.Errorf("no device path specified")
		return err
	}

	// set the device path where attached
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

	if d.InUseByOther() {
		return fmt.Errorf("%s is still in use", d.DatastoreURI)
	}

	return nil
}

func (d *VirtualDisk) setDetached() error {
	if !d.Attached() {
		return fmt.Errorf("%s is already detached", d.DatastoreURI)
	}

	if d.Mounted() {
		return fmt.Errorf("%s is still mounted (%s)", d.DatastoreURI, d.mountPath)
	}

	if !d.AttachedByOther() {
		d.DevicePath = ""
	} else {
		log.Warnf("%s is still in use", d.DatastoreURI)
	}
	d.attachedRefs--
	log.Debugf("Decreased attach references for %s to %d", d.DatastoreURI, d.attachedRefs)

	return nil
}

// Mkfs formats the disk with Filesystem and sets the disk label
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

// SetLabel sets this disk's label
func (d *VirtualDisk) SetLabel(labelName string) error {
	d.lock()
	defer d.unlock()

	if !d.Attached() {
		return fmt.Errorf("%s isn't attached", d.DatastoreURI)
	}

	return d.fs.SetLabel(d.DevicePath, labelName)
}

// Attached returns true if this disk is attached, false otherwise
func (d *VirtualDisk) Attached() bool {
	return d.DevicePath != ""
}

// AttachedByOther returns true if the attached references are > 1
func (d *VirtualDisk) AttachedByOther() bool {
	return d.attachedRefs > 1
}

// MountedByOther returns true if the mounted references are > 1
func (d *VirtualDisk) MountedByOther() bool {
	return d.mountedRefs > 1
}

// InUseByOther returns true if the disk is currently attached or
// mounted by someone else
func (d *VirtualDisk) InUseByOther() bool {
	return d.MountedByOther() || d.AttachedByOther()
}

// Mount attempts to mount this disk. A NOP occurs if the disk is already mounted
func (d *VirtualDisk) Mount(mountPath string, options []string) (err error) {
	d.lock()
	defer d.unlock()

	defer func() {
		// bump mounted reference count
		d.mountedRefs++
		log.Debugf("Increased mount references for %s to %d", d.DatastoreURI, d.mountedRefs)
	}()

	if d.Mounted() {
		p, _ := d.MountPath()
		log.Warnf("%s already mounted at %s", d.DatastoreURI, p)
		return nil
	}

	if !d.Attached() {
		err = fmt.Errorf("%s isn't attached", d.DatastoreURI)
		return err
	}

	if err = d.fs.Mount(d.DevicePath, mountPath, options); err != nil {
		return err
	}

	d.mountPath = mountPath
	return nil
}

// Unmount attempts to unmount a virtual disk
func (d *VirtualDisk) Unmount() error {
	d.lock()
	defer d.unlock()

	if !d.Mounted() {
		return fmt.Errorf("%s already unmounted", d.DatastoreURI)
	}

	d.mountedRefs--
	log.Debugf("Decreased mount references for %s to %d", d.DatastoreURI, d.mountedRefs)

	// no more mount references to this disk, so actually unmount
	if d.mountedRefs == 0 {
		if err := d.fs.Unmount(d.mountPath); err != nil {
			return err
		}
		d.mountPath = ""
	}

	return nil
}

// MountPath returns the path on which the virtual disk is mounted,
// or an error if the disk is not mounted
func (d *VirtualDisk) MountPath() (string, error) {
	if !d.Mounted() {
		return "", fmt.Errorf("%s isn't mounted", d.DatastoreURI)
	}

	return d.mountPath, nil
}

// DiskPath returns a URL referencing the path of the virtual disk
// on the datastore
func (d *VirtualDisk) DiskPath() url.URL {

	return url.URL{
		Scheme: "ds",
		Path:   d.DatastoreURI.String(),
	}
}

// Mounted returns true if the virtual disk is mounted, false otherwise
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

// VerifyDatastoreDiskURI ensures the disk name ends in ".vmdk"
func VerifyDatastoreDiskURI(name string) error {
	if !strings.HasSuffix(name, ".vmdk") {
		return fmt.Errorf("%s isn't a vmdk", name)
	}
	return nil
}
