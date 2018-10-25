// Copyright 2017 VMware, Inc. All Rights Reserved.
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

package common

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/vic/pkg/flags"
)

type Proxies struct {
	HTTPSProxy *string
	HTTPProxy  *string
	NoProxy    *string
	IsSet      bool
}

func (p *Proxies) ProxyFlags() []cli.Flag {
	return []cli.Flag{
		// proxies
		cli.GenericFlag{
			Name:   "https-proxy",
			Value:  flags.NewOptionalString(&p.HTTPSProxy),
			Usage:  "An HTTPS proxy for use when fetching images, in the form http(s)://fqdn_or_ip:port",
			Hidden: true,
		},
		cli.GenericFlag{
			Name:   "http-proxy",
			Value:  flags.NewOptionalString(&p.HTTPProxy),
			Usage:  "An HTTP proxy for use when fetching images, in the form http(s)://fqdn_or_ip:port",
			Hidden: true,
		},
		cli.GenericFlag{
			Name:   "no-proxy",
			Value:  flags.NewOptionalString(&p.NoProxy),
			Usage:  "URLs that should be excluded from proxying. This should be * or a comma-separated list of hostnames, domain names, or a mixture of both, e.g. localhost,.example.com",
			Hidden: true,
		},
	}
}

func (p *Proxies) ProcessProxies() (hproxy, sproxy *url.URL, nproxy *string, err error) {
	if p.HTTPProxy != nil || p.HTTPSProxy != nil {
		p.IsSet = true
	}

	if p.HTTPProxy != nil && *p.HTTPProxy != "" {
		hproxy, err = p.validate(*p.HTTPProxy)
		if err != nil {
			return
		}
	}

	if p.HTTPSProxy != nil && *p.HTTPSProxy != "" {
		sproxy, err = p.validate(*p.HTTPSProxy)
		if err != nil {
			return
		}
	}

	if p.NoProxy != nil && *p.NoProxy != "" {
		nproxy, err = p.validateURIs(strings.Split(*p.NoProxy, ","))
	}
	return
}

func (p *Proxies) validate(ref string) (proxy *url.URL, err error) {
	proxy, err = url.Parse(ref)
	if err != nil {
		return
	}
	if proxy.Host == "" || (proxy.Scheme != "http" && proxy.Scheme != "https") {
		err = cli.NewExitError(fmt.Sprintf("Could not parse HTTP(S) proxy - expected format http(s)://fqnd_or_ip:port: %s", ref), 1)
	}
	return
}

func (p *Proxies) validateURIs(refs []string) (trimedNproxy *string, err error) {
	var buffer bytes.Buffer
	buffer.WriteString(strings.TrimSpace(refs[0]))
	for _, ref := range refs[1:] {
		trimedRef := strings.TrimSpace(ref)
		_, err = url.Parse(trimedRef)
		if err != nil {
			err = cli.NewExitError(fmt.Sprintf("Could not parse no-proxy - expected format is * or a comma-separated list of hostnames, domain names, or a mixture of both, e.g. localhost,.example.com: %s", ref), 1)
			return
		}
		buffer.WriteString(",")
		buffer.WriteString(trimedRef)
	}
	trimedStr := buffer.String()
	trimedNproxy = &trimedStr
	return
}
