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
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"

	"github.com/vmware/vic/pkg/trace"
)

const (
	// time to wait for ws writes
	waitTime = time.Second * 5
)

var (
	// lock for single command execution
	cmdDone   = make(chan error, 1)
	logStream = NewLogStream()
)

// LogStream streams a command's execution over a websocket connection
type LogStream struct {
	cmd *exec.Cmd
	ws  *websocket.Conn
}

// NewLogStream returns an empty LogSream struct (no ws connection or exec command)
func NewLogStream() *LogStream {
	return &LogStream{}
}

func (ls *LogStream) setCmd(command string) {
	defer trace.End(trace.Begin(""))

	args := strings.Split(command, " ")
	ls.cmd = exec.Command(args[0], args[1:]...)
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

	go func() {
		s := bufio.NewScanner(logReader)
		for s.Scan() {
			log.Infoln("scanning...")
			ls.send(string(s.Bytes()))
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

func (ls *LogStream) wsServer(resp http.ResponseWriter, req *http.Request) {
	defer trace.End(trace.Begin(""))
	//defer ls.ws.Close()

	//turn http requests into websockets
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Infoln("ERROR")
		log.Infoln(err)
		panic(err)
	}

	//set logstrem websocket for use by start() and send()
	ls.ws = ws

	//create the command
	ls.setCmd(engineInstaller.CreateCommand)
	ls.start()
}

func (ls *LogStream) send(msg string) {
	ls.ws.SetWriteDeadline(time.Now().Add(waitTime))
	if err := ls.ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Infof("ERROR: %v\n", err)
	} else {
		log.Infof("SENT: %s\n", msg)
	}
}
