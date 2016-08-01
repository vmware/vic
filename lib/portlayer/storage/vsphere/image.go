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
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/archive"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	portlayer "github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"golang.org/x/net/context"
)

// All paths on the datastore for images are relative to <datastore>/VIC/
var StorageParentDir = "VIC"

const (
	StorageImageDir  = "images"
	defaultDiskLabel = "containerfs"
	defaultDiskSize  = 8388608
	metaDataDir      = "imageMetadata"
)

type ImageStore struct {
	dm *disk.Manager

	// govmomi session
	s *session.Session

	ds *datastore.Helper

	// Parent relationships
	// This will go away when First Class Disk support is added to vsphere.
	// Currently, we can't get a disk spec for a disk outside of creating the
	// disk (and the spec).  This spec has the parent relationship for the
	// disk.  So, for now, persist this data in the datastore and look it up
	// when we need it.
	parents *parentM
}

func NewImageStore(ctx context.Context, s *session.Session) (*ImageStore, error) {
	dm, err := disk.NewDiskManager(ctx, s)
	if err != nil {
		return nil, err
	}

	// Currently using the datastore associated with the session which is not
	// ideal.  This should be passed in via the config.  The datastore need not
	// be the same datastore used for the rest of the system.
	ds, err := datastore.NewHelper(ctx, s, s.Datastore, StorageParentDir)
	if err != nil {
		return nil, err
	}

	vis := &ImageStore{
		dm: dm,
		ds: ds,
		s:  s,
	}

	return vis, nil
}

// Returns the path to a given image store.  Currently this is the UUID of the VCH.
// `/VIC/imageStoreName (currently the vch uuid)/images`
func (v *ImageStore) imageStorePath(storeName string) string {
	return path.Join(storeName, StorageImageDir)
}

// Returns the path to the image relative to the given
// store.  The dir structure for an image in the datastore is
// `/VIC/imageStoreName (currently the vch uuid)/imageName/imageName.vmkd`
func (v *ImageStore) imageDirPath(storeName, imageName string) string {
	return path.Join(v.imageStorePath(storeName), imageName)
}

// Returns the path to the vmdk itself in datastore url format
func (v *ImageStore) imageDiskPath(storeName, imageName string) string {
	return path.Join(v.ds.RootURL, v.imageDirPath(storeName, imageName), imageName+".vmdk")
}

// Returns the path to the metadata directory for an image
func (v *ImageStore) imageMetadataDirPath(storeName, imageName string) string {
	return path.Join(v.imageDirPath(storeName, imageName), metaDataDir)
}

func (v *ImageStore) CreateImageStore(ctx context.Context, storeName string) (*url.URL, error) {
	// convert the store name to a port layer url.
	u, err := util.ImageStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	if _, err = v.ds.Mkdir(ctx, true, v.imageStorePath(storeName)); err != nil {
		return nil, err
	}

	if v.parents == nil {
		pm, err := restoreParentMap(ctx, v.ds, storeName)
		if err != nil {
			return nil, err
		}
		v.parents = pm
	}
	return u, nil
}

// GetImageStore checks to see if the image store exists on disk and returns an
// error or the store's URL.
func (v *ImageStore) GetImageStore(ctx context.Context, storeName string) (*url.URL, error) {
	u, err := util.ImageStoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	p := v.imageStorePath(storeName)
	info, err := v.ds.Stat(ctx, p)
	if err != nil {
		return nil, err
	}

	_, ok := info.(*types.FolderFileInfo)
	if !ok {
		return nil, fmt.Errorf("Stat error:  path doesn't exist (%s)", p)
	}

	if v.parents == nil {
		pm, err := restoreParentMap(ctx, v.ds, storeName)
		if err != nil {
			return nil, err
		}
		v.parents = pm
	}

	return u, nil
}

