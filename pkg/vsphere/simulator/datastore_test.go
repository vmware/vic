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

	"github.com/vmware/govmomi/vim25/types"
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
