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

package metadata

import (
	"crypto/tls"
	"net/mail"
	"net/url"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

// Can we just treat the VCH appliance as a containerVM booting off a specific bootstrap image
// It has many of the same requirements (around networks being attached, version recorded,
// volumes mounted, et al). Each of the components can easily be captured as a Session given they
// are simply processes.
// This would require that the bootstrap read session record for the VM and relaunch them - that
// actually aligns very well with containerVMs restarting their processes if restarted directly
// (this is obviously a behaviour we'd want to toggles for in regular containers).

// VirtualContainerHostConfigSpec holds the metadata for a Virtual Container Host that should be visible inside the appliance VM.
type VirtualContainerHostConfigSpec struct {
	// The base config for the appliance. This includes the networks that are to be attached
	// and disks to be mounted.
	// Networks are keyed by interface name
	ExecutorConfig

	// The set of components to launch
	ComponentConfig []SessionConfig

	// Virtual Container Host version
	Version string

	// Administrative contact for the Virtual Container Host
	Admin []mail.Address
	// Administrative contact for hosting infrastructure
	InfrastructureAdmin []mail.Address

	// Certificates for user authentication - this needs to be expanded to allow for directory server auth
	UserCertificates []tls.Certificate
	// Certificates for general outgoing network access, keyed by CIDR (IPNet.String())
	NetworkCertificates map[string]tls.Certificate
	// Certificates for specific system access, keyed by FQDN
	HostCertificates map[string]tls.Certificate

	// Port Layer - storage
	// Datastore URLs for image stores - the top layer is [0], the bottom layer is [len-1]
	ImageStores []url.URL
	// Permitted datastore URL roots for volumes
	VolumeLocations []url.URL
	// Permitted datastore URLs for container storage for this virtual container host
	ContainerStores []url.URL
	// Compute resource roots under which all containers will be created
	ComputeResource []object.ComputeResource

	// Port Layer - network
	// The default bridge network supplied for the Virtual Container Host
	BridgeNetwork object.Network
	// Published networks available for containers to join, keyed by consumption name
	ContainerNetworks map[string]ContainerNetwork

	// Port Layer - exec
	// Default containerVM capacity
	ContainerVMSize Resources

	// Allow custom naming convention for containerVMs
	ContainerNameConvention string

	// Imagec
	// Whitelist of registries
	RegistryWhitelist []url.URL
	// Blacklist of registries
	RegistryBlacklist []url.URL
}

// CustomerExperienceImprovementProgram provides configuration for the phone home mechanism
// This is broken out so that we can have more granular configuration in here in the future
// and so that it is insulated from changes to Virtual Container Host structure
type CustomerExperienceImprovementProgram struct {
	// The server target is as follows, where the uuid is the raw number, no dashes
	// "https://vcsa.vmware.com/ph-stg/api/hyper/send?_v=1.0&_c=vic.1_0&_i="+vc.uuid
	// If this is non-nil then it's enabled
	CEIPGateway url.URL
}

// Resources is used instead of the ResourceAllocation structs in govmomi as
// those don't currently hold IO or storage related data.
type Resources struct {
	CPU     types.ResourceAllocationInfo
	Memory  types.ResourceAllocationInfo
	IO      types.ResourceAllocationInfo
	Storage types.ResourceAllocationInfo
}

// VirtualContainerHostTargetConfigSpec holds the metadata for a Virtual Container Host that should be visible only to the infrastructure on which it is running.
type VirtualContainerHostTargetConfigSpec struct {
	// Target and credentials for the environment - experimenting with storing this as an OAuth2 config
	// TargetConfig oauth2.Config

	// Virtual Container Host capacity
	VCHSize Resources
	// Appliance capacity
	ApplianceSize Resources
	// Freeform notes for target administrator user
	Notes string
}
