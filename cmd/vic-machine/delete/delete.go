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
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"

	"golang.org/x/net/context"
)

// Delete has all input parameters for vic-machine delete command
type Uninstall struct {
	*data.Data

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
		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "Force the uninstall",
			Destination: &d.Force,
		},
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for appliance initialization",
			Destination: &d.Timeout,
		},
	}
	preFlags := append(d.TargetFlags(), d.IDFlags()...)
	preFlags = append(preFlags, d.ComputeFlags()...)
	flags = append(preFlags, flags...)
	flags = append(flags, d.DebugFlags()...)

	return flags
}

func (d *Uninstall) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := d.HasCredentials(); err != nil {
		return err
	}

	if err := d.ProcessID(); err != nil {
		return err
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
	if d.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	log.Infof("### Removing VCH ####")

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout)
	defer cancel()

	validator, err := validate.NewValidator(ctx, d.Data)
	if err != nil {
		err = errors.Errorf("%s. Exiting...", err)
		return err
	}
	executor := management.NewDispatcher(validator.Context, validator.Session, nil, d.Force)

	vch, _, err := executor.NewVCHFromComputePath(d.Data.ComputeResourcePath, d.Data.DisplayName)
	if err != nil {
		log.Errorf("Failed to get Virtual Container Host %s", d.DisplayName)
		return err
	}
	vchConfig, err := executor.GetVCHConfig(vch)
	if err != nil {
		log.Errorf("Failed to get Virtual Container Host configuration")
		return err
	}
	executor.InitDiagnosticLogs(vchConfig)

	if validator.IsVC() {
		log.Infoln("Removing VCH vSphere extension")
		if err = executor.GenerateExtensionName(vchConfig); err != nil {
			log.Warnf("Wasn't able to get extension name during VCH deletion. Failed with error: %s", err)
		}
		if err = executor.UnregisterExtension(vchConfig.ExtensionName); err != nil {
			log.Warnf("Wasn't able to remove extension %s due to error: %s", vchConfig.ExtensionName, err)
		}
	}

	if err = executor.DeleteVCH(vchConfig); err != nil {
		executor.CollectDiagnosticLogs()
		return err
	}

	log.Infof("Completed successfully")

	return nil
}
