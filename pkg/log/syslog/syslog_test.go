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

// +build !windows,!nacl,!plan9

package syslog

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/RackSec/srslog"
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var cfg = &SyslogConfig{
	Network:   "tcp",
	RAddr:     "localhost:514",
	Tag:       "test",
	Priority:  srslog.LOG_INFO | srslog.LOG_DAEMON,
	Formatter: RFC3164,
}

func TestNewSyslogHook(t *testing.T) {
	// other errors should still result in a hook being created
	d := &mockDialer{}
	d.On("dial", cfg).Return(nil, assert.AnError)
	h, err := NewHook(cfg, d)
	assert.Nil(t, h)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
	d.AssertCalled(t, "dial", cfg)
	d.AssertNumberOfCalls(t, "dial", 1)

	d = &mockDialer{}
	w := &MockWriter{}
	d.On("dial", cfg).Return(w, nil)
	h, err = NewHook(cfg, d)
	assert.NotNil(t, h)
	assert.NoError(t, err)
	assert.Equal(t, w, h.writer)
	d.AssertCalled(t, "dial", cfg)
	d.AssertNumberOfCalls(t, "dial", 1)
}

func TestLevels(t *testing.T) {
	m := &MockWriter{}
	d := &mockDialer{}
	d.On("dial", mock.Anything).Return(m, nil)
	h, err := NewHook(cfg, d)

	assert.NotNil(t, h)
	assert.NoError(t, err)

	m.On("Crit", mock.Anything).Return(nil)
	m.On("Err", mock.Anything).Return(nil)
	m.On("Warning", mock.Anything).Return(nil)
	m.On("Debug", mock.Anything).Return(nil)
	m.On("Info", mock.Anything).Return(nil)

	var tests = []struct {
		entry *logrus.Entry
		f     string
	}{
		{
			entry: &logrus.Entry{Message: "panic", Level: logrus.PanicLevel},
			f:     "Crit",
		},
		{
			entry: &logrus.Entry{Message: "fatal", Level: logrus.FatalLevel},
			f:     "Crit",
		},
		{
			entry: &logrus.Entry{Message: "error", Level: logrus.ErrorLevel},
			f:     "Err",
		},
		{
			entry: &logrus.Entry{Message: "warn", Level: logrus.WarnLevel},
			f:     "Warning",
		},
		{
			entry: &logrus.Entry{Message: "info", Level: logrus.InfoLevel},
			f:     "Info",
		},
		{
			entry: &logrus.Entry{Message: "debug", Level: logrus.DebugLevel},
			f:     "Debug",
		},
	}

	calls := make(map[string]int)
	for _, te := range tests {
		calls[te.f] = 0
	}

	for _, te := range tests {
		assert.NoError(t, h.writeEntry(te.entry))
		calls[te.f]++
		m.AssertCalled(t, te.f, te.entry.Message)
		m.AssertNumberOfCalls(t, te.f, calls[te.f])
	}
}

func TestWriterWrapperReconnect(t *testing.T) {
	d := &mockDialer{}
	w := &writerWrapper{
		dialer: d,
		cfg:    *cfg,
	}

	d.On("dial", &w.cfg).Return(nil, assert.AnError)
	calls := []struct {
		fm string
		f  func(string) error
	}{
		{"Emerg", w.Emerg},
		{"Crit", w.Crit},
		{"Err", w.Err},
		{"Warning", w.Warning},
		{"Info", w.Info},
		{"Debug", w.Debug},
	}
	i := 0
	for _, c := range calls {
		err := c.f("test")
		assert.NoError(t, err)
		d.AssertCalled(t, "dial", &w.cfg)
		i++
		d.AssertNumberOfCalls(t, "dial", i)
	}

	mw := &MockWriter{}
	d = &mockDialer{}
	d.On("dial", &w.cfg).Return(mw, nil)
	w.dialer = d
	for _, c := range calls {
		mw.On(c.fm, "test").Return(nil)
		err := c.f("test")
		assert.NoError(t, err)
		d.AssertCalled(t, "dial", &w.cfg)
		d.AssertNumberOfCalls(t, "dial", 1)
	}

	// no reconnect in the case of Close()
	d = &mockDialer{}
	d.On("dial", &w.cfg).Return(mw, nil)
	w = &writerWrapper{
		dialer: d,
		cfg:    *cfg,
	}

	mw.On("Close").Return(nil)
	w.Close()
	d.AssertNotCalled(t, "dial", &w.cfg)

	w.writer = mw
	w.Close()
	mw.AssertCalled(t, "Close")
}

