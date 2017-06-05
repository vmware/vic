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

package network

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vmware/vic/lib/tether-ng/types"
)

type Network struct {
	uuid uuid.UUID
	ctx  context.Context
}

func NewNetwork(ctx context.Context) *Network {
	return &Network{
		uuid: uuid.New(),
		ctx:  ctx,
	}
}

func (n *Network) Configure(ctx context.Context, config *types.ExecutorConfig) error {
	fmt.Printf("network\n")

	return nil
}

func (n *Network) Start(ctx context.Context) error    { return nil }
func (n *Network) Stop(ctx context.Context) error     { return nil }
func (n *Network) UUID(ctx context.Context) uuid.UUID { return n.uuid }

func (n *Network) Network() error {
	fmt.Printf("networking\n")

	return nil
}
