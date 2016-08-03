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
	"bytes"
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"

	"golang.org/x/net/context"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/ip"
)

var testBridgeNetwork object.NetworkReference

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
	// not from default pool, dns, gateway, and ipam specified
	{params{"bridge", "bar51", &net.IPNet{IP: net.ParseIP("10.33.0.0"), Mask: net.CIDRMask(16, 32)}, net.IPv4(10, 33, 0, 1), []net.IP{net.ParseIP("10.10.1.1")}, []string{"10.33.0.0/16"}},
		&params{"bridge", "bar51", &net.IPNet{IP: net.ParseIP("10.33.0.0"), Mask: net.CIDRMask(16, 32)}, net.ParseIP("10.33.0.1"), []net.IP{net.ParseIP("10.10.1.1")}, nil},
		nil},
	// from default pool, subnet specified
	{params{"bridge", "bar6", &net.IPNet{IP: net.IPv4(172, 19, 0, 0), Mask: net.CIDRMask(16, 32)}, nil, nil, nil},
		&params{"bridge", "bar6", &net.IPNet{IP: net.IPv4(172, 19, 0, 0), Mask: net.CIDRMask(16, 32)}, net.ParseIP("172.19.0.1"), nil, nil},
		nil},
}

type mockLink struct{}

func (l *mockLink) AddrAdd(_ net.IPNet) error {
	return nil
}

func (l *mockLink) AddrDel(_ net.IPNet) error {
	return nil
}

func (l *mockLink) Attrs() *LinkAttrs {
	return &LinkAttrs{Name: "foo"}
}

func TestMain(m *testing.M) {
	n := object.NewNetwork(nil, types.ManagedObjectReference{})
	n.InventoryPath = "testBridge"
	testBridgeNetwork = n

	n = object.NewNetwork(nil, types.ManagedObjectReference{})
	n.InventoryPath = "testExternal"
	testExternalNetwork := n

	Config = Configuration{
		BridgeLink:    &mockLink{},
		BridgeNetwork: "bridge",
		ContainerNetworks: map[string]*ContainerNetwork{
			"bridge": &ContainerNetwork{
				PortGroup: testBridgeNetwork,
			},
			"bar7": &ContainerNetwork{
				Common: executor.Common{
					Name: "external",
				},
				Gateway:     net.IPNet{IP: net.ParseIP("10.13.0.1"), Mask: net.CIDRMask(16, 32)},
				Nameservers: []net.IP{net.ParseIP("10.10.1.1")},
				Pools:       []ip.Range{*ip.ParseRange("10.13.1.0-255"), *ip.ParseRange("10.13.2.0-10.13.2.15")},
				PortGroup:   testExternalNetwork,
			},
			"bar71": &ContainerNetwork{
				Common: executor.Common{
					Name: "external",
				},
				Gateway:     net.IPNet{IP: net.ParseIP("10.131.0.1"), Mask: net.CIDRMask(16, 32)},
				Nameservers: []net.IP{net.ParseIP("10.131.0.1"), net.ParseIP("10.131.0.2")},
				Pools:       []ip.Range{*ip.ParseRange("10.131.1.0/16")},
				PortGroup:   testExternalNetwork,
			},
			"bar72": &ContainerNetwork{
				Common: executor.Common{
					Name: "external",
				},
				PortGroup: testExternalNetwork,
			},
			"bar73": &ContainerNetwork{
				Common: executor.Common{
					Name: "external",
				},
				Gateway:   net.IPNet{IP: net.ParseIP("10.133.0.1"), Mask: net.CIDRMask(16, 32)},
				PortGroup: testExternalNetwork,
			},
		},
	}

	rc := m.Run()

	os.Exit(rc)
}

