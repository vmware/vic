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

package toolbox

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vmware/vic/lib/tether-ng/types"
)

type Toolbox struct {
	uuid uuid.UUID
	ctx  context.Context
}

func NewToolbox(ctx context.Context) *Toolbox {
	return &Toolbox{
		uuid: uuid.New(),
		ctx:  ctx,
	}
}

func (t *Toolbox) Configure(ctx context.Context, config *types.ExecutorConfig) error {
	fmt.Printf("toolbox\n")

	return nil

}

func (t *Toolbox) Start(ctx context.Context) error    { return nil }
func (t *Toolbox) Stop(ctx context.Context) error     { return nil }
func (t *Toolbox) UUID(ctx context.Context) uuid.UUID { return t.uuid }

func (t *Toolbox) Toolbox() error {
	fmt.Printf("boxing\n")

	return nil
}
