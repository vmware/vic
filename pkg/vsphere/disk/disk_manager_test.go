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
	"path"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/test/env"
	"golang.org/x/net/context"
)

func Session(ctx context.Context, t *testing.T) *session.Session {
	config := &session.Config{
		Service: env.URL(t),

		DatastorePath: env.DS(t),

		Insecure:  true,
		Keepalive: time.Duration(5) * time.Minute,
	}

	s, err := session.NewSession(config).Create(ctx)
	// Vsan has special UUID / URI handling of top level directories which
	// we've implemented at another level.  We can't import them here or it'd
	// be a circular dependency.  Also, we already have tests that test these
	// cases but from a higher level.
	if err != nil || s.IsVSAN(ctx) {
		t.SkipNow()
	}

	return s
}

// Create a lineage of disks inheriting from eachother, write portion of a
// string to each, the confirm the result is the whole string
func TestCreateAndDetach(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	client := Session(context.Background(), t)
	if client == nil {
		return
	}

	imagestore := client.Datastore.Path(datastore.TestName("diskManagerTest"))

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
	parent, err := vdm.Create(context.TODO(), path.Join(imagestore, "scratch.vmdk"), diskSize)
	if !assert.NoError(t, err) {
		return
	}

	numChildren := 3
	children := make([]*VirtualDisk, numChildren)

	testString := "Ground control to Major Tom"
	writeSize := len(testString) / numChildren
	// Create children which inherit from eachother
	for i := 0; i < numChildren; i++ {

		p := path.Join(imagestore, fmt.Sprintf("child%d.vmdk", i))
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
	_, err = tasks.WaitForResult(context.TODO(), func(ctx context.Context) (tasks.Task, error) {
		return fm.DeleteDatastoreFile(ctx, imagestore, nil)
	})

	if !assert.NoError(t, err) {
		return
	}
}
