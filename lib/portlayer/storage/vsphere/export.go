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
		return nil, errors.New("Unavailable")
	}

	// TODO(jzt): tweak this when online export is available
	for _, o := range owners {
		// o is a VM
		_, _ = v.newOnlineDataSource(op, o)
		// if a != nil && a.available() {
		// 	return a, nil
		// }
	}

	return nil, errors.New("Unavailable")
}

func (v *VolumeStore) newDataSource(op trace.Operation, url *url.URL) (storage.DataSource, error) {
	mountPath, cleanFunc, err := v.Mount(op, url, false)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(mountPath)
	if err != nil {
		return nil, err
	}

	return &storage.MountDataSource{
		Path:  f,
		Clean: cleanFunc,
	}, nil
}

func (v *VolumeStore) newOnlineDataSource(op trace.Operation, vm *vm.VirtualMachine) (storage.DataSource, error) {
	return nil, errors.New("online source not yet supported - expecting this to be a common toolbox implementaiton")
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

	// for now we assume ancetor instead of entirely generic left/right
	// this allows us to assume it's an image
	r, err := i.NewDataSource(op, ancestor)
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
		return nil, fmt.Errorf("mismatched datasource types: %T, %T", ls, rs)
	}

	tar, err := archive.Diff(op, fl.Name(), fr.Name(), spec, data)
	if err != nil {
		go closers()
		return nil, err
	}

	return &storage.CleanupReader{
		ReadCloser: tar,
		Clean:      closers,
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
		return nil, err
	}

	return &storage.MountDataSource{
		Path:  f,
		Clean: cleanFunc,
	}, nil
}
