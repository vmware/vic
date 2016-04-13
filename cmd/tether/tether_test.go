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
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"syscall"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"

	"github.com/docker/docker/pkg/stringid"
	"github.com/vmware/vic/cmd/tether/serial"
	"github.com/vmware/vic/metadata"
)

var mockOps osopsMock

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)

	retCode := m.Run()

	// call with result of m.Run()
	os.Exit(retCode)
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

	updated chan bool
}

func (t *testAttachServer) start() error {
	err := t.attachServerSSH.start()

	t.updated <- true
	log.Info("Started test attach server")
	return err
}

func (t *testAttachServer) stop() {
	t.attachServerSSH.stop()

	log.Info("Stopped test attach server")
	t.updated <- true
}

func testSetup(t *testing.T) {
	var err error

	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	// use the mock ops - fresh one each time
	mockOps = osopsMock{
		// there's 4 util functions that will write to this
		updated: make(chan bool, 5),
	}
	ops = &mockOps

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

	if server == nil {
		server = &attachServerSSH{}
	}
}

func testTeardown(t *testing.T) {
	// let the main tether loop exit
	if reload != nil {
		close(reload)
		reload = nil
	}
	// cleanup
	os.RemoveAll(pathPrefix)
	log.SetOutput(os.Stdout)
}

func clientBackchannel(ctx context.Context) (net.Conn, error) {
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

/////////////////////////////////////////////////////////////////////////////////////
// TestPathLookup constructs the spec for a Session where the binary path must be
// resolved from the PATH environment variable - this is a variation from normal
// Cmd handling where that is done during creation of Cmd
//

type TestPathLookupConfig struct{}

func (c *TestPathLookupConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestPathLookupConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "pathlookup"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"pathlookup": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "pathlookup",
				Name: "tether_test_session",
			},
			Tty: false,
			Cmd: metadata.Cmd{
				// test relative path
				Path: "date",
				Args: []string{"date", "--reference=/"},
				Env:  []string{"PATH=/bin"},
				Dir:  "/bin",
			},
		},
	}

	return &config, nil
}

func TestPathLookup(t *testing.T) {
	t.Skip("PATH based lookup not yet implemented")

	testSetup(t)
	defer testTeardown(t)

	if err := run(&TestPathLookupConfig{}); err != nil {
		t.Error(err)
		return
	}
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestRelativePathConfig constructs the spec for a Session with relative path to existing
// binary - this tests resolution vs the working directory
//

type TestRelativePathConfig struct{}

func (c *TestRelativePathConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestRelativePathConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "relpath"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"relpath": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "relpath",
				Name: "tether_test_session",
			},
			Tty: false,
			Cmd: metadata.Cmd{
				// test relative path
				Path: "date",
				Args: []string{"date", "--reference=/"},
				Env:  []string{"PATH=/bin"},
				Dir:  "/bin",
			},
		},
	}

	return &config, nil
}

func TestRelativePath(t *testing.T) {
	t.Skip("Relative path resolution not yet implemented")

	testSetup(t)
	defer testTeardown(t)

	if err := run(&TestRelativePathConfig{}); err != nil {
		t.Error(err)
		return
	}
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestAbsPathConfig constructs the spec for a Session with absolute path to existing
// binary
//
type TestAbsPathConfig struct{}

func (c *TestAbsPathConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestAbsPathConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "abspath"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"abspath": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "abspath",
				Name: "tether_test_session",
			},
			Tty: false,
			Cmd: metadata.Cmd{
				// test abs path
				Path: "/bin/date",
				Args: []string{"date", "--reference=/"},
				Env:  []string{},
				Dir:  "/",
			},
		},
	}

	return &config, nil
}

func TestAbsPath(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestAbsPathConfig{}

	if err := run(cfg); err != nil {
		t.Error(err)
	}

	// check the exit code was set
	status := sessions["abspath"].exitStatus
	if status != 0 {
		t.Error("reference process 'data --reference=/' did not exit cleanly: %d", status)
		return
	}

	// read the output from the session
	log, err := ioutil.ReadFile(pathPrefix + "/ttyS2")
	if err != nil {
		fmt.Printf("Failed to open log file for command: %s", err)
		t.Error(err)
		return
	}

	// run the command directly
	out, err := exec.Command("/bin/date", "--reference=/").Output()
	if err != nil {
		fmt.Printf("Failed to run date for comparison data: %s", err)
		t.Error(err)
		return
	}

	if !bytes.Equal(out, log) {
		err := fmt.Errorf("Actual and expected output did not match\nExpected: %s\nActual:   %s\n", out, log)
		t.Error(err)
		return
	}

}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestMissingBinaryConfig constructs the spec for a Session with invalid binary path
//
type TestMissingBinaryConfig struct{}

