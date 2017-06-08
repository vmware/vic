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

package create

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGatewaySpec(t *testing.T) {
	var tests = []struct {
		in   string
		dest []string
		gw   string
		err  error
	}{
		{
			in: "10.10.10.10",
			gw: "10.10.10.10",
		},
		{
			in:   "10.12.0.0/16:10.10.10.10",
			dest: []string{"10.12.0.0/16"},
			gw:   "10.10.10.10",
		},
		{
			in:   "10.13.0.0/16,10.12.0.0/16:10.10.10.10",
			dest: []string{"10.13.0.0/16", "10.12.0.0/16"},
			gw:   "10.10.10.10",
		},
	}

	for _, te := range tests {
		dest, gw, err := parseGatewaySpec(te.in)
		if te.err != nil {
			assert.EqualError(t, err, te.err.Error())
		} else {
			assert.NoError(t, err)
		}

		assert.NotNil(t, gw)
		assert.Equal(t, te.gw, gw.IP.String())

		assert.Equal(t, len(te.dest), len(dest))
		for _, d := range te.dest {
			found := false
			for _, d2 := range dest {
				if d2.String() == d {
					found = true
					break
				}
			}

			assert.True(t, found)
		}
	}
}

func TestProcessVolumeStoreParam(t *testing.T) {
	positiveTestCases := []string{
		"nfs://Shared.Volumes.Org/path/to/store:nfs-volumes",
		"ds://vsphere.target.here/:root-path",
		"no.scheme.target:/with/path:ds-store",
		"looooooooooooooooooooooooooooooong.hoooooooooooooooooooooooooooooooost/short/path:long-check",
		"nfs://0.0.0.0/ip/check:simple-target",
		"nfs://prod.shared.storage/vch_prod/volumes:test-label",
		"ds://0.0.0.0/ip/check?myArg=simple&complex=anotherArg:simple-target:test-label",
	}

	negativeTestCases := []string{
		"ds://vsphere.rocks.com/no/label/here",
		"junk-text-%^()!@#:with-label",
		"junk-text-%^()!@#-no-label",
		":no-text",
		"no-label:",
		"no-label/with/path",
	}

	for _, v := range positiveTestCases {
		target, rawString, label, err := processVolumeStoreParam(v)

		assert.NotNil(t, target, v)
		assert.NotEqual(t, "", rawString, v)
		assert.NotEqual(t, "", label, v)
		assert.Nil(t, err, v)
	}

	for _, v := range negativeTestCases {
		target, _, _, err := processVolumeStoreParam(v)

		// here "" is possible for rawString and label so we check for err and nil target.
		assert.Nil(t, target, v)
		assert.NotNil(t, err, v)
	}
}
