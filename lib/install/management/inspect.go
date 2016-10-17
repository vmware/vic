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
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

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
		log.Errorf("Failed to get VM power state, service might not be available at this moment.")
	}
	if state != types.VirtualMachinePowerStatePoweredOn {
		err = errors.Errorf("VCH is not powered on, state %s", state)
		log.Errorf("%s", err)
		return err
	}

	clientIP := conf.ExecutorConfig.Networks["client"].Assigned.IP

	if ip.IsUnspecifiedIP(clientIP) {
		err = errors.Errorf("No client IP address assigned")
		log.Errorf("%s", err)
		return err
	}

	d.HostIP = clientIP.String()
	log.Debugf("IP address for client interface: %s", d.HostIP)
	if !conf.HostCertificate.IsNil() {
		d.VICAdminProto = "https"
		d.DockerPort = fmt.Sprintf("%d", opts.DefaultTLSHTTPPort)
	} else {
		d.VICAdminProto = "http"
		d.DockerPort = fmt.Sprintf("%d", opts.DefaultHTTPPort)
	}

	// try looking up preferred name, irrespective of CAs
	if cert, err := conf.HostCertificate.X509Certificate(); err == nil {
		name, _ := viableHostAddress([]net.IP{clientIP}, cert, conf.CertificateAuthorities)
		if name != "" {
			log.Debugf("Retrieved proposed name from host certificate: %q", name)
			log.Debugf("Assigning first name from set: %s", name)

			if name != d.HostIP {
				log.Infof("Using address from host certificate over allocated IP: %s", d.HostIP)
				// reassign
				d.HostIP = name
			}
		} else {
			log.Warnf("Unable to identify address acceptable to host certificate")
		}
	} else {
		log.Debugf("Failed to load host cert: %s", err)
	}

	d.ShowVCH(conf, "", "", "", "")
	return nil
}

func (d *Dispatcher) ShowVCH(conf *config.VirtualContainerHostConfigSpec, key string, cert string, cacert string, envfile string) {
	if d.sshEnabled {
		log.Infof("")
		log.Infof("SSH to appliance")
		log.Infof("ssh root@%s", d.HostIP)
	}

	log.Infof("")
	log.Infof("vic-admin portal:")
	log.Infof("%s://%s:2378", d.VICAdminProto, d.HostIP)
	log.Infof("")
	tls := ""

	dEnv := " "
	if !conf.HostCertificate.IsNil() {
		// if we're generating then there's no CA currently
		if len(conf.CertificateAuthorities) > 0 {
			// find the name to use
			if key != "" {
				tls = fmt.Sprintf(" --tlsverify --tlscacert=%q --tlscert=%q --tlskey=%q", cacert, cert, key)
			} else {
				tls = fmt.Sprintf(" --tlsverify ")
			}

			dEnv = fmt.Sprintf("%s DOCKER_TLS_VERIFY=1", dEnv)
			info, err := os.Stat(conf.ExecutorConfig.Name)
			if err == nil && info.IsDir() {
				if abs, err := filepath.Abs(info.Name()); err == nil {
					dEnv = fmt.Sprintf("%s DOCKER_CERT_PATH=%s", dEnv, abs)
				}
			}
		} else {
			tls = " --tls"
		}
	}

	dEnv = fmt.Sprintf("%s DOCKER_HOST=%s:%s", dEnv, d.HostIP, d.DockerPort)
	log.Infof("Docker environment variables:")
	log.Info(dEnv)
	log.Infof("")

	if envfile != "" {
		log.Infof("Environment saved in %s", envfile)
		ioutil.WriteFile(envfile, []byte(dEnv), 0644)
	}

	log.Infof("")
	log.Infof("Connect to docker:")
	log.Infof("docker -H %s:%s%s info", d.HostIP, d.DockerPort, tls)
}
