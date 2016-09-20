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
	"fmt"
	"net"

	"github.com/vmware/vic/pkg/uid"
)

type Endpoint struct {
	container *Container
	scope     *Scope
	ip        net.IP
	gateway   net.IP
	subnet    net.IPNet
	static    bool
	ports     map[Port]interface{} // exposed ports
}

func newEndpoint(container *Container, scope *Scope, ip *net.IP, subnet net.IPNet, gateway net.IP, pciSlot *int32) *Endpoint {
	e := &Endpoint{
		container: container,
		scope:     scope,
		gateway:   gateway,
		subnet:    subnet,
		ip:        net.IPv4(0, 0, 0, 0),
		static:    false,
		ports:     make(map[Port]interface{}),
	}

	if ip != nil {
		e.ip = *ip
	}
	if !e.ip.IsUnspecified() {
		e.static = true
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

func (e *Endpoint) addPort(p Port) error {
	if _, ok := e.ports[p]; ok {
		return fmt.Errorf("port %s already exposed", p)
	}

	e.ports[p] = nil
	return nil
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

func (e *Endpoint) Container() *Container {
	return e.container
}

func (e *Endpoint) ID() uid.UID {
	return e.container.ID()
}

func (e *Endpoint) Name() string {
	return e.container.Name()
}

func (e *Endpoint) Gateway() net.IP {
	return e.gateway
}

func (e *Endpoint) Ports() []Port {
	ports := make([]Port, len(e.ports))
	i := 0
	for p := range e.ports {
		ports[i] = p
		i++
	}

	return ports
}

func (e *Endpoint) copy() *Endpoint {
	other := &Endpoint{}
	*other = *e
	other.ports = make(map[Port]interface{})
	for p := range e.ports {
		other.ports[p] = nil
	}

	return other
}
