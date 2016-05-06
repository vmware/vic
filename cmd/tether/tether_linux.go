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
	"unsafe"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/kr/pty"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/serial"
	"github.com/vmware/vic/pkg/trace"
)

// allow us to pick up some of the osops implementations when mocking
// allowing it to be less all or nothing
func init() {
	ops = &osopsLinux{}
	utils = &osopsLinux{}
}

var backchannelMode = os.ModePerm

// Mkdev will hopefully get rolled into go.sys at some point
func Mkdev(majorNumber int, minorNumber int) int {
	return (majorNumber << 8) | (minorNumber & 0xff) | ((minorNumber & 0xfff00) << 12)
}

var incoming chan os.Signal

// childReaper is used to handle events from child processes, including child exit.
// If running as pid=1 then this means it handles zombie process reaping for orphaned children
// as well as direct child processes.
func childReaper() {
	incoming = make(chan os.Signal, 10)
	signal.Notify(incoming, syscall.SIGCHLD)

	log.Info("Started reaping child processes")

	go func() {
		for _ = range incoming {
			var status syscall.WaitStatus

			func() {
				// general resiliency
				defer recover()

				// reap until no more children to process
				for {
					log.Debugf("Inspecting children with status change")
					pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)
					if pid == 0 || err == syscall.ECHILD {
						log.Debug("No more child processes to reap")
						break
					}
					if err == nil {
						if !status.Exited() {
							log.Debugf("Received notifcation about non-exit status change for %d:", pid)
							// no reaping or exit handling required
							continue
						}

						log.Debugf("Reaped process %d, return code: %d", pid, status.ExitStatus())

						session, ok := RemoveChildPid(pid)
						if ok {
							session.ExitStatus = status.ExitStatus()
							handleSessionExit(session)
						} else {
							// This is an adopted zombie. The Wait4 call
							// already clean it up from the kernel
							log.Infof("Reaped zombie process PID %d\n", pid)
						}
					} else {
						log.Warnf("Wait4 got error: %v\n", err)
					}
				}
			}()
		}
	}()
}

func (t *osopsLinux) setup() error {
	defer trace.End(trace.Begin("run OS specific tether setup"))

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
	childReaper()

	return nil
}

// cleanup does exactly that - this is called when the tether is being shut down.
// This should log errors, but no error is returned as this is a path of not return and
// there's not likely to be a remediation available
func (t *osopsLinux) cleanup() {
	defer trace.End(trace.Begin("running OS specific tether cleanup"))

	// stop child reaping
	log.Info("Shutting down reaper")
	signal.Reset(syscall.SIGCHLD)
	close(incoming)
	incoming = nil
}

func (t *osopsLinux) backchannel(ctx context.Context) (net.Conn, error) {
	defer trace.End(trace.Begin("establish tether backchannel"))

	log.Info("opening ttyS0 for backchannel")
	f, err := os.OpenFile(pathPrefix+"/ttyS0", os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, backchannelMode)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for backchannel: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	log.Infof("creating raw connection from ttyS0 (fd=%d)\n", f.Fd())
	conn, err := serial.NewFileConn(f)

	if err != nil {
		detail := fmt.Sprintf("failed to create raw connection from ttyS0 file handle: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// HACK: currently RawConn dosn't implement timeout so throttle the spinning
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

// processEnvOS does OS specific checking and munging on the process environment prior to launch
func (t *osopsLinux) processEnvOS(env []string) []string {
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
func (t *osopsLinux) sessionLogWriter() (dio.DynamicMultiWriter, error) {
	defer trace.End(trace.Begin("configure tether session log writer"))

	// open SttyS2 for session logging
	log.Info("opening ttyS2 for session logging")
	f, err := os.OpenFile(pathPrefix+"/ttyS2", os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, 777)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for session log: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	// use multi-writer so it goes to both screen and session log
	return dio.MultiWriter(f, os.Stdout), nil
}

func (t *osopsLinux) establishPty(session *SessionConfig) error {
	defer trace.End(trace.Begin("initializing pty handling for session " + session.ID))

	// TODO: if we want to allow raw output to the log so that subsequent tty enabled
	// processing receives the control characters then we should be binding the PTY
	// during attach, and using the same path we have for non-tty here
	var err error
	session.pty, err = pty.Start(&session.Cmd)
	if session.pty != nil {
		// TODO: do we need to ensure all reads have completed before calling Wait on the process?
		// it frees up all resources - does that mean it frees the output buffers?
		go func() {
			_, gerr := io.Copy(session.outwriter, session.pty)
			log.Debug(gerr)
		}()
		go func() {
			_, gerr := io.Copy(session.pty, session.reader)
			log.Debug(gerr)
		}()
	}

	return err
}

// The syscall struct
type winsize struct {
	wsRow    uint16
	wsCol    uint16
	wsXpixel uint16
	wsYpixel uint16
}

func (t *osopsLinux) resizePty(pty uintptr, winSize *WindowChangeMsg) error {
	defer trace.End(trace.Begin("resize pty"))

	ws := &winsize{uint16(winSize.Rows), uint16(winSize.Columns), uint16(winSize.WidthPx), uint16(winSize.HeightPx)}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		pty,
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func (t *osopsLinux) signalProcess(process *os.Process, sig ssh.Signal) error {
	signal := Signals[sig]
	defer trace.End(trace.Begin(fmt.Sprintf("signal process %d: %d", process.Pid, signal)))

	s := syscall.Signal(signal)
	return process.Signal(s)
}
