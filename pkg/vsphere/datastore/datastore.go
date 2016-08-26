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

package datastore

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"golang.org/x/net/context"
)

// Helper gives access to the datastore regardless of type (esx, esx + vc,
// or esx + vc + vsan).  Also wraps paths to a given root directory
type Helper struct {
	// The Datastore API likes everything in "path/to/thing" format.
	ds *object.Datastore

	s *session.Session

	// The FileManager API likes everything in "[dsname] path/to/thing" format.
	fm *object.FileManager

	// The datastore url (including root) in "[dsname] /path" format.
	RootURL string
}

// NewDatastore returns a Datastore.
// ctx is a context,
// s is an authenticated session
// ds is the vsphere datastore
// rootdir is the top level directory to root all data.  If root does not exist,
// it will be created.  If it already exists, NOOP. This cannot be empty.
func NewHelper(ctx context.Context, s *session.Session, ds *object.Datastore, rootdir string) (*Helper, error) {

	d := &Helper{
		ds: ds,
		s:  s,
		fm: object.NewFileManager(s.Vim25()),
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
			if !os.IsExist(err) {
				return nil, err
			}
			log.Debugf("%s already exists", d.RootURL)
		}
		d.RootURL = r
	}

	log.Infof("Datastore path is %s", d.RootURL)
	return d, nil
}

// GetDatastores returns a map of datastores given a map of names and urls
func GetDatastores(ctx context.Context, s *session.Session, dsURLs map[string]url.URL) (map[string]*Helper, error) {
	stores := make(map[string]*Helper)

	fm := object.NewFileManager(s.Vim25())
	for name, dsURL := range dsURLs {

		vsDs, err := s.Finder.DatastoreOrDefault(ctx, s.DatastorePath)
		if err != nil {
			return nil, err
		}

		d := &Helper{
			ds:      vsDs,
			s:       s,
			fm:      fm,
			RootURL: dsURL.Path,
		}

		stores[name] = d
	}

	return stores, nil
}

func (d *Helper) Summary(ctx context.Context) (*types.DatastoreSummary, error) {

	var mds mo.Datastore
	if err := d.ds.Properties(ctx, d.ds.Reference(), []string{"info", "summary"}, &mds); err != nil {
		return nil, err
	}

	return &mds.Summary, nil
}

