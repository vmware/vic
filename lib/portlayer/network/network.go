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

	log "github.com/Sirupsen/logrus"
)

var (
	DefaultContext *Context
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

func Init() error {
	var err error

	bridgeRange := Config.BridgeIPRange
	if bridgeRange == nil || len(bridgeRange.IP) == 0 || bridgeRange.IP.IsUnspecified() {
		_, bridgeRange, err = net.ParseCIDR("172.16.0.0/12")
		if err != nil {
			return err
		}
	}

	// make sure a NIC attached to the bridge network exists
	Config.BridgeLink, err = getBridgeLink()
	if err != nil {
		return err
	}

	bridgeWidth := Config.BridgeNetworkWidth
	if bridgeWidth == nil || len(*bridgeWidth) == 0 {
		w := net.CIDRMask(16, 32)
		bridgeWidth = &w
	}

	DefaultContext, err = NewContext(*bridgeRange, *bridgeWidth)
	if err == nil {
		log.Infof("Default network context allocated: %s", bridgeRange.String())
	}

	return err
}

func getBridgeLink() (Link, error) {
	// add the gateway address to the bridge interface
	link, err := LinkByName(Config.BridgeNetwork)
	if err != nil {
		// lookup by alias
		return LinkByAlias(Config.BridgeNetwork)
	}

	return link, nil
}
