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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/archive"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/exec"
	portlayer "github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/lib/portlayer/util"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

// All paths on the datastore for images are relative to <datastore>/VIC/
var StorageParentDir = "VIC"

const (
	StorageImageDir  = "images"
	defaultDiskLabel = "containerfs"
	defaultDiskSize  = 8388608
	metaDataDir      = "imageMetadata"
	manifest         = "manifest"
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

func NewImageStore(ctx context.Context, s *session.Session, u *url.URL) (*ImageStore, error) {
	dm, err := disk.NewDiskManager(ctx, s)
	if err != nil {
		return nil, err
	}

	datastores, err := s.Finder.DatastoreList(ctx, u.Host)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Host returned error when trying to locate provided datastore %s: %s", u.String(), err.Error()))
	}

	if len(datastores) != 1 {
		return nil, errors.New(fmt.Sprintf("Found %d datastores with provided datastore path %s. Cannot create image store.", len(datastores), u.String()))
	}

	ds, err := datastore.NewHelper(ctx, s, datastores[0], path.Join(u.Path, StorageParentDir))
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

func (v *ImageStore) imageDiskPath(storeName, imageName string) string {
	return path.Join(v.imageDirPath(storeName, imageName), imageName+".vmdk")
}

// Returns the path to the vmdk itself in datastore url format
func (v *ImageStore) imageDiskDSPath(storeName, imageName string) string {
	return path.Join(v.ds.RootURL, v.imageDiskPath(storeName, imageName))
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
		// This is startup.  Look for image directories without manifest files and
		// nuke them.
		if err := v.cleanup(ctx, u); err != nil {
			return nil, err
		}

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

		// persist the relationship
		v.parents.Add(ID, parent.ID)

		if err := v.parents.Save(ctx); err != nil {
			return nil, err
		}

		if err := v.writeImage(ctx, storeName, parent.ID, ID, meta, sum, r); err != nil {
			return nil, err
		}

	}

	newImage := &portlayer.Image{
		ID:         ID,
		SelfLink:   imageURL,
		ParentLink: parent.SelfLink,
		Store:      parent.Store,
		Metadata:   meta,
	}

	return newImage, nil
}

// Create the image directory, create a temp vmdk in this directory,
// attach/mount the disk, unpack the tar, check the checksum.  If the data
// doesn't match the expected checksum, abort by nuking the image directory.
// If everything matches, move the tmp vmdk to ID.vmdk.  The unwind path is a
// bit convoluted here;  we need to clean up on the way out in the error case
func (v *ImageStore) writeImage(ctx context.Context, storeName, parentID, ID string, meta map[string][]byte,
	sum string, r io.Reader) error {

	// Create a temp image directory in the store.
	imageDir := v.imageDirPath(storeName, ID)
	_, err := v.ds.Mkdir(ctx, true, imageDir)
	if err != nil {
		return err
	}

	// Write the metadata to the datastore
	metaDataDir := v.imageMetadataDirPath(storeName, ID)
	err = writeMetadata(ctx, v.ds, metaDataDir, meta)
	if err != nil {
		return err
	}

	// datastore path to the parent
	parentDiskDsURI := v.imageDiskDSPath(storeName, parentID)

	// datastore path to the disk we're creating
	diskDsURI := v.imageDiskDSPath(storeName, ID)
	log.Infof("Creating image %s (%s)", ID, diskDsURI)

	var vmdisk *disk.VirtualDisk

	// On error, unmount if mounted, detach if attached, and nuke the image directory
	defer func() {
		if err != nil {
			log.Errorf("Cleaning up failed WriteImage directory %s", imageDir)

			if vmdisk != nil {
				if vmdisk.Mounted() {
					log.Debugf("Unmounting abandoned disk")
					vmdisk.Unmount()
				}

				if vmdisk.Attached() {
					log.Debugf("Detaching abandoned disk")
					v.dm.Detach(ctx, vmdisk)
				}
			}

			v.ds.Rm(ctx, imageDir)
		}
	}()

	// Create the disk
	vmdisk, err = v.dm.CreateAndAttach(ctx, diskDsURI, parentDiskDsURI, 0, os.O_RDWR)
	if err != nil {
		return err
	}
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
		err = fmt.Errorf("Failed to validate image checksum. Expected %s, got %s", sum, actualSum)
		return err
	}

	if err = vmdisk.Unmount(); err != nil {
		return err
	}

	if err = v.dm.Detach(ctx, vmdisk); err != nil {
		return err
	}

	// Write our own bookkeeping manifest file to the image's directory.  We
	// treat the manifest file like a done file.  Its existence means this vmdk
	// is consistent.
	if err = v.writeManifest(ctx, storeName, ID, nil); err != nil {
		return err
	}

	return nil
}

