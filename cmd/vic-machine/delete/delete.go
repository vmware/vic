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

package delete

import (
	"io"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/urfave/cli"
	"github.com/vmware/vic/cmd/vic-machine/data"
	"github.com/vmware/vic/cmd/vic-machine/validate"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/pkg/errors"

	"golang.org/x/net/context"
)

// Delete has all input parameters for vic-machine delete command
type Uninstall struct {
	*data.Data

	id      string
	logfile string

	executor *management.Dispatcher
}

func NewUninstall() *Uninstall {
	d := &Uninstall{}
	d.Data = data.NewData()
	return d
}

// Flags return all cli flags for delete
func (d *Uninstall) Flags() []cli.Flag {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:        "vch-id",
			Value:       "",
			Usage:       "The ID of the Virtual Container Host - not supported until vic-machine ls is ready",
			Destination: &d.id,
		},
		cli.StringFlag{
			Name:        "compute-resource",
			Value:       "",
			Usage:       "Compute resource path, e.g. /ha-datacenter/host/myCluster/Resources/myRP",
			Destination: &d.ComputeResourcePath,
		},
		cli.BoolFlag{
			Name:        "force",
			Usage:       "Force the uninstall",
			Destination: &d.Force,
		},
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for appliance initialization",
			Destination: &d.Timeout,
		},
		cli.StringFlag{
			Name:        "name",
			Value:       "",
			Usage:       "The name of the Virtual Container Host",
			Destination: &d.DisplayName,
		},
	}
	flags = append(d.TargetFlags(), flags...)
	return flags
}

func (d *Uninstall) processParams() error {
	if err := d.HasCredentials(); err != nil {
		return err
	}

	if (d.ComputeResourcePath == "" || d.DisplayName == "") && d.id == "" {
		return cli.NewExitError("--compute-resource, --name or --vch-id should be specified", 1)
	}

	d.logfile = "delete.log"
	d.Insecure = true
	return nil
}

func (d *Uninstall) Run(cli *cli.Context) error {
	var err error
	if err = d.processParams(); err != nil {
		return err
	}

	// Open log file
	f, err := os.OpenFile(d.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		err = errors.Errorf("Error opening logfile %s: %v", d.logfile, err)
		return err
	}
	defer f.Close()

	// Initiliaze logger with default TextFormatter
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	// SetOutput to io.MultiWriter so that we can log to stdout and a file
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	log.Infof("### Removing VCH ####")

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout)
	defer cancel()

	validator, err := validate.NewValidator(ctx, d.Data)
	if err != nil {
		err = errors.Errorf("%s. Exiting...", err)
		return err
	}

	vch, err := validator.GetVCH(d.Data)
	if err != nil {
		log.Errorf("Failed to get Virtual Container Host %s", d.DisplayName)
		return err
	}
	vchConfig, err := validator.GetVCHConfig(vch, d.DisplayName)
	if err != nil {
		log.Errorf("Failed to get Virtual Container Host configuration")
		return err
	}

	executor := management.NewDispatcher(validator.Context, validator.Session, vchConfig, false)
	if err = executor.DeleteVCH(vchConfig); err != nil {
		executor.CollectDiagnosticLogs()
		return err
	}

	log.Infof("Completed successfully")

	return nil
}
