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
	"io"
	"testing"
	"time"
)

func TestDefaultIP(t *testing.T) {
	ip := DefaultIP()
	if ip == "" {
		t.Error("failed to get a default IP address")
	}
	t.Logf("DefaultIP=%s", ip)
}

type testRPC struct {
	cmd    string
	expect string
}

type mockChannelIn struct {
	t       *testing.T
	service *Service
	rpc     []*testRPC
}

func (c *mockChannelIn) Start() error {
	return nil
}

func (c *mockChannelIn) Stop() error {
	return nil
}

func (c *mockChannelIn) Receive() ([]byte, error) {
	if len(c.rpc) == 0 {
		// Stop the service after all test RPC requests have been consumed
		defer c.service.Stop()
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
		c.t.Errorf("expected '%s' reply for request '%s', got: '%s'", expect, c.rpc[0].cmd, string(buf))
	}

	c.rpc = c.rpc[1:]

	return nil
}

// discard rpc out for now
type mockChannelOut struct{}

func (c *mockChannelOut) Start() error {
	return nil
}

func (c *mockChannelOut) Stop() error {
	return nil
}

func (c *mockChannelOut) Receive() ([]byte, error) {
	panic("receive on out channel")
}

func (c *mockChannelOut) Send(buf []byte) error {
	return nil
}

func TestServiceRun(t *testing.T) {
	Trace = testing.Verbose()

	in := new(mockChannelIn)
	out := new(mockChannelOut)

	service := NewService(in, out)
	service.Interval = time.Millisecond

	in.rpc = []*testRPC{
		{"reset", "OK ATR toolbox"},
		{"ping", "OK "},
		{"Set_Option synctime 0", "OK "},
		{"NOPE", "Unknown Command"},
		{"Set_Option broadcastIP 1", "OK "},
	}

	in.service = service

	in.t = t

	err := service.Start()
	if err != nil {
		t.Fatal(err)
	}

	service.Wait()
}
