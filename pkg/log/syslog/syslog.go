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
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
)

// The Priority is a combination of the syslog facility and
// severity. For example, LOG_ALERT | LOG_FTP sends an alert severity
// message from the FTP facility. The default severity is LOG_EMERG;
// the default facility is LOG_KERN.
type Priority int

const severityMask = 0x07
const facilityMask = 0xf8

const (
	// Severity.

	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

const (
	// Facility.

	// From /usr/include/sys/syslog.h.
	// These are the same up to LOG_FTP on Linux, BSD, and OS X.
	LOG_KERN Priority = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
	LOG_LOCAL2
	LOG_LOCAL3
	LOG_LOCAL4
	LOG_LOCAL5
	LOG_LOCAL6
	LOG_LOCAL7
)

// New establishes a new connection to the system log daemon. Each
// write to the returned writer sends a log message with the given
// priority and prefix.
func New(priority Priority, tag string) (Writer, error) {
	return Dial("", "", priority, tag)
}

// Dial establishes a connection to a log daemon by connecting to
// address raddr on the specified network. Each write to the returned
// writer sends a log message with the given facility, severity and
// tag.
// If network is empty, Dial will connect to the local syslog server.
func Dial(network, raddr string, priority Priority, tag string) (Writer, error) {
	if priority < 0 || priority > LOG_LOCAL7|LOG_DEBUG {
		return nil, errors.New("log/syslog: invalid priority")
	}

	return dial(priority, tag, newDialer(network, raddr))
}

const maxLogBuffer = 100

func dial(priority Priority, tag string, d dialer) (Writer, error) {
	tag = MakeTag("", tag)
	hostname, _ := os.Hostname()

	w := &writer{
		priority: priority,
		tag:      tag,
		hostname: hostname,
		dialer:   d,
		msgs:     make(chan *msg, maxLogBuffer),
		done:     make(chan struct{}),
	}

	go func() {
		defer func() {
			logger.Infof("exiting syslog writer loop")
			if w.conn != nil {
				w.conn.Close()
			}
			close(w.done)
		}()

		if err := w.connect(); err != nil {
			switch err.(type) {
			case *net.ParseError, *net.AddrError:
				logger.Errorf("could not connec to syslog server (will not try again): %s", err)
				return
			}
		}

		for m := range w.msgs {
			if m == nil {
				// writer closed
				return
			}

			for _, s := range strings.SplitAfter(m.msg, "\n") {
				if _, err := w.writeAndRetry(m.p, m.tag, s); err != nil {
					logger.Debugf("could not write syslog message: %s", err)
				}
			}
		}
	}()

	return w, nil
}

const sep = "/"

// MakeTag returns prfeix + sep + proc if prefix is not empty.
// If proc is empty, proc is set to filepath.Base(os.Args[0]).
// If prefix is empty, MakeTag returns proc.
func MakeTag(prefix, proc string) string {
	if len(proc) == 0 {
		proc = filepath.Base(os.Args[0])
	}

	if len(prefix) > 0 {
		return prefix + sep + proc
	}

	return proc
}

var logger = logrus.New()

func init() {
	logger.Level = logrus.DebugLevel
}
