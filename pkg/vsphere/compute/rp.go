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

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/context"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
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

func (rp *ResourcePool) GetChildrenVMs(ctx context.Context, s *session.Session) ([]*vm.VirtualMachine, error) {
	var err error
	var mrp mo.ResourcePool
	var vms []*vm.VirtualMachine

	if err = rp.Properties(ctx, rp.Reference(), []string{"vm"}, &mrp); err != nil {
		log.Errorf("Unable to get children vm of resource pool %s: %s", rp.Name(), err)
		return vms, err
	}

	for _, o := range mrp.Vm {
		v := vm.NewVirtualMachine(ctx, s, o)
		vms = append(vms, v)
	}
	return vms, nil
}

func (rp *ResourcePool) GetChildVM(ctx context.Context, s *session.Session, name string) (*vm.VirtualMachine, error) {
	searchIndex := object.NewSearchIndex(s.Client.Client)
	child, err := searchIndex.FindChild(ctx, rp.Reference(), name)
	if err != nil {
		return nil, errors.Errorf("Unable to find VM(%s): %s", name, err.Error())
	}
	if child == nil {
		return nil, nil
	}
	// instantiate the vm object
	return vm.NewVirtualMachine(ctx, s, child.Reference()), nil
}

func (rp *ResourcePool) GetCluster(ctx context.Context) (*object.ComputeResource, error) {
	var err error
	var mrp mo.ResourcePool

	if err = rp.Properties(ctx, rp.Reference(), []string{"owner"}, &mrp); err != nil {
		log.Errorf("Unable to get cluster of resource pool %s: %s", rp.Name(), err)
		return nil, err
	}

	return object.NewComputeResource(rp.Client.Client, mrp.Owner), nil
}
