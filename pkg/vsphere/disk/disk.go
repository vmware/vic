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

// VirtualDisk represents a VMDK in the datastore, the device node it may be
// attached at (if it's attached), the mountpoint it is mounted at (if
// mounted), and other configuration.
type VirtualDisk struct {
	*VirtualDiskConfig

	// The device node the disk is attached to
	DevicePath string

	// The path on the filesystem this device is attached to.
	mountPath string

	// To avoid attach/detach races, this lock serializes operations to the disk.
	l sync.Mutex

	mountedRefs int

	attachedRefs int
}

// NewVirtualDisk creates and returns a new VirtualDisk object associated with the
// given datastore formatted with the specified FilesystemType
func NewVirtualDisk(config *VirtualDiskConfig, disks map[uint64]*VirtualDisk) (*VirtualDisk, error) {
	if !strings.HasSuffix(config.DatastoreURI.String(), ".vmdk") {
		return nil, fmt.Errorf("%s isn't a vmdk", config.DatastoreURI.String())
	}

	if d, ok := disks[config.Hash()]; ok {
		log.Debugf("Found the disk %s in the DiskManager cache, using it", config.DatastoreURI)
		return d, nil
	}
	log.Debugf("Didn't find the disk %s in the DiskManager cache, creating it", config.DatastoreURI)

	d := &VirtualDisk{
		VirtualDiskConfig: config,
	}
	disks[config.Hash()] = d

	return d, nil
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

func (d *VirtualDisk) setDetached(disks map[uint64]*VirtualDisk) error {
	defer func() {
		if d.attachedRefs == 0 {
			log.Debugf("Dropping %s from the DiskManager cache", d.DatastoreURI)

			delete(disks, d.Hash())
		}
	}()

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
	d.l.Lock()
	defer d.l.Unlock()

	if !d.Attached() {
		return fmt.Errorf("%s isn't attached", d.DatastoreURI)
	}

	if d.Mounted() {
		return fmt.Errorf("%s is still mounted (%s)", d.DatastoreURI, d.mountPath)
	}

	return d.Filesystem.Mkfs(d.DevicePath, labelName)
}

// SetLabel sets this disk's label
func (d *VirtualDisk) SetLabel(labelName string) error {
	d.l.Lock()
	defer d.l.Unlock()

	if !d.Attached() {
		return fmt.Errorf("%s isn't attached", d.DatastoreURI)
	}

	return d.Filesystem.SetLabel(d.DevicePath, labelName)
}

// Attached returns true if this disk is attached, false otherwise
func (d *VirtualDisk) Attached() bool {
	d.l.Lock()
	defer d.l.Unlock()

	return d.DevicePath != ""
}

// AttachedByOther returns true if the attached references are > 1
func (d *VirtualDisk) AttachedByOther() bool {
	d.l.Lock()
	defer d.l.Unlock()

	return d.attachedRefs > 1
}

// MountedByOther returns true if the mounted references are > 1
func (d *VirtualDisk) MountedByOther() bool {
	d.l.Lock()
	defer d.l.Unlock()

	return d.mountedRefs > 1
}

// InUseByOther returns true if the disk is currently attached or
// mounted by someone else
func (d *VirtualDisk) InUseByOther() bool {
	d.l.Lock()
	defer d.l.Unlock()

	return d.MountedByOther() || d.AttachedByOther()
}

// Mount attempts to mount this disk. A NOP occurs if the disk is already mounted
func (d *VirtualDisk) Mount(mountPath string, options []string) (err error) {
	d.l.Lock()
	defer d.l.Unlock()

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

	if err = d.Filesystem.Mount(d.DevicePath, mountPath, options); err != nil {
		return err
	}

	d.mountPath = mountPath
	return nil
}

// Unmount attempts to unmount a virtual disk
func (d *VirtualDisk) Unmount() error {
	d.l.Lock()
	defer d.l.Unlock()

	if !d.Mounted() {
		return fmt.Errorf("%s already unmounted", d.DatastoreURI)
	}

	d.mountedRefs--
	log.Debugf("Decreased mount references for %s to %d", d.DatastoreURI, d.mountedRefs)

	// no more mount references to this disk, so actually unmount
	if d.mountedRefs == 0 {
		if err := d.Filesystem.Unmount(d.mountPath); err != nil {
			return err
		}
		d.mountPath = ""
	}

	return nil
}

// MountPath returns the path on which the virtual disk is mounted,
// or an error if the disk is not mounted
func (d *VirtualDisk) MountPath() (string, error) {
	d.l.Lock()
	defer d.l.Unlock()

	if !d.Mounted() {
		return "", fmt.Errorf("%s isn't mounted", d.DatastoreURI)
	}

	return d.mountPath, nil
}

// DiskPath returns a URL referencing the path of the virtual disk
// on the datastore
func (d *VirtualDisk) DiskPath() url.URL {
	d.l.Lock()
	defer d.l.Unlock()

	return url.URL{
		Scheme: "ds",
		Path:   d.DatastoreURI.String(),
	}
}

// Mounted returns true if the virtual disk is mounted, false otherwise
func (d *VirtualDisk) Mounted() bool {
	d.l.Lock()
	defer d.l.Unlock()

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
