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
	"os"
	"reflect"
	"testing"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/vsphere/session"
)

const (
	testBridgeName = "testBridge"
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

// mockBridgeNetworkName mocks getBridgeNetworkName so that tests don't
// need to query guestInfo
func mockBridgeNetworkName(sess *session.Session) (string, error) {
	return testBridgeName, nil
}

func TestMain(m *testing.M) {
	origBridgeNetworkName := getBridgeNetworkName
	getBridgeNetworkName = mockBridgeNetworkName

	rc := m.Run()

	getBridgeNetworkName = origBridgeNetworkName

	os.Exit(rc)
}

func TestContext(t *testing.T) {
	sess := &session.Session{}

	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32), sess)
	if err != nil {
		t.Errorf("NewContext() => (nil, %s), want (ctx, nil)", err)
		return
	}

	if ctx.BridgeNetworkName != testBridgeName {
		t.Errorf("ctx.BridgeNetworkName => %v, want %s", ctx.BridgeNetworkName, testBridgeName)
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

		if s.Type() == "bridge" && s.NetworkName != testBridgeName {
			t.Errorf("s.NetworkName => %v, want %s", s.NetworkName, testBridgeName)
			continue
		}

		if s.Type() == "external" && s.NetworkName != "" {
			t.Errorf("s.NetworkName => %v, want %s", s.NetworkName, "")
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
	origBridgeNetworkName := getBridgeNetworkName
	getBridgeNetworkName = mockBridgeNetworkName
	defer func() { getBridgeNetworkName = origBridgeNetworkName }()

	sess := &session.Session{}

	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32), sess)
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