func (v *ImageStore) ListImageStores(ctx context.Context) ([]*url.URL, error) {
	res, err := v.ds.Ls(ctx, v.imageStorePath(""))
	if err != nil {
		return nil, err
	}

	stores := []*url.URL{}
	for _, f := range res.File {
		folder, ok := f.(*types.FolderFileInfo)
		if !ok {
			continue
		}
		u, err := util.ImageStoreNameToURL(folder.Path)
		if err != nil {
			return nil, err
		}
		stores = append(stores, u)

	}

	return stores, nil
}

// WriteImage creates a new image layer from the given parent.
// Eg parentImage + newLayer = new Image built from parent
//
// parent - The parent image to create the new image from.
// ID - textual ID for the image to be written
// meta - metadata associated with the image
// Tag - the tag of the image to be written
func (v *ImageStore) WriteImage(ctx context.Context, parent *portlayer.Image, ID string, meta map[string][]byte, sum string,
	r io.Reader) (*portlayer.Image, error) {

	storeName, err := util.ImageStoreName(parent.Store)
	if err != nil {
		return nil, err
	}

	imageURL, err := util.ImageURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	// If this is scratch, then it's the root of the image store.  All images
	// will be descended from this created and prepared fs.
	if ID == portlayer.Scratch.ID {
		// Create the scratch layer
		if err := v.scratch(ctx, storeName); err != nil {
			return nil, err
		}
	} else {

		if parent.ID == "" {
			return nil, fmt.Errorf("parent ID is empty")
		}

		if err := v.writeImage(ctx, storeName, parent.ID, ID, sum, r); err != nil {
			return nil, err
		}

		// persist the relationship
		v.parents.Add(ID, parent.ID)

		if err := v.parents.Save(ctx); err != nil {
			return nil, err
		}
	}

	// Write the metadata to the datastore
	metaDataDir := v.imageMetadataDirPath(storeName, ID)
	err = writeMetadata(ctx, v.ds, metaDataDir, meta)
	if err != nil {
		return nil, err
	}

	newImage := &portlayer.Image{
		ID:       ID,
		SelfLink: imageURL,
		Parent:   parent.SelfLink,
		Store:    parent.Store,
		Metadata: meta,
	}

	return newImage, nil
}

// Create a temporary directory, create a vmdk in this directory, attach/mount
// the disk, unpack the tar, check the checksum.  If the data doesn't match the
// expected checksum, abort by nuking the tempdir.  If everything matches, move
// the tmpdir to the expected location (with the vmdk inside it).  The unwind
// path is a bit convoluted here;  we need to clean up on the way out in the
// error case (using the tmpdir).
func (v *ImageStore) writeImage(ctx context.Context, storeName, parentID, ID,
	sum string, r io.Reader) error {

	// Create a temp image directory in the store.
	tmpImageDir := v.imageDirPath(storeName, ID+"inprogress")
	_, err := v.ds.Mkdir(ctx, true, tmpImageDir)
	if err != nil {
		return err
	}

	// datastore path to the parent
	parentDiskDsURI := v.imageDiskPath(storeName, parentID)

	// datastore path to the disk we're creating
	imageDiskDsURI := path.Join(v.ds.RootURL, tmpImageDir, ID+".vmdk")
	log.Infof("Creating image %s (%s)", ID, imageDiskDsURI)

	// Create the disk
	vmdisk, err := v.dm.CreateAndAttach(ctx, imageDiskDsURI, parentDiskDsURI, 0, os.O_RDWR)
	if err != nil {
		return err
	}

	defer func() {
		if vmdisk.Mounted() {
			log.Debugf("Unmounting abandonned disk")
			vmdisk.Unmount()
		}
		if vmdisk.Attached() {
			log.Debugf("Detaching abandonned disk")
			v.dm.Detach(ctx, vmdisk)
		}

		v.ds.Rm(ctx, tmpImageDir)
	}()

	// tmp dir to mount the disk
	dir, err := ioutil.TempDir("", "mnt-"+ID)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if err := vmdisk.Mount(dir, nil); err != nil {
		return err
	}

	h := sha256.New()
	t := io.TeeReader(r, h)

	// Untar the archive
	if err = archive.Untar(t, dir, &archive.TarOptions{}); err != nil {
		return err
	}

	actualSum := fmt.Sprintf("sha256:%x", h.Sum(nil))
	if actualSum != sum {
		return fmt.Errorf("Failed to validate image checksum. Expected %s, got %s", sum, actualSum)
	}

	if err = vmdisk.Unmount(); err != nil {
		return err
	}

	if err = v.dm.Detach(ctx, vmdisk); err != nil {
		return err
	}

	// If we've gotten here, prepare the final location for the image
	imageDir := v.imageDirPath(storeName, ID)
	_, err = v.ds.Mkdir(ctx, true, imageDir)
	if err != nil {
		return err
	}

	// Move the disk to it's proper directory
	vdm := object.NewVirtualDiskManager(v.s.Vim25())
	return tasks.Wait(ctx, func(context.Context) (tasks.Waiter, error) {
		dest := v.imageDiskPath(storeName, ID)
		log.Infof("Moving disk %s to %s", imageDiskDsURI, dest)
		t, err := vdm.MoveVirtualDisk(ctx, imageDiskDsURI, v.s.Datacenter, dest, v.s.Datacenter, true)
		log.Infof("move task = %s", t)
		return t, err
	})
}

