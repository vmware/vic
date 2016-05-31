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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vmw-guestinfo/rpcvmx"
)

const (
	bridgeNetworkKey         = "guestinfo.vch/networks/bridge"
	pciSlotNumberBegin int32 = 0xc0
	pciSlotNumberEnd   int32 = 1 << 10
	pciSlotNumberInc   int32 = 1 << 5
)

type Context struct {
	sync.Mutex

	defaultBridgePool *AddressSpace
	defaultBridgeMask net.IPMask

	scopes       map[string]*Scope
	containers   map[exec.ID]*Container
	defaultScope *Scope

	BridgeNetworkName string // Portgroup name of the bridge network
}

var getBridgeNetworkName = func() (string, error) {
	config := rpcvmx.NewConfig()
	bnn, err := config.String(fmt.Sprintf("%s/portgroup", bridgeNetworkKey), "")
	if err != nil {
		return "", err
	}

	if bnn == "" {
		return bnn, fmt.Errorf("bridge network name not set")
	}

	return bnn, nil
}

func NewContext(bridgePool net.IPNet, bridgeMask net.IPMask) (*Context, error) {
	pones, pbits := bridgePool.Mask.Size()
	mones, mbits := bridgeMask.Size()
	if pbits != mbits || mones < pones {
		return nil, fmt.Errorf("bridge mask is not compatiable with bridge pool mask")
	}

	bnn, err := getBridgeNetworkName()
	if err != nil {
		return nil, err
	}

	ctx := &Context{
		defaultBridgeMask: bridgeMask,
		defaultBridgePool: NewAddressSpaceFromNetwork(&bridgePool),
		scopes:            make(map[string]*Scope),
		containers:        make(map[exec.ID]*Container),
		BridgeNetworkName: bnn,
	}

	s, err := ctx.NewScope("bridge", bridgeScopeType, nil, net.IPv4(0, 0, 0, 0), nil, nil)
	if err != nil {
		return nil, err
	}

	s.builtin = true
	ctx.defaultScope = s
	return ctx, nil
}

func reserveBroadcastAndNetwork(space *AddressSpace) error {
	if space.Network == nil {
		return nil
	}

	if err := space.ReserveIP4(space.Network.IP); err != nil {
		return err
	}

	if err := space.ReserveIP4(highestIP4(space.Network)); err != nil {
		return err
	}

	return nil
}

func isUnspecifiedSubnet(n *net.IPNet) bool {
	if n == nil {
		return true
	}

	ones, bits := n.Mask.Size()
	return bits == 0 || ones == 0
}

func (c *Context) newScopeCommon(id, name, scopeType string, subnet *net.IPNet, gateway net.IP, dns []net.IP, ipam *IPAM, networkName string) (*Scope, error) {
	if isUnspecifiedSubnet(subnet) {
		subnet = defaultSubnet
	}

	var err error

	// allocate the subnet
	space, defaultPool, err := c.reserveSubnet(subnet)
	defer func() {
		if err != nil && space != nil && defaultPool {
			c.defaultBridgePool.ReleaseIP4Range(space)
		}
	}()

	if err != nil {
		return nil, err
	}

	subnet = space.Network

	// reserve the network and broadcast addresses
	err = reserveBroadcastAndNetwork(space)
	defer func() {
		if err == nil || space.Network == nil {
			return
		}

		lo := incrementIP4(space.Network.IP)
		hi := decrementIP4(highestIP4(space.Network))
		space.ReleaseIP4(lo)
		space.ReleaseIP4(hi)
	}()

	if err != nil {
		return nil, err
	}

	subSpaces, err := reservePools(space, ipam)
	if err != nil {
		return nil, err
	}

	ipam.spaces = subSpaces

	if gateway.IsUnspecified() {
		gateway, err = ipam.spaces[0].ReserveNextIP4()
		defer func() {
			if err != nil && !gateway.IsUnspecified() {
				ipam.spaces[0].ReleaseIP4(gateway)
			}
		}()

		if err != nil {
			return nil, err
		}
	}

	newScope := &Scope{
		id:          id,
		name:        name,
		subnet:      *subnet,
		gateway:     gateway,
		ipam:        ipam,
		containers:  make(map[exec.ID]*Container),
		scopeType:   scopeType,
		space:       space,
		dns:         dns,
		builtin:     false,
		NetworkName: networkName,
	}

	c.scopes[name] = newScope

	return newScope, nil
}

