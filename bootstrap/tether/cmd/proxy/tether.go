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

// +build linux

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/vmware/vic/bootstrap/tether"
	"github.com/vmware/vic/bootstrap/tether/handlers"
	"github.com/vmware/vic/bootstrap/tether/serial"
	"github.com/vmware/vic/bootstrap/tether/utils"

	"golang.org/x/crypto/ssh"
)

var (
	tetherKey string
	private   ssh.Signer
)

const sttyArg0 = "/bin/stty"
const magicPrompt = "Exiting DOS " // the prompt is magic because when we see the prompt, we disconnect from the container - it also serves as a useful message

var (
	sttyArgvEOff []string = []string{"stty", "raw", "-echo"}
	sttyArgvEOn  []string = []string{"stty", "cooked", "echo"}
)

func privateKey(file string) ssh.Signer {
	privateBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("failed to load private key: %v", tetherKey)
	}

	priv, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", tetherKey)
	}
	return priv
}

// Hopefully this gets rolled into go.sys at some point
func Mkdev(majorNumber int64, minorNumber int64) int {
	return int((majorNumber << 8) | (minorNumber & 0xff) | ((minorNumber & 0xfff00) << 12))
}

func init() {
	flag.StringVar(&tetherKey, "key", "/opt/proxy_tether_init_key", "tetherd control channel private key")
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		// mustn't listen on the container portgroup
		if iface.Name == "docker0" {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func handleConnectionFromGuest(guestConn *net.TCPConn) {
	defer guestConn.Close()

	// create our per connection state
	handler := handlers.GlobalProxyHandler{}
	handler.SetConn(guestConn)
	handler.SetAllowExec(true)

	// Establish connection to daemon
	log.Println("Establishing the connection to docker daemon")
	dconn, err := net.Dial("tcp", "localhost:2377")
	if err != nil {
		log.Printf("failed to create loopback connection to daemon: %s", err)
		return
	}
	defer dconn.Close()

	log.Println("Handshaking with daemon")
	serial.HandshakeServer(dconn, time.Duration(10*time.Second))

	handler.DiscardUntilPrompt()

	output, _ := handler.CmdCombinedOutput("type C:\\TETHER\\DOCKERID")
	log.Printf("Docker ID passed back from container = \"%s\"", output)

	handler.SetContainerId(utils.StripCommandOutput(output))

	err = handler.CmdStart("C:\\AUTOEXEC.BAT")
	if err != nil {
		log.Printf("error running C:\\AUTOEXEC.BAT: %s", err)
		return
	}
	// Set the magic prompt which is the signal to disconnect from the container
	//   expectation is that command.com is run from the magic prompt and resets to a standard prompt
	//   this is done as the entry point of the base Dockerfile
	// Note that magic prompt must be run after autoexec as that could reset the prompt
	err = handler.CmdStart(fmt.Sprintf("SET PROMPT=%s$g", magicPrompt))
	if err != nil {
		log.Printf("error setting initial prompt: %s", err)
		return
	}

	log.Printf("Initiating ssh command logic with handler.id = %s", handler.ContainerId())
	tether.StartTether(dconn, private, &handler)
}

func main() {
	// seems necessary given rand.Reader access
	var err error

	flag.Parse()

	// Parse ssh private key
	private = privateKey(tetherKey)

	// start our main connection loop
	ext_ip, err := externalIP()

	if err != nil {
		log.Printf("unable to get external IP address: %s", err)
	}

	addr, err := net.ResolveTCPAddr("tcp", ext_ip+":2377")
	if err != nil {
		log.Printf("unable to construct TCP address to listen: %s", err)
	}

	log.Println("Listening for guest connections")
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Printf("unable listen %s: %s", addr, err)
	}
	defer ln.Close()

	// main listen loop
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Println("accept error")
		} else {
			log.Println("handing guest connection off to handler thread")
			go handleConnectionFromGuest(conn)
		}
	}
}
