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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/archive"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
	portlayer "github.com/vmware/vic/portlayer/storage"
	"github.com/vmware/vic/portlayer/util"
	"golang.org/x/net/context"
)

// All paths on the datastore for images are relative to <datastore>/VIC/
var datastoreParentPath = "VIC"

const (
	defaultDiskLabel = "containerfs"
	defaultDiskSize  = 8388608
	metaDataDir      = "imageMetadata"
)

type ImageStore struct {
	dm *disk.Manager
	fm *object.FileManager

	// govmomi session
	s *session.Session

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

	pm, err := restoreParentMap(ctx, s)
	if err != nil {
		return nil, err
	}

	vis := &ImageStore{
		dm:      dm,
		fm:      object.NewFileManager(s.Vim25()),
		s:       s,
		parents: pm,
	}

	err = vis.makeImageStoreParentDir(ctx)
	if err != nil {
		return nil, err
	}

	return vis, nil
}

// Takes a path and returns the datastore path by prepending the datastore name
// to the path.
func (v *ImageStore) datastorePath(p string) string {
	return v.s.Datastore.Path(p)
}

// Returns the path to a given image store
func (v *ImageStore) imageStorePath(storeName string) string {
	return path.Join(datastoreParentPath, storeName)
}

// Returns the path to the image relative to the given
// store.  The dir structure for an image in the datastore is
// `/VIC/imageStoreName/imageName/imageName.vmkd`
func (v *ImageStore) imageDirPath(storeName, imageName string) string {
	return path.Join(datastoreParentPath, storeName, imageName)
}

// Returns the path to the vmdk itself
func (v *ImageStore) imageDiskPath(storeName, imageName string) string {
	return path.Join(v.imageDirPath(storeName, imageName), imageName+".vmdk")
}

// Returns the path to the metadata directory for an image
func (v *ImageStore) imageMetadataDirPath(storeName, imageName string) string {
	return path.Join(v.imageDirPath(storeName, imageName), metaDataDir)
}

func (v *ImageStore) CreateImageStore(ctx context.Context, storeName string) (*url.URL, error) {
	// convert the store name to a port layer url.
	u, err := util.StoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	// Create a vsphere datastore imagestore directury structure url from the
	// storename.  We create scratch since it's the root of the image store.
	// All images inherit from this root image.
	imagestore := v.datastorePath(v.imageStorePath(storeName))

	log.Infof("Creating imagestore directory %s", imagestore)
	if err = v.fm.MakeDirectory(ctx, imagestore, v.s.Datacenter, false); err != nil {
		soapFault := soap.ToSoapFault(err)
		if _, ok := soapFault.VimFault().(types.FileAlreadyExists); ok {
			// Rest API expects this error
			log.Debugf("File exists. %s", err)
			err = os.ErrExist
		}
		return nil, err
	}
	return u, nil
}

// GetImageStore checks to see if the image store exists on disk and returns an
// error or the store's URL.
func (v *ImageStore) GetImageStore(ctx context.Context, storeName string) (*url.URL, error) {
	u, err := util.StoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	// Since we're statting the datastore itself, this need not be in datastore
	// path format.
	p := path.Join(datastoreParentPath, storeName)
	info, err := v.s.Datastore.Stat(ctx, p)
	if err != nil {
		return nil, err
	}

	_, ok := info.(*types.FolderFileInfo)
	if !ok {
		return nil, fmt.Errorf("Stat error:  path doesn't exist (%s)", p)
	}

	return u, nil
}

