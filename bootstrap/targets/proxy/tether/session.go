package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"

	"enatai-gerrit.eng.vmware.com/bonneville-container/tether"
)

func (ch *SessionHandler) SetChannel(channel *ssh.Channel) {
	log.Println("Called SetChannel")
	ch.channel = channel
	log.Printf("Set channel to %v\n", ch.channel)
}

func (ch *SessionHandler) Setenv(name, value string) (ok bool, payload []byte) {
	log.Println("Called Setenv")

	cmd := fmt.Sprintf("set %s=\"%s\"", name, value)

	ch.cmdCombinedOutput(cmd)

	ch.env[name] = value
	fmt.Printf("Set environment variable: %s=%s\n", name, value)

	return true, nil
}

func (ch *SessionHandler) AssignPty() {
	log.Println("Called AssignPty")
	ch.assignTty = true
}

func (ch *SessionHandler) ResizePty(winSize *tether.WindowChangeMsg) error {
	// returning nil so we fail soft
	log.Println("Called Resizepty")
	return nil
}

func (ch *SessionHandler) Shell() (ok bool, payload []byte) {
	log.Println("Called Shell")
	//TODO: implement
	return false, []byte("shell request is not implemented")
}

func (ch *SessionHandler) Signal(sig ssh.Signal) error {
	log.Println("Called Signal")
	detail := fmt.Sprintf("Unable to signal process: ", "signal not supported")
	log.Print(detail)
	return errors.New(detail)
}

func (ch *SessionHandler) Kill() error {
	log.Println("Called Kill")
	detail := fmt.Sprintf("Unable to kill process: ", "kill not supported")
	log.Print(detail)
	return errors.New(detail)
}

func (ch *SessionHandler) Exec(command string, args []string, config map[string]string) (ok bool, payload []byte) {
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
	output, _ := ch.cmdCombinedOutput("type C:\\TETHER\\GRAPHCOM")
	sshChannel.Write([]byte(StripCommandOutput(output) + "\r\n\r\n"))

	log.Println("Sending command %+q", cmd_str)
	ch.cmdStart(cmd_str)
	log.Println("Ccommand sent: %+q", cmd_str)

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

func (ch *SessionHandler) GetPendingWork() func() {
	log.Println("Called GetPendingWork")
	return ch.pendingFn
}

func (ch *SessionHandler) ClearPendingWork() {
	log.Println("Called ClearPendingWork")
	ch.pendingFn = nil
}
