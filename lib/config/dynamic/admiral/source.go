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

package admiral

import (
	"context"

	vchcfg "github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/config/dynamic"
)

func NewSource() (dynamic.Source, error) {
	return &source{c: newClient()}, nil
}

type source struct {
	lastCfg *vchcfg.VirtualContainerHostConfigSpec
	c       client
}

// Get returns the dynamic config portion from an Admiral instance. For now,
// this is empty pending details from the Admiral team.
func (a *source) Get(ctx context.Context) (*vchcfg.VirtualContainerHostConfigSpec, error) {
	return nil, dynamic.ErrSourceUnavailable
}
