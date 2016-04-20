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

	"golang.org/x/net/context"

	"github.com/vmware/vic/pkg/vsphere/guest"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/vm"
)

const (
	bridgeNetworkKey = "guestinfo.vch/networks/bridge"
)

type Context struct {
	sync.Mutex

	defaultBridgePool *AddressSpace
	defaultBridgeMask net.IPMask

	scopes       map[string]*Scope
	defaultScope *Scope

	BridgeNetworkName string // Portgroup name of the bridge network
}

var getBridgeNetworkName = func(sess *session.Session) (string, error) {
	c := context.Background()
	vch, err := guest.GetSelf(c, sess)
	if err != nil {
		return "", fmt.Errorf("unable to get VCH ref")
	}

	vchVM := vm.NewVirtualMachine(c, sess, vch.Reference())
	guestInfo, err := vchVM.FetchExtraConfig(c)
	if err != nil {
		return "", err
	}
	return guestInfo[bridgeNetworkKey], nil
}

func NewContext(bridgePool net.IPNet, bridgeMask net.IPMask, sess *session.Session) (*Context, error) {
	pones, pbits := bridgePool.Mask.Size()
	mones, mbits := bridgeMask.Size()
	if pbits != mbits || mones < pones {
		return nil, fmt.Errorf("bridge mask is not compatiable with bridge pool mask")
	}

	bnn, err := getBridgeNetworkName(sess)
	if err != nil {
		return nil, err
	}

	ctx := &Context{
		defaultBridgeMask: bridgeMask,
		defaultBridgePool: NewAddressSpaceFromNetwork(&bridgePool),
		scopes:            make(map[string]*Scope),
		BridgeNetworkName: bnn,
	}

	s, err := ctx.NewScope("bridge", "bridge", nil, net.IPv4(0, 0, 0, 0), nil, nil)
	if err != nil {
		return nil, err
	}

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
		containers:  make(map[string]*Container),
		scopeType:   scopeType,
		space:       space,
		dns:         dns,
		NetworkName: networkName,
	}

	c.scopes[name] = newScope

	return newScope, nil
}

func (c *Context) newBridgeScope(id, name string, subnet *net.IPNet, gateway net.IP, dns []net.IP, ipam *IPAM) (newScope *Scope, err error) {
	s, err := c.newScopeCommon(id, name, "bridge", subnet, gateway, dns, ipam, c.BridgeNetworkName)
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
	return c.newScopeCommon(id, name, "external", subnet, gateway, dns, ipam, "")
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
	case "bridge":
		return c.newBridgeScope(generateID(), name, subnet, gateway, dns, &IPAM{pools: pools})

	case "external":
		return c.newExternalScope(generateID(), name, subnet, gateway, dns, &IPAM{pools: pools})

	default:
		return nil, fmt.Errorf("scope type not supported")
	}
}

func (c *Context) Scopes(idName *string) ([]*Scope, error) {
	c.Lock()
	defer c.Unlock()

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

		return nil, ResourceNotFoundError{}
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

func (c *Context) DefaultScope() *Scope {
	return c.defaultScope
}
