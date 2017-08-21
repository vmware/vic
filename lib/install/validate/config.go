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
	"net"

	"context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/govc/host/esxcli"
	"github.com/vmware/govmomi/license"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

type FirewallStatus struct {
	Rule                          types.HostFirewallRule
	MisconfiguredEnabled          []string
	MisconfiguredDisabled         []string
	UnknownEnabled                []string
	UnknownDisabled               []string
	MisconfiguredAllowedIPEnabled []string
	Correct                       []string
}

type FirewallConfigUnavailableError struct {
	Host string
}

func (e *FirewallConfigUnavailableError) Error() string {
	return fmt.Sprintf("Firewall configuration unavailable on %q", e.Host)
}

type FirewallMisconfiguredError struct {
	Host string
	Rule types.HostFirewallRule
}

func (e *FirewallMisconfiguredError) Error() string {
	return fmt.Sprintf("Firewall configuration on %q does not permit %s %d/%s %s",
		e.Host, e.Rule.PortType, e.Rule.Port, e.Rule.Protocol, e.Rule.Direction)
}

type FirewallUnknownDHCPAllowedIPError struct {
	AllowedIPs []string
	Host       string
	Rule       types.HostFirewallRule
	TargetIP   net.IPNet
}

func (e *FirewallUnknownDHCPAllowedIPError) Error() string {
	return fmt.Sprintf("Firewall configuration on %q may prevent connection on %s %d/%s %s with allowed IPs: %s",
		e.Host, e.Rule.PortType, e.Rule.Port, e.Rule.Protocol, e.Rule.Direction, e.AllowedIPs)
}

type FirewallMisconfiguredAllowedIPError struct {
	AllowedIPs []string
	Host       string
	Rule       types.HostFirewallRule
	TargetIP   net.IPNet
}

func (e *FirewallMisconfiguredAllowedIPError) Error() string {
	return fmt.Sprintf("Firewall configuration on %q does not permit %s %d/%s %s for %s with allowed IPs: %s",
		e.Host, e.Rule.PortType, e.Rule.Port, e.Rule.Protocol, e.Rule.Direction, e.TargetIP.IP, e.AllowedIPs)
}

// CheckFirewall verifies that host firewall configuration allows tether traffic and outputs results
func (v *Validator) CheckFirewall(ctx context.Context, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	mgmtIP := v.GetMgmtIP(conf)
	log.Debugf("Checking firewall with management network IP %s", mgmtIP)
	fwStatus := v.CheckFirewallForTether(ctx, mgmtIP)

	// log the results
	v.FirewallCheckOutput(fwStatus)
}

// CheckFirewallForTether which host firewalls are configured to allow tether traffic
func (v *Validator) CheckFirewallForTether(ctx context.Context, mgmtIP net.IPNet) FirewallStatus {

	var hosts []*object.HostSystem
	var err error

	requiredRule := types.HostFirewallRule{
		Port:      constants.SerialOverLANPort,
		PortType:  types.HostFirewallRulePortTypeDst,
		Protocol:  string(types.HostFirewallRuleProtocolTcp),
		Direction: types.HostFirewallRuleDirectionOutbound,
	}

	status := FirewallStatus{Rule: requiredRule}

	errMsg := "Firewall check SKIPPED"
	if !v.sessionValid(errMsg) {
		return status
	}

	if hosts, err = v.Session.Datastore.AttachedClusterHosts(ctx, v.Session.Cluster); err != nil {
		log.Errorf("Unable to get the list of hosts attached to given storage: %s", err)
		v.NoteIssue(err)
		return status
	}

	for _, host := range hosts {
		firewallEnabled, err := v.FirewallEnabled(host)
		if err != nil {
			v.NoteIssue(err)
			break
		}

		mgmtAllowed, err := v.ManagementNetAllowed(ctx, mgmtIP, host, requiredRule)
		if mgmtAllowed && err == nil {
			status.Correct = append(status.Correct, host.InventoryPath)
		}
		if err != nil {
			switch err.(type) {
			case *FirewallMisconfiguredError:
				if firewallEnabled {
					log.Debugf("fw misconfigured with fw enabled %q", host.InventoryPath)
					status.MisconfiguredEnabled = append(status.MisconfiguredEnabled, host.InventoryPath)
				} else {
					log.Debugf("fw misconfigured with fw disabled %q", host.InventoryPath)
					status.MisconfiguredDisabled = append(status.MisconfiguredDisabled, host.InventoryPath)
				}
			case *FirewallUnknownDHCPAllowedIPError:
				if firewallEnabled {
					log.Debugf("fw unknown (dhcp) with fw enabled %q", host.InventoryPath)
					status.UnknownEnabled = append(status.UnknownEnabled, host.InventoryPath)
				} else {
					log.Debugf("fw unknown (dhcp) with fw disabled %q", host.InventoryPath)
					status.UnknownDisabled = append(status.UnknownDisabled, host.InventoryPath)
				}
				log.Warn(err)
			case *FirewallMisconfiguredAllowedIPError:
				if firewallEnabled {
					log.Debugf("fw misconfigured allowed IP with fw enabled %q", host.InventoryPath)
					status.MisconfiguredAllowedIPEnabled = append(status.MisconfiguredAllowedIPEnabled, host.InventoryPath)
					log.Error(err)
				} else {
					log.Debugf("fw misconfigured allowed IP with fw disabled %q", host.InventoryPath)
					status.MisconfiguredDisabled = append(status.MisconfiguredDisabled, host.InventoryPath)
					log.Warn(err)
				}
			case *FirewallConfigUnavailableError:
				if firewallEnabled {
					log.Debugf("fw configuration unavailable %q", host.InventoryPath)
					status.UnknownEnabled = append(status.UnknownEnabled, host.InventoryPath)
					log.Error(err)
				} else {
					log.Debugf("fw configuration unavailable %q", host.InventoryPath)
					status.UnknownDisabled = append(status.UnknownDisabled, host.InventoryPath)
					log.Warn(err)
				}
			default:
				v.NoteIssue(err)
			}
			continue
		}
	}

	return status
}

