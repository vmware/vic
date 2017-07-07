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
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/test/env"
)

func Session(ctx context.Context, t *testing.T) *session.Session {
	config := &session.Config{
		Service: env.URL(t),

		DatastorePath: env.DS(t),

		Insecure:  true,
		Keepalive: time.Duration(5) * time.Minute,
	}

	s := session.NewSession(config)

	_, err := s.Connect(ctx)
	if err != nil {
		t.Skip(err.Error())
		return nil
	}

	// we're treating this as an atomic behaviour, so log out if we failed
	defer func() {
		if err != nil {
			t.Skip(err.Error())
			s.Client.Logout(ctx)
		}
	}()

	_, err = s.Populate(ctx)
	if err != nil {
		t.Skip(err.Error())
		return nil
	}

	// Vsan has special UUID / URI handling of top level directories which
	// we've implemented at another level.  We can't import them here or it'd
	// be a circular dependency.  Also, we already have tests that test these
	// cases but from a higher level.
	if err != nil || s.IsVSAN(ctx) {
		t.Logf("error: %s", err.Error())
		t.SkipNow()
	}

	return s
}

// Create a lineage of disks inheriting from eachother, write portion of a
// string to each, the confirm the result is the whole string
func TestCreateAndDetach(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	session := Session(context.Background(), t)
	if session == nil {
		return
	}

	op := trace.NewOperation(context.TODO(), t.Name())

	_, err := guest.GetSelf(op, session)
	if err != nil {
		t.Skip("Not in a vm")
	}

	imagestore := &object.DatastorePath{
		Datastore: session.Datastore.Name(),
		Path:      datastore.TestName(t.Name()),
	}

	// file manager
	fm := object.NewFileManager(session.Vim25())
	// create a directory in the datastore
	err = fm.MakeDirectory(context.TODO(), imagestore.String(), nil, true)
	if !assert.NoError(t, err) {
		return
	}

	// Nuke the image store
	defer func() {
		task, err := fm.DeleteDatastoreFile(context.TODO(), imagestore.String(), nil)
		if !assert.NoError(t, err) {
			return
		}
		_, err = task.WaitForResult(context.TODO(), nil)
		if !assert.NoError(t, err) {
			return
		}
	}()

	// create a diskmanager
	vdm, err := NewDiskManager(op, session)
	if !assert.NoError(t, err) || !assert.NotNil(t, vdm) {
		return
	}

	diskSize := int64(1 << 10)
	scratch := &object.DatastorePath{
		Datastore: session.Datastore.Name(),
		Path:      path.Join(imagestore.Path, "scratch.vmdk"),
	}
	config := NewPersistentDisk(scratch).WithCapacity(diskSize)
	parent, err := vdm.Create(op, config)
	if !assert.NoError(t, err) {
		return
	}

	numChildren := 3
	children := make([]*VirtualDisk, numChildren)

	testString := "Ground control to Major Tom"
	writeSize := len(testString) / numChildren
	// Create children which inherit from each other
	for i := 0; i < numChildren; i++ {

		p := &object.DatastorePath{
			Datastore: imagestore.Datastore,
			Path:      path.Join(imagestore.Path, fmt.Sprintf("child%d.vmdk", i)),
		}

		config := NewPersistentDisk(p).WithParent(parent.DatastoreURI)
		child, cerr := vdm.CreateAndAttach(op, config)
		if !assert.NoError(t, cerr) {
			return
		}
		refs := child.attachedRefs.Count()
		assert.EqualValues(t, 1, refs, "Expected %d attach references, found %d", refs)

		// attempt to recreate and attach the already attached disk
		config = NewPersistentDisk(p).WithParent(parent.DatastoreURI)
		child, cerr = vdm.CreateAndAttach(op, config)
		if !assert.NoError(t, cerr) {
			return
		}
		refs = child.attachedRefs.Count()
		assert.EqualValues(t, 2, refs, "Expected %d attach references, found %d", refs)

		// attempt detach
		cerr = vdm.Detach(op, config)
		if !assert.NoError(t, cerr) {
			return
		}
		// should not actually detach, and should still have 1 reference
		refs = child.attachedRefs.Count()
		assert.EqualValues(t, 1, refs, "Expected %d attach references, found %d", refs)

		children[i] = child

		// Write directly to the disk
		f, cerr := os.OpenFile(child.DevicePath, os.O_RDWR, os.FileMode(0777))
		if !assert.NoError(t, cerr) {
			return
		}

		start := i * writeSize
		end := start + writeSize

		if i == numChildren-1 {
			// last chunk, write to the end.
			_, cerr = f.WriteAt([]byte(testString[start:]), int64(start))
			if !assert.NoError(t, cerr) || !assert.NoError(t, f.Sync()) {
				return
			}

			// Try to read the whole string
			b := make([]byte, len(testString))
			f.Seek(0, 0)
			_, cerr = f.Read(b)
			if !assert.NoError(t, cerr) {
				return
			}

			//check against the test string
			if !assert.Equal(t, testString, string(b)) {
				return
			}

		} else {
			_, cerr = f.WriteAt([]byte(testString[start:end]), int64(start))
			if !assert.NoError(t, cerr) || !assert.NoError(t, f.Sync()) {
				return
			}
		}

		f.Close()

		cerr = vdm.Detach(op, config)
		if !assert.NoError(t, cerr) {
			return
		}

		// use this image as the next parent
		parent = child
	}
}