func (c *Context) newBridgeScope(id, name string, subnet *net.IPNet, gateway net.IP, dns []net.IP, ipam *IPAM) (newScope *Scope, err error) {
	s, err := c.newScopeCommon(id, name, bridgeScopeType, subnet, gateway, dns, ipam, c.BridgeNetworkName)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Context) newExternalScope(id, name string, subnet *net.IPNet, gateway net.IP, dns []net.IP, ipam *IPAM) (*Scope, error) {
	// have to specify IPAM
	if ipam == nil || len(ipam.pools) == 0 {
		return nil, fmt.Errorf("no ipam spec for external network")
	}

	if isUnspecifiedSubnet(subnet) || gateway.IsUnspecified() {
		return nil, fmt.Errorf("neither subnet nor gateway specified for external network")
	}

	// cannot overlap with the default bridge pool
	if c.defaultBridgePool.Network.Contains(subnet.IP) ||
		c.defaultBridgePool.Network.Contains(highestIP4(subnet)) {
		return nil, fmt.Errorf("external network cannot overlap with default bridge network")
	}

	// TODO Get the correct networkName
	return c.newScopeCommon(id, name, externalScopeType, subnet, gateway, dns, ipam, "")
}

func isDefaultSubnet(subnet *net.IPNet) bool {
	return subnet.IP == nil || subnet.IP.Equal(net.ParseIP("0.0.0.0"))
}

func (c *Context) reserveSubnet(subnet *net.IPNet) (space *AddressSpace, defaultPool bool, err error) {
	defaultPool = true
	if isDefaultSubnet(subnet) {
		space, err = c.defaultBridgePool.ReserveNextIP4Net(subnet.Mask)
		return
	}

	err = c.checkNetOverlap(subnet)
	if err != nil {
		return
	}

	// reserve from the default pool first
	space, err = c.defaultBridgePool.ReserveIP4Net(subnet)
	if err == nil {
		return
	}
	err = nil

	defaultPool = false
	space = NewAddressSpaceFromNetwork(subnet)
	return
}

func (c *Context) checkNetOverlap(subnet *net.IPNet) error {
	// check if the requested subnet is available
	highestIP := highestIP4(subnet)
	for _, scope := range c.scopes {
		if scope.subnet.Contains(subnet.IP) || scope.subnet.Contains(highestIP) {
			return fmt.Errorf("could not allocate subnet for scope")
		}
	}

	return nil
}

