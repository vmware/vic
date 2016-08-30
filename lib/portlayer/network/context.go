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
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/ip"
	"github.com/vmware/vic/pkg/uid"
)

const (
	pciSlotNumberBegin int32 = 0xc0
	pciSlotNumberEnd   int32 = 1 << 10
	pciSlotNumberInc   int32 = 1 << 5

	DefaultScopeName = "bridge"
)

// Context denotes a networking context that represents a set of scopes, endpoints,
// and containers. Each context has its own separate IPAM.
type Context struct {
	sync.Mutex

	defaultBridgePool *AddressSpace
	defaultBridgeMask net.IPMask

	scopes       map[string]*Scope
	containers   map[uid.UID]*Container
	defaultScope *Scope
}

type AddContainerOptions struct {
	Scope   string
	IP      *net.IP
	Aliases []string
	Ports   []string
}

func NewContext(bridgePool net.IPNet, bridgeMask net.IPMask) (*Context, error) {
	pones, pbits := bridgePool.Mask.Size()
	mones, mbits := bridgeMask.Size()
	if pbits != mbits || mones < pones {
		return nil, fmt.Errorf("bridge mask is not compatible with bridge pool mask")
	}

	ctx := &Context{
		defaultBridgeMask: bridgeMask,
		defaultBridgePool: NewAddressSpaceFromNetwork(&bridgePool),
		scopes:            make(map[string]*Scope),
		containers:        make(map[uid.UID]*Container),
	}

	s, err := ctx.NewScope(DefaultScopeName, BridgeScopeType, nil, net.IPv4(0, 0, 0, 0), nil, nil)
	if err != nil {
		return nil, err
	}
	s.builtin = true
	s.dns = []net.IP{s.gateway}
	ctx.defaultScope = s

	// add any external networks
	for nn, n := range Config.ContainerNetworks {
		if nn == "bridge" {
			continue
		}

		pools := make([]string, len(n.Pools))
		for i, p := range n.Pools {
			pools[i] = p.String()
		}

		s, err := ctx.NewScope(ExternalScopeType, nn, &net.IPNet{IP: n.Gateway.IP.Mask(n.Gateway.Mask), Mask: n.Gateway.Mask}, n.Gateway.IP, n.Nameservers, pools)
		if err != nil {
			return nil, err
		}

		s.builtin = true
	}

	return ctx, nil
}

func reserveGateway(gateway net.IP, subnet *net.IPNet, ipam *IPAM) (net.IP, error) {
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
	bn, ok := Config.ContainerNetworks[Config.BridgeNetwork]
	if !ok || bn == nil {
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

	s, err := c.newScopeCommon(id, name, BridgeScopeType, subnet, gateway, dns, ipam, bn.PortGroup)
	if err != nil {
		return nil, err
	}

	// add the gateway address to the bridge interface
	if err = Config.BridgeLink.AddrAdd(net.IPNet{IP: s.Gateway(), Mask: s.Subnet().Mask}); err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EEXIST {
			log.Warnf("failed to add gateway address %s to bridge interface: %s", s.Gateway(), err)
		}
	}

	return s, nil
}

func (c *Context) newExternalScope(id uid.UID, name string, subnet *net.IPNet, gateway net.IP, dns []net.IP, ipam *IPAM) (*Scope, error) {
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

	n := Config.ContainerNetworks[name]
	if n == nil {
		return nil, fmt.Errorf("no network info for external scope %s", name)
	}

	return c.newScopeCommon(id, name, ExternalScopeType, subnet, gateway, dns, ipam, n.PortGroup)
}

func (c *Context) reserveSubnet(subnet *net.IPNet) (*AddressSpace, bool, error) {
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
	case BridgeScopeType:
		return c.newBridgeScope(uid.New(), name, subnet, gateway, dns, &IPAM{pools: pools})

	case ExternalScopeType:
		return c.newExternalScope(uid.New(), name, subnet, gateway, dns, &IPAM{pools: pools})

	default:
		return nil, fmt.Errorf("scope type not supported")
	}
}

