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

package nfs

import (
	"github.com/vmware/vic/lib/portlayer/storage"
)

type NFSVolumeBacking struct {

	//This is the path where the volume directory exists on the nfs server
	NFSPath string

	//This is the target filesystem path for mounting into a container
	MountPath string
}

func NewNFSVolumeBacking(NFSPath string) NFSVolumeBacking {
	v := NFSVolumeBacking{
		NFSPath:   NFSPath,
		MountPath: MountTargetPath,
	}
	return v
}

func (v *NFSVolumeBacking) MountPath() (string, error) {
	return nil
}

func (v *NFSVolumeBacking) DiskPath() string {
	return v.NFSPath
}
