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

package vicadmin

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	// "github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
	// "github.com/vmware/vic/pkg/vsphere/session"
	"golang.org/x/net/context"
)

type Validator struct {
	Hostname       string
	Version        string
	FirewallErrors []string
}

func NewValidator(ctx context.Context, vch metadata.VirtualContainerHostConfigSpec) *Validator {
	defer trace.End(trace.Begin(""))
	log.Infof("Creating new validator")
	v := new(Validator)
	v.checkFirewall(ctx, vch)
	v.FirewallErrors = nil
	v.Version = vch.Version
	log.Info(fmt.Sprintf("Setting version to %s", v.Version))

	v.Hostname, _ = os.Hostname()

	return v
}

func (v *Validator) checkFirewall(ctx context.Context, vch metadata.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))
	log.Infof("Checking firewall rules")
	//rule := types.HostFirewallRule{
	//Port:      8080, // serialOverLANPort
	//PortType:  types.HostFirewallRulePortTypeDst,
	//Protocol:  string(types.HostFirewallRuleProtocolTcp),
	//Direction: types.HostFirewallRuleDirectionOutbound,
	//}
}
