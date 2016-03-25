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
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
	portlayer "github.com/vmware/vic/portlayer/storage"
	"github.com/vmware/vic/portlayer/util"
	"golang.org/x/net/context"
)

// All paths on the datastore for images are relative to <datastore>/VIC/
var datastoreParentPath = "VIC"

const (
	defaultDiskSize = 8388608
)

type ImageStore struct {
	dm *disk.Manager
	fm *object.FileManager

	// govmomi session
	s *session.Session
}

func NewImageStore(ctx context.Context, s *session.Session) (*ImageStore, error) {
	dm, err := disk.NewDiskManager(ctx, s)
	if err != nil {
		return nil, err
	}

	vis := &ImageStore{
		dm: dm,
		fm: object.NewFileManager(s.Vim25()),
		s:  s,
	}

	// Ignore the error if the parent path already exists
	err = vis.makeImageStoreParentDir(ctx)
	if err != nil {
		return nil, err
	}

	return vis, nil
}

// Returns the URI in the datastore for a given path
func (v *ImageStore) imageStoreDatastoreURI(storeName string) string {
	return v.s.Datastore.Path(path.Join(datastoreParentPath, storeName))
}

// Returns the URI in the datastore for a given image relative to the given
// store.  The dir structure for an image in the datastore is
// `/VIC/imageStoreName/imageName/imageName.vmkd`
func (v *ImageStore) imageDirDatastoreURI(storeName, imageName string) string {
	return v.s.Datastore.Path(path.Join(datastoreParentPath, storeName, imageName))
}

// Uri to the vmdk itself
func (v *ImageStore) imageDiskDatastoreURI(storeName, imageName string) string {
	return path.Join(v.imageDirDatastoreURI(storeName, imageName), imageName+".vmdk")
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
	imagestore := v.imageStoreDatastoreURI(storeName)

	log.Infof("Creating imagestore directory %s", imagestore)
	if err = v.fm.MakeDirectory(ctx, imagestore, nil, false); err != nil {
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
	// URI format.
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
	res, err := lsDir(ctx, v.s.Datastore, v.imageStoreDatastoreURI(""))
	if err != nil {
		return nil, err
	}

	images := []*url.URL{}
	for _, f := range res.File {
		folder, ok := f.(*types.FolderFileInfo)
		if !ok {
			continue
		}
		u, err := util.StoreNameToURL(folder.Path)
		if err != nil {
			return nil, err
		}
		images = append(images, u)

	}

	return images, nil
}

// WriteImage creates a new image layer from the given parent.
// Eg parentImage + newLayer = new Image built from parent
//
// parent - The parent image to create the new image from.
// ID - textual ID for the image to be written
// Tag - the tag of the image to be written
func (v *ImageStore) WriteImage(ctx context.Context, parent *portlayer.Image, ID string, r io.Reader) (*portlayer.Image, error) {

	storeName, err := util.StoreName(parent.Store)
	if err != nil {
		return nil, err
	}

	imageURL, err := util.ImageURL(storeName, ID)
	if err != nil {
		return nil, err
	}

	// Create the image directory in the store.
	imageDirDsURI := v.imageDirDatastoreURI(storeName, ID)
	if err = v.fm.MakeDirectory(ctx, imageDirDsURI, nil, false); err != nil {
		return nil, err
	}

	ImageDiskDsURI := v.imageDiskDatastoreURI(storeName, ID)
	log.Infof("Creating image %s", ID)

	// If this is scratch, then it's the root of the image store.  All images
	// will be descended from this created and prepared fs.
	if ID == portlayer.Scratch.ID {
		// Create the disk
		vmdisk, err := v.dm.CreateAndAttach(ctx, ImageDiskDsURI, "", defaultDiskSize, os.O_RDWR)
		if err != nil {
			return nil, err
		}
		defer v.dm.Detach(ctx, vmdisk)

		// Make the filesystem
		if err = vmdisk.Mkfs(ID); err != nil {
			return nil, err
		}
	} else {

		if parent.ID == "" {
			return nil, fmt.Errorf("parent ID is empty")
		}

		// Create the disk
		parentDiskDsURI := v.imageDiskDatastoreURI(storeName, parent.ID)
		vmdisk, err := v.dm.CreateAndAttach(ctx, ImageDiskDsURI, parentDiskDsURI, 0, os.O_RDWR)
		if err != nil {
			return nil, err
		}
		defer v.dm.Detach(ctx, vmdisk)

		dir, err := ioutil.TempDir("", "mnt-"+ID)
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(dir)

		if err := vmdisk.Mount(dir, nil); err != nil {
			return nil, err
		}
		defer vmdisk.Unmount()

		// Untar the archive
		err = archive.Untar(r, dir, &archive.TarOptions{})
		if err != nil {
			return nil, err
		}
	}

	newImage := &portlayer.Image{
		ID:       ID,
		SelfLink: imageURL,
		Parent:   parent.SelfLink,
		Store:    parent.Store,
	}

	return newImage, nil
}

func (v *ImageStore) GetImage(ctx context.Context, store *url.URL, ID string) (*portlayer.Image, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (v *ImageStore) ListImages(ctx context.Context, store *url.URL, IDs []string) ([]*portlayer.Image, error) {
	return nil, fmt.Errorf("not yet implemented")
}

// Create the top level directory the image storeas are created under
func (v *ImageStore) makeImageStoreParentDir(ctx context.Context) error {

	// check if it already exists
	res, err := lsDir(ctx, v.s.Datastore, v.s.Datastore.Path(""))
	if err != nil {
		return err
	}

	for _, f := range res.File {
		folder, ok := f.(*types.FolderFileInfo)
		if !ok {
			continue
		}

		if folder.Path == datastoreParentPath {
			return nil
		}
	}

	log.Infof("Creating image store parent directory %s", datastoreParentPath)
	return v.fm.MakeDirectory(ctx, v.imageStoreDatastoreURI(""), nil, true)
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
