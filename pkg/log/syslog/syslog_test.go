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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	network  = "tcp"
	raddr    = "localhost:514"
	tag      = "test"
	priority = LOG_INFO | LOG_DAEMON
)

func TestNewSyslogHook(t *testing.T) {
	// other errors should still result in a hook being created
	d := &mockDialer{}
	d.On("dial").Return(nil, assert.AnError)
	h, err := NewHook(network, raddr, priority, tag)
	assert.Nil(t, h)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
	d.AssertCalled(t, "dial")
	d.AssertNumberOfCalls(t, "dial", 1)

	d = &mockDialer{}
	w := &MockWriter{}
	d.On("dial").Return(w, nil)
	h, err = NewHook(network, raddr, priority, tag)
	assert.NotNil(t, h)
	assert.NoError(t, err)
	assert.Equal(t, w, h.writer)
	d.AssertCalled(t, "dial")
	d.AssertNumberOfCalls(t, "dial", 1)
}

func TestLevels(t *testing.T) {
	m := &MockWriter{}
	d := &mockDialer{}
	d.On("dial").Return(m, nil)
	h, err := NewHook(network, raddr, priority, tag)

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

func TestWriterReconnect(t *testing.T) {
	d := &mockDialer{}
	w := &writer{
		dialer: d,
	}

	d.On("dial").Return(nil, assert.AnError)
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
		d.AssertCalled(t, "dial")
		i++
		d.AssertNumberOfCalls(t, "dial", i)
	}

	mw := &MockWriter{}
	d = &mockDialer{}
	d.On("dial").Return(mw, nil)
	w.dialer = d
	for _, c := range calls {
		mw.On(c.fm, "test").Return(nil)
		err := c.f("test")
		assert.NoError(t, err)
		d.AssertCalled(t, "dial")
		d.AssertNumberOfCalls(t, "dial", 1)
	}

	// no reconnect in the case of Close()
	d = &mockDialer{}
	d.On("dial").Return(mw, nil)
	mw.On("Close").Return(nil)
	mw.Close()
	d.AssertNotCalled(t, "dial")
}

func TestWriterWrite(t *testing.T) {
	d := &mockDialer{}
	w := &writer{
		dialer: d,
	}

	mw := &MockWriter{}
	d.On("dial").Return(mw, nil)

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
	d.On("dial").Return(w, nil)
	h, err := NewHook(network, raddr, priority, tag)
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

	select {
	case <-time.After(5 * time.Second):
		for i := 0; i < maxLogBuffer; i++ {
			w.AssertCalled(t, "Debug", fmt.Sprintf("%d", i))
		}
		w.AssertNumberOfCalls(t, "Debug", maxLogBuffer)
		w.AssertNotCalled(t, "Debug", fmt.Sprintf("%d", maxLogBuffer))
	}
}

func TestConnect(t *testing.T) {
	// attempt a connection to a server that
	// does not exist
	h, err := NewHook(
		"tcp",
		"foo:514",
		LOG_INFO,
		"test",
	)

	assert.NoError(t, err)
	assert.NotNil(t, h)

	h.Fire(&logrus.Entry{
		Message: "foo",
		Level:   logrus.InfoLevel,
	})

	<-time.After(5 * time.Second)
}

func TestMakeTag(t *testing.T) {
	p := filepath.Base(os.Args[0])
	var tests = []struct {
		prefix string
		proc   string
		out    string
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
			out:    "foo" + sep + p,
		},
		{
			prefix: "bar",
			proc:   "foo",
			out:    "bar" + sep + "foo",
		},
	}

	for _, te := range tests {
		out := MakeTag(te.prefix, te.proc)
		assert.Equal(t, te.out, out)
	}
}
