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
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stringid"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

const (
	// The maximum number of records to keep for restarting processes
	MaxDeathRecords = 5
)

type tether struct {
	// the implementation to use for tailored operations
	ops Operations

	// the reload channel is used to block reloading of the config
	reload chan bool

	// config holds the main configuration for the executor
	config *ExecutorConfig

	// a set of extensions that get to operate on the config
	extensions map[string]Extension

	src  extraconfig.DataSource
	sink extraconfig.DataSink

	incoming chan os.Signal
}

func New(src extraconfig.DataSource, sink extraconfig.DataSink, ops Operations) Tether {
	t := &tether{
		ops:    ops,
		reload: make(chan bool, 1),
		config: &ExecutorConfig{
			pids: make(map[int]*SessionConfig),
		},
		extensions: make(map[string]Extension),
		src:        src,
		sink:       sink,
		incoming:   make(chan os.Signal, 10),
	}

	// HACK: workaround file descriptor conflict in pipe2 return from the exec.Command.Start
	// it's not clear whether this is a cross platform issue, or still an issue as of this commit
	// keeping it until there's time to verify and fix properly with a Go PR.
	_, _, _ = os.Pipe()

	return t
}

// removeChildPid is a synchronized accessor for the pid map the deletes the entry and returns the value
func (t *tether) removeChildPid(pid int) (*SessionConfig, bool) {
	t.config.pidMutex.Lock()
	defer t.config.pidMutex.Unlock()

	session, ok := t.config.pids[pid]
	delete(t.config.pids, pid)
	return session, ok
}

// lenChildPid returns the number of entries
func (t *tether) lenChildPid() int {
	t.config.pidMutex.Lock()
	defer t.config.pidMutex.Unlock()

	return len(t.config.pids)
}

func (t *tether) setup() error {
	defer trace.End(trace.Begin("main tether setup"))

	// set up tether logging destination
	out, err := t.ops.Log()
	if err != nil {
		log.Errorf("failed to open tether log: %s", err)
		return err
	}
	if out != nil {
		log.SetOutput(io.MultiWriter(out, os.Stdout))
	}

	t.reload = make(chan bool, 1)
	t.config = &ExecutorConfig{
		pids: make(map[int]*SessionConfig),
	}

	t.childReaper()

	t.ops.Setup(t)

	for name, ext := range t.extensions {
		log.Infof("Starting extension %s", name)
		err := ext.Start()
		if err != nil {
			log.Errorf("Failed to start extension %s: %s", name, err)
			return err
		}
	}

	return nil
}

func (t *tether) cleanup() {
	defer trace.End(trace.Begin("main tether cleanup"))

	// stop child reaping
	t.stopReaper()

	// stop the extensions first as they may use the config
	for name, ext := range t.extensions {
		log.Infof("Stopping extension %s", name)
		err := ext.Stop()
		if err != nil {
			log.Warnf("Failed to cleanly stop extension %s", name)
		}
	}

	// return logging to standard location
	log.SetOutput(os.Stdout)

	// perform basic cleanup
	t.reload = nil
	t.ops.Cleanup()
}

