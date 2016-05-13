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

package network

import (
	"net"

	"github.com/vmware/vic/metadata"
)

type Endpoint struct {
	id        string
	container *Container
	scope     *Scope
	ip        net.IP
	pciSlot   int32
	gateway   net.IP
	subnet    net.IPNet
	static    bool
	bound     bool
}

func newEndpoint(container *Container, scope *Scope, ip *net.IP, subnet net.IPNet, gateway net.IP, pciSlot *int32) *Endpoint {
	e := &Endpoint{
		id:        generateID(),
		container: container,
		scope:     scope,
		gateway:   gateway,
		subnet:    subnet,
		ip:        net.IPv4(0, 0, 0, 0),
		static:    false,
		bound:     false,
	}

	if ip != nil {
		e.ip = *ip
	}
	if pciSlot != nil {
		e.pciSlot = *pciSlot
	}

	return e
}

func removeEndpointHelper(ep *Endpoint, eps []*Endpoint) []*Endpoint {
	for i, e := range eps {
		if ep != e {
			continue
		}

		return append(eps[:i], eps[i+1:]...)
	}

	return eps
}

func (e *Endpoint) IP() net.IP {
	return e.ip
}

func (e *Endpoint) Scope() *Scope {
	return e.scope
}

func (e *Endpoint) Subnet() *net.IPNet {
	return &e.subnet
}

func (e *Endpoint) PciSlot() int32 {
	return e.pciSlot
}

func (e *Endpoint) Container() *Container {
	return e.container
}

func (e *Endpoint) ID() string {
	return e.id
}

func (e *Endpoint) Gateway() net.IP {
	return e.gateway
}

func (e *Endpoint) metadataEndpoint() *metadata.NetworkEndpoint {
	ne := &metadata.NetworkEndpoint{
		IP: net.IPNet{
			IP:   e.ip,
			Mask: e.subnet.Mask,
		},
		PCISlot: e.pciSlot,
		Network: metadata.ContainerNetwork{
			Name: e.scope.name,
		},
	}
	ne.Network.Gateway = net.IPNet{IP: e.gateway, Mask: e.subnet.Mask}
	// ip should be unspecified if endpoint is not bound
	if !e.bound {
		ne.IP.IP = net.IPv4zero
	}

	return ne
}

func (e *Endpoint) IsBound() bool {
	return e.bound
}

func (e *Endpoint) bind(s *Scope) error {
	if e.bound {
		return nil // already bound
	}

	if err := s.reserveEndpointIP(e); err != nil {
		return err
	}

	e.bound = true
	return nil
}

func (e *Endpoint) unbind(s *Scope) error {
	if !e.IsBound() {
		return nil // not bound
	}

	if err := s.releaseEndpointIP(e); err != nil {
		return err
	}

	e.bound = false
	return nil
}
