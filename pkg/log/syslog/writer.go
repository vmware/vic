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
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Writer interface {
	io.WriteCloser

	Emerg(string) error
	Crit(string) error
	Err(string) error
	Warning(string) error
	Info(string) error
	Debug(string) error

	WithTag(tag string) Writer
}

type writer struct {
	priority Priority
	tag      string
	hostname string

	msgs          chan *msg
	once          sync.Once
	done, running chan struct{}

	dialer    netDialer
	conn      net.Conn
	formatter formatter

	parent *writer
}

type msg struct {
	p   Priority
	tag string
	msg string
}

func newWriter(priority Priority, tag, hostname string, dialer netDialer, f formatter) *writer {
	return &writer{
		priority:  priority,
		tag:       tag,
		hostname:  hostname,
		dialer:    dialer,
		msgs:      make(chan *msg, maxLogBuffer),
		done:      make(chan struct{}),
		running:   make(chan struct{}),
		formatter: f,
	}
}

// connect makes a connection to the syslog server.
// It must be called with w.mu held.
func (w *writer) connect() (err error) {
	if w.conn != nil {
		// ignore err from close, it makes sense to continue anyway
		w.conn.Close()
		w.conn = nil
	}

	w.conn, err = w.dialer.dial()
	if err == nil {
		Logger.Info("successfully connected to syslog server")
		if w.hostname == "" {
			w.hostname, _, _ = net.SplitHostPort(w.conn.LocalAddr().String())
		}
	}
	return
}

// Write sends a log message to the syslog daemon.
func (w *writer) Write(b []byte) (int, error) {
	w.queueWrite(w.priority, w.tag, string(b))
	return len(b), nil
}

// Close closes a connection to the syslog daemon.
func (w *writer) Close() error {
	w.once.Do(func() {
		close(w.msgs)
		<-w.done
	})
	return nil
}

// Emerg logs a message with severity LOG_EMERG, ignoring the severity
// passed to New.
func (w *writer) Emerg(m string) error {
	return w.queueWrite(LOG_EMERG, w.tag, m)
}

// Alert logs a message with severity LOG_ALERT, ignoring the severity
// passed to New.
func (w *writer) Alert(m string) error {
	return w.queueWrite(LOG_ALERT, w.tag, m)
}

// Crit logs a message with severity LOG_CRIT, ignoring the severity
// passed to New.
func (w *writer) Crit(m string) error {
	return w.queueWrite(LOG_CRIT, w.tag, m)
}

// Err logs a message with severity LOG_ERR, ignoring the severity
// passed to New.
func (w *writer) Err(m string) error {
	return w.queueWrite(LOG_ERR, w.tag, m)
}

// Warning logs a message with severity LOG_WARNING, ignoring the
// severity passed to New.
func (w *writer) Warning(m string) error {
	return w.queueWrite(LOG_WARNING, w.tag, m)
}

// Notice logs a message with severity LOG_NOTICE, ignoring the
// severity passed to New.
func (w *writer) Notice(m string) error {
	return w.queueWrite(LOG_NOTICE, w.tag, m)
}

// Info logs a message with severity LOG_INFO, ignoring the severity
// passed to New.
func (w *writer) Info(m string) error {
	return w.queueWrite(LOG_INFO, w.tag, m)
}

// Debug logs a message with severity LOG_DEBUG, ignoring the severity
// passed to New.
func (w *writer) Debug(m string) error {
	return w.queueWrite(LOG_DEBUG, w.tag, m)
}

func (w *writer) queueWrite(p Priority, tag, s string) error {
	for w.parent != nil {
		w = w.parent
	}

	select {
	case w.msgs <- &msg{p: p, tag: tag, msg: s}:
	default:
		return errors.New("queue full or writer closed")
	}

	return nil
}

func (w *writer) writeAndRetry(p Priority, tag, s string) (int, error) {
	if len(s) == 0 {
		return 0, nil
	}

	pr := (w.priority & facilityMask) | (p & severityMask)

	if w.conn != nil {
		if n, err := w.write(pr, tag, s); err == nil {
			return n, err
		}
	}
	if err := w.connect(); err != nil {
		return 0, err
	}
	return w.write(pr, tag, s)
}

// write generates and writes a syslog formatted string. The
// format is as follows: <PRI>TIMESTAMP HOSTNAME TAG[PID]: MSG
func (w *writer) write(p Priority, tag, msg string) (int, error) {
	s := w.formatter.Format(p, time.Now(), w.hostname, tag, msg)
	// ensure it ends in a \n
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}
	_, err := w.conn.Write([]byte(s))
	if err != nil {
		return 0, err
	}

	// return len(msg), since we want to behave as an io.Writer
	return len(msg), nil
}

func (w *writer) WithTag(tag string) Writer {
	return &writer{
		hostname: w.hostname,
		tag:      tag,
		priority: w.priority,
		parent:   w,
	}
}

func (w *writer) run() {
	defer func() {
		Logger.Infof("exiting syslog writer loop")
		if w.conn != nil {
			w.conn.Close()
		}
		close(w.done)
	}()

	if err := w.connect(); err != nil {
		switch err.(type) {
		case *net.ParseError, *net.AddrError:
			Logger.Errorf("could not connec to syslog server (will not try again): %s", err)
			return
		}
	}

	close(w.running)

	for m := range w.msgs {
		for _, s := range strings.SplitAfter(m.msg, "\n") {
			if _, err := w.writeAndRetry(m.p, m.tag, s); err != nil {
				Logger.Debugf("could not write syslog message: %s", err)
			}
		}
	}

}
