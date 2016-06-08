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
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/vmware/vic/cmd/vic-machine/create"
)

var (
	MajorVersion string
	BuildID      string
)

func main() {
	app := cli.NewApp()

	app.Name = filepath.Base(os.Args[0])
	app.Usage = "Create and manage Virtual Container Hosts"
	app.EnableBashCompletion = true

	create := create.NewCreate()
	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Deploy VCH",
			Action: create.Run,
			Flags:  create.Flags(),
		},
	}
	app.Version = fmt.Sprintf("%s.%s", MajorVersion, BuildID)
	app.Run(os.Args)
}
