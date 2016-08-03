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

package validate

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/govc/host/esxcli"
	"github.com/vmware/govmomi/license"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"golang.org/x/net/context"
)

func (v *Validator) CheckFirewall(ctx context.Context) {
	if v.DisableFirewallCheck {
		return
	}
	defer trace.End(trace.Begin(""))

	var hosts []*object.HostSystem
	var err error

	rule := types.HostFirewallRule{
		Port:      8080, // serialOverLANPort
		PortType:  types.HostFirewallRulePortTypeDst,
		Protocol:  string(types.HostFirewallRuleProtocolTcp),
		Direction: types.HostFirewallRuleDirectionOutbound,
	}

	errMsg := "Firewall check SKIPPED"
	if !v.sessionValid(errMsg) {
		return
	}

	if hosts, err = v.Session.Datastore.AttachedClusterHosts(ctx, v.Session.Cluster); err != nil {
		log.Errorf("Unable to get the list of hosts attached to given storage: %s", err)
		v.NoteIssue(err)
		return
	}

	var misconfiguredEnabled []string
	var misconfiguredDisabled []string
	var correct []string

	for _, host := range hosts {
		fs, err := host.ConfigManager().FirewallSystem(ctx)
		if err != nil {
			v.NoteIssue(err)
			break
		}

		disabled := false
		esxfw, err := esxcli.GetFirewallInfo(host)
		if err != nil {
			v.NoteIssue(err)
			break
		}
		if !esxfw.Enabled {
			disabled = true
			log.Infof("Firewall status: DISABLED on %q", host.InventoryPath)
		} else {
			log.Infof("Firewall status: ENABLED on %q", host.InventoryPath)
		}

		info, err := fs.Info(ctx)
		if err != nil {
			v.NoteIssue(err)
			break
		}

		rs := object.HostFirewallRulesetList(info.Ruleset)
		_, err = rs.EnabledByRule(rule, true)
		if err != nil {
			if !disabled {
				misconfiguredEnabled = append(misconfiguredEnabled, host.InventoryPath)
			} else {
				misconfiguredDisabled = append(misconfiguredDisabled, host.InventoryPath)
			}
		} else {
			correct = append(correct, host.InventoryPath)
		}
	}

	if len(correct) > 0 {
		log.Info("Firewall configuration OK on hosts:")
		for _, h := range correct {
			log.Infof("  %q", h)
		}
	}
	if len(misconfiguredEnabled) > 0 {
		log.Error("Firewall configuration incorrect on hosts:")
		for _, h := range misconfiguredEnabled {
			log.Errorf("  %q", h)
		}
		// TODO: when we can intelligently place containerVMs on hosts with proper config, install
		// can proceed if there is at least one host properly configured. For now this prevents install.
		msg := "Firewall must permit 8080/tcp outbound to use VIC"
		log.Error(msg)
		v.NoteIssue(errors.New(msg))
	}
	if len(misconfiguredDisabled) > 0 {
		log.Warning("Firewall configuration will be incorrect if firewall is reenabled on hosts:")
		for _, h := range misconfiguredDisabled {
			log.Warningf("  %q", h)
		}
		log.Warning("Firewall must permit 8080/tcp outbound if firewall is reenabled")
	}
	return
}

func (v *Validator) CheckLicense(ctx context.Context) {
	var err error

	errMsg := "License check SKIPPED"
	if !v.sessionValid(errMsg) {
		return
	}

	if v.IsVC() {
		if err = v.checkAssignedLicenses(ctx); err != nil {
			v.NoteIssue(err)
			return
		}
	} else {
		if err = v.checkLicense(ctx); err != nil {
			v.NoteIssue(err)
			return
		}
	}
}

func (v *Validator) assignedLicenseHasFeature(la []types.LicenseAssignmentManagerLicenseAssignment, feature string) bool {
	for _, a := range la {
		if license.HasFeature(a.AssignedLicense, feature) {
			return true
		}
	}
	return false
}

