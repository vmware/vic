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
	"strings"

	glob "github.com/ryanuber/go-glob"
)

type Entry interface {
	Contains(e Entry) bool
	Match(e string) bool
	Equal(other Entry) bool
	String() string
}

func ParseEntry(s string) Entry {
	_, ipnet, err := net.ParseCIDR(s)
	if err == nil {
		return &cidrEntry{ipnet: ipnet}
	}

	u, err := url.Parse(s)
	if err == nil && len(u.Host) > 0 {
		return &urlEntry{u: u}
	}

	return &strEntry{e: s}
}

type cidrEntry struct {
	ipnet *net.IPNet
}

func (c *cidrEntry) Contains(e Entry) bool {
	if ip := net.ParseIP(e.String()); ip != nil {
		return c.ipnet.Contains(ip)
	}

	if e, ok := e.(*cidrEntry); ok {
		return c.ipnet.Contains(e.ipnet.IP.Mask(e.ipnet.Mask))
	}

	return false
}

func (c *cidrEntry) Match(s string) bool {
	h := getHost(s)
	ip := net.ParseIP(h)
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

type urlEntry struct {
	u *url.URL
}

func (u *urlEntry) Contains(e Entry) bool {
	return u.Match(e.String())
}

func (u *urlEntry) Match(s string) bool {
	return strings.HasPrefix(s, u.u.String())
}

func (u *urlEntry) String() string {
	return u.u.String()
}

func (u *urlEntry) Equal(other Entry) bool {
	return other.String() == u.u.String()
}

type strEntry struct {
	e string
}

func (w *strEntry) Contains(e Entry) bool {
	return w.Match(e.String())
}

func (w *strEntry) Match(s string) bool {
	// url?
	if u, err := url.Parse(s); err == nil && len(u.Host) > 0 {
		return glob.Glob(w.e, u.Host) ||
			glob.Glob(w.e, u.Hostname())
	}

	// host:port ?
	h, _, err := net.SplitHostPort(s)
	if err == nil && glob.Glob(w.e, h) {
		return true
	}

	return glob.Glob(w.e, s)
}

func (w *strEntry) String() string {
	return w.e
}

func (w *strEntry) Equal(other Entry) bool {
	return other.String() == w.String()
}

func getHost(s string) string {
	// url?
	if u, err := url.Parse(s); err == nil && len(u.Host) > 0 {
		return u.Hostname()
	}

	// host:port?
	if h, _, err := net.SplitHostPort(s); err == nil {
		return h
	}

	return s
}
