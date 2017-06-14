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
	"errors"

	"github.com/vmware/vic/lib/config"
)

var SourceUnavailableErr = errors.New("source unavailable")
var SourceNotModifiedErr = errors.New("source not modified")
var AccessDeniedErr = errors.New("access denied")

type Source interface {
	Get(ctx context.Context) (*config.VirtualContainerHostConfigSpec, error)
}

type Merger interface {
	Merge(orig, other *config.VirtualContainerHostConfigSpec) (*config.VirtualContainerHostConfigSpec, error)
}
