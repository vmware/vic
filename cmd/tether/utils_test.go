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

package main

import (
	"path"
	"testing"
)

func TestSlotToPciPath(t *testing.T) {
	var tests = []struct {
		slot int32
		p    string
		err  error
	}{
		{0, path.Join(pciDevPath, "0000:00:00.0"), nil},
		{32, path.Join(pciDevPath, "0000:00:11.0", "0000:*:00.0"), nil},
		{33, path.Join(pciDevPath, "0000:00:11.0", "0000:*:01.0"), nil},
		{192, path.Join(pciDevPath, "0000:00:16.0", "0000:*:00.0"), nil},
	}

	for _, te := range tests {
		p, err := slotToPCIPath(te.slot)
		if te.err != nil {
			if err == nil {
				t.Fatalf("slotToPCIPath(%d) => (%#v, %#v), want (%s, nil)", te.slot, p, err, te.p)
			}

			continue
		}

		if p != te.p {
			t.Fatalf("slotToPCIPath(%d) => (%#v, %#v), want (%s, nil)", te.slot, p, err, te.p)
		}
	}
}
