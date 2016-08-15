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
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
)

var Logger = log.New()

type tr struct {
	msg       string
	frameName string
	startTime time.Time
}

func newTrace(msg string) *tr {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	return &tr{
		msg:       msg,
		frameName: name,
		startTime: time.Now(),
	}
}

// Begin starts the trace.  Msg is the msg to log.
func Begin(msg string) *tr {
	t := newTrace(msg)

	if msg == "" {
		Logger.Debugf("[BEGIN] [%s]", t.frameName)
	} else {
		Logger.Debugf("[BEGIN] [%s] %s", t.frameName, t.msg)
	}

	return t
}

// End ends the trace.
func End(t *tr) {
	endTime := time.Now()
	Logger.Debugf("[ END ] [%s] [%s] %s", t.frameName, endTime.Sub(t.startTime), t.msg)
}
