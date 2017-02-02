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

package interaction

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kr/pty"

	"github.com/vmware/vic/lib/tether-ng/types"
	"github.com/vmware/vic/pkg/dio"
)

const (
	ioCopyBufferSize = 32 * 1024
)

type Interaction struct {
	uuid uuid.UUID
	ctx  context.Context

	config types.ExecutorConfig
}

func NewInteraction(ctx context.Context) *Interaction {
	return &Interaction{
		uuid: uuid.New(),
		ctx:  ctx,
	}
}

func (i *Interaction) Configure(ctx context.Context, config *types.ExecutorConfig) error {
	// create our own copy
	i.config = *config

	return nil
}

func (i *Interaction) Start(ctx context.Context) error    { return nil }
func (i *Interaction) Stop(ctx context.Context) error     { return nil }
func (i *Interaction) UUID(ctx context.Context) uuid.UUID { return i.uuid }

func (i *Interaction) Release(ctx context.Context, out chan<- chan struct{}) {
	fmt.Printf("Releasing\n")
	if out != nil {
		// make a new response chan
		release := make(chan struct{})

		fmt.Printf("Really releasing %#v\n", out)

		fmt.Printf("Sleeping some before unblocking\n")

		time.Sleep(3 * time.Second)

		// send the response
		out <- release
	}
	fmt.Printf("Done Releasing\n")
}

func (i *Interaction) PseudoTerminal(ctx context.Context, in <-chan *types.Session) <-chan struct{} {
	var err error

	var wg sync.WaitGroup
	c := make(chan struct{})

	// get the session
	session := <-in

	// add SSH channels
	session.Reader = dio.MultiReader(os.Stdin)
	session.Outwriter = dio.MultiWriter(os.Stdout)

	session.Pty, err = pty.Start(&session.Cmd)
	if err != nil {
		close(c)
		return c
	}

	wg.Add(1)
	go func() {
		_, gerr := io.CopyBuffer(session.Pty, session.Reader, make([]byte, ioCopyBufferSize))
		fmt.Printf("stdin returned: %s\n", gerr)
	}()

	go func() {
		_, gerr := io.CopyBuffer(session.Outwriter, session.Pty, make([]byte, ioCopyBufferSize))
		fmt.Printf("stdout returned: %s\n", gerr)

		wg.Done()
	}()

	// wait all and close the channel
	go func() {
		wg.Wait()
		close(c)
	}()

	return c
}

func (i *Interaction) NonInteract(ctx context.Context, in <-chan *types.Session) <-chan struct{} {
	var wg sync.WaitGroup
	c := make(chan struct{})

	// get the session
	session := <-in

	// add SSH channels
	session.Reader = dio.MultiReader(os.Stdin)
	session.Outwriter = dio.MultiWriter(os.Stdout)
	session.Errwriter = dio.MultiWriter(os.Stderr)

	// add SSH channels
	session.Outwriter.Add(os.Stdout)
	session.Errwriter.Add(os.Stderr)

	// get pipes
	stdin, _ := session.Cmd.StdinPipe()
	stdout, _ := session.Cmd.StdoutPipe()
	stderr, _ := session.Cmd.StderrPipe()

	wg.Add(2)
	go func() {
		if session.OpenStdin {
			_, gerr := io.CopyBuffer(stdin, session.Reader, make([]byte, ioCopyBufferSize))
			fmt.Printf("stdin returned: %s\n", gerr)
		}
	}()

	go func() {
		_, gerr := io.CopyBuffer(session.Outwriter, stdout, make([]byte, ioCopyBufferSize))
		fmt.Printf("stdout returned: %s\n", gerr)
		wg.Done()
	}()

	go func() {
		_, gerr := io.CopyBuffer(session.Errwriter, stderr, make([]byte, ioCopyBufferSize))
		fmt.Printf("stderr returned: %s\n", gerr)
		wg.Done()
	}()

	// wait all and close the channel
	go func() {
		wg.Wait()
		close(c)
	}()

	return c
}

func (i *Interaction) Close(ctx context.Context, in <-chan *types.Session) <-chan struct{} {
	var wg sync.WaitGroup
	c := make(chan struct{})

	// get the session
	session := <-in

	fmt.Printf("Closing for %s\n", session.ID)
	if session.Pty != nil {
		fmt.Printf("PTY close\n")
		// https://github.com/golang/go/issues/7970
		session.Pty.Close()

		close(c)
		return c
	}

	wg.Add(3)
	go func() {
		if session.OpenStdin {
			session.Reader.Close()
		}
		wg.Done()
	}()

	go func() {
		session.Outwriter.Close()
		wg.Done()
	}()

	go func() {
		session.Errwriter.Close()
		wg.Done()
	}()

	// wait all and close the channel
	go func() {
		wg.Wait()

		fmt.Printf("NONPTY close\n")

		close(c)
	}()

	return c
}
