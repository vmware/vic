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

	// imported for the side effect
	_ "expvar"
	"fmt"
	"net"
	"net/http"
	"time"
	// imported for the side effect
	_ "net/http/pprof"
	"net/url"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/config"
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
	debugLevel int
)

func init() {
	// load the vch config
	// TODO: Optimize this to just pull the fields we need...
	src, err := extraconfig.GuestInfoSource()
	if err != nil {
		log.Errorf("Unable to load configuration from guestinfo")
		return
	}

	vchConfig := new(config.VirtualContainerHostConfigSpec)
	extraconfig.Decode(src, vchConfig)
	debugLevel = vchConfig.ExecutorConfig.Diagnostics.DebugLevel
}

// GetPprofEndpoint determines on which interface the pprof server should be launched. This makes
// an explicit assumption that there will be a `client.localhost` name when running in debug modes.
func GetPprofEndpoint(component PprofPort) *url.URL {
	return getEndpoint(component, debugLevel)
}

func getEndpoint(component PprofPort, debug int) *url.URL {
	if component >= maxPort {
		return nil
	}
	port := component + basePort

	ip := "127.0.0.1"
	// exposing this data on an external port definitely counts as a change of behaviour,
	// so this is > 1, just debug on/off.
	if debug > 1 {
		log.Debugf("Waiting for client.localhost to resolve before launching pprof")
		for {
			ips, err := net.LookupIP("client.localhost")
			if err != nil || len(ips) == 0 {
				log.Warnf("Unable to resolve 'client.localhost': ", err)
				time.Sleep(time.Second)
				continue
			}

			ip = ips[0].String()
			break
		}
	}

	endpoint, err := url.Parse(fmt.Sprintf("http://%s:%d", ip, port))
	if err != nil {
		return nil
	}
	return endpoint
}

func StartPprof(name string, component PprofPort) error {
	go func() {
		url := getEndpoint(component, debugLevel)
		if url == nil {
			log.Errorf("Unable to get pprof endpoint for %s, reverting to localhost.", name)
			url = getEndpoint(component, 0)
			if url == nil {
				log.Errorf("Unable to construct localhost pprof endpoint for %s, giving up", name)
				return
			}
		}

		location := url.String()[7:] // Strip off leading "http://"

		log.Info(fmt.Sprintf("Launching %s pprof server on %s", name, location))
		log.Info(http.ListenAndServe(location, nil))
	}()

	return nil
}
