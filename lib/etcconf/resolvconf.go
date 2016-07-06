// Copyright 2016 VMware, Inc. All Rights Reserved.
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
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

const ResolvConfPath = "/etc/resolv.conf"

type ResolvConf interface {
	Conf

	AddNameservers(...net.IP)
	RemoveNameservers(...net.IP)
	Nameservers() []net.IP
}

type resolvConf struct {
	sync.Mutex

	EntryConsumer

	dirty       bool
	path        string
	nameservers []net.IP
}

type resolvConfWalker struct {
	nameservers []net.IP
	i           int
}

func (w *resolvConfWalker) HasNext() bool {
	return w.i < len(w.nameservers)
}

func (w *resolvConfWalker) Next() string {
	s := fmt.Sprintf("nameserver %s", w.nameservers[w.i].String())
	w.i++
	return s
}

func NewResolvConf(path string) ResolvConf {
	if path == "" {
		path = ResolvConfPath
	}

	return &resolvConf{path: path}
}

func (r *resolvConf) ConsumeEntry(t string) error {
	r.Lock()
	defer r.Unlock()

	fs := strings.Fields(t)
	if len(fs) != 2 {
		log.Warnf("skipping invalid line \"%s\"", t)
		return nil
	}

	if fs[0] == "nameserver" {
		ip := net.ParseIP(fs[1])
		if ip == nil {
			log.Warnf("skipping invalid line \"%s\": invalid ip address", t)
			return nil
		}

		r.addNameservers(ip)
	}

	return nil
}

func (r *resolvConf) Load() error {
	r.Lock()
	defer r.Unlock()

	rc := &resolvConf{}
	if err := load(r.path, rc); err != nil {
		return err
	}

	r.nameservers = rc.nameservers
	return nil
}

func (r *resolvConf) Save() error {
	r.Lock()
	defer r.Unlock()

	log.Debugf("%+v", r)
	if !r.dirty {
		return nil
	}

	walker := &resolvConfWalker{nameservers: r.nameservers}
	log.Debugf("%+v", walker)
	if err := save(r.path, walker); err != nil {
		return err
	}

	r.dirty = false
	return nil
}

func (r *resolvConf) AddNameservers(nss ...net.IP) {
	r.Lock()
	defer r.Unlock()

	r.addNameservers(nss...)
}

func (r *resolvConf) addNameservers(nss ...net.IP) {
	for _, n := range nss {
		if n == nil {
			continue
		}

		found := false
		for _, rn := range r.nameservers {
			if rn.Equal(n) {
				found = true
				break
			}
		}

		if !found {
			r.nameservers = append(r.nameservers, n)
			r.dirty = true
		}
	}
}

func (r *resolvConf) RemoveNameservers(nss ...net.IP) {
	r.Lock()
	defer r.Unlock()

	for _, n := range nss {
		if n == nil {
			continue
		}

		for i, rn := range r.nameservers {
			if n.Equal(rn) {
				r.nameservers = append(r.nameservers[:i], r.nameservers[i+1:]...)
				r.dirty = true
				break
			}
		}
	}
}

func (r *resolvConf) Nameservers() []net.IP {
	r.Lock()
	defer r.Unlock()

	return r.nameservers
}