func (t *tether) Start() error {
	defer trace.End(trace.Begin("main tether loop"))

	t.setup()
	defer t.cleanup()

	// initial entry, so seed this
	t.reload <- true
	for _ = range t.reload {
		log.Info("Loading main configuration")
		// load the config - this modifies the structure values in place
		extraconfig.Decode(t.src, t.config)
		logConfig(t.config)

		if err := t.ops.SetHostname(stringid.TruncateID(t.config.ID), t.config.Name); err != nil {
			detail := fmt.Sprintf("failed to set hostname: %s", err)
			log.Error(detail)
			// we don't attempt to recover from this - it's a fundemental misconfiguration
			// so just exit
			return errors.New(detail)
		}

		// process the networks then publish any dynamic data
		for _, v := range t.config.Networks {
			if err := t.ops.Apply(v); err != nil {
				detail := fmt.Sprintf("failed to apply network endpoint config: %s", err)
				log.Error(detail)
				return errors.New(detail)
			}
		}
		extraconfig.Encode(t.sink, t.config)

		//process the filesystem mounts - this is performed after networks to allow for network mounts
		for k, v := range t.config.Mounts {
			if v.Source.Scheme != "label" {
				detail := fmt.Sprintf("unsupported volume mount type for %s: %s", k, v.Source.Scheme)
				log.Error(detail)
				return errors.New(detail)
			}

			// this could block indefinitely while waiting for a volume to present
			t.ops.MountLabel(v.Source.Path, v.Path, context.Background())
		}

		// process the sessions and launch if needed
		for id, session := range t.config.Sessions {
			log.Debugf("Processing config for session %s", session.ID)
			var proc = session.Cmd.Process

			// check if session is alive and well
			if proc != nil && proc.Signal(syscall.Signal(0)) == nil {
				log.Debugf("Process for session %s is already running (pid: %d)", session.ID, proc.Pid)
				continue
			}

			// check if session has never been started or is configured for restart
			if proc == nil || session.Restart {
				if proc == nil {
					log.Infof("Launching process for session %s", session.ID)
				} else {
					session.Diagnostics.ResurrectionCount++

					// FIXME: we cannot have this embedded knowledge of the extraconfig encoding pattern, but not
					// currently sure how to expose it neatly via a utility function
					extraconfig.EncodeWithPrefix(t.sink, session, fmt.Sprintf("guestinfo..sessions|%s", session.ID))
					log.Warnf("Re-launching process for session %s (count: %d)", session.ID, session.Diagnostics.ResurrectionCount)
					session.Cmd = *restartableCmd(&session.Cmd)
				}

				err := t.launch(session)
				if err != nil {
					detail := fmt.Sprintf("failed to launch %s for %s: %s", session.Cmd.Path, id, err)
					log.Error(detail)

					// TODO: check if failure to launch this is fatal to everything in this containerVM
					// 		for now failure to launch at all is terminal
					return errors.New(detail)
				}

				continue
			}

			log.Warnf("Process for session %s has exited (%d) and is not configured for restart", session.ID, session.ExitStatus)
		}

		for name, ext := range t.extensions {
			log.Info("Passing config to " + name)
			err := ext.Reload(t.config)
			if err != nil {
				log.Errorf("Failed to cleanly reload config for extension %s: %s", name, err)
				return err
			}
		}
	}

	return nil
}

func (t *tether) Stop() error {
	defer trace.End(trace.Begin(""))

	// TODO: kill all the children
	if t.reload != nil {
		close(t.reload)
	}

	return nil
}

func (t *tether) Reload() {
	log.Infof("Reload triggered")

	t.reload <- true
}

func (t *tether) Register(name string, extension Extension) {
	log.Infof("Registering tether extension " + name)

	t.extensions[name] = extension
}

// handleSessionExit processes the result from the session command, records it in persistent
// maner and determines if the Executor should exit
func (t *tether) handleSessionExit(session *SessionConfig) {
	defer trace.End(trace.Begin("handling exit of session " + session.ID))

	// close down the IO
	session.Reader.Close()
	// live.outwriter.Close()
	// live.errwriter.Close()

	// flush session log output

	// Log a death record trimming records if need be
	logs := session.Diagnostics.ExitLogs
	logCount := len(logs)
	if logCount >= MaxDeathRecords {
		logs = logs[logCount-MaxDeathRecords+1:]
	}

	session.Diagnostics.ExitLogs = append(logs, &metadata.ExitLog{
		Time:       time.Now(),
		ExitStatus: session.ExitStatus,
		// We don't have any message for now
	})

	// this returns an arbitrary closure for invocation after the session status update
	f := t.ops.HandleSessionExit(t.config, session)

	log.Infof("%s exit code: %d", session.ID, session.ExitStatus)
	// record exit status
	// FIXME: we cannot have this embedded knowledge of the extraconfig encoding pattern, but not
	// currently sure how to expose it neatly via a utility function
	extraconfig.EncodeWithPrefix(t.sink, session, fmt.Sprintf("guestinfo..sessions|%s", session.ID))

	if f != nil {
		f()
	}
}

