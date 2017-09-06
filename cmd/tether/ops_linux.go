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
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"

	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/constants"
	"github.com/vmware/vic/lib/iolog"
	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/lib/tether/netfilter"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/trace"
)

const (
	runMountPoint = "/run"

	// default values to set for ulimit fields
	defaultNOFILE = 1024 * 1024
	defaultULimit = math.MaxUint64
)

type operations struct {
	tether.BaseOperations

	logging bool
}

func (t *operations) Log() (io.Writer, error) {
	defer trace.End(trace.Begin("operations.Log"))

	// redirect logging to the serial log
	log.Infof("opening %s/ttyS1 for debug log", pathPrefix)
	f, err := os.OpenFile(path.Join(pathPrefix, "ttyS1"), os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, 0)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for debug log: %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	if err := setTerminalSpeed(f.Fd()); err != nil {
		log.Errorf("Setting terminal speed failed with %s", err)
	}

	// enable raw mode
	_, err = terminal.MakeRaw(int(f.Fd()))
	if err != nil {
		detail := fmt.Sprintf("Making ttyS1 raw failed with %s", err)
		log.Error(detail)
		return nil, errors.New(detail)
	}

	return io.MultiWriter(f, os.Stdout), nil
}

// sessionLogWriter returns a writer that will persist the session output
func (t *operations) SessionLog(session *tether.SessionConfig) (dio.DynamicMultiWriter, dio.DynamicMultiWriter, error) {
	defer trace.End(trace.Begin("configure session log writer"))

	if t.logging {
		detail := "unable to log more than one session concurrently to persistent logging"
		log.Warn(detail)
		// use multi-writer so it's still viable for attach
		return dio.MultiWriter(), dio.MultiWriter(), nil
	}

	t.logging = true

	// open SttyS2 for session logging
	log.Info("opening ttyS2 for session logging")
	f, err := os.OpenFile(path.Join(pathPrefix, "ttyS2"), os.O_RDWR|os.O_SYNC|syscall.O_NOCTTY, 0)
	if err != nil {
		detail := fmt.Sprintf("failed to open serial port for session log: %s", err)
		log.Error(detail)
		return nil, nil, errors.New(detail)
	}

	if err := setTerminalSpeed(f.Fd()); err != nil {
		log.Errorf("Setting terminal speed failed with %s", err)
	}

	// enable raw mode
	_, err = terminal.MakeRaw(int(f.Fd()))
	if err != nil {
		detail := fmt.Sprintf("Making ttyS2 raw failed with %s", err)
		log.Error(detail)
		return nil, nil, errors.New(detail)
	}

	// wrap output in a LogWriter to serialize it into our persisted
	// containerVM output format, using iolog.LogClock for timestamps
	lw := iolog.NewLogWriter(f, iolog.LogClock{})

	// use multi-writer so it goes to both screen and session log
	return dio.MultiWriter(lw, os.Stdout), dio.MultiWriter(lw, os.Stderr), nil
}

func (t *operations) Setup(sink tether.Config) error {
	if err := t.BaseOperations.Setup(sink); err != nil {
		return err
	}

	// unmount /run - https://github.com/vmware/vic/issues/1643
	if err := tether.Sys.Syscall.Unmount(runMountPoint, syscall.MNT_DETACH); err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EINVAL {
			return err
		}
	}

	// NOTE: ulimit default values should change when we support ulimit configuration
	ApplyDefaultULimit()

	return nil
}

// ApplyDefaultULimit sets ulimit fields to their defined default value
func ApplyDefaultULimit() {
	var rLimit syscall.Rlimit

	// NOFILE does not support defaultULimit as a value due to kernel restriction on number of open files
	rLimit.Max = defaultNOFILE
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Errorf("Cannot set ulimit for nofile: %s", err.Error())
	}

	rLimit.Max = defaultULimit
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_STACK, &rLimit); err != nil {
		log.Errorf("Cannot set ulimit for stack: %s ", err.Error())
	}

	if err := syscall.Setrlimit(syscall.RLIMIT_CORE, &rLimit); err != nil {
		log.Errorf("Cannot set ulimit for core blocks: %s", err.Error())
	}

	if err := syscall.Setrlimit(unix.RLIMIT_MEMLOCK, &rLimit); err != nil {
		log.Errorf("Cannot set ulimit for memlock: %s", err.Error())
	}

	if err := syscall.Setrlimit(unix.RLIMIT_NPROC, &rLimit); err != nil {
		log.Errorf("Cannot set ulimit for nproc: %s", err.Error())
	}
}

