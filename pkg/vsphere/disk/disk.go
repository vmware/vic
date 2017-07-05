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
	"sync/atomic"

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

// Semaphore represents the number of references to a disk
type Semaphore struct {
	resource string
	refname  string
	count    uint64
}

// NewSemaphore creates and returns a Semaphore initialized to 0
func NewSemaphore(r, n string) *Semaphore {
	return &Semaphore{
		resource: r,
		refname:  n,
		count:    0,
	}
}

// Increment increases the reference count by one
func (r *Semaphore) Increment() uint64 {
	return atomic.AddUint64(&r.count, 1)
}

// Decrement decreases the reference count by one
func (r *Semaphore) Decrement() uint64 {
	return atomic.AddUint64(&r.count, ^uint64(0))
}

// Count returns the current reference count
func (r *Semaphore) Count() uint64 {
	return atomic.LoadUint64(&r.count)
}

// InUseError is returned when a detach is attempted on a disk that is
// still in use
type InUseError struct {
	error
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

	mountedRefs *Semaphore

	attachedRefs *Semaphore
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

	uri := config.DatastoreURI.String()
	d := &VirtualDisk{
		VirtualDiskConfig: config,
		mountedRefs:       NewSemaphore(uri, "mount"),
		attachedRefs:      NewSemaphore(uri, "attach"),
	}
	disks[config.Hash()] = d

	return d, nil
}

func (d *VirtualDisk) setAttached(devicePath string) (err error) {
	defer func() {
		if err == nil {
			// bump the attached reference count
			d.attachedRefs.Increment()
		}
	}()

	if d.attached() {
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
	if !d.attached() {
		return fmt.Errorf("%s is already detached", d.DatastoreURI)
	}

	if d.mounted() {
		return fmt.Errorf("%s is mounted (%s)", d.DatastoreURI, d.mountPath)
	}

	if d.inUseByOther() {
		return fmt.Errorf("Detach skipped - %s is still in use", d.DatastoreURI)
	}

	return nil
}

func (d *VirtualDisk) setDetached(disks map[uint64]*VirtualDisk) error {
	defer func() {
		if d.attachedRefs.Count() == 0 {
			log.Debugf("Dropping %s from the DiskManager cache", d.DatastoreURI)

			delete(disks, d.Hash())
		}
	}()

	if !d.attached() {
		return fmt.Errorf("%s is already detached", d.DatastoreURI)
	}

	if d.mounted() {
		return fmt.Errorf("%s is still mounted (%s)", d.DatastoreURI, d.mountPath)
	}

	if !d.attachedByOther() {
		d.DevicePath = ""
	} else {
		log.Warnf("%s is still in use", d.DatastoreURI)
	}

	return nil
}

// Mkfs formats the disk with Filesystem and sets the disk label
func (d *VirtualDisk) Mkfs(labelName string) error {
	d.l.Lock()
	defer d.l.Unlock()

	if !d.attached() {
		return fmt.Errorf("%s isn't attached", d.DatastoreURI)
	}

	if d.mounted() {
		return fmt.Errorf("%s is still mounted (%s)", d.DatastoreURI, d.mountPath)
	}

	return d.Filesystem.Mkfs(d.DevicePath, labelName)
}

// SetLabel sets this disk's label
func (d *VirtualDisk) SetLabel(labelName string) error {
	d.l.Lock()
	defer d.l.Unlock()

	if !d.attached() {
		return fmt.Errorf("%s isn't attached", d.DatastoreURI)
	}

	return d.Filesystem.SetLabel(d.DevicePath, labelName)
}

func (d *VirtualDisk) attached() bool {
	return d.DevicePath != ""
}

// Attached returns true if this disk is attached, false otherwise
func (d *VirtualDisk) Attached() bool {
	d.l.Lock()
	defer d.l.Unlock()

	return d.attached()
}

func (d *VirtualDisk) attachedByOther() bool {
	return d.attachedRefs.Count() > 1
}

// AttachedByOther returns true if the attached references are > 1
func (d *VirtualDisk) AttachedByOther() bool {
	d.l.Lock()
	defer d.l.Unlock()

	return d.attachedByOther()
}

func (d *VirtualDisk) mountedByOther() bool {
	return d.mountedRefs.Count() > 1
}

// MountedByOther returns true if the mounted references are > 1
func (d *VirtualDisk) MountedByOther() bool {
	d.l.Lock()
	defer d.l.Unlock()

	return d.mountedByOther()
}

func (d *VirtualDisk) inUseByOther() bool {
	return d.mountedByOther() || d.attachedByOther()
}

// InUseByOther returns true if the disk is currently attached or
// mounted by someone else
func (d *VirtualDisk) InUseByOther() bool {
	d.l.Lock()
	defer d.l.Unlock()

	return d.inUseByOther()
}

// Mount attempts to mount this disk. A NOP occurs if the disk is already mounted
func (d *VirtualDisk) Mount(mountPath string, options []string) (err error) {
	d.l.Lock()
	defer d.l.Unlock()

	defer func() {
		// bump mounted reference count
		d.mountedRefs.Increment()
	}()

	if d.mounted() {
		p, _ := d.mountPathFn()
		log.Warnf("%s already mounted at %s", d.DatastoreURI, p)
		return nil
	}

	if !d.attached() {
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

	if !d.mounted() {
		return fmt.Errorf("%s already unmounted", d.DatastoreURI)
	}

	d.mountedRefs.Decrement()

	// no more mount references to this disk, so actually unmount
	if d.mountedRefs.Count() == 0 {
		if err := d.Filesystem.Unmount(d.mountPath); err != nil {
			return err
		}
		d.mountPath = ""
	}

	return nil
}

func (d *VirtualDisk) mountPathFn() (string, error) {
	if !d.mounted() {
		return "", fmt.Errorf("%s isn't mounted", d.DatastoreURI)
	}

	return d.mountPath, nil
}

// MountPath returns the path on which the virtual disk is mounted,
// or an error if the disk is not mounted
func (d *VirtualDisk) MountPath() (string, error) {
	d.l.Lock()
	defer d.l.Unlock()

	return d.mountPathFn()
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

func (d *VirtualDisk) mounted() bool {
	return d.mountPath != ""
}

// Mounted returns true if the virtual disk is mounted, false otherwise
func (d *VirtualDisk) Mounted() bool {
	d.l.Lock()
	defer d.l.Unlock()

	return d.mounted()
}

func (d *VirtualDisk) canBeUnmounted() error {
	if !d.attached() {
		return fmt.Errorf("%s is detached", d.DatastoreURI)
	}

	if !d.mounted() {
		return fmt.Errorf("%s is unmounted", d.DatastoreURI)
	}

	return nil
}

func (d *VirtualDisk) setUmounted() error {
	if !d.mounted() {
		return fmt.Errorf("%s already unmounted", d.DatastoreURI)
	}

	d.mountPath = ""
	return nil
}
