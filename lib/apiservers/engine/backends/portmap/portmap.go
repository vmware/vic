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

package portmap

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/libnetwork/iptables"
)

type Operation int

const (
	Map Operation = iota
	Unmap
)

type PortMapper interface {
	MapPort(op Operation, ip net.IP, port int, proto string, destIP string, destPort int, srcIface, destIface string) error
}

type bindKey struct {
	ip   string
	port int
}

type portMapper struct {
	sync.Mutex

	bindings map[bindKey]interface{}
}

func NewPortMapper() PortMapper {
	return &portMapper{bindings: make(map[bindKey]interface{})}
}

func (p *portMapper) isPortAvailable(proto string, ip net.IP, port int) bool {
	addr := ""
	if ip != nil && !ip.IsUnspecified() {
		addr = ip.String()
	}

	if _, ok := p.bindings[bindKey{addr, port}]; ok {
		return false
	}

	c, err := net.Dial(proto, net.JoinHostPort(addr, strconv.Itoa(port)))
	defer func() {
		if c != nil {
			c.Close()
		}
	}()

	if err != nil {
		return true
	}

	return false
}

func (p *portMapper) MapPort(op Operation, ip net.IP, port int, proto string, destIP string, destPort int, srcIface, destIface string) error {
	p.Lock()
	defer p.Unlock()

	if port <= 0 {
		return fmt.Errorf("source port must be specified")
	}

	if destPort <= 0 {
		log.Infof("destination port not specified, using source port %d", port)
		destPort = port
	}

	if destIP == "" {
		return fmt.Errorf("destination IP is not specified")
	}

	var action iptables.Action
	switch op {
	case Map:
		if !p.isPortAvailable(proto, ip, port) {
			return fmt.Errorf("port %d not available", port)
		}
		action = iptables.Append

	case Unmap:
		action = iptables.Delete

	default:
		return fmt.Errorf("invalid port mapping operation %d", op)
	}

	return p.forward(action, ip, port, proto, destIP, destPort, srcIface, destIface)
}

// adapted from https://github.com/docker/libnetwork/blob/master/iptables/iptables.go
func (p *portMapper) forward(action iptables.Action, ip net.IP, port int, proto, destAddr string, destPort int, srcIface, destIface string) error {
	daddr := ip.String()
	if ip == nil || ip.IsUnspecified() {
		// iptables interprets "0.0.0.0" as "0.0.0.0/32", whereas we
		// want "0.0.0.0/0". "0/0" is correctly interpreted as "any
		// value" by both iptables and ip6tables.
		daddr = "0/0"
	}
	args := []string{"-t", string(iptables.Nat), string(action), "VIC",
		"-i", srcIface,
		"-p", proto,
		"-d", daddr,
		"--dport", strconv.Itoa(port),
		"-j", "DNAT",
		"--to-destination", net.JoinHostPort(destAddr, strconv.Itoa(destPort))}
	if output, err := iptables.Raw(args...); err != nil {
		return err
	} else if len(output) != 0 {
		return iptables.ChainError{Chain: "FORWARD", Output: output}
	}

	ipStr := ""
	if ip != nil && !ip.IsUnspecified() {
		ipStr = ip.String()
	}

	switch action {
	case iptables.Append:
		p.bindings[bindKey{ipStr, port}] = nil

	case iptables.Delete:
		delete(p.bindings, bindKey{ipStr, port})
	}

	if output, err := iptables.Raw("-t", string(iptables.Filter), string(action), "VIC",
		"-i", srcIface,
		"-o", destIface,
		"-p", proto,
		"-d", destAddr,
		"--dport", strconv.Itoa(destPort),
		"-j", "ACCEPT"); err != nil {
		return err
	} else if len(output) != 0 {
		return iptables.ChainError{Chain: "FORWARD", Output: output}
	}

	if output, err := iptables.Raw("-t", string(iptables.Nat), string(action), "POSTROUTING",
		"-p", proto,
		"-d", destAddr,
		"--dport", strconv.Itoa(destPort),
		"-j", "MASQUERADE"); err != nil {
		return err
	} else if len(output) != 0 {
		return iptables.ChainError{Chain: "FORWARD", Output: output}
	}

	return nil
}
