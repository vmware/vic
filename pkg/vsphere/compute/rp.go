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

package compute

import (
	"path"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"

	"github.com/vmware/vic/pkg/vsphere/session"
)

// ResourcePool struct defines the ResourcePool which provides additional
// VIC specific methods over object.ResourcePool as well as keeps some state
type ResourcePool struct {
	*object.ResourcePool

	*session.Session
}

// NewResourcePool returns a New ResourcePool object
func NewResourcePool(ctx context.Context, session *session.Session, moref types.ManagedObjectReference) *ResourcePool {
	return &ResourcePool{
		ResourcePool: object.NewResourcePool(
			session.Vim25(),
			moref,
		),
		Session: session,
	}
}

func FindResourcePool(ctx context.Context, s *session.Session, name string) (*ResourcePool, error) {
	var err error

	if name != "" {
		if !path.IsAbs(name) {
			name = path.Join(s.Cluster.InventoryPath, "Resources", name)
		}
	} else {
		name = path.Join(s.Cluster.InventoryPath, "Resources")
	}

	pool, err := s.Finder.ResourcePoolOrDefault(ctx, name)
	if err != nil {
		return nil, err
	}
	return NewResourcePool(ctx, s, pool.Reference()), nil
}
