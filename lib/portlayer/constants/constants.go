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

package constants

import (
	"fmt"
	"time"

	"github.com/vmware/vic/pkg/version"
)

/* VCH constants */
const (
	SerialOverLANPort  = 2377
	AttachServerPort   = 2379
	ManagementHostName = "management.localhost"
	// BridgeScopeType denotes a scope that is of type bridge
	BridgeScopeType = "bridge"
	// ExternalScopeType denotes a scope that is of type external
	ExternalScopeType = "external"
	// DefaultBridgeRange is the default pool for bridge networks
	DefaultBridgeRange = "172.16.0.0/12"
	// Constants for assemble the VM display name on vSphere
	MaxVMNameLength = 80
	ShortIDLen      = 12
	// vSphere Display name for the VCH's Guest Name and for VAC support
	defaultAltVCHGuestName       = "Photon - VCH"
	defaultAltContainerGuestName = "Photon - Container"

	PropertyCollectorTimeout = 3 * time.Minute
)

func DefaultAltVCHGuestName() string {
	return fmt.Sprintf("%s %s, %s, %7s", defaultAltVCHGuestName, version.Version, version.BuildNumber, version.GitCommit)
}

func DefaultAltContainerGuestName() string {
	return fmt.Sprintf("%s %s, %s, %7s", defaultAltContainerGuestName, version.Version, version.BuildNumber, version.GitCommit)
}
