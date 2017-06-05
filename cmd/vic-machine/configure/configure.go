// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package configure

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/vic/cmd/vic-machine/common"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/lib/install/management"
	"github.com/vmware/vic/lib/install/validate"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

// Configure has all input parameters for vic-machine configure command
type Configure struct {
	*data.Data

	proxies   common.Proxies
	cNetworks common.CNetworks

	upgrade  bool
	executor *management.Dispatcher
}

func NewConfigure() *Configure {
	configure := &Configure{}
	configure.Data = data.NewData()

	return configure
}

// Flags return all cli flags for configure
func (c *Configure) Flags() []cli.Flag {
	util := []cli.Flag{
		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "Force the configure operation",
			Destination: &c.Force,
		},
		cli.DurationFlag{
			Name:        "timeout",
			Value:       3 * time.Minute,
			Usage:       "Time to wait for configure",
			Destination: &c.Timeout,
		},
		cli.BoolFlag{
			Name:        "reset-progress",
			Usage:       "Reset the UpdateInProgress flag. Warning: Do not reset this flag if another upgrade/configure process is running",
			Destination: &c.ResetInProgressFlag,
		},
		cli.BoolFlag{
			Name:        "rollback",
			Usage:       "Roll back VCH configuration to before the current upgrade/configure",
			Destination: &c.Rollback,
			Hidden:      true,
		},
		cli.BoolFlag{
			Name:        "upgrade",
			Usage:       "Upgrade VCH to latest version together with configure",
			Destination: &c.upgrade,
			Hidden:      true,
		},
	}

	proxies := c.proxies.ProxyFlags(false)

	target := c.TargetFlags()
	id := c.IDFlags()
	compute := c.ComputeFlags()
	iso := c.ImageFlags(false)
	debug := c.DebugFlags()
	cNetwork := c.cNetworks.CNetworkFlags(false)

	// flag arrays are declared, now combined
	var flags []cli.Flag
	for _, f := range [][]cli.Flag{target, id, compute, iso, cNetwork, proxies, util, debug} {
		flags = append(flags, f...)
	}

	return flags
}

func (c *Configure) processParams() error {
	defer trace.End(trace.Begin(""))

	if err := c.HasCredentials(); err != nil {
		return err
	}

	hproxy, sproxy, err := c.proxies.ProcessProxies()
	if err != nil {
		return err
	}
	c.HTTPProxy = hproxy
	c.HTTPSProxy = sproxy

	c.ContainerNetworks, err = c.cNetworks.ProcessContainerNetworks()
	if err != nil {
		return err
	}

	return nil
}

// copyChangedConf should copy all configuration change specified by user to old configuration
// currently we cannot automatically override old configuration with any difference in the new configuration, because there are some configuration set in VCH
// creation process, for example, image store path, volume store path, network slot id, etc. So we'll copy changes based on user input
func (c *Configure) copyChangedConf(o *config.VirtualContainerHostConfigSpec, n *config.VirtualContainerHostConfigSpec) {
	//TODO: copy changed data
	personaSession := o.ExecutorConfig.Sessions[config.PersonaService]
	vicAdminSession := o.ExecutorConfig.Sessions[config.VicAdminService]
	if c.proxies.IsSet {
		hProxy := ""
		if c.HTTPProxy != nil {
			hProxy = c.HTTPProxy.String()
		}
		sProxy := ""
		if c.HTTPSProxy != nil {
			sProxy = c.HTTPSProxy.String()
		}
		updateSessionEnv(personaSession, config.GeneralHTTPProxy, hProxy)
		updateSessionEnv(personaSession, config.GeneralHTTPSProxy, sProxy)
		updateSessionEnv(vicAdminSession, config.VICAdminHTTPProxy, hProxy)
		updateSessionEnv(vicAdminSession, config.VICAdminHTTPSProxy, sProxy)
	}

	if c.Debug.Debug != nil {
		o.SetDebug(n.Diagnostics.DebugLevel)
	}

	if c.cNetworks.IsSet {
		o.ContainerNetworks = n.ContainerNetworks
	}
}

func updateSessionEnv(sess *executor.SessionConfig, envName, envValue string) {
	envs := sess.Cmd.Env
	var newEnvs []string
	for _, env := range envs {
		if strings.HasPrefix(env, envName+"=") {
			continue
		}
		newEnvs = append(newEnvs, env)
	}
	if envValue != "" {
		newEnvs = append(newEnvs, fmt.Sprintf("%s=%s", envName, envValue))
	}
	sess.Cmd.Env = newEnvs
}

