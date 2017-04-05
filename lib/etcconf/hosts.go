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

package etcconf

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type hostEntry struct {
	IP        net.IP
	Hostnames []string
}

func (e *hostEntry) String() string {
	return fmt.Sprintf("%s %s", e.IP, strings.Join(e.Hostnames, " "))
}

type Hosts interface {
	Conf

	SetHost(hostname string, ip net.IP)
	RemoveHost(hostname string)
	RemoveAll()

	HostIP(hostname string) net.IP
}

type hosts struct {
	sync.Mutex

	EntryConsumer

	hosts map[string]net.IP
	dirty bool
	path  string
}

type hostsWalker struct {
	entries []*hostEntry
	i       int
}

func (w *hostsWalker) HasNext() bool {
	return w.i < len(w.entries)
}

func (w *hostsWalker) Next() string {
	s := w.entries[w.i].String()
	w.i++
	return s
}

func NewHosts(path string) Hosts {
	if path == "" {
		path = hostsPath
	}

	return &hosts{
		path:  path,
		hosts: make(map[string]net.IP),
	}
}

func (h *hosts) ConsumeEntry(t string) error {
	h.Lock()
	defer h.Unlock()

	fs := strings.Fields(t)
	if len(fs) < 2 {
		log.Warnf("ignoring incomplete line %q", t)
		return nil
	}

	ip := net.ParseIP(fs[0])
	if ip == nil {
		log.Warnf("ignoring line %q due to invalid ip address", t)
		return nil
	}

	for _, hs := range fs[1:] {
		h.setHost(hs, ip)
	}

	return nil
}

func (h *hosts) Load() error {
	h.Lock()
	defer h.Unlock()

	newHosts := &hosts{hosts: make(map[string]net.IP)}
	if err := load(h.path, newHosts); err != nil {
		return err
	}

	h.hosts = newHosts.hosts
	h.dirty = false
	return nil
}

func (h *hosts) Save() error {
	h.Lock()
	defer h.Unlock()

	if !h.dirty {
		log.Debugf("skipping writing file since there are no new entries")
		return nil
	}

	var entries []*hostEntry
	for host, ip := range h.hosts {
		found := false
		for _, e := range entries {
			if e.IP.Equal(ip) {
				e.Hostnames = append(e.Hostnames, host)
				found = true
				break
			}
		}

		if !found {
			entries = append(entries, &hostEntry{IP: ip, Hostnames: []string{host}})
		}
	}

	if err := save(h.path, &hostsWalker{entries: entries}); err != nil {
		return err
	}

	// make sure the file is readable #nosec
	if err := os.Chmod(h.path, 0644); err != nil {
		return err
	}

	h.dirty = false
	return nil
}

func (h *hosts) SetHost(hostname string, ip net.IP) {
	h.Lock()
	defer h.Unlock()

	h.setHost(hostname, ip)
}

func (h *hosts) setHost(hostname string, ip net.IP) {
	h.hosts[hostname] = ip
	h.dirty = true
}

func (h *hosts) RemoveHost(hostname string) {
	h.Lock()
	defer h.Unlock()

	delete(h.hosts, hostname)
	h.dirty = true
}

func (h *hosts) RemoveAll() {
	h.Lock()
	defer h.Unlock()

	h.hosts = make(map[string]net.IP)
	h.dirty = true
}

func (h *hosts) HostIP(hostname string) net.IP {
	h.Lock()
	defer h.Unlock()

	return h.hosts[hostname]
}
