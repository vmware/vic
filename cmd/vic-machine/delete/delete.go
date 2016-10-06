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
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/urfave/cli"

	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/vm"

	"golang.org/x/net/context"
)

// Delete has all input parameters for vic-machine delete command
type Uninstall struct {
	*data.Data

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
			Usage:       "Time to wait for delete",
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

	d.Insecure = true
	return nil
}

func (d *Uninstall) Run(cli *cli.Context) (err error) {
	if err = d.processParams(); err != nil {
		return err
	}

	if d.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if len(cli.Args()) > 0 {
		log.Errorf("Unknown argument: %s", cli.Args()[0])
		return errors.New("invalid CLI arguments")
	}

	log.Infof("### Removing VCH ####")

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout)
	defer cancel()
	defer func() {
		if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
			//context deadline exceeded, replace returned error message
			err = errors.Errorf("Delete timed out: use --timeout to add more time")
		}
	}()

	validator, err := validate.NewValidator(ctx, d.Data)
	if err != nil {
		log.Errorf("Delete cannot continue - failed to create validator: %s", err)
		return errors.New("delete failed")
	}
	executor := management.NewDispatcher(validator.Context, validator.Session, nil, d.Force)

	var vch *vm.VirtualMachine
	if d.Data.ID != "" {
		vch, err = executor.NewVCHFromID(d.Data.ID)
	} else {
		vch, err = executor.NewVCHFromComputePath(d.Data.ComputeResourcePath, d.Data.DisplayName, validator)
	}
	if err != nil {
		log.Errorf("Failed to get Virtual Container Host %s", d.DisplayName)
		log.Error(err)
		return errors.New("delete failed")
	}

	log.Infof("")
	log.Infof("VCH ID: %s", vch.Reference().String())

	vchConfig, err := executor.GetVCHConfig(vch)
	if err != nil {
		log.Error("Failed to get Virtual Container Host configuration")
		log.Error(err)
		return errors.New("delete failed")
	}
	executor.InitDiagnosticLogs(vchConfig)

	if err = executor.DeleteVCH(vchConfig); err != nil {
		executor.CollectDiagnosticLogs()
		log.Errorf("%s", err)
		return errors.New("delete failed")
	}

	log.Infof("Completed successfully")

	return nil
}
