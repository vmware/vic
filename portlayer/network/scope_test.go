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
	"errors"
	"net"
	"reflect"
	"testing"

	"github.com/vmware/vic/pkg/vsphere/session"
)

func makeIP(a, b, c, d byte) *net.IP {
	i := net.IPv4(a, b, c, d)
	return &i
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
		name string
		ip   *net.IP
		out  *Endpoint
		err  error
	}{
		// empty container id
		{"", nil, nil, errors.New("")},
		// add a new container to scope
		{"foo", nil, &Endpoint{ip: net.IPv4(172, 16, 0, 2), subnet: s.subnet, gateway: s.gateway}, nil},
		// container already part of scope
		{"foo", nil, nil, DuplicateResourceError{}},
		// container with ip
		{"bar", makeIP(172, 16, 0, 3), &Endpoint{ip: net.IPv4(172, 16, 0, 3), subnet: s.subnet, gateway: s.gateway}, nil},
		// container with ip that is not available
		{"baz", makeIP(172, 16, 0, 3), nil, errors.New("")},
		{"baz", makeIP(172, 16, 0, 0), nil, errors.New("")},
		{"baz", makeIP(172, 16, 255, 255), nil, errors.New("")},
	}

	for _, te := range tests1 {
		t.Logf("testing name = %s, ip = %v", te.name, te.ip)
		e, err := s.AddContainer(te.name, te.ip)
		if te.err != nil {
			if err == nil {
				t.Errorf("s.AddContainer() => (_, nil), want (_, err)")
				continue
			}

			if reflect.TypeOf(err) != reflect.TypeOf(te.err) {
				t.Errorf("s.AddContainer() => (_, %v), want (_, %v)", reflect.TypeOf(err), reflect.TypeOf(te.err))
				continue
			}

			// for any other error other than DuplicateResourcError
			// verify that the container was not added
			if _, ok := err.(DuplicateResourceError); !ok {
				c, err := s.Container(te.name)
				if _, ok := err.(ResourceNotFoundError); !ok || c != nil {
					t.Errorf("s.Container(%s) => (%v, %v), want (nil, err)", te.name, c, err)
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

		if e.container.Name() != te.name {
			t.Errorf("s.AddContainer() => e.container == %s, want e.container == %s", e.container.Name(), te.name)
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

		c, err := s.Container(te.name)
		if err != nil {
			t.Errorf("s.Container(%s) => (nil, %s), want (c, nil)", te.name, err)
			continue
		}

		if c.Endpoint() != e {
			t.Errorf("container %s does not contain %v", te.name, e)
		}
	}

	// test RemoveContainer
	var tests2 = []struct {
		name string
		err  error
	}{
		{"", ResourceNotFoundError{}},
		{"c1", ResourceNotFoundError{}},
		{"foo", nil},
	}

	for _, te := range tests2 {
		err = s.RemoveContainer(te.name)
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

		_, ok := s.containers[te.name]
		if ok {
			t.Errorf("s.RemoveContainer() did not remove container %s", te.name)
			continue
		}

		for _, e := range s.endpoints {
			if e.container.Name() == te.name {
				t.Errorf("s.RemoveContainer() did not remove endpoint for container %s", te.name)
				break
			}
		}

	}
}
