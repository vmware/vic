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

// StatPath will use Guest Tools to stat a given path in the container

package compute

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"
)

func StatPath(op trace.Operation, vc *exec.Container, path string) (*types.GuestFileInfo, error) {
	defer trace.End(trace.Begin(""))

	// err returned if file does not exist on the host
	filemgr, err := guest.NewOperationsManager(vc.VIM25Reference(), vc.VMReference()).FileManager(op)
	if err != nil {
		return nil, err
	}

	auth := types.NamePasswordAuthentication{
		Username: vc.ExecConfig.ID,
	}

	// List the files on the container path. If the path represents a regular
	// file, it's file info are returned. If the path represents a directory,
	// ListFiles will return all the files in that directory. This means that if
	// the for loop does not return becuase the path does not exist in files.Filed.
	// Consequently, we can safely assume the target is a directory.
	var offset int32
	files, err := filemgr.ListFiles(op, &auth, path, offset, 0, "")
	if err != nil {
		return nil, fmt.Errorf("file listing for container %s failed\n: %s", vc.ExecConfig.ID, err.Error())
	}

	for _, file := range files.Files {
		op.Debugf("Stats for file %s --- %s\n", path, file)
		if file.Path == filepath.Base(path) {
			return &file, nil
		}
	}

	time := time.Now()
	return &types.GuestFileInfo{
		Path: path,
		Type: string(types.GuestFileTypeDirectory),
		Size: int64(4096),
		Attributes: &types.GuestFileAttributes{
			ModificationTime: &time,
		},
	}, nil
}
