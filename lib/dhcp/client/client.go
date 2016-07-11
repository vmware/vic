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

package client

import (
	"fmt"
	"net"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/d2g/dhcp4"
	"github.com/d2g/dhcp4client"
	"github.com/vmware/vic/lib/dhcp"
	"github.com/vmware/vic/pkg/ip"
)

// Client represents a DHCP client
type Client interface {
	// SetTimeout sets the timeout for a subsequent DHCP request
	SetTimeout(t time.Duration)

	// Request sends a full DHCP request, resulting in a DHCP lease.
	// On a successful lease, returns a DHCP acknowledgment packet
	Request(ID) (*dhcp.Packet, error)

	// Renew renews an existing DHCP lease. Returns a new acknowledgment
	// packet on success.
	Renew(ID, *dhcp.Packet) (*dhcp.Packet, error)

	// Release releases an existing DHCP lease.
	Release(*dhcp.Packet) error
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

func isCompletePacket(p *dhcp.Packet) bool {
	complete := !ip.IsUnspecifiedIP(p.Gateway()) &&
		!ip.IsUnspecifiedIP(p.YourIP()) &&
		!ip.IsUnspecifiedIP(p.ServerIP())

	if !complete {
		return false
	}

	ones, bits := p.SubnetMask().Size()
	if ones == 0 || bits == 0 {
		return false
	}

	if p.LeaseTime().Seconds() == 0 {
		return false
	}

	return true
}

func (c *client) appendOptions(p *dhcp4.Packet, id ID) (*dhcp4.Packet, error) {
	p.AddOption(
		dhcp4.OptionParameterRequestList,
		[]byte{
			byte(dhcp4.OptionSubnetMask),
			byte(dhcp4.OptionRouter),
			byte(dhcp4.OptionDomainNameServer),
		},
	)
	b, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}

	p.AddOption(dhcp4.OptionClientIdentifier, b)
	return p, nil
}

func (c *client) discoverPacket(id ID, cl *dhcp4client.Client) (*dhcp4.Packet, error) {
	dp := cl.DiscoverPacket()
	return c.appendOptions(&dp, id)
}

func (c *client) requestPacket(id ID, cl *dhcp4client.Client, op *dhcp4.Packet) (*dhcp4.Packet, error) {
	rp := cl.RequestPacket(op)
	return c.appendOptions(&rp, id)
}

func (c *client) request(id ID, cl *dhcp4client.Client) (bool, *dhcp.Packet, error) {
	dp, err := c.discoverPacket(id, cl)
	if err != nil {
		return false, nil, err
	}

	dp.PadToMinSize()
	if err = cl.SendPacket(*dp); err != nil {
		return false, nil, err
	}

	var op dhcp4.Packet
	for {
		op, err = cl.GetOffer(dp)
		if err != nil {
			return false, nil, err
		}

		if isCompletePacket(dhcp.NewPacket([]byte(op))) {
			break
		}
	}

	rp, err := c.requestPacket(id, cl, &op)
	if err != nil {
		return false, nil, err
	}

	rp.PadToMinSize()
	if err = cl.SendPacket(*rp); err != nil {
		return false, nil, err
	}

	ack, err := cl.GetAcknowledgement(rp)
	if err != nil {
		return false, nil, err
	}

	opts := ack.ParseOptions()
	if dhcp4.MessageType(opts[dhcp4.OptionDHCPMessageType][0]) == dhcp4.NAK {
		return false, nil, fmt.Errorf("Got NAK from DHCP server")
	}

	return true, dhcp.NewPacket([]byte(ack)), nil
}

func (c *client) Request(id ID) (*dhcp.Packet, error) {
	log.Debugf("id: %+v", id)
	// send the request over a raw socket
	raw, err := dhcp4client.NewPacketSock(id.IfIndex)
	if err != nil {
		return nil, err
	}

	rawc, err := dhcp4client.New(dhcp4client.Connection(raw), dhcp4client.Timeout(c.timeout), dhcp4client.HardwareAddr(id.HardwareAddr))
	if err != nil {
		return nil, err
	}
	defer rawc.Close()

	success := false
	var p *dhcp.Packet
	err = withRetry(func() error {
		var err error
		success, p, err = c.request(id, rawc)
		return err
	})

	if err != nil {
		return nil, err
	}

	if !success {
		return nil, fmt.Errorf("failed dhcp request")
	}

	log.Debugf("%+v", p)
	return p, nil
}

func (c *client) newClient(ack *dhcp.Packet) (*dhcp4client.Client, error) {
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

func (c *client) renew(id ID, ack dhcp4.Packet, cl *dhcp4client.Client) (dhcp4.Packet, error) {
	rp := cl.RenewalRequestPacket(&ack)
	_, err := c.appendOptions(&rp, id)
	if err != nil {
		return nil, err
	}

	rp.PadToMinSize()
	if err = cl.SendPacket(rp); err != nil {
		return nil, err
	}

	newack, err := cl.GetAcknowledgement(&rp)
	if err != nil {
		return nil, err
	}

	opts := newack.ParseOptions()
	if dhcp4.MessageType(opts[dhcp4.OptionDHCPMessageType][0]) == dhcp4.NAK {
		return nil, fmt.Errorf("received NAK from DHCP server")
	}

	return newack, nil
}

func (c *client) Renew(id ID, ack *dhcp.Packet) (*dhcp.Packet, error) {
	c.Lock()
	defer c.Unlock()

	log.Debugf("renewing IP %s", ack.YourIP())

	cl, err := c.newClient(ack)
	if err != nil {
		return nil, err
	}
	defer cl.Close()

	p := dhcp4.Packet(ack.Packet)
	var newack dhcp4.Packet
	err = withRetry(func() error {
		var err error
		newack, err = c.renew(id, p, cl)
		return err
	})

	if err != nil {
		return nil, err
	}

	return dhcp.NewPacket([]byte(newack)), nil
}

func (c *client) Release(ack *dhcp.Packet) error {
	c.Lock()
	defer c.Unlock()

	log.Debugf("releasing IP %s", ack.YourIP())

	cl, err := c.newClient(ack)
	if err != nil {
		return err
	}
	defer cl.Close()

	return withRetry(func() error {
		return cl.Release(dhcp4.Packet(ack.Packet))
	})
}