func (v *ImageStore) scratch(ctx context.Context, storeName string) error {

	// Create the image directory in the store.
	imageDir := v.imageDirPath(storeName, portlayer.Scratch.ID)
	if _, err := v.ds.Mkdir(ctx, false, imageDir); err != nil {
		return err
	}

	// Write the metadata to the datastore
	metaDataDir := v.imageMetadataDirPath(storeName, portlayer.Scratch.ID)
	if err := writeMetadata(ctx, v.ds, metaDataDir, nil); err != nil {
		return err
	}

	imageDiskDsURI := v.imageDiskDSPath(storeName, portlayer.Scratch.ID)
	log.Infof("Creating image %s (%s)", portlayer.Scratch.ID, imageDiskDsURI)

	var size int64
	size = defaultDiskSize
	if portlayer.Config.ScratchSize != 0 {
		size = portlayer.Config.ScratchSize
	}

	// Create the disk
	vmdisk, err := v.dm.CreateAndAttach(ctx, imageDiskDsURI, "", size, os.O_RDWR)
	if err != nil {
		return err
	}
	defer func() {
		if vmdisk.Attached() {
			v.dm.Detach(ctx, vmdisk)
		}
	}()
	log.Debugf("Scratch disk created with size %d", portlayer.Config.ScratchSize)

	// Make the filesystem and set its label to defaultDiskLabel
	if err = vmdisk.Mkfs(defaultDiskLabel); err != nil {
		return err
	}

	if err = v.dm.Detach(ctx, vmdisk); err != nil {
		return err
	}

	if err = v.writeManifest(ctx, storeName, portlayer.Scratch.ID, nil); err != nil {
		return err
	}

	return nil
}

func (v *ImageStore) GetImage(ctx context.Context, store *url.URL, ID string) (*portlayer.Image, error) {

	defer trace.End(trace.Begin(store.String()))
	storeName, err := util.ImageStoreName(store)
	if err != nil {
		return nil, err
	}

	imageURL, err := util.ImageURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	if err = v.verifyImage(ctx, storeName, ID); err != nil {
		return nil, err
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
		Store:      &s,
		ParentLink: parentURL,
		Metadata:   meta,
	}

	log.Debugf("Returning image from location %s with parent url %s", newImage.SelfLink, newImage.Parent())
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

// DeleteImage deletes an image from the image store.  If the image is in
// use either by way of inheritence or because it's attached to a
// container, this will return an error.
func (v *ImageStore) DeleteImage(ctx context.Context, image *portlayer.Image) error {
	//  check if the image is in use.
	if err := inUse(ctx, image.ID); err != nil {
		log.Errorf("ImageStore: delete image error: %s", err.Error())
		return err
	}

	storeName, err := util.ImageStoreName(image.Store)
	if err != nil {
		return err
	}

	imageDir := v.imageDirPath(storeName, image.ID)
	log.Infof("ImageStore: Deleting %s", imageDir)
	if err := v.ds.Rm(ctx, imageDir); err != nil {
		log.Errorf("ImageStore: delete image error: %s", err.Error())
		return err
	}

	return nil
}

// Find any image directories without the manifest file and remove them.
func (v *ImageStore) cleanup(ctx context.Context, store *url.URL) error {
	log.Infof("Checking for inconsistent images on %s", store.String())

	storeName, err := util.ImageStoreName(store)
	if err != nil {
		return err
	}

	res, err := v.ds.Ls(ctx, v.imageStorePath(storeName))
	if err != nil {
		return err
	}

	for _, f := range res.File {
		file, ok := f.(*types.FileInfo)
		if !ok {
			continue
		}

		ID := file.Path

		if ID == portlayer.Scratch.ID {
			continue
		}

		if err := v.verifyImage(ctx, storeName, ID); err != nil {
			imageDir := v.imageDirPath(storeName, ID)
			log.Infof("Removing inconsistent image (%s) %s", ID, imageDir)

			// Eat the error so we can continue cleaning up.  The tasks package will log the error if there is one.
			_ = v.ds.Rm(ctx, imageDir)

		}
	}

	return nil
}

// Manifest file for the image.
func (v *ImageStore) writeManifest(ctx context.Context, storeName, ID string, r io.Reader) error {
	pth := path.Join(v.imageDirPath(storeName, ID), manifest)
	if err := v.ds.Upload(ctx, r, pth); err != nil {
		return err
	}

	return nil
}

// check for the manifest file AND the vmdk
func (v *ImageStore) verifyImage(ctx context.Context, storeName, ID string) error {
	imageDir := v.imageDirPath(storeName, ID)

	// Check for teh manifiest file and the vmdk
	for _, p := range []string{path.Join(imageDir, manifest), v.imageDiskPath(storeName, ID)} {
		if _, err := v.ds.Stat(ctx, p); err != nil {
			return err
		}
	}

	return nil
}

// XXX TODO This should be tied to an interface so we don't have to import exec
// here (or wherever the cache lives).
func inUse(ctx context.Context, ID string) error {
	// XXX why doesnt this ever return an error?  Strange.
	// Gather all the running containers and the images they are refering to.
	conts := exec.Containers(true)
	if len(conts) == 0 {
		return nil
	}

	for _, cont := range conts {
		layerID := cont.ExecConfig.LayerID
		if layerID == ID {
			return &portlayer.ErrImageInUse{
				fmt.Sprintf("image %s in use by %s", layerID, cont.ExecConfig.ID),
			}
		}
	}

	return nil
}