func TestMapExternalNetworks(t *testing.T) {
	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32))
	if err != nil {
		t.Fatalf("NewContext() => (nil, %s), want (ctx, nil)", err)
	}

	// check if external networks were loaded
	for n, nn := range Config.ContainerNetworks {
		scopes, err := ctx.Scopes(&n)
		if err != nil || len(scopes) != 1 {
			t.Fatalf("external network %s was not loaded", n)
		}

		s := scopes[0]
		pools := s.IPAM().Pools()
		if !ip.IsUnspecifiedIP(nn.Gateway.IP) {
			subnet := &net.IPNet{IP: nn.Gateway.IP.Mask(nn.Gateway.Mask), Mask: nn.Gateway.Mask}
			if ip.IsUnspecifiedSubnet(s.Subnet()) || !s.Subnet().IP.Equal(subnet.IP) || !bytes.Equal(s.Subnet().Mask, subnet.Mask) {
				t.Fatalf("external network %s was loaded with wrong subnet, got: %s, want: %s", n, s.Subnet(), subnet)
			}

			if ip.IsUnspecifiedIP(s.Gateway()) || !s.Gateway().Equal(nn.Gateway.IP) {
				t.Fatalf("external network %s was loaded with wrong gateway, got: %s, want: %s", n, s.Gateway(), nn.Gateway.IP)
			}

			if len(nn.Pools) == 0 {
				// only one pool corresponding to the subnet
				if len(pools) != 1 || pools[0] != subnet.String() {
					t.Fatalf("external network %s was loaded with wrong pool, got: %+v, want %+v", n, pools, []*net.IPNet{subnet})
				}
			}
		}

		for _, d := range nn.Nameservers {
			found := false
			for _, d2 := range s.DNS() {
				if d2.Equal(d) {
					found = true
					break
				}
			}

			if !found {
				t.Fatalf("external network %s was loaded with wrong nameservers, got: %+v, want: %+v", n, s.DNS(), nn.Nameservers)
			}
		}

		for _, p := range nn.Pools {
			found := false
			for _, p2 := range pools {
				if p2 == p.String() {
					found = true
					break
				}
			}

			if !found {
				t.Fatalf("external network %s was loaded with wrong pools, got: %+v, want: %+v", n, s.IPAM().Pools(), nn.Pools)
			}
		}
	}
}

