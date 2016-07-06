// Copyright 2016 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pprof

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/vsphere/extraconfig"
)

type PprofPort int

const basePort = 6060

const (
	VCHInitPort PprofPort = iota
	VicadminPort
	DockerPort
	PortlayerPort
	maxPort
)

var (
	vchConfig metadata.VirtualContainerHostConfigSpec
)

func init() {
	// load the vch config
	// TODO: Optimize this to just pull the fields we need...
	src, err := extraconfig.GuestInfoSource()
	if err != nil {
		log.Errorf("Unable to load configuration from guestinfo")
		return
	}
	extraconfig.Decode(src, &vchConfig)
}

func GetPprofEndpoint(component PprofPort) *url.URL {
	if component >= maxPort {
		return nil
	}
	port := component + basePort

	ip := "127.0.0.1"
	if vchConfig.ExecutorConfig.Diagnostics.DebugLevel > 1 {
		ips, err := net.LookupIP("client.localhost")
		if err != nil || len(ips) == 0 {
			log.Warnf("Unable to resolve 'client.localhost': ", err)
		} else {
			ip = ips[0].String()
		}
	}

	endpoint, err := url.Parse(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil
	}
	return endpoint
}

func StartPprof(name string, component PprofPort) error {
	url := GetPprofEndpoint(component)
	if url == nil {
		err := errors.New(fmt.Sprintf("Unable to get pprof endpoint for %s.", name))
		log.Error(err.Error())
		return err
	}

	log.Info(fmt.Sprintf("Launching %s server on %s", name, url.String()))
	go func() {
		log.Info(http.ListenAndServe(url.String(), nil))
	}()

	return nil
}
