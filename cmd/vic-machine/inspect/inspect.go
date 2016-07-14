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

package inspect

import (
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
type Inspect struct {
	*data.Data

	logfile string

	executor *management.Dispatcher
}

func NewInspect() *Inspect {
	d := &Inspect{}
	d.Data = data.NewData()
	return d
}

// Flags return all cli flags for delete
func (i *Inspect) Flags() []cli.Flag {
	preFlags := append(i.TargetFlags(), i.IDFlags()...)
	preFlags = append(preFlags, i.ComputeFlags()...)
	flags := append(preFlags, i.DebugFlags()...)

	return flags
}

func (i *Inspect) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := i.HasCredentials(); err != nil {
		return err
	}

	if err := i.ProcessID(); err != nil {
		return err
	}

	i.logfile = "inspect.log"
	i.Insecure = true
	return nil
}

func (i *Inspect) Run(cli *cli.Context) error {
	var err error
	if err = i.processParams(); err != nil {
		return err
	}

	if i.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if len(cli.Args()) > 0 {
		log.Error("Inspect cannot continue: invalid CLI arguments")
		log.Errorf("Unknown argument: %s", cli.Args()[0])
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
	executor := management.NewDispatcher(validator.Context, validator.Session, nil, i.Force)

	vch, path, err := executor.NewVCHFromComputePath(i.Data.ComputeResourcePath, i.Data.DisplayName, validator)
	if err != nil {
		log.Errorf("Failed to get Virtual Container Host %s", i.DisplayName)
		log.Error(err)
		return errors.New("inspect failed")
	}

	log.Infof("")
	log.Infof("VCH: %s", path)

	vchConfig, err := executor.GetVCHConfig(vch)
	if err != nil {
		log.Error("Failed to get Virtual Container Host configuration")
		log.Error(err)
		return errors.New("inspect failed")
	}
	executor.InitDiagnosticLogs(vchConfig)

	if err = executor.InspectVCH(vch, vchConfig); err != nil {
		executor.CollectDiagnosticLogs()
		log.Errorf("%s", err)
		return errors.New("inspect failed")
	}

	log.Infof("Completed successfully")

	return nil
}