// FirewallCheckOutput outputs firewall status messages associated
// with the hosts in each of the status categories
func (v *Validator) FirewallCheckOutput(status FirewallStatus) {
	var err error
	// TODO: when we can intelligently place containerVMs on hosts with proper config, install
	// can proceed if there is at least one host properly configured. For now this prevents install.
	if len(status.Correct) > 0 {
		log.Info("Firewall configuration OK on hosts:")
		for _, h := range status.Correct {
			log.Infof("\t%q", h)
		}
	}
	if len(status.MisconfiguredEnabled) > 0 {
		log.Error("Firewall configuration incorrect on hosts:")
		for _, h := range status.MisconfiguredEnabled {
			log.Errorf("\t%q", h)
		}
		err = fmt.Errorf("Firewall must permit %s %d/%s %s to the VCH management interface",
			status.Rule.PortType, status.Rule.Port, status.Rule.Protocol, status.Rule.Direction)
		log.Error(err)
		v.NoteIssue(err)
	}
	if len(status.MisconfiguredAllowedIPEnabled) > 0 {
		log.Error("Firewall configuration incorrect due to allowed IP restrictions on hosts:")
		for _, h := range status.MisconfiguredAllowedIPEnabled {
			log.Errorf("\t%q", h)
		}
		err = fmt.Errorf("Firewall must permit %s %d/%s %s to the VCH management interface",
			status.Rule.PortType, status.Rule.Port, status.Rule.Protocol, status.Rule.Direction)
		log.Error(err)
		v.NoteIssue(err)
	}
	if len(status.MisconfiguredDisabled) > 0 {
		log.Warning("Firewall configuration will be incorrect if firewall is reenabled on hosts:")
		for _, h := range status.MisconfiguredDisabled {
			log.Warningf("\t%q", h)
		}
		log.Infof("Firewall must permit %s %d/%s %s to VCH management interface if firewall is reenabled",
			status.Rule.PortType, status.Rule.Port, status.Rule.Protocol, status.Rule.Direction)
	}

	preMsg := "Firewall allowed IP configuration may prevent required connection on hosts:"
	postMsg := fmt.Sprintf("Firewall must permit %s %d/%s %s to the VCH management interface",
		status.Rule.PortType, status.Rule.Port, status.Rule.Protocol, status.Rule.Direction)
	v.FirewallCheckDHCPMessage(status.UnknownEnabled, preMsg, postMsg)

	preMsg = "Firewall configuration may be incorrect if firewall is reenabled on hosts:"
	postMsg = fmt.Sprintf("Firewall must permit %s %d/%s %s to the VCH management interface if firewall is reenabled",
		status.Rule.PortType, status.Rule.Port, status.Rule.Protocol, status.Rule.Direction)
	v.FirewallCheckDHCPMessage(status.UnknownDisabled, preMsg, postMsg)
}

// FirewallCheckDHCPMessage outputs warning message when we are unable to check
// that the management interface is allowed by the host firewall allowed IP rules due to DHCP
func (v *Validator) FirewallCheckDHCPMessage(hosts []string, preMsg, postMsg string) {
	if len(hosts) > 0 {
		log.Warn("Unable to fully verify firewall configuration due to DHCP use on management network")
		log.Warn("VCH management interface IP assigned by DHCP must be permitted by allowed IP settings")
		log.Warn(preMsg)
		for _, h := range hosts {
			log.Warningf("\t%q", h)
		}
		log.Infof(postMsg)
	}
}

