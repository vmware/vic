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
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/kr/pty"
	"github.com/vmware/vic/pkg/trace"
)

const (
	//https://github.com/golang/go/blob/master/src/syscall/zerrors_linux_arm64.go#L919
	SetChildSubreaper = 0x24
)

// Mkdev will hopefully get rolled into go.sys at some point
func Mkdev(majorNumber int, minorNumber int) int {
	return (majorNumber << 8) | (minorNumber & 0xff) | ((minorNumber & 0xfff00) << 12)
}

// childReaper is used to handle events from child processes, including child exit.
// If running as pid=1 then this means it handles zombie process reaping for orphaned children
// as well as direct child processes.
func (t *tether) childReaper() error {
	signal.Notify(t.incoming, syscall.SIGCHLD)

	/*
	   PR_SET_CHILD_SUBREAPER (since Linux 3.4)
	          If arg2 is nonzero, set the "child subreaper" attribute of the
	          calling process; if arg2 is zero, unset the attribute.  When a
	          process is marked as a child subreaper, all of the children
	          that it creates, and their descendants, will be marked as
	          having a subreaper.  In effect, a subreaper fulfills the role
	          of init(1) for its descendant processes.  Upon termination of
	          a process that is orphaned (i.e., its immediate parent has
	          already terminated) and marked as having a subreaper, the
	          nearest still living ancestor subreaper will receive a SIGCHLD
	          signal and be able to wait(2) on the process to discover its
	          termination status.
	*/
	if _, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, SetChildSubreaper, uintptr(1), 0); err != 0 {
		return err
	}

	log.Info("Started reaping child processes")

	go func() {
		for _ = range t.incoming {
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
						if !status.Exited() && !status.Signaled() {
							log.Debugf("Received notifcation about non-exit status change for %d: %d", pid, status)
							// no reaping or exit handling required
							continue
						}

						log.Debugf("Reaped process %d, return code: %d", pid, status.ExitStatus())

						session, ok := t.removeChildPid(pid)
						if ok {
							// will be -1 if terminated due to signal
							session.ExitStatus = status.ExitStatus()
							t.handleSessionExit(session)
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

	return nil
}

func (t *tether) stopReaper() {
	defer trace.End(trace.Begin("Shutting down child reaping"))

	// stop child reaping
	log.Info("Shutting down reaper")
	signal.Reset(syscall.SIGCHLD)
	close(t.incoming)
}

func findExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return os.ErrPermission
}

// lookPath searches for an executable binary named file in the directories
// specified by the path argument.
// This is a direct modification of the unix os/exec core libary impl
func lookPath(file string, env []string, dir string) (string, error) {
	// if it starts with a ./ or ../ it's a relative path
	// need to check explicitly to allow execution of .hidden files

	if strings.HasPrefix(file, "./") || strings.HasPrefix(file, "../") {
		file = fmt.Sprintf("%s%c%s", dir, os.PathSeparator, file)
		err := findExecutable(file)
		if err == nil {
			return file, nil
		}
		return "", err
	}

	// check if it's already a path spec
	if strings.Contains(file, "/") {
		err := findExecutable(file)
		if err == nil {
			return file, nil
		}
		return "", err
	}

	// extract path from the env
	var pathenv string
	for _, value := range env {
		if strings.HasPrefix(value, "PATH=") {
			pathenv = value
			break
		}
	}

	pathval := strings.TrimPrefix(pathenv, "PATH=")

	dirs := filepath.SplitList(pathval)
	for _, dir := range dirs {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		path := dir + "/" + file
		if err := findExecutable(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("%s: no such executable in PATH", file)
}

func establishPty(session *SessionConfig) error {
	defer trace.End(trace.Begin("initializing pty handling for session " + session.ID))

	// TODO: if we want to allow raw output to the log so that subsequent tty enabled
	// processing receives the control characters then we should be binding the PTY
	// during attach, and using the same path we have for non-tty here
	var err error
	session.Pty, err = pty.Start(&session.Cmd)
	if session.Pty != nil {
		// TODO: do we need to ensure all reads have completed before calling Wait on the process?
		// it frees up all resources - does that mean it frees the output buffers?
		go func() {
			_, gerr := io.Copy(session.Outwriter, session.Pty)
			log.Debug(gerr)
		}()
		go func() {
			_, gerr := io.Copy(session.Pty, session.Reader)
			log.Debug(gerr)
		}()
	}

	return err
}
