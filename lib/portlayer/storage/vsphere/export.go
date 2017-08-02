// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/lib/guest"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// Export reads the delta between child and parent volume layers, returning
// the difference as a tar archive.
//
// store - the volume store containing the two layers
// id - must inherit from ancestor if ancestor is specified
// ancestor - the volume layer up the chain against which to diff
// spec - describes filters on paths found in the data (include, exclude, strip)
// data - set to true to include file data in the tar archive, false to include headers only
func (v *VolumeStore) Export(op trace.Operation, id, ancestor string, spec *archive.FilterSpec, data bool) (io.ReadCloser, error) {
	if ancestor != "" {
		return nil, fmt.Errorf("volume diff is not supported in this volume store: %s", v.SelfLink.String())
	}

	l, err := v.NewDataSource(op, id)
	if err != nil {
		return nil, err
	}

	return l.Export(op, spec, data)
}

// NewDataSource creates and returns an DataSource associated with container storage
func (v *VolumeStore) NewDataSource(op trace.Operation, id string) (storage.DataSource, error) {
	uri, err := v.URL(op, id)
	if err != nil {
		return nil, err
	}

	offlineAttempt := 0
offline:
	offlineAttempt++

	// offline disk attempt
	source, err := v.newDataSource(op, uri)
	if err == nil {
		return source, err
	}

	// check for vmdk locked error here
	if !disk.IsLockedError(err) {
		op.Warnf("Unable to mount %s and do not know how to recover from error")
		// continue anyway because maybe there's an online option
	}

	// online - Owners() should filter out the appliance VM
	owners, _ := v.Owners(op, uri, disk.LockedVMDKFilter)
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
			self, _ := guest.IsSelf(op, uuid)
			if self && offlineAttempt < 2 {
				op.Infof("Appliance is owner of online vmdk - retrying offline source path")
				goto offline
			}
		}

		online, err := v.newOnlineDataSource(op, o, id)
		if online != nil {
			return online, err
		}

		op.Debugf("Failed to create online datasource with owner %s: %s", o.Reference(), err)
	}

	return nil, errors.New("unable to create online or offline data source")
}

func (v *VolumeStore) newDataSource(op trace.Operation, url *url.URL) (storage.DataSource, error) {
	mountPath, cleanFunc, err := v.Mount(op, url, false)
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

func (v *VolumeStore) newOnlineDataSource(op trace.Operation, owner *vm.VirtualMachine, id string) (storage.DataSource, error) {
	op.Debugf("Constructing toolbox data source: %s.%s", owner.Reference(), id)

	return &ToolboxDataSource{
		VM:    owner,
		ID:    storage.Label(id),
		Clean: func() { return },
	}, nil
}

// Export reads the delta between child and parent image layers, returning
// the difference as a tar archive.
//
// id - must inherit from ancestor if ancestor is specified
// ancestor - the layer up the chain against which to diff
// spec - describes filters on paths found in the data (include, exclude, rebase, strip)
// data - set to true to include file data in the tar archive, false to include headers only
func (i *ImageStore) Export(op trace.Operation, id, ancestor string, spec *archive.FilterSpec, data bool) (io.ReadCloser, error) {
	l, err := i.NewDataSource(op, id)
	if err != nil {
		return nil, err
	}

	if ancestor == "" {
		return l.Export(op, spec, data)
	}

	// for now we assume ancestor instead of entirely generic left/right
	// this allows us to assume it's an image
	r, err := i.NewDataSource(op, ancestor)
	if err != nil {
		op.Debugf("Unable to get datasource for ancestor: %s", err)

		l.Close()
		return nil, err
	}

	closers := func() error {
		op.Debugf("Callback to io.Closer function for image delta export")

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
		return nil, fmt.Errorf("mismatched datasource types: %T, %T", ls, rs)
	}

	// if we want data, exclude the xattrs, otherwise assume diff
	tar, err := archive.Diff(op, fl.Name(), fr.Name(), spec, data, !data)
	if err != nil {
		go closers()
		return nil, err
	}

	return &storage.ProxyReadCloser{
		ReadCloser: tar,
		Closer:     closers,
	}, nil
}

func (i *ImageStore) NewDataSource(op trace.Operation, id string) (storage.DataSource, error) {
	url, err := i.URL(op, id)
	if err != nil {
		return nil, err
	}

	return i.newDataSource(op, url)
}

func (i *ImageStore) newDataSource(op trace.Operation, url *url.URL) (storage.DataSource, error) {
	mountPath, cleanFunc, err := i.Mount(op, url, false)
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
