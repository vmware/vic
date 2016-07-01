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

package toolbox

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// Service receives and dispatches incoming RPC requests from the vmx
type Service struct {
	name     string
	in       Channel
	out      Channel
	handlers map[string]Handler
	stop     chan struct{}
	wg       *sync.WaitGroup

	Interval  time.Duration
	PrimaryIP func() string
}

// NewService initializes a Service instance
func NewService(rpcIn Channel, rpcOut Channel) *Service {
	s := &Service{
		name:     "toolbox", // Same name used by vmtoolsd
		in:       NewTraceChannel(rpcIn),
		out:      NewTraceChannel(rpcOut),
		handlers: make(map[string]Handler),
		wg:       new(sync.WaitGroup),
		stop:     make(chan struct{}, 1),

		Interval:  time.Second,
		PrimaryIP: DefaultIP,
	}

	s.RegisterHandler("reset", s.Reset)
	s.RegisterHandler("ping", s.Ping)
	s.RegisterHandler("Set_Option", s.SetOption)
	s.RegisterHandler("Capabilities_Register", s.CapabilitiesRegister)

	return s
}

// Start initializes the RPC channels and starts a goroutine to listen for incoming RPC requests
func (s *Service) Start() error {
	err := s.in.Start()
	if err != nil {
		return err
	}

	err = s.out.Start()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(s.Interval)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		for {
			select {
			case <-ticker.C:
				_ = s.in.Send(nil) // POKE

				request, _ := s.in.Receive()

				if len(request) > 0 {
					response := s.Dispatch(request)

					_ = s.in.Send(response)
				}
			case <-s.stop:
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

// Stop cancels the RPC listener routine created via Start
func (s *Service) Stop() {
	s.stop <- struct{}{}

	_ = s.in.Stop()
	_ = s.out.Stop()
}

// Wait blocks until Start returns
func (s *Service) Wait() {
	s.wg.Wait()
}

// Handler is given the raw argument portion of an RPC request and returns a response
type Handler func([]byte) ([]byte, error)

// RegisterHandler for the given RPC name
func (s *Service) RegisterHandler(name string, handler Handler) {
	s.handlers[name] = handler
}

// Dispatch an incoming RPC request to a Handler
func (s *Service) Dispatch(request []byte) []byte {
	msg := bytes.SplitN(request, []byte{' '}, 2)
	name := msg[0]

	// Trim NULL byte terminator
	name = bytes.TrimRight(name, "\x00")

	handler, ok := s.handlers[string(name)]

	if !ok {
		log.Printf("unknown command: '%s'\n", name)
		return []byte("Unknown Command")
	}

	var args []byte
	if len(msg) == 2 {
		args = msg[1]
	}

	response, err := handler(args)
	if err == nil {
		response = append([]byte("OK "), response...)
	} else {
		log.Printf("error calling %s: %s\n", name, err)
		response = append([]byte("ERR "), response...)
	}

	return response
}

// Reset is the default Handler for reset requests
func (s *Service) Reset([]byte) ([]byte, error) {
	return []byte("ATR " + s.name), nil
}

// Ping is the default Handler for ping requests
func (s *Service) Ping([]byte) ([]byte, error) {
	return nil, nil
}

// SetOption is the default Handler for Set_Option requests
func (s *Service) SetOption(args []byte) ([]byte, error) {
	opts := bytes.SplitN(args, []byte{' '}, 2)
	key := string(opts[0])
	val := string(opts[1])

	if Trace {
		fmt.Fprintf(os.Stderr, "set option '%s'='%s'\n", key, val)
	}

	switch key {
	case "broadcastIP": // TODO: const-ify
		if val == "1" {
			ip := s.PrimaryIP()
			msg := fmt.Sprintf("info-set guestinfo.ip %s", ip)
			return nil, s.out.Send([]byte(msg))
		}
	default:
		// TODO: handle other options...
	}

	return nil, nil
}

// DefaultIP is used by default when responding to a Set_Option broadcastIP request
// It can be overridden with the Service.PrimaryIP field
func DefaultIP() string {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
				if ip.IP.To4() != nil {
					return ip.IP.String()
				}
			}
		}
	}

	return ""
}

func (s *Service) CapabilitiesRegister([]byte) ([]byte, error) {
	// TODO: this is here just to make Fusion happy.  ESX doesn't seem to mind if we don't support this RPC
	return nil, nil
}
