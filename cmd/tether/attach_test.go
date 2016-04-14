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

package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net"
	"testing"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"

	"github.com/vmware/vic/metadata"
)

/////////////////////////////////////////////////////////////////////////////////////
// TestAttachConfig sets up the config for attach testing - the grep will echo anything
// sent and adds colour which is useful for tty testing
//
type TestAttachConfig struct{}

func (c *TestAttachConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestAttachConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "attach"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"attach": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "attach",
				Name: "tether_test_session",
			},
			Tty:    false,
			Attach: true,
			Cmd: metadata.Cmd{
				Path: "/usr/bin/tee",
				// grep, matching everything, reading from stdin
				Args: []string{"/usr/bin/tee", pathPrefix + "/tee.out"},
				Env:  []string{},
				Dir:  "/",
			},
		},
	}

	// generate a host key for the tether
	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		return nil, err
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}

	config.Key = pem.EncodeToMemory(&privateKeyBlock)

	return &config, nil
}

func TestAttach(t *testing.T) {
	// supply custom attach server so we can inspect its state
	testServer := &testAttachServer{
		updated: make(chan bool, 10),
	}
	server = testServer

	testSetup(t)
	// defer testTeardown(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestAttachConfig{}
	testConfig, err := cfg.LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		err := run(cfg)
		if err != nil {
			t.Error(err)
		}
	}()

	// wait for updates to occur
	<-testServer.updated

	if !testServer.enabled {
		t.Error("attach server was not enabled")
		return
	}

	// create client on the mock pipe
	conn, err := clientBackchannel(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	cconfig := &ssh.ClientConfig{
		User: "daemon",
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// create the SSH client
	sConn, chans, reqs, err := ssh.NewClientConn(conn, "notappliable", cconfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer sConn.Close()
	client := ssh.NewClient(sConn, chans, reqs)

	session, err := SSHAttach(client, testConfig.ID)
	if err != nil {
		t.Error(err)
		return
	}

	stdout := session.Stdout()

	// FIXME: the pipe pair are line buffered - how do I disable that so we don't have odd hangs to diagnose
	// when the trailing \n is missed
	testBytes := []byte("hello world!\n")
	// read from session into buffer
	buf := &bytes.Buffer{}
	done := make(chan bool)
	go func() { io.CopyN(buf, stdout, int64(len(testBytes))); done <- true }()

	// write something to echo
	log.Debug("sending test data")
	session.Stdin().Write(testBytes)
	log.Debug("sent test data")

	// wait for the close to propogate
	<-done
	session.Stdin().Close()

	if !bytes.Equal(buf.Bytes(), testBytes) {
		t.Errorf("expected: \"%s\", actual: \"%s\"", string(testBytes), buf.String())
		return
	}
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestAttachTTYConfig sets up the config for attach testing
//
type TestAttachTTYConfig struct{}

func (c *TestAttachTTYConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestAttachTTYConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "attach"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"attach": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "attach",
				Name: "tether_test_session",
			},
			Tty:    true,
			Attach: true,
			Cmd: metadata.Cmd{
				Path: "/usr/bin/tee",
				// grep, matching everything, reading from stdin
				Args: []string{"/usr/bin/tee", pathPrefix + "/tee.out"},
				Env:  []string{},
				Dir:  "/",
			},
		},
	}

	// generate a host key for the tether
	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		return nil, err
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}

	config.Key = pem.EncodeToMemory(&privateKeyBlock)

	return &config, nil
}

func TestAttachTTY(t *testing.T) {
	t.Skip("not sure how to test TTY yet")

	// supply custom attach server so we can inspect its state
	testServer := &testAttachServer{
		updated: make(chan bool, 10),
	}
	server = testServer

	testSetup(t)
	// defer testTeardown(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestAttachTTYConfig{}
	testConfig, err := cfg.LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		err := run(cfg)
		if err != nil {
			t.Error(err)
		}
	}()

	// wait for updates to occur
	<-testServer.updated

	if !testServer.enabled {
		t.Error("attach server was not enabled")
		return
	}

	// create client on the mock pipe
	conn, err := clientBackchannel(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	cconfig := &ssh.ClientConfig{
		User: "daemon",
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// create the SSH client
	sConn, chans, reqs, err := ssh.NewClientConn(conn, "notappliable", cconfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer sConn.Close()
	client := ssh.NewClient(sConn, chans, reqs)

	session, err := SSHAttach(client, testConfig.ID)
	if err != nil {
		t.Error(err)
		return
	}

	stdout := session.Stdout()

	// FIXME: this is line buffered - how do I disable that so we don't have odd hangs to diagnose
	// when the trailing \n is missed
	testBytes := []byte("hello world!\n")
	// read from session into buffer
	buf := &bytes.Buffer{}
	done := make(chan bool)
	go func() { io.CopyN(buf, stdout, int64(len(testBytes))); done <- true }()

	// write something to echo
	log.Debug("sending test data")
	session.Stdin().Write(testBytes)
	log.Debug("sent test data")

	// wait for the close to propogate
	<-done
	session.Stdin().Close()

	if !bytes.Equal(buf.Bytes(), testBytes) {
		t.Errorf("expected: \"%s\", actual: \"%s\"", string(testBytes), buf.String())
		return
	}
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestAttachTwoConfig sets up the config for attach testing - tests launching and
// attaching to two different processes simultaneously
//
type TestAttachTwoConfig struct{}

func (c *TestAttachTwoConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestAttachTwoConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "attachtwo"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"tee1": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "tee1",
				Name: "tether_test_session1",
			},
			Tty:    false,
			Attach: true,
			Cmd: metadata.Cmd{
				Path: "/usr/bin/tee",
				// grep, matching everything, reading from stdin
				Args: []string{"/usr/bin/tee", pathPrefix + "/tee.out"},
				Env:  []string{},
				Dir:  "/",
			},
		},
		"tee2": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "tee2",
				Name: "tether_test_session2",
			},
			Tty:    false,
			Attach: true,
			Cmd: metadata.Cmd{
				Path: "/usr/bin/tee",
				// grep, matching everything, reading from stdin
				Args: []string{"/usr/bin/tee", pathPrefix + "/tee2.out"},
				Env:  []string{},
				Dir:  "/",
			},
		},
	}

	// generate a host key for the tether
	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		return nil, err
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}

	config.Key = pem.EncodeToMemory(&privateKeyBlock)

	return &config, nil
}

