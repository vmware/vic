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
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stringid"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

// pathPrefix is used for testing - it allows for creating and manupulating files outside of
// a full containerVM environment
var pathPrefix string

// the reload channel is used to block reloading of the config
// there will only be something on this channel on two occasions:
// 1. initial start
// 2. post-vmfork
var reload chan bool

// config holds the main configuration for the executor
var config *ExecutorConfig

var dataSource extraconfig.DataSource
var dataSink extraconfig.DataSink

// RemoveChildPid is a synchronized accessor for the pid map the deletes the entry and returns the value
func RemoveChildPid(pid int) (*SessionConfig, bool) {
	config.pidMutex.Lock()
	defer config.pidMutex.Unlock()

	session, ok := config.pids[pid]
	delete(config.pids, pid)
	return session, ok
}

// LenChildPid returns the number of entries
func LenChildPid() int {
	config.pidMutex.Lock()
	defer config.pidMutex.Unlock()

	return len(config.pids)
}

func run(src extraconfig.DataSource, sink extraconfig.DataSink) error {
	defer trace.End(trace.Begin("main tether loop"))

	// remake all of the main management structures so there's no cross contamination between tests
	reload = make(chan bool, 1)
	config = &ExecutorConfig{
		pids: make(map[int]*SessionConfig),
	}

	dataSource = src
	dataSink = sink

	// HACK: workaround file descriptor conflict in pipe2 return from the exec.Command.Start
	// it's not clear whether this is a cross platform issue, or still an issue as of this commit
	// keeping it until there's time to verify and fix properly with a Go PR.
	_, _, _ = os.Pipe()

	// perform basic one off OS specific setup
	err := utils.setup()
	if err != nil {
		detail := fmt.Sprintf("failed during setup: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}

	defer func() {
		// perform basic cleanup
		reload = nil
		// FIXME: Cannot clean up sessions until we are persisting exit status elsewhere for test validation
		//    also referenced in handleSessionExit
		// config = nil

		utils.cleanup()
	}()

	// initial setup, so seed this
	reload <- true
	for _ = range reload {
		// load the config - this modifies the structure values in place
		extraconfig.Decode(src, config)
		if err != nil {
			detail := fmt.Sprintf("failed to load config: %s", err)
			log.Error(detail)
			// we don't attempt to recover from this - our async config and reporting channel isn't working
			// as expected so just exit
			return errors.New(detail)
		}

		logConfig(config)

		if err := ops.SetHostname(stringid.TruncateID(config.ID)); err != nil {
			detail := fmt.Sprintf("failed to set hostname: %s", err)
			log.Error(detail)
			// we don't attempt to recover from this - it's a fundemental misconfiguration
			// so just exit
			return errors.New(detail)
		}

		/*
			for _, v := range config.Networks {
				if err := ops.Apply(v); err != nil {
					detail := fmt.Sprintf("failed to apply network endpoint config: %s", err)
					log.Error(detail)
					return errors.New(detail)
				}
			}
		*/

		// process the sessions and launch if needed
		attach := false
		for id, session := range config.Sessions {
			log.Debugf("Processing config for session %s", session.ID)
			var proc = session.Cmd.Process

			if session.Attach {
				attach = true
				log.Debugf("Session %s is configured for attach", session.ID)
				// this will return nil if already running
				err := server.start()
				if err != nil {
					detail := fmt.Sprintf("unable to start attach server: %s", err)
					log.Error(detail)
					continue
				}
			}

			// check if session is alive and well
			if proc != nil && proc.Signal(syscall.Signal(0)) != nil {
				log.Debugf("Process for session %s is already running", session.ID)
				continue
			}

			// check if session has never been started
			if proc == nil {
				log.Infof("Launching process for session %s\n", session.ID)
				err := launch(session)
				if err != nil {
					detail := fmt.Sprintf("failed to launch %s for %s: %s", session.Cmd.Path, id, err)
					log.Error(detail)

					// TODO: check if failure to launch this is fatal to everything in this containerVM
					return errors.New(detail)
				}

				// TODO: decide how to handle restart - probably needs to glue into the child reaping
			}

			// handle exited session
			// TODO
		}

		// none of the sessions allows attach, so stop the server
		if !attach {
			server.stop()
		}
	}

	return nil
}

// handleSessionExit processes the result from the session command, records it in persistent
// maner and determines if the Executor should exit
func handleSessionExit(session *SessionConfig) error {
	defer trace.End(trace.Begin("handling exit of session " + session.ID))

	// close down the IO
	session.reader.Close()
	// live.outwriter.Close()
	// live.errwriter.Close()

	// flush session log output

	// record exit status
	// FIXME: we cannot have this embedded knowledge of the extraconfig encoding pattern, but not
	// currently sure how to expose it neatly via a utility function
	extraconfig.EncodeWithPrefix(dataSink, session.ExitStatus, fmt.Sprintf("sessions|%s/status", session.ID))
	log.Infof("%s exit code: %d", session.ID, session.ExitStatus)

	// check for executor behaviour
	if LenChildPid() == 0 {
		// let the main loop exit if there's no more sessions to wait on
		if reload != nil {
			close(reload)
			reload = nil
		}
	}

	return nil
}

// launch will launch the command defined in the session.
// This will return an error if the session fails to launch
func launch(session *SessionConfig) error {
	defer trace.End(trace.Begin("launching session " + session.ID))

	logwriter, err := utils.sessionLogWriter()
	if err != nil {
		detail := fmt.Sprintf("failed to get log writer for session: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}

	// we store these outside of the session.Cmd struct so that there's consistent
	// handling between tty & non-tty paths
	session.outwriter = logwriter
	session.errwriter = logwriter
	session.reader = dio.MultiReader()

	session.Cmd.Env = utils.processEnvOS(session.Cmd.Env)
	session.Cmd.Stdout = session.outwriter
	session.Cmd.Stderr = session.errwriter
	session.Cmd.Stdin = session.reader

	// Use the mutex to make creating a child and adding the child pid into the
	// childPidTable appear atomic to the reaper function. Use a anonymous function
	// so we can defer unlocking locally
	err = func() error {
		config.pidMutex.Lock()
		defer config.pidMutex.Unlock()

		log.Infof("Launching command %#v\n", session.Cmd.Args)
		if !session.Tty {
			err = session.Cmd.Start()
		} else {
			err = utils.establishPty(session)
		}

		if err != nil {
			detail := fmt.Sprintf("failed to start container process: %s", err)
			log.Error(detail)
			return errors.New(detail)
		}

		// ChildReaper will use this channel to inform us the wait status of the child.
		config.pids[session.Cmd.Process.Pid] = session

		log.Debugf("Launched command with pid %d", session.Cmd.Process.Pid)

		return nil
	}()

	return err
}

func logConfig(config *ExecutorConfig) {
	// just pretty print the json for now
	log.Info("Loaded executor config")

	// TODO: investigate whether it's the govmomi types package cause the binary size
	// inflation - if so we need an alternative approach here or in extraconfig
	if log.GetLevel() == log.DebugLevel {
		sink := map[string]string{}
		extraconfig.Encode(extraconfig.MapSink(sink), config)

		for k, v := range sink {
			log.Debugf("%s: %s", k, v)
		}
	}
}

func forkHandler() {
	defer trace.End(trace.Begin("start fork trigger handler"))

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in StartConnectionManager", r)
		}
	}()

	incoming := make(chan os.Signal, 1)
	signal.Notify(incoming, syscall.SIGABRT)

	log.Info("SIGABRT handling initialized for fork support")
	for _ = range incoming {
		// validate that this is a fork trigger and not just a random signal from
		// container processes
		log.Info("Received SIGABRT - preparing to transition to fork parent")

		// TODO: record fork trigger in Config and persist

		// TODO: do we need to rebind session executions stdio to /dev/null or to files?
		err := ops.Fork(config)
		if err != nil {
			log.Errorf("vmfork failed:%s\n", err)
			// TODO: how do we handle fork failure behaviour at a container level?
			// Does it differ if triggered manually vs via pointcut conditions in a build file
			continue
		}

		// trigger a reload of the configuration
		log.Info("Triggering reload of config after fork")
		reload <- true
	}
}
