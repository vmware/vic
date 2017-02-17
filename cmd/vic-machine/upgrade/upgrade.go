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

package upgrade

import (
	"context"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// Upgrade has all input parameters for vic-machine upgrade command
type Upgrade struct {
	*data.Data

	executor *management.Dispatcher
}

func NewUpgrade() *Upgrade {
	upgrade := &Upgrade{}
	upgrade.Data = data.NewData()

	return upgrade
}

// Flags return all cli flags for upgrade
func (u *Upgrade) Flags() []cli.Flag {
	util := []cli.Flag{
		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "Force the upgrade (ignores version checks)",
			Destination: &u.Force,
		},
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for upgrade",
			Destination: &u.Timeout,
		},
	}

	target := u.TargetFlags()
	id := u.IDFlags()
	compute := u.ComputeFlags()
	iso := u.ImageFlags(false)
	debug := u.DebugFlags()

	// flag arrays are declared, now combined
	var flags []cli.Flag
	for _, f := range [][]cli.Flag{target, id, compute, iso, util, debug} {
		flags = append(flags, f...)
	}

	return flags
}

func (u *Upgrade) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := u.HasCredentials(); err != nil {
		return err
	}

	return nil
}

func (u *Upgrade) Run(clic *cli.Context) (err error) {
	// urfave/cli will print out exit in error handling, so no more information in main method can be printed out.
	defer func() {
		err = common.LogErrorIfAny(clic, err)
	}()

	if err = u.processParams(); err != nil {
		return err
	}

	if u.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if len(clic.Args()) > 0 {
		log.Errorf("Unknown argument: %s", clic.Args()[0])
		return errors.New("invalid CLI arguments")
	}

	var images map[string]string
	if images, err = u.CheckImagesFiles(u.Force); err != nil {
		return err
	}

	log.Infof("### Upgrading VCH ####")

	ctx, cancel := context.WithTimeout(context.Background(), u.Timeout)
	defer cancel()

	validator, err := validate.NewValidator(ctx, u.Data)
	if err != nil {
		log.Errorf("Upgrade cannot continue - failed to create validator: %s", err)
		return errors.New("upgrade failed")
	}
	_, err = validator.ValidateTarget(ctx, u.Data)
	if err != nil {
		log.Errorf("Upgrade cannot continue - target validation failed: %s", err)
		return errors.New("upgrade failed")
	}
	executor := management.NewDispatcher(validator.Context, validator.Session, nil, u.Force)

	var vch *vm.VirtualMachine
	if u.Data.ID != "" {
		vch, err = executor.NewVCHFromID(u.Data.ID)
	} else {
		vch, err = executor.NewVCHFromComputePath(u.Data.ComputeResourcePath, u.Data.DisplayName, validator)
	}
	if err != nil {
		log.Errorf("Failed to get Virtual Container Host %s", u.DisplayName)
		log.Error(err)
		return errors.New("upgrade failed")
	}

	log.Infof("")
	log.Infof("VCH ID: %s", vch.Reference().String())

	vchConfig, err := executor.FetchAndMigrateVCHConfig(vch)
	if err != nil {
		log.Error("Failed to get Virtual Container Host configuration")
		log.Error(err)
		return errors.New("upgrade failed")
	}

	executor.InitDiagnosticLogs(vchConfig)

	vConfig := validator.AddDeprecatedFields(ctx, vchConfig, u.Data)
	vConfig.ImageFiles = images
	vConfig.ApplianceISO = path.Base(u.ApplianceISO)
	vConfig.BootstrapISO = path.Base(u.BootstrapISO)
	vConfig.RollbackTimeout = u.Timeout

	if vchConfig, err = validator.ValidateMigratedConfig(ctx, vchConfig); err != nil {
		log.Errorf("Failed to migrate Virtual Container Host configuration %s", u.DisplayName)
		log.Error(err)
		return errors.New("upgrade failed")
	}

	if err = executor.Upgrade(vch, vchConfig, vConfig); err != nil {
		// upgrade failed
		executor.CollectDiagnosticLogs()
		if err == nil {
			err = errors.New("upgrade failed")
		}
		return err
	}

	// check the docker endpoint is responsive
	if err = executor.CheckDockerAPI(vchConfig, nil); err != nil {

		executor.CollectDiagnosticLogs()
		return err
	}

	log.Infof("Completed successfully")

	return nil
}
