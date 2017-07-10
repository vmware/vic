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

package netfilter

import (
	"context"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"syscall"

	"github.com/Sirupsen/logrus"
)

//
// # default
// # iptables -A INPUT -m state --state ESTABLISHED -i eth0 -j ACCEPT

// # Expose()
// # iptables -A INPUT -p tcp --dport 7 -i eth0 -j ACCEPT
//
// # default on non-bridge
// # iptables -A INPUT -i eth0 -j REJECT
//
// # iptables --list
//
// Chain INPUT (policy ACCEPT)
// target     prot opt source               destination
// ACCEPT     all  --  anywhere             anywhere             state ESTABLISHED
// ACCEPT     tcp  --  anywhere             anywhere             tcp dpt:echo
// REJECT     all  --  anywhere             anywhere             reject-with icmp-port-unreachable
//
// Chain FORWARD (policy ACCEPT)
// target     prot opt source               destination
//
// Chain OUTPUT (policy ACCEPT)
// target     prot opt source               destination
//

type Chain string
type State string
type Protocol string
type Target string
type Table string
type ICMPType string

const (
	Prerouting = Chain("PREROUTING")
	Input      = Chain("INPUT")
	Forward    = Chain("FORWARD")
	Output     = Chain("OUTPUT")

	Invalid     = State("INVALID")
	Established = State("ESTABLISHED")
	New         = State("NEW")
	Related     = State("RELATED")
	Untracked   = State("UNTRACKED")

	TCP  = Protocol("tcp")
	UDP  = Protocol("udp")
	ICMP = Protocol("icmp")

	Drop     = Target("DROP")
	Accept   = Target("ACCEPT")
	Reject   = Target("REJECT")
	Redirect = Target("REDIRECT")

	NAT = Table("nat")

	EchoRequest = ICMPType("echo-request")
	EchoReply   = ICMPType("echo-reply")
)

type Rule struct {
	Table
	Chain
	States []State
	ICMPType

	Protocol
	Target

	Interface        string
	SourceAddresses  []string
	FromPort, ToPort int
}

func (r *Rule) Commit(ctx context.Context) error {
	args, err := r.args()
	if err != nil {
		return err
	}

	return iptables(ctx, args)
}

func (r *Rule) args() ([]string, error) {
	var args []string

	if r.Table != "" {
		args = append(args, "-t", string(r.Table))
	}

	args = append(args, "-A", string(r.Chain))

	if r.Protocol != "" {
		args = append(args, "-p", string(r.Protocol))
	}

	if r.ICMPType != "" {
		args = append(args, "--icmp-type", string(r.ICMPType))
	}

	if len(r.SourceAddresses) > 0 {
		args = append(args, "-s", strings.Join(r.SourceAddresses, ","))
	}

	if r.FromPort != 0 {
		args = append(args, "--dport", strconv.Itoa(r.FromPort))
	}

	if len(r.States) > 0 {
		args = append(args, "-m", "state", "--state", joinStates(r.States))
	}

	if r.Interface != "" {
		if r.Chain == Output {
			args = append(args, "-o", r.Interface)
		} else {
			args = append(args, "-i", r.Interface)
		}
	}

	if r.Target == "" {
		return nil, errors.New("target cannot be empty")
	}
	args = append(args, "-j", string(r.Target))

	if r.ToPort != 0 {
		args = append(args, "--to-port", strconv.Itoa(r.ToPort))
	}

	return args, nil
}

func iptables(ctx context.Context, args []string) error {
	logrus.Infof("Execing iptables %q", args)

	// #nosec: Subprocess launching with variable
	cmd := &exec.Cmd{
		Path: "/lib64/ld-linux-x86-64.so.2",
		Dir:  "/",
		Args: append([]string{"/lib64/ld-linux-x86-64.so.2", "/iptables"}, args...),
		SysProcAttr: &syscall.SysProcAttr{
			Chroot: "/.tether",
		},
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("iptables error: %s", err.Error())

		exitErr, ok := err.(*exec.ExitError)
		if ok && len(exitErr.Stderr) > 0 {
			logrus.Errorf("iptables error: %s", string(exitErr.Stderr))
		}
	}

	if len(b) > 0 {
		logrus.Infof("iptables: %s", string(b))
	}

	return err
}

func Flush(ctx context.Context, table string) error {
	args := []string{"-F"}
	if table != "" {
		args = append(args, "-t", table)
	}

	return iptables(ctx, args)
}

func Policy(ctx context.Context, chain Chain, target Target) error {
	return iptables(ctx, []string{"-P", string(chain), string(target)})
}

func joinStates(states []State) string {
	tmp := make([]string, len(states))
	for i, v := range states {
		tmp[i] = string(v)
	}
	return strings.Join(tmp, ",")
}
