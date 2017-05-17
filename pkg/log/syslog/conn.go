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

package syslog

import (
	"fmt"
	"net"
	"os"
	"time"
)

type conn interface {
	Local() bool
	Write(p Priority, ts time.Time, hostname, tag, msg string) (int, error)
	Close() error
}

type localConn struct {
	c net.Conn
}

func (l *localConn) Write(p Priority, ts time.Time, _, tag, msg string) (int, error) {
	return fmt.Fprintf(l.c, "<%d>%s %s[%d]: %s", p, ts.Format(time.Stamp), tag, os.Getpid(), msg)
}

func (l *localConn) Local() bool {
	return true
}

func (l *localConn) Close() error {
	return l.c.Close()
}

type rfc3164Conn struct {
	c net.Conn
}

func (c *rfc3164Conn) Write(p Priority, ts time.Time, hostname, tag, msg string) (int, error) {
	if hostname == "" {
		hostname, _, _ = net.SplitHostPort(c.c.LocalAddr().String())
	}
	return fmt.Fprintf(c.c, "<%d>%s %s %s[%d]: %s", p, ts.Format(time.RFC3339), hostname, tag, os.Getpid(), msg)
}

func (c *rfc3164Conn) Local() bool {
	return false
}

func (c *rfc3164Conn) Close() error {
	return c.c.Close()
}

type netDialer interface {
	dial() (conn, error)
}

type defaultNetDialer struct {
	network, address string
}

func (d *defaultNetDialer) dial() (conn, error) {
	c, err := net.DialTimeout(d.network, d.address, defaultDialTimeout)
	if err != nil {
		return nil, err
	}

	return &rfc3164Conn{c: c}, nil
}
