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

// IP address managment

package network

import (
	"fmt"
	"net"
)

// A tuple signifying an IP address range
type addressRange struct {
	firstIP net.IP
	lastIP  net.IP
}

// An AddressSpace is a collection of
// IP address ranges
type AddressSpace struct {
	parent          *AddressSpace
	availableRanges []addressRange
}

// compareIPv4 compares two IPv4 addresses.
// Returns -1 if ip1 < ip2, 0 if they are equal,
// and 1 if ip1 > ip2
func compareIP4(ip1 net.IP, ip2 net.IP) int {
	ip1 = ip1.To16()
	ip2 = ip2.To16()
	for i := 0; i < len(ip1); i++ {
		if ip1[i] < ip2[i] {
			return -1
		}

		if ip1[i] > ip2[i] {
			return 1
		}
	}

	return 0
}

func incrementIP4(ip net.IP) net.IP {
	if !isIP4(ip) {
		return nil
	}

	newIP := copyIP(ip)
	s := 0
	if len(ip) == net.IPv6len {
		s = 12
	}
	for i := len(newIP) - 1; i >= s; i-- {
		newIP[i]++
		if newIP[i] > 0 {
			break
		}
	}

	return newIP
}

func decrementIP4(ip net.IP) net.IP {
	if !isIP4(ip) {
		return nil
	}

	newIP := copyIP(ip)
	s := 0
	if len(ip) == net.IPv6len {
		s = 12
	}
	for i := len(newIP) - 1; i >= s; i-- {
		newIP[i]--
		if newIP[i] != 0xff {
			break
		}
	}

	return newIP
}

func copyIP(ip net.IP) net.IP {
	newIP := make([]byte, len(ip))
	copy(newIP, ip)
	return newIP
}

func isIP4(ip net.IP) bool {
	if ip = ip.To4(); ip != nil {
		return true
	}

	return false
}

func lowestIP4(ipRange net.IPNet) net.IP {
	return ipRange.IP.Mask(ipRange.Mask).To16()
}

func highestIP4(ipRange net.IPNet) net.IP {
	if !isIP4(ipRange.IP) {
		return nil
	}

	newIP := net.IPv4(0, 0, 0, 0)
	ipRange.IP = ipRange.IP.To4()
	for i := 0; i < len(ipRange.Mask); i++ {
		newIP[i+12] = ipRange.IP[i] | ^ipRange.Mask[i]
	}

	return newIP
}

// Creates a new AddressSpace from a network specification.
func NewAddressSpaceFromNetwork(ipRange net.IPNet) *AddressSpace {
	return &AddressSpace{
		availableRanges: []addressRange{{firstIP: lowestIP4(ipRange),
			lastIP: highestIP4(ipRange)}}}
}

// Creates a new AddressSpace from a range of IP addresses.
func NewAddressSpaceFromRange(firstIP net.IP, lastIP net.IP) *AddressSpace {
	return &AddressSpace{availableRanges: []addressRange{{firstIP: firstIP, lastIP: lastIP}}}
}

// Reserves a new sub address space within the given address space, given a
// bitmask specifying the "width" of the requested space.
func (s *AddressSpace) ReserveNextIP4Net(mask net.IPMask) (*AddressSpace, error) {
	for i, r := range s.availableRanges {
		network := r.firstIP.Mask(mask).To16()
		var firstIP net.IP
		if compareIP4(network, r.firstIP) >= 0 {
			// found the start of the range
			firstIP = network
		} else {
			// range does not start on the
			// mask boundary
			ones, _ := mask.Size()
			for i := len(network) - 1; i >= 12; i-- {
				partialByteIndex := ones/8 + 12
				if i == partialByteIndex {
					inc := (byte)(2 << (uint)(8-ones%8))
					if network[i]+inc > 0 {
						network[i] += inc
						firstIP = network
						break
					}
				} else if i < partialByteIndex {
					if network[i]+1 > 0 {
						network[i]++
						firstIP = network
						break
					}
				}
			}
		}

		if firstIP != nil {
			// check if the available range can accommodate the highest address
			// given the mask
			lastIP := highestIP4(net.IPNet{IP: firstIP, Mask: mask})
			if compareIP4(lastIP, r.lastIP) <= 0 {
				s.reserveSubRange(firstIP, lastIP, i)
				return NewAddressSpaceFromRange(firstIP, lastIP), nil
			}
		}
	}

	return nil, fmt.Errorf("could not find IP range for mask %s", mask)
}

func newAddressRange(firstIP net.IP, lastIP net.IP) *addressRange {
	return &addressRange{firstIP: firstIP.To16(), lastIP: lastIP.To16()}
}

func splitRange(parentRange addressRange, firstIP net.IP, lastIP net.IP) (before, reserved, after *addressRange) {
	if !firstIP.Equal(parentRange.firstIP) {
		before = newAddressRange(parentRange.firstIP, decrementIP4(firstIP))
	}
	if !lastIP.Equal(parentRange.lastIP) {
		after = newAddressRange(incrementIP4(lastIP), parentRange.lastIP)
	}
	reserved = newAddressRange(firstIP, lastIP)
	return
}

