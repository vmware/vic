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

package scp

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/testdata"
)

type Mode string

const (
	SOURCE = Mode("source")
	DEST   = Mode("dest")
)

var (
	TestUser       = "user"
	TestPassword   = "password"
	TestPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAw899jLRFFRErKsERYk6TWM/nfm7Gbp5WeEGP5Y6tIeuRCINm
KyVCac1c8TfaARKKOFbLHKhkp1ajIsidQ2tXL3VZQFVY2nD1QJ187t9nRHK2D9F5
AacrCq4EnNF6hHArPLAgG40fRfeqBQ1QDjQYsdIMhJRoYo6qQW2npqpeZqZ/tMxE
PnNjDcoBFF7JdSJi8j9/hznolXWmz+RXj1G8tIF9dlWd2lEiodLdU+U0/WVOs2Fi
Ycm1a59mYtu2XV8MfYpdfNogyH8VypUPK1cY0er0IfacEo0OxcO8cDv0WvNsZ+Pg
3pn8+Hjlcr8I1AXC1yM5MnA1AivuV00Ucnb2YwIDAQABAoIBAD9gfSZ5gpKbB/nC
m7nR7OcmA6tsd1V+ckZiEg0e7PK5qCu1O/BjEufjzF1W0nzeX1Z04TDZYBq5c/vi
KuSTbZiyxryH40ZwoTDUyIcYT/hbmInuJtheHxRJ1rxbIOiU1anC5+GC/8hJrQIN
mQe/3O5RKjQPROoBeyHKOCU9p2hZkrRw72ofJUTnuJVJTPtG/tXMthjS7dcK1AMh
Hqo50Qjth9iaNEMdxlRleFP7kkxFfgbREgDtUxbaW5ehuz0QC8Qdk7JLAfhc0s9T
6pNasET/WPHohpA3y3lQnNqzP8qgc+ZobLeg4cuRXC5s4HAP1X0ZJtW2in4OWD2l
vumfjgECgYEA4SlnEYcLjOBJ7jsjWyWgIijcKOI3ErKT5DSBYagMNBtAr7caYjYc
bDKsOLgSsPKvh/2e2VwEapNiSPUwDoDy90S+Pr4LlU8DAO+GThX700SYsoQNpEfU
oGucDVDFk8gBvndOwWX/4QRBBcnO+TcMJwIcqdak1zO3ifzFx+JHdaMCgYEA3qD6
dvzwODkGpSgXDYNbIAwDLUE3Fa5ymGQx8LKZJ2/mLtmEKlfb928Oi5W9rpLYy7oR
Nn9JLaLtGDDUH5iznaTR0ZqqMN7skjkQRqr+DtxZ9PERwRZ7qDhQQb20hlmv/tyX
yeExk9zQtsYhjl7wXpvWMI+f6rZzvJmAuhsXCEECgYATpuFwDjFb4leRi2fSlL4d
PSO6DcRwxVVTHaINO/WUtqw2qeyLld11NBcD/EzlVMktPV2X6wgXpTV22K+RFIAg
RMe2AjBQn8zLUByQxCpujhYlvpDSPdK5DatZHiugclx0m0UsbBKhORXTw4FlDwDo
hq6pxCou/jyOtpkskPtbsQKBgFV1S8/DFl8unLtnITpBswghNFL51rBO75RJ2dXA
aQP3c0+GlbI/WaOokNfKGi7aFbhWa2cVAz0ubn67t6GNV11rOFOSYEQ0PnF+0B2g
Y7fGpA1fQGZzP/J582zY6mQsJ1/Yw7dt5z8QI1oVwinJjdFzVov7hfJuKQ07i2tl
HLlBAoGASUFpIwz+7cwCnMHXwBsZwujHOtGm0AM2dpsPO0eJJwDNxjbIGcnDNwHB
R5KSZ6i7FRCsXClB3SLfRH9AqT9UgpPE8FANUhJMvaWL1CAKYvayum7fNGvqGyfY
pZo4v2rNNUvnNPS3g7x0ss3hTxBseY8LoEt5G9rhziitSEn1LkQ=
-----END RSA PRIVATE KEY-----
`
	testPublicKey []byte
)

func init() {
	signer, err := ssh.ParsePrivateKey([]byte(TestPrivateKey))
	if err != nil {
		panic(err)
	}
	testPublicKey = signer.PublicKey().Marshal()
}

var TestCertChecker = ssh.CertChecker{
	IsAuthority: func(k ssh.PublicKey) bool {
		return true
	},
	IsRevoked: func(c *ssh.Certificate) bool {
		return false
	},
	UserKeyFallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		if conn.User() == TestUser && bytes.Equal(key.Marshal(), testPublicKey) {
			return nil, nil
		}
		return nil, errors.New("unauthorized key")
	},
}

type SessionHandler interface {
	SetChannel(channel *ssh.Channel)

	Source(dirent string) (ok bool, payload []byte)

	Destination(dirent string) (ok bool, payload []byte)

	// Retrieve closure for any pending work - this is necessary as data cannot be returned
	// via ssh before request replys are sent so exec, et al, must be async.
	// This will be called by the tether in a goroutine after replys have been processed
	GetPendingWork() func()

	// Called after GetPendingWork closure has been invoked
	ClearPendingWork()
}

// StartSSHExecServer runs an ssh server that accepts a single connection and handles only "exec" requests.
// The handler callback is given the command string and stderr io.Writer.  The handler return
// value is propagated to the client via "exit-status".
func StartSSHExecServer(port *int, sh SessionHandler) (*sync.WaitGroup, error) {
	config := &ssh.ServerConfig{
		PublicKeyCallback: TestCertChecker.Authenticate,
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == TestUser && string(pass) == TestPassword {
				return nil, nil
			}
			return nil, errors.New("unauthorized user")
		},
	}

	signer, err := ssh.ParsePrivateKey(testdata.PEMBytes["dsa"])
	if err != nil {
		panic(err)
	}
	config.AddHostKey(signer)

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		panic(err)
	}
	*port = listener.Addr().(*net.TCPAddr).Port

	var wg sync.WaitGroup
	wg.Add(1)

	// service 1 request
	go func() {
		defer listener.Close()
		defer wg.Done()

		tcpConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			return
		}
		// Before use, a handshake must be performed on the incoming net.Conn.
		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			log.Printf("Failed to handshake (%s)", err)
			return
		}

		log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())
		// Discard all global out-of-band Requests
		go ssh.DiscardRequests(reqs)
		// Accept all channels
		handleChannels(sshConn, chans, sh)

	}()

	return &wg, nil
}

func handleChannels(sshCon *ssh.ServerConn, chans <-chan ssh.NewChannel, sh SessionHandler) {
	// Service the incoming Channel channel in go routine
	for newChannel := range chans {
		handleChannel(newChannel, sh)
		sshCon.Close()
	}
}

func handleChannel(newChannel ssh.NewChannel, sh SessionHandler) {
	// At this point, we have the opportunity to reject the client's
	// request for another logical connection
	channel, requests, err := newChannel.Accept()
	if err != nil {
		log.Printf("Could not accept channel (%s)", err)
		return
	}

	// Prepare teardown function
	close := func() {
		channel.Close()
		log.Printf("Session closed")
	}
	defer close()

	sh.SetChannel(&channel)
	var payload []byte
	ok := true

	for req := range requests {

		if len(req.Payload) == 0 {
			req.Reply(true, nil)
		}

		path := string(req.Payload)
		typ := Mode(req.Type)

		switch typ {
		// in reality, this would be an scp with the appropriate args,
		// but for testing this is fine.
		case SOURCE:
			ok, payload = sh.Source(path)
		case DEST:
			ok, payload = sh.Destination(path)
		default:
			return
		}

		// make sure that errors get send back if we failed
		if req.WantReply {
			log.Printf("Wants reply %s", string(payload))
			req.Reply(ok, payload)
		}

		// run any pending work now that a reply has been sent
		pendingFn := sh.GetPendingWork()
		if pendingFn != nil {
			pendingFn()
		}
	}
}
