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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWriterReconnect(t *testing.T) {
	dn := &mockNetDialer{}
	dn.On("dial").Return(nil, assert.AnError)
	w := newWriter(priority, tag, "", dn, nil)

	go w.run()
	<-w.running

	calls := []func(string) error{
		w.Emerg,
		w.Crit,
		w.Err,
		w.Warning,
		w.Info,
		w.Debug,
	}
	for _, f := range calls {
		err := f("test")
		assert.NoError(t, err)
	}

	assertExpectations(t, func() bool {
		return dn.AssertNumberOfCalls(t, "dial", 1+len(calls))
	})

	// no reconnect in the case of Close()
	dn = &mockNetDialer{}
	w.dialer = dn
	w.Close()
	dn.AssertNotCalled(t, "dial")
}

func assertExpectations(t *testing.T, f func() bool) {
	ticker := time.NewTicker(1 * time.Second)
	for i := 0; i < 5; i++ {
		select {
		case <-ticker.C:
			if f() {
				ticker.Stop()
				return
			}
		}
	}

	// fail after 5 iterations
	assert.FailNow(t, "timed out waiting for write")
}

func TestWriterWrite(t *testing.T) {
	msg := "foo"

	f := &mockFormatter{}
	f.On("Format", priority, mock.Anything, "host", tag, msg).Return("test")

	a := &MockAddr{}
	a.On("String").Return("host:123")

	c := &MockNetConn{}
	c.On("LocalAddr").Return(a)
	c.On("Write", []byte("test\n")).Return(len(msg), nil)

	dn := &mockNetDialer{}
	dn.On("dial").Return(c, nil)

	w := newWriter(priority, tag, "", dn, f)
	n, err := w.Write([]byte(msg))
	assert.NoError(t, err)
	assert.Equal(t, len(msg), n)

	go w.run()

	<-w.running

	assertExpectations(t, func() bool {
		return c.AssertExpectations(t) && dn.AssertNumberOfCalls(t, "dial", 1)
	})
}

func TestMaxLogBuffer(t *testing.T) {
	f := &mockFormatter{}

	dn := &mockNetDialer{}
	c := &MockNetConn{}
	a := &MockAddr{}
	a.On("String").Return("foo")
	c.On("LocalAddr").Return(a)

	dn.On("dial").Return(c, nil)
	w := newWriter(priority, tag, "", dn, f)

	for i := 0; i < maxLogBuffer+1; i++ {
		msg := fmt.Sprintf("%d", i)
		f.On("Format", priority, mock.Anything, "", tag, msg).Return(msg)
		c.On("Write", []byte(msg+"\n")).Return(len(msg), nil)
		w.Write([]byte(msg))
	}

	go w.run()

	<-w.running

	assertExpectations(t, func() bool {
		for i := 0; i < maxLogBuffer; i++ {
			if !f.AssertCalled(t, "Format", priority, mock.Anything, "", tag, fmt.Sprintf("%d", i)) ||
				!c.AssertCalled(t, "Write", []byte(fmt.Sprintf("%d\n", i))) {
				return false
			}
		}
		return f.AssertNumberOfCalls(t, "Format", maxLogBuffer) && f.AssertNotCalled(t, "Format", priority, mock.Anything, "", tag, fmt.Sprintf("%d", maxLogBuffer))
	})
}
