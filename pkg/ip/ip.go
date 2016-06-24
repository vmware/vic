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

package ip

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

type Range struct {
	FirstIP net.IP `vic:"0.1" scope:"read-only" key:"first"`
	LastIP  net.IP `vic:"0.1" scope:"read-only" key:"last"`
}

func NewRange(first, last net.IP) *Range {
	return &Range{FirstIP: first, LastIP: last}
}

func (i *Range) Overlaps(other Range) bool {
	if (bytes.Compare(i.FirstIP, other.FirstIP) <= 0 && bytes.Compare(other.FirstIP, i.LastIP) <= 0) ||
		(bytes.Compare(i.FirstIP, other.LastIP) <= 0 && bytes.Compare(other.FirstIP, i.LastIP) <= 0) {
		return true
	}

	return false
}

func (i *Range) String() string {
	return fmt.Sprintf("%s-%s", i.FirstIP, i.LastIP)
}

func (i *Range) Equal(other *Range) bool {
	return i.FirstIP.Equal(other.FirstIP) && i.LastIP.Equal(other.LastIP)
}

func ParseRange(r string) *Range {
	var first, last net.IP
	// check if its a CIDR
	_, ipnet, _ := net.ParseCIDR(r)
	if ipnet != nil {
		first = ipnet.IP
		last := make(net.IP, len(first))
		for i, f := range first {
			last[i] = f | ^ipnet.Mask[i]
		}

		return &Range{first, last}
	}

	comps := strings.Split(r, "-")
	if len(comps) != 2 {
		return nil
	}

	first = net.ParseIP(comps[0])
	if first == nil {
		return nil
	}

	last = net.ParseIP(comps[1])
	if last == nil {
		var end int
		end, err := strconv.Atoi(comps[1])
		if err != nil || end <= int(first[15]) || end > math.MaxUint8 {
			return nil
		}

		last = net.IPv4(first[12], first[13], first[14], byte(end))
	}

	if bytes.Compare(first, last) > 0 {
		return nil
	}

	return &Range{first, last}
}

// MarshalText implements the encoding.TextMarshaler interface
func (i *Range) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UmarshalText implements the encoding.TextUnmarshaler interface
func (i *Range) UnmarshalText(text []byte) error {
	s := string(text)
	r := ParseRange(s)
	if r == nil {
		return fmt.Errorf("parse error: %s", s)
	}

	*i = *r
	return nil
}

func IsUnspecifiedIP(ip net.IP) bool {
	return len(ip) == 0 || ip.IsUnspecified()
}