func (v *ImageStore) ListImageStores(ctx context.Context) ([]*url.URL, error) {
	res, err := lsDir(ctx, v.s.Datastore, v.datastorePath(v.imageStorePath("")))
	if err != nil {
		return nil, err
	}

	stores := []*url.URL{}
	for _, f := range res.File {
		folder, ok := f.(*types.FolderFileInfo)
		if !ok {
			continue
		}
		u, err := util.StoreNameToURL(folder.Path)
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
func (v *ImageStore) WriteImage(ctx context.Context, parent *portlayer.Image, ID string, meta map[string][]byte,
	r io.Reader) (*portlayer.Image, error) {

	storeName, err := util.StoreName(parent.Store)
	if err != nil {
		return nil, err
	}

	imageURL, err := util.ImageURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	// Create the image directory in the store.
	imageDirDsURI := v.datastorePath(v.imageDirPath(storeName, ID))
	if err = v.fm.MakeDirectory(ctx, imageDirDsURI, v.s.Datacenter, false); err != nil {
		return nil, err
	}

	ImageDiskDsURI := v.datastorePath(v.imageDiskPath(storeName, ID))
	log.Infof("Creating image %s", ID)

	// If this is scratch, then it's the root of the image store.  All images
	// will be descended from this created and prepared fs.
	if ID == portlayer.Scratch.ID {
		// Create the disk
		vmdisk, cerr := v.dm.CreateAndAttach(ctx, ImageDiskDsURI, "", defaultDiskSize, os.O_RDWR)
		if cerr != nil {
			return nil, cerr
		}
		defer v.dm.Detach(ctx, vmdisk)

		// Make the filesystem and set its label to defaultDiskLabel
		if cerr = vmdisk.Mkfs(defaultDiskLabel); cerr != nil {
			return nil, cerr
		}
	} else {

		if parent.ID == "" {
			return nil, fmt.Errorf("parent ID is empty")
		}

		// Create the disk
		parentDiskDsURI := v.datastorePath(v.imageDiskPath(storeName, parent.ID))
		vmdisk, cerr := v.dm.CreateAndAttach(ctx, ImageDiskDsURI, parentDiskDsURI, 0, os.O_RDWR)
		if cerr != nil {
			return nil, cerr
		}
		defer v.dm.Detach(ctx, vmdisk)

		dir, cerr := ioutil.TempDir("", "mnt-"+ID)
		if cerr != nil {
			return nil, cerr
		}
		defer os.RemoveAll(dir)

		if merr := vmdisk.Mount(dir, nil); merr != nil {
			return nil, merr
		}
		defer vmdisk.Unmount()

		// Untar the archive
		cerr = archive.Untar(r, dir, &archive.TarOptions{})
		if cerr != nil {
			return nil, cerr
		}

		// persist the relationship
		v.parents.Add(ID, parent.ID)

		if cerr = v.parents.Save(ctx); cerr != nil {
			return nil, cerr
		}
	}

	// Write the metadata to the datastore
	err = v.writeMeta(ctx, storeName, ID, meta)
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

func (v *ImageStore) GetImage(ctx context.Context, store *url.URL, ID string) (*portlayer.Image, error) {

	storeName, err := util.StoreName(store)
	if err != nil {
		return nil, err
	}

	imageURL, err := util.ImageURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	p := v.imageDirPath(storeName, ID)
	info, err := v.s.Datastore.Stat(ctx, p)
	if err != nil {
		return nil, err
	}

	_, ok := info.(*types.FolderFileInfo)
	if !ok {
		return nil, fmt.Errorf("Stat error:  image doesn't exist (%s)", p)
	}

	meta, err := v.getMeta(ctx, storeName, ID)
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

	storeName, err := util.StoreName(store)
	if err != nil {
		return nil, err
	}

	res, err := lsDir(ctx, v.s.Datastore, v.datastorePath(v.imageStorePath(storeName)))
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

// Write the opaque metadata blobs (by name) for an image.  We create a
// directory under the image's parent directory.  Each blob in the metadata map
// is written to a file with the corresponding name.  Likewise, when we read it
// back (on restart) we populate the map accordingly.
func (v *ImageStore) writeMeta(ctx context.Context, storeName string, ID string,
	meta map[string][]byte) error {
	// XXX this should be done via disklib so this meta follows the disk in
	// case of motion.

	metaDataDir := v.imageMetadataDirPath(storeName, ID)

	if meta != nil && len(meta) != 0 {
		for name, value := range meta {
			r := bytes.NewReader(value)
			pth := path.Join(metaDataDir, name)
			log.Infof("Writing metadata %s", pth)
			if err := v.s.Datastore.Upload(ctx, r, pth, &soap.DefaultUpload); err != nil {
				return err
			}
		}
	} else {
		if err := v.fm.MakeDirectory(ctx, v.datastorePath(metaDataDir), v.s.Datacenter, false); err != nil {
			return err
		}
	}

	return nil
}

func (v *ImageStore) getMeta(ctx context.Context, storeName string, ID string) (map[string][]byte, error) {
	metaDataDir := v.datastorePath(v.imageMetadataDirPath(storeName, ID))
	res, err := lsDir(ctx, v.s.Datastore, metaDataDir)
	if err != nil {
		return nil, err
	}

	meta := make(map[string][]byte)
	for _, f := range res.File {
		finfo, ok := f.(*types.FileInfo)
		if !ok {
			continue
		}

		p := path.Join(v.imageMetadataDirPath(storeName, ID), finfo.Path)
		log.Infof("Getting meta for image (%s) %s", ID, p)
		rc, _, err := v.s.Datastore.Download(ctx, p, &soap.DefaultDownload)
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		buf, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, err
		}

		meta[finfo.Path] = buf
	}

	return meta, nil
}

// Create the top level directory the image storeas are created under
func (v *ImageStore) makeImageStoreParentDir(ctx context.Context) error {

	// check if it already exists
	res, err := lsDir(ctx, v.s.Datastore, v.s.Datastore.Path(""))
	if err != nil {
		return err
	}

	for _, f := range res.File {
		folder, ok := f.(*types.FileInfo)
		if !ok {
			continue
		}

		if folder.Path == datastoreParentPath {
			return nil
		}
	}

	log.Infof("Creating image store parent directory %s", datastoreParentPath)
	return v.fm.MakeDirectory(ctx, v.datastorePath(v.imageStorePath("")), v.s.Datacenter, false)
}

func lsDir(ctx context.Context, d *object.Datastore, p string) (*types.HostDatastoreBrowserSearchResults, error) {
	spec := types.HostDatastoreBrowserSearchSpec{
		MatchPattern: []string{"*"},
	}

	b, err := d.Browser(ctx)
	if err != nil {
		return nil, err
	}

	task, err := b.SearchDatastore(context.TODO(), p, &spec)
	if err != nil {
		return nil, err
	}

	info, err := task.WaitForResult(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	res := info.Result.(types.HostDatastoreBrowserSearchResults)
	return &res, nil
}
