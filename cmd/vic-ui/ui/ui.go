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

package ui

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"context"

	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/install/plugin"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

// Plugin has all input parameters for vic-ui ui command
type Plugin struct {
	*common.Target
	common.Debug

	Force    bool
	Insecure bool

	Company               string
	HideInSolutionManager bool
	Key                   string
	Name                  string
	ServerThumbprint      string
	Summary               string
	Type                  string
	URL                   string
	Version               string
}

func NewUI() *Plugin {
	p := &Plugin{Target: common.NewTarget()}
	return p
}

// Flags return all cli flags for ui
func (p *Plugin) Flags() []cli.Flag {
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "Force install",
			Destination: &p.Force,
		},
		cli.StringFlag{
			Name:        "company",
			Value:       "",
			Usage:       "Plugin company name (required)",
			Destination: &p.Company,
		},
		cli.StringFlag{
			Name:        "key",
			Value:       "",
			Usage:       "Plugin key (required)",
			Destination: &p.Key,
		},
		cli.StringFlag{
			Name:        "name",
			Value:       "",
			Usage:       "Plugin name (required)",
			Destination: &p.Name,
		},
		cli.BoolFlag{
			Name:        "no-show",
			Usage:       "Hide plugin in UI",
			Destination: &p.HideInSolutionManager,
		},
		cli.StringFlag{
			Name:        "server-thumbprint",
			Value:       "",
			Usage:       "Plugin server thumbprint (required for HTTPS plugin URL)",
			Destination: &p.ServerThumbprint,
		},
		cli.StringFlag{
			Name:        "summary",
			Value:       "",
			Usage:       "Plugin summary (required)",
			Destination: &p.Summary,
		},
		cli.StringFlag{
			Name:        "url",
			Value:       "",
			Usage:       "Plugin URL (required)",
			Destination: &p.URL,
		},
		cli.StringFlag{
			Name:        "version",
			Value:       "",
			Usage:       "Plugin version (required)",
			Destination: &p.Version,
		},
	}
	flags = append(p.TargetFlags(), flags...)
	flags = append(flags, p.DebugFlags()...)

	return flags
}

func (p *Plugin) processInstallParams() error {
	defer trace.End(trace.Begin(""))

	if err := p.HasCredentials(); err != nil {
		return err
	}

	if p.Company == "" {
		return cli.NewExitError("--company must be specified", 1)
	}

	if p.Key == "" {
		return cli.NewExitError("--key must be specified", 1)
	}

	if p.Name == "" {
		return cli.NewExitError("--name must be specified", 1)
	}

	if p.Summary == "" {
		return cli.NewExitError("--summary must be specified", 1)
	}

	if p.URL == "" {
		return cli.NewExitError("--url must be specified", 1)
	}

	if p.Version == "" {
		return cli.NewExitError("--version must be specified", 1)
	}

	if strings.HasPrefix(strings.ToLower(p.URL), "https://") && p.ServerThumbprint == "" {
		return cli.NewExitError("--server-thumbprint must be specified when using HTTPS plugin URL", 1)
	}

	p.Insecure = true
	return nil
}

func (p *Plugin) processRemoveParams() error {
	defer trace.End(trace.Begin(""))

	if err := p.HasCredentials(); err != nil {
		return err
	}

	if p.Key == "" {
		return cli.NewExitError("--key must be specified", 1)
	}

	p.Insecure = true
	return nil
}

func (p *Plugin) Install(cli *cli.Context) error {
	var err error
	if err = p.processInstallParams(); err != nil {
		return err
	}

	if p.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if len(cli.Args()) > 0 {
		log.Error("Install cannot continue: invalid CLI arguments")
		log.Errorf("Unknown argument: %s", cli.Args()[0])
		return errors.New("invalid CLI arguments")
	}

	log.Infof("### Installing UI Plugin ####")

	pInfo := &plugin.Info{
		Company:               p.Company,
		Key:                   p.Key,
		Name:                  p.Name,
		ServerThumbprint:      p.ServerThumbprint,
		ShowInSolutionManager: !p.HideInSolutionManager,
		Summary:               p.Summary,
		Type:                  "vsphere-client-serenity",
		URL:                   p.URL,
		Version:               p.Version,
	}

	pl, err := plugin.NewPluginator(context.TODO(), p.Target.URL, pInfo)
	if err != nil {
		return err
	}

	reg, err := pl.IsRegistered(pInfo.Key)
	if err != nil {
		return err
	}
	if reg {
		if p.Force {
			log.Info("Removing existing plugin to force install")
			err = pl.Unregister(pInfo.Key)
			if err != nil {
				return err
			}
			log.Info("Removed existing plugin")
		} else {
			msg := fmt.Sprintf("plugin (%s) is already registered", pInfo.Key)
			log.Errorf("Install failed: %s", msg)
			return errors.New(msg)
		}
	}

	log.Info("Installing plugin")
	err = pl.Register()
	if err != nil {
		return err
	}

	reg, err = pl.IsRegistered(pInfo.Key)
	if err != nil {
		return err
	}
	if !reg {
		msg := fmt.Sprintf("post-install check failed to find %s registered", pInfo.Key)
		log.Errorf("Install failed: %s", msg)
		return errors.New(msg)
	}

	log.Info("Installed UI plugin")
	return nil
}

func (p *Plugin) Remove(cli *cli.Context) error {
	var err error
	if err = p.processRemoveParams(); err != nil {
		return err
	}
	if p.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if len(cli.Args()) > 0 {
		log.Error("Remove cannot continue: invalid CLI arguments")
		log.Errorf("Unknown argument: %s", cli.Args()[0])
		return errors.New("invalid CLI arguments")
	}

	if p.Force {
		log.Info("Ignoring --force")
	}

	log.Infof("### Removing UI Plugin ####")

	pInfo := &plugin.Info{
		Key: p.Key,
	}

	pl, err := plugin.NewPluginator(context.TODO(), p.Target.URL, pInfo)
	if err != nil {
		return err
	}
	reg, err := pl.IsRegistered(pInfo.Key)
	if err != nil {
		return err
	}
	if reg {
		log.Infof("Found target plugin: %s", pInfo.Key)
	} else {
		msg := fmt.Sprintf("failed to find target plugin (%s)", pInfo.Key)
		log.Errorf("Remove failed: %s", msg)
		return errors.New(msg)
	}

	log.Info("Removing plugin")
	err = pl.Unregister(pInfo.Key)
	if err != nil {
		return err
	}

	reg, err = pl.IsRegistered(pInfo.Key)
	if err != nil {
		return err
	}
	if reg {
		msg := fmt.Sprintf("post-remove check found %s still registered", pInfo.Key)
		log.Errorf("Remove failed: %s", msg)
		return errors.New(msg)
	}

	log.Info("Removed UI plugin")
	return nil
}
