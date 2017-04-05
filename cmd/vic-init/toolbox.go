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

// +build !windows,!darwin

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/diag"
	"github.com/vmware/vic/pkg/vsphere/toolbox"

	log "github.com/Sirupsen/logrus"
)

// startCommand is the switch for the synthetic commands that are permitted within the appliance.
// This is not intended to allow arbitrary commands to be executed.
// returns:
//  pid: we return -1 as this is a synthetic command
//  error
func startCommand(r *toolbox.VixMsgStartProgramRequest) (int, error) {
	defer trace.End(trace.Begin(r.ProgramPath))

	switch r.ProgramPath {
	case "enable-ssh":
		return -1, enableSSH(r.Arguments)
	case "passwd":
		return -1, passwd(r.Arguments)
	case "test-vc-api":
		return diag.CheckAPIAvailability(r.Arguments), nil
	default:
		return -1, fmt.Errorf("unknown command %q", r.ProgramPath)
	}
}

// enableShell changes the root shell from /bin/false to /bin/bash
func enableShell() error {
	defer trace.End(trace.Begin(""))

	// #nosec
	chsh := exec.Command("/bin/chsh", "-s", "/bin/bash", "root")
	err := chsh.Start()
	if err != nil {
		err := fmt.Errorf("Failed to launch chsh: %s", err)
		log.Error(err)
		return err
	}

	// ignore the error - it's likely raced with child reaper, we just want to make sure
	// that it's exited by the time we pass this point
	chsh.Wait()

	// confirm the change
	file, err := os.Open("/etc/passwd")
	if err != nil {
		err := fmt.Errorf("Failed to open file to confirm change: %s", err)
		log.Error(err)
		return err
	}

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		err := fmt.Errorf("Failed to read line from file to confirm change: %s", err)
		log.Error(err)
		return err
	}

	// assert that first line is root
	if !strings.HasPrefix(line, "root") {
		err := fmt.Errorf("Expected line to start with root: %s", line)
		log.Error(err)
		return err
	}

	// assert that first line is root
	if !strings.HasSuffix(line, "/bin/bash\n") {
		err := fmt.Errorf("Expected line to end with /bin/bash: %s", line)
		log.Error(err)
		return err
	}

	return nil
}

// passwd sets the password for the root user to that provided as an argument
func passwd(pass string) error {
	defer trace.End(trace.Begin(""))

	err := enableShell()
	if err != nil {
		err := fmt.Errorf("Failed to enable shell: %s", err)
		log.Error(err)
		// continue anyway as people may be able to get something useful
	}

	// #nosec
	setPasswd := exec.Command("/sbin/chpasswd")
	stdin, err := setPasswd.StdinPipe()
	if err != nil {
		err := fmt.Errorf("Failed to create stdin pipe for chpasswd: %s", err)
		log.Error(err)
		return err
	}

	err = setPasswd.Start()
	if err != nil {
		err := fmt.Errorf("Failed to launch chpasswd: %s", err)
		log.Error(err)
		return err
	}

	_, err = stdin.Write([]byte("root:" + pass))

	// so that we're actively waiting when the process exits, or we'll race (and lose) to child reaper
	go func() {
		setPasswd.Wait()
	}()

	err = stdin.Close()
	if err != nil {
		err := fmt.Errorf("Failed to close input to chpasswd: %s", err)
		log.Error(err)

		// fire and forget as we're already on error path
		setPasswd.Process.Kill()

		return err
	}

	return nil
}

// enableSSH receives a key as an argument
func enableSSH(key string) error {
	defer trace.End(trace.Begin(""))

	err := enableShell()
	if err != nil {
		err := fmt.Errorf("Failed to enable shell: %s", err)
		log.Error(err)
		// continue anyway as people may be able to get something useful
	}

	// basic sanity check for args - we don't bother validating it's a key
	if len(key) != 0 {
		err = os.MkdirAll("/root/.ssh", 0700)
		if err != nil {
			err := fmt.Errorf("unable to create path for keys: %s", err)
			log.Error(err)
			return err
		}

		err = ioutil.WriteFile("/root/.ssh/authorized_keys", []byte(key), 0600)
		if err != nil {
			err := fmt.Errorf("unable to create authorized_keys: %s", err)
			log.Error(err)
			return err
		}
	}

	return startSSH()
}

// startSSH launches the sshd server
func startSSH() error {
	// #nosec
	c := exec.Command("/usr/bin/systemctl", "start", "sshd")

	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = &b

	if err := c.Start(); err != nil {
		return err
	}

	go func() {
		// because init is explicitly reaping child processes we cannot use simple
		// exec commands to gather status
		_ = c.Wait()
		log.Info("Attempted to start ssh service:\n %s", b)
	}()

	return nil
}
