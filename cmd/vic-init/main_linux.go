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

package main

import (
	"os"
	"runtime/debug"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/vishvananda/netlink"
	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
	"github.com/vmware/vic/pkg/vsphere/toolbox"
)

var tthr tether.Tether

func init() {
	trace.Logger.Level = log.DebugLevel
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("run time panic: %s : %s", r, debug.Stack())
		}

		reboot()
	}()

	logFile, err := os.OpenFile("/dev/ttyS1", os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		log.Errorf("Could not pipe stderr to serial for debugging info. Some debug info may be lost! Error reported was %s", err)
	}
	err = syscall.Dup3(int(logFile.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		log.Errorf("Could not pipe logfile to standard error due to error %s", err)
	}

	_, err = os.Stderr.WriteString("all stderr redirected to debug log")
	if err != nil {
		log.Errorf("Could not write to Stderr due to error %s", err)
	}
	if strings.HasSuffix(os.Args[0], "-debug") {
		extraconfig.DecodeLogLevel = log.DebugLevel
		extraconfig.EncodeLogLevel = log.DebugLevel
	}
	log.SetLevel(log.DebugLevel)

	src, err := extraconfig.GuestInfoSourceWithPrefix("init")
	if err != nil {
		log.Error(err)
		return
	}

	sink, err := extraconfig.GuestInfoSinkWithPrefix("init")
	if err != nil {
		log.Error(err)
		return
	}

	// create the tether
	tthr = tether.New(src, sink, &operations{})

	// register the toolbox extension and configure for appliance
	toolbox := configureToolbox(tether.NewToolbox())
	toolbox.PrimaryIP = externalIP
	tthr.Register("Toolbox", toolbox)

	err = tthr.Start()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Clean exit from init")
}

// exit cleanly shuts down the system
func halt() {
	log.Infof("Powering off the system")
	if strings.HasSuffix(os.Args[0], "-debug") {
		log.Info("Squashing power off for debug init")
		return
	}

	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}

func reboot() {
	log.Infof("Rebooting the system")
	if strings.HasSuffix(os.Args[0], "-debug") {
		log.Info("Squashing reboot for debug init")
		return
	}

	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

func configureToolbox(t *tether.Toolbox) *tether.Toolbox {
	vix := t.Service.VixCommand
	vix.ProcessStartCommand = startCommand

	return t
}

// externalIP attempts to find an external IP to be reported as the guest IP
func externalIP() string {
	l, err := netlink.LinkByName("client")
	if err != nil && !os.IsNotExist(err) {
		log.Errorf("error looking up client interface: %s", err)
		return ""
	}

	if l == nil {
		l, err = netlink.LinkByAlias("client")
		if err != nil {
			log.Errorf("error looking up client interface: %s", err)
			return ""
		}
	}

	addrs, err := netlink.AddrList(l, netlink.FAMILY_V4)
	if err != nil {
		log.Errorf("error getting address list for client interface: %s", err)
		return ""
	}

	if len(addrs) == 0 {
		log.Warnf("no addresses set on client interface")
		return ""
	}

	return addrs[0].IP.String()
}

// defaultIP tries externalIP, falling back to toolbox.DefaultIP()
func defaultIP() string {
	ip := externalIP()
	if ip != "" {
		return ip
	}

	return toolbox.DefaultIP()
}