// CheckIPInNets checks that an IP is within allowedIPs or allowedNets
func (v *Validator) CheckIPInNets(checkIP net.IPNet, allowedIPs []string, allowedNets []string) bool {
	for _, a := range allowedIPs {
		aIP := net.ParseIP(a)
		if aIP != nil && checkIP.IP.Equal(aIP) {
			return true
		}
	}

	for _, n := range allowedNets {
		_, aNet, err := net.ParseCIDR(n)
		if err != nil {
			continue
		}
		if aNet.Contains(checkIP.IP) {
			return true
		}
	}
	return false
}

func isMethodNotFoundError(err error) bool {
	if soap.IsSoapFault(err) {
		_, ok := soap.ToSoapFault(err).VimFault().(types.MethodNotFound)
		return ok
	}

	return false
}

// FirewallEnabled checks if the host firewall is enabled
func (v *Validator) FirewallEnabled(host *object.HostSystem) (bool, error) {
	esxfw, err := esxcli.GetFirewallInfo(host)
	if err != nil {
		if isMethodNotFoundError(err) {
			return true, nil // vcsim does not support esxcli; assume firewall is enabled in this case
		}
		return false, err
	}
	if esxfw.Enabled {
		log.Infof("Firewall status: ENABLED on %q", host.InventoryPath)
		return true, nil
	}
	log.Infof("Firewall status: DISABLED on %q", host.InventoryPath)
	return false, nil
}

// GetMgmtIP finds the management network IP in config
func (v *Validator) GetMgmtIP(conf *config.VirtualContainerHostConfigSpec) net.IPNet {
	var mgmtIP net.IPNet
	if conf != nil {
		n := conf.ExecutorConfig.Networks[config.ManagementNetworkName]
		if n != nil && n.Network.Common.Name == config.ManagementNetworkName {
			if n.IP != nil {
				mgmtIP = *n.IP
			}
			return mgmtIP
		}
	}
	return mgmtIP
}

// ManagementNetAllowed checks if the management network is allowed based
// on the host firewall's allowed IP settings
func (v *Validator) ManagementNetAllowed(ctx context.Context, mgmtIP net.IPNet,
	host *object.HostSystem, requiredRule types.HostFirewallRule) (bool, error) {

	fs, err := host.ConfigManager().FirewallSystem(ctx)
	if err != nil {
		return false, err
	}
	info, err := fs.Info(ctx)
	if err != nil {
		return false, err
	}

	// we've seen cases where the firewall config isn't available
	if info == nil {
		return false, &FirewallConfigUnavailableError{Host: host.InventoryPath}
	}

	rs := object.HostFirewallRulesetList(info.Ruleset)
	filteredRules, err := rs.EnabledByRule(requiredRule, true) // find matching rules that are enabled
	if err != nil {                                            // rule not enabled (fw is misconfigured)
		return false, &FirewallMisconfiguredError{Host: host.InventoryPath, Rule: requiredRule}
	}
	log.Debugf("filtered rules: %v", filteredRules)

	// check that allowed IPs permit management IP
	var allowedIPs []string
	var allowedNets []string
	for _, r := range filteredRules {
		log.Debugf("filtered IPs: %v networks: %v allIP: %v rule: %v",
			r.AllowedHosts.IpAddress, r.AllowedHosts.IpNetwork, r.AllowedHosts.AllIp, r.Key)

		if r.AllowedHosts.AllIp { // this rule allows all hosts
			return true, nil
		}
		for _, h := range r.AllowedHosts.IpAddress {
			allowedIPs = append(allowedIPs, h)
		}
		for _, n := range r.AllowedHosts.IpNetwork {
			s := fmt.Sprintf("%s/%d", n.Network, n.PrefixLength)
			allowedNets = append(allowedNets, s)
		}
	}

	if mgmtIP.IP == nil { // DHCP
		if len(allowedIPs) > 0 || len(allowedNets) > 0 {
			return false, &FirewallUnknownDHCPAllowedIPError{AllowedIPs: append(allowedNets, allowedIPs...),
				Host:     host.InventoryPath,
				Rule:     requiredRule,
				TargetIP: mgmtIP}
		}
		// no allowed IPs
		return false, &FirewallMisconfiguredError{Host: host.InventoryPath, Rule: requiredRule}
	}

	// static management IP, check that it is allowed
	mgmtAllowed := v.CheckIPInNets(mgmtIP, allowedIPs, allowedNets)
	if mgmtAllowed {
		return true, nil
	}
	return false, &FirewallMisconfiguredAllowedIPError{AllowedIPs: append(allowedNets, allowedIPs...),
		Host:     host.InventoryPath,
		Rule:     requiredRule,
		TargetIP: mgmtIP}
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
		log.Errorf("  DRS must be enabled on cluster %q", v.Session.Cluster.InventoryPath)
		v.NoteIssue(errors.New("DRS must be enabled to use VIC"))
		return
	}
	log.Info("DRS check OK on:")
	log.Infof("  %q", v.Session.Cluster.InventoryPath)
}
