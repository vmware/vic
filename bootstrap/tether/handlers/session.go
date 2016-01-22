package handlers

import (
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"enatai-gerrit.eng.vmware.com/bonneville-container/tether"

	"golang.org/x/crypto/ssh"
)

var (
	// Set of child PIDs created by us.
	childPidTable = make(map[int]chan syscall.WaitStatus)
	// Exclusive access to childPidTable
	childPidTableMutex = &sync.Mutex{}
)

type SessionHandler struct {
	tether.GlobalContext
	// TODO: add some locking in here if non exec requests can touch exec or shell related items
	cmd       *exec.Cmd
	channel   *ssh.Channel
	env       map[string]string
	assignTty bool
	pty       *os.File
	waitGroup sync.WaitGroup
	pendingFn func()
}

// The syscall struct
type winsize struct {
	ws_row    uint16
	ws_col    uint16
	ws_xpixel uint16
	ws_ypixel uint16
}

func (ch *SessionHandler) SetChannel(channel *ssh.Channel) {
	ch.channel = channel
}

func (ch *SessionHandler) Setenv(name, value string) (ok bool, payload []byte) {
	ch.env[name] = value
	log.Printf("Set environment variable: %s=%s\n", name, value)

	return true, nil
}

func (ch *SessionHandler) Shell() (ok bool, payload []byte) {
	//TODO: implement
	return false, nil
}

func (ch *SessionHandler) Kill() error {
	return ch.cmd.Process.Kill()
}

func (ch *SessionHandler) GetPendingWork() func() {
	return ch.pendingFn
}

func (ch *SessionHandler) ClearPendingWork() {
	ch.pendingFn = nil
}

func GetChildPid(pid int) (chan syscall.WaitStatus, bool) {
	childPidTableMutex.Lock()
	defer childPidTableMutex.Unlock()
	c, ok := childPidTable[pid]
	return c, ok
}