func TestContextAddContainer(t *testing.T) {
	sess := &session.Session{}

	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32), sess)
	if err != nil {
		t.Errorf("NewContext() => (nil, %s), want (ctx, nil)", err)
		return
	}

	h := exec.NewContainer("foo")

	var devices object.VirtualDeviceList
	backing := &types.VirtualEthernetCardNetworkBackingInfo{
		VirtualDeviceDeviceBackingInfo: types.VirtualDeviceDeviceBackingInfo{
			DeviceName: ctx.DefaultScope().NetworkName,
		},
	}

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

	aecErr := func(_ *exec.Handle, _ *Scope, _ *Container) (types.BaseVirtualDevice, error) {
		return nil, fmt.Errorf("error")
	}

	otherScope, err := ctx.NewScope(bridgeScopeType, "other", nil, net.IPv4(0, 0, 0, 0), nil, nil)
	if err != nil {
		t.Fatalf("failed to add scope")
	}

	var tests = []struct {
		aec      func(h *exec.Handle, s *Scope, con *Container) (types.BaseVirtualDevice, error)
		h        *exec.Handle
		s        *spec.VirtualMachineConfigSpec
		scope    string
		ip       *net.IP
		ethAdded bool
		e        *Endpoint
		err      error
	}{
		// nil handle
		{nil, nil, nil, "", nil, false, nil, fmt.Errorf("")},
		// scope not found
		{nil, h, nil, "foo", nil, false, nil, ResourceNotFoundError{}},
		// addEthernetCard returns error
		{aecErr, h, nil, "default", nil, false, nil, fmt.Errorf("")},
		// add a container
		{nil, h, nil, "default", nil, true, &Endpoint{ip: net.IPv4(0, 0, 0, 0), scope: ctx.DefaultScope(), gateway: ctx.DefaultScope().gateway, subnet: ctx.DefaultScope().subnet, static: false}, nil},
		// container already added
		{nil, h, nil, "default", nil, false, nil, DuplicateResourceError{}},
		{nil, exec.NewContainer("bar"), specWithEthCard, "default", nil, true, &Endpoint{ip: net.IPv4(0, 0, 0, 0), scope: ctx.DefaultScope(), gateway: ctx.DefaultScope().Gateway(), subnet: *ctx.DefaultScope().Subnet(), static: false}, nil},
		{nil, exec.GetContainer(exec.ParseID("bar")), nil, otherScope.Name(), nil, false, &Endpoint{ip: net.IPv4(0, 0, 0, 0), scope: otherScope, gateway: otherScope.Gateway(), subnet: *otherScope.Subnet(), static: false}, nil},
	}

	origAEC := addEthernetCard
	defer func() { addEthernetCard = origAEC }()

	for i, te := range tests {
		// setup
		addEthernetCard = origAEC
		if te.h != nil {
			te.h.SetSpec(te.s)
		}
		scopy := &spec.VirtualMachineConfigSpec{}
		if te.s != nil {
			*scopy = *te.s
		}
		if te.aec != nil {
			addEthernetCard = te.aec
		}

		e, err := ctx.AddContainer(te.h, te.scope, te.ip)
		if te.err != nil {
			// expect an error
			if err == nil || te.e != e {
				t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) => (%v, %s) want (%v, err)", i, te.h, te.scope, te.ip, e, err, te.e)
			}

			if reflect.TypeOf(err) != reflect.TypeOf(te.err) {
				t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) => (%v, %v) want (%v, %v)", i, te.h, te.scope, te.ip, err, te.err, err, te.err)
			}

			if _, ok := te.err.(DuplicateResourceError); ok {
				continue
			}

			// verify the container was not added to the scope
			s, _ := ctx.resolveScope(te.scope)
			if s != nil && te.h != nil {
				c := s.Container(te.h.Container.ID)
				if c != nil {
					t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) added container", i, te.h, te.scope, te.ip)
				}
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
			t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) => (%v, %s) want (%v, nil)", i, te.h, te.scope, te.ip, e, err, te.e)
		}

		if te.e.scope != e.scope ||
			!te.e.gateway.Equal(e.gateway) ||
			te.e.subnet.String() != e.subnet.String() ||
			te.e.static != e.static {
			t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) => (%v, %s) want (%v, nil)", i, te.h, te.scope, te.ip, e, err, te.e)
		}

		// spec should have a nic attached to the scope's network
		found := false
		for _, d := range te.h.Spec.DeviceChange {
			if d.GetVirtualDeviceConfigSpec().Operation != types.VirtualDeviceConfigSpecOperationAdd {
				continue
			}

			dev := d.GetVirtualDeviceConfigSpec().Device
			if backing, ok := dev.GetVirtualDevice().Backing.(*types.VirtualEthernetCardNetworkBackingInfo); ok {
				if backing.DeviceName == e.Scope().NetworkName {
					if e.pciSlot != spec.VirtualDeviceSlotNumber(dev) {
						t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) => pciSlot == %d, want pciSlot == %d", i, te.h, te.scope, te.ip, e.pciSlot, spec.VirtualDeviceSlotNumber(dev))
					}
					found = true
					break
				}
			}
		}

		if found != te.ethAdded {
			t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) ethAdded == %v, want %v", i, te.h, te.scope, te.ip, found, te.ethAdded)
		}

		// spec metadata should be updated with endpoint info
		ne, ok := te.h.ExecConfig.Networks[e.Scope().Name()]
		if !ok {
			t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) no network endpoint info added", i, te.h, te.scope, te.ip)
		}

		ip := net.IPNet{IP: e.IP(), Mask: e.Subnet().Mask}
		gw := net.IPNet{IP: e.Gateway(), Mask: e.Subnet().Mask}
		if ip.String() != ne.IP.String() ||
			e.pciSlot != ne.PCISlot ||
			e.Scope().Name() != ne.Network.Name ||
			gw.String() != ne.Network.Gateway.String() {
			t.Fatalf("case %d: ctx.AddContainer(%v, %s, %s) metadata endpoint = %v, want %v", i, te.h, te.scope, te.ip, ne, e.metadataEndpoint())
		}
	}
}

