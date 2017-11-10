// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"os/user"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/opts"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/constants"
	"github.com/vmware/vic/pkg/certificate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

func (d *Dispatcher) InspectVCH(vch *vm.VirtualMachine, conf *config.VirtualContainerHostConfigSpec, certPath string) error {
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

	var clientIP net.IP
	var publicIP net.IP

	clientNet := conf.ExecutorConfig.Networks["client"]
	if clientNet != nil {
		clientIP = clientNet.Assigned.IP
	}
	publicNet := conf.ExecutorConfig.Networks["public"]
	if publicNet != nil {
		publicIP = publicNet.Assigned.IP
	}

	if ip.IsUnspecifiedIP(clientIP) {
		err = errors.Errorf("No client IP address assigned")
		log.Errorf("%s", err)
		return err
	}

	if ip.IsUnspecifiedIP(publicIP) {
		err = errors.Errorf("No public IP address assigned")
		log.Errorf("%s", err)
		return err
	}

	d.HostIP = clientIP.String()
	log.Debugf("IP address for client interface: %s", d.HostIP)
	if !conf.HostCertificate.IsNil() {
		d.DockerPort = fmt.Sprintf("%d", opts.DefaultTLSHTTPPort)
	} else {
		d.DockerPort = fmt.Sprintf("%d", opts.DefaultHTTPPort)
	}

	// try looking up preferred name, irrespective of CAs
	if cert, err := conf.HostCertificate.X509Certificate(); err == nil {
		// #nosec: Errors unhandled.
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
		log.Debugf("No host certificates provided")
	}

	// Check for valid client cert for a tls-verify configuration
	if len(conf.CertificateAuthorities) > 0 {
		possibleCertPaths := findCertPaths(conf.Name, certPath)

		// Check if a valid client cert exists in one of possibleCertPaths
		certPath = ""
		for _, path := range possibleCertPaths {
			certFile := filepath.Join(path, certificate.ClientCert)
			keyFile := filepath.Join(path, certificate.ClientKey)
			ckp := certificate.NewKeyPair(certFile, keyFile, nil, nil)
			if err = ckp.LoadCertificate(); err != nil {
				log.Debugf("Unable to load client cert in %s: %s", path, err)
				continue
			}

			if _, err = certificate.VerifyClientCert(conf.CertificateAuthorities, ckp); err != nil {
				log.Debugf(err.Error())
				continue
			}

			certPath = path
			break
		}
	}

	d.ShowVCH(conf, "", "", "", "", certPath)
	return nil
}

// findCertPaths returns candidate paths for client certs depending on whether
// a certPath was specified in the CLI.
func findCertPaths(vchName, certPath string) []string {
	var possibleCertPaths []string

	if certPath != "" {
		log.Infof("--tls-cert-path supplied - only checking for certs in %s/", certPath)
		possibleCertPaths = append(possibleCertPaths, certPath)
		return possibleCertPaths
	}

	possibleCertPaths = append(possibleCertPaths, vchName, ".")
	logMsg := fmt.Sprintf("--tls-cert-path not supplied - checking for certs in current directory, %s/", vchName)

	dockerConfPath := ""
	user, err := user.Current()
	if err == nil {
		dockerConfPath = filepath.Join(user.HomeDir, ".docker")
		possibleCertPaths = append(possibleCertPaths, dockerConfPath)
		logMsg = fmt.Sprintf("%s and %s/", logMsg, dockerConfPath)
	}
	log.Infof(logMsg)

	return possibleCertPaths
}

func (d *Dispatcher) ShowVCH(conf *config.VirtualContainerHostConfigSpec, key string, cert string, cacert string, envfile string, certpath string) {
	if d.sshEnabled {
		log.Infof("")
		log.Infof("SSH to appliance:")
		log.Infof("ssh root@%s", d.HostIP)
	}

	log.Infof("")
	log.Infof("VCH Admin Portal:")
	log.Infof("https://%s:%d", d.HostIP, constants.VchAdminPortalPort)

	log.Infof("")
	publicIP := conf.ExecutorConfig.Networks["public"].Assigned.IP
	log.Infof("Published ports can be reached at:")
	log.Infof("%s", publicIP.String())

	cmd, env := d.GetDockerAPICommand(conf, key, cert, cacert, certpath)

	log.Info("")
	log.Infof("Docker environment variables:")
	log.Info(env)

	if envfile != "" {
		if err := ioutil.WriteFile(envfile, []byte(env), 0644); err == nil {
			log.Infof("")
			log.Infof("Environment saved in %s", envfile)
		}
	}

	log.Infof("")
	log.Infof("Connect to docker:")
	log.Infof(cmd)
}

// GetDockerAPICommand generates values to display for usage of a deployed VCH
func (d *Dispatcher) GetDockerAPICommand(conf *config.VirtualContainerHostConfigSpec, key string, cert string, cacert string, certpath string) (cmd, env string) {
	var dEnv []string
	tls := ""

	if d.HostIP == "" {
		return "", ""
	}

	if !conf.HostCertificate.IsNil() {
		// if we're generating then there's no CA currently
		if len(conf.CertificateAuthorities) > 0 {
			// find the name to use
			if key != "" {
				tls = fmt.Sprintf(" --tlsverify --tlscacert=%q --tlscert=%q --tlskey=%q", cacert, cert, key)
			} else {
				tls = fmt.Sprintf(" --tlsverify")
			}

			dEnv = append(dEnv, "DOCKER_TLS_VERIFY=1")
			info, err := os.Stat(certpath)
			if err == nil && info.IsDir() {
				if abs, err := filepath.Abs(info.Name()); err == nil {
					dEnv = append(dEnv, fmt.Sprintf("DOCKER_CERT_PATH=%s", abs))
				}
			} else {
				log.Infof("")
				log.Warnf("Unable to find valid client certs")
				log.Warnf("DOCKER_CERT_PATH must be provided in environment or certificates specified individually via CLI arguments")
			}
		} else {
			tls = " --tls"
		}
	}
	dEnv = append(dEnv, fmt.Sprintf("DOCKER_HOST=%s:%s", d.HostIP, d.DockerPort))

	cmd = fmt.Sprintf("docker -H %s:%s%s info", d.HostIP, d.DockerPort, tls)
	env = strings.Join(dEnv, " ")

	return cmd, env
}
