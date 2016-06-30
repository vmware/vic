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

package etcconf

import (
	"net"
	"testing"
)

func TestHostEntryMerge(t *testing.T) {
	var tests = []struct {
		e1, e2 HostEntry
		res    HostEntry
		mod    bool
	}{
		{HostEntry{IP: net.ParseIP("10.10.10.10"), Hostnames: []string{"bar"}}, HostEntry{IP: net.ParseIP("10.10.10.10"), Hostnames: []string{"foo"}}, HostEntry{IP: net.ParseIP("10.10.10.10"), Hostnames: []string{"foo", "bar"}}, true},
		{HostEntry{IP: net.ParseIP("10.10.10.10"), Hostnames: []string{"bar"}}, HostEntry{IP: net.ParseIP("10.10.10.10"), Hostnames: []string{"bar"}}, HostEntry{IP: net.ParseIP("10.10.10.10"), Hostnames: []string{"bar"}}, false},
	}

	for _, te := range tests {
		mod := te.e1.merge(te.e2)
		if !te.e1.Equal(te.res) || mod != te.mod {
			t.Fatalf("got %+v, want %+v", te.e1, te.res)
		}
	}
}
