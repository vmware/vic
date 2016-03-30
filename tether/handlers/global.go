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
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/vmware/vic/tether"

	"golang.org/x/crypto/ssh"
)

type GlobalHandler struct {
	id string
}

func NewGlobalHandler(id string) *GlobalHandler {
	return &GlobalHandler{
		id: id,
	}
}

func (ch *GlobalHandler) StartConnectionManager(conn *ssh.ServerConn) {
	log.Println("Registering fork trigger signal handler")
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in StartConnectionManager", r)
		}
	}()

	var incoming = make(chan os.Signal, 1)
	signal.Notify(incoming, syscall.SIGABRT)

	log.Println("SIGABRT handling initialized for fork support")
	for _ = range incoming {
		// validate that this is a fork trigger and not just a random signal from
		// container processes
		log.Println("Received SIGABRT - preparing to transition to fork parent")
		break
	}

	// tell client that we're disconnecting
	if ok, _, err := conn.SendRequest("fork", true, nil); !ok || err != nil {
		if err != nil {
			log.Printf("Unable to inform remote about fork (channel error): %s\n", err)
		} else {
			log.Println("Unable to register fork with remote - remote error")
		}
	} else {

		log.Println("Closing control connections")

		// regardless of errors we have to continue if externally driven
		conn.Close()

		// TODO: do we need to rebind session executions stdio to /dev/null or to files?
		// run the /.tether/vmfork.sh script
		log.Println("Running vmfork.sh")
		cmd := exec.Command("/.tether/vmfork.sh")
		// FORK HAPPENS DURING CALL, BEFORE RETURN FROM COMBINEDOUTPUT
		out, err := cmd.CombinedOutput()
		log.Printf("vmfork:%s\n%s\n", err, string(out))

		return
	}

	log.Println("Closing control connections")

	// regardless of errors we have to continue if externally driven
	conn.Close()

	// the StartTether loop will now exit and we'll fall back into waiting for SIGHUP in main
}

func (ch *GlobalHandler) ContainerId() string {
	return ch.id
}

func (c *GlobalHandler) NewSessionContext() tether.SessionContext {
	return &SessionHandler{
		GlobalContext: c,
		env:           make(map[string]string),
	}
}
