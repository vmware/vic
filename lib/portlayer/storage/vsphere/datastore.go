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
	"io"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"golang.org/x/net/context"
)

// Datastore gives access to the datastore regardless of type (esx, esx + vc,
// or esx + vc + vsan).  Also wraps paths to a given root directory
type datastore struct {
	// The Datastore API likes everything in "path/to/thing" format.
	ds *object.Datastore

	s *session.Session

	// The FileManager API likes everything in "[dsname] path/to/thing" format.
	fm *object.FileManager

	// the name of the datastore path (in case of vsan, this is the human readable name)
	rootname string

	// The datastore url (including root)
	rooturl string

	// datastore path (not in url form)
	rootdir string
}

// newDatastore returns a Datastore.
// ctx is a context,
// s is an authenticated session
// ds is the vsphere datastore
// rootdir is the top level directory to root all data.  If root does not exist,
// it will be created.  If it already exists, NOOP. This cannot be empty.
func newDatastore(ctx context.Context, s *session.Session, ds *object.Datastore, rootdir string) (*datastore, error) {

	d := &datastore{
		ds:       ds,
		s:        s,
		rootname: rootdir,
		fm:       object.NewFileManager(s.Vim25()),
	}

	if strings.HasPrefix(rootdir, "/") {
		rootdir = strings.TrimPrefix(rootdir, "/")
	}

	// Get the root directory element split from the rest of the path (if there is one)
	root := strings.SplitN(rootdir, "/", 2)

	// Create the first element.  This handles vsan vmfs top level dirs.
	if err := d.mkRootDir(ctx, root[0]); err != nil {
		log.Infof("error creating root directory %s: %s", rootdir, err)
		return nil, err
	}

	// Create the rest conventionally
	if len(root) > 1 {
		r, err := d.Mkdir(ctx, true, root[1])
		if err != nil {
			return nil, err
		}
		d.rooturl = r
	}

	log.Infof("Datastore path is %s", d.rooturl)
	return d, nil
}

func (d *datastore) Summary(ctx context.Context) (*types.DatastoreSummary, error) {

	var mds mo.Datastore
	if err := d.ds.Properties(ctx, d.ds.Reference(), []string{"info", "summary"}, &mds); err != nil {
		return nil, err
	}

	return &mds.Summary, nil
}

// Mkdir creates directories.
func (d *datastore) Mkdir(ctx context.Context, createParentDirectories bool, dirs ...string) (string, error) {

	upth := path.Join(dirs...)

	// check if it already exists
	_, err := d.Ls(ctx, upth)
	if err != nil && !types.IsFileNotFound(err) {
		return "", err
	}

	upth = path.Join(d.rooturl, upth)

	// dir already exists
	if err == nil {
		return upth, nil
	}

	log.Infof("Creating directory %s", upth)

	if err := d.fm.MakeDirectory(ctx, upth, d.s.Datacenter, createParentDirectories); err != nil {
		return "", err
	}

	return upth, nil
}

// Ls returns a list of dirents at the given path (relative to root)
//
// A note aboutpaths and the datastore browser.
// None of these work paths work
// r, err := ds.Ls(ctx, "ds:///vmfs/volumes/vsan:52a67632ac3497a3-411916fd50bedc27/0ea65357-0494-d42d-2ede-000c292dc5b5")
// r, err := ds.Ls(ctx, "[vsanDatastore] ds:///vmfs/volumes/vsan:52a67632ac3497a3-411916fd50bedc27/")
// r, err := ds.Ls(ctx, "[vsanDatastore] //vmfs/volumes/vsan:52a67632ac3497a3-411916fd50bedc27/")
// r, err := ds.Ls(ctx, "[] ds:///vmfs/volumes/vsan:52a67632ac3497a3-411916fd50bedc27/0ea65357-0494-d42d-2ede-000c292dc5b5")
// r, err := ds.Ls(ctx, "[] /vmfs/volumes/vsan:52a67632ac3497a3-411916fd50bedc27/0ea65357-0494-d42d-2ede-000c292dc5b5")
// r, err := ds.Ls(ctx, "[] ../vmfs/volumes/vsan:52a67632ac3497a3-411916fd50bedc27/0ea65357-0494-d42d-2ede-000c292dc5b5")
// r, err := ds.Ls(ctx, "[] ./vmfs/volumes/vsan:52a67632ac3497a3-411916fd50bedc27/0ea65357-0494-d42d-2ede-000c292dc5b5")
// r, err := ds.Ls(ctx, "[52a67632ac3497a3-411916fd50bedc27] /0ea65357-0494-d42d-2ede-000c292dc5b5")
// r, err := ds.Ls(ctx, "[vsan:52a67632ac3497a3-411916fd50bedc27] /0ea65357-0494-d42d-2ede-000c292dc5b5")
// r, err := ds.Ls(ctx, "[vsan:52a67632ac3497a3-411916fd50bedc27] 0ea65357-0494-d42d-2ede-000c292dc5b5")
// r, err := ds.Ls(ctx, "[vsanDatastore] /vmfs/volumes/vsan:52a67632ac3497a3-411916fd50bedc27/0ea65357-0494-d42d-2ede-000c292dc5b5")

