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

package tether

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof" // allow enabling pprof in contianerVM
	"os"
	"os/exec"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/cmd/tether/msgs"
	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/system"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/serial"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

const (
	// MaxDeathRecords The maximum number of records to keep for restarting processes
	MaxDeathRecords = 5

	// the length of a truncated ID for use as hostname
	shortLen = 12

	// temp directory to copy existing data to mounts
	bindDir = "/.tether/.bind"
)

var Sys = system.New()
var once sync.Once

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

	// Cancelable context and its cancel func.
	ctx    context.Context
	cancel context.CancelFunc

	incoming chan os.Signal
}

func New(src extraconfig.DataSource, sink extraconfig.DataSink, ops Operations) Tether {
	ctx, cancel := context.WithCancel(context.Background())
	return &tether{
		ops:    ops,
		reload: make(chan bool, 1),
		config: &ExecutorConfig{
			pids: make(map[int]*SessionConfig),
		},
		extensions: make(map[string]Extension),
		src:        src,
		sink:       sink,
		ctx:        ctx,
		cancel:     cancel,
		incoming:   make(chan os.Signal, 32),
	}
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
		log.SetOutput(out)
	}

	t.reload = make(chan bool, 1)
	t.config = &ExecutorConfig{
		pids: make(map[int]*SessionConfig),
	}

	if err := t.childReaper(); err != nil {
		log.Errorf("Failed to start reaper %s", err)
		return err
	}

	if err := t.ops.Setup(t); err != nil {
		log.Errorf("Failed tether setup: %s", err)
		return err
	}

	for name, ext := range t.extensions {
		log.Infof("Starting extension %s", name)
		err := ext.Start()
		if err != nil {
			log.Errorf("Failed to start extension %s: %s", name, err)
			return err
		}
	}

	// #nosec: Expect directory permissions to be 0700 or less
	if err = os.MkdirAll(PIDFileDir(), 0755); err != nil {
		log.Errorf("could not create pid file directory %s: %s", PIDFileDir(), err)
	}

	// Create PID file for tether
	tname := path.Base(os.Args[0])
	err = ioutil.WriteFile(fmt.Sprintf("%s.pid", path.Join(PIDFileDir(), tname)),
		[]byte(fmt.Sprintf("%d", os.Getpid())),
		0644)
	if err != nil {
		log.Errorf("Unable to open PID file for %s : %s", os.Args[0], err)
	}

	// seed the incoming channel once to trigger child reaper. This is required to collect the zombies created by switch-root
	t.triggerReaper()

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

func (t *tether) setLogLevel() {
	// TODO: move all of this into an extension.Pre() block when we move to that model
	// adjust the logging level appropriately
	log.SetLevel(log.InfoLevel)
	// TODO: do not echo application output to console without debug enabled
	serial.DisableTracing()

	if t.config.DebugLevel > 0 {
		log.SetLevel(log.DebugLevel)

		logConfig(t.config)
	}

	if t.config.DebugLevel > 1 {
		serial.EnableTracing()

		log.Info("Launching pprof server on port 6060")
		fn := func() {
			go http.ListenAndServe("0.0.0.0:6060", nil)
		}

		once.Do(fn)
	}
}

func (t *tether) setHostname() error {
	short := t.config.ID
	if len(short) > shortLen {
		short = short[:shortLen]
	}

	if err := t.ops.SetHostname(short, t.config.Name); err != nil {
		// we don't attempt to recover from this - it's a fundamental misconfiguration
		// so just exit
		return fmt.Errorf("failed to set hostname: %s", err)
	}
	return nil
}

func (t *tether) setNetworks() error {
	for _, v := range t.config.Networks {
		if err := t.ops.Apply(v); err != nil {
			return fmt.Errorf("failed to apply network endpoint config: %s", err)
		}
	}
	return nil
}

func (t *tether) setMounts() error {
	for k, v := range t.config.Mounts {
		switch v.Source.Scheme {
		case "label":
			// this could block indefinitely while waiting for a volume to present
			t.ops.MountLabel(context.Background(), v.Source.Path, v.Path)

		case "nfs":
			t.ops.MountTarget(context.Background(), v.Source, v.Path, v.Mode)

		default:
			return fmt.Errorf("unsupported volume mount type for %s: %s", k, v.Source.Scheme)
		}
	}
	return t.populateVolumes()
}

