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
	"context"
	"fmt"
	"path"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/pllib"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

func (d *Dispatcher) CreateVCH(conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) error {
	defer trace.End(trace.Begin(conf.Name))

	var err error

	if err = d.checkExistence(conf, settings); err != nil {
		return err
	}

	if err = d.createPool(conf, settings); err != nil {
		return err
	}

	if err = d.createBridgeNetwork(conf); err != nil {
		return err
	}

	appID, err := d.createAppliance(conf, settings)
	if err != nil {
		return errors.Errorf("Creating the appliance failed with %s. Exiting...", err)
	}
	d.applianceID = appID

	// TODO: create image store
	if err = d.createVolumeStores(conf); err != nil {
		return errors.Errorf("Exiting because we could not create volume stores due to error: %s", err)
	}

	if err = d.reconfigureAppliance(conf, settings); err != nil {
		return errors.Errorf("Reconfiguring the appliance failed with %s. Exiting...", err)
	}
	if err = d.uploadImages(settings.ImageFiles); err != nil {
		return errors.Errorf("Uploading images failed with %s. Exiting...", err)
	}

	return nil
}

func (d *Dispatcher) createPool(conf *config.VirtualContainerHostConfigSpec, settings *data.InstallerData) error {
	defer trace.End(trace.Begin(""))

	var err error

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
		} else {
			d.vchPool = d.vchVapp.ResourcePool
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

	return nil
}

func (d *Dispatcher) StartAppliance(ctx context.Context) error {
	defer trace.End(trace.Begin(""))

	var h interface{}
	var err error

	if h = d.pl.NewHandle(ctx, d.applianceID); h == nil {
		err = errors.Errorf("Unable to get handle %s", d.applianceID)
		return err
	}

	h, err = d.pl.ChangeState(ctx, h, pllib.Running)
	if err != nil {
		return errors.Errorf("Failed to set state %s. Exiting...", err)
	}
	err = d.pl.Commit(ctx, h)
	if err != nil {
		return errors.Errorf("Failed to power on appliance %s. Exiting...", err)
	}

	return nil
}

func (d *Dispatcher) uploadImages(files map[string]string) error {
	defer trace.End(trace.Begin(""))

	var wg sync.WaitGroup

	// upload the images
	log.Infof("Uploading images for container")

	wg.Add(len(files))
	results := make(chan error, len(files))
	for key, image := range files {
		go func(key string, image string) {
			defer wg.Done()

			log.Infof("\t%q", image)
			err := d.session.Datastore.UploadFile(d.ctx, image, path.Join(d.vmPathName, key), nil)
			if err != nil {
				log.Errorf("\t\tUpload failed for %q: %s", image, err)
				if d.force {
					log.Warnf("\t\tContinuing despite failures (due to --force option)")
					log.Warnf("\t\tNote: The VCH will not function without %q...", image)
				} else {
					results <- err
				}
			}
		}(key, image)
	}
	wg.Wait()
	close(results)

	for err := range results {
		return err
	}
	return nil
}
