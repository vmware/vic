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
	"testing"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
	"github.com/vmware/vic/pkg/vsphere/simulator/vc"
	"golang.org/x/net/context"
)

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

	_, err = ff.CreateDatacenter(ctx, "biz")
	if err != nil {
		t.Error(err)
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
