// Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/constants"
	viclog "github.com/vmware/vic/pkg/log"
	"github.com/vmware/vic/pkg/log/syslog"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

var (
	vchConfig config.VirtualContainerHostConfigSpec
)

const (
	KubeletConfigFile = "/etc/kubelet.conf"
	KubeletPath       = "/sbin/virtual-kubelet"
)

func main() {
	op := trace.NewOperation(context.Background(), "kubelet-starter")

	// load the vch config
	src, err := extraconfig.GuestInfoSource()
	if err != nil {
		op.Fatalf("Unable to load configuration from guestinfo: %s", err)
	}
	extraconfig.Decode(src, &vchConfig)

	logcfg := viclog.NewLoggingConfig()
	if vchConfig.Diagnostics.DebugLevel > 0 {
		logcfg.Level = log.DebugLevel
		trace.Logger.Level = log.DebugLevel
		syslog.Logger.Level = log.DebugLevel
	}

	if vchConfig.Diagnostics.DebugLevel > 3 {
		// extraconfig is very, very verbose
		extraconfig.SetLogLevel(log.DebugLevel)
	}

	if vchConfig.Diagnostics.SysLogConfig != nil {
		logcfg.Syslog = &viclog.SyslogConfig{
			Network:  vchConfig.Diagnostics.SysLogConfig.Network,
			RAddr:    vchConfig.Diagnostics.SysLogConfig.RAddr,
			Priority: syslog.Info | syslog.Daemon,
		}
	}

	op.Infof("%+v", *logcfg)
	// #nosec: Errors unhandled.
	viclog.Init(logcfg)
	trace.InitLogger(logcfg)

	if vchConfig.Diagnostics.DebugLevel > 2 {
		// expose portlayer service on client interface
		plPort := constants.DebugPortLayerPort
		os.Setenv("PORTLAYER_ADDR", strconv.Itoa(plPort))
	}

	op.Infof("KUBELET_NAME = %s", os.Getenv("KUBELET_NAME"))
	op.Infof("KUBERNETES_SERVICE_HOST = %s", os.Getenv("KUBERNETES_SERVICE_HOST"))
	op.Infof("KUBERNETES_SERVICE_PORT = %s", os.Getenv("KUBERNETES_SERVICE_PORT"))
	op.Infof("PERSONA_ADDR = %s", os.Getenv("PERSONA_ADDR"))
	op.Infof("PORTLAYER_ADDR = %s", os.Getenv("PORTLAYER_ADDR"))

	// Create Kubelet config file
	content := []byte(vchConfig.Kubelet.KubeletConfigContent)
	err = ioutil.WriteFile(KubeletConfigFile, content, 0644)

	if err != nil {
		op.Fatalf("Cannot write Kubelet config file to: %s, err: %s", KubeletConfigFile, err)
	}

	kubeletName := os.Getenv("KUBELET_NAME")

	op.Infof("Executing kubelet: %s %s %s %s %s %s %s", KubeletPath, "--provider", "mock", "--kubeconfig", KubeletConfigFile, "--nodename", kubeletName)
	/* #nosec */
	kubeletCmd := exec.Command(KubeletPath, "--provider", "mock", "--kubeconfig", KubeletConfigFile, "--nodename", kubeletName)
	output, err := kubeletCmd.CombinedOutput()
	op.Infof("Output: %s, Error: %s", string(output), err)
}