func (t *tether) populateVolumes() error {
	defer trace.End(trace.Begin(fmt.Sprintf("populateVolumes")))
	// skip if no mounts present
	if len(t.config.Mounts) == 0 {
		return nil
	}

	for _, mnt := range t.config.Mounts {
		if mnt.Path == "" {
			continue
		}
		if mnt.CopyMode == executor.CopyNew {
			err := t.ops.CopyExistingContent(mnt.Path)
			if err != nil {
				log.Errorf("error copyExistingContent for mount %s: %+v", mnt.Path, err)
				return err
			}
		}
	}

	return nil
}

func (t *tether) initializeSessions() error {

	maps := map[string]map[string]*SessionConfig{
		"Sessions": t.config.Sessions,
		"Execs":    t.config.Execs,
	}

	// we need to iterate over both sessions and execs
	for name, m := range maps {

		// Iterate over the Sessions and initialize them if needed
		for id, session := range m {
			// make it a func so that we can use defer
			err := func() error {
				session.Lock()
				defer session.Unlock()

				if session.wait != nil {
					log.Warnf("Session %s already initialized", id)
					return nil
				}
				log.Debugf("Initializing session %s", id)

				if session.RunBlock {
					log.Infof("Session %s wants attach capabilities. Creating its channel", id)
					session.ClearToLaunch = make(chan struct{})
				}

				stdout, stderr, err := t.ops.SessionLog(session)
				if err != nil {
					detail := fmt.Errorf("failed to get log writer for session: %s", err)
					session.Started = detail.Error()

					return detail
				}
				session.Outwriter = stdout
				session.Errwriter = stderr
				session.Reader = dio.MultiReader()

				session.wait = &sync.WaitGroup{}
				session.extraconfigKey = name

				return nil
			}()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *tether) reloadExtensions() error {
	// reload the extensions
	for name, ext := range t.extensions {
		log.Debugf("Passing config to %s", name)
		err := ext.Reload(t.config)
		if err != nil {
			return fmt.Errorf("Failed to cleanly reload config for extension %s: %s", name, err)
		}
	}
	return nil
}

func (t *tether) processSessions() error {
	type results struct {
		id    string
		path  string
		err   error
		fatal bool
	}

	// so that we can launch multiple sessions in parallel
	var wg sync.WaitGroup
	// to collect the errors back from them
	resultsCh := make(chan results, len(t.config.Sessions))

	maps := []struct {
		sessions map[string]*SessionConfig
		fatal    bool
	}{
		{t.config.Sessions, true},
		{t.config.Execs, false},
	}

	// we need to iterate over both sessions and execs
	for i := range maps {
		m := maps[i]

		// process the sessions and launch if needed
		for id, session := range m.sessions {
			session.Lock()

			log.Debugf("Processing config for session %s", id)
			var proc = session.Cmd.Process

			// check if session is alive and well
			if proc != nil && proc.Signal(syscall.Signal(0)) == nil {
				log.Debugf("Process for session %s is running (pid: %d)", id, proc.Pid)
				if !session.Active {
					// stop process - for now this doesn't do any staged levels of aggression
					log.Infof("Running session %s has been deactivated (pid: %d)", id, proc.Pid)

					killHelper(session)
				}

				session.Unlock()
				continue
			}

			// if we're not activating this session and it's not running, then skip
			if !session.Active {
				log.Debugf("Skipping inactive session %s", id)
				session.Unlock()
				continue
			}

			// check if session has never been started or is configured for restart
			if proc == nil || session.Restart {
				if proc == nil {
					log.Infof("Launching process for session %s", id)
					log.Debugf("Launch failures are fatal: %t", m.fatal)
				} else {
					session.Diagnostics.ResurrectionCount++

					// FIXME: we cannot have this embedded knowledge of the extraconfig encoding pattern, but not
					// currently sure how to expose it neatly via a utility function
					extraconfig.EncodeWithPrefix(t.sink, session, extraconfig.CalculateKeys(t.config, fmt.Sprintf("%s.%s", session.extraconfigKey, id), "")[0])
					log.Warnf("Re-launching process for session %s (count: %d)", id, session.Diagnostics.ResurrectionCount)
					session.Cmd = *restartableCmd(&session.Cmd)
				}

				wg.Add(1)
				go func(session *SessionConfig) {
					defer wg.Done()
					resultsCh <- results{
						id:    session.ID,
						path:  session.Cmd.Path,
						err:   t.launch(session),
						fatal: m.fatal,
					}
				}(session)

			}
			session.Unlock()
		}
	}

	wg.Wait()
	// close the channel
	close(resultsCh)

	// iterate over the results
	for result := range resultsCh {
		if result.err != nil {
			detail := fmt.Errorf("failed to launch %s for %s: %s", result.path, result.id, result.err)
			if result.fatal {
				log.Error(detail)
				return detail
			}

			log.Warn(detail)
			return nil
		}
	}
	return nil
}

func (t *tether) Start() error {
	defer trace.End(trace.Begin("main tether loop"))

	// do the initial setup and start the extensions
	if err := t.setup(); err != nil {
		log.Errorf("Failed to run setup: %s", err)
		return err
	}
	defer t.cleanup()

	// initial entry, so seed this
	t.reload <- true
	for range t.reload {
		log.Info("Loading main configuration")

		// load the config - this modifies the structure values in place
		extraconfig.Decode(t.src, t.config)

		t.setLogLevel()

		if err := t.setHostname(); err != nil {
			log.Error(err)
			return err
		}

		// process the networks then publish any dynamic data
		if err := t.setNetworks(); err != nil {
			log.Error(err)
			return err
		}
		extraconfig.Encode(t.sink, t.config)

		// setup the firewall
		if err := t.ops.SetupFirewall(t.config); err != nil {
			log.Warnf("Failed to setup firewall: %s", err)
		}

		//process the filesystem mounts - this is performed after networks to allow for network mounts
		if err := t.setMounts(); err != nil {
			log.Error(err)
			return err
		}

		if err := t.initializeSessions(); err != nil {
			log.Error(err)
			return err
		}

		if err := t.reloadExtensions(); err != nil {
			log.Error(err)
			return err
		}

		if err := t.processSessions(); err != nil {
			log.Error(err)
			return err
		}
	}

	log.Info("Finished processing sessions")

	return nil
}

func (t *tether) Stop() error {
	defer trace.End(trace.Begin(""))

	// TODO: kill all the children
	if t.reload != nil {
		close(t.reload)
	}

	// cancel the context to unblock waiters
	t.cancel()
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
// caller needs to hold session Lock
func (t *tether) handleSessionExit(session *SessionConfig) {
	defer trace.End(trace.Begin("handling exit of session " + session.ID))

	log.Debugf("Waiting on session.wait")
	session.wait.Wait()
	log.Debugf("Wait on session.wait completed")

	log.Debugf("Calling wait on cmd")
	if err := session.Cmd.Wait(); err != nil {
		log.Warnf("Wait returned %s", err)
	}

	log.Debugf("Calling close on reader")
	if err := session.Reader.Close(); err != nil {
		log.Warnf("Close for Reader returned %s", err)
	}

	// close down the outputs
	log.Debugf("Calling close on writers")
	if err := session.Outwriter.Close(); err != nil {
		log.Warnf("Close for Outwriter returned %s", err)
	}
	if err := session.Errwriter.Close(); err != nil {
		log.Warnf("Close for Errwriter returned %s", err)
	}

	// close the signaling channel (it is nil for detached sessions) and set it to nil (for restart)
	if session.ClearToLaunch != nil {
		log.Debugf("Calling close chan")
		close(session.ClearToLaunch)
		session.ClearToLaunch = nil
	}

	// Remove associated PID file
	cmdname := path.Base(session.Cmd.Path)
	_ = os.Remove(fmt.Sprintf("%s.pid", path.Join(PIDFileDir(), cmdname)))

	// set the stop time
	session.StopTime = time.Now().UTC().Unix()

	// this returns an arbitrary closure for invocation after the session status update
	f := t.ops.HandleSessionExit(t.config, session)

	extraconfig.EncodeWithPrefix(t.sink, session, extraconfig.CalculateKeys(t.config, fmt.Sprintf("%s.%s", session.extraconfigKey, session.ID), "")[0])

	if f != nil {
		log.Debugf("Calling t.ops.HandleSessionExit")
		f()
	}
}

// launch will launch the command defined in the session.
// This will return an error if the session fails to launch
func (t *tether) launch(session *SessionConfig) error {
	defer trace.End(trace.Begin("launching session " + session.ID))

	session.Lock()
	defer session.Unlock()

	// encode the result whether success or error
	defer func() {
		prefix := extraconfig.CalculateKeys(t.config, fmt.Sprintf("%s.%s", session.extraconfigKey, session.ID), "")[0]
		log.Debugf("Encoding result of launch for session %s under key: %s", session.ID, prefix)
		extraconfig.EncodeWithPrefix(t.sink, session, prefix)
	}()

	if len(session.User) > 0 || len(session.Group) > 0 {
		user, err := getUserSysProcAttr(session.User, session.Group)
		if err != nil {
			log.Errorf("user lookup failed %s:%s, %s", session.User, session.Group, err)
			session.Started = err.Error()
			return err
		}
		session.Cmd.SysProcAttr = user
	}

	session.Cmd.Env = t.ops.ProcessEnv(session.Cmd.Env)
	// Set Std{in|out|err} to nil, we will control pipes
	session.Cmd.Stdin = nil
	session.Cmd.Stdout = nil
	session.Cmd.Stderr = nil

	resolved, err := lookPath(session.Cmd.Path, session.Cmd.Env, session.Cmd.Dir)
	if err != nil {
		log.Errorf("Path lookup failed for %s: %s", session.Cmd.Path, err)
		session.Started = err.Error()
		return err
	}
	log.Debugf("Resolved %s to %s", session.Cmd.Path, resolved)
	session.Cmd.Path = resolved

	// block until we have a connection
	if session.RunBlock && session.ClearToLaunch != nil {
		log.Debugf("Waiting clear signal to launch %s", session.ID)
		select {
		case <-t.ctx.Done():
			log.Warnf("Waiting to launch %s canceled, bailing out", session.ID)
			return nil
		case <-session.ClearToLaunch:
			log.Debugf("Received the clear signal to launch %s", session.ID)
		}
	}

	pid := 0
	// Use the mutex to make creating a child and adding the child pid into the
	// childPidTable appear atomic to the reaper function. Use a anonymous function
	// so we can defer unlocking locally
	// logging is done after the function to keep the locked time as low as possible
	err = func() error {
		t.config.pidMutex.Lock()
		defer t.config.pidMutex.Unlock()

		if !session.Tty {
			err = establishNonPty(session)
		} else {
			err = establishPty(session)
		}
		if err != nil {
			return err
		}

		pid = session.Cmd.Process.Pid
		t.config.pids[pid] = session

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

	// Write the PID to the associated PID file
	cmdname := path.Base(session.Cmd.Path)
	err = ioutil.WriteFile(fmt.Sprintf("%s.pid", path.Join(PIDFileDir(), cmdname)),
		[]byte(fmt.Sprintf("%d", pid)),
		0644)
	if err != nil {
		log.Errorf("Unable to write PID file for %s: %s", cmdname, err)
	}
	log.Debugf("Launched command with pid %d", pid)

	return nil
}

func logConfig(config *ExecutorConfig) {
	// just pretty print the json for now
	log.Info("Loaded executor config")

	// figure out the keys to filter
	keys := make(map[string]interface{})
	if config.DebugLevel < 2 {
		for _, f := range []string{
			"Sessions.*.Cmd.Args",
			"Sessions.*.Cmd.Args.*",
			"Sessions.*.Cmd.Env",
			"Sessions.*.Cmd.Env.*",
			"Key"} {
			for _, k := range extraconfig.CalculateKeys(config, f, "") {
				keys[k] = nil
			}
		}
	}

	sink := map[string]string{}
	extraconfig.Encode(
		func(k, v string) error {
			if _, ok := keys[k]; !ok {
				sink[k] = v
			}

			return nil
		},
		config,
	)

	for k, v := range sink {
		log.Debugf("%s: %s", k, v)
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
	for range incoming {
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

// Config interface
func (t *tether) UpdateNetworkEndpoint(e *NetworkEndpoint) error {
	defer trace.End(trace.Begin("tether.UpdateNetworkEndpoint"))

	if e == nil {
		return fmt.Errorf("endpoint must be specified")
	}

	if _, ok := t.config.Networks[e.Network.Name]; !ok {
		return fmt.Errorf("network endpoint not found in config")
	}

	t.config.Networks[e.Network.Name] = e
	return nil
}

func (t *tether) Flush() error {
	defer trace.End(trace.Begin("tether.Flush"))

	extraconfig.Encode(t.sink, t.config)
	return nil
}

func PIDFileDir() string {
	return path.Join(Sys.Root, pidFilePath)
}

// killHelper was pulled from toolbox, and that variant should be directed at this
// one eventually
func killHelper(session *SessionConfig) error {
	sig := new(msgs.SignalMsg)
	name := session.StopSignal
	if name == "" {
		name = string(ssh.SIGTERM)
	}

	err := sig.FromString(name)
	if err != nil {
		return err
	}

	num := syscall.Signal(sig.Signum())

	log.Infof("sending signal %s (%d) to %s", sig.Signal, num, session.ID)

	if err := session.Cmd.Process.Signal(num); err != nil {
		return fmt.Errorf("failed to signal %s: %s", session.ID, err)
	}

	return nil
}
