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
	"log"
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
	res *types.HostDatastoreBrowserSearchResults
}

func (s *searchDatastoreTask) addFile(dir string, name string) {
	details := s.req.SearchSpec.Details
	file := path.Join(dir, name)
	st, err := os.Stat(file)
	if err != nil {
		log.Printf("stat(%s): %s", file, err)
		return
	}

	info := types.FileInfo{
		Path: name,
	}

	var finfo types.BaseFileInfo

	if details.FileSize {
		info.FileSize = st.Size()
	}

	if details.Modification {
		mtime := st.ModTime()
		info.Modification = &mtime
	}

	if isTrue(details.FileOwner) {
		// Assume for now this process created all files in the datastore
		user := os.Getenv("USER")

		info.Owner = user
	}

	if st.IsDir() {
		finfo = &types.FolderFileInfo{FileInfo: info}
	} else if details.FileType {
		switch path.Ext(name) {
		case ".img":
			finfo = &types.FloppyImageFileInfo{FileInfo: info}
		case ".iso":
			finfo = &types.IsoImageFileInfo{FileInfo: info}
		case ".log":
			finfo = &types.VmLogFileInfo{FileInfo: info}
		case ".nvram":
			finfo = &types.VmNvramFileInfo{FileInfo: info}
		case ".vmdk":
			// TODO: lookup device to set other fields
			finfo = &types.VmDiskFileInfo{FileInfo: info}
		case ".vmx":
			finfo = &types.VmConfigFileInfo{FileInfo: info}
		default:
			finfo = &info
		}
	}

	s.res.File = append(s.res.File, finfo)
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

	dir := path.Join(ds.Info.GetDatastoreInfo().Url, p.Path)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		ff := types.FileFault{
			File: p.Path,
		}
		if os.IsNotExist(err) {
			return nil, &types.FileNotFound{FileFault: ff}
		}

		return nil, &ff
	}

	s.res = &types.HostDatastoreBrowserSearchResults{
		Datastore:  &ds.Self,
		FolderPath: s.req.DatastorePath,
	}

	for _, file := range files {
		for _, m := range s.req.SearchSpec.MatchPattern {
			if ok, _ := path.Match(m, file.Name()); ok {
				s.addFile(dir, file.Name())
			}
		}
	}

	return s.res, nil
}

func (b *HostDatastoreBrowser) SearchDatastoreTask(s *types.SearchDatastore_Task) soap.HasFault {
	task := NewTask(&searchDatastoreTask{b, s, nil})

	task.Run()

	return &methods.SearchDatastore_TaskBody{
		Res: &types.SearchDatastore_TaskResponse{
			Returnval: task.Self,
		},
	}
}
