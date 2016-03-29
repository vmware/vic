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
	"reflect"
	"testing"
)

type params struct {
	scopeType, name string
	subnet          *net.IPNet
	gateway         net.IP
	dns             []net.IP
	ipam            []string
}

var validScopeTests = []struct {
	in  params
	out *params
	err error
}{
	// bridge scopes

	// default bridge pool, only name specified
	{params{"bridge", "bar1", nil, net.IPv4(0, 0, 0, 0), nil, nil},
		&params{"bridge", "bar1", &net.IPNet{IP: net.IPv4(172, 17, 0, 0), Mask: net.CIDRMask(16, 32)}, net.ParseIP("172.17.0.1"), nil, nil},
		nil},
	// default bridge pool with gateway specified
	{params{"bridge", "bar2", nil, net.IPv4(172, 18, 0, 2), nil, nil},
		&params{"bridge", "bar2", &net.IPNet{IP: net.IPv4(172, 18, 0, 0), Mask: net.CIDRMask(16, 32)}, net.ParseIP("172.18.0.2"), nil, nil},
		nil},
	// not from default bridge pool
	{params{"bridge", "bar3", &net.IPNet{IP: net.ParseIP("10.10.0.0"), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, nil},
		&params{"bridge", "bar3", &net.IPNet{IP: net.ParseIP("10.10.0.0"), Mask: net.CIDRMask(16, 32)}, net.ParseIP("10.10.0.1"), nil, nil},
		nil},
	// not from default bridge pool, dns specified
	{params{"bridge", "bar4", &net.IPNet{IP: net.ParseIP("10.11.0.0"), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), []net.IP{net.ParseIP("10.10.1.1")}, nil},
		&params{"bridge", "bar4", &net.IPNet{IP: net.ParseIP("10.11.0.0"), Mask: net.CIDRMask(16, 32)}, net.ParseIP("10.11.0.1"), []net.IP{net.ParseIP("10.10.1.1")}, nil},
		nil},
	// not from default pool, dns and ipam specified
	{params{"bridge", "bar5", &net.IPNet{IP: net.ParseIP("10.12.0.0"), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), []net.IP{net.ParseIP("10.10.1.1")}, []string{"10.12.1.0/24", "10.12.2.0-10.12.2.15"}},
		&params{"bridge", "bar5", &net.IPNet{IP: net.ParseIP("10.12.0.0"), Mask: net.CIDRMask(16, 32)}, net.ParseIP("10.12.1.0"), []net.IP{net.ParseIP("10.10.1.1")}, nil},
		nil},
	// from default pool, subnet specified
	{params{"bridge", "bar6", &net.IPNet{IP: net.IPv4(172, 19, 0, 0), Mask: net.CIDRMask(16, 32)}, nil, nil, nil},
		&params{"bridge", "bar6", &net.IPNet{IP: net.IPv4(172, 19, 0, 0), Mask: net.CIDRMask(16, 32)}, net.ParseIP("172.19.0.1"), nil, nil},
		nil},

	// external scopes
	{params{"external", "bar7", &net.IPNet{IP: net.ParseIP("10.13.0.0"), Mask: net.CIDRMask(16, 32)}, net.ParseIP("10.13.0.1"), []net.IP{net.ParseIP("10.10.1.1")}, []string{"10.13.1.0/24", "10.13.2.0-10.13.2.15"}},
		&params{"external", "bar7", &net.IPNet{IP: net.ParseIP("10.13.0.0"), Mask: net.CIDRMask(16, 32)}, net.ParseIP("10.13.0.1"), []net.IP{net.ParseIP("10.10.1.1")}, nil},
		nil},
}

func TestContext(t *testing.T) {
	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32))
	if err != nil {
		t.Errorf("NewContext() => (nil, %s), want (ctx, nil)", err)
		return
	}

	var tests = []struct {
		in  params
		out *params
		err error
	}{
		// empty name
		{params{"bridge", "", nil, net.IPv4(0, 0, 0, 0), nil, nil}, nil, nil},
		// unsupported network type
		{params{"foo", "bar8", nil, net.IPv4(0, 0, 0, 0), nil, nil}, nil, nil},
		// duplicate name
		{params{"bridge", "bar6", nil, net.IPv4(0, 0, 0, 0), nil, nil}, nil, DuplicateResourceError{}},
		// ip range already allocated
		{params{"bridge", "bar9", &net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, nil}, nil, nil},
		// ipam out of range of network
		{params{"bridge", "bar10", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, []string{"10.14.1.0/24", "10.15.1.0/24"}}, nil, nil},
		// this should succeed now
		{params{"bridge", "bar11", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, []string{"10.14.1.0/24"}},
			&params{"bridge", "bar11", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(10, 14, 1, 0), nil, nil},
			nil},
		// bad ipam
		{params{"bridge", "bar12", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, []string{"10.14.1.0/24", "10.15.1"}}, nil, nil},
		// bad ipam, default bridge pool
		{params{"bridge", "bar12", &net.IPNet{IP: net.IPv4(172, 21, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, []string{"172.21.1.0/24", "10.15.1"}}, nil, nil},
		// external networks must have subnet specified
		{params{"external", "bar13", nil, net.IPv4(0, 0, 0, 0), nil, []string{"10.15.0.0/24"}}, nil, nil},
		// external networks must have gateway specified
		{params{"external", "bar14", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, []string{"10.15.0.0/24"}}, nil, nil},
		// external networks must have ipam specified
		{params{"external", "bar15", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(10, 14, 0, 1), nil, nil}, nil, nil},
		// external networks cannot overlap bridge pool
		{params{"external", "bar16", &net.IPNet{IP: net.IPv4(172, 20, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(10, 14, 0, 1), nil, []string{"172.20.0.0/16"}}, nil, nil},
	}

	tests = append(validScopeTests, tests...)

	for _, te := range tests {
		s, err := ctx.NewScope(te.in.scopeType,
			te.in.name,
			te.in.subnet,
			te.in.gateway,
			te.in.dns,
			te.in.ipam)

		if te.out == nil {
			// error case
			if s != nil || err == nil {
				t.Errorf("NewScope() => (s, nil), want (nil, err)")
			}

			// if there is an error specified, check if we got that error
			if te.err != nil &&
				reflect.TypeOf(err) != reflect.TypeOf(te.err) {
				t.Errorf("NewScope() => (nil, %s), want (nil, %s)", reflect.TypeOf(err), reflect.TypeOf(te.err))
			}

			if _, o := err.(DuplicateResourceError); !o {
				// sanity check
				if _, ok := ctx.scopes[te.in.name]; ok {
					t.Errorf("scope %s added on error", te.in.name)
				}
			}

			continue
		}

		if err != nil {
			t.Errorf("got: %s, expected: nil", err)
			continue
		}

		if s.Type() != te.out.scopeType {
			t.Errorf("s.Type() => %s, want %s", s.Type(), te.out.scopeType)
			continue
		}

		if s.Name() != te.out.name {
			t.Errorf("s.Name() => %s, want %s", s.Name(), te.out.name)
		}

		if s.Subnet().String() != te.out.subnet.String() {
			t.Errorf("s.Subnet() => %s, want %s", s.Subnet(), te.out.subnet)
		}

		if !s.Gateway().Equal(te.out.gateway) {
			t.Errorf("s.Gateway() => %s, want %s", s.Gateway(), te.out.gateway)
		}

		for _, d1 := range s.DNS() {
			found := false
			for _, d2 := range te.out.dns {
				if d2.Equal(d1) {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("s.DNS() => %q, want %q", s.DNS(), te.out.dns)
				break
			}
		}

		ipam := s.IPAM()
		if ipam == nil {
			t.Errorf("s.IPAM() == nil, want %q", te.in.ipam)
			continue
		}

		if te.in.ipam != nil && len(ipam.spaces) != len(te.in.ipam) {
			t.Errorf("len(ipam.spaces) => %d != len(te.in.ipam) => %d", len(ipam.spaces), len(te.in.ipam))
		}

		for i, p := range ipam.spaces {
			if te.in.ipam == nil {
				if p != s.space {
					t.Errorf("got %v, want %v", p, s.space)
				}
				continue
			}

			if p.Parent != s.space {
				t.Errorf("p.Parent => %v, want %v", p.Parent, s.space)
				continue
			}

			if p.Network != nil {
				if p.Network.String() != te.in.ipam[i] {
					t.Errorf("p.Network => %s, want %s", p.Network, te.in.ipam[i])
				}
			} else if p.Pool.String() != te.in.ipam[i] {
				t.Errorf("p.Pool => %s, want %s", p.Pool, te.in.ipam[i])
			}
		}
	}
}

func TestScopes(t *testing.T) {
	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32))
	if err != nil {
		t.Errorf("NewContext() => (nil, %s), want (ctx, nil)", err)
		return
	}

	scopes := make([]*Scope, 0, 0)
	scopesByID := make(map[string]*Scope)
	scopesByName := make(map[string]*Scope)
	for _, te := range validScopeTests {
		s, err := ctx.NewScope(
			te.in.scopeType,
			te.in.name,
			te.in.subnet,
			te.in.gateway,
			te.in.dns,
			te.in.ipam)

		if err != nil {
			t.Errorf("NewScope() => (_, %s), want (_, nil)", err)
		}

		scopesByID[s.ID()] = s
		scopesByName[s.Name()] = s
		scopes = append(scopes, s)
	}

	id := scopesByName[validScopeTests[0].in.name].ID()
	partialID := scopesByName[validScopeTests[0].in.name].ID()[:8]
	partialID2 := partialID[1:]
	badName := "foo"

	var tests = []struct {
		in  *string
		out []*Scope
	}{
		// name match
		{&validScopeTests[0].in.name, []*Scope{scopesByName[validScopeTests[0].in.name]}},
		// id match
		{&id, []*Scope{scopesByName[validScopeTests[0].in.name]}},
		// partial id match
		{&partialID, []*Scope{scopesByName[validScopeTests[0].in.name]}},
		// all scopes
		{nil, scopes},
		// partial id match only matches prefix
		{&partialID2, nil},
		// no match
		{&badName, nil},
	}

	for _, te := range tests {
		l, err := ctx.Scopes(te.in)
		if te.out == nil {
			if err == nil {
				t.Errorf("Scopes() => (_, nil), want (_, err)")
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Scopes() => (_, %s), want (_, nil)", err)
				continue
			}
		}

		// +1 for the default bridge scope
		if te.in == nil {
			if len(l) != len(te.out)+1 {
				t.Errorf("len(scopes) => %d != %d", len(l), len(te.out)+1)
				continue
			}
		} else {
			if len(l) != len(te.out) {
				t.Errorf("len(scopes) => %d != %d", len(l), len(te.out))
				continue
			}
		}

		for _, s1 := range te.out {
			found := false
			for _, s2 := range l {
				if s1 == s2 {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("got=%v, want=%v", l, te.out)
				break
			}
		}
	}
}