func TestAttachTwo(t *testing.T) {
	// supply custom attach server so we can inspect its state
	testServer := &testAttachServer{
		updated: make(chan bool, 10),
	}
	server = testServer

	testSetup(t)
	// defer testTeardown(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestAttachTwoConfig{}
	_, err := cfg.LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		err := run(cfg)
		if err != nil {
			t.Error(err)
		}
	}()

	// wait for updates to occur
	<-testServer.updated

	if !testServer.enabled {
		t.Error("attach server was not enabled")
		return
	}

	// create client on the mock pipe
	conn, err := clientBackchannel(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	cconfig := &ssh.ClientConfig{
		User: "daemon",
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// create the SSH client
	sConn, chans, reqs, err := ssh.NewClientConn(conn, "notappliable", cconfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer sConn.Close()
	client := ssh.NewClient(sConn, chans, reqs)

	sessionA, err := SSHAttach(client, "tee1")
	if err != nil {
		t.Error(err)
		return
	}

	sessionB, err := SSHAttach(client, "tee2")
	if err != nil {
		t.Error(err)
		return
	}

	stdoutA := sessionA.Stdout()
	stdoutB := sessionB.Stdout()

	// FIXME: this is line buffered - how do I disable that so we don't have odd hangs to diagnose
	// when the trailing \n is missed
	testBytesA := []byte("hello world!\n")
	testBytesB := []byte("goodbye world!\n")
	// read from session into buffer
	bufA := &bytes.Buffer{}
	bufB := &bytes.Buffer{}

	doneA := make(chan bool)
	doneB := make(chan bool)
	go func() { io.CopyN(bufA, stdoutA, int64(len(testBytesA))); doneA <- true }()
	go func() { io.CopyN(bufB, stdoutB, int64(len(testBytesB))); doneB <- true }()

	// write something to echo
	log.Debug("sending test data")
	sessionA.Stdin().Write(testBytesA)
	sessionB.Stdin().Write(testBytesB)
	log.Debug("sent test data")

	// wait for the close to propogate
	<-doneA
	<-doneB
	sessionA.Stdin().Close()
	sessionB.Stdin().Close()

	if !bytes.Equal(bufA.Bytes(), testBytesA) {
		t.Errorf("expected: \"%s\", actual: \"%s\"", string(testBytesA), bufA.String())
		return
	}

	if !bytes.Equal(bufB.Bytes(), testBytesB) {
		t.Errorf("expected: \"%s\", actual: \"%s\"", string(testBytesB), bufB.String())
		return
	}
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestAttachInvalid sets up the config for attach testing - launches a process but
// tries to attach to an invalid session id
//
type TestAttachInvalidConfig struct{}

func (c *TestAttachInvalidConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestAttachInvalidConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "attachinvalid"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"valid": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "valid",
				Name: "tether_test_session",
			},
			Tty:    true,
			Attach: true,
			Cmd: metadata.Cmd{
				Path: "/usr/bin/tee",
				// grep, matching everything, reading from stdin
				Args: []string{"/usr/bin/tee", pathPrefix + "/tee.out"},
				Env:  []string{},
				Dir:  "/",
			},
		},
	}

	// generate a host key for the tether
	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		return nil, err
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}

	config.Key = pem.EncodeToMemory(&privateKeyBlock)

	return &config, nil
}

func TestAttachInvalid(t *testing.T) {
	// supply custom attach server so we can inspect its state
	testServer := &testAttachServer{
		updated: make(chan bool, 10),
	}
	server = testServer

	testSetup(t)
	// defer testTeardown(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestAttachInvalidConfig{}
	_, err := cfg.LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		err := run(cfg)
		if err != nil {
			t.Error(err)
		}
	}()

	// wait for updates to occur
	<-testServer.updated

	if !testServer.enabled {
		t.Error("attach server was not enabled")
		return
	}

	// create client on the mock pipe
	conn, err := clientBackchannel(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	cconfig := &ssh.ClientConfig{
		User: "daemon",
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// create the SSH client
	sConn, chans, reqs, err := ssh.NewClientConn(conn, "notappliable", cconfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer sConn.Close()
	client := ssh.NewClient(sConn, chans, reqs)

	_, err = SSHAttach(client, "invalidID")
	if err != nil {
		t.Log(err)
		return
	}

	t.Error("Expected to fail on attempt to attach to invalid session")
}

//
/////////////////////////////////////////////////////////////////////////////////////
