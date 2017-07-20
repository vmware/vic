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

package disk

import (
	"net/url"

	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
	"github.com/vmware/vic/pkg/vsphere/session"
)

// Vmdk is intended to be embedded by stores that manage VMDK-based data resources
type Vmdk struct {
	*Manager
	*datastore.Helper
	*session.Session
}

const (
	DiskBackendKey = "msg.disk.noBackEnd"
	LockedFileKey  = "msg.fileio.lock"
)

// Mount mounts the disk, returning the mount path and the function used to unmount/detaches
// when no longer in use
func (v *Vmdk) Mount(op trace.Operation, uri *url.URL, persistent bool) (string, func(), error) {
	if uri.Scheme != "ds" {
		return "", nil, errors.New("vmdk path must be a datastore url with \"ds\" scheme")
	}

	dsPath, err := datastore.PathFromString(uri.Path)
	if err != nil {
		return "", nil, err
	}

	cleanFunc := func() {
		if err := v.UnmountAndDetach(op, dsPath, persistent); err != nil {
			op.Errorf("Error cleaning up disk: %s", err.Error())
		}
	}

	mountPath, err := v.AttachAndMount(op, dsPath, persistent)
	return mountPath, cleanFunc, err
}

func LockedVMDKFilter(vm *mo.VirtualMachine) bool {
	return vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn
}

func IsLockedError(op trace.Operation, err error) bool {
	switch err := err.(type) {
	case task.Error:
		fault := err.Fault().GetMethodFault()

		// indicates it's specific to a disk
		diskBackend := false
		// indicates it's specific to locking
		fileLock := false

		for i := range fault.FaultMessage {
			message := &fault.FaultMessage[i]
			switch message.Key {
			case DiskBackendKey:
				op.Debugf("diskbackend is true")
				diskBackend = true
			case LockedFileKey:
				op.Debugf("lockedfilekey is true")
				fileLock = true
			}

			if diskBackend && fileLock {
				return true
			}
		}
	default:
		op.Debugf("IsLockedError: this is not a task error")
	}

	return false
}

// LockedDisks returns locked devices path in the error if it's device lock error
func LockedDisks(err error) []string {
	var faultMessage []types.LocalizableMessage

	if soap.IsSoapFault(err) {
		switch f := soap.ToSoapFault(err).VimFault().(type) {
		case *types.GenericVmConfigFault:
			faultMessage = f.FaultMessage
		}
	} else if soap.IsVimFault(err) {
		faultMessage = soap.ToVimFault(err).GetMethodFault().FaultMessage
	} else {
		switch err := err.(type) {
		case task.Error:
			faultMessage = err.Fault().GetMethodFault().FaultMessage
		}
	}

	if faultMessage == nil {
		return nil
	}

	lockedFile := false
	var devices []string
	for _, message := range faultMessage {
		switch message.Key {
		case LockedFileKey:
			lockedFile = true
		case DiskBackendKey:
			for _, arg := range message.Arg {
				if device, ok := arg.Value.(string); ok {
					devices = append(devices, device)
					continue
				}
			}
		}
	}
	if lockedFile {
		// make sure locked devices are returned only when both keys appear in the error
		return devices
	}
	return nil
}
