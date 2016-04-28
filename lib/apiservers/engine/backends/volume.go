// Copyright 2016 VMware, Inc. All Rights Reserved.
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

package vicbackends

import (
	"fmt"

	"github.com/docker/engine-api/types"
)

type Volume struct {
	ProductName string
}

func (v *Volume) Volumes(filter string) ([]*types.Volume, []string, error) {
	return nil, make([]string, 0), fmt.Errorf("%s does not implement volume.Volumes", v.ProductName)
}

func (v *Volume) VolumeInspect(name string) (*types.Volume, error) {
	return nil, fmt.Errorf("%s does not implement volume.VolumeInspect", v.ProductName)
}

func (v *Volume) VolumeCreate(name, driverName string, opts, labels map[string]string) (*types.Volume, error) {
	return nil, fmt.Errorf("%s does not implement volume.VolumeCreate", v.ProductName)
}

func (v *Volume) VolumeRm(name string) error {
	return fmt.Errorf("%s does not implement volume.VolumeRm", v.ProductName)
}