func TestContextNewScope(t *testing.T) {
	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32))
	if err != nil {
		t.Fatalf("NewContext() => (nil, %s), want (ctx, nil)", err)
	}

	var tests = []struct {
		in  params
		out *params
		err error
	}{
		// empty name
		{params{"bridge", "", nil, net.IPv4(0, 0, 0, 0), nil, nil}, nil, fmt.Errorf("")},
		// unsupported network type
		{params{"foo", "bar8", nil, net.IPv4(0, 0, 0, 0), nil, nil}, nil, fmt.Errorf("")},
		// duplicate name
		{params{"bridge", "bar6", nil, net.IPv4(0, 0, 0, 0), nil, nil}, nil, DuplicateResourceError{}},
		// ip range already allocated
		{params{"bridge", "bar9", &net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, nil}, nil, fmt.Errorf("")},
		// ipam out of range of network
		{params{"bridge", "bar10", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, []string{"10.14.1.0/24", "10.15.1.0/24"}}, nil, fmt.Errorf("")},
		// gateway not on subnet
		{params{"bridge", "bar101", &net.IPNet{IP: net.IPv4(10, 141, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(10, 14, 0, 1), nil, nil}, nil, fmt.Errorf("")},
		// gateway is allzeros address
		{params{"bridge", "bar102", &net.IPNet{IP: net.IPv4(10, 142, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(10, 142, 0, 0), nil, nil}, nil, fmt.Errorf("")},
		// gateway is allones address
		{params{"bridge", "bar103", &net.IPNet{IP: net.IPv4(10, 143, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(10, 143, 255, 255), nil, nil}, nil, fmt.Errorf("")},
		// this should succeed now
		{params{"bridge", "bar11", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, []string{"10.14.1.0/24"}},
			&params{"bridge", "bar11", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(10, 14, 1, 0), nil, nil},
			nil},
		// bad ipam
		{params{"bridge", "bar12", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, []string{"10.14.1.0/24", "10.15.1"}}, nil, fmt.Errorf("")},
		// bad ipam, default bridge pool
		{params{"bridge", "bar12", &net.IPNet{IP: net.IPv4(172, 21, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, []string{"172.21.1.0/24", "10.15.1"}}, nil, fmt.Errorf("")},
		// external networks must have subnet specified, if pool is specified
		{params{"external", "bar13", nil, net.IPv4(0, 0, 0, 0), nil, []string{"10.15.0.0/24"}}, nil, fmt.Errorf("")},
		// external networks must have gateway specified, if pool is specified
		{params{"external", "bar14", &net.IPNet{IP: net.IPv4(10, 14, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(0, 0, 0, 0), nil, []string{"10.15.0.0/24"}}, nil, fmt.Errorf("")},
		// external networks cannot overlap bridge pool
		{params{"external", "bar16", &net.IPNet{IP: net.IPv4(172, 20, 0, 0), Mask: net.CIDRMask(16, 32)}, net.IPv4(10, 14, 0, 1), nil, []string{"172.20.0.0/16"}}, nil, fmt.Errorf("")},
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
				t.Fatalf("NewScope(%s, %s, %s, %s, %+v, %+v) => (s, nil), want (nil, err)",
					te.in.scopeType,
					te.in.name,
					te.in.subnet,
					te.in.gateway,
					te.in.dns,
					te.in.ipam)
			}

			// if there is an error specified, check if we got that error
			if te.err != nil &&
				reflect.TypeOf(err) != reflect.TypeOf(te.err) {
				t.Fatalf("NewScope() => (nil, %s), want (nil, %s)", reflect.TypeOf(err), reflect.TypeOf(te.err))
			}

			if _, o := err.(DuplicateResourceError); !o {
				// sanity check
				if _, ok := ctx.scopes[te.in.name]; ok {
					t.Fatalf("scope %s added on error", te.in.name)
				}
			}

			continue
		}

		if err != nil {
			t.Fatalf("got: %s, expected: nil", err)
			continue
		}

		if s.Type() != te.out.scopeType {
			t.Fatalf("s.Type() => %s, want %s", s.Type(), te.out.scopeType)
			continue
		}

		if s.Name() != te.out.name {
			t.Fatalf("s.Name() => %s, want %s", s.Name(), te.out.name)
		}

		if s.Subnet().String() != te.out.subnet.String() {
			t.Fatalf("s.Subnet() => %s, want %s", s.Subnet(), te.out.subnet)
		}

		if !s.Gateway().Equal(te.out.gateway) {
			t.Fatalf("s.Gateway() => %s, want %s", s.Gateway(), te.out.gateway)
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
				t.Fatalf("s.DNS() => %q, want %q", s.DNS(), te.out.dns)
				break
			}
		}

		ipam := s.IPAM()
		if ipam == nil {
			t.Fatalf("s.IPAM() == nil, want %q", te.in.ipam)
			continue
		}

		if s.Type() == BridgeScopeType && s.Network() != testBridgeNetwork {
			t.Fatalf("s.NetworkName => %v, want %s", s.Network(), testBridgeNetwork)
			continue
		}

		if te.in.ipam != nil && len(ipam.spaces) != len(te.in.ipam) {
			t.Fatalf("len(ipam.spaces) => %d != len(te.in.ipam) => %d", len(ipam.spaces), len(te.in.ipam))
		}

		for i, p := range ipam.spaces {
			if te.in.ipam == nil {
				if p != s.space {
					t.Fatalf("got %v, want %v", p, s.space)
				}
				continue
			}

			if p.Parent != s.space {
				t.Fatalf("p.Parent => %v, want %v", p.Parent, s.space)
				continue
			}

			if p.Network != nil {
				if p.Network.String() != te.in.ipam[i] {
					t.Fatalf("p.Network => %s, want %s", p.Network, te.in.ipam[i])
				}
			} else if p.Pool.String() != te.in.ipam[i] {
				t.Fatalf("p.Pool => %s, want %s", p.Pool, te.in.ipam[i])
			}
		}
	}
}

func TestScopes(t *testing.T) {
	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32))
	if err != nil {
		t.Fatalf("NewContext() => (nil, %s), want (ctx, nil)", err)
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
			t.Fatalf("NewScope() => (_, %s), want (_, nil)", err)
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
				t.Fatalf("Scopes() => (_, nil), want (_, err)")
				continue
			}
		} else {
			if err != nil {
				t.Fatalf("Scopes() => (_, %s), want (_, nil)", err)
				continue
			}
		}

		// +5 for the default bridge scope, and 4 external networks
		if te.in == nil {
			if len(l) != len(te.out)+5 {
				t.Fatalf("len(scopes) => %d != %d", len(l), len(te.out)+5)
				continue
			}
		} else {
			if len(l) != len(te.out) {
				t.Fatalf("len(scopes) => %d != %d", len(l), len(te.out))
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
				t.Fatalf("got=%v, want=%v", l, te.out)
				break
			}
		}
	}
}

func TestContextAddContainer(t *testing.T) {
	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32))
	if err != nil {
		t.Fatalf("NewContext() => (nil, %s), want (ctx, nil)", err)
		return
	}

	h := exec.NewContainer("foo")

	var devices object.VirtualDeviceList
	backing, _ := ctx.DefaultScope().Network().EthernetCardBackingInfo(context.TODO())

	specWithEthCard := &spec.VirtualMachineConfigSpec{
		VirtualMachineConfigSpec: &types.VirtualMachineConfigSpec{},
	}

	var d types.BaseVirtualDevice
	if d, err = devices.CreateEthernetCard("vmxnet3", backing); err == nil {
		d.GetVirtualDevice().SlotInfo = &types.VirtualDevicePciBusSlotInfo{
			PciSlotNumber: 1111,
		}
		devices = append(devices, d)
		var cs []types.BaseVirtualDeviceConfigSpec
		if cs, err = devices.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd); err == nil {
			specWithEthCard.DeviceChange = cs
		}
	}

	if err != nil {
		t.Fatalf(err.Error())
	}

	aecErr := func(_ *exec.Handle, _ *Scope) (types.BaseVirtualDevice, error) {
		return nil, fmt.Errorf("error")
	}

	otherScope, err := ctx.NewScope(BridgeScopeType, "other", nil, net.IPv4(0, 0, 0, 0), nil, nil)
	if err != nil {
		t.Fatalf("failed to add scope")
	}

	hBar := exec.NewContainer("bar")

	var tests = []struct {
		aec   func(h *exec.Handle, s *Scope) (types.BaseVirtualDevice, error)
		h     *exec.Handle
		s     *spec.VirtualMachineConfigSpec
		scope string
		ip    *net.IP
		err   error
	}{
		// nil handle
		{nil, nil, nil, "", nil, fmt.Errorf("")},
		// scope not found
		{nil, h, nil, "foo", nil, ResourceNotFoundError{}},
		// addEthernetCard returns error
		{aecErr, h, nil, "default", nil, fmt.Errorf("")},
		// add a container
		{nil, h, nil, "default", nil, nil},
		// container already added
		{nil, h, nil, "default", nil, nil},
		{nil, hBar, specWithEthCard, "default", nil, nil},
		{nil, hBar, nil, otherScope.Name(), nil, nil},
	}

	origAEC := addEthernetCard
	defer func() { addEthernetCard = origAEC }()

	for i, te := range tests {
		// setup
		addEthernetCard = origAEC
		scopy := &spec.VirtualMachineConfigSpec{}
		if te.h != nil {
			te.h.SetSpec(te.s)
			if te.h.Spec != nil {
				*scopy = *te.h.Spec
			}
		}

		if te.aec != nil {
			addEthernetCard = te.aec
		}

		options := &AddContainerOptions{
			Scope: te.scope,
			IP:    te.ip,
		}
		err := ctx.AddContainer(te.h, options)
		if te.err != nil {
			// expect an error
			if err == nil {
				t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) => nil want err", i, te.h, te.scope, te.ip)
			}

			if reflect.TypeOf(err) != reflect.TypeOf(te.err) {
				t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) => (%v, %v) want (%v, %v)", i, te.h, te.scope, te.ip, err, te.err, err, te.err)
			}

			if _, ok := te.err.(DuplicateResourceError); ok {
				continue
			}

			// verify no device changes in the spec
			if te.s != nil {
				if len(scopy.DeviceChange) != len(h.Spec.DeviceChange) {
					t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) added device", i, te.h, te.scope, te.ip)
				}
			}

			continue
		}

		if err != nil {
			t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) => %s want nil", i, te.h, te.scope, te.ip, err)
		}

		// verify the container was not added to the scope
		s, _ := ctx.resolveScope(te.scope)
		if s != nil && te.h != nil {
			c := s.Container(exec.ParseID(te.h.Container.ExecConfig.ID))
			if c != nil {
				t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) added container", i, te.h, te.scope, te.ip)
			}
		}

		// spec should have a nic attached to the scope's network
		var dev types.BaseVirtualDevice
		dcs, err := te.h.Spec.FindNICs(context.TODO(), s.Network())
		if len(dcs) != 1 {
			t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) more than one NIC added for scope %s", i, te.h, te.scope, te.ip, s.Network())
		}
		dev = dcs[0].GetVirtualDeviceConfigSpec().Device
		if spec.VirtualDeviceSlotNumber(dev) == spec.NilSlot {
			t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) NIC added has nil pci slot", i, te.h, te.scope, te.ip)
		}

		// spec metadata should be updated with endpoint info
		ne, ok := te.h.ExecConfig.Networks[s.Name()]
		if !ok {
			t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) no network endpoint info added", i, te.h, te.scope, te.ip)
		}

		if spec.VirtualDeviceSlotNumber(dev) != atoiOrZero(ne.ID) {
			t.Fatalf("case %d; ctx.AddContainer(%v, %s, %s) => ne.ID == %d, want %d", i, te.h, te.scope, te.ip, atoiOrZero(ne.ID), spec.VirtualDeviceSlotNumber(dev))
		}

		if ne.Network.Name != s.Name() {
			t.Fatalf("case %d; ctx.AddContainer(%v, %s, %s) => ne.NetworkName == %s, want %s", i, te.h, te.scope, te.ip, ne.Network.Name, s.Name())
		}

		if te.ip != nil && !te.ip.Equal(ne.Static.IP) {
			t.Fatalf("case %d; ctx.AddContainer(%v, %s, %s) => ne.Static.IP == %s, want %s", i, te.h, te.scope, te.ip, ne.Static.IP, te.ip)
		}

		if te.ip == nil && ne.Static != nil {
			t.Fatalf("case %d; ctx.AddContainer(%v, %s, %s) => ne.Static.IP == %s, want %s", i, te.h, te.scope, te.ip, ne.Static.IP, net.IPv4zero)
		}
	}
}

