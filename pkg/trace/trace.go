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

package trace

import (
	"os"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/vic/pkg/log"
)

var tracingEnabled = true

// Enable global tracing.
func EnableTracing() {
	tracingEnabled = true
}

// Disable global tracing.
func DisableTracing() {
	tracingEnabled = false
}

var Logger = &logrus.Logger{
	Out: os.Stderr,
	// We're using our own text formatter to skip the \n and \t escaping logrus
	// was doing on non TTY Out (we redirect to a file) descriptors.
	Formatter: log.NewTextFormatter(),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.InfoLevel,
}

// trace object used to grab run-time state
type Message struct {
	msg      string
	funcName string
	lineNo   int

	startTime time.Time
}

func (t *Message) delta() time.Duration {
	if t == nil {
		return 0
	}
	return time.Now().Sub(t.startTime)
}

// begin a trace from this stack frame less the skip.
func newTrace(msg string, skip int) *Message {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return nil
	}

	name := runtime.FuncForPC(pc).Name()

	return &Message{
		msg:       msg,
		funcName:  name,
		lineNo:    line,
		startTime: time.Now(),
	}
}

// Begin starts the trace.  Msg is the msg to log.
func Begin(msg string) *Message {
	if tracingEnabled && Logger.Level >= logrus.DebugLevel {
		if t := newTrace(msg, 2); t != nil {
			if msg == "" {
				Logger.Debugf("[BEGIN] [%s:%d]", t.funcName, t.lineNo)
			} else {
				Logger.Debugf("[BEGIN] [%s:%d] %s", t.funcName, t.lineNo, t.msg)
			}
			return t

		}
	}
	return nil
}

// End ends the trace.
func End(t *Message) {
	if t == nil {
		return
	}
	Logger.Debugf("[ END ] [%s:%d] [%s] %s", t.funcName, t.lineNo, t.delta(), t.msg)
}
