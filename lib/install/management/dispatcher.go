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
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"strings"

	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/diagnostic"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"

	log "github.com/Sirupsen/logrus"
)

type Dispatcher struct {
	session *session.Session
	ctx     context.Context
	force   bool
	secret  *extraconfig.SecretKey

	isVC          bool
	vchPoolPath   string
	vmPathName    string
	dockertlsargs string

	DockerPort string
	HostIP     string

	vchPool   *object.ResourcePool
	vchVapp   *object.VirtualApp
	appliance *vm.VirtualMachine

	oldApplianceISO string

	sshEnabled         bool
	parentResourcepool *compute.ResourcePool
}

type diagnosticLog struct {
	key     string
	name    string
	start   int32
	host    *object.HostSystem
	collect bool
}

var diagnosticLogs = make(map[string]*diagnosticLog)

// NewDispatcher creates a dispatcher that can act upon VIC management operations.
// clientCert is an optional client certificate to allow interaction with the Docker API for verification
// force will ignore some errors
func NewDispatcher(ctx context.Context, s *session.Session, conf *config.VirtualContainerHostConfigSpec, force bool) *Dispatcher {
	defer trace.End(trace.Begin(""))
	isVC := s.IsVC()
	e := &Dispatcher{
		session: s,
		ctx:     ctx,
		isVC:    isVC,
		force:   force,
	}
	if conf != nil {
		e.InitDiagnosticLogs(conf)
	}
	return e
}

// Get the current log header LineEnd of the hostd/vpxd logs.
// With this we avoid collecting log file data that existed prior to install.
func (d *Dispatcher) InitDiagnosticLogs(conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	if d.isVC {
		diagnosticLogs[d.session.ServiceContent.About.InstanceUuid] =
			&diagnosticLog{"vpxd:vpxd.log", "vpxd.log", 0, nil, true}
	}

	var err error
	if d.session.Datastore == nil {
		if len(conf.ImageStores) == 0 {
			log.Errorf("Image datastore is empty")
			return

		}
		if d.session.Datastore, err = d.session.Finder.DatastoreOrDefault(d.ctx, conf.ImageStores[0].Host); err != nil {
			log.Errorf("Failure finding image store from VCH config (%s): %s", conf.ImageStores[0].Host, err.Error())
			return
		}
		log.Debugf("Found ds: %s", conf.ImageStores[0].Host)
	}
	// find the host(s) attached to given storage
	if d.session.Cluster == nil {
		if len(conf.ComputeResources) > 0 {
			rp := compute.NewResourcePool(d.ctx, d.session, conf.ComputeResources[0])
			if d.session.Cluster, err = rp.GetCluster(d.ctx); err != nil {
				log.Errorf("Unable to get cluster for given resource pool %s: %s", conf.ComputeResources[0], err)
				return
			}
		} else {
			log.Errorf("Compute resource is empty")
			return
		}
	}
	hosts, err := d.session.Datastore.AttachedClusterHosts(d.ctx, d.session.Cluster)
	if err != nil {
		log.Errorf("Unable to get the list of hosts attached to given storage: %s", err)
		return
	}

	if d.session.Host == nil {
		// vCenter w/ auto DRS.
		// Set collect=false here as we do not want to collect all hosts logs,
		// just the hostd log where the VM is placed.
		for _, host := range hosts {
			diagnosticLogs[host.Reference().Value] =
				&diagnosticLog{"hostd", "hostd.log", 0, host, false}
		}
	} else {
		// vCenter w/ manual DRS or standalone ESXi
		var host *object.HostSystem
		if d.isVC {
			host = d.session.Host
		}

		diagnosticLogs[d.session.Host.Reference().Value] =
			&diagnosticLog{"hostd", "hostd.log", 0, host, true}
	}

	m := diagnostic.NewDiagnosticManager(d.session)

	for k, l := range diagnosticLogs {
		// get LineEnd without any LineText
		h, err := m.BrowseLog(d.ctx, l.host, l.key, math.MaxInt32, 0)

		if err != nil {
			log.Warnf("Disabling %s %s collection (%s)", k, l.name, err)
			diagnosticLogs[k] = nil
			continue
		}

		l.start = h.LineEnd
	}
}

