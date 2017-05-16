// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!nacl,!plan9

package syslog

import (
	"errors"
	"net"
)

type unixSyslogDialer struct{}

// unixSyslog opens a connection to the syslog daemon running on the
// local machine using a Unix domain socket.
func (u *unixSyslogDialer) dial() (conn, error) {
	logTypes := []string{"unixgram", "unix"}
	logPaths := []string{"/dev/log", "/var/run/syslog", "/var/run/log"}
	for _, network := range logTypes {
		for _, path := range logPaths {
			conn, err := net.Dial(network, path)
			if err != nil {
				continue
			} else {
				return &localConn{c: conn}, nil
			}
		}
	}
	return nil, errors.New("Unix syslog delivery error")
}

func newDialer(network, address string) dialer {
	if network == "" {
		return &unixSyslogDialer{}
	}

	return &defaultDialer{
		network: network,
		address: address,
	}
}