func TestFindSlotNumber(t *testing.T) {
	allSlots := make(map[int32]bool)
	for s := pciSlotNumberBegin; s != pciSlotNumberEnd; s += pciSlotNumberInc {
		allSlots[s] = true
	}

	// missing first slot
	missingFirstSlot := make(map[int32]bool)
	for s := pciSlotNumberBegin + pciSlotNumberInc; s != pciSlotNumberEnd; s += pciSlotNumberInc {
		missingFirstSlot[s] = true
	}

	// missing last slot
	missingLastSlot := make(map[int32]bool)
	for s := pciSlotNumberBegin; s != pciSlotNumberEnd-pciSlotNumberInc; s += pciSlotNumberInc {
		missingLastSlot[s] = true
	}

	// missing a slot in the middle
	var missingSlot int32
	missingMiddleSlot := make(map[int32]bool)
	for s := pciSlotNumberBegin; s != pciSlotNumberEnd-pciSlotNumberInc; s += pciSlotNumberInc {
		if pciSlotNumberBegin+(2*pciSlotNumberInc) == s {
			missingSlot = s
			continue
		}
		missingMiddleSlot[s] = true
	}

	var tests = []struct {
		slots map[int32]bool
		out   int32
	}{
		{make(map[int32]bool), pciSlotNumberBegin},
		{allSlots, spec.NilSlot},
		{missingFirstSlot, pciSlotNumberBegin},
		{missingLastSlot, pciSlotNumberEnd - pciSlotNumberInc},
		{missingMiddleSlot, missingSlot},
	}

	for _, te := range tests {
		if s := findSlotNumber(te.slots); s != te.out {
			t.Fatalf("findSlotNumber(%v) => %d, want %d", te.slots, s, te.out)
		}
	}
}

