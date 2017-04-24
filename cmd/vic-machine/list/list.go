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

package list

import (
	"context"
	"fmt"
	"path"
	"text/tabwriter"
	"text/template"
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
)

type items struct {
	ID            string
	Path          string
	Name          string
	Version       string
	UpgradeStatus string
}

// templ is parsed by text/template package
const templ = `{{range .}}
{{.ID}}	{{.Path}}	{{.Name}}	{{.Version}}	{{.UpgradeStatus}}{{end}}
`

// List has all input parameters for vic-machine ls command
type List struct {
	*data.Data

	executor *management.Dispatcher
}

func NewList() *List {
	d := &List{}
	d.Data = data.NewData()
	return d
}

// Flags return all cli flags for ls
func (l *List) Flags() []cli.Flag {
	util := []cli.Flag{
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for list",
			Destination: &l.Timeout,
		},
	}

	target := l.TargetFlags()
	// TODO: why not allow name as a filter, like most list operations
	compute := l.ComputeFlagsNoName()
	debug := l.DebugFlags()

	// flag arrays are declared, now combined
	var flags []cli.Flag
	for _, f := range [][]cli.Flag{target, compute, util, debug} {
		flags = append(flags, f...)
	}

	return flags
}

func (l *List) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := l.HasCredentials(); err != nil {
		return err
	}

	return nil
}

func (l *List) prettyPrint(ctx context.Context, cli *cli.Context, vchs []*vm.VirtualMachine, executor *management.Dispatcher) {
	data := []items{
		{"ID", "PATH", "NAME", "VERSION", "UPGRADE STATUS"},
	}
	installerVer := version.GetBuild()
	for _, vch := range vchs {

		vchConfig, err := executor.GetVCHConfig(vch)
		var version string
		if err != nil {
			log.Error("Failed to get Virtual Container Host configuration")
			log.Error(err)
			version = "unknown"
		} else {
			version = vchConfig.Version.ShortVersion()
		}

		parentPath := path.Dir(path.Dir(vch.InventoryPath))
		name := path.Base(vch.InventoryPath)
		upgradeStatus := l.upgradeStatusMessage(ctx, vch, installerVer, vchConfig.Version)
		data = append(data,
			items{vch.Reference().Value, parentPath, name, version, upgradeStatus})
	}
	t := template.New("vic-machine ls")
	t, _ = t.Parse(templ)
	w := tabwriter.NewWriter(cli.App.Writer, 8, 8, 8, ' ', 0)
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
	}
	w.Flush()
}

func (l *List) Run(clic *cli.Context) (err error) {
	// urfave/cli will print out exit in error handling, so no more information in main method can be printed out.
	defer func() {
		err = common.LogErrorIfAny(clic, err)
	}()

	if err = l.processParams(); err != nil {
		return err
	}

	if l.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if len(clic.Args()) > 0 {
		log.Errorf("Unknown argument: %s", clic.Args()[0])
		return errors.New("invalid CLI arguments")
	}

	log.Infof("### Listing VCHs ####")

	ctx, cancel := context.WithTimeout(context.Background(), l.Timeout)
	defer cancel()
	defer func() {
		if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
			//context deadline exceeded, replace returned error message
			err = errors.Errorf("List timed out: use --timeout to add more time")
		}
	}()

	var validator *validate.Validator
	if validator, err = validate.NewValidator(ctx, l.Data); err != nil {
		log.Errorf("List cannot continue - failed to create validator: %s", err)
		return errors.New("list failed")
	}
	// If dc is not set, and multiple datacenter is available, vic-machine ls will list VCHs under all datacenters.
	validator.AllowEmptyDC()

	_, err = validator.ValidateTarget(ctx, l.Data)
	if err != nil {
		log.Errorf("List cannot continue - target validation failed: %s", err)
		return errors.New("list failed")
	}
	_, err = validator.ValidateCompute(ctx, l.Data, false)
	if err != nil {
		log.Errorf("List cannot continue - compute resource validation failed: %s", err)
		return errors.New("list failed")
	}
	executor := management.NewDispatcher(validator.Context, validator.Session, nil, false)
	vchs, err := executor.SearchVCHs(validator.ClusterPath)
	if err != nil {
		log.Errorf("List cannot continue - failed to search VCHs in %s: %s", validator.ResourcePoolPath, err)
	}
	l.prettyPrint(ctx, clic, vchs, executor)
	return nil
}

// upgradeStatusMessage generates a user facing status string about upgrade progress and status
func (l *List) upgradeStatusMessage(ctx context.Context, vch *vm.VirtualMachine, installerVer *version.Build, vchVer *version.Build) string {
	if sameVer := installerVer.Equal(vchVer); sameVer {
		return "Up to date"
	}

	upgrading, err := vch.VCHUpdateStatus(ctx)
	if err != nil {
		return fmt.Sprintf("Unknown: %s", err)
	}
	if upgrading {
		return "Upgrade in progress"
	}

	canUpgrade, err := installerVer.IsNewer(vchVer)
	if err != nil {
		return fmt.Sprintf("Unknown: %s", err)
	}
	if canUpgrade {
		return fmt.Sprintf("Upgradeable to %s", installerVer.ShortVersion())
	}

	oldInstaller, err := installerVer.IsOlder(vchVer)
	if err != nil {
		return fmt.Sprintf("Unknown: %s", err)
	}
	if oldInstaller {
		return fmt.Sprintf("VCH has newer version")
	}

	// can't get here
	return "Invalid upgrade status"
}
