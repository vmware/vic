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
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestScpFromHost(t *testing.T) {
	scpTest(t, SOURCE)
}

func TestScpToHost(t *testing.T) {
	scpTest(t, DEST)
}

func scpTest(t *testing.T, mode Mode) {
	port := 0
	req := &ScpRequest{}
	serverWg, err := StartSSHExecServer(&port, req)
	if err != nil {
		t.Error(err)
		return
	}

	sourceFile := "scp_test.go"
	tmpFile, err := ioutil.TempFile("", sourceFile)
	if err != nil {
		t.Fatal(err)
	}
	destFile := tmpFile.Name()

	var sourceFileMd5 []byte
	var md5wg sync.WaitGroup
	md5wg.Add(1)
	go func() {
		defer md5wg.Done()
		sourceFileMd5, err = ComputeMd5(sourceFile)
		if err != nil {
			t.Error(err)
		}
	}()

	if err = scp(port, mode, sourceFile, destFile); err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(destFile)

	destFileMd5, err := ComputeMd5(destFile)
	if err != nil {
		t.Error(err)
		return
	}

	serverWg.Wait()
	md5wg.Wait()
	if bytes.Compare(sourceFileMd5, destFileMd5) != 0 {
		t.Fatalf("files don't match %s %s", hex.EncodeToString(sourceFileMd5), hex.EncodeToString(destFileMd5))
		return
	}

}

func scp(port int, serverMode Mode, sourceFile, destFile string) error {
	log.Printf("Copying %s to %s in %s mode", sourceFile, destFile, string(serverMode))

	// Create client config
	config := &ssh.ClientConfig{
		User: TestUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(TestPassword),
		},
	}
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", "localhost:"+strconv.Itoa(port), config)
	if err != nil {
		return fmt.Errorf("unable to connect: %s", err)
	}
	defer conn.Close()

	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("unable to create session: %s", err)
	}

	var op *Operation
	// copy from the server
	if serverMode == SOURCE {

		sshPipe, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("error openning pipe %s", err)
		}

		ok, err := session.SendRequest(string(serverMode), true, []byte(sourceFile))
		if !ok {
			return fmt.Errorf("not ok")
		}
		if err != nil {
			return fmt.Errorf("error sending request %s", err)
		}

		// write the result locally using the scp protocol operation
		op, err = Write(ioutil.NopCloser(sshPipe), destFile)
		if err != nil {
			return fmt.Errorf("error writing %s", err)
		}

	} else {

		sshPipe, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("error openning pipe %s", err)
		}

		// copy to the server
		ok, err := session.SendRequest(string(serverMode), true, []byte(destFile))
		if !ok {
			return fmt.Errorf("not ok")
		}
		if err != nil {
			return fmt.Errorf("error sending request %s", err)
		}

		// open the file for reading locally
		op, err = OpenFile(sourceFile, os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		// read the file into the pipe
		n, err := op.Read(sshPipe)
		if err != nil {
			return fmt.Errorf("error reading %s", err)
		}

		if err = session.Close(); err != nil {
			return fmt.Errorf("error closing sesson %s", err)
		}

		log.Printf("scp copied %d bytes to server", n)
	}

	if err = op.Close(); err != nil {
		return fmt.Errorf("error closing %s", err)
	}

	if err = conn.Wait(); err != nil && err != io.EOF {
		return fmt.Errorf("error waiting for conn close %s", err)
	}

	return nil
}

func ComputeMd5(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}