func (c *Context) findScopes(idName *string) ([]*Scope, error) {
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
	c.Lock()
	defer c.Unlock()

	var con *Container
	var err error

	if len(h.ExecConfig.Networks) == 0 {
		return nil, fmt.Errorf("nothing to bind")
	}

	con, ok := c.containers[uid.Parse(h.ExecConfig.ID)]
	if ok {
		return nil, fmt.Errorf("container %s already bound", h.ExecConfig.ID)
	}

	con = &Container{
		id:   uid.Parse(h.ExecConfig.ID),
		name: h.ExecConfig.Name,
	}
	defaultMarked := false
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

		for _, p := range ne.Ports {
			var port Port
			if port, err = ParsePort(p); err != nil {
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
		if !defaultMarked && e.Scope().Type() == ExternalScopeType {
			defaultMarked = true
			ne.Network.Default = true
		}

		endpoints = append(endpoints, e)
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

	// local map to hold the container mapping
	containers := make(map[uid.UID]*Container)

	// Adding long id, short id and common name to the map to point same container
	// Last two is needed by DNS subsystem
	containers[con.id] = con

	tid := con.id.Truncate()
	cname := h.ExecConfig.Common.Name

	var key string
	// network scoped entries
	for i := range endpoints {
		e := endpoints[i]
		// scope name
		sname := e.Scope().Name()

		// SCOPE:SHORT ID
		key = fmt.Sprintf("%s:%s", sname, tid)
		log.Debugf("Adding %s to the containers", key)
		containers[uid.Parse(key)] = con

		// SCOPE:NAME
		key = fmt.Sprintf("%s:%s", sname, cname)
		log.Debugf("Adding %s to the containers", key)
		containers[uid.Parse(key)] = con

		ne, ok := h.ExecConfig.Networks[sname]
		if !ok {
			err := fmt.Errorf("Failed to find Network %s", sname)
			log.Errorf(err.Error())
			return nil, err
		}

		// Aliases/Links
		for i := range ne.Network.Aliases {
			l := strings.Split(ne.Network.Aliases[i], ":")
			if len(l) != 2 {
				err := fmt.Errorf("Parsing %s failed", l)
				log.Errorf(err.Error())
				return nil, err
			}
			who, what := l[0], l[1]
			// if who is empty string that means it is a alias
			// which points to the container itself
			if who == "" {
				who = cname
			}
			// Find the scope:who container
			key = fmt.Sprintf("%s:%s", sname, who)
			// search global map
			con, ok := c.containers[uid.Parse(key)]
			if !ok {
				// search local map
				con, ok = containers[uid.Parse(key)]
				if !ok {
					err := fmt.Errorf("Failed to find container %s", key)
					log.Errorf(err.Error())
					return nil, err
				}
			}
			log.Debugf("Found container %s", key)

			// Set scope:what to scope:who
			key = fmt.Sprintf("%s:%s", sname, what)
			log.Debugf("Adding %s to the containers", key)
			containers[uid.Parse(key)] = con
		}
	}

	// set the real map now that we are err free
	for k, v := range containers {
		c.containers[k] = v
	}

	return endpoints, nil
}

func (c *Context) UnbindContainer(h *exec.Handle) ([]*Endpoint, error) {
	c.Lock()
	defer c.Unlock()

	con, ok := c.containers[uid.Parse(h.ExecConfig.ID)]
	if !ok {
		return nil, ResourceNotFoundError{error: fmt.Errorf("container %s not found", h.ExecConfig.ID)}
	}

	// local map to hold the container mapping
	var containers []uid.UID

	// Removing long id, short id and common name from the map
	containers = append(containers, uid.Parse(h.ExecConfig.ID))

	tid := con.id.Truncate()
	cname := h.ExecConfig.Common.Name

	var key string
	var endpoints []*Endpoint
	var err error
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

		// scope name
		sname := e.Scope().Name()

		// delete scope:short id
		key = fmt.Sprintf("%s:%s", sname, tid)
		log.Debugf("Removing %s from the containers", key)
		containers = append(containers, uid.Parse(key))

		// delete scope:name
		key = fmt.Sprintf("%s:%s", sname, cname)
		log.Debugf("Removing %s from the containers", key)
		containers = append(containers, uid.Parse(key))

		// delete aliases
		for i := range ne.Network.Aliases {
			l := strings.Split(ne.Network.Aliases[i], ":")
			if len(l) != 2 {
				err := fmt.Errorf("Parsing %s failed", l)
				log.Errorf(err.Error())
				return nil, err
			}

			_, what := l[0], l[1]

			// delete scope:what
			key = fmt.Sprintf("%s:%s", sname, what)
			log.Debugf("Removing %s from the containers", key)
			containers = append(containers, uid.Parse(key))
		}

		endpoints = append(endpoints, e)
	}

	// delete from real map now that we are err free
	for i := range containers {
		delete(c.containers, containers[i])
	}

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
		if s.Type() == ExternalScopeType {
			for name := range h.ExecConfig.Networks {
				sc, _ := c.resolveScope(name)
				if sc.Type() == ExternalScopeType {
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
	if s.Type() == BridgeScopeType {
		for _, ne := range h.ExecConfig.Networks {
			sc, err := c.resolveScope(ne.Network.Name)
			if err != nil {
				return err
			}

			if sc.Type() != BridgeScopeType {
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
	c.Lock()
	defer c.Unlock()

	if h == nil {
		return fmt.Errorf("handle is required")
	}

	if _, ok := c.containers[uid.Parse(h.ExecConfig.ID)]; ok {
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

func (c *Context) Container(id uid.UID) *Container {
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

	if len(s.Endpoints()) != 0 {
		return fmt.Errorf("%s has active endpoints", s.Name())
	}

	if s.Type() == BridgeScopeType {

		// remove gateway ip from bridge interface
		addr := net.IPNet{IP: s.Gateway(), Mask: s.Subnet().Mask}
		if err := Config.BridgeLink.AddrDel(addr); err != nil {
			if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EADDRNOTAVAIL {
				log.Warnf("could not remove gateway address %s for scope %s on link %s: %s", addr, s.Name(), Config.BridgeLink.Attrs().Name, err)
			}

			err = nil
		}
	}

	delete(c.scopes, s.Name())
	return nil
}

func (c *Context) UpdateContainer(h *exec.Handle) error {
	c.Lock()
	defer c.Unlock()

	con := c.containers[uid.Parse(h.ExecConfig.ID)]
	if con == nil {
		return ResourceNotFoundError{}
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
