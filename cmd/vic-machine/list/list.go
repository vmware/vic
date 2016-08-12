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
	"path"
	"text/tabwriter"
	"text/template"

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

type items struct {
	ID   string
	Path string
	Name string
}

// templ is parsed by text/template package
const templ = `{{range .}}
{{.ID}}	{{.Path}}	{{.Name}}{{end}}
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
	flags := append(l.TargetFlags(), l.ComputeFlagsNoName()...)
	flags = append(flags, l.DebugFlags()...)

	return flags
}

func (l *List) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := l.HasCredentials(); err != nil {
		return err
	}

	l.Insecure = true
	return nil
}

func (l *List) prettyPrint(cli *cli.Context, ctx context.Context, vchs []*vm.VirtualMachine) {
	data := []items{
		{"ID", "PATH", "NAME"},
	}
	for _, vch := range vchs {
		parentPath := path.Dir(path.Dir(vch.InventoryPath))
		name := path.Base(vch.InventoryPath)
		data = append(data,
			items{vch.Reference().Value, parentPath, name})
	}
	t := template.New("vic-machine ls")
	t, _ = t.Parse(templ)
	w := tabwriter.NewWriter(cli.App.Writer, 8, 8, 8, ' ', 0)
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
	}
	w.Flush()
}

func (l *List) Run(cli *cli.Context) error {
	var err error
	if err = l.processParams(); err != nil {
		return err
	}

	if l.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if len(cli.Args()) > 0 {
		log.Errorf("Unknown argument: %s", cli.Args()[0])
		return errors.New("invalid CLI arguments")
	}

	log.Infof("### Listing VCHs ####")

	ctx, cancel := context.WithTimeout(context.Background(), l.Timeout)
	defer cancel()

	var validator *validate.Validator
	if l.Data.ComputeResourcePath == "" {
		validator, err = validate.CreateNoDCCheck(ctx, l.Data)
	} else {
		validator, err = validate.NewValidator(ctx, l.Data)
	}
	if err != nil {
		log.Errorf("List cannot continue - failed to create validator: %s", err)
		return errors.New("list failed")
	}

	_, err = validator.ValidateCompute(ctx, l.Data)
	if err != nil {
		log.Errorf("List cannot continue - compute resource validation failed: %s", err)
		return err
	}
	executor := management.NewDispatcher(validator.Context, validator.Session, nil, false)
	vchs, err := executor.SearchVCHs(validator.ResourcePoolPath)
	if err != nil {
		log.Errorf("List cannot continue - failed to search VCHs in %s: %s", validator.ResourcePoolPath, err)
	}
	l.prettyPrint(cli, ctx, vchs)
	return nil
}
