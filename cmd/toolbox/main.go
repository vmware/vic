// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"

	viclog "github.com/vmware/vic/pkg/log"
	"github.com/vmware/vic/pkg/vsphere/toolbox"
)

// This example can be run on a VM hosted by ESX, Fusion or Workstation
func main() {
	log.SetFormatter(viclog.NewTextFormatter())

	flag.Parse()

	in := toolbox.NewBackdoorChannelIn()
	out := toolbox.NewBackdoorChannelOut()

	service := toolbox.NewService(in, out)

	// Trigger a command start, for example:
	// govc guest.start -vm vm-name kill SIGHUP
	service.VixCommand.ProcessStartCommand = func(r *toolbox.VixMsgStartProgramRequest) (int, error) {
		fmt.Fprintf(os.Stderr, "guest-command: %s %s\n", r.ProgramPath, r.Arguments)
		return -1, nil
	}

	if os.Getuid() == 0 {
		power := service.PowerCommand
		power.Halt.Handler = toolbox.Halt
		power.Reboot.Handler = toolbox.Reboot
	}

	err := service.Start()
	if err != nil {
		log.Fatal(err)
	}

	// handle the signals and gracefully shutdown the service
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("signal %s received", <-sig)
		service.Stop()
	}()

	service.Wait()
}
