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
	"math"
	"os"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/diagnostic"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	"golang.org/x/net/context"
)

type Dispatcher struct {
	session *session.Session
	ctx     context.Context
	force   bool

	isVC          bool
	vchPoolPath   string
	vmPathName    string
	dockertlsargs string
	DockerPort    string
	HostIP        string
	VICAdminProto string

	vchPool   *compute.ResourcePool
	appliance *vm.VirtualMachine
}

type diagnosticLog struct {
	key     string
	name    string
	start   int32
	host    *object.HostSystem
	collect bool
}

var diagnosticLogs = make(map[string]*diagnosticLog)

func NewDispatcher(ctx context.Context, s *session.Session,
	conf *metadata.VirtualContainerHostConfigSpec, force bool) *Dispatcher {
	isVC := s.IsVC()
	e := &Dispatcher{
		session: s,
		ctx:     ctx,
		isVC:    isVC,
		force:   force,
	}
	e.initDiagnosticLogs(conf)
	return e
}

// Get the current log header LineEnd of the hostd/vpxd logs.
// With this we avoid collecting log file data that existed prior to install.
func (d *Dispatcher) initDiagnosticLogs(conf *metadata.VirtualContainerHostConfigSpec) {
	if d.isVC {
		diagnosticLogs[d.session.ServiceContent.About.InstanceUuid] =
			&diagnosticLog{"vpxd:vpxd.log", "vpxd.log", 0, nil, true}
	}

	// find the host(s) attached to given storage
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

func (d *Dispatcher) Dispatch(conf *metadata.VirtualContainerHostConfigSpec) error {
	var err error
	if d.vchPool, err = d.createResourcePool(conf); err != nil && !d.force {
		return errors.Errorf("Creating resource pool failed with %s. Exiting...", err)
	}

	if err = d.removeApplianceIfForced(conf); err != nil {
		return errors.Errorf("%s", err)
	}
	if err = d.createAppliance(conf); err != nil {
		return errors.Errorf("Creating the appliance failed with %s. Exiting...", err)
	}
	if err = d.uploadImages(conf); err != nil {
		return errors.Errorf("Uploading images failed with %s. Exiting...", err)
	}

	_, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return d.appliance.PowerOn(ctx)
	})
	if err != nil {
		return errors.Errorf("Failed to power on appliance %s. Exiting...", err)
	}

	if err = d.setMacToGuestInfo(conf); err != nil {
		return errors.Errorf("Failed to set Mac address %s. Exiting...", err)
	}
	if err = d.makeSureApplianceRuns(); err != nil {
		return errors.Errorf("%s. Exiting...", err)
	}
	return nil
}

func (d *Dispatcher) uploadImages(conf *metadata.VirtualContainerHostConfigSpec) error {
	var err error
	var wg sync.WaitGroup

	// upload the images
	log.Infof("Uploading images for container")
	wg.Add(len(conf.ImageFiles))
	results := make(chan error, len(conf.ImageFiles))
	for _, image := range conf.ImageFiles {
		go func(image string) {
			defer wg.Done()

			log.Infof("\t%s", image)
			base := filepath.Base(image)
			err = d.session.Datastore.UploadFile(d.ctx, image, d.vmPathName+"/"+base, nil)
			if err != nil {
				log.Errorf("\t\tUpload failed for %s", image)
				if d.force {
					log.Warnf("\t\tSkipping %s...", image)
					results <- nil
				} else {
					results <- err
				}
				return
			}
			results <- nil
		}(image)
	}
	wg.Wait()
	close(results)

	for err := range results {
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dispatcher) CollectDiagnosticLogs() {
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