func (v *ImageStore) scratch(ctx context.Context, storeName string) error {

	// Create the image directory in the store.
	imageDir := v.imageDirPath(storeName, portlayer.Scratch.ID)
	if _, err := v.ds.Mkdir(ctx, false, imageDir); err != nil {
		return err
	}

	imageDiskDsURI := v.imageDiskPath(storeName, portlayer.Scratch.ID)
	log.Infof("Creating image %s (%s)", portlayer.Scratch.ID, imageDiskDsURI)

	// Create the disk
	vmdisk, err := v.dm.CreateAndAttach(ctx, imageDiskDsURI, "", defaultDiskSize, os.O_RDWR)
	if err != nil {
		return err
	}
	defer v.dm.Detach(ctx, vmdisk)

	// Make the filesystem and set its label to defaultDiskLabel
	if err := vmdisk.Mkfs(defaultDiskLabel); err != nil {
		return err
	}

	return nil
}

func (v *ImageStore) GetImage(ctx context.Context, store *url.URL, ID string) (*portlayer.Image, error) {

	storeName, err := util.ImageStoreName(store)
	if err != nil {
		return nil, err
	}

	imageURL, err := util.ImageURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	p := v.imageDirPath(storeName, ID)
	info, err := v.ds.Stat(ctx, p)
	if err != nil {
		return nil, err
	}

	_, ok := info.(*types.FolderFileInfo)
	if !ok {
		return nil, fmt.Errorf("Stat error:  image doesn't exist (%s)", p)
	}

	// get the metadata
	metaDataDir := v.imageMetadataDirPath(storeName, ID)
	meta, err := getMetadata(ctx, v.ds, metaDataDir)
	if err != nil {
		return nil, err
	}

	var s = *store
	var parentURL *url.URL

	parentID := v.parents.Get(ID)
	if parentID != "" {
		parentURL, _ = util.ImageURL(storeName, parentID)
	}

	newImage := &portlayer.Image{
		ID:       ID,
		SelfLink: imageURL,
		// We're relying on the parent map for this since we don't currently have a
		// way to get the disk's spec.  See VIC #482 for details.  Parent:
		// parent.SelfLink,
		Store:    &s,
		Parent:   parentURL,
		Metadata: meta,
	}

	return newImage, nil
}

func (v *ImageStore) ListImages(ctx context.Context, store *url.URL, IDs []string) ([]*portlayer.Image, error) {

	storeName, err := util.ImageStoreName(store)
	if err != nil {
		return nil, err
	}

	res, err := v.ds.Ls(ctx, v.imageStorePath(storeName))
	if err != nil {
		return nil, err
	}

	images := []*portlayer.Image{}
	for _, f := range res.File {
		file, ok := f.(*types.FileInfo)
		if !ok {
			continue
		}

		ID := file.Path

		img, err := v.GetImage(ctx, store, ID)
		if err != nil {
			return nil, err
		}

		images = append(images, img)
	}

	return images, nil
}
