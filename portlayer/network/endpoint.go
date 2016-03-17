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

import "net"

type Endpoint struct {
	id        string
	container *Container
	scope     *Scope
	ip        net.IP
	link      string
	mac       string
	gateway   net.IP
	subnet    net.IPNet
}

func newEndpoint(container *Container, scope *Scope, ip *net.IP, subnet net.IPNet, gateway net.IP, link *string, mac *string) *Endpoint {
	e := &Endpoint{
		id:        generateID(),
		container: container,
		scope:     scope,
		gateway:   gateway,
		subnet:    subnet,
		ip:        net.IPv4(0, 0, 0, 0),
	}

	if ip != nil {
		e.ip = *ip
	}
	if link != nil {
		e.link = *link
	}
	if mac != nil {
		e.mac = *mac
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

func (e *Endpoint) Link() string {
	return e.link
}

func (e *Endpoint) Scope() *Scope {
	return e.scope
}

func (e *Endpoint) Subnet() *net.IPNet {
	return &e.subnet
}

func (e *Endpoint) Mac() string {
	return e.mac
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
