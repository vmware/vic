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

package common

import "github.com/urfave/cli"

type Compute struct {
	ComputeResourcePath string
	DisplayName         string
}

func (c *Compute) ComputeFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "compute-resource, r",
			Value:       "",
			Usage:       "Compute resource path, e.g. myCluster/Resources/myRP. Default to <default cluster>/Resources",
			Destination: &c.ComputeResourcePath,
		},
		cli.StringFlag{
			Name:        "name, n",
			Value:       "docker-appliance",
			Usage:       "The name of the Virtual Container Host",
			Destination: &c.DisplayName,
		},
	}
}
