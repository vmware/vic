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
	"path"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/tasks"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"

	"golang.org/x/net/context"
)

func (d *Dispatcher) CreateVCH(conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) error {
	defer trace.End(trace.Begin(conf.Name))

	var err error

	if err = d.checkExistence(conf, settings); err != nil {
		return err
	}

	if d.isVC && !settings.UseRP {
		if d.vchVapp, err = d.createVApp(conf, settings); err != nil {
			detail := fmt.Sprintf("Creating virtual app failed: %s", err)
			if !d.force {
				return errors.New(detail)
			}

			log.Error(detail)
			log.Errorf("Deploying vch under parent pool %q, (--force=true)", settings.ResourcePoolPath)
			d.vchPool = d.session.Pool
			conf.ComputeResources = append(conf.ComputeResources, d.vchPool.Reference())
		}
	} else {
		if d.vchPool, err = d.createResourcePool(conf, settings); err != nil {
			detail := fmt.Sprintf("Creating resource pool failed: %s", err)
			if !d.force {
				return errors.New(detail)
			}

			log.Error(detail)
			log.Errorf("Deploying vch under parent pool %q, (--force=true)", settings.ResourcePoolPath)
			d.vchPool = d.session.Pool
			conf.ComputeResources = append(conf.ComputeResources, d.vchPool.Reference())
		}
	}

	if err = d.createBridgeNetwork(conf); err != nil {
		return err
	}

	if err = d.createVolumeStores(conf); err != nil {
		return errors.Errorf("Exiting because we could not create volume stores due to error: %s", err)
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
	_, err = tasks.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
		return d.appliance.PowerOn(ctx)
	})

	if err != nil {
		return errors.Errorf("Failed to power on appliance %s. Exiting...", err)
	}

	if err = d.ensureApplianceInitializes(conf); err != nil {
		return errors.Errorf("%s. Exiting...", err)
	}

	// wait till the appliance components are fully initialized
	if err = d.ensureComponentsInitialize(conf); err != nil {
		return errors.Errorf("%s. Exiting...", err)
	}

	return nil
}

func (d *Dispatcher) uploadImages(files map[string]string) error {
	defer trace.End(trace.Begin(""))

	var err error
	var wg sync.WaitGroup

	// upload the images
	log.Infof("Uploading images for container")
	wg.Add(len(files))
	results := make(chan error, len(files))
	for key, image := range files {
		go func(key string, image string) {
			defer wg.Done()

			log.Infof("\t%q", image)
			err = d.session.Datastore.UploadFile(d.ctx, image, path.Join(d.vmPathName, key), nil)
			if err != nil {
				log.Errorf("\t\tUpload failed for %q: %s", image, err)
				if d.force {
					log.Warnf("\t\tContinuing despite failures (due to --force option)")
					log.Warnf("\t\tNote: The VCH will not function without %q...", image)
					results <- nil
				} else {
					results <- err
				}
				return
			}
			results <- nil
		}(key, image)
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

func (d *Dispatcher) RegisterExtension(conf *config.VirtualContainerHostConfigSpec, extension types.Extension) error {
	defer trace.End(trace.Begin(conf.ExtensionName))

	log.Infoln("Registering VCH as a vSphere extension")

	// vSphere confusingly calls the 'name' of the extension a 'key'
	// This variable is named IdKey as to not confuse it with its private key
	if conf.ExtensionCert == "" {
		return errors.Errorf("Extension certificate does not exist")
	}

	extensionManager := object.NewExtensionManager(d.session.Vim25())

	extension.LastHeartbeatTime = time.Now().UTC()
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
