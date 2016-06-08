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
	"io/ioutil"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/mount"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	"golang.org/x/net/context"
)

// Create a disk, make an ext filesystem on it, set the label, mount it,
// unmount it, then clean up.
func TestCreateFS(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	client := Session(context.Background(), t)
	if client == nil {
		return
	}

	imagestore := client.Datastore.Path("imagestore")

	fm := object.NewFileManager(client.Vim25())

	// create a directory in the datastore
	// eat the error because we dont care if it exists
	fm.MakeDirectory(context.TODO(), imagestore, nil, true)

	// Nuke the image store
	defer func() {
		task, err := fm.DeleteDatastoreFile(context.TODO(), imagestore, nil)
		if err != nil && err.Error() == "can't find the hosting vm" {
			t.Skip("Skipping: test must be run in a VM")
		}
		if !assert.NoError(t, err) {
			return
		}
		_, err = task.WaitForResult(context.TODO(), nil)
		if !assert.NoError(t, err) {
			return
		}
	}()

	vdm, err := NewDiskManager(context.TODO(), client)
	if err != nil && err.Error() == "can't find the hosting vm" {
		t.Skip("Skipping: test must be run in a VM")
	}
	if !assert.NoError(t, err) || !assert.NotNil(t, vdm) {
		return
	}

	diskSize := int64(1 << 10)
	d, err := vdm.CreateAndAttach(context.TODO(), client.Datastore.Path("imagestore/scratch.vmdk"), "", diskSize, os.O_RDWR)
	if !assert.NoError(t, err) {
		return
	}
	defer vdm.Detach(context.TODO(), d)

	// make the filesysetem
	if err = d.Mkfs("foo"); !assert.NoError(t, err) {
		return
	}

	// set the label
	if err = d.SetLabel("foo"); !assert.NoError(t, err) {
		return
	}

	// make a tempdir to mount the fs to
	dir, err := ioutil.TempDir("", "mnt")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(dir)

	// do the mount
	err = d.Mount(dir, nil)
	if !assert.NoError(t, err) {
		return
	}

	// boom
	if mounted, err := mount.Mounted(dir); !assert.NoError(t, err) || !assert.True(t, mounted) {
		return
	}

	//  clean up
	err = d.Unmount()
	if !assert.NoError(t, err) {
		return
	}
}
