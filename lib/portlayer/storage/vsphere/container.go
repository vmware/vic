// Copyright 2017 VMware, Inc. All Rights Reserved.
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
	"errors"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// ContainerStorer defines the interface contract expected to allow import and export
// against containers
type ContainerStorer interface {
	storage.Resolver
	storage.Importer
	storage.Exporter
}

// ContainerStore stores container storage information
type ContainerStore struct {
	disk.Vmdk

	// used to resolve images when diffing
	images storage.Resolver
}

// NewContainerStore creates and returns a new container store
func NewContainerStore(op trace.Operation, s *session.Session, imageResolver storage.Resolver) (*ContainerStore, error) {
	dm, err := disk.NewDiskManager(op, s, storage.Config.ContainerView)
	if err != nil {
		return nil, err
	}

	cs := &ContainerStore{
		Vmdk: disk.Vmdk{
			Manager: dm,
			//ds: ds,
			Session: s,
		},

		images: imageResolver,
	}
	return cs, nil
}

// URL converts the id of a resource to a URL
func (c *ContainerStore) URL(op trace.Operation, id string) (*url.URL, error) {
	// using diskfinder with a basic suffix match is an inefficient and potentially error prone way of doing this
	// mapping, but until the container store has a structured means of knowing this information it's at least
	// not going to be incorrect without an ID collision.
	dsPath, err := c.DiskFinder(op, func(filename string) bool {
		return strings.HasSuffix(filename, id+".vmdk")
	})
	if err != nil {
		return nil, err
	}

	return &url.URL{
		Scheme: "ds",
		Path:   dsPath,
	}, nil
}

// Owners returns a list of VMs that are using the resource specified by `url`
func (c *ContainerStore) Owners(op trace.Operation, url *url.URL, filter func(vm *mo.VirtualMachine) bool) ([]*vm.VirtualMachine, error) {
	if url.Scheme != "ds" {
		return nil, errors.New("vmdk path must be a datastore url with \"ds\" scheme")
	}

	return c.Vmdk.Owners(op, url, filter)
}

// NewDataSource creates and returns an DataSource associated with container storage
func (c *ContainerStore) NewDataSource(op trace.Operation, id string) (storage.DataSource, error) {
	uri, err := c.URL(op, id)
	if err != nil {
		return nil, err
	}

	offlineAttempt := 0
offline:
	offlineAttempt++

	// This is persistent to avoid issues with concurrent Stat/Import calls
	source, err := c.newDataSource(op, uri, true)
	if err == nil {
		return source, err
	}

	// check for vmdk locked error here
	if !disk.IsLockedError(err) {
		op.Warnf("Unable to mount %s and do not know how to recover from error")
		// continue anyway because maybe there's an online option
	}

	// online - Owners() should filter out the appliance VM
	// #nosec: Errors unhandled.
	owners, _ := c.Owners(op, uri, disk.LockedVMDKFilter)
	if len(owners) == 0 {
		op.Infof("No online owners were found for %s", id)
		return nil, errors.New("unable to create offline data source and no online owners found")
	}

	for _, o := range owners {
		// sanity check to see if we are the owner - this should catch transitions
		// from container running to diff or commit for example between the offline attempt and here
		uuid, err := o.UUID(op)
		if err == nil {
			// check if the vm is appliance VM if we can successfully get its UUID
			// #nosec: Errors unhandled.
			self, _ := guest.IsSelf(op, uuid)
			if self && offlineAttempt < 2 {
				op.Infof("Appliance is owner of online vmdk - retrying offline source path")
				goto offline
			}
		}

		online, err := c.newOnlineDataSource(op, o, id)
		if online != nil {
			return online, err
		}

		op.Debugf("Failed to create online datasource with owner %s: %s", o.Reference(), err)
	}

	return nil, errors.New("unable to create online or offline data source")
}

func (c *ContainerStore) newDataSource(op trace.Operation, url *url.URL, persistent bool) (storage.DataSource, error) {
	mountPath, cleanFunc, err := c.Mount(op, url, persistent)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(mountPath)
	if err != nil {
		cleanFunc()
		return nil, err
	}

	op.Debugf("Created mount data source for access to %s at %s", url, mountPath)
	return storage.NewMountDataSource(op, f, cleanFunc), nil
}

func (c *ContainerStore) newOnlineDataSource(op trace.Operation, owner *vm.VirtualMachine, id string) (storage.DataSource, error) {
	op.Debugf("Constructing toolbox data source: %s.%s", owner.Reference(), id)

	return &ToolboxDataSource{
		VM:    owner,
		ID:    id,
		Clean: func() { return },
	}, nil
}

