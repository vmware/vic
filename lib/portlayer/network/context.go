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
	"strconv"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-connections/nat"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/portlayer/event/events"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
)

const (
	pciSlotNumberBegin int32 = 0xc0
	pciSlotNumberEnd   int32 = 1 << 10
	pciSlotNumberInc   int32 = 1 << 5

	DefaultBridgeName = "bridge"
)

// Context denotes a networking context that represents a set of scopes, endpoints,
// and containers. Each context has its own separate IPAM.
type Context struct {
	sync.Mutex

	config *Configuration

	defaultBridgePool *AddressSpace
	defaultBridgeMask net.IPMask

	scopes       map[string]*Scope
	containers   map[string]*Container
	defaultScope *Scope
}

type AddContainerOptions struct {
	Scope   string
	IP      *net.IP
	Aliases []string
	Ports   []string
}

func NewContext(bridgePool net.IPNet, bridgeMask net.IPMask, config *Configuration) (*Context, error) {
	defer trace.End(trace.Begin(""))
	if config == nil {
		return nil, fmt.Errorf("missing config")
	}

	pones, pbits := bridgePool.Mask.Size()
	mones, mbits := bridgeMask.Size()
	if pbits != mbits || mones < pones {
		return nil, fmt.Errorf("bridge mask is not compatible with bridge pool mask")
	}

	ctx := &Context{
		config:            config,
		defaultBridgeMask: bridgeMask,
		defaultBridgePool: NewAddressSpaceFromNetwork(&bridgePool),
		scopes:            make(map[string]*Scope),
		containers:        make(map[string]*Container),
	}

	s, err := ctx.newBridgeScope(uid.New(), DefaultBridgeName, nil, net.IPv4(0, 0, 0, 0), nil, &IPAM{})
	if err != nil {
		return nil, err
	}
	s.builtin = true
	s.dns = []net.IP{s.gateway}
	ctx.defaultScope = s

	// add any bridge/external networks
	for nn, n := range ctx.config.ContainerNetworks {
		if nn == DefaultBridgeName {
			continue
		}

		pools := make([]string, len(n.Pools))
		for i, p := range n.Pools {
			pools[i] = p.String()
		}

		s, err := ctx.NewScope(n.Type, nn, &net.IPNet{IP: n.Gateway.IP.Mask(n.Gateway.Mask), Mask: n.Gateway.Mask}, n.Gateway.IP, n.Nameservers, pools)
		if err != nil {
			return nil, err
		}

		if n.Type == constants.ExternalScopeType {
			s.builtin = true
		}
	}

	// subscribe to the event stream for Vm events
	sub := fmt.Sprintf("%s(%p)", "netCtx", ctx)
	if exec.Config.EventManager != nil {
		exec.Config.EventManager.Subscribe(events.NewEventType(events.ContainerEvent{}).Topic(), sub, ctx.handleEvent)
	}

	return ctx, nil
}

func reserveGateway(gateway net.IP, subnet *net.IPNet, ipam *IPAM) (net.IP, error) {
	defer trace.End(trace.Begin(""))
	if ip.IsUnspecifiedSubnet(subnet) {
		return nil, fmt.Errorf("cannot reserve gateway for nil subnet")
	}

	if !ip.IsUnspecifiedIP(gateway) {
		// verify gateway is routable address
		if !ip.IsRoutableIP(gateway, subnet) {
			return nil, fmt.Errorf("gateway address %s is not routable on network %s", gateway, subnet)
		}

		// optionally reserve it in one of the pools
		for _, p := range ipam.spaces {
			if err := p.ReserveIP4(gateway); err == nil {
				break
			}
		}

		return gateway, nil
	}

	// gateway is not specified, pick one from the available pools
	if len(ipam.spaces) > 0 {
		var err error
		if gateway, err = ipam.spaces[0].ReserveNextIP4(); err != nil {
			return nil, err
		}

		if !ip.IsRoutableIP(gateway, subnet) {
			return nil, fmt.Errorf("gateway address %s is not routable on network %s", gateway, subnet)
		}

		return gateway, nil
	}

	return nil, fmt.Errorf("could not reserve gateway address for network %s", subnet)
}