func reservePools(space *AddressSpace, ipam *IPAM) ([]*AddressSpace, error) {
	if ipam.pools == nil || len(ipam.pools) == 0 {
		// pool not specified so use the default
		ipam.pools = []string{space.Network.String()}
		return []*AddressSpace{space}, nil
	}

	var err error
	subSpaces := make([]*AddressSpace, len(ipam.pools))
	defer func() {
		if err == nil {
			return
		}

		for _, s := range subSpaces {
			if s == nil {
				continue
			}
			space.ReleaseIP4Range(s)

		}
	}()

	for i, p := range ipam.pools {
		var nw *net.IPNet
		_, nw, err = net.ParseCIDR(p)
		if err == nil {
			subSpaces[i], err = space.ReserveIP4Net(nw)
			if err != nil {
				break
			}

			continue
		}

		// ip range
		r := ParseIPRange(p)
		if r == nil {
			err = fmt.Errorf("error in pool spec")
			break
		}

		var ss *AddressSpace
		ss, err = space.ReserveIP4Range(r.FirstIP, r.LastIP)
		if err != nil {
			break
		}

		subSpaces[i] = ss
	}

	if err != nil {
		return nil, err
	}

	return subSpaces, nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (c *Context) NewScope(scopeType, name string, subnet *net.IPNet, gateway net.IP, dns []net.IP, pools []string) (*Scope, error) {
	// sanity checks
	if name == "" {
		return nil, fmt.Errorf("scope name must not be empty")
	}

	if gateway == nil {
		gateway = net.IPv4(0, 0, 0, 0)
	}

	c.Lock()
	defer c.Unlock()

	if _, ok := c.scopes[name]; ok {
		return nil, DuplicateResourceError{resID: name}
	}

	switch scopeType {
	case bridgeScopeType:
		return c.newBridgeScope(generateID(), name, subnet, gateway, dns, &IPAM{pools: pools})

	case externalScopeType:
		return c.newExternalScope(generateID(), name, subnet, gateway, dns, &IPAM{pools: pools})

	default:
		return nil, fmt.Errorf("scope type not supported")
	}
}

func (c *Context) findScope(idName *string) ([]*Scope, error) {
	if idName != nil && *idName != "" {
		// search by name
		scope, ok := c.scopes[*idName]
		if ok {
			return []*Scope{scope}, nil
		}

		// search by id or partial id
		for _, s := range c.scopes {
			if strings.HasPrefix(s.id, *idName) {
				return []*Scope{s}, nil
			}
		}

		return nil, ResourceNotFoundError{error: fmt.Errorf("scope %s not found", *idName)}
	}

	_scopes := make([]*Scope, len(c.scopes))
	// list all scopes
	i := 0
	for _, scope := range c.scopes {
		_scopes[i] = scope
		i++
	}

	return _scopes, nil
}

func (c *Context) Scopes(idName *string) ([]*Scope, error) {
	c.Lock()
	defer c.Unlock()

	return c.findScope(idName)
}

func (c *Context) DefaultScope() *Scope {
	return c.defaultScope
}

func (c *Context) BindContainer(h *exec.Handle) error {
	c.Lock()
	defer c.Unlock()

	con, ok := c.containers[h.Container.ID]
	if !ok {
		return ResourceNotFoundError{error: fmt.Errorf("container %s not found", h.Container.ID)}
	}

	scopes := con.Scopes()
	if len(scopes) == 0 {
		return nil // container not part of any scopes
	}

	var err error
	var bound []*Scope
	defer func() {
		if err == nil {
			return
		}

		for _, s := range bound {
			s.unbindContainer(con)
		}
	}()

	for _, s := range scopes {
		if err = s.bindContainer(con); err != nil {
			return err
		}

		bound = append(bound, s)
	}

	updateMetadata(h, con.Endpoints())
	return nil
}

func (c *Context) UnbindContainer(h *exec.Handle) error {
	c.Lock()
	defer c.Unlock()

	con, ok := c.containers[h.Container.ID]
	if !ok {
		return ResourceNotFoundError{error: fmt.Errorf("container %s not found", h.Container.ID)}
	}

	scopes := con.Scopes()
	if len(scopes) == 0 {
		return nil // container not part of any scopes
	}

	var err error
	var unbound []*Scope
	defer func() {
		if err == nil {
			return
		}

		for _, s := range unbound {
			s.bindContainer(con)
		}
	}()

	for _, s := range scopes {
		if err = s.unbindContainer(con); err != nil {
			return err
		}

		unbound = append(unbound, s)
	}

	updateMetadata(h, con.Endpoints())
	return nil
}

func updateMetadata(h *exec.Handle, endpoints []*Endpoint) {
	h.SetSpec(nil)
	if h.ExecConfig.Networks == nil {
		h.ExecConfig.Networks = make(map[string]*metadata.NetworkEndpoint)
	}

	for _, e := range endpoints {
		h.ExecConfig.Networks[e.Scope().Name()] = e.metadataEndpoint()
	}
}

func findSlotNumber(slots map[int32]bool) int32 {
	// see https://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=2047927
	slot := pciSlotNumberBegin
	for _, ok := slots[slot]; ok && slot != pciSlotNumberEnd; {
		slot += pciSlotNumberInc
		_, ok = slots[slot]
	}

	if slot == pciSlotNumberEnd {
		return spec.NilSlot
	}

	return slot
}

var addEthernetCard = func(h *exec.Handle, s *Scope) (types.BaseVirtualDevice, error) {
	var err error
	var devices object.VirtualDeviceList
	var d types.BaseVirtualDevice
	var dc types.BaseVirtualDeviceConfigSpec

	dcs, err := h.Spec.FindNICs(s.NetworkName)
	for _, ds := range dcs {
		if ds.GetVirtualDeviceConfigSpec().Operation == types.VirtualDeviceConfigSpecOperationAdd {
			d = ds.GetVirtualDeviceConfigSpec().Device
		}
	}

	if d == nil {
		backing := &types.VirtualEthernetCardNetworkBackingInfo{
			VirtualDeviceDeviceBackingInfo: types.VirtualDeviceDeviceBackingInfo{
				DeviceName: s.NetworkName,
			},
		}

		if d, err = devices.CreateEthernetCard("vmxnet3", backing); err != nil {
			return nil, err
		}
	}

	if spec.VirtualDeviceSlotNumber(d) == spec.NilSlot {
		slots := make(map[int32]bool)
		for _, e := range h.ExecConfig.Networks {
			if e.PCISlot > 0 {
				slots[e.PCISlot] = true
			}
		}

		for _, slot := range h.Spec.CollectSlotNumbers() {
			slots[slot] = true
		}

		slot := findSlotNumber(slots)
		if slot == spec.NilSlot {
			return nil, fmt.Errorf("out of slots")
		}

		d.GetVirtualDevice().SlotInfo = &types.VirtualDevicePciBusSlotInfo{PciSlotNumber: slot}
	}

	if dc == nil {
		devices = append(devices, d)
		deviceSpecs, err := devices.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
		if err != nil {
			return nil, err
		}

		h.Spec.DeviceChange = append(h.Spec.DeviceChange, deviceSpecs...)
	}

	return d, nil
}

func (c *Context) resolveScope(scope string) (*Scope, error) {
	var s *Scope
	switch scope {
	// docker's default network, usually maps to the default bridge network
	case "default":
		s = c.DefaultScope()

	default:
		scopes, err := c.findScope(&scope)
		if err != nil || len(scopes) != 1 {
			return nil, err
		}

		// should have only one match at this point
		s = scopes[0]
	}

	return s, nil
}

// AddContainer add a container to the specified scope, optionally specifying an ip address
// for the container in the scope
func (c *Context) AddContainer(h *exec.Handle, scope string, ip *net.IP) (*Endpoint, error) {
	c.Lock()
	defer c.Unlock()

	if h == nil {
		return nil, fmt.Errorf("handle is required")
	}

	var err error
	s, err := c.resolveScope(scope)
	if err != nil {
		return nil, err
	}

	con, ok := c.containers[h.Container.ID]
	if !ok {
		con = &Container{id: h.Container.ID}
	}

	e, err := s.addContainer(con, ip)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			s.removeContainer(con)
		}
	}()

	addNIC := true
	if s.Type() == bridgeScopeType {
		for _, ec := range con.endpoints {
			if s == ec.Scope() {
				continue
			}

			if ec.Scope().Type() == bridgeScopeType {
				e.pciSlot = ec.pciSlot
				addNIC = false
				break
			}
		}
	}

	h.SetSpec(nil)

	if addNIC {
		var d types.BaseVirtualDevice
		d, err = addEthernetCard(h, s)
		if err != nil {
			return nil, err
		}

		e.pciSlot = spec.VirtualDeviceSlotNumber(d)
	}

	c.containers[con.ID()] = con
	updateMetadata(h, []*Endpoint{e})
	return e, nil
}

