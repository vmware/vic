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

package common

import (
	"gopkg.in/urfave/cli.v1"
)

type Registries struct {
	RegistryCAs         cli.StringSlice `arg:"registry-ca"`
	InsecureRegistries  cli.StringSlice `arg:"insecure-registry"`
	WhitelistRegistries cli.StringSlice `arg:"whitelist-registry"`
}

func (r *Registries) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:   "registry-ca, rc",
			Usage:  "Specify a list of additional certificate authority files to use to verify secure registry servers",
			Value:  &r.RegistryCAs,
			Hidden: true,
		},
		cli.StringSliceFlag{
			Name:  "insecure-registry, dir",
			Value: &r.InsecureRegistries,
			Usage: "Specify a list of permitted insecure registry server addresses",
		},
		cli.StringSliceFlag{
			Name:  "whitelist-registry, wr",
			Value: &r.WhitelistRegistries,
			Usage: "Specify a list of permitted whitelist registry server addresses (insecure addresses still require the --insecure-registry option in addition)",
		},
	}
}
