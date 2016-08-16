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
	"net"
	"os"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	// "github.com/vmware/govmomi/vim25/types"

	"github.com/docker/docker/opts"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

type Validator struct {
	Hostname         string
	Version          string
	FirewallStatus   template.HTML
	FirewallIssues   template.HTML
	LicenseStatus    template.HTML
	LicenseIssues    template.HTML
	NetworkStatus    template.HTML
	NetworkIssues    template.HTML
	StorageRemaining template.HTML
	HostIP           string
	DockerPort       string
}

const (
	GoodStatus = template.HTML(`<i class="icon-ok"></i>`)
	BadStatus  = template.HTML(`<i class="icon-attention"></i>`)
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

	//Firewall Status Check
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

	//License Check
	v2.ClearIssues()
	v2.CheckLicense(ctx)
	licenseIssues := v2.GetIssues()

	if len(licenseIssues) == 0 {
		v.LicenseStatus = GoodStatus
		v.LicenseIssues = template.HTML("")
	} else {
		v.LicenseStatus = BadStatus
		for _, err := range licenseIssues {
			v.LicenseIssues = template.HTML(fmt.Sprintf("%s<span class=\"error-message\">%s</span>\n", v.LicenseIssues, err))
		}
	}

	//Network Connection Check
	hosts := []string{
		"google.com:80",
		"docker.io:443",
	}
	nwErrors := []error{}

	for _, host := range hosts {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			nwErrors = append(nwErrors, err)
		} else {
			conn.Close()
		}
	}

	if len(nwErrors) > 0 {
		v.NetworkStatus = BadStatus
		for _, err := range nwErrors {
			v.NetworkIssues = template.HTML(fmt.Sprintf("%s<span class=\"error-message\">%s</span>\n", v.NetworkIssues, err))
		}
	} else {
		v.NetworkStatus = GoodStatus
		v.NetworkIssues = template.HTML("")

	}

	//Retrieve Host IP Information and Set Docker Endpoint
	v.HostIP = vch.ExecutorConfig.Networks["client"].Assigned.IP.String()

	if vch.HostCertificate.IsNil() {
		v.DockerPort = fmt.Sprintf("%d", opts.DefaultHTTPPort)
	} else {
		v.DockerPort = fmt.Sprintf("%d", opts.DefaultTLSHTTPPort)
	}

	v.QueryDatastore(ctx, vch, sess)
	return v
}

type dsList []mo.Datastore

func (d dsList) Len() int { return len(d) }
func (d dsList) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d dsList) Less(i, j int) bool { return d[i].Name < d[j].Name }

func (v *Validator) QueryDatastore(ctx context.Context, vch *config.VirtualContainerHostConfigSpec, sess *session.Session) {
	var dataStores dsList
	dsNames := make(map[string]bool)

	for _, url := range vch.ImageStores {
		dsNames[url.Host] = true
	}

	for _, url := range vch.VolumeLocations {
		dsNames[url.Host] = true
	}

	for _, url := range vch.ContainerStores {
		dsNames[url.Host] = true
	}

	refs := []types.ManagedObjectReference{}
	for dsName, _ := range dsNames {
		ds, err := sess.Finder.DatastoreOrDefault(ctx, dsName)
		if err != nil {
			log.Errorf("Unable to collect information for datastore %s: %s", dsName, err)
		} else {
			refs = append(refs, ds.Reference())
		}
	}

	pc := property.DefaultCollector(sess.Client.Client)
	err := pc.Retrieve(ctx, refs, nil, &dataStores)

	sort.Sort(dataStores)
	if err != nil {
		log.Errorf("Error while accessing datastore: %s", err)
		return
	}
	for _, ds := range dataStores {
		log.Infof("Datastore %s Status: %s", ds.Name, ds.OverallStatus)
		log.Infof("Datastore %s Free Space: %.1fGB", ds.Name, float64(ds.Summary.FreeSpace)/(1<<30))
		log.Infof("Datastore %s Capacity: %.1fGB", ds.Name, float64(ds.Summary.Capacity)/(1<<30))

		v.StorageRemaining = template.HTML(fmt.Sprintf(`%s
			<div class="row card-text">
			  <div class="sixty">%s:</div>
			  <div class="fourty">%.1f GB remaining</div>
			</div>`, v.StorageRemaining, ds.Name, float64(ds.Summary.FreeSpace)/(1<<30)))
	}
}