func (c *TestMissingBinaryConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestMissingBinaryConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "missing"
	config.Name = "tether_test_executor"
	config.Sessions = map[string]metadata.SessionConfig{
		"missing": metadata.SessionConfig{
			Common: metadata.Common{
				ID:   "missing",
				Name: "tether_test_session",
			},
			Tty: false,
			Cmd: metadata.Cmd{
				// test relative path
				Path: "/not/there",
				Args: []string{"/not/there"},
				Env:  []string{"PATH=/not"},
				Dir:  "/",
			},
		},
	}

	return &config, nil
}

func TestMissingBinary(t *testing.T) {
	testSetup(t)

	err := run(&TestMissingBinaryConfig{})
	if err == nil {
		t.Error("Expected error from missing binary")
	}

	// read the output from the session
	log, err := ioutil.ReadFile(pathPrefix + "/ttyS2")
	if err != nil {
		fmt.Printf("Failed to open log file for command: %s", err)
	}
	if len(log) > 0 {
		fmt.Printf("Command output: %s", string(log))
	}

	testTeardown(t)
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestHostnameConfig constructs the spec with no sesisons specifically for testing
// hostname setting - this checks the value passed to the mocked SetHostname
//
type TestSetHostnameConfig struct{}

func (c *TestSetHostnameConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestSetHostnameConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "sethostname"
	config.Name = "tether_test_executor"

	return &config, nil
}

func TestSetHostname(t *testing.T) {
	testSetup(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestSetHostnameConfig{}
	testConfig, _ := cfg.LoadConfig()
	go func() {
		err := run(cfg)
		if err != nil {
			t.Error(err)
		}
	}()

	// wait for updates to occur
	<-mockOps.updated

	expected := stringid.TruncateID(testConfig.ID)
	if mockOps.hostname != expected {
		t.Errorf("expected: %s, actual: %s", expected, mockOps.hostname)
	}

	testTeardown(t)
}

//
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// TestSetIpAddressConfig constructs the spec for setting IP addresses - this checks
// the values passed to the Apply mock match those from the test config
//
type TestIPConfig struct{}

func (c *TestIPConfig) StoreConfig(*metadata.ExecutorConfig) (string, error) {
	return "", errors.New("not implemented")
}
func (c *TestIPConfig) LoadConfig() (*metadata.ExecutorConfig, error) {
	config := metadata.ExecutorConfig{}

	config.ID = "ipconfig"
	config.Name = "tether_test_executor"

	return &config, nil
}

func TestSetIpAddress(t *testing.T) {
	t.Skip("Network configuration processing not yet implemented")

	testSetup(t)

	// if there's no session command with guaranteed exit then tether needs to run in the background
	cfg := &TestIPConfig{}
	testConfig, _ := cfg.LoadConfig()
	go func() {
		err := run(cfg)
		if err != nil {
			t.Error(err)
		}
	}()

	// wait for updates to occur
	<-mockOps.updated

	// TEST LOGIC GOES HERE
	_ = testConfig

	testTeardown(t)
}

//
/////////////////////////////////////////////////////////////////////////////////////

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
// TestAttachConfig sets up the config for attach testing - the grep will echo anything
// sent and adds colour which is useful for tty testing
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
			// Cmd: metadata.Cmd{
			// 	Path: "/bin/grep",
			// 	// grep, matching everything, reading from stdin
			// 	Args: []string{"/bin/grep", ".", "-"},
			// 	Env:  []string{},
			// 	Dir:  "/",
			// },
			// Cmd: metadata.Cmd{
			// 	Path: "/bin/bash",
			// 	// grep, matching everything, reading from stdin
			// 	Args: []string{},
			// 	Env:  []string{},
			// 	Dir:  "/",
			// },
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
