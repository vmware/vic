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
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/cmd/vic-machine/converter"
	"github.com/vmware/vic/cmd/vic-machine/create"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// Inspect has all input parameters for vic-machine inspect command
type Inspect struct {
	*data.Data

	CertPath string

	executor *management.Dispatcher

	Format string
}

type state struct {
	i         *Inspect
	ctx       context.Context
	validator *validate.Validator
	vchConfig *config.VirtualContainerHostConfigSpec
	vch       *vm.VirtualMachine
	executor  *management.Dispatcher
}

type command func(state) error

func NewInspect() *Inspect {
	d := &Inspect{}
	d.Data = data.NewData()
	return d
}

// Flags returns all cli flags for inspect
func (i *Inspect) Flags() []cli.Flag {
	util := []cli.Flag{
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for inspect",
			Destination: &i.Timeout,
		},
		cli.StringFlag{
			Name:        "tls-cert-path",
			Value:       "",
			Usage:       "The path to check for existing certificates. Defaults to './<vch name>/'",
			Destination: &i.CertPath,
		},
	}

	target := i.TargetFlags()
	id := i.IDFlags()
	compute := i.ComputeFlags()
	debug := i.DebugFlags(true)

	// flag arrays are declared, now combined
	var flags []cli.Flag
	for _, f := range [][]cli.Flag{target, id, compute, util, debug} {
		flags = append(flags, f...)
	}

	return flags
}

func (i *Inspect) ConfigFlags() []cli.Flag {
	output := cli.StringFlag{
		Name:        "format",
		Value:       "verbose",
		Usage:       "Determine the format of configuration output. Supported formats: raw, verbose",
		Destination: &i.Format,
	}

	return []cli.Flag{output}
}

func (i *Inspect) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := i.HasCredentials(); err != nil {
		return err
	}

	return nil
}

func (i *Inspect) run(clic *cli.Context, cmd command) (err error) {
	// urfave/cli will print out exit in error handling, so no more information in main method can be printed out.
	defer func() {
		err = common.LogErrorIfAny(clic, err)
	}()

	if i.Debug.Debug != nil && *i.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if err = i.processParams(); err != nil {
		return err
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
	defer validator.Session.Logout(ctx)

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

	vchConfig, err := executor.GetNoSecretVCHConfig(vch)
	if err != nil {
		log.Error("Failed to get Virtual Container Host configuration")
		log.Error(err)
		return errors.New("inspect failed")
	}

	return cmd(state{i, ctx, validator, vchConfig, vch, executor})
}

func (i *Inspect) RunConfig(clic *cli.Context) (err error) {
	if i.Format == "raw" {
		log.SetLevel(log.ErrorLevel)
		log.SetOutput(os.Stderr)
	} else if i.Format != "verbose" {
		log.Warnf("Invalid configuration output format '%s'. Valid options are raw, verbose.", i.Format)
		log.Warn("Using verbose configuration format")
		i.Format = "verbose"
	}

	return i.run(clic, func(s state) error {
		err = i.showConfiguration(s.ctx, s.validator.Session.Finder, s.vchConfig, s.vch)
		if err != nil {
			log.Error("Failed to print Virtual Container Host configuration")
			log.Error(err)
			return errors.New("inspect failed")
		}
		return nil
	})
}

func (i *Inspect) Run(clic *cli.Context) (err error) {
	return i.run(clic, func(s state) error {
		log.Infof("")
		log.Infof("VCH ID: %s", s.vch.Reference().String())

		installerVer := version.GetBuild()

		log.Info("")
		log.Infof("Installer version: %s", installerVer.ShortVersion())
		log.Infof("VCH version: %s", s.vchConfig.Version.ShortVersion())
		log.Info("")
		log.Info("VCH upgrade status:")
		i.upgradeStatusMessage(s.ctx, s.vch, installerVer, s.vchConfig.Version)

		if err = s.executor.InspectVCH(s.vch, s.vchConfig, i.CertPath); err != nil {
			s.executor.CollectDiagnosticLogs()
			log.Errorf("%s", err)
			return errors.New("inspect failed")
		}

		log.Infof("Completed successfully")

		return nil
	})
}

func retrieveMapOptions(ctx context.Context, finder validate.Finder,
	conf *config.VirtualContainerHostConfigSpec, vm *vm.VirtualMachine) (map[string][]string, error) {
	data, err := validate.NewDataFromConfig(ctx, finder, conf)
	if err != nil {
		return nil, err
	}
	if err = validate.SetDataFromVM(ctx, finder, vm, data); err != nil {
		return nil, err
	}
	return converter.DataToOption(data)
}

func (i Inspect) showConfiguration(ctx context.Context, finder validate.Finder, conf *config.VirtualContainerHostConfigSpec, vm *vm.VirtualMachine) error {
	mapOptions, err := retrieveMapOptions(ctx, finder, conf, vm)
	if err != nil {
		return err
	}
	options := i.sortedOutput(mapOptions)
	if i.Format == "raw" {
		strOptions := strings.Join(options, " ")
		fmt.Println(strOptions)
	} else if i.Format == "verbose" {
		strOptions := strings.Join(options, "\n\t")
		log.Info("")
		log.Infof("Target VCH created with the following options: \n\n\t%s\n", strOptions)
	}

	return nil
}

func (i *Inspect) sortedOutput(mapOptions map[string][]string) (output []string) {
	create := create.NewCreate()
	cFlags := create.Flags()
	for _, f := range cFlags {
		key := f.GetName()
		// change multiple option name to long name: e.g. from target,t => target
		s := strings.Split(key, ",")
		if len(s) > 1 {
			key = s[0]
		}

		values, ok := mapOptions[key]
		if !ok {
			continue
		}

		defaultValue := ""
		switch t := f.(type) {
		case cli.StringFlag:
			defaultValue = t.Value
		case cli.IntFlag:
			defaultValue = strconv.Itoa(t.Value)
		}
		for _, val := range values {
			if val == defaultValue {
				// do not print command option if it's same to default
				continue
			}
			output = append(output, fmt.Sprintf("--%s=%s", key, val))
		}
	}
	return
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
		log.Errorf("Unable to determine if upgrade is available: %s", err)
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
