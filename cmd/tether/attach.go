package main

import (
	"errors"
	"fmt"
	"net"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/dio"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

var (
	Signals = map[ssh.Signal]int{
		ssh.SIGABRT: 6,
		ssh.SIGALRM: 14,
		ssh.SIGFPE:  8,
		ssh.SIGHUP:  1,
		ssh.SIGILL:  4,
		ssh.SIGINT:  2,
		ssh.SIGKILL: 9,
		ssh.SIGPIPE: 13,
		ssh.SIGQUIT: 3,
		ssh.SIGSEGV: 11,
		ssh.SIGTERM: 15,
		ssh.SIGUSR1: 10,
		ssh.SIGUSR2: 12,
	}
)

// PtyRequestMsg the RFC4254 struct
type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

// WindowChangeMsg the RFC4254 struct
type WindowChangeMsg struct {
	Columns  uint32
	Rows     uint32
	WidthPx  uint32
	HeightPx uint32
}

type signalMsg struct {
	Signal string
}

var attachServer *attachServerSSH

// conn is held directly as it's how we stop the attach server
type attachServerSSH struct {
	conn   *net.Conn
	config *ssh.ServerConfig

	enabled bool
}

// start is not thread safe with stop
func (t *attachServerSSH) start() error {
	if t.enabled {
		return nil
	}

	// don't assume that the key hasn't changed
	pkey, err := ssh.ParsePrivateKey(config.Key)
	if err != nil {
		detail := fmt.Sprintf("failed to load key for attach: %s", err)
		log.Error(detail)
		return errors.New(detail)
	}

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	// TODO: update this with generated credentials for the appliance
	t.config = &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if c.User() == "daemon" {
				return &ssh.Permissions{}, nil
			}
			return nil, fmt.Errorf("expected daemon user")
		},
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == "daemon" {
				return &ssh.Permissions{}, nil
			}
			return nil, fmt.Errorf("expected daemon user")
		},
		NoClientAuth: true,
	}
	t.config.AddHostKey(pkey)

	go t.run()

	return nil
}

// stop is not thread safe with start
func (t *attachServerSSH) stop() {
	if !t.enabled {
		return
	}

	t.enabled = false
	if t.conn != nil {
		(*t.conn).Close()
		t.conn = nil
	}
}

