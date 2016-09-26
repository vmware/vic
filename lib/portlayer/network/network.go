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
	"context"
	"fmt"
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/portlayer/event"
	"github.com/vmware/vic/lib/portlayer/event/events"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
)

var (
	DefaultContext *Context

	initializer struct {
		err  error
		once sync.Once
	}
)

type DuplicateResourceError struct {
	resID string
}

type ResourceNotFoundError struct {
	error
}

func (e DuplicateResourceError) Error() string {
	return fmt.Sprintf("%s already exists", e.resID)
}

func Init(ctx context.Context, sess *session.Session, source extraconfig.DataSource, sink extraconfig.DataSink) error {
	trace.End(trace.Begin("network.Init"))

	initializer.once.Do(func() {
		var err error
		defer func() {
			initializer.err = err
		}()

		f := find.NewFinder(sess.Vim25(), false)

		var config Configuration
		config.sink = sink
		config.source = source
		config.Decode()

		log.Debugf("Decoded VCH config for network: %#v", config)
		for nn, n := range config.ContainerNetworks {
			pgref := new(types.ManagedObjectReference)
			if !pgref.FromString(n.ID) {
				log.Warnf("Could not reacquire object reference from id for network %s: %s", nn, n.ID)
			}

			r, err := f.ObjectReference(ctx, *pgref)
			if err != nil {
				log.Warnf("could not get network reference for %s network", nn)
				continue
			}

			n.PortGroup = r.(object.NetworkReference)
		}

		bridgeRange := config.BridgeIPRange
		if bridgeRange == nil || len(bridgeRange.IP) == 0 || bridgeRange.IP.IsUnspecified() {
			_, bridgeRange, err = net.ParseCIDR("172.16.0.0/12")
			if err != nil {
				return
			}
		}

		// make sure a NIC attached to the bridge network exists
		config.BridgeLink, err = getBridgeLink(&config)
		if err != nil {
			return
		}

		bridgeWidth := config.BridgeNetworkWidth
		if bridgeWidth == nil || len(*bridgeWidth) == 0 {
			w := net.CIDRMask(16, 32)
			bridgeWidth = &w
		}

		var netctx *Context
		netctx, err = NewContext(*bridgeRange, *bridgeWidth, &config)
		if err != nil {
			return
		}

		if err = engageContext(netctx, exec.Config.EventManager); err == nil {
			DefaultContext = netctx
			log.Infof("Default network context allocated: %s", bridgeRange.String())
		}
	})

	return initializer.err
}

// handleEvent processes events
func handleEvent(netctx *Context, ie events.Event) {
	switch ie.String() {
	case events.ContainerPoweredOff:
		handle := exec.GetContainer(uid.Parse(ie.Reference()))
		if handle == nil {
			log.Errorf("Container %s not found - unable to UnbindContainer", ie.Reference())
			return
		}
		_, err := netctx.UnbindContainer(handle)
		if err != nil {
			log.Warnf("Failed to unbind container %s", ie.Reference())
		}
	}
	return
}

// engageContext connects the given network context into a vsphere environment
// using an event manager, and a container cache. This hooks up a callback to
// react to vsphere events, as well as populate the context with any containers
// that are present.
func engageContext(netctx *Context, em event.EventManager) error {
	var err error

	// grab the context lock so that we do not unbind any containers
	// that stop out of band. this could cause, for example, for us
	// to bind a container when it has already been unbound by an
	// event callback
	netctx.Lock()
	defer netctx.Unlock()

	// subscribe to the event stream for Vm events
	if em == nil {
		return fmt.Errorf("event manager is required for default network context")
	}

	sub := fmt.Sprintf("%s(%p)", "netCtx", netctx)
	topic := events.NewEventType(events.ContainerEvent{}).Topic()
	em.Subscribe(topic, sub, func(ie events.Event) {
		handleEvent(netctx, ie)
	})

	defer func() {
		if err != nil {
			em.Unsubscribe(topic, sub)
		}
	}()

	state := exec.StateRunning
	for _, c := range exec.Containers.Containers(&state) {
		log.Debugf("adding container %s", c.ExecConfig.ID)
		h := c.NewHandle()
		defer h.Close()

		// add any user created networks that show up in container's config
		for n, ne := range h.ExecConfig.Networks {
			var s []*Scope
			s, err = netctx.findScopes(&n)
			if err != nil {
				if _, ok := err.(ResourceNotFoundError); !ok {
					return err
				}
			}

			if len(s) > 0 {
				continue
			}

			pools := make([]string, len(ne.Network.Pools))
			for i := range ne.Network.Pools {
				pools[i] = ne.Network.Pools[i].String()
			}

			log.Debugf("adding scope %s", n)
			if _, err = netctx.newScope(ne.Network.Type, n, nil, ne.Network.Gateway.IP, ne.Network.Nameservers, pools); err != nil {
				return err
			}
		}

		if _, err = netctx.bindContainer(h); err != nil {
			return err
		}
	}

	return nil
}

func getBridgeLink(config *Configuration) (Link, error) {
	// add the gateway address to the bridge interface
	link, err := LinkByName(config.BridgeNetwork)
	if err != nil {
		// lookup by alias
		return LinkByAlias(config.BridgeNetwork)
	}

	return link, nil
}
