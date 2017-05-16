// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/vmware/vic/lib/iolog"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/netfilter"
	"github.com/vmware/vic/pkg/trace"
)

const runMountPoint = "/run"

type operations struct {
	tether.BaseOperations

	logging bool
}

func (t *operations) Log() (io.Writer, error) {
	defer trace.End(trace.Begin("operations.Log"))

	// redirect logging to the serial log
	log.Infof("opening %s/ttyS1 for debug log", pathPrefix)
	f, err := os.OpenFile(pathPrefix+"/ttyS1", os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, 0)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for debug log: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	if err := setTerminalSpeed(f.Fd()); err != nil {
		log.Errorf("Setting terminal speed failed with %s", err)
	}

	// enable raw mode
	_, err = terminal.MakeRaw(int(f.Fd()))
	if err != nil {
		detail := fmt.Sprintf("Making ttyS1 raw failed with %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	return io.MultiWriter(f, os.Stdout), nil
}

// sessionLogWriter returns a writer that will persist the session output
func (t *operations) SessionLog(session *tether.SessionConfig) (dio.DynamicMultiWriter, dio.DynamicMultiWriter, error) {
	defer trace.End(trace.Begin("configure session log writer"))

	if t.logging {
		detail := "unable to log more than one session concurrently to persistent logging"
		log.Warn(detail)
		// use multi-writer so it's still viable for attach
		return dio.MultiWriter(), dio.MultiWriter(), nil
	}

	t.logging = true

	// open SttyS2 for session logging
	log.Info("opening ttyS2 for session logging")
	f, err := os.OpenFile(pathPrefix+"/ttyS2", os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, 0)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for session log: %s", err)
		log.Error(detail)
		return nil, nil, errors.New(detail)
	}

	if err := setTerminalSpeed(f.Fd()); err != nil {
		log.Errorf("Setting terminal speed failed with %s", err)
	}

	// enable raw mode
	_, err = terminal.MakeRaw(int(f.Fd()))
	if err != nil {
		detail := fmt.Sprintf("Making ttyS2 raw failed with %s", err)
		log.Error(detail)
		return nil, nil, errors.New(detail)
	}

	// wrap output in a LogWriter to serialize it into our persisted
	// containerVM output format, using iolog.LogClock for timestamps
	lw := iolog.NewLogWriter(f, iolog.LogClock{})

	// use multi-writer so it goes to both screen and session log
	return dio.MultiWriter(lw, os.Stdout), dio.MultiWriter(lw, os.Stderr), nil
}

func (t *operations) Setup(sink tether.Config) error {
	if err := t.BaseOperations.Setup(sink); err != nil {
		return err
	}

	// symlink /etc/mtab to /proc/mounts
	var err error
	if err = tether.Sys.Syscall.Symlink("/proc/mounts", "/etc/mtab"); err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EEXIST {
			return err
		}
	}

	// unmount /run - https://github.com/vmware/vic/issues/1643
	if err = tether.Sys.Syscall.Unmount(runMountPoint, syscall.MNT_DETACH); err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EINVAL {
			return err
		}
	}

	return nil
}

// SetupFirewall sets up firewall rules on the external scope only.  Any
// portmaps are honored as are port exposes.
func (t *operations) SetupFirewall(config *tether.ExecutorConfig) error {
	// XXX It looks like we'd want to collect the errors here, but we
	// can't.  Since this is running inside init (tether) and tether
	// reaps all children, the os.exec package won't be able to collect
	// the error code in time before the reaper does.  The exec package
	// calls wait and attempts to collect its child, but the reaper will
	// have raptured the pid before that.  So, best effort, just keep going.
	_ = netfilter.Flush(context.Background(), "")

	// default rule set
	established := &netfilter.Rule{
		Chain:  netfilter.Input,
		States: []netfilter.State{netfilter.Established},
		Target: netfilter.Accept,
	}

	reject := &netfilter.Rule{
		Chain:  netfilter.Input,
		Target: netfilter.Reject,
	}

	for _, endpoint := range config.Networks {
		if endpoint.Network.Type == constants.ExternalScopeType {

			id, err := strconv.Atoi(endpoint.ID)
			if err != nil {
				log.Errorf("can't apply port rules: %s", err.Error())
				continue
			}

			iface, err := t.LinkBySlot(int32(id))
			if err != nil {
				log.Errorf("can't apply rules: %s", err.Error())
				continue
			}

			ifaceName := iface.Attrs().Name
			log.Debugf("slot %d -> %s", endpoint.ID, ifaceName)

			established.Interface = ifaceName
			_ = established.Commit(context.TODO())

			// handle the ports
			for _, p := range endpoint.Ports {
				// parse the port maps
				r, err := portToRule(p)
				if err != nil {
					log.Errorf("can't apply port rule (%s): %s", p, err.Error())
					continue
				}

				log.Infof("Applying rule for port %s", p)
				r.Interface = ifaceName
				_ = r.Commit(context.TODO())
			}

			reject.Interface = ifaceName
			_ = reject.Commit(context.TODO())

			break
		}
	}

	return nil
}

func portToRule(p string) (*netfilter.Rule, error) {
	if strings.Contains(p, ":") {
		return nil, errors.New("port maps are TBD")
	}

	// 9999/tcp
	s := strings.Split(p, "/")
	if len(s) != 2 {
		return nil, errors.New("can't parse port spec: " + p)
	}

	rule := &netfilter.Rule{
		Chain:     netfilter.Input,
		Interface: "external",
		Target:    netfilter.Accept,
	}

	switch netfilter.Protocol(s[1]) {
	case netfilter.UDP:
		rule.Protocol = netfilter.UDP
	case netfilter.TCP:
		rule.Protocol = netfilter.TCP

	default:
		return nil, errors.New("unknown protocol")
	}

	port, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, err
	}

	rule.FromPort = port

	return rule, nil
}