// The only URI that works on VC + VSAN.
// r, err := ds.Ls(ctx, "[vsanDatastore] /0ea65357-0494-d42d-2ede-000c292dc5b5")
//
func (d *datastore) Ls(ctx context.Context, p string) (*types.HostDatastoreBrowserSearchResults, error) {
	spec := types.HostDatastoreBrowserSearchSpec{
		MatchPattern: []string{"*"},
	}

	b, err := d.ds.Browser(ctx)
	if err != nil {
		return nil, err
	}

	task, err := b.SearchDatastore(ctx, path.Join(d.rooturl, p), &spec)
	if err != nil {
		return nil, err
	}

	info, err := task.WaitForResult(ctx, nil)
	if err != nil {
		return nil, err
	}

	res := info.Result.(types.HostDatastoreBrowserSearchResults)
	return &res, nil
}

// LsDirs returns a list of dirents at the given path (relative to root)
func (d *datastore) LsDirs(ctx context.Context, p string) (*types.ArrayOfHostDatastoreBrowserSearchResults, error) {
	spec := types.HostDatastoreBrowserSearchSpec{
		MatchPattern: []string{"*"},
	}

	b, err := d.ds.Browser(ctx)
	if err != nil {
		return nil, err
	}

	task, err := b.SearchDatastoreSubFolders(ctx, path.Join(d.rooturl, p), &spec)
	if err != nil {
		return nil, err
	}

	info, err := task.WaitForResult(ctx, nil)
	if err != nil {
		return nil, err
	}

	res := info.Result.(types.ArrayOfHostDatastoreBrowserSearchResults)
	return &res, nil
}

func (d *datastore) Upload(ctx context.Context, r io.Reader, pth string) error {
	return d.ds.Upload(ctx, r, path.Join(d.rootdir, pth), &soap.DefaultUpload)
}

func (d *datastore) Download(ctx context.Context, pth string) (io.ReadCloser, error) {
	rc, _, err := d.ds.Download(ctx, path.Join(d.rootdir, pth), &soap.DefaultDownload)
	return rc, err
}

func (d *datastore) Stat(ctx context.Context, pth string) (types.BaseFileInfo, error) {
	return d.ds.Stat(ctx, path.Join(d.rootdir, pth))
}

func (d *datastore) Mv(ctx context.Context, fromPath, toPath string) error {
	from := path.Join(d.rooturl, fromPath)
	to := path.Join(d.rooturl, toPath)
	err := tasks.Wait(ctx, func(context.Context) (tasks.Waiter, error) {
		return d.fm.MoveDatastoreFile(ctx, from, d.s.Datacenter, to, d.s.Datacenter, true)
	})

	return err
}

func (d *datastore) IsVSAN(ctx context.Context) bool {
	dsType, _ := d.ds.Type(ctx)
	return dsType == types.HostFileSystemVolumeFileSystemTypeVsan
}

// This creates the root directory in the datastore and sets the rooturl and
// rootdir in the datastore struct so we can reuse it for other routines.  This
// handles vsan + vc, vsan + esx, and esx.  The URI conventions are not the
// same for each and this tries to create the directory and stash the relevant
// result so the URI doesn't need to be recomputed for every datastore
// operation.
func (d *datastore) mkRootDir(ctx context.Context, rootdir string) error {

	// Handle vsan
	// Vsan will complain if the root dir exists.  Just call it directly and
	// swallow the error if it's already there.
	if d.IsVSAN(ctx) {
		nm := object.NewDatastoreNamespaceManager(d.s.Vim25())

		// This returns the vmfs path (including the datastore and directory
		// UUIDs).  Use the directory UUID in future operations because it is
		// the stable path which we can use regardless of vsan state.
		uuid, err := nm.CreateDirectory(ctx, d.ds, rootdir, "")
		if err != nil {
			if !soap.IsSoapFault(err) {
				return err
			}

			soapFault := soap.ToSoapFault(err)
			_, ok := soapFault.VimFault().(types.FileAlreadyExists)
			if ok {
				// XXX UGLY HACK until we move this into the installer.  Use the
				// display name if the dir exists since we can't get the UUID after the
				// directory is created.

				uuid = rootdir
				err = nil
			} else {
				return err
			}
		}

		// set the root url to the UUID of the dir we created
		d.rootdir = path.Base(uuid)
		d.rooturl = d.ds.Path(d.rootdir)
		log.Infof("Created store parent directory (%s) at %s", rootdir, d.rooturl)
		return nil
	}

	// Handle regular local datastore
	// check if it already exists

	d.rootdir = rootdir
	d.rooturl = d.ds.Path(rootdir)

	if _, err := d.Mkdir(ctx, true); err != nil {
		return err
	}

	return nil
}
