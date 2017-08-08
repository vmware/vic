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
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/vmware/vic/lib/config/executor"
	"github.com/vmware/vic/lib/etcconf"
	"github.com/vmware/vic/lib/iolog"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/lib/tether"
	"github.com/vmware/vic/lib/tether/netfilter"
	"github.com/vmware/vic/pkg/dio"
	"github.com/vmware/vic/pkg/trace"
)

const (
	runMountPoint         = "/run"
	hostnameFileBindSrc   = "/.tether/etc/hostname"
	hostsPathBindSrc      = "/.tether/etc/hosts"
	resolvConfPathBindSrc = "/.tether/etc/resolv.conf"
)

var (
	filesForMinOSLinux = map[string]os.FileMode{
		"/etc/hostname":            0644,
		"/etc/hosts":               0644,
		"/etc/resolv.conf":         0644,
		"/.tether/etc/hostname":    0644,
		"/.tether/etc/hosts":       0644,
		"/.tether/etc/resolv.conf": 0644,
	}
)

type operations struct {
	tether.BaseOperations

	logging bool
}

func init() {
	tether.Sys.Hosts = etcconf.NewHosts(hostsPathBindSrc)
	tether.Sys.ResolvConf = etcconf.NewResolvConf(resolvConfPathBindSrc)
	tether.Sys.Hostname = hostnameFileBindSrc
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
	if err := createBindSrcTarget(filesForMinOSLinux); err != nil {
		return err
	}

	if err := t.BaseOperations.Setup(sink); err != nil {
		return err
	}

	// unmount /run - https://github.com/vmware/vic/issues/1643
	if err := tether.Sys.Syscall.Unmount(runMountPoint, syscall.MNT_DETACH); err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.EINVAL {
			return err
		}
	}

	return nil
}

// SetupFirewall sets up firewall rules on the external scope only.  Any
// portmaps are honored as are port exposes.
func (t *operations) SetupFirewall(config *tether.ExecutorConfig) error {
	// XXX It looks like we'd want to collect the errors here, but we
	// can't.  Since this is running inside init (tether) and tether
	// reaps all children, the os.exec package won't be able to collect
	// the error code in time before the reaper does.  The exec package
	// calls wait and attempts to collect its child, but the reaper will
	// have raptured the pid before that.  So, best effort, just keep going.
	if err := netfilter.Flush(context.Background(), ""); err != nil {
		return err
	}
	if err := generalPolicy(netfilter.Drop); err != nil {
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
					if err := (&netfilter.Rule{
						Chain:     chain,
						Target:    netfilter.Accept,
						Interface: ifaceName,
					}).Commit(context.TODO()); err != nil {
						return err
					}
				}

			case executor.Closed:
				// Reject all incoming and outgoing traffic
				// Since our default policy is to drop traffic, nothing is needed here.

			case executor.Outbound:
				// Reject all incoming traffic, but allow outgoing
				if err := setupOutboundFirewall(ifaceName); err != nil {
					return err
				}

			case executor.Peers:
				// Outbound + all ports open to source addresses in --container-network-ip-range
				if err := setupOutboundFirewall(ifaceName); err != nil {
					return err
				}
				sourceAddresses := make([]string, len(endpoint.Network.Pools))
				for i, v := range endpoint.Network.Pools {
					sourceAddresses[i] = v.String()
				}
				if err := (&netfilter.Rule{
					Chain:           netfilter.Input,
					Target:          netfilter.Accept,
					SourceAddresses: sourceAddresses,
					Interface:       ifaceName,
				}).Commit(context.TODO()); err != nil {
					return err
				}
				if err := allowPingTraffic(ifaceName, sourceAddresses); err != nil {
					return err
				}

			case executor.Published:
				if err := setupPublishedFirewall(endpoint, ifaceName); err != nil {
					return err
				}

			case executor.Unspecified:
				// Unspecified network firewalls default to published.
				if err := setupPublishedFirewall(endpoint, ifaceName); err != nil {
					return err
				}

			default:
				log.Warningf("Received invalid firewall configuration %v: defaulting to published.",
					endpoint.Network.TrustLevel)
				if err := setupPublishedFirewall(endpoint, ifaceName); err != nil {
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
			if err := setupOutboundFirewall(ifaceName); err != nil {
				return err
			}
			sourceAddresses := make([]string, len(endpoint.Network.Pools))

			for i, v := range endpoint.Network.Pools {
				sourceAddresses[i] = v.String()
			}
			if err := (&netfilter.Rule{
				Chain:           netfilter.Input,
				Target:          netfilter.Accept,
				SourceAddresses: sourceAddresses,
				Interface:       ifaceName,
			}).Commit(context.TODO()); err != nil {
				return err
			}
		}
	}

	return nil
}

func setupOutboundFirewall(ifaceName string) error {
	// All already established inputs are accepted
	if err := (&netfilter.Rule{
		Chain:     netfilter.Input,
		States:    []netfilter.State{netfilter.Established, netfilter.Related},
		Target:    netfilter.Accept,
		Interface: ifaceName,
	}).Commit(context.TODO()); err != nil {
		return err
	}
	// All output is accepted
	if err := (&netfilter.Rule{
		Chain:     netfilter.Output,
		Target:    netfilter.Accept,
		Interface: ifaceName,
	}).Commit(context.TODO()); err != nil {
		return err
	}
	return nil
}

func setupPublishedFirewall(endpoint *tether.NetworkEndpoint, ifaceName string) error {
	if err := setupOutboundFirewall(ifaceName); err != nil {
		return err
	}
	if err := allowPingTraffic(ifaceName, nil); err != nil {
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
		if err := r.Commit(context.TODO()); err != nil {
			return err
		}
	}
	return nil
}

func generalPolicy(target netfilter.Target) error {
	for _, chain := range []netfilter.Chain{netfilter.Input, netfilter.Output, netfilter.Forward} {
		if err := netfilter.Policy(context.TODO(), chain, target); err != nil {
			return err
		}
	}
	return nil
}

func allowPingTraffic(ifaceName string, sourceAddresses []string) error {
	if err := (&netfilter.Rule{
		Chain:           netfilter.Input,
		Target:          netfilter.Accept,
		Interface:       ifaceName,
		Protocol:        netfilter.ICMP,
		ICMPType:        netfilter.EchoRequest,
		SourceAddresses: sourceAddresses,
	}).Commit(context.TODO()); err != nil {
		return err
	}
	if err := (&netfilter.Rule{
		Chain:           netfilter.Output,
		Target:          netfilter.Accept,
		Interface:       ifaceName,
		Protocol:        netfilter.ICMP,
		ICMPType:        netfilter.EchoReply,
		SourceAddresses: sourceAddresses,
	}).Commit(context.TODO()); err != nil {
		return err
	}
	return nil
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

// Create necessary directories/files as the src/target for bind mount.
// See https://github.com/vmware/vic/issues/489
func createBindSrcTarget(files map[string]os.FileMode) error {
	log.Infof("set up bind mount src and target")

	// The directory has to exist before creating the new file
	for filePath, fmode := range files {
		dir := path.Dir(filePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// #nosec: Expect file permissions to be 0600 or less
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %s", dir, err)
			}
		}
		f, err := os.OpenFile(filePath, os.O_CREATE, fmode)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %s", filePath, err)
		}

		if err = f.Close(); err != nil {
			return fmt.Errorf("failed to close file %s: %s", filePath, err)
		}
	}

	return nil
}