// Mkdir creates directories.
func (d *Helper) Mkdir(ctx context.Context, createParentDirectories bool, dirs ...string) (string, error) {

	upth := path.Join(dirs...)
	upth = path.Join(d.RootURL, upth)

	log.Infof("Creating directory %s", upth)

	if err := d.fm.MakeDirectory(ctx, upth, d.s.Datacenter, createParentDirectories); err != nil {

		log.Debugf("Creating %s error: %s", upth, err)

		if err != nil {
			if soap.IsSoapFault(err) {
				soapFault := soap.ToSoapFault(err)
				if _, ok := soapFault.VimFault().(types.FileAlreadyExists); ok {
					return "", os.ErrExist
				}
			}
		}

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
func (d *Helper) Ls(ctx context.Context, p string) (*types.HostDatastoreBrowserSearchResults, error) {
	spec := types.HostDatastoreBrowserSearchSpec{
		MatchPattern: []string{"*"},
	}

	b, err := d.ds.Browser(ctx)
	if err != nil {
		return nil, err
	}

	task, err := b.SearchDatastore(ctx, path.Join(d.RootURL, p), &spec)
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
func (d *Helper) LsDirs(ctx context.Context, p string) (*types.ArrayOfHostDatastoreBrowserSearchResults, error) {
	spec := &types.HostDatastoreBrowserSearchSpec{
		MatchPattern: []string{"*"},
	}

	b, err := d.ds.Browser(ctx)
	if err != nil {
		return nil, err
	}

	task, err := b.SearchDatastoreSubFolders(ctx, path.Join(d.RootURL, p), spec)
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

func (d *Helper) Upload(ctx context.Context, r io.Reader, pth string) error {
	return d.ds.Upload(ctx, r, path.Join(d.rootDir(), pth), &soap.DefaultUpload)
}

func (d *Helper) Download(ctx context.Context, pth string) (io.ReadCloser, error) {
	rc, _, err := d.ds.Download(ctx, path.Join(d.rootDir(), pth), &soap.DefaultDownload)
	return rc, err
}

func (d *Helper) Stat(ctx context.Context, pth string) (types.BaseFileInfo, error) {
	return d.ds.Stat(ctx, path.Join(d.rootDir(), pth))
}

func (d *Helper) Mv(ctx context.Context, fromPath, toPath string) error {
	from := path.Join(d.RootURL, fromPath)
	to := path.Join(d.RootURL, toPath)
	log.Infof("Moving %s to %s", from, to)
	err := tasks.Wait(ctx, func(context.Context) (tasks.ResultWaiter, error) {
		return d.fm.MoveDatastoreFile(ctx, from, d.s.Datacenter, to, d.s.Datacenter, true)
	})

	return err
}

func (d *Helper) Rm(ctx context.Context, pth string) error {
	f := path.Join(d.RootURL, pth)
	log.Infof("Removing %s", pth)
	return tasks.Wait(context.TODO(), func(ctx context.Context) (tasks.ResultWaiter, error) {
		return d.fm.DeleteDatastoreFile(ctx, f, d.s.Datacenter)
	})
}

func (d *Helper) IsVSAN(ctx context.Context) bool {
	dsType, _ := d.ds.Type(ctx)
	return dsType == types.HostFileSystemVolumeFileSystemTypeVsan
}

// This creates the root directory in the datastore and sets the rooturl and
// rootdir in the datastore struct so we can reuse it for other routines.  This
// handles vsan + vc, vsan + esx, and esx.  The URI conventions are not the
// same for each and this tries to create the directory and stash the relevant
// result so the URI doesn't need to be recomputed for every datastore
// operation.
func (d *Helper) mkRootDir(ctx context.Context, rootdir string) error {

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
		d.RootURL = d.ds.Path(path.Base(uuid))
		log.Infof("Created store parent directory (%s) at %s", rootdir, d.RootURL)
	} else {

		// Handle regular local datastore
		// check if it already exists

		d.RootURL = d.ds.Path(rootdir)
		if _, err := d.Mkdir(ctx, true); err != nil {
			if os.IsExist(err) {
				log.Debugf("%s already exists", d.RootURL)
				return nil
			}
			return err
		}
	}

	return nil
}

// Return the root of the datastore path (without the [datastore] portion)
func (d *Helper) rootDir() string {
	return strings.SplitN(d.RootURL, " ", 2)[1]
}

// Parse the datastore format ([datastore1] /path/to/thing) to groups.
var datastoreFormat = regexp.MustCompile(`^\[([\w\d\(\)-_\.\s]+)\]`)
var pathFormat = regexp.MustCompile(`\s([\/\w-_\.]*$)`)

// Converts `[datastore] /path` to ds:// URL
func ToURL(ds string) (*url.URL, error) {
	u := new(url.URL)
	var matches []string
	if matches = datastoreFormat.FindStringSubmatch(ds); len(matches) != 2 {
		return nil, fmt.Errorf("Ambiguous datastore hostname format encountered from input: %s.", ds)
	}
	u.Host = matches[1]
	if matches = pathFormat.FindStringSubmatch(ds); len(matches) != 2 {
		return nil, fmt.Errorf("Ambiguous datastore path format encountered from input: %s.", ds)
	}

	u.Path = path.Clean(matches[1])
	u.Scheme = "ds"

	return u, nil
}

// Converts ds:// URL for datastores to datastore format ([datastore1] /path/to/thing)
func URLtoDatastore(u *url.URL) (string, error) {
	scheme := "ds"
	if u.Scheme != scheme {
		return "", fmt.Errorf("url (%s) is not a datastore", u.String())
	}
	return fmt.Sprintf("[%s] %s", u.Host, u.Path), nil
}

// TestName builds a unique datastore name
func TestName(suffix string) string {
	return uuid.New().String()[0:16] + "-" + suffix
}
