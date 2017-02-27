// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package nfs

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func listenAndServe(t *testing.T, port int) (*net.TCPListener, *sync.WaitGroup, error) {

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		l.Accept()
		t.Logf("Accepted conn")
		l.Accept()
		t.Logf("Accepted conn")
		wg.Done()
	}()

	return l, wg, nil
}

// test we can bind without colliding
func TestDialService(t *testing.T) {
	listener, wg, err := listenAndServe(t, 6666)
	if err != nil {
		t.Logf("error starting listener: %s", err.Error())
		t.Fail()
		return
	}
	defer listener.Close()

	_, err = dialService("127.0.0.1", 6666)
	if err != nil {
		t.Logf("error dialing: %s", err.Error())
		t.FailNow()
	}

	_, err = dialService("127.0.0.1", 6666)
	if err != nil {
		t.Logf("error dialing: %s", err.Error())
		t.FailNow()
	}

	wg.Wait()
}
