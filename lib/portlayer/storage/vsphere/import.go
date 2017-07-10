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
	"io"
	"net/url"
	"os"

	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/lib/portlayer/storage"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/disk"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

func (v *VolumeStore) Import(op trace.Operation, id string, spec *archive.FilterSpec, tarstream io.ReadCloser) error {
	l, err := v.NewDataSink(op, id)
	if err != nil {
		return err
	}

	return l.Import(op, spec, tarstream)
}

// NewDataSource creates and returns an DataSource associated with container storage
func (v *VolumeStore) NewDataSink(op trace.Operation, id string) (storage.DataSink, error) {
	uri, err := v.URL(op, id)
	if err != nil {
		return nil, err
	}

	source, err := v.newDataSink(op, uri)
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
		_, _ = v.newOnlineDataSink(op, o)
		// if a != nil && a.available() {
		// 	return a, nil
		// }
	}

	return nil, errors.New("Unavailable")
}

func (v *VolumeStore) newDataSink(op trace.Operation, url *url.URL) (storage.DataSink, error) {
	mountPath, cleanFunc, err := v.Mount(op, url, true)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(mountPath)
	if err != nil {
		return nil, err
	}

	return &storage.MountDataSink{
		Path:  f,
		Clean: cleanFunc,
	}, nil
}

func (v *VolumeStore) newOnlineDataSink(op trace.Operation, vm *vm.VirtualMachine) (storage.DataSink, error) {
	return nil, errors.New("online sink not yet supported - expecting this to be a common toolbox implementaiton")
}

func (i *ImageStore) Import(op trace.Operation, id string, spec *archive.FilterSpec, tarStream io.ReadCloser) error {
	l, err := i.NewDataSink(op, id)
	if err != nil {
		return err
	}

	return l.Import(op, spec, tarStream)
}

// NewDataSink creates and returns an DataSource associated with image storage
func (i *ImageStore) NewDataSink(op trace.Operation, id string) (storage.DataSink, error) {
	uri, err := i.URL(op, id)
	if err != nil {
		return nil, err
	}

	// there is no online fail over path for images
	// we should probably have a check in here as to whether the image is "sealed" and can no longer
	// be modified.
	return i.newDataSink(op, uri)
}

func (i *ImageStore) newDataSink(op trace.Operation, url *url.URL) (storage.DataSink, error) {
	mountPath, cleanFunc, err := i.Mount(op, url, true)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(mountPath)
	if err != nil {
		return nil, err
	}

	return &storage.MountDataSink{
		Path:  f,
		Clean: cleanFunc,
	}, nil
}
