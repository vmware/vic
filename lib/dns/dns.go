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

package dns

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/trace"

	mdns "github.com/miekg/dns"
)

const (
	DefaultPort      = 53
	DefaultTTL       = 600 * time.Second
	DefaultCacheSize = 1024
)

var (
	options = ServerOptions{}
)

// ServerOptions represents the server options
type ServerOptions struct {
	IP   string
	Port int

	Nameservers flagMultipleVar

	Timeout time.Duration

	TTL       time.Duration
	CacheSize int

	Debug bool

	Profiling string
}

// Server represents udp/tcp server and clients
type Server struct {
	ServerOptions

	// used for serving dns
	udpserver *mdns.Server
	udpconn   *net.UDPConn

	tcpserver *mdns.Server
	tcplisten *net.TCPListener

	// used for forwarding queries
	udpclient *mdns.Client
	tcpclient *mdns.Client

	// used for speeding up external lookups
	cache *Cache
	wg    *sync.WaitGroup
}

type flagMultipleVar []string

func (i *flagMultipleVar) String() string {
	return fmt.Sprint(*i)
}

func (i *flagMultipleVar) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// NewServer returns a new Server
func NewServer(options ServerOptions) *Server {
	var err error

	// Default TTL
	if options.TTL == 0 {
		options.TTL = DefaultTTL
	}

	// Default port
	if options.Port == 0 {
		options.Port = DefaultPort
	}

	// Default nameservers
	if len(options.Nameservers) == 0 {
		options.Nameservers = resolvconf()
	}

	// Default cache size
	if options.CacheSize == 0 {
		options.CacheSize = DefaultCacheSize
	}

	server := &Server{
		ServerOptions: options,
		cache:         NewCache(CacheOptions{options.CacheSize, options.TTL}),
		wg:            new(sync.WaitGroup),
	}

	udpaddr := &net.UDPAddr{
		IP:   net.ParseIP(server.IP),
		Port: server.Port,
	}

	server.udpconn, err = net.ListenUDP("udp", udpaddr)
	if err != nil {
		log.Errorf("ListenUDP failed %s", err)
		return nil
	}

	tcpaddr := &net.TCPAddr{
		IP:   net.ParseIP(server.IP),
		Port: server.Port,
	}

	server.tcplisten, err = net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Errorf("ListenTCP failed %s", err)
		return nil
	}

	server.udpclient = &mdns.Client{
		Net:            "udp",
		ReadTimeout:    server.Timeout,
		WriteTimeout:   server.Timeout,
		SingleInflight: true,
	}

	server.tcpclient = &mdns.Client{
		Net:            "tcp",
		ReadTimeout:    server.Timeout,
		WriteTimeout:   server.Timeout,
		SingleInflight: true,
	}

	return server
}

// Addr returns the ip:port of the server
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

func respServerFailure(w mdns.ResponseWriter, r *mdns.Msg) error {
	m := new(mdns.Msg)
	m.SetRcode(r, mdns.RcodeServerFailure)
	// Does not matter if this write fails
	return w.WriteMsg(m)
}

func respNotImplemented(w mdns.ResponseWriter, r *mdns.Msg) error {
	m := &mdns.Msg{}
	m.SetReply(r)
	m.Compress = true

	m.Authoritative = false
	m.RecursionDesired = false
	m.RecursionAvailable = false
	m.Rcode = mdns.RcodeNotImplemented

	if err := w.WriteMsg(m); err != nil {
		log.Errorf("Error writing RcodeNotImplemented response, %s", err)
		return err
	}
	return nil
}

