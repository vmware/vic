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

package vicadmin

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	// "github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

type Validator struct {
	Hostname       string
	Version        string
	FirewallStatus template.HTML
	FirewallIssues template.HTML
	HostIP         string
	DockerPort     string
}

const (
	GoodStatus = template.HTML(`<span class="right"><i class="fa fa-check"></i></span>`)
	BadStatus  = template.HTML(`<span class="right warning"><i class="fa fa-exclamation-triangle"></i></span>`)
)

func NewValidator(ctx context.Context, vch *config.VirtualContainerHostConfigSpec, sess *session.Session) *Validator {
	defer trace.End(trace.Begin(""))
	log.Infof("Creating new validator")
	v := &Validator{}
	v.Version = vch.Version
	log.Info(fmt.Sprintf("Setting version to %s", v.Version))

	//VCH Name
	v.Hostname, _ = os.Hostname()
	v.Hostname = strings.Title(v.Hostname)
	log.Info(fmt.Sprintf("Setting hostname to %s", v.Hostname))


	//Retrieve Host IP Information and Set socker Endpoint
	v.HostIP = vch.ExecutorConfig.Networks["client"].Assigned.IP.String()

	if !vch.HostCertificate.IsNil() {
		v.DockerPort = "2376"
	} else {
		v.DockerPort = "2375"
	}

	v2, _ := validate.CreateFromVCHConfig(ctx, vch, sess)
	v2.CheckFirewall(ctx)
	firewallIssues := v2.GetIssues()

	if len(firewallIssues) == 0 {
		v.FirewallStatus = GoodStatus
		v.FirewallIssues = template.HTML("")
	} else {
		v.FirewallStatus = BadStatus
		for _, err := range firewallIssues {
			v.FirewallIssues = template.HTML(fmt.Sprintf("%s<span class=\"error-message\">%s</span>\n", v.FirewallIssues, err))
		}
	}
	log.Info(fmt.Sprintf("FirewallStatus set to: %s", v.FirewallStatus))
	log.Info(fmt.Sprintf("FirewallIssues set to: %s", v.FirewallIssues))

	return v
}
