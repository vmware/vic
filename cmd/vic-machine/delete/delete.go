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

package delete

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/vm"

	"context"
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
	util := []cli.Flag{
		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "Force the deletion",
			Destination: &d.Force,
		},
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for delete",
			Destination: &d.Timeout,
		},
	}

	target := d.TargetFlags()
	id := d.IDFlags()
	compute := d.ComputeFlags()
	debug := d.DebugFlags(true)

	// flag arrays are declared, now combined
	var flags []cli.Flag
	for _, f := range [][]cli.Flag{target, id, compute, util, debug} {
		flags = append(flags, f...)
	}

	return flags
}

func (d *Uninstall) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := d.HasCredentials(); err != nil {
		return err
	}

	return nil
}

func (d *Uninstall) Run(clic *cli.Context) (err error) {
	// urfave/cli will print out exit in error handling, so no more information in main method can be printed out.
	defer func() {
		err = common.LogErrorIfAny(clic, err)
	}()

	if err = d.processParams(); err != nil {
		return err
	}

	if d.Debug.Debug != nil && *d.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if len(clic.Args()) > 0 {
		log.Errorf("Unknown argument: %s", clic.Args()[0])
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
	defer validator.Session.Logout(ctx)
	_, err = validator.ValidateTarget(ctx, d.Data)
	if err != nil {
		log.Errorf("Delete cannot continue - target validation failed: %s", err)
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

	vchConfig, err := executor.GetNoSecretVCHConfig(vch)
	if err != nil {
		log.Error("Failed to get Virtual Container Host configuration")
		log.Error(err)
		return errors.New("delete failed")
	}

	// compare vch version and vic-machine version
	installerBuild := version.GetBuild()
	if vchConfig.Version == nil || !installerBuild.Equal(vchConfig.Version) {
		if !d.Data.Force {
			log.Errorf("VCH version %q is different than installer version %s. Upgrade VCH before deleting or specify --force to force delete", vchConfig.Version.ShortVersion(), installerBuild.ShortVersion())
			return errors.New("delete failed")
		}

		log.Warnf("VCH version %q is different than installer version %s. Force delete will attempt to remove everything related to the installed VCH", vchConfig.Version.ShortVersion(), installerBuild.ShortVersion())
	}

	if err = executor.DeleteVCH(vchConfig); err != nil {
		executor.CollectDiagnosticLogs()
		log.Errorf("%s", err)
		return errors.New("delete failed")
	}

	log.Info("----------")
	log.Info("If firewall changes were made for VIC during install, they were not reverted during delete")
	log.Info("To modify firewall rules see vic-machine update firewall --help")
	log.Info("----------")
	log.Info("Completed successfully")

	return nil
}
