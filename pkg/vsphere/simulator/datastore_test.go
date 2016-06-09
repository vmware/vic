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
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"

	"golang.org/x/net/context"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/simulator/esx"
	"github.com/vmware/vic/pkg/vsphere/simulator/vc"
)

func TestParseDatastorePath(t *testing.T) {
	tests := []struct {
		dsPath string
		dsFile string
		fail   bool
	}{
		{"", "", true},
		{"x", "", true},
		{"[", "", true},
		{"[nope", "", true},
		{"[test]", "", false},
		{"[test] foo", "foo", false},
		{"[test] foo/foo.vmx", "foo/foo.vmx", false},
		{"[test]foo bar/foo bar.vmx", "foo bar/foo bar.vmx", false},
	}

	for _, test := range tests {
		p, err := parseDatastorePath(test.dsPath)
		if test.fail {
			if err == nil {
				t.Errorf("expected error for: %s", test.dsPath)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error '%#v' for: %s", err, test.dsPath)
			} else {
				if test.dsFile != p.Path {
					t.Errorf("dsFile=%s", p.Path)
				}
				if p.Datastore != "test" {
					t.Errorf("ds=%s", p.Datastore)
				}
			}
		}
	}
}

func TestRefreshDatastore(t *testing.T) {
	tests := []struct {
		dir  string
		fail bool
	}{
		{".", false},
		{"-", true},
	}

	for _, test := range tests {
		ds := &Datastore{}
		ds.Info = &types.LocalDatastoreInfo{
			DatastoreInfo: types.DatastoreInfo{
				Url: test.dir,
			},
		}

		res := ds.RefreshDatastore(nil)
		err := res.Fault()

		if test.fail {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestDatastoreHTTP(t *testing.T) {
	models := []func() (*browseDatastoreModel, error){
		browseDatastoreModelWithESX,
		browseDatastoreModelWithVC,
	}

	ctx := context.Background()
	src := "datastore_test.go"
	dst := "tmp.go"

	for _, model := range models {
		m, err := model()
		defer m.Server.Close()

		download := func(name string, fail bool) {
			_, _, err = m.Datastore.Download(ctx, name, nil)
			if fail {
				if err == nil {
					t.Fatal("expected error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		}

		upload := func(name string, fail bool, method string) {
			f, err := os.Open(src)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			p := soap.DefaultUpload
			p.Method = method

			err = m.Datastore.Upload(ctx, f, name, &p)
			if fail {
				if err == nil {
					t.Fatalf("%s %s: expected error", method, name)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
			}
		}

		// GET file does not exist = fail
		download(dst, true)

		// PUT file = ok
		upload(dst, false, "PUT")

		// GET file exists = ok
		download(dst, false)

		// POST file exists = fail
		upload(dst, true, "POST")

		// POST file does not exist = ok
		upload(dst+".post", false, "POST")

		// PATCH method not supported = fail
		upload(dst+".patch", true, "PATCH")

		// PUT path is directory = fail
		upload("", true, "PUT")

		// GET datastore does not exist = fail
		download(dst, false)

		// cover the case where datacenter or datastore lookup fails
		for _, q := range []string{"dcName=nope", "dsName=nope"} {
			u := *m.Server.URL
			u.RawQuery = q
			u.Path = path.Join(folderPrefix, dst)

			r, err := http.Get(u.String())
			if err != nil {
				t.Fatal(err)
			}

			if r.StatusCode == http.StatusOK {
				t.Error("expected failure")
			}
		}
	}
}

// TODO: this is a copy of vmCreateModel; we should make these helpers public elsewhere
type browseDatastoreModel struct {
	Service *Service
	Server  *Server

	Client     *govmomi.Client
	Finder     *find.Finder
	Datacenter *object.Datacenter
	Folders    *object.DatacenterFolders
	Cluster    *object.ClusterComputeResource
	Pool       *object.ResourcePool
	Host       *object.HostSystem
	Datastore  *object.Datastore
}

func browseDatastoreModelWithESX() (*browseDatastoreModel, error) {
	m := new(browseDatastoreModel)

	m.Service = New(NewServiceInstance(esx.ServiceContent, esx.RootFolder))
	m.Server = m.Service.NewServer()

	ctx := context.Background()

	c, err := govmomi.NewClient(ctx, m.Server.URL, true)
	if err != nil {
		return nil, err
	}

	m.Client = c

	m.Finder = find.NewFinder(c.Client, false)

	m.Datacenter, err = m.Finder.DefaultDatacenter(ctx)
	if err != nil {
		return nil, err
	}

	m.Finder.SetDatacenter(m.Datacenter)

	m.Folders, err = m.Datacenter.Folders(ctx)
	if err != nil {
		return nil, err
	}

	m.Host, err = m.Finder.DefaultHostSystem(ctx)
	if err != nil {
		return nil, err
	}

	m.Pool, err = m.Finder.DefaultResourcePool(ctx)
	if err != nil {
		return nil, err
	}

	_, err = m.Server.TempDatastore(ctx, m.Host)
	if err != nil {
		return nil, err
	}

	m.Datastore, err = m.Finder.DefaultDatastore(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func browseDatastoreModelWithVC() (*browseDatastoreModel, error) {
	m := new(browseDatastoreModel)

	m.Service = New(NewServiceInstance(vc.ServiceContent, vc.RootFolder))

	m.Server = m.Service.NewServer()

	ctx := context.Background()

	c, err := govmomi.NewClient(ctx, m.Server.URL, true)
	if err != nil {
		return nil, err
	}

	m.Client = c

	f := object.NewRootFolder(c.Client)

	m.Datacenter, err = f.CreateDatacenter(ctx, "dc1")
	if err != nil {
		return nil, err
	}

	m.Finder = find.NewFinder(c.Client, false)

	m.Datacenter, err = m.Finder.DefaultDatacenter(ctx)
	if err != nil {
		return nil, err
	}

	m.Finder.SetDatacenter(m.Datacenter)

	m.Folders, err = m.Datacenter.Folders(ctx)
	if err != nil {
		return nil, err
	}

	m.Cluster, err = m.Folders.HostFolder.CreateCluster(ctx, "cluster1", types.ClusterConfigSpecEx{})
	if err != nil {
		return nil, err
	}

	m.Pool, err = m.Finder.ResourcePool(ctx, "*/*")
	if err != nil {
		return nil, err
	}

	var hosts []*object.HostSystem

	for i := 0; i < 3; i++ {
		spec := types.HostConnectSpec{
			HostName: fmt.Sprintf("host-%d", i),
		}

		task, cerr := m.Cluster.AddHost(ctx, spec, true, nil, nil)
		if cerr != nil {
			return nil, cerr
		}

		info, cerr := task.WaitForResult(ctx, nil)
		if cerr != nil {
			return nil, cerr
		}

		hosts = append(hosts, object.NewHostSystem(c.Client, info.Result.(types.ManagedObjectReference)))
	}

	_, err = m.Server.TempDatastore(ctx, hosts...)
	if err != nil {
		return nil, err
	}

	m.Datastore, err = m.Finder.DefaultDatastore(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}
