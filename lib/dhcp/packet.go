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

package dhcp

import (
	"bytes"
	"net"
	"time"

	"encoding/binary"

	"github.com/d2g/dhcp4"
)

const DefaultLeaseTime = 1 * time.Hour

// Packet is a representation of a DHCP packet
type Packet struct {
	packet  []byte
	options dhcp4.Options
}

// YourIP returns the YIP field in the packet (this is IP assigned by server to the client)
func (p *Packet) YourIP() net.IP {
	if len(p.packet) == 0 {
		return nil
	}

	return dhcp4.Packet(p.packet).YIAddr()
}

// Gateway return the GIP field in the packet (server assigned gateway)
func (p *Packet) Gateway() net.IP {
	if len(p.packet) == 0 {
		return nil
	}

	b := p.options[dhcp4.OptionRouter]
	if len(b) >= 4 {
		return net.IP(b[:4])
	}

	return nil
}

// SubnetMask returns the subnet mask option in the packet
func (p *Packet) SubnetMask() net.IPMask {
	return p.options[dhcp4.OptionSubnetMask]
}

// LeaseTime returns the lease time (in seconds) in the packet
func (p *Packet) LeaseTime() time.Duration {
	b := p.options[dhcp4.OptionIPAddressLeaseTime]
	if b == nil {
		return DefaultLeaseTime
	}

	var t uint32
	if err := binary.Read(bytes.NewReader(b), binary.BigEndian, &t); err != nil {
		return 0 * time.Second
	}

	return time.Duration(t) * time.Second
}

// DNS returns the name server entries in the dhcp packet
func (p *Packet) DNS() []net.IP {
	b := p.options[dhcp4.OptionDomainNameServer]
	if b == nil {
		return nil
	}

	var dns []net.IP
	for i := 0; i < len(b); i += 4 {
		dns = append(dns, net.IP(b[i:i+4]))
	}

	return dns
}

// ServerIP returns the DHCP server's IP address
func (p *Packet) ServerIP() net.IP {
	if len(p.packet) == 0 {
		return nil
	}

	b := p.options[dhcp4.OptionServerIdentifier]
	if len(b) < net.IPv4len {
		return nil
	}

	return net.IP(b[:net.IPv4len])
}
