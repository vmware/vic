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

package trace

import (
	"fmt"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
)

var Logger = log.New()

// trace object used to grab run-time state
type Message struct {
	msg      string
	funcName string
	lineNo   int

	startTime time.Time
}

func (t *Message) delta() time.Duration {
	return time.Now().Sub(t.startTime)
}

func (t *Message) beginHdr() string {
	return fmt.Sprintf("[BEGIN] [%s:%d]", t.funcName, t.lineNo)
}

func (t *Message) endHdr() string {
	return fmt.Sprintf("[ END ] [%s:%d]", t.funcName, t.lineNo)
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
	t := newTrace(msg, 2)

	if msg == "" {
		Logger.Debugf(t.beginHdr())
	} else {
		Logger.Debugf("%s %s", t.beginHdr(), t.msg)
	}

	return t
}

// End ends the trace.
func End(t *Message) {
	Logger.Debugf("%s [%s] %s", t.endHdr(), t.delta(), t.msg)
}
