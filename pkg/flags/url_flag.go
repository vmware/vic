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

package flags

import (
	"flag"
	"net/url"

	"github.com/vmware/govmomi/vim25/soap"
)

type URLFlag struct {
	u **url.URL
}

func (f *URLFlag) Set(s string) error {
	var err error
	url, err := soap.ParseURL(s)
	*f.u = url
	return err
}

func (f *URLFlag) Get() interface{} {
	return *f.u
}

func (f *URLFlag) String() string {
	if *f.u == nil {
		return "<nil>"
	}
	return (*f.u).String()
}

func (f *URLFlag) IsBoolFlag() bool { return false }

// NewURLFlag returns a flag.Value.
func NewURLFlag(u **url.URL) flag.Value {
	return &URLFlag{u}
}
