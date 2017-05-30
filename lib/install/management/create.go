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
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
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

	results := make(chan error)
	logMsgs := make(chan string)
	uploadCount := len(files)
	for key, image := range files {
		go func(key string, image string) {
			finalMessage := ""
			log.Infof("\t%q", image)

			// upload function that is passed to retry
			operationForRetry := func() error {
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
					logMsgs <- fmt.Sprintf("Attempted upload a total of %d times without success, Upload process failed.", uploadRetryLimit)
					return false
				}
				logMsgs <- fmt.Sprintf("failed an attempt to upload isos with err (%s), %d retries remain", err.Error(), retryCount)
				return true
			}

			// Build retry config
			backoffConf := &retry.BackoffConfig{
				InitialInterval:     uploadInitialInterval,
				RandomizationFactor: retry.DefaultRandomizationFactor,
				Multiplier:          retry.DefaultMultiplier,
				MaxInterval:         uploadMaxInterval,
				MaxElapsedTime:      uploadMaxElapsedTime,
			}

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
			results <- nil
		}(key, image)
	}

	var resultErrs []string
	for {
		select {
		case msg := <-logMsgs:
			log.Infof("%s", msg)

		case uploadResult := <-results:
			if uploadResult != nil {
				resultErrs = append(resultErrs, uploadResult.Error())
			}

			uploadCount--

			// once this is hit the upload attempts are over
			if uploadCount == 0 {
				if len(logMsgs) > 0 {
					log.Infof("%s", <-logMsgs)
				}

				for uploadErr := range resultErrs {
					log.Error(uploadErr)
				}

				close(logMsgs)
				close(results)

				if len(resultErrs) > 0 {
					return errors.New("Unable to upload isos.")
				}

				return nil
			}
		}
	}

}