func (c *Context) RemoveContainer(h *exec.Handle, scope string) error {
	c.Lock()
	defer c.Unlock()

	if h == nil {
		return fmt.Errorf("handle is required")
	}

	var err error
	s, err := c.resolveScope(scope)
	if err != nil {
		return err
	}

	con, ok := c.containers[h.Container.ID]
	if !ok {
		return ResourceNotFoundError{error: fmt.Errorf("container %s not found", h.Container.ID)}
	}

	e := con.Endpoint(s)
	if e == nil {
		return ResourceNotFoundError{error: fmt.Errorf("endpoint for container %s not found in scope %s", con.ID(), s.Name())}
	}

	if err = s.removeContainer(con); err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		var ip *net.IP
		if e.static {
			i := e.IP()
			ip = &i
		}

		s.addContainer(con, ip)
	}()

	// remove NIC if no other scopes in the container need
	// it
	removeNIC := true
	if s.Type() == bridgeScopeType {
		for _, ec := range con.Endpoints() {
			if ec.Scope().Type() == bridgeScopeType {
				removeNIC = false
				break
			}
		}
	}

	if removeNIC {
		// ensure spec is not nil
		h.SetSpec(nil)

		var devices object.VirtualDeviceList
		backing := &types.VirtualEthernetCardNetworkBackingInfo{
			VirtualDeviceDeviceBackingInfo: types.VirtualDeviceDeviceBackingInfo{
				DeviceName: s.NetworkName,
			},
		}

		d, err := devices.CreateEthernetCard("vmxnet3", backing)
		if err != nil {
			return err
		}

		devices = append(devices, d)
		spec, err := devices.ConfigSpec(types.VirtualDeviceConfigSpecOperationRemove)
		if err != nil {
			return err
		}
		h.Spec.DeviceChange = append(h.Spec.DeviceChange, spec...)
	}

	// remove metadata
	delete(h.ExecConfig.Networks, scope)

	return nil
}

func (c *Context) Container(id exec.ID) *Container {
	c.Lock()
	defer c.Unlock()

	if con, ok := c.containers[id]; ok {
		return con
	}

	return nil
}

func (c *Context) DeleteScope(name string) error {
	c.Lock()
	defer c.Unlock()

	s, err := c.resolveScope(name)
	if err != nil {
		return err
	}

	if s == nil {
		return ResourceNotFoundError{}
	}

	if s.builtin {
		return fmt.Errorf("cannot remove builtin scope")
	}

	// check if any of the scope's endpoints are bound
	for _, e := range s.Endpoints() {
		if e.IsBound() {
			return fmt.Errorf("scope has bound endpoints")
		}
	}

	delete(c.scopes, name)
	return nil
}