func resolvconf() []string {
	var servers []string

	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// comment.
		if len(line) > 0 && (line[0] == ';' || line[0] == '#') {
			continue
		}
		f := strings.SplitN(line, " ", 2)

		if len(f) < 1 {
			continue
		}

		if f[0] == "nameserver" {
			if len(f) > 1 && len(servers) < 3 { // small, but the standard limit
				// One more check: make sure server name is
				// just an IP address.  Otherwise we need DNS
				// to look it up.
				if net.ParseIP(f[1]) != nil {
					servers = append(servers, f[1])
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil
	}
	return servers
}

// HandleForwarding forwards a request to the nameservers and returns the response
func (s *Server) HandleForwarding(w mdns.ResponseWriter, request *mdns.Msg) error {
	defer trace.End(trace.Begin(request.String()))

	var r *mdns.Msg
	var err error
	var try int

	if len(s.Nameservers) == 0 {
		log.Errorf("No nameservers defined, can not forward")

		return respServerFailure(w, request)
	}

	// Do we have it in the cache
	if m := s.cache.Get(request); m != nil {
		log.Debugf("Cache hit %#v", m)

		// Overwrite the ID with the request's ID
		m.Id = request.Id
		m.Compress = true
		m.Truncated = false

		if err := w.WriteMsg(m); err != nil {
			log.Errorf("Error writing response, %s", err)
			return err
		}
		return nil
	}

	// which protocol  are they talking
	tcp := false
	if _, ok := w.RemoteAddr().(*net.TCPAddr); ok {
		tcp = true
	}

	// Use request ID for "random" nameserver selection.
	nsid := int(request.Id) % len(s.Nameservers)

Redo:
	nameserver := s.Nameservers[nsid]
	if i := strings.Index(nameserver, ":"); i < 0 {
		nameserver += ":53"
	}

	if tcp {
		r, _, err = s.tcpclient.Exchange(request, nameserver)
	} else {
		r, _, err = s.udpclient.Exchange(request, nameserver)
	}
	if err != nil {
		// Seen an error, this can only mean, "server not reached", try again but only if we have not exausted our nameservers.
		if try < len(s.Nameservers) {
			try++
			nsid = (nsid + 1) % len(s.Nameservers)
			goto Redo
		}

		log.Errorf("Failure to forward request %q", err)
		return respServerFailure(w, request)
	}

	// We have a response so cache it
	s.cache.Add(r)

	r.Compress = true
	if err := w.WriteMsg(r); err != nil {
		log.Errorf("Error writing response, %s", err)
		return err
	}
	return nil
}

// HandleVIC returns a response to a container name/id request
func (s *Server) HandleVIC(question mdns.Question) []mdns.RR {
	defer trace.End(trace.Begin(question.String()))

	//TODO
	log.Warnf("NOT IMPLEMENTED")

	return nil
}

// ServeDNS implements the handler interface
func (s *Server) ServeDNS(w mdns.ResponseWriter, r *mdns.Msg) {
	defer trace.End(trace.Begin(r.String()))

	clientIP, _, err := net.SplitHostPort(w.RemoteAddr().String())
	if err != nil {
		log.Infof("Request from %s", clientIP)
	}

	if r == nil || len(r.Question) == 0 {
		return
	}

	// Reject multi-question query
	if len(r.Question) != 1 {
		log.Errorf("Rejected multi-question query")

		respServerFailure(w, r)
		return
	}
	q := r.Question[0]

	if q.Qclass != mdns.ClassINET {
		log.Errorf("Rejected non-inet query")

		respNotImplemented(w, r)
		return
	}

	// Reject any query
	if q.Qtype == mdns.TypeANY {
		log.Errorf("Rejected ANY query")

		respNotImplemented(w, r)
		return
	}

	// Start crafting reply msg
	m := &mdns.Msg{}
	m.SetReply(r)

	m.Authoritative = true
	m.RecursionAvailable = true
	m.Compress = true

	// VIC
	answer := s.HandleVIC(q)
	if answer != nil {
		m.Answer = append(m.Answer, answer...)
	}

	if answer == nil {
		s.HandleForwarding(w, r)
		return
	}

	// which protocol we are talking
	tcp := false
	if _, ok := w.LocalAddr().(*net.TCPAddr); ok {
		tcp = true
	}

	// 512 byte payload guarantees that DNS packets can be reassembled if fragmented in transit.
	bufsize := 512

	// With EDNS0 in use a larger payload size can be specified.
	if o := r.IsEdns0(); o != nil {
		bufsize = int(o.UDPSize())
	}

	// Make sure we are not smaller than 512
	if bufsize < 512 {
		bufsize = 512
	}

	// with TCP we can send up to 64K
	if tcp {
		bufsize = mdns.MaxMsgSize - 1
	}

	// trim the answer RRs one by one till the whole message fits within the reply size
	if m.Len() > bufsize {
		if tcp {
			m.Truncated = true
		}

		for m.Len() > bufsize {
			m.Answer = m.Answer[:len(m.Answer)-1]
		}
	}

	if err := w.WriteMsg(m); err != nil {
		log.Errorf("Error writing response, %s", err)
	}
	w.Close()
}

// Start starts the DNS server
func (s *Server) Start() {
	udpserver := &mdns.Server{
		Handler:    s,
		PacketConn: s.udpconn,
	}
	s.udpserver = udpserver

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		udpserver.ActivateAndServe()
		log.Debugf("UDP server exited")
	}()
	log.Infof("Ready for queries on udp://%s", s.Addr())

	tcpserver := &mdns.Server{Handler: s, Listener: s.tcplisten}
	s.tcpserver = tcpserver

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		tcpserver.ActivateAndServe()
		log.Debugf("TCP server exited")
	}()
	log.Infof("Ready for queries on tcp://%s", s.Addr())

}

// Stop stops the DNS server gracefully
func (s *Server) Stop() {
	if s.udpserver != nil {
		log.Debugf("Shutting down udpserver")
		s.udpserver.Shutdown()
	}
	s.udpconn = nil
	s.udpserver = nil

	if s.tcpserver != nil {
		log.Debugf("Shutting down tcpserver")
		s.tcpserver.Shutdown()
	}
	s.tcplisten = nil
	s.tcpserver = nil
}

// Wait block until wg returns
func (s *Server) Wait() {
	s.wg.Wait()
}
