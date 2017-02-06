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
	"net/url"
)

//device information for interaction with the tether and portlayer
type NFSVolume struct {

	//This is the nfs host the the volume belongs to
	Host *url.URL

	//Path on the Host where the volume is located
	NFSPath string
}

func NewNFSVolumeDevice(host *url.URL, NFSPath string) NFSVolume {
	v := NFSVolume{
		Host:    host,
		NFSPath: NFSPath,
	}
	return v
}

func (v NFSVolume) MountPath() (string, error) {
	return "", nil
}

// includes url to nfs directory for container to mount,
func (v NFSVolume) DiskPath() url.URL {
	if v.Host == nil {
		return url.URL{}
	}
	return *v.Host
}
