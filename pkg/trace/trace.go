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

func Begin(msg string) (string, string, time.Time) {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	if msg == "" {
		log.Debugf("[BEGIN] [%s]", name)
	} else {
		log.Debugf("[BEGIN] [%s] %s", name, msg)
	}
	return msg, name, time.Now()
}

func End(msg string, name string, startTime time.Time) {
	endTime := time.Now()
	log.Debugf("[ END ] [%s] [%s] %s", name, endTime.Sub(startTime), msg)
}
