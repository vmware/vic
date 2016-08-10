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
	"errors"
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
	stop   chan struct{}
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

	t.stop = make(chan struct{})
	on := make(chan struct{})

	t.Service.PowerCommand.PowerOn.Handler = func() error {
		log.Info("toolbox: service is ready (power on event received)")
		close(on)
		return nil
	}

	err := t.Service.Start()
	if err != nil {
		return err
	}

	// Wait for the vmx to send the OS_PowerOn message,
	// at which point it will be ready to service vix command requests.
	log.Info("toolbox: waiting for initialization")

	select {
	case <-on:
	case <-time.After(time.Second):
		log.Warn("toolbox: timeout waiting for power on event")
	}

	return nil
}

// Stop implementation of the tether.Extension interface
func (t *Toolbox) Stop() error {
	t.Service.Stop()

	t.Service.Wait()

	close(t.stop)

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

	vix := t.Service.VixCommand
	vix.Authenticate = t.containerAuthenticate
	vix.ProcessStartCommand = t.containerStartCommand

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

func (t *Toolbox) containerAuthenticate(_ toolbox.VixCommandRequestHeader, data []byte) error {
	var c toolbox.VixUserCredentialNamePassword
	if err := c.UnmarshalBinary(data); err != nil {
		return err
	}
	// no authentication yet, just using container ID as a sanity check for now
	if c.Name != t.config.ID {
		return errors.New("failed to verify container ID")
	}

	return nil
}

func (t *Toolbox) containerStartCommand(r *toolbox.VixMsgStartProgramRequest) (int, error) {
	switch r.ProgramPath {
	case "kill":
		return -1, t.kill(r.Arguments)
	default:
		return -1, fmt.Errorf("unknown command %q", r.ProgramPath)
	}
}

func (t *Toolbox) halt() error {
	session := t.config.Sessions[t.config.ID]
	log.Infof("stopping %s", session.ID)

	if err := t.kill(t.config.StopSignal); err != nil {
		return err
	}

	// Killing the executor session in the container VM will stop the tether and its extensions.
	// If that doesn't happen within the timeout, send a SIGKILL.
	select {
	case <-t.stop:
		log.Infof("%s has stopped", session.ID)
		return nil
	case <-time.After(time.Second * 10):
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