func (v *Validator) checkAssignedLicenses(ctx context.Context) error {
	var hosts []*object.HostSystem
	var invalidLic []string
	var validLic []string
	var err error
	client := v.Session.Client.Client

	if hosts, err = v.Session.Datastore.AttachedClusterHosts(ctx, v.Session.Cluster); err != nil {
		log.Errorf("Unable to get the list of hosts attached to given storage: %s", err)
		return err
	}

	lm := license.NewManager(client)

	am, err := lm.AssignmentManager(ctx)
	if err != nil {
		return err
	}

	features := []string{"serialuri", "dvs"}

	for _, host := range hosts {
		valid := true
		la, err := am.QueryAssigned(ctx, host.Reference().Value)
		if err != nil {
			return err
		}

		for _, feature := range features {
			if !v.assignedLicenseHasFeature(la, feature) {
				valid = false
				msg := fmt.Sprintf("%q - license missing feature %q", host.InventoryPath, feature)
				invalidLic = append(invalidLic, msg)
			}
		}

		if valid == true {
			validLic = append(validLic, host.InventoryPath)
		}
	}

	if len(validLic) > 0 {
		log.Infof("License check OK on hosts:")
		for _, h := range validLic {
			log.Infof("  %q", h)
		}
	}
	if len(invalidLic) > 0 {
		log.Errorf("License check FAILED on hosts:")
		for _, h := range invalidLic {
			log.Errorf("  %q", h)
		}
		msg := "License does not meet minimum requirements to use VIC"
		return errors.New(msg)
	}
	return nil
}

func (v *Validator) checkLicense(ctx context.Context) error {
	var invalidLic []string
	client := v.Session.Client.Client

	lm := license.NewManager(client)
	licenses, err := lm.List(ctx)
	if err != nil {
		return err
	}
	v.checkEvalLicense(licenses)

	features := []string{"serialuri"}

	for _, feature := range features {
		if len(licenses.WithFeature(feature)) == 0 {
			msg := fmt.Sprintf("Host license missing feature %q", feature)
			invalidLic = append(invalidLic, msg)
		}
	}

	if len(invalidLic) > 0 {
		log.Errorf("License check FAILED:")
		for _, h := range invalidLic {
			log.Errorf("  %q", h)
		}
		msg := "License does not meet minimum requirements to use VIC"
		return errors.New(msg)
	}
	log.Infof("License check OK")
	return nil
}

func (v *Validator) checkEvalLicense(licenses []types.LicenseManagerLicenseInfo) {
	for _, l := range licenses {
		if l.EditionKey == "eval" {
			log.Warning("Evaluation license detected. VIC may not function if evaluation expires or insufficient license is later assigned.")
		}
	}
}

// isStandaloneHost checks if host is ESX or vCenter with single host
func (v *Validator) isStandaloneHost() bool {
	cl := v.Session.Cluster.Reference()

	if cl.Type != "ClusterComputeResource" {
		return true
	}
	return false
}

// drs checks that DRS is enabled
func (v *Validator) CheckDrs(ctx context.Context) {
	if v.DisableDRSCheck {
		return
	}
	defer trace.End(trace.Begin(""))

	errMsg := "DRS check SKIPPED"
	if !v.sessionValid(errMsg) {
		return
	}

	cl := v.Session.Cluster
	ref := cl.Reference()

	if v.isStandaloneHost() {
		log.Info("DRS check SKIPPED - target is standalone host")
		return
	}

	var ccr mo.ClusterComputeResource

	err := cl.Properties(ctx, ref, []string{"configurationEx"}, &ccr)
	if err != nil {
		msg := fmt.Sprintf("Failed to validate DRS config: %s", err)
		v.NoteIssue(errors.New(msg))
		return
	}

	z := ccr.ConfigurationEx.(*types.ClusterConfigInfoEx).DrsConfig

	if !(*z.Enabled) {
		log.Error("DRS check FAILED")
		log.Errorf("  DRS must be enabled on cluster %q", v.Session.Pool.InventoryPath)
		v.NoteIssue(errors.New("DRS must be enabled to use VIC"))
		return
	}
	log.Info("DRS check OK on:")
	log.Infof("  %q", v.Session.Pool.InventoryPath)
}