func (c *Configure) Run(clic *cli.Context) (err error) {
	// urfave/cli will print out exit in error handling, so no more information in main method can be printed out.
	defer func() {
		err = common.LogErrorIfAny(clic, err)
	}()

	// process input parameters, this should reuse same code with create command, to make sure same options are provided
	if err = c.processParams(); err != nil {
		return err
	}

	if c.Debug.Debug != nil && *c.Debug.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		trace.Logger.Level = log.DebugLevel
	}

	if len(clic.Args()) > 0 {
		log.Errorf("Unknown argument: %s", clic.Args()[0])
		return errors.New("invalid CLI arguments")
	}

	// TODO: add additional parameter processing, reuse same code with create command as well

	if c.upgrade {
		// verify upgrade required parameters here
	}

	log.Infof("### Configuring VCH ####")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	validator, err := validate.NewValidator(ctx, c.Data)
	if err != nil {
		log.Errorf("Configuring cannot continue - failed to create validator: %s", err)
		return errors.New("configure failed")
	}
	_, err = validator.ValidateTarget(ctx, c.Data)
	if err != nil {
		log.Errorf("Configuring cannot continue - target validation failed: %s", err)
		return errors.New("configure failed")
	}
	executor := management.NewDispatcher(validator.Context, validator.Session, nil, c.Force)

	var vch *vm.VirtualMachine
	if c.Data.ID != "" {
		vch, err = executor.NewVCHFromID(c.Data.ID)
	} else {
		vch, err = executor.NewVCHFromComputePath(c.Data.ComputeResourcePath, c.Data.DisplayName, validator)
	}
	if err != nil {
		log.Errorf("Failed to get Virtual Container Host %s", c.DisplayName)
		log.Error(err)
		return errors.New("configure failed")
	}

	log.Infof("")
	log.Infof("VCH ID: %s", vch.Reference().String())

	if c.ResetInProgressFlag {
		if err = vch.SetVCHUpdateStatus(ctx, false); err != nil {
			log.Error("Failed to reset UpdateInProgress flag")
			log.Error(err)
			return errors.New("configure failed")
		}
		log.Infof("Reset UpdateInProgress flag successfully")
		return nil
	}

	var vchConfig *config.VirtualContainerHostConfigSpec
	if c.upgrade {
		vchConfig, err = executor.FetchAndMigrateVCHConfig(vch)
	} else {
		vchConfig, err = executor.GetVCHConfig(vch)
	}
	if err != nil {
		log.Error("Failed to get Virtual Container Host configuration")
		log.Error(err)
		return errors.New("configure failed")
	}

	installerVer := version.GetBuild().PluginVersion
	if vchConfig.ExecutorConfig.Version == nil {
		log.Error("Cannot configure VCH with version unavailable")
		return errors.New("configure failed")
	}
	if vchConfig.ExecutorConfig.Version.PluginVersion < installerVer {
		log.Error(fmt.Sprintf("Cannot configure VCH with version %s, specify --upgrade to upgrade VCH at the same time", vchConfig.ExecutorConfig.Version.ShortVersion()))
		return errors.New("configure failed")
	}

	// Convert guestinfo *VirtualContainerHost back to *Data, decrypt secret data
	oldData, err := validate.NewDataFromConfig(ctx, validator.Session.Finder, vchConfig)
	if err != nil {
		log.Error("Configuring cannot continue: configuration conversion failed")
		log.Error(err)
		return err
	}

	// using new configuration override configuration query from guestinfo
	if err = oldData.CopyNonEmpty(c.Data); err != nil {
		log.Error("Configuring cannot continue: copying configuration failed")
		return err
	}
	c.Data = oldData

	// evaluate merged configuration
	newConfig, err := validator.Validate(ctx, c.Data)
	if err != nil {
		log.Error("Configuring cannot continue: configuration validation failed")
		return err
	}

	// TODO: copy changed configuration here. https://github.com/vmware/vic/issues/2911
	c.copyChangedConf(vchConfig, newConfig)

	vConfig := validator.AddDeprecatedFields(ctx, vchConfig, c.Data)
	vConfig.Timeout = c.Timeout

	updating, err := vch.VCHUpdateStatus(ctx)
	if err != nil {
		log.Error("Unable to determine if upgrade/configure is in progress")
		log.Error(err)
		return errors.New("configure failed")
	}
	if updating {
		log.Error("Configure failed: another upgrade/configure operation is in progress")
		log.Error("If no other upgrade/configure process is running, use --reset-progress to reset the VCH upgrade/configure status")
		return errors.New("configure failed")
	}

	if err = vch.SetVCHUpdateStatus(ctx, true); err != nil {
		log.Error("Failed to set UpdateInProgress flag to true")
		log.Error(err)
		return errors.New("configure failed")
	}

	defer func() {
		if err = vch.SetVCHUpdateStatus(ctx, false); err != nil {
			log.Error("Failed to reset UpdateInProgress")
			log.Error(err)
		}
	}()

	if !c.Data.Rollback {
		err = executor.Configure(vch, vchConfig, vConfig)
	} else {
		err = executor.Rollback(vch, vchConfig, vConfig)
	}

	if err != nil {
		// configure failed
		executor.CollectDiagnosticLogs()
		return errors.New("configure failed")
	}

	log.Infof("Completed successfully")

	return nil
}
