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
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/vmware/vic/lib/iolog"
	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/pkg/dio"
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
