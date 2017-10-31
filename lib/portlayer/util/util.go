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

package util

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/vic/lib/constants"
	"github.com/vmware/vic/lib/spec"
)

const (
	// XXX leaving this as http for now.  We probably want to make this unix://
	scheme = "http://"
)

var (
	DefaultHost = Host()
)

func Host() *url.URL {
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	thisHost, err := url.Parse(scheme + name)
	if err != nil {
		log.Fatal(err)
	}

	return thisHost
}

// ServiceURL returns the URL for a given service relative to this host
func ServiceURL(serviceName string) *url.URL {
	s, err := DefaultHost.Parse(serviceName)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

// Update the VM display name on vSphere UI
func DisplayName(config *spec.VirtualMachineConfigSpecConfig, namingConvention string) string {
	var name string
	shortID := config.ID[:constants.ShortIDLen]
	nameMaxLen := constants.MaxVMNameLength - len(shortID)
	prettyName := config.Name

	if namingConvention != "" {
		//TODO: need to respect max length -- potentially enforce '-' separator
		if strings.Contains(namingConvention, "{id}") {
			name = strings.Replace(namingConvention, "{id}", shortID, -1)
		} else {
			//assume name
			name = strings.Replace(namingConvention, "{name}", prettyName, -1)
		}
		log.Infof("Applied naming convention: %s resulting %s", namingConvention, name)
	} else {
		// no naming convention specified during VCH creation / reconfigure so apply
		// standard VM display name convention
		if len(prettyName) > nameMaxLen-1 {
			name = fmt.Sprintf("%s-%s", prettyName[:nameMaxLen-1], shortID)
		} else {
			name = fmt.Sprintf("%s-%s", prettyName, shortID)
		}
	}
	return name
}

func ClientIP() (net.IP, error) {
	ips, err := net.LookupIP(constants.ClientHostName)
	if err != nil {
		return nil, err
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("No IP found on %s", constants.ClientHostName)
	}

	if len(ips) > 1 {
		return nil, fmt.Errorf("Multiple IPs found on %s: %#v", constants.ClientHostName, ips)
	}
	return ips[0], nil
}
