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

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/uid"
)

type Scope struct {
	sync.RWMutex

	id         uid.UID
	name       string
	scopeType  string
	subnet     net.IPNet
	gateway    net.IP
	dns        []net.IP
	ipam       *IPAM
	containers map[uid.UID]*Container
	endpoints  []*Endpoint
	space      *AddressSpace
	builtin    bool
	network    object.NetworkReference
}

type IPAM struct {
	pools  []string
	spaces []*AddressSpace
}

func (s *Scope) Name() string {
	s.RLock()
	defer s.RUnlock()

	return s.name
}

func (s *Scope) ID() uid.UID {
	s.RLock()
	defer s.RUnlock()

	return s.id
}

func (s *Scope) Type() string {
	s.RLock()
	defer s.RUnlock()

	return s.scopeType
}

func (s *Scope) IPAM() *IPAM {
	s.RLock()
	defer s.RUnlock()

	return s.ipam
}

func (s *Scope) Network() object.NetworkReference {
	s.RLock()
	defer s.RUnlock()

	return s.network
}

func (s *Scope) isDynamic() bool {
	return s.scopeType != constants.BridgeScopeType && s.ipam.spaces == nil
}

func (s *Scope) reserveEndpointIP(e *Endpoint) error {
	if s.isDynamic() {
		return nil
	}

	// reserve an ip address
	var err error
	for _, p := range s.ipam.spaces {
		if !ip.IsUnspecifiedIP(e.ip) {
			if err = p.ReserveIP4(e.ip); err == nil {
				return nil
			}
		} else {
			var eip net.IP
			if eip, err = p.ReserveNextIP4(); err == nil {
				e.ip = eip
				return nil
			}
		}
	}

	return err
}

func (s *Scope) releaseEndpointIP(e *Endpoint) error {
	if s.isDynamic() {
		return nil
	}

	for _, p := range s.ipam.spaces {
		if err := p.ReleaseIP4(e.ip); err == nil {
			if !e.static {
				e.ip = net.IPv4(0, 0, 0, 0)
			}
			return nil
		}
	}

	return fmt.Errorf("could not release IP for endpoint")
}

func (s *Scope) AddContainer(con *Container, e *Endpoint) error {
	s.Lock()
	defer s.Unlock()

	if con == nil {
		return fmt.Errorf("container is nil")
	}

	_, ok := s.containers[con.id]
	if ok {
		return DuplicateResourceError{resID: con.id.String()}
	}

	if err := s.reserveEndpointIP(e); err != nil {
		return err
	}

	con.addEndpoint(e)
	s.endpoints = append(s.endpoints, e)
	s.containers[con.id] = con
	return nil
}

func (s *Scope) RemoveContainer(con *Container) error {
	s.Lock()
	defer s.Unlock()

	c, ok := s.containers[con.id]
	if !ok || c != con {
		return ResourceNotFoundError{}
	}

	e := c.Endpoint(s)
	if e == nil {
		return ResourceNotFoundError{}
	}

	if err := s.releaseEndpointIP(e); err != nil {
		return err
	}

	delete(s.containers, c.id)
	s.endpoints = removeEndpointHelper(e, s.endpoints)
	c.removeEndpoint(e)
	return nil
}

func (s *Scope) Containers() []*Container {
	s.RLock()
	defer s.RUnlock()

	containers := make([]*Container, len(s.containers))
	i := 0
	for _, c := range s.containers {
		containers[i] = c
		i++
	}

	return containers
}

func (s *Scope) Container(id uid.UID) *Container {
	s.RLock()
	defer s.RUnlock()

	if c, ok := s.containers[id]; ok {
		return c
	}

	return nil
}

func (s *Scope) ContainerByAddr(addr net.IP) *Endpoint {
	s.RLock()
	defer s.RUnlock()

	if addr == nil || addr.IsUnspecified() {
		return nil
	}

	for _, e := range s.endpoints {
		if addr.Equal(e.IP()) {
			return e
		}
	}

	return nil
}

func (s *Scope) Endpoints() []*Endpoint {
	s.RLock()
	defer s.RUnlock()

	eps := make([]*Endpoint, len(s.endpoints))
	copy(eps, s.endpoints)
	return eps
}

func (s *Scope) Subnet() *net.IPNet {
	s.RLock()
	defer s.RUnlock()

	return &s.subnet
}

func (s *Scope) Gateway() net.IP {
	s.RLock()
	defer s.RUnlock()

	return s.gateway
}

func (s *Scope) DNS() []net.IP {
	s.RLock()
	defer s.RUnlock()

	return s.dns
}

func (s *Scope) Refresh(h *exec.Handle) error {
	s.Lock()
	defer s.Unlock()

	if !s.isDynamic() {
		return nil
	}

	ne := h.ExecConfig.Networks[s.name]
	if ip.IsUnspecifiedSubnet(&ne.Network.Gateway) {
		return fmt.Errorf("updating container %s: gateway not present for scope %s", h.ExecConfig.ID, s.name)
	}

	gw, snet, err := net.ParseCIDR(ne.Network.Gateway.String())
	if err != nil {
		return err
	}

	s.gateway = gw
	s.subnet = *snet

	return nil
}

func (i *IPAM) Pools() []ip.Range {
	var pools []ip.Range
	for _, s := range i.spaces {
		pools = append(pools, *s.Pool)
	}

	return pools
}
