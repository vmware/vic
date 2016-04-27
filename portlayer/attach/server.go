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

package attach

import (
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/errors"
)

const (
	// docker official ports are 2375 and 2376
	serialOverLANPort = 2377
)

// AttachServer waits for TCP client connections on serialOverLANPort, then
// once connected, attemps to negotiate an SSH connection to the attached
// client.  The client is the ssh server.
type Server struct {
	port int
	ip   string
	l    *net.TCPListener

	connServer *Connector
}

func NewAttachServer(ip string, port int) *Server {
	if port == 0 {
		port = serialOverLANPort
	}

	return &Server{ip: ip, port: port}
}

// Start starts the TCP listener.
func (n *Server) Start() error {

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", n.ip, n.port))

	n.l, err = net.ListenTCP("tcp", addr)
	if err != nil {
		err = fmt.Errorf("vmware exec driver unable listen for spawned container VMs on %s: %s", addr, errors.ErrorStack(err))
		log.Errorf("%s", err)
		return err
	}

	// starts serving requests immediately
	n.connServer = NewConnector(n.l)

	return nil
}

func (n *Server) Stop() error {
	err := n.l.Close()
	n.connServer.Stop()
	return err
}