// NewDataSink creates and returns an DataSink associated with container storage
func (c *ContainerStore) NewDataSink(op trace.Operation, id string) (storage.DataSink, error) {
	uri, err := c.URL(op, id)
	if err != nil {
		return nil, err
	}

	offlineAttempt := 0
offline:
	offlineAttempt++

	sink, err := c.newDataSink(op, uri)
	if err == nil {
		return sink, err
	}

	// check for vmdk locked error here
	if !disk.IsLockedError(err) {
		op.Warnf("Unable to mount %s and do not know how to recover from error")
		// continue anyway because maybe there's an online option
	}

	// online - Owners() should filter out the appliance VM
	// #nosec: Errors unhandled.
	owners, _ := c.Owners(op, uri, disk.LockedVMDKFilter)
	if len(owners) == 0 {
		op.Infof("No online owners were found for %s", id)
		return nil, errors.New("unable to create offline data sink and no online owners found")
	}

	for _, o := range owners {
		// sanity check to see if we are the owner - this should catch transitions
		// from container running to diff or commit for example between the offline attempt and here
		uuid, err := o.UUID(op)
		if err == nil {
			// check if the vm is appliance VM if we can successfully get its UUID
			// #nosec: Errors unhandled.
			self, _ := guest.IsSelf(op, uuid)
			if self && offlineAttempt < 2 {
				op.Infof("Appliance is owner of online vmdk - retrying offline source path")
				goto offline
			}
		}

		online, err := c.newOnlineDataSink(op, o, id)
		if online != nil {
			return online, err
		}

		op.Debugf("Failed to create online datasink with owner %s: %s", o.Reference(), err)
	}

	return nil, errors.New("unable to create online or offline data sink")
}

func (c *ContainerStore) newDataSink(op trace.Operation, url *url.URL) (storage.DataSink, error) {
	mountPath, cleanFunc, err := c.Mount(op, url, true)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(mountPath)
	if err != nil {
		cleanFunc()
		return nil, err
	}

	op.Debugf("Created mount data sink for access to %s at %s", url, mountPath)
	return storage.NewMountDataSink(op, f, cleanFunc), nil
}

func (c *ContainerStore) newOnlineDataSink(op trace.Operation, owner *vm.VirtualMachine, id string) (storage.DataSink, error) {
	op.Debugf("Constructing toolbox data sink: %s.%s", owner.Reference(), id)

	return &ToolboxDataSink{
		VM:    owner,
		ID:    id,
		Clean: func() { return },
	}, nil
}

func (c *ContainerStore) Import(op trace.Operation, id string, spec *archive.FilterSpec, tarstream io.ReadCloser) error {
	l, err := c.NewDataSink(op, id)
	if err != nil {
		return err
	}

	return l.Import(op, spec, tarstream)
}

func (c *ContainerStore) Export(op trace.Operation, id, ancestor string, spec *archive.FilterSpec, data bool) (io.ReadCloser, error) {
	l, err := c.NewDataSource(op, id)
	if err != nil {
		return nil, err
	}

	if ancestor == "" {
		op.Infof("No ancestor specified so following basic export path")
		return l.Export(op, spec, data)
	}

	// for now we assume ancetor instead of entirely generic left/right
	// this allows us to assume it's an image
	img, err := c.images.URL(op, ancestor)
	if err != nil {
		op.Errorf("Failed to map ancestor %s to image: %s", ancestor, err)

		l.Close()
		return nil, err
	}
	op.Debugf("Mapped ancestor %s to %s", ancestor, img.String())

	r, err := c.newDataSource(op, img, false)
	if err != nil {
		op.Debugf("Unable to get datasource for ancestor: %s", err)

		l.Close()
		return nil, err
	}

	closers := func() error {
		op.Debugf("Callback to io.Closer function for container export")

		l.Close()
		r.Close()

		return nil
	}

	ls := l.Source()
	rs := r.Source()

	fl, lok := ls.(*os.File)
	fr, rok := rs.(*os.File)

	if !lok || !rok {
		go closers()
		return nil, errors.New("mismatched datasource types")
	}

	// if we want data, exclude the xattrs, otherwise assume diff
	xattrs := !data

	tar, err := archive.Diff(op, fl.Name(), fr.Name(), spec, data, xattrs)
	if err != nil {
		go closers()
		return nil, err
	}

	return &storage.ProxyReadCloser{
		ReadCloser: tar,
		Closer:     closers,
	}, nil
}
