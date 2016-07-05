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

package dhcp

import (
	"fmt"
	"net"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/d2g/dhcp4"
	"github.com/d2g/dhcp4client"
)

// Client represents a DHCP client
type Client interface {
	// SetTimeout sets the timeout for a subsequent DHCP request
	SetTimeout(t time.Duration)

	// Request sends a full DHCP request, resulting in a DHCP lease.
	// On a successful lease, returns a DHCP acknowledgment packet
	Request(int, net.HardwareAddr) (*Packet, error)

	// Renew renews an existing DHCP lease. Returns a new acknowledgment
	// packet on success.
	Renew(*Packet) (*Packet, error)

	// Release releases an existing DHCP lease.
	Release(*Packet) error
}

type client struct {
	sync.Mutex

	timeout time.Duration
}

// The default timeout for the client
const defaultTimeout = 10 * time.Second

// NewClient creates a new DHCP client. Note the returned object is not thread-safe.
func NewClient() (Client, error) {
	return &client{timeout: defaultTimeout}, nil
}

func (c *client) SetTimeout(t time.Duration) {
	c.timeout = t
}

func withRetry(op func() error) error {
	for {
		if err := op(); err != nil {
			if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EAGAIN {
				return err
			}
		} else {
			return nil
		}
	}
}

func (c *client) Request(linkIndex int, hw net.HardwareAddr) (*Packet, error) {
	// send the request over a raw socket
	raw, err := dhcp4client.NewPacketSock(linkIndex)
	if err != nil {
		return nil, err
	}

	rawc, err := dhcp4client.New(dhcp4client.Connection(raw), dhcp4client.Timeout(c.timeout), dhcp4client.HardwareAddr(hw))
	if err != nil {
		return nil, err
	}
	defer rawc.Close()

	success := false
	var p dhcp4.Packet
	err = withRetry(func() error {
		var err error
		success, p, err = rawc.Request()
		return err
	})

	if err != nil {
		return nil, err
	}

	if !success {
		return nil, fmt.Errorf("failed dhcp request")
	}

	return &Packet{
		packet:  p,
		options: p.ParseOptions(),
	}, nil
}

func (c *client) newClient(ack *Packet) (*dhcp4client.Client, error) {
	conn, err := dhcp4client.NewInetSock(dhcp4client.SetRemoteAddr(net.UDPAddr{IP: ack.ServerIP(), Port: 67}))
	if err != nil {
		return nil, err
	}

	cl, err := dhcp4client.New(dhcp4client.Connection(conn), dhcp4client.Timeout(c.timeout))
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (c *client) Renew(ack *Packet) (*Packet, error) {
	c.Lock()
	defer c.Unlock()

	log.Debugf("renewing IP %s", ack.YourIP())

	cl, err := c.newClient(ack)
	if err != nil {
		return nil, err
	}
	defer cl.Close()

	success := false
	var p dhcp4.Packet
	err = withRetry(func() error {
		var err error
		success, p, err = cl.Renew(dhcp4.Packet(ack.packet))
		return err
	})

	if err != nil {
		return nil, err
	}

	if !success {
		return nil, fmt.Errorf("failed dhcp request")
	}

	return &Packet{
		packet:  p,
		options: p.ParseOptions(),
	}, nil
}

func (c *client) Release(ack *Packet) error {
	c.Lock()
	defer c.Unlock()

	log.Debugf("releasing IP %s", ack.YourIP())

	cl, err := c.newClient(ack)
	if err != nil {
		return err
	}
	defer cl.Close()

	return withRetry(func() error {
		return cl.Release(dhcp4.Packet(ack.packet))
	})
}
