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
	log "github.com/Sirupsen/logrus"

	"context"

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

func (rp *ResourcePool) GetDatacenter(ctx context.Context) (*object.Datacenter, error) {
	dcRef, err := rp.getLowestAncestor(ctx, "Datacenter")
	if err != nil || dcRef == nil {
		log.Errorf("Unable to get datacenter ancestor of rp %s: %s", rp.Name(), err)
		return nil, errors.Errorf("Unable to get datacenter ancestor of rp %s: %s", rp.Name(), err)
	}

	return object.NewDatacenter(rp.Client.Client, *dcRef), nil
}

func (rp *ResourcePool) getAncestors(ctx context.Context, inType string) ([]types.ManagedObjectReference, error) {
	client := rp.Session.Vim25()

	ancestors, err := mo.Ancestors(ctx, client, client.ServiceContent.PropertyCollector, rp.Reference())
	if err != nil {
		log.Errorf("Unable to get ancestors of rp %s: %s", rp.Name(), err)
		return nil, err
	}

	outAncestors := make([]types.ManagedObjectReference, 0, len(ancestors))
	for _, ancestor := range ancestors {
		if ancestor.Self.Type == inType {
			a := ancestor.Self
			outAncestors = append(outAncestors, a)
		}
	}

	return outAncestors, nil
}

func (rp *ResourcePool) getLowestAncestor(ctx context.Context, inType string) (*types.ManagedObjectReference, error) {
	ancestors, err := rp.getAncestors(ctx, inType)
	if err != nil {
		log.Errorf("Unable to get ancestors of rp %s: %s", rp.Name(), err)
		return nil, err
	}

	if len(ancestors) == 0 {
		return nil, nil
	}

	index := len(ancestors) - 1
	return &ancestors[index], nil
}
