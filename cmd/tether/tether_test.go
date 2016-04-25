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
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"runtime"
	"syscall"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"

	"github.com/vmware/vic/metadata"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/serial"
	"github.com/vmware/vic/pkg/trace"
)

var mocked mocker

// store the OS specific ops
var specificOps osops

// store the OS specific utils
var specificUtils utilities

type mocker struct {
	// so that we can call through to the core methods where viable
	ops   osops
	utils utilities

	// allow tests to tell when the tether has finished setup
	started chan bool
	// allow tests to tell when the tether has finished
	cleaned chan bool

	// the hostname of the system
	hostname string
	// the ip configuration for mac indexed interfaces
	ips map[string]net.IPNet
	// filesystem mounts, indexed by disk label
	mounts map[string]string
}

func (t *mocker) setup() error {
	err := t.utils.setup()
	close(t.started)
	return err
}

func (t *mocker) cleanup() {
	t.utils.cleanup()
	close(t.cleaned)
}

func (t *mocker) sessionLogWriter() (dio.DynamicMultiWriter, error) {
	return t.utils.sessionLogWriter()
}

func (t *mocker) processEnvOS(env []string) []string {
	return t.utils.processEnvOS(env)
}

func (t *mocker) establishPty(live *liveSession) error {
	return t.utils.establishPty(live)
}

func (t *mocker) resizePty(pty uintptr, winSize *WindowChangeMsg) error {
	return t.utils.resizePty(pty, winSize)
}

func (t *mocker) signalProcess(process *os.Process, sig ssh.Signal) error {
	return t.utils.signalProcess(process, sig)
}

// SetHostname sets both the kernel hostname and /etc/hostname to the specified string
func (t *mocker) SetHostname(hostname string) error {
	defer trace.End(trace.Begin("mocking hostname to " + hostname))

	// TODO: we could mock at a much finer granularity, only extracting the syscall
	// that would exercise the file modification paths, however it's much less generalizable
	t.hostname = hostname
	return nil
}

// Apply takes the network endpoint configuration and applies it to the system
func (t *mocker) Apply(endpoint *metadata.NetworkEndpoint) error {
	defer trace.End(trace.Begin("mocking endpoint configuration for " + endpoint.Network.Name))
	return errors.New("Apply test not implemented")
}

// MountLabel performs a mount with the source treated as a disk label
// This assumes that /dev/disk/by-label is being populated, probably by udev
func (t *mocker) MountLabel(label, target string, ctx context.Context) error {
	defer trace.End(trace.Begin(fmt.Sprintf("mocking mounting %s on %s", label, target)))

	if t.mounts == nil {
		t.mounts = make(map[string]string)
	}

	t.mounts[label] = target
	return nil
}

// Fork triggers vmfork and handles the necessary pre/post OS level operations
func (t *mocker) Fork(config *metadata.ExecutorConfig) error {
	defer trace.End(trace.Begin("mocking fork"))
	return errors.New("Fork test not implemented")
}