func TestWriterWrapperWrite(t *testing.T) {
	d := &mockDialer{}
	w := &writerWrapper{
		dialer: d,
		cfg:    *cfg,
	}

	mw := &MockWriter{}
	d.On("dial", &w.cfg).Return(mw, nil)

	var tests = []struct {
		b   []byte
		n   int
		err error
	}{
		{[]byte("foo"), len("foo"), nil},
		{[]byte("bar"), 0, assert.AnError},
	}

	for _, te := range tests {
		mw.On("Write", te.b).Return(te.n, te.err)
		n, err := w.Write(te.b)
		if te.err != nil {
			assert.IsType(t, te.err, err)
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, te.n, n)
		mw.AssertCalled(t, "Write", te.b)

		// dial should be called only once
		d.AssertNumberOfCalls(t, "dial", 1)
	}
}

func TestMaxLogBuffer(t *testing.T) {
	d := &mockDialer{}
	w := &MockWriter{}
	d.On("dial", cfg).Return(w, nil)
	h, err := NewHook(cfg, d)
	assert.NotNil(t, h)
	assert.NoError(t, err)
	assert.Equal(t, h.writer, w)

	for i := 0; i < maxLogBuffer+1; i++ {
		h.Fire(&logrus.Entry{
			Message: fmt.Sprintf("%d", i),
			Level:   logrus.DebugLevel,
		})
	}

	w.On("Debug", mock.Anything).Return(nil)

	go h.Run()

	<-h.running

	select {
	case <-time.After(5 * time.Second):
		for i := 0; i < maxLogBuffer; i++ {
			w.AssertCalled(t, "Debug", fmt.Sprintf("%d", i))
		}
		w.AssertNumberOfCalls(t, "Debug", maxLogBuffer)
		w.AssertNotCalled(t, "Debug", fmt.Sprintf("%d", maxLogBuffer))
	}
}

func TestDefaultDialer(t *testing.T) {
	var tests = []error{
		&net.AddrError{},
		&net.ParseError{},
	}

	for _, te := range tests {
		md := &mockDialer{}
		dd := &defaultDialer{
			d: md,
		}

		md.On("dial", cfg).Return(nil, te)
		w, err := dd.dial(cfg)
		assert.Nil(t, w)
		assert.Error(t, err)
		assert.IsType(t, te, err)
		md.AssertCalled(t, "dial", cfg)
		md.AssertNumberOfCalls(t, "dial", 1)
	}

	tests = []error{
		assert.AnError,
	}

	for _, te := range tests {
		md := &mockDialer{}
		dd := &defaultDialer{
			d: md,
		}

		md.On("dial", cfg).Return(nil, te)
		w, err := dd.dial(cfg)
		assert.NotNil(t, w)
		assert.NoError(t, err)
		md.AssertCalled(t, "dial", cfg)
		md.AssertNumberOfCalls(t, "dial", 1)
	}
}

func TestConnect(t *testing.T) {
	// attempt a connection to a server that
	// does not exist
	h, err := NewHook(&SyslogConfig{
		Network:   "tcp",
		RAddr:     "foo:514",
		Tag:       "test",
		Formatter: RFC3164,
		Priority:  srslog.LOG_INFO,
	}, nil)

	assert.NoError(t, err)
	assert.NotNil(t, h)

	go h.Run()

	<-h.running

	h.Fire(&logrus.Entry{
		Message: "foo",
		Level:   logrus.InfoLevel,
	})

	<-time.After(5 * time.Second)
}

func TestMakeTag(t *testing.T) {
	p := filepath.Base(os.Args[0])
	if len(p) > maxTagLen {
		p = p[:maxTagLen]
	}

	var tests = []struct {
		prefix, proc string
		out          string
	}{
		{
			prefix: "",
			proc:   "",
			out:    p,
		},
		{
			prefix: "",
			proc:   "foo",
			out:    "foo",
		},
		{
			prefix: "foo",
			proc:   "",
			out:    "foo" + ":" + p,
		},
		{
			prefix: "foo",
			proc:   "bar",
			out:    "foo:bar",
		},
		{
			prefix: "",
			proc:   strings.Repeat("a", maxTagLen),
			out:    strings.Repeat("a", maxTagLen),
		},
		{
			prefix: "",
			proc:   strings.Repeat("a", maxTagLen) + "c",
			out:    strings.Repeat("a", maxTagLen),
		},
		{
			prefix: "pre",
			proc:   strings.Repeat("a", maxTagLen-2) + "c",
			out:    strings.Repeat("a", maxTagLen-2) + "c",
		},
		{
			prefix: strings.Repeat("a", maxTagLen-1) + "c",
			proc:   "bar",
			out:    strings.Repeat("a", maxTagLen-len(":bar")) + "-bar",
		},
		{
			prefix: "bar",
			proc:   strings.Repeat("a", maxTagLen) + "c",
			out:    (strings.Repeat("a", maxTagLen) + "c")[:maxTagLen],
		},
	}

	for _, te := range tests {
		out := MakeTag(te.prefix, te.proc)
		assert.Equal(t, te.out, out)
		assert.Condition(t, func() bool { return len(out) <= maxTagLen })
	}
}
