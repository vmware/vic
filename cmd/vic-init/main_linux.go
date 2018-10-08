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
	"context"
	"errors"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"

	"github.com/vmware/govmomi/toolbox"
	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/lib/tether/shared"
	viclog "github.com/vmware/vic/pkg/log"
	"github.com/vmware/vic/pkg/log/syslog"
	"github.com/vmware/vic/pkg/logmgr"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

var (
	tthr       tether.Tether
	config     ExecutorConfig
	debugLevel int
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("run time panic: %s : %s", r, debug.Stack())
		}

		reboot()
	}()

	// ensure that panics and error output are persisted
	logFile, err := os.OpenFile("/dev/ttyS0", os.O_WRONLY|os.O_SYNC, 0)
	if err != nil {
		log.Errorf("Could not redirect outputs to serial for debugging info, some debug info may be lost! Error reported was %s", err)
	}

	err = syscall.Dup3(int(logFile.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		log.Errorf("Could not pipe standard error to logfile: %s", err)
	}
	_, err = os.Stderr.WriteString("all stderr redirected to debug log")
	if err != nil {
		log.Errorf("Could not write to Stderr due to error %s", err)
	}

	err = syscall.Dup3(int(logFile.Fd()), int(os.Stdout.Fd()), 0)
	if err != nil {
		log.Errorf("Could not pipe standard out to logfile: %s", err)
	}
	_, err = os.Stderr.WriteString("all stdout redirected to debug log")
	if err != nil {
		log.Errorf("Could not write to stdout due to error %s", err)
	}

	src, err := extraconfig.GuestInfoSourceWithPrefix("init")
	if err != nil {
		log.Fatal(err)
	}

	extraconfig.Decode(src, &config)
	debugLevel = config.Diagnostics.DebugLevel

	startSignalHandler()

	logcfg := viclog.NewLoggingConfig()
	if debugLevel > 0 {
		logcfg.Level = log.DebugLevel
		trace.Logger.Level = log.DebugLevel
		syslog.Logger.Level = log.DebugLevel
	}

	if config.Diagnostics.SysLogConfig != nil {
		logcfg.Syslog = &viclog.SyslogConfig{
			Network:  config.Diagnostics.SysLogConfig.Network,
			RAddr:    config.Diagnostics.SysLogConfig.RAddr,
			Priority: syslog.Info | syslog.Daemon,
		}
	}

	viclog.Init(logcfg)

	if debugLevel > 2 {
		enableShell()
	}

	sink, err := extraconfig.GuestInfoSinkWithPrefix("init")
	if err != nil {
		log.Fatal(err)
	}

	// create the tether
	tthr = tether.New(src, sink, &operations{})

	// register the toolbox extension and configure for appliance
	toolbox := configureToolbox(tether.NewToolbox())
	toolbox.PrimaryIP = externalIP
	tthr.Register("Toolbox", toolbox)

	// Check logs every 5 minutes and rotate them if their size exceeds 20MB.
	// The history size we keep is 2 previous files in a compressed form.
	// TODO: Check available memory to tune log size and history length for log files.
	logrotate, err := logmgr.NewLogManager(time.Second * 300)
	const maxLogSizeBytes = 20 * 1024 * 1024
	if err == nil {
		logrotate.AddLogRotate("/var/log/vic/port-layer.log", logmgr.Daily, maxLogSizeBytes, 2, true)
		logrotate.AddLogRotate("/var/log/vic/init.log", logmgr.Daily, maxLogSizeBytes, 2, true)
		logrotate.AddLogRotate("/var/log/vic/docker-personality.log", logmgr.Daily, maxLogSizeBytes, 2, true)
		logrotate.AddLogRotate("/var/log/vic/vicadmin.log", logmgr.Daily, maxLogSizeBytes, 2, true)
		tthr.Register("logrotate", logrotate)
	} else {
		log.Error(err)
	}

	err = tthr.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Clean exit from init")
}

// exitTether signals the current process, which triggers tether.Stop and the killing of its children.
// NOTE: I don't like having this here and it really needs to be moved into an interface that
// can be provided to toolbox for system callbacks. While this could be part of the Operations
// interface I think I'd rather have a separate one specifically for the possible toolbox interactions.
func exitTether() error {
	defer trace.End(trace.Begin(""))

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}

	if err = p.Signal(syscall.SIGUSR2); err != nil {
		return err
	}

	return err
}

// exit cleanly shuts down the system
func halt() error {
	log.Infof("Powering off the system")

	err := exitTether()
	if err != nil {
		log.Warn(err)
	}

	if debugLevel > 2 {
		log.Info("Squashing power off for debug init")
		return errors.New("debug config suppresses shutdown")
	}

	timeout, cancel := context.WithTimeout(context.Background(), shared.GuestShutdownTimeout)
	err = tthr.Wait(timeout)
	cancel()

	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)

	return err
}

func reboot() error {
	log.Infof("Rebooting the system")

	err := exitTether()
	if err != nil {
		log.Warn(err)
	}

	if debugLevel > 2 {
		log.Info("Squashing reboot for debug init")
		return errors.New("debug config suppresses reboot")
	}

	timeout, cancel := context.WithTimeout(context.Background(), shared.GuestRebootTimeout)
	err = tthr.Wait(timeout)
	cancel()

	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)

	return err
}

func configureToolbox(t *tether.Toolbox) *tether.Toolbox {
	cmd := t.Service.Command
	cmd.ProcessStartCommand = startCommand

	t.Power.Halt.Handler = halt
	t.Power.Reboot.Handler = reboot
	t.Power.Suspend.Handler = exitTether

	return t
}

// externalIP attempts to find an external IP to be reported as the guest IP
func externalIP() string {
	l, err := netlink.LinkByName("client")
	if err != nil {
		log.Debugf("error looking up client interface by name: %s", err)
		l, err = netlink.LinkByAlias("client")
		if err != nil {
			log.Errorf("error looking up client interface by alias: %s", err)
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

// This code is mirrored in cmd/tether/main_linux.go and should be de-duped
func startSignalHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGPWR, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for s := range sigs {
			switch s {
			case syscall.SIGHUP:
				log.Infof("Reloading tether configuration")
				tthr.Reload()
			case syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGPWR:
				log.Infof("Stopping tether via signal %s", s.String())
				tthr.Stop()
				return
			case syscall.SIGTERM, syscall.SIGINT:
				log.Infof("Stopping system in lieu of restart handling via signal %s", s.String())
				// TODO: update this to adjust power off handling for reboot
				// this should be in guest reboot rather than power cycle
				tthr.Stop()
				return
			default:
				log.Infof("%s signal not defined", s.String())
			}
		}
	}()
}
