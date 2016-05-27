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

package tether

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"testing"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"

	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

var mocked mocker

type mocker struct {
	base BaseOperations

	// allow tests to tell when the tether has finished setup
	started chan bool
	// allow tests to tell when the tether has finished
	cleaned chan bool

	// session output gets logged here
	log bytes.Buffer

	// the hostname of the system
	hostname string
	// the ip configuration for mac indexed interfaces
	ips map[string]net.IPNet
	// filesystem mounts, indexed by disk label
	mounts map[string]string

	windowCol uint32
	windowRow uint32
	signal    ssh.Signal
}

// Start implements the extension method
func (t *mocker) Start() error {
	return nil
}

// Stop implements the extension method
func (t *mocker) Stop() error {
	close(t.cleaned)
	return nil
}

// Reload implements the extension method
func (t *mocker) Reload(config *ExecutorConfig) error {
	// the tether has definitely finished it's startup by the time we hit this
	close(t.started)
	return nil
}

func (t *mocker) Setup() error {
	return nil
}

func (t *mocker) Cleanup() error {
	return nil
}

func (t *mocker) Log() (io.Writer, error) {
	return nil, nil
}

func (t *mocker) SessionLog(session *SessionConfig) (dio.DynamicMultiWriter, error) {
	return dio.MultiWriter(&t.log, os.Stdout), nil
}

func (t *mocker) HandleSessionExit(config *ExecutorConfig, session *SessionConfig) bool {
	// check for executor behaviour
	return session.ID == config.ID
}

func (t *mocker) ProcessEnv(env []string) []string {
	return t.base.ProcessEnv(env)
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
func (t *mocker) Fork() error {
	defer trace.End(trace.Begin("mocking fork"))
	return errors.New("Fork test not implemented")
}

func startTether(t *testing.T, cfg *metadata.ExecutorConfig) (Tether, extraconfig.DataSource) {
	store := map[string]string{}
	sink := extraconfig.MapSink(store)
	src := extraconfig.MapSource(store)
	extraconfig.Encode(sink, cfg)
	log.Debugf("Test configuration: %#v", sink)

	tthr := New(src, sink, &mocked)
	tthr.Register("mocker", &mocked)

	// run the tether to service the attach
	go func() {
		erR := tthr.Start()
		if erR != nil {
			t.Error(erR)
		}
	}()

	return tthr, src
}

func runTether(t *testing.T, cfg *metadata.ExecutorConfig) (Tether, extraconfig.DataSource, error) {
	store := map[string]string{}
	sink := extraconfig.MapSink(store)
	src := extraconfig.MapSource(store)
	extraconfig.Encode(sink, cfg)
	log.Debugf("Test configuration: %#v", sink)

	tthr := New(src, sink, &mocked)
	tthr.Register("mocker", &mocked)

	// run the tether to service the attach
	erR := tthr.Start()

	return tthr, src, erR
}

func OptionValueArrayToString(options []types.BaseOptionValue) string {
	// create the key/value store from the extraconfig slice for lookups
	kv := make(map[string]string)
	for i := range options {
		k := options[i].GetOptionValue().Key
		v := options[i].GetOptionValue().Value.(string)
		kv[k] = v
	}

	return fmt.Sprintf("%#v", kv)
}

func testSetup(t *testing.T) {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	log.Infof("Started test setup for %s", name)

	// use the mock ops - fresh one each time as tests might apply different mocked calls
	mocked = mocker{
		started: make(chan bool, 0),
		cleaned: make(chan bool, 0),
	}

	// pathPrefix, err := ioutil.TempDir("", path.Base(name))
	// if err != nil {
	// 	fmt.Println(err)
	// 	t.Error(err)
	// }

	// err = os.MkdirAll(pathPrefix, 0777)
	// if err != nil {
	// 	fmt.Println(err)
	// 	t.Error(err)
	// }
	// log.Infof("Using %s as test prefix", pathPrefix)
}

func testTeardown(t *testing.T) {
	// cleanup
	// os.RemoveAll(pathPrefix)
	log.SetOutput(os.Stdout)

	<-mocked.cleaned

	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	log.Infof("Finished test teardown for %s", name)
}
