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
	"io/ioutil"
	"os"
	"path"

	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type HostDatastoreBrowser struct {
	mo.HostDatastoreBrowser
}

type searchDatastoreTask struct {
	*HostDatastoreBrowser

	req *types.SearchDatastore_Task
}

func (s *searchDatastoreTask) Run(Task *Task) (types.AnyType, types.BaseMethodFault) {
	p, fault := parseDatastorePath(s.req.DatastorePath)
	if fault != nil {
		return nil, fault
	}

	ref := Map.FindByName(p.Datastore, s.Datastore)
	if ref == nil {
		return nil, &types.InvalidDatastore{Name: p.Datastore}
	}

	ds := ref.(*Datastore)

	dir := ds.Info.GetDatastoreInfo().Url

	files, err := ioutil.ReadDir(path.Join(dir, p.Path))
	if err != nil {
		ff := types.FileFault{
			File: p.Path,
		}
		if os.IsNotExist(err) {
			return nil, &types.FileNotFound{FileFault: ff}
		}

		return nil, &ff
	}

	res := &types.HostDatastoreBrowserSearchResults{
		Datastore:  &ds.Self,
		FolderPath: s.req.DatastorePath,
	}

	for _, file := range files {
		for _, m := range s.req.SearchSpec.MatchPattern {
			if ok, _ := path.Match(m, file.Name()); ok {
				info := &types.FileInfo{
					Path: file.Name(),
				}
				res.File = append(res.File, info)
			}
		}
	}

	return res, nil
}

func (b *HostDatastoreBrowser) SearchDatastoreTask(s *types.SearchDatastore_Task) soap.HasFault {
	task := NewTask(&searchDatastoreTask{b, s})

	task.Run()

	return &methods.SearchDatastore_TaskBody{
		Res: &types.SearchDatastore_TaskResponse{
			Returnval: task.Self,
		},
	}
}
