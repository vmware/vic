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

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
	"github.com/vmware/vic/pkg/vsphere/simulator/vc"
)

func addStandaloneHost(folder *object.Folder, spec types.HostConnectSpec) (*object.Task, error) {
	// TODO: add govmomi wrapper
	req := types.AddStandaloneHost_Task{
		This:         folder.Reference(),
		Spec:         spec,
		AddConnected: true,
	}

	res, err := methods.AddStandaloneHost_Task(context.TODO(), folder.Client(), &req)
	if err != nil {
		return nil, err
	}

	task := object.NewTask(folder.Client(), res.Returnval)
	return task, nil
}

func TestFolderESX(t *testing.T) {
	content := esx.ServiceContent
	s := New(NewServiceInstance(content, esx.RootFolder))

	ts := s.NewServer()
	defer ts.Close()

	ctx := context.Background()
	c, err := govmomi.NewClient(ctx, ts.URL, true)
	if err != nil {
		t.Fatal(err)
	}

	f := object.NewRootFolder(c.Client)

	_, err = f.CreateFolder(ctx, "foo")
	if err == nil {
		t.Error("expected error")
	}

	_, err = f.CreateDatacenter(ctx, "foo")
	if err == nil {
		t.Error("expected error")
	}

	finder := find.NewFinder(c.Client, false)
	dc, err := finder.DatacenterOrDefault(ctx, "")
	if err != nil {
		t.Fatal(err)
	}

	folders, err := dc.Folders(ctx)
	if err != nil {
		t.Fatal(err)
	}

	spec := types.HostConnectSpec{}
	_, err = addStandaloneHost(folders.HostFolder, spec)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFolderVC(t *testing.T) {
	content := vc.ServiceContent
	s := New(NewServiceInstance(content, vc.RootFolder))

	ts := s.NewServer()
	defer ts.Close()

	ctx := context.Background()
	c, err := govmomi.NewClient(ctx, ts.URL, true)
	if err != nil {
		t.Fatal(err)
	}

	f := object.NewRootFolder(c.Client)

	ff, err := f.CreateFolder(ctx, "foo")
	if err != nil {
		t.Error(err)
	}

	dc, err := f.CreateDatacenter(ctx, "bar")
	if err != nil {
		t.Error(err)
	}

	for _, ref := range []object.Reference{ff, dc} {
		o := Map.Get(ref.Reference())
		if o == nil {
			t.Fatalf("failed to find %#v", ref)
		}

		e := o.(mo.Entity).Entity()
		if *e.Parent != f.Reference() {
			t.Fail()
		}
	}

	dc, err = ff.CreateDatacenter(ctx, "biz")
	if err != nil {
		t.Error(err)
	}

	folders, err := dc.Folders(ctx)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		state types.TaskInfoState
	}{
		{"", types.TaskInfoStateError},
		{"foo.local", types.TaskInfoStateSuccess},
	}

	for _, test := range tests {
		spec := types.HostConnectSpec{
			HostName: test.name,
		}

		task, err := addStandaloneHost(folders.HostFolder, spec)
		if err != nil {
			t.Fatal(err)
		}

		res, err := task.WaitForResult(ctx, nil)
		if test.state == types.TaskInfoStateError {
			if err == nil {
				t.Error("expected error")
			}

			if res.Result != nil {
				t.Error("expected nil")
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}

			ref, ok := res.Result.(types.ManagedObjectReference)
			if !ok {
				t.Errorf("expected moref, got type=%T", res.Result)
			}
			host := Map.Get(ref).(*HostSystem)
			if host.Name != test.name {
				t.Fail()
			}

			if ref == esx.HostSystem.Self {
				t.Error("expected new host Self reference")
			}

			pool := Map.Get(*host.Parent).(*mo.ComputeResource).ResourcePool
			if *pool == esx.ResourcePool.Self {
				t.Error("expected new pool Self reference")
			}
		}

		if res.State != test.state {
			t.Fatalf("%s", res.State)
		}
	}
}

func TestFolderFaults(t *testing.T) {
	f := Folder{}
	f.ChildType = []string{"VirtualMachine"}

	if f.CreateFolder(nil).Fault() == nil {
		t.Error("expected fault")
	}

	if f.CreateDatacenter(nil).Fault() == nil {
		t.Error("expected fault")
	}
}
