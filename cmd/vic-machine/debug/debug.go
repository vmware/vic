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

package debug

import (
	"io/ioutil"
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

// Debug has all input parameters for vic-machine Debug command
type Debug struct {
	*data.Data

	executor *management.Dispatcher

	enableSSH     bool
	password      string
	authorizedKey string
}

func NewDebug() *Debug {
	d := &Debug{}
	d.Data = data.NewData()
	return d
}

// Flags return all cli flags for Debug
func (d *Debug) Flags() []cli.Flag {
	preFlags := append(d.TargetFlags(), d.IDFlags()...)
	preFlags = append(preFlags, d.ComputeFlags()...)

	ssh := []cli.Flag{
		cli.BoolFlag{
			Name:        "enable-ssh, ssh",
			Usage:       "Enable SSH server within appliance VM",
			Destination: &d.enableSSH,
		},
		cli.StringFlag{
			Name:        "authorized-key, key",
			Value:       "",
			Usage:       "File with public key to place as /root/.ssh/authorized_keys",
			Destination: &d.authorizedKey,
		},
		cli.StringFlag{
			Name:        "rootpw, pw",
			Value:       "",
			Usage:       "Password to set for root user (non-persistent over reboots)",
			Destination: &d.password,
		},
	}

	util := []cli.Flag{
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for operation to complete",
			Destination: &d.Timeout,
		},
	}

	target := d.TargetFlags()
	id := d.IDFlags()
	compute := d.ComputeFlags()
	debug := d.DebugFlags(true)

	// flag arrays are declared, now combined
	var flags []cli.Flag
	for _, f := range [][]cli.Flag{target, id, compute, ssh, util, debug} {
		flags = append(flags, f...)
	}

	return flags
}

func (d *Debug) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := d.HasCredentials(); err != nil {
		return err
	}

	return nil
}

func (d *Debug) Run(clic *cli.Context) (err error) {
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

	log.Infof("### Configuring VCH for debug ####")

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout)
	defer cancel()

	validator, err := validate.NewValidator(ctx, d.Data)

	if err != nil {
		log.Errorf("Debug cannot continue - failed to create validator: %s", err)
		return errors.New("debug failed")
	}
	defer validator.Session.Logout(ctx)

	_, err = validator.ValidateTarget(ctx, d.Data)
	if err != nil {
		log.Errorf("Debug cannot continue - target validation failed: %s", err)
		return errors.New("debug failed")
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
		return errors.New("debug failed")
	}

	log.Infof("")
	log.Infof("VCH ID: %s", vch.Reference().String())

	vchConfig, err := executor.GetVCHConfig(vch)
	if err != nil {
		log.Error("Failed to get Virtual Container Host configuration")
		log.Error(err)
		return errors.New("debug failed")
	}

	installerVer := version.GetBuild()

	log.Info("")
	log.Infof("Installer version: %s", installerVer.ShortVersion())
	log.Infof("VCH version: %s", vchConfig.Version.ShortVersion())

	// load the key file if set
	var key []byte
	if d.authorizedKey != "" {
		key, err = ioutil.ReadFile(d.authorizedKey)
		if err != nil {
			log.Errorf("Unable to read public key from %s: %s", d.authorizedKey, err)
			return errors.New("unable to load public key")
		}
	}

	if err = executor.DebugVCH(vch, vchConfig, d.password, string(key)); err != nil {
		executor.CollectDiagnosticLogs()
		log.Errorf("%s", err)
		return errors.New("debug failed")
	}

	// display the VCH endpoints again for convenience
	if err = executor.InspectVCH(vch, vchConfig, ""); err != nil {
		executor.CollectDiagnosticLogs()
		log.Errorf("%s", err)
		return errors.New("debug failed")
	}

	log.Infof("Completed successfully")

	return nil
}
