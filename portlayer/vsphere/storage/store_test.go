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

package storage

import (
	"archive/tar"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/test"
	portlayer "github.com/vmware/vic/portlayer/storage"
	"golang.org/x/net/context"
)

func setup(t *testing.T) (*portlayer.NameLookupCache, *session.Session, error) {
	datastoreParentPath = "testingParentDirectory"

	client := test.Session(context.TODO(), t)
	if client == nil {
		return nil, nil, fmt.Errorf("skip")
	}

	vsImageStore, err := NewImageStore(context.TODO(), client)
	if err != nil {
		if err.Error() == "can't find the hosting vm" {
			t.Skip("Skipping: test must be run in a VM")
		}
		return nil, nil, err
	}

	s := &portlayer.NameLookupCache{
		DataStore: vsImageStore,
	}

	return s, client, nil
}

func TestRestartImageStore(t *testing.T) {
	// Start the image store once
	_, client, err := setup(t)
	if !assert.NoError(t, err) {
		return
	}
	defer rm(t, client, "")

	// now start it again
	vsImageStore, err := NewImageStore(context.TODO(), client)
	if !assert.NoError(t, err) || !assert.NotNil(t, vsImageStore) {
		return
	}
}

// Create an image store then test it exists
func TestCreateAndGetImageStore(t *testing.T) {
	vsis, client, err := setup(t)
	if !assert.NoError(t, err) {
		return
	}

	// Nuke the parent image store directory
	defer rm(t, client, "")

	storeName := "storeName"
	u, err := vsis.CreateImageStore(context.TODO(), storeName)
	if !assert.NoError(t, err) || !assert.NotNil(t, u) {
		return
	}

	u, err = vsis.GetImageStore(context.TODO(), storeName)
	if !assert.NoError(t, err) || !assert.NotNil(t, u) {
		return
	}

	// Negative test.  Check for a dir that doesnt exist
	u, err = vsis.GetImageStore(context.TODO(), storeName+"garbage")
	if !assert.Error(t, err) || !assert.Nil(t, u) {
		return
	}
}

func TestListImageStore(t *testing.T) {
	vsis, client, err := setup(t)
	if !assert.NoError(t, err) {
		return
	}

	// Nuke the parent image store directory
	defer rm(t, client, "")

	count := 3
	for i := 0; i < count; i++ {
		storeName := fmt.Sprintf("storeName%d", i)
		u, err := vsis.CreateImageStore(context.TODO(), storeName)
		if !assert.NoError(t, err) || !assert.NotNil(t, u) {
			return
		}
	}

	images, err := vsis.ListImageStores(context.TODO())
	if !assert.NoError(t, err) || !assert.Equal(t, len(images), count) {
		return
	}
}

// Creates a tar archive in memory for each layer and uses this to test image creation of layers
func TestCreateImageLayers(t *testing.T) {
	numLayers := 3

	vsis, client, err := setup(t)
	if !assert.NoError(t, err) {
		return
	}

	// Nuke the parent image store directory
	defer rm(t, client, "")

	storeURL, err := vsis.CreateImageStore(context.TODO(), "testStore")
	if !assert.NoError(t, err) {
		return
	}
	// base this image off scratch
	parent, err := vsis.GetImage(context.TODO(), storeURL, portlayer.Scratch.ID)
	if !assert.NoError(t, err) {
		return
	}

	// Keep a list of all files we're extracting via layers so we can verify
	// they exist in the leaf layer.  Ext adds lost+found, so add it here.
	expected := []string{"lost+found"}

	for layer := 0; layer < numLayers; layer++ {
		// Create a buffer to write our archive to.
		buf := new(bytes.Buffer)

		// Create a new tar archive.
		tw := tar.NewWriter(buf)

		dirName := fmt.Sprintf("dir%d", layer)

		// Add some files to the archive.
		var files = []struct {
			Name string
			Type byte
			Body string
		}{
			{dirName, tar.TypeDir, ""},
			{dirName + "/readme.txt", tar.TypeReg, "This archive contains some text files."},
			{dirName + "/gopher.txt", tar.TypeReg, "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
			{dirName + "/todo.txt", tar.TypeReg, "Get animal handling license."},
		}

		for _, file := range files {
			hdr := &tar.Header{
				Name:     file.Name,
				Mode:     0777,
				Typeflag: file.Type,
				Size:     int64(len(file.Body)),
			}

			if err := tw.WriteHeader(hdr); err != nil {
				log.Fatalln(err)
			}

			expected = append(expected, file.Name)

			if file.Type == tar.TypeDir {
				continue
			}

			if _, err := tw.Write([]byte(file.Body)); err != nil {
				log.Fatalln(err)
			}
		}

		// Make sure to check the error on Close.
		if err := tw.Close(); err != nil {
			log.Fatalln(err)
		}

		h := sha256.New()
		h.Write(buf.Bytes())
		sum := fmt.Sprintf("sha256:%x", h.Sum(nil))

		newImage, err := vsis.WriteImage(context.TODO(), parent, dirName, sum, buf)
		if !assert.NoError(t, err) || !assert.NotNil(t, newImage) {
			return
		}

		// make the next image a child of the one we just created
		parent = newImage
	}

	// verify we did anything by attaching the last layer rdonly
	v := vsis.DataStore.(*ImageStore)

	roDisk, err := mountLayerRO(v, parent)
	if !assert.NoError(t, err) {
		return
	}

	p, err := roDisk.MountPath()
	if !assert.NoError(t, err) {
		return
	}

	defer v.dm.Detach(context.TODO(), roDisk)
	defer os.RemoveAll(p)
	defer roDisk.Unmount()

	actual := []string{}
	// Diff the contents
	err = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		f := path[len(p):]
		if f != "" {
			// strip the slash
			actual = append(actual, f[1:])
		}
		return nil
	})
	if !assert.NoError(t, err) {
		return
	}

	sort.Strings(actual)
	sort.Strings(expected)

	log.Infof("actual = %s", actual)
	log.Infof("expected = %s", expected)
	if !assert.Equal(t, expected, actual) {
		return
	}
}

func mountLayerRO(v *ImageStore, parent *portlayer.Image) (*disk.VirtualDisk, error) {
	roName := v.imageStoreDatastoreURI("testStore") + "/" + parent.ID + "-ro.vmdk"
	parentDsURI := v.imageDiskDatastoreURI("testStore", parent.ID)
	roDisk, err := v.dm.CreateAndAttach(context.TODO(), roName, parentDsURI, 0, os.O_RDONLY)
	if err != nil {
		return nil, err
	}

	dir, err := ioutil.TempDir("", parent.ID+"ro")
	if err != nil {
		return nil, err
	}

	if err := roDisk.Mount(dir, nil); err != nil {
		return nil, err
	}

	return roDisk, nil
}

func rm(t *testing.T, client *session.Session, name string) {
	fm := object.NewFileManager(client.Vim25())
	p := client.Datastore.Path(path.Join(datastoreParentPath, name))
	task, err := fm.DeleteDatastoreFile(context.TODO(), p, nil)
	if !assert.NoError(t, err) {
		return
	}
	_, err = task.WaitForResult(context.TODO(), nil)
	if !assert.NoError(t, err) {
		return
	}
}