func (c *Context) newScopeCommon(id uid.UID, name, scopeType string, subnet *net.IPNet, gateway net.IP, dns []net.IP, ipam *IPAM, network object.NetworkReference) (*Scope, error) {
	defer trace.End(trace.Begin(""))
	var err error
	var space *AddressSpace
	var defaultPool bool
	var allzeros, allones net.IP

	// cleanup
	defer func() {
		if err == nil || space == nil || !defaultPool {
			return
		}

		for _, p := range ipam.spaces {
			// release DNS IPs
			for _, d := range dns {
				p.ReleaseIP4(d)
			}

			// release gateway
			if !ip.IsUnspecifiedIP(gateway) {
				p.ReleaseIP4(gateway)
			}

			// release all-ones and all-zeros addresses
			if !ip.IsUnspecifiedIP(allzeros) {
				p.ReleaseIP4(allzeros)
			}
			if !ip.IsUnspecifiedIP(allones) {
				p.ReleaseIP4(allones)
			}
		}

		c.defaultBridgePool.ReleaseIP4Range(space)
	}()

	// subnet may not be specified, e.g. for "external" networks
	if !ip.IsUnspecifiedSubnet(subnet) {
		// allocate the subnet
		space, defaultPool, err = c.reserveSubnet(subnet)
		if err != nil {
			return nil, err
		}

		subnet = space.Network

		ipam.spaces, err = reservePools(space, ipam)
		if err != nil {
			return nil, err
		}

		// reserve all-ones and all-zeros addresses, which are not routable and so
		// should not be handed out
		allones = ip.AllOnesAddr(subnet)
		allzeros = ip.AllZerosAddr(subnet)
		for _, p := range ipam.spaces {
			p.ReserveIP4(allones)
			p.ReserveIP4(allzeros)

			// reserve DNS IPs
			for _, d := range dns {
				if d.Equal(gateway) {
					continue // gateway will be reserved later
				}

				p.ReserveIP4(d)
			}
		}

		if gateway, err = reserveGateway(gateway, subnet, ipam); err != nil {
			return nil, err
		}

	}

	newScope := &Scope{
		id:         id,
		name:       name,
		subnet:     *subnet,
		gateway:    gateway,
		ipam:       ipam,
		containers: make(map[uid.UID]*Container),
		scopeType:  scopeType,
		space:      space,
		dns:        dns,
		builtin:    false,
		network:    network,
	}

	c.scopes[name] = newScope

	return newScope, nil
}

func (c *Context) newBridgeScope(id uid.UID, name string, subnet *net.IPNet, gateway net.IP, dns []net.IP, ipam *IPAM) (newScope *Scope, err error) {
	defer trace.End(trace.Begin(""))
	bnPG, ok := c.config.PortGroups[c.config.BridgeNetwork]
	if !ok || bnPG == nil {
		return nil, fmt.Errorf("bridge network not set")
	}

	if ip.IsUnspecifiedSubnet(subnet) {
		// get the next available subnet from the default bridge pool
		var err error
		subnet, err = c.defaultBridgePool.NextIP4Net(c.defaultBridgeMask)
		if err != nil {
			return nil, err
		}
	}

	s, err := c.newScopeCommon(id, name, constants.BridgeScopeType, subnet, gateway, dns, ipam, bnPG)
	if err != nil {
		return nil, err
	}

	// add the gateway address to the bridge interface
	if err = c.config.BridgeLink.AddrAdd(net.IPNet{IP: s.Gateway(), Mask: s.Subnet().Mask}); err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EEXIST {
			log.Warnf("failed to add gateway address %s to bridge interface: %s", s.Gateway(), err)
		}
	}

	return s, nil
}

func (c *Context) newExternalScope(id uid.UID, name string, subnet *net.IPNet, gateway net.IP, dns []net.IP, ipam *IPAM) (*Scope, error) {
	defer trace.End(trace.Begin(""))
	// ipam cannot be specified without gateway and subnet
	if ipam != nil && len(ipam.pools) > 0 {
		if ip.IsUnspecifiedSubnet(subnet) || gateway.IsUnspecified() {
			return nil, fmt.Errorf("ipam cannot be specified without gateway and subnet for external network")
		}
	}

	if !ip.IsUnspecifiedSubnet(subnet) {
		// cannot overlap with the default bridge pool
		if c.defaultBridgePool.Network.Contains(subnet.IP) ||
			c.defaultBridgePool.Network.Contains(highestIP4(subnet)) {
			return nil, fmt.Errorf("external network cannot overlap with default bridge network")
		}
	}

	nPG := c.config.PortGroups[name]
	if nPG == nil {
		return nil, fmt.Errorf("no network info for external scope %s", name)
	}

	return c.newScopeCommon(id, name, constants.ExternalScopeType, subnet, gateway, dns, ipam, nPG)
}

