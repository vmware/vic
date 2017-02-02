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
	"context"
	"fmt"
	"time"

	tetherng "github.com/vmware/vic/lib/tether-ng"

	"github.com/vmware/vic/lib/tether-ng/interaction"
	"github.com/vmware/vic/lib/tether-ng/process"
	"github.com/vmware/vic/lib/tether-ng/types"
)

func main() {
	ctx := context.Background()

	tether := tetherng.NewTether(ctx)

	config := &types.ExecutorConfig{
		Sessions: map[string]*types.SessionConfig{
			"ls": &types.SessionConfig{
				Env:        []string{},
				Cmd:        []string{"/bin/ls", "-l"},
				WorkingDir: "/",

				Attach:    true,
				OpenStdin: false,
				RunBlock:  true,
				Tty:       false,
				Restart:   false,
			},
			"uname": &types.SessionConfig{
				Env:        []string{},
				Cmd:        []string{"/bin/uname"},
				WorkingDir: "/",

				Attach:    true,
				OpenStdin: false,
				RunBlock:  false,
				Tty:       false,
				Restart:   false,
			},
			"ps": &types.SessionConfig{
				Env:        []string{},
				Cmd:        []string{"/bin/ps"},
				WorkingDir: "/",

				Attach:    true,
				OpenStdin: false,
				RunBlock:  false,
				Tty:       false,
				Restart:   false,
			},
			"bash": &types.SessionConfig{
				Env:        []string{},
				Cmd:        []string{"/bin/bash"},
				WorkingDir: "/",

				Attach:    true,
				OpenStdin: true,
				RunBlock:  false,
				Tty:       true,
				Restart:   false,
			},
		},
	}

	callctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// order is important
	interaction := interaction.NewInteraction(callctx)
	process := process.NewProcess(callctx)

	process.SetReleaser(ctx, interaction)
	process.SetInteractor(ctx, interaction)

	if err := tether.Register(callctx, interaction); err != nil {
		fmt.Printf("Err: %s\n", err)
	}
	if err := tether.Register(callctx, process); err != nil {
		fmt.Printf("Err: %s\n", err)
	}

	for _, i := range tether.Plugins(callctx) {
		if err := i.Configure(callctx, config); err != nil {
			fmt.Printf("Err: %s\n", err)
		}
	}

	for _, i := range tether.Plugins(callctx) {
		if err := i.Start(callctx); err != nil {
			fmt.Printf("Err: %s\n", err)
		}
	}

	for _, i := range tether.Plugins(callctx) {
		fmt.Printf("%s\n", i.UUID(callctx))
	}

	for _, i := range tether.Plugins(callctx) {
		if s, ok := i.(tetherng.Signaler); ok {
			for {
				running := false
				for k := range config.Sessions {
					r := s.Running(ctx, k)
					if r {
						fmt.Printf("%s is running\n", k)
					}
					running = running || r
				}
				if !running {
					break
				}
				time.Sleep(5 * time.Second)
			}
		}
	}

	for _, i := range tether.Plugins(callctx) {
		if err := i.Stop(callctx); err != nil {
			fmt.Printf("Err: %s\n", err)
		}
	}
}
