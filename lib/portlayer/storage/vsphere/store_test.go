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

package vsphere

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

	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	portlayer "github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

func setup(t *testing.T) (*portlayer.NameLookupCache, *session.Session, error) {
	StorageParentDir = RandomString(10) + "imageTests"

	client := Session(context.TODO(), t)
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

	s := portlayer.NewLookupCache(vsImageStore)

	return s, client, nil
}

func TestRestartImageStore(t *testing.T) {
	// Start the image store once
	cacheStore, client, err := setup(t)
	if !assert.NoError(t, err) {
		return
	}
	defer rm(t, client, client.Datastore.Path(StorageParentDir))

	origVsStore := cacheStore.DataStore.(*ImageStore)

	// now start it again
	restartedVsStore, err := NewImageStore(context.TODO(), client)
	if !assert.NoError(t, err) || !assert.NotNil(t, restartedVsStore) {
		return
	}

	// Check we didn't create a new UUID directory (relevant if vsan)
	if !assert.NotEqual(t, origVsStore.ds.rooturl, restartedVsStore.ds.rootdir) {
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
	defer rm(t, client, client.Datastore.Path(StorageParentDir))

	storeName := "bogusStoreName"
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

	// Test for a store that already exists
	u, err = vsis.CreateImageStore(context.TODO(), storeName)
	if !assert.Error(t, err) || !assert.Nil(t, u) || !assert.Equal(t, err, os.ErrExist) {
		return
	}
}

func TestListImageStore(t *testing.T) {
	vsis, client, err := setup(t)
	if !assert.NoError(t, err) {
		return
	}

	// Nuke the parent image store directory
	defer rm(t, client, client.Datastore.Path(StorageParentDir))

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
	numLayers := 4

	cacheStore, client, err := setup(t)
	if !assert.NoError(t, err) {
		return
	}

	vsStore := cacheStore.DataStore.(*ImageStore)

	// Nuke the files and then the parent dir.  Unfortunately, because this is
	// vsan, we need to delete the files in the directories first (maybe
	// because they're linked vmkds) before we can delete the parent directory.
	defer func() {
		res, err := vsStore.ds.LsDirs(context.TODO(), "")
		if err != nil {
			t.Logf("error: %s", err)
			return
		}

		for _, dir := range res.HostDatastoreBrowserSearchResults {
			for _, f := range dir.File {
				file, ok := f.(*types.FileInfo)
				if !ok {
					continue
				}

				rm(t, client, path.Join(dir.FolderPath, file.Path))
			}
			rm(t, client, dir.FolderPath)
		}

		rm(t, client, client.Datastore.Path(StorageParentDir))
	}()

	storeURL, err := cacheStore.CreateImageStore(context.TODO(), "testStore")
	if !assert.NoError(t, err) {
		return
	}

	// Get an image that doesn't exist and check for error
	grbg, err := cacheStore.GetImage(context.TODO(), storeURL, "garbage")
	if !assert.Error(t, err) || !assert.Nil(t, grbg) {
		return
	}

	// base this image off scratch
	parent, err := cacheStore.GetImage(context.TODO(), storeURL, portlayer.Scratch.ID)
	if !assert.NoError(t, err) {
		return
	}
	parent.Metadata = make(map[string][]byte)

	// Keep a list of all files we're extracting via layers so we can verify
	// they exist in the leaf layer.  Ext adds lost+found, so add it here.
	expectedFilesOnDisk := []string{"lost+found"}

	// Keep a list of images we created
	expectedImages := make(map[string]*portlayer.Image)
	expectedImages[parent.ID] = parent

	for layer := 0; layer < numLayers; layer++ {

		dirName := fmt.Sprintf("dir%d", layer)
		// Add some files to the archive.
		var files = []tarFile{
			{dirName, tar.TypeDir, ""},
			{dirName + "/readme.txt", tar.TypeReg, "This archive contains some text files."},
			{dirName + "/gopher.txt", tar.TypeReg, "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
			{dirName + "/todo.txt", tar.TypeReg, "Get animal handling license."},
		}

		for _, i := range files {
			expectedFilesOnDisk = append(expectedFilesOnDisk, i.Name)
		}

		// meta for the image
		meta := make(map[string][]byte)
		meta[dirName+"_meta"] = []byte("Some Meta")
		meta[dirName+"_moreMeta"] = []byte("Some More Meta")
		meta[dirName+"_scorpions"] = []byte("Here I am, rock you like a hurricane")

		// Tar the files
		buf, terr := tarFiles(files, meta)
		if !assert.NoError(t, terr) {
			return
		}

		// Calculate the checksum
		h := sha256.New()
		h.Write(buf.Bytes())
		sum := fmt.Sprintf("sha256:%x", h.Sum(nil))

		// Write the image via the cache (which writes to the vsphere impl)
		writtenImage, terr := cacheStore.WriteImage(context.TODO(), parent, dirName, meta, sum, buf)
		if !assert.NoError(t, terr) || !assert.NotNil(t, writtenImage) {
			return
		}

		expectedImages[dirName] = writtenImage

		// Get the image directly via the vsphere image store impl.
		vsImage, terr := vsStore.GetImage(context.TODO(), parent.Store, dirName)
		if !assert.NoError(t, terr) || !assert.NotNil(t, vsImage) {
			return
		}

		assert.Equal(t, writtenImage, vsImage)

		// make the next image a child of the one we just created
		parent = writtenImage
	}

	// Test list images on the datastore
	listedImages, err := vsStore.ListImages(context.TODO(), parent.Store, nil)
	if !assert.NoError(t, err) || !assert.NotNil(t, listedImages) {
		return
	}
	for _, img := range listedImages {
		if !assert.Equal(t, expectedImages[img.ID], img) {
			return
		}
	}

	// verify the disk's data by attaching the last layer rdonly
	roDisk, err := mountLayerRO(vsStore, parent)
	if !assert.NoError(t, err) {
		return
	}

	p, err := roDisk.MountPath()
	if !assert.NoError(t, err) {
		return
	}

	defer vsStore.dm.Detach(context.TODO(), roDisk)
	defer os.RemoveAll(p)
	defer roDisk.Unmount()

	filesFoundOnDisk := []string{}
	// Diff the contents of the RO file of the last child (with all of the contents)
	err = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		f := path[len(p):]
		if f != "" {
			// strip the slash
			filesFoundOnDisk = append(filesFoundOnDisk, f[1:])
		}
		return nil
	})
	if !assert.NoError(t, err) {
		return
	}

	sort.Strings(filesFoundOnDisk)
	sort.Strings(expectedFilesOnDisk)

	if !assert.Equal(t, expectedFilesOnDisk, filesFoundOnDisk) {
		return
	}
}

type tarFile struct {
	Name string
	Type byte
	Body string
}

func tarFiles(files []tarFile, meta map[string][]byte) (*bytes.Buffer, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new tar archive.
	tw := tar.NewWriter(buf)

	// Write data to the tar as if it came from the hub
	for _, file := range files {
		hdr := &tar.Header{
			Name:     file.Name,
			Mode:     0777,
			Typeflag: file.Type,
			Size:     int64(len(file.Body)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}

		if file.Type == tar.TypeDir {
			continue
		}

		if _, err := tw.Write([]byte(file.Body)); err != nil {
			return nil, err
		}
	}

	// Make sure to check the error on Close.
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func mountLayerRO(v *ImageStore, parent *portlayer.Image) (*disk.VirtualDisk, error) {
	roName := v.imageDiskPath("testStore", parent.ID) + "-ro.vmdk"
	parentDsURI := v.imageDiskPath("testStore", parent.ID)

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
	task, err := fm.DeleteDatastoreFile(context.TODO(), name, client.Datacenter)
	if !assert.NoError(t, err) {
		return
	}
	_, _ = task.WaitForResult(context.TODO(), nil)
}
