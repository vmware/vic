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

package inspect

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

// Inspect has all input parameters for vic-machine inspect command
type Inspect struct {
	*data.Data

	executor *management.Dispatcher
}

func NewInspect() *Inspect {
	d := &Inspect{}
	d.Data = data.NewData()
	return d
}

// Flags return all cli flags for inspect
func (i *Inspect) Flags() []cli.Flag {
	util := []cli.Flag{
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for inspect",
			Destination: &i.Timeout,
		},
	}

	target := i.TargetFlags()
	id := i.IDFlags()
	compute := i.ComputeFlags()
	debug := i.DebugFlags()

	// flag arrays are declared, now combined
	var flags []cli.Flag
	for _, f := range [][]cli.Flag{target, id, compute, util, debug} {
		flags = append(flags, f...)
	}

	return flags
}

func (i *Inspect) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := i.HasCredentials(); err != nil {
		return err
	}

	return nil
}

func (i *Inspect) Run(clic *cli.Context) (err error) {
	// urfave/cli will print out exit in error handling, so no more information in main method can be printed out.
	defer func() {
		err = common.LogErrorIfAny(clic, err)
	}()

	if err = i.processParams(); err != nil {
		return err
	}

	if i.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if len(clic.Args()) > 0 {
		log.Errorf("Unknown argument: %s", clic.Args()[0])
		return errors.New("invalid CLI arguments")
	}

	log.Infof("### Inspecting VCH ####")

	ctx, cancel := context.WithTimeout(context.Background(), i.Timeout)
	defer cancel()

	validator, err := validate.NewValidator(ctx, i.Data)
	if err != nil {
		log.Errorf("Inspect cannot continue - failed to create validator: %s", err)
		return errors.New("inspect failed")
	}
	_, err = validator.ValidateTarget(ctx, i.Data)
	if err != nil {
		log.Errorf("Inspect cannot continue - target validation failed: %s", err)
		return errors.New("inspect failed")
	}

	executor := management.NewDispatcher(validator.Context, validator.Session, nil, i.Force)

	var vch *vm.VirtualMachine
	if i.Data.ID != "" {
		vch, err = executor.NewVCHFromID(i.Data.ID)
	} else {
		vch, err = executor.NewVCHFromComputePath(i.Data.ComputeResourcePath, i.Data.DisplayName, validator)
	}
	if err != nil {
		log.Errorf("Failed to get Virtual Container Host %s", i.DisplayName)
		log.Error(err)
		return errors.New("inspect failed")
	}

	log.Infof("")
	log.Infof("VCH ID: %s", vch.Reference().String())

	vchConfig, err := executor.GetVCHConfig(vch)
	if err != nil {
		log.Error("Failed to get Virtual Container Host configuration")
		log.Error(err)
		return errors.New("inspect failed")
	}

	installerVer := version.GetBuild()

	log.Info("")
	log.Infof("Installer version: %s", installerVer.ShortVersion())
	log.Infof("VCH version: %s", vchConfig.Version.ShortVersion())
	log.Info("")
	log.Info("VCH upgrade status:")
	i.upgradeStatusMessage(ctx, vch, installerVer, vchConfig.Version)

	if err = executor.InspectVCH(vch, vchConfig); err != nil {
		executor.CollectDiagnosticLogs()
		log.Errorf("%s", err)
		return errors.New("inspect failed")
	}

	log.Infof("Completed successfully")

	return nil
}

// upgradeStatusMessage generates a user facing status string about upgrade progress and status
func (i *Inspect) upgradeStatusMessage(ctx context.Context, vch *vm.VirtualMachine, installerVer *version.Build, vchVer *version.Build) {
	if sameVer := installerVer.Equal(vchVer); sameVer {
		log.Info("Installer has same version as VCH")
		log.Info("No upgrade available with this installer version")
		return
	}

	upgrading, err := vch.VCHUpdateStatus(ctx)
	if err != nil {
		log.Errorf("Unable to determine if upgrade/configure is in progress: %s", err)
		return
	}
	if upgrading {
		log.Info("Upgrade/configure in progress")
		return
	}

	canUpgrade, err := installerVer.IsNewer(vchVer)
	if err != nil {
		log.Errorf("Unable to determine if upgrade is availabile: %s", err)
		return
	}
	if canUpgrade {
		log.Info("Upgrade available")
		return
	}

	oldInstaller, err := installerVer.IsOlder(vchVer)
	if err != nil {
		log.Errorf("Unable to determine if upgrade is available: %s", err)
		return
	}
	if oldInstaller {
		log.Info("Installer has older version than VCH")
		log.Info("No upgrade available with this installer version")
		return
	}

	// can't get here
	log.Warn("Invalid upgrade status")
	return
}
