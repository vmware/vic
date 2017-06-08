// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"

	"github.com/vmware/vic/pkg/trace"
)

const (
	// time to wait for websocket writes
	waitTime = time.Minute * 3
)

var (
	// lock for single command execution
	cmdDone   = make(chan error, 1)
	logStream = NewLogStream()
	mu        sync.Mutex
)

// LogStream streams a command's execution over a websocket connection
type LogStream struct {
	cmd       *exec.Cmd
	websocket *websocket.Conn
}

// NewLogStream returns an empty LogSream struct (no websocket connection or exec command)
func NewLogStream() *LogStream {
	defer trace.End(trace.Begin(""))

	return &LogStream{}
}

func (ls *LogStream) setCmd(command []string) {
	defer trace.End(trace.Begin(""))

	ls.cmd = exec.Command(command[0], command[1:]...)
}

func (ls *LogStream) start() {
	defer trace.End(trace.Begin(""))

	// log reader and writer.
	logReader, logWriter, err := os.Pipe()
	if err != nil {
		log.Infoln("ERROR")
		log.Infoln(err)
	}
	defer logReader.Close()
	defer logWriter.Close()

	// attach cmd std out and std err to log Writer
	ls.cmd.Stderr = logWriter
	ls.cmd.Stdout = logWriter
	r, _ := regexp.Compile(`DOCKER_HOST=(\d{1,3}\.){3}(\d{1,3}):\d{4}`)

	go func() {
		s := bufio.NewScanner(logReader)
		for s.Scan() {
			ls.send(string(s.Bytes()))
			//if we get a docker endpoint the setup is complete and we should attach this vch to admiral
			match := r.Find(s.Bytes())
			stringMatch := string(match)
			if err == nil && strings.Contains(stringMatch, "=") {
				dockerIP := strings.Split(stringMatch, "=")[1]
				log.Infof("Docker endpoint is: %s\n", dockerIP)
				go setupDefaultAdmiral(string(dockerIP))
			}
		}
	}()

	go func() {
		cmdDone <- ls.cmd.Run()
	}()

	select {
	case <-time.After(waitTime):
		if err := ls.cmd.Process.Kill(); err != nil {
			ls.send(fmt.Sprintf("failed to kill create: %v ", err))
		}
		ls.send("Create exited after timeout")
	case err := <-cmdDone:
		if err != nil {
			ls.send(fmt.Sprintf("Create failed with error: %v\n", err))
		} else {
			ls.send("Execution complete.")
		}
	}
}

func (ls *LogStream) websocketServer(resp http.ResponseWriter, req *http.Request) {
	defer trace.End(trace.Begin(""))

	//turn http requests into websockets
	upgrader := websocket.Upgrader{}
	websocket, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Infoln("ERROR")
		log.Infoln(err)
		panic(err)
	}

	//set logstrem websocket for use by start() and send()
	ls.websocket = websocket
	defer ls.websocket.Close()

	//create the command
	ls.setCmd(engineInstaller.CreateCommand)
	ls.start()
}

func (ls *LogStream) send(msg string) {
	defer trace.End(trace.Begin(""))

	mu.Lock()
	defer mu.Unlock()

	ls.websocket.SetWriteDeadline(time.Now().Add(waitTime))
	if err := ls.websocket.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Infof("ERROR SENDING --------\n%s\n--------: %v\n", msg, err)
	} else {
		log.Infof("SENT: %s\n", msg)
	}
}
