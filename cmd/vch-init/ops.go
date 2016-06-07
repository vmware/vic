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
	"io"
	"net/http"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/trace"
)

// pathPrefix is present to allow the various files referenced by tether to be placed
// in specific directories, primarily for testing.
var pathPrefix string

const (
	logDir  = "/var/log/vic"
	initLog = "init.log"
)

type operations struct {
	tether.BaseOperations
}

func (t *operations) Setup() error {
	// make the logging directory
	os.MkdirAll(fmt.Sprintf("%s%c%s", pathPrefix, os.PathSeparator, logDir), 0777)

	// TODO: enabled for initial dev debugging only
	log.Info("Launching pprof server on port 6060")
	go func() {
		log.Info(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	return nil
}

func (t *operations) Cleanup() error {
	return nil
}

// HandleSessionExit controls the behaviour on session exit - for the tether if the session exiting
// is the primary session (i.e. SessionID matches ExecutorID) then we exit everything.
func (t *operations) HandleSessionExit(config *tether.ExecutorConfig, session *tether.SessionConfig) bool {
	// This is the appliance so relaunch
	tthr.Reload()
	return false
}

func (t *operations) Log() (io.Writer, error) {
	defer trace.End(trace.Begin("operations.Log"))

	logPath := fmt.Sprintf("%s%c%s%[2]c%s", pathPrefix, os.PathSeparator, logDir, initLog)

	// redirect logging to /var/log/vic/init
	log.Infof("opening %s for debug log", logPath)
	out, err := os.OpenFile(logPath, os.O_RDWR|os.O_APPEND|os.O_SYNC|syscall.O_NOCTTY, 0777)
	if err != nil {
		detail := fmt.Sprintf("failed to open file port for debug log: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	return out, nil
}

// sessionLogWriter returns a writer that will persist the session output
func (t *operations) SessionLog(session *tether.SessionConfig) (dio.DynamicMultiWriter, error) {
	defer trace.End(trace.Begin("configure session log writer"))

	logPath := fmt.Sprintf("%s%c%s%[2]c%s", pathPrefix, os.PathSeparator, logDir, session.Name)

	// open SttyS2 for session logging
	log.Infof("opening %s for session logging", logPath)
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_APPEND|os.O_SYNC|syscall.O_NOCTTY, 777)
	if err != nil {
		detail := fmt.Sprintf("failed to open file for session log: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// use multi-writer so it goes to both screen and session log
	return dio.MultiWriter(f, os.Stdout), nil
}