func (t *mocker) backchannel(ctx context.Context) (net.Conn, error) {
	log.Info("opening ttyS0 pipe pair for backchannel")
	c, err := os.OpenFile(pathPrefix+"/ttyS0c", os.O_WRONLY|syscall.O_NOCTTY, 0777)
	if err != nil {
		detail := fmt.Sprintf("failed to open cpipe for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	s, err := os.OpenFile(pathPrefix+"/ttyS0s", os.O_RDONLY|syscall.O_NOCTTY, 0777)
	if err != nil {
		detail := fmt.Sprintf("failed to open spipe for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	log.Infof("creating raw connection from ttyS0 pipe pair (c=%d, s=%d)\n", c.Fd(), s.Fd())
	conn, err := serial.NewHalfDuplixFileConn(s, c, pathPrefix+"/ttyS0", "file")

	if err != nil {
		detail := fmt.Sprintf("failed to create raw connection from ttyS0 pipe pair: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// still run handshake over it to test that
	ticker := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			err := serial.HandshakeServer(ctx, conn)
			if err == nil {
				return conn, nil
			}
		case <-ctx.Done():
			conn.Close()
			ticker.Stop()
			return nil, ctx.Err()
		}
	}
}

// createFakeDevices creates regular files or pipes in place of the char devices used
// in a full VM
func createFakeDevices() error {
	var err error
	// create control channel
	path := fmt.Sprintf("%s/ttyS0", pathPrefix)
	err = MkNamedPipe(path+"s", os.ModePerm)
	if err != nil {
		detail := fmt.Sprintf("failed to create fifo pipe %ss for com0: %s", path, err)
		return errors.New(detail)
	}
	err = MkNamedPipe(path+"c", os.ModePerm)
	if err != nil {
		detail := fmt.Sprintf("failed to create fifo pipe %sc for com0: %s", path, err)
		return errors.New(detail)
	}
	log.Debugf("created %s/ttyS0{c,s} as raw conn pipes", pathPrefix)

	// others are non-interactive
	for i := 1; i < 3; i++ {
		path = fmt.Sprintf("%s/ttyS%d", pathPrefix, i)
		_, err = os.Create(path)
		if err != nil {
			detail := fmt.Sprintf("failed to create %s for com%d: %s", path, i+1, err)
			return errors.New(detail)
		}
		log.Debugf("created %s as persistent log destinations", path)
	}

	// make an access to urandom
	path = fmt.Sprintf("%s/urandom", pathPrefix)
	err = os.Symlink("/dev/urandom", path)
	if err != nil {
		detail := fmt.Sprintf("failed to create urandom access: %s", err)
		return errors.New(detail)
	}

	return nil
}

type testAttachServer struct {
	attachServerSSH
	enabled bool
	updated chan bool
}

func (t *testAttachServer) start() error {
	err := t.attachServerSSH.start()

	if err == nil {
		t.updated <- true
		t.enabled = true
	}

	log.Info("Started test attach server")
	return err
}

func (t *testAttachServer) stop() {
	if t.enabled {
		t.attachServerSSH.stop()

		log.Info("Stopped test attach server")
		t.updated <- true
		t.enabled = false
	}
}

// TestMain simply so we have control of debugging level and somewhere to call package wide test setup
func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)

	// save the base os specific structures
	specificOps = ops
	specificUtils = utils

	retCode := m.Run()

	// call with result of m.Run()
	os.Exit(retCode)
}

func testSetup(t *testing.T) {
	var err error

	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	log.Infof("Started test setup for %s", name)

	// use the mock ops - fresh one each time as tests might apply different mocked calls
	mocked = mocker{
		ops:     specificOps,
		utils:   specificUtils,
		started: make(chan bool, 0),
		cleaned: make(chan bool, 0),
	}
	ops = &mocked
	utils = &mocked

	pathPrefix, err = ioutil.TempDir("", path.Base(name))
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	err = os.MkdirAll(pathPrefix, 0777)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}
	log.Infof("Using %s as test prefix", pathPrefix)

	backchannelMode = os.ModeNamedPipe | os.ModePerm
	err = createFakeDevices()
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	// supply custom attach server so we can inspect its state
	testServer := &testAttachServer{
		updated: make(chan bool, 10),
	}
	server = testServer
}

func testTeardown(t *testing.T) {
	// let the main tether loop exit
	r := reload
	reload = nil
	if r != nil {
		close(r)
	}
	// cleanup
	os.RemoveAll(pathPrefix)
	log.SetOutput(os.Stdout)

	<-mocked.cleaned

	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	log.Infof("Finished test teardown for %s", name)
}

// create client on the mock pipe
func mockSerialConnection(ctx context.Context) (net.Conn, error) {
	log.Info("opening ttyS0 pipe pair for backchannel")
	c, err := os.OpenFile(pathPrefix+"/ttyS0c", os.O_RDONLY|syscall.O_NOCTTY, 0777)
	if err != nil {
		detail := fmt.Sprintf("failed to open cpipe for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	s, err := os.OpenFile(pathPrefix+"/ttyS0s", os.O_WRONLY|syscall.O_NOCTTY, 0777)
	if err != nil {
		detail := fmt.Sprintf("failed to open spipe for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	log.Infof("creating raw connection from ttyS0 pipe pair (c=%d, s=%d)\n", c.Fd(), s.Fd())
	conn, err := serial.NewHalfDuplixFileConn(c, s, pathPrefix+"/ttyS0", "file")

	if err != nil {
		detail := fmt.Sprintf("failed to create raw connection from ttyS0 pipe pair: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// HACK: currently RawConn dosn't implement timeout so throttle the spinning
	ticker := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			// FIXME: need to implement timeout of purging hangs with no content
			// on the pipe
			// serial.PurgeIncoming(ctx, conn)
			err := serial.HandshakeClient(ctx, conn)
			if err == nil {
				return conn, nil
			}
		case <-ctx.Done():
			conn.Close()
			ticker.Stop()
			return nil, ctx.Err()
		}
	}
}