// invoke will invoke the closure returned from the tether netfilter prep and
// block until complete. It handles both potential preparation errors and invocation
// errors.
// the 'task' specified is used to construct error messages with the specific operation
// embedded
func invoke(t *tether.BaseOperations, fn tether.UtilityFn, task string) error {
	exitChan, err := t.LaunchUtility(fn)
	if err != nil {
		return fmt.Errorf("%s failed: %s", task, err)
	}

	exitCode := <-exitChan
	if exitCode != 0 {
		return fmt.Errorf("%s returned non-zero: %d", task, exitCode)
	}

	return nil
}

// SetupFirewall sets up firewall rules on the external scope only.  Any
// portmaps are honored as are port exposes.
func (t *operations) SetupFirewall(ctx context.Context, config *tether.ExecutorConfig) error {
	return setupFirewall(ctx, &t.BaseOperations, config)
}

// setupFirewall is broken out from SetupFirewall so that it can be referenced from the test code
func setupFirewall(ctx context.Context, t *tether.BaseOperations, config *tether.ExecutorConfig) error {
	fn := netfilter.Flush(ctx, "VIC")
	if err := invoke(t, fn, "flush"); err != nil {
		return err
	}

	for _, endpoint := range config.Networks {
		switch endpoint.Network.Type {
		case constants.ExternalScopeType:
			id, err := strconv.Atoi(endpoint.ID)
			if err != nil {
				log.Errorf("can't apply port rules: %s", err.Error())
				continue
			}

			iface, err := t.LinkBySlot(int32(id))
			if err != nil {
				log.Errorf("can't apply rules: %s", err.Error())
				continue
			}

			ifaceName := iface.Attrs().Name
			log.Debugf("slot %d -> %s", endpoint.ID, ifaceName)

			switch endpoint.Network.TrustLevel {
			case executor.Open:
				// Accept all incoming and outgoing traffic
				for _, chain := range []netfilter.Chain{netfilter.Input, netfilter.Output, netfilter.Forward} {
					fn := (&netfilter.Rule{
						Chain:     chain,
						Target:    netfilter.Accept,
						Interface: ifaceName,
					}).Commit(ctx)
					if err := invoke(t, fn, "accept all"); err != nil {
						return err
					}
				}

			case executor.Closed:
				// Reject all incoming and outgoing traffic
				// Since our default policy is to drop traffic, nothing is needed here.

			case executor.Outbound:
				// Reject all incoming traffic, but allow outgoing
				if err := setupOutboundFirewall(ctx, t, ifaceName); err != nil {
					return err
				}

			case executor.Peers:
				// Outbound + all ports open to source addresses in --container-network-ip-range
				if err := setupOutboundFirewall(ctx, t, ifaceName); err != nil {
					return err
				}
				sourceAddresses := make([]string, len(endpoint.Network.Pools))
				for i, v := range endpoint.Network.Pools {
					sourceAddresses[i] = v.String()
				}
				fn := (&netfilter.Rule{
					Chain:           netfilter.Input,
					Target:          netfilter.Accept,
					SourceAddresses: sourceAddresses,
					Interface:       ifaceName,
				}).Commit(ctx)
				if err := invoke(t, fn, "allow outbound and peers"); err != nil {
					return err
				}
				if err := allowPingTraffic(ctx, t, ifaceName, sourceAddresses); err != nil {
					return err
				}

			case executor.Published:
				if err := setupPublishedFirewall(ctx, t, endpoint, ifaceName); err != nil {
					return err
				}

			case executor.Unspecified:
				// Unspecified network firewalls default to published.
				if err := setupPublishedFirewall(ctx, t, endpoint, ifaceName); err != nil {
					return err
				}

			default:
				log.Warningf("Received invalid firewall configuration %v: defaulting to published.",
					endpoint.Network.TrustLevel)
				if err := setupPublishedFirewall(ctx, t, endpoint, ifaceName); err != nil {
					return err
				}
			}
		case constants.BridgeScopeType:
			id, err := strconv.Atoi(endpoint.ID)
			if err != nil {
				log.Errorf("can't apply port rules: %s", err.Error())
				continue
			}

			iface, err := t.LinkBySlot(int32(id))
			if err != nil {
				log.Errorf("can't apply rules: %s", err.Error())
				continue
			}

			ifaceName := iface.Attrs().Name
			log.Debugf("slot %d -> %s", endpoint.ID, ifaceName)

			// Traffic over container bridge network should be peers+outbound.
			if err := setupOutboundFirewall(ctx, t, ifaceName); err != nil {
				return err
			}
			sourceAddresses := make([]string, len(endpoint.Network.Pools))

			for i, v := range endpoint.Network.Pools {
				sourceAddresses[i] = v.String()
			}
			fn := (&netfilter.Rule{
				Chain:           netfilter.Input,
				Target:          netfilter.Accept,
				SourceAddresses: sourceAddresses,
				Interface:       ifaceName,
			}).Commit(ctx)
			if err := invoke(t, fn, "configure for bridge scope"); err != nil {
				return err
			}
		}
	}

	return invoke(t, netfilter.Return(ctx, "VIC"), "return from VIC chain")
}

