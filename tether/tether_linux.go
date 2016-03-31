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
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/tether/serial"
)

// Mkdev will hopefully get rolled into go.sys at some point
func Mkdev(majorNumber int, minorNumber int) int {
	return (majorNumber << 8) | (minorNumber & 0xff) | ((minorNumber & 0xfff00) << 12)
}

func childReaper() {
	var incoming = make(chan os.Signal, 10)
	signal.Notify(incoming, syscall.SIGCHLD)

	log.Info("Started reaping child processes")

	for _ = range incoming {
		var status syscall.WaitStatus

		// reap until no more children to process
		for {
			log.Debugf("Inspecting children with status change")
			pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)
			if pid == 0 || err == syscall.ECHILD {
				log.Debug("No more child processes to reap")
				break
			}
			if err == nil {
				log.Debugf("Reaped process %d, return code: %d\n", pid, status.ExitStatus())

				cmd, ok := RemoveChildPid(pid)
				if ok {
					handleSessionExit(cmd)
				} else {
					// This is an adopted zombie. The Wait4 call
					// already clean it up from the kernel
					log.Infof("Reaped zombie process PID %d\n", pid)
				}
			} else {
				log.Warnf("Wait4 got error: %v\n", err)
			}
		}
	}
}

func setup() error {
	// seems necessary given rand.Reader access
	var err error

	// redirect logging to the serial log
	log.Infof("opening %s/ttyS1 for debug log", pathPrefix)
	out, err := os.OpenFile(pathPrefix+"/ttyS1", os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, 0777)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for debug log: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}
	log.SetOutput(io.MultiWriter(out, os.Stdout))

	// TODO: enabled for initial dev debugging only
	go func() {
		log.Info(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	rand.Reader, err = os.Open(pathPrefix + "/urandom")
	if err != nil {
		detail := fmt.Sprintf("failed to open new urandom device: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}

	// TODO: Call prctl with PR_SET_CHILD_SUBREAPER so that we reap regardless of pid 1 or not
	// we already get our direct children, but not lower in the hierarchy
	go childReaper()

	return nil
}

func backchannel() (net.Conn, error) {
	log.Info("opening ttyS0 for backchannel")
	f, err := os.OpenFile(pathPrefix+"/ttyS0", os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, 0777)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	log.Errorf("creating raw connection from ttyS0 (fd=%d)\n", f.Fd())
	conn, err := serial.NewFileConn(f)

	if err != nil {
		detail := fmt.Sprintf("failed to create raw connection from ttyS0 file handle: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// HACK: currently RawConn dosn't implement timeout
	serial.HandshakeServer(conn, time.Duration(10*time.Second))

	return conn, nil
}

// processEnvOS does OS specific checking and munging on the process environment prior to launch
func processEnvOS(env []string) []string {
	// TODO: figure out how we're going to specify user and pass all the settings along
	// in the meantime, hardcode HOME to /root
	homeIndex := -1
	for i, tuple := range env {
		if strings.HasPrefix(tuple, "HOME=") {
			homeIndex = i
			break
		}
	}
	if homeIndex == -1 {
		return append(env, "HOME=/root")
	}

	return env
}

// sessionLogWriter returns a writer that will persist the session output
func sessionLogWriter() (io.Writer, error) {
	// open SttyS2 for session logging
	log.Info("opening ttyS2 for session logging")
	f, err := os.OpenFile(pathPrefix+"/ttyS2", os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, 777)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for session log: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// use multi-writer so it goes to both screen and session log
	return io.MultiWriter(f, os.Stdout), nil
}
