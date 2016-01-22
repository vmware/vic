package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"enatai-gerrit.eng.vmware.com/bonneville-container/tether"
	"golang.org/x/crypto/ssh"
)

type GlobalProxyHandler struct {
	GlobalHandler
	dosconn *net.TCPConn

	allowExec bool
}

func (ch *GlobalProxyHandler) SetContainerId(id string) {
	ch.id = id
}

func (ch *GlobalProxyHandler) StartConnectionManager(conn *ssh.ServerConn) {
}

func (ch *GlobalProxyHandler) StaticIPAddress(cidr, gateway string) error {
	var ip net.IP
	var err error

	// first return is ip address
	if ip, _, err = net.ParseCIDR(cidr); err != nil {
		return err
	}

	log.Printf("Set IP address to %s", ip.String())

	return nil
}

func (ch *GlobalProxyHandler) DynamicIPAddress() (string, error) {
	return "", nil
}

func (ch *GlobalProxyHandler) MountLabel(label, target string) error {
	detail := fmt.Sprintf("Unable to mount %s: ", label, "mount not implemented")
	log.Print(detail)
	return errors.New(detail)
}

func (ch *GlobalProxyHandler) Sync() {

}

func (ch *GlobalProxyHandler) NewSessionContext() tether.SessionContext {
	return &SessionProxyHandler{
		GlobalProxyHandler: ch,
		env:                make(map[string]string),
	}
}

func (ch *GlobalProxyHandler) SetConn(conn *net.TCPConn) {
	ch.dosconn = conn
}

func (ch *GlobalProxyHandler) SetAllowExec(s bool) {
	ch.allowExec = s
}

// Reads from the dos connection until it sees a command prompt
// We're presuming that means end of output
func (ch *GlobalProxyHandler) DiscardUntilPrompt() error {
	tmp := make([]byte, 16)
	for {
		n, err := ch.dosconn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
				return err
			}
			break
		}
		/* consume prompt text */
		if tmp[n-1] == '>' {
			break
		}

	}

	return nil
}

// Reads from the dos connection until it sees a command prompt
// We're presuming that means end of output
func (ch *GlobalProxyHandler) ReadUntilPrompt() (string, error) {
	buf := make([]byte, 0, 512)
	tmp := make([]byte, 16)

	for {
		n, err := ch.dosconn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
				return string(buf), err
			}
			break
		}

		buf = append(buf, tmp[:n]...)

		if tmp[n-1] == '>' {
			break
		}

	}

	// TODO: back track to last \n and strip entire prompt
	return string(buf), nil
}

// Run the command and return the output as a string
func (ch *GlobalProxyHandler) CmdCombinedOutput(cmd string) (string, error) {
	_, _ = ch.dosconn.Write([]byte(cmd + "\r\n"))

	return ch.ReadUntilPrompt()
}

// Run the command and return
func (ch *GlobalProxyHandler) CmdStart(cmd string) error {
	_, _ = ch.dosconn.Write([]byte(cmd + "\r\n"))

	// check for "command not found" without taking
	// it off the stream

	return nil
}