func setupOutboundFirewall(ctx context.Context, t *tether.BaseOperations, ifaceName string) error {
	// All already established inputs are accepted
	fn := (&netfilter.Rule{
		Chain:     netfilter.Input,
		States:    []netfilter.State{netfilter.Established, netfilter.Related},
		Target:    netfilter.Accept,
		Interface: ifaceName,
	}).Commit(ctx)
	if err := invoke(t, fn, "permit established inbound"); err != nil {
		return err
	}

	// All output is accepted
	fn = (&netfilter.Rule{
		Chain:     netfilter.Output,
		Target:    netfilter.Accept,
		Interface: ifaceName,
	}).Commit(ctx)
	return invoke(t, fn, "permit all outbound")
}

func setupPublishedFirewall(ctx context.Context, t *tether.BaseOperations, endpoint *tether.NetworkEndpoint, ifaceName string) error {
	if err := setupOutboundFirewall(ctx, t, ifaceName); err != nil {
		return err
	}
	if err := allowPingTraffic(ctx, t, ifaceName, nil); err != nil {
		return err
	}

	// handle the ports
	for _, p := range endpoint.Ports {
		// parse the port maps
		r, err := portToRule(p)
		if err != nil {
			log.Errorf("can't apply port rule (%s): %s", p, err.Error())
			continue
		}

		log.Infof("Applying rule for port %s", p)
		r.Interface = ifaceName
		fn := r.Commit(ctx)
		if err := invoke(t, fn, "allow incoming on published port"); err != nil {
			return err
		}
	}

	return nil
}

func allowPingTraffic(ctx context.Context, t *tether.BaseOperations, ifaceName string, sourceAddresses []string) error {
	fn := (&netfilter.Rule{
		Chain:           netfilter.Input,
		Target:          netfilter.Accept,
		Interface:       ifaceName,
		Protocol:        netfilter.ICMP,
		ICMPType:        netfilter.EchoRequest,
		SourceAddresses: sourceAddresses,
	}).Commit(ctx)
	if err := invoke(t, fn, "allow ping inbound"); err != nil {
		return err
	}

	fn = (&netfilter.Rule{
		Chain:           netfilter.Output,
		Target:          netfilter.Accept,
		Interface:       ifaceName,
		Protocol:        netfilter.ICMP,
		ICMPType:        netfilter.EchoReply,
		SourceAddresses: sourceAddresses,
	}).Commit(ctx)
	return invoke(t, fn, "allow ping outbound")
}

func portToRule(p string) (*netfilter.Rule, error) {
	if strings.Contains(p, ":") {
		return nil, errors.New("port maps are TBD")
	}

	// 9999/tcp
	s := strings.Split(p, "/")
	if len(s) != 2 {
		return nil, errors.New("can't parse port spec: " + p)
	}

	rule := &netfilter.Rule{
		Chain:     netfilter.Input,
		Interface: "external",
		Target:    netfilter.Accept,
	}

	switch netfilter.Protocol(s[1]) {
	case netfilter.UDP:
		rule.Protocol = netfilter.UDP
	case netfilter.TCP:
		rule.Protocol = netfilter.TCP

	default:
		return nil, errors.New("unknown protocol")
	}

	port, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, err
	}

	rule.FromPort = port

	return rule, nil
}
