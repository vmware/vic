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
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"testing"
)

func TestDefaultIP(t *testing.T) {
	ip := DefaultIP()
	if ip == "" {
		t.Error("failed to get a default IP address")
	}
}

type testRPC struct {
	cmd    string
	expect string
}

type mockChannelIn struct {
	t       *testing.T
	service *Service
	rpc     []*testRPC
	wg      sync.WaitGroup
	start   error
}

func (c *mockChannelIn) Start() error {
	return c.start
}

func (c *mockChannelIn) Stop() error {
	return nil
}

func (c *mockChannelIn) Receive() ([]byte, error) {
	if len(c.rpc) == 0 {
		if c.rpc != nil {
			// All test RPC requests have been consumed
			c.wg.Done()
			c.rpc = nil
		}
		return nil, io.EOF
	}

	return []byte(c.rpc[0].cmd), nil
}

func (c *mockChannelIn) Send(buf []byte) error {
	if buf == nil {
		return nil
	}

	expect := c.rpc[0].expect
	if string(buf) != expect {
		c.t.Errorf("expected %q reply for request %q, got: %q", expect, c.rpc[0].cmd, buf)
	}

	c.rpc = c.rpc[1:]

	return nil
}

// discard rpc out for now
type mockChannelOut struct {
	reply [][]byte
	start error
}

func (c *mockChannelOut) Start() error {
	return c.start
}

func (c *mockChannelOut) Stop() error {
	return nil
}

func (c *mockChannelOut) Receive() ([]byte, error) {
	if len(c.reply) == 0 {
		return nil, io.EOF
	}
	reply := c.reply[0]
	c.reply = c.reply[1:]
	return reply, nil
}

func (c *mockChannelOut) Send(buf []byte) error {
	if len(buf) == 0 {
		return io.ErrShortBuffer
	}
	return nil
}

func TestServiceRun(t *testing.T) {
	in := new(mockChannelIn)
	out := new(mockChannelOut)

	service := NewService(in, out)

	in.rpc = []*testRPC{
		{"reset", "OK ATR toolbox"},
		{"ping", "OK "},
		{"Set_Option synctime 0", "OK "},
		{"Capabilities_Register", "OK "},
		{"Set_Option broadcastIP 1", "OK "},
	}

	in.wg.Add(1)

	// replies to register capabilities
	for i := 0; i < len(capabilities); i++ {
		out.reply = append(out.reply, rpciOK)
	}

	out.reply = append(out.reply,
		rpciOK, // reply to IP broadcast
	)

	in.service = service

	in.t = t

	err := service.Start()
	if err != nil {
		t.Fatal(err)
	}

	in.wg.Wait()

	service.Stop()
	service.Wait()

	// verify we don't set delay > maxDelay
	for i := 0; i <= maxDelay+1; i++ {
		service.backoff()
	}

	if service.delay != maxDelay {
		t.Errorf("delay=%d", service.delay)
	}
}

func TestServiceErrors(t *testing.T) {
	Trace = true
	if !testing.Verbose() {
		// cover TraceChannel but discard output
		traceLog = ioutil.Discard
	}

	netInterfaceAddrs = func() ([]net.Addr, error) {
		return nil, io.EOF
	}

	in := new(mockChannelIn)
	out := new(mockChannelOut)

	service := NewService(in, out)

	service.RegisterHandler("Sorry", func([]byte) ([]byte, error) {
		return nil, errors.New("i am so sorry")
	})

	in.rpc = []*testRPC{
		{"Capabilities_Register", "OK "},
		{"Set_Option broadcastIP 1", "ERR "},
		{"NOPE", "Unknown Command"},
		{"Sorry", "ERR "},
	}

	in.wg.Add(1)

	// replies to register capabilities
	for i := 0; i < len(capabilities); i++ {
		out.reply = append(out.reply, rpciERR)
	}

	foo := []byte("foo")
	out.reply = append(
		out.reply,
		rpciERR,
		append(rpciOK, foo...),
		rpciERR,
	)

	in.service = service

	in.t = t

	err := service.Start()
	if err != nil {
		t.Fatal(err)
	}

	in.wg.Wait()

	// Done serving RPCs, test ChannelOut errors
	reply, err := service.out.Request(rpciOK)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(reply, foo) {
		t.Errorf("reply=%q", foo)
	}

	_, err = service.out.Request(rpciOK)
	if err == nil {
		t.Error("expected error")
	}

	_, err = service.out.Request(nil)
	if err == nil {
		t.Error("expected error")
	}

	service.Stop()
	service.Wait()

	// cover service start error paths
	start := errors.New("fail")

	in.start = start
	err = service.Start()
	if err != start {
		t.Error("expected error")
	}

	in.start = nil
	out.start = start
	err = service.Start()
	if err != start {
		t.Error("expected error")
	}
}

var (
	testESX = flag.Bool("toolbox.testesx", false, "Test toolbox service against ESX (vmtoolsd must not be running)")
	testPID = flag.Int("toolbox.testpid", 0, "PID to return from toolbox start command")
)

func TestServiceRunESX(t *testing.T) {
	if *testESX == false {
		t.SkipNow()
	}

	Trace = testing.Verbose()

	var wg sync.WaitGroup

	in := NewBackdoorChannelIn()
	out := NewBackdoorChannelOut()

	service := NewService(in, out)

	// assert that reset, ping, Set_Option and Capabilities_Register are called at least once
	for name, handler := range service.handlers {
		n := name
		h := handler
		wg.Add(1)

		service.handlers[name] = func(b []byte) ([]byte, error) {
			defer wg.Done()

			service.handlers[n] = h // reset

			return h(b)
		}
	}

	vix := RegisterVixRelayedCommandHandler(service)

	if *testPID != 0 {
		wg.Add(1)
		vix.ProcessStartCommand = func(r *VixMsgStartProgramRequest) (int, error) {
			defer wg.Done()

			if r.ProgramPath != "/bin/date" {
				t.Errorf("ProgramPath=%q", r.ProgramPath)
			}

			return *testPID, nil
		}
	}

	wg.Add(1)
	service.PrimaryIP = func() string {
		defer wg.Done()
		log.Print("broadcasting IP")
		return DefaultIP()
	}

	err := service.Start()
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	service.Stop()
	service.Wait()
}
