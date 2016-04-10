// This is the set of commands that the tether will attempt to invoke

package tether

import "golang.org/x/crypto/ssh"

type GlobalContext interface {
	// global commands - affects all sessions

	// kicks off a server connection manager if it's needed
	StartConnectionManager(conn *ssh.ServerConn)

	// Return the container ID (long form)
	ContainerId() string

	// bring up the network interface
	// cidr is an IP addresses specified in CIDR format (x.x.x.x/mask)
	// gateway is the default gateway for external traffic
	StaticIPAddress(cidr, gateway string) error

	// bring up the network interface
	DynamicIPAddress() (string, error)

	// mount a labeled volume on the target directory
	MountLabel(label, target string) error

	// force the filesystem to commit reads (if necessary on the target)
	// this is called prior to power off
	Sync()

	// session commands - affects single session

	// Used to create a session channel handler
	NewSessionContext() SessionContext
}

type SessionContext interface {
	// session commands - affects a singular session

	// Sets the channel for the session. This is the i/o channel the pty
	// will be bound to if specified
	SetChannel(channel *ssh.Channel)

	// Set an environment variable for any subsequent Exec
	// Returns true if successful, false otherwise
	// payload contains details in error case
	Setenv(name, value string) (bool, []byte)

	// Requests assignment of a PTY
	AssignPty()

	// SSH window resize command - see RFC4254
	ResizePty(winSize *WindowChangeMsg) error

	// Helper call to execute a shell of default type for the target system
	// This conceptually wraps Exec, so all Exec related behaviours should be preserved
	// Returns true if successful, false otherwise
	// Returns a payload that will be returned to the remote caller
	Shell() (bool, []byte)

	// Signal the executing process - meaningless if Exec has not yet been called
	Signal(signal ssh.Signal) error

	// Forcibly terminate the executing process - meaning if Exec has not yet been called
	Kill() error

	// Exec should prep the execution synchronously, but place the actual execution into a
	// closure presented by GetPendingWork
	Exec(command string, args []string, config map[string]string) (ok bool, payload []byte)

	// force the filesystem to commit reads (if necessary on the target)
	// this is called prior to power off
	Sync()

	// Retrieve closure for any pending work - this is necessary as data cannot be returned
	// via ssh before request replys are sent so exec, et al, must be async.
	// This will be called by the tether in a goroutine after replys have been processed
	GetPendingWork() func()

	// Called after GetPendingWork closure has been invoked
	ClearPendingWork()
}
