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

package handlers

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/kr/pty"
	"github.com/vmware/vic/tether"
	"golang.org/x/crypto/ssh"
)

func (ch *SessionHandler) AssignPty() {
	ch.assignTty = true
}

func (ch *SessionHandler) ResizePty(winSize *tether.WindowChangeMsg) error {
	if !ch.assignTty {
		// not sure what we should return in this circumstance but we cannot act on it
		log.Println("Received windows resize request for non-pty session")
		return nil
	}

	ws := &winsize{uint16(winSize.Rows), uint16(winSize.Columns), uint16(winSize.Width_px), uint16(winSize.Height_px)}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		ch.pty.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func (ch *SessionHandler) Signal(sig ssh.Signal) error {
	log.Printf("Sending signal %s to main process %s (%d)\n", string(sig), ch.cmd.Path, ch.cmd.Process.Pid)
	err := ch.cmd.Process.Signal(syscall.Signal(tether.Signals[sig]))
	if err != nil {
		log.Printf("Failed to dispatch signal to process: %s\n", err)
	}
	return err
}

func (ch *SessionHandler) Exec(command string, args []string, config map[string]string) (ok bool, payload []byte) {
	// strip quotes from the args if they are first AND last positions in an arg element
	for k, v := range args {
		if v[0] == '"' && v[len(v)-1] == '"' {
			args[k] = v[1 : len(v)-1]
		}
	}

	// TODO: figure out how we're going to specify user and pass all the settings along
	// in the meantime, hardcode HOME to /root

	// we process the environment before creating the command so that LookPath has a chance to do it's thing
	env := []string{"HOME=/root"}
	for k, v := range ch.env {
		val := k + "=" + v

		switch k {
		case "HOME":
			env[0] = val
		case "PATH":
			os.Setenv(k, v)
			fallthrough
		default:
			env = append(env, val)
		}
	}

	// TODO: fork & exec with chroot (legacy support)
	// or use mnt namespace, or similar
	ch.cmd = exec.Command(command, args...)

	// set the working directory if specified
	if val, ok := config["dir"]; ok {
		ch.cmd.Dir = val
	}

	// publish the environment we constructed above
	ch.cmd.Env = env

	// Use the mutex to make creating a child and adding the child pid into the
	// childPidTable appear atomic to the reaper function. Use a anonymous function
	// so we can defer unlocking locally
	statusChan, err := func() (statusChan chan syscall.WaitStatus, err error) {
		err = nil
		childPidTableMutex.Lock()
		defer childPidTableMutex.Unlock()

		log.Printf("Exec'ing command %+q\n", ch.cmd.Args)
		if !ch.assignTty {
			ch.cmd.Stdin = *ch.channel
			ch.cmd.Stdout = *ch.channel
			ch.cmd.Stderr = (*ch.channel).Stderr()
			err = ch.cmd.Start()
		} else {
			ch.pty, err = pty.Start(ch.cmd)
		}

		if err == nil {
			// ChildReaper will use this channel to inform us the wait status of the child.
			statusChan = make(chan syscall.WaitStatus, 1)
			childPidTable[ch.cmd.Process.Pid] = statusChan
		}

		return statusChan, err
	}()

	if err == nil {
		log.Printf("Started child process %d\n", ch.cmd.Process.Pid)
	} else {
		log.Printf("Failed to start process: %s\n", err)
		ch.pendingFn = func() {
			log.Println("Closing channel after failed exec")
			(*ch.channel).CloseWrite()
			(*ch.channel).Close()
		}
		return false, []byte(err.Error())
	}

	// we need to send a reply to the exec request now, monitoring of the process must continue
	// in a different thread

	ch.pendingFn = func() {
		// wait for the exec reply to be sent
		if ch.pty != nil {
			ch.waitGroup.Add(1)
			go func() { io.Copy(*ch.channel, ch.pty); ch.waitGroup.Done() }()
			// we shouldn't wait on stdin to close or we'll be here forever
			go func() { io.Copy(ch.pty, *ch.channel) }()

			log.Println("Waiting for I/O streams to close")
			ch.waitGroup.Wait()
		}

		// shouldn't call wait until all reads from stdout/stderr have completed
		log.Println("Waiting for command to complete")
		// ChildReaper waits for all child processes, and it sends the status
		// through statusChan.
		var waitStatus syscall.WaitStatus
		for waitStatus = range statusChan {
			if waitStatus.Exited() || waitStatus.Signaled() {
				break
			} else {
				// STOP, CONT or TRAP occurs when program is traced in a debugger
			}
		}
		log.Printf("Got wait status of child %d: %d\n",
			ch.cmd.Process.Pid, waitStatus)
		// Now remove this child pid from the childPidTable
		childPidTableMutex.Lock()
		delete(childPidTable, ch.cmd.Process.Pid)
		childPidTableMutex.Unlock()
		// close the channel after the pid is removed
		close(statusChan)

		exitStatus := uint32(waitStatus.ExitStatus())
		log.Printf("Process completed with exit status %d\n", exitStatus)
		// See "man 2 waitpid" for explanation
		if waitStatus.Exited() {
			log.Println("Command exited normally")
		} else if waitStatus.Signaled() {
			log.Printf("Command terminated with Signal: %v\n",
				waitStatus.Signal())
			if waitStatus.CoreDump() {
				// SIGSEGV
				log.Println("Command core dumped")
			}
		} else if waitStatus.Stopped() {
			log.Printf("Command Stopped by Signal: %v\n",
				waitStatus.StopSignal())
		} else if waitStatus.Continued() {
			log.Println("Command continued")
		} else {
			log.Printf("Trap cause = %d\n", waitStatus.TrapCause())
		}

		// ensure that changes are flushed to disk before we report exit
		syscall.Sync()

		if err := (*ch.channel).CloseWrite(); err != nil {
			log.Println("Error sending channel EOF: ", err)
		}

		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, exitStatus)
		if _, err := (*ch.channel).SendRequest("exit-status", false, bytes); err != nil {
			log.Println("Error sending exit status: ", err)
		}

		if err := (*ch.channel).Close(); err != nil {
			log.Println("Error sending channel close: ", err)
		}

		log.Println("Returned exit status and closed channel")
	}

	// send the immediate reply to the exec request
	log.Println("Started process successfully")
	return true, nil
}
