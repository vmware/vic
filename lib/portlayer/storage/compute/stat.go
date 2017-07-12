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

// StatPath will use Guest Tools to stat a given path in the container

package compute

import (
	"os"
	"time"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/errors"
)

type FileStat struct {
	LinkTarget string
	Mode       uint32
	Name       string
	Size       int64
	ModTime	   time.Time
}

// interface for offline container statpath
type ContainerStatPath interface {
	StatPath(op trace.Operation, storeId, deviceId, target string) (*FileStat, error)
}

// InspectFileStat runs lstat on the target
func InspectFileStat(target string) (*FileStat, error) {
	fileInfo, err := os.Lstat(target)
	if err != nil {
		return nil, errors.Errorf("error returned from %s, target %s", err.Error(), target)
	}

	var linkTarget string
	// check for symlink
	if fileInfo.Mode() & os.ModeSymlink != 0 {
		linkTarget, err = os.Readlink(target)
		if err != nil {
			return nil, err
		}
	}

	return &FileStat{linkTarget, uint32(fileInfo.Mode()), fileInfo.Name(), fileInfo.Size(), fileInfo.ModTime()}, nil
}
