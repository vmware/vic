// +build linux
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"

	winserial "github.com/tarm/serial"

	"enatai-gerrit.eng.vmware.com/bonneville-container/tether"
	"enatai-gerrit.eng.vmware.com/bonneville-container/tether/serial"
)

var (
	tetherKey string
	port      string
	hup       chan os.Signal

	// Set of child PIDs created by us.
	childPidTable = make(map[int]chan syscall.WaitStatus)
	// Exclusive access to childPidTable
	childPidTableMutex = &sync.Mutex{}
)

type GlobalHandler struct {
	id string
}

type SessionHandler struct {
	*GlobalHandler
	// TODO: add some locking in here if non exec requests can touch exec or shell related items
	cmd       *exec.Cmd
	channel   *ssh.Channel
	env       map[string]string
	assignTty bool
	pty       *os.File
	waitGroup sync.WaitGroup
	pendingFn func()
}

type NamedPort struct {
	*winserial.Port

	config winserial.Config
	fd     uintptr
}

func (p *NamedPort) Name() string {
	return p.config.Name
}

func (p *NamedPort) Fd() uintptr {
	return p.fd
}

func OpenPort(config *winserial.Config) (*NamedPort, error) {
	port, err := winserial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	return &NamedPort{Port: port, config: *config, fd: 0}, nil
}

func privateKey(file string) ssh.Signer {
	privateBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("failed to load private key: %v", tetherKey)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", tetherKey)
	}
	return private
}

// Hopefully this gets rolled into go.sys at some point
func Mkdev(majorNumber int64, minorNumber int64) int {
	return int((majorNumber << 8) | (minorNumber & 0xff) | ((minorNumber & 0xfff00) << 12))
}

func init() {
	flag.StringVar(&tetherKey, "key", "/Windows/tether-init/init_key", "tetherd control channel private key")
	flag.StringVar(&port, "port", "COM2", "com port to use for control")
}

// load the ID from the file system
func id() string {
	id, err := ioutil.ReadFile("/Windows/tether-init/docker-id")
	if err != nil {
		log.Fatalf("failed to read ID from file: %s", err)
	}

	// strip any trailing garbage
	if len(id) > 64 {
		return string(id[:64])
	}
	return string(id)
}

func main() {
	// get the windows service logic running so that we can play well in that mode
	runService("VMware Tether", false)

	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	// Parse ssh private key
	private := privateKey(tetherKey)

	// HACK: workaround file descriptor conflict in pipe2 return from the exec.Command.Start
	_, _, _ = os.Pipe()

	/*
		pid := os.Getpid()

		// register the signal handling that we use to determine when the tether should initialize runtime data
		hup := make(chan os.Signal, 1)
		signal.Notify(hup, syscall.SIGHUP)
		syscall.Kill(pid, syscall.SIGHUP)
	*/

	for {
		// block until HUP
		/*
			log.Printf("Waiting for HUP signal - blocking tether initialization")
			<-hup
			log.Printf("Received HUP signal - initializing tether")
		*/

		c := &winserial.Config{Name: port, Baud: 115200}
		s, err := OpenPort(c)
		if err != nil {
			log.Printf("failed to open com1 for ssh server: %s", err)
			return
		}

		defer s.Close()

		log.Printf("creating raw connection from %s (fd=%d)\n", s.Name(), s.Fd())
		conn, err := serial.NewTypedConn(s, "file")

		if err != nil {
			log.Printf("failed to create raw connection from ttyS0 file handle: %s", err)
			return
		}

		// Delete ourselves from the filesystem so we're not polluting the container
		// TODO: the deletion routine should be a closure passed to tether as we don't know what filesystem
		// access is still needed for basic setup

		handler := GlobalHandler{
			id: id(),
		}

		// ensure we don't have significant obsolete data built up
		s.Port.Flush()

		// HACK: currently RawConn dosn't implement timeout
		serial.HandshakeServer(conn, time.Duration(10*time.Second))

		log.Println("creating ssh connection over %s", s.Name())
		tether.StartTether(conn, private, &handler)

		s.Close()
	}
}
