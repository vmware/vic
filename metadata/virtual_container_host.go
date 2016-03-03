package metadata

import (
	"crypto/tls"
	"net/mail"
	"net/url"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/oauth2"
)

// VirtualContainerHostConfigSpec holds the metadata for a Virtual Container Host that should be visible inside the appliance VM.
type VirtualContainerHostConfigSpec struct {
	// Virtual Container Host version
	Version string

	// The unambiguous identifer for the Virtual Container Host within a vSphere environment
	MOID types.ManagedObjectReference
	// Descriptive display name
	Name string
	// Administrative contact for the Virtual Container Host
	VCHAdmin []mail.Address
	// Administrative contact for vSphere
	VSphereAdmin []mail.Address

	// Certificates for user authentication
	UserCertificates []tls.Certificate
	// Certificates for general network access, keyed by CIDR (IPNet.String())
	NetworkCertificates map[string]tls.Certificate
	// Certificates for specific system access, keyed by FQDN
	HostCertificates map[string]tls.Certificate

	// Datastore URL for image store
	ImageStore url.URL
	// Datastore URL roots for volumes
	VolumeLocations []url.URL
	// Compute resource root
	ComputeResource object.ComputeResource

	// Networks to be attached to the Virtual Container Host, keyed by interface name
	ApplianceNetworks map[string]NetworkEndpoint
	// Published networks available for containers to join, keyed by consumption name
	ContainerNetworks map[string]NetworkEndpoint
	// The default bridge network supplied for the Virtual Container Host
	BridgeNetwork object.Network

	// Whitelist of registries
	RegistryWhitelist []url.URL
	// Blacklist of registries
	RegistryBlacklist []url.URL
}

// CustomerExperienceImprovementProgram provides configuration for the phone home mechanism
// This is broken out so that we can have more granular configuration in here in the future
// and so that it is insulated from changes to Virtual Container Host structure
type CustomerExperienceImprovementProgram struct {
	// Is phone home reporting enabled at all
	Enabled bool
}

// Credentials is an exceptionally basic container for credentials
type Credentials struct {
	User        string
	Password    []byte
	Certificate []byte
	Key         []byte
}

// VirtualContainerHostTargetConfigSpec holds the metadata for a Virtual Container Host that should be visible only to the infrastructure on which it is running.
type VirtualContainerHostTargetConfigSpec struct {
	// Target and credentials for the environment - experimenting with storing this as an OAuth2 config
	TargetConfig oauth2.Config
	// Freeform notes for target administrator user
	Notes string
}
