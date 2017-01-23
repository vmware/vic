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

package scp

import (
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type Request struct {
	ch        ssh.Channel
	pendingFn func()
}

// From remote to local.  Remote is server (source)
// $ { scp -vvvv luser@hasselhoff:/etc/fstab /tmp/f 2>&1; } | grep -e scp -e Sink
// Executing: program /usr/bin/ssh host hasselhoff, user luser, command scp -v -f /etc/fstab
// debug1: Sending command: scp -v -f /etc/fstab
// Sink: C0644 968 fstab

// To remote from local.  Remote is server (but this time is dest)
// $ { scp -vvvv /tmp/f luser@hasselhoff:/tmp/f 2>&1; } | grep -e scp -e Sink
// Executing: program /usr/bin/ssh host hasselhoff, user luser, command scp -v -t /tmp/f
// debug1: Sending command: scp -v -t /tmp/f
// Sink: C0664 968 f

// For copying from a host.  This acts as the server (source).
func (scp *Request) Source(path string) (ok bool, payload []byte) {

	op, err := OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return false, []byte(err.Error())
	}
	f := func() {
		defer op.Close()
		defer scp.ch.Close()

		n, err := op.Read(scp.ch)
		if err != nil {
			log.Printf("Source: error reading %s", err)
			return
		}

		log.Printf("Source: copied %s (%d bytes) to client", op.File.Name(), n)
	}

	scp.pendingFn = f
	return true, nil
}

func (scp *Request) Destination(path string) (ok bool, payload []byte) {
	f := func() {

		op, err := Write(scp.ch, path)
		if err != nil && err != io.EOF {
			log.Printf("Source: error writing %s", err)
			return
		}

		if err = op.Close(); err != nil {
			log.Printf("error closing %s", err)
			return
		}

		scp.ch.Close()
		log.Printf("Destination: copied %s (%d bytes) from client", op.File.Name(), op.Size)
	}

	scp.pendingFn = f
	return true, nil
}

func (scp *Request) SetChannel(channel *ssh.Channel) {
	scp.ch = *channel
}

func (scp *Request) GetPendingWork() func() {
	return scp.pendingFn
}

func (scp *Request) ClearPendingWork() {
}
