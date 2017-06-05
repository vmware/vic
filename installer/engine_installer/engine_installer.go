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
	"os"
	"strings"
)

var (
	engineInstaller = NewEngineInstaller()
)

// EngineInstallerConfigOptions contains resource options for selection by user in options.html
type EngineInstallerConfigOptions struct {
	Networks      []string
	Datastores    []string
	ResourcePools []string
}

// EnginerInstallerHTML contains html select boxes for use in options.html template
type EnginerInstallerHTML struct {
	BridgeNetwork   template.HTML
	PublicNetwork   template.HTML
	ImageStore      template.HTML
	ComputeResource template.HTML
}

//
type EngineInstaller struct {
	BridgeNetwork   string
	PublicNetwork   string
	ImageStore      string
	ComputeResource string
	Target          string
	User            string
	Password        string
	Name            string
	Thumbprint      string
	CreateCommand   string
}

// NewEngineInstaller returns a new EngineInstaller struct with empty parameters
func NewEngineInstaller() *EngineInstaller {
	return &EngineInstaller{}
}

func (ei *EngineInstaller) populateConfigOptions() *EngineInstallerConfigOptions {
	return &EngineInstallerConfigOptions{
		Networks:      []string{"1", "2", "3"},
		Datastores:    []string{"4", "5", "6"},
		ResourcePools: []string{"7", "8", "9"},
	}
}

func (ei *EngineInstaller) buildCreateCommand() {
	gopath := os.Getenv("GOPATH")
	var createCommand []string

	createCommand = append(createCommand, gopath+"/src/github.com/vmware/vic/bin/vic-machine-linux")
	createCommand = append(createCommand, "create")
	createCommand = append(createCommand, "--no-tlsverify")
	createCommand = append(createCommand, fmt.Sprintf("--target %s", ei.Target))
	createCommand = append(createCommand, fmt.Sprintf("--user %s", ei.User))
	createCommand = append(createCommand, fmt.Sprintf("--password %s", ei.Password))
	createCommand = append(createCommand, fmt.Sprintf("--name %s", ei.Name))
	createCommand = append(createCommand, fmt.Sprintf("--public-network %s", ei.PublicNetwork))
	createCommand = append(createCommand, fmt.Sprintf("--bridge-network %s", ei.BridgeNetwork))
	createCommand = append(createCommand, fmt.Sprintf("--compute-resource %s", ei.ComputeResource))
	createCommand = append(createCommand, fmt.Sprintf("--image-store %s", ei.ImageStore))
	createCommand = append(createCommand, fmt.Sprintf("--thumbprint %s", ei.Thumbprint))

	ei.CreateCommand = strings.Join(createCommand, " ")
}