func TestContextBindUnbindContainer(t *testing.T) {
	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32))
	if err != nil {
		t.Fatalf("NewContext() => (nil, %s), want (ctx, nil)", err)
	}

	scope, err := ctx.NewScope(BridgeScopeType, "scope", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("ctx.NewScope(%s, %s, nil, nil, nil) => (nil, %s)", BridgeScopeType, "scope", err)
	}

	foo := exec.NewContainer("foo")
	added := exec.NewContainer("added")
	staticIP := exec.NewContainer("staticIP")
	ipErr := exec.NewContainer("ipErr")
	alias := exec.NewContainer("alias")
	aliasErr := exec.NewContainer("aliasErr")

	options := &AddContainerOptions{
		Scope: ctx.DefaultScope().Name(),
	}
	// add a container to the default scope
	if err = ctx.AddContainer(added, options); err != nil {
		t.Fatalf("ctx.AddContainer(%s, %s, nil) => %s", added, ctx.DefaultScope().Name(), err)
	}

	// add a container with a static IP
	ip := net.IPv4(172, 16, 0, 10)
	options = &AddContainerOptions{
		Scope: ctx.DefaultScope().Name(),
		IP:    &ip,
	}
	if err = ctx.AddContainer(staticIP, options); err != nil {
		t.Fatalf("ctx.AddContainer(%s, %s, nil) => %s", staticIP, ctx.DefaultScope().Name(), err)
	}

	options = &AddContainerOptions{
		Scope: scope.Name(),
	}
	if err = ctx.AddContainer(added, options); err != nil {
		t.Fatalf("ctx.AddContainer(%s, %s, nil) => %s", added, scope.Name(), err)
	}

	// add a container with an ip that is already taken,
	// causing Scope.BindContainer call to fail
	gw := ctx.DefaultScope().Gateway()
	options = &AddContainerOptions{
		Scope: scope.Name(),
	}
	ctx.AddContainer(ipErr, options)

	options = &AddContainerOptions{
		Scope: ctx.DefaultScope().Name(),
		IP:    &gw,
	}
	ctx.AddContainer(ipErr, options)

	// add a container with correct aliases
	options = &AddContainerOptions{
		Scope:   ctx.DefaultScope().Name(),
		Aliases: []string{"added:foo", ":bar"},
	}
	if err = ctx.AddContainer(alias, options); err != nil {
		t.Fatalf("ctx.AddContainer(%s, %s, nil) => %s", alias, ctx.DefaultScope().Name(), err)
	}

	// add a container with incorrect aliases
	options = &AddContainerOptions{
		Scope:   ctx.DefaultScope().Name(),
		Aliases: []string{"cloud:foo", "bar"},
	}
	if err = ctx.AddContainer(aliasErr, options); err != nil {
		t.Fatalf("ctx.AddContainer(%s, %s, nil) => %s", aliasErr, ctx.DefaultScope().Name(), err)
	}

	var tests = []struct {
		i      int
		h      *exec.Handle
		scopes []string
		ips    []net.IP
		static bool
		err    error
	}{
		// container not added to scope
		{0, foo, []string{}, []net.IP{}, false, fmt.Errorf("")},
		// container has bad ip address
		{1, ipErr, []string{}, nil, false, fmt.Errorf("")},
		// successful container bind
		{2, added, []string{ctx.DefaultScope().Name(), scope.Name()}, []net.IP{net.IPv4(172, 16, 0, 2), net.IPv4(172, 17, 0, 2)}, false, nil},
		{3, staticIP, []string{ctx.DefaultScope().Name()}, []net.IP{net.IPv4(172, 16, 0, 10)}, true, nil},
		{4, alias, []string{ctx.DefaultScope().Name()}, []net.IP{net.IPv4(172, 16, 0, 3)}, false, nil},
		{5, aliasErr, []string{ctx.DefaultScope().Name()}, []net.IP{}, false, fmt.Errorf("")},
	}

	for _, te := range tests {
		eps, err := ctx.BindContainer(te.h)
		if te.err != nil {
			// expect an error
			if err == nil || eps != nil {
				t.Fatalf("%d: ctx.BindContainer(%s) => (%#v, %#v), want (%#v, %#v)", te.i, te.h, eps, err, nil, te.err)
			}

			con := ctx.Container(exec.ParseID(te.h.Container.ExecConfig.ID))
			if con != nil {
				t.Fatalf("%d: ctx.BindContainer(%s) added container %#v", te.i, te.h, con)
			}

			continue
		}

		// check if the correct endpoints were added
		con := ctx.Container(exec.ParseID(te.h.Container.ExecConfig.ID))
		if con == nil {
			t.Fatalf("%d: ctx.Container(%s) => nil, want %s", te.i, te.h.Container.ExecConfig.ID, te.h.Container.ExecConfig.ID)
		}

		if len(con.Scopes()) != len(te.scopes) {
			t.Fatalf("%d: len(con.Scopes()) %#v != len(te.scopes) %#v", te.i, con.Scopes(), te.scopes)
		}

		// check endpoints
		for i, s := range te.scopes {
			found := false
			for _, e := range eps {
				if e.Scope().Name() != s {
					continue
				}

				found = true
				if !e.Gateway().Equal(e.Scope().Gateway()) {
					t.Fatalf("%d: ctx.BindContainer(%s) => endpoint gateway %s, want %s", te.i, te.h, e.Gateway(), e.Scope().Gateway())
				}
				if !e.IP().Equal(te.ips[i]) {
					t.Fatalf("%d: ctx.BindContainer(%s) => endpoint IP %s, want %s", te.i, te.h, e.IP(), te.ips[i])
				}
				if e.Subnet().String() != e.Scope().Subnet().String() {
					t.Fatalf("%d: ctx.BindContainer(%s) => endpoint subnet %s, want %s", te.i, te.h, e.Subnet(), e.Scope().Subnet())
				}

				ne := te.h.ExecConfig.Networks[s]
				if !ne.Static.IP.Equal(te.ips[i]) {
					t.Fatalf("%d: ctx.BindContainer(%s) => metadata endpoint IP %s, want %s", te.i, te.h, ne.Static.IP, te.ips[i])
				}
				if ne.Static.Mask.String() != e.Scope().Subnet().Mask.String() {
					t.Fatalf("%d: ctx.BindContainer(%s) => metadata endpoint IP mask %s, want %s", te.i, te.h, ne.Static.Mask.String(), e.Scope().Subnet().Mask.String())
				}
				if !ne.Network.Gateway.IP.Equal(e.Scope().Gateway()) {
					t.Fatalf("%d: ctx.BindContainer(%s) => metadata endpoint gateway %s, want %s", te.i, te.h, ne.Network.Gateway.IP, e.Scope().Gateway())
				}
				if ne.Network.Gateway.Mask.String() != e.Scope().Subnet().Mask.String() {
					t.Fatalf("%d: ctx.BindContainer(%s) => metadata endpoint gateway mask %s, want %s", te.i, te.h, ne.Network.Gateway.Mask.String(), e.Scope().Subnet().Mask.String())
				}

				break
			}

			if !found {
				t.Fatalf("%d: ctx.BindContainer(%s) => endpoint for scope %s not added", te.i, te.h, s)
			}
		}
	}

	tests = []struct {
		i      int
		h      *exec.Handle
		scopes []string
		ips    []net.IP
		static bool
		err    error
	}{
		// container not found
		{0, foo, []string{}, nil, false, fmt.Errorf("")},
		// container has bad ip address
		{1, ipErr, []string{ctx.DefaultScope().Name(), scope.Name()}, nil, false, fmt.Errorf("")},
		// successful container unbind
		{2, added, []string{ctx.DefaultScope().Name(), scope.Name()}, nil, false, nil},
		{3, staticIP, []string{ctx.DefaultScope().Name()}, nil, true, nil},
		{4, alias, []string{ctx.DefaultScope().Name()}, nil, false, nil},
		{5, aliasErr, []string{ctx.DefaultScope().Name()}, nil, false, fmt.Errorf("")},
	}

	// test UnbindContainer
	for _, te := range tests {
		eps, err := ctx.UnbindContainer(te.h)
		if te.err != nil {
			if err == nil {
				t.Fatalf("%d: ctx.UnbindContainer(%s) => nil, want err", te.i, te.h)
			}

			continue
		}

		// container should not be there
		con := ctx.Container(exec.ParseID(te.h.Container.ExecConfig.ID))
		if con != nil {
			t.Fatalf("%d: ctx.Container(%s) => %#v, want nil", te.i, te.h, con)
		}

		for _, s := range te.scopes {
			found := false
			for _, e := range eps {
				if e.Scope().Name() == s {
					found = true
				}
			}

			if !found {
				t.Fatalf("%d: ctx.UnbindContainer(%s) did not return endpoint for scope %s. Endpoints: %+v", te.i, te.h, s, eps)
			}

			// container should not be part of scope
			scopes, err := ctx.Scopes(&s)
			if err != nil || len(scopes) != 1 {
				t.Fatalf("%d: ctx.Scopes(%s) => (%#v, %#v)", te.i, s, scopes, err)
			}
			if scopes[0].Container(exec.ParseID(te.h.Container.ExecConfig.ID)) != nil {
				t.Fatalf("%d: container %s is still part of scope %s", te.i, te.h.Container.ExecConfig.ID, s)
			}

			// check if endpoint is still there, but without the ip
			ne, ok := te.h.ExecConfig.Networks[s]
			if !ok {
				t.Fatalf("%d: container endpoint not present in %v", te.i, te.h.ExecConfig)
			}

			if !te.static && ne.Static != nil {
				t.Fatalf("%d: endpoint IP should be nil in %v", te.i, ne)
			}

			if te.static && (ne.Static == nil || ne.Static.IP.Equal(net.IPv4zero)) {
				t.Fatalf("%d: endpoint IP should not be zero in %v", te.i, ne)
			}
		}
	}
}

