// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package registry

import "testing"
import "github.com/stretchr/testify/assert"

func TestContains(t *testing.T) {
	var tests = []struct {
		first, second Entry
		res           bool
	}{
		{
			first:  &ipEntry{e: "192.168.0.1"},
			second: &ipEntry{e: "192.168.0.1"},
			res:    true,
		},
		{
			first:  &ipEntry{e: "192.168.0.1"},
			second: &ipEntry{e: "192.168.0.1/32"},
			res:    true,
		},
		{
			first:  &ipEntry{e: "192.168.0.1"},
			second: &ipEntry{e: "192.168.0.1/16"},
			res:    false,
		},
		{
			first:  &ipEntry{e: "192.168.0.1"},
			second: &ipEntry{e: "192.168.0.2"},
			res:    false,
		},
		{
			first:  &domainEntry{e: "*.google.com"},
			second: &domainEntry{e: "*.com"},
			res:    false,
		},
		{
			first:  &domainEntry{e: "mail.google.com"},
			second: &domainEntry{e: "*.google.com"},
			res:    false,
		},
		{
			first:  &domainEntry{e: "*.google.com"},
			second: &domainEntry{e: "mail.google.com"},
			res:    true,
		},
		{
			first:  &domainEntry{e: "*.com"},
			second: &domainEntry{e: "*.google.com"},
			res:    true,
		},
	}

	for _, te := range tests {
		assert.Equal(t, te.res, te.first.Contains(te.second))
	}

}
