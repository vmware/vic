// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package process

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	"github.com/google/uuid"

	tether "github.com/vmware/vic/lib/tether-ng"
	"github.com/vmware/vic/lib/tether-ng/types"
)

const (
	//https://github.com/golang/go/blob/master/src/syscall/zerrors_linux_arm64.go#L919
	SetChildSubreaper = 0x24
)

type Process struct {
	uuid uuid.UUID
	ctx  context.Context

	config types.ExecutorConfig

	sm       sync.Mutex
	sessions map[string]*types.Session

	m    sync.RWMutex
	pids map[int]*types.Session

	incoming chan os.Signal

	tether.Releaser
	tether.Interactor
}

func NewProcess(ctx context.Context) *Process {
	return &Process{
		uuid:     uuid.New(),
		ctx:      ctx,
		sessions: make(map[string]*types.Session),
		pids:     make(map[int]*types.Session),
		incoming: make(chan os.Signal, 32),
	}
}

func (p *Process) SetReleaser(ctx context.Context, releaser tether.Releaser) {
	p.Releaser = releaser
}

func (p *Process) SetInteractor(ctx context.Context, interactor tether.Interactor) {
	p.Interactor = interactor
}

func (p *Process) Configure(ctx context.Context, config *types.ExecutorConfig) error {
	// create our own copy
	p.config = *config

	for k, i := range p.config.Sessions {
		session := &types.Session{
			ID: k,
			Cmd: exec.Cmd{
				Path:   i.Cmd[0],
				Args:   i.Cmd,
				Env:    i.Env,
				Dir:    i.WorkingDir,
				Stdin:  nil,
				Stdout: nil,
				Stderr: nil,
			},
			SessionConfig: i,
		}
		p.sessions[k] = session

	}
	return nil
}

func (p *Process) Start(ctx context.Context) error {
	p.Reap(ctx)

	for _, i := range p.sessions {
		go func(i *types.Session) {
			p.StartSession(ctx, i)
		}(i)
	}

	return nil
}

func (p *Process) Stop(ctx context.Context) error {
	p.StopReaper(ctx)

	return nil
}

func (p *Process) StartSession(ctx context.Context, i *types.Session) error {
	in := make(chan *types.Session, 1)

	// send session over and receive Done channel
	in <- i
	if i.Tty {
		fmt.Printf("PseudoTerminal %s\n", i.ID)
		i.Done = p.PseudoTerminal(ctx, in)
	} else {
		fmt.Printf("Non-PseudoTermina %s\n", i.ID)
		i.Done = p.NonInteract(ctx, in)
	}

	if i.RunBlock {
		fmt.Printf("Runblock\n")
		requestChan := make(chan chan struct{})
		go p.Release(ctx, requestChan)
		p.Wait(ctx, requestChan)
	}

	p.m.Lock()
	i.Cmd.Start()
	p.pids[i.Cmd.Process.Pid] = i
	p.m.Unlock()

	return nil
}

func (p *Process) StopSession(ctx context.Context, i *types.Session) error {
	in := make(chan *types.Session, 1)

	// make sure reads/writes finishes before we call cmd.Wait
	<-i.Done
	i.Cmd.Wait()

	// send session over and wait till Close returns
	in <- i
	<-p.Close(ctx, in)

	fmt.Printf("Deleting %s\n", i.ID)

	p.sm.Lock()
	delete(p.sessions, i.ID)
	p.sm.Unlock()

	return nil
}

func (p *Process) PidToSession(ctx context.Context, pid int) *types.Session {
	p.m.Lock()
	defer p.m.Unlock()

	session, ok := p.pids[pid]
	if ok {
		return session
	}
	return nil
}

func (p *Process) UUID(ctx context.Context) uuid.UUID { return p.uuid }

func (p *Process) Running(ctx context.Context, sessionID string) bool {
	p.sm.Lock()
	_, ok := p.sessions[sessionID]
	p.sm.Unlock()

	return ok
}

func (p *Process) Kill(ctx context.Context, sessionID string) error { return nil }

func (p *Process) Wait(ctx context.Context, in <-chan chan struct{}) {
	fmt.Printf("WAITING\n")
	if in != nil {
		<-in
	}
}

func (p *Process) Reap(ctx context.Context) error {
	signal.Notify(p.incoming, syscall.SIGCHLD)

	if _, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, SetChildSubreaper, uintptr(1), 0); err != 0 {
		return err
	}

	flag := syscall.WNOHANG | syscall.WUNTRACED | syscall.WCONTINUED

	go func() {
		var status syscall.WaitStatus

		for range p.incoming {
			// wrap it so that we can use recover
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Fprintf(os.Stderr, "Recovered in childReaper %s", debug.Stack())
					}
				}()

				// reap until no more children to process
				for {
					fmt.Printf("Inspecting children with status change\n")

					select {
					case <-ctx.Done():
						fmt.Printf("Someone called shutdown, returning from child reaper\n")
						return
					default:
					}

					pid, err := syscall.Wait4(-1, &status, flag, nil)
					// pid 0 means no processes wish to report status
					if pid == 0 || err == syscall.ECHILD {
						fmt.Printf("No more child processes to reap\n")
						break
					}

					if err != nil {
						fmt.Printf("Wait4 got error: %v\n", err)
						break
					}

					if !status.Exited() && !status.Signaled() {
						fmt.Printf("Received notifcation about non-exit status change for %d: %d\n", pid, status)
						continue
					}

					fmt.Printf("PID is %d\n", pid)

					session := p.PidToSession(ctx, pid)
					if session == nil {
						// This is an adopted zombie. The Wait4 call already clean it up from the kernel
						fmt.Printf("Reaped zombie process PID %d", pid)
						continue
					}
					session.Lock()
					session.ExitStatus = status.ExitStatus()
					p.StopSession(ctx, session)
					session.Unlock()
				}
			}()
		}
		fmt.Printf("Stopped reaping child processes\n")
	}()

	return nil
}

func (p *Process) StopReaper(ctx context.Context) error {
	signal.Reset(syscall.SIGCHLD)
	close(p.incoming)

	return nil
}
