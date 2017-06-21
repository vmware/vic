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

package dynamic

import (
	"context"

	"github.com/vmware/vic/lib/config"
)

type AdmiralSource struct {
}

// Get returns the dynamic config portion from an Admiral instance. For now,
// this is empty pending details from the Admiral team.
func (a *AdmiralSource) Get(context.Context) (*config.VirtualContainerHostConfigSpec, error) {
	return nil, ErrSourceUnavailable
}