// start will establish an ssh server listening on the backchannel
func (t *attachServerSSH) run() error {
	var sConn *ssh.ServerConn
	var chans <-chan ssh.NewChannel
	var reqs <-chan *ssh.Request
	var err error

	// keep waiting for the connection to establish
	for t.enabled && sConn == nil {
		// wait for backchannel to establish
		conn, err := backchannel(context.Background())
		t.conn = &conn

		// create the SSH server
		sConn, chans, reqs, err = ssh.NewServerConn(*t.conn, t.config)
		if err != nil {
			detail := fmt.Sprintf("failed to establish ssh handshake: %s", err)
			log.Error(detail)
			continue
		}

		defer sConn.Close()
	}
	if err != nil {
		detail := fmt.Sprintf("abandoning attempt to start attach server: %s", err)
		log.Error(detail)
		return err
	}

	// Global requests
	go t.globalMux(reqs)

	log.Println("ready to service attach requests")
	// Service the incoming channels
	for attachchan := range chans {
		// The only channel type we'll support is "attach"
		if attachchan.ChannelType() != "attach" {
			detail := fmt.Sprintf("unknown channel type %s", attachchan.ChannelType())
			log.Error(detail)
			attachchan.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		// check we have a Session matching the requested ID
		bytes := attachchan.ExtraData()
		if bytes == nil {
			detail := "attach channel requires ID in ExtraData"
			log.Error(detail)
			attachchan.Reject(ssh.Prohibited, detail)
			continue
		}

		session, ok := config.Sessions[string(bytes)]
		if !ok || session.Cmd.Cmd == nil || session.Cmd.Cmd.ProcessState.Exited() {
			detail := fmt.Sprintf("specified ID for attach is unavailable: %s", string(bytes))
			log.Error(detail)
			attachchan.Reject(ssh.Prohibited, detail)
			continue
		}

		channel, requests, err := attachchan.Accept()
		if err != nil {
			detail := fmt.Sprintf("could not accept channel: %s", err)
			log.Println(detail)
			continue
		}

		// bind the channel to the Session
		if !session.Tty {
			// if it's not a TTY then bind the channel directly to the multiwriter that's already associated with the process
			dmwStdout, okA := session.Cmd.Cmd.Stdout.(dio.DynamicMultiWriter)
			dmwStderr, okB := session.Cmd.Cmd.Stdout.(dio.DynamicMultiWriter)
			if !okA || !okB {
				detail := fmt.Sprintf("target session IO cannot be duplicated to attach streams: %s", string(bytes))
				log.Error(detail)
				attachchan.Reject(ssh.ConnectionFailed, detail)
				continue
			}

			dmwStdout.Add(channel)
			dmwStderr.Add(channel.Stderr())

			go t.channelMux(requests, session.Cmd.Cmd.Process, nil)
			continue
		}

		// if it's a TTY bind the channel to the multiwriter that's on the far side of the PTY from the process
		// this is done so the logging is done with processed output
		ptysession, ok := ptys[session.ID]
		ptysession.writer.Add(channel)
		// PTY merges stdout & stderr so the two are the same

		go t.channelMux(requests, session.Cmd.Cmd.Process, ptysession.pty)
	}

	log.Println("incoming attach channel closed")

	return nil
}

func (t *attachServerSSH) globalMux(reqchan <-chan *ssh.Request) {
	for req := range reqchan {
		var pendingFn func()
		var payload []byte
		ok := true

		log.Printf("received global request type %v", req.Type)

		switch req.Type {
		case "container-ids":
			keys := make([]string, len(config.Sessions))
			i := 0
			for k := range config.Sessions {
				keys[i] = k
				i++
			}

			payload = []byte(ssh.Marshal(keys))
		default:
			ok = false
			payload = []byte("unknown global request type: " + req.Type)
		}

		log.Debugf("Returning payload: %s", string(payload))

		// make sure that errors get send back if we failed
		if req.WantReply {
			req.Reply(ok, payload)
		}

		// run any pending work now that a reply has been sent
		if pendingFn != nil {
			log.Debug("Invoking pending work")
			go pendingFn()
			pendingFn = nil
		}
	}
}

func (t *attachServerSSH) channelMux(in <-chan *ssh.Request, process *os.Process, pty *os.File) {
	var err error
	for req := range in {
		var pendingFn func()
		var payload []byte
		ok := true

		switch req.Type {
		case "window-change":
			msg := WindowChangeMsg{}
			if pty == nil {
				ok = false
				payload = []byte("illegal window-change request for non-tty")
			} else if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else if err := resizePty(pty.Fd(), &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			}
		case "signal":
			msg := signalMsg{}
			if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
				ok = false
				payload = []byte(err.Error())
			} else {
				log.Printf("Sending signal %s to container process, pid=%d\n", string(msg.Signal), process.Pid)
				err := signalProcess(process, ssh.Signal(msg.Signal))
				if err != nil {
					log.Printf("Failed to dispatch signal to process: %s\n", err)
				}
				payload = []byte{}
			}
		default:
			ok = false
			err = fmt.Errorf("ssh request type %s is not supported", req.Type)
			log.Println(err.Error())
		}

		// make sure that errors get send back if we failed
		if req.WantReply {
			req.Reply(ok, payload)
		}

		// run any pending work now that a reply has been sent
		if pendingFn != nil {
			log.Debug("Invoking pending work")
			go pendingFn()
			pendingFn = nil
		}
	}
	log.Println("incoming attach request channel closed")
}
