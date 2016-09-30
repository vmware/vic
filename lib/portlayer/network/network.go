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
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/session"
)

var (
	DefaultContext *Context

	initializer sync.Once
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
	var err error
	initializer.Do(func() {
		f := find.NewFinder(sess.Vim25(), false)

		var config Configuration
		config.sink = sink
		config.source = source
		config.Decode()
		config.PortGroups = make(map[string]object.NetworkReference)

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

			config.PortGroups[nn] = r.(object.NetworkReference)
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

		log.Infof("Default network context allocated: %s", bridgeRange.String())

		// populate existing containers
		state := exec.StateRunning
		for _, c := range exec.Containers.Containers(&state) {
			h := c.NewHandle()
			defer h.Close()

			if _, err = netctx.BindContainer(h); err != nil {
				return
			}
		}

		if err == nil {
			DefaultContext = netctx
		}
	})

	return err
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
