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
	"sync"
)

var (
	defaultSubnet *net.IPNet
)

type Scope struct {
	sync.Mutex

	id         string
	name       string
	scopeType  string
	subnet     net.IPNet
	gateway    net.IP
	dns        []net.IP
	ipam       *IPAM
	containers map[string]*Container
	endpoints  []*Endpoint
	space      *AddressSpace
}

type IPAM struct {
	pools  []string
	spaces []*AddressSpace
}

func init() {
	_, defaultSubnet, _ = net.ParseCIDR("0.0.0.0/16")
}

func (s *Scope) Name() string {
	return s.name
}

func (s *Scope) ID() string {
	return s.id
}

func (s *Scope) Type() string {
	return s.scopeType
}

func (s *Scope) IPAM() *IPAM {
	return s.ipam
}

func (s *Scope) reserveEndpointIP(e *Endpoint) error {
	// reserve an ip address
	var err error
	for _, p := range s.ipam.spaces {
		if !e.ip.IsUnspecified() {
			if err = p.ReserveIP4(e.ip); err == nil {
				break
			}
		} else {
			var ip net.IP
			if ip, err = p.ReserveNextIP4(); err == nil {
				e.ip = ip
				break
			}
		}
	}

	return err
}

func (s *Scope) releaseEndpointIP(e *Endpoint) error {
	for _, p := range s.ipam.spaces {
		if err := p.ReleaseIP4(e.ip); err == nil {
			return nil
		}
	}

	return fmt.Errorf("could not release IP for endpoint")
}

func (s *Scope) AddContainer(name string, ip *net.IP) (*Endpoint, error) {
	if name == "" {
		return nil, fmt.Errorf("empty container name")
	}

	s.Lock()
	defer s.Unlock()

	c, ok := s.containers[name]
	if ok {
		return nil, DuplicateResourceError{resID: name}
	}

	e := newEndpoint(c, s, ip, s.subnet, s.gateway, nil, nil)

	var err error
	err = s.reserveEndpointIP(e)
	defer func() {
		if err != nil {
			s.releaseEndpointIP(e)
		}
	}()

	if err != nil {
		return nil, err
	}

	s.endpoints = append(s.endpoints, e)

	c = &Container{
		name:     name,
		endpoint: e,
	}

	e.container = c
	s.containers[name] = c
	return e, nil
}

func (s *Scope) RemoveContainer(name string) error {
	s.Lock()
	defer s.Unlock()

	c, ok := s.containers[name]
	if !ok {
		return ResourceNotFoundError{}
	}

	var e *Endpoint
	for _, e = range s.endpoints {
		if e.container == c {
			break
		}
	}

	if e == nil || e.container != c {
		return ResourceNotFoundError{}
	}

	err := s.releaseEndpointIP(e)
	if err != nil {
		return err
	}

	s.endpoints = removeEndpointHelper(e, s.endpoints)
	delete(s.containers, name)
	return nil
}

func (s *Scope) Containers() []*Container {
	s.Lock()
	defer s.Unlock()

	containers := make([]*Container, len(s.containers))
	i := 0
	for _, c := range s.containers {
		containers[i] = c
		i++
	}

	return containers
}

func (s *Scope) Container(name string) (*Container, error) {
	s.Lock()
	defer s.Unlock()

	c, ok := s.containers[name]
	if !ok {
		return nil, ResourceNotFoundError{}
	}

	return c, nil
}

func (s *Scope) Endpoints() []*Endpoint {
	return s.endpoints
}

func (s *Scope) Subnet() *net.IPNet {
	return &s.subnet
}

func (s *Scope) Gateway() net.IP {
	return s.gateway
}

func (s *Scope) DNS() []net.IP {
	return s.dns
}

func (i *IPAM) Pools() []string {
	return i.pools
}
