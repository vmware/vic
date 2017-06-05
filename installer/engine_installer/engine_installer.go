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
	"html/template"
	"strings"

	"context"
	"fmt"
	"net/url"

	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/validate"
)

// EngineInstallerConfigOptions contains resource options for selection by user in options.html
type EngineInstallerConfigOptions struct {
	Networks      []string
	Datastores    []string
	ResourcePools []string
}

// EngineInstaller contains all options to be used in the vic-machine create command
type EngineInstaller struct {
	BridgeNetwork   string `json:"bridge-net"`
	PublicNetwork   string `json:"public-net"`
	ImageStore      string `json:"img-store"`
	ComputeResource string `json:"compute"`
	Target          string `json:"target"`
	User            string `json:"user"`
	Password        string `json:"password"`
	Name            string `json:"name"`
	Thumbprint      string `json:"thumbprint"`
	CreateCommand   string
}

// AutHTML holds the invalid login variable
type AuthHTML struct {
	InvalidLogin bool
}

type ExecHTMLOptions struct {
	BridgeNetwork   template.HTML
	PublicNetwork   template.HTML
	ImageStore      template.HTML
	ComputeResource template.HTML
	Target          string
	User            string
	Password        string
	Name            string
	Thumbprint      string
	CreateCommand   string
}

var (
	validator *validate.Validator
)

// NewEngineInstaller returns a new EngineInstaller struct with empty parameters
func NewEngineInstaller() *EngineInstaller {
	return &EngineInstaller{Name: "defualt-vch"}
}

func (ei *EngineInstaller) populateConfigOptions() *EngineInstallerConfigOptions {

	vc := validator.IsVC()
	fmt.Printf("Is VC: %t\n", vc)

	dcs, err := validator.ListDatacenters()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, d := range dcs {
		fmt.Printf("DC: %s\n", d)
	}

	comp, err := validator.ListComputeResource()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, c := range comp {
		fmt.Printf("compute: %s\n", c)
	}

	rp, err := validator.ListResourcePool("*")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, p := range rp {
		fmt.Printf("rp: %s\n", p)
	}

	nets, err := validator.ListNetworks(!vc) // set to false for vC
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, n := range nets {
		fmt.Printf("net: %s\n", n)
	}

	dss, err := validator.ListDatastores()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, d := range dss {
		fmt.Printf("ds: %s\n", d)
	}

	return &EngineInstallerConfigOptions{
		Networks:      nets,
		Datastores:    dss,
		ResourcePools: rp,
	}
}

func (ei *EngineInstaller) buildCreateCommand(binaryPath string) {
	var createCommand []string

	createCommand = append(createCommand, binaryPath)
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

func (ei *EngineInstaller) verifyLogin() error {
	ctx := context.TODO()

	var u url.URL
	u.User = url.UserPassword(ei.User, ei.Password)
	u.Host = ei.Target
	u.Path = ""
	fmt.Printf("server URL: %v\n", u)

	input := data.NewData()

	input.OpsUser = u.User.Username()
	passwd, _ := u.User.Password()
	input.OpsPassword = &passwd
	input.URL = &u
	input.Force = true

	input.User = u.User.Username()
	input.Password = &passwd

	v, err := validate.NewValidator(ctx, input)
	if err != nil {
		fmt.Printf("validator: %s", err)
		return err
	}
	// if !v.IsVC() {
	// 	fmt.Printf("validator : %v\n", v)
	// 	return errors.New("target is not a vCenter instance")
	// }
	validator = v

	return nil
}
