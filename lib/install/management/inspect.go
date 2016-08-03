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

package management

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/opts"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

func (d *Dispatcher) InspectVCH(vch *vm.VirtualMachine, conf *config.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin(conf.Name))

	state, err := vch.PowerState(d.ctx)
	if err != nil {
		log.Errorf("Failed to get VM power state, service might not be avaialble at this moment.")
	}
	if state != types.VirtualMachinePowerStatePoweredOn {
		err = errors.Errorf("VCH is not powered on, state %s", state)
		log.Errorf("%s", err)
		return err
	}
	if ip.IsUnspecifiedIP(conf.ExecutorConfig.Networks["client"].Assigned.IP) {
		err = errors.Errorf("No client IP address assigned")
		log.Errorf("%s", err)
		return err
	}

	d.HostIP = conf.ExecutorConfig.Networks["client"].Assigned.IP.String()
	log.Debug("IP address for client interface: %s", d.HostIP)
	if !conf.HostCertificate.IsNil() {
		d.VICAdminProto = "https"
		d.DockerPort = fmt.Sprintf("%d", opts.DefaultTLSHTTPPort)
	} else {
		d.VICAdminProto = "http"
		d.DockerPort = fmt.Sprintf("%d", opts.DefaultHTTPPort)
	}
	d.ShowVCH(conf, "", "")
	return nil
}

func (d *Dispatcher) ShowVCH(conf *config.VirtualContainerHostConfigSpec, key string, cert string) {
	// #1218: Temporarily disable SSH access for TP3
	//	log.Infof("")
	//	log.Infof("SSH to appliance (default=root:password)")
	//	log.Infof("ssh root@%s", d.HostIP)
	log.Infof("")
	log.Infof("vic-admin portal:")
	log.Infof("%s://%s:2378", d.VICAdminProto, d.HostIP)
	log.Infof("")
	tls := ""

	if !conf.HostCertificate.IsNil() {
		// if we're generating then there's no CA currently
		if len(conf.CertificateAuthorities) > 0 && key != "" {
			tls = fmt.Sprintf(" --tls --tlscert='%s' --tlskey='%s'", cert, key)
		} else {
			tls = " --tls"
		}
	}
	log.Infof("DOCKER_HOST=%s:%s", d.HostIP, d.DockerPort)
	log.Infof("")
	log.Infof("Connect to docker:")
	log.Infof("docker -H %s:%s%s info", d.HostIP, d.DockerPort, tls)
}
