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
)

func TestSearchIndex(t *testing.T) {
	ctx := context.Background()

	for _, model := range []*Model{ESX(), VPX()} {
		defer model.Remove()
		err := model.Create()
		if err != nil {
			t.Fatal(err)
		}

		s := model.Service.NewServer()
		defer s.Close()

		c, err := govmomi.NewClient(ctx, s.URL, true)
		if err != nil {
			t.Fatal(err)
		}

		finder := find.NewFinder(c.Client, false)
		dc, err := finder.DefaultDatacenter(ctx)
		if err != nil {
			t.Fatal(err)
		}

		finder.SetDatacenter(dc)

		vms, err := finder.VirtualMachineList(ctx, "*")
		if err != nil {
			t.Fatal(err)
		}

		vm := Map.Get(vms[0].Reference()).(*VirtualMachine)

		si := object.NewSearchIndex(c.Client)

		ref, err := si.FindByDatastorePath(ctx, dc, vm.Config.Files.VmPathName)
		if err != nil {
			t.Fatal(err)
		}

		if ref.Reference() != vm.Reference() {
			t.Errorf("moref mismatch %s != %s", ref, vm.Reference())
		}

		ref, err = si.FindByDatastorePath(ctx, dc, vm.Config.Files.VmPathName+"enoent")
		if err != nil {
			t.Fatal(err)
		}

		if ref != nil {
			t.Errorf("ref=%s", ref)
		}
	}
}
