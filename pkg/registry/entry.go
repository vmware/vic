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

	// assume domain name
	return &strEntry{e: s}
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

type strEntry struct {
	e string
}

func (w *strEntry) Contains(e Entry) bool {
	return w.Match(e.String())
}

func (w *strEntry) Match(s string) bool {
	return glob.Glob(w.e, s)
}

func (w *strEntry) String() string {
	return w.e
}

func (w *strEntry) Equal(other Entry) bool {
	return other.String() == w.String()
}
