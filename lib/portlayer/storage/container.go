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

package storage

import (
	"errors"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// ContainerStorer defines the interface contract expected to allow import and export
// against containers
type ContainerStorer interface {
	Resolver
	Importer
	Exporter
}

// ContainerStore stores container storage information
type ContainerStore struct {
	disk.Vmdk

	// used to resolve images when diffing
	images Resolver
}

// NewContainerStore creates and returns a new container store
func NewContainerStore(op trace.Operation, s *session.Session, imageResolver Resolver) (*ContainerStore, error) {
	dm, err := disk.NewDiskManager(op, s, Config.ContainerView)
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

	dsPath, _ := datastore.PathFromString(url.Path)
	config := disk.NewPersistentDisk(dsPath)

	return c.InUse(op, config, disk.LockedVMDKFilter)
}

// NewDataSource creates and returns an DataSource associated with container storage
func (c *ContainerStore) NewDataSource(op trace.Operation, id string) (DataSource, error) {
	uri, err := c.URL(op, id)
	if err != nil {
		return nil, err
	}

	source, err := c.newDataSource(op, uri)
	if err == nil {
		return source, err
	}

	// check for vmdk locked error here
	if !disk.IsLockedError(err) {
		op.Warnf("Unable to mount %s and do not know how to recover from error")
		// continue anyway because maybe there's an online option
	}

	// online - Owners() should filter out the appliance VM
	owners, _ := c.Owners(op, uri, disk.LockedVMDKFilter)
	if len(owners) == 0 {
		return nil, errors.New("Unavailable")
	}

	// TODO(jzt): tweak this when online export is available
	for _, o := range owners {
		// o is a VM
		_, _ = c.newOnlineDataSource(op, o)
		// if a != nil && a.available() {
		// 	return a, nil
		// }
	}

	return nil, errors.New("Unavailable")
}

func (c *ContainerStore) newDataSource(op trace.Operation, url *url.URL) (DataSource, error) {
	mountPath, cleanFunc, err := c.Mount(op, url, false)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(mountPath)
	if err != nil {
		return nil, err
	}

	return &MountDataSource{
		Path:  f,
		Clean: cleanFunc,
	}, nil
}

func (c *ContainerStore) newOnlineDataSource(op trace.Operation, vm *vm.VirtualMachine) (DataSource, error) {
	return nil, errors.New("online source not yet supported - expecting this to be a common toolbox implementaiton")
}

// NewDataSink creates and returns an DataSink associated with container storage
func (c *ContainerStore) NewDataSink(op trace.Operation, id string) (DataSink, error) {
	uri, err := c.URL(op, id)
	if err != nil {
		return nil, err
	}

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
	owners, _ := c.Owners(op, uri, disk.LockedVMDKFilter)
	if len(owners) == 0 {
		return nil, errors.New("Unavailable")
	}

	// TODO(jzt): tweak this when online export is available
	for _, o := range owners {
		// o is a VM
		_, _ = c.newOnlineDataSink(op, o)
		// if a != nil && a.available() {
		// 	return a, nil
		// }
	}

	return nil, errors.New("Unavailable")
}

func (c *ContainerStore) newDataSink(op trace.Operation, url *url.URL) (DataSink, error) {
	mountPath, cleanFunc, err := c.Mount(op, url, true)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(mountPath)
	if err != nil {
		return nil, err
	}

	return &MountDataSink{
		Path:  f,
		Clean: cleanFunc,
	}, nil
}

func (c *ContainerStore) newOnlineDataSink(op trace.Operation, owner *vm.VirtualMachine) (DataSink, error) {
	return nil, errors.New("online sink not yet supported - expecting this to be a common toolbox implementaiton")
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
		return l.Export(op, spec, data)
	}

	// for now we assume ancetor instead of entirely generic left/right
	// this allows us to assume it's an image
	img, err := c.images.URL(op, ancestor)
	if err != nil {
		return nil, err
	}

	r, err := c.newDataSource(op, img)
	if err != nil {
		l.Close()
		return nil, err
	}

	closers := func() {
		l.Close()
		r.Close()
	}

	ls := l.Source()
	rs := r.Source()

	fl, lok := ls.(*os.File)
	fr, rok := rs.(*os.File)

	if !lok || !rok {
		go closers()
		return nil, errors.New("Mismatched datasource types")
	}

	tar, err := archive.Diff(op, fl.Name(), fr.Name(), spec, data)
	if err != nil {
		go closers()
		return nil, err
	}

	return &CleanupReader{
		ReadCloser: tar,
		Clean:      closers,
	}, nil
}
