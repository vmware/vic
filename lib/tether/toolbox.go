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
	"fmt"
	"net"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/cmd/tether/msgs"
	"github.com/vmware/vic/pkg/vsphere/toolbox"
)

// Toolbox is a tether extension that wraps toolbox.Service
type Toolbox struct {
	*toolbox.Service

	config *ExecutorConfig
}

// NewToolbox returns a tether.Extension that wraps the vsphere/toolbox service
func NewToolbox() *Toolbox {
	in := toolbox.NewBackdoorChannelIn()
	out := toolbox.NewBackdoorChannelOut()

	service := toolbox.NewService(in, out)

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
func (t *Toolbox) Reload(config *ExecutorConfig) error {
	t.config = config
	return nil
}

// InContainer configures the toolbox to run within a container VM
func (t *Toolbox) InContainer() *Toolbox {
	t.PowerCommand.Halt.Handler = t.halt

	return t
}

func (t *Toolbox) kill(name string) error {
	session := t.config.Sessions[t.config.ID]

	if name == "" {
		name = string(ssh.SIGTERM)
	}

	sig := new(msgs.SignalMsg)
	err := sig.FromString(name)
	if err != nil {
		return err
	}

	num := syscall.Signal(sig.Signum())

	log.Infof("sending signal %s (%d) to %s", sig.Signal, num, session.ID)

	if err := session.Cmd.Process.Signal(num); err != nil {
		return fmt.Errorf("failed to signal %s: %s", session.ID, err)
	}

	return nil
}

func (t *Toolbox) halt() error {
	session := t.config.Sessions[t.config.ID]
	log.Infof("stopping %s", session.ID)

	if err := t.kill(t.config.StopSignal); err != nil {
		return err
	}

	select {
	case <-session.exit:
		log.Infof("%s has stopped", session.ID)
		return nil
	case <-time.After(time.Second * 10): // TODO: honor -t flag from docker stop
	}

	log.Warnf("killing %s", session.ID)

	return session.Cmd.Process.Kill()
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