func TestContextRemoveContainer(t *testing.T) {

	hFoo := exec.NewContainer("foo")

	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32))
	if err != nil {
		t.Fatalf("NewContext() => (nil, %s), want (ctx, nil)", err)
	}

	scope, err := ctx.NewScope(BridgeScopeType, "scope", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("ctx.NewScope() => (nil, %s), want (scope, nil)", err)
	}

	options := &AddContainerOptions{
		Scope: scope.Name(),
	}
	ctx.AddContainer(hFoo, options)
	ctx.BindContainer(hFoo)

	// container that is added to multiple bridge scopes
	hBar := exec.NewContainer("bar")
	options.Scope = "default"
	ctx.AddContainer(hBar, options)
	options.Scope = scope.Name()
	ctx.AddContainer(hBar, options)

	var tests = []struct {
		h     *exec.Handle
		scope string
		err   error
	}{
		{nil, "", fmt.Errorf("")},                             // nil handle
		{hBar, "bar", fmt.Errorf("")},                         // scope not found
		{hFoo, scope.Name(), fmt.Errorf("")},                  // bound container
		{exec.NewContainer("baz"), "default", fmt.Errorf("")}, // container not part of scope
		{hBar, "default", nil},
		{hBar, scope.Name(), nil},
	}

	for i, te := range tests {
		var ne *executor.NetworkEndpoint
		if te.h != nil && te.h.ExecConfig.Networks != nil {
			ne = te.h.ExecConfig.Networks[te.scope]
		}

		err = ctx.RemoveContainer(te.h, te.scope)
		if te.err != nil {
			// expect error
			if err == nil {
				t.Fatalf("%d: ctx.RemoveContainer(%#v, %s) => nil want err", i, te.h, te.scope)
			}

			continue
		}

		s, err := ctx.resolveScope(te.scope)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if s.Container(exec.ParseID(te.h.Container.ExecConfig.ID)) != nil {
			t.Fatalf("container %s is part of scope %s", te.h, s.Name())
		}

		// should have a remove spec for NIC, if container was only part of one bridge scope
		dcs, err := te.h.Spec.FindNICs(context.TODO(), s.Network())
		if err != nil {
			t.Fatalf(err.Error())
		}

		found := false
		var d types.BaseVirtualDevice
		for _, dc := range dcs {
			if dc.GetVirtualDeviceConfigSpec().Operation != types.VirtualDeviceConfigSpecOperationRemove {
				continue
			}

			d = dc.GetVirtualDeviceConfigSpec().Device
			found = true
			break
		}

		// if a remove spec for the NIC was found, check if any other
		// network endpoints are still using it
		if found {
			for _, ne := range te.h.ExecConfig.Networks {
				if atoiOrZero(ne.ID) == spec.VirtualDeviceSlotNumber(d) {
					t.Fatalf("%d: NIC with pci slot %d is still in use by a network endpoint %#v", i, spec.VirtualDeviceSlotNumber(d), ne)
				}
			}
		} else if ne != nil {
			// check if remove spec for NIC should have been there
			for _, ne2 := range te.h.ExecConfig.Networks {
				if ne.ID == ne2.ID {
					t.Fatalf("%d: NIC with pci slot %s should have been removed", i, ne.ID)
				}
			}
		}

		// metadata should be gone
		if _, ok := te.h.ExecConfig.Networks[te.scope]; ok {
			t.Fatalf("%d: endpoint metadata for container still present in handle %#v", i, te.h.ExecConfig)
		}
	}
}

