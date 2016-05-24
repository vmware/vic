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
	"reflect"
	"testing"

	"golang.org/x/net/context"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
	"github.com/vmware/vic/pkg/vsphere/simulator/vc"
)

func TestRetrieveProperties(t *testing.T) {
	configs := []struct {
		folder  mo.Folder
		content types.ServiceContent
		dc      *types.ManagedObjectReference
	}{
		{esx.RootFolder, esx.ServiceContent, &esx.Datacenter.Self},
		{vc.RootFolder, vc.ServiceContent, nil},
	}

	for _, config := range configs {
		s := New(NewServiceInstance(config.content, config.folder))

		ts := s.NewServer()
		defer ts.Close()

		ctx := context.Background()

		client, err := govmomi.NewClient(ctx, ts.URL, true)
		if err != nil {
			t.Fatal(err)
		}

		if config.dc == nil {
			dc, cerr := object.NewRootFolder(client.Client).CreateDatacenter(ctx, "dc1")
			if cerr != nil {
				t.Fatal(cerr)
			}
			ref := dc.Reference()
			config.dc = &ref
		}

		// Retrieve a specific property
		f := mo.Folder{}
		err = client.RetrieveOne(ctx, config.content.RootFolder, []string{"name"}, &f)
		if err != nil {
			t.Fatal(err)
		}

		if f.Name != config.folder.Name {
			t.Fail()
		}

		// Retrieve all properties
		f = mo.Folder{}
		err = client.RetrieveOne(ctx, config.content.RootFolder, nil, &f)
		if err != nil {
			t.Fatal(err)
		}

		if f.Name != config.folder.Name {
			t.Fatalf("'%s' vs '%s'", f.Name, config.folder.Name)
		}

		// Retrieve an ArrayOf property
		f = mo.Folder{}
		err = client.RetrieveOne(ctx, config.content.RootFolder, []string{"childEntity"}, &f)
		if err != nil {
			t.Fatal(err)
		}

		if len(f.ChildEntity) != 1 {
			t.Fail()
		}

		es, err := mo.Ancestors(ctx, client.Client, config.content.PropertyCollector, config.content.RootFolder)
		if err != nil {
			t.Fatal(err)
		}

		if len(es) != 1 {
			t.Fail()
		}

		finder := find.NewFinder(client.Client, false)
		dc, err := finder.DatacenterOrDefault(ctx, "")
		if err != nil {
			t.Fatal(err)
		}

		if dc.Reference() != *config.dc {
			t.Fail()
		}

		finder.SetDatacenter(dc)

		es, err = mo.Ancestors(ctx, client.Client, config.content.PropertyCollector, dc.Reference())
		if err != nil {
			t.Fatal(err)
		}

		expect := map[string]types.ManagedObjectReference{
			"Folder":     config.folder.Reference(),
			"Datacenter": dc.Reference(),
		}

		if len(es) != len(expect) {
			t.Fail()
		}

		for _, e := range es {
			ref := e.Reference()
			if r, ok := expect[ref.Type]; ok {
				if r != ref {
					t.Errorf("%#v vs %#v", r, ref)
				}
			} else {
				t.Errorf("unexpected object %#v", e.Reference())
			}
		}

		// finder tests
		ls, err := finder.ManagedObjectListChildren(ctx, ".")
		if err != nil {
			t.Error(err)
		}

		folders, err := dc.Folders(ctx)
		if err != nil {
			t.Fatal(err)
		}

		// Validated name properties are recursively retrieved for the datacenter and its folder children
		ipaths := []string{
			folders.VmFolder.InventoryPath,
			folders.HostFolder.InventoryPath,
			folders.DatastoreFolder.InventoryPath,
			folders.NetworkFolder.InventoryPath,
		}

		var lpaths []string
		for _, p := range ls {
			lpaths = append(lpaths, p.Path)
		}

		if !reflect.DeepEqual(ipaths, lpaths) {
			t.Errorf("%#v != %#v\n", ipaths, lpaths)
		}

		// We have no VMs, expect NotFoundError
		_, err = finder.VirtualMachineList(ctx, "*")
		if err == nil {
			t.Error("expected error")
		} else {
			if _, ok := err.(*find.NotFoundError); !ok {
				t.Error(err)
			}
		}

		// Retrieve a missing property
		mdc := mo.Datacenter{}
		err = client.RetrieveOne(ctx, dc.Reference(), []string{"enoent"}, &mdc)
		if err == nil {
			t.Error("expected error")
		} else {
			switch fault := soap.ToVimFault(err).(type) {
			case *types.InvalidProperty:
				// ok
			default:
				t.Errorf("unexpected fault: %#v", fault)
			}
		}

		// Retrieve a nested property
		Map.Get(dc.Reference()).(*mo.Datacenter).Configuration.DefaultHardwareVersionKey = "foo"
		mdc = mo.Datacenter{}
		err = client.RetrieveOne(ctx, dc.Reference(), []string{"configuration.defaultHardwareVersionKey"}, &mdc)
		if err != nil {
			t.Fatal(err)
		}
		if mdc.Configuration.DefaultHardwareVersionKey != "foo" {
			t.Fail()
		}

		// Retrieve a missing nested property
		mdc = mo.Datacenter{}
		err = client.RetrieveOne(ctx, dc.Reference(), []string{"configuration.enoent"}, &mdc)
		if err == nil {
			t.Error("expected error")
		} else {
			switch fault := soap.ToVimFault(err).(type) {
			case *types.InvalidProperty:
				// ok
			default:
				t.Errorf("unexpected fault: %#v", fault)
			}
		}

		// Retrieve an empty property
		err = client.RetrieveOne(ctx, dc.Reference(), []string{""}, &mdc)
		if err != nil {
			t.Error(err)
		}

		// Expect ManagedObjectNotFoundError
		Map.Remove(dc.Reference())
		err = client.RetrieveOne(ctx, dc.Reference(), []string{"name"}, &mdc)
		if err == nil {
			t.Fatal("expected error")
		}
	}
}
