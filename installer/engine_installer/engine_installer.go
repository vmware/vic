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

package main

import (
	"fmt"
	"html/template"
)

type EngineInstallerConfigOptions struct {
	Networks      []string
	Datastores    []string
	ResourcePools []string
}

type EnginerInstallerHTML struct {
	BridgeNetwork   template.HTML
	PublicNetwork   template.HTML
	ImageStore      template.HTML
	ComputeResource template.HTML
}

type EngineInstaller struct {
	BridgeNetwork   string
	PublicNetwork   string
	ImageStore      string
	ComputeResource string
	Target          string
	User            string
	Password        string
	Name            string
	CreateCommand   string
}

func NewEngineInstaller() *EngineInstaller {
	return &EngineInstaller{
		User:     "root",
		Password: "password",
		Name:     "dev",
		Target:   "192.168.1.1",
	}
}

func (ei *EngineInstaller) populateConfigOptions() *EngineInstallerConfigOptions {
	return &EngineInstallerConfigOptions{
		Networks:      []string{"1", "2", "3"},
		Datastores:    []string{"4", "5", "6"},
		ResourcePools: []string{"7", "8", "9"},
	}
}

func (ei *EngineInstaller) buildCreateCommand() {
	createCommand := "create"

	createCommand = fmt.Sprintf("%s --target=%s", createCommand, ei.Target)
	createCommand = fmt.Sprintf("%s --user=%s", createCommand, ei.User)
	createCommand = fmt.Sprintf("%s --password=%s", createCommand, ei.Password)
	createCommand = fmt.Sprintf("%s --name=%s", createCommand, ei.Name)
	createCommand = fmt.Sprintf("%s --public-network=%s", createCommand, ei.PublicNetwork)
	createCommand = fmt.Sprintf("%s --bridge-network=%s", createCommand, ei.BridgeNetwork)
	createCommand = fmt.Sprintf("%s --compute-resource=%s", createCommand, ei.ComputeResource)
	createCommand = fmt.Sprintf("%s --image-store=%s", createCommand, ei.ImageStore)

	ei.CreateCommand = createCommand
}
