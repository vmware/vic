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
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/vmware/vic/tether"
	"github.com/vmware/vic/tether/utils"

	"golang.org/x/crypto/ssh"
)

const magicPrompt = "Exiting DOS " // the prompt is magic because when we see the prompt, we disconnect from the container - it also serves as a useful message

type SessionProxyHandler struct {
	*GlobalProxyHandler

	// TODO: add some locking in here if non exec requests can touch exec or shell related items
	channel   *ssh.Channel
	env       map[string]string
	assignTty bool
	waitGroup sync.WaitGroup
	pendingFn func()
}

func (ch *SessionProxyHandler) SetChannel(channel *ssh.Channel) {
	log.Println("Called SetChannel")
	ch.channel = channel
	log.Printf("Set channel to %v\n", ch.channel)
}

func (ch *SessionProxyHandler) Setenv(name, value string) (ok bool, payload []byte) {
	log.Println("Called Setenv")

	cmd := fmt.Sprintf("set %s=\"%s\"", name, value)

	ch.CmdCombinedOutput(cmd)

	ch.env[name] = value
	fmt.Printf("Set environment variable: %s=%s\n", name, value)

	return true, nil
}

func (ch *SessionProxyHandler) AssignPty() {
	log.Println("Called AssignPty")
	ch.assignTty = true
}

func (ch *SessionProxyHandler) ResizePty(winSize *tether.WindowChangeMsg) error {
	// returning nil so we fail soft
	log.Println("Called Resizepty")
	return nil
}

func (ch *SessionProxyHandler) Shell() (ok bool, payload []byte) {
	log.Println("Called Shell")
	//TODO: implement
	return false, []byte("shell request is not implemented")
}

func (ch *SessionProxyHandler) Signal(sig ssh.Signal) error {
	log.Println("Called Signal")
	detail := "Unable to signal process: signal not supported"
	log.Print(detail)
	return errors.New(detail)
}

func (ch *SessionProxyHandler) Kill() error {
	log.Println("Called Kill")
	detail := "Unable to kill process: kill not supported"
	log.Print(detail)
	return errors.New(detail)
}

func (ch *SessionProxyHandler) Exec(command string, args []string, config map[string]string) (ok bool, payload []byte) {
	if !ch.allowExec {
		detail := "Multiple execs not supported"
		log.Println(detail)
		return false, []byte(detail)
	}

	log.Println("Called Exec")
	// strip quotes from the args if they are first AND last positions in an arg element
	cmd_str := command
	for k, v := range args {
		if v[0] == '"' && v[len(v)-1] == '"' {
			args[k] = v[1 : len(v)-1]
		}
		cmd_str = cmd_str + " " + v
	}

	// print a welcome message on Exec along with info on how to connect to the graphics
	sshChannel := *ch.channel
	sshChannel.Write([]byte("\r\nWelcome to MS-DOS on Bonneville! To connect graphics to this VM, type: "))
	output, _ := ch.CmdCombinedOutput("type C:\\TETHER\\GRAPHCOM")
	sshChannel.Write([]byte(utils.StripCommandOutput(output) + "\r\n\r\n"))

	log.Printf("Sending command %+q", cmd_str)
	ch.CmdStart(cmd_str)
	log.Printf("Command sent: %+q", cmd_str)

	log.Printf("Copying to channel to %v\n", ch.channel)

	ch.pendingFn = func() {
		if err := ch.copyUntilPrompt(); err != nil {
			fmt.Println("Unexpected read/write error in copyUntilPrompt: ", err)
		}
		// ensure that changes are flushed to disk before we report exit
		ch.Sync()

		var exitStatus uint32 = 0

		if err := (*ch.channel).CloseWrite(); err != nil {
			fmt.Println("Error sending channel EOF: ", err)
		}

		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, exitStatus)
		if _, err := (*ch.channel).SendRequest("exit-status", false, bytes); err != nil {
			fmt.Println("Error sending exit status: ", err)
		}

		if err := (*ch.channel).Close(); err != nil {
			fmt.Println("Error sending channel close: ", err)
		}

		fmt.Println("Returned exit status and closed channel")
	}

	// send the immediate reply to the exec request
	fmt.Println("Started process successfully")
	return true, nil
}

func (ch *SessionProxyHandler) GetPendingWork() func() {
	log.Println("Called GetPendingWork")
	return ch.pendingFn
}

func (ch *SessionProxyHandler) ClearPendingWork() {
	log.Println("Called ClearPendingWork")
	ch.pendingFn = nil
}

// Reads from the dos connection until it sees a command prompt
// We're presuming that means end of output
// Copies the data to the handler.chennel
func (ch *SessionProxyHandler) copyUntilPrompt() error {
	// wait for the exec reply to be sent
	ch.waitGroup.Add(1)
	var fatalerr error = nil
	log.Println("copying the data from ssh to dosconn and vice versa")
	go func() {
		buf := make([]byte, 16)
		sshChannel := *ch.channel
		exitTrigger := magicPrompt + ">" // We're monitoring the output looking for this string
		triggerPos := 0
	ReadLoop:
		for {
			nr, er := ch.dosconn.Read(buf)
			// log.Printf(">\"%s\"", string(buf[:]))
			if nr > 0 {
				for i := 0; i < nr; i++ {
					// log.Printf("comparing %c with %c: %d, %d", exitTrigger[triggerPos], buf[i], triggerPos, i)
					if exitTrigger[triggerPos] == buf[i] {
						triggerPos++
						if triggerPos == len(exitTrigger) {
							sshChannel.Write([]byte{13, 10}) // crlf on the end
							break ReadLoop
						}
					} else {
						triggerPos = 0
					}
				}
				nw, ew := sshChannel.Write(buf[0:nr])
				if ew != nil {
					fatalerr = fmt.Errorf("Error writing to ssh channel: %s", ew)
					break
				}
				if nw != nr {
					fatalerr = fmt.Errorf("Short write! %d != %d", nr, nw)
					break
				}
			}
			if er == io.EOF {
				break
			} else if er != nil {
				if e, ok := er.(*net.OpError); ok { // We can't distinguish between different types of OpError, but most are fatal
					log.Printf("Lost read connection: %s. Exiting", e)
					fatalerr = er
					break
				}
			}
		}
		ch.waitGroup.Done()
	}()

	// we shouldn't wait on stdin to close or we'll be here forever
	go func() { log.Println("Copying from daemon to dosconn"); io.Copy(ch.dosconn, *ch.channel) }()

	log.Println("Waiting for I/O streams to close")
	ch.waitGroup.Wait()
	if fatalerr != nil {
		log.Println("Returning read/write error: ", fatalerr)
	}
	return fatalerr
}

// Run the command and copy the output until it's done
func (ch *SessionProxyHandler) cmdRun(cmd string) error {
	_, _ = ch.dosconn.Write([]byte(cmd + "\r\n"))

	return ch.copyUntilPrompt()
}