func (c *Context) reserveSubnet(subnet *net.IPNet) (*AddressSpace, bool, error) {
	defer trace.End(trace.Begin(""))
	err := c.checkNetOverlap(subnet)
	if err != nil {
		return nil, false, err
	}

	// reserve from the default pool first
	space, err := c.defaultBridgePool.ReserveIP4Net(subnet)
	if err == nil {
		return space, true, nil
	}

	space = NewAddressSpaceFromNetwork(subnet)
	return space, false, nil
}

func (c *Context) checkNetOverlap(subnet *net.IPNet) error {
	// check if the requested subnet is available
	highestIP := highestIP4(subnet)
	for _, scope := range c.scopes {
		if scope.subnet.Contains(subnet.IP) || scope.subnet.Contains(highestIP) {
			return fmt.Errorf("subnet %s overlaps with scope %s subnet %s", subnet, scope.Name(), scope.Subnet())
		}
	}

	return nil
}

func reservePools(space *AddressSpace, ipam *IPAM) ([]*AddressSpace, error) {
	defer trace.End(trace.Begin(""))
	if len(ipam.pools) == 0 {
		// pool not specified so use the entire space
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
		r := ip.ParseRange(p)
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

func (c *Context) NewScope(scopeType, name string, subnet *net.IPNet, gateway net.IP, dns []net.IP, pools []string) (*Scope, error) {
	defer trace.End(trace.Begin(""))
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

	var s *Scope
	var err error
	switch scopeType {
	case constants.BridgeScopeType:
		s, err = c.newBridgeScope(uid.New(), name, subnet, gateway, dns, &IPAM{pools: pools})

	case constants.ExternalScopeType:
		s, err = c.newExternalScope(uid.New(), name, subnet, gateway, dns, &IPAM{pools: pools})

	default:
		return nil, fmt.Errorf("scope type not supported")
	}

	if err != nil {
		return nil, err
	}

	// add the new scope to the config
	c.config.ContainerNetworks[s.Name()] = &executor.ContainerNetwork{
		Common: executor.Common{
			ID:   s.ID().String(),
			Name: s.Name(),
		},
		Type:        s.Type(),
		Gateway:     net.IPNet{IP: s.Gateway(), Mask: s.Subnet().Mask},
		Nameservers: s.DNS(),
		Pools:       s.IPAM().Pools(),
	}
	c.config.PortGroups[s.Name()] = s.network

	// write config
	c.config.Encode()

	return s, nil
}

func (c *Context) findScopes(idName *string) ([]*Scope, error) {
	defer trace.End(trace.Begin(""))
	if idName != nil && *idName != "" {
		if *idName == "default" {
			return []*Scope{c.DefaultScope()}, nil
		}

		// search by name
		scope, ok := c.scopes[*idName]
		if ok {
			return []*Scope{scope}, nil
		}

		// search by id or partial id
		for _, s := range c.scopes {
			if strings.HasPrefix(s.id.String(), *idName) {
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

	return c.findScopes(idName)
}

func (c *Context) DefaultScope() *Scope {
	return c.defaultScope
}

func (c *Context) BindContainer(h *exec.Handle) ([]*Endpoint, error) {
	defer trace.End(trace.Begin(""))
	c.Lock()
	defer c.Unlock()

	con, err := c.container(h)
	if con != nil {
		return con.Endpoints(), nil // already bound
	}

	if _, ok := err.(ResourceNotFoundError); !ok {
		return nil, err
	}

	con = &Container{
		id:   uid.Parse(h.ExecConfig.ID),
		name: h.ExecConfig.Name,
	}

	defaultMarked := false
	aliases := make(map[string]*Container)
	var endpoints []*Endpoint
	for _, ne := range h.ExecConfig.Networks {
		var s *Scope
		s, ok := c.scopes[ne.Network.Name]
		if !ok {
			return nil, &ResourceNotFoundError{}
		}

		defer func() {
			if err == nil {
				return
			}

			s.removeContainer(con)
		}()

		var ip *net.IP
		if ne.Static != nil {
			ip = &ne.Static.IP
		}

		var e *Endpoint
		if e, err = s.addContainer(con, ip); err != nil {
			return nil, err
		}

		ports, _, err := nat.ParsePortSpecs(ne.Ports)
		if err != nil {
			return nil, err
		}
		for p := range ports {
			var port Port
			if port, err = ParsePort(string(p)); err != nil {
				return nil, err
			}

			if err = e.addPort(port); err != nil {
				return nil, err
			}
		}

		eip := e.IP()
		if eip != nil && !eip.IsUnspecified() {
			ne.Static = &net.IPNet{
				IP:   eip,
				Mask: e.Scope().Subnet().Mask,
			}
		}
		ne.Network.Gateway = net.IPNet{IP: e.gateway, Mask: e.subnet.Mask}
		ne.Network.Nameservers = make([]net.IP, len(s.dns))
		copy(ne.Network.Nameservers, s.dns)

		// mark the external network as default
		if !defaultMarked && e.Scope().Type() == constants.ExternalScopeType {
			defaultMarked = true
			ne.Network.Default = true
		}

		// dns lookup aliases
		aliases[fmt.Sprintf("%s:%s", s.Name(), con.name)] = con
		aliases[fmt.Sprintf("%s:%s", s.Name(), con.id.Truncate())] = con

		// container specific aliases
		for _, a := range ne.Network.Aliases {
			log.Debugf("adding alias %s", a)
			l := strings.Split(a, ":")
			if len(l) != 2 {
				err = fmt.Errorf("Parsing network alias %s failed", a)
				return nil, err
			}

			who, what := l[0], l[1]
			if who == "" {
				who = con.name
			}
			if a, exists := e.addAlias(who, what); a != badAlias && !exists {
				whoc := con
				// if the alias is not for this container, then
				// find it in the container collection
				if who != con.name {
					whoc = c.containers[who]
				}

				// whoc may be nil here, which means that the aliased
				// container is not bound yet; this is OK, and will be
				// fixed up when "who" is bound
				if whoc != nil {
					aliases[a.scopedName()] = whoc
				}
			}
		}

		// fix up the aliases to this container
		// from other containers
		for _, e := range s.Endpoints() {
			if e.Container() == con {
				continue
			}

			for _, a := range e.getAliases(con.name) {
				aliases[a.scopedName()] = con
			}
		}

		endpoints = append(endpoints, e)
	}

	// verify all the aliases to be added do not conflict with
	// existing container keys
	for a := range aliases {
		if _, ok := c.containers[a]; ok {
			return nil, fmt.Errorf("duplicate alias %s for container %s", a, con.ID())
		}
	}

	// FIXME: if there was no external network to mark as default,
	// then just pick the first network to mark as default
	if !defaultMarked {
		defaultMarked = true
		for _, ne := range h.ExecConfig.Networks {
			ne.Network.Default = true
			break
		}
	}

	// long id
	c.containers[con.id.String()] = con
	// short id
	c.containers[con.id.Truncate().String()] = con
	// name
	c.containers[con.name] = con
	// aliases
	for k, v := range aliases {
		log.Debugf("adding alias %s -> %s", k, v.Name())
		c.containers[k] = v
	}

	return endpoints, nil
}

func (c *Context) container(h *exec.Handle) (*Container, error) {
	defer trace.End(trace.Begin(""))
	id := uid.Parse(h.ExecConfig.ID)
	if id == uid.NilUID {
		return nil, fmt.Errorf("invalid container id %s", h.ExecConfig.ID)
	}

	if con, ok := c.containers[id.String()]; ok {
		return con, nil
	}

	return nil, ResourceNotFoundError{error: fmt.Errorf("container %s not found", id.String())}
}

func (c *Context) UnbindContainer(h *exec.Handle) ([]*Endpoint, error) {
	defer trace.End(trace.Begin(""))
	c.Lock()
	defer c.Unlock()

	con, err := c.container(h)
	if err != nil {
		if _, ok := err.(ResourceNotFoundError); ok {
			return nil, nil // not bound
		}

		return nil, err
	}

	// aliases to remove
	var aliases []string
	var endpoints []*Endpoint
	for _, ne := range h.ExecConfig.Networks {
		var s *Scope
		s, ok := c.scopes[ne.Network.Name]
		if !ok {
			return nil, &ResourceNotFoundError{}
		}

		defer func() {
			if err == nil {
				return
			}

			var ip *net.IP
			if ne.Static != nil {
				ip = &ne.Static.IP
			}
			s.addContainer(con, ip)
		}()

		// save the endpoint info
		e := con.Endpoint(s).copy()

		if err = s.removeContainer(con); err != nil {
			return nil, err
		}

		if !e.static {
			ne.Static = nil
		}

		// aliases to remove
		// name for dns lookup
		aliases = append(aliases, fmt.Sprintf("%s:%s", s.Name(), con.name))
		aliases = append(aliases, fmt.Sprintf("%s:%s", s.Name(), con.id.Truncate()))
		for _, as := range e.aliases {
			for _, a := range as {
				aliases = append(aliases, a.scopedName())
			}
		}

		// aliases from other containers
		for _, e := range s.Endpoints() {
			if e.Container() == con {
				continue
			}

			for _, a := range e.getAliases(con.name) {
				aliases = append(aliases, a.scopedName())
			}
		}

		endpoints = append(endpoints, e)
	}

	// remove aliases
	for _, a := range aliases {
		delete(c.containers, a)
	}

	// long id
	delete(c.containers, con.ID().String())
	// short id
	delete(c.containers, con.ID().Truncate().String())
	// name
	delete(c.containers, con.Name())

	return endpoints, nil
}

var addEthernetCard = func(h *exec.Handle, s *Scope) (types.BaseVirtualDevice, error) {
	var devices object.VirtualDeviceList
	var d types.BaseVirtualDevice
	var dc types.BaseVirtualDeviceConfigSpec

	ctx := context.Background()
	dcs, err := h.Spec.FindNICs(ctx, s.network)
	if err != nil {
		return nil, err
	}

	for _, ds := range dcs {
		if ds.GetVirtualDeviceConfigSpec().Operation == types.VirtualDeviceConfigSpecOperationAdd {
			d = ds.GetVirtualDeviceConfigSpec().Device
			dc = ds
			break
		}
	}

	if d == nil {
		backing, err := s.network.EthernetCardBackingInfo(ctx)
		if err != nil {
			return nil, err
		}

		if d, err = devices.CreateEthernetCard("vmxnet3", backing); err != nil {
			return nil, err
		}

		d.GetVirtualDevice().DeviceInfo = &types.Description{
			Label: s.name,
		}
	}

	if spec.VirtualDeviceSlotNumber(d) == spec.NilSlot {
		slots := make(map[int32]bool)
		for _, e := range h.ExecConfig.Networks {
			if e.Common.ID != "" {
				slot, err := strconv.Atoi(e.Common.ID)
				if err == nil {
					slots[int32(slot)] = true
				}
			}
		}

		h.Spec.AssignSlotNumber(d, slots)
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
	scopes, err := c.findScopes(&scope)
	if err != nil || len(scopes) != 1 {
		return nil, err
	}

	return scopes[0], nil
}

// AddContainer add a container to the specified scope, optionally specifying an ip address
// for the container in the scope
func (c *Context) AddContainer(h *exec.Handle, options *AddContainerOptions) error {
	defer trace.End(trace.Begin(""))
	c.Lock()
	defer c.Unlock()

	if h == nil {
		return fmt.Errorf("handle is required")
	}

	var err error
	s, err := c.resolveScope(options.Scope)
	if err != nil {
		return err
	}

	if h.ExecConfig.Networks != nil {
		if _, ok := h.ExecConfig.Networks[s.Name()]; ok {
			// already part of this scope
			return nil
		}

		// check if container is already part of an "external" scope;
		// only one "external" scope per container is allowed
		if s.Type() == constants.ExternalScopeType {
			for name := range h.ExecConfig.Networks {
				sc, _ := c.resolveScope(name)
				if sc.Type() == constants.ExternalScopeType {
					return fmt.Errorf("container can only be added to at most one mapped network")
				}
			}
		}
	}

	if err := h.SetSpec(nil); err != nil {
		return err
	}

	// figure out if we need to add a new NIC
	// if there is already a NIC connected to a
	// bridge network and we are adding the container
	// to a bridge network, we just reuse that
	// NIC
	var pciSlot int32
	if s.Type() == constants.BridgeScopeType {
		for _, ne := range h.ExecConfig.Networks {
			sc, err := c.resolveScope(ne.Network.Name)
			if err != nil {
				return err
			}

			if sc.Type() != constants.BridgeScopeType {
				continue
			}

			if ne.ID != "" {
				pciSlot = atoiOrZero(ne.ID)
				if pciSlot != 0 {
					break
				}
			}
		}
	}

	if pciSlot == 0 {
		d, err := addEthernetCard(h, s)
		if err != nil {
			return err
		}

		pciSlot = spec.VirtualDeviceSlotNumber(d)
	}

	if h.ExecConfig.Networks == nil {
		h.ExecConfig.Networks = make(map[string]*executor.NetworkEndpoint)
	}

	ne := &executor.NetworkEndpoint{
		Common: executor.Common{
			ID: strconv.Itoa(int(pciSlot)),
			// Name: this would cause NIC renaming if uncommented
		},
		Network: executor.ContainerNetwork{
			Common: executor.Common{
				Name: s.Name(),
			},
			Aliases: options.Aliases,
			Pools:   s.IPAM().Pools(),
		},
		Ports: options.Ports,
	}

	if options.IP != nil && !options.IP.IsUnspecified() {
		ne.Static = &net.IPNet{
			IP:   *options.IP,
			Mask: s.Subnet().Mask,
		}
	}

	h.ExecConfig.Networks[s.Name()] = ne
	return nil
}

func (c *Context) RemoveContainer(h *exec.Handle, scope string) error {
	defer trace.End(trace.Begin(""))
	c.Lock()
	defer c.Unlock()

	if h == nil {
		return fmt.Errorf("handle is required")
	}

	if con, _ := c.container(h); con != nil {
		return fmt.Errorf("container is bound")
	}

	var err error
	s, err := c.resolveScope(scope)
	if err != nil {
		return err
	}

	var ne *executor.NetworkEndpoint
	ne, ok := h.ExecConfig.Networks[s.Name()]
	if !ok {
		return fmt.Errorf("container %s not part of network %s", h.ExecConfig.ID, s.Name())
	}

	// figure out if any other networks are using the NIC
	removeNIC := true
	for _, ne2 := range h.ExecConfig.Networks {
		if ne2 == ne {
			continue
		}
		if ne2.ID == ne.ID {
			removeNIC = false
			break
		}
	}

	if removeNIC {
		// ensure spec is not nil
		h.SetSpec(nil)

		var devices object.VirtualDeviceList
		backing, err := s.network.EthernetCardBackingInfo(context.Background())
		if err != nil {
			return err
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

	delete(h.ExecConfig.Networks, s.Name())

	return nil
}

func (c *Context) Container(key string) *Container {
	c.Lock()
	defer c.Unlock()

	log.Debugf("container lookup for %s", key)
	if con, ok := c.containers[key]; ok {
		return con
	}

	return nil
}

func (c *Context) ContainerByAddr(addr net.IP) *Endpoint {
	c.Lock()
	defer c.Unlock()

	for _, s := range c.scopes {
		if e := s.ContainerByAddr(addr); e != nil {
			return e
		}
	}

	return nil
}

func (c *Context) DeleteScope(name string) error {
	defer trace.End(trace.Begin(""))
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

	if len(s.Endpoints()) != 0 {
		return fmt.Errorf("%s has active endpoints", s.Name())
	}

	if s.Type() == constants.BridgeScopeType {

		// remove gateway ip from bridge interface
		addr := net.IPNet{IP: s.Gateway(), Mask: s.Subnet().Mask}
		if err := c.config.BridgeLink.AddrDel(addr); err != nil {
			if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EADDRNOTAVAIL {
				log.Warnf("could not remove gateway address %s for scope %s on link %s: %s", addr, s.Name(), c.config.BridgeLink.Attrs().Name, err)
			}

			err = nil
		}
	}

	delete(c.scopes, s.Name())
	return nil
}

func (c *Context) UpdateContainer(h *exec.Handle) error {
	defer trace.End(trace.Begin(""))
	c.Lock()
	defer c.Unlock()

	con, err := c.container(h)
	if err != nil {
		return err
	}

	for _, s := range con.Scopes() {
		if !s.isDynamic() {
			continue
		}

		ne := h.ExecConfig.Networks[s.Name()]
		if ne == nil {
			return fmt.Errorf("container config does not have info for network scope %s", s.Name())
		}

		e := con.Endpoint(s)
		e.ip = ne.Assigned.IP
		gw, snet, err := net.ParseCIDR(ne.Network.Gateway.String())
		if err != nil {
			return err
		}

		e.gateway = gw
		e.subnet = *snet

		s.gateway = gw
		s.subnet = *snet
	}

	return nil
}

func atoiOrZero(a string) int32 {
	i, _ := strconv.Atoi(a)
	return int32(i)
}

// handleEvent processes events
func (c *Context) handleEvent(ie events.Event) {
	defer trace.End(trace.Begin(ie.String()))
	switch ie.String() {
	case events.ContainerPoweredOff:
		handle := exec.GetContainer(uid.Parse(ie.Reference()))
		if handle == nil {
			log.Errorf("Container %s not found - unable to UnbindContainer", ie.Reference())
			return
		}
		_, err := c.UnbindContainer(handle)
		if err != nil {
			log.Warnf("Failed to unbind container %s", ie.Reference())
		}
	}
	return
}
