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

// Interaction implements the interaction plugin
type Interaction struct {
	uuid uuid.UUID
	ctx  context.Context

	err chan error

	config types.ExecutorConfig
}

// NewInteraction returns a new Interaction instance
func NewInteraction(ctx context.Context) *Interaction {
	return &Interaction{
		uuid: uuid.New(),
		ctx:  ctx,
		err:  make(chan error),
	}
}

// Configure sets the config
func (i *Interaction) Configure(ctx context.Context, config *types.ExecutorConfig) error {
	// create our own copy
	i.config = *config

	return nil
}

// Start starts the plugin
func (i *Interaction) Start(ctx context.Context) error { return nil }

// Stop stops the plugin
func (i *Interaction) Stop(ctx context.Context) error {
	// close the err chan as we are reporter
	close(i.err)

	return nil
}

func (i *Interaction) UUID(ctx context.Context) uuid.UUID { return i.uuid }

// Release releases the caller
func (i *Interaction) Release(ctx context.Context, out chan<- chan struct{}) {
	if out != nil {
		release := make(chan struct{})

		// simulating some work
		i.err <- fmt.Errorf("Something happened (not really)\n")
		time.Sleep(3 * time.Second)
		i.err <- fmt.Errorf("Something else happened (not really)\n")

		// send the response
		out <- release
	}
	i.err <- fmt.Errorf("Nothing happened (not really)\n")

	fmt.Printf("Done Releasing\n")
}

// PseudoTerminal implements pty
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

// NonInteract implements non-pty
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

// Close closes the readers/writers
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

// Report implements the Reporter interface
func (i *Interaction) Report(ctx context.Context, err chan<- error) {
	for {
		select {
		case msg := <-i.err:
			if msg == nil {
				close(err)
				return
			}
			err <- msg
		}
	}
}