// launch will launch the command defined in the session.
// This will return an error if the session fails to launch
func (t *tether) launch(session *SessionConfig) error {
	defer trace.End(trace.Begin("launching session " + session.ID))

	// encode the result whether success or error
	defer func() {
		extraconfig.EncodeWithPrefix(t.sink, session, fmt.Sprintf("guestinfo..sessions|%s", session.ID))
	}()

	logwriter, err := t.ops.SessionLog(session)
	if err != nil {
		detail := fmt.Sprintf("failed to get log writer for session: %s", err)
		log.Error(detail)
		session.Started = detail

		return errors.New(detail)
	}

	// we store these outside of the session.Cmd struct so that there's consistent
	// handling between tty & non-tty paths
	session.Outwriter = logwriter
	session.Errwriter = logwriter
	session.Reader = dio.MultiReader()

	session.Cmd.Env = t.ops.ProcessEnv(session.Cmd.Env)
	session.Cmd.Stdout = session.Outwriter
	session.Cmd.Stderr = session.Errwriter
	session.Cmd.Stdin = session.Reader

	resolved, err := lookPath(session.Cmd.Path, session.Cmd.Env, session.Cmd.Dir)
	if err != nil {
		log.Errorf("Path lookup failed for %s: %s", session.Cmd.Path, err)
		session.Started = err.Error()
		return err
	}
	log.Debugf("Resolved %s to %s", session.Cmd.Path, resolved)
	session.Cmd.Path = resolved

	// Use the mutex to make creating a child and adding the child pid into the
	// childPidTable appear atomic to the reaper function. Use a anonymous function
	// so we can defer unlocking locally
	// logging is done after the function to keep the locked time as low as possible
	err = func() error {
		t.config.pidMutex.Lock()
		defer t.config.pidMutex.Unlock()

		if !session.Tty {
			err = session.Cmd.Start()
		} else {
			err = establishPty(session)
		}

		if err != nil {
			return err
		}

		// ChildReaper will use this channel to inform us the wait status of the child.
		t.config.pids[session.Cmd.Process.Pid] = session

		return nil
	}()

	if err != nil {
		detail := fmt.Sprintf("failed to start container process: %s", err)
		log.Error(detail)

		// Set the Started key to the undecorated error message
		session.Started = err.Error()

		return errors.New(detail)
	}

	// Set the Started key to "true" - this indicates a successful launch
	session.Started = "true"
	log.Debugf("Launched command with pid %d", session.Cmd.Process.Pid)

	return nil
}

func logConfig(config *ExecutorConfig) {
	// just pretty print the json for now
	log.Info("Loaded executor config")

	// TODO: investigate whether it's the govmomi types package cause the binary size
	// inflation - if so we need an alternative approach here or in extraconfig
	if log.GetLevel() == log.DebugLevel && config.DebugLevel > 1 {
		sink := map[string]string{}
		extraconfig.Encode(extraconfig.MapSink(sink), config)

		for k, v := range sink {
			log.Debugf("%s: %s", k, v)
		}
	}
}

func (t *tether) forkHandler() {
	defer trace.End(trace.Begin("start fork trigger handler"))

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in forkHandler", r)
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
		err := t.ops.Fork()
		if err != nil {
			log.Errorf("vmfork failed:%s\n", err)
			// TODO: how do we handle fork failure behaviour at a container level?
			// Does it differ if triggered manually vs via pointcut conditions in a build file
			continue
		}

		// trigger a reload of the configuration
		log.Info("Triggering reload of config after fork")
		t.reload <- true
	}
}

// restartableCmd takes the Cmd struct for a process that has been run and creates a new
// one that can be lauched again. Stdin/out will need to be set up again.
func restartableCmd(cmd *exec.Cmd) *exec.Cmd {
	return &exec.Cmd{
		Path:        cmd.Path,
		Args:        cmd.Args,
		Env:         cmd.Env,
		Dir:         cmd.Dir,
		ExtraFiles:  cmd.ExtraFiles,
		SysProcAttr: cmd.SysProcAttr,
	}
}

// ConfigSink interface
func (t *tether) WriteKey(key string, value interface{}) error {
	extraconfig.EncodeWithPrefix(t.sink, value, key)
	return nil
}
