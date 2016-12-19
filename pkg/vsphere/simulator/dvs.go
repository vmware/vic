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
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type VmwareDistributedVirtualSwitch struct {
	mo.VmwareDistributedVirtualSwitch
}

func (s *VmwareDistributedVirtualSwitch) AddDVPortgroupTask(c *types.AddDVPortgroup_Task) soap.HasFault {
	task := CreateTask(s, "addDVPortgroupTask", func(t *Task) (types.AnyType, types.BaseMethodFault) {
		f := Map.getEntityParent(s, "Folder").(*Folder)

		for _, spec := range c.Spec {
			pg := &mo.DistributedVirtualPortgroup{}
			pg.Name = spec.Name
			pg.Entity().Name = pg.Name

			if obj := Map.FindByName(pg.Name, f.ChildEntity); obj != nil {
				return nil, &types.DuplicateName{
					Name:   pg.Name,
					Object: obj.Reference(),
				}
			}

			f.putChild(pg)
		}

		return nil, nil
	})

	task.Run()

	return &methods.AddDVPortgroup_TaskBody{
		Res: &types.AddDVPortgroup_TaskResponse{
			Returnval: task.Self,
		},
	}
}
