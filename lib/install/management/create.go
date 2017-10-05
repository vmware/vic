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
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/vchlog"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/retry"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/tasks"
)

const (
	uploadRetryLimit      = 5
	uploadMaxElapsedTime  = 30 * time.Minute
	uploadMaxInterval     = 1 * time.Minute
	uploadInitialInterval = 10 * time.Second
	timeFormat            = "2006-01-02T15:04:05-0700"
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

	if err = d.createAppliance(conf, settings); err != nil {
		return errors.Errorf("Creating the appliance failed with %s. Exiting...", err)
	}

	// send the signal to VCH logger to indicate VCH datastore path is ready
	datastoreReadySignal := vchlog.DatastoreReadySignal{
		Datastore:   d.session.Datastore,
		Name:        "create",
		Operation:   trace.NewOperation(d.ctx, "create"),
		VMPathName:  d.vmPathName,
		Timestamp:   time.Now().UTC().Format(timeFormat),
	}
	vchlog.Signal(datastoreReadySignal)

	if err = d.uploadImages(settings.ImageFiles); err != nil {
		return errors.Errorf("Uploading images failed with %s. Exiting...", err)
	}

	return d.startAppliance(conf)
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
		}
	} else {
		if d.vchPool, err = d.createResourcePool(conf, settings); err != nil {
			detail := fmt.Sprintf("Creating resource pool failed: %s", err)
			return errors.New(detail)
		}
	}

	return nil
}

func (d *Dispatcher) startAppliance(conf *config.VirtualContainerHostConfigSpec) error {
	defer trace.End(trace.Begin(""))

	var err error
	_, err = d.appliance.WaitForResult(d.ctx, func(ctx context.Context) (tasks.Task, error) {
		return d.appliance.PowerOn(ctx)
	})

	if err != nil {
		return errors.Errorf("Failed to power on appliance %s. Exiting...", err)
	}

	return nil
}

func (d *Dispatcher) uploadImages(files map[string]string) error {
	defer trace.End(trace.Begin(""))

	// upload the images
	log.Infof("Uploading images for container")

	results := make(chan error, len(files))
	var wg sync.WaitGroup

	for key, image := range files {

		wg.Add(1)
		go func(key string, image string) {
			finalMessage := ""
			log.Infof("\t%q", image)

			// upload function that is passed to retry
			operationForRetry := func() error {
				// attempt to delete the iso image first in case of failed upload
				dc := d.session.Datacenter
				fm := d.session.Datastore.NewFileManager(dc, false)
				ds := d.session.Datastore

				isoTargetPath := path.Join(d.vmPathName, key)
				// check iso first
				_, err := ds.Stat(d.ctx, isoTargetPath)
				if err != nil {
					switch err.(type) {
					// if not found, do nothing
					case object.DatastoreNoSuchFileError:
						// otherwise force delete
					default:
						log.Debugf("target delete path = %s", isoTargetPath)
						err := fm.Delete(d.ctx, isoTargetPath)
						if err != nil {
							log.Debugf("Failed to delete image (%s) with error (%s)", image, err.Error())
							return err
						}
					}
				}

				return d.session.Datastore.UploadFile(d.ctx, image, path.Join(d.vmPathName, key), nil)
			}

			// counter for retry decider
			retryCount := uploadRetryLimit

			// decider for our retry, will retry the upload uploadRetryLimit times
			uploadRetryDecider := func(err error) bool {
				if err == nil {
					return false
				}

				retryCount--
				if retryCount < 0 {
					log.Warnf("Attempted upload a total of %d times without success, Upload process failed.", uploadRetryLimit)
					return false
				}
				log.Warnf("failed an attempt to upload isos with err (%s), %d retries remain", err.Error(), retryCount)
				return true
			}

			// Build retry config
			backoffConf := retry.NewBackoffConfig()
			backoffConf.InitialInterval = uploadInitialInterval
			backoffConf.MaxInterval = uploadMaxInterval
			backoffConf.MaxElapsedTime = uploadMaxElapsedTime

			uploadErr := retry.DoWithConfig(operationForRetry, uploadRetryDecider, backoffConf)
			if uploadErr != nil {
				finalMessage = fmt.Sprintf("\t\tUpload failed for %q: %s\n", image, uploadErr)
				if d.force {
					finalMessage = fmt.Sprintf("%s\t\tContinuing despite failures (due to --force option)\n", finalMessage)
					finalMessage = fmt.Sprintf("%s\t\tNote: The VCH will not function without %q...", finalMessage, image)
					results <- errors.New(finalMessage)
				} else {
					results <- errors.New(finalMessage)
				}
			}
			wg.Done()
		}(key, image)
	}

	wg.Wait()
	close(results)

	uploadFailed := false
	for err := range results {
		if err != nil {
			log.Error(err.Error())
			uploadFailed = true
		}
	}

	if uploadFailed {
		return errors.New("Failed to upload iso images successfully.")
	}
	return nil
}