func TestRefCounting(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	session := Session(context.Background(), t)
	if session == nil {
		return
	}

	op := trace.NewOperation(context.TODO(), t.Name())

	_, err := guest.GetSelf(op, session)
	if err != nil {
		t.Skip("Not in a vm")
	}

	imagestore := &object.DatastorePath{
		Datastore: session.Datastore.Name(),
		Path:      datastore.TestName(t.Name()),
	}

	// file manager
	fm := object.NewFileManager(session.Vim25())
	// create a directory in the datastore
	err = fm.MakeDirectory(context.TODO(), imagestore.String(), nil, true)
	if !assert.NoError(t, err) {
		return
	}

	// Nuke the image store
	defer func() {
		task, err := fm.DeleteDatastoreFile(context.TODO(), imagestore.String(), nil)
		if !assert.NoError(t, err) {
			return
		}
		_, err = task.WaitForResult(context.TODO(), nil)
		if !assert.NoError(t, err) {
			return
		}
	}()

	// create a diskmanager
	vdm, err := NewDiskManager(op, session)
	if !assert.NoError(t, err) || !assert.NotNil(t, vdm) {
		return
	}

	if !assert.NoError(t, err) || !assert.NotNil(t, vdm) {
		return
	}

	diskSize := int64(1 << 10)
	scratch := &object.DatastorePath{
		Datastore: session.Datastore.Name(),
		Path:      path.Join(imagestore.Path, "scratch.vmdk"),
	}
	config := NewPersistentDisk(scratch).WithCapacity(diskSize)
	p, err := vdm.Create(op, config)
	if !assert.NoError(t, err) {
		return
	}

	assert.False(t, p.Attached(), "%s is attached but should not be", p.DatastoreURI)

	child := &object.DatastorePath{
		Datastore: imagestore.Datastore,
		Path:      path.Join(imagestore.Path, "testDisk.vmdk"),
	}
	config = NewPersistentDisk(child).WithParent(scratch)

	// attempt attach
	assert.NoError(t, vdm.attach(op, config), "Error attempting to attach %s", config)

	devicePath, err := vdm.devicePathByURI(op, child)
	if !assert.NoError(t, err) {
		return
	}

	d, err := NewVirtualDisk(config, vdm.Disks)
	if !assert.NoError(t, err) {
		return
	}

	blockDev, err := waitForDevice(op, devicePath)
	if !assert.NoError(t, err) {
		return
	}

	assert.False(t, d.Attached(), "%s is attached but should not be", d.DatastoreURI)

	// Attach the disk
	assert.NoError(t, d.setAttached(blockDev), "Error attempting to mark %s as attached", d.DatastoreURI)

	assert.True(t, d.Attached(), "%s is not attached but should be", d.DatastoreURI)
	assert.NoError(t, d.canBeDetached(), "%s should be detachable but is not", d.DatastoreURI)
	assert.False(t, d.InUseByOther(), "%s is in use but should not be", d.DatastoreURI)
	assert.Equal(t, 1, d.attachedRefs.Count(), "%s has %d attach references but should have 1", d.DatastoreURI, d.attachedRefs.Count())

	// attempt another attach at disk level to increase reference count
	// TODO(jzt): This should probably eventually use the attach code coming in
	// https://github.com/vmware/vic/issues/5422
	assert.NoError(t, d.setAttached(blockDev), "Error attempting to mark %s as attached", d.DatastoreURI)

	assert.True(t, d.Attached(), "%s is not attached but should be", d.DatastoreURI)
	assert.Error(t, d.canBeDetached(), "%s should not be detachable but is", d.DatastoreURI)
	assert.True(t, d.InUseByOther(), "%s is not in use but should be", d.DatastoreURI)
	assert.Equal(t, 2, d.attachedRefs.Count(), "%s has %d attach references but should have 2", d.DatastoreURI, d.attachedRefs.Count())

	// reduce reference count by calling detach
	assert.NoError(t, d.setDetached(vdm.Disks), "Error attempting to mark %s as detached", d.DatastoreURI)

	assert.True(t, d.Attached(), "%s is not attached but should be", d.DatastoreURI)
	assert.NoError(t, d.canBeDetached(), "%s should be detachable but is not", d.DatastoreURI)
	assert.False(t, d.InUseByOther(), "%s is in use but should not be", d.DatastoreURI)
	assert.Equal(t, 1, d.attachedRefs.Count(), "%s has %d attach references but should have 1", d.DatastoreURI, d.attachedRefs.Count())

	// test mount reference counting
	assert.NoError(t, d.Mkfs("testDisk"), "Error attempting to format %s", d.DatastoreURI)

	// create temp mount path
	dir, err := ioutil.TempDir("", "mnt")
	if !assert.NoError(t, err) {
		return
	}

	// cleanup
	defer func() {
		assert.NoError(t, os.RemoveAll(dir), "Error cleaning up mount path %s", dir)
	}()

	// initial mount
	assert.NoError(t, d.Mount(dir, nil), "Error attempting to mount %s at %s", d.DatastoreURI, dir)

	mountPath, err := d.MountPath()
	if !assert.NoError(t, err) {
		return
	}

	assert.True(t, d.Mounted(), "%s is not mounted but should be", d.DatastoreURI)
	assert.Error(t, d.canBeDetached(), "%s should not be detachable but is", d.DatastoreURI)
	assert.False(t, d.InUseByOther(), "%s is in use but should not be", d.DatastoreURI)
	assert.Equal(t, 1, d.mountedRefs.Count(), "%s has %d mount references but should have 1", d.DatastoreURI, d.mountedRefs.Count())
	assert.Equal(t, dir, mountPath, "%s is mounted at %s but should be mounted at %s", d.DatastoreURI, mountPath, dir)

	// attempt another mount
	assert.NoError(t, d.Mount(dir, nil), "Error attempting to mount %s at %s", d.DatastoreURI, dir)

	assert.True(t, d.Mounted(), "%s is not mounted but should be", d.DatastoreURI)
	assert.Error(t, d.canBeDetached(), "%s should not be detachable but is", d.DatastoreURI)
	assert.True(t, d.InUseByOther(), "%s is not in use but should be", d.DatastoreURI)
	assert.Equal(t, 2, d.mountedRefs.Count(), "%s has %d mount references but should have 2", d.DatastoreURI, d.mountedRefs.Count())

	// attempt unmount
	assert.NoError(t, d.Unmount(), "Error attempting to unmount %s", d.DatastoreURI)

	assert.True(t, d.Mounted(), "%s is not mounted but should be", d.DatastoreURI)
	assert.Error(t, d.canBeDetached(), "%s should not be detachable but is", d.DatastoreURI)
	assert.False(t, d.InUseByOther(), "%s is in use but should not be", d.DatastoreURI)
	assert.Equal(t, 1, d.mountedRefs.Count(), "%s has %d mount references but should have 1", d.DatastoreURI, d.mountedRefs.Count())

	// actually unmount
	assert.NoError(t, d.Unmount(), "Error attempting to unmount %s", d.DatastoreURI)

	assert.False(t, d.Mounted(), "%s is mounted but should not be", d.DatastoreURI)
	assert.NoError(t, d.canBeDetached(), "%s should be detachable but is not", d.DatastoreURI)
	assert.False(t, d.InUseByOther(), "%s is in use but should not be", d.DatastoreURI)
	assert.Equal(t, 0, d.mountedRefs.Count(), "%s has %d mount references but should have 0", d.DatastoreURI, d.mountedRefs.Count())

	// detach
	assert.NoError(t, vdm.Detach(op, config), "Error attempting to detach %s", d.DatastoreURI)

	assert.False(t, d.Attached(), "%s is attached but should not be", d.DatastoreURI)
	assert.False(t, d.Mounted(), "%s is mounted but should not be", d.DatastoreURI)
	assert.Error(t, d.canBeDetached(), "%s should not be detachable but is", d.DatastoreURI)
	assert.False(t, d.InUseByOther(), "%s is in use but should not be", d.DatastoreURI)
	assert.Equal(t, 0, d.attachedRefs.Count(), "%s has %d attach references but should have 0", d.DatastoreURI, d.attachedRefs.Count())
	assert.Equal(t, 0, d.mountedRefs.Count(), "%s has %d mount references but should have 0", d.DatastoreURI, d.mountedRefs.Count())

	if !assert.NoError(t, err) {
		return
	}
}

