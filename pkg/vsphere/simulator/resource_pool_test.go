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

package simulator

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
)

func TestResourcePool(t *testing.T) {
	ctx := context.Background()

	m := &Model{
		ServiceContent: esx.ServiceContent,
		RootFolder:     esx.RootFolder,
	}

	err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	c := m.Service.client

	finder := find.NewFinder(c, false)
	finder.SetDatacenter(object.NewDatacenter(c, esx.Datacenter.Reference()))

	spec := NewResourceConfigSpec()

	parent := object.NewResourcePool(c, esx.ResourcePool.Self)

	// can't destroy a root pool
	task, err := parent.Destroy(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err = task.Wait(ctx); err == nil {
		t.Fatal("expected error destroying a root pool")
	}

	// create a child pool
	childName := uuid.New().String()

	child, err := parent.Create(ctx, childName, spec)
	if err != nil {
		t.Fatal(err)
	}

	if child.Reference() == esx.ResourcePool.Self {
		t.Error("expected new pool Self reference")
	}

	// create a grandchild pool
	grandChildName := uuid.New().String()
	_, err = child.Create(ctx, grandChildName, spec)
	if err != nil {
		t.Fatal(err)
	}

	// create sibling (of the grand child) pool
	siblingName := uuid.New().String()
	_, err = child.Create(ctx, siblingName, spec)
	if err != nil {
		t.Fatal(err)
	}

	// finder should return the 2 grand children
	pools, err := finder.ResourcePoolList(ctx, "*/Resources/"+childName+"/*")
	if err != nil {
		t.Fatal(err)
	}
	if len(pools) != 2 {
		t.Fatalf("len(pools) == %d", len(pools))
	}

	// destroy the child
	task, err = child.Destroy(ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = task.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// finder should error not found after destroying the child
	_, err = finder.ResourcePoolList(ctx, "*/Resources/"+childName+"/*")
	if err == nil {
		t.Fatal("expected not found error")
	}

	// since the child was destroyed, grand child pools should now be children of the root pool
	pools, err = finder.ResourcePoolList(ctx, "*/Resources/*")
	if err != nil {
		t.Fatal(err)
	}

	if len(pools) != 2 {
		t.Fatalf("len(pools) == %d", len(pools))
	}
}
