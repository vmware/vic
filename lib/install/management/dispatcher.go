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
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/vsphere/compute"
	"github.com/vmware/vic/pkg/vsphere/diagnostic"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/tasks"
	"github.com/vmware/vic/pkg/vsphere/vm"

	"github.com/vmware/govmomi/vim25/types"
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

	vchPool   *object.ResourcePool
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
	if conf != nil {
		e.InitDiagnosticLogs(conf)
	}
	return e
}

// Get the current log header LineEnd of the hostd/vpxd logs.
// With this we avoid collecting log file data that existed prior to install.
func (d *Dispatcher) InitDiagnosticLogs(conf *metadata.VirtualContainerHostConfigSpec) {
	if d.isVC {
		diagnosticLogs[d.session.ServiceContent.About.InstanceUuid] =
			&diagnosticLog{"vpxd:vpxd.log", "vpxd.log", 0, nil, true}
	}

	var err error
	if d.session.Datastore == nil {
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

func (d *Dispatcher) RegisterExtension(conf *metadata.VirtualContainerHostConfigSpec, extension types.Extension) error {
	log.Infoln("Registering VCH as a vSphere extension")

	// vSphere confusingly calls the 'name' of the extension a 'key'
	// This variable is named IdKey as to not confuse it with its private key
	if conf.ExtensionCert == "" {
		return errors.Errorf("Extension certificate does not exist")
	}

	extensionManager := object.NewExtensionManager(d.session.Vim25())

	if err := extensionManager.Register(d.ctx, extension); err != nil {
		log.Errorf("Could not register the vSphere extension due to err: %s", err)
		return err
	}

	if err := extensionManager.SetCertificate(d.ctx, conf.ExtensionName, conf.ExtensionCert); err != nil {
		log.Errorf("Could not set the certificate on the vSphere extension due to error: %s", err)
		return err
	}

	return nil
}

func (d *Dispatcher) UnregisterExtension(name string) error {
	extensionManager := object.NewExtensionManager(d.session.Vim25())
	if err := extensionManager.Unregister(d.ctx, name); err != nil {
		return errors.Errorf("Failed to remove extension w/ name %s due to error: %s", name, err)
	}
	return nil
}

func (d *Dispatcher) Dispatch(conf *metadata.VirtualContainerHostConfigSpec, settings *data.InstallerData) error {
	var err error
	if d.vchPool, err = d.createResourcePool(conf, settings); err != nil {
		detail := fmt.Sprintf("Creating resource pool failed: %s", err)
		if d.force {
			return errors.New(detail)
		}

		log.Error(detail)
	}

	if err = d.createBridgeNetwork(conf); err != nil {
		return err
	}

	if err = d.checkExistence(conf); err != nil {
		return err
	}

	if err = d.createAppliance(conf, settings); err != nil {
		return errors.Errorf("Creating the appliance failed with %s. Exiting...", err)
	}

	if err = d.uploadImages(settings.ImageFiles); err != nil {
		return errors.Errorf("Uploading images failed with %s. Exiting...", err)
	}

	if d.session.IsVC() {
		if err = d.RegisterExtension(conf, settings.Extension); err != nil {
			return errors.Errorf("Error registering VCH vSphere extension: %s", err)
		}
	}
	_, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.ResultWaiter, error) {
		return d.appliance.PowerOn(ctx)
	})

	if err != nil {
		return errors.Errorf("Failed to power on appliance %s. Exiting...", err)
	}

	if err = d.makeSureApplianceRuns(conf); err != nil {
		return errors.Errorf("%s. Exiting...", err)
	}
	return nil
}

func (d *Dispatcher) uploadImages(files []string) error {
	var err error
	var wg sync.WaitGroup

	// upload the images
	log.Infof("Uploading images for container")
	wg.Add(len(files))
	results := make(chan error, len(files))
	for _, image := range files {
		go func(image string) {
			defer wg.Done()

			log.Infof("\t%s", image)
			base := filepath.Base(image)
			err = d.session.Datastore.UploadFile(d.ctx, image, d.vmPathName+"/"+base, nil)
			if err != nil {
				log.Errorf("\t\tUpload failed for %s, %s", image, err)
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
