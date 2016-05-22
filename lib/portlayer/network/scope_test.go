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
	"reflect"
	"testing"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/vsphere/session"
)

func makeIP(a, b, c, d byte) *net.IP {
	i := net.IPv4(a, b, c, d)
	return &i
}

var addEthernetCardOrig = addEthernetCard
var addEthernetCardErr = func(_ *exec.Handle, _ *Scope) (types.BaseVirtualDevice, error) {
	return nil, fmt.Errorf("")
}

func TestScopeAddRemoveContainer(t *testing.T) {
	origBridgeNetworkName := getBridgeNetworkName
	getBridgeNetworkName = mockBridgeNetworkName
	defer func() { getBridgeNetworkName = origBridgeNetworkName }()

	var err error
	sess := &session.Session{}

	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32), sess)
	if err != nil {
		t.Errorf("NewContext() => (nil, %s), want (ctx, nil)", err)
		return
	}

	s := ctx.defaultScope

	var tests1 = []struct {
		c   *Container
		ip  *net.IP
		out *Endpoint
		err error
	}{
		// no container
		{nil, nil, nil, fmt.Errorf("")},
		// add a new container to scope
		{&Container{id: "foo"}, nil, &Endpoint{ip: net.IPv4(0, 0, 0, 0), subnet: s.subnet, gateway: s.gateway}, nil},
		// container already part of scope
		{&Container{id: "foo"}, nil, nil, DuplicateResourceError{}},
		// container with ip
		{&Container{id: "bar"}, makeIP(172, 16, 0, 3), &Endpoint{ip: net.IPv4(172, 16, 0, 3), subnet: s.subnet, gateway: s.gateway, static: true}, nil},
	}

	for _, te := range tests1 {
		t.Logf("testing id = %v, ip = %v", te.c, te.ip)
		var e *Endpoint
		e, err = s.addContainer(te.c, te.ip)
		if te.err != nil {
			if err == nil {
				t.Errorf("s.AddContainer() => (_, nil), want (_, err)")
				continue
			}

			if reflect.TypeOf(err) != reflect.TypeOf(te.err) {
				t.Errorf("s.AddContainer() => (_, %v), want (_, %v)", reflect.TypeOf(err), reflect.TypeOf(te.err))
				continue
			}

			if te.c == nil {
				continue
			}

			// for any other error other than DuplicateResourcError
			// verify that the container was not added
			if _, ok := err.(DuplicateResourceError); !ok {
				c := s.Container(te.c.ID())
				if c != nil {
					t.Errorf("s.Container(%s) => (%v, %v), want (nil, err)", te.c.ID(), c, err)
				}
			}

			continue
		}

		if !e.IP().Equal(te.out.IP()) {
			t.Errorf("s.AddContainer() => e.IP() == %v, want e.IP() == %v", e.IP(), te.out.IP())
			continue
		}

		if !e.Gateway().Equal(te.out.Gateway()) {
			t.Errorf("s.AddContainer() => e.Gateway() == %v, want e.Gateway() == %v", e.Gateway(), te.out.Gateway())
			continue
		}

		if e.subnet.String() != s.subnet.String() {
			t.Errorf("s.AddContainer() => e.subnet == %s, want e.subnet == %s", e.subnet, s.subnet)
			continue
		}

		if e.static != te.out.static {
			t.Errorf("s.AddContainer() => e.static == %#v, want e.static == %#v", e.static, te.out.static)
		}

		if e.container.ID() != te.c.ID() {
			t.Errorf("s.AddContainer() => e.container == %s, want e.container == %s", e.container.ID(), te.c.ID())
			continue
		}

		found := false
		for _, e1 := range s.Endpoints() {
			if e1 == e {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("s.endpoints does not contain %v", e)
		}

		c := s.Container(te.c.id)
		if c == nil {
			t.Errorf("s.Container(%s) => nil, want %v", te.c.ID(), te.c)
			continue
		}

		if c.Endpoint(s) != e {
			t.Errorf("container %s does not contain %v", te.c.ID(), e)
		}
	}

	bound := exec.NewContainer("bound")
	ctx.AddContainer(bound, ctx.defaultScope.Name(), nil)
	ctx.BindContainer(bound)

	// test RemoveContainer
	var tests2 = []struct {
		c   *Container
		err error
	}{
		// container not found
		{&Container{id: "c1"}, ResourceNotFoundError{}},
		// try to remove a bound container
		{s.Container(bound.Container.ID), fmt.Errorf("")},
		// remove a container
		{s.Container(exec.ParseID("foo")), nil},
	}

	for _, te := range tests2 {
		err = s.removeContainer(te.c)
		if te.err != nil {
			if err == nil {
				t.Errorf("s.RemoveContainer() => nil, want %v", te.err)
			}

			continue
		}

		// container was removed, verify
		if err != nil {
			t.Errorf("s.RemoveContainer() => %s, want nil", err)
			continue
		}

		c := s.Container(te.c.ID())
		if c != nil {
			t.Errorf("s.RemoveContainer() did not remove container %s", te.c.ID())
			continue
		}

		for _, e := range s.endpoints {
			if e.container.ID() == te.c.ID() {
				t.Errorf("s.RemoveContainer() did not remove endpoint for container %s", te.c.ID())
				break
			}
		}

	}
}

func TestScopeBindUnbindContainer(t *testing.T) {
	var err error
	sess := &session.Session{}

	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32), sess)
	if err != nil {
		t.Fatalf("NewContext() => (nil, %s), want (ctx, nil)", err)
	}

	// add a container that is not part of the scope
	e, err := ctx.AddContainer(exec.NewContainer("notAdded"), ctx.DefaultScope().Name(), nil)
	notAdded := e.Container()
	ctx.DefaultScope().removeContainer(notAdded)

	var tests = []struct {
		i   int
		c   *Container
		err error
	}{
		{0, nil, fmt.Errorf("")},
		// bind a container that is not part of scope
		{1, notAdded, fmt.Errorf("")},
	}

	for _, te := range tests {
		err = ctx.DefaultScope().bindContainer(te.c)
		if te.err != nil {
			if te.err == nil {
				t.Fatalf("%d: Scope.bindContainer(%s) => nil, want err", te.i, te.c.ID())
			}

			continue
		}
	}

	tests = []struct {
		i   int
		c   *Container
		err error
	}{
		{0, nil, fmt.Errorf("")},
		{1, notAdded, fmt.Errorf("")},
	}

	for _, te := range tests {
		err = ctx.DefaultScope().unbindContainer(te.c)
		if te.err != nil {
			if te.err == nil {
				t.Fatalf("%d: Scope.unbindContainer(%s) => nil, want err", te.i, te.c.ID())
			}

			continue
		}
	}
}
