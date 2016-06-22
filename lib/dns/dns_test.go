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

package dns

import (
	"testing"

	log "github.com/Sirupsen/logrus"

	mdns "github.com/miekg/dns"
)

var (
	types = []uint16{
		mdns.TypeA,
		mdns.TypeTXT,
		mdns.TypeAAAA,
	}

	names = []string{
		"facebook.com.",
		"google.com.",
	}
)

func TestForwarding(t *testing.T) {
	log.SetLevel(log.PanicLevel)

	options.IP = "127.0.0.1"
	options.Port = 5354

	server := NewServer(options)
	if server != nil {
		server.Start()
	}

	c := new(mdns.Client)

	size := 1024
	for i := 0; i < size; i++ {
		for _, Δ := range types {
			for _, Θ := range names {
				m := new(mdns.Msg)

				m.SetQuestion(Θ, Δ)

				r, _, err := c.Exchange(m, server.Addr())
				if err != nil || len(r.Answer) == 0 {
					t.Fatalf("Exchange failed: %s", err)
				}
			}
		}
	}

	n := len(types) * len(names)
	if server.cache.Hits() != uint64(n*size-n) && server.cache.Misses() != uint64(size) {
		t.Fatalf("Cache hits %d misses %d", server.cache.Hits(), server.cache.Misses())
	}

	server.Stop()
	server.Wait()
}