func TestRefCountingParallel(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	session := Session(context.Background(), t)
	if session == nil {
		return
	}

	op := trace.NewOperation(context.TODO(), t.Name())

	_, err := guest.GetSelf(op, session)
	if err != nil {
		t.Skip("Not in a vm")
	}

	imagestore := &object.DatastorePath{
		Datastore: session.Datastore.Name(),
		Path:      datastore.TestName(t.Name()),
	}

	// file manager
	fm := object.NewFileManager(session.Vim25())
	// create a directory in the datastore
	err = fm.MakeDirectory(context.TODO(), imagestore.String(), nil, true)
	if !assert.NoError(t, err) {
		return
	}

	// Nuke the image store
	defer func() {
		task, err := fm.DeleteDatastoreFile(context.TODO(), imagestore.String(), nil)
		if !assert.NoError(t, err) {
			return
		}
		_, err = task.WaitForResult(context.TODO(), nil)
		if !assert.NoError(t, err) {
			return
		}
	}()

	// create a diskmanager
	vdm, err := NewDiskManager(op, session)
	if !assert.NoError(t, err) || !assert.NotNil(t, vdm) {
		return
	}

	if !assert.NoError(t, err) || !assert.NotNil(t, vdm) {
		return
	}

	diskSize := int64(1 << 10)
	uri := &object.DatastorePath{
		Datastore: session.Datastore.Name(),
		Path:      path.Join(imagestore.Path, "test.vmdk"),
	}
	config := NewPersistentDisk(uri).WithCapacity(diskSize)
	d, err := vdm.CreateAndAttach(op, config)
	if !assert.NoError(t, err) {
		return
	}

	assert.True(t, d.Attached(), "%s is not attached but should be", d.DatastoreURI)
	assert.NoError(t, d.canBeDetached(), "%s should be detachable but is not", d.DatastoreURI)
	assert.False(t, d.InUseByOther(), "%s is in use but should not be", d.DatastoreURI)
	assert.EqualValues(t, 1, d.attachedRefs.Count(), "%s has %d attach references but should have 1", d.DatastoreURI, d.attachedRefs.Count())

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			var err error
			defer wg.Done()

			for j := 0; j < 5; j++ {

				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

				d, err = vdm.CreateAndAttach(op, config)
				if !assert.NoError(t, err) {
					return
				}

				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

				err = vdm.Detach(op, config)
				if !assert.NoError(t, err) {
					return
				}
			}
		}()
	}
	wg.Wait()

	assert.True(t, d.Attached(), "%s is not attached but should be", d.DatastoreURI)
	assert.NoError(t, d.canBeDetached(), "%s should be detachable but is not", d.DatastoreURI)
	assert.False(t, d.InUseByOther(), "%s is in use but should not be", d.DatastoreURI)
	assert.EqualValues(t, 1, d.attachedRefs.Count(), "%s has %d attach references but should have 1", d.DatastoreURI, d.attachedRefs.Count())

	err = vdm.Detach(op, config)
	if !assert.NoError(t, err) {
		log.Error("Error detaching disk: %s", err.Error())
		return
	}

	assert.False(t, d.Attached(), "%s is attached but should not be", d.DatastoreURI)
	assert.Error(t, d.canBeDetached(), "%s should be detachable but is not", d.DatastoreURI)
	assert.False(t, d.InUseByOther(), "%s is in use but should not be", d.DatastoreURI)
	assert.EqualValues(t, 0, d.attachedRefs.Count(), "%s has %d attach references but should have 0", d.DatastoreURI, d.attachedRefs.Count())
}
