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

// +build !windows,!darwin

package tether

import (
	"net"

	"github.com/vmware/vic/pkg/vsphere/toolbox"
)

// Toolbox is a tether extension that wraps toolbox.Service
type Toolbox struct {
	*toolbox.Service
}

// NewToolbox returns a tether.Extension that wraps the vsphere/toolbox service
func NewToolbox() *Toolbox {
	in := toolbox.NewBackdoorChannelIn()
	out := toolbox.NewBackdoorChannelOut()

	service := toolbox.NewService(in, out)

	toolbox.RegisterVixRelayedCommandHandler(service)

	toolbox.RegisterPowerCommandHandler(service)

	return &Toolbox{Service: service}
}

// Start implementation of the tether.Extension interface
func (t *Toolbox) Start() error {
	t.Service.PrimaryIP = t.defaultIP

	return t.Service.Start()
}

// Stop implementation of the tether.Extension interface
func (t *Toolbox) Stop() error {
	t.Service.Stop()

	return nil
}

// Reload implementation of the tether.Extension interface
func (*Toolbox) Reload(config *ExecutorConfig) error {
	return nil
}

// externalIP attempts to find an external IP to be reported as the guest IP
func (t *Toolbox) externalIP() string {
	netif, err := net.InterfaceByName("client")
	if err != nil {
		return ""
	}

	addrs, err := netif.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}

	return ""
}

// defaultIP tries externalIP, falling back to toolbox.DefaultIP()
func (t *Toolbox) defaultIP() string {
	ip := t.externalIP()
	if ip != "" {
		return ip
	}

	return toolbox.DefaultIP()
}
