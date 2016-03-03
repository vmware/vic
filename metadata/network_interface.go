package metadata

import (
	"net"

	"github.com/vmware/govmomi/object"
)

// NetworkEndpoint describes a network presence in the form a vNIC in sufficient detail that it can be:
// a. created - the vNIC added to a VM
// b. identified - the guestOS can determine which interface it corresponds to
// c. configured - the guestOS can configure the interface correctly
type NetworkEndpoint struct {
	// IP addresses to assign - may be empty if DHCP
	IP net.IP
	// The network scope the IP belongs to
	Network net.IPNet
	// Default gateway if any
	Gateway net.IP
	// The set of nameservers associated with this network - may be empty
	Nameservers []net.IP
	// The MAC address for the vNIC - this allows for interface idenitifcaton
	MAC string
	// The backing network identifier
	Backing object.Network
}
