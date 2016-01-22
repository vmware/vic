package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"enatai-gerrit.eng.vmware.com/bonneville-container/tether"
	"golang.org/x/crypto/ssh"
)

func (ch *GlobalHandler) StartConnectionManager(conn *ssh.ServerConn) {
}

func (ch *GlobalHandler) ContainerId() string {
	return ch.id
}

func (ch *GlobalHandler) StaticIPAddress(cidr, gateway string) error {
	var ip net.IP
	var err error

	// first return is ip address
	if ip, _, err = net.ParseCIDR(cidr); err != nil {
		return err
	}

	log.Printf("Set IP address to %s", ip.String())

	return nil
}

func (ch *GlobalHandler) DynamicIPAddress() (string, error) {
	return "", nil
}

func (c *GlobalHandler) MountLabel(label, target string) error {
	detail := fmt.Sprintf("Unable to mount %s: ", label, "mount not implemented")
	log.Print(detail)
	return errors.New(detail)
}

func (c *GlobalHandler) Sync() {

}

func (c *GlobalHandler) NewSessionContext() tether.SessionContext {
	return &SessionHandler{
		GlobalHandler: c,
		env:           make(map[string]string),
	}
}
