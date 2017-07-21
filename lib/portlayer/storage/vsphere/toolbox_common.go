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
	"fmt"
	"path"

	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/guest/toolbox"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

const (
	DiskLabelQueryName  = "disk-label"
	FilterSpecQueryName = "filter-spec"
)

// Parse Archive does something.
func BuildArchiveURL(op trace.Operation, disklabel, target string, fs *archive.FilterSpec) (string, error) {
	encodedSpec, err := archive.EncodeFilterSpec(op, fs)
	if err != nil {
		return "", err
	}
	target = path.Join("/archive:/", target)

	// if diskLabel is longer than 16 characters, then the function was passed a containerID
	// use containerfs as the diskLabel
	if len(disklabel) > 16 {
		disklabel = "containerfs"
	}

	target += fmt.Sprintf("?%s=%s&%s=%s", DiskLabelQueryName, disklabel, FilterSpecQueryName, *encodedSpec)
	op.Debugf("OnlineData* Url: %s", target)
	return target, nil
}

// GetToolboxClient returns a toolbox client given a vm and id
func GetToolboxClient(op trace.Operation, vm *vm.VirtualMachine, id string) (*toolbox.Client, error) {
	opmgr := guest.NewOperationsManager(vm.Session.Client.Client, vm.Reference())
	pm, err := opmgr.ProcessManager(op)
	if err != nil {
		op.Debugf("Failed to create new process manager ")
		return nil, err
	}
	fm, err := opmgr.FileManager(op)
	if err != nil {
		op.Debugf("Failed to create new file manager ")
		return nil, err
	}

	return &toolbox.Client{
		ProcessManager: pm,
		FileManager:    fm,
		Authentication: &types.NamePasswordAuthentication{
			Username: id,
		},
	}, nil
}
