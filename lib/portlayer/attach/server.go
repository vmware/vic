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
	"context"
	"fmt"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

// AttachServer waits for TCP client connections on serialOverLANPort, then
// once connected, attempts to negotiate an SSH connection to the attached
// client.  The client is the ssh server.
type Server struct {
	port int
	ip   string
	l    *net.TCPListener

	connServer *Connector
}

func NewAttachServer(ip string, port int) *Server {
	defer trace.End(trace.Begin(""))

	return &Server{ip: "localhost", port: port}
}

// Start starts the TCP listener.
func (n *Server) Start(debug bool) error {
	defer trace.End(trace.Begin(""))

	log.Infof("Attach server listening on %s:%d", n.ip, n.port)

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", n.ip, n.port))
	if err != nil {
		return fmt.Errorf("Attach server error %s:%d: %s", n.ip, n.port, errors.ErrorStack(err))
	}

	n.l, err = net.ListenTCP("tcp", addr)
	if err != nil {
		err = fmt.Errorf("Attach server error %s: %s", addr, errors.ErrorStack(err))
		log.Errorf("%s", err)
		return err
	}

	// starts serving requests immediately
	n.connServer = NewConnector(n.l, debug)

	return nil
}

func (n *Server) Stop() error {
	defer trace.End(trace.Begin(""))

	err := n.l.Close()
	n.connServer.Stop()
	return err
}

func (n *Server) Addr() string {
	defer trace.End(trace.Begin(""))

	return n.l.Addr().String()
}

// Get returns the session interface for the given container.  If the container
// cannot be found, this call will wait for the given timeout.
// id is ID of the container.
func (n *Server) Get(ctx context.Context, id string, timeout time.Duration) (SessionInteraction, error) {
	defer trace.End(trace.Begin(id))

	return n.connServer.Get(ctx, id, timeout)
}

func (n *Server) Remove(id string) error {
	defer trace.End(trace.Begin(id))

	return n.connServer.Remove(id)
}
