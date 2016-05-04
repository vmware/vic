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
	"os"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/test"
	"golang.org/x/net/context"
)

// Create a lineage of disks inheriting from eachother, write portion of a
// string to each, the confirm the result is the whole string
func TestCreateAndDetach(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	client := test.Session(context.Background(), t)
	if client == nil {
		return
	}

	imagestore := client.Datastore.Path("imagestore")

	fm := object.NewFileManager(client.Vim25())

	// create a directory in the datastore
	// eat the error because we dont care if it exists
	fm.MakeDirectory(context.TODO(), imagestore, nil, true)

	vdm, err := NewDiskManager(context.TODO(), client)
	if err != nil && err.Error() == "can't find the hosting vm" {
		t.Skip("Skipping: test must be run in a VM")
	}

	if !assert.NoError(t, err) || !assert.NotNil(t, vdm) {
		return
	}

	diskSize := int64(1 << 10)
	parent, err := vdm.Create(context.TODO(), client.Datastore.Path("imagestore/scratch.vmdk"), diskSize)
	if !assert.NoError(t, err) {
		return
	}

	numChildren := 3
	children := make([]*VirtualDisk, numChildren)

	testString := "Ground control to Major Tom"
	writeSize := len(testString) / numChildren
	// Create children which inherit from eachother
	for i := 0; i < numChildren; i++ {

		p := client.Datastore.Path(fmt.Sprintf("imagestore/child%d.vmdk", i))
		child, cerr := vdm.CreateAndAttach(context.TODO(), p, parent.DatastoreURI, 0, os.O_RDWR)
		if !assert.NoError(t, cerr) {
			return
		}

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
			if !assert.NoError(t, cerr) {
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
			if !assert.True(t, strings.Compare(testString, string(b)) == 0) {
				return
			}

		} else {
			_, cerr = f.WriteAt([]byte(testString[start:end]), int64(start))
			if !assert.NoError(t, cerr) {
				return
			}
		}

		f.Close()

		cerr = vdm.Detach(context.TODO(), child)
		if !assert.NoError(t, cerr) {
			return
		}

		// use this image as the next parent
		parent = child
	}

	//	// Nuke the images
	//	for i := len(children) - 1; i >= 0; i-- {
	//		err = vdm.Delete(context.TODO(), children[i])
	//		if !assert.NoError(t, err) {
	//			return
	//		}
	//	}

	// Nuke the image store
	_, err = tasks.WaitForResult(context.TODO(), func(ctx context.Context) (tasks.ResultWaiter, error) {
		return fm.DeleteDatastoreFile(ctx, imagestore, nil)
	})

	if !assert.NoError(t, err) {
		return
	}
}