func TestDeleteScope(t *testing.T) {
	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32))
	if err != nil {
		t.Fatalf("NewContext() => (nil, %s), want (ctx, nil)", err)
	}

	foo, err := ctx.NewScope(BridgeScopeType, "foo", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("ctx.NewScope(%s, \"foo\", nil, nil, nil, nil) => (nil, %#v), want (foo, nil)", BridgeScopeType, err)
	}
	h := exec.NewContainer("container")
	options := &AddContainerOptions{
		Scope: foo.Name(),
	}
	ctx.AddContainer(h, options)

	// bar is a scope with bound endpoints
	bar, err := ctx.NewScope(BridgeScopeType, "bar", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("ctx.NewScope(%s, \"bar\", nil, nil, nil, nil) => (nil, %#v), want (bar, nil)", BridgeScopeType, err)
	}

	h = exec.NewContainer("container2")
	options.Scope = bar.Name()
	ctx.AddContainer(h, options)
	ctx.BindContainer(h)

	baz, err := ctx.NewScope(BridgeScopeType, "bazScope", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("ctx.NewScope(%s, \"bazScope\", nil, nil, nil, nil) => (nil, %#v), want (baz, nil)", BridgeScopeType, err)
	}

	qux, err := ctx.NewScope(BridgeScopeType, "quxScope", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("ctx.NewScope(%s, \"quxScope\", nil, nil, nil, nil) => (nil, %#v), want (qux, nil)", BridgeScopeType, err)
	}

	var tests = []struct {
		name string
		err  error
	}{
		{"", ResourceNotFoundError{}},
		{ctx.DefaultScope().Name(), fmt.Errorf("cannot delete builtin scopes")},
		{bar.Name(), fmt.Errorf("cannot delete scope with bound endpoints")},
		// full name
		{foo.Name(), nil},
		// full id
		{baz.ID(), nil},
		// partial id
		{qux.ID()[:6], nil},
	}

	for _, te := range tests {
		err := ctx.DeleteScope(te.name)
		if te.err != nil {
			if err == nil {
				t.Fatalf("DeleteScope(%s) => nil, expected err", te.name)
			}

			if reflect.TypeOf(te.err) != reflect.TypeOf(err) {
				t.Fatalf("DeleteScope(%s) => %#v, want %#v", te.name, err, te.err)
			}

			continue
		}

		scopes, err := ctx.Scopes(&te.name)
		if _, ok := err.(ResourceNotFoundError); !ok || len(scopes) != 0 {
			t.Fatalf("scope %s not deleted", te.name)
		}
	}
}
