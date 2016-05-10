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
	"net"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/testdata"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/pkg/serial"
)

// Start the server, make 200 client connections, test they connect, then Stop.
func TestAttachStartStop(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	s := NewAttachServer("", -1)

	wg := &sync.WaitGroup{}

	dial := func() {
		wg.Add(1)
		c, err := net.Dial("tcp", s.l.Addr().String())
		assert.NoError(t, err)
		assert.NotNil(t, c)

		buf := make([]byte, 1)
		c.Read(buf)
		if !assert.Error(t, serial.HandshakeServer(context.Background(), c)) {
			return
		}

		if !assert.NoError(t, serial.HandshakeServer(context.Background(), c)) {
			return
		}

		wg.Done()
	}

	assert.NoError(t, s.Start())

	for i := 0; i < 200; i++ {
		go dial()
	}

	done := make(chan bool)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fail()
	}
	assert.NoError(t, s.Stop())

	_, err := net.Dial("tcp", s.l.Addr().String())
	assert.Error(t, err)
}

func TestAttachSshSession(t *testing.T) {
	log.SetLevel(log.InfoLevel)

	s := NewAttachServer("", -1)
	assert.NoError(t, s.Start())
	defer s.Stop()

	expectedID := "foo"

	// This should block until the ssh server returns its container ID
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		_, err := s.connServer.Get(context.Background(), expectedID, 5*time.Second)
		if !assert.NoError(t, err) {
			return
		}
	}()

	// Dial the attach server.  This is a TCP client
	networkClientCon, err := net.Dial("tcp", s.l.Addr().String())
	if !assert.NoError(t, err) {
		return
	}

	if !assert.NoError(t, serial.HandshakeServer(context.Background(), networkClientCon)) {
		return
	}

	containerConfig := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	signer, err := ssh.ParsePrivateKey(testdata.PEMBytes["dsa"])
	if !assert.NoError(t, err) {
		return
	}
	containerConfig.AddHostKey(signer)

	// create the SSH server on the client.  The attach server will ssh connect to this.
	sshConn, chans, reqs, err := ssh.NewServerConn(networkClientCon, containerConfig)
	if !assert.NoError(t, err) {
		return
	}
	defer sshConn.Close()

	// Service the incoming Channel channel.
	go func() {
		wg.Add(1)
		defer wg.Done()
		for req := range reqs {
			if req.Type == requestContainerIDs {
				req.Reply(true, ssh.Marshal(struct {
					Strings []string
				}{[]string{expectedID}}))
				break
			}
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		for ch := range chans {
			assert.Equal(t, ch.ChannelType(), attachChannelType)
			_, _, _ = ch.Accept()
			break
		}
	}()

	wg.Wait()
}
