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
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/validate"
)

// EngineInstallerConfigOptions contains resource options for selection by user in exec.html
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

// AuthHTML holds the invalid login variable
type AuthHTML struct {
	InvalidLogin bool
}

// ExecHTMLOptions contains fields for html templating in exec.html
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
	return &EngineInstaller{Name: "default-vch"}
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

	createCommand = append(createCommand, binaryPath+"/vic/vic-machine-linux")
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

	validator = v

	return nil
}

func setupDefaultAdmiral(vchIp string) {
	admiral := "https://localhost:8282"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// validate vch host
	sslTrustPayload := fmt.Sprintf("{\"hostState\":{\"id\":\"%s\",\"address\":\"https://%s\",\"customProperties\":{\"__adapterDockerType\":\"API\",\"__containerHostType\":\"VCH\"}}}", vchIp, vchIp)
	sslTrustReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/resources/hosts?validate=true", admiral), bytes.NewBuffer([]byte(sslTrustPayload)))
	sslTrustResp, err := client.Do(sslTrustReq)
	if err != nil || sslTrustResp.StatusCode != http.StatusOK {
		log.Infoln(err, sslTrustResp.StatusCode)
		log.Infoln("Cannot add vch to Admiral.")
		return
	}

	// trust vch host on admiral
	sslCert, err := ioutil.ReadAll(sslTrustResp.Body)
	if err != nil {
		log.Infoln(err)
		log.Infoln("Cannot add vch to Admiral.")
		return
	}
	if len(sslCert) > 0 {
		sslCertReq, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/config/trust-certs", admiral), bytes.NewBuffer([]byte(sslCert)))
		sslCertResp, err := client.Do(sslCertReq)
		if err != nil || sslCertResp.StatusCode != http.StatusOK {
			log.Infoln(err, sslCertResp.StatusCode)
			log.Infoln("Admiral cannot trust host certificate.")
			return
		}
	}

	// add host to admiral
	addHostReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/resources/hosts", admiral), bytes.NewBuffer([]byte(sslTrustPayload)))
	addHostResp, err := client.Do(addHostReq)
	if err != nil || addHostResp.StatusCode != http.StatusNoContent {
		log.Infoln(err, addHostResp.StatusCode)
		log.Infoln("Error adding host to Admiral.")
		return
	}

	log.Infoln("Host added to admiral.")
}
