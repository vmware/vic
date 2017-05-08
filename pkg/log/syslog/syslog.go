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
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/RackSec/srslog"
	"github.com/Sirupsen/logrus"
)

type SyslogConfig struct {
	Network   string
	RAddr     string
	Tag       string
	Priority  srslog.Priority
	Formatter Formatter
}

type dialer interface {
	dial(*SyslogConfig) (Writer, error)
}

type Hook struct {
	entries chan *logrus.Entry
	writer  Writer
	cfg     SyslogConfig

	// for testing
	running chan struct{}
}

type Writer interface {
	io.WriteCloser

	Emerg(string) error
	Crit(string) error
	Err(string) error
	Warning(string) error
	Info(string) error
	Debug(string) error
}

const maxLogBuffer = 100

type Formatter int

const (
	RFC3164 Formatter = iota
)

func NewHook(cfg *SyslogConfig, d dialer) (*Hook, error) {
	hook := &Hook{
		entries: make(chan *logrus.Entry, maxLogBuffer),
		running: make(chan struct{}),
	}

	if d == nil {
		d = &defaultDialer{
			d: &syslogDialer{},
		}
	}

	var err error
	hook.writer, err = d.dial(cfg)
	if err != nil {
		return nil, err
	}

	return hook, nil
}

type defaultDialer struct {
	// for testing purposes
	d dialer
}

func (d *defaultDialer) dial(cfg *SyslogConfig) (Writer, error) {
	w, err := d.d.dial(cfg)
	if err != nil {
		switch err.(type) {
		case *net.AddrError, *net.ParseError:
			return nil, err
		}
	}

	return &writerWrapper{
		writer: w,
		dialer: d.d,
		cfg:    *cfg,
	}, nil
}

type syslogDialer struct{}

func (d *syslogDialer) dial(cfg *SyslogConfig) (Writer, error) {
	w, err := srslog.Dial(cfg.Network, cfg.RAddr, cfg.Priority, cfg.Tag)
	if w != nil {
		switch cfg.Formatter {
		case RFC3164:
			w.SetFormatter(srslog.RFC3164Formatter)
		}

		return w, nil
	}

	return nil, err
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	select {
	case hook.entries <- entry:
	default:
		// drop log entry
	}

	return nil
}

func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *Hook) Run() {
	close(hook.running)
	for entry := range hook.entries {
		hook.writeEntry(entry)
	}

	logrus.Warnf("exited syslog loop")
}

func (hook *Hook) writeEntry(entry *logrus.Entry) error {
	// just use the message since the timestamp
	// is added by the syslog package
	line := entry.Message

	switch entry.Level {
	case logrus.PanicLevel:
		return hook.writer.Crit(line)
	case logrus.FatalLevel:
		return hook.writer.Crit(line)
	case logrus.ErrorLevel:
		return hook.writer.Err(line)
	case logrus.WarnLevel:
		return hook.writer.Warning(line)
	case logrus.InfoLevel:
		return hook.writer.Info(line)
	case logrus.DebugLevel:
		return hook.writer.Debug(line)
	}

	return nil
}

type writerWrapper struct {
	writer Writer
	cfg    SyslogConfig
	dialer dialer
}

func (w *writerWrapper) withRetry(f func()) {
	if w.writer == nil {
		w.writer, _ = w.dialer.dial(&w.cfg)
	}

	if w.writer != nil {
		f()
	}
}

func (w *writerWrapper) Write(b []byte) (n int, err error) {
	w.withRetry(func() {
		n, err = w.writer.Write(b)
	})

	return
}

func (w *writerWrapper) Emerg(m string) (err error) {
	w.withRetry(func() {
		err = w.writer.Emerg(m)
	})

	return
}

func (w *writerWrapper) Crit(m string) (err error) {
	w.withRetry(func() {
		err = w.writer.Crit(m)
	})

	return
}

func (w *writerWrapper) Err(m string) (err error) {
	w.withRetry(func() {
		err = w.writer.Err(m)
	})

	return
}

func (w *writerWrapper) Warning(m string) (err error) {
	w.withRetry(func() {
		err = w.writer.Warning(m)
	})

	return
}

func (w *writerWrapper) Info(m string) (err error) {
	w.withRetry(func() {
		err = w.writer.Info(m)
	})

	return
}

func (w *writerWrapper) Debug(m string) (err error) {
	w.withRetry(func() {
		err = w.writer.Debug(m)
	})

	return
}

func (w *writerWrapper) Close() error {
	if w.writer == nil {
		return nil
	}

	return w.writer.Close()
}

const maxTagLen = 32

// MakeTag makes an RFC 3164 compliant tag (32 characters or less)
// using the provided proc. If proc is empty, the name of the current
// executable is used, trauncated to maxTagLen characters if
// necessary.
func MakeTag(proc string) string {
	if len(proc) == 0 {
		proc = filepath.Base(os.Args[0])
	}
	proc = strings.TrimSpace(proc)
	if len(proc) >= maxTagLen {
		return proc[:maxTagLen]
	}

	return proc
}