func TestContextBindUnbindContainer(t *testing.T) {
	sess := &session.Session{}

	ctx, err := NewContext(net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)}, net.CIDRMask(16, 32), sess)
	if err != nil {
		t.Fatalf("NewContext() => (nil, %s), want (ctx, nil)", err)
	}

	scope, err := ctx.NewScope(bridgeScopeType, "scope", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("ctx.NewScope(%s, %s, nil, nil, nil) => (nil, %s)", bridgeScopeType, "scope", err)
	}

	foo := exec.NewContainer("foo")
	added := exec.NewContainer("added")
	ipErr := exec.NewContainer("ipErr")

	// add a container to the default scope
	if _, err = ctx.AddContainer(added, ctx.DefaultScope().Name(), nil); err != nil {
		t.Fatalf("ctx.AddContainer(%s, %s, nil) => %s", added, ctx.DefaultScope().Name(), err)
	}

	if _, err = ctx.AddContainer(added, scope.Name(), nil); err != nil {
		t.Fatalf("ctx.AddContainer(%s, %s, nil) => %s", added, scope.Name(), err)
	}

	// add a container with an ip that is already taken,
	// causing Scope.BindContainer call to fail
	gw := ctx.DefaultScope().Gateway()
	ctx.AddContainer(ipErr, scope.Name(), nil)
	ctx.AddContainer(ipErr, ctx.DefaultScope().Name(), &gw)

	var tests = []struct {
		i                int
		h                *exec.Handle
		containerPresent bool
		scopes           []string
		err              error
	}{
		// container not found
		{0, foo, false, []string{}, fmt.Errorf("")},
		// container has bad ip address
		{1, ipErr, true, []string{}, fmt.Errorf("")},
		// successful container bind
		{2, added, true, []string{ctx.DefaultScope().Name(), scope.Name()}, nil},
	}

	for _, te := range tests {
		err = ctx.BindContainer(te.h)
		if te.err != nil {
			// expect an error
			if err == nil {
				t.Fatalf("%d: ctx.BindContainer(%s) => nil, want err", te.i, te.h)
			}

			if te.containerPresent {
				con := ctx.Container(te.h.Container.ID)
				if con == nil {
					t.Fatalf("%d: ctx.Container(%s) => nil, want %s", te.i, te.h.Container.ID, te.h.Container.ID)
				}

				// check if the con is unbound
				for _, e := range con.Endpoints() {
					if e.IsBound() {
						t.Fatalf("%d: con %s endpoint is bound, want unbound", te.i, te.h.Container.ID)
					}
				}
			}

			continue
		}

		// check if the correct endpoints were added
		con := ctx.Container(te.h.Container.ID)
		if con == nil {
			t.Fatalf("%d: ctx.Container(%s) => nil, want %s", te.i, te.h.Container.ID, te.h.Container.ID)
		}

		if len(con.Scopes()) != len(te.scopes) {
			t.Fatalf("%d: len(con.Scopes()) %v != len(te.scopes) %v", te.i, con.Scopes(), te.scopes)
		}

		scopeNames := make([]string, len(con.Scopes()))
		i := 0
		for _, s := range con.Scopes() {
			scopeNames[i] = s.Name()
			i++
		}

		if !reflect.DeepEqual(scopeNames, te.scopes) {
			t.Fatalf("%d: %v, want %v", te.i, scopeNames, te.scopes)
		}

		// check all endpoints are bound
		for _, e := range con.Endpoints() {
			if !e.IsBound() {
				t.Fatalf("%d: con %s endpoint is unbound, want bound", te.i, te.h.Container.ID)
			}
		}

		for _, s := range con.Scopes() {
			ne, ok := te.h.ExecConfig.Networks[s.Name()]
			if !ok {
				t.Fatalf("%d: endpoint for scope %s not present in %v", te.i, s.Name(), te.h.ExecConfig)
			}

			// check if there is an IP in the endpoint
			if ne.IP.IP.IsUnspecified() {
				t.Fatalf("%d: ne.IP.IP is unspecified in %v", te.i, ne)
			}
		}

	}

	tests = []struct {
		i                int
		h                *exec.Handle
		containerPresent bool
		scopes           []string
		err              error
	}{
		// container not found
		{0, foo, false, []string{}, fmt.Errorf("")},
		// container has bad ip address
		{1, ipErr, true, []string{ctx.DefaultScope().Name(), scope.Name()}, nil},
		// successful container bind
		{2, added, true, []string{ctx.DefaultScope().Name(), scope.Name()}, nil},
	}

	// test UnbindContainer
	for _, te := range tests {
		err = ctx.UnbindContainer(te.h)
		if te.err != nil {
			if err == nil {
				t.Fatalf("%d: ctx.UnbindContainer(%s) => nil, want err", te.i, te.h)
			}

			continue
		}

		con := ctx.Container(te.h.Container.ID)
		if con == nil {
			t.Fatalf("%d: ctx.Container(%s) => nil, want %s", te.i, te.h.Container.ID, te.h.Container.ID)
		}

		for _, s := range te.scopes {
			// container should still be part of scopes
			scopes, err := ctx.Scopes(&s)
			if err != nil || len(scopes) != 1 {
				t.Fatalf("%d: could not find scope %s", te.i, s)
			}

			sc := scopes[0]
			if c := sc.Container(te.h.Container.ID); c == nil {
				t.Fatalf("%d: container %s not part of scope %s", te.i, te.h.Container.ID, s)
			}

			e := con.Endpoint(sc)
			if e.IsBound() {
				t.Fatalf("%d: container %s is still bound to scope %s", te.i, con.ID(), s)
			}

			// check if endpoint is still there, but without the ip
			ne, ok := te.h.ExecConfig.Networks[s]
			if !ok {
				t.Fatalf("%d: container endpoint not present in %v", te.i, te.h.ExecConfig)
			}

			if !ne.IP.IP.IsUnspecified() {
				t.Fatalf("%d: endpoint IP should be unspecified in %v", te.i, ne)
			}
		}
	}
}
