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
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/trace"
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
	CreateCommand   []string
	validator       *validate.Validator
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

func NewEngineInstaller() *EngineInstaller {
	return &EngineInstaller{Name: "default-vch"}
}

func (ei *EngineInstaller) populateConfigOptions() *EngineInstallerConfigOptions {
	defer trace.End(trace.Begin(""))

	vc := ei.validator.IsVC()
	log.Infof("Is VC: %t\n", vc)

	dcs, err := ei.validator.ListDatacenters()
	if err != nil {
		log.Infoln(err)
		return nil
	}
	for _, d := range dcs {
		log.Infof("DC: %s\n", d)
	}

	comp, err := ei.validator.ListComputeResource()
	if err != nil {
		log.Infoln(err)
		return nil
	}
	for _, c := range comp {
		log.Infof("compute: %s\n", c)
	}

	rp, err := ei.validator.ListResourcePool("*")
	if err != nil {
		log.Infoln(err)
		return nil
	}
	for _, p := range rp {
		log.Infof("rp: %s\n", p)
	}

	nets, err := ei.validator.ListNetworks(!vc) // set to false for vC
	if err != nil {
		log.Infoln(err)
		return nil
	}
	for _, n := range nets {
		log.Infof("net: %s\n", n)
	}

	dss, err := ei.validator.ListDatastores()
	if err != nil {
		log.Infoln(err)
		return nil
	}
	for _, d := range dss {
		log.Infof("ds: %s\n", d)
	}

	return &EngineInstallerConfigOptions{
		Networks:      nets,
		Datastores:    dss,
		ResourcePools: rp,
	}
}

func (ei *EngineInstaller) buildCreateCommand(binaryPath string) {
	defer trace.End(trace.Begin(""))

	var createCommand []string

	createCommand = append(createCommand, binaryPath+"/vic/vic-machine-linux")
	createCommand = append(createCommand, "create")
	createCommand = append(createCommand, "--no-tlsverify")
	createCommand = append(createCommand, []string{"--target", ei.Target}...)
	createCommand = append(createCommand, []string{"--user", ei.User}...)
	createCommand = append(createCommand, []string{"--password", ei.Password}...)
	createCommand = append(createCommand, []string{"--name", ei.Name}...)
	createCommand = append(createCommand, []string{"--public-network", ei.PublicNetwork}...)
	createCommand = append(createCommand, []string{"--bridge-network", ei.BridgeNetwork}...)
	createCommand = append(createCommand, []string{"--compute-resource", ei.ComputeResource}...)
	createCommand = append(createCommand, []string{"--image-store", ei.ImageStore}...)
	createCommand = append(createCommand, []string{"--thumbprint", ei.Thumbprint}...)

	ei.CreateCommand = createCommand
}

func (ei *EngineInstaller) verifyLogin() error {
	defer trace.End(trace.Begin(""))

	ctx := context.TODO()

	var u url.URL
	u.User = url.UserPassword(ei.User, ei.Password)
	u.Host = ei.Target
	u.Path = ""
	log.Infof("server URL: %v\n", u)

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
		log.Infof("validator: %s", err)
		return err
	}

	ei.validator = v

	return nil
}

func setupDefaultAdmiral(vchIP string) {
	defer trace.End(trace.Begin(""))

	admiral := "https://localhost:8282"
	client := &http.Client{}

	// test if admiral is available on https. OVA uses https, but local development uses http.
	_, err := client.Head(admiral)
	if err != nil {
		// https failed, use http
		admiral = "http://localhost:8282"
	}

	// validate vch host
	sslTrustPayload := fmt.Sprintf("{\"hostState\":{\"id\":\"%s\",\"address\":\"https://%s\",\"customProperties\":{\"__adapterDockerType\":\"API\",\"__containerHostType\":\"VCH\"}}}", vchIP, vchIP)
	sslTrustReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/resources/hosts?validate=true", admiral), bytes.NewBuffer([]byte(sslTrustPayload)))
	sslTrustResp, err := client.Do(sslTrustReq)
	if err != nil || sslTrustResp.StatusCode != http.StatusOK {
		log.Infof("error: %v\nresponse: %v\n", err, sslTrustResp)
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
			log.Infof("error: %v\nresponse: %v\n", err, sslCertResp)
			log.Infoln("Admiral cannot trust host certificate.")
			return
		}
	}

	// add host to admiral
	addHostReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/resources/hosts", admiral), bytes.NewBuffer([]byte(sslTrustPayload)))
	addHostResp, err := client.Do(addHostReq)
	if err != nil || addHostResp.StatusCode != http.StatusNoContent {
		log.Infof("error: %v\nresponse: %v\n", err, addHostResp)
		log.Infoln("Error adding host to Admiral.")
		return
	}

	log.Infoln("Host added to admiral.")
}
