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

package registry

import (
	"net"
	"net/url"

	glob "github.com/ryanuber/go-glob"
)

type Entry interface {
	Contains(e Entry) bool
	Match(e string) bool
	Equal(other Entry) bool
	String() string
}

func ParseEntry(s string) Entry {
	ip := net.ParseIP(s)
	if ip != nil {
		return &ipEntry{e: s}
	}

	_, ipnet, err := net.ParseCIDR(s)
	if err == nil {
		return &cidrEntry{ipnet: ipnet}
	}

	// check if url
	u, err := url.Parse(s)
	if err == nil {
		// only use the hostname
		h := u.Hostname()
		if len(h) > 0 {
			return &domainEntry{e: h}
		}
	}

	// assume domain name
	return &domainEntry{e: s}
}

type ipEntry struct {
	e string
}

func (e *ipEntry) Contains(other Entry) bool {
	return e.Match(other.String())
}

func (e *ipEntry) Match(s string) bool {
	return e.e == s || e.e+"/32" == s
}

func (e *ipEntry) Equal(other Entry) bool {
	return e.Match(other.String())
}

func (e *ipEntry) String() string {
	return e.e
}

type cidrEntry struct {
	ipnet *net.IPNet
}

func (c *cidrEntry) Contains(e Entry) bool {
	if c.Match(e.String()) {
		return true
	}

	if e, ok := e.(*cidrEntry); ok {
		return c.ipnet.Contains(e.ipnet.IP.Mask(e.ipnet.Mask))
	}

	return false
}

func (c *cidrEntry) Match(s string) bool {
	ip := net.ParseIP(s)
	if ip != nil {
		return c.ipnet.Contains(ip)
	}

	return false
}

func (c *cidrEntry) Equal(other Entry) bool {
	return other.String() == c.ipnet.String()
}

func (c *cidrEntry) String() string {
	return c.ipnet.String()
}

type domainEntry struct {
	e string
}

func (w *domainEntry) Contains(e Entry) bool {
	return w.Match(e.String())
}

func (w *domainEntry) Match(s string) bool {
	return glob.Glob(w.e, s)
}

func (w *domainEntry) String() string {
	return w.e
}

func (w *domainEntry) Equal(other Entry) bool {
	return other.String() == w.String()
}
