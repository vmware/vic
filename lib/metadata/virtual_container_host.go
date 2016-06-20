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
	"fmt"
	"net/mail"
	"net/url"
	"time"

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
	ExecutorConfig `vic:"0.1" scope:"read-only" key:"init"`

	// The sdk URL
	Target url.URL `vic:"0.1" scope:"read-only" key:"target"`
	// Whether the session connection is secure
	Insecure bool `vic:"0.1" scope:"read-only" key:"insecure"`
	// The session timeout
	Keepalive time.Duration `vic:"0.1" scope:"read-only" key:"keepalive"`
	// Turn on debug logging
	Debug bool `vic:"0.1" scope:"read-only" key:"debug"`
	// Virtual Container Host version
	Version string `vic:"0.1" scope:"read-only" key:"version"`

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
	ImageStores []url.URL `vic:"0.1" scope:"read-only" key:"image_stores"`
	// Permitted datastore URL roots for volumes
	VolumeLocations map[string]url.URL `vic:"0.1" scope:"read-only"`
	// Permitted datastore URLs for container storage for this virtual container host
	ContainerStores []url.URL `vic:"0.1" scope:"read-only" recurse:"depth=0"`
	// Resource pools under which all containers will be created
	ComputeResources []types.ManagedObjectReference `vic:"0.1" scope:"read-only"`

	// Port Layer - network
	// The default bridge network supplied for the Virtual Container Host
	BridgeNetwork string `vic:"0.1" scope:"read-only" key:"bridge_network"`
	// Published networks available for containers to join, keyed by consumption name
	ContainerNetworks map[string]*ContainerNetwork `vic:"0.1" scope:"read-only" key:"container_networks"`

	// Virtual Container Host capacity
	VCHSize Resources `vic:"0.1" scope:"read-only" recurse:"depth=0"`
	// Appliance capacity
	ApplianceSize Resources `vic:"0.1" scope:"read-only" recurse:"depth=0"`

	// Port Layer - exec
	// Default containerVM capacity
	ContainerVMSize Resources `vic:"0.1" scope:"read-only" recurse:"depth=0"`

	// Allow custom naming convention for containerVMs
	ContainerNameConvention string

	// Imagec
	// Whitelist of registries
	RegistryWhitelist []url.URL `vic:"0.1" scope:"read-only" recurse:"depth=0"`
	// Blacklist of registries
	RegistryBlacklist []url.URL `vic:"0.1" scope:"read-only" recurse:"depth=0"`

	KeyPEM  string `vic:"0.1" scope:"read-only" key:"key_pem"`
	CertPEM string `vic:"0.1" scope:"read-only" key:"cert_pem"`

	//FIXME: remove following attributes after port-layer-server read config from guestinfo
	DatacenterName         string `vic:"0.1" scope:"read-only" key:"datacenter_name"`
	ClusterPath            string `vic:"0.1" scope:"read-only" key:"cluster_path"`
	ImageStoreName         string `vic:"0.1" scope:"read-only" key:"image_store_name"`
	ResourcePoolPath       string `vic:"0.1" scope:"read-only" key:"resource_pool_path"`
	ApplianceInventoryPath string `vic:"0.1" scope:"read-only" key:"appliance_path"`

	Datacenter types.ManagedObjectReference `vic:"0.1" scope:"read-only" key:"datacenter"`
	Cluster    types.ManagedObjectReference `vic:"0.1" scope:"read-only" key:"cluster"`

	// FIXME: remove following attributes after change to launch through tether
	// Networks represents mapping between nic name and network info object. For example: bridge: vmomi object
	Networks map[string]*NetworkInfo `vic:"0.1" scope:"read-only" key:"networks2"`

	ImageFiles []string `vic:"0.1" scope:"read-only" recurse:"depth=0"`
}

type NetworkInfo struct {
	InventoryPath string                       `vic:"0.1" scope:"read-only" recurse:"depth=0"`
	Mac           string                       `vic:"0.1" scope:"read-only" key:"mac"`
	PortGroup     object.NetworkReference      `vic:"0.1" scope:"read-only" recurse:"depth=0"`
	PortGroupName string                       `vic:"0.1" scope:"read-only" key:"portgroup"`
	PortGroupRef  types.ManagedObjectReference `vic:"0.1" scope:"read-only" key:"portgroup_ref"`
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

// SetHostCertificate sets the certificate for authenticting with the appliance itself
func (t *VirtualContainerHostConfigSpec) SetHostCertificate(key *[]byte) {
	t.ExecutorConfig.Key = *key
}

// SetName sets the name of the VCH - this will be used as the hostname for the appliance
func (t *VirtualContainerHostConfigSpec) SetName(name string) {
	t.ExecutorConfig.Common.Name = name
}

// SetMoref sets the moref of the VCH - this allows components to acquire a handle to
// the appliance VM.
func (t *VirtualContainerHostConfigSpec) SetMoref(moref *types.ManagedObjectReference) {
	if moref != nil {
		t.ExecutorConfig.Common.ID = fmt.Sprintf("%s-%s", moref.Type, moref.Value)
	}
}

// AddNetwork adds a network that will be configured on the appliance VM
func (t *VirtualContainerHostConfigSpec) AddNetwork(net *NetworkEndpoint) {
	if net != nil {
		if t.ExecutorConfig.Networks == nil {
			t.ExecutorConfig.Networks = make(map[string]*NetworkEndpoint)
		}

		t.ExecutorConfig.Networks[net.Network.Common.Name] = net
	}
}

func (t *VirtualContainerHostConfigSpec) AddComponent(name string, component *SessionConfig) {
	if component != nil {
		if t.ExecutorConfig.Sessions == nil {
			t.ExecutorConfig.Sessions = make(map[string]SessionConfig)
		}

		t.ExecutorConfig.Sessions[name] = *component
	}
}

func (t *VirtualContainerHostConfigSpec) AddImageStore(url *url.URL) {
	if url != nil {
		t.ImageStores = append(t.ImageStores, *url)
	}
}

func (t *VirtualContainerHostConfigSpec) AddVolumeLocation(name string, url *url.URL) {
	if url != nil {
		t.VolumeLocations[name] = *url
	}
}

// AddComputeResource adds a moref to the set of permitted root pools. It takes a ResourcePool rather than
// an inventory path to encourage validation.
func (t *VirtualContainerHostConfigSpec) AddComputeResource(pool *types.ManagedObjectReference) {
	if pool != nil {
		t.ComputeResources = append(t.ComputeResources, *pool)
	}
}

func (t *VirtualContainerHostConfigSpec) SetvSphereTarget(url *url.URL) {
	if url != nil {
		t.Target = *url
	}
}

func CreateSession(cmd string, args ...string) *SessionConfig {
	cfg := &SessionConfig{
		Cmd: Cmd{
			Path: cmd,
			Args: []string{
				cmd,
			},
		},
	}

	cfg.Cmd.Args = append(cfg.Cmd.Args, args...)

	return cfg
}