// Reserve a new sub address space given an IP and mask.
// Mask is required.
// If IP is nil or "0.0.0.0", same as calling ReserveNextIP4Net
// with the mask.
func (s *AddressSpace) ReserveIP4Net(ipNet net.IPNet) (*AddressSpace, error) {
	if ipNet.Mask == nil {
		return nil, fmt.Errorf("network mask not specified")
	}

	if ipNet.IP == nil || ipNet.IP.Equal(net.ParseIP("0.0.0.0")) {
		return s.ReserveNextIP4Net(ipNet.Mask)
	}

	return s.ReserveIP4Range(lowestIP4(ipNet), highestIP4(ipNet))
}

func (s *AddressSpace) reserveSubRange(firstIP net.IP, lastIP net.IP, index int) {
	before, _, after := splitRange(s.availableRanges[index], firstIP, lastIP)
	s.availableRanges = append(s.availableRanges[:index], s.availableRanges[index+1:]...)
	if before != nil {
		s.availableRanges = insertAddressRanges(s.availableRanges, index, *before)
		index++
	}
	if after != nil {
		s.availableRanges = insertAddressRanges(s.availableRanges, index, *after)
	}
}

// Reserve a sub address space given a first and last IP.
func (s *AddressSpace) ReserveIP4Range(firstIP net.IP, lastIP net.IP) (*AddressSpace, error) {
	for i, r := range s.availableRanges {
		if compareIP4(firstIP, r.firstIP) < 0 ||
			compareIP4(lastIP, r.lastIP) > 0 {
			continue
		}

		// found range
		s.reserveSubRange(firstIP, lastIP, i)
		return NewAddressSpaceFromRange(firstIP, lastIP), nil
	}

	return nil, fmt.Errorf("could not find IP range")
}

func insertAddressRanges(r []addressRange, index int, ranges ...addressRange) []addressRange {
	if index == len(r) {
		return append(r, ranges...)
	}

	for i := 0; i < len(ranges); i++ {
		r = append(r, addressRange{})
	}

	copy(r[index+len(ranges):], r[index:])
	for i := 0; i < len(ranges); i++ {
		r[index+i] = ranges[i]
	}

	return r
}

// Reserve the next available IPv4 address.
func (s *AddressSpace) ReserveNextIP4() (net.IP, error) {
	space, err := s.ReserveIP4Net(net.IPNet{Mask: net.CIDRMask(32, 32)})
	if err != nil {
		return nil, err
	}

	return space.availableRanges[0].firstIP, nil
}

// Reserve the given IPv4 address.
func (s *AddressSpace) ReserveIP4(ip net.IP) error {
	_, err := s.ReserveIP4Range(ip, ip)
	return err
}

// Release a sub address space into the parent address space.
// Sub address space has to have only a single available range.
func (s *AddressSpace) ReleaseIP4Range(space *AddressSpace) error {
	// nothing to release
	if space == nil || len(space.availableRanges) == 0 {
		return nil
	}

	// cannot release a range if it has more than one available sub range
	if len(space.availableRanges) > 1 {
		return fmt.Errorf("can not release an address space with more than one available range")
	}

	firstIP := space.availableRanges[0].firstIP
	lastIP := space.availableRanges[0].lastIP
	if compareIP4(firstIP, lastIP) > 0 {
		return fmt.Errorf("address space first ip %s is greater than last ip %s", firstIP, lastIP)
	}

	i := 0
	for ; i < len(s.availableRanges); i++ {
		if compareIP4(lastIP, s.availableRanges[i].firstIP) < 0 {
			if i == 0 {
				break
			}

			if i > 0 && compareIP4(firstIP, s.availableRanges[i-1].lastIP) > 0 {
				break
			}
		}
	}

	if i > 0 && i == len(s.availableRanges) {
		if compareIP4(firstIP, s.availableRanges[i-1].lastIP) <= 0 {
			return fmt.Errorf("Could not release IP range")
		}
	}

	s.availableRanges = insertAddressRanges(s.availableRanges, i, space.availableRanges...)
	return nil
}

// Release the given IPv4 address.
func (s *AddressSpace) ReleaseIP4(ip net.IP) error {
	return s.ReleaseIP4Range(NewAddressSpaceFromRange(ip, ip))
}

// "Defragments" the available IP address ranges.
func (s *AddressSpace) Defragment() error {
	for i := 1; i < len(s.availableRanges); {
		first := &s.availableRanges[i-1]
		second := &s.availableRanges[i]
		if incrementIP4(first.lastIP).Equal(second.firstIP) {
			first.lastIP = second.lastIP
			s.availableRanges = append(s.availableRanges[:i], s.availableRanges[i+1:]...)
		} else {
			i++
		}
	}

	return nil
}

// Compares two address spaces for equality.
func (s *AddressSpace) Equal(other *AddressSpace) bool {
	if len(s.availableRanges) != len(other.availableRanges) {
		return false
	}

	for i := 0; i < len(s.availableRanges); i++ {
		if compareIP4(s.availableRanges[i].firstIP, other.availableRanges[i].firstIP) != 0 ||
			compareIP4(s.availableRanges[i].lastIP, other.availableRanges[i].lastIP) != 0 {
			return false
		}
	}

	return true
}