func (d *Dispatcher) CollectDiagnosticLogs() {
	defer trace.End(trace.Begin(""))

	m := diagnostic.NewDiagnosticManager(d.session)

	for k, l := range diagnosticLogs {
		if l == nil || !l.collect {
			continue
		}

		log.Infof("Collecting %s %s", k, l.name)

		var lines []string
		start := l.start

		for i := 0; i < 2; i++ {
			h, err := m.BrowseLog(d.ctx, l.host, l.key, start, 0)
			if err != nil {
				log.Errorf("Failed to collect %s %s: %s", k, l.name, err)
				break
			}

			lines = h.LineText
			if len(lines) != 0 {
				break // l.start was still valid, log was not rolled over
			}

			// log rolled over, start at the beginning.
			// TODO: If this actually happens we will have missed some log data,
			// it is possible to get data from the previous log too.
			start = 0
			log.Infof("%s %s rolled over", k, l.name)
		}

		if len(lines) == 0 {
			log.Warnf("No log data for %s %s", k, l.name)
			continue
		}

		f, err := os.Create(l.name)
		if err != nil {
			log.Errorf("Failed to create local %s: %s", l.name, err)
			continue
		}
		defer f.Close()

		for _, line := range lines {
			fmt.Fprintln(f, line)
		}
	}
}

func (d *Dispatcher) opManager(ctx context.Context, vch *vm.VirtualMachine) (*guest.ProcessManager, error) {
	state, err := vch.PowerState(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to get appliance power state, service might not be available at this moment.")
	}
	if state != types.VirtualMachinePowerStatePoweredOn {
		return nil, fmt.Errorf("VCH appliance is not powered on, state %s", state)
	}

	running, err := vch.IsToolsRunning(ctx)
	if err != nil || !running {
		return nil, errors.New("Tools are not running in the appliance, unable to continue")
	}

	manager := guest.NewOperationsManager(d.session.Client.Client, vch.Reference())
	processManager, err := manager.ProcessManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to manage processes in appliance VM: %s", err)
	}
	return processManager, nil
}

func (d *Dispatcher) CheckAccessToVCAPI(ctx context.Context, vch *vm.VirtualMachine, target string) (int64, error) {
	pm, err := d.opManager(ctx, vch)
	if err != nil {
		return -1, err
	}
	auth := types.NamePasswordAuthentication{}
	spec := types.GuestProgramSpec{
		ProgramPath:      "test-vc-api",
		Arguments:        target,
		WorkingDirectory: "/",
		EnvVariables:     []string{},
	}
	return pm.StartProgram(ctx, &auth, &spec)
}

// given a set of IP addresses this will determine what address, if any, can be used to
// connect to the host certificate
// if none can be found, will return empty string and an err
func addrToUse(candidateIPs []net.IP, cert *x509.Certificate, cas []byte) (string, error) {
	if cert == nil {
		return "", errors.New("unable to determine suitable address with nil certificate")
	}

	log.Debug("Loading CAs for client auth")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(cas)

	// update target to use FQDN
	for _, ip := range candidateIPs {
		names, err := net.LookupAddr(ip.String())
		if err != nil {
			log.Debugf("Unable to perform reverse lookup of IP address: %s", err)
		}

		// check all the returned names, and lastly the raw IP
		for _, n := range append(names, ip.String()) {
			opts := x509.VerifyOptions{
				Roots:   pool,
				DNSName: n,
			}

			_, err := cert.Verify(opts)
			if err == nil {
				// this identifier will work
				log.Debugf("Matched %s for use against host certificate", n)
				// trim '.' fqdn suffix if fqdn
				return strings.TrimSuffix(n, "."), nil
			}

			log.Debugf("Checked %s, no match for host cert", n)
		}
	}

	// no viable address
	return "", errors.New("unable to determine viable address")
}

/// viableHostAddresses attempts to determine which possibles addresses in the certificate
// are viable from the current location.
// This will return all IP addresses - it attempts to validate DNS names via resolution.
// This does NOT check connectivity
func viableHostAddress(candidateIPs []net.IP, cert *x509.Certificate, cas []byte) (string, error) {
	if cert == nil {
		return "", fmt.Errorf("unable to determine suitable address with nil certificate")
	}

	log.Debug("Loading CAs for client auth")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(cas)

	dnsnames := cert.DNSNames

	// assemble the common name and alt names
	ip := net.ParseIP(cert.Subject.CommonName)
	if ip != nil {
		candidateIPs = append(candidateIPs, ip)
	} else {
		// assume it's dns
		dnsnames = append([]string{cert.Subject.CommonName}, dnsnames...)
	}

	// turn the DNS names into IPs
	for _, n := range dnsnames {
		// see which resolve from here
		ips, _ := net.LookupIP(n)
		if len(ips) == 0 {
			log.Debugf("Discarding name from viable set: %s", n)
			continue

		}

		candidateIPs = append(candidateIPs, ips...)
	}

	// always add all the altname IPs - we're not checking for connectivity
	candidateIPs = append(candidateIPs, cert.IPAddresses...)

	return addrToUse(candidateIPs, cert, cas)
}
