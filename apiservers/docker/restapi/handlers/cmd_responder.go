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
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/go-swagger/go-swagger/httpkit"
)

type CmdResponder struct {
	writer  io.Writer
	flusher http.Flusher
	cmdPath string
	cmdArgs []string
}

func NewCmdResponder(path string, args []string) *CmdResponder {
	responder := &CmdResponder{cmdPath: path, cmdArgs: args}

	return responder
}

func (cr *CmdResponder) Write(data []byte) (int, error) {
	n, err := cr.writer.Write(data)

	cr.flusher.Flush()

	return n, err
}

// WriteResponse to the client
func (cr *CmdResponder) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {
	var exist bool

	rw.Header().Set("Content-Type", "application/json")

	cr.flusher, exist = rw.(http.Flusher)

	if exist {
		cr.writer = rw

		cmd := exec.Command(cr.cmdPath, cr.cmdArgs...)
		cmd.Stdout = cr
		cmd.Stderr = cr

		// Execute
		err := cmd.Start()

		if err != nil {
			log.Printf("Error starting %s - %s\n", cr.cmdPath, err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Wait for the fetcher to finish.  Should effectively close the stdio pipe above.
		err = cmd.Wait()

		if err != nil {
			log.Println("imagec exit code:", err)
		}

		rw.WriteHeader(http.StatusOK)
		return
	}

	log.Println("CmdResponder failed to get the HTTP flusher")

	rw.WriteHeader(http.StatusInternalServerError)
}
